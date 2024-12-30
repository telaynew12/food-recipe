<template>
  <div>
    <!-- Categories Navigation -->
    <nav class="fixed top-20 left-0 w-full bg-gray-500 p-4 md:p-6 z-10">
      <div class="container mx-auto flex justify-between items-center">
        <ul class="flex flex-wrap gap-4 md:gap-6 text-lg text-white">
          <!-- Loading Categories -->
          <li v-if="loadingCategories">Loading categories...</li>
          <!-- No Categories Found -->
          <li v-else-if="categories.length === 0">No categories available.</li>
          <!-- Render Categories -->
          <li
            v-else
            v-for="category in categories"
            :key="category.id"
            class="hover:text-blue-300 cursor-pointer"
            @click="fetchRecipesByCategory(category.name)"
          >
            {{ category.name }}
          </li>
        </ul>
      </div>
    </nav>

    <!-- Recipes Grid -->
    <div class="mt-24 container mx-auto p-4 grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      <!-- Loading Recipes -->
      <div v-if="loadingRecipes" class="col-span-full text-center text-gray-700">
        Loading recipes...
      </div>
      <!-- No Recipes Found -->
      <div v-else-if="recipes.length === 0" class="col-span-full text-center text-gray-700">
        No recipes found for the selected category.
      </div>
      <!-- Render Recipes -->
      <div
        v-else
        v-for="recipe in recipes"
        :key="recipe.id"
        class="bg-white shadow-md rounded-md p-4"
      >
        <img
          :src="recipe.featured_image || '/placeholder.png'"
          alt="Recipe Image"
          class="w-full h-48 object-cover rounded-md"
        />
        <h3 class="mt-4 text-lg font-semibold text-gray-800">
          {{ recipe.title }}
        </h3>
        <p class="mt-2 text-gray-600">
          {{ recipe.description }}
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useQuery, useLazyQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { computed, watch } from 'vue';

// GraphQL Queries
const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`;

const GET_RECIPES_BY_CATEGORY = gql`
  query GetRecipesByCategory($categoryName: String!) {
    recipes(where: { category: { name: { _eq: $categoryName } } }) {
      id
      title
      description
      featured_image
    }
  }
`;

// Fetch Categories
const { result: categoriesResult, loading: loadingCategories, error: categoriesError } = useQuery(GET_CATEGORIES);
const categories = computed(() => categoriesResult.value?.categories || []);

// Fetch Recipes
const { execute: fetchRecipes, result: recipesResult, loading: loadingRecipes, error: recipesError } = useLazyQuery(GET_RECIPES_BY_CATEGORY);
const recipes = computed(() => recipesResult.value?.recipes || []);

// Watch for Debugging
watch(categoriesResult, (newValue) => console.log("Categories fetched:", newValue));
watch(recipesResult, (newValue) => console.log("Recipes fetched:", newValue));

// Handle Fetch Recipes
const fetchRecipesByCategory = (categoryName) => {
  fetchRecipes({ categoryName })
    .catch((err) => console.error("Error fetching recipes by category:", err.message));
};

// Log Errors
if (categoriesError) console.error("Error fetching categories:", categoriesError.message);
if (recipesError) console.error("Error fetching recipes:", recipesError.message);
</script>

<style scoped>
/* Add any additional custom styles if needed */
</style>
