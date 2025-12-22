import SignUp from './SignUp';
import Login from './Login';

const AuthPage: React.FC = () => {
  const env = import.meta.env.VITE_ENV || 'local';
  const testEmail = import.meta.env.VITE_TEST_EMAIL || '';
  const testPassword = import.meta.env.VITE_TEST_PASSWORD || '';

  const isProdOrLocal = env === 'production' || env === 'local';

  return (
    <div className="w-full max-w-4xl mx-auto space-y-8">
      {isProdOrLocal ? (
        <>
          <SignUp />
          <Login />
        </>
      ) : (
        <div className="bg-white p-6 rounded-lg shadow space-y-4">
          <h2 className="text-2xl font-bold text-center">Login (Test Environment)</h2>
          <p className="text-center text-gray-500 mb-2">
            Use preset test credentials below:
            なんかログインにめっちゃ時間かかる時やログイン失敗する時ある。原因調査中。
          </p>
          <div className="text-center text-sm bg-gray-100 p-3 rounded">
            <p>Email: <code>{testEmail}</code></p>
            <p>Password: <code>{testPassword}</code></p>
          </div>
          <Login defaultEmail={testEmail} defaultPassword={testPassword} />
        </div>
      )}
    </div>
  );
};

export default AuthPage;
