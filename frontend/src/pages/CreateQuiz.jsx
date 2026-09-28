import { useState } from 'react';
import { useNavigate } from 'react-router-dom';
import { api } from '../api/client';
import { useQuizQuestions } from '../hooks/useQuizQuestions';
import TeacherLayout from '../components/TeacherLayout';
import Card from '../components/Card';
import Input from '../components/Input';
import Button from '../components/Button';

export default function CreateQuiz() {
      const navigate = useNavigate();

      const [quizId, setQuizId] = useState(null);
      const [title, setTitle] = useState('');
      const [description, setDescription] = useState('');
      const [duration, setDuration] = useState(10);
      const [quizError, setQuizError] = useState('');
      const [savingQuiz, setSavingQuiz] = useState(false);

      const { questions, loading: questionsLoading, reload } = useQuizQuestions(quizId);

      const [qText, setQText] = useState('');
      const [options, setOptions] = useState([
            { text: '', isCorrect: false },
            { text: '', isCorrect: false },
      ]);
      const [questionError, setQuestionError] = useState('');
      const [savingQuestion, setSavingQuestion] = useState(false);
      const [actionNotice, setActionNotice] = useState('');

      async function handleCreateQuiz(e) {
            e.preventDefault();
            setQuizError('');
            setSavingQuiz(true);
            try {
                  const quiz = await api.post('/quizzes', {
                        title,
                        description,
                        duration_minutes: Number(duration),
                  });
                  setQuizId(quiz.id);
            } catch (err) {
                  setQuizError(err.message);
            } finally {
                  setSavingQuiz(false);
            }
      }

      function updateOptionText(idx, value) {
            setOptions((prev) => prev.map((o, i) => (i === idx ? { ...o, text: value } : o)));
      }

      function markCorrect(idx) {
            setOptions((prev) => prev.map((o, i) => ({ ...o, isCorrect: i === idx })));
      }

      function addOption() {
            setOptions((prev) => [...prev, { text: '', isCorrect: false }]);
      }

      function removeOption(idx) {
            setOptions((prev) => prev.filter((_, i) => i !== idx));
      }

      function resetQuestionForm() {
            setQText('');
            setOptions([{ text: '', isCorrect: false }, { text: '', isCorrect: false }]);
      }

      async function handleSaveQuestion(e) {
            e.preventDefault();
            setQuestionError('');

            const filledOptions = options.filter((o) => o.text.trim() !== '');
            if (qText.trim() === '') {
                  setQuestionError('Question text is required.');
                  return;
            }
            if (filledOptions.length < 2) {
                  setQuestionError('At least two options are required.');
                  return;
            }
            if (!filledOptions.some((o) => o.isCorrect)) {
                  setQuestionError('Please mark one option as correct.');
                  return;
            }

            setSavingQuestion(true);
            try {
                  const question = await api.post(`/quizzes/${quizId}/questions`, { text: qText });

                  for (const opt of filledOptions) {
                        await api.post(`/questions/${question.id}/options`, {
                              text: opt.text,
                              is_correct: opt.isCorrect,
                        });
                  }

                  resetQuestionForm();
                  await reload();
            } catch (err) {
                  setQuestionError(err.message);
            } finally {
                  setSavingQuestion(false);
            }
      }

      function handleUnsupported(action) {
            setActionNotice(`${action} isn't supported yet.`);
            setTimeout(() => setActionNotice(''), 3000);
      }

      if (!quizId) {
            return (
                  <TeacherLayout>
                        <Card className="max-w-lg">
                              <h1 className="text-xl font-semibold mb-6">Create Quiz</h1>
                              <form onSubmit={handleCreateQuiz} className="flex flex-col gap-4">
                                    <Input label="Title" value={title} onChange={(e) => setTitle(e.target.value)} required />
                                    <div className="flex flex-col gap-1">
                                          <label className="text-sm font-medium text-slate-700">Description</label>
                                          <textarea
                                                className="px-3 py-2 rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-slate-900"
                                                rows={3}
                                                value={description}
                                                onChange={(e) => setDescription(e.target.value)}
                                          />
                                    </div>
                                    <Input
                                          label="Duration (minutes)"
                                          type="number"
                                          min="1"
                                          value={duration}
                                          onChange={(e) => setDuration(e.target.value)}
                                          required
                                    />
                                    {quizError && <p className="text-sm text-red-500">{quizError}</p>}
                                    <Button type="submit" variant="accent" disabled={savingQuiz}>
                                          {savingQuiz ? 'Saving...' : 'Save & Continue'}
                                    </Button>
                              </form>
                        </Card>
                  </TeacherLayout>
            );
      }

      return (
            <TeacherLayout>
                  <h1 className="text-xl font-semibold mb-1">Quiz: {title}</h1>
                  <p className="text-sm text-slate-500 mb-6">
                        Add questions below, then return to your dashboard when done.
                  </p>

                  <Card className="max-w-2xl mb-8">
                        <h2 className="font-semibold mb-4">Question {questions.length + 1}</h2>
                        <form onSubmit={handleSaveQuestion} className="flex flex-col gap-4">
                              <Input label="Question text" value={qText} onChange={(e) => setQText(e.target.value)} />

                              <div className="flex flex-col gap-2">
                                    <label className="text-sm font-medium text-slate-700">Options</label>
                                    {options.map((opt, idx) => (
                                          <div key={idx} className="flex items-center gap-2">
                                                <input
                                                      type="radio"
                                                      name="correct-option"
                                                      checked={opt.isCorrect}
                                                      onChange={() => markCorrect(idx)}
                                                />
                                                <input
                                                      className="flex-1 px-3 py-2 rounded-xl border border-slate-300 focus:outline-none focus:ring-2 focus:ring-slate-900"
                                                      placeholder={`Option ${idx + 1}`}
                                                      value={opt.text}
                                                      onChange={(e) => updateOptionText(idx, e.target.value)}
                                                />
                                                {options.length > 2 && (
                                                      <button
                                                            type="button"
                                                            onClick={() => removeOption(idx)}
                                                            className="text-sm text-slate-400 hover:text-red-500"
                                                      >
                                                            ✕
                                                      </button>
                                                )}
                                          </div>
                                    ))}
                                    <Button type="button" variant="secondary" onClick={addOption} className="self-start">
                                          Add Option
                                    </Button>
                              </div>

                              {questionError && <p className="text-sm text-red-500">{questionError}</p>}

                              <Button type="submit" variant="accent" disabled={savingQuestion}>
                                    {savingQuestion ? 'Saving...' : 'Save Question'}
                              </Button>
                        </form>
                  </Card>

                  <h2 className="text-lg font-semibold mb-4">Questions</h2>

                  {actionNotice && <p className="text-sm text-amber-600 mb-4">{actionNotice}</p>}

                  {questionsLoading ? (
                        <p className="text-slate-500">Loading questions...</p>
                  ) : questions.length === 0 ? (
                        <p className="text-slate-500">No questions added yet.</p>
                  ) : (
                        <div className="flex flex-col gap-4 max-w-2xl">
                              {questions.map((q, idx) => (
                                    <Card key={q.id}>
                                          <p className="font-medium mb-2">
                                                {idx + 1}. {q.text}
                                          </p>
                                          <ul className="text-sm text-slate-600 mb-3 space-y-1">
                                                {q.options.map((opt, oIdx) => (
                                                      <li key={opt.id}>
                                                            {String.fromCharCode(65 + oIdx)}. {opt.text}
                                                            {opt.is_correct && (
                                                                  <span className="text-emerald-500 font-medium"> ✓ Correct</span>
                                                            )}
                                                      </li>
                                                ))}
                                          </ul>
                                          <div className="flex gap-2">
                                                <Button variant="secondary" onClick={() => handleUnsupported('Editing questions')}>
                                                      Edit
                                                </Button>
                                                <Button variant="danger" onClick={() => handleUnsupported('Deleting questions')}>
                                                      Delete
                                                </Button>
                                          </div>
                                    </Card>
                              ))}
                        </div>
                  )}

                  <div className="mt-8">
                        <Button onClick={() => navigate('/teacher/dashboard')}>Done — Back to Dashboard</Button>
                  </div>
            </TeacherLayout>
      );
}