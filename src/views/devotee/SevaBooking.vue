<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-6">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-3xl font-bold text-gray-900">Book Seva</h1>
              <p class="mt-2 text-sm text-gray-600">Choose from available sevas and book your preferred time slots</p>
            </div>
            <div class="flex items-center space-x-4">
              <div class="text-right">
                <p class="text-sm text-gray-500">Current Temple</p>
                <p class="font-semibold text-indigo-600">{{ templeInfo.name }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filters and Search -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-8">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Search Sevas</label>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search by seva name..."
        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
      />
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Seva Type</label>
      <select
        v-model="selectedCategory"
        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
      >
        <option value="">Select Seva Type</option>
        <option value="daily">Daily Seva</option>
        <option value="special">Special Seva</option>
        <option value="festival">Festival Seva</option>
        <option value="personal">Personal Seva</option>
      </select>
    </div>
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">Price Range</label>
      <select
        v-model="priceRange"
        class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
      >
        <option value="">All Prices</option>
        <option value="0-500">₹0 - ₹500</option>
        <option value="500-1000">₹500 - ₹1,000</option>
        <option value="1000-5000">₹1,000 - ₹5,000</option>
        <option value="5000+">₹5,000+</option>
      </select>
    </div>
  </div>
</div>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
      </div>

      <!-- No Sevas Available -->
      <div v-else-if="filteredSevas.length === 0" class="text-center py-12">
        <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 11V9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"/>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No Sevas Available</h3>
        <p class="text-gray-500 mb-4">No sevas match your current filters or none have been added by the temple yet.</p>
        <button
          @click="clearFilters"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
        >
          Clear Filters
        </button>
      </div>

      <!-- Seva Cards Grid -->
      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        <div
          v-for="seva in filteredSevas"
          :key="seva.id"
          class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden hover:shadow-md transition-shadow"
        >
          <!-- Seva Image -->
          <div class="h-48 bg-gradient-to-br from-indigo-500 to-purple-600 relative">
            <img
              v-if="seva.image"
              :src="seva.image"
              :alt="seva.name"
              class="w-full h-full object-cover"
            />
            <div class="absolute inset-0 bg-black bg-opacity-20"></div>
            <div class="absolute top-4 right-4">
              <span
                :class="[
                  'px-3 py-1 rounded-full text-xs font-medium',
                  seva.seva_type === 'daily' ? 'bg-green-100 text-green-800' :
                  seva.seva_type === 'special' ? 'bg-blue-100 text-blue-800' :
                  seva.seva_type === 'festival' ? 'bg-purple-100 text-purple-800' :
                  'bg-yellow-100 text-yellow-800'
                ]"
              >
                {{ seva.seva_type }}
              </span>
            </div>
          </div>

          <!-- Seva Content -->
          <div class="p-6">
            <div class="flex justify-between items-start mb-3">
              <h3 class="text-lg font-semibold text-gray-900">{{ seva.name }}</h3>
              <div class="text-right">
                <p class="text-xl font-bold text-indigo-600">₹{{ seva.price.toLocaleString() }}</p>
                <p v-if="seva.originalPrice && seva.originalPrice > seva.price" class="text-sm text-gray-500 line-through">
                  ₹{{ seva.originalPrice.toLocaleString() }}
                </p>
              </div>
            </div>

            <p class="text-gray-600 text-sm mb-4 line-clamp-2">{{ seva.description }}</p>

            <!-- Seva Details -->
            <div class="space-y-2 mb-4">
              <div class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Duration: {{ seva.duration }} minutes
              </div>
              <div class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                Available: {{ seva.max_bookings_per_day - (bookingCounts[seva.id] || 0) }} slots
              </div>
              <div v-if="seva.date" class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                Date: {{ formatDate(seva.date) }}
              </div>
              <div v-if="seva.start_time && seva.end_time" class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                Time: {{ seva.start_time }} - {{ seva.end_time }}
              </div>
            </div>

            <!-- Action Button -->
            <button
              @click="directBookSeva(seva)"
              :disabled="(seva.max_bookings_per_day - (bookingCounts[seva.id] || 0)) === 0 || bookingInProgress[seva.id] || bookedSevas[seva.id]"
              :class="[
                'w-full py-3 px-4 rounded-lg font-medium transition-colors',
                bookingInProgress[seva.id] ? 'bg-gray-400 text-white cursor-wait' :
                bookedSevas[seva.id] ? 'bg-green-600 text-white cursor-default' :
                (seva.max_bookings_per_day - (bookingCounts[seva.id] || 0)) > 0
                  ? 'bg-indigo-600 text-white hover:bg-indigo-700'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              ]"
            >
              {{ bookingInProgress[seva.id] ? 'Booking...' :
                 bookedSevas[seva.id] ? 'Seva Booked' : 
                 (seva.max_bookings_per_day - (bookingCounts[seva.id] || 0)) > 0 ? 'Book Now' : 'Fully Booked' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSevaStore } from '@/stores/seva'
import { useToast } from '@/composables/useToast'
import sevaService from '@/services/seva.service'

const route = useRoute()
const router = useRouter()
const sevaStore = useSevaStore()
const toast = useToast()

// Reactive data
const loading = ref(true)
const searchQuery = ref('')
const selectedCategory = ref('')
const selectedDate = ref('')
const priceRange = ref('')
const bookingCounts = ref({}) // To track available slots
const bookedSevas = reactive({}) // Track booked sevas using reactive
const bookingInProgress = reactive({}) // Track booking in progress

// Temple info
const templeInfo = ref({
  name: 'Temple Name',
  location: 'Location'
})

// Computed properties
const today = computed(() => {
  const date = new Date()
  return date.toISOString().split('T')[0]
})

const filteredSevas = computed(() => {
  let filtered = sevaStore.sevas
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(seva =>
      seva.name?.toLowerCase().includes(query) ||
      seva.description?.toLowerCase().includes(query)
    )
  }
  
  if (selectedCategory.value) {
    filtered = filtered.filter(seva => seva.seva_type === selectedCategory.value)
  }
  
  if (priceRange.value) {
    filtered = filtered.filter(seva => {
      const price = seva.price
      switch (priceRange.value) {
        case '0-500': return price <= 500
        case '500-1000': return price > 500 && price <= 1000
        case '1000-5000': return price > 1000 && price <= 5000
        case '5000+': return price > 5000
        default: return true
      }
    })
  }
  
  return filtered
})

