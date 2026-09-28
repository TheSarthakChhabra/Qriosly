import { Routes, Route, Navigate } from 'react-router-dom';
import Login from './pages/Login';
import Register from './pages/Register';
import StudentDashboard from './pages/StudentDashboard';
import TeacherDashboard from './pages/TeacherDashboard';
import AdminDashboard from './pages/AdminDashboard';
import QuizDetails from './pages/QuizDetails';
import ProtectedRoute from './components/ProtectedRoute';
import AttemptPage from './pages/AttemptPage';
import ReviewAttempt from './pages/ReviewAttempt';
import ResultPage from './pages/ResultPage';
import MyAttempts from './pages/MyAttempts';
import CreateQuiz from './pages/CreateQuiz';
import QuizResultsPlaceholder from './pages/QuizResultsPlaceholder';
import QuizResults from './pages/QuizResults';
import TeacherAttemptDetails from './pages/TeacherAttemptDetails';
import AdminUsers from './pages/AdminUsers';

function App() {
  return (
    <Routes>
      <Route path="/" element={<Navigate to="/login" replace />} />
      <Route path="/login" element={<Login />} />
      <Route path="/register" element={<Register />} />

      <Route
        path="/student/dashboard"
        element={
          <ProtectedRoute allowedRoles={['student']}>
            <StudentDashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/quizzes/:id"
        element={
          <ProtectedRoute allowedRoles={['student']}>
            <QuizDetails />
          </ProtectedRoute>
        }
      />
      <Route
        path="/teacher/dashboard"
        element={
          <ProtectedRoute allowedRoles={['teacher']}>
            <TeacherDashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/admin/dashboard"
        element={
          <ProtectedRoute allowedRoles={['admin']}>
            <AdminDashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/attempt/:attemptId"
        element={
          <ProtectedRoute allowedRoles={['student']}>
            <AttemptPage />
          </ProtectedRoute>
        }
      />
      <Route
        path="/attempt/:attemptId/review"
        element={<ProtectedRoute allowedRoles={['student']}><ReviewAttempt /></ProtectedRoute>}
      />
      <Route
        path="/result/:attemptId"
        element={<ProtectedRoute allowedRoles={['student']}><ResultPage /></ProtectedRoute>}
      />
      <Route
        path='/student/attempts'
        element={<ProtectedRoute allowedRoles={['student']}><MyAttempts /></ProtectedRoute>}
      />
      <Route path="/teacher/quizzes" element={<ProtectedRoute allowedRoles={['teacher']}><TeacherDashboard /></ProtectedRoute>} />
      <Route path="/teacher/quizzes/create" element={<ProtectedRoute allowedRoles={['teacher']}><CreateQuiz /></ProtectedRoute>} />
      <Route path="/teacher/quizzes/:id/edit" element={<ProtectedRoute allowedRoles={['teacher']}><QuizResultsPlaceholder /></ProtectedRoute>} />
      <Route path="/teacher/quizzes/:id/results" element={<ProtectedRoute allowedRoles={['teacher']}><QuizResultsPlaceholder /></ProtectedRoute>} />
      <Route path="/teacher/quizzes/:id/results" element={<ProtectedRoute allowedRoles={['teacher']}><QuizResults /></ProtectedRoute>} />
      <Route path="/teacher/attempts/:attemptId" element={<ProtectedRoute allowedRoles={['teacher']}><TeacherAttemptDetails /></ProtectedRoute>} />
      <Route path="/admin/users" element={<ProtectedRoute allowedRoles={['admin']}><AdminUsers /></ProtectedRoute>} />
    </Routes>
  );
}

export default App;