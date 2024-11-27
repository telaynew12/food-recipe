import { defineNuxtPlugin } from '#app'
import { ApolloClient, InMemoryCache, HttpLink } from '@apollo/client/core'

export default defineNuxtPlugin((nuxtApp) => {
  // Create an Apollo Client instance
  const apolloClient = new ApolloClient({
    cache: new InMemoryCache(),
    link: new HttpLink({
      uri: 'http://localhost:8082/v1/graphql', // Replace with your GraphQL server URI
    }),
  })

  // Inject Apollo Client globally, available as 'apolloClient'
  nuxtApp.provide('apolloClient', apolloClient)
})
