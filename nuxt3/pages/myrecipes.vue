<template>
  <div>
    <h1 class="text-2xl font-bold mb-4">My Recipes</h1>
    <div v-if="loading" class="text-gray-500">Loading...</div>
    <div v-else-if="error" class="text-red-500">Error: {{ error }}</div>
    <div v-else>
      <div
        v-for="recipe in recipes"
        :key="recipe.id"
        class="border rounded-lg p-4 mb-4 shadow-sm relative"
      >
        <h2 class="text-xl font-semibold">{{ recipe.title }}</h2>
        <p class="text-gray-600">{{ recipe.description }}</p>
        <p class="text-sm text-gray-500">
          Preparation Time: {{ recipe.preparation_time }} mins
        </p>
        <p class="text-sm text-gray-500">
          Category: {{ recipe.category.name }}
        </p>
        <img
          v-if="recipe.featured_image"
          :src="getImageUrl(recipe.featured_image)"
          alt="Recipe Image"
          class="w-[1000px] h-[500px] object-cover rounded-lg mt-2"
        />
        <div class="mt-4">
          <h3 class="font-semibold">Ingredients:</h3>
          <ul class="list-disc list-inside">
            <li v-for="ingredient in recipe.ingredients" :key="ingredient.id">
              {{ ingredient.quantity }} {{ ingredient.name }}
            </li>
          </ul>
        </div>
        <div class="mt-4">
          <h3 class="font-semibold">Steps:</h3>
          <ol class="list-decimal list-inside">
            <li v-for="step in recipe.steps" :key="step.id">
              {{ step.step_number }}. {{ step.description }}
            </li>
          </ol>
        </div>

         <div class="absolute top-4 right-4 flex gap-2">
          <button
            @click="navigateTo('/editrecipes')"
            class="px-4 py-2 bg-blue-500 text-white rounded-md"
          >
            Edit
          </button>
          <button
            @click="deleteRecipe(recipe.id)"
            class="px-4 py-2 bg-red-500 text-white rounded-md"
          >
            Delete
          </button>
        </div>
        <!--
      </div>
    </div>

    <div
      v-if="editingRecipe"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center"
    >
      <div class="bg-white p-6 rounded-lg">
        <h2 class="text-xl font-bold mb-4">Edit Recipe</h2>
        <form @submit.prevent="updateRecipe">
          <div class="mb-4">
            <label class="block font-semibold mb-2">Title</label>
            <input
              v-model="editingRecipe.title"
              class="border rounded-md p-2 w-full"
              type="text"
            />
          </div>
          <div class="mb-4">
            <label class="block font-semibold mb-2">Description</label>
            <textarea
              v-model="editingRecipe.description"
              class="border rounded-md p-2 w-full"
            ></textarea>
          </div>
          <div class="mb-4">
            <label class="block font-semibold mb-2">Preparation Time</label>
            <input
              v-model="editingRecipe.preparation_time"
              class="border rounded-md p-2 w-full"
              type="number"
            />
          </div>

          <div class="mb-4">
            <h3 class="font-semibold">Ingredients</h3>
            <div
              v-for="(ingredient, index) in editingRecipe.ingredients"
              :key="ingredient.id || index"
              class="flex items-center gap-2"
            >
              <input
                v-model="ingredient.quantity"
                class="border rounded-md p-2 w-1/4"
                placeholder="Quantity"
                type="text"
              />
              <input
                v-model="ingredient.name"
                class="border rounded-md p-2 w-3/4"
                placeholder="Name"
                type="text"
              />
              <button
                type="button"
                @click="removeIngredient(index)"
                class="px-2 py-1 bg-red-500 text-white rounded-md"
              >
                Remove
              </button>
            </div>
            <button
              type="button"
              @click="addIngredient"
              class="mt-2 px-4 py-2 bg-green-500 text-white rounded-md"
            >
              Add Ingredient
            </button>
          </div>

          <div class="mb-4">
            <h3 class="font-semibold">Steps</h3>
            <div
              v-for="(step, index) in editingRecipe.steps"
              :key="step.id || index"
              class="flex items-center gap-2"
            >
              <input
                v-model="step.step_number"
                class="border rounded-md p-2 w-1/4"
                placeholder="Step Number"
                type="number"
              />
              <textarea
                v-model="step.description"
                class="border rounded-md p-2 w-3/4"
                placehold er="Step Description"
              ></textarea>
              <button
                type="button"
                @click="removeStep(index)"
                class="px-2 py-1 bg-red-500 text-white rounded-md"
              >
                Remove
              </button>
            </div>
            <button
              type="button"
              @click="addStep"
              class="mt-2 px-4 py-2 bg-green-500 text-white rounded-md"
            >
              Add Step
            </button>
          </div>

          <div class="flex gap-2">
            <button
              type="submit"
              class="px-4 py-2 bg-green-500 text-white rounded-md"
            >
              Save
            </button>
            <button
              @click="cancelEditing"
              type="button"
              class="px-4 py-2 bg-gray-500 text-white rounded-md"
            >
              Cancel
            </button>
          </div> 
        </form>
      -->
      </div>
    </div>
  </div> 
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';

