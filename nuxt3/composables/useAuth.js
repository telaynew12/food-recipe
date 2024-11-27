// composables/useAuth.js

import { useState } from '#imports';
import { useCookie } from '#app';
import { useRouter } from 'vue-router';

export const useAuth = () => {
  // Check if there's a stored token in the cookie or localStorage
  const authToken = useCookie('authToken').value;
  const isAuthenticated = useState('isAuthenticated', () => !!authToken); // Boolean value (true/false)
  const router = useRouter();

  // Function to logout
  const logout = () => {
    // Remove the token from cookie/localStorage
    useCookie('authToken').value = null; // Clear the cookie
    isAuthenticated.value = false; // Update the auth state to false
    router.push('/login'); // Redirect to login page
  };

  return {
    isAuthenticated,
    logout,
  };
};
