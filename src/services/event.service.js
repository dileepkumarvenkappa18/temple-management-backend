// src/services/event.service.js
import { apiClient } from '@/plugins/axios';

const eventService = {
  async getEvents() {
  try {
    const response = await apiClient.event.getAll();
    console.log("🔍 Full Response:", response); // Confirmed from your logs
    return response.data || []; // ✅ FIXED: return directly
  } catch (error) {
    throw this.handleError(error);
  }
},

  async getUpcomingEvents() {
    try {
      const response = await apiClient.event.getUpcoming();
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async getEventById(id) {
    try {
      const response = await apiClient.event.getById(id);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async getEventStats() {
    try {
      const response = await apiClient.event.getStats();
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async createEvent(eventData) {
    try {
      if (eventData instanceof FormData) {
        const dataJson = eventData.get('data');
        if (dataJson) {
          const parsedData = JSON.parse(dataJson);
          const date = new Date(parsedData.eventDate);
          const dateStr = date.toISOString().split('T')[0];
          const timeStr = date.toTimeString().slice(0, 5);

          const apiData = {
            title: parsedData.title,
            description: parsedData.description || '',
            // FIX: Don't default to 'other' - use event_type directly
            event_type: parsedData.event_type || parsedData.type || parsedData.eventType,
            event_date: dateStr,
            event_time: timeStr,
            location: parsedData.location || 'Temple Premises',
            is_active: parsedData.isActive !== undefined ? parsedData.isActive : true
          };

          console.log('Creating event with data:', apiData); // Debug log
          const response = await apiClient.event.create(apiData);
          return response.data;
        }
      } else {
        const apiData = {
          title: eventData.title,
          description: eventData.description || '',
          // FIX: Don't default to 'other' - use event_type directly
          event_type: eventData.event_type || eventData.type || eventData.eventType,
          event_date: eventData.event_date || eventData.date,
          event_time: eventData.event_time || eventData.time,
          location: eventData.location || 'Temple Premises',
          is_active: eventData.isActive !== undefined ? eventData.isActive : true
        };

        console.log('Creating event with data:', apiData); // Debug log
        const response = await apiClient.event.create(apiData);
        return response.data;
      }
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async updateEvent(id, eventData) {
    try {
      const apiData = {
        title: eventData.title,
        description: eventData.description || '',
        // FIX: Don't default to 'other' - use event_type directly
        event_type: eventData.event_type || eventData.type || eventData.eventType,
        event_date: eventData.event_date || eventData.date,
        event_time: eventData.event_time || eventData.time,
        location: eventData.location || 'Temple Premises',
        is_active: eventData.isActive !== undefined ? eventData.isActive : true
      };

      console.log('Updating event with data:', apiData); // Debug log
      const response = await apiClient.event.update(id, apiData);

      return {
        ...response.data,
        message: 'Event updated successfully',
        wasUpdate: true,
        id
      };
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async deleteEvent(id) {
    try {
      const response = await apiClient.event.delete(id);
      return {
        message: 'Event deleted successfully',
        id
      };
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async registerForEvent(eventId) {
    try {
      const response = await apiClient.event.createRSVP(eventId);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async cancelRegistration(eventId) {
    try {
      console.warn('CANCEL RSVP: No endpoint exists.');
      throw new Error('RSVP cancellation not supported by backend API');
    } catch (error) {
      throw this.handleError(error);
    }
  },

  async getEventRSVPs(eventId) {
    try {
      const response = await apiClient.event.getRSVPs(eventId);
      return response.data;
    } catch (error) {
      throw this.handleError(error);
    }
  },

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

export default eventService;