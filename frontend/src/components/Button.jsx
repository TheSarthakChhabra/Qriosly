export default function Button({ children, variant = 'primary', ...props }) {
      const base = 'px-4 py-2 rounded-xl font-medium transition-colors';
      const variants = {
            primary: 'bg-slate-900 text-white hover:bg-slate-800',
            secondary: 'bg-white text-slate-900 border border-slate-900 hover:bg-slate-50',
            accent: 'bg-amber-500 text-white hover:bg-amber-400',
            danger: 'bg-red-500 text-white hover:bg-red-600',
      };

      return (
            <button className={`${base} ${variants[variant]}`} {...props}>
                  {children}
            </button>
      );
}