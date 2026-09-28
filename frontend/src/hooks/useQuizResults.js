import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useQuizResults(quizId) {
      const [quiz, setQuiz] = useState(null);
      const [summaries, setSummaries] = useState([]);
      const [state, setState] = useState('loading'); // loading | ready | not_found | forbidden | error

      useEffect(() => {
            async function load() {
                  try {
                        const [quizData, summaryList] = await Promise.all([
                              api.get(`/quizzes/${quizId}`),
                              api.get(`/quizzes/${quizId}/attempts`),
                        ]);
                        setQuiz(quizData);
                        setSummaries(summaryList || []);
                        setState('ready');
                  } catch (err) {
                        if (err.status === 404) setState('not_found');
                        else if (err.status === 403) setState('forbidden');
                        else setState('error');
                  }
            }

            load();
      }, [quizId]);

      return { quiz, summaries, state };
}