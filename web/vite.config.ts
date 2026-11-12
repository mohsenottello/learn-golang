import { defineConfig } from "vite";
import react from "@vitejs/plugin-react";

export default defineConfig({
  plugins: [react()],
  server: {
    port: 5173,
    host: true,                                 // listen on all interfaces, not just 127.0.0.1
    allowedHosts: ["host.docker.internal"],     // accept Gotenberg's Host header
    proxy: { "/api": "http://localhost:8080" },
  },
});
