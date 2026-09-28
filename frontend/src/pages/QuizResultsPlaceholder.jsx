import { useParams } from 'react-router-dom';
import Card from '../components/Card';
import TeacherLayout from '../components/TeacherLayout';

export default function QuizResultsPlaceholder() {
      const { id } = useParams();
      return (
            <TeacherLayout>
                  <Card className="max-w-md">
                        <h1 className="text-lg font-semibold mb-2">Quiz {id}</h1>
                        <p className="text-slate-500 text-sm">This page is coming in a later step.</p>
                  </Card>
            </TeacherLayout>
      );
}