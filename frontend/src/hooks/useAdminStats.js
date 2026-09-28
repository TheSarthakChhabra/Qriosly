import { useState, useEffect } from 'react';
import { api } from '../api/client';

export function useAdminStats() {
      const [stats, setStats] = useState({ students: 0, teachers: 0, quizzes: 0 });
      const [loading, setLoading] = useState(true);

      useEffect(() => {
            async function load() {
                  try {
                        const [users, quizzes] = await Promise.all([api.get('/admin/users'), api.get('/quizzes')]);
                        setStats({
                              students: users.filter((u) => u.role === 'student').length,
                              teachers: users.filter((u) => u.role === 'teacher').length,
                              quizzes: quizzes.length,
                        });
                  } catch (err) {
                        console.error('Failed to load admin stats:', err);
                  } finally {
                        setLoading(false);
                  }
            }
            load();
      }, []);

      return { stats, loading };
}