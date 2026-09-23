import { useState } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import Card from '../components/Card';
import Input from '../components/Input';
import Button from '../components/Button';

export default function Login() {
      const [email, setEmail] = useState('');
      const [password, setPassword] = useState('');
      const [error, setError] = useState('');
      const [loading, setLoading] = useState(false);

      const { login } = useAuth();
      const navigate = useNavigate();

      async function handleSubmit(e) {
            e.preventDefault();
            setError('');
            setLoading(true);

            try {
                  const user = await login(email, password);
                  navigate(`/${user.role}/dashboard`);
            } catch (err) {
                  setError(err.message);
            } finally {
                  setLoading(false);
            }
      }

      return (
            <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
                  <Card className="w-full max-w-sm">
                        <h1 className="text-2xl font-semibold mb-6 text-center">Log In</h1>

                        <form onSubmit={handleSubmit} className="flex flex-col gap-4">
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
                                    {loading ? 'Logging in...' : 'Login'}
                              </Button>
                        </form>

                        <p className="text-sm text-slate-500 text-center mt-4">
                              Don't have an account?{' '}
                              <Link to="/register" className="text-indigo-600 hover:underline">
                                    Register
                              </Link>
                        </p>
                  </Card>
            </div>
      );
}