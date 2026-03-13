/** @type {import('tailwindcss').Config} */
module.exports = {
    content: [
        './src/**/*.{js,ts,jsx,tsx,mdx}',
    ],
    theme: {
        extend: {
            colors: {
                cf: {
                    bg: '#0a0a0f',
                    surface: '#12121a',
                    border: '#1e1e2e',
                    primary: '#6366f1',
                    accent: '#22d3ee',
                    success: '#22c55e',
                    warning: '#f59e0b',
                    error: '#ef4444',
                    text: '#e2e8f0',
                    muted: '#64748b',
                },
            },
        },
    },
    plugins: [],
};
