import { defineConfig } from "vite";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import path from "node:path";

// https://vite.dev/config/
export default defineConfig({
  plugins: [tailwindcss(), svelte()],
  envPrefix: "CLAW_",
  base: "", // Ensures relative paths are used for assets. So server can serve from any path.
  resolve: {
    alias: {
      "#": path.resolve(__dirname, "src"),
    },
  },
  build: {
    outDir: "../cmd/claw/internal/webui",
    manifest: true,
    emptyOutDir: true,
  },
});
