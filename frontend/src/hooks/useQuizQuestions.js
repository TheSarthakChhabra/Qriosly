import { useState, useEffect, useCallback } from 'react';
import { api } from '../api/client';

export function useQuizQuestions(quizId) {
      const [questions, setQuestions] = useState([]);
      const [loading, setLoading] = useState(true);

      const reload = useCallback(async () => {
            if (!quizId) return;
            setLoading(true);
            const questionList = (await api.get(`/quizzes/${quizId}/questions`)) || [];
            const withOptions = await Promise.all(
                  questionList.map(async (q) => {
                        const options = (await api.get(`/questions/${q.id}/options`)) || [];
                        return { ...q, options };
                  }),
            );
            setQuestions(withOptions);
            setLoading(false);
      }, [quizId]);

      useEffect(() => { reload(); }, [reload]);

      return { questions, loading, reload };
}