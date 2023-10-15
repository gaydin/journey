import { fileURLToPath, URL } from 'node:url'

import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// https://vitejs.dev/config/
export default defineConfig({
  server: {
    proxy: {
      "/admin/v1": {
        target: "https://localhost:8085",
        changeOrigin: true,
        secure: false,
      },
      "/images": {
        target: "https://localhost:8085",
        changeOrigin: true,
        secure: false,
      },
    },
  },
  plugins: [
    vue(),
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url))
    }
  },
  base: "/admin"
})
