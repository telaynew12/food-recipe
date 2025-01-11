<template>
  <div class="container mx-auto py-8 px-4 max-w-4xl">
    <h1 class="text-3xl font-bold text-center text-gray-800 mb-8">
      Add Recipe with Details
    </h1>

    <form @submit.prevent="submitRecipe" class="bg-white p-6 rounded-lg shadow-md space-y-6">
      <div>
        <label for="title" class="block font-semibold text-gray-700 mb-2">Title</label>
        <input
          v-model="form.title"
          id="title"
          type="text"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
          placeholder="Enter recipe title"
          required
        />
      </div>

      <div>
        <label for="description" class="block font-semibold text-gray-700 mb-2">Description</label>
        <textarea
          v-model="form.description"
          id="description"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
          placeholder="Enter recipe description"
          rows="4"
          required
        ></textarea>
      </div>

      <div>
        <label for="category" class="block font-semibold text-gray-700 mb-2">Category</label>
        <select
          v-model="selectedCategory"
          id="category"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
          @change="handleCategoryChange"
          required
        >
          <option value="" disabled>Select a category</option>
          <option v-for="category in categories" :key="category.id" :value="category.id">
            {{ category.name }}
          </option>
          <option value="new">+ Add New Category</option>
        </select>

        <div v-if="isAddingCategory" class="mt-4">
          <input
            v-model="newCategory"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-green-300 focus:outline-none"
            placeholder="Enter new category name"
            required
          />
          <button
            @click.prevent="addNewCategory"
            class="mt-3 bg-green-500 text-white px-6 py-2 rounded-lg hover:bg-green-600 transition"
          >
            Add Category
          </button>
        </div>
      </div>

      <div>
        <label for="time" class="block font-semibold text-gray-700 mb-2">Preparation Time (minutes)</label>
        <input
          v-model="form.preparation_time"
          id="time"
          type="number"
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
          placeholder="Enter preparation time"
          required
        />
      </div>

      <div>
        <label class="block font-semibold text-gray-700 mb-2">Ingredients</label>
        <div v-for="(ingredient, index) in form.ingredients" :key="index" class="flex items-center gap-4 mb-3">
          <input
            v-model="ingredient.name"
            type="text"
            class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
            placeholder="Ingredient name"
            required
          />
          <button
            @click.prevent="removeIngredient(index)"
            class="text-red-500 hover:text-red-600 transition"
          >
            Remove
          </button>
        </div>
        <button @click.prevent="addIngredient" class="text-blue-500 hover:text-blue-600 transition">
          + Add Ingredient
        </button>
      </div>

      <div>
        <label class="block font-semibold text-gray-700 mb-2">Steps</label>
        <div v-for="(step, index) in form.steps" :key="index" class="flex items-center gap-4 mb-3">
          <input
            v-model="step.step_number"
            type="number"
            class="w-1/6 px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
            placeholder="Step #"
            required
          />
          <textarea
            v-model="step.description"
            class="w-5/6 px-4 py-2 border border-gray-300 rounded-lg focus:ring focus:ring-blue-300 focus:outline-none"
            placeholder="Step description"
            rows="2"
            required
          ></textarea>
          <button
            @click.prevent="removeStep(index)"
            class="text-red-500 hover:text-red-600 transition"
          >
            Remove
          </button>
        </div>
        <button @click.prevent="addStep" class="text-blue-500 hover:text-blue-600 transition">
          + Add Step
        </button>
      </div>

      <button
        type="submit"
        class="w-full bg-blue-500 text-white py-3 rounded-lg shadow-lg hover:bg-blue-600 transition"
      >
        Submit Recipe
      </button>

      <p v-if="successMessage" class="text-green-500 mt-6">{{ successMessage }}</p>
      <p v-if="errorMessage" class="text-red-500 mt-6">{{ errorMessage }}</p>
    </form>

    <div class="mt-12 p-6 bg-gray-50 rounded-lg shadow-md">
      <h2 class="text-xl font-bold text-gray-700 text-center mb-6">File Upload</h2>

      <label
        for="fileInput"
        class="block w-full p-4 text-center bg-blue-500 text-white rounded-lg cursor-pointer hover:bg-blue-600 transition"
      >
        Choose File  
        <input id="fileInput" type="file" class="hidden" @change="handleFileChange" />
      </label>

      <button
        @click="uploadFile"
        :disabled="!selectedFile"
        class="w-full mt-4 bg-green-500 text-white py-3 rounded-lg hover:bg-green-600 transition disabled:bg-gray-300 disabled:cursor-not-allowed"
      >
        Upload
      </button>

      <p v-if="uploadStatus" class="text-gray-600 mt-4">{{ uploadStatus }}</p>
    </div>
  </div>
</template>
 n                         pkkiojjj

<script setup>
import { ref, onMounted } from "vue";
import { useApolloClient } from "@vue/apollo-composable";
import gql from "graphql-tag";

let recipeId = null;

