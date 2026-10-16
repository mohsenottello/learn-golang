import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  build: { emptyOutDir: false },
  server: {
    port: 5173,
    // In dev, forward API calls to the Go server so cookies stay same-origin.
    proxy: { "/api": "http://localhost:8080" },
  },
});
