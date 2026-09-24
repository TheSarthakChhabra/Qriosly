import { useParams, useNavigate } from 'react-router-dom';
import { useAttemptResult } from '../hooks/useAttemptResult';
import Card from '../components/Card';
import Button from '../components/Button';

function formatDate(isoString) {
      return new Date(isoString).toLocaleDateString('en-US', {
            day: 'numeric',
            month: 'short',
            year: 'numeric',
      });
}

function formatDuration(startedAt, submittedAt) {
      const ms = new Date(submittedAt) - new Date(startedAt);
      const minutes = Math.round(ms / 60000);
      return `${minutes} min`;
}

export default function ResultPage() {
      const { attemptId } = useParams();
      const navigate = useNavigate();
      const { attempt, quiz, totalQuestions, state } = useAttemptResult(attemptId);

      if (state === 'loading') {
            return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
      }

      if (state === 'not_found') {
            return (
                  <div className="min-h-screen flex items-center justify-center p-4">
                        <Card className="max-w-md text-center">
                              <h1 className="text-lg font-semibold mb-2">Attempt Not Found</h1>
                              <p className="text-slate-500 mb-4">This attempt doesn't exist or has been removed.</p>
                              <Button onClick={() => navigate('/student/dashboard')}>Back to Dashboard</Button>
                        </Card>
                  </div>
            );
      }

      if (state === 'forbidden') {
            return (
                  <div className="min-h-screen flex items-center justify-center p-4">
                        <Card className="max-w-md text-center">
                              <h1 className="text-lg font-semibold mb-2">Access Denied</h1>
                              <p className="text-slate-500 mb-4">You don't have permission to view this attempt.</p>
                              <Button onClick={() => navigate('/student/dashboard')}>Back to Dashboard</Button>
                        </Card>
                  </div>
            );
      }

      if (state === 'error') {
            return (
                  <div className="min-h-screen flex items-center justify-center p-4">
                        <Card className="max-w-md text-center">
                              <h1 className="text-lg font-semibold mb-2">Something Went Wrong</h1>
                              <p className="text-slate-500 mb-4">We couldn't load this result. Please try again.</p>
                              <Button onClick={() => navigate('/student/dashboard')}>Back to Dashboard</Button>
                        </Card>
                  </div>
            );
      }

      // state === 'ready' from here on — attempt/quiz/totalQuestions are guaranteed present.
      const correct = attempt.score;
      const incorrect = totalQuestions - correct;
      const percentage = totalQuestions > 0 ? Math.round((correct / totalQuestions) * 100) : 0;

      return (
            <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
                  <Card className="max-w-lg w-full text-center">
                        <p className="text-emerald-500 font-semibold mb-1">✓ Quiz Submitted</p>
                        <h1 className="text-xl font-semibold mb-6">{quiz.title}</h1>

                        <div className="bg-indigo-50 rounded-xl py-8 mb-6">
                              <p className="text-4xl font-bold text-indigo-700">
                                    {correct} / {totalQuestions}
                              </p>
                              <p className="text-lg text-indigo-600 mt-1">{percentage}%</p>
                        </div>

                        <p className="text-sm text-slate-500 mb-1">
                              Completed: {formatDate(attempt.submitted_at)}
                        </p>
                        <p className="text-sm text-slate-500 mb-6">
                              Time taken: {formatDuration(attempt.started_at, attempt.submitted_at)}
                        </p>

                        <hr className="border-slate-200 mb-6" />

                        <p className="text-sm font-medium text-slate-600 mb-3">Questions</p>
                        <div className="flex justify-center gap-6 mb-8">
                              <div>
                                    <p className="text-lg font-semibold">{totalQuestions}</p>
                                    <p className="text-xs text-slate-500">Total</p>
                              </div>
                              <div>
                                    <p className="text-lg font-semibold text-emerald-500">{correct}</p>
                                    <p className="text-xs text-slate-500">Correct</p>
                              </div>
                              <div>
                                    <p className="text-lg font-semibold text-red-500">{incorrect}</p>
                                    <p className="text-xs text-slate-500">Incorrect</p>
                              </div>
                        </div>

                        <div className="flex flex-col gap-2">
                              <Button onClick={() => navigate('/student/attempts')}>View My Attempts</Button>
                              <Button variant="secondary" onClick={() => navigate('/student/dashboard')}>
                                    Back to Dashboard
                              </Button>
                        </div>
                  </Card>
            </div>
      );
}