/** @type {import('tailwindcss').Config} */
module.exports = {
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
            blockquote: false,
            'blockquote p:first-of-type::before': false,
            'blockquote p:last-of-type::after': false,
            // light colors for prose
            "--tw-prose-body": theme("colors.black"),
            "--tw-prose-headings": theme("colors.black"),
            "--tw-prose-lead": theme("colors.gray.light.600"),
            "--tw-prose-links": theme("colors.blue.light.500"),
            "--tw-prose-bold": theme("colors.black"),
            "--tw-prose-counters": theme("colors.black"),
            "--tw-prose-bullets": theme("colors.black"),
            "--tw-prose-hr": theme("colors.divider.light"),
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
            "--tw-prose-invert-captions": theme("colors.gray.dark.600"),
            "--tw-prose-invert-th-borders": theme("colors.gray.dark.200"),
            "--tw-prose-invert-td-borders": theme("colors.gray.dark.300"),
          },
        },
      }),
    },
  },
};
