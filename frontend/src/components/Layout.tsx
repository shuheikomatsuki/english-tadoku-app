import React from 'react';
import { Outlet } from 'react-router-dom';
import Header from './Header';
import Footer from './Footer';

const Layout: React.FC = () => {
  return (
    <div className="flex flex-col min-h-screen bg-gray-100">
      <Header />

      <main className="flex-grow mx-auto px-0 md:px-4 py-6 max-w-screen-lg">
        <Outlet />
      </main>

      <Footer />
    </div>
  );
};

export default Layout;