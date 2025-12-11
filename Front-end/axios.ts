import axios from 'axios';
import { useAuth } from './auth/AuthProvider';

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_AXIOS_URL
});

export const setupInterceptors = (auth: ReturnType<typeof useAuth>) => {
  axiosInstance.interceptors.request.use(config => {
    if (auth.token?.access_token) {
      config.headers['Authorization'] = `Bearer ${auth.token.access_token}`;
      config.headers['Refresh-Token'] = auth.token.refresh_token;
    }
    return config;
  });

  axiosInstance.interceptors.response.use(
    response => {
      const newAccess = response.headers['authorization'];
      const newRefresh = response.headers['refresh-token'];

      if (newAccess) {
        auth.setToken({
          access_token: newAccess.replace('Bearer ', ''),
          refresh_token: newRefresh ?? auth.token?.refresh_token,
          expires_in: auth.token?.expires_in,
          profile: auth.token?.profile
        });
      }

      return response;
    },
    error => {
      if (error.response?.status === 401) {
        auth.logout();
      }
      return Promise.reject(error);
    }
  );

};

export default axiosInstance;
