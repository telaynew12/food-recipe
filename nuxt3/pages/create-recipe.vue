<template>
  <div class="container mx-auto py-8 px-4 center">
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

     <!-- // Dynamic Inputs for Images -->
      <!-- <div>
        <label class="block font-medium mb-1">Images</label>
        <div v-for="(image, index) in form.images" :key="index" class="flex gap-4 mb-2">
          <input
            v-model="image.image_url"
            type="text"
            class="w-3/4 px-3 py-2 border rounded"
            placeholder="Image URL"
            required
          />
          <input
            v-model="image.is_featured"
            type="checkbox"
            class="w-1/4"
          />
          <label>Featured</label>
          <button @click.prevent="removeImage(index)" class="text-red-500">Remove</button>
        </div>
        <button @click.prevent="addImage" class="text-blue-500">+ Add Image</button>
      </div> -->

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
  <div class="flex flex-col items-center justify-center space-y-4 p-6 bg-gray-100 rounded-lg shadow-md">
  <h1 class="text-2xl font-bold text-gray-700">File Upload</h1>
  
  <!-- File Input -->
  <label 
    for="fileInput" 
    class="w-20 h-20 flex items-center justify-center bg-blue-500 text-white font-medium  cursor-pointer hover:transition">
    Choose File
    <input 
      id="fileInput" 
      type="file" 
      class="hidden" 
      @change="handleFileChange" />
  </label>
  
  <!-- Upload Button -->
  <button 
    @click="uploadFile" 
    :disabled="!selectedFile" 
    class="px-6 py-2 bg-green-500 text-white font-medium rounded-md hover:bg-green-600 transition disabled:bg-gray-400 disabled:cursor-not-allowed">
    Upload
  </button>

  <!-- Status Message -->
  <p v-if="uploadStatus" class="text-sm text-gray-600">{{ uploadStatus }}</p>
</div>

</template>

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
  ingredients: [{ name: "", quantity: "" }],
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

const client = useApolloClient().client;

const fetchUserIdFromToken = () => {
  if (typeof window !== "undefined") {
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
  form.value.ingredients.push({ name: "", quantity: "" });
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

    const variables = {
      input: {
        title: form.value.title,
        description: form.value.description,
        category_id: form.value.category_id,
        user_id: form.value.user_id,
        preparation_time: parseInt(form.value.preparation_time, 10),
        ingredients: {
          data: form.value.ingredients.map(({ name, quantity }) => ({ name, quantity })),
        },
        steps: {
          data: form.value.steps.map(({ step_number, description }) => ({
            step_number,
            description,
          })),
        },
        shares: {
          data: form.value.shares.map(({ is_shared }) => ({ is_shared })),
        },
      },
    };

    const { data } = await client.mutate({
      mutation,
      variables,
    });
    recipeId = data.insert_recipes_one.id
     successMessage.value = `Recipe "${data.insert_recipes_one.title}" added successfully!`;
    form.value = {
      title: "",
      description: "",
      category_id: "",
      user_id: "",
      preparation_time: "",
      ingredients: [{ name: "", quantity: "" }],
      steps: [{ step_number: 1, description: "" }],
      images: [{ image_url: "", is_featured: false }],
    
    };
  } catch (error) {
    console.error("Failed to submit recipe:", error);
    errorMessage.value = "Failed to submit recipe. Please try again.";
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
  fetchCategories();
});
</script>

