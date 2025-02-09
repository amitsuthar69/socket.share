import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";
import path from "path";

export default defineConfig({
  plugins: [react()],
  resolve: {
    alias: {
      "@": path.resolve(__dirname, "./src"),
      wailsjs: path.resolve(__dirname, "./wailsjs"),
    },
  },
  build: {
    rollupOptions: {
      external: ["wailsjs/runtime", "wailsjs/go/models", "wailsjs/go/main/App"],
      output: {
        globals: {
          "wailsjs/runtime": "Wails",
          "wailsjs/go/models": "Models",
          "wailsjs/go/main/App": "App",
        },
      },
    },
  },
});
