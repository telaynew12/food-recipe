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
          <li>
            <button @click="logout" class="block px-4 py-2 hover:bg-blue-700">
              Logout
            </button>
          </li>
        </ul>
      </nav>
    </aside>

    <!-- Main Content -->
    <div class="flex-1 flex flex-col">
      <header class="bg-white shadow py-4 px-6">
        <div class="flex justify-between items-center">
          <h1 class="text-lg font-semibold">
            Welcome Back, {{ user.name }}
          </h1>
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

        <div class="bg-white shadow rounded-lg p-4 mb-4">
          <h3 class="text-lg font-semibold mb-2">Your User Details</h3>
          <p><strong>User ID:</strong> {{ user.userId }}</p>
          <p><strong>Name:</strong> {{ user.name }}</p>
          <p><strong>Email:</strong> {{ user.email }}</p>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
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
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth'; // Pinia store for user data
import { ref, onMounted } from 'vue';

// Get route and auth store
const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

// State for user data
const user = ref({
  name: '',
  email: '',
  userId: ''
});

// On mounted, fetch user details from the store, query params, or localStorage
onMounted(() => {
  // Check for user details in the query params
  const userIdFromQuery = route.query.userId;
  const nameFromQuery = route.query.name;
  const emailFromQuery = route.query.email;

  if (userIdFromQuery && nameFromQuery && emailFromQuery) {
    // If userId, name, and email are passed via query, set them
    user.value.userId = userIdFromQuery;
    user.value.name = nameFromQuery;
    user.value.email = emailFromQuery;
  } else if (auth.user) {
    // Otherwise, use stored user info from Pinia
    user.value = auth.user;
  } else {
    // If no user is found, check localStorage
    const storedUser = JSON.parse(localStorage.getItem('user'));
    if (storedUser) {
      user.value = storedUser;
    }
  }
});

// Handle logout action
const logout = () => {
  auth.clearUser();
  router.push('/login'); // Redirect to login page
};
</script>

<style scoped>
/* Add custom styles here */
</style>
