import './App.css'
import { useAuth } from './contexts/AuthContext';
import AuthPage from './components/AuthPage';
import Dashboard from './components/Dashboard';
import { Routes, Route, Navigate } from 'react-router-dom';
import StoryDetail from './components/StoryDetail';

function App() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto p-4 py-10">
        <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">
          English Tadoku App
        </h1>
        {/* {isAuthenticated ? <Dashboard /> : <AuthPage />} */}
        <Routes>
          {!isAuthenticated && (
            <>
              <Route path="/auth" element={<AuthPage />} />
              <Route path="*" element={<Navigate to="/auth" replace />} />
            </>
          )}
          {isAuthenticated && (
            <>
              <Route path="/" element={<Dashboard />} />
              <Route path="/stories/:id" element={<StoryDetail />} />
              <Route path="*" element={<Navigate to="/" replace />} />
            </>
          )}
        </Routes>
      </div>
    </div>
  );
};

export default App;
