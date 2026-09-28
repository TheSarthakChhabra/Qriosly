import { useParams, useNavigate } from 'react-router-dom';
import { useQuizResults } from '../hooks/useQuizResults';
import TeacherLayout from '../components/TeacherLayout';
import Card from '../components/Card';
import Button from '../components/Button';

function formatDateTime(isoString) {
      return new Date(isoString).toLocaleString('en-US', {
            day: 'numeric', month: 'short', hour: '2-digit', minute: '2-digit',
      });
}

export default function QuizResults() {
      const { id } = useParams();
      const navigate = useNavigate();
      const { quiz, summaries, state } = useQuizResults(id);

      if (state === 'loading') {
            return <TeacherLayout><p className="text-slate-500">Loading results...</p></TeacherLayout>;
      }

      if (state === 'not_found') {
            return (
                  <TeacherLayout>
                        <Card className="max-w-md">
                              <p className="text-slate-500 mb-4">This quiz doesn't exist.</p>
                              <Button onClick={() => navigate('/teacher/dashboard')}>Back to Dashboard</Button>
                        </Card>
                  </TeacherLayout>
            );
      }

      if (state === 'forbidden') {
            return (
                  <TeacherLayout>
                        <Card className="max-w-md">
                              <p className="text-slate-500 mb-4">You don't have permission to view this quiz's results.</p>
                              <Button onClick={() => navigate('/teacher/dashboard')}>Back to Dashboard</Button>
                        </Card>
                  </TeacherLayout>
            );
      }

      if (state === 'error') {
            return (
                  <TeacherLayout>
                        <Card className="max-w-md">
                              <p className="text-red-500">Failed to load results. Please try again.</p>
                        </Card>
                  </TeacherLayout>
            );
      }

      return (
            <TeacherLayout>
                  <h1 className="text-xl font-semibold mb-6">{quiz.title} — Results</h1>

                  {summaries.length === 0 ? (
                        <Card className="max-w-md text-center">
                              <p className="text-slate-500 mb-1">No attempts yet.</p>
                              <p className="text-sm text-slate-400">
                                    Students' results will appear here once they start submitting this quiz.
                              </p>
                        </Card>
                  ) : (
                        <Card className="p-0 overflow-hidden">
                              <table className="w-full text-left text-sm">
                                    <thead className="bg-slate-50 text-slate-500 border-b border-slate-200">
                                          <tr>
                                                <th className="px-4 py-3 font-medium">Student</th>
                                                <th className="px-4 py-3 font-medium">Submitted</th>
                                                <th className="px-4 py-3 font-medium">Score</th>
                                                <th className="px-4 py-3 font-medium">Status</th>
                                                <th className="px-4 py-3 font-medium"></th>
                                          </tr>
                                    </thead>
                                    <tbody>
                                          {summaries.map((s) => {
                                                const isSubmitted = s.status === 'submitted';
                                                return (
                                                      <tr key={s.attempt_id} className="border-b border-slate-100 last:border-0">
                                                            <td className="px-4 py-3">{s.user_name}</td>
                                                            <td className="px-4 py-3 text-slate-500">
                                                                  {isSubmitted ? formatDateTime(s.submitted_at) : '—'}
                                                            </td>
                                                            <td className="px-4 py-3">
                                                                  {isSubmitted && s.score !== null ? s.score : '—'}
                                                            </td>
                                                            <td className="px-4 py-3">
                                                                  <span
                                                                        className={`px-2 py-1 rounded-full text-xs font-medium ${isSubmitted ? 'bg-emerald-50 text-emerald-600' : 'bg-amber-50 text-amber-600'
                                                                              }`}
                                                                  >
                                                                        {isSubmitted ? 'Done' : 'In Progress'}
                                                                  </span>
                                                            </td>
                                                            <td className="px-4 py-3 text-right">
                                                                  <Button
                                                                        variant="secondary"
                                                                        onClick={() => navigate(`/teacher/attempts/${s.attempt_id}`)}
                                                                  >
                                                                        View
                                                                  </Button>
                                                            </td>
                                                      </tr>
                                                );
                                          })}
                                    </tbody>
                              </table>
                        </Card>
                  )}
            </TeacherLayout>
      );
}