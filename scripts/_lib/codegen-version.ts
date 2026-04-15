import { createHash } from "node:crypto";
import { readFile } from "node:fs/promises";

import { readJson } from "./fs.ts";

type PackageLike = {
  dependencies?: Record<string, string>;
  devDependencies?: Record<string, string>;
};

export async function computeCodegenVersion(input: {
  packageJsonPath: string;
  goModPath: string;
  rawOpenApiPath: string;
  sdkOpenApiPath: string;
  generatorConfigPath: string;
}): Promise<string> {
  const packageJson = await readJson<PackageLike>(input.packageJsonPath);
  const [goMod, raw, sdk, generatorConfig] = await Promise.all([
    readFile(input.goModPath, "utf8"),
    readFile(input.rawOpenApiPath, "utf8"),
    readFile(input.sdkOpenApiPath, "utf8"),
    readFile(input.generatorConfigPath, "utf8"),
  ]);

  const hash = createHash("sha256");
  hash.update(
    JSON.stringify({
      codegenCore:
        packageJson.dependencies?.["@imgwire/codegen-core"] ??
        packageJson.devDependencies?.["@imgwire/codegen-core"] ??
        "missing",
      openapiGeneratorCli:
        packageJson.devDependencies?.["@openapitools/openapi-generator-cli"] ??
        "missing",
      typescript: packageJson.devDependencies?.typescript ?? "missing",
      goMod,
      raw,
      sdk,
      generatorConfig,
    }),
  );

  return hash.digest("hex");
}
