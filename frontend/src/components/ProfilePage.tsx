import React, { useState, useEffect } from 'react';
import apiClient from '../apiClient';

interface UserStats {
  total_word_count: number;
  today_word_count: number;
  weekly_word_count: number;
  monthly_word_count: number;
  yearly_word_count: number;
  last_7_days_word_count: { [date: string]: number };
}

const StatCard: React.FC<{ label: string; value: number }> = ({ label, value }) => (
  <div className="bg-gray-100 p-4 rounded-lg shadow-md text-center">
    <p className="text-sm text-gray-600">{label}</p>
    <p className="text-2xl font-bold mt-2">{value.toLocaleString()}</p>
  </div>
);

const BarChart: React.FC<{ data: { [date: string]: number } }> = ({ data }) => {
  if (!data || Object.keys(data).length === 0) {
    return <p className="text-center">No recent activity.</p>;
  }

  const entries = Object.entries(data);
  const values = entries.map(([, value]) => value);
  const maxValue = Math.max(...values, 1);

  return (
    <div className="mt-8">
      <h3 className="text-lg font-semibold mb-4 text-center">過去7日間</h3>
      <div className="flex justify-around items-end h-40 bg-gray-50 p-4 rounded-lg">
        {entries.map(([date, value], index) => (
          <div key={index} className="relative flex flex-col items-center justify-end w-10 h-full">
            
            {value > 0 && (
              <div className="absolute -top-5 text-xs font-semibold text-gray-600">
                {value.toLocaleString()}
              </div>
            )}

            <div
              className="w-full bg-blue-500 rounded-t-md transition-all duration-300"
              style={{ height: `${(value / maxValue) * 100}%` }}
              title={`${value.toLocaleString()} words on ${new Date(date).toLocaleDateString()}`}
            ></div>
            <p className="text-xs mt-2">
              {new Date(date).toLocaleDateString('en-US', { weekday: 'short' })}
            </p>
          </div>
        ))}
      </div>
    </div>
  );
};

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
        setError('統計情報の読み込みに失敗しました。');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    
    fetchStats();
  }, []);

  if (isLoading) return <p className="text-center">統計情報を読み込み中...</p>
  if (error) return <p className="text-center text-red-500">{error}</p>

  return (
    <div className="bg-white p-8 rounded-lg shadow-md max-w-2xl mx-auto">
      <h1 className="text-3xl font-bold text-center mb-8">読書統計</h1>
      
      {/* --- 累計単語数 (メイン) --- */}
      <div className="text-center mb-8">
        <p className="text-lg text-gray-600">累計単語数</p>
        <p className="text-6xl font-bold text-blue-600 mt-2">
          {stats?.total_word_count.toLocaleString()}
        </p>
      </div>

      {/* --- 期間別統計 (サブ) --- */}
      <div className="grid grid-cols-2 md:grid-cols-4 gap-4 mb-8">
        <StatCard label="今日" value={stats?.today_word_count || 0} />
        <StatCard label="今週" value={stats?.weekly_word_count || 0} />
        <StatCard label="今月" value={stats?.monthly_word_count || 0} />
        <StatCard label="今年" value={stats?.yearly_word_count || 0} />
      </div>

      {/* --- 過去7日間のグラフ --- */}
      {stats && <BarChart data={stats.last_7_days_word_count} />}
    </div>
  );
};

export default ProfilePage;