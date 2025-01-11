const createPayment = async (req, res) => {
    // const { amount } = req.body;
    const { buyer,buys,address,phone,
            totalPrice,paymentDate,status,
            fname,lname} =req.body;
    try {
        // if(!address || !phone) return  res.json({error:"address and phone number are required"})
  const payment = new Payment({ amount:totalPrice, reference: generateReference() });
    await payment.save();

      const buyerBuysProduct =await BuyerBuysProduct.create({buyer,buys,address,phone,totalPrice,status})
      if(buys){
        buys.map(async(Element)=>{
        const product=await Product.findOne({_id:Element.product})
            await  Product.findByIdAndUpdate({_id:Element.product},{quantity:Element.quantity<=product.quantity?product.quantity-Element.quantity:product.quantity}) 
        })
      }
     const reference=generateReference()
      const response = await axios.post('https://api.chapa.co/v1/transaction/initialize', {
        amount: totalPrice,
        currency: 'ETB',
        email: buyer,
        first_name: fname,
        last_name: lname,
        tx_ref: payment.reference,
        callback_url: 'http://localhost:3000/payment-result' // Chapa's callback URLhttp://localhost:3000/payment-result
      }, {
        headers: {
          Authorization:`Bearer ${process.env.CHAPA_SECRET_KEY}`
        }
      });
  
      res.status(200).json({
        chapaUrl: response.data.data.checkout_url,
        reference: payment.reference
      });
    } catch (error) {
      console.error('Payment initialization error:', error);
      res.status(500).json({ error: 'Payment initialization failed' });
    }
  };

  const confirmPayment = async (req, res) => {
    const { reference } = req.params;
  
    try {
      const response = await axios.get(`https://api.chapa.co/v1/transaction/verify/${reference}`, {
        headers: {
          Authorization: `Bearer ${process.env.CHAPA_SECRET_KEY}`
        }
      });
  
      const payment = await Payment.findOne({ reference });
      if (!payment) return res.status(404).json({ error: 'Payment not found' });
  
      if (response.data.data.status === 'success') {
        payment.status = 'completed';
      } else {
        payment.status = 'failed';
      }
  
      await payment.save();
      res.status(200).json(payment);
    } catch (error) {
      console.error('Payment verification error:', error);
      res.status(500).json({ error: 'Payment verification failed' });
    }
  };
  const generateReference = () => {
    return 'tx-' + Math.random().toString(36).substr(2, 9);
  };
  const categoryAndType=async(req,res)=>{
try {
    const {category,a}=req.body;
    const categoryAndType=await CategoryAndType.create({category:category.category,types:[...a]})
    if(categoryAndType){
       return res.json({success:'successfully created'})
    }
    if(!categoryAndType){
        return res.json({error:"The category does not created."})
    }
} catch (error) {
    console.log(error);
}
  }