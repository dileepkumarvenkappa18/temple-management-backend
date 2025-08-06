import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import devoteeService from '@/services/devotee.service'
import donationService from '@/services/donation.service'
import { useToast } from '@/composables/useToast'
import { useAuth } from '@/composables/useAuth'

export const useDevoteeStore = defineStore('devotee', () => {
  const toast = useToast()
  const { user, currentEntity } = useAuth()

  // ===== STATE =====

  // Profile structured data
  const profile = ref({
    personal: {},
    spiritual: {},
    lineage: {},
    preferences: {},
    family: {},
    notes: {}
  })

  // Dashboard data for devotee home view
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

  // Donations state
  const donationHistory = ref([])
  const currentDonation = ref(null)
  const isProcessingDonation = ref(false)

  // Profile completion trackers
  const completedSteps = ref([])
  const incompleteSteps = ref([])
  const totalSteps = ref(6)
  const completionPercentage = computed(() =>
    totalSteps.value === 0 ? 0 : Math.round((completedSteps.value.length / totalSteps.value) * 100)
  )
  const isProfileComplete = computed(() => completionPercentage.value === 100)

  // Selected temple for context
  const selectedTempleId = ref(null)

  // Devotee management state
  const devotees = ref([])
  const devoteeStats = ref({
    total: 0,
    active: 0,
    inactive: 0,
    newThisMonth: 0,
    maleCount: 0,
    femaleCount: 0,
    otherCount: 0
  })
  const currentDevotee = ref(null)

  // Filters and sorting for devotees list
  const devoteeFilters = ref({
    status: 'all',
    search: '',
    gender: 'all',
    sortBy: 'name',
    sortOrder: 'asc'
  })

  // Computed filtered and sorted devotee list
  const filteredDevotees = computed(() => {
    let result = [...devotees.value]

    // Filter by status
    if (devoteeFilters.value.status !== 'all') {
      result = result.filter(devotee => devotee.status === devoteeFilters.value.status)
    }

    // Filter by gender
    if (devoteeFilters.value.gender !== 'all') {
      result = result.filter(devotee => devotee.gender === devoteeFilters.value.gender)
    }

    // Search filter (case-insensitive)
    if (devoteeFilters.value.search) {
      const searchTerm = devoteeFilters.value.search.toLowerCase()
      result = result.filter(devotee =>
        (devotee.full_name?.toLowerCase().includes(searchTerm)) ||
        (devotee.email?.toLowerCase().includes(searchTerm)) ||
        (devotee.phone?.includes(searchTerm)) ||
        (devotee.gotra?.toLowerCase().includes(searchTerm))
      )
    }

    // Sorting logic
    result.sort((a, b) => {
      let valueA, valueB
      switch (devoteeFilters.value.sortBy) {
        case 'name':
          valueA = a.full_name || ''
          valueB = b.full_name || ''
          break
        case 'email':
          valueA = a.email || ''
          valueB = b.email || ''
          break
        case 'gotra':
          valueA = a.gotra || ''
          valueB = b.gotra || ''
          break
        case 'date':
          valueA = a.created_at ? new Date(a.created_at) : new Date(0)
          valueB = b.created_at ? new Date(b.created_at) : new Date(0)
          break
        default:
          valueA = a.full_name || ''
          valueB = b.full_name || ''
      }

      if (valueA instanceof Date) {
        return devoteeFilters.value.sortOrder === 'asc' ? valueA - valueB : valueB - valueA
      }

      // String comparison
      return devoteeFilters.value.sortOrder === 'asc'
        ? valueA.localeCompare(valueB)
        : valueB.localeCompare(valueA)
    })

    return result
  })

  // General state
  const loading = ref(false)
  const error = ref(null)

  // ===== HELPERS =====

  // Safely get nested properties with default fallback
  const safeGet = (obj, path, defaultValue = '') => {
    try {
      return path.split('.').reduce((current, key) => current?.[key], obj) || defaultValue
    } catch {
      return defaultValue
    }
  }

  // Ensure value is an array or empty array
  const safeArray = (value) => Array.isArray(value) ? value : []

  // Get current entity ID from currentEntity or selected temple
  const getCurrentEntityId = () => currentEntity.value?.id || selectedTempleId.value

  // ===== PROFILE ACTIONS =====

  // Load devotee profile from API and map to state
   const loadProfileData = async () => {
    try {
      loading.value = true
      error.value = null
      
      const response = await devoteeService.getProfile()
      
      if (response.data) {
        const data = response.data
        // Map API response to our profile structure
        const profileData = {
          personal: {
            fullName: data.full_name || '',
            dateOfBirth: data.dob || '',
            gender: data.gender || '',
            phone: '', // Not directly in API
            email: '', // Not directly in API
            address: {
              street: data.street_address || '',
              city: data.city || '',
              state: data.state || '',
              pincode: data.pincode || '',
              country: data.country || 'India'
            }
          },
          spiritual: {
            gotra: data.gotra || '',
            nakshatra: data.nakshatra || '',
            rashi: data.rashi || '',
            lagna: data.lagna || '',
            vedaShaka: data.veda_shaka || ''
          },
          lineage: {
            father: {
              name: data.father_name || '',
              gotra: data.father_gotra || '',
              nativePlace: data.father_native_place || '',
              vedaShaka: data.father_veda_shaka || ''
            },
            mother: {
              name: data.mother_name || '',
              maidenGotra: data.maiden_gotra || '',
              nativePlace: data.mother_native_place || '',
              fatherName: data.maternal_grandfather_name || ''
            },
            paternalGrandfather: data.paternal_grandfather_name || '',
            paternalGrandmother: data.paternal_grandmother_name || '',
            maternalGrandfather: data.maternal_grandfather_name || '',
            maternalGrandmother: data.maternal_grandmother_name || ''
          },
          preferences: {
            favoriteSevas: getSevaPreferences(data),
            donationInterests: getDonationInterests(data),
            preferredFrequency: '',
            specialInterests: data.special_interests_or_notes || ''
          },
          family: {
            spouse: {
              name: data.spouse_name || '',
              email: data.spouse_email || '',
              phone: data.spouse_phone || '',
              dateOfBirth: data.spouse_dob || '',
              gotra: data.spouse_gotra || '',
              nakshatra: data.spouse_nakshatra || ''
            },
            children: (data.children || []).map(child => ({
              name: child.child_name || '',
              dateOfBirth: child.child_dob || '',
              gender: child.child_gender || '',
              education: child.child_education || '',
              interests: child.child_interests || ''
            })),
            emergencyContacts: (data.emergency_contacts || []).map(contact => ({
              name: contact.contact_name || '',
              relationship: contact.contact_relationship || '',
              phone: contact.contact_phone || '',
              address: contact.contact_address || ''
            }))
          },
          notes: {
            healthNotes: data.health_notes || '',
            allergies: data.allergies_or_conditions || '',
            dietaryRestrictions: data.dietary_restrictions || '',
            sankalpa: data.personal_sankalpa || '',
            additionalNotes: data.additional_notes || '',
            documents: [],
            profilePhoto: null,
            familyPhoto: null,
            shareProfile: true,
            receiveNotifications: true
          }
        }
        
        // Update the main profile state
        profile.value = profileData
        
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

  // Extract Seva preferences from data
  function getSevaPreferences(data) {
    const sevas = []
    if (data.seva_abhisheka) sevas.push('abhisheka')
    if (data.seva_arti) sevas.push('arti')
    if (data.seva_annadana) sevas.push('annadana')
    if (data.seva_archana) sevas.push('archana')
    if (data.seva_kalyanam) sevas.push('kalyanam')
    if (data.seva_homam) sevas.push('homam')
    return sevas
  }

  // Extract Donation interests from data
  function getDonationInterests(data) {
    const interests = []
    if (data.donate_temple_maintenance) interests.push('temple-maintenance')
    if (data.donate_annadana_program) interests.push('annadana')
    if (data.donate_festival_celebrations) interests.push('festivals')
    if (data.donate_religious_education) interests.push('education')
    if (data.donate_temple_construction) interests.push('construction')
    if (data.donate_general) interests.push('general')
    return interests
  }

  // Update a specific profile step
  const updateProfileStep = async (stepId, data) => {
     console.log(`📝 STORE UPDATING STEP ${stepId}:`, data);
    try {
      loading.value = true
      error.value = null
      
      // Update local state
      profile.value[stepId] = data
      
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

  // Submit full profile to backend to complete profile
  const completeProfile = async () => {
  console.log(`🔄 PROFILE STATE BEFORE API CONVERSION:`, JSON.stringify(profile.value, null, 2));

  try {
    loading.value = true
    error.value = null
    
    console.log('Starting profile completion...')
    
    // Helper for date conversion
    const convertToISODate = (dateString) => {
      if (!dateString) return null
      
      try {
        const date = new Date(dateString)
        if (isNaN(date.getTime())) return null
        return date.toISOString()
      } catch (error) {
        console.error('Date conversion error:', error)
        return null
      }
    }
    
    // Transform frontend data to backend structure
 const apiData = {
  // Personal Details
  full_name: profile.value?.personal?.fullName || null,
  dob: convertToISODate(profile.value?.personal?.dateOfBirth),
  gender: profile.value?.personal?.gender || null,
  street_address: profile.value?.personal?.address?.street || null,
  city: profile.value?.personal?.address?.city || null,
  state: profile.value?.personal?.address?.state || null,
  pincode: profile.value?.personal?.address?.pincode || null,
  country: profile.value?.personal?.address?.country || 'India',
  
  // Spiritual Info
  gotra: profile.value?.spiritual?.gotra || null,
  nakshatra: profile.value?.spiritual?.nakshatra || null,
  rashi: profile.value?.spiritual?.rashi || null,
  lagna: profile.value?.spiritual?.lagna || null,
  veda_shaka: profile.value?.spiritual?.vedaShaka || null,
  
  // Family Lineage - Added additional null checks here
  father_name: profile.value?.lineage?.father?.name || null,
  father_gotra: profile.value?.lineage?.father?.gotra || null,
  father_native_place: profile.value?.lineage?.father?.nativePlace || null,
  father_veda_shaka: profile.value?.lineage?.father?.vedaShaka || null,
  
  mother_name: profile.value?.lineage?.mother?.name || null,
  maiden_gotra: profile.value?.lineage?.mother?.maidenGotra || null,
  mother_native_place: profile.value?.lineage?.mother?.nativePlace || null,
  maternal_grandfather_name: profile.value?.lineage?.mother?.fatherName || null,
  
  paternal_grandfather_name: profile.value?.lineage?.paternalGrandfather || null,
  paternal_grandmother_name: profile.value?.lineage?.paternalGrandmother || null,
  maternal_grandmother_name: profile.value?.lineage?.maternalGrandmother || null,
  
  // Seva Preferences - Added null checks for preferences object
  seva_abhisheka: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
                 profile.value?.preferences?.favoriteSevas.includes('abhisheka') : false,
  seva_arti: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
            profile.value?.preferences?.favoriteSevas.includes('arti') : false,
  seva_annadana: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
                profile.value?.preferences?.favoriteSevas.includes('annadana') : false,
  seva_archana: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
               profile.value?.preferences?.favoriteSevas.includes('archana') : false,
  seva_kalyanam: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
                profile.value?.preferences?.favoriteSevas.includes('kalyanam') : false,
  seva_homam: Array.isArray(profile.value?.preferences?.favoriteSevas) ? 
             profile.value?.preferences?.favoriteSevas.includes('homam') : false,
  
  // Donation Interests - Added null checks
  donate_temple_maintenance: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                           profile.value?.preferences?.donationInterests.includes('temple-maintenance') : false,
  donate_annadana_program: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                          profile.value?.preferences?.donationInterests.includes('annadana') : false,
  donate_festival_celebrations: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                              profile.value?.preferences?.donationInterests.includes('festivals') : false,
  donate_religious_education: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                            profile.value?.preferences?.donationInterests.includes('education') : false,
  donate_temple_construction: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                            profile.value?.preferences?.donationInterests.includes('construction') : false,
  donate_general: Array.isArray(profile.value?.preferences?.donationInterests) ? 
                 profile.value?.preferences?.donationInterests.includes('general') : false,
  
  special_interests_or_notes: profile.value?.preferences?.specialInterests || null,
  
  // Family Members - Added null checks
  spouse_name: profile.value?.family?.spouse?.name || null,
  spouse_email: profile.value?.family?.spouse?.email || null,
  spouse_phone: profile.value?.family?.spouse?.phone || null,
  spouse_dob: convertToISODate(profile.value?.family?.spouse?.dateOfBirth),
  spouse_gotra: profile.value?.family?.spouse?.gotra || null,
  spouse_nakshatra: profile.value?.family?.spouse?.nakshatra || null,
  
  // Children - Enhanced null checks
  children: Array.isArray(profile.value?.family?.children) ? 
          profile.value?.family?.children.map(child => ({
            child_name: child?.name || '',
            child_dob: convertToISODate(child?.dateOfBirth),
            child_gender: child?.gender || '',
            child_education: child?.education || '',
            child_interests: child?.interests || ''
          })) : [],
  
  // Emergency Contacts - Enhanced null checks
  emergency_contacts: Array.isArray(profile.value?.family?.emergencyContacts) ? 
                    profile.value?.family?.emergencyContacts.map(contact => ({
                      contact_name: contact?.name || '',
                      contact_relationship: contact?.relationship || '',
                      contact_phone: contact?.phone || '',
                      contact_address: contact?.address || ''
                    })) : [],
  
  // Notes & Attachments - Added null checks
  health_notes: profile.value?.notes?.healthNotes || null,
  allergies_or_conditions: profile.value?.notes?.allergies || null,
  dietary_restrictions: profile.value?.notes?.dietaryRestrictions || null,
  personal_sankalpa: profile.value?.notes?.sankalpa || null,
  additional_notes: profile.value?.notes?.additionalNotes || null
}
    
    console.log('Sending API data:', apiData)
    
    // Send to backend
    const response = await devoteeService.createOrUpdateProfile(apiData)
    console.log('Profile completed:', response.data)
    toast.success('Profile completed successfully!')
    
    // Get entity ID from response or fallback to current
    const entityId = response.data?.entity_id || getCurrentEntityId()
    if (entityId) {
      return { 
        success: true, 
        data: response.data, 
        redirectPath: `/entity/${entityId}/devotee/dashboard` 
      }
    } else {
      return { 
        success: true, 
        data: response.data, 
        redirectPath: '/devotee/dashboard' 
      }
    }
  } catch (err) {
    error.value = err.response?.data?.error || err.message || 'Failed to complete profile'
    console.error('Error completing profile:', err)
    console.error('Error response:', err.response?.data)
    toast.error(error.value)
    throw err
  } finally {
    loading.value = false
  }
}


  // Check if profile has any files to upload
  const hasFiles = () => {
    return (
      safeArray(safeGet(profile.value, 'notes.documents', [])).length > 0 ||
      profile.value.notes?.profilePhoto ||
      profile.value.notes?.familyPhoto
    )
  }

  // Convert base64 data URL to File object
  const dataURLtoFile = (dataurl, filename) => {
    if (!dataurl) return null
    try {
      const arr = dataurl.split(',')
      if (arr.length < 2) return null
      const mimeMatch = arr[0].match(/:(.*?);/)
      if (!mimeMatch) return null
      const mime = mimeMatch[1]
      const bstr = atob(arr[1])
      let n = bstr.length
      const u8arr = new Uint8Array(n)
      while (n--) {
        u8arr[n] = bstr.charCodeAt(n)
      }
      return new File([u8arr], filename, { type: mime })
    } catch (error) {
      console.error('Error converting data URL to file:', error)
      return null
    }
  }


  // ===== TEMPLE ACTIONS =====

  // Join temple by ID
  const joinTemple = async (templeId) => {
    try {
      loading.value = true
      error.value = null

      const response = await devoteeService.joinTemple(templeId)
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

  // Fetch dashboard data for temple
  const fetchDashboardData = async (templeId) => {
    try {
      loading.value = true

      const response = await devoteeService.getDashboardData(templeId)
      if (response.data) {
        dashboardData.value = response.data
        return response.data
      }
      throw new Error('No dashboard data returned')
    } catch (error) {
      console.error('Error fetching dashboard data:', error)
      toast.error('Failed to load dashboard data')
      throw error
    } finally {
      loading.value = false
    }
  }

  // Update completed steps array, based on profile content
  const updateCompletedSteps = (profileData) => {
    completedSteps.value = []

    if (profileData.personal?.fullName) completedSteps.value.push('personal')
    if (profileData.spiritual?.gotra || profileData.spiritual?.nakshatra) completedSteps.value.push('spiritual')
    if (profileData.lineage?.father?.name || profileData.lineage?.mother?.name) completedSteps.value.push('lineage')
    if ((profileData.preferences?.favoriteSevas?.length || 0) > 0 || (profileData.preferences?.donationInterests?.length || 0) > 0) {
      completedSteps.value.push('preferences')
    }
    if (profileData.family?.spouse?.name || (profileData.family.children?.length || 0) > 0) completedSteps.value.push('family')
    if (profileData.notes?.sankalpa || profileData.notes?.allergies) completedSteps.value.push('notes')

    const allSteps = ['personal', 'spiritual', 'lineage', 'preferences', 'family', 'notes']
    incompleteSteps.value = allSteps.filter(step => !completedSteps.value.includes(step))
  }

  // Set currently selected temple ID
  const setSelectedTempleId = (templeId) => {
    selectedTempleId.value = templeId
  }

  // ===== DEVOTEE MANAGEMENT FUNCTIONS =====

  // Fetch devotees for an entity
  const fetchDevoteesByEntity = async (entityId) => {
    try {
      loading.value = true
      error.value = null

      if (!entityId) throw new Error('Entity ID is required')

      console.log('🔍 Fetching devotees for entity:', entityId)
      const response = await devoteeService.getAllDevotees(entityId)

      if (response.data) {
        devotees.value = response.data
        console.log('✅ Devotees loaded:', response.data.length)
        return response.data
      }

      devotees.value = []
      return []
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to fetch devotees'
      console.error('❌ Error fetching devotees:', err)
      devotees.value = []
      return []
    } finally {
      loading.value = false
    }
  }

  // Alias for compatibility
  const fetchDevotees = async (entityId) => fetchDevoteesByEntity(entityId)

  // Fetch devotee stats or fallback to calculation
  const fetchDevoteeStats = async (entityId) => {
    try {
      loading.value = true
      error.value = null

      if (!entityId) throw new Error('Entity ID is required')

      try {
        const response = await devoteeService.getDevoteeStats(entityId)
        if (response.data) {
          devoteeStats.value = response.data
          console.log('✅ Devotee stats loaded from API:', response.data)
          return response.data
        }
      } catch {
        // Fallback: calculate from local devotees data
        if (devotees.value.length === 0) {
          await fetchDevoteesByEntity(entityId)
        }
        const stats = {
          total: devotees.value.length,
          active: devotees.value.filter(d => d.status === 'active').length,
          inactive: devotees.value.filter(d => d.status === 'inactive').length,
          newThisMonth: calculateNewDevoteesThisMonth(),
          maleCount: devotees.value.filter(d => d.gender === 'male').length,
          femaleCount: devotees.value.filter(d => d.gender === 'female').length,
          otherCount: devotees.value.filter(d => d.gender !== 'male' && d.gender !== 'female').length
        }
        devoteeStats.value = stats
        console.log('✅ Devotee stats calculated:', stats)
        return stats
      }
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to fetch devotee statistics'
      console.error('❌ Error fetching devotee stats:', err)

      const defaultStats = {
        total: 0,
        active: 0,
        inactive: 0,
        newThisMonth: 0,
        maleCount: 0,
        femaleCount: 0,
        otherCount: 0
      }
      devoteeStats.value = defaultStats
      return defaultStats
    } finally {
      loading.value = false
    }
  }

  // Calculate number of new devotees added in the current month
  const calculateNewDevoteesThisMonth = () => {
    const now = new Date()
    const firstDayOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    return devotees.value.filter(devotee => {
      if (!devotee.created_at) return false
      const createdDate = new Date(devotee.created_at)
      return createdDate >= firstDayOfMonth
    }).length
  }

  // Fetch a single devotee's details by ID (try local cache first)
  const getDevoteeById = async (devoteeId) => {
    try {
      loading.value = true
      error.value = null

      // Ensure devoteeId is valid
      if (!devoteeId) {
        throw new Error('Devotee ID is required')
      }

      const existing = devotees.value.find(d => d.id === Number(devoteeId))
      if (existing) {
        currentDevotee.value = existing
        return existing
      }

      const response = await devoteeService.getDevoteeById(devoteeId)
      if (response.data) {
        currentDevotee.value = response.data
        return response.data
      }

      throw new Error('Devotee not found')
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to fetch devotee'
      console.error('Error fetching devotee:', err)
      toast.error(error.value)
      return null
    } finally {
      loading.value = false
    }
  }

  // Update devotee status (active/inactive) - FIXED VERSION
  const updateDevoteeStatus = async (entityId, devoteeId, status) => {
    try {
      loading.value = true
      error.value = null

      // Input validation
      if (!entityId) {
        throw new Error('Entity ID is required')
      }
      
      if (!devoteeId) {
        throw new Error('Devotee ID is required')
      }

      if (!status || !['active', 'inactive'].includes(status)) {
        throw new Error('Status must be either "active" or "inactive"')
      }

      // Convert IDs to proper format
      const entityIdNum = Number(entityId)
      const devoteeIdNum = Number(devoteeId)

      if (isNaN(entityIdNum) || isNaN(devoteeIdNum)) {
        throw new Error('Invalid ID format')
      }

      console.log('🔄 Updating devotee status:', { 
        entityId: entityIdNum, 
        devoteeId: devoteeIdNum, 
        status 
      })

      // Check if devotee exists in local store
      const devoteeIndex = devotees.value.findIndex(d => d.id === devoteeIdNum)
      if (devoteeIndex === -1) {
        console.warn('⚠️ Devotee not found in local store, proceeding with API call')
      }

      const response = await devoteeService.updateDevoteeStatus(entityIdNum, devoteeIdNum, status)

      if (!response || !response.data) {
        throw new Error('Invalid response from server')
      }

      // Update local devotee in the store
      if (devoteeIndex !== -1) {
        devotees.value[devoteeIndex] = {
          ...devotees.value[devoteeIndex],
          status: status,
          updated_at: new Date().toISOString()
        }
        console.log('✅ Local devotee updated:', devotees.value[devoteeIndex])
      }

      // Update currentDevotee if it matches
      if (currentDevotee.value && currentDevotee.value.id === devoteeIdNum) {
        currentDevotee.value = {
          ...currentDevotee.value,
          status: status,
          updated_at: new Date().toISOString()
        }
      }

      // Refresh stats after successful update
      try {
        await fetchDevoteeStats(entityIdNum)
      } catch (statsError) {
        console.warn('⚠️ Failed to refresh stats after status update:', statsError)
        // Don't throw here, as the main operation succeeded
      }

      const statusText = status === 'active' ? 'activated' : 'deactivated'
      toast.success(`Devotee ${statusText} successfully`)
      
      console.log('✅ Devotee status updated successfully')
      return response.data

    } catch (err) {
      const errorMessage = err.response?.data?.error || 
                          err.response?.data?.message || 
                          err.message || 
                          'Failed to update devotee status'
      
      error.value = errorMessage
      console.error('❌ Error updating devotee status:', {
        error: err,
        entityId,
        devoteeId,
        status,
        response: err.response?.data
      })
      
      toast.error(errorMessage)
      throw new Error(errorMessage)
    } finally {
      loading.value = false
    }
  }

  // Convenience functions for activating/deactivating devotees - IMPROVED
  const activateDevotee = async (entityId, devoteeId) => {
    try {
      return await updateDevoteeStatus(entityId, devoteeId, 'active')
    } catch (err) {
      console.error('❌ Error activating devotee:', err)
      throw err
    }
  }

  const deactivateDevotee = async (entityId, devoteeId) => {
    try {
      return await updateDevoteeStatus(entityId, devoteeId, 'inactive')
    } catch (err) {
      console.error('❌ Error deactivating devotee:', err)
      throw err
    }
  }

  // Update devotee filters
  const setDevoteeFilter = (key, value) => {
    if (devoteeFilters.value.hasOwnProperty(key)) {
      devoteeFilters.value[key] = value
    } else {
      console.warn(`⚠️ Invalid filter key: ${key}`)
    }
  }

  // Reset all devotee filters to default
  const resetDevoteeFilters = () => {
    devoteeFilters.value = {
      status: 'all',
      search: '',
      gender: 'all',
      sortBy: 'name',
      sortOrder: 'asc'
    }
  }

  // ===== DONATION FUNCTIONS =====

  // Create a new donation order for Razorpay
  const createDonation = async (donationData) => {
    try {
      isProcessingDonation.value = true
      error.value = null

      // Validate donation data
      if (!donationData || !donationData.amount) {
        throw new Error('Invalid donation data')
      }

      // Razorpay expects amount in paise (INR subunits)
      const dataToSend = {
        ...donationData,
        amount: Math.round(donationData.amount * 100), // Ensure integer
        currency: 'INR'
      }

      const response = await donationService.createDonation(dataToSend)
      if (response.data) {
        currentDonation.value = response.data
        return response.data
      }

      throw new Error('Failed to create donation order')
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to create donation'
      toast.error(error.value)
      console.error('Error creating donation:', err)
      throw err
    } finally {
      isProcessingDonation.value = false
    }
  }

  // Verify a completed Razorpay payment
  const verifyDonation = async (verificationData) => {
    try {
      isProcessingDonation.value = true
      error.value = null

      if (!verificationData) {
        throw new Error('Verification data is required')
      }

      const response = await donationService.verifyDonation(verificationData)
      if (response.data) {
        toast.success('Donation verified successfully!')

        // Refresh donation history and dashboard
        try {
          await fetchDonationHistory()
          if (selectedTempleId.value) {
            await fetchDashboardData(selectedTempleId.value)
          }
        } catch (refreshError) {
          console.warn('⚠️ Failed to refresh data after donation verification:', refreshError)
        }

        return response.data
      }

      throw new Error('Failed to verify donation')
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to verify donation'
      toast.error(error.value)
      console.error('Error verifying donation:', err)
      throw err
    } finally {
      isProcessingDonation.value = false
    }
  }

  // Fetch donation history for current user
  const fetchDonationHistory = async () => {
    try {
      loading.value = true
      error.value = null

      const response = await donationService.getMyDonations()
      if (response.data) {
        donationHistory.value = response.data
        return response.data
      }

      donationHistory.value = []
      return []
    } catch (err) {
      error.value = err.response?.data?.error || err.message || 'Failed to fetch donation history'
      console.error('Error fetching donation history:', err)
      donationHistory.value = []
      return []
    } finally {
      loading.value = false
    }
  }

  // ===== UTILITY FUNCTIONS =====

  // Clear all errors
  const clearError = () => {
    error.value = null
  }

  // Reset store to initial state
  const resetStore = () => {
    // Reset profile data
    profile.value = {
      personal: {},
      spiritual: {},
      lineage: {},
      preferences: {},
      family: {},
      notes: {}
    }

    // Reset dashboard data
    dashboardData.value = {
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
    }

    // Reset other state
    donationHistory.value = []
    currentDonation.value = null
    devotees.value = []
    currentDevotee.value = null
    selectedTempleId.value = null
    completedSteps.value = []
    incompleteSteps.value = []
    
    // Reset flags
    loading.value = false
    isProcessingDonation.value = false
    error.value = null

    // Reset filters
    resetDevoteeFilters()
  }

  // ===== RETURN STATE & ACTIONS =====

  return {
    // State
    profile,
    dashboardData,
    loading,
    error,
    selectedTempleId,
    donationHistory,
    currentDonation,
    isProcessingDonation,

    // Profile completion
    completedSteps,
    incompleteSteps,
    totalSteps,
    completionPercentage,
    isProfileComplete,

    // Devotee management
    devotees,
    devoteeStats,
    currentDevotee,
    devoteeFilters,
    filteredDevotees,

    // Profile actions
    loadProfileData,
    updateProfileStep,
    completeProfile,

    // Temple actions
    fetchDashboardData,
    joinTemple,
    setSelectedTempleId,
    getCurrentEntityId,

    // Devotee management actions
    fetchDevoteesByEntity,
    fetchDevotees, // alias
    fetchDevoteeStats,
    getDevoteeById,
    updateDevoteeStatus,
    activateDevotee,
    deactivateDevotee,
    setDevoteeFilter,
    resetDevoteeFilters,

    // Donation actions
    createDonation,
    verifyDonation,
    fetchDonationHistory,

    // Utility actions
    clearError,
    resetStore
  }
})