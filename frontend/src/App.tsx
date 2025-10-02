import './App.css'
import { useAuth } from './contexts/AuthContext';
import AuthPage from './components/AuthPage';
import Dashboard from './components/Dashboard';
import { Routes, Route, Navigate } from 'react-router-dom';
import StoryDetail from './components/StoryDetail';
import Layout from './components/Layout';

function App() {
  const { isAuthenticated } = useAuth();

  return (
    <Routes>
      <Route path="/" element={<Layout />}>
        {isAuthenticated ? (
          <>
            <Route index element={<Dashboard />} />
            <Route path="stories/:id" element={<StoryDetail />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </>
        ) : (
          <>
            <Route path="auth" element={<AuthPage />} />
            <Route path="*" element={<Navigate to="/auth" replace />} />
          </>
        )}
      </Route>
    </Routes>
  );
};

export default App;
