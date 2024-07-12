/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.{html,js}"],
  theme: {
    extend: {
      width: {
        '200%': '200%',
      },
      aspectRatio: {
        '16/4': '16 / 4',
        '16/5': '16 / 5',
        '16/6': '16 / 6'
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["light"],
  },
}

