from fastapi import FastAPI
from pydantic import BaseModel
from fastapi.middleware.cors import CORSMiddleware
from app.payment import Payment

# FastAPI app initialization
app = FastAPI()

# Add CORS middleware to allow cross-origin requests from the frontend
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:3000"],  # Frontend URL
    allow_credentials=True,
    allow_methods=["*"],  # All HTTP methods
    allow_headers=["*"],  # All headers
)

# Root endpoint to test the API
@app.get('/')
def main():
    return {'message': 'Hi'}

# Pydantic models to validate request payloads
class Payit(BaseModel):
    email: str
    fname: str
    lname: str
    amount: int
    rdurl: str

class Txnum(BaseModel):
    ref_num: str

# Route to initiate payment
@app.post('/pay')
def pay(payit: Payit):
    data = Payment.pay(
        email=payit.email, 
        fname=payit.fname, 
        lname=payit.lname, 
        amount=payit.amount, 
        rdurl=payit.rdurl
    )
    return data

# Route to verify the payment
@app.post('/verify')
def verify(txnum: Txnum):
    ver = Payment.verify(tx_num=txnum.ref_num)
    return ver

# Route to handle payment completion and verification
@app.get("/payment-complete")
def payment_complete(tx_ref: str):
    """
    Handles payment completion verification and returns the appropriate message.
    """
    try:
        # Verify the payment with the tx_ref
        verification_response = Payment.verify(tx_ref)
        
        # If the verification is successful
        if verification_response['status'] == 'success':
            return {
                "message": "Payment successful!",
                "tx_id": tx_ref,
                "payment_details": verification_response
            }
        else:
            return {
                "message": "Payment verification failed.",
                "tx_id": tx_ref,
                "details": verification_response
            }
    except Exception as e:
        return {"error": "An error occurred during payment verification.", "details": str(e)}
