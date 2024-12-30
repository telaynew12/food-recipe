<template>
  <div class="min-h-screen flex items-center justify-center bg-gray-100">
    <div class="bg-white p-8 rounded-lg shadow-lg w-96">
      <h1 class="text-2xl font-bold text-center mb-6">Login</h1>
      <form @submit.prevent="submitLogin">
        <!-- Email -->
        <div class="mb-4">
          <label for="email" class="block text-gray-700 font-medium mb-2">Email</label>
          <Field
            id="email"
            name="email"
            type="email"
            class="input-field"
            placeholder="Enter your email"
            :class="{ 'border-red-500': errors.email }"
          />
          <ErrorMessage name="email" class="text-red-500 text-sm mt-1" />
        </div>

        <!-- Password -->
        <div class="mb-6">
          <label for="password" class="block text-gray-700 font-medium mb-2">Password</label>
          <Field
            id="password"
            name="password"
            type="password"
            class="input-field"
            placeholder="Enter your password"
            :class="{ 'border-red-500': errors.password }"
          />
          <ErrorMessage name="password" class="text-red-500 text-sm mt-1" />
        </div>

        <!-- Submit Button -->
        <button
          type="submit"
          class="w-full bg-blue-500 hover:bg-blue-600 text-white py-2 px-4 rounded font-medium"
        >
          Login
        </button>

        <!-- Error Alert -->
        <div v-if="loginError" class="mt-4 text-red-500 text-center">
          {{ loginError }}
        </div>
      </form>
          <!-- Login Link -->
    <p v-if="!userVerified" class="text-center text-gray-600 mt-4">
      Don't have account? 
      <router-link to="/signup" class="text-blue-500 hover:underline">
        Signup here
      </router-link>
    </p>
    </div>
  </div>
</template>

<script setup>
import { useField, useForm, Field, ErrorMessage } from 'vee-validate';
import { object, string } from 'yup';
import { useApolloClient } from '@vue/apollo-composable';
import { useRouter } from 'vue-router';
import gql from 'graphql-tag';
import { ref } from 'vue';
import { useAuthStore } from '@/stores/auth'; // Pinia store

// Login GraphQL Mutation
const LOGIN_MUTATION = gql`
  mutation LoginUser($email: String!, $password: String!) {
    login(input: { credentials: { email: $email, password: $password } }) {
      accessToken
      userId
      role
    }
  }
`;

// Form Validation Schema
const loginSchema = object({
  email: string().required('Email is required').email('Invalid email format'),
  password: string().required('Password is required'),
});

// Vee Validate Form Setup
const { handleSubmit, errors } = useForm({ validationSchema: loginSchema });
const { value: email } = useField('email');
const { value: password } = useField('password');

// Apollo and Router
const apolloClient = useApolloClient().client;
const router = useRouter();

// Pinia Auth Store
const auth = useAuthStore();

// State for error handling
const loginError = ref(null);
// Form Submit Handler

const submitLogin = handleSubmit(async (values) => {
  loginError.value = null; // Reset error message
  console.log("Attempting to login with:", values);

  try {
    const { data } = await apolloClient.mutate({
      mutation: LOGIN_MUTATION,
      variables: { email: values.email, password: values.password },
    });

    console.log("Login successful:", data); // Handle success

    // Store the token in Pinia
    auth.setUser({
      token: data.login.accessToken,
      userId: data.login.userId,
      role: data.login.role,
    });

    // Store token in localStorage (optional)
    localStorage.setItem('token', data.login.accessToken);

    // Redirect to the dashboard with user ID as query parameter
   router.push({
      path: '/',
      query: { userId: data.login.userId

      },
     });
  } catch (error) {
    console.error("Login error:", error);
    loginError.value = error.message || 'Login failed';
  }
});

</script>

<style scoped>
.input-field {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 0.375rem;
  transition: border-color 0.2s;
}
.input-field:focus {
  border-color: #2563eb;
  outline: none;
}
</style>
