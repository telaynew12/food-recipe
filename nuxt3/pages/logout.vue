<template>
  <div class="max-w-lg mx-auto p-6 mt-12 border border-gray-300 rounded-lg shadow-md bg-white">
    <h2 class="text-2xl font-semibold text-center mb-6 text-gray-800">Register</h2>

    <!-- Registration Form -->
    <form @submit.prevent="onSubmit" class="space-y-4">
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
        <label for="confirmPassword" class="block text-sm font-medium text-gray-700">Confirm Password</label>
        <input
          id="confirmPassword"
          v-model="form.confirmPassword"
          type="password"
          name="confirmPassword"
          class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.confirmPassword}"
        />
        <p v-if="errors.confirmPassword" class="text-red-500 text-xs mt-1">{{ errors.confirmPassword }}</p>
      </div>

      <div>
        <button
          type="submit"
          class="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
          :disabled="loading || hasPasswordErrors"
        >
          Register
        </button>
      </div>
    </form>

    <!-- Verification Form -->
    <div v-if="registrationSuccess && !userVerified" class="mt-6">
      <h3 class="text-xl font-semibold text-center mb-4 text-gray-800">Enter Verification Code</h3>
      
      <form @submit.prevent="verifyCodeMethod" class="space-y-4">
        <div>
          <label for="verificationCode" class="block text-sm font-medium text-gray-700">Verification Code</label>
          <input
            id="verificationCode"
            v-model="verificationCode"
            type="text"
            class="mt-1 block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          />
        </div>

        <div>
          <button
            type="submit"
            class="w-full py-2 px-4 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500"
          >
            Verify Code
          </button>
        </div>
      </form>
    </div>

    <p v-if="message" class="text-green-500 mt-4">{{ message }}</p>
  </div>
</template>
<script setup>
import { ref, computed } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import * as yup from 'yup'

const form = ref({
  email: '',
  name: '',
  password: '',
  confirmPassword: ''
})
const message = ref('')
const loading = ref(false)
const errors = ref({})
const registrationSuccess = ref(false)
const userVerified = ref(false)
const verificationCode = ref('')

const hasPasswordErrors = computed(() => form.value.password !== form.value.confirmPassword)

const schema = yup.object({
  email: yup.string().email('Invalid email').required('Email is required'),
  name: yup.string().matches(/^[A-Za-z]/, 'Name must start with a letter').required('Name is required'),
  password: yup.string()
    .min(6, 'Password must be at least 6 characters')
    .matches(/[A-Z]/, 'Password must contain at least one uppercase letter')
    .matches(/[a-z]/, 'Password must contain at least one lowercase letter')
    .matches(/[0-9]/, 'Password must contain at least one number')
    .required('Password is required'),
  confirmPassword: yup.string()
    .oneOf([yup.ref('password'), null], 'Passwords must match')
    .required('Confirm Password is required')
})

const REGISTER_USER = gql`
  mutation Register($input: RegisterInput!) {
    register(input: $input) {
      message
    }
  }
`

const VERIFY_CODE = gql`
  mutation VerifyCode($email: String!, $code: String!) {
    verifyCode(email: $email, code: $code) {
      message
      userVerified
    }
  }
`

const { mutate: registerUser } = useMutation(REGISTER_USER)
const { mutate: verifyCode } = useMutation(VERIFY_CODE)

const onSubmit = async () => {
  loading.value = true
  try {
    await schema.validate(form.value, { abortEarly: false })
    const { data } = await registerUser({ input: form.value })
    message.value = data.register.message
    registrationSuccess.value = true
  } catch (err) {
    if (err instanceof yup.ValidationError) {
      err.inner.forEach(e => errors.value[e.path] = e.message)
    } else {
      message.value = 'Registration failed.'
    }
  } finally {
    loading.value = false
  }
}

const verifyCodeMethod = async () => {
  loading.value = true
  try {
    const { data } = await verifyCode({ email: form.value.email, code: verificationCode.value })
    userVerified.value = data.verifyCode.userVerified
    message.value = userVerified.value ? 'Account verified successfully!' : 'Invalid verification code.'
  } catch {
    message.value = 'Verification failed.'
  } finally {
    loading.value = false
  }
}
</script>
