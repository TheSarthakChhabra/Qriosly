import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useQuizzes() {
      const [quizzes, setQuizzes] = useState([]);
      const [loading, setLoading] = useState(true);

      useEffect(() => {
            async function load() {
                  try {
                        const quizList = await api.get('/quizzes');

                        const withCounts = await Promise.all(
                              quizList.map(async (quiz) => {
                                    const questions = await api.get(`/quizzes/${quiz.id}/questions`);
                                    return { ...quiz, questionCount: questions.length };
                              }),
                        );

                        setQuizzes(withCounts);
                  } catch (err) {
                        console.error('Failed to load quizzes:', err);
                  } finally {
                        setLoading(false);
                  }
            }

            load();
      }, []);

      return { quizzes, loading };
}