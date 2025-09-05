import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";

// https://vite.dev/config/
export default defineConfig({
  plugins: [svelte()],
  envPrefix: "CLAW_",
  base: "./", // Ensures relative paths are used for assets. So server can serve from any path.
  build: {
    outDir: "../cmd/claw/internal/webui",
    manifest: true,
    emptyOutDir: true,
  },
});
