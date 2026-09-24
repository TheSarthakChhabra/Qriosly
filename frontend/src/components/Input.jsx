export default function Input({ label, ...props }) {
      return (
            <div className="flex flex-col gap-1">
                  {label && <label className="text-sm font-medium text-slate-700">{label}</label>}
                  <input
                        className="px-3 py-2 rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-slate-900"
                        {...props}
                  />
            </div>
      );
}