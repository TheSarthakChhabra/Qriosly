import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useMyAttempts() {
      const [attempts, setAttempts] = useState([]);
      const [quizzesById, setQuizzesById] = useState({});
      const [state, setState] = useState('loading'); // loading | ready | error

      useEffect(() => {
            async function load() {
                  try {
                        const attemptList = await api.get('/my-attempts');

                        const uniqueQuizIds = [...new Set(attemptList.map((a) => a.quiz_id))];
                        const quizzes = await Promise.all(uniqueQuizIds.map((id) => api.get(`/quizzes/${id}`)));

                        const byId = {};
                        quizzes.forEach((q) => { byId[q.id] = q; });

                        setAttempts(attemptList);
                        setQuizzesById(byId);
                        setState('ready');
                  } catch (err) {
                        setState('error');
                  }
            }

            load();
      }, []);

      return { attempts, quizzesById, state };
}