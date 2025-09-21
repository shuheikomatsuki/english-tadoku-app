import React, { useState } from 'react';
import apiClient from '../apiClient';
import { SparklesIcon } from '@heroicons/react/24/outline';

interface Story {
  id: number;
  userId: number;
  title: string;
  content: string;
  createdAt: string;
  updatedAt: string;
}

const StoryGenerator: React.FC = () => {
  const [prompt, setPrompt] = useState('');
  const [generatedStory, setGeneratedStory] = useState<Story | null>(null);
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState('');

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
      setGeneratedStory(response.data);
    } catch (err) {
      setError('Failed to generate story. Please try again.');
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
        <button
          className="w-full mt-4 py-3 bg-black text-white rounded-lg hover:bg-gray-800 font-bold flex justify-center"
          onClick={handleGenerate}
          disabled={isLoading}
        >
          <SparklesIcon className="h-5 w-5 mr-2" />
          {isLoading ? 'Generating...' : 'Generate Story'}
        </button>
      </div>

      {/* --- エラーメッセージ表示 --- */}
      {error && (
        <div className="bg-red-100 border border-red-400 px-4 py-3 rounded mb-6" role="alert">
          <span className="">{error}</span>
        </div>
      )}

      {/* --- テスト用エラーメッセージ --- */}
      <div className="bg-red-100 border border-red-400 px-4 py-3 rounded mb-6" role="alert">
        <p>This is a sample error message.</p>
      </div>
      
      {/* --- 生成されたストーリー表示 --- */}
      {generatedStory && (
        <div className="bg-white p-6 rounded-lg shadow-md animate-fade-in">
          <h3 className="text-lg font-bold p-2">{generatedStory.title}</h3>
          <p className="whitespace-pre-wrap p-2">{generatedStory.content}</p>
        </div>
      )}

      {/* --- テスト用生成されたストーリー表示 --- */}
      <div className="bg-white p-6 rounded-lg shadow-md animate-fade-in">
        <h3 className="text-lg font-bold p-2">The Mysterious Forest</h3>
        <p className="whitespace-pre-wrap p-2">
          Once upon a time, in a land far away, there was a mysterious forest that no one dared to enter. The trees were tall and twisted, and strange noises echoed through the air. One day, a brave young adventurer decided to explore the forest and uncover its secrets.
        </p>
      </div>
    </div>
  );
};

export default StoryGenerator;