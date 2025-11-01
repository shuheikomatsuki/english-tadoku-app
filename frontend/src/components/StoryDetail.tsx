import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import apiClient from '../apiClient';
import type { Story } from '../types';
import { PencilIcon, CheckIcon, XMarkIcon, BookOpenIcon } from '@heroicons/react/24/outline';
import { format, parseISO } from 'date-fns';

interface StoryDetailData extends Story {
  read_count: number;
}


const StoryDetail: React.FC = () => {
  const { id } = useParams<{ id: string }>();
  const [storyDetail, setStoryDetail] = useState<StoryDetailData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [updateError, setUpdateError] = useState('');
  const [isEditing, setIsEditing] = useState(false);
  const [editedTitle, setEditedTitle] = useState('');
  const [isSubmitting, setIsSubmitting] = useState(false);

  useEffect(() => {
    const fetchStory = async () => {
      try {
        const response = await apiClient.get<StoryDetailData>(`/stories/${id}`);
        setStoryDetail(response.data);
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
    if (!storyDetail) return;

    if (!editedTitle.trim()) {
      setUpdateError('Title cannot be empty.');
      return;
    }

    setUpdateError('');

    try {
      const response = await apiClient.patch<StoryDetailData>(`/stories/${storyDetail.id}`, { title: editedTitle });
      setStoryDetail(response.data);
      setIsEditing(false);
    } catch (error) {
      console.error('Failed to update title:', error);
      // setError('Failed to update title.');
      setUpdateError('Failed to update title.');
    }
  };

  const handleCancel = () => {
    if (!storyDetail) return;
    setEditedTitle(storyDetail.title);
    setIsEditing(false);
    setUpdateError('');
  };

  const handleMarkAsRead = async () => {
    if (!storyDetail) return;
    setIsSubmitting(true);
    try {
      await apiClient.post(`/stories/${storyDetail.id}/read`);
      setStoryDetail(prev => prev ? { ...prev, read_count: prev.read_count + 1} : null);
      alert('Your reading has been recorded!');
    } catch (error) {
      console.error('Failed to mark as read', error);
      alert('Failed to record your reading.');
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUndoRead = async () => {
    if (!storyDetail) return;
    setIsSubmitting(true);
    try {
      await apiClient.delete(`/stories/${storyDetail.id}/read/latest`);
      setStoryDetail(prev => prev ? { ...prev, read_count: Math.max(0, prev.read_count - 1) } : null);
      alert('Last reading record has been removed.');
    } catch (err) {
      console.error('Failed to undo last read', err);
      alert('Failed to undo last reading record.');
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
          <h1 className="text-3xl font-bold">{storyDetail?.title}</h1>
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

      <div className="flex justify-between items-center text-sm text-gray-500 mb-4 pt-4">
        <p>
          Created on: {storyDetail ? format(parseISO(storyDetail.created_at), 'PPPPp') : ''}
        </p>
        <p className="font-bold">
          Read Count: {storyDetail?.read_count || 0}
        </p>
      </div>

      <div className="prose max-w-none mb-6 pt-10">
        <p className="whitespace-pre-wrap">{storyDetail?.content}</p>
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

        {storyDetail && storyDetail.read_count > 0 && (
          <button
            onClick={handleUndoRead}
            disabled={isSubmitting}
            className="w-full p-2 bg-gray-200 hover:bg-gray-300 rounded-full mt-6 font-bold"
            title="Undo last read"
          >
            Undo Last Read
          </button>
        )}
      </div>
    </div>
  );
};

export default StoryDetail;

