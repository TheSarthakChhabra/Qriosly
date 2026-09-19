import { useAuth } from '../context/AuthContext';
import Button from '../components/Button';

export default function StudentDashboard() {
      const { user, logout } = useAuth();

      return (
            <div className="min-h-screen bg-slate-50 p-8">
                  <h1 className="text-2xl font-semibold mb-2">Student Dashboard</h1>
                  <p className="text-slate-600 mb-4">Logged in as user ID: {user.id}</p>
                  <Button variant="secondary" onClick={logout}>Logout</Button>
            </div>
      );
}