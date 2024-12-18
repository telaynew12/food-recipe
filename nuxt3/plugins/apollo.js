
import { defineNuxtPlugin } from '#app';
import { ApolloClient, InMemoryCache } from '@apollo/client/core';
import createUploadLink from "apollo-upload-client/createUploadLink.mjs";
// import createUploadLink from "apollo-upload-client/formDataAppendFile.mjs";
// import createUploadLink from "apollo-upload-client/isExtractableFile.mjs";
export default defineNuxtPlugin((nuxtApp) => {
  const uploadLink = createUploadLink({
    uri: 'http://localhost:8083/v1/graphql', // Your GraphQL endpoint
    headers: {
      'x-hasura-admin-secret': process.env.HASURA_ADMIN_SECRET || 'Telay5870@',
    },
  });

  const apolloClient = new ApolloClient({
    link: uploadLink, // Use the upload link for file handling
    cache: new InMemoryCache(),
  });

  nuxtApp.provide('apolloClient', apolloClient); // Make Apollo Client available globally
});
