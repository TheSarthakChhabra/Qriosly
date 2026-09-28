import { useNavigate } from 'react-router-dom';
import { useTeacherQuizzes } from '../hooks/useTeacherQuizzes';
import TeacherLayout from '../components/TeacherLayout';
import Card from '../components/Card';
import Button from '../components/Button';
import { useState } from 'react';

function StatCard({ label, value }) {
  return (
    <Card className="text-center">
      <p className="text-sm text-slate-500 mb-1">{label}</p>
      <p className="text-3xl font-semibold text-slate-900">{value}</p>
    </Card>
  );
}

export default function TeacherDashboard() {
  const navigate = useNavigate();
  const { quizzes, totalAttempts, uniqueStudents, state } = useTeacherQuizzes();
  const [deleteNotice, setDeleteNotice] = useState('');

  function handleDelete(quizId) {
    setDeleteNotice(`Deleting quizzes isn't supported yet.`);
    setTimeout(() => setDeleteNotice(''), 3000);
  }

  return (
    <TeacherLayout>
      <h1 className="text-2xl font-semibold mb-6">Welcome back 👋</h1>

      <div className="grid grid-cols-3 gap-4 mb-8">
        <StatCard label="My Quizzes" value={state === 'loading' ? '...' : quizzes.length} />
        <StatCard label="Attempts" value={state === 'loading' ? '...' : totalAttempts} />
        <StatCard label="Students" value={state === 'loading' ? '...' : uniqueStudents} />
      </div>

      <div className="flex items-center justify-between mb-4">
        <h2 className="text-lg font-semibold">My Quizzes</h2>
        <Button variant="accent" onClick={() => navigate('/teacher/quizzes/create')}>
          + Create Quiz
        </Button>
      </div>

      {deleteNotice && <p className="text-sm text-amber-600 mb-4">{deleteNotice}</p>}

      {state === 'loading' ? (
        <p className="text-slate-500">Loading your quizzes...</p>
      ) : state === 'error' ? (
        <p className="text-red-500">Failed to load your quizzes.</p>
      ) : quizzes.length === 0 ? (
        <Card className="max-w-md text-center">
          <p className="text-slate-500 mb-4">You haven't created any quizzes yet.</p>
          <Button variant="accent" onClick={() => navigate('/teacher/quizzes/create')}>
            + Create Quiz
          </Button>
        </Card>
      ) : (
        <div className="grid grid-cols-3 gap-4">
          {quizzes.map((quiz) => (
            <Card key={quiz.id}>
              <h3 className="font-semibold text-lg mb-1">{quiz.title}</h3>
              <p className="text-sm text-slate-500 mb-4">
                {quiz.questionCount} Questions · {quiz.duration_minutes} min
              </p>
              <div className="flex gap-2 flex-wrap">
                <Button variant="secondary" onClick={() => navigate(`/teacher/quizzes/${quiz.id}/edit`)}>
                  Edit
                </Button>
                <Button variant="secondary" onClick={() => navigate(`/teacher/quizzes/${quiz.id}/results`)}>
                  Results
                </Button>
                <Button variant="danger" onClick={() => handleDelete(quiz.id)}>
                  Delete
                </Button>
              </div>
            </Card>
          ))}
        </div>
      )}
    </TeacherLayout>
  );
}