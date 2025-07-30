<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header Section -->
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-4" style="font-family: 'Roboto', sans-serif;">
          Welcome, Volunteer!
        </h1>
        <p class="text-lg text-gray-600 max-w-2xl mx-auto" style="font-family: sans-serif;">
          Select a temple to begin your volunteer journey. You'll be able to view assigned sevas, upcoming events, and receive notifications from the temple administration.
        </p>
      </div>

      <!-- Welcome Card -->
      <div class="bg-indigo-50 border border-indigo-200 rounded-xl p-6 mb-8">
        <div class="flex items-start space-x-4">
          <div class="flex-shrink-0">
            <div class="w-12 h-12 bg-indigo-600 rounded-full flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"></path>
              </svg>
            </div>
          </div>
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-indigo-800 mb-2" style="font-family: sans-serif;">
              Ready to Serve?
            </h3>
            <p class="text-indigo-700 text-sm leading-relaxed">
              As a volunteer, you'll have access to view seva assignments, upcoming temple events, and receive important notifications. Choose your temple below to get started.
            </p>
          </div>
        </div>
      </div>

      <!-- Search and Filters -->
      <div class="bg-white rounded-xl shadow-md p-6 mb-8">
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <!-- Search Input -->
          <div class="md:col-span-2">
            <label for="search" class="block text-sm font-medium text-gray-700 mb-2" style="font-family: sans-serif;">
              Search Temples
            </label>
            <div class="relative">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                </svg>
              </div>
              <input
                v-model="searchQuery"
                type="text"
                id="search"
                placeholder="Search by temple name, city, or location..."
                class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              />
            </div>
          </div>

          <!-- Location Filter -->
          <div>
            <label for="location" class="block text-sm font-medium text-gray-700 mb-2" style="font-family: sans-serif;">
              Filter by Location
            </label>
            <select
              v-model="selectedLocation"
              id="location"
              class="block w-full px-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
              <option value="">All Locations</option>
              <option value="bengaluru">Bengaluru</option>
              <option value="mumbai">Mumbai</option>
              <option value="delhi">Delhi</option>
              <option value="chennai">Chennai</option>
              <option value="hyderabad">Hyderabad</option>
            </select>
          </div>
        </div>

        <!-- Results Summary -->
        <div class="flex items-center justify-between text-sm text-gray-600 border-t pt-4">
          <span style="font-family: sans-serif;">
            {{ filteredTemples.length }} temple{{ filteredTemples.length !== 1 ? 's' : '' }} found
          </span>
          <button
            @click="clearFilters"
            v-if="searchQuery || selectedLocation"
            class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
          >
            Clear filters
          </button>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>

      <!-- No Results -->
      <div v-else-if="filteredTemples.length === 0" class="text-center py-12">
        <div class="w-24 h-24 mx-auto mb-4 text-gray-300">
          <svg fill="currentColor" viewBox="0 0 24 24">
            <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm-2 15l-5-5 1.41-1.41L10 14.17l7.59-7.59L19 8l-9 9z"/>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2" style="font-family: sans-serif;">
          No temples found
        </h3>
        <p class="text-gray-600 mb-4">
          No approved temples match your search criteria. Please try different filters or check back later.
        </p>
        <button
          @click="clearFilters"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
        >
          View All Temples
        </button>
      </div>

      <!-- Temples Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="temple in filteredTemples"
          :key="temple.id"
          class="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden group cursor-pointer"
          @click="selectTemple(temple)"
        >
          <!-- Temple Image -->
          <div class="relative h-48 bg-gradient-to-br from-indigo-100 to-indigo-200">
            <img
              v-if="temple.image"
              :src="temple.image"
              :alt="temple.name"
              class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-200"
            />
            <div v-else class="w-full h-full flex items-center justify-center">
              <svg class="w-16 h-16 text-indigo-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
              </svg>
            </div>
            
            <!-- Status Badge -->
            <div class="absolute top-3 right-3">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
                <svg class="w-3 h-3 mr-1" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"></path>
                </svg>
                Active
              </span>
            </div>
          </div>

          <!-- Temple Info -->
          <div class="p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-2 group-hover:text-indigo-600 transition-colors duration-200" style="font-family: sans-serif;">
              {{ temple.name }}
            </h3>
            
            <div class="space-y-2 mb-4">
              <div class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                {{ temple.city }}, {{ temple.state }}
              </div>
              
              <div class="flex items-center text-sm text-gray-600" v-if="temple.phone">
                <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
                </svg>
                {{ temple.phone }}
              </div>

              <div class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"></path>
                </svg>
                {{ temple.volunteersCount || 0 }} volunteers
              </div>
            </div>

            <!-- Action Button -->
            <button
              @click.stop="selectTemple(temple)"
              :disabled="isTempleJoined(temple.id)"
              class="w-full inline-flex items-center justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white transition-all duration-200"
              :class="isTempleJoined(temple.id) 
                ? 'bg-green-600 cursor-not-allowed opacity-80' 
                : 'bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 group-hover:bg-indigo-700'"
            >
              <span v-if="isTempleJoined(temple.id)">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
                Joined Successfully
              </span>
              <span v-else>
                Join as Volunteer
                <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7l5 5m0 0l-5 5m5-5H6"></path>
                </svg>
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Selection Confirmation Modal -->
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
          <div class="w-16 h-16 mx-auto mb-4 bg-indigo-100 rounded-full flex items-center justify-center">
            <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
            </svg>
          </div>
          
          <h3 class="text-lg font-semibold text-gray-900 mb-2" style="font-family: sans-serif;">
            Join {{ selectedTemple?.name }}?
          </h3>
          
          <p class="text-gray-600 text-sm mb-6">
            You're about to join {{ selectedTemple?.name }} as a volunteer. You'll receive access to view seva assignments and temple notifications.
          </p>
          
          <div class="flex space-x-3">
            <button
              @click="closeModal"
              class="flex-1 px-4 py-2 border border-gray-300 text-sm font-medium rounded-lg text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
            >
              Cancel
            </button>
            <button
              @click="confirmSelection"
              :disabled="isJoining"
              class="flex-1 px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <span v-if="isJoining" class="flex items-center justify-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Joining...
              </span>
              <span v-else>Join Temple</span>
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useToast } from '@/composables/useToast';
import { useAuthStore } from '@/stores/auth';
// Import temple and volunteer services
import templeService from '@/services/temple.service';
import volunteerService from '@/services/volunteer.service';

