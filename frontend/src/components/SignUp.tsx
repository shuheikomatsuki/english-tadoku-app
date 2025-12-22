import React, { useState } from "react";
import apiClient from "../apiClient";
import { Eye, EyeOff } from 'lucide-react';
import axios from 'axios';

interface ErrorResponse {
  error: string;
}

const SignUp: React.FC = () => {
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [message, setMessage] = useState("");
  const [showPassword, setShowPassword] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setMessage("");

    try {
      const response = await apiClient.post("/signup", { email, password });
      setMessage("ユーザーが正常に作成されました！");
      console.log(response.data);
      setEmail("");
      setPassword("");
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        const errorData = error.response.data as ErrorResponse;
        setMessage(errorData.error || "予期しないエラーが発生しました。");
      } else {
        setMessage("ユーザーの作成に失敗しました。もう一度お試しください。");
      }
      console.error(error);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow-md">
      <h2 className="text-2xl font-bold mb-5 text-center">Sign Up</h2>
      <form onSubmit={handleSubmit} className="space-y-6">
        <div>
          <label className="block font-bold mb-2" htmlFor="signup-email">
            Email
          </label>
          <input
            className="border rounded w-full px-3 py-2 focus:outline-none focus:ring-2 focus:ring-gray-500"
            id="signup-email"
            type="email"
            placeholder="Email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            required
          />
        </div>
        <div>
          <label className="block font-bold mb-2" htmlFor="signup-password">
            Password
          </label>
          <div className="relative">
            <input
              className="border rounded w-full px-3 py-2 pr-10 focus:outline-none focus:ring-2 focus:ring-gray-500"
              id="signup-password"
              type={showPassword ? "text" : "password"}
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              required
            />
            <button
              type="button"
              onClick={() => setShowPassword(!showPassword)}
              className="absolute right-3 top-1/2 transform -translate-y-1/2 text-gray-600 hover:text-gray-800 focus:outline-none"
              style={{
                background: 'none',
                border: 'none',
                padding: 0,
                margin: 0,
                cursor: 'pointer',
                outline: 'none'
              }}
            >
              {showPassword ? <EyeOff size={20} /> : <Eye size={20} />}
            </button>
          </div>
        </div>
        <div>
          <button
            type="submit"
            className="w-full mt-4 py-2 px-4 bg-black text-white rounded-lg hover:bg-gray-800 transition font-medium"
            style={{
              backgroundColor: '#000000',
              color: '#ffffff',
            }}
          >
            Sign Up
          </button>
        </div>
      </form>
      {message && <p className="mt-4 text-center">{message}</p>}
    </div>
  );
};

export default SignUp;