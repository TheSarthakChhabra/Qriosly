import { useState, useEffect, useCallback } from 'react';
import { api } from '../api/client';

export function useAdminUsers() {
  const [users, setUsers] = useState([]);
  const [state, setState] = useState('loading');

  const reload = useCallback(async () => {
    setState('loading');
    try {
      const data = await api.get('/admin/users');
      setUsers(data || []);
      setState('ready');
    } catch (err) {
      setState('error');
    }
  }, []);

  useEffect(() => { reload(); }, [reload]);

  async function updateRole(userId, newRole) {
    const updated = await api.put(`/admin/users/${userId}/role`, { role: newRole });
    setUsers((prev) => prev.map((u) => (u.id === updated.id ? updated : u)));
  }

  return { users, state, updateRole };
}