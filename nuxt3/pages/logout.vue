<template>
  <div class="w-full h-screen p-6 bg-white shadow-md rounded-lg mt-10 overflow-auto">
    <!-- Display All Recipes -->
    <div v-if="filteredRecipes.length" class="space-y-8">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div
          v-for="(recipe, index) in filteredRecipes"
          :key="recipe.id"
          class="border p-4 rounded-lg shadow-sm"
        >
          <!-- Recipe Title -->
          <h2 class="text-2xl font-bold mb-4">{{ recipe.title }}</h2>

          <!-- Recipe Image -->
          <img
            v-if="recipe.featured_image"
            :src="getImageUrl(recipe.featured_image)"
            alt="Recipe Image"
  class="w-full h-48 object-cover rounded-lg mb-4 cursor-pointer"
            @click="navigateToDetails(recipe.id)"
          />

          <!-- Buttons Section -->
          <div class="flex justify-start gap-6 mb-4 items-center">
            <!-- Like Button -->
            <button
              @click="handleLike(index)"
              :class="{
                'text-blue-500 font-bold': recipe.liked,
                'text-gray-600': !recipe.liked
              }"
              class="flex items-center space-x-2"
            >
              <span>👍</span>
              <span>{{ recipe.liked ? 'Liked' : 'Like' }}</span>
              <span>({{ recipe.likesCount }})</span>
            </button>

            <!-- Comments Button -->
            <button
              @click="toggleComments(index)"
              class="flex items-center space-x-2 text-gray-600 hover:text-gray-800"
            >
              <span>💬</span>
              <span>Comments</span>
              <span>({{ recipe.commentsCount }})</span>
            </button>

            <!-- Bookmark Button -->
            <button
              @click="toggleBookmark(recipe.id, index)"
              :class="{
                'text-yellow-500': recipe.bookmarked,
                'text-gray-600': !recipe.bookmarked
              }"
            >
              {{ recipe.bookmarked ? 'Bookmarked' : 'Bookmark' }}
            </button>
          </div>

          <!-- Rating Section -->
          <div class="flex items-center space-x-1 mb-4">
            <span class="font-semibold">Rate:</span>
            <div class="flex">
              <span
                v-for="star in 5"
                :key="star"
                :class="{
                  'text-yellow-400': recipe.currentRating >= star,
                  'text-gray-400': recipe.currentRating < star
                }"
                class="cursor-pointer text-2xl"
                @click="handleRating(index, star)"
              >
                ★
              </span>
            </div>
            <span
              v-if="recipe.hasRated"
              class="ml-2 text-sm text-green-500"
            >
              You rated this recipe: {{ recipe.currentRating }}
            </span>
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
            <div v-if="isAuthenticated" class="mt-2">
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
            <div v-else class="text-sm text-gray-600">
              Log in to add a comment.
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- No Recipes Available -->
    <div v-else>
      <p class="text-center">No recipes available.</p>
    </div>
  </div>
</template>



<script setup>
import { useRouter } from 'vue-router';
import { onMounted, computed, ref } from 'vue';
import { useNuxtApp } from '#app';
import { useAuthStore } from '@/stores/auth'; // Import the auth store
import { useRecipeStore } from '@/stores/recipe'; // Import the recipe store
import gql from 'graphql-tag';
import { useDetailStore } from '@/stores/detail'; // Import the detail store
const detailStore = useDetailStore(); // Access the detail store

const router = useRouter();
const userId = ref('');
const { $apolloClient } = useNuxtApp();
const recipes = ref([]);
const authStore = useAuthStore();
const recipeStore = useRecipeStore(); // Access the recipe store
const getImageUrl = (path) => {
  return path ? `${backendBaseUrl}${path}` : null;
};

const isAuthenticated = computed(() => authStore.isAuthenticated); // Reactive check for authentication

const backendBaseUrl = 'http://localhost:8085/';

const searchQuery = computed({
  get: () => recipeStore.searchQuery,
  set: (query) => recipeStore.setSearchQuery(query),
});

// Computed property for filtered recipes
const filteredRecipes = computed(() => {
  const query = searchQuery.value.toLowerCase();
  return recipes.value.filter((recipe) =>
    recipe.title.toLowerCase().includes(query)
  );
});
// Function to navigate to the recipe details page and store the recipeId
const navigateToDetails = (recipeId) => {
  // Store the recipeId in the detail store
  detailStore.setRecipeId(recipeId);
  
  // Navigate to the details page
  router.push({ name: 'recipeDetail', params: { id: recipeId } });
};


