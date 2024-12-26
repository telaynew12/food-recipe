<template>
    <div>
      <h2>Payment Status</h2>
      <div v-if="loading">Verifying payment...</div>
      <div v-if="paymentSuccess">
        <h3>Payment Successful!</h3>
        <p>Transaction ID: {{ paymentSuccess.tx_id }}</p>
        <button @click="redirectToHome">Go to Home</button>
      </div>
      <div v-if="paymentError" class="error">
        <p>{{ paymentError }}</p>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from "vue";
  import { useRoute } from "vue-router";
  import axios from "axios";
  
  const route = useRoute();
  const txRef = route.query.tx_ref; // Get transaction reference from query params
  const loading = ref(true);
  const paymentSuccess = ref(null);
  const paymentError = ref("");
  
  onMounted(async () => {
    try {
      const response = await axios.get(`http://127.0.0.1:8000/payment-complete?tx_ref=${txRef}`);
      paymentSuccess.value = response.data;
    } catch (error) {
      console.error("Payment verification failed:", error);
      paymentError.value =
        error.response?.data?.message || "An error occurred while verifying payment.";
    } finally {
      loading.value = false;
    }
  });
  
  const redirectToHome = () => {
    window.location.href = "http://localhost:3000/";
  };
  </script>
         
  <style>
  .error {
    color: red;
    margin-top: 10px;
  }
  </style>
  