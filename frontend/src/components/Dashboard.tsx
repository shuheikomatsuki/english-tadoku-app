import React, {useState, useEffect } from 'react';
import { useAuth } from '../contexts/AuthContext';
import StoryGenerator from './StoryGenerator';
import StoryList from './StoryList';
import apiClient from '../apiClient';
import type { Story } from '../types';


const Dashboard: React.FC = () => {
  const { logout } = useAuth();

  const [stories, setStories] = useState<Story[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchStories = async () => {
      try {
        const response = await apiClient.get<Story[]>('/stories');
        setStories(response.data || []);
      } catch (err) {
        setError('Failed to load stories.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchStories();
  }, []);

  const handleStoryGenerated = (newStory: Story) => {
    setStories(prevStories => [newStory, ...prevStories]);
  };

  const handleDeleteStory = async (id: number) => {
    if (!window.confirm('Are you sure you want to delete this story?')) {
      return;
    }

    try {
      await apiClient.delete(`/stories/${id}`);
      setStories(prevStories => prevStories.filter(story => story.id !== id));
    } catch (err) {
      console.error('Failed to delete story:', err);
      alert('Failed to delete the story. Please try again.');
    }
  };

  return (
    <div className="w-full max-w-4xl mx-auto">
      <div className="flex justify-between items-center mb-8">
        <h1 className="text-3xl font-bold">Dashboard</h1>
        {/* <button 
          onClick={logout}
          className="bg-red-500 hover:bg-red-600 text-white font-bold py-2 px-4 rounded transition-colors"
          
        >
          Logout
        </button> */}
      </div>

      <StoryGenerator onStoryGenerated={handleStoryGenerated} />

      {/* TODO: 他のダッシュボードコンテンツをここに追加可能 */}

      {isLoading ? (
        <p className="mt-8">Loading stories.</p>
      ) : error ? (
        <p className="mt-8 text-red-500">{error}</p>
      ) : (
        <StoryList stories={stories} onDeleteStory={handleDeleteStory} />
      )}

    </div>
  );
};

export default Dashboard;