// Methods
const formatDate = (dateString) => {
  if (!dateString) return ''
  // Format: YYYY-MM-DD to DD MMM YYYY
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const clearFilters = () => {
  searchQuery.value = ''
  selectedCategory.value = ''
  selectedDate.value = ''
  priceRange.value = ''
}

// Direct booking function - works with the new service method
const directBookSeva = async (seva) => {
  // Don't do anything if already booked, in progress, or fully booked
  if (bookedSevas[seva.id] || 
      bookingInProgress[seva.id] || 
      (seva.max_bookings_per_day - (bookingCounts.value[seva.id] || 0)) <= 0) {
    return
  }
  
  // Set booking in progress for this seva
  bookingInProgress[seva.id] = true
  
  try {
    console.log('Sending booking request:', { seva_id: seva.id });
    
    // Pass only the seva ID to the service
    const response = await sevaService.bookSeva(seva.id);
    console.log('Booking response:', response);
    
    if (response.success) {
      // Mark this seva as booked
      bookedSevas[seva.id] = true;
      
      // Update booking count
      if (!bookingCounts.value[seva.id]) {
        bookingCounts.value[seva.id] = 0;
      }
      bookingCounts.value[seva.id] = bookingCounts.value[seva.id] + 1;
      
      toast.success('Seva booked successfully!');
    } else {
      toast.error(response.error || 'Failed to book seva');
    }
  } catch (error) {
    console.error('Booking failed:', error);
    console.error('Error details:', error.response?.data || 'No response data');
    
    toast.error('Booking failed. Please try again.');
  } finally {
    // Clear booking in progress state
    bookingInProgress[seva.id] = false;
  }
}

const fetchSevas = async () => {
  loading.value = true
  try {
    // Apply current filters when fetching
    const params = {
      search: searchQuery.value,
      seva_type: selectedCategory.value
    }
    
    await sevaStore.fetchSevas(params)
    
    // Only update booking counts if sevas were loaded successfully
    if (sevaStore.sevas.length > 0) {
      sevaStore.sevas.forEach(seva => {
        bookingCounts.value[seva.id] = Math.floor(Math.random() * seva.max_bookings_per_day)
      })
    }
    
    // Check if user has any booked sevas
    fetchUserBookings()
    
  } catch (error) {
    console.error('Failed to load sevas:', error)
    
    // Display more specific error message
    if (error.toString && error.toString().includes('Availability: unsupported relations')) {
      toast.error('Database schema issue. Please contact the administrator.')
    } else {
      toast.error('Failed to load sevas. Please try again.')
    }
  } finally {
    loading.value = false
  }
}

// Fetch user's bookings and update the UI
const fetchUserBookings = async () => {
  try {
    const response = await sevaService.getMyBookings()
    
    if (response.success && response.data) {
      // Mark sevas as booked if they exist in the user's bookings
      response.data.forEach(booking => {
        bookedSevas[booking.seva_id] = true
      })
    }
  } catch (error) {
    console.error('Failed to fetch user bookings:', error)
  }
}

// Lifecycle
onMounted(async () => {
  await fetchSevas()
})
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;

  /* Future standard support */
  line-clamp: 2;

  overflow: hidden;
}
</style>