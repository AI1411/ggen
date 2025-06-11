export default {
  backend: {
    output: {
      mode: "split",
      target: "src/service/generated/client.ts",
      schemas: "src/service/generated/model",
      client: "react-query",
      mock: true,
      prettier: {
        semi: false,
        singleQuote: true,
        trailingComma: "es5",
      },
    },
    input: {
      target: "../backend/docs/api/swagger.yaml", // OpenAPI仕様ファイル
    },
  },
}
