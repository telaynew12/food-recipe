// stores/detail.js
import { defineStore } from 'pinia';

export const useDetailStore = defineStore('detail', {
  state: () => ({
    recipeId: null,
  }),
  actions: {
    setRecipeId(id) {
      this.recipeId = id;
    },
  },
});
