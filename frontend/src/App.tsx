import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import { ArrowRight, Link2 } from 'lucide-react';

import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Analytics from './pages/Analytics';

const Landing = () => {
  return (
    <div className="min-h-screen flex flex-col items-center justify-center p-4">
      <div className="text-center max-w-2xl w-full">
        <div className="flex justify-center mb-6">
          <div className="bg-black text-white p-3 rounded-xl">
            <Link2 size={32} />
          </div>
        </div>
        
        <h1 className="text-5xl font-extrabold tracking-tight text-black mb-6">
          ALRU
        </h1>
        <p className="text-xl text-gray-500 mb-10 max-w-md mx-auto">
          A sleek, high-performance tool to shorten links and track analytics without the visual clutter.
        </p>
        
        <div className="flex flex-col sm:flex-row items-center justify-center gap-4">
          <Link 
            to="/register" 
            className="w-full sm:w-auto flex items-center justify-center gap-2 bg-black hover:bg-gray-800 text-white font-medium py-3 px-8 rounded-lg transition-colors duration-200"
          >
            Get Started <ArrowRight size={18} />
          </Link>
          <Link 
            to="/login" 
            className="w-full sm:w-auto flex items-center justify-center gap-2 bg-black hover:bg-gray-800 text-white font-medium py-3 px-8 rounded-lg transition-colors duration-200"
          >
            Sign In
          </Link>
        </div>
      </div>
    </div>
  );
};

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Landing />} />
        <Route path="/login" element={<Login />} />
        <Route path="/register" element={<Register />} />
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/analytics/:shortCode" element={<Analytics />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;