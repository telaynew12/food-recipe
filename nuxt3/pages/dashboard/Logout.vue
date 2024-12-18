<template>
  <button @click="logout" class="logout-btn">Logout</button>
</template>

<script setup>
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth' // Pinia store for authentication
import { useRuntimeConfig } from '#app'

const router = useRouter()
const authStore = useAuthStore()
const runtimeConfig = useRuntimeConfig()

const logout = async () => {
  try {
    // Retrieve the token from Pinia or localStorage
    const token = authStore.token || localStorage.getItem('authToken')

    // Send a logout request to the backend if the token exists
    if (token) {
      await $fetch(`${runtimeConfig.public.apiBase}/logout`, {
        method: 'POST',
        headers: {
          Authorization: `Bearer ${token}`,
        },
      })
    }

    // Clear token in both Pinia and localStorage
    authStore.clearToken()
    localStorage.removeItem('authToken')

    // Redirect the user to the login page
    router.push('/login')
  } catch (error) {
    console.error('Error during logout:', error)
    alert('Failed to log out. Please try again.')
  }
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
  font-size: 16px;
  transition: background-color 0.3s ease;
}

.logout-btn:hover {
  background-color: #e63946;
}
</style>
