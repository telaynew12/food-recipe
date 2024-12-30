<template>
  <nav class="fixed top-12 left-0 w-full bg-green-500 p-4 md:p-6 z-10">
    <div class="container mx-auto flex justify-between items-center">
      <ul class="flex flex-wrap gap-4 md:gap-6 text-lg text-white">
        <li
          v-for="category in categories"
          :key="category.id"
          class="hover:text-blue-300 cursor-pointer"
          @click="handleCategoryClick(category.name)"
        >
          {{ category.name }}
        </li>
      </ul>
    </div>
  </nav>
</template>

<script setup>
  import { useRouter } from 'vue-router';
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';
import { useCategoryStore } from '@/stores/categoryStore'; // Assuming Pinia is used
const router = useRouter();

// GraphQL query to fetch categories
const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`;

// Fetch categories from the GraphQL server
const { result, loading, error } = useQuery(GET_CATEGORIES);

// Pinia store
const categoryStore = useCategoryStore();

// Watch the query result and update the store
watch(
  () => result.value?.categories,
  (newCategories) => {
    if (newCategories) {
      categoryStore.setCategories(newCategories.map((category) => category.name));
    }
  },
  { immediate: true }
);

// Reactive categories list for this component
const categories = computed(() => result.value?.categories || []);

// Handle category clicks
const handleCategoryClick = (categoryName) => {
  console.log('Selected Category:', categoryName);
  categoryStore.setSelectedCategory(categoryName); // Assuming your store has a `setSelectedCategory` method.
  router.push({ name: 'CategoryPage', params: { category: categoryName } });
};

</script>

<style scoped>
/* Add any additional custom styles if needed */
</style>
