<template>
    <div>
      <h1>User Recipes</h1>
      <div v-if="loading">Loading...</div>
      <div v-else-if="error">Error: {{ error.message }}</div>
      <div v-else>
        <div v-for="recipe in recipes" :key="recipe.id" class="recipe">
          <h2><strong>Title:</strong> {{ recipe.title }}</h2>
          <img v-if="recipe.image_url" :src="recipe.image_url" alt="Recipe Image" class="recipe-image" />
          <p><strong>Description:</strong> {{ recipe.description }}</p>
          <p><strong>Category:</strong> {{ recipe.category?.name || 'Uncategorized' }}</p>
          <!-- <p><strong>Preparation Time:</strong> {{ recipe.prep_time ? recipe.prep_time + ' mins' : 'N/A' }}</p>
          <p><strong>Cooking Time:</strong> {{ recipe.cook_time ? recipe.cook_time + ' mins' : 'N/A' }}</p>
          <p><strong>Servings:</strong> {{ recipe.servings ? recipe.servings : 'N/A' }}</p>
          <p><small>Created at: {{ recipe.created_at }}</small></p>
          <p><small>Last updated at: {{ recipe.updated_at }}</small></p> -->
        </div>
      </div>
    </div>
  </template>
<script setup>
import { ref, onMounted } from 'vue';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';

// Reactive state
const recipes = ref([]);
const loading = ref(true);
const error = ref(null);
const userId = ref(''); // Holds the decoded userId from the JWT token

const { $apolloClient } = useNuxtApp();

// GraphQL Query for Fetching Recipes
const FETCH_RECIPES_QUERY = gql`
  query FetchRecipes($user_id: uuid!) {
    recipes(where: { user_id: { _eq: $user_id } }) {
      id
      title
      description
      prep_time
      cook_time
      servings
      image_url
      created_at
      updated_at
      category {
        name
      }
    }
  }
`;

// Decode userId and Fetch Recipes
onMounted(async () => {
  try {
    const token = localStorage.getItem('token'); // JWT token from localStorage

    if (token) {
      try {
        const decodedToken = JSON.parse(atob(token.split('.')[1])); // Decode JWT
        userId.value = decodedToken.userId; // Get userId from token payload
      } catch (decodeError) {
        console.error('Failed to decode token', decodeError);
        throw decodeError;
      }
    }

    if (userId.value) {
      // Fetch recipes for the user
      const response = await $apolloClient.query({
        query: FETCH_RECIPES_QUERY,
        variables: { user_id: userId.value },
      });
      recipes.value = response.data.recipes;
    } else {
      throw new Error('User ID not found.');
    }
  } catch (err) {
    error.value = err;
  } finally {
    loading.value = false;
  }
});
</script>
  