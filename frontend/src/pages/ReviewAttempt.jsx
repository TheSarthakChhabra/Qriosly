import { useParams, useNavigate, useLocation } from 'react-router-dom';
import { useState } from 'react';
import { api } from '../api/client';
import Card from '../components/Card';
import Button from '../components/Button';
import Modal from '../components/Modal';
import QuestionPalette from '../components/QuestionPalette';

export default function ReviewAttempt() {
      const { attemptId } = useParams();
      const navigate = useNavigate();
      const location = useLocation();

      // Passed forward from AttemptPage via navigate state — no re-fetching needed.
      const { questions, answers } = location.state || { questions: [], answers: {} };

      const [confirmOpen, setConfirmOpen] = useState(false);
      const [submitting, setSubmitting] = useState(false);
      const [submitError, setSubmitError] = useState('');

      const answeredCount = questions.filter((q) => !!answers[q.id]).length;
      const unansweredCount = questions.length - answeredCount;

      async function handleSubmit() {
            setSubmitting(true);
            setSubmitError('');

            try {
                  await api.post(`/attempts/${attemptId}/submit`, {});
                  navigate(`/result/${attemptId}`);
            } catch (err) {
                  // Handle each real backend outcome distinctly, not just a generic failure.
                  if (err.message.toLowerCase().includes('already') || err.message.toLowerCase().includes('submitted')) {
                        setSubmitError('This attempt has already been submitted.');
                  } else if (err.message.toLowerCase().includes('expired') || err.message.toLowerCase().includes('time limit')) {
                        setSubmitError('The time limit for this attempt has passed. Your answers were saved but the attempt could not be submitted.');
                  } else if (err.message === 'Failed to fetch') {
                        setSubmitError('Network error — please check your connection and try again.');
                  } else {
                        setSubmitError(err.message || 'Something went wrong while submitting.');
                  }
                  setSubmitting(false);
                  setConfirmOpen(false);
            }
      }

      if (questions.length === 0) {
            return (
                  <div className="min-h-screen flex items-center justify-center p-4">
                        <Card className="max-w-md text-center">
                              <p className="text-slate-500 mb-4">No review data available.</p>
                              <Button variant="secondary" onClick={() => navigate(`/attempt/${attemptId}`)}>
                                    Back to Quiz
                              </Button>
                        </Card>
                  </div>
            );
      }

      return (
            <div className="min-h-screen bg-slate-50 p-4 md:p-8">
                  <Card className="max-w-2xl mx-auto">
                        <h1 className="text-xl font-semibold text-center mb-6">Review Your Attempt</h1>

                        <p className="text-slate-600 mb-4">{questions.length} Questions</p>

                        <div className="flex gap-6 mb-6">
                              <p className="text-sm">
                                    Answered: <span className="font-semibold text-emerald-500">{answeredCount}</span>
                              </p>
                              <p className="text-sm">
                                    Unanswered: <span className="font-semibold text-red-500">{unansweredCount}</span>
                              </p>
                        </div>

                        <QuestionPalette
                              questions={questions}
                              answers={answers}
                              currentIndex={-1}
                              onSelect={(idx) => navigate(`/attempt/${attemptId}`, { state: { jumpTo: idx } })}
                        />

                        {unansweredCount > 0 && (
                              <p className="text-sm text-amber-500 font-medium mt-6">
                                    ⚠ You have {unansweredCount} unanswered question{unansweredCount !== 1 ? 's' : ''}.
                              </p>
                        )}

                        {submitError && <p className="text-sm text-red-500 mt-4">{submitError}</p>}

                        <div className="flex justify-between mt-8">
                              <Button
                                    variant="secondary"
                                    onClick={() => navigate(`/attempt/${attemptId}`, { state: { answers, resumeQuestions: questions } })}
                              >
                                    Back to Quiz
                              </Button>
                              <Button variant="accent" onClick={() => setConfirmOpen(true)}>
                                    Submit Quiz
                              </Button>
                        </div>
                  </Card>

                  <Modal isOpen={confirmOpen} onClose={() => setConfirmOpen(false)}>
                        <h2 className="text-lg font-semibold mb-2">Are you sure you want to submit?</h2>
                        <p className="text-slate-500 text-sm mb-6">
                              You won't be able to change your answers after submission.
                        </p>
                        <div className="flex justify-end gap-2">
                              <Button variant="secondary" onClick={() => setConfirmOpen(false)} disabled={submitting}>
                                    Cancel
                              </Button>
                              <Button variant="danger" onClick={handleSubmit} disabled={submitting}>
                                    {submitting ? 'Submitting...' : 'Submit'}
                              </Button>
                        </div>
                  </Modal>
            </div>
      );
}