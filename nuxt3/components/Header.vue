// components/Header.vue
<script setup>
import { useAuth } from '~/composables/useAuth'; // Import the auth state

const isAuthenticated = useAuth(); // Access global authentication state

function logout() {
  localStorage.removeItem('token'); // Clear token
  isAuthenticated.value = false; // Update auth state
}
</script>

<template>
  <nav class="p-4 bg-gray-200">
    <ul class="flex space-x-4">
      <!-- Show login and signup links if not authenticated -->
      <li v-if="!isAuthenticated"><a href="/login">Login</a></li>
      <li v-if="!isAuthenticated"><a href="/signup">Signup</a></li>

      <!-- Show dashboard and logout links if authenticated -->
      <li v-if="isAuthenticated"><a href="/dashboard">Dashboard</a></li>
      <li v-if="isAuthenticated"><a href="#" @click="logout">Logout</a></li>
    </ul>
  </nav>
</template>
