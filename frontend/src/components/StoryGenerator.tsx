import React, { useState, useEffect } from 'react';
import apiClient from '../apiClient';
import { SparklesIcon, ArrowPathIcon } from '@heroicons/react/24/outline';
import type { Story } from '../types';
import { useNavigate } from 'react-router-dom';
import axios from 'axios';


interface StoryGeneratorProps {
  onStoryGenerated: (story: Story) => void;
}

interface GenerationStatus {
  current_count: number;
  limit: number;
}

const StoryGenerator: React.FC<StoryGeneratorProps> = ( { onStoryGenerated }) => {
  const navigate = useNavigate();
  const [prompt, setPrompt] = useState('');
  const [generatedStory, setGeneratedStory] = useState<Story | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');
  const [generationStatus, setGenerationStatus] = useState<GenerationStatus | null>(null);

  useEffect(() => {
    const fetchGenerationStatus = async () => {
      try {
        const response = await apiClient.get<GenerationStatus>('/users/me/generation-status');
        setGenerationStatus(response.data);
      } catch (err) {
        console.error('Failed to fetch generation status:', err);
      }
    };

    fetchGenerationStatus();
  }, []);

  const handleGenerate = async () => {
    if (!prompt.trim()) {
      setError('Please enter a prompt.');
      return;
    }

    setIsLoading(true);
    setError('');
    setGeneratedStory(null);

    try {
      const response = await apiClient.post<Story>('/stories', { prompt });
      const newStory = response.data;
      setGeneratedStory(response.data);
      onStoryGenerated(response.data);

      if (generationStatus) {
        setGenerationStatus({
          ...generationStatus,
          current_count: generationStatus.current_count + 1,
        });
      }

      setPrompt('');
      navigate(`/stories/${newStory.id}`);
    } catch (err) {
      if (axios.isAxiosError(err) && err.response) {
        if (err.response.status === 429) {
          setError('You have reached your daily story generation limit. Please try again tomorrow.');
        } else {
          setError('Failed to generate story. Please try again.');
        }
      } else {
        setError('An unexpected error occurred.');
      }
      console.error(err);
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="">
      {/* --- ストーリー生成フォーム --- */}
      <div className="bg-white p-6 rounded-lg shadow-md mb-6">
        <h2 className="text-xl font-bold mb-4">Generate a New Story</h2>
        <textarea
          className="w-full p-2 mb-4 border border-gray-300 rounded-md"
          rows={3}
          placeholder="Enter a prompt for your story..."
          value={prompt}
          onChange={(e) => setPrompt(e.target.value)}
          disabled={isLoading}
        />

        {/* --- 生成回数の表示 --- */}
        {generationStatus && (
          <div className="text-lg text-right mb-2">
            Today's generations: {generationStatus.current_count} / {generationStatus.limit}
          </div>
        )}

        {/* --- 制限到達時のメッセージ --- */}
        {generationStatus && generationStatus.current_count >= generationStatus.limit && (
          <div className="bg-yellow-100 border border-yellow-400 text-yellow-700 px-4 py-3 rounded text-center" role="alert">
            <span className="block sm:inline">You have reached your daily generation limit. Please try again tomorrow.</span>
          </div>
        )}

        <button
          className="w-full mt-4 py-3 bg-black text-white rounded-lg hover:bg-gray-800 font-bold flex justify-center"
          onClick={handleGenerate}
          // disabled={isLoading}
          disabled={isLoading || (generationStatus ? generationStatus.current_count >= generationStatus.limit : false)}
        >
          {isLoading ? (
            // --- ローディング中の表示 ---
            <>
              <ArrowPathIcon className="h-5 w-5 mr-3 animate-spin" />
              <span>Generating...</span>
            </>
          ) : (
            // --- 通常時の表示 (変更なし) ---
            <>
              <SparklesIcon className="h-5 w-5 mr-2" />
              <span>Generate Story</span>
            </>
          )}
        </button>
      </div>

      {/* --- エラーメッセージ表示 --- */}
      {error && (
        <div className="bg-red-100 border border-red-400 px-4 py-3 rounded mb-6" role="alert">
          <span className="">{error}</span>
        </div>
      )}

    </div>
  );
};

export default StoryGenerator;