<template>
    <div>
      <h1>File Upload</h1>
      <input type="file" @change="handleFileChange" />
      <button @click="uploadFile" :disabled="!selectedFile">Upload</button>
      <p v-if="uploadStatus">{{ uploadStatus }}</p>
    </div>
  </template>
  
  <script setup>
  import { ref } from 'vue'
  import { useMutation } from '@vue/apollo-composable'
  import gql from 'graphql-tag'
  
  // GraphQL mutation for file upload
  const FILE_UPLOAD = gql`
    mutation fileUpload($name: String!, $type: String!, $base64str: String!) {
      fileUpload(name: $name, type: $type, base64str: $base64str) {
        file_path
      }
    }
  `
  
  const { mutate } = useMutation(FILE_UPLOAD)
  
  const selectedFile = ref(null)
  const uploadStatus = ref('')
  
  const handleFileChange = async (event) => {
    const file = event.target.files[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        selectedFile.value = {
          name: file.name,
          type: file.type,
          base64str: e.target.result.split(',')[1],
        }
      }
      reader.readAsDataURL(file)
    }
  }
  
  const uploadFile = async () => {
    if (!selectedFile.value) {
      uploadStatus.value = 'No file selected!'
      return
    }
  
    try {
      uploadStatus.value = 'Uploading...'
      const { data } = await mutate({
        name: selectedFile.value.name,
        type: selectedFile.value.type,
        base64str: selectedFile.value.base64str,
      })
      uploadStatus.value = `File uploaded successfully: ${data.fileUpload.file_path}`
    } catch (error) {
      uploadStatus.value = `Error: ${error.message}`
    }
  }
  </script>
  