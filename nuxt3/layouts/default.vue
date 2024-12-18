<template>
  <div>
    <!-- Header with Navbar -->
    <header class="bg-white shadow-md fixed w-full z-5">
      <nav class="container mx-auto px-6 py-4 flex items-center justify-between">
        <!-- Logo -->
        <div>
          <a
            href="#"
            class="text-xl font-bold text-gray-800 hover:text-gray-600"
            @click.prevent="navigateToHome"
          >
            FoodRecipes
          </a>
        </div>

        <!-- Filter + Search Bar -->
        <div class="flex items-center flex-1 mx-4">
          <!-- Filter Dropdown -->
          <div class="relative mr-4">
            <button
              @click="toggleFilter"
              class="flex items-center px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring focus:ring-blue-300 hover:bg-gray-100"
            >
              <span class="mr-2">Filter</span>
              <svg
                xmlns="http://www.w3.org/2000/svg"
                fill="none"
                viewBox="0 0 24 24"
                stroke-width="1.5"
                stroke="currentColor"
                class="w-4 h-4"
              >
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  d="M12 3c.828 0 1.5.672 1.5 1.5v15c0 .828-.672 1.5-1.5 1.5H9c-.828 0-1.5-.672-1.5-1.5V4.5C7.5 3.672 8.172 3 9 3h3zM15 3c.828 0 1.5.672 1.5 1.5v15c0 .828-.672 1.5-1.5 1.5H9c-.828 0-1.5-.672-1.5-1.5V4.5C7.5 3.672 8.172 3 9 3h6z"
                />
              </svg>
            </button>
            <!-- Dropdown Menu -->
            <div
              v-if="showFilter"
              class="absolute left-0 mt-2 w-48 bg-white border border-gray-300 rounded-lg shadow-lg z-10"
            >
              <ul class="py-2 text-sm text-gray-700">
                <li>
                  <button
                    class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                    @click="applyFilter('Vegetarian')"
                  >
                    Vegetarian
                  </button>
                </li>
                <li>
                  <button
                    class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                    @click="applyFilter('Non-Vegetarian')"
                  >
                    Non-Vegetarian
                  </button>
                </li>
                <li>
                  <button
                    class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                    @click="applyFilter('Vegan')"
                  >
                    Vegan
                  </button>
                </li>
              </ul>
            </div>
          </div>

          <!-- Search Bar -->
          <form class="relative flex-1" @submit.prevent="searchRecipes">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search recipes..."
              class="w-full px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring focus:ring-blue-300"
            />
            <button
              type="submit"
              class="absolute right-2 top-1/2 transform -translate-y-1/2 text-blue-500 hover:text-gray-800"
            >
              🔍
            </button>
          </form>
        </div>

        <!-- Authentication Section -->
        <div class="flex items-center space-x-4">
          <!-- Before Login: Signup and Login -->
          <template v-if="!authStore.isAuthenticated">
            <button
              class="px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800 focus:outline-none"
              @click="goToSignup"
            >
              Signup
            </button>
            <button
              class="px-4 py-2 text-sm font-medium text-white bg-blue-500 rounded-lg shadow hover:bg-blue-600 focus:outline-none"
              @click="goToLogin"
            >
              Login
            </button>
          </template>

          <!-- After Login: My Account Dropdown -->
          <template v-else>
            <div class="relative group">
              <button
                class="flex items-center px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800 focus:outline-none"
              >
                <img
                  :src="authStore.user?.profilePicture || '/default-avatar.png'"
                  alt="Profile"
                  class="w-8 h-8 rounded-full mr-2"
                />
                My Account
              </button>
              <!-- Dropdown Menu -->
              <div
                class="absolute right-0 mt-2 w-48 bg-white border border-gray-300 rounded-lg shadow-lg z-10 opacity-0 group-hover:opacity-100 transition-opacity"
              >
                <ul class="py-2 text-sm text-gray-700">
                  <li>
                    <button
                      class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                      @click="navigateTo('/myrecipes')"
                    >
                      My Recipes
                    </button>
                  </li>
                  <li>
                    <button
                      class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                      @click="navigateTo('/create-recipe')"
                    >
                      Create Recipe
                    </button>
                  </li>
                  <li>
                    <button
                      class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                      @click="navigateTo('/profile')"
                    >
                      My Profile
                    </button>
                  </li>
                  <li>
                    <button
                      class="w-full px-4 py-2 text-red-500 hover:bg-gray-100 text-left"
                      @click="navigateTo('/login')"
                    >
                      Logout
                    </button>
                  </li>
                </ul>
              </div>
            </div>
          </template>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto py-24">
        <NuxtPage />
    </main>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth';

const authStore = useAuthStore();
const showFilter = ref(false);
const searchQuery = ref('');

const navigateToHome = () => {
  if (authStore.isAuthenticated) {
    navigateTo('/');
  }
};

const goToSignup = () => {
  navigateTo('/signup');
};

const goToLogin = () => {
  navigateTo('/login');
};

const toggleFilter = () => {
  showFilter.value = !showFilter.value;
};

const applyFilter = (filterType) => {
  console.log(`Filter applied: ${filterType}`);
  showFilter.value = false;
};

const searchRecipes = () => {
  console.log(`Searching for: ${searchQuery.value}`);
};
</script>
