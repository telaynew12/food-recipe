import { defineStore } from 'pinia';

export const useRecipeStore = defineStore('searchrecipe', {
  state: () => ({
    searchQuery: '', // State to store search query
  }),
  actions: {
    setSearchQuery(query) {
      this.searchQuery = query;
    },
  },
});
