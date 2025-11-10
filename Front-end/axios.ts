import axios from 'axios';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_AXIOS_URL,
});

axiosInstance.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token');
    if (token && token.trim() !== '' && token !== "null") {
      config.headers['Authorization'] = `Bearer ${token}`;
    }
    return config;
  },
  (error) => Promise.reject(error)
);

export default axiosInstance;
