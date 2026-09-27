import tailwindcss from "@tailwindcss/vite";
import adapter from "@sveltejs/adapter-static";
import { sveltekit } from "@sveltejs/kit/vite";
import { defineConfig } from "vitest/config";

export default defineConfig({
  plugins: [
    tailwindcss(),
    sveltekit({
      // SPA: every route falls back to the client shell, which the Go server serves as index.html.
      adapter: adapter({ fallback: "index.html" }),
      compilerOptions: {
        // Runes everywhere except in libraries.
        runes: ({ filename }) => (filename.split(/[/\\]/).includes("node_modules") ? undefined : true),
      },
    }),
  ],
  server: {
    // Development is same-origin like production: the browser talks only to Vite, which forwards
    // /api to the Go server. The object form keeps the Host header, which the backend's csrf
    // middleware compares with Origin; the string form would rewrite it and writes would fail.
    proxy: {
      "/api": { target: "http://localhost:3000" },
    },
  },
  test: {
    include: ["src/**/*.test.ts"],
    environment: "node",
  },
});
