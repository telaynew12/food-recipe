<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-4">Recipes</h1>
  
    <!-- Loading, Error, or Recipes -->
    <div v-if="loading" class="text-center">Loading recipes...</div>
    <div v-else-if="error" class="text-red-500 text-center">{{ error.message }}</div>
    <div v-else>
      <div v-if="filteredRecipes.length === 0" class="text-center text-gray-500">
        No recipes found matching your search.
      </div>
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="recipe in filteredRecipes"
          :key="recipe.id"
          class="bg-white rounded-lg shadow p-4 hover:shadow-lg transition"
        >
          <img
            :src="getImageUrl(recipe.featured_image)"     
            alt="Recipe Image"
            class="w-full h-40 object-cover rounded-t-lg"
          />
          <h2 class="text-xl font-semibold mt-4">{{ recipe.title }}</h2>
          <p class="text-gray-600 mt-2">{{ recipe.description }}</p>
          <p class="text-gray-500 mt-2">
            Category: <span class="font-bold">{{ recipe.category.name }}</span>
          </p>
          <p class="text-gray-500">Preparation Time: {{ recipe.preparation_time }} mins</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useQuery } from '@vue/apollo-composable';
import { computed, ref, watchEffect, onMounted } from 'vue';
import gql from 'graphql-tag';
import { useCategoryStore } from '@/stores/categoryStore';
import { useRecipeStore } from '@/stores/recipe'; // Search query store

const backendBaseUrl = 'http://localhost:8085/';

// Function to get image URL
const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : null;
};

// GraphQL query
const GET_RECIPES_BY_CATEGORIES = gql`
  query GetRecipesByCategory($categoryName: String!) {
    recipes(where: { category: { name: { _eq: $categoryName } } }) {
      id
      title
      description
      category {
        id
        name
      }
      preparation_time
      featured_image
    }
  }
`;

// Access stores
const categoryStore = useCategoryStore();
const recipeStore = useRecipeStore(); // Search query store

// Search query state
const searchQuery = computed({
  get: () => recipeStore.searchQuery,
  set: (value) => recipeStore.setSearchQuery(value),
});

// Provide a default category value if selectedCategory is null
const currentCategory = computed(() => categoryStore.selectedCategory || "DefaultCategory");

// Reactive property to track whether the category is loaded
const isCategoryLoaded = ref(false);

// Query recipes using the currentCategory value
const { result, loading, error, refetch } = useQuery(GET_RECIPES_BY_CATEGORIES, {
  categoryName: currentCategory.value, // Use currentCategory which has a default value
  skip: !isCategoryLoaded.value, // Skip the query until the category is loaded
});

// Recipes
const recipes = computed(() => result.value?.recipes || []);

// Filter recipes by search query
const filteredRecipes = computed(() => {
  if (!searchQuery.value) {
    return recipes.value;
  }
  return recipes.value.filter((recipe) =>
    recipe.title.toLowerCase().includes(searchQuery.value.toLowerCase())
  );
});

// Watch for changes in categoryStore.selectedCategory and refetch recipes
watchEffect(() => {
  if (categoryStore.selectedCategory) {
    isCategoryLoaded.value = true; // Mark category as loaded
    refetch({ categoryName: categoryStore.selectedCategory });
  } else {
    isCategoryLoaded.value = false; // If no category is selected, skip fetching
  }
});

// Ensure category is loaded on page load
onMounted(() => {
  if (categoryStore.selectedCategory) {
    isCategoryLoaded.value = true;
  }
});
</script>





