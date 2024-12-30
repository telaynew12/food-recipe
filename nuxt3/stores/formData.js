// stores/formData.js
import { defineStore } from 'pinia';

export const useFormDataStore = defineStore('formData', {
  state: () => ({
    formData: {
      name: '',
      email: '',
    },
  }),

  actions: {
    updateFormData(key, value) {
      this.formData[key] = value; // Dynamically update the form data
    },
    resetFormData() {
      this.formData = {
        name: '',
        email: '',
      };
    },
  },
});
