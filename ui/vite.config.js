import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'
import eslint from '@nabla/vite-plugin-eslint'
import vueI18n from '@intlify/unplugin-vue-i18n/vite'
import compression from 'vite-plugin-compression2'
import { fileURLToPath, URL } from 'node:url'

// https://vite.dev/config/
export default defineConfig({
  base: '/websh/',
  plugins: [
    vue(),
    eslint(),
    compression({
      algorithms: ['gzip'],
      deleteOriginalAssets: false,
      threshold: 10240
    }),
    vueI18n()
  ],
  resolve: {
    alias: {
      '@': fileURLToPath(new URL('./src', import.meta.url)),
    },
  },
  server: {
    proxy: {
      '/groups': {
        target: 'http://127.0.0.1:5913'
      },
      '/devs': {
        target: 'http://127.0.0.1:5913'
      },
      '/signin': {
        target: 'http://127.0.0.1:5913'
      },
      '/signout': {
        target: 'http://127.0.0.1:5913'
      },
      '/alive': {
        target: 'http://127.0.0.1:5913'
      },
      '^/cmd/.*': {
        target: 'http://127.0.0.1:5913'
      },
      '^/connect/.*': {
        ws: true,
        target: 'http://127.0.0.1:5913'
      },
      '^/web/*': {
        target: 'http://127.0.0.1:5913'
      }
    }
  }
})
