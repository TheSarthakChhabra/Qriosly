import { useParams, useNavigate } from 'react-router-dom';
import { useState } from 'react';
import { useQuizDetails } from '../hooks/useQuizDetails';
import { api } from '../api/client';
import Card from '../components/Card';
import Button from '../components/Button';

export default function QuizDetails() {
      const { id } = useParams();
      const navigate = useNavigate();
      const { quiz, questionCount, loading, error } = useQuizDetails(id);
      const [starting, setStarting] = useState(false);
      const [startError, setStartError] = useState('');

      async function handleStart() {
            setStarting(true);
            setStartError('');
            try {
                  const attempt = await api.post(`/quizzes/${id}/attempts`, {});
                  navigate(`/attempt/${attempt.id}`);
            } catch (err) {
                  setStartError(err.message);
                  setStarting(false);
            }
      }

      if (loading) {
            return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
      }

      if (error) {
            return (
                  <div className="min-h-screen flex items-center justify-center">
                        <p className="text-red-500">{error}</p>
                  </div>
            );
      }

      return (
            <div className="min-h-screen bg-slate-50 p-8">
                  <button
                        onClick={() => navigate('/student/dashboard')}
                        className="text-sm text-slate-500 hover:text-slate-700 mb-6"
                  >
                        ← Back to Quizzes
                  </button>

                  <Card className="max-w-2xl mx-auto">
                        <h1 className="text-2xl font-semibold text-center mb-6">{quiz.title}</h1>

                        <hr className="border-slate-200 mb-6" />

                        <h2 className="font-semibold mb-2">Description</h2>
                        <p className="text-slate-600 mb-6">{quiz.description || 'No description provided.'}</p>

                        <div className="grid grid-cols-2 gap-4 mb-6">
                              <Card className="text-center">
                                    <p className="text-sm text-slate-500 mb-1">Questions</p>
                                    <p className="text-2xl font-semibold text-slate-900">{questionCount}</p>
                              </Card>
                              <Card className="text-center">
                                    <p className="text-sm text-slate-500 mb-1">Duration</p>
                                    <p className="text-2xl font-semibold text-slate-900">{quiz.duration_minutes} min</p>
                              </Card>
                        </div>

                        <hr className="border-slate-200 mb-6" />

                        <h2 className="font-semibold mb-2">Instructions</h2>
                        <ul className="list-disc list-inside text-slate-600 mb-6 space-y-1">
                              <li>You have {quiz.duration_minutes} minutes to complete this quiz.</li>
                              <li>Once submitted, the attempt cannot be changed.</li>
                              <li>Make sure you have a stable connection.</li>
                        </ul>

                        {startError && <p className="text-sm text-red-600 mb-4">{startError}</p>}

                        <div className="flex justify-end">
                              <Button variant="accent" onClick={handleStart} disabled={starting}>
                                    {starting ? 'Starting...' : 'Start Quiz'}
                              </Button>
                        </div>
                  </Card>
            </div>
      );
}