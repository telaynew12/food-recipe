<template>
  <header class="bg-gray-800 text-white p-4">
    <nav class="container mx-auto flex justify-between items-center">
      <h1 class="text-lg font-bold">Recipe Categories</h1>
      <ul class="flex gap-4">
        <li 
          v-for="category in categories" 
          :key="category.id" 
          class="hover:underline cursor-pointer"
        >
          {{ category.name }}
        </li>
      </ul>
    </nav>
  </header>
</template>

<script setup>
import { useQuery } from '@vue/apollo-composable';
import gql from 'graphql-tag';

const GET_CATEGORIES = gql`
  query GetCategories {
    categories {
      id
      name
    }
  }
`;

const { result, loading, error } = useQuery(GET_CATEGORIES);

// Reactive categories list
const categories = computed(() => result.value?.Categories || []);
</script>

<style>
/* Add custom styles here if needed */
</style>
