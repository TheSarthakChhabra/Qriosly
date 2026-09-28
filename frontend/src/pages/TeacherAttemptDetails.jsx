import { useParams, useNavigate } from 'react-router-dom';
import { useState, useEffect } from 'react';
import { api } from '../api/client';
import TeacherLayout from '../components/TeacherLayout';
import Card from '../components/Card';
import Button from '../components/Button';

export default function TeacherAttemptDetails() {
      const { attemptId } = useParams();
      const navigate = useNavigate();

      const [attempt, setAttempt] = useState(null);
      const [quiz, setQuiz] = useState(null);
      const [state, setState] = useState('loading');

      useEffect(() => {
            async function load() {
                  try {
                        const attemptData = await api.get(`/attempts/${attemptId}`);
                        const quizData = await api.get(`/quizzes/${attemptData.quiz_id}`);
                        setAttempt(attemptData);
                        setQuiz(quizData);
                        setState('ready');
                  } catch (err) {
                        if (err.status === 404) setState('not_found');
                        else if (err.status === 403) setState('forbidden');
                        else setState('error');
                  }
            }
            load();
      }, [attemptId]);

      if (state === 'loading') {
            return <TeacherLayout><p className="text-slate-500">Loading attempt...</p></TeacherLayout>;
      }

      if (state === 'not_found' || state === 'forbidden' || state === 'error') {
            const messages = {
                  not_found: "This attempt doesn't exist.",
                  forbidden: "You don't have permission to view this attempt.",
                  error: 'Failed to load this attempt. Please try again.',
            };
            return (
                  <TeacherLayout>
                        <Card className="max-w-md">
                              <p className="text-slate-500 mb-4">{messages[state]}</p>
                              <Button onClick={() => navigate(-1)}>Back</Button>
                        </Card>
                  </TeacherLayout>
            );
      }

      const isSubmitted = attempt.status === 'submitted';

      return (
            <TeacherLayout>
                  <Card className="max-w-lg">
                        <h1 className="text-lg font-semibold mb-1">Quiz: {quiz.title}</h1>
                        <p className="text-sm text-slate-500 mb-6">Attempt ID: {attempt.id}</p>

                        <div className="flex flex-col gap-2 text-sm mb-6">
                              <p>
                                    <span className="text-slate-500">Score: </span>
                                    {isSubmitted && attempt.score !== null ? attempt.score : '—'}
                              </p>
                              <p>
                                    <span className="text-slate-500">Status: </span>
                                    {isSubmitted ? 'Submitted' : 'In Progress'}
                              </p>
                              <p>
                                    <span className="text-slate-500">Started: </span>
                                    {new Date(attempt.started_at).toLocaleString()}
                              </p>
                              <p>
                                    <span className="text-slate-500">Submitted: </span>
                                    {isSubmitted ? new Date(attempt.submitted_at).toLocaleString() : '—'}
                              </p>
                        </div>

                        <Button variant="secondary" onClick={() => navigate(-1)}>
                              Back to Results
                        </Button>
                  </Card>
            </TeacherLayout>
      );
}