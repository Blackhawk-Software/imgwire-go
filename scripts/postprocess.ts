import { readdir, rm } from "node:fs/promises";
import { resolve } from "node:path";

import { GENERATED_DIR } from "./_lib/paths.ts";

const GENERATED_ROOT_ENTRIES_TO_REMOVE = new Set([
  ".gitignore",
  ".openapi-generator-ignore",
  ".travis.yml",
  "README.md",
  "api",
  "docs",
  "git_push.sh",
  "go.mod",
  "go.sum",
  "test"
]);

export async function runPostprocess(options?: {
  generatedDir?: string;
}): Promise<void> {
  const generatedDir = options?.generatedDir ?? GENERATED_DIR;

  for (const name of GENERATED_ROOT_ENTRIES_TO_REMOVE) {
    await rm(resolve(generatedDir, name), { force: true, recursive: true });
  }

  await rm(resolve(generatedDir, ".openapi-generator"), {
    force: true,
    recursive: true
  });

  const entries = await readdir(generatedDir, { withFileTypes: true });
  for (const entry of entries) {
    if (entry.isDirectory() && entry.name.startsWith(".swagger-codegen")) {
      await rm(resolve(generatedDir, entry.name), {
        force: true,
        recursive: true
      });
    }
  }
}

if (import.meta.url === `file://${process.argv[1]}`) {
  await runPostprocess();
}
