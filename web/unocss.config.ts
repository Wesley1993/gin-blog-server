import { defineConfig, presetUno, presetAttributify } from 'unocss';

export default defineConfig({
  presets: [presetUno(), presetAttributify()],
  theme: {
    breakpoints: {
      sm: '640px',
      md: '768px',
      lg: '1024px',
      xl: '1280px',
      '2xl': '1536px',
      '4k': '2560px',
    },
    colors: {
      ink: '#1a1815',
      paper: '#f4efe6',
      parchment: '#faf7f0',
      amber: {
        500: '#c2701d',
        600: '#a85a10',
      },
    },
  },
  shortcuts: {
    'page-container': 'p-6 min-h-full w-full max-w-[1600px] mx-auto 4k:max-w-[1800px]',
    'toolbar-row': 'flex items-center gap-3 mb-4 flex-wrap',
    'card-panel': 'bg-white rounded-lg border border-solid border-[#e7e0d2]',
  },
});
