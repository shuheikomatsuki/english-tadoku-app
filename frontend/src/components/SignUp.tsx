import React, { useState } from 'react';
import apiClient from '../apiClient';

const SignUp: React.FC = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');
    const [message, setMessage] = useState('');

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        setMessage('');

        try {
            const response = await apiClient.post('/signup', { email, password });
            setMessage('User created successfully!');
            console.log(response.data);
            setEmail('');
            setPassword('');

        } catch (error) {
            setMessage('Failed to create user. Please try again.');
            console.error(error);
        }
    }

    return (
        <div　className="max-w-md mx-auto mt-10 p-6">
            <h2 className="text-2xl font-bold mb-5 text-center">Sign Up</h2>
            <form onSubmit={handleSubmit}>
                <div className="mb-4">
                    <label className="block font-bold mb-2 px-2" htmlFor="email">
                        Email
                    </label>
                    <input 
                        className="border rounded w-full px-3 py-2 focus:outline-none"
                        id="email" 
                        type="email" 
                        placeholder="Email" 
                        value={email} 
                        onChange={(e) => setEmail(e.target.value)} 
                        required 
                    />
                </div>
                <div className="mb-6">
                    <label className="block font-bold mb-2 px-2" htmlFor="password">
                        Password
                    </label>
                    <input 
                        className="border rounded w-full px-3 py-2 focus:outline-none"
                        id="password" 
                        type="password" 
                        placeholder="Password" 
                        value={password} 
                        onChange={(e) => setPassword(e.target.value)} 
                        required 
                    />
                </div>
                <div className="flex items-center justify-center">
                    <button className="bg-blue-500" type="submit">
                    <span className="text-xl font-bold">
                        Sign Up
                    </span>
                </button>
                </div>
            </form>
            {message && <p className="mt-4 text-center">{message}</p>}
        </div>
    );
};

export default SignUp;