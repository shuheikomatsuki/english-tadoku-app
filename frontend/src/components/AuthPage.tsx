import SignUp from './SignUp';
import Login from './Login';

const AuthPage: React.FC = () => {
    return (
        <div className="w-full max-w-4xl mx-auto space-y-8">
            <SignUp />
            <Login />
        </div>
    );
};

export default AuthPage;