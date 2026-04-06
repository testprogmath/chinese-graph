import { ApolloClient, InMemoryCache, createHttpLink } from '@apollo/client'

// Get the GraphQL endpoint URL
const getGraphQLEndpoint = () => {
  // Use the Next.js rewrite path which handles routing to backend
  if (typeof window !== 'undefined') {
    // Client-side: use current origin with rewrite path
    return `${window.location.origin}/graphql`
  }
  
  // Server-side: use environment-specific URLs
  if (process.env.NODE_ENV === 'production') {
    return 'http://backend:8080/graphql'
  }
  
  return 'http://localhost:8080/graphql'
}

const httpLink = createHttpLink({
  uri: getGraphQLEndpoint(),
  credentials: 'include', // Include cookies for session management if needed
})

export const apolloClient = new ApolloClient({
  link: httpLink,
  cache: new InMemoryCache({
    typePolicies: {
      Word: {
        keyFields: ['id'],
      },
      WordNode: {
        keyFields: ['id'],
      },
      WordRelation: {
        keyFields: false, // WordRelations don't have IDs, use object identity
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      errorPolicy: 'all',
      fetchPolicy: 'cache-and-network',
    },
    query: {
      errorPolicy: 'all',
      fetchPolicy: 'cache-first',
    },
  },
})

export default apolloClient