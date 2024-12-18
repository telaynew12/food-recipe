<template>
  <div>
    <h1>User Details</h1>
    <div v-if="loading">Loading...</div>
    <div v-else-if="error">Error: {{ error.message }}</div>
    <div v-else>
      <div v-for="user in users" :key="user.id" class="user">
        <h2><strong>Name:</strong>{{ user.name }}</h2>
        <p><strong>Email:</strong> {{ user.email }}</p>
        <!-- <p><strong>Role:</strong> {{ user.role }}</p> -->
        <!-- <p><strong>Joined at: {{ user.created_at }}</strong></p>
        <p><strong>Last updated at: {{ user.updated_at }}</strong></p> -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';

// State to hold user information
const userId = ref(''); // Hold the user ID
const users = ref([]);
const loading = ref(true);
const error = ref(null);

const { $apolloClient } = useNuxtApp();

// GraphQL query to fetch user details
const FETCH_USERS_QUERY = gql`
  query FetchUsers($id: uuid!) {
    users(where: { id: { _eq: $id } }) {
      id
      name
      email
      verification_code
      role
      created_at
      updated_at
    }
  }
`;

// Logic to retrieve userId from JWT token
onMounted(async () => {
  const token = localStorage.getItem('token'); // Retrieve the token from localStorage
  
  if (token) {
    try {
      const decodedToken = JSON.parse(atob(token.split('.')[1])); // Decode the JWT token
      userId.value = decodedToken.userId; // Set the userId from the token payload
    } catch (err) {
      console.error('Failed to decode token', err);
    }
  }

  // Fetch user details if userId is available
  if (userId.value) {
    try {
      const response = await $apolloClient.query({
        query: FETCH_USERS_QUERY,
        variables: { id: userId.value },
      });
      users.value = response.data.users;
    } catch (err) {
      error.value = err;
    } finally {
      loading.value = false;
    }
  } else {
    error.value = new Error('User ID not found');
    loading.value = false;
  }
});
</script>

<style scoped>
.user {
  margin-bottom: 20px;
  border: 1px solid #ccc;
  padding: 10px;
  border-radius: 5px;
}
</style>
