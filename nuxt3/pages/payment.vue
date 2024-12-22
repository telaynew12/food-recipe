<template>
  <div>
    <h2>Make a Payment</h2>
    <form @submit.prevent="submitPayment">
      <input v-model="email" type="email" placeholder="Email" required />
      <input v-model="fname" type="text" placeholder="First Name" required />
      <input v-model="lname" type="text" placeholder="Last Name" required />
      <input v-model="amount" type="number" placeholder="Amount" required min="1" />
      <button type="submit" :disabled="loading">Pay Now</button>
    </form>

    <div v-if="loading">Processing payment...</div>

    <div v-if="paymentSuccess">
      <h3>Payment Successful!</h3>
      <p>Thank you for your payment. Transaction ID: {{ paymentSuccess.tx_id }}</p>
      <button @click="redirectToHome">Go to Home</button>
    </div>

    <div v-if="paymentError" class="error">
      <p>{{ paymentError }}</p>
    </div>
  </div>
</template>

<script setup>
import { ref } from "vue";
import axios from "axios";

const email = ref("");
const fname = ref("");
const lname = ref("");
const amount = ref(0);
const loading = ref(false);
const paymentSuccess = ref(null);
const paymentError = ref("");

const submitPayment = async () => {
  if (!email.value || !fname.value || !lname.value || amount.value <= 0) {
    alert("Please fill in all required fields correctly.");
    return;
  }

  loading.value = true;
  paymentError.value = "";
  paymentSuccess.value = null;

  try {
    const response = await axios.post("http://127.0.0.1:8000/pay", {
      email: email.value,
      fname: fname.value,
      lname: lname.value,
      amount: amount.value,
      rdurl: "http://localhost:3000/paymentComplete",
    });

    const paymentData = response.data;
    const checkoutUrl = paymentData.detail.data.checkout_url;

    if (checkoutUrl) {
      window.location.replace(checkoutUrl); // Redirect to Chapa checkout
    } else {
      paymentError.value = "Checkout URL not received.";
    }
  } catch (error) {
    console.error("Payment initialization failed:", error);
    paymentError.value =
      error.response?.data?.message || "An error occurred while initializing payment.";
  } finally {
    loading.value = false;
  }
};

const redirectToHome = () => {
  window.location.href = "http://localhost:3000/"; // Redirect to home page
};
</script>

<style>
.error {
  color: red;
  margin-top: 10px;
}
</style>
