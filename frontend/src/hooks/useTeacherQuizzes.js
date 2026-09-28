import { useState, useEffect } from 'react';
import { api } from '../api/client';
import { useAuth } from '../context/AuthContext';

export function useTeacherQuizzes() {
      const { user } = useAuth();
      const [quizzes, setQuizzes] = useState([]);
      const [totalAttempts, setTotalAttempts] = useState(0);
      const [uniqueStudents, setUniqueStudents] = useState(0);
      const [state, setState] = useState('loading');

      useEffect(() => {
            async function load() {
                  try {
                        const allQuizzes = await api.get('/quizzes');
                        const myQuizzes = allQuizzes.filter((q) => q.created_by === user.id);

                        const withDetails = await Promise.all(
                              myQuizzes.map(async (quiz) => {
                                    const questions = (await api.get(`/quizzes/${quiz.id}/questions`)) || [];
                                    return { ...quiz, questionCount: questions.length };
                              }),
                        );

                        const studentIds = new Set();
                        let attemptCount = 0;
                        if (myQuizzes.length > 0) {
                              const summariesPerQuiz = await Promise.all(
                                    myQuizzes.map((quiz) => api.get(`/quizzes/${quiz.id}/attempts`).catch(() => [])),
                              );
                              summariesPerQuiz.forEach((summaries) => {
                                    (summaries || []).forEach((s) => {
                                          attemptCount++;
                                          studentIds.add(s.user_id);
                                    });
                              });
                        }

                        setQuizzes(withDetails);
                        setTotalAttempts(attemptCount);
                        setUniqueStudents(studentIds.size);
                        setState('ready');
                  } catch (err) {
                        setState('error');
                  }
            }

            load();
      }, [user.id]);

      return { quizzes, totalAttempts, uniqueStudents, state };
}