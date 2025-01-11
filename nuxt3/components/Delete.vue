<template>
  <div v-if="userId">
    <h1 class="font-bold"> Liked Rcipes</h1>
    <h1>Recipes for User ID: {{ userId }}</h1>
    <div v-if="recipes && recipes.length > 0">
      <ul>
        <li v-for="recipe in recipes" :key="recipe.id">
          <h2> <strong>Recipe Title:</strong>{{ recipe.title }}</h2>
      
          <div v-if="recipe.likes.length > 0">
            <h3 class="font-bold">Liked by:</h3>
            <ul>
  <li v-for="(like, index) in recipe.likes" :key="like.id">
                 {{ index + 1 }}.{{ like.user.name }}
              </li>
            </ul>
          </div>
          <div v-else>
            <p>No likes for this recipe yet.</p>
          </div>
        </li>
      </ul>
    </div>
    <div v-else>
      <p>No recipes found for this user. Please check again later.</p>
    </div>
  </div>
  <div v-else>
    <p>User ID not available. Please log in to view your recipes.</p>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';
// Access the authentication store
const authStore = useAuthStore();
const userId = computed(() => authStore.userId);

// Reactive variable to store recipes
const recipes = ref([]);

// GraphQL query for fetching recipes with likes and user names
const FETCH_RECIPES_QUERY = gql`
  query FetchRecipes($userId: uuid!) {
    recipes(where: { user_id: { _eq: $userId } }) {
      id
      title
      description
      featured_image
      likes {
        id
        user {
          id
          name
        }
      }
    }
  }
`;

// Fetch recipes function
const fetchRecipes = async () => {
  try {
    const { client } = useApolloClient();
    const { data } = await client.query({
      query: FETCH_RECIPES_QUERY,
      variables: { userId: userId.value },
    });

    if (data && data.recipes) {
      recipes.value = data.recipes;
    } else {
      console.log('No recipes found.');
    }
  } catch (error) {
    console.error('Error fetching recipes:', error.message);
  }
};

// Fetch recipes when the component is mounted
onMounted(() => {
  if (userId.value) {
    fetchRecipes();
  } else {
    console.error('User ID is not available.');
  }
});
</script>
