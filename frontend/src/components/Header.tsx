import React from 'react';
import { Link, NavLink } from 'react-router-dom';
import { useAuth } from '../contexts/AuthContext';

const Header: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto p-4 flex justify-between items-center">
        <Link to="/" className="text-2xl sm:text-3xl md:text-4xl font-bold text-gray-800">
          English Tadoku App
        </Link>

        <nav className="flex items-center">
          {isAuthenticated ? (
            <>
              <NavLink
                to="/profile"
                className={({ isActive }) =>
                  `font-medium ${isActive ? 'text-blue-600' : 'text-gray-600 hover:text-blue-600'} mx-4 text-xl`
                }
              >
                Profile
              </NavLink>

              <button
                onClick={logout}
                className="bg-red-500 hover:bg-red-600 text-white text-xl font-bold m-4 py-2 px-2 rounded"
              >
                Logout
              </button>
            </>
          ) : (
            <Link
              to="/auth"
              className="bg-blue-500 hover:bg-blue-600 text-white text-xl font-bold py-2 px-4 rounded"
            >
              Login / Sign Up
            </Link>
          )}
        </nav>
      </div>
    </header>
  );
};

export default Header;
