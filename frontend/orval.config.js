module.exports = {
  api: {
    output: {
      mode: 'tags-split',
      target: 'src/api/generated.ts',
      schemas: 'src/api/model',
      client: 'react-query',
      override: {
        mutator: {
          path: 'src/api/mutator/custom-instance.ts',
          name: 'customInstance',
        },
        query: {
          useQuery: true,
          useInfinite: true,
          useInfiniteQueryParam: 'cursor',
          options: {
            staleTime: 10000,
          },
        },
      },
    },
    input: {
      target: './openapi.yaml',
      // If you don't have an OpenAPI spec yet, you can uncomment the following line
      // to use a mock server during development
      // target: 'https://stoplight.io/mocks/your-project/your-api/12345/api',
    },
  },
};