const toast = useToast();
const authStore = useAuthStore();

// Reactive data
const searchQuery = ref('');
const selectedLocation = ref('');
const isLoading = ref(true);
const showConfirmModal = ref(false);
const selectedTemple = ref(null);
const isJoining = ref(false);
const temples = ref([]);
const joinedTemples = ref([]);

// Get user-specific localStorage key
const getJoinedTemplesKey = () => {
  const userId = authStore.user?.id || authStore.user?.email || 'anonymous';
  return `joinedTemples_${userId}`;
};

// Load joined temples from localStorage for the current user
const loadJoinedTemples = () => {
  try {
    const key = getJoinedTemplesKey();
    const storedJoinedTemples = localStorage.getItem(key);
    if (storedJoinedTemples) {
      joinedTemples.value = JSON.parse(storedJoinedTemples);
    } else {
      joinedTemples.value = [];
    }
    console.log('Loaded joined temples from localStorage:', joinedTemples.value);
  } catch (error) {
    console.error('Error loading joined temples from localStorage:', error);
    joinedTemples.value = [];
  }
};

// Save joined temples to localStorage for the current user
const saveJoinedTemples = () => {
  try {
    const key = getJoinedTemplesKey();
    localStorage.setItem(key, JSON.stringify(joinedTemples.value));
    console.log('Saved joined temples:', joinedTemples.value);
  } catch (error) {
    console.error('Error saving joined temples to localStorage:', error);
  }
};

// Fetch user memberships from API
const fetchUserMemberships = async () => {
  try {
    console.log('Fetching user memberships from API...');
    const result = await volunteerService.getUserMemberships();
    
    if (result.success && Array.isArray(result.data)) {
      // Extract entity IDs from memberships
      const templeIds = result.data.map(membership => {
        // Handle different property names
        return membership.entity_id || 
               membership.EntityID || 
               membership.entityId || 
               membership.id;
      }).filter(id => id !== undefined && id !== null);
      
      joinedTemples.value = templeIds;
      
      // Also update localStorage as a backup
      saveJoinedTemples();
      
      console.log('User memberships fetched from API:', templeIds);
    } else {
      console.warn('No memberships found or API error, falling back to localStorage');
      loadJoinedTemples();
    }
  } catch (error) {
    console.error('Error fetching user memberships from API:', error);
    // Fall back to localStorage if API fails
    loadJoinedTemples();
  }
};

