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
import MyAttempts from './pages/MyAttempts'
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
    </Routes>
  );
}

export default App;