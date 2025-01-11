<template>
  <div class="max-w-md mx-auto p-6 bg-white shadow-md rounded-lg">
    <h1 class="text-2xl font-bold mb-4">Payment Information</h1>
    <form @submit.prevent="handleSubmit">
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="amount">Amount</label>
        <input
          v-model="form.amount"
          type="number"
          id="amount"
          class="w-full p-2 border rounded-md"
          placeholder="Enter amount"
          required
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="email">Email</label>
        <input
          v-model="form.email"
          type="email"
          id="email"
          class="w-full p-2 border rounded-md"
          placeholder="Enter email"
          required
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="firstName">First Name</label>
        <input
          v-model="form.firstName"
          type="text"
          id="firstName"
          class="w-full p-2 border rounded-md"
          placeholder="Enter first name"
          required
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="lastName">Last Name</label>
        <input
          v-model="form.lastName"
          type="text"
          id="lastName"
          class="w-full p-2 border rounded-md"
          placeholder="Enter last name"
          required
        />
      </div>
      <div class="mb-4">
        <label class="block text-gray-700 font-bold mb-2" for="phoneNumber">Phone Number</label>
        <input
          v-model="form.phoneNumber"
          type="text"
          id="phoneNumber"
          class="w-full p-2 border rounded-md"
          placeholder="Enter phone number"
          required
        />
      </div>
      <button
        type="submit"
        class="w-full bg-blue-500 text-white p-2 rounded-md hover:bg-blue-600"
      >
       Pay
      </button>
    </form>

    <div v-if="result" class="mt-4 p-4 bg-yellow-100 text-yellow-700 rounded-md">
      <p><strong>Waiting...</strong></p>
    </div>
  </div>
</template>

<script setup>
import { useMutation } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { ref, watch } from 'vue';

const START_PAYMENT_MUTATION = gql`
  mutation startPayment($input: StartPaymentInput!) {
    start_payment(input: $input) {
      message
      status
      data
    }
  }
`;

const form = ref({
  amount: 100.0,
  email: 'telaynew11@gmail.com',
  firstName: 'Telaynew',
  lastName: 'Ambachew',
  phoneNumber: '0900123456',
});

const result = ref(null);

const { mutate, onDone } = useMutation(START_PAYMENT_MUTATION);

const handleSubmit = () => {
  mutate({
    input: {
      ...form.value,
    },
  });

  onDone(({ data }) => {
    result.value = data.start_payment;
  });
};

// Watcher to automatically navigate if checkout_url exists
watch(result, (newResult) => {
  if (newResult && newResult.data && newResult.data.checkout_url) {
    // Redirect to checkout_url automatically
    window.location.href = newResult.data.checkout_url;
  }
});
</script>

<style scoped>
/* Add any scoped styles if needed */
</style>