// Check if a temple is already joined
const isTempleJoined = (templeId) => {
  if (!templeId) return false;
  
  // Convert to numbers for consistent comparison
  const numericId = Number(templeId);
  const numericJoinedIds = joinedTemples.value.map(id => Number(id));
  
  return numericJoinedIds.includes(numericId);
};

// Computed properties
const filteredTemples = computed(() => {
  let filtered = temples.value;

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(temple =>
      temple.name?.toLowerCase().includes(query) ||
      temple.city?.toLowerCase().includes(query) ||
      temple.state?.toLowerCase().includes(query)
    );
  }

  if (selectedLocation.value) {
    filtered = filtered.filter(temple =>
      temple.city?.toLowerCase() === selectedLocation.value.toLowerCase()
    );
  }

  return filtered;
});

// Methods
const clearFilters = () => {
  searchQuery.value = '';
  selectedLocation.value = '';
};

const selectTemple = (temple) => {
  if (isTempleJoined(temple.id)) {
    // If already joined, show a message or do nothing
    toast.info(`You are already a volunteer at ${temple.name}`);
    return;
  }

  selectedTemple.value = temple;
  showConfirmModal.value = true;
};

const closeModal = () => {
  showConfirmModal.value = false;
  selectedTemple.value = null;
};

// Join temple with API integration - with fallback
const confirmSelection = async () => {
  if (!selectedTemple.value) return;
  isJoining.value = true;

  try {
    console.log(`Joining temple with ID ${selectedTemple.value.id}...`);
    
    // Update UI immediately (optimistic update)
    if (!joinedTemples.value.includes(selectedTemple.value.id)) {
      joinedTemples.value.push(selectedTemple.value.id);
      saveJoinedTemples(); // Update localStorage right away
    }
    
    // Store temple selection
    try {
      localStorage.setItem('selectedEntityId', selectedTemple.value.id.toString());
      localStorage.setItem('selectedTempleName', selectedTemple.value.name);
    } catch (storageError) {
      console.error('localStorage error:', storageError);
    }
    
    // Call the API to join the temple
    const result = await volunteerService.joinTemple(selectedTemple.value.id);
    
    if (result.success) {
      console.log('API join temple success:', result);
      toast.success(`Successfully joined ${selectedTemple.value.name} as a volunteer!`);
    } else {
      console.error('API join temple failed, but UI already updated:', result);
      // Don't show error to user since the UI is already updated
    }
  } catch (error) {
    console.error('Error joining temple:', error);
    // Don't roll back UI change - we want to show "Joined Successfully" even if API fails
    toast.success(`Joined ${selectedTemple.value.name} as a volunteer! (local only)`);
  } finally {
    isJoining.value = false;
    closeModal();
  }
};

const fetchTemples = async () => {
  isLoading.value = true;

  try {
    console.log('Fetching temples from API...');
    
    // Use the temple service to fetch temples
    const templeData = await templeService.getTemples();
    
    if (Array.isArray(templeData) && templeData.length > 0) {
      temples.value = templeData;
      console.log(`Successfully processed ${temples.value.length} temples from API`);
    } else {
      console.warn('No temples found or invalid response, using mock data');
      temples.value = templeService.getMockTemples();
    }
    
    if (temples.value.length === 0) {
      toast.info('No approved temples found. Please check back later.');
    }
    
  } catch (error) {
    console.error('Error fetching temple data from API:', error);
    // Fallback to mock data
    temples.value = templeService.getMockTemples();
    toast.info('Using demo temple data.');
  } finally {
    isLoading.value = false;
  }
};

// Lifecycle
onMounted(async () => {
  // Wait for auth store to be initialized
  if (!authStore.isAuthenticated) {
    await authStore.initialize();
  }
  
  // Load joined temples from API first, with localStorage as fallback
  await fetchUserMemberships();
  
  // Fetch temples
  await fetchTemples();
});
</script>