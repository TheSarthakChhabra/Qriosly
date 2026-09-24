import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useAttemptResult(attemptId) {
      const [attempt, setAttempt] = useState(null);
      const [quiz, setQuiz] = useState(null);
      const [totalQuestions, setTotalQuestions] = useState(null);
      const [state, setState] = useState('loading'); // loading | ready | not_found | forbidden | error

      useEffect(() => {
            async function load() {
                  try {
                        const attemptData = await api.get(`/attempts/${attemptId}`);
                        const [quizData, questions] = await Promise.all([
                              api.get(`/quizzes/${attemptData.quiz_id}`),
                              api.get(`/quizzes/${attemptData.quiz_id}/questions`),
                        ]);

                        setAttempt(attemptData);
                        setQuiz(quizData);
                        setTotalQuestions((questions || []).length);
                        setState('ready');
                  } catch (err) {
                        if (err.status === 404) setState('not_found');
                        else if (err.status === 403) setState('forbidden');
                        else setState('error');
                  }
            }

            load();
      }, [attemptId]);

      return { attempt, quiz, totalQuestions, state };
}