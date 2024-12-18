<script setup>
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth'; // Pinia store for user data
import { ref, onMounted } from 'vue';

// Get route and auth store
const route = useRoute();
const router = useRouter();
const auth = useAuthStore();

// State for user data
const user = ref({
  name: '',
  email: '',
  userId: ''
});

// Dashboard stats data
const totalRecipes = ref(0);
const bookmarkedRecipes = ref(0);
const newComments = ref(0);

// On mounted, fetch user details from the store, query params, or localStorage
onMounted(() => {
  // Fetch user data from store or localStorage or query params
  const userIdFromQuery = route.query.userId;
  const nameFromQuery = route.query.name;
  const emailFromQuery = route.query.email;

  if (userIdFromQuery && nameFromQuery && emailFromQuery) {
    // If userId, name, and email are passed via query, set them
    user.value.userId = userIdFromQuery;
    user.value.name = nameFromQuery;
    user.value.email = emailFromQuery;
  } else if (auth.user) {
    // Otherwise, use stored user info from Pinia
    user.value = auth.user;
  } else {
    // If no user is found, check localStorage
    const storedUser = JSON.parse(localStorage.getItem('user'));
    if (storedUser) {
      user.value = storedUser;
    }
  }

  // Fetch data based on user ID
  if (user.value.userId) {
    fetchDashboardData(user.value.userId);
  } else {
    router.push('/login'); // Redirect to login if no user ID
  }
});

// Function to fetch dashboard data using userId
const fetchDashboardData = async (userId) => {
  try {
    const response = await fetch(`/api/dashboard/stats?userId=${userId}`);
    const data = await response.json();
    // Update the dashboard stats
    totalRecipes.value = data.totalRecipes;
    bookmarkedRecipes.value = data.bookmarkedRecipes;
    newComments.value = data.newComments;
  } catch (error) {
    console.error('Error fetching dashboard data:', error);
  }
};

// Handle logout action
const logout = () => {
  auth.clearUser();
  router.push('/login'); // Redirect to login page
};
</script>
