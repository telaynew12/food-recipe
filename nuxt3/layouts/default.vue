<template>
  <div class="flex flex-col min-h-screen">
    <!-- Always Visible FoodRecipes Link -->
    <header class="left-0 bg-gray-300 shadow-md fixed w-full z-50">
      <nav class="container mx-auto px-6 py-4 flex items-center justify-between">
        <!-- Logo: Always Visible -->
        <div>
          <a
            href="#"
            class="text-xl font-bold italic text-red-500 hover:text-blue-600"
            @click.prevent="navigateToHome"
          >
            FoodRecipes
          </a>
        </div>
        <!-- Conditionally Visible Header Content -->
        <div v-if="isHomePage" class="relative mr-4 group" @mouseenter="showDropdown = true" @mouseleave="startHideTimeout">
          <!-- Button to trigger dropdown -->
          <button
            class="flex items-center px-4 py-2 border border-gray-300 rounded-lg shadow-sm focus:outline-none focus:ring-2 focus:ring-blue-500 hover:bg-gray-100 hover:shadow-md transition duration-300 ease-in-out"
            @click="toggleDropdown"
          >
            <span class="mr-2">Categories</span>
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
          <!-- Dropdown menu -->
          <div
            v-if="showDropdown"
            class="absolute left-0 mt-2 w-48 bg-white border border-gray-300 rounded-lg shadow-lg z-10 opacity-100 pointer-events-auto transition-opacity duration-300 ease-in-out"
          >
            <ul class="py-2 text-sm text-gray-700">
              <li
                v-for="category in categories"
                :key="category.id"
                class="transition-colors duration-200 ease-in-out"
              >
                <button
                  class="w-full px-4 py-2 text-left hover:bg-blue-50 focus:outline-none focus:bg-blue-100"
                  @click="handleCategoryClick(category.name)"
                >
                  {{ category.name }}
                </button>
              </li>
            </ul>
          </div>
        </div>
        <div v-if="!isLoginOrSignupPage" class="flex-1 flex justify-between items-center">
          <!-- Filter + Search Bar -->
          <div class="flex items-center flex-1 mx-4">
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
                  <!-- Circle with the first letter of the user's name -->
                  <div
                    class="w-8 h-8 rounded-full bg-red-500 flex items-center justify-center text-white font-medium mr-2"
                  >
                    {{ capitalizeFirstLetter(users[0]?.name?.charAt(0) || 'U') }}
                  </div>
                  {{ capitalizeFirstLetter(users[0]?.name || 'User') }}
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
                          @click="navigateTo('/Bookmarked')"
                        >
                          Bookmarked Recipes
                        </button>
                      </li>
                      <li>
                        <button
                          class="w-full px-4 py-2 hover:bg-gray-100 text-left"
                          @click="navigateTo('recipemangment')"
                        >
                          MY Recipes Interaction </button>
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
    <main class="container mx-auto py-24 flex-1">
      <NuxtPage />
    </main>

    <!-- Footer -->
    <footer class="bg-gray-800 text-white py-6 mt-auto">
      <div class="container mx-auto text-center">
        <p>&copy; 2025 FoodRecipes. All rights reserved.</p>
        <div class="flex justify-center space-x-6 mt-4">
          <a href="https://www.linkedin.com/in/telaynew-ambachew-44950825a/" target="_blank" class="text-gray-400 hover:text-gray-200">
            LinkedIn
          </a>
          <a href="https://www.github.com/telaynew12" target="_blank" class="text-gray-400 hover:text-gray-200">
            GitHub
          </a>
          <a href="mailto:contact@foodrecipes.com" class="text-gray-400 hover:text-gray-200">
            telaynew11@gmail.com
          </a>
        </div>
        <p class="mt-4 text-sm text-gray-500">Developed by Telaynew Ambachew</p>
      </div>
    </footer>
  </div>
</template>


<script setup>
import { ref, computed } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRecipeStore } from '@/stores/recipe';
import { useRouter, useRoute } from 'vue-router';
import { useCategoryStore } from '@/stores/categoryStore'; // Assuming Pinia is used

import { useNuxtApp } from '#app';

// GraphQL query for categories
const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`;

// Initialize reactive data
const { result, loading, error } = useQuery(GET_CATEGORIES);
const authStore = useAuthStore();
const recipeStore = useRecipeStore();
const categoryStore = useCategoryStore();
const router = useRouter();
const route = useRoute();

const showDropdown = ref(false);
const showAccountDropdown = ref(false);
const users = ref([]);

// Computed properties
const categories = computed(() => result.value?.categories || []);
const isHomePage = computed(() => route.path === '/');
const isLoginOrSignupPage = computed(() =>
  ['/login', '/signup', '/recipeDetail','/payment-success'].includes(route.path)
);

// Functions
// Capitalize the first letter of a string
const capitalizeFirstLetter = (str) => {
  return str.charAt(0).toUpperCase() + str.slice(1);
};
const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value;
};

const toggleAccountDropdown = () => {
  showAccountDropdown.value = !showAccountDropdown.value;
};

const handleCategoryClick = (categoryName) => {
  categoryStore.setSelectedCategory(categoryName);
  router.push({ name: 'CategoryPage', params: { category: categoryName } });
};

const navigateToHome = () => router.push('/');
const goToSignup = () => router.push('/signup');
const goToLogin = () => router.push('/login');
const logout = () => {
  authStore.clearUser();
  localStorage.removeItem('token');
  router.push('/login');
};

onMounted(async () => {
  const userId = authStore.userId;
  if (userId) {
    try {
      const { $apolloClient } = useNuxtApp();
      const response = await $apolloClient.query({
        query: gql`
          query FetchUsers($id: uuid!) {
            users(where: { id: { _eq: $id } }) {
              id
              name
              
            }
          }
        `,
        variables: { id: userId },
      });
      users.value = response.data.users;
    } catch (err) {
      console.error('Error fetching user details:', err);
    }
  }
});
</script>

<style scoped>
.footer {
      background-color: #2d3748; /* Tailwind's bg-gray-800 color */
      color: white;
      text-align: center;
      padding: 20px 0;
    }

/* Additional scoped styles can go here */
</style>
