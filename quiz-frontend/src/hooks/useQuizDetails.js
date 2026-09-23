import { useState, useEffect } from "react";
import { api } from '../api/client'

export function useQuizDetails(quizId) {
      const [quiz, setQuiz] = useState(null);
      const [questionCount, setQuestionCount] = useState(null);
      const [loading, setLoading] = useState(true);
      const [error, setError] = useState('');

      useEffect(() => {
            async function load() {
                  try {
                        const [quizData, questions] = await Promise.all([
                              api.get(`/quizzes/${quizId}`),
                              api.get(`/quizzes/${quizId}/quuestions`),
                        ]);
                        setQuiz(quizData);
                        setQuestionCount(questions.length);
                  } catch (err) {
                        setError(err.message);
                  } finally {
                        setLoading(false);
                  }
            }
            load();
      }, [quizId]);
      return { quiz, questionCount, loading, error };
}