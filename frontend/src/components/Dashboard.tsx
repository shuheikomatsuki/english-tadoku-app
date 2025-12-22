import React, {useState, useEffect } from 'react';
// import { useAuth } from '../contexts/AuthContext';
import StoryGenerator from './StoryGenerator';
import StoryList from './StoryList';
import apiClient from '../apiClient';
import type { Story } from '../types';
import ConfirmModal from './ConfirmModal';

interface StoriesResponse {
  stories: Story[];
  total_pages: number;
}

const storiesPerPage = 5;

const Dashboard: React.FC = () => {
  // const { logout } = useAuth();

  const [stories, setStories] = useState<Story[]>([]);
  const [currentPage, setCurrentPage] = useState(1);
  const [totalPages, setTotalPages] = useState(0);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  const [isDeleteModalOpen, setIsDeleteModalOpen] = useState(false);
  const [storyToDelete, setStoryToDelete] = useState<number | null>(null);


  useEffect(() => {
    const fetchStories = async () => {
      setIsLoading(true);
      try {
        const response = await apiClient.get<StoriesResponse>(`/stories?page=${currentPage}&limit=${storiesPerPage}`);
        setStories(response.data.stories || []);
        setTotalPages(response.data.total_pages);
      } catch (err) {
        setError('Failed to load stories.');
        console.error(err);
      } finally {
        setIsLoading(false);
      }
    };
    fetchStories();
  }, [currentPage]);

  const handlePageChange = (newPage: number) => {
    setCurrentPage(newPage);
  };

  const handleStoryGenerated = (newStory: Story) => {
    if (currentPage === 1) {
      setStories(prevStories => [newStory, ...prevStories]);
    } else {
      setCurrentPage(1);
    }
  };

  // 削除ボタンが押されたときに呼ばれる関数
  const handleDeleteClick = (id: number) => {
    setStoryToDelete(id); // どのストーリーを削除するかIDを記憶
    setIsDeleteModalOpen(true); // モーダルを開く
  };

  // モーダルの「確認」ボタンが押されたときに呼ばれる関数 (APIを呼び出す)
  const executeDeleteStory = async () => {
    if (storyToDelete === null) return; // 対象がなければ何もしない

    try {
      await apiClient.delete(`/stories/${storyToDelete}`);
      setStories(prevStories => prevStories.filter(story => story.id !== storyToDelete));
      // TODO: 総ページ数が変わる可能性があるので、リストを再取得するのがより堅牢
      // setRefetchTrigger(prev => prev + 1); // のような仕組みを再導入する
    } catch (err) {
      console.error('Failed to delete story:', err);
      // TODO: トースト通知でエラーを表示
    } finally {
      // 処理が終わったら、モーダルを閉じて対象IDをリセット
      setIsDeleteModalOpen(false);
      setStoryToDelete(null);
    }
  };

  return (
    <>
      <div className="w-full max-w-4xl mx-auto">
        <div className="flex justify-between items-center mb-8">
          <h1 className="text-3xl font-bold text-center w-full">Dashboard</h1>
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
          <StoryList 
            stories={stories} 
            onDeleteStory={handleDeleteClick} 
            currentPage={currentPage}
            totalPages={totalPages}
            onPageChange={handlePageChange}
            refetchTrigger={0}
          />
        )}
      </div>

      <ConfirmModal
        isOpen={isDeleteModalOpen}
        onClose={() => setIsDeleteModalOpen(false)}
        onConfirm={executeDeleteStory}
        title="Delete Story"
        message="Are you sure you want to delete this story? This action cannot be undone."
        confirmText="Yes, Delete"
        intent="warning"
      />
    </>
  );
};

export default Dashboard;