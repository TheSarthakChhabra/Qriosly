import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useAttempt(attemptId) {
      const [attempt, setAttempt] = useState(null);
      const [quiz, setQuiz] = useState(null);
      const [questions, setQuestions] = useState([]);
      const [loading, setLoading] = useState(true);
      const [error, setError] = useState('');

      useEffect(() => {
            async function load() {
                  try {
                        const attemptData = await api.get(`/attempts/${attemptId}`);
                        const quizData = await api.get(`/quizzes/${attemptData.quiz_id}`);
                        const questionList = await api.get(`/quizzes/${attemptData.quiz_id}/questions`);

                        const withOptions = await Promise.all(
                              questionList.map(async (q) => {
                                    const options = await api.get(`/questions/${q.id}/options`);
                                    return { ...q, options };
                              }),
                        );

                        setAttempt(attemptData);
                        setQuiz(quizData);
                        setQuestions(withOptions);
                  } catch (err) {
                        setError(err.message);
                  } finally {
                        setLoading(false);
                  }
            }

            load();
      }, [attemptId]);

      return { attempt, quiz, questions, loading, error };
}