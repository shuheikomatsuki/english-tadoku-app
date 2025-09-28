import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import apiClient from '../apiClient';
import type { Story } from '../types';


const StoryDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [story, setStory] = useState<Story | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchStory = async () => {
      try {
        const response = await apiClient.get<Story>(`/stories/${id}`);
        setStory(response.data);
      } catch (err) {
        setError('Failed to load the story.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStory();
  }, [id]);

  if (isLoading) return <p>Loading story...</p>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div className="bg-white p-6 rounded-lg shadow-md mx-auto max-w-2xl">
      <Link to="/" className="text-blue-500 text-2xl hover:underline mb-4">Back to Dashboard</Link>

      <h1 className="text-2xl font-bold p-4">{story?.title}</h1>
      <p className="text-sm text-gray-500 p-4">
        Created on: {story ? new Date(story.created_at).toLocaleString() : ''}
      </p>
      <div className="">
        <p className="">{story?.content}</p>
      </div>
    </div>
  );
};

export default StoryDetail;

