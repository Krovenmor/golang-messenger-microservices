import { defineConfig, type Plugin } from 'vite'
import react from '@vitejs/plugin-react'

const WOFF_FALLBACK = /,\s*url\([^)]*\.woff\)\s*format\('woff'\)/g

function dropLegacyWoffFallback(): Plugin {
  return {
    name: 'drop-legacy-woff-fallback',
    enforce: 'pre',
    transform(code, id) {
      if (!id.includes('@fontsource') || !id.split('?')[0].endsWith('.css')) return null
      return code.replace(WOFF_FALLBACK, '')
    },
  }
}

export default defineConfig({
  plugins: [react(), dropLegacyWoffFallback()],
  server: {
    host: true,
    port: 5173,
  },
  preview: {
    host: true,
    port: 4173,
  },
})
