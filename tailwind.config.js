module.exports = {
  content: ["./internal/**/*.{go,js,templ,html}"],
  theme: {
    extend: {}
  },
  plugins: [require("daisyui")],
  daisyui: {
    themes: ["synthwave --default", "dark"]
  }
};
