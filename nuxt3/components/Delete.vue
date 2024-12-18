<template>
    <div class="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h1 class="text-2xl font-semibold mb-6">Manage Recipe</h1>
      <div v-if="recipe">
        <h2 class="text-xl font-semibold">{{ recipe.title }}</h2>
        <p>{{ recipe.description }}</p>
        <p><strong>Preparation Time:</strong> {{ recipe.prep_time }} mins</p>
        <p><strong>Cooking Time:</strong> {{ recipe.cook_time }} mins</p>
        <p><strong>Servings:</strong> {{ recipe.servings }}</p>
        <p><strong>Category:</strong> {{ recipe.category_id }}</p>
        <img v-if="recipe.image_url" :src="recipe.image_url" alt="Recipe Image" class="mt-4 w-full max-w-xs">
  
        <!-- Delete Button -->
        <div class="mt-6">
          <button
            @click="deleteRecipe"
            class="py-2 px-4 bg-red-600 text-white font-semibold rounded-md shadow-md hover:bg-red-700 focus:ring-2 focus:ring-red-500"
          >
            Delete Recipe
          </button>
        </div>
  
        <!-- Success or Error Messages -->
        <div v-if="success" class="mt-4 text-green-500">Recipe deleted successfully!</div>
        <div v-if="error" class="mt-4 text-red-500">Error: {{ error.message }}</div>
      </div>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue';
  import { gql } from '@apollo/client/core';
  import { useNuxtApp } from '#app';
  import { useRoute } from 'vue-router';
  
  // Get recipe ID from route
  const route = useRoute();
  const recipeId = route.params.id; // Assuming the recipe ID is passed as a route parameter
  
  const recipe = ref(null);
  const error = ref(null);
  const success = ref(false);
  
  const { $apolloClient } = useNuxtApp();
  
  // GraphQL Mutation for deleting a recipe
  const DELETE_RECIPE_MUTATION = gql`
    mutation DeleteRecipe($id: uuid!) {
      delete_recipes(where: { id: { _eq: $id } }) {
        affected_rows
      }
    }
  `;
  
  // Fetch recipe data (You can also add the recipe query here for the full details)
  const fetchRecipe = async () => {
    try {
      const response = await $apolloClient.query({
        query: gql`
          query GetRecipe($id: uuid!) {
            recipes_by_pk(id: $id) {
              id
              title
              description
              prep_time
              cook_time
              servings
              image_url
              category_id
            }
          }
        `,
        variables: { id: recipeId }
      });
      recipe.value = response.data.recipes_by_pk;
    } catch (err) {
      error.value = err;
      console.error('Error fetching recipe:', err);
    }
  };
  
  // Call to delete the recipe
  const deleteRecipe = async () => {
    try {
      const response = await $apolloClient.mutate({
        mutation: DELETE_RECIPE_MUTATION,
        variables: { id: recipeId }
      });
  
      if (response.data.delete_recipes.affected_rows > 0) {
        success.value = true;
        // Optionally, you can redirect the user after deletion
        // You can use `useRouter` from 'vue-router' to navigate to another page, like a recipe list page
      } else {
        error.value = { message: 'Failed to delete recipe.' };
      }
    } catch (err) {
      error.value = err;
      console.error('Error deleting recipe:', err);
    }
  };
  
  // Fetch recipe data when the component is mounted
  fetchRecipe();
  </script>
  
  <style scoped>
  /* Custom styles for delete button and error/success messages */
  </style>
  