const recipes = ref([]);
const editingRecipe = ref(null);
const loading = ref(true);
const error = ref(null);
let userId = ''; // Initially empty

const { $apolloClient } = useNuxtApp();
const backendBaseUrl = 'http://localhost:8085';

// Fetch the logged-in user's ID from the token
const fetchUserIdFromToken = () => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("token");
    if (token) {
      try {
        const decodedToken = JSON.parse(atob(token.split(".")[1]));
        userId = decodedToken.userId;
      } catch (error) {
        console.error("Failed to decode token", error);
        error.value = "Invalid token. Please log in again.";
      }
    }
  }
};

const FETCH_RECIPES_QUERY = gql`
  query FetchRecipes($userId: uuid!) {
    recipes(where: { user_id: { _eq: $userId } }) {
      id
      title
      description
      preparation_time
      featured_image
      category {
        name
      }
      ingredients {
        id
        name
        quantity
      }
      steps {
        id
        step_number
        description
      }
    }
  }
`;

const EDIT_RECIPE_MUTATION = gql`
  mutation EditRecipe(
    $id: uuid!
    $input: recipes_set_input!
    $ingredients: [ingredients_insert_input!]!
    $steps: [steps_insert_input!]!
  ) {
    update_recipes_by_pk(pk_columns: { id: $id }, _set: $input) {
      id
      title
      description
      preparation_time
      ingredients {
        id
        name
        quantity
      }
      steps {
        id
        step_number
        description
      }
    }
    delete_ingredients(where: { recipe_id: { _eq: $id } })
    insert_ingredients(objects: $ingredients) {
      returning {
        id
        name
        quantity
      }
    }
    delete_steps(where: { recipe_id: { _eq: $id } })
    insert_steps(objects: $steps) {
      returning {
        id
        step_number
        description
      }
    }
  }
`;

const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : 'default-image.jpg';
};

const fetchRecipes = async () => {
  try {
    await fetchUserIdFromToken(); // Ensure user ID is fetched
    const response = await $apolloClient.query({
      query: FETCH_RECIPES_QUERY,
      variables: { userId },
    });
    recipes.value = response.data.recipes;
  } catch (err) {
    error.value = err.message;
  } finally {
    loading.value = false;
  }
};

const startEditing = (recipe) => {
  editingRecipe.value = structuredClone(recipe);
};

const cancelEditing = () => {
  editingRecipe.value = null;
};

// const updateRecipe = async () => {
//   if (!editingRecipe.value) return;

//   try {
//     const variables = {
//       id: editingRecipe.value.id,
//       input: {
//         title: editingRecipe.value.title,
//         description: editingRecipe.value.description,
//         preparation_time: editingRecipe.value.preparation_time,
//       },
//       ingredients: editingRecipe.value.ingredients.map((ingredient) => ({
//         id: ingredient.id || null,
//         recipe_id: editingRecipe.value.id,
//         name: ingredient.name,
//         quantity: ingredient.quantity,
//       })),
//       steps: editingRecipe.value.steps.map((step) => ({
//         id: step.id || null,
//         recipe_id: editingRecipe.value.id,
//         step_number: parseInt(step.step_number, 10),
//         description: step.description,
//       })),
//     };

//     await $apolloClient.mutate({
//       mutation: EDIT_RECIPE_MUTATION,
//       variables,
//     });

//     await fetchRecipes(); // Refetch to sync
//     cancelEditing();
//   } catch (err) {
//     console.error('Update Recipe Error:', err);
//     alert('Failed to update the recipe. Check console for details.');
//   }
// };

// const addIngredient = () => {
//   editingRecipe.value.ingredients.push({ id: null, name: '', quantity: '' });
// };

// const removeIngredient = (index) => {
//   editingRecipe.value.ingredients.splice(index, 1);
// };

// const addStep = () => {
//   editingRecipe.value.steps.push({ id: null, step_number: '', description: '' });
// };

// const removeStep = (index) => {
//   editingRecipe.value.steps.splice(index, 1);
// };

onMounted(fetchRecipes);
</script>

<style scoped>
/* Add any additional styling */
</style>
