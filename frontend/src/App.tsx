import './App.css'
import { useAuth } from './contexts/AuthContext';
import AuthPage from './components/AuthPage';
import Dashboard from './components/Dashboard';

function App() {
  const { isAuthenticated } = useAuth();

  return (
    <div className="min-h-screen bg-gray-100">
      <div className="container mx-auto p-4 py-10">
        <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">
          English Tadoku App
        </h1>
        {isAuthenticated ? <Dashboard /> : <AuthPage />}
      </div>
    </div>
  );
};

export default App;
