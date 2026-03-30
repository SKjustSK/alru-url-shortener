import axios from 'axios';

// Vite exposes env variables via import.meta.env
const baseURL = import.meta.env.VITE_API_URL || 'http://localhost:1323';

const api = axios.create({
  baseURL: baseURL, 
});

api.interceptors.request.use((config) => {
  const token = localStorage.getItem('alru_token');
  if (token && config.headers) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;