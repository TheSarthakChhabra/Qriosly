import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Navbar from './Navbar';
import Sidebar from './Sidebar';
import Button from './Button';

export default function AdminLayout({ children }) {
      const { user, logout } = useAuth();
      const location = useLocation();

      const links = [
            { to: '/admin/dashboard', label: 'Dashboard' },
            { to: '/admin/users', label: 'Users' },
      ];

      return (
            <div className="min-h-screen flex flex-col">
                  <Navbar><span className="text-sm text-white">Admin</span></Navbar>
                  <div className="flex flex-1">
                        <Sidebar>
                              {links.map((link) => (
                                    <Link
                                          key={link.to}
                                          to={link.to}
                                          className={`px-3 py-2 rounded-lg text-sm font-medium ${location.pathname === link.to
                                                      ? 'bg-slate-100 text-slate-900 font-semibold'
                                                      : 'text-slate-600 hover:bg-slate-50'
                                                }`}
                                    >
                                          {link.label}
                                    </Link>
                              ))}
                              <div className="mt-auto pt-4 border-t border-slate-200">
                                    <p className="text-sm font-medium text-slate-700 px-3">{user.id}</p>
                                    <Button variant="secondary" className="w-full mt-2" onClick={logout}>Logout</Button>
                              </div>
                        </Sidebar>
                        <main className="flex-1 p-8 bg-slate-50">{children}</main>
                  </div>
            </div>
      );
}