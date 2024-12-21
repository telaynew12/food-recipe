<template>
  <div class="p-6 max-w-4xl mx-auto bg-white shadow-md rounded-lg">
    <!-- Display All Recipes -->
    <div v-if="recipes.length" class="space-y-8">
      <div v-for="(recipe, index) in recipes" :key="recipe.id" class="border-b pb-6">
        <!-- Recipe Title -->
        <h2 class="text-2xl font-bold mb-4">{{ recipe.title }}</h2>

        <!-- Recipe Image -->
        <img 
          v-if="recipe.featured_image"
          :src="getImageUrl(recipe.featured_image)"
          alt="Recipe Image"
          class="w-full h-500 object-cover rounded-lg mb-4"
        />

        <!-- Buttons Section (flex container) -->
        <div class="flex justify-start gap-6 mb-4 items-center">
          <!-- Like Button -->
<!-- Like Button -->
<button
  @click="handleLike(index)"
  :class="{ 'text-blue-500 font-bold': recipe.liked, 'text-gray-600': !recipe.liked }"
  class="flex items-center space-x-2"
>
  <span>👍</span>
  <span>{{ recipe.liked ? 'liked' : 'Like' }}</span>
  <span>({{ recipe.likesCount }})</span>
</button>


<button 
  @click="toggleComments(index)"
  class="flex items-center space-x-2 text-gray-600 hover:text-gray-800"
>
  <span>💬</span>
  <span>Comments</span>
  <span>({{ recipe.commentsCount }})</span> <!-- Use commentsCount here -->
</button>


          <!-- Bookmark Button -->
          <button
  @click="handleBookmark(index)"
  :class="{ 'text-yellow-500': recipe.bookmarked, 'text-gray-600': !recipe.bookmarked }"
  class="flex items-center space-x-2"
>
  <span>🔖</span>
  <span>{{ recipe.bookmarked ? 'Bookmarked' : 'Bookmark' }}</span>
</button>


          <!-- Rating Section -->
          <div class="flex items-center space-x-1">
            <span class="font-semibold">Rate:</span>
  <div class="flex">
    <span
      v-for="star in 5"
      :key="star"
      :class="{
        'text-yellow-400': recipe.rating >= star,
        'text-gray-400': recipe.rating < star
      }"
      class="cursor-pointer text-2xl"
      @click="handleRating(index, star)"
    >
      ★
    </span>
            </div>
            <span class="ml-2 text-sm text-gray-600">({{ recipe.rating.toFixed(1) }} / 5)</span>
          </div>
        </div>

        <!-- Comments Section -->
        <div v-if="recipe.showComments" class="mt-4">
          <h3 class="text-lg font-semibold mb-2">Comments</h3>
          <div
            v-for="(comment, i) in recipe.comments"
            :key="i"
            class="p-2 border-b"
          >
            <p>{{ comment.text }}</p>
          </div>
          <!-- Add Comment -->
          <div class="mt-2">
            <textarea
              v-model="recipe.newComment"
              placeholder="Write a comment..."
              class="w-full p-2 border rounded"
            ></textarea>
            <button
              @click="addComment(index)"
              class="bg-blue-500 text-white px-4 py-2 rounded mt-2 hover:bg-blue-600"
            >
              Submit Comment
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- No Recipes Available -->
    <div v-else>
      <p class="text-center">No recipes available.</p>
    </div>

    <!-- Log In Prompt -->
    <div v-if="!isAuthenticated" class="text-center mt-4">
      <p>Please log in to like, comment, rate, or buy recipes.</p>
      <button @click="loginPrompt" class="bg-blue-500 text-white px-4 py-2 rounded mt-2 hover:bg-blue-600">
        Log In
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { useNuxtApp } from '#app';
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import gql from 'graphql-tag';
const router = useRouter();
const userId = ref('');


const { $apolloClient } = useNuxtApp();
const recipes = ref([]);
const authStore = useAuthStore(); // Access auth store

const isAuthenticated = computed(() => authStore.isAuthenticated); // Reactive check for authentication

const backendBaseUrl = 'http://localhost:8085/';

const FETCH_SHARED_RECIPES_QUERY = gql`
  query FetchSharedRecipes($userId: uuid!) {
    recipes(where: { shares: { is_shared: { _eq: true } } }) {
      id
      title
      featured_image
      ratings_aggregate {
        aggregate {
          avg {
            rating
          }
        }
      }
      user_ratings: ratings(where: { user_id: { _eq: $userId } }) {
        rating
      }
      likes_aggregate {
        aggregate {
          count
        }
      }
      likes(where: { user_id: { _eq: $userId } }) {
        user_id
      }
      comments(order_by: { created_at: desc }) {
        id
        comment
        created_at
        user_id
      }
    }
  }
`;


// Helper function to generate the full image URL
const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : null;
};

