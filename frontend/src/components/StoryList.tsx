import React from 'react';
import { Link } from 'react-router-dom';
import type { Story } from '../types';
import { TrashIcon } from '@heroicons/react/24/outline';
import { format, parseISO } from 'date-fns';


interface StoryListProps {
  stories: Story[];
  onDeleteStory: (id: number) => void;
  currentPage: number;
  totalPages: number;
  onPageChange: (newPage: number) => void;
  refetchTrigger: number;
}

const StoryList: React.FC<StoryListProps> = ({ stories, onDeleteStory, currentPage, totalPages, onPageChange }) => {
  return (
    <div className="mt-8 p-4 bg-white rounded-lg shadow-md">
      <h2 className="text-2xl font-semibold mb-4">文章の履歴</h2>
      {stories.length === 0 ? (
        <p className="text-gray-500">まだ文章が生成されていません。</p>
      ) : (
        <ul className="space-y-4">
          {stories.map(story => (
            <li key={story.id} className="flex justify-between items-center">
              <Link to={`/stories/${story.id}`} className="flex-grow">
                <h3 className="font-bold text-lg hover:text-blue-600">{story.title}</h3>
                <p className="text-sm text-gray-500">
                  作成日: {format(parseISO(story.created_at), 'PPpp')}
                </p>
              </Link>

              <button
                onClick={() => onDeleteStory(story.id)}
                className="ml-4 p-2 text-gray-500 hover:text-red-600 hover:bg-red-100 rounded-full transition-colors"
                aria-label="文章を削除"
              >
                <TrashIcon className="h-5 w-5" />
              </button>
            </li>
          ))}
        </ul>
      )}

      {totalPages > 1 && (
        <div className="flex justify-center items-center mt-8 space-x-4">
          <button
            onClick={() => onPageChange(currentPage - 1)}
            disabled={currentPage === 1}
            className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          >
            前へ
          </button>

          <button
            onClick={() => onPageChange(currentPage + 1)}
            disabled={currentPage === totalPages}
            className="px-4 py-2 bg-gray-200 rounded disabled:opacity-50"
          >
            次へ
          </button>
        </div>
      )}
    </div>
  );
};

export default StoryList;