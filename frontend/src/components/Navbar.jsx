export default function Navbar({ children }) {
      return (
            <nav className="bg-white border-b border-slate-200 px-6 py-4 flex items-center justify-between">
                  <span className="text-lg font-semibold text-indigo-700">Quiz Platform</span>
                  <div className="flex items-center gap-4">{children}</div>
            </nav>
      );
}