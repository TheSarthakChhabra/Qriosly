import { useParams } from 'react-router-dom';
import { useState } from 'react';
import { useAttempt } from '../hooks/useAttempt';
import { useCountdown } from '../hooks/useCountdown';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useLocation } from 'react-router-dom';
import Card from '../components/Card';
import Button from '../components/Button';
import QuestionPalette from '../components/QuestionPalette';
import Timer from '../components/Timer';

export default function AttemptPage() {
      const { attemptId } = useParams();
      const { attempt, quiz, questions, loading, error } = useAttempt(attemptId);
      const navigate = useNavigate();
      const [currentIndex, setCurrentIndex] = useState(0);
      const [answeredIds, setAnsweredIds] = useState(new Set());
      const [saveError, setSaveError] = useState('');
      const location = useLocation();
      const [answers, setAnswers] = useState(location.state?.answers || {});
      const deadline = attempt && quiz
            ? new Date(attempt.started_at).getTime() + quiz.duration_minutes * 60000
            : null;
      const remainingMs = useCountdown(deadline);
      const expired = remainingMs !== null && remainingMs <= 0;

      if (loading) return <div className="min-h-screen flex items-center justify-center">Loading...</div>;
      if (error) return <div className="min-h-screen flex items-center justify-center text-red-500">{error}</div>;

      const currentQuestion = questions[currentIndex];

      async function handleSelect(optionId) {
            if (expired) return;
            const questionId = currentQuestion.id;
            const previousSelection = answers[questionId];
            const alreadyAnswered = !!previousSelection;
            setAnswers((prev) => ({ ...prev, [questionId]: optionId }));
            setSaveError('');
            try {
                  if (alreadyAnswered) {
                        await api.put(`/attempts/${attemptId}/answers/${questionId}`, { selected_option_id: optionId });
                  } else {
                        await api.post(`/attempts/${attemptId}/answers`, {
                              question_id: questionId,
                              selected_option_id: optionId,
                        });
                  }
            } catch (err) {
                  setSaveError(err.message);
                  setAnswers((prev) => ({ ...prev, [questionId]: previousSelection }));
            }
      }

      return (
            <div className="min-h-screen bg-slate-50 p-4 md:p-8">
                  <Card className="max-w-3xl mx-auto">
                        <div className="flex items-center justify-between mb-4">
                              <h1 className="text-lg font-semibold">{quiz.title}</h1>
                              <Timer remainingMs={remainingMs} />
                        </div>

                        <hr className="border-slate-200 mb-4" />

                        <p className="text-sm text-slate-500 mb-2">
                              Question {currentIndex + 1} of {questions.length}
                        </p>
                        <p className="text-lg font-medium mb-4">{currentQuestion.text}</p>

                        <div className="flex flex-col gap-2 mb-6">
                              {currentQuestion.options.map((opt) => {
                                    const selected = answers[currentQuestion.id] === opt.id;
                                    return (
                                          <button
                                                key={opt.id}
                                                onClick={() => handleSelect(opt.id)}
                                                disabled={expired}
                                                className={`text-left px-4 py-3 rounded-xl border transition-colors flex items-center ${selected ? 'border-slate-900 bg-slate-100' : 'border-slate-200 bg-white hover:bg-slate-50'
                                                      } ${expired ? 'opacity-50 cursor-not-allowed' : ''}`}
                                          >
                                                <span
                                                      className={`inline-block w-4 h-4 rounded-full border mr-3 ${selected ? 'bg-indigo-700 border-indigo-700' : 'border-slate-300'
                                                            }`}
                                                />
                                                {opt.text}
                                          </button>
                                    );
                              })}
                        </div>

                        {saveError && <p className="text-sm text-red-500 mb-4">{saveError}</p>}

                        <div className="flex justify-between mb-6">
                              <Button
                                    variant="secondary"
                                    disabled={currentIndex === 0}
                                    onClick={() => setCurrentIndex((i) => i - 1)}
                              >
                                    Previous
                              </Button>
                              <Button
                                    disabled={currentIndex === questions.length - 1}
                                    onClick={() => setCurrentIndex((i) => i + 1)}
                              >
                                    Next
                              </Button>
                              <Button
                                    variant="accent"
                                    onClick={() => navigate(`/attempt/${attemptId}/review`, { state: { questions, answers } })}
                              >
                                    Review & Submit
                              </Button>
                        </div>

                        <hr className="border-slate-200 mb-4" />

                        <p className="text-sm font-medium text-slate-600 mb-2">Questions</p>
                        <QuestionPalette
                              questions={questions}
                              answers={answers}
                              currentIndex={currentIndex}
                              onSelect={setCurrentIndex}
                        />
                  </Card>
            </div>
      );
}
