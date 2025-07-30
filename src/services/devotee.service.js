// src/services/devotee.service.js
import axios from 'axios';

// Get the API base URL from environment variables
const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080/api/v1';

// Get the auth token from localStorage
const getAuthToken = () => {
  return localStorage.getItem('temple_auth_token');
};

// Create axios instance with auth headers
const createAuthHeader = () => {
  return {
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'application/json'
    }
  };
};

// Create multipart form data header with auth
const createMultipartHeader = () => {
  return {
    headers: {
      'Authorization': `Bearer ${getAuthToken()}`,
      'Content-Type': 'multipart/form-data'
    }
  };
};

export default {
  /**
   * Get devotee profile
   */
  async getProfile() {
    try {
      const response = await axios.get(`${API_URL}/userprofiles/me`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get devotee by ID (for admin purposes)
   */
  async getDevoteeById(id) {
    try {
      const response = await axios.get(`${API_URL}/userprofiles/${id}`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Create or update devotee profile
   */
  async createOrUpdateProfile(profileData) {
    try {
      // Check if we're dealing with FormData or plain object
      let response;
      if (profileData instanceof FormData) {
        response = await axios.post(`${API_URL}/userprofiles`, profileData, createMultipartHeader());
      } else {
        // The API expects a specific format - format the data appropriately
        const apiData = {
          date_of_birth: profileData.dateOfBirth,
          gender: profileData.gender,
          address: profileData.address?.street,
          gotra: profileData.gotra,
          nakshatra: profileData.nakshatra,
          rashi: profileData.rashi,
          lagna: profileData.lagna,
          veda_shaka: profileData.vedaShaka,
          father_name: profileData.fatherName,
          mother_name: profileData.motherName,
          grandfather_name: profileData.grandfatherName,
          grandmother_name: profileData.grandmotherName,
          family_tree_url: profileData.familyTreeUrl,
          favorite_sevas: profileData.favoriteSevas ? JSON.stringify(profileData.favoriteSevas) : null,
          donation_interests: profileData.donationInterests ? JSON.stringify(profileData.donationInterests) : null,
          spouse_info: profileData.spouseInfo ? JSON.stringify(profileData.spouseInfo) : null,
          children_info: profileData.childrenInfo ? JSON.stringify(profileData.childrenInfo) : null,
          emergency_contacts: profileData.emergencyContacts ? JSON.stringify(profileData.emergencyContacts) : null,
          health_info: profileData.healthInfo,
          sankalpa_notes: profileData.sankalpaNotes
        };

        response = await axios.post(`${API_URL}/userprofiles`, apiData, createAuthHeader());
      }
      
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Update devotee profile
   */
  async updateProfile(id, profileData) {
    try {
      let response;
      if (profileData instanceof FormData) {
        response = await axios.put(`${API_URL}/userprofiles/${id}`, profileData, createMultipartHeader());
      } else {
        // Format profile data for the API
        const apiData = {
          date_of_birth: profileData.dateOfBirth,
          gender: profileData.gender,
          address: profileData.address?.street,
          gotra: profileData.gotra,
          nakshatra: profileData.nakshatra,
          rashi: profileData.rashi,
          lagna: profileData.lagna,
          veda_shaka: profileData.vedaShaka,
          father_name: profileData.fatherName,
          mother_name: profileData.motherName,
          grandfather_name: profileData.grandfatherName,
          grandmother_name: profileData.grandmotherName,
          family_tree_url: profileData.familyTreeUrl,
          favorite_sevas: profileData.favoriteSevas ? JSON.stringify(profileData.favoriteSevas) : null,
          donation_interests: profileData.donationInterests ? JSON.stringify(profileData.donationInterests) : null,
          spouse_info: profileData.spouseInfo ? JSON.stringify(profileData.spouseInfo) : null,
          children_info: profileData.childrenInfo ? JSON.stringify(profileData.childrenInfo) : null,
          emergency_contacts: profileData.emergencyContacts ? JSON.stringify(profileData.emergencyContacts) : null,
          health_info: profileData.healthInfo,
          sankalpa_notes: profileData.sankalpaNotes
        };
        
        response = await axios.put(`${API_URL}/userprofiles/${id}`, apiData, createAuthHeader());
      }
      
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Join a temple - Updated to use the correct endpoint from backend
   */
  async joinTemple(entityId) {
    try {
      const response = await axios.post(`${API_URL}/memberships`, {
        entity_id: entityId
      }, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get temple memberships
   */
  async getMemberships() {
    try {
      const response = await axios.get(`${API_URL}/memberships`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get all devotees (for admin purposes)
   */
  async getAllDevotees() {
    try {
      const response = await axios.get(`${API_URL}/userprofiles`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Add family member to devotee profile
   */
  async addFamilyMember(profileId, familyMemberData) {
    try {
      const response = await axios.post(`${API_URL}/userprofiles/${profileId}/family`, familyMemberData, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Update family member
   */
  async updateFamilyMember(profileId, familyMemberId, familyMemberData) {
    try {
      const response = await axios.put(`${API_URL}/userprofiles/${profileId}/family/${familyMemberId}`, familyMemberData, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Remove family member
   */
  async removeFamilyMember(profileId, familyMemberId) {
    try {
      const response = await axios.delete(`${API_URL}/userprofiles/${profileId}/family/${familyMemberId}`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  /**
   * Get dashboard data
   */
  async getDashboardData() {
    try {
      const response = await axios.get(`${API_URL}/devotee/dashboard`, createAuthHeader());
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },
  
  /**
   * Handle API errors consistently
   */
  handleError(error) {
    if (error.response?.data?.message) {
      return new Error(error.response.data.message);
    } else if (error.response?.data?.error) {
      return new Error(error.response.data.error);
    } else if (error.message) {
      return new Error(error.message);
    } else {
      return new Error('An unexpected error occurred');
    }
  }
};