import React from 'react';
import { Link } from 'react-router-dom';
import type { Story } from '../types';


interface StoryListProps {
  stories: Story[];
}

const StoryList: React.FC<StoryListProps> = ({ stories }) => {
  return (
    <div className="mt-8">
      <h2 className="text-2xl font-semibold mb-4">Your Stories</h2>
      {stories.length === 0 ? (
        <p className="text-gray-500">You haven't generated any stories yet.</p>
      ) : (
        <ul className="space-y-4">
          {stories.map(story => (
            <li key={story.id}>
              <Link to={`/stories/${story.id}`} className="block bg-white p-4 rounded-lg shadow transition hover:shadow-lg">
                <h3 className="font-bold text-lg">{story.title}</h3>
                <p className="text-sm text-gray-500">
                  Created on: {new Date(story.createdAt).toLocaleString()}
                </p>
              </Link>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

export default StoryList;