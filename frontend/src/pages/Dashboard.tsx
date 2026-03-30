import { useState, useEffect } from 'react';
import { useNavigate, Link } from 'react-router-dom';
import { Link2, LogOut, Copy, BarChart2, Plus, Check } from 'lucide-react';
import api from '../lib/axios';

interface LinkItem {
  short_url: string;
  long_url: string;
  short_code: string;
  created_at: string;
  expires_at: string;
}

export default function Dashboard() {
  const [links, setLinks] = useState<LinkItem[]>([]);
  const [longUrl, setLongUrl] = useState('');
  const [customCode, setCustomCode] = useState('');
  const [isLoading, setIsLoading] = useState(true);
  const [isCreating, setIsCreating] = useState(false);
  const [error, setError] = useState('');
  const [copiedCode, setCopiedCode] = useState<string | null>(null);
  
  const navigate = useNavigate();

  const displayBaseUrl = (import.meta.env.VITE_API_URL || 'http://localhost:1323').replace(/^https?:\/\//, '');

  useEffect(() => {
    const token = localStorage.getItem('alru_token');
    if (!token) {
      navigate('/login');
      return;
    }
    fetchLinks();
  }, [navigate]);

  const fetchLinks = async () => {
    try {
      const res = await api.get('/api/links');
      setLinks(res.data.links || []);
    } catch (err: any) {
      if (err.response?.status === 401) {
        handleLogout();
      } else {
        setError(err.response?.data?.error || 'Failed to fetch links');
      }
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreate = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setIsCreating(true);

    try {
      const res = await api.post('/api/links', {
        long_url: longUrl,
        short_code: customCode.trim() !== '' ? customCode.trim() : undefined
      });

      const newLink: LinkItem = {
        short_url: res.data.short_url,
        long_url: res.data.long_url,
        short_code: res.data.short_url.split('/').pop() || '',
        created_at: new Date().toISOString(),
        expires_at: res.data.expires_on,
      };

      setLinks([newLink, ...links]);
      setLongUrl('');
      setCustomCode('');
    } catch (err: any) {
      setError(err.response?.data?.error || 'Failed to create link');
    } finally {
      setIsCreating(false);
    }
  };

  const handleLogout = () => {
    localStorage.removeItem('alru_token');
    navigate('/login');
  };

  const handleCopy = (url: string, code: string) => {
    navigator.clipboard.writeText(url);
    setCopiedCode(code);
    setTimeout(() => setCopiedCode(null), 2000);
  };

  // Helper function to format date as DD/MM/YY
  const formatDate = (isoString: string) => {
    const date = new Date(isoString);
    const day = String(date.getDate()).padStart(2, '0');
    const month = String(date.getMonth() + 1).padStart(2, '0');
    const year = String(date.getFullYear()).slice(-2);
    return `${day}/${month}/${year}`;
  };

  if (isLoading) {
    return <div className="min-h-screen flex items-center justify-center text-gray-500">Loading your workspace...</div>;
  }

  return (
    <div className="min-h-screen bg-gray-50 pb-12">
      <nav className="bg-white border-b border-gray-200 px-4 py-3 sm:px-6 lg:px-8 flex justify-between items-center">
        <div className="flex items-center gap-2">
          <div className="bg-black text-white p-1.5 rounded-md">
            <Link2 size={20} />
          </div>
          <span className="font-bold text-xl tracking-tight text-black">ALRU</span>
        </div>
        <button 
          onClick={handleLogout}
          className="flex items-center gap-2 text-sm font-medium text-gray-600 hover:text-black transition-colors"
        >
          <LogOut size={16} />
          Sign out
        </button>
      </nav>

      <main className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 mt-8">
        {error && (
          <div className="mb-6 p-4 text-sm text-red-600 bg-red-50 rounded-lg border border-red-100">
            {error}
          </div>
        )}

        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-8">
          <h2 className="text-lg font-bold text-black mb-4 flex items-center gap-2">
            <Plus size={20} /> Create New Link
          </h2>
          
          {/* UPDATED FORM LAYOUT */}
          <form onSubmit={handleCreate} className="flex flex-col gap-4">
            
            {/* Top Row: Long URL Input */}
            <div>
              <input
                type="url"
                required
                placeholder="https://your-long-url.com/something-very-long"
                value={longUrl}
                onChange={(e) => setLongUrl(e.target.value)}
                className="w-full px-4 py-3 border border-gray-300 rounded-lg shadow-sm placeholder-gray-400 focus:outline-none focus:ring-2 focus:ring-black focus:border-transparent transition-all"
              />
            </div>
            
            {/* Bottom Row: Custom Code + Button */}
            <div className="flex flex-col sm:flex-row gap-4">
              
              {/* Custom Code Input with URL Prefix Preview */}
              <div className="flex-grow flex items-center border border-gray-300 rounded-lg shadow-sm focus-within:ring-2 focus-within:ring-black focus-within:border-transparent overflow-hidden transition-all bg-white">
                <span className="pl-4 pr-1 text-gray-400 select-none whitespace-nowrap">
                  {displayBaseUrl}/
                </span>
                <input
                  type="text"
                  placeholder="alias (opt)"
                  value={customCode}
                  onChange={(e) => setCustomCode(e.target.value)}
                  className="w-full py-3 pr-4 outline-none bg-transparent placeholder-gray-400 text-black"
                />
              </div>

              {/* Shorten Button aligned to the right */}
              <button
                type="submit"
                disabled={isCreating}
                className="px-8 py-3 bg-black hover:bg-gray-800 text-white font-medium rounded-lg transition-colors disabled:opacity-50 whitespace-nowrap"
              >
                {isCreating ? 'Shortening...' : 'Shorten URL'}
              </button>
            </div>

          </form>
        </div>

        <div className="space-y-4">
          <h3 className="text-sm font-bold text-gray-500 uppercase tracking-wider mb-4">Your Links</h3>
          
          {links.length === 0 ? (
            <div className="text-center py-12 bg-white rounded-xl border border-gray-200 border-dashed">
              <p className="text-gray-500">You haven't created any links yet.</p>
            </div>
          ) : (
            links.map((link) => (
              <div key={link.short_code} className="bg-white p-5 rounded-xl shadow-sm border border-gray-200 flex flex-col sm:flex-row sm:items-center justify-between gap-4 transition-all hover:border-gray-300">
                <div className="overflow-hidden flex-grow">
                  <div className="flex items-center gap-2 mb-1">
                    <a href={link.short_url} target="_blank" rel="noopener noreferrer" className="font-bold text-black hover:underline text-lg truncate">
                      {link.short_url.replace(/^https?:\/\//, '')}
                    </a>
                  </div>
                  <p className="text-sm text-gray-500 truncate" title={link.long_url}>
                    {link.long_url}
                  </p>
                  <p className="text-xs text-gray-400 mt-2">
                    Created: {formatDate(link.created_at)}
                  </p>
                </div>
                
                <div className="flex items-center gap-2 sm:shrink-0">
                  <button 
                    onClick={() => handleCopy(link.short_url, link.short_code)}
                    className="p-2.5 text-gray-600 hover:text-black hover:bg-gray-100 rounded-lg transition-colors flex items-center justify-center"
                    title="Copy to clipboard"
                  >
                    {copiedCode === link.short_code ? <Check size={18} className="text-green-600" /> : <Copy size={18} />}
                  </button>
                  <Link 
                    to={`/analytics/${link.short_code}`}
                    className="flex items-center gap-2 px-4 py-2.5 bg-gray-100 hover:bg-gray-200 text-black font-medium text-sm rounded-lg transition-colors"
                  >
                    <BarChart2 size={16} />
                    Analytics
                  </Link>
                </div>
              </div>
            ))
          )}
        </div>
      </main>
    </div>
  );
}