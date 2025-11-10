/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/*.{html,js}"],
  theme: {
    extend: {
      width: {
        '200%': '200%',
        '98%' : '98%'
      },
      aspectRatio: {
        '16/4': '16 / 4',
        '16/5': '16 / 5',
        '16/6': '16 / 6',
        '16/7': '16 / 7',
        '16/8': '16 / 8',
      },
    },
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: [
      {
        myTheme: {
          "primary": "hsl(43 96% 56% / 1)",
          "secondary": "#1a73e8",
          "accent": "#37cdbe",
          "neutral": "#3d4451",
          "base-100": "#ffffff",
        },
      }
    ]
  },
}

