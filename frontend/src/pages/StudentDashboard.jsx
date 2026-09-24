import { useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useDashboardStats } from '../hooks/useDashboardStats';
import { useQuizzes } from '../hooks/useQuizzes';
import DashboardLayout from '../components/DashboardLayout';
import Card from '../components/Card';
import Button from '../components/Button';

function StatCard({ label, value }) {
  return (
    <Card className="text-center">
      <p className="text-sm text-slate-900 mb-1">{label}</p>
      <p className="text-3xl font-semibold text-indigo-700">{value}</p>
    </Card>
  );
}

export default function StudentDashboard() {
  const { user } = useAuth();
  const { stats, loading: statsLoading } = useDashboardStats();
  const { quizzes, loading: quizzesLoading } = useQuizzes();
  const navigate = useNavigate();

  return (
    <DashboardLayout>
      <h1 className="text-2xl font-semibold mb-6">Welcome back 👋</h1>

      <div className="grid grid-cols-3 gap-4 mb-8">
        <StatCard label="Available Quizzes" value={statsLoading ? '...' : stats.availableQuizzes} />
        <StatCard label="Completed Attempts" value={statsLoading ? '...' : stats.completedAttempts} />
        <StatCard
          label="Average Score"
          value={statsLoading ? '...' : (stats.averageScore ?? '—')}
        />
      </div>

      <h2 className="text-lg font-semibold mb-4">Available Quizzes</h2>

      {quizzesLoading ? (
        <p className="text-slate-500">Loading quizzes...</p>
      ) : quizzes.length === 0 ? (
        <p className="text-slate-500">No quizzes available yet.</p>
      ) : (
        <div className="grid grid-cols-3 gap-4">
          {quizzes.map((quiz) => (
            <Card key={quiz.id} className="flex flex-col justify-between">
              <div>
                <h3 className="font-semibold text-lg mb-2">{quiz.title}</h3>
                <p className="text-sm text-slate-500">{quiz.questionCount} Questions</p>
                <p className="text-sm text-slate-500">Duration: {quiz.duration_minutes} min</p>
              </div>
              <Button className="mt-4" onClick={() => navigate(`/quizzes/${quiz.id}`)}>
                View Quiz
              </Button>
            </Card>
          ))}
        </div>
      )}
    </DashboardLayout>
  );
}