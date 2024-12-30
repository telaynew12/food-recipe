


<template>
  <div class="p-4">
    

    <!-- Recipe Details -->
    <div v-if="recipe" class="rounded-lg shadow-lg bg-white p-6">
      <img
        v-if="recipe.featured_image"
        :src="getImageUrl(recipe.featured_image)"
        alt="Recipe Image"
        class="w-300 h-90 object-cover rounded-t-lg mb-4"
      />
      <h1 class="text-2xl font-bold mb-2">{{ recipe.title }}</h1>
      <p class="text-gray-700 mb-2">{{ recipe.description }}</p>
      <p class="text-gray-600 text-sm">
        <strong>Preparation Time:</strong> {{ recipe.preparation_time }} mins
      </p>
      <p class="text-gray-600 text-sm">
        <strong>Created At:</strong> {{ new Date(recipe.created_at).toLocaleDateString() }}
      </p>

      <!-- Steps Section -->
      <h2 class="text-xl font-bold mt-6 mb-4">Steps</h2>
      <ul class="list-decimal list-inside">
        <li
          v-for="step in recipe.steps"
          :key="step.step_number"
          class="text-gray-700 mb-2"
        >
          <strong>Step {{ step.step_number }}:</strong> {{ step.description }}
        </li>
      </ul>

      <!-- Ingredients Section -->
      <h2 class="text-xl font-bold mt-6 mb-4">Ingredients</h2>
      <ul class="list-disc list-inside">
        <li
          v-for="ingredient in recipe.ingredients"
          :key="ingredient.name"
          class="text-gray-700 mb-2"
        >
          <strong>{{ ingredient.name }}:</strong> {{ ingredient.quantity }}
        </li>
      </ul>
    </div>

    <!-- Error Handling -->
    <div v-else-if="errorMessage" class="text-red-500">
      Failed to fetch recipe details: {{ errorMessage.message }}
    </div>

    <!-- Loading or Empty State -->
    <div v-else class="text-gray-500">
      <p v-if="!recipeId">Please select a recipe to view its details.</p>
      <p v-else>Loading recipe details...</p>
    </div>
    <!-- Print Button -->
    <div class="mb-4">
      <button
        @click="printPage"
        class="bg-blue-500 text-white px-4 py-2 rounded shadow hover:bg-blue-600"
      >
        Print Recipe
      </button>
    </div>
  </div>

</template>







<script setup>
import { computed, watchEffect } from 'vue';
import { useQuery, useResult } from '@vue/apollo-composable';
import { gql } from 'graphql-tag';
import { useDetailStore } from '@/stores/detail';

const backendBaseUrl = 'http://localhost:8085/';
const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : null;
};

const GET_RECIPE_QUERY = gql`
  query GetRecipe($id: uuid!) {
    recipes_by_pk(id: $id) {
      id
      title
      description
      preparation_time
      featured_image
      created_at
      steps {
        step_number
        description
      }
      ingredients {
        name
        quantity
      }
    }
  }
`;

const recipeDetailStore = useDetailStore();
const recipeId = computed(() => recipeDetailStore.recipeId);

// Only run the query if recipeId is a valid UUID
const { result, error } = useQuery(GET_RECIPE_QUERY, { id: recipeId.value || undefined }, { enabled: !!recipeId.value });

const recipe = useResult(result, null, (data) => data?.recipes_by_pk);
const errorMessage = error;

const printPage = () => {
  window.print();
};
</script>