const FETCH_SHARED_RECIPES_QUERY = gql`
  query FetchSharedRecipes($userId: uuid!) {
    recipes(where: { shares: { is_shared: { _eq: true } } }) {
      id
      title
      featured_image
      user_id 
      bookmarks(where: { user_id: { _eq: $userId } }) {
        user_id
      }
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


const fetchRecipes = async () => {
  try {
    const response = await $apolloClient.query({
      query: FETCH_SHARED_RECIPES_QUERY,
      variables: { userId: authStore.user.id },
    });

    recipes.value = response.data.recipes.map((recipe) => {
      const userRating = recipe.user_ratings.length > 0 ? recipe.user_ratings[0].rating : null;

      return {
        ...recipe,
        bookmarked: recipe.bookmarks.length > 0,
        rating: recipe.ratings_aggregate.aggregate.avg.rating || 0,
        likesCount: recipe.likes_aggregate.aggregate.count,
        liked: recipe.likes.length > 0,
        hasRated: userRating !== null,
        currentRating: userRating || 0,
        showComments: false,
        comments: recipe.comments.map((comment) => ({
          text: comment.comment,
          created_at: comment.created_at,
        })),
        commentsCount: recipe.comments.length,
        newComment: '',
      };
    });
  } catch (error) {
    console.error('Error fetching recipes:', error);
  }
};

const loginPrompt = () => {
  router.push('/login');
};

const handleLike = async (index) => {
  if (!isAuthenticated.value) {
    alert('Please log in to like recipes.');
    return;
  }

  const recipe = filteredRecipes.value[index]; // Access filtered recipe

  // Check if the logged-in user is the same as the recipe owner
  if (recipe.user_id === authStore.user.id) {
    alert('You cannot like your own recipe.');
    return;
  }

  try {
    if (!recipe.liked) {
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
  const newComment = recipe.newComment.trim();

  if (newComment) {
    try {
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation AddComment($recipeId: uuid!, $userId: uuid!, $comment: String!) {
            insert_comments_one(object: { recipe_id: $recipeId, user_id: $userId, comment: $comment }) {
              id
            }
          }
        `,
        variables: {
          recipeId: recipe.id,
          userId: authStore.user.id,
          comment: newComment,
        },
      });

      if (response.data.insert_comments_one) {
        recipe.comments.push({ text: newComment, created_at: new Date().toISOString() });
        recipe.commentsCount += 1;
        recipe.newComment = ''; // Reset comment input field
      }
    } catch (error) {
      console.error('Error adding comment:', error);
    }
  }
};




const handleRating = async (index, selectedRating) => {
  if (!isAuthenticated.value) {
    loginPrompt(); // Trigger login if not authenticated
    return;
  }

  const recipe = filteredRecipes.value[index];  // Use 


    // Check if the logged-in user is the owner of the recipe
  if (recipe.user_id === authStore.user.id) {
    alert("You cannot rate your own recipe.");
    return;
  }

  try {
    const response = await $apolloClient.mutate({
      mutation: gql`
        mutation UpsertRating($recipeId: uuid!, $userId: uuid!, $rating: numeric!) {
          insert_ratings_one(
            object: { recipe_id: $recipeId, user_id: $userId, rating: $rating }
            on_conflict: {
              constraint: ratings_recipe_id_user_id_key,
              update_columns: [rating]
            }
          ) {
            recipe {
              ratings_aggregate {
                aggregate {
                  avg {
                    rating
                  }
                }
              }
            }
          }
        }
      `,
      variables: {
        recipeId: recipe.id,
        userId: authStore.user.id,
        rating: selectedRating,
      },
    });

    // Update frontend state
    if (response.data.insert_ratings_one) {
      const updatedAvgRating =
        response.data.insert_ratings_one.recipe.ratings_aggregate.aggregate.avg.rating;

      recipe.rating = updatedAvgRating || 0; // Update average rating
      recipe.hasRated = true; // Mark as rated
      recipe.currentRating = selectedRating; // Update user's specific rating
    }
  } catch (error) {
    console.error('Error updating rating:', error);
  }
};

const toggleBookmark = async (recipeId, index) => {
  if (!isAuthenticated.value) {
    loginPrompt();
    return;
  }

  const recipe = filteredRecipes.value[index];  // Use filteredRecipes
    // Check if the logged-in user is the owner of the recipe
  if (recipe.user_id === authStore.user.id) {
    alert("You cannot  bookmark your own recipe.");
    return;
  }

  try {
    if (!recipe.bookmarked) {
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation BookmarkRecipe($recipeId: uuid!, $userId: uuid!) {
            insert_bookmarks_one(object: { recipe_id: $recipeId, user_id: $userId }) {
              id
            }
          }
        `,
        variables: { recipeId, userId: authStore.user.id },
      });

      if (response.data.insert_bookmarks_one) {
        recipe.bookmarked = true;
      }
    } else {
      const response = await $apolloClient.mutate({
        mutation: gql`
          mutation RemoveBookmark($recipeId: uuid!, $userId: uuid!) {
            delete_bookmarks(
              where: { recipe_id: { _eq: $recipeId }, user_id: { _eq: $userId } }
            ) {
              affected_rows
            }
          }
        `,
        variables: { recipeId, userId: authStore.user.id },
      });

      if (response.data.delete_bookmarks.affected_rows > 0) {
        recipe.bookmarked = false;
      }
    }
  } catch (error) {
    console.error('Error toggling bookmark:', error);
  }
};

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
    if (isAuthenticated.value === true || isAuthenticated.value === false) {
    await fetchRecipes();
}

   
  } catch (error) {
    console.error('Error during setup:', error);
  }
});
</script>

<style scoped>
/* Add any custom styling here */
</style>
