import { defineStore } from 'pinia';

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: null,
    isAuthenticated: false,
  }),
  actions: {
    setUser(user) {
      this.user = user;
      this.isAuthenticated = !!user;
      if (typeof window !== 'undefined') {
        localStorage.setItem('user', JSON.stringify(user));
      }
    },
    clearUser() {
      this.user = null;
      this.isAuthenticated = false;
      if (typeof window !== 'undefined') {
        localStorage.removeItem('user');
      }
    },
    initializeAuth() {
      if (typeof window !== 'undefined') {
        const user = localStorage.getItem('user');
        this.user = user ? JSON.parse(user) : null;
        this.isAuthenticated = !!user;
      }
    },
  },
  getters: {
    userId(state) {
      return state.user?.id || null; // Return the user ID if available
    },
  },
});
