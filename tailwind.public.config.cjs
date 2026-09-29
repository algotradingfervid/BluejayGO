module.exports = {
  content: ['./templates/public/**/*.html', './templates/partials/**/*.html', './public/js/*.js'],
  theme: {
    extend: {
      colors: {
        primary: '#0066CC', 'primary-dark': '#004499', 'navy-tech': '#004499',
        'background-light': '#F8F9FA', 'background-dark': '#1A1A2E',
        'text-primary': '#333333', 'text-secondary': '#666666',
      },
      fontFamily: {display: ['Inter', 'sans-serif'], mono: ['JetBrains Mono', 'monospace']},
      borderRadius: {DEFAULT: '0px', lg: '0px', xl: '0px', full: '9999px'},
    },
  },
  plugins: [require('@tailwindcss/forms'), require('@tailwindcss/container-queries')],
};
