import './App.css'
import SignUp from './components/SignUp';

function App() {
  // return (
  //   <div className="p-6 bg-blue-500 text-white rounded-xl">
  //     <h1 className="text-2xl font-bold">Hello Tailwind v4 + Vite!</h1>
  //   </div>
  // )
  return (
    <div className="min-h-screen bg-gray-100 flex flex-col items-center pt-10">
      <h1 className="text-4xl font-bold text-gray-800 mb-8">
        English Tadoku App
      </h1>
      <SignUp />
    </div>
  );
}

export default App
