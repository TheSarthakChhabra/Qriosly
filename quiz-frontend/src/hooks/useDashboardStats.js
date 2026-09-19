import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useDashboardStats() {
      const [stats, setStats] = useState({ availableQuizzes: 0, completedAttempts: 0, averageScore: null });
      const [loading, setLoading] = useState(true);

      useEffect(() => {
            async function load() {
                  try {
                        const [quizzes, attempts] = await Promise.all([
                              api.get('/quizzes'),
                              api.get('/my-attempts'),
                        ]);

                        const completed = attempts.filter((a) => a.status === 'submitted');
                        const scores = completed.map((a) => a.score).filter((s) => s !== null && s !== undefined);
                        const average = scores.length > 0
                              ? Math.round(scores.reduce((sum, s) => sum + s, 0) / scores.length)
                              : null;

                        setStats({
                              availableQuizzes: quizzes.length,
                              completedAttempts: completed.length,
                              averageScore: average,
                        });
                  } catch (err) {
                        console.error('Failed to load dashboard stats:', err);
                  } finally {
                        setLoading(false);
                  }
            }

            load();
      }, []);

      return { stats, loading };
}