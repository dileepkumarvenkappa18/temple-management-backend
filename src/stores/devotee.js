// src/stores/devotee.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import devoteeService from '@/services/devotee.service'
import { useToast } from '@/composables/useToast'

export const useDevoteeStore = defineStore('devotee', () => {
  const toast = useToast()

  // State
  const profile = ref({
    personalDetails: {
      fullName: '',
      dateOfBirth: '',
      gender: '',
      phone: '',
      email: '',
      address: {
        street: '',
        city: '',
        state: '',
        pincode: '',
        country: 'India'
      }
    }
  })

  // Dashboard data state
  const dashboardData = ref({
    summary: {
      totalDonations: 0,
      totalSevaBookings: 0,
      upcomingEvents: 0,
      profileCompletion: 0
    },
    recentDonations: [],
    upcomingSevas: [],
    announcements: [],
    quickStats: {
      thisMonthDonations: 0,
      thisMonthSevas: 0,
      favoriteDeity: '',
      memberSince: ''
    }
  })

  // Profile completion state
  const completedSteps = ref([])
  const incompleteSteps = ref([])
  const totalSteps = ref(6)
  const completionPercentage = computed(() =>
    totalSteps.value === 0
      ? 0
      : Math.round((completedSteps.value.length / totalSteps.value) * 100)
  )
  const isProfileComplete = computed(() => completionPercentage.value === 100)
  const selectedTempleId = ref(null)

  // General state
  const loading = ref(false)
  const error = ref(null)

  // Actions
  // Load profile data from API
  const loadProfileData = async () => {
    try {
      loading.value = true
      error.value = null
      
      const response = await devoteeService.getProfile()
      
      if (response.data) {
        // Map API response to our profile structure
        // This depends on the exact API response format
        // This is a simplified example
        const profileData = {
          personal: {
            fullName: response.data.fullName || '',
            dateOfBirth: response.data.dateOfBirth || '',
            gender: response.data.gender || '',
            phone: response.data.phone || '',
            email: response.data.email || '',
            address: {
              street: response.data.address || '',
              city: response.data.city || '',
              state: response.data.state || '',
              pincode: response.data.pincode || '',
              country: 'India'
            }
          },
          spiritual: {
            gotra: response.data.gotra || '',
            nakshatra: response.data.nakshatra || '',
            rashi: response.data.rashi || '',
            lagna: response.data.lagna || '',
            vedaShaka: response.data.vedaShaka || ''
          },
          lineage: {
            fatherName: response.data.fatherName || '',
            motherName: response.data.motherName || '',
            grandfatherName: response.data.grandfatherName || '',
            grandmotherName: response.data.grandmotherName || ''
          }
          // Add other profile sections based on your API response
        }
        
        // Update completed steps based on data availability
        updateCompletedSteps(profileData)
        
        return profileData
      }
      
      return null
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to load profile data'
      console.error('Error loading profile data:', err)
      // Don't show error toast here as this might be normal for new users
      return null
    } finally {
      loading.value = false
    }
  }

  // Update a specific profile step
  const updateProfileStep = async (stepId, data) => {
    try {
      loading.value = true
      error.value = null
      
      // Update local state
      switch (stepId) {
        case 'personal':
          profile.value.personalDetails = data
          break
        case 'spiritual':
          profile.value.spiritualInfo = data
          break
        case 'lineage':
          profile.value.familyLineage = data
          break
        case 'preferences':
          profile.value.sevaPreferences = data
          break
        case 'family':
          profile.value.familyMembers = data
          break
        case 'notes':
          profile.value.notesAttachments = data
          break
      }
      
      // Don't send API request for each step to avoid multiple calls
      // Just update local state and mark as completed
      if (!completedSteps.value.includes(stepId)) {
        completedSteps.value.push(stepId)
        incompleteSteps.value = incompleteSteps.value.filter(step => step !== stepId)
      }
      
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.error || `Error updating profile step: ${stepId}`
      console.error(`Error updating profile step ${stepId}:`, err)
      toast.error(error.value)
      return { success: false }
    } finally {
      loading.value = false
    }
  }

  // Complete the entire profile by sending all data to API
  const completeProfile = async () => {
    try {
      loading.value = true
      error.value = null
      
      // Prepare data for API in the required format
      const apiData = {
        date_of_birth: profile.value.personalDetails?.dateOfBirth,
        gender: profile.value.personalDetails?.gender,
        address: profile.value.personalDetails?.address?.street,
        gotra: profile.value.spiritualInfo?.gotra,
        nakshatra: profile.value.spiritualInfo?.nakshatra,
        rashi: profile.value.spiritualInfo?.rashi,
        lagna: profile.value.spiritualInfo?.lagna,
        veda_shaka: profile.value.spiritualInfo?.vedaShaka,
        father_name: profile.value.familyLineage?.fatherName,
        mother_name: profile.value.familyLineage?.motherName,
        grandfather_name: profile.value.familyLineage?.grandfatherName,
        grandmother_name: profile.value.familyLineage?.grandmotherName,
        family_tree_url: profile.value.familyLineage?.familyTreeUrl,
        favorite_sevas: JSON.stringify(profile.value.sevaPreferences?.interests || []),
        donation_interests: JSON.stringify(profile.value.sevaPreferences?.donationInterests || []),
        spouse_info: JSON.stringify(profile.value.familyMembers?.spouse || null),
        children_info: JSON.stringify(profile.value.familyMembers?.children || []),
        emergency_contacts: JSON.stringify(profile.value.familyMembers?.emergencyContacts || []),
        health_info: profile.value.notesAttachments?.healthInfo,
        sankalpa_notes: profile.value.notesAttachments?.notes,
        photos_url: JSON.stringify(profile.value.notesAttachments?.photos || [])
      }
      
      // Prepare FormData if we have files
      const formData = new FormData()
      formData.append('data', JSON.stringify(apiData))
      
      if (profile.value.notesAttachments?.documents) {
        profile.value.notesAttachments.documents.forEach((doc, index) => {
          formData.append(`document_${index}`, doc.file)
        })
      }
      
      if (profile.value.personalDetails?.photo) {
        formData.append('profile_photo', profile.value.personalDetails.photo)
      }
      
      const response = await devoteeService.createOrUpdateProfile(formData)
      
      toast.success('Profile completed successfully!')
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to complete profile'
      toast.error(error.value)
      console.error('Error completing profile:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Join a temple
  const joinTemple = async (templeId) => {
    try {
      loading.value = true
      error.value = null
      
      const response = await devoteeService.joinTemple(templeId)
      
      // Set the selected temple ID after joining
      selectedTempleId.value = templeId
      
      toast.success('Successfully joined temple')
      return response.data
    } catch (err) {
      error.value = err.response?.data?.error || 'Failed to join temple'
      toast.error(error.value)
      console.error('Error joining temple:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Fetch dashboard data
  const fetchDashboardData = async (templeId) => {
    try {
      loading.value = true
      
      // In a real implementation, you would call an API endpoint here
      // For now, we'll just use dummy data
      const dummyData = {
        summary: {
          totalDonations: 45000,
          totalSevaBookings: 15,
          upcomingEvents: 3,
          profileCompletion: 85
        },
        recentDonations: [
          {
            id: 1,
            amount: 2500,
            purpose: 'Abhishekam Seva',
            date: '2024-06-20',
            status: 'completed'
          },
          {
            id: 2,
            amount: 1000,
            purpose: 'Annadanam',
            date: '2024-06-18',
            status: 'completed'
          }
        ],
        upcomingSevas: [
          {
            id: 1,
            name: 'Satyanarayana Vratam',
            date: '2024-06-28',
            time: '10:00 AM',
            status: 'confirmed'
          }
        ],
        announcements: [
          {
            id: 1,
            title: 'Brahmotsavam Festival',
            message: 'Annual Brahmotsavam festival starts from July 10th.',
            date: '2024-06-25',
            priority: 'high'
          }
        ],
        quickStats: {
          thisMonthDonations: 8500,
          thisMonthSevas: 4,
          favoriteDeity: 'Lord Venkateswara',
          memberSince: '2024-01-15'
        }
      }

      // Simulate API call delay
      await new Promise(resolve => setTimeout(resolve, 500))
      
      dashboardData.value = dummyData
      return dummyData
    } catch (error) {
      console.error('Error fetching dashboard data:', error)
      toast.error('Failed to load dashboard data')
      throw error
    } finally {
      loading.value = false
    }
  }

  // Helper function to update completed steps based on data
  const updateCompletedSteps = (profileData) => {
    completedSteps.value = []
    
    // Check personal details
    if (profileData.personal && profileData.personal.fullName) {
      completedSteps.value.push('personal')
    }
    
    // Check spiritual info
    if (profileData.spiritual && (profileData.spiritual.gotra || profileData.spiritual.nakshatra)) {
      completedSteps.value.push('spiritual')
    }
    
    // Check lineage
    if (profileData.lineage && (profileData.lineage.fatherName || profileData.lineage.motherName)) {
      completedSteps.value.push('lineage')
    }
    
    // Update incomplete steps
    const allSteps = ['personal', 'spiritual', 'lineage', 'preferences', 'family', 'notes']
    incompleteSteps.value = allSteps.filter(step => !completedSteps.value.includes(step))
  }

  // Set selected temple ID
  const setSelectedTempleId = (templeId) => {
    selectedTempleId.value = templeId
  }

  return {
    // State
    profile,
    dashboardData, 
    loading,
    error,
    selectedTempleId,
    
    // Profile completion
    completedSteps,
    incompleteSteps,
    totalSteps,
    completionPercentage,
    isProfileComplete,
    
    // Actions
    loadProfileData,
    updateProfileStep,
    completeProfile,
    fetchDashboardData,
    joinTemple,
    setSelectedTempleId
  }
})