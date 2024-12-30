// stores/categoryStore.js


import { defineStore } from 'pinia';

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
});
