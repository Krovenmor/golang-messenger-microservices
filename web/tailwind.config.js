/** @type {import('tailwindcss').Config} */
export default {
  content: ['./index.html', './src/**/*.{js,ts,jsx,tsx}'],
  theme: {
    extend: {
      colors: {
        void: '#0B0E14',
        surface: '#12161F',
        raised: '#1A2029',
        line: '#242B38',
        ink: '#E8EAED',
        muted: '#8B93A7',
        faint: '#565F72',
        signal: '#5EEAD4',
        signalDim: '#1F3A38',
        warn: '#F2555A',
        warnDim: '#3A1F22',
      },
      fontFamily: {
        display: ['"Space Grotesk"', 'sans-serif'],
        body: ['Inter', 'sans-serif'],
        mono: ['"JetBrains Mono"', 'monospace'],
      },
      keyframes: {
        'pulse-signal': {
          '0%, 100%': { opacity: 1, transform: 'scale(1)' },
          '50%': { opacity: 0.4, transform: 'scale(0.85)' },
        },
        rise: {
          '0%': { opacity: 0, transform: 'translateY(6px)' },
          '100%': { opacity: 1, transform: 'translateY(0)' },
        },
      },
      animation: {
        'pulse-signal': 'pulse-signal 2s ease-in-out infinite',
        rise: 'rise 0.18s ease-out',
      },
    },
  },
  plugins: [],
}
