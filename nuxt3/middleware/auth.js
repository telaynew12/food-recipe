// // middleware/auth.js
// export default defineNuxtRouteMiddleware((to, from) => {
//     const isAuthenticated = useState('isAuthenticated').value;
  
//     if (!isAuthenticated && to.path !== '/login' && to.path !== '/signup') {
//       return navigateTo('/login'); // Redirect if not authenticated
//     }
//   });
  