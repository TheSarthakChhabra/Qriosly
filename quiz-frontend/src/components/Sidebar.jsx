export default function Sidebar({ children }) {
      return (
            <aside className="w-56 bg-white border-r border-slate-200 h-screen p-4 flex flex-col gap-2">
                  {children}
            </aside>
      );
}