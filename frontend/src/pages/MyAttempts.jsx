import { useNavigate } from 'react-router-dom';
import { useMyAttempts } from '../hooks/useMyAttempts';
import DashboardLayout from '../components/DashboardLayout';
import Card from '../components/Card';
import Button from '../components/Button';

function formatDate(isoString) {
      return new Date(isoString).toLocaleDateString('en-US', { day: 'numeric', month: 'short' });
}

export default function MyAttempts() {
      const navigate = useNavigate();
      const { attempts, quizzesById, state } = useMyAttempts();

      if (state === 'loading') {
            return (
                  <DashboardLayout>
                        <p className="text-slate-500">Loading your attempts...</p>
                  </DashboardLayout>
            );
      }

      if (state === 'error') {
            return (
                  <DashboardLayout>
                        <Card className="max-w-md">
                              <p className="text-red-500">Failed to load your attempts. Please try again.</p>
                        </Card>
                  </DashboardLayout>
            );
      }

      if (attempts.length === 0) {
            return (
                  <DashboardLayout>
                        <Card className="max-w-md text-center">
                              <p className="text-slate-500 mb-4">You haven't taken any quizzes yet.</p>
                              <Button onClick={() => navigate('/student/dashboard')}>Browse Quizzes</Button>
                        </Card>
                  </DashboardLayout>
            );
      }

      return (
            <DashboardLayout>
                  <h1 className="text-xl font-semibold mb-6">My Attempts</h1>

                  <Card className="p-0 overflow-hidden">
                        <table className="w-full text-left text-sm">
                              <thead className="bg-slate-50 text-slate-500 border-b border-slate-200">
                                    <tr>
                                          <th className="px-4 py-3 font-medium">Quiz</th>
                                          <th className="px-4 py-3 font-medium">Date</th>
                                          <th className="px-4 py-3 font-medium">Score</th>
                                          <th className="px-4 py-3 font-medium">Status</th>
                                          <th className="px-4 py-3 font-medium"></th>
                                    </tr>
                              </thead>
                              <tbody>
                                    {attempts.map((attempt) => {
                                          const quiz = quizzesById[attempt.quiz_id];
                                          const isSubmitted = attempt.status === 'submitted';

                                          return (
                                                <tr key={attempt.id} className="border-b border-slate-100 last:border-0">
                                                      <td className="px-4 py-3">{quiz ? quiz.title : 'Unknown Quiz'}</td>
                                                      <td className="px-4 py-3 text-slate-500">{formatDate(attempt.started_at)}</td>
                                                      <td className="px-4 py-3">
                                                            {isSubmitted && attempt.score !== null && quiz
                                                                  ? `${attempt.score}`
                                                                  : '—'}
                                                      </td>
                                                      <td className="px-4 py-3">
                                                            <span
                                                                  className={`px-2 py-1 rounded-full text-xs font-medium ${isSubmitted
                                                                              ? 'bg-emerald-50 text-emerald-600'
                                                                              : 'bg-amber-50 text-amber-600'
                                                                        }`}
                                                            >
                                                                  {isSubmitted ? 'Done' : 'In Progress'}
                                                            </span>
                                                      </td>
                                                      <td className="px-4 py-3 text-right">
                                                            {isSubmitted ? (
                                                                  <Button variant="secondary" onClick={() => navigate(`/result/${attempt.id}`)}>
                                                                        View Result
                                                                  </Button>
                                                            ) : (
                                                                  <Button variant="secondary" onClick={() => navigate(`/attempt/${attempt.id}`)}>
                                                                        Resume
                                                                  </Button>
                                                            )}
                                                      </td>
                                                </tr>
                                          );
                                    })}
                              </tbody>
                        </table>
                  </Card>
            </DashboardLayout>
      );
}