import { useState } from 'react';
import { useAuth } from '../context/AuthContext';
import { useAdminUsers } from '../hooks/useAdminUsers';
import AdminLayout from '../components/AdminLayout';
import Card from '../components/Card';
import Button from '../components/Button';

export default function AdminUsers() {
      const { user: currentUser } = useAuth();
      const { users, state, updateRole } = useAdminUsers();
      const [error, setError] = useState('');
      const [pendingId, setPendingId] = useState(null);

      async function handleToggle(targetUser) {
            const newRole = targetUser.role === 'student' ? 'teacher' : 'student';
            setError('');
            setPendingId(targetUser.id);
            try {
                  await updateRole(targetUser.id, newRole);
            } catch (err) {
                  setError(err.message);
            } finally {
                  setPendingId(null);
            }
      }

      if (state === 'loading') {
            return <AdminLayout><p className="text-slate-500">Loading users...</p></AdminLayout>;
      }

      if (state === 'error') {
            return (
                  <AdminLayout>
                        <Card className="max-w-md"><p className="text-red-500">Failed to load users.</p></Card>
                  </AdminLayout>
            );
      }

      return (
            <AdminLayout>
                  <h1 className="text-xl font-semibold mb-6">Users</h1>

                  {error && <p className="text-sm text-red-500 mb-4">{error}</p>}

                  <Card className="p-0 overflow-hidden">
                        <table className="w-full text-left text-sm">
                              <thead className="bg-slate-50 text-slate-500 border-b border-slate-200">
                                    <tr>
                                          <th className="px-4 py-3 font-medium">Name</th>
                                          <th className="px-4 py-3 font-medium">Email</th>
                                          <th className="px-4 py-3 font-medium">Role</th>
                                          <th className="px-4 py-3 font-medium"></th>
                                    </tr>
                              </thead>
                              <tbody>
                                    {users.map((u) => {
                                          const isSelf = u.id === currentUser.id;
                                          const canToggle = !isSelf && (u.role === 'student' || u.role === 'teacher');

                                          return (
                                                <tr key={u.id} className="border-b border-slate-100 last:border-0">
                                                      <td className="px-4 py-3">{u.name}</td>
                                                      <td className="px-4 py-3 text-slate-500">{u.email}</td>
                                                      <td className="px-4 py-3 capitalize">{u.role}</td>
                                                      <td className="px-4 py-3 text-right">
                                                            {canToggle && (
                                                                  <Button
                                                                        variant="secondary"
                                                                        disabled={pendingId === u.id}
                                                                        onClick={() => handleToggle(u)}
                                                                  >
                                                                        {pendingId === u.id
                                                                              ? 'Updating...'
                                                                              : u.role === 'student'
                                                                                    ? 'Make Teacher'
                                                                                    : 'Make Student'}
                                                                  </Button>
                                                            )}
                                                            {isSelf && <span className="text-xs text-slate-400">(you)</span>}
                                                      </td>
                                                </tr>
                                          );
                                    })}
                              </tbody>
                        </table>
                  </Card>
            </AdminLayout>
      );
}