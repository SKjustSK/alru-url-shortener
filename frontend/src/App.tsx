import { BrowserRouter, Routes, Route, Link } from 'react-router-dom';
import { ArrowRight, Link2 } from 'lucide-react';

import Login from './pages/Login';
import Register from './pages/Register';
import Dashboard from './pages/Dashboard';
import Analytics from './pages/Analytics';

// Custom inline SVG for GitHub since Lucide removed brand icons
const GithubIcon = ({ size = 24, className = "" }) => (
  <svg
    xmlns="http://www.w3.org/2000/svg"
    width={size}
    height={size}
    viewBox="0 0 24 24"
    fill="none"
    stroke="currentColor"
    strokeWidth="2"
    strokeLinecap="round"
    strokeLinejoin="round"
    className={className}
  >
    <path d="M15 22v-4a4.8 4.8 0 0 0-1-3.03c3.18-.3 6.5-1.5 6.5-7.07a4.6 4.6 0 0 0-1.3-3.2 4.2 4.2 0 0 0-.1-3.2s-1.1-.3-3.5 1.3a12.3 12.3 0 0 0-6.2 0C6.5 2.8 5.4 3.1 5.4 3.1a4.2 4.2 0 0 0-.1 3.2 4.6 4.6 0 0 0-1.3 3.2c0 5.6 3.3 6.8 6.5 7.07a4.8 4.8 0 0 0-1 3.03V22"></path>
    <path d="M9 18c-3.14 1.5-4.64-1.5-6-2"></path>
  </svg>
);

const Landing = () => {
  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center p-4">
      
      {/* GitHub Link in Top Right */}
      <div className="absolute top-6 right-6 sm:top-8 sm:right-8">
        <a 
          href="https://github.com/SKjustSK/alru-url-shortener" 
          target="_blank" 
          rel="noopener noreferrer"
          className="text-gray-400 hover:text-black transition-colors"
          title="View source on GitHub"
        >
          <GithubIcon size={24} />
        </a>
      </div>

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