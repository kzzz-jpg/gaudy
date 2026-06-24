import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

// 開發時把 /api 反向代理到後端 :8080，避免跨網域問題。
export default defineConfig({
  plugins: [vue()],
  server: {
    port: 5173,
    proxy: {
      '/api': {
        target: 'http://localhost:8080',
        changeOrigin: true
      }
    }
  }
})