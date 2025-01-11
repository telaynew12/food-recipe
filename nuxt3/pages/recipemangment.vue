<template>
  <div v-if="user" class="p-6 max-w-7xl mx-auto">
    <h2 class="text-lg mb-8 text-gray-600">Welcome:<strong>{{ user.name }}</strong> </h2>
    <div v-if="recipes && recipes.length > 0">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-6">
        <div v-for="recipe in recipes" :key="recipe.id" class="p-4 border rounded-lg shadow-sm hover:shadow-md">
          <h3 class="text-lg font-semibold mb-2 text-gray-800">{{ recipe.title }}</h3>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Likes -->
            <div>
              <h4 class="font-bold text-gray-700">Liked by:</h4>
              <ul class="text-sm text-gray-600">
                <li v-for="(like, index) in recipe.likes" :key="like.id">
                  {{ index + 1 }}. {{ like.user.name }}
                </li>
              </ul>
            </div>

            <!-- Ratings -->
            <div>
              <h4 class="font-bold text-gray-700">Rated by:</h4>
              <ul class="text-sm text-gray-600">
                <li v-for="(rate, index) in recipe.ratings" :key="rate.id">
                  {{ index + 1 }}. {{ rate.user.name }}
                </li>
              </ul>
            </div>

            <!-- Bookmarks -->
            <div>
              <h4 class="font-bold text-gray-700">Bookmarked by:</h4>
              <ul class="text-sm text-gray-600">
                <li v-for="(bookmark, index) in recipe.bookmarks" :key="bookmark.id">
                  {{ index + 1 }}. {{ bookmark.user.name }}
                </li>
              </ul>
            </div>

            <!-- Comments -->
            <div>
              <h4 class="font-bold text-gray-700">Comments:</h4>
              <ul class="text-sm text-gray-600">
                <li v-for="(comment, index) in recipe.comments" :key="comment.id">
                  <p class="text-gray-800">{{ index + 1 }}. {{ comment.user.name }}</p>
                  <p class="italic text-gray-500">"{{ comment.comment }}"</p>
                </li>
              </ul>
            </div>
          </div>
        </div>
      </div>
    </div>
    <div v-else class="text-center text-gray-500 mt-10">
      <p>No recipes found for this user. Please check again later.</p>
    </div>
  </div>
  <div v-else class="text-center text-gray-500">
    <p>User ID not available. Please log in to view your recipes.</p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { gql } from '@apollo/client/core';
import { useNuxtApp } from '#app';

// Access the authentication store
const authStore = useAuthStore();
const userId = computed(() => authStore.userId);

// Reactive variable to store user and recipes
const user = ref(null);
const recipes = ref([]);

// GraphQL query to fetch recipes with likes and user names
const FETCH_RECIPES_QUERY = gql`
  query FetchRecipes($userId: uuid!) {
    recipes(where: { user_id: { _eq: $userId } }) {
      id
      title
      description
      featured_image
      comments {
        id
        comment
        user {
          id
          name
        }
      }
      bookmarks {
        id
        user {
          id
          name
        }
      }
      likes {
        id
        user {
          id
          name
        }
      }
      ratings {
        id
        user {
          id
          name
        }
      }
    }
  }
`;

// Fetch recipes and user data function
const fetchData = async () => {
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

    // Fetch user details
    const { data: userData } = await client.query({
      query: gql`
        query FetchUser($userId: uuid!) {
          users(where: { id: { _eq: $userId } }) {
            id
            name
          }
        }
      `,
      variables: { userId: userId.value },
    });

    if (userData && userData.users.length > 0) {
      user.value = userData.users[0];
    } else {
      console.error('User not found.');
    }
  } catch (error) {
    console.error('Error fetching data:', error.message);
  }
};

// Fetch data when the component is mounted
onMounted(() => {
  if (userId.value) {
    fetchData();
  } else {
    console.error('User ID is not available.');
  }
});
</script>
