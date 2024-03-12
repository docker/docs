module.exports = {
  plugins: {
    "postcss-import": {},
    "tailwindcss/nesting": {},
    tailwindcss: {},
    ...(process.env.NODE_ENV === "production" ? { autoprefixer: {} } : {}),
  },
};
