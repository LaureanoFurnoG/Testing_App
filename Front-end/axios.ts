import axios from 'axios';
import { useAuth } from './auth/AuthProvider';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_AXIOS_URL
});
axiosInstance.defaults.withCredentials = true; //send cookies

export const setupInterceptors = (auth: ReturnType<typeof useAuth>) => {
  axiosInstance.interceptors.request.use(config => {
    const stored = sessionStorage.getItem('Token');
    if (stored) {
      const token = JSON.parse(stored);
      config.headers.Authorization = `Bearer ${token.access_token}`;
    }
    return config;
  });

  axiosInstance.interceptors.response.use(response => {
    const authHeader = response.headers['authorization'];

    if (authHeader?.startsWith('Bearer ')) {
    auth.setToken(prev => {
      if (!prev) return prev;
      return {
        ...prev,
        access_token: authHeader.slice(7),
      };
    });
    }

    return response;
  }
  ,
  error => {
    if (error.response?.status === 401) {
      auth.logout();
    }
    return Promise.reject(error);
  }
  );

};

export default axiosInstance;
