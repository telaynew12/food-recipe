<template>
  <div class="flex min-h-screen bg-gray-100">
    <!-- Sidebar -->
    <aside class="w-64 bg-blue-600 text-white flex flex-col">
      <div class="py-6 px-4 text-center">
        <h2 class="text-xl font-bold">User Dashboard</h2>
      </div>
      <nav class="flex-1">
        <ul>
          <li>
            <NuxtLink
              to="/dashboard/overview"
              class="block px-4 py-2 hover:bg-blue-700"
            >
              Overview
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/dashboard/profile"
              class="block px-4 py-2 hover:bg-blue-700"
            >
              Profile
            </NuxtLink>
          </li>
          <li>
            <NuxtLink
              to="/dashboard/settings"
              class="block px-4 py-2 hover:bg-blue-700"
            >
              Settings
            </NuxtLink>
          </li>
          <!-- Logout -->
          <li>
            <Logout /> <!-- Use the Logout component -->
          </li>
        </ul>
      </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <!-- Header -->
      <header class="bg-white shadow py-4 px-6">
        <div class="flex justify-between items-center">
          <h1 class="text-lg font-semibold">Welcome Back, [User Name]</h1>
          <NuxtLink
            to="/dashboard/notifications"
            class="text-blue-600 hover:underline"
          >
            Notifications
          </NuxtLink>
        </div>
      </header>

      <!-- Content -->
      <main class="flex-1 p-6">
        <h2 class="text-2xl font-bold mb-4">Dashboard Overview</h2>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <!-- Example Card -->
          <div class="bg-white shadow rounded-lg p-4">
            <h3 class="text-lg font-semibold mb-2">Total Recipes</h3>
            <p class="text-4xl font-bold text-blue-600">42</p>
          </div>
          <div class="bg-white shadow rounded-lg p-4">
            <h3 class="text-lg font-semibold mb-2">Bookmarked Recipes</h3>
            <p class="text-4xl font-bold text-green-600">10</p>
          </div>
          <div class="bg-white shadow rounded-lg p-4">
            <h3 class="text-lg font-semibold mb-2">New Comments</h3>
            <p class="text-4xl font-bold text-red-600">5</p>
          </div>
        </div>
      </main>
    </div>
  </div>
</template>

<script setup>
// Import the Logout component from components
import Logout from '@/components/Logout.vue';
</script>

<style scoped>
/* Add custom styles here if needed */
</style>
