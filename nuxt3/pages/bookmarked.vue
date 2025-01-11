<template>
  <div>
    <p class="flex items-center-top justify-center">
      <h1 v-if="user">Welcome <strong>{{ user.name }}</strong></h1>
    </p>

    <div class="p-4">
      <h1 class="text-2xl font-bold mb-4">Recipes</h1>

      <!-- Loading State -->
      <div v-if="loadingBookmarks" class="text-center">
        <p>Loading recipes...</p>
      </div>

      <!-- Error State -->
      <div v-if="errorBookmarks" class="text-red-500">
        <p>Error: {{ errorBookmarks.message }}</p>
      
        <p v-if="errorBookmarks.message.includes('User ID not found')">It seems the user ID is invalid or missing.</p>
      </div>

      <!-- Recipes List -->
      <div v-if="bookmarkedRecipes.length" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="recipe in bookmarkedRecipes"
          :key="recipe.id"
          class="bg-white rounded-lg shadow-md p-4"
        >
         
          <!-- Recipe Details -->
          <div class="mt-4">
            <p > <strong>Title:</strong     >  {{ recipe.title }}</p>
            <p class="text-lg font-bold"> Category:{{recipe.category.name}}</p>
          </div>

      
        </div>
      </div>

      <!-- No Bookmarked Recipes -->
      <div v-if="!loadingBookmarks && !bookmarkedRecipes.length" class="text-center text-gray-500">
        <p>No bookmarked recipes available.</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';

// State to hold user and recipe information
const authStore = useAuthStore();
const user = ref(null); // Single user object
const bookmarkedRecipes = ref([]); // Array of bookmarked recipes
const loadingUser = ref(true); // Loading state for user data
const loadingBookmarks = ref(true); // Loading state for bookmarks
const errorUser = ref(null); // Error for user data
const errorBookmarks = ref(null); // Error for bookmarks

// Access userId directly from the auth store
const userId = authStore.userId;

// Access Apollo client from Nuxt app
const { $apolloClient } = useNuxtApp();

// GraphQL query to fetch user details
const FETCH_USER_DETAILS = gql`
  query FetchUserDetails($id: uuid!) {
    users(where: { id: { _eq: $id } }) {
      id
      name
      email
      verification_code
      role
      created_at
      updated_at
    }
  }
`;

// GraphQL query to fetch bookmarked recipes
const FETCH_BOOKMARKED_RECIPES = gql`
  query FetchBookmarkedRecipes($userId: uuid!) {
    bookmarks(where: { user_id: { _eq: $userId } }) {
      id
      recipe_id
      recipe {
        id
        title
        category {
        id
        name
      }
        
      }
    }
  }
`;

const fetchUserDetails = async () => {
  try {
    const response = await $apolloClient.query({
      query: FETCH_USER_DETAILS,
      variables: { id: userId },
    });
    user.value = response.data.users[0];
  } catch (err) {
    console.error("Error fetching user details:", err);
    errorUser.value = err;
  } finally {
    loadingUser.value = false;
  }
};

const fetchBookmarkedRecipes = async () => {
  try {
    const response = await $apolloClient.query({
      query: FETCH_BOOKMARKED_RECIPES,
      variables: { userId },
    });
    bookmarkedRecipes.value = response.data.bookmarks.map((bookmark) => bookmark.recipe) || [];
  } catch (err) {
    console.error("Error fetching bookmarked recipes:", err);
    errorBookmarks.value = err;
  } finally {
    loadingBookmarks.value = false;
  }
};

// Fetch both data sets on component mount
onMounted(() => {
  if (userId) {
    fetchUserDetails();
    fetchBookmarkedRecipes();
  } else {
    errorUser.value = new Error('User ID not found');
    errorBookmarks.value = new Error('User ID not found');
    loadingUser.value = false;
    loadingBookmarks.value = false;
  }
});
</script>
