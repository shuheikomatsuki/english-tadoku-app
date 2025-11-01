import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import apiClient from '../apiClient';
import type { Story } from '../types';
import { PencilIcon, CheckIcon, XMarkIcon, BookOpenIcon } from '@heroicons/react/24/outline';
import { format, parseISO } from 'date-fns';


const StoryDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [story, setStory] = useState<Story | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [updateError, setUpdateError] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchStory = async () => {
      try {
        const response = await apiClient.get<Story>(`/stories/${id}`);
        setStory(response.data);
        setEditedTitle(response.data.title);
      } catch (err) {
        setError('Failed to load the story.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };

    fetchStory();
  }, [id]);

  const handleSave = async () => {
    if (!story) return;

    if (!editedTitle.trim()) {
      setUpdateError('Title cannot be empty.');
      return;
    }

    setUpdateError('');

    try {
      const response = await apiClient.patch<Story>(`/stories/${story.id}`, { title: editedTitle });
      setStory(response.data);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to update title:', error);
      // setError('Failed to update title.');
      setUpdateError('Failed to update title.');
    }
  };

  const handleCancel = () => {
    if (!story) return;
    setEditedTitle(story.title);
    setIsEditing(false);
    setUpdateError('');
  };

  const handleMarkAsRead = async () => {
    if (!story) return;
    setIsSubmitting(true);
    try {
      await apiClient.post(`/stories/${story.id}/read`);
      alert('Your reading has been recorded!');
    } catch (error) {
      console.error('Failed to mark as read', error);
      alert('Failed to record your reading.');
    } finally {
      setIsSubmitting(false);
    }
  };

  if (isLoading) return <p>Loading story...</p>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <div className="bg-white p-6 rounded-lg shadow-md mx-auto max-w-2xl">
      <Link to="/" className="text-blue-500 text-2xl hover:underline mb-4">&larr; Back to Dashboard</Link>

      <div className="mb-4">
        {isEditing ? (
          // --- 編集モード ---
          <input
            type="text"
            value={editedTitle}
            onChange={(e) => setEditedTitle(e.target.value)}
            className="text-3xl font-bold border-b-2 border-blue-500 focus:outline-none w-full p-2 mb-4"
          />
        ) : (
          // --- 表示モード ---
          <h1 className="text-3xl font-bold">{story?.title}</h1>
        )}

        {isEditing ? (
          <div className="flex space-x-2">
            <button onClick={handleSave} className="bg-blue-500 text-white px-4 py-2 rounded mt-2"><CheckIcon className="h-5 w-5" /></button>
            <button onClick={handleCancel} className="bg-gray-300 text-gray-700 px-4 py-2 rounded"><XMarkIcon className="h-5 w-5" /></button>
          </div>
        ) : (
            <button onClick={() => setIsEditing(true)} className="bg-blue-500 text-white px-4 py-2 rounded mt-2"><PencilIcon className="h-5 w-5" /></button>
        )}

      </div>

      {updateError && (
        <div className="bg-red-100 border border-red-400 text-red-700 text-xl px-4 py-3 rounded mb-4">
          {updateError}
        </div>
      )}

      <p className="text-sm text-gray-500 mb-4">
        Created on: {story ? format(parseISO(story.created_at), 'PPPPp') : ''}
      </p>
      <div className="prose max-w-none mb-6">
        <p className="whitespace-pre-wrap">{story?.content}</p>
      </div>

      {/* 読了マークボタン */}
      <div className="pt-6 mt-6">
        <button
          onClick={handleMarkAsRead}
          disabled={isSubmitting}
          className="w-full flex items-center justify-center py-2 px-4 bg-green-500 text-white font-bold rounded-lg shadow-md hover:bg-green-600"
        >
          <BookOpenIcon className="h-5 w-5 mr-2" />
          {isSubmitting ? 'Recording...' : 'Mark as Read'}
        </button>

      </div>
    </div>
  );
};

export default StoryDetail;

