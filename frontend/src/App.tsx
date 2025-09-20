import './App.css'
import SignUp from './components/SignUp';
import Login from './components/Login';

function App() {
  return (
    <div className="min-h-screen bg-gray-100 py-10">
      
      <div className="container mx-auto px-4 max-w-lg"> 
        
        <h1 className="text-4xl font-bold text-gray-800 mb-8 text-center">
          English Tadoku App
        </h1>
        
        <div className="space-y-8">
          <SignUp />
          <Login />
        </div>

      </div>
    </div>
  );
}

export default App
