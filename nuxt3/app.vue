<template>
  <NuxtLayout>
    
  <NuxtPage />
</NuxtLayout>
</template>


	<script setup>
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
authStore.initializeAuth(); // This will now run after Pinia is initialized
</script>


<style>
/* Add global styles here if needed */
</style>
