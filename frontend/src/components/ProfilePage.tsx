import React, { useState, useEffect } from 'react';
import apiClient from '../apiClient';

interface UserStats {
  total_word_count: number;
}

const ProfilePage: React.FC = () => {
  const [stats, setStats] = useState<UserStats | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const response = await apiClient.get<UserStats>('/users/me/stats');
        setStats(response.data);
      } catch (err) {
        setError('Failed to load user stats.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchStats();
  }, []);

  if (isLoading) return <p className="text-center">Loading stats...</p>
  if (error) return <p className="text-center text-red-500">{error}</p>

  return (
    <div className="bg-white p-8 rounded-lg shadow-md">
      <h1 className="text-3xl font-bold mb-6 text-center">
        <div className="text-center">
          <p className="text-lg text-gray-600">Total Words Read</p>
          <p className="text-5xl font-bold text-blue-600 mt-2">
            {stats ? stats.total_word_count.toLocaleString() : 0}
          </p>
        </div>
      </h1>
    </div>
  );
};

export default ProfilePage;