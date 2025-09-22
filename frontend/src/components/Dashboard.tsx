import React from 'react';
import { useAuth } from '../contexts/AuthContext';
import StoryGenerator from './StoryGenerator';

const Dashboard: React.FC = () => {
  const { logout } = useAuth();

  return (
    <div className="w-full max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        <button 
          onClick={logout}
          className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition-colors"
          
        >
          Logout
        </button>
      </div>

      <StoryGenerator />

      {/* TODO: 他のダッシュボードコンテンツをここに追加可能 */}
    </div>
  );
};

export default Dashboard;