/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./hugo_stats.json", "layouts/**/*.html", "assets/**/*.js"],
  darkMode: "class",
  theme: {
    extend: {
      typography: (theme) => ({
        DEFAULT: {
          css: {
            pre: false,
            code: false,
            'pre code': false,
            'code::before': false,
            'code::after': false,
            // light colors for prose
            "--tw-prose-body": theme("colors.black"),
            "--tw-prose-headings": theme("colors.black"),
            "--tw-prose-lead": theme("colors.gray.light.600"),
            "--tw-prose-links": theme("colors.blue.light.500"),
            "--tw-prose-bold": theme("colors.black"),
            "--tw-prose-counters": theme("colors.black"),
            "--tw-prose-bullets": theme("colors.black"),
            "--tw-prose-hr": theme("colors.divider.light"),
            "--tw-prose-quotes": theme("colors.black"),
            "--tw-prose-quote-borders": theme("colors.blue.light.500"),
            "--tw-prose-captions": theme("colors.gray.light.600"),
            "--tw-prose-th-borders": theme("colors.gray.light.200"),
            "--tw-prose-td-borders": theme("colors.gray.light.200"),

            // dark colors for prose
            "--tw-prose-invert-body": theme("colors.white"),
            "--tw-prose-invert-headings": theme("colors.white"),
            "--tw-prose-invert-lead": theme("colors.gray.dark.600"),
            "--tw-prose-invert-links": theme("colors.blue.dark.500"),
            "--tw-prose-invert-bold": theme("colors.white"),
            "--tw-prose-invert-counters": theme("colors.white"),
            "--tw-prose-invert-bullets": theme("colors.white"),
            "--tw-prose-invert-hr": theme("colors.divider.dark"),
            "--tw-prose-invert-quotes": theme("colors.white"),
            "--tw-prose-invert-quote-borders": theme("colors.blue.dark.500"),
            "--tw-prose-invert-captions": theme("colors.gray.dark.600"),
            "--tw-prose-invert-th-borders": theme("colors.gray.dark.200"),
            "--tw-prose-invert-td-borders": theme("colors.gray.dark.200"),
          },
        },
      }),

      gridTemplateColumns: {
        'main-xl': 'minmax(300px, 1fr) minmax(100ch, 1fr) 1fr',
        'main-lg': '300px minmax(75ch, 2fr) 1fr',
        'main-md': '250px 1fr',
      },
    },

    // theme values
    fontSize: {
      xs: ["0.7143rem", { letterSpacing: "0.015em", fontWeight: 500 }],
      sm: "0.851rem",
      base: ["14px", {}],
      lg: ["1.1429rem", 1.75],
      xl: ["1.2857rem", { letterSpacing: "-0.015em", fontWeight: 500 }],
      "2xl": ["1.5rem", { letterSpacing: "-0.015em", fontWeight: 500 }],
      "3xl": ["2rem", { fontWeight: 500 }],
      "4xl": ["2.5rem", { letterSpacing: "-0.015em", fontWeight: 500 }],
    },

    colors: {
      white: "#fff",
      black: "#000",
      transparent: 'transparent',

      accent: {
        light: "#677285",
        dark: "#2D404E",
      },

      background: {
        light: "#f9f9fa",
        dark: "#141b1f",
      },

      divider: {
        light: "hsla(0, 0%, 0%, 0.1)",
        dark: "hsla(0, 0%, 100%, 0.05)",
      },

      amber: {
        light: {
          DEFAULT: "#b85504",
          100: "#fff4dc",
          200: "#fce1a9",
          300: "#fbb552",
          400: "#dd7805",
          500: "#b85504",
          600: "#a9470f",
          700: "#893607",
          800: "#421a02",
        },
        dark: {
          DEFAULT: "#ed8d25",
          100: "#753715",
          200: "#944307",
          300: "#af560a",
          400: "#cd6a0a",
          500: "#ed8d25",
          600: "#f6a650",
          700: "#f8b974",
          800: "#fac892",
        },
      },

      red: {
        light: {
          DEFAULT: "#d52536",
          100: "#fdeaea",
          200: "#f6cfd0",
          300: "#eea3a5",
          400: "#e25d68",
          500: "#d52536",
          600: "#b72132",
          700: "#8b1924",
          800: "#350a10",
        },
        dark: {
          DEFAULT: "#dd4659",
          100: "#741624",
          200: "#951c2f",
          300: "#bc233c",
          400: "#d1334c",
          500: "#dd4659",
          600: "#e96c7c",
          700: "#ea8e9a",
          800: "#f0aab4",
        },
      },

      violet: {
        light: {
          DEFAULT: "#7d2eff",
          100: "#f7ecff",
          200: "#e9d4ff",
          300: "#c9a6ff",
          400: "#9860ff",
          500: "#7d2eff",
          600: "#6d00eb",
          700: "#5700bb",
          800: "#220041",
        },
        dark: {
          DEFAULT: "#a371fc",
          100: "#380093",
          200: "#4F00B4",
          300: "#6D1CDB",
          400: "#8a53ec",
          500: "#a371fc",
          600: "#b38bfc",
          700: "#c5a6fd",
          800: "#d4bdfe",
        },
      },

      magenta: {
        light: {
          DEFAULT: "#C918C0",
          100: "#FFE6FB",
          200: "#FFC9F6",
          300: "#FFA6F0",
          400: "#E950E2",
          500: "#C918C0",
          600: "#AB00A4",
          700: "#830080",
          800: "#440040",
        },
        dark: {
          DEFAULT: "#E950E2",
          100: "#7E0078",
          200: "#92008B",
          300: "#AB00A4",
          400: "#CC18C4",
          500: "#E950E2",
          600: "#FF6FF9",
          700: "#FF8AFA",
          800: "#FFA4FB",
        },
      },


      blue: {
        light: {
          DEFAULT: "#086dd7",
          100: "#e5f2fc",
          200: "#c0e0fa",
          300: "#8bc7f5",
          400: "#1c90ed",
          500: "#1D63ED",
          600: "#0C49C2",
          700: "#00308D",
          800: "#00084d",
        },
        dark: {
          DEFAULT: "#3391ee",
          100: "#002EA3",
          200: "#063BB7",
          300: "#1351D4",
          400: "#1D63ED",
          500: "#3391ee",
          600: "#55a4f1",
          700: "#7cb9f4",
          800: "#98c8f6",
        },
      },

      green: {
        light: {
          DEFAULT: "#2e7f74",
          100: "#e6f5f3",
          200: "#c6eae1",
          300: "#88d5c0",
          400: "#3ba08d",
          500: "#2e7f74",
          600: "#1e6c5f",
          700: "#185a51",
          800: "#0c2c28",
        },
        dark: {
          DEFAULT: "#2aa391",
          100: "#003F36",
          200: "#005045",
          300: "#006256",
          400: "#008471",
          500: "#00A58C",
          600: "#3cc1ad",
          700: "#7accc3",
          800: "#a5ddd6",
        },
      },

      gray: {
        light: {
          DEFAULT: "#677285",
          100: "#F4F4F6",
          200: "#e1e2e6",
          300: "#c4c8d1",
          400: "#8993a5",
          500: "#677285",
          600: "#505968",
          700: "#393f49",
          800: "#17191e",
        },
        dark: {
          DEFAULT: "#7794ab",
          100: "#080B0E",
          200: "#1C262D",
          300: "#2D404E",
          400: "#4E6A81",
          500: "#7794ab",
          600: "#94abbc",
          700: "#adbecb",
          800: "#c4d0da",
        },
      },
    },

    fontFamily: {
      sans: [
        "Roboto Flex",
        "system-ui",
        "-apple-system",
        "BlinkMacSystemFont",
        "Segoe UI",
        "Oxygen",
        "Ubuntu",
        "Cantarell",
        "Open Sans",
        "Helvetica Neue",
        "sans-serif",
      ],
      mono: ["Roboto Mono", "monospace"],
      icons: ["Material Symbols Rounded"],
    },
  },
  plugins: [require("@tailwindcss/typography")],
};
