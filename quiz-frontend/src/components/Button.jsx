export default function Button({ children, variant = 'primary', ...props }) {
      const base = 'px-4 py-2 rounded-xl font-medium transition-shadow';
      const variants = {
            primary: 'bg-indigo-600 text-white hover:shadow-md',
            secondary: 'bg-white text-indigo-600 border border-indigo-600 hover:shadow-md',
            danger: 'bg-red-600 text-white hover:shadow-md',
      };

      return (
            <button className={`${base} ${variants[variant]}`} {...props}>
                  {children}
            </button>
      );
}