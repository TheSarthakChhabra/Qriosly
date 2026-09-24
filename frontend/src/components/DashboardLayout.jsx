import { Link, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Navbar from './Navbar';
import Sidebar from './Sidebar';
import Button from './Button';

export default function DashboardLayout({ children }) {
      const { user, logout } = useAuth();
      const location = useLocation();

      const links = [
            { to: '/student/dashboard', label: 'Dashboard' },
            { to: '/student/quizzes', label: 'Available Quizzes' },
            { to: '/student/attempts', label: 'My Attempts' },
            { to: '/student/results', label: 'Results' },
      ];

      return (
            <div className="min-h-screen flex flex-col">
                  <Navbar>
                        <span className="text-sm text-slate-600">Student</span>
                  </Navbar>

                  <div className="flex flex-1">
                        <Sidebar>
                              {links.map((link) => (
                                    <Link
                                          key={link.to}
                                          to={link.to}
                                          className={`px-3 py-2 rounded-lg text-sm font-medium ${location.pathname === link.to
                                                      ? 'bg-indigo-50 text-indigo-700'
                                                      : 'text-slate-600 hover:bg-slate-50'
                                                }`}
                                    >
                                          {link.label}
                                    </Link>
                              ))}

                              <div className="mt-auto pt-4 border-t border-slate-200">
                                    <p className="text-sm font-medium text-slate-700 px-3">{user.id}</p>
                                    <Button variant="secondary" className="w-full mt-2" onClick={logout}>
                                          Logout
                                    </Button>
                              </div>
                        </Sidebar>

                        <main className="flex-1 p-8 bg-slate-50">{children}</main>
                  </div>
            </div>
      );
}