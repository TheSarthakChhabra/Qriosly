import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Card from '../components/Card';
import Input from '../components/Input';
import Button from '../components/Button';

export default function Register() {
      const [name, setName] = useState('');
      const [email, setEmail] = useState('');
      const [password, setPassword] = useState('');
      const [error, setError] = useState('');
      const [loading, setLoading] = useState(false);

      const { register } = useAuth();
      const navigate = useNavigate();

      async function handleSubmit(e) {
            e.preventDefault();
            setError('');
            setLoading(true);

            try {
                  await register(name, email, password);
                  navigate('/login');
            } catch (err) {
                  setError(err.message);
            } finally {
                  setLoading(false);
            }
      }

      return (
            <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
                  <Card className="w-full max-w-sm">
                        <h1 className="text-2xl font-semibold mb-6 text-center">Register</h1>

                        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
                              <Input label="Name" value={name} onChange={(e) => setName(e.target.value)} required />
                              <Input
                                    label="Email"
                                    type="email"
                                    value={email}
                                    onChange={(e) => setEmail(e.target.value)}
                                    required
                              />
                              <Input
                                    label="Password"
                                    type="password"
                                    value={password}
                                    onChange={(e) => setPassword(e.target.value)}
                                    required
                              />

                              {error && <p className="text-sm text-red-500">{error}</p>}

                              <Button type="submit" disabled={loading}>
                                    {loading ? 'Registering...' : 'Register'}
                              </Button>
                        </form>

                        <p className="text-sm text-slate-500 text-center mt-4">
                              Already have an account?{' '}
                              <Link to="/login" className="text-indigo-600 hover:underline">
                                    Login
                              </Link>
                        </p>
                  </Card>
            </div>
      );
}