const form = ref({
  title: "",
  description: "",
  category_id: "",
  user_id: "",
  preparation_time: "",
  ingredients: [{ name: "" }],
  steps: [{ step_number: 1, description: "" }],
  images: [{ image_url: "", is_featured: false }],
  shares: [{ is_shared: false }], // Adjusted for compatibility
});

const categories = ref([]);
const selectedCategory = ref("");
const newCategory = ref("");
const isAddingCategory = ref(false);

const successMessage = ref("");
const errorMessage = ref("");
const uploadStatus = ref("");
const selectedFile = ref(null);
const isAuthenticated = ref(false);

const client = useApolloClient().client;

const fetchUserIdFromToken = () => {
  if (typeof window !== "undefined") {
    const token = localStorage.getItem("token");
    if (token) {
      try {
        const decodedToken = JSON.parse(atob(token.split(".")[1]));
        form.value.user_id = decodedToken.userId;
        isAuthenticated.value = true; // User is authenticated
      } catch (error) {
        console.error("Failed to decode token", error);
        errorMessage.value = "Invalid token. Please log in again.";
        isAuthenticated.value = false;
      }
    } else {
      isAuthenticated.value = false;
    }
  }
};

const ensureAuthentication = () => {
  if (!isAuthenticated.value) {
    router.push("/login");
  }
};

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
  isAddingCategory.value = selectedCategory.value === "new";
  form.value.category_id = isAddingCategory.value ? "" : selectedCategory.value;
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

const addIngredient = () => {
  if (form.value.ingredients.some((ingredient) => ingredient.name.trim() === "")) {
    errorMessage.value = "Please fill the existing ingredient before adding a new one.";
    return;
  }
  form.value.ingredients.push({ name: "" });
};


const removeIngredient = (index) => {
  form.value.ingredients.splice(index, 1);
};

const addStep = () => {
  const nextStepNumber = form.value.steps.length + 1;
  form.value.steps.push({ step_number: nextStepNumber, description: "" });
};

const removeStep = (index) => {
  form.value.steps.splice(index, 1);
  form.value.steps.forEach((step, i) => (step.step_number = i + 1));
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

    // Ensure the structure of ingredients and steps is valid
    const ingredients = form.value.ingredients
      ?.filter((ingredient) => ingredient.name?.trim()) // Ensure valid names
      .map((ingredient) => ({ name: ingredient.name })) || [];

    const steps = form.value.steps
      ?.filter((step) => step.description?.trim()) // Ensure valid steps
      .map(({ step_number, description }) => ({
        step_number,
        description,
      })) || [];

    const shares = form.value.shares
      ?.map(({ is_shared }) => ({ is_shared })) || [];

    // Build the variables object
    const variables = {
      input: {
        title: form.value.title,
        description: form.value.description,
        category_id: form.value.category_id,
        user_id: form.value.user_id || null, // Ensure this is nullable if not mandatory
        preparation_time: parseInt(form.value.preparation_time, 10) || null,
        ingredients: { data: ingredients },
        steps: { data: steps },
        shares: { data: shares },
      },
    };

    const { data } = await client.mutate({
      mutation,
      variables,
    });

    recipeId = data.insert_recipes_one.id;
    successMessage.value = `Recipe "${data.insert_recipes_one.title}" added successfully!`;

    // Reset the form
    form.value = {
      title: "",
      description: "",
      category_id: "",
      user_id: form.value.user_id,
      preparation_time: "",
      ingredients: [{ name: "" }],
      steps: [{ step_number: 1, description: "" }],
      images: [{ image_url: "", is_featured: false }],
      shares: [{ is_shared: false }],
    };
  } catch (error) {
    console.error("Failed to submit recipe:", error);
    errorMessage.value = `Failed to submit recipe: ${error.message}`;
  }
};


// GraphQL mutation for file upload
const FILE_UPLOAD = gql`
  mutation fileUpload($name: String!, $type: String!, $base64str: String!,$RecipeID:String!) {
    fileUpload(name: $name, type: $type, base64str: $base64str,RecipeID:$RecipeID) {
      file_path
      RecipeID
    }
  }
`

const { mutate } = useMutation(FILE_UPLOAD)

const handleFileChange = async (event) => {
  const file = event.target.files[0]
  if (file) {
    const reader = new FileReader()
    reader.onload = (e) => {
      selectedFile.value = {
        name: file.name,
        type: file.type,
        base64str: e.target.result.split(',')[1],
      }
    }
    reader.readAsDataURL(file)
  }
}

const uploadFile = async () => {
  if (!selectedFile.value) {
    uploadStatus.value = 'No file selected!'
    return
  }

  try {
    uploadStatus.value = 'Uploading...'
    const { data } = await mutate({
      name: selectedFile.value.name,
      type: selectedFile.value.type,
      base64str: selectedFile.value.base64str,
      RecipeID:recipeId,
    })
    uploadStatus.value = `File uploaded successfully: ${data.fileUpload.file_path}`
  } catch (error) {
    uploadStatus.value = `Error: ${error.message}`
  }
}




onMounted(() => {
  fetchUserIdFromToken();
  ensureAuthentication(); // Redirects to login if not authenticated
  fetchCategories();
});
</script>

