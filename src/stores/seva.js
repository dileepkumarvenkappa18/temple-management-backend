// src/stores/seva.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { sevaService } from '@/services/seva.service'

export const useSevaStore = defineStore('seva', () => {
  // State
  const sevas = ref([])
  const sevaBookings = ref([])
  const loading = ref(false)
  const error = ref(null)
  const searchQuery = ref('')
  const filters = ref({
    status: 'all',
    type: 'all',
    dateRange: 'all'
  })
  const selectedSeva = ref(null)
  
  // New state for recent sevas
  const recentSevas = ref([])
  const loadingRecentSevas = ref(false)
  
  // New state for seva catalog (all available sevas)
  const sevaCatalog = ref([])
  const loadingCatalog = ref(false)
  
  // Computed properties
  const filteredSevas = computed(() => {
    let filtered = sevas.value
    
    // Search filter
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(seva =>
        seva.name?.toLowerCase().includes(query) ||
        seva.description?.toLowerCase().includes(query)
      )
    }
    
    // Status filter
    if (filters.value.status !== 'all') {
      filtered = filtered.filter(seva => seva.status === filters.value.status)
    }
    
    // Type filter
    if (filters.value.type !== 'all') {
      filtered = filtered.filter(seva => seva.type === filters.value.type)
    }
    
    return filtered
  })
  
  // Computed stats
  const sevaStats = computed(() => ({
    total: sevas.value.length,
    active: sevas.value.filter(s => s.is_active).length,
    pending: sevas.value.filter(s => s.status === 'pending').length,
  }))
  
  // Actions - Connect to backend API
  const fetchSevas = async (params = {}) => {
    loading.value = true
    error.value = null
    
    try {
      // Only pass filtering parameters
      const response = await sevaService.getSevas({
        page: params.page || 1,
        limit: params.limit || 10,
        seva_type: params.seva_type || '',
        search: params.search || ''
      })
      
      if (response.success) {
        sevas.value = response.data || []
        return sevas.value
      } else {
        error.value = response.error
        return []
      }
    } catch (err) {
      console.error('Error fetching sevas:', err)
      error.value = err.message || 'Failed to fetch sevas'
      return []
    } finally {
      loading.value = false
    }
  }
  
  // New method to fetch seva catalog (all sevas for mapping)
  const fetchSevaCatalog = async () => {
    loadingCatalog.value = true
    try {
      const response = await sevaService.getSevas({
        page: 1,
        limit: 1000, // Get all sevas
        search: ''
      })
      
      if (response.success) {
        sevaCatalog.value = response.data || []
        console.log('Seva catalog loaded:', sevaCatalog.value)
        return sevaCatalog.value
      } else {
        console.error('Failed to fetch seva catalog:', response.error)
        sevaCatalog.value = []
        return []
      }
    } catch (error) {
      console.error('Error fetching seva catalog:', error)
      sevaCatalog.value = []
      return []
    } finally {
      loadingCatalog.value = false
    }
  }
  
  const fetchEntityBookings = async () => {
    loading.value = true
    error.value = null
    
    try {
      const response = await sevaService.getEntityBookings()
      sevaBookings.value = response.data || []
      return sevaBookings.value
    } catch (err) {
      console.error('Error fetching seva bookings:', err)
      error.value = err.response?.data?.error || 'Failed to fetch bookings'
      return []
    } finally {
      loading.value = false
    }
  }
  
  // Fetch recent sevas for the devotee dashboard with seva names
  const fetchRecentSevas = async () => {
    loadingRecentSevas.value = true
    try {
      // First get the bookings
      const response = await sevaService.getMyBookings()
      
      if (response.success) {
        let bookings = response.data || []
        
        // If we have bookings, ensure seva catalog is loaded
        if (bookings.length > 0) {
          // Fetch seva catalog if not already loaded
          if (sevaCatalog.value.length === 0) {
            await fetchSevaCatalog()
          }
          
          // Map seva names to bookings
          bookings = bookings.map(booking => {
            const sevaId = booking.seva_id || booking.SevaID
            const seva = sevaCatalog.value.find(s => s.id === sevaId || s.ID === sevaId)
            
            return {
              ...booking,
              seva_name: seva?.name || seva?.Name || `Seva ${sevaId}`,
              seva_type: seva?.type || seva?.Type || seva?.seva_type || '',
              seva_description: seva?.description || seva?.Description || '',
              seva: seva ? {
                id: seva.id || seva.ID,
                name: seva.name || seva.Name,
                type: seva.type || seva.Type || seva.seva_type,
                description: seva.description || seva.Description
              } : null
            }
          })
          
          console.log('Bookings with seva names:', bookings)
        }
        
        // Sort by booking time, newest first
        const sorted = [...bookings].sort((a, b) => {
          const dateA = new Date(a.booking_time || a.BookingTime || a.created_at || Date.now())
          const dateB = new Date(b.booking_time || b.BookingTime || b.created_at || Date.now())
          return dateB - dateA
        })
        
        recentSevas.value = sorted
      } else {
        recentSevas.value = []
      }
    } catch (error) {
      console.error('Failed to fetch recent sevas:', error)
      recentSevas.value = []
    } finally {
      loadingRecentSevas.value = false
    }
  }
  
  const createSeva = async (sevaData) => {
    loading.value = true
    error.value = null
    
    try {
      // Map the form data to match the Go backend model
      console.log('Creating seva with data:', sevaData)
      
      const response = await sevaService.createSeva(sevaData)
      
      // Add the new seva to the list if successful
      if (response.success && response.data) {
        sevas.value.push(response.data)
        // Also add to catalog
        sevaCatalog.value.push(response.data)
      }
      
      return { 
        success: response.success, 
        message: response.message || 'Seva created successfully',
        data: response.data
      }
    } catch (err) {
      console.error('Error creating seva:', err)
      error.value = err.response?.data?.error || 'Failed to create seva'
      return { 
        success: false, 
        message: error.value
      }
    } finally {
      loading.value = false
    }
  }
  
  const updateSeva = async (sevaId, sevaData) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await sevaService.updateSeva(sevaId, sevaData)
      
      // Update the seva in the list
      if (response.success && response.data) {
        const index = sevas.value.findIndex(s => s.id === sevaId)
        if (index !== -1) {
          sevas.value[index] = response.data
        }
        
        // Also update in catalog
        const catalogIndex = sevaCatalog.value.findIndex(s => s.id === sevaId)
        if (catalogIndex !== -1) {
          sevaCatalog.value[catalogIndex] = response.data
        }
      }
      
      return { 
        success: response.success, 
        message: response.message || 'Seva updated successfully',
        data: response.data
      }
    } catch (err) {
      console.error('Error updating seva:', err)
      error.value = err.response?.data?.error || 'Failed to update seva'
      return { 
        success: false, 
        message: error.value
      }
    } finally {
      loading.value = false
    }
  }
  
  const deleteSeva = async (sevaId) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await sevaService.deleteSeva(sevaId)
      
      if (response.success) {
        // Remove the seva from the list
        sevas.value = sevas.value.filter(s => s.id !== sevaId)
        // Also remove from catalog
        sevaCatalog.value = sevaCatalog.value.filter(s => s.id !== sevaId)
      }
      
      return { 
        success: response.success, 
        message: response.message || 'Seva deleted successfully' 
      }
    } catch (err) {
      console.error('Error deleting seva:', err)
      error.value = err.response?.data?.error || 'Failed to delete seva'
      return { 
        success: false, 
        message: error.value 
      }
    } finally {
      loading.value = false
    }
  }
  
  const updateSevaStatus = async (sevaId, status) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await sevaService.updateBookingStatus(sevaId, status)
      
      if (response.success) {
        // Update the seva status in the list
        const index = sevas.value.findIndex(s => s.id === sevaId)
        if (index !== -1) {
          sevas.value[index].status = status
        }
      }
      
      return { 
        success: response.success, 
        message: response.message || `Seva ${status} successfully` 
      }
    } catch (err) {
      console.error('Error updating seva status:', err)
      error.value = err.response?.data?.error || 'Failed to update seva status'
      return { 
        success: false, 
        message: error.value 
      }
    } finally {
      loading.value = false
    }
  }
  
  // Utility functions
  const setSearchQuery = (query) => {
    searchQuery.value = query
  }
  
  const setFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }
  
  const clearFilters = () => {
    searchQuery.value = ''
    filters.value = {
      status: 'all',
      type: 'all',
      dateRange: 'all'
    }
  }
  
  const setSelectedSeva = (seva) => {
    selectedSeva.value = seva
  }
  
  return {
    // State
    sevas,
    sevaBookings,
    loading,
    error,
    searchQuery,
    filters,
    selectedSeva,
    recentSevas,
    loadingRecentSevas,
    sevaCatalog,
    loadingCatalog,
    
    // Getters
    filteredSevas,
    sevaStats,
    
    // Actions
    fetchSevas,
    fetchSevaCatalog,
    fetchEntityBookings,
    fetchRecentSevas,
    createSeva,
    updateSeva,
    deleteSeva,
    updateSevaStatus,
    setSearchQuery,
    setFilters,
    clearFilters,
    setSelectedSeva
  }
})