// src/services/devotee.service.js
import api from '@/plugins/axios';
import router from '@/router';

// Get the auth token from localStorage with validation
const getAuthToken = () => {
  const token = localStorage.getItem('temple_auth_token') || localStorage.getItem('auth_token');
  if (!token) {
    console.log('No auth token found in localStorage');
    return null;
  }
  return token;
};

// Check if user is authenticated
const isAuthenticated = () => {
  const token = getAuthToken();
  return token !== null && token !== undefined && token !== '';
};

// Clear authentication data
const clearAuth = () => {
  localStorage.removeItem('temple_auth_token');
  localStorage.removeItem('temple_user_data');
  localStorage.removeItem('temple_user_role');
  // Clear legacy keys too
  localStorage.removeItem('auth_token');
  localStorage.removeItem('user_data');
};

// Redirect to login
const redirectToLogin = () => {
  clearAuth();
  if (router.currentRoute.value.path !== '/login') {
    router.push('/login');
  }
};

export default {
  /**
   * Get devotee profile (current user)
   */
  async getProfile() {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get('/v1/profiles/me');
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get devotee by ID (for admin purposes)
   */
  async getDevoteeById(id) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get(`/v1/profiles/${id}`);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Create or update devotee profile
   */
  async createOrUpdateProfile(profileData) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }
      
      let response;
      
      if (profileData instanceof FormData) {
        response = await api.post('/v1/profiles', profileData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
      } else {
        response = await api.post('/v1/profiles', profileData);
      }
      
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },


  /**
   * Update devotee profile
   */
  async updateProfile(id, profileData) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }
      
      let response;
      if (profileData instanceof FormData) {
        response = await api.put('/v1/profiles', profileData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        });
      } else {
        response = await api.put('/v1/profiles', profileData);
      }
      
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Join a temple - Updated to use the correct endpoint from backend
   */
  async joinTemple(entityId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.post('/v1/memberships', {
        entity_id: entityId
      });
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get temple memberships
   */
  async getMemberships() {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get('/v1/memberships');
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Check if devotee has joined a specific temple
   */
  async checkTempleJoined(templeId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get(`/v1/memberships/check/${templeId}`);
      return response;
    } catch (error) {
      // For 404 (not joined), return a valid response with joined: false
      if (error.response?.status === 404) {
        return { data: { joined: false } };
      }
      throw this.handleError(error);
    }
  },

  /**
   * FIXED: Get all devotees for an entity (matches backend route)
   * This matches the backend route: GET /v1/entities/:id/devotees
   */
  async getAllDevotees(entityId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      if (!entityId) {
        throw new Error('Entity ID is required');
      }

      const response = await api.get(`/v1/entities/${entityId}/devotees`);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * FIXED: Get devotee stats for an entity
   * This matches the backend route: GET /v1/entities/:id/devotee-stats
   */
  async getDevoteeStats(entityId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      if (!entityId) {
        throw new Error('Entity ID is required');
      }

      const response = await api.get(`/v1/entities/${entityId}/devotee-stats`);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * FIXED: Update devotee membership status
   * This matches the backend route: PATCH /v1/entities/:entityID/devotees/:userID/status
   */
  async updateDevoteeStatus(entityId, devoteeId, status) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.patch(`/v1/entities/${entityId}/devotees/${devoteeId}/status`, {
        status: status
      });
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get all profiles/devotees (fallback method)
   */
  async getAllProfiles() {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get('/v1/profiles');
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Add family member to devotee profile
   */
  async addFamilyMember(profileId, familyMemberData) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.post(`/v1/profiles/${profileId}/family`, familyMemberData);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Update family member
   */
  async updateFamilyMember(profileId, familyMemberId, familyMemberData) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.put(`/v1/profiles/${profileId}/family/${familyMemberId}`, familyMemberData);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Remove family member
   */
  async removeFamilyMember(profileId, familyMemberId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.delete(`/v1/profiles/${profileId}/family/${familyMemberId}`);
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get dashboard data - Uses entity-specific devotee dashboard endpoint
   */
  async getDashboardData(entityId) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      // This matches the backend route for devotee dashboard
      const response = await api.get(`/v1/entity/${entityId}/devotee/dashboard`);
      return response;
    } catch (error) {
      console.warn('Dashboard endpoint not available');
      throw this.handleError(error);
    }
  },

  /**
   * Search temples - matches backend route
   */
  async searchTemples(searchParams = {}) {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get('/v1/temples/search', { params: searchParams });
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get recent temples - matches backend route
   */
  async getRecentTemples() {
    try {
      if (!isAuthenticated()) {
        throw new Error('Authentication required');
      }

      const response = await api.get('/v1/temples/recent');
      return response;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Handle API errors consistently
   */
  handleError(error) {
    // Handle authentication errors specifically
    if (error.response?.status === 401) {
      console.error('Authentication error:', error.response?.data);
      redirectToLogin();
      return new Error('Authentication failed. Please login again.');
    }

    // Handle other HTTP errors
    if (error.response?.status === 403) {
      return new Error('Access denied. You do not have permission to perform this action.');
    }

    if (error.response?.status === 404) {
      return new Error('The requested resource was not found.');
    }

    if (error.response?.status >= 500) {
      return new Error('Server error. Please try again later.');
    }

    // Handle network errors
    if (error.code === 'NETWORK_ERROR' || !error.response) {
      return new Error('Network error. Please check your internet connection.');
    }

    // Handle timeout errors
    if (error.code === 'ECONNABORTED') {
      return new Error('Request timeout. Please try again.');
    }

    // Extract error messages from response
    if (error.response?.data?.message) {
      return new Error(error.response.data.message);
    } else if (error.response?.data?.error) {
      return new Error(error.response.data.error);
    } else if (error.message) {
      return new Error(error.message);
    } else {
      return new Error('An unexpected error occurred');
    }
  },

  /**
   * Utility methods for authentication
   */
  isAuthenticated,
  getAuthToken,
  clearAuth,
  redirectToLogin
};