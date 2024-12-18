<template>
  <div class="container mx-auto py-8 px-4">
    <h1 class="text-2xl font-bold mb-6">Add Recipe with Details</h1>

    <form @submit.prevent="submitRecipe" class="space-y-4">
      <div>
        <label for="title" class="block font-medium mb-1">Title</label>
        <input
          v-model="form.title"
          id="title"
          type="text"
          class="w-full px-3 py-2 border rounded"
          placeholder="Enter recipe title"
          required
        />
      </div>

      <div>
        <label for="description" class="block font-medium mb-1">Description</label>
        <textarea
          v-model="form.description"
          id="description"
          class="w-full px-3 py-2 border rounded"
          placeholder="Enter recipe description"
          required
        ></textarea>
      </div>

      <div>
        <label for="category" class="block font-medium mb-1">Category</label>
        <select
          v-model="selectedCategory"
          id="category"
          class="w-full px-3 py-2 border rounded"
          @change="handleCategoryChange"
          required
        >
          <option value="" disabled>Select a category</option>
          <option v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
          <option value="new">+ Add New Category</option>
        </select>
        <div v-if="isAddingCategory" class="mt-2">
          <input
            v-model="newCategory"
            type="text"
            class="w-full px-3 py-2 border rounded"
            placeholder="Enter new category name"
            required
          />
          <button @click.prevent="addNewCategory" class="mt-2 bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600">
            Add Category
          </button>
        </div>
      </div>

      <div>
        <label for="time" class="block font-medium mb-1">Preparation Time (minutes)</label>
        <input
          v-model="form.preparation_time"
          id="time"
          type="number"
          class="w-full px-3 py-2 border rounded"
          placeholder="Enter preparation time"
          required
        />
      </div>

      <!-- Dynamic Inputs for Ingredients -->
      <div>
        <label class="block font-medium mb-1">Ingredients</label>
        <div v-for="(ingredient, index) in form.ingredients" :key="index" class="flex gap-4 mb-2">
          <input
            v-model="ingredient.name"
            type="text"
            class="w-1/2 px-3 py-2 border rounded"
            placeholder="Ingredient name"
            required
          />
          <input
            v-model="ingredient.quantity"
            type="text"
            class="w-1/2 px-3 py-2 border rounded"
            placeholder="Quantity (e.g., 1 tsp)"
            required
          />
          <button @click.prevent="removeIngredient(index)" class="text-red-500">Remove</button>
        </div>
        <button @click.prevent="addIngredient" class="text-blue-500">+ Add Ingredient</button>
      </div>

      <!-- Dynamic Inputs for Steps -->
      <div>
        <label class="block font-medium mb-1">Steps</label>
        <div v-for="(step, index) in form.steps" :key="index" class="flex gap-4 mb-2">
          <input
            v-model="step.step_number"
            type="number"
            class="w-1/6 px-3 py-2 border rounded"
            placeholder="Step #"
            required
          />
          <textarea
            v-model="step.description"
            class="w-5/6 px-3 py-2 border rounded"
            placeholder="Step description"
            required
          ></textarea>
          <button @click.prevent="removeStep(index)" class="text-red-500">Remove</button>
        </div>
        <button @click.prevent="addStep" class="text-blue-500">+ Add Step</button>
      </div>

      <button
        type="submit"
        class="bg-blue-500 text-white px-4 py-2 rounded hover:bg-blue-600"
      >
        Submit Recipe
      </button>

      <p v-if="successMessage" class="text-green-500 mt-4">{{ successMessage }}</p>
      <p v-if="errorMessage" class="text-red-500 mt-4">{{ errorMessage }}</p>
    </form>
  </div>
</template>
<script setup>
import { ref, onMounted } from "vue";
import { useApolloClient } from "@vue/apollo-composable";
import gql from "graphql-tag";

const form = ref({
  title: "",
  description: "",
  category_id: "",
  user_id: "",
  preparation_time: "",
  ingredients: [{ name: "", quantity: "" }],
  steps: [{ step_number: 1, description: "" }],
});

const categories = ref([]);
const selectedCategory = ref("");
const newCategory = ref("");
const isAddingCategory = ref(false);

const successMessage = ref("");
const errorMessage = ref("");

const client = useApolloClient().client;

// Automatically set user_id from token
const token = localStorage.getItem("token");
if (token) {
  try {
    const decodedToken = JSON.parse(atob(token.split(".")[1]));
    form.value.user_id = decodedToken.userId;
  } catch (error) {
    console.error("Failed to decode token", error);
    errorMessage.value = "Invalid token. Please log in again.";
  }
}

// Fetch categories on component mount
const fetchCategories = async () => {
  try {
    const query = gql`
      query GetCategories {
        categories {
          id
          name
        }
      }
    `;

    const { data } = await client.query({ query, fetchPolicy: "no-cache" });
    categories.value = data.categories;
  } catch (error) {
    console.error("Failed to fetch categories:", error);
    errorMessage.value = "Failed to load categories. Please try again.";
  }
};

const handleCategoryChange = () => {
  if (selectedCategory.value === "new") {
    isAddingCategory.value = true;
    form.value.category_id = ""; // Clear category_id until the new category is added
  } else {
    isAddingCategory.value = false;
    form.value.category_id = selectedCategory.value;
  }
};

const addNewCategory = async () => {
  try {
    const mutation = gql`
      mutation AddCategory($name: String!) {
        insert_categories_one(object: { name: $name }) {
          id
          name
        }
      }
    `;

    const { data } = await client.mutate({
      mutation,
      variables: { name: newCategory.value },
      fetchPolicy: "no-cache",
    });

    categories.value.push(data.insert_categories_one);
    selectedCategory.value = data.insert_categories_one.id;
    form.value.category_id = data.insert_categories_one.id;

    isAddingCategory.value = false;
    newCategory.value = "";
    successMessage.value = "Category added successfully!";
  } catch (error) {
    console.error("Failed to add new category:", error);
    errorMessage.value = "Failed to add new category. Please try again.";
  }
};

const submitRecipe = async () => {
  if (!form.value.category_id) {
    errorMessage.value = "Please select or add a category.";
    return;
  }

  try {
    const mutation = gql`
      mutation AddRecipe($input: recipes_insert_input!) {
        insert_recipes_one(object: $input) {
          id
          title
        }
      }
    `;

    const variables = {
      input: {
        title: form.value.title,
        description: form.value.description,
        category_id: form.value.category_id,
        user_id: form.value.user_id,
        preparation_time: parseInt(form.value.preparation_time, 10),
        ingredients: {
          data: form.value.ingredients.map((ingredient) => ({
            name: ingredient.name,
            quantity: ingredient.quantity,
          })),
        },
        steps: {
          data: form.value.steps.map((step) => ({
            step_number: parseInt(step.step_number, 10),
            description: step.description,
          })),
        },
      },
    };

    const { data } = await client.mutate({ mutation, variables });
    successMessage.value = `Recipe "${data.insert_recipes_one.title}" added successfully!`;
    resetForm();
  } catch (error) {
    console.error("Failed to submit recipe:", error);
    errorMessage.value = "Failed to submit recipe. Please try again.";
  }
};

const resetForm = () => {
  form.value = {
    title: "",
    description: "",
    category_id: "",
    user_id: form.value.user_id, // Preserve the logged-in user's ID
    preparation_time: "",
    ingredients: [{ name: "", quantity: "" }],
    steps: [{ step_number: 1, description: "" }],
  };
  selectedCategory.value = "";
};

onMounted(fetchCategories);
</script>
