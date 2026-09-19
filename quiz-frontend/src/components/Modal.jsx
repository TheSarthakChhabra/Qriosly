export default function Modal({ isOpen, onClose, children }) {
      if (!isOpen) return null;

      return (
            <div className="fixed inset-0 bg-black/40 flex items-center justify-center z-50">
                  <div className="bg-white rounded-xl shadow-lg p-6 max-w-md w-full mx-4">
                        {children}
                        <button
                              onClick={onClose}
                              className="mt-4 text-sm text-slate-500 hover:text-slate-700"
                        >
                              Close
                        </button>
                  </div>
            </div>
      );
}