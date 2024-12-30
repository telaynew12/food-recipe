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
import { computed, watch } from 'vue';
import gql from 'graphql-tag';
import { useCategoryStore } from '@/stores/categoryStore';
import { useRecipeStore } from '@/stores/recipe'; // Search query store

const backendBaseUrl = 'http://localhost:8085/';

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
const recipeStore = useRecipeStore(); // Search store
const selectedCategory = computed(() => categoryStore.selectedCategory);

// Search query state
const searchQuery = computed({
  get: () => recipeStore.searchQuery,
  set: (value) => recipeStore.setSearchQuery(value),
});

// Query recipes
const { result, loading, error, refetch } = useQuery(GET_RECIPES_BY_CATEGORIES, {
  categoryName: selectedCategory.value,
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

// Watch for category changes and refetch
watch(
  () => selectedCategory.value,
  (newCategory) => {
    if (newCategory) {
      refetch({ categoryName: newCategory });
    }
  },
  { immediate: true }
);
</script>
