/** @type {import('tailwindcss').Config} */
export default {
  content: ['./src/**/*.{html,js,svelte,ts}'],
  theme: {
    extend: {
      colors: {
        'news-dark': '#0f1729',
        'news-gray': '#1e283b',
        'news-light': '#38bdf8',
        'news-separator': '#172f47',
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        dark: {
          'primary': '#38bdf8',
          'base-content': '#38bdf8',
        },
      }
    ],
  },
}

