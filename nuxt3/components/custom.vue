<template>
  <div class="p-6 max-w-4xl mx-auto bg-white shadow-md rounded-lg">
    <!-- Display All Recipes -->
    <div v-if="recipes.length" class="space-y-8">
      <div v-for="(recipe, index) in recipes" :key="recipe.id" class="border-b pb-6">
        <!-- Recipe Title -->
        <h2 class="text-2xl font-bold mb-4">{{ recipe.title }}</h2>

        <!-- Recipe Image -->
        <img
          v-if="recipe.featured_image"
          :src="getImageUrl(recipe.featured_image)"
          alt="Recipe Image"
          class="w-full h-500 object-cover rounded-lg mb-4"
        />

        <!-- Buttons Section (flex container) -->
        <div class="flex justify-start gap-6 mb-4 items-center">
          <!-- Like Button -->
          <button
            @click="handleLike(index)"
            :class="{ 'text-blue-500 font-bold': recipe.liked, 'text-gray-600': !recipe.liked }"
            class="flex items-center space-x-2"
          >
            <span>👍</span>
            <span>Like</span>
            <span>({{ recipe.likesCount }})</span>
          </button>

          <!-- Comment Button -->
          <button
            @click="toggleComments(index)"
            class="flex items-center space-x-2 text-gray-600 hover:text-gray-800"
          >
            <span>💬</span>
            <span>Comments</span>
            <span>({{ recipe.comments.length }})</span>
          </button>

          <!-- Bookmark Button -->
          <button
            @click="handleBookmark(index)"
            :class="{ 'text-yellow-500': recipe.bookmarked, 'text-gray-600': !recipe.bookmarked }"
            class="flex items-center space-x-2"
          >
            <span>🔖</span>
            <span>{{ recipe.bookmarked ? 'Bookmarked' : 'Bookmark' }}</span>
          </button>

          <!-- Rating Section -->
          <div class="flex items-center space-x-1">
            <span class="font-semibold">  </span>
            <div class="flex">
              <span
                v-for="star in 5"
                :key="star"
                @click="handleRating(index, star)"
                :class="{
                  'text-yellow-400': recipe.rating >= star,
                  'text-gray-400': recipe.rating < star
                }"
                class="cursor-pointer text-2xl"
              >
                ★
              </span>
            </div>
            <span class="ml-2 text-sm text-gray-600">({{ recipe.rating }} / 5)</span>
          </div>
        </div>

        <!-- Comments Section -->
        <div v-if="recipe.showComments" class="mt-4">
          <h3 class="text-lg font-semibold mb-2">Comments</h3>
          <div
            v-for="(comment, i) in recipe.comments"
            :key="i"
            class="p-2 border-b"
          >
            <p>{{ comment.text }}</p>
          </div>
          <!-- Add Comment -->
          <div class="mt-2">
            <textarea
              v-model="recipe.newComment"
              placeholder="Write a comment..."
              class="w-full p-2 border rounded"
            ></textarea>
            <button
              @click="addComment(index)"
              class="bg-blue-500 text-white px-4 py-2 rounded mt-2 hover:bg-blue-600"
            >
              Submit Comment
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- No Recipes Available -->
    <div v-else>
      <p class="text-center">No recipes available.</p>
    </div>

    <!-- Log In Prompt -->
    <div v-if="!isAuthenticated" class="text-center mt-4">
      <p>Please log in to like, comment, rate, or buy recipes.</p>
      <button @click="loginPrompt" class="bg-blue-500 text-white px-4 py-2 rounded mt-2 hover:bg-blue-600">
        Log In
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useNuxtApp } from '#app';
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import gql from 'graphql-tag';
const router = useRouter();

const { $apolloClient } = useNuxtApp();
const recipes = ref([]);
const authStore = useAuthStore(); // Access auth store

const isAuthenticated = computed(() => authStore.isAuthenticated); // Reactive check for authentication

const backendBaseUrl = 'http://localhost:8085/';

// GraphQL Query to fetch shared recipes
const FETCH_SHARED_RECIPES_QUERY = gql`
  query FetchSharedRecipes {
    recipes(where: { share_statuses: { is_shared: { _eq: true } } }) {
      id
      title
      featured_image
    }
  }
`;

// Helper function to generate the full image URL
const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : null;
};

// Fetch shared recipes from the backend
const fetchRecipes = async () => {
  try {
    const response = await $apolloClient.query({
      query: FETCH_SHARED_RECIPES_QUERY,
    });

    // Initialize state for each recipe (likes, comments, etc.)
    recipes.value = response.data.recipes.map((recipe) => ({
      ...recipe,
      likesCount: 0,
      liked: false,
      bookmarked: false,
      rating: 0, // Added for rating
      showComments: false,
      comments: [],
      newComment: '',
    }));
  } catch (error) {
    console.error('Error fetching recipes:', error);
  }
};

// Login prompt or login action
const loginPrompt = () => {
  // You can redirect to login page, open a modal, or trigger a login function
  alert('Redirecting to login page...');
  // Simulate login for now:
  // On successful login, you should call:
  // authStore.setUser(user); to store the user and update authentication status.
};

// Interactive Handlers
const handleLike = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = recipes.value[index];
  recipe.likesCount += recipe.liked ? -1 : 1;
  recipe.liked = !recipe.liked;
};

const toggleComments = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  recipes.value[index].showComments = !recipes.value[index].showComments;
};

const addComment = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = recipes.value[index];
  if (recipe.newComment.trim()) {
    recipe.comments.push({ text: recipe.newComment });
    recipe.newComment = '';
  }
};

const handleBookmark = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  recipes.value[index].bookmarked = !recipes.value[index].bookmarked;
};

const handleRating = (index, rating) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = recipes.value[index];
  recipe.rating = rating; // Set the rating when a star is clicked
};

const handleBuy = (title) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  alert(`Redirecting to purchase page for "${title}"...`);
};

// Fetch recipes when the component is mounted
onMounted(() => {
  fetchRecipes();
});
</script>

<style scoped>
/* Add custom styles if needed */
</style>
