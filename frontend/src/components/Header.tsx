import React from 'react';
import { Link, NavLink } from 'react-router-dom';
import { useAuth } from '../contexts/authContext';

const Header: React.FC = () => {
  const { isAuthenticated, logout } = useAuth();

  return (
    <header className="bg-white shadow-md">
      <div className="container mx-auto p-4 flex justify-between items-center">
        <Link to="/" className="text-xl sm:text-2xl md:text-3xl font-bold text-gray-800">
          Readoku
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
              className="bg-blue-500 text-white rounded whitespace-nowrap
             py-2 px-3 text-sm
             md:py-3 md:px-4 md:text-base
             hover:bg-blue-600 transition"
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
