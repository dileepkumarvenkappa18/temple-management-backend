import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import eventService from '@/services/event.service'
import { useToast } from '@/composables/useToast'

export const useEventStore = defineStore('event', () => {
  // State
  const events = ref([])
  const loading = ref(false)
  const error = ref(null)
  const selectedEvent = ref(null)
  const filters = ref({
    category: 'all',
    status: 'all',
    dateRange: 'upcoming',
    rsvpStatus: 'all'
  })
  const eventStats = ref({
    total: 0,
    upcoming: 0,
    thisMonth: 0,
    totalRSVPs: 0
  })
  const toast = useToast()

  // Getters
  const upcomingEvents = computed(() => {
    const today = new Date().toISOString().split('T')[0]
    return events.value
      .filter(event => {
        const eventDate = new Date(event.event_date || event.eventDate)
        return eventDate >= new Date(today) && event.is_active !== false
      })
      .sort((a, b) => {
        const dateA = new Date(a.event_date || a.eventDate)
        const dateB = new Date(b.event_date || b.eventDate)
        return dateA - dateB
      })
  })

  const pastEvents = computed(() => {
    const today = new Date().toISOString().split('T')[0]
    return events.value
      .filter(event => {
        const eventDate = new Date(event.event_date || event.eventDate)
        return eventDate < new Date(today)
      })
      .sort((a, b) => {
        const dateA = new Date(a.event_date || a.eventDate)
        const dateB = new Date(b.event_date || b.eventDate)
        return dateB - dateA
      })
  })

  const filteredEvents = computed(() => {
    let filtered = events.value

    // Category filter
    if (filters.value.category !== 'all') {
      filtered = filtered.filter(event => 
        (event.event_type || event.type || '').toLowerCase() === filters.value.category.toLowerCase()
      )
    }

    // Status filter
    if (filters.value.status !== 'all') {
      filtered = filtered.filter(event => {
        const eventDate = new Date(event.event_date || event.eventDate)
        const today = new Date()
        
        if (filters.value.status === 'upcoming') {
          return eventDate >= today
        } else if (filters.value.status === 'ongoing') {
          // Assuming an event is "ongoing" if it's today
          return eventDate.toDateString() === today.toDateString()
        } else if (filters.value.status === 'completed') {
          return eventDate < today
        }
        
        return true
      })
    }

    return filtered.sort((a, b) => {
      const dateA = new Date(a.event_date || a.eventDate)
      const dateB = new Date(b.event_date || b.eventDate)
      return dateA - dateB
    })
  })

  // Actions
  const fetchEvents = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await eventService.getEvents()
      events.value = normalizeEvents(response)
      await fetchEventStats()
      return events.value
    } catch (err) {
      error.value = err.message || 'Failed to fetch events'
      toast.error(error.value)
      console.error('Error fetching events:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchEventStats = async () => {
    try {
      const response = await eventService.getEventStats()
      eventStats.value = {
        total: response.total_events || 0,
        upcoming: response.upcoming_events || 0,
        thisMonth: response.this_month_events || 0,
        totalRSVPs: response.total_rsvps || 0
      }
      return eventStats.value
    } catch (err) {
      console.error('Error fetching event stats:', err)
      // Don't throw error, as this is a secondary feature
      return eventStats.value
    }
  }

  const fetchUpcomingEvents = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await eventService.getUpcomingEvents()
      const upcomingEvents = normalizeEvents(response)
      return upcomingEvents
    } catch (err) {
      error.value = err.message || 'Failed to fetch upcoming events'
      console.error('Error fetching upcoming events:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  const fetchEventById = async (id) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await eventService.getEventById(id)
      selectedEvent.value = normalizeEvent(response)
      return selectedEvent.value
    } catch (err) {
      error.value = err.message || `Failed to fetch event with ID: ${id}`
      toast.error(error.value)
      console.error('Error fetching event:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  const createEvent = async (eventData) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await eventService.createEvent(eventData)
      
      // Refresh events list
      await fetchEvents()
      
      toast.success('Event created successfully')
      
      return response
    } catch (err) {
      error.value = err.message || 'Failed to create event'
      toast.error(error.value)
      console.error('Error creating event:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

const updateEvent = async (id, updates) => {
  loading.value = true
  error.value = null

  try {
    const response = await eventService.updateEvent(id, updates)

    // Replace the updated event in the list
    const updatedEvent = normalizeEvent(response)
    const index = events.value.findIndex(e => e.id === Number(id))
    if (index !== -1) {
      events.value[index] = updatedEvent
    }

    toast.success('Event updated successfully')
    return response
  } catch (err) {
    error.value = err.message || `Failed to update event with ID: ${id}`
    toast.error(error.value)
    console.error('Error updating event:', err)
    throw err
  } finally {
    loading.value = false
  }
}


const deleteEvent = async (id) => {
  loading.value = true
  error.value = null

  try {
    const response = await eventService.deleteEvent(id)

    events.value = events.value.filter(e => e.id !== Number(id))

    toast.success('Event deleted successfully')
    return response
  } catch (err) {
    error.value = err.message || `Failed to delete event with ID: ${id}`
    toast.error(error.value)
    console.error('Error deleting event:', err)
    throw err
  } finally {
    loading.value = false
  }
}


  const setFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  const clearFilters = () => {
    filters.value = {
      category: 'all',
      status: 'all',
      dateRange: 'upcoming',
      rsvpStatus: 'all'
    }
  }

  const setSelectedEvent = (event) => {
    selectedEvent.value = event
  }

  // Helper function to normalize event data from backend
const normalizeEvent = (event) => {
  if (!event) return null

  // FIX 1: Debug logging to understand what we're receiving
  // console.log('Raw event data:', JSON.stringify(event, null, 2))

  // FIX 2: Better type handling - ensure we extract the correct event_type
  let eventType = '';
  if (event.event_type) {
    eventType = event.event_type.toLowerCase();
  } else if (event.eventType) {
    eventType = event.eventType.toLowerCase();
  } else if (event.type) {
    eventType = event.type.toLowerCase();
  } else {
    eventType = 'other'; // Default if nothing is found
  }

  let rawDate = event.event_date || event.eventDate || '';
  let rawTime = '';
  
  // FIX 3: Better time handling - prioritize event_time field
  if (event.event_time) {
    if (typeof event.event_time === 'string') {
      rawTime = event.event_time;
    } else {
      // It could be a date object or something else
      try {
        const timeObj = new Date(event.event_time);
        if (!isNaN(timeObj.getTime())) {
          const hours = String(timeObj.getHours()).padStart(2, '0');
          const minutes = String(timeObj.getMinutes()).padStart(2, '0');
          rawTime = `${hours}:${minutes}`;
        }
      } catch (e) {
        console.warn('Invalid event_time value:', event.event_time);
      }
    }
  } else if (event.eventTime) {
    rawTime = event.eventTime;
  }

  let fullDateTimeString = '';

  // Safely build ISO datetime
  if (typeof rawDate === 'string' && rawDate.includes('T')) {
    fullDateTimeString = rawDate;
  } else if (rawDate && rawTime && /^\d{2}:\d{2}/.test(rawTime)) {
    fullDateTimeString = `${rawDate}T${rawTime}`;
  } else if (rawDate) {
    fullDateTimeString = `${rawDate}T00:00`;
  } else {
    console.warn('Missing event_date for event:', event);
    fullDateTimeString = new Date().toISOString();
  }

  let isoDateString = '';
  try {
    const date = new Date(fullDateTimeString);
    if (isNaN(date.getTime())) throw new Error('Invalid Date');
    isoDateString = date.toISOString();
  } catch (err) {
    console.warn('Invalid date-time during normalizeEvent:', fullDateTimeString, err);
    isoDateString = new Date().toISOString(); // Fallback to now
  }

  return {
    id: event.id,
    title: event.title,
    description: event.description,
    type: eventType,
    eventDate: isoDateString,
    location: event.location,
    isActive: event.is_active !== false,
    createdBy: event.created_by || event.createdBy,
    createdAt: event.created_at || event.createdAt,
    entityId: event.entity_id || event.entityId,
    currentAttendees: event.rsvp_count || event.RSVPCount || 0,
    status: getEventStatus(isoDateString)
  }
}

  
  const normalizeEvents = (eventsData) => {
    if (!eventsData || !Array.isArray(eventsData)) return []
    return eventsData.map(event => normalizeEvent(event))
  }
  
  // Helper function to determine event status
  const getEventStatus = (dateStr) => {
    const eventDate = new Date(dateStr)
    const now = new Date()
    
    // Set time to midnight for date comparison
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const eventDay = new Date(eventDate.getFullYear(), eventDate.getMonth(), eventDate.getDate())
    
    if (eventDay > today) return 'upcoming'
    if (eventDay.getTime() === today.getTime()) return 'ongoing'
    return 'completed'
  }

  return {
    // State
    events,
    loading,
    error,
    selectedEvent,
    filters,
    eventStats,
    
    // Getters
    upcomingEvents,
    pastEvents,
    filteredEvents,
    
    // Actions
    fetchEvents,
    fetchUpcomingEvents,
    fetchEventById,
    fetchEventStats,
    createEvent,
    updateEvent,
    deleteEvent,
    setFilters,
    clearFilters,
    setSelectedEvent
  }
})