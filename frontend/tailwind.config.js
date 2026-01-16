/** @type {import('tailwindcss').Config} */
export default {
  content: [
    "./index.html",
    "./src/**/*.{js,ts,jsx,tsx}",
  ],
  theme: {
    extend: {
      colors: {
        background: "#0a0a0a",
        foreground: "#ededed",
        card: "#121212",
        "card-foreground": "#ffffff",
        primary: "#3b82f6",
        secondary: "#1f2937",
        accent: "#10b981",
        muted: "#737373",
      },
    },
  },
  plugins: [],
}
