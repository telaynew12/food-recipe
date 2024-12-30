<template>
  <div class="max-w-lg mx-auto p-6 mt-12 border border-gray-300 rounded-lg shadow-md bg-white">
    <h2 class="text-2xl font-semibold text-center mb-6 text-gray-800">Register</h2>

    <!-- Registration Form -->
    <form v-if="!verificationPending" @submit.prevent="onSubmit" class="space-y-4">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="email"
          v-model="form.email"
          type="email"
          name="email"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.email}"
        />
        <p v-if="errors.email" class="text-red-500 text-xs mt-1">{{ errors.email }}</p>
      </div>

      <div>
        <label for="name" class="block text-sm font-medium text-gray-700">Name</label>
        <input
          id="name"
          v-model="form.name"
          type="text"
          name="name"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.name}"
        />
        <p v-if="errors.name" class="text-red-500 text-xs mt-1">{{ errors.name }}</p>
      </div>

      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
        <input
          id="password"
          v-model="form.password"
          type="password"
          name="password"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.password}"
        />
        <p v-if="errors.password" class="text-red-500 text-xs mt-1">{{ errors.password }}</p>
      </div>

      <div>
        <button
          type="submit"
          class="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          :disabled="loading"
        >
          Register
        </button>
      </div>
    </form>
        <!-- Login Link -->
    <p v-if="!userVerified" class="text-center text-gray-600 mt-4">
      Already have an account? 
      <router-link to="/login" class="text-blue-500 hover:underline">
        Login here
      </router-link>
    </p>

    <!-- Verification Form -->
    <form v-if="verificationPending" @submit.prevent="onVerify" class="space-y-4">
      <div>
        <label for="verificationCode" class="block text-sm font-medium text-gray-700">Verification Code</label>
        <input
          id="verificationCode"
          v-model="verificationCode"
          type="text"
          name="verificationCode"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.verificationCode}"
        />
        <p v-if="errors.verificationCode" class="text-red-500 text-xs mt-1">{{ errors.verificationCode }}</p>
      </div>

      <div>
        <button
          type="submit"
          class="w-full py-2 px-4 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500"
          :disabled="loading"
        >
          Verify Code
        </button>
      </div>
    </form>

    <p v-if="message" class="text-green-500 mt-4">{{ message }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { useForm } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from 'vue-router';
const router = useRouter();

// Form data and states
const form = ref({
  email: '',
  name: '',
  password: ''
})

const verificationCode = ref('')
const message = ref('')
const loading = ref(false)
const errors = ref({})
const verificationPending = ref(false) // Indicates if the user is in the verification step

// Form validation schema
const schema = yup.object({
  email: yup.string().email('Invalid email').required('Email is required'),
  name: yup.string().required('Name is required'),
  password: yup.string().min(6, 'Password must be at least 6 characters').required('Password is required')
})

const { handleSubmit, resetForm, setErrors } = useForm({
  validationSchema: schema,
  initialValues: form.value
})

// GraphQL mutations
const REGISTER_USER = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      message
    }
  }
`

const VERIFY_USER = gql`
  mutation Verify($input: VerifyInput!) {
    verify(input: $input) {
      message
    }
  }
`


const { mutate: registerUser, loading: mutationLoading } = useMutation(REGISTER_USER)
const { mutate: verifyUser, loading: verificationLoading } = useMutation(VERIFY_USER)

const onSubmit = async () => {
  loading.value = true
  try {
    const { data } = await registerUser({
      input: form.value
    })

    message.value = data.register.message
    verificationPending.value = true // Show the verification form
    resetForm()
  } catch (err) {
    // Handle GraphQL or network errors
    setErrors(err.graphQLErrors[0]?.extensions?.exception || {})
    errors.value.general = 'Registration failed, please try again.'
  } finally {
    loading.value = false
  }
}

const onVerify = async () => {
  loading.value = true
  try {
    const { data } = await verifyUser({
      input: { email: form.value.email, code: verificationCode.value }
    })

    // Success handling
    message.value = data.verify.message
    verificationPending.value = false // Verification completed, hide the form
        router.push('/login');

  } catch (err) {
    // Handle verification errors
    errors.value.verificationCode = 'Invalid verification code'
  } finally {
    loading.value = false
  }
}

</script>

<style scoped>
/* Additional custom styles if necessary */
</style>
