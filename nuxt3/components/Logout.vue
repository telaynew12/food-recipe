<template>
  <button @click="logout" class="logout-btn">Logout</button>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const logout = () => {
  // Clear any authentication data, for example:
  localStorage.removeItem('authToken')

  // Redirect to the login page
  router.push('/login')
}
</script>

<style scoped>
.logout-btn {
  padding: 10px 20px;
  background-color: #ff4d4f;
  color: white;
  border: none;
  border-radius: 5px;
  cursor: pointer;
}

.logout-btn:hover {
  background-color: #e63946;
}
</style>
