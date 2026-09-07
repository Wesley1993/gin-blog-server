/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_BASE_URL: string
  readonly VITE_BLOG_URL: string
}
interface ImportMeta {
  readonly env: ImportMetaEnv
}

declare module 'virtual:uno.css' {
  const css: string;
  export default css;
}
