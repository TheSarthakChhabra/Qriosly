import { useAdminStats } from '../hooks/useAdminStats';
import AdminLayout from '../components/AdminLayout';
import Card from '../components/Card';

function StatCard({ label, value }) {
      return (
            <Card className="text-center">
                  <p className="text-sm text-slate-500 mb-1">{label}</p>
                  <p className="text-3xl font-semibold text-slate-900">{value}</p>
            </Card>
      );
}

export default function AdminDashboard() {
      const { stats, loading } = useAdminStats();

      return (
            <AdminLayout>
                  <h1 className="text-2xl font-semibold mb-6">Admin Dashboard</h1>
                  <div className="grid grid-cols-3 gap-4">
                        <StatCard label="Students" value={loading ? '...' : stats.students} />
                        <StatCard label="Teachers" value={loading ? '...' : stats.teachers} />
                        <StatCard label="Quizzes" value={loading ? '...' : stats.quizzes} />
                  </div>
            </AdminLayout>
      );
}