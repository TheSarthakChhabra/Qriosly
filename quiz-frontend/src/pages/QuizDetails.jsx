import { useParams, useNavigate } from 'react-router-dom';
import Button from '../components/Button';
import Card from '../components/Card';

export default function QuizDetails() {
      const { id } = useParams();
      const navigate = useNavigate();

      return (
            <div className="min-h-screen bg-slate-50 p-8">
                  <Card className="max-w-md">
                        <h1 className="text-xl font-semibold mb-2">Quiz Details</h1>
                        <p className="text-slate-500 mb-4">Quiz ID: {id}</p>
                        <p className="text-sm text-slate-400 mb-4">
                              Quiz-taking interface coming in a later step.
                        </p>
                        <Button variant="secondary" onClick={() => navigate(-1)}>
                              Back
                        </Button>
                  </Card>
            </div>
      );
}