import QurioslyLogo from "../images/QurioslyLogo.png";

export default function Navbar({ children }) {
      return (
            <nav className="h-20 bg-slate-900 border-b border-slate-800 px-8 flex items-center justify-between">
                  <img
                        src={QurioslyLogo}
                        alt="Quriosly"
                        className="w-36 h-auto object-contain"
                  />

                  <div className="flex items-center gap-4 text-white">
                        {children}
                  </div>
            </nav>
      );
}