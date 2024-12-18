<template>
    <div class="max-w-2xl mx-auto p-6 bg-white rounded-lg shadow-md">
      <h1 class="text-2xl font-semibold mb-6">Edit Recipe</h1>
      <form @submit.prevent="submitRecipe">
        <div class="mb-4">
          <label for="title" class="block text-sm font-semibold mb-2">Title</label>
          <input
            id="title"
            v-model="title"
            type="text"
            class="w-full p-2 border border-gray-300 rounded-lg"
            required
          />
        </div>
        <div class="mb-4">
          <label for="description" class="block text-sm font-semibold mb-2">Description</label>
          <textarea
            id="description"
            v-model="description"
            class="w-full p-2 border border-gray-300 rounded-lg"
            required
          ></textarea>
        </div>
        <div class="mb-4">
          <label for="prep_time" class="block text-sm font-semibold mb-2">Preparation Time (mins)</label>
          <input
            id="prep_time"
            v-model.number="prepTime"
            type="number"
            class="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>
        <div class="mb-4">
          <label for="cook_time" class="block text-sm font-semibold mb-2">Cooking Time (mins)</label>
          <input
            id="cook_time"
            v-model.number="cookTime"
            type="number"
            class="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>
        <div class="mb-4">
          <label for="servings" class="block text-sm font-semibold mb-2">Servings</label>
          <input
            id="servings"
            v-model.number="servings"
            type="number"
            class="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>
        <div class="mb-4">
          <label for="image_url" class="block text-sm font-semibold mb-2">Image URL</label>
          <input
            id="image_url"
            v-model="imageUrl"
            type="text"
            class="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>
        <div class="mb-4">
          <label for="category_id" class="block text-sm font-semibold mb-2">Category ID</label>
          <input
            id="category_id"
            v-model="categoryId"
            type="text"
            class="w-full p-2 border border-gray-300 rounded-lg"
          />
        </div>
        <button
          type="submit"
          class="w-full py-2 bg-blue-600 text-white font-semibold rounded-lg hover:bg-blue-700 focus:ring-2 focus:ring-blue-500"
        >
          Update Recipe
        </button>
      </form>
  
      <!-- Success or Error Messages -->
      <div v-if="success" class="mt-4 text-green-500">Recipe updated successfully!</div>
      <div v-if="error" class="mt-4 text-red-500">Error: {{ error.message }}</div>
    </div>
  </template>
  
  <script setup>
  import { ref, onMounted } from 'vue';
  import { gql } from '@apollo/client/core';
  import { useNuxtApp } from '#app';
  import { useRoute, useRouter } from 'vue-router';
  
  // Get recipe ID from the route params
  const route = useRoute();
  const router = useRouter();
  const recipeId = route.params.id; // Recipe ID passed in route params
  
  // Reactive state for form fields
  const title = ref('');
  const description = ref('');
  const prepTime = ref(null);
  const cookTime = ref(null);
  const servings = ref(null);
  const imageUrl = ref('');
  const categoryId = ref('');
  const success = ref(false);
  const error = ref(null);
  
  // Apollo Client
  const { $apolloClient } = useNuxtApp();
  
  // GraphQL Query for Fetching Recipe Details
  const FETCH_RECIPE_QUERY = gql`
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
  `;
  
  // GraphQL Mutation for Updating a Recipe
  const UPDATE_RECIPE_MUTATION = gql`
    mutation UpdateRecipe(
      $id: uuid!
      $title: String!
      $description: String!
      $prep_time: Int
      $cook_time: Int
      $servings: Int
      $image_url: String
      $category_id: uuid
    ) {
      update_recipes(
        where: { id: { _eq: $id } }
        _set: {
          title: $title
          description: $description
          prep_time: $prep_time
          cook_time: $cook_time
          servings: $servings
          image_url: $image_url
          category_id: $category_id
        }
      ) {
        returning {
          id
          title
        }
      }
    }
  `;
  
  // Fetch the recipe data when the component is mounted
  onMounted(async () => {
    try {
      const response = await $apolloClient.query({
        query: FETCH_RECIPE_QUERY,
        variables: { id: recipeId },
      });
      const recipe = response.data.recipes_by_pk;
      title.value = recipe.title;
      description.value = recipe.description;
      prepTime.value = recipe.prep_time;
      cookTime.value = recipe.cook_time;
      servings.value = recipe.servings;
      imageUrl.value = recipe.image_url;
      categoryId.value = recipe.category_id;
    } catch (err) {
      error.value = err;
      console.error('Error fetching recipe:', err);
    }
  });
  
  // Submit the form and update the recipe
  const submitRecipe = async () => {
    try {
      const response = await $apolloClient.mutate({
        mutation: UPDATE_RECIPE_MUTATION,
        variables: {
          id: recipeId,
          title: title.value,
          description: description.value,
          prep_time: prepTime.value,
          cook_time: cookTime.value,
          servings: servings.value,
          image_url: imageUrl.value,
          category_id: categoryId.value,
        },
      });
  
      if (response.data.update_recipes.returning.length > 0) {
        success.value = true;
        setTimeout(() => {
          router.push('/recipes'); // Redirect to recipes list page after successful update
        }, 1500); // Wait for 1.5 seconds before redirect
      }
    } catch (err) {
      error.value = err;
      console.error('Error updating recipe:', err);
    }
  };
  </script>
  
  <style scoped>
  /* Custom styles for form and buttons */
  </style>
  