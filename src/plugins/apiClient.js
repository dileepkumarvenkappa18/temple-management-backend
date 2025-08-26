import axios from './axios';
import { endpoints } from './api';

// Strip the leading /api from endpoints in development
const processEndpoints = (endpoints) => {
  if (!import.meta.env.DEV) return endpoints;
  
  const processed = {};
  
  for (const [key, value] of Object.entries(endpoints)) {
    if (typeof value === 'string') {
      processed[key] = value.replace(/^\/api/, '');
    } else if (typeof value === 'function') {
      processed[key] = (...args) => value(...args).replace(/^\/api/, '');
    } else if (typeof value === 'object') {
      processed[key] = processEndpoints(value);
    } else {
      processed[key] = value;
    }
  }
  
  return processed;
};

const apiClient = {
  client: axios,
  endpoints: processEndpoints(endpoints),
  
  // Add convenience methods if needed
  get: (path, config) => axios.get(path, config),
  post: (path, data, config) => axios.post(path, data, config),
  put: (path, data, config) => axios.put(path, data, config),
  delete: (path, config) => axios.delete(path, config)
};

export default apiClient;