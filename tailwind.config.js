/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./**/*.{html,js,go}"],
  darkMode: "selector",
  theme: {
    extend: {
      screens: {
        sm: "640px", // Small devices (mobile)
        md: "768px", // Medium devices (tablets)
        lg: "1024px", // Large devices (desktops)
        xl: "1280px", // Extra large devices (large desktops)
        "2xl": "1536px", // Extra extra large devices
        "3xl": "1920px", // Very large screens
        "4xl": "2560px", // Ultra wide screens
      },
      colors: {
        dracula: {
          background: "#282a36",
          current: "#44475a",
          foreground: "#f8f8f2",
          comment: "#6272a4",
          cyan: "#8be9fd",
          green: "#50fa7b",
          orange: "#ffb86c",
          pink: "#ff79c6",
          purple: "#bd93f9",
          red: "#ff5555",
          yellow: "#f1fa8c",
        },
      },
    },
  },
  plugins: [],
};
