// stores/categoryStore.js


import { defineStore } from 'pinia';
import piniaPersist from 'pinia-plugin-persistedstate';

export const useCategoryStore = defineStore('categoryStore', {
  state: () => ({
    categories: [],
    selectedCategory: null,
  }),
  actions: {
    setCategories(categories) {
      this.categories = categories;
    },
    setSelectedCategory(category) {
      this.selectedCategory = category;
    },
  },
  persist: true, // This will persist the store state across page refreshes
});
