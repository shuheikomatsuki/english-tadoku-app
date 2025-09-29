import React from 'react';

const Footer: React.FC = () => {
  return (
    <footer className="bg-gray-800 text-white py-4">
      <div className="container mx-auto text-center text-sm">
        <p>&copy; {new Date().getFullYear()} English Tadoku App. All Rights Reserved.</p>
      </div>
    </footer>
  );
};

export default Footer;