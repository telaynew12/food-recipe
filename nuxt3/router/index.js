import { createRouter, createWebHistory } from 'vue-router';
import Home from '../views/Home.vue';
import RecipeDetail from '../views/RecipeDetail.vue'; // Adjust path as needed

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home,
  },
  {
    path: '/recipe/:id', // Recipe detail page with dynamic id parameter
    name: 'RecipeDetail',
    component: RecipeDetail,
    props: true, // This ensures the 'id' is passed as a prop to the component
  },
];

const router = createRouter({
  history: createWebHistory(process.env.BASE_URL),
  routes,
});

export default router;
