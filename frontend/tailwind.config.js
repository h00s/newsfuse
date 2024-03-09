/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        'news-dark': '#0f1729',
        'news-gray': '#1e283b',
        'news-light': '#38bdf8',
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      'light',
      {
        'dark': {
          'base-content': '#38bdf8',
        },
      }
    ],
  },
}

