import api from './api.js'

class VolunteerService {
  constructor() {
    // These base URLs should match your actual API structure
    this.baseURL = '/api/v1/volunteers'
    this.entityBaseURL = '/api/v1/entities'
    this.membershipURL = '/api/v1/memberships' // Adjust based on your actual API structure
  }

  // Join a temple as a volunteer
  async joinTemple(entityId) {
    try {
      console.log(`Attempting to join temple with ID: ${entityId}`);
      
      // Try multiple possible endpoint formats to increase chances of success
      let response;
      try {
        // First attempt - standard format
        response = await api.post(this.membershipURL, {
          entity_id: entityId
        });
      } catch (firstError) {
        console.log('First join attempt failed, trying alternative endpoint...');
        try {
          // Second attempt - userprofile join endpoint format
          response = await api.post('/api/v1/userprofile/join-temple', {
            entity_id: entityId
          });
        } catch (secondError) {
          // Third attempt - entity-specific endpoint
          response = await api.post(`${this.entityBaseURL}/${entityId}/join`, {});
        }
      }
      
      console.log('Join temple response:', response.data);
      
      return {
        success: true,
        data: response.data,
        message: 'Temple joined successfully'
      }
    } catch (error) {
      console.error('Error joining temple:', error.response || error);
      return {
        success: false,
        message: error.response?.data?.error || 'Failed to join temple',
        error: error.response?.data || error.message
      }
    }
  }

  // Get all temples the current volunteer has joined
  async getUserMemberships() {
    try {
      console.log('Fetching user memberships...');
      
      // Try multiple possible endpoints to increase chances of success
      let response;
      try {
        // First attempt - standard format
        response = await api.get(this.membershipURL);
      } catch (firstError) {
        console.log('First membership attempt failed, trying alternative endpoint...');
        try {
          // Second attempt - user-based endpoint
          response = await api.get('/api/v1/userprofile/memberships');
        } catch (secondError) {
          // Third attempt - entities joined endpoint
          response = await api.get('/api/v1/user/joined-entities');
        }
      }
      
      console.log('User memberships response:', response.data);
      
      // Transform the response to a standard format regardless of API structure
      let memberships = response.data;
      
      // Handle different possible response formats
      if (Array.isArray(response.data.memberships)) {
        memberships = response.data.memberships;
      } else if (Array.isArray(response.data.data)) {
        memberships = response.data.data;
      } else if (typeof response.data === 'object' && !Array.isArray(response.data)) {
        // If it's an object with entity IDs as properties
        memberships = Object.keys(response.data).map(key => ({
          entity_id: parseInt(key),
          ...response.data[key]
        }));
      }
      
      return {
        success: true,
        data: memberships || [],
        message: 'User memberships fetched successfully'
      }
    } catch (error) {
      console.error('Error fetching memberships:', error.response || error);
      return {
        success: false,
        message: error.response?.data?.error || 'Failed to fetch user memberships',
        error: error.response?.data || error.message
      }
    }
  }

  // Get all volunteers for a specific entity/temple
  async getVolunteersByEntity(entityId, params = {}) {
    try {
      const queryParams = new URLSearchParams({
        page: params.page || 1,
        limit: params.limit || 10,
        search: params.search || '',
        status: params.status || '',
        joinedFrom: params.joinedFrom || '',
        joinedTo: params.joinedTo || '',
        sortBy: params.sortBy || 'createdAt',
        sortOrder: params.sortOrder || 'desc'
      }).toString();

      const response = await api.get(`${this.entityBaseURL}/${entityId}/volunteers?${queryParams}`);
      return {
        success: true,
        data: response.data.volunteers || [],
        meta: response.data.meta || {},
        message: 'Volunteers fetched successfully'
      }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.message || 'Failed to fetch volunteers',
        error: error.response?.data || error.message
      }
    }
  }

  // Get volunteer dashboard data
  async getVolunteerDashboard(volunteerId) {
    try {
      const response = await api.get(`${this.baseURL}/${volunteerId}/dashboard`);
      return {
        success: true,
        data: response.data.dashboard,
        message: 'Volunteer dashboard data fetched successfully'
      }
    } catch (error) {
      return {
        success: false,
        message: error.response?.data?.message || 'Failed to fetch volunteer dashboard data',
        error: error.response?.data || error.message
      }
    }
  }
}

export default new VolunteerService();