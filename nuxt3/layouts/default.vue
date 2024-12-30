<template>
  <div>
    <!-- Always Visible FoodRecipes Link -->
    <header class="left-0 bg-green-500 shadow-md fixed w-full z-0">
      <nav class="container mx-auto px-6 py-4 flex items-center justify-between">
        <!-- Logo: Always Visible -->
        <div>
          <a
            href="#"
            class="text-xl font-bold italic text-red-800 hover:text-blue-600"
            @click.prevent="navigateToHome"
          >
            FoodRecipes
          </a>
        </div>

        <!-- Conditionally Visible Header Content -->
        <div v-if="!isLoginOrSignupPage" class="flex-1 flex justify-between items-center">
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
                v-model="recipeStore.searchQuery"
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

            <template v-else>
              <div
                class="relative group"
                @mouseenter="showAccountDropdown = true"
                @mouseleave="showAccountDropdown = false"
              >
                <button
                  @click="toggleAccountDropdown"
                  class="flex items-center px-4 py-2 text-sm font-medium text-gray-600 hover:text-gray-800 focus:outline-none"
                >
                  <img
                    :src="authStore.user?.profilePicture || '/default-avatar.png'"
                    alt="Profile"
                    class="w-8 h-8 rounded-full mr-2"
                  />
                  My Account
                </button>
                <!-- Dropdown Menu with Transition -->
                <transition name="fade" @before-enter="beforeEnter" @enter="enter" @leave="leave">
                  <div
                    v-if="showAccountDropdown"
                    class="absolute right-0 mt-2 w-48 bg-white border border-gray-300 rounded-lg shadow-lg z-10"
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
                          @click="logout"
                        >
                          Logout
                        </button>
                      </li>
                    </ul>
                  </div>
                </transition>
              </div>
            </template>
          </div>
        </div>
      </nav>
    </header>

    <!-- Main Content -->
    <main class="container mx-auto py-24">
      <NuxtPage />
    </main>
<!-- Footer -->
<footer class="bg-gray-800 text-white py-6">
  <div class="container mx-auto">
    <!-- Footer Content -->
    <div class="flex flex-col md:flex-row justify-between items-center">
      <!-- About Us -->
      <div class="mb-4 md:mb-0">
        <p class="text-sm">&copy; 2024 FoodRecipes. All rights reserved.</p>
        <p class="text-sm">
          <a href="#about" class="hover:text-blue-400">About Us</a>
        </p>
      </div>

      <!-- Social Media Links -->
      <div class="flex space-x-4 mb-4 md:mb-0">
        <a href="https://facebook.com" target="_blank" class="hover:text-blue-400">Facebook</a>
        <a href="https://linkedin.com/in/telaynew-ambachew" target="_blank" class="hover:text-blue-400">LinkedIn</a>
        <a href="https://instagram.com" target="_blank" class="hover:text-pink-400">Instagram</a>
        <a href="https://twitter.com" target="_blank" class="hover:text-blue-300">Twitter</a>
      </div>

      <!-- Navigation Links -->
      <div class="flex space-x-4">
        <a href="#recipes" class="hover:text-blue-400">FoodRecipes</a>
        <a href="#top" class="hover:text-blue-400">Back to Top</a>
      </div>
    </div>
  </div>
</footer>


  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRecipeStore } from '@/stores/recipe';
import { useRouter, useRoute } from 'vue-router';

const authStore = useAuthStore();
const recipeStore = useRecipeStore();
const router = useRouter();
const route = useRoute();

const showFilter = ref(false);
const showAccountDropdown = ref(false);

const navigateToHome = () => router.push('/');
const goToSignup = () => router.push('/signup');
const goToLogin = () => router.push('/login');
const toggleFilter = () => (showFilter.value = !showFilter.value);
const applyFilter = (filterType) => (showFilter.value = false);
const searchRecipes = () => console.log('Search query:', recipeStore.searchQuery);
const logout = () => {
  authStore.setUser(null);
  authStore.isAuthenticated = false;
  localStorage.removeItem('token');
  sessionStorage.clear();
  router.push('/');
};

const toggleAccountDropdown = () => {
  showAccountDropdown.value = !showAccountDropdown.value;
};

const isLoginOrSignupPage = computed(() =>
  ['/login', '/signup', '/recipeDetail'].includes(route.path)
);

// Transition Methods for Smooth Dropdown
const beforeEnter = (el) => {
  el.style.opacity = 0;
  el.style.transition = 'opacity 0.5s ease';
};

const enter = (el) => {
  el.offsetHeight; // Trigger reflow
  el.style.opacity = 1;
};

const leave = (el) => {
  el.style.opacity = 0;
};
</script>

<style scoped>
/* Optional: Additional styles for smooth transitions */
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.000005s ease;
}
.fade-enter, .fade-leave-to {
  opacity: 0;
}
</style>
