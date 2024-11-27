<template>
  <div>
    <!-- Header -->
    <header class="bg-blue-600 text-white py-4">
      <div class="container mx-auto flex justify-between items-center px-4">
        <h1 class="text-xl font-bold">My Nuxt 3 App</h1>
        <nav>
          <ul class="flex space-x-4">
            <li><NuxtLink to="/login" class="hover:underline">Login</NuxtLink></li>
            <li><NuxtLink to="/signup" class="hover:underline">Signup</NuxtLink></li>
            <li><NuxtLink to="/logout" class="hover:underline">Logout</NuxtLink></li>
          </ul>
        </nav>
      </div>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto py-8 px-4">
      <NuxtPage />
    </main>

    <!-- Footer -->
    <footer class="bg-gray-800 text-white py-4 mt-8">
      <div class="container mx-auto text-center">
        <p>&copy; 2024 My Nuxt 3 App. All rights reserved.</p>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { NuxtLink } from "#components";
</script>

<style scoped>
/* Add any custom styles here if needed */
</style>
