// src/stores/volunteer.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useVolunteerStore = defineStore('volunteer', () => {
  // State
  const volunteers = ref([])
  const currentVolunteer = ref(null)
  const loading = ref(false)
  const error = ref(null)
  const filters = ref({
    search: '',
    status: 'all',
    temple: '',
    skills: []
  })
  const pagination = ref({
    page: 1,
    limit: 10,
    total: 0
  })

  // Mock data for development
  const mockVolunteers = ref([
    {
      id: 1,
      name: 'Ramesh Kumar',
      email: 'ramesh@example.com',
      phone: '+91 9876543210',
      templeId: 1,
      templeName: 'Sri Venkateshwara Temple',
      status: 'active',
      skills: ['cooking', 'decoration', 'crowd_management'],
      joinedDate: '2024-01-15',
      avatar: null,
      assignments: [
        { id: 1, sevaName: 'Morning Pooja', date: '2024-06-15', status: 'confirmed' },
        { id: 2, sevaName: 'Evening Aarti', date: '2024-06-16', status: 'pending' }
      ],
      availability: {
        monday: true,
        tuesday: true,
        wednesday: false,
        thursday: true,
        friday: true,
        saturday: true,
        sunday: true
      }
    },
    {
      id: 2,
      name: 'Priya Sharma',
      email: 'priya@example.com',
      phone: '+91 9876543211',
      templeId: 1,
      templeName: 'Sri Venkateshwara Temple',
      status: 'active',
      skills: ['music', 'teaching', 'coordination'],
      joinedDate: '2024-02-20',
      avatar: null,
      assignments: [
        { id: 3, sevaName: 'Bhajan Session', date: '2024-06-17', status: 'confirmed' }
      ],
      availability: {
        monday: false,
        tuesday: true,
        wednesday: true,
        thursday: false,
        friday: true,
        saturday: true,
        sunday: true
      }
    },
    {
      id: 3,
      name: 'Arjun Patel',
      email: 'arjun@example.com',
      phone: '+91 9876543212',
      templeId: 2,
      templeName: 'Krishna Mandir',
      status: 'inactive',
      skills: ['security', 'maintenance'],
      joinedDate: '2024-03-10',
      avatar: null,
      assignments: [],
      availability: {
        monday: true,
        tuesday: true,
        wednesday: true,
        thursday: true,
        friday: false,
        saturday: false,
        sunday: true
      }
    }
  ])

  // Getters
  const filteredVolunteers = computed(() => {
    let result = mockVolunteers.value

    // Apply search filter
    if (filters.value.search) {
      const search = filters.value.search.toLowerCase()
      result = result.filter(volunteer => 
        volunteer.name.toLowerCase().includes(search) ||
        volunteer.email.toLowerCase().includes(search) ||
        volunteer.templeName.toLowerCase().includes(search)
      )
    }

    // Apply status filter
    if (filters.value.status !== 'all') {
      result = result.filter(volunteer => volunteer.status === filters.value.status)
    }

    // Apply temple filter
    if (filters.value.temple) {
      result = result.filter(volunteer => volunteer.templeId.toString() === filters.value.temple)
    }

    // Apply skills filter
    if (filters.value.skills.length > 0) {
      result = result.filter(volunteer => 
        filters.value.skills.some(skill => volunteer.skills.includes(skill))
      )
    }

    return result
  })

  const volunteerStats = computed(() => ({
    total: mockVolunteers.value.length,
    active: mockVolunteers.value.filter(v => v.status === 'active').length,
    inactive: mockVolunteers.value.filter(v => v.status === 'inactive').length,
    withAssignments: mockVolunteers.value.filter(v => v.assignments.length > 0).length
  }))

  const availableSkills = computed(() => {
    const skills = new Set()
    mockVolunteers.value.forEach(volunteer => {
      volunteer.skills.forEach(skill => skills.add(skill))
    })
    return Array.from(skills).sort()
  })

  // Actions
  const fetchVolunteers = async (templeId = null) => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 500))
      
      if (templeId) {
        volunteers.value = mockVolunteers.value.filter(v => v.templeId === templeId)
      } else {
        volunteers.value = mockVolunteers.value
      }
      
      pagination.value.total = volunteers.value.length
    } catch (err) {
      error.value = 'Failed to fetch volunteers'
      console.error('Error fetching volunteers:', err)
    } finally {
      loading.value = false
    }
  }

  const getVolunteerById = async (id) => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 300))
      
      const volunteer = mockVolunteers.value.find(v => v.id === parseInt(id))
      if (volunteer) {
        currentVolunteer.value = volunteer
        return volunteer
      } else {
        throw new Error('Volunteer not found')
      }
    } catch (err) {
      error.value = 'Failed to fetch volunteer details'
      console.error('Error fetching volunteer:', err)
      return null
    } finally {
      loading.value = false
    }
  }

  const updateVolunteerStatus = async (id, status) => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 300))
      
      const volunteerIndex = mockVolunteers.value.findIndex(v => v.id === id)
      if (volunteerIndex !== -1) {
        mockVolunteers.value[volunteerIndex].status = status
        
        // Update current volunteer if it's the same one
        if (currentVolunteer.value && currentVolunteer.value.id === id) {
          currentVolunteer.value.status = status
        }
        
        return true
      } else {
        throw new Error('Volunteer not found')
      }
    } catch (err) {
      error.value = 'Failed to update volunteer status'
      console.error('Error updating volunteer status:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  const assignVolunteerToSeva = async (volunteerId, sevaId, sevaName, date) => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 300))
      
      const volunteerIndex = mockVolunteers.value.findIndex(v => v.id === volunteerId)
      if (volunteerIndex !== -1) {
        const newAssignment = {
          id: Date.now(),
          sevaId,
          sevaName,
          date,
          status: 'pending'
        }
        
        mockVolunteers.value[volunteerIndex].assignments.push(newAssignment)
        
        // Update current volunteer if it's the same one
        if (currentVolunteer.value && currentVolunteer.value.id === volunteerId) {
          currentVolunteer.value.assignments.push(newAssignment)
        }
        
        return true
      } else {
        throw new Error('Volunteer not found')
      }
    } catch (err) {
      error.value = 'Failed to assign volunteer to seva'
      console.error('Error assigning volunteer:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  const updateVolunteerAvailability = async (id, availability) => {
    loading.value = true
    error.value = null
    
    try {
      // Simulate API call
      await new Promise(resolve => setTimeout(resolve, 300))
      
      const volunteerIndex = mockVolunteers.value.findIndex(v => v.id === id)
      if (volunteerIndex !== -1) {
        mockVolunteers.value[volunteerIndex].availability = availability
        
        // Update current volunteer if it's the same one
        if (currentVolunteer.value && currentVolunteer.value.id === id) {
          currentVolunteer.value.availability = availability
        }
        
        return true
      } else {
        throw new Error('Volunteer not found')
      }
    } catch (err) {
      error.value = 'Failed to update volunteer availability'
      console.error('Error updating availability:', err)
      return false
    } finally {
      loading.value = false
    }
  }

  const setFilters = (newFilters) => {
    filters.value = { ...filters.value, ...newFilters }
  }

  const clearFilters = () => {
    filters.value = {
      search: '',
      status: 'all',
      temple: '',
      skills: []
    }
  }

  const setPagination = (newPagination) => {
    pagination.value = { ...pagination.value, ...newPagination }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    // State
    volunteers,
    currentVolunteer,
    loading,
    error,
    filters,
    pagination,
    
    // Getters
    filteredVolunteers,
    volunteerStats,
    availableSkills,
    
    // Actions
    fetchVolunteers,
    getVolunteerById,
    updateVolunteerStatus,
    assignVolunteerToSeva,
    updateVolunteerAvailability,
    setFilters,
    clearFilters,
    setPagination,
    clearError
  }
})