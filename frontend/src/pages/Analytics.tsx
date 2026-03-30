import { useState, useEffect } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { ArrowLeft, ExternalLink, Activity, LayoutGrid, Globe, MousePointerClick } from 'lucide-react';
import { 
  LineChart, Line, XAxis, YAxis, Tooltip, ResponsiveContainer, 
  BarChart, Bar, PieChart, Pie, Cell 
} from 'recharts';
import api from '../lib/axios';

interface TimeSeriesPoint {
  date: string;
  count: number;
}

interface HourlyPoint {
  date: string;
  hour: string;
  count: number;
}

interface StatItem {
  name: string;
  count: number;
}

interface AnalyticsData {
  short_code: string;
  long_url: string;
  total_clicks: number;
  timeline: TimeSeriesPoint[];
  hourly: HourlyPoint[]; // NEW
  os: StatItem[];
  browsers: StatItem[];
  devices: StatItem[];
  countries: StatItem[];
  referrers: StatItem[];
}

const PIE_COLORS = ['#000000', '#374151', '#6B7280', '#9CA3AF', '#D1D5DB'];

export default function Analytics() {
  const { shortCode } = useParams<{ shortCode: string }>();
  const navigate = useNavigate();
  
  const [data, setData] = useState<AnalyticsData | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState('');
  
  const [selectedDate, setSelectedDate] = useState<string | 'all'>('all');

  const displayBaseUrl = (import.meta.env.VITE_API_URL || 'http://localhost:1323').replace(/^https?:\/\//, '');

  useEffect(() => {
    const fetchAnalytics = async () => {
      try {
        const res = await api.get(`/api/links/${shortCode}/analytics`);
        setData(res.data);
      } catch (err: any) {
        if (err.response?.status === 401) {
          localStorage.removeItem('alru_token');
          navigate('/login');
        } else {
          setError(err.response?.data?.error || 'Failed to load analytics');
        }
      } finally {
        setIsLoading(false);
      }
    };

    fetchAnalytics();
  }, [shortCode, navigate]);

  if (isLoading) {
    return <div className="min-h-screen flex items-center justify-center text-gray-500">Loading analytics...</div>;
  }

  if (error || !data) {
    return (
      <div className="min-h-screen flex flex-col items-center justify-center text-center px-4">
        <p className="text-red-600 mb-4">{error || 'Data not found'}</p>
        <Link to="/dashboard" className="text-black font-medium hover:underline flex items-center gap-2">
          <ArrowLeft size={16} /> Back to Dashboard
        </Link>
      </div>
    );
  }

  const chartData = selectedDate === 'all' 
    ? data.timeline.map(d => ({ label: d.date, count: d.count }))
    : data.hourly
        .filter(h => h.date === selectedDate)
        .map(h => ({ label: h.hour, count: h.count }));

  const CustomTooltip = ({ active, payload, label }: any) => {
    if (active && payload && payload.length) {
      return (
        <div className="bg-black text-white text-xs py-1.5 px-3 rounded shadow-lg font-medium">
          {`${label}: ${payload[0].value} clicks`}
        </div>
      );
    }
    return null;
  };

  return (
    <div className="min-h-screen bg-gray-50 pb-12">
      <nav className="bg-white border-b border-gray-200 px-4 py-3 sm:px-6 lg:px-8">
        <div className="max-w-6xl mx-auto flex items-center justify-between">
          <Link to="/dashboard" className="flex items-center gap-2 text-gray-500 hover:text-black transition-colors font-medium text-sm">
            <ArrowLeft size={16} /> Dashboard
          </Link>
          <div className="text-sm font-bold tracking-widest text-black">ALRU</div>
        </div>
      </nav>

      <main className="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 mt-8 space-y-6">
        
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4">
          <div className="overflow-hidden">
            <h1 className="text-2xl font-extrabold text-black flex items-center gap-2 mb-1">
              <Activity size={24} /> Analytics Overview
            </h1>
            <div className="flex items-center gap-2 text-sm">
              <span className="font-semibold text-gray-900 bg-gray-100 px-2 py-0.5 rounded">
                {displayBaseUrl}/{data.short_code}
              </span>
              <span className="text-gray-400">→</span>
              <a href={data.long_url} target="_blank" rel="noopener noreferrer" className="text-gray-500 hover:text-black hover:underline truncate flex items-center gap-1 max-w-[200px] sm:max-w-md">
                {data.long_url} <ExternalLink size={14} />
              </a>
            </div>
          </div>
          
          <div className="bg-black text-white px-6 py-3 rounded-lg text-center shrink-0">
            <div className="text-xs uppercase tracking-wider text-gray-300 font-medium mb-1">Total Clicks</div>
            <div className="text-3xl font-black leading-none">{data.total_clicks}</div>
          </div>
        </div>

        {data.total_clicks === 0 ? (
          <div className="text-center py-20 bg-white rounded-xl border border-gray-200 border-dashed">
            <p className="text-gray-500">No clicks recorded yet. Share your link to generate data!</p>
          </div>
        ) : (
          <>
            <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              
              <div className="flex justify-between items-center mb-6">
                <h2 className="text-sm font-bold text-gray-500 uppercase tracking-wider">Click Timeline</h2>
                {data.timeline.length > 0 && (
                  <select 
                    value={selectedDate}
                    onChange={(e) => setSelectedDate(e.target.value)}
                    className="bg-gray-50 border border-gray-200 text-gray-900 text-xs rounded-lg focus:ring-black focus:border-black block py-1.5 px-3 outline-none cursor-pointer transition-colors"
                  >
                    <option value="all">Daily Overview</option>
                    {data.timeline.map(t => (
                      <option key={t.date} value={t.date}>{t.date}</option>
                    ))}
                  </select>
                )}
              </div>

              <div className="h-64 w-full">
                <ResponsiveContainer width="100%" height="100%">
                  <LineChart data={chartData}>
                    {/* UPDATED: dataKey is now 'label' so it handles both dates and hours automatically */}
                    <XAxis dataKey="label" axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: '#9CA3AF' }} dy={10} />
                    <YAxis axisLine={false} tickLine={false} tick={{ fontSize: 12, fill: '#9CA3AF' }} allowDecimals={false} />
                    <Tooltip content={<CustomTooltip />} cursor={{ stroke: '#E5E7EB', strokeWidth: 2 }} />
                    <Line type="monotone" dataKey="count" stroke="#000000" strokeWidth={3} dot={{ r: 4, fill: '#000', strokeWidth: 0 }} activeDot={{ r: 6 }} />
                  </LineChart>
                </ResponsiveContainer>
              </div>
            </div>

            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
              
              <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
                <h2 className="text-sm font-bold text-gray-500 uppercase tracking-wider mb-6 flex items-center gap-2">
                  <LayoutGrid size={16} /> Devices
                </h2>
                <div className="h-48 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <PieChart>
                      <Pie data={data.devices} cx="50%" cy="50%" innerRadius={50} outerRadius={70} paddingAngle={2} dataKey="count">
                        {data.devices.map((entry, index) => (
                          <Cell key={`cell-${index}`} fill={PIE_COLORS[index % PIE_COLORS.length]} />
                        ))}
                      </Pie>
                      <Tooltip content={<CustomTooltip />} />
                    </PieChart>
                  </ResponsiveContainer>
                </div>
                <div className="mt-4 space-y-2">
                  {data.devices.slice(0, 3).map((item, i) => (
                    <div key={item.name} className="flex justify-between text-sm">
                      <span className="flex items-center gap-2"><div className="w-2 h-2 rounded-full" style={{ backgroundColor: PIE_COLORS[i % PIE_COLORS.length]}}></div>{item.name}</span>
                      <span className="font-medium">{item.count}</span>
                    </div>
                  ))}
                </div>
              </div>

              <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6 lg:col-span-2">
                <h2 className="text-sm font-bold text-gray-500 uppercase tracking-wider mb-6 flex items-center gap-2">
                  <MousePointerClick size={16} /> Top Referrers
                </h2>
                <div className="h-48 w-full">
                  <ResponsiveContainer width="100%" height="100%">
                    <BarChart data={data.referrers.slice(0, 5)} layout="vertical" margin={{ top: 0, right: 0, bottom: 0, left: 0 }}>
                      <XAxis type="number" hide />
                      <YAxis dataKey="name" type="category" axisLine={false} tickLine={false} width={100} tick={{ fontSize: 12, fill: '#4B5563' }} />
                      <Tooltip content={<CustomTooltip />} cursor={{ fill: '#F3F4F6' }} />
                      <Bar dataKey="count" fill="#000000" radius={[0, 4, 4, 0]} barSize={24} />
                    </BarChart>
                  </ResponsiveContainer>
                </div>
              </div>

              <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
                <h2 className="text-sm font-bold text-gray-500 uppercase tracking-wider mb-6 flex items-center gap-2">
                  <Globe size={16} /> Top Locations
                </h2>
                <div className="space-y-4">
                  {data.countries.length === 0 ? (
                    <p className="text-sm text-gray-400">No location data</p>
                  ) : (
                    data.countries.slice(0, 5).map((country) => (
                      <div key={country.name} className="flex items-center justify-between text-sm">
                        <span className="font-medium text-black">{country.name}</span>
                        <div className="flex items-center gap-3">
                          <div className="w-16 h-1.5 bg-gray-100 rounded-full overflow-hidden">
                            <div className="h-full bg-black rounded-full" style={{ width: `${(country.count / data.total_clicks) * 100}%` }}></div>
                          </div>
                          <span className="text-gray-500 w-6 text-right">{country.count}</span>
                        </div>
                      </div>
                    ))
                  )}
                </div>
              </div>

            </div>
          </>
        )}
      </main>
    </div>
  );
}