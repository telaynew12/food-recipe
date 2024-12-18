import { defineNuxtConfig } from 'nuxt/config';
import dotenv from 'dotenv';

// Load environment variables from .env file
dotenv.config();

export default defineNuxtConfig({
  // Modules
  modules: ['@nuxtjs/apollo'],

  // Apollo Client Configuration
  apollo: {
    clients: {
      default: {
        httpEndpoint: 'http://localhost:8083/v1/graphql',
        // Include the Admin Secret in the headers
        httpLinkOptions: {
          headers: {
            'x-hasura-admin-secret': process.env.HASURA_ADMIN_SECRET || 'Telay5870@', // Use environment variable
          },
        },
      },
    },
  },

  // Global CSS
  css: ['~/assets/css/tailwind.css'],
  
  // Router Configuration
  

  // Plugins
  plugins: ['~/plugins/apollo.js'], // Ensure this file exists

  // Build Configuration
  build: {
    transpile: ['@apollo/client'], // Ensure Apollo Client is transpiled
  },

  // PostCSS Configuration
  postcss: {
    plugins: {
      tailwindcss: {},
      autoprefixer: {},
    },
  },

  // App Metadata
  app: {
    head: {
      title: 'FOOD Recipe',
    },
  },

  // Compatibility
  compatibilityDate: '2024-11-25',
});
