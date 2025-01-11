<template>
  <div class="max-w-lg mx-auto p-6 mt-12 border border-gray-300 rounded-lg shadow-md bg-white">
    <h2 class="text-2xl font-semibold text-center mb-6 text-gray-800">Register</h2>

    <!-- General Error Banner -->
    <div v-if="errors.general" class="mb-4 p-4 bg-red-100 text-red-700 border border-red-500 rounded-md">
      {{ errors.general }}
    </div>

    <!-- Registration Form -->
    <form v-if="!verificationPending" @submit.prevent="onSubmit" class="space-y-4">
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <input
          id="email"
          v-model="form.email"
          type="email"
          name="email"
          class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
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
          class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
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
          class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.password}"
        />
        <p v-if="errors.password" class="text-red-500 text-xs mt-1">{{ errors.password }}</p>
        <p v-if="passwordStrengthMessage" class="text-yellow-500 text-xs mt-1">{{ passwordStrengthMessage }}</p>
      </div>

      <div>
        <button
          type="submit"
          class="w-full py-2 px-4 bg-blue-600 text-white rounded-md hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 disabled:bg-blue-300 disabled:cursor-not-allowed"
          :disabled="loading"
        >
          Register
        </button>
      </div>
    </form>

    <!-- Verification Form -->
    <form v-if="verificationPending" @submit.prevent="onVerify" class="space-y-4">
      <div>
        <label for="verificationCode" class="block text-sm font-medium text-gray-700">Verification Code</label>
        <input
          id="verificationCode"
          v-model="verificationCode"
          type="text"
          name="verificationCode"
          class="mt-1 block w-full px-3 py-2 border rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm"
          :class="{'border-red-500': errors.verificationCode}"
        />
        <p v-if="errors.verificationCode" class="text-red-500 text-xs mt-1">{{ errors.verificationCode }}</p>
      </div>

      <div>
        <button
          type="submit"
          class="w-full py-2 px-4 bg-green-600 text-white rounded-md hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-green-500 disabled:bg-green-300 disabled:cursor-not-allowed"
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
import { ref, watch } from 'vue'
import { useMutation } from '@vue/apollo-composable'
import gql from 'graphql-tag'
import { useForm, defineRule, configure } from 'vee-validate'
import * as yup from 'yup'
import { useRouter } from 'vue-router';
const router = useRouter();

// Define custom validation rules
defineRule('strong_password', value => {
  const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/;
  if (!strongPasswordRegex.test(value)) {
    return 'Password Should be at least 8 characters long, include at least one uppercase letter, one lowercase letter, and one number.';
  }
  return true;
});

defineRule('letters_only', value => {
  const lettersOnlyRegex = /^[A-Za-z]+$/;
  if (!lettersOnlyRegex.test(value)) {
    return 'The name must contain letters only.';
  }
  return true;
});

defineRule('capitalized', value => {
  if (value && value[0] !== value[0].toUpperCase()) {
    return 'The name must start with a capital letter.';
  }
  return true;
});

// Configure Vee Validate to show all errors
configure({
  generateMessage: (context) => {
    const messages = {
      required: `${context.field} is required`,
      email: `Invalid email address`,
      min: `${context.field} must be at least ${context.rule.params[0]} characters`,
      strong_password: 'Password must be at least 8 characters long, include at least one uppercase letter, one lowercase letter, and one number',
      letters_only: 'The name must contain letters only',
      capitalized: 'The name must start with a capital letter'
    };

    return messages[context.rule.name]
      ? messages[context.rule.name]
      : `${context.field} is not valid`;
  },
  validateOnBlur: true,
  validateOnChange: true,
  validateOnInput: true,
  validateOnModelUpdate: true
});

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
const passwordStrengthMessage = ref('')
const verificationPending = ref(false) // Indicates if the user is in the verification step

// Form validation schema
const schema = yup.object({
  email: yup.string().email().required(),
  name: yup.string().required().matches(/^[A-Za-z]+$/, 'The name must contain letters only').matches(/^[A-Z][a-zA-Z]*$/, 'The name must start with a capital letter'),
  password: yup.string().required().matches(/^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/, 'Password must be at least 8 characters long, include at least one uppercase letter, one lowercase letter, and one number.')
})

const { handleSubmit, resetForm, setErrors } = useForm({
  validationSchema: schema,
  initialValues: form.value,
  validateOnMount: true
})

// Automatically capitalize the first letter of the name and validate letters only
watch(() => form.value.name, (newValue) => {
  if (newValue && /^[a-zA-Z]+$/.test(newValue)) {
    form.value.name = newValue.charAt(0).toUpperCase() + newValue.slice(1);
  }
})

// Watch for password changes to update strength message
watch(() => form.value.password, (newValue) => {
  const strongPasswordRegex = /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)[a-zA-Z\d]{8,}$/;
  passwordStrengthMessage.value = strongPasswordRegex.test(newValue)
    ? ''
    : 'Password must be at least 8 characters long, include at least one uppercase letter, one lowercase letter, and one number.';
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
  loading.value = true;
  try {
    const { data } = await verifyUser({
      input: { email: form.value.email, code: verificationCode.value }
    });

    // Success handling
    message.value = data.verify.message;
    verificationPending.value = false; // Verification completed, hide the form

    // Redirect to the login page after successful verification
    router.push({ name: 'login' }); // Adjust the route name/path as per your setup
  } catch (err) {
    // Handle verification errors
    errors.value.verificationCode = 'Invalid verification code';
  } finally {
    loading.value = false;
  }
};

</script>


<style scoped>
/* Additional custom styles if necessary */
</style>
