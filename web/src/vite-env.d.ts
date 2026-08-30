/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL?: string
  readonly VITE_WS_URL?: string
  readonly VITE_STATIC_URL?: string
}

interface ImportMeta {
  readonly env: ImportMetaEnv
}

/** Injected at build time from package.json — see vite.config.ts. */
declare const __APP_VERSION__: string
