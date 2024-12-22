<template>
    <div>
      <h2>Payment</h2>
      <form @submit.prevent="submitPayment">
        <input v-model="email" type="email" placeholder="Email" required />
        <input v-model="fname" type="text" placeholder="First Name" required />
        <input v-model="lname" type="text" placeholder="Last Name" required />
        <input v-model="amount" type="number" placeholder="Amount" required />
        <input v-model="rdurl" type="url" placeholder="Return URL" required />
        <button type="submit">Pay Now</button>
      </form>
  
      <div v-if="paymentResponse">
        <h3>Payment Details</h3>
        <p>Transaction ID: {{ paymentResponse.tx_id }}</p>
        <p>Payment Status: {{ paymentResponse.detail.status }}</p>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import axios from 'axios'
  
  const email = ref('')
  const fname = ref('')
  const lname = ref('')
  const amount = ref(0)
  const rdurl = ref('')
  
  const paymentResponse = ref(null)
  
  const submitPayment = async () => {
    try {
      const response = await axios.post('http://127.0.0.1:8000/pay', {
        email: email.value,
        fname: fname.value,
        lname: lname.value,
        amount: amount.value,
        rdurl: rdurl.value
      })
      paymentResponse.value = response.data
    } catch (error) {
      console.error('Payment failed:', error)
    }
  }
  </script>
  