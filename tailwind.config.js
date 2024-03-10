const plugin = require('tailwindcss/plugin')

// Let's create a plugin that adds utilities!
const capitalizeFirst = plugin(function ({ addUtilities }) {
  const newUtilities = {
    '.capitalize-first:first-letter': {
      textTransform: 'uppercase',
    },
  }
  addUtilities(newUtilities, ['responsive', 'hover'])
})

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./templates/**/*.templ"],
  theme: {
    extend: {},
  },
  plugins: [capitalizeFirst, require("daisyui")],
  daisyui: {
    themes: ["dracula"],
  },
}

