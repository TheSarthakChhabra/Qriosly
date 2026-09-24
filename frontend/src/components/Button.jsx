export default function Button({ children, variant = 'primary', ...props }) {
      const base = 'px-4 py-2 rounded-xl font-medium transition-colors';
      const variants = {
            primary: 'bg-indigo-700 text-white hover:bg-indigo-800',
            secondary: 'bg-white text-indigo-700 border border-indigo-700 hover:bg-indigo-50',
            accent: 'bg-amber-500 text-white hover:bg-amber-400',
            danger: 'bg-red-500 text-white hover:bg-red-600',
      };

      return (
            <button className={`${base} ${variants[variant]}`} {...props}>
                  {children}
            </button>
      );
}