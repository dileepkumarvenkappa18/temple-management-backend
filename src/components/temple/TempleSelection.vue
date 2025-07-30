<template>
  <div class="min-h-screen bg-gradient-to-b from-indigo-50 via-white to-indigo-50">
    <!-- Header -->
    <div class="bg-indigo-600">
      <div class="max-w-5xl mx-auto px-4 sm:px-6 lg:px-8 py-8 relative z-10">
        <div class="text-center">
          <h1 class="text-3xl md:text-4xl font-bold text-white mb-3 font-heading">Find Your Spiritual Home</h1>
          <p class="text-indigo-100 text-base md:text-lg font-light font-sans">
            Begin your spiritual journey by joining an approved temple.
          </p>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-6xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Search Section -->
      <div class="bg-white rounded-xl shadow-lg border border-indigo-100 p-5 mb-8">
        <div class="flex flex-col sm:flex-row gap-4">
          <!-- Search Input -->
          <div class="flex-1">
            <label for="search" class="block text-sm font-semibold text-indigo-900 mb-2 font-side">
              Find Your Temple
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
              </div>
              <input
                id="search"
                v-model="searchQuery"
                type="text"
                placeholder="Search by temple name, city, or state..."
                class="block w-full pl-10 pr-3 py-3 border border-indigo-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-indigo-800"
              />
            </div>
          </div>

          <!-- Filters -->
          <div class="flex flex-col sm:flex-row gap-3">
            <div>
              <label class="block text-sm font-semibold text-indigo-900 mb-2 font-side">State</label>
              <select
                v-model="selectedState"
                class="px-4 py-3 border border-indigo-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-indigo-800 min-w-[140px] appearance-none"
              >
                <option value="">All States</option>
                <option value="Karnataka">Karnataka</option>
                <option value="Tamil Nadu">Tamil Nadu</option>
                <option value="Kerala">Kerala</option>
                <option value="Andhra Pradesh">Andhra Pradesh</option>
                <option value="Maharashtra">Maharashtra</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-semibold text-indigo-900 mb-2 font-side">Type</label>
              <select
                v-model="selectedCategory"
                class="px-4 py-3 border border-indigo-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-indigo-800 min-w-[140px] appearance-none"
              >
                <option value="">All Types</option>
                <option value="Traditional">Traditional</option>
                <option value="Modern">Modern</option>
                <option value="Heritage">Heritage</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="loading" class="flex flex-col items-center justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-100 border-t-indigo-600"></div>
        <p class="text-indigo-600 font-medium mt-4">Discovering sacred spaces...</p>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredTemples.length === 0 && !loading" class="text-center py-12 bg-white rounded-xl shadow-lg">
        <div class="max-w-md mx-auto px-6">
          <svg class="h-12 w-12 text-indigo-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
          </svg>
          <h3 class="text-xl font-bold text-indigo-900 mb-3 font-heading">No Temples Found</h3>
          <p class="text-indigo-700 mb-6 font-side">
            We couldn't find any temples matching your criteria. Please try different search terms.
          </p>
          <button 
            @click="resetFilters" 
            class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700"
          >
            Reset Filters
          </button>
        </div>
      </div>

      <!-- Temple Results -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="temple in filteredTemples"
          :key="temple.id || temple.ID"
          class="bg-white rounded-xl shadow-lg border border-gray-200 overflow-hidden"
        >
          <!-- Temple Header -->
          <div class="h-28 bg-gray-50 border-b border-gray-100 flex items-center justify-center">
            <div class="bg-white rounded-full p-3 shadow-sm border border-gray-200">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-10 w-10 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
              </svg>
            </div>
          </div>

          <!-- Temple Info -->
          <div class="p-5">
            <h3 class="text-lg font-bold text-gray-900 mb-2 font-heading">{{ temple.name || temple.Name }}</h3>

            <div class="space-y-2 mb-4">
              <div class="flex items-center text-sm text-gray-700 font-side">
                <svg class="h-4 w-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                {{ temple.city || temple.City || 'Unknown' }}, {{ temple.state || temple.State || 'Unknown' }}
              </div>
              
              <div v-if="temple.category || temple.templeType || temple.TempleType" class="flex items-center text-sm text-gray-700 font-side">
                <svg class="h-4 w-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"></path>
                </svg>
                <span class="px-2 py-0.5 bg-gray-100 rounded-md text-gray-700 text-xs font-medium">
                  {{ temple.category || temple.templeType || temple.TempleType }}
                </span>
              </div>

              <div class="flex items-center text-sm text-gray-700 font-side">
                <svg class="h-4 w-4 mr-2 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
                </svg>
                <span class="font-medium">{{ temple.devoteeCount || temple.DevoteeCount || 0 }}</span> devotees
              </div>
            </div>

            <!-- Description -->
            <p v-if="temple.description || temple.Description" class="text-sm text-gray-600 mb-4 line-clamp-2 font-side">
              {{ temple.description || temple.Description }}
            </p>

            <!-- Join Button -->
            <button
              @click="selectTemple(temple)"
              :disabled="joiningTemple === (temple.id || temple.ID)"
              class="w-full bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2.5 px-4 rounded-lg flex items-center justify-center"
            >
              <div v-if="joiningTemple === (temple.id || temple.ID)" class="flex items-center">
                <div class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2"></div>
                Joining...
              </div>
              <span v-else>Join Temple</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Join Confirmation Modal -->
    <div
      v-if="showConfirmModal"
      class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50"
      @click="closeModal"
    >
      <div
        class="bg-white rounded-xl shadow-xl max-w-md w-full p-6"
        @click.stop
      >
        <div class="text-center">
          <div class="mx-auto flex items-center justify-center h-14 w-14 rounded-full bg-indigo-100 mb-4">
            <svg class="h-8 w-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
            </svg>
          </div>
          <h3 class="text-xl font-bold text-indigo-900 mb-2 font-heading">
            Join {{ selectedTemple?.name || selectedTemple?.Name }}?
          </h3>
          <p class="text-indigo-700 mb-6 font-side">
            You're about to join {{ selectedTemple?.name || selectedTemple?.Name }} as a devotee. 
            You'll be able to book sevas, make donations, and participate in temple events.
          </p>
          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="closeModal"
              class="flex-1 bg-white hover:bg-gray-50 text-indigo-700 font-medium py-2 px-4 rounded-lg border border-indigo-200"
            >
              Cancel
            </button>
            <button
              @click="confirmJoinTemple"
              :disabled="isJoining"
              class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-lg"
            >
              <span v-if="isJoining" class="flex items-center justify-center">
                <div class="animate-spin rounded-full h-4 w-4 border-2 border-white border-t-transparent mr-2"></div>
                Joining...
              </span>
              <span v-else>Join Temple</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Modal -->
    <div
      v-if="showSuccessModal"
      class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50"
    >
      <div
        class="bg-white rounded-xl shadow-xl max-w-md w-full p-6"
        @click.stop
      >
        <div class="text-center">
          <!-- Success Icon -->
          <div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
            <svg class="h-10 w-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
            </svg>
          </div>

          <!-- Success Message -->
          <h3 class="text-xl font-bold text-indigo-900 mb-2 font-heading">
            Successfully Joined!
          </h3>
          <p class="text-indigo-700 mb-6 font-side">
            You've successfully joined <span class="font-semibold">{{ joinedTempleName }}</span>!
            Start building your profile to get the most out of your temple experience.
          </p>

          <!-- Action Buttons -->
          <div class="flex flex-col sm:flex-row gap-3">
            <button
              @click="goToProfileCreation"
              class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-lg"
            >
              Complete Profile
            </button>
            <button
              @click="goToDashboard"
              class="flex-1 bg-white hover:bg-gray-50 text-indigo-700 font-medium py-2 px-4 rounded-lg border border-indigo-200"
            >
              Go to Dashboard
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'
import { useAuth } from '@/composables/useAuth'
import templeService from '@/services/temple.service'

export default {
  name: 'TempleSelection',
  setup() {
    const router = useRouter()
    const { showToast } = useToast()
    const { user } = useAuth()

    // Reactive data
    const searchQuery = ref('')
    const selectedState = ref('')
    const selectedCategory = ref('')
    const loading = ref(true)
    const joiningTemple = ref(null)
    const showConfirmModal = ref(false)
    const isJoining = ref(false)
    const selectedTemple = ref(null)
    const temples = ref([])
    const joinedTemples = ref([])

    // Success modal state
    const showSuccessModal = ref(false)
    const joinedTempleId = ref(null)
    const joinedTempleName = ref('')

    // Load joined temples from localStorage
    const loadJoinedTemples = () => {
      try {
        const storedJoinedTemples = localStorage.getItem('devoteeJoinedTemples')
        if (storedJoinedTemples) {
          joinedTemples.value = JSON.parse(storedJoinedTemples)
        }
      } catch (error) {
        console.error('Error loading joined temples from localStorage:', error)
      }
    }

    // Save joined temples to localStorage
    const saveJoinedTemples = () => {
      try {
        localStorage.setItem('devoteeJoinedTemples', JSON.stringify(joinedTemples.value))
      } catch (error) {
        console.error('Error saving joined temples to localStorage:', error)
      }
    }

    // Check if a temple is already joined
    const isTempleJoined = (templeId) => {
      return joinedTemples.value.includes(templeId)
    }

    // Computed
    const filteredTemples = computed(() => {
      if (!temples.value) return []
      
      let filtered = temples.value.filter(temple => 
        temple.status === 'approved' || 
        temple.status === 'APPROVED' || 
        temple.Status === 'approved' ||
        temple.Status === 'APPROVED' ||
        !temple.status // Include temples without status for testing
      )
      
      // Search filter
      if (searchQuery.value.trim()) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(temple => {
          const name = (temple.name || temple.Name || '').toLowerCase()
          const city = (temple.city || temple.City || '').toLowerCase()
          const state = (temple.state || temple.State || '').toLowerCase()
          
          return name.includes(query) || city.includes(query) || state.includes(query)
        })
      }
      
      // State filter
      if (selectedState.value) {
        filtered = filtered.filter(temple => 
          (temple.state || temple.State) === selectedState.value
        )
      }
      
      // Category filter
      if (selectedCategory.value) {
        filtered = filtered.filter(temple => 
          (temple.category || temple.templeType || temple.TempleType) === selectedCategory.value
        )
      }
      
      return filtered
    })

    // Methods
    const resetFilters = () => {
      searchQuery.value = ''
      selectedState.value = ''
      selectedCategory.value = ''
    }

    // Open confirmation modal
    const selectTemple = (temple) => {
      selectedTemple.value = temple
      showConfirmModal.value = true
    }

    // Close modal
    const closeModal = () => {
      showConfirmModal.value = false
    }

    // Close success modal
    const closeSuccessModal = () => {
      showSuccessModal.value = false
    }

    // Go to profile creation
    const goToProfileCreation = () => {
      showSuccessModal.value = false
      
      // Use session storage instead of localStorage for temporary data
      sessionStorage.setItem('tempEntityId', joinedTempleId.value.toString())
      sessionStorage.setItem('tempTempleName', joinedTempleName.value)
      
      router.push('/devotee/profile/create')
    }

    // Go to dashboard
    const goToDashboard = () => {
      showSuccessModal.value = false
      
      // Use session storage for the current session only
      sessionStorage.setItem('tempEntityId', joinedTempleId.value.toString())
      sessionStorage.setItem('tempTempleName', joinedTempleName.value)
      
      router.push(`/entity/${joinedTempleId.value}/devotee/dashboard`)
    }

    // Join temple after confirmation
    const confirmJoinTemple = async () => {
      if (!selectedTemple.value) return
      
      const templeId = selectedTemple.value.id || selectedTemple.value.ID
      isJoining.value = true
      
      try {
        console.log(`Joining temple with ID: ${templeId}`)
        
        // Call the updated service method
        await templeService.joinTemple(templeId);
        
        // Update local state
        if (!joinedTemples.value.includes(templeId)) {
          joinedTemples.value.push(templeId)
          saveJoinedTemples()
        }
        
        showToast('Successfully joined temple!', 'success')
        
        // Close the confirmation modal
        closeModal()
        
        // Show success modal with options
        joinedTempleId.value = templeId
        joinedTempleName.value = selectedTemple.value.name || selectedTemple.value.Name
        showSuccessModal.value = true
        
      } catch (error) {
        console.error('Error joining temple:', error)
        showToast('Failed to join temple. Please try again.', 'error')
        closeModal()
      } finally {
        isJoining.value = false
      }
    }

    // Fetch temples from the backend
    const fetchTemples = async () => {
      loading.value = true
      
      try {
        console.log('Fetching temples from API')
        
        // Build search params from filter state
        const searchParams = {
          query: searchQuery.value,
          state: selectedState.value,
          type: selectedCategory.value
        }
        
        // Call the updated service method
        const templeData = await templeService.getTemples(searchParams)
        console.log('Temples fetched from API:', templeData)
        
        temples.value = templeData || []
        
        console.log(`Successfully processed ${temples.value.length} temples`)
      } catch (error) {
        console.error('Error fetching temple data:', error)
        showToast('Error loading temples. Please try again later.', 'error')
        temples.value = []
      } finally {
        loading.value = false
      }
    }

    // Clear all temple-related storage keys on component mount
    const clearTempleSelectionStorage = () => {
      // Clear all potential storage keys that might cause automatic redirection
      localStorage.removeItem('selectedEntityId');
      localStorage.removeItem('currentEntityId');
      localStorage.removeItem('entityId');
      localStorage.removeItem('selectedTempleId');
      localStorage.removeItem('defaultEntityId');
      localStorage.removeItem('lastVisitedEntityId');
      
      // Check session storage too
      sessionStorage.removeItem('selectedEntityId');
      sessionStorage.removeItem('currentEntityId');
      sessionStorage.removeItem('entityId');
      sessionStorage.removeItem('selectedTempleId');
      
      // This line is crucial - need to tell the router we're on the selection page
      sessionStorage.setItem('onTempleSelectionPage', 'true');
      
      console.log('🧹 Cleared all temple selection storage keys to prevent auto-redirect');
    }

    // On component mount
    onMounted(() => {
      // Clear ALL storage keys that might cause automatic redirection
      clearTempleSelectionStorage();
      
      // Load joined temples list (but don't use it for redirection)
      loadJoinedTemples();
      
      // Fetch temples from API
      fetchTemples();
      
      // Add protection against future redirects by watching storage changes
      window.addEventListener('storage', (event) => {
        if (event.key && 
            (event.key.includes('Entity') || 
             event.key.includes('entity') || 
             event.key.includes('Temple') || 
             event.key.includes('temple'))) {
          // Someone is trying to set a value that might cause redirection
          // Clear it immediately
          clearTempleSelectionStorage();
        }
      });
    });

    return {
      searchQuery,
      selectedState,
      selectedCategory,
      loading,
      joiningTemple,
      showConfirmModal,
      isJoining,
      selectedTemple,
      temples,
      joinedTemples,
      showSuccessModal,
      joinedTempleId,
      joinedTempleName,
      filteredTemples,
      resetFilters,
      selectTemple,
      closeModal,
      confirmJoinTemple,
      fetchTemples,
      goToProfileCreation,
      goToDashboard
    }
  }
}
</script>

<style scoped>
/* Import Google Fonts */
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap');

/* Custom font classes */
.font-heading {
  font-family: 'Roboto', sans-serif;
}

.font-side {
  font-family: sans-serif;
}

.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  line-clamp: 2; /* Standard property for future compatibility */
  overflow: hidden;
}

/* Animation for loading spinner */
@keyframes spin {
  to { transform: rotate(360deg); }
}
.animate-spin {
  animation: spin 1s linear infinite;
}
</style>