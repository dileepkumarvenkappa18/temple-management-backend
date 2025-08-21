import { defineStore } from 'pinia';
import { formatDate } from '@/utils/date';
import axios from 'axios'; // Import axios directly

export const useAuditLogStore = defineStore('auditLog', {
  state: () => ({
    logs: [],
    selectedLog: null,
    isLoading: false,
    error: null,
    fetchError: null, // Additional error field for detailed debugging
    currentPage: 1,
    totalPages: 1,
    limit: 10,
    total: 0,
    rawResponse: null, // For debugging
    lastFetchMethod: null // Track which method was used
  }),
  
  actions: {
    async fetchAuditLogs() {
      this.isLoading = true;
      this.error = null;
      this.fetchError = null;
      
      // Try multiple approaches to get the data
      await this.fetchWithFetch();
      
      // If first approach failed, try with axios
      if (this.logs.length === 0 && !this.error) {
        await this.fetchWithAxios();
      }
      
      // If both approaches failed, try with hardcoded URL
      if (this.logs.length === 0 && !this.error) {
        await this.fetchWithHardcodedUrl();
      }
      
      this.isLoading = false;
    },
    
    // Approach 1: Using fetch
    async fetchWithFetch() {
      try {
        console.log('Fetching audit logs with native fetch...');
        
        const response = await fetch(`/api/v1/auditlogs?page=${this.currentPage}&limit=${this.limit}`, {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          credentials: 'include' // Include cookies for authentication
        });
        
        console.log('Fetch response status:', response.status);
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        this.rawResponse = data;
        
        console.log('Fetch API Response:', data);
        
        if (data && Array.isArray(data.data)) {
          console.log('Setting logs from fetch data.data:', data.data);
          this.logs = data.data;
          this.totalPages = data.total_pages || 1;
          this.total = data.total || 0;
          this.lastFetchMethod = 'fetch';
          return true; // Successfully fetched data
        } else {
          console.warn('Invalid response format from fetch:', data);
          return false;
        }
      } catch (error) {
        console.error('Error in fetchWithFetch:', error);
        this.fetchError = `Fetch error: ${error.message}`;
        return false;
      }
    },
    
    // Approach 2: Using axios
    async fetchWithAxios() {
      try {
        console.log('Fetching audit logs with axios...');
        
        const response = await axios.get('/api/v1/auditlogs', {
          params: {
            page: this.currentPage,
            limit: this.limit
          }
        });
        
        console.log('Axios response:', response);
        
        const data = response.data;
        this.rawResponse = data;
        
        if (data && Array.isArray(data.data)) {
          console.log('Setting logs from axios data.data:', data.data);
          this.logs = data.data;
          this.totalPages = data.total_pages || 1;
          this.total = data.total || 0;
          this.lastFetchMethod = 'axios';
          return true;
        } else {
          console.warn('Invalid response format from axios:', data);
          return false;
        }
      } catch (error) {
        console.error('Error in fetchWithAxios:', error);
        this.fetchError = `${this.fetchError || ''}\nAxios error: ${error.message}`;
        return false;
      }
    },
    
    // Approach 3: Try with hardcoded full URL
    async fetchWithHardcodedUrl() {
      try {
        console.log('Fetching audit logs with hardcoded URL...');
        
        // Determine the base URL from the current location
        const baseUrl = window.location.origin;
        const url = `${baseUrl}/api/v1/auditlogs?page=${this.currentPage}&limit=${this.limit}`;
        
        console.log('Attempting with URL:', url);
        
        const response = await fetch(url, {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const data = await response.json();
        this.rawResponse = data;
        
        if (data && Array.isArray(data.data)) {
          console.log('Setting logs from hardcoded URL data.data:', data.data);
          this.logs = data.data;
          this.totalPages = data.total_pages || 1;
          this.total = data.total || 0;
          this.lastFetchMethod = 'hardcoded URL';
          return true;
        } else {
          console.warn('Invalid response format from hardcoded URL:', data);
          return false;
        }
      } catch (error) {
        console.error('Error in fetchWithHardcodedUrl:', error);
        this.fetchError = `${this.fetchError || ''}\nHardcoded URL error: ${error.message}`;
        this.error = `Failed to fetch audit logs after multiple attempts. Please check the console for details.`;
        return false;
      }
    },
    
    async fetchAuditLogDetails(id) {
      this.isDetailLoading = true;
      
      try {
        const response = await fetch(`/api/v1/auditlogs/${id}`, {
          method: 'GET',
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json'
          },
          credentials: 'include'
        });
        
        if (!response.ok) {
          throw new Error(`API error: ${response.status}`);
        }
        
        const logData = await response.json();
        
        // Try to parse details if it's a string
        let parsedDetails = {};
        if (typeof logData.details === 'string' && logData.details) {
          try {
            parsedDetails = JSON.parse(logData.details);
          } catch (e) {
            parsedDetails = { raw: logData.details };
          }
        }
        
        // Extract user information
        let userName = 'System';
        if (parsedDetails && parsedDetails.target_user_name) {
          userName = parsedDetails.target_user_name;
        }
        
        this.selectedLog = {
          id: logData.id,
          userName: userName,
          entityName: logData.entity_id ? `Entity ${logData.entity_id}` : '-',
          action: logData.action || '',
          status: logData.status || 'unknown',
          ipAddress: logData.ip_address || '-',
          details: parsedDetails,
          createdAt: logData.created_at || new Date().toISOString(),
          formattedDate: formatDate(logData.created_at || new Date(), 'DD MMM YYYY, HH:mm:ss')
        };
      } catch (error) {
        console.error('Error fetching log details:', error);
        this.error = error.message;
        this.selectedLog = null;
      } finally {
        this.isDetailLoading = false;
      }
    },
    
    setPage(page) {
      this.currentPage = page;
      this.fetchAuditLogs();
    },
    
    setLimit(limit) {
      this.limit = limit;
      this.currentPage = 1;
      this.fetchAuditLogs();
    },
    
    clearSelectedLog() {
      this.selectedLog = null;
    }
  }
});