import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import apiClient from '../apiClient';
import type { Story } from '../types';
import { PencilIcon, CheckIcon, XMarkIcon } from '@heroicons/react/24/outline';
import { ClockIcon, BookOpenIcon as ReadCountIcon, HashtagIcon } from '@heroicons/react/24/outline';
import { format, parseISO } from 'date-fns';
import ReactMarkdown from 'react-markdown';
import remarkGfm from 'remark-gfm';
import ConfirmModal from './ConfirmModal';

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
  const [isMarkAsReadModalOpen, setIsMarkAsReadModalOpen] = useState(false);
  const [isUndoModalOpen, setIsUndoModalOpen] = useState(false);

  useEffect(() => {
    const fetchStory = async () => {
      try {
        const response = await apiClient.get<StoryDetailData>(`/stories/${id}`);
        setStoryDetail(response.data);
        setEditedTitle(response.data.title);
      } catch (err) {
        setError('ストーリーの読み込みに失敗しました。');
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
      setUpdateError('タイトルを入力してください。');
      return;
    }

    setUpdateError('');

    try {
      const response = await apiClient.patch<StoryDetailData>(`/stories/${storyDetail.id}`, { title: editedTitle });
      setStoryDetail(response.data);
      setIsEditing(false);
    } catch (error) {
      console.error('タイトルの更新に失敗しました:', error);
      // setError('Failed to update title.');
      setUpdateError('タイトルの更新に失敗しました。もう一度お試しください。');
    }
  };

  const handleCancel = () => {
    if (!storyDetail) return;
    setEditedTitle(storyDetail.title);
    setIsEditing(false);
    setUpdateError('');
  };

  const executeMarkAsRead = async () => {
    if (!storyDetail) return;
    setIsSubmitting(true);
    try {
      await apiClient.post(`/stories/${storyDetail.id}/read`);
      setStoryDetail(prev => prev ? { ...prev, read_count: prev.read_count + 1 } : null);
      // alertは削除
    } catch (error) {
      console.error('既読マークの登録に失敗しました:', error);
      // TODO: トースト通知でエラーを表示
    } finally {
      setIsSubmitting(false);
    }
  };

  const executeUndoRead = async () => {
    if (!storyDetail) return;
    setIsSubmitting(true);
    try {
      await apiClient.delete(`/stories/${storyDetail.id}/read/latest`);
      setStoryDetail(prev => prev ? { ...prev, read_count: Math.max(0, prev.read_count - 1) } : null);
      // alertは削除
    } catch (err) {
      console.error('最後の既読マークの取り消しに失敗しました:', err);
      // TODO: トースト通知でエラーを表示
    } finally {
      setIsSubmitting(false);
    }
  };
  
  // --- ボタンクリック時に実行されるハンドラ ---
  const handleMarkAsReadClick = () => {
    if (!storyDetail) return;
    // 初めて読むときは確認なしで実行
    if (storyDetail.read_count === 0) {
      executeMarkAsRead();
    } else {
      // 2回目以降はモーダルを開く
      setIsMarkAsReadModalOpen(true);
    }
    // setIsMarkAsReadModalOpen(true);
  };

  const handleUndoReadClick = () => {
    setIsUndoModalOpen(true);
  };

  if (isLoading) return <p>文章を読み込んでいます...</p>;
  if (error) return <p className="text-red-500">{error}</p>;

  return (
    <>
      <div className="bg-white p-1 md:p-8 rounded-lg shadow-md w-full max-w-full mx-auto">
        <Link to="/" className="block text-blue-500 text-2xl hover:underline mt-3 mb-4 md:mt-0">&larr; Dashboardに戻る</Link>

        <div className="my-4 pt-6">
          {isEditing ? (
            // --- 編集モード ---
            <div className="flex items-center space-x-2">
              <input
                type="text"
                value={editedTitle}
                onChange={(e) => setEditedTitle(e.target.value)}
                className="text-3xl font-bold border-b-2 border-blue-500 focus:outline-none w-full"
                autoFocus
              />
              <button onClick={handleSave} className="p-2 text-green-600 hover:bg-green-100 rounded-full"><CheckIcon className="h-6 w-6" /></button>
              <button onClick={handleCancel} className="p-2 text-red-600 hover:bg-red-100 rounded-full"><XMarkIcon className="h-6 w-6" /></button>
            </div>
          ) : (
            // --- 表示モード ---
            <div className="flex justify-between items-center">
              <h1 className="text-3xl font-bold">{storyDetail?.title}</h1>
              <button 
                onClick={() => setIsEditing(true)} 
                className="p-2 text-gray-500 hover:text-blue-600 hover:bg-blue-100 rounded-full transition-colors"
                aria-label="Edit title"
              >
                <PencilIcon className="h-6 w-6" />
              </button>
            </div>
          )}
        </div>

        {updateError && (
          <div className="bg-red-100 border border-red-400 text-red-700 text-xl px-4 py-3 rounded mb-4">
            {updateError}
          </div>
        )}

        <div className="flex justify-between items-center mb-6 border-t border-b py-3">
          {/* 作成日 */}
          <div className="flex items-center">
            <ClockIcon className="h-4 w-4 mr-1.5" />
            <span>
              {storyDetail ? format(parseISO(storyDetail.created_at), 'MMM d, yyyy') : ''}
            </span>
          </div>
          
          {/* 単語数 */}
          <div className="flex items-center">
            <HashtagIcon className="h-4 w-4 mr-1.5" />
            <span className="font-semibold">
              {storyDetail?.word_count.toLocaleString()} 単語
            </span>
          </div>

          {/* 読んだ回数 */}
          <div className="flex items-center">
            <ReadCountIcon className="h-4 w-4 mr-1.5" />
            <span className="font-semibold">
              読んだ回数: {storyDetail?.read_count || 0} 回
            </span>
          </div>
        </div>

        <div className="prose max-w-none text-sm md:text-base lg:text-lg mb-6 pt-6 text-left px-2 md:px-0">
          <ReactMarkdown remarkPlugins={[remarkGfm]}>
            {storyDetail?.content || ''}
          </ReactMarkdown>
        </div>

        {/* 読了マークボタン */}
        <div className="pt-6 mt-6">
          <button
            onClick={handleMarkAsReadClick}
            disabled={isSubmitting}
            className="w-full flex items-center justify-center py-2 px-4 bg-green-500 text-white font-bold rounded-lg shadow-md hover:bg-green-600"
          >
            <CheckIcon className="h-5 w-5 mr-2" />
            {isSubmitting ? '記録中...' : '読了を記録'}
          </button>

          {storyDetail && storyDetail.read_count > 0 && (
            <button
              onClick={handleUndoReadClick}
              disabled={isSubmitting}
              className="w-full p-2 bg-gray-200 hover:bg-gray-300 rounded-full mt-6 font-bold"
              title="直近の読了を取り消す"
            >
              直近の読了を取り消す
            </button>
          )}
        </div>
      </div>

      {/* --- モーダルコンポーネントを配置 --- */}
        <ConfirmModal
          isOpen={isMarkAsReadModalOpen}
          onClose={() => setIsMarkAsReadModalOpen(false)}
          onConfirm={executeMarkAsRead}
          title="再度読了を記録"
          message="この文章は以前に読了しています。再度記録しますか？"
          confirmText="はい、記録します"
          intent="success"
        />

        <ConfirmModal
          isOpen={isUndoModalOpen}
          onClose={() => setIsUndoModalOpen(false)}
          onConfirm={executeUndoRead}
          title="直近の読了を取り消す"
          message="直近の読了記録を削除してもよろしいですか？この操作は元に戻せません。"
          confirmText="はい、取り消します"
          intent="warning"
        />
    </>
  );
};

export default StoryDetail;