// Fetch shared recipes from the backend
const fetchRecipes = async () => {
  try {
    const response = await $apolloClient.query({
      query: FETCH_SHARED_RECIPES_QUERY,
      variables: { userId: authStore.user.id },
    });

    recipes.value = response.data.recipes.map((recipe) => ({
      ...recipe,
        rating: recipe.ratings_aggregate.aggregate.avg.rating || 0,

      likesCount: recipe.likes_aggregate.aggregate.count,
      liked: recipe.likes.length > 0, // If the user has liked the recipe
      bookmarked: false,
      rating: 0,
      showComments: false,
      comments: recipe.comments.map((comment) => ({
        text: comment.comment,
        created_at: comment.created_at,
      })),
      commentsCount: recipe.comments.length, // Store the number of comments
      newComment: '',
    }));
  } catch (error) {
    console.error('Error fetching recipes:', error);
  }
};



// Login prompt or login action
const loginPrompt = () => {
  // You can redirect to login page, open a modal, or trigger a login function
  alert('Redirecting to login page...');
  // Simulate login for now:
  // On successful login, you should call:
  // authStore.setUser(user); to store the user and update authentication status.
};

  // Like a Recipe
  const handleLike = async (index) => {
  if (!isAuthenticated.value) {
    alert('Please log in to like recipes.');
    return;
  }

  const recipe = recipes.value[index];

  try {
    if (!recipe.liked) {
      // Like the recipe
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation LikeRecipe($recipeId: uuid!, $userId: uuid!) {
            insert_likes_one(object: { recipe_id: $recipeId, user_id: $userId }) {
              id
            }
          }
        `,
        variables: { recipeId: recipe.id, userId: authStore.user.id },
      });

      if (response.data.insert_likes_one) {
        recipe.likesCount += 1;
        recipe.liked = true;
      }
    } else {
      // Unlike the recipe
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation UnlikeRecipe($recipeId: uuid!, $userId: uuid!) {
            delete_likes(
              where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
            ) {
              affected_rows
            }
          }
        `,
        variables: { recipeId: recipe.id, userId: authStore.user.id },
      });

      if (response.data.delete_likes.affected_rows > 0) {
        recipe.likesCount -= 1;
        recipe.liked = false;
      }
    }
  } catch (error) {
    console.error('Error toggling like:', error);
  }
};

const toggleComments = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  recipes.value[index].showComments = !recipes.value[index].showComments;
};

const addComment = async (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = recipes.value[index];
  if (recipe.newComment.trim()) {
    try {
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation AddComment($recipeId: uuid!, $userId: uuid!, $comment: String!) {
            insert_comments_one(object: { recipe_id: $recipeId, user_id: $userId, comment: $comment }) {
              id
              comment
              created_at
            }
          }
        `,
        variables: {
          recipeId: recipe.id,
          userId: authStore.user.id,
          comment: recipe.newComment,
        },
      });

      if (response.data.insert_comments_one) {
        // Optionally, add the new comment to the UI directly
        recipe.comments.push({
          text: response.data.insert_comments_one.comment,
          created_at: response.data.insert_comments_one.created_at,
        });
        recipe.newComment = ''; // Clear the input field
        // Update the comment count
        recipe.commentsCount += 1;
      }
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  }
};


const handleBookmark = (index) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  recipes.value[index].bookmarked = !recipes.value[index].bookmarked;
};

const handleRating = async (index, selectedRating) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = recipes.value[index];

  try {
    // GraphQL mutation for inserting or updating a rating
    const response = await $apolloClient.mutate({
      mutation: gql`
        mutation UpsertRating($recipeId: uuid!, $userId: uuid!, $rating: numeric!) {
          insert_ratings_one(
            object: { recipe_id: $recipeId, user_id: $userId, rating: $rating }
            on_conflict: {
              constraint: ratings_recipe_id_user_id_key
              update_columns: [rating]
            }
          ) {
            id
            rating
          }
        }
      `,
      variables: {
        recipeId: recipe.id,
        userId: authStore.user.id,
        rating: selectedRating,
      },
    });

    if (response.data.insert_ratings_one) {
      // Update the recipe's rating in the UI
      recipe.rating = selectedRating;
    }
  } catch (error) {
    console.error('Error updating rating:', error);
  }
};


const handleBuy = (title) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  alert(`Redirecting to purchase page for "${title}"...`);
};

// Fetch recipes when the component is mounted


onMounted(async () => {
  try {
    const token = localStorage.getItem('token');

    if (!token) {
      console.error('Token not found. Please log in again.');
      return;
    }

    const decodedToken = JSON.parse(atob(token.split('.')[1]));

    if (!decodedToken.userId) {
      console.error('User ID not found in token. Please log in again.');
      return;
    }

    userId.value = decodedToken.userId;
    authStore.setUser({ id: userId.value }); // Sync with the auth store

    await fetchRecipes();
  } catch (error) {
    console.error('Error during setup:', error);
  }
});





const login = (user) => {
  authStore.setUser(user); // Ensure this function updates the store with user details
  userId.value = user.id; // Sync `userId` ref
};




</script>

<style scoped>
/* Add custom styles if needed */
</style>
