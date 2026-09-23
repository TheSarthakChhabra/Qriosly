import { useParams } from 'react-router-dom';
import Card from '../components/Card';

export default function AttemptPage() {
      const { attemptId } = useParams();

      return (
            <div className="min-h-screen bg-slate-50 flex items-center justify-center p-4">
                  <Card className="max-w-md text-center">
                        <h1 className="text-xl font-semibold mb-2">Attempt Started</h1>
                        <p className="text-slate-500 mb-1">Attempt ID: {attemptId}</p>
                        <p className="text-sm text-slate-400">
                              The question-taking interface is coming in a later step.
                        </p>
                  </Card>
            </div>
      );
}