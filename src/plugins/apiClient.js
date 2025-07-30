// src/plugins/apiClient.js
import api from '@/services/api';
import superAdminService from '@/services/superadmin.service';

// Create and export a composite API client
const apiClient = {
  ...api,
  admin: superAdminService
};

export default apiClient;