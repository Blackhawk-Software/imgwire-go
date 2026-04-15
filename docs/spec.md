# Spec Pipeline

The Go SDK checks in two OpenAPI artifacts:

- `openapi/raw.openapi.json`: the raw backend snapshot.
- `openapi/sdk.openapi.json`: the server-shaped SDK contract emitted by `@imgwire/codegen-core` with `buildSdkSpec({ target: "go" })`.

Generation flow:

1. Load the raw OpenAPI source.
2. Shape it for the Go server SDK with `@imgwire/codegen-core`.
3. Generate the disposable Go client with OpenAPI Generator.
4. Post-process generated output to remove non-runtime files.
5. Update `CODEGEN_VERSION`.

Handwritten packages wrap the generated client but never modify generated sources in place.
