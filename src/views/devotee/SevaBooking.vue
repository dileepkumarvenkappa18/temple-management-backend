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
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
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
            <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
            <select
              v-model="selectedCategory"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            >
              <option value="">All Categories</option>
              <option value="daily">Daily Pooja</option>
              <option value="special">Special Occasions</option>
              <option value="festival">Festival</option>
              <option value="donation">Donation</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date</label>
            <input
              v-model="selectedDate"
              type="date"
              :min="today"
              class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            />
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
                  seva.category === 'daily' ? 'bg-green-100 text-green-800' :
                  seva.category === 'special' ? 'bg-blue-100 text-blue-800' :
                  seva.category === 'festival' ? 'bg-purple-100 text-purple-800' :
                  'bg-yellow-100 text-yellow-800'
                ]"
              >
                {{ seva.category }}
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
                Duration: {{ seva.duration }}
              </div>
              <div class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                </svg>
                Available: {{ seva.availableSlots }} slots
              </div>
              <div v-if="seva.nextAvailable" class="flex items-center text-sm text-gray-600">
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                Next: {{ formatDate(seva.nextAvailable) }}
              </div>
            </div>

            <!-- Action Button -->
            <button
              @click="openBookingModal(seva)"
              :disabled="seva.availableSlots === 0"
              :class="[
                'w-full py-3 px-4 rounded-lg font-medium transition-colors',
                seva.availableSlots > 0
                  ? 'bg-indigo-600 text-white hover:bg-indigo-700'
                  : 'bg-gray-300 text-gray-500 cursor-not-allowed'
              ]"
            >
              {{ seva.availableSlots > 0 ? 'Book Now' : 'Fully Booked' }}
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Booking Modal -->
    <div
      v-if="showBookingModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      @click="closeBookingModal"
    >
      <div
        class="bg-white rounded-xl max-w-md w-full p-6"
        @click.stop
      >
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-semibold text-gray-900">Book {{ selectedSeva.name }}</h3>
          <button
            @click="closeBookingModal"
            class="text-gray-400 hover:text-gray-600 transition-colors"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <form @submit.prevent="bookSeva" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Preferred Date</label>
            <input
              v-model="bookingForm.date"
              type="date"
              :min="today"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Time Slot</label>
            <select
              v-model="bookingForm.timeSlot"
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Select time slot</option>
              <option v-for="slot in availableTimeSlots" :key="slot" :value="slot">
                {{ slot }}
              </option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Special Requests (Optional)</label>
            <textarea
              v-model="bookingForm.notes"
              rows="3"
              placeholder="Any special requests or notes..."
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none"
            ></textarea>
          </div>

          <!-- Price Summary -->
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex justify-between items-center">
              <span class="font-medium text-gray-900">Total Amount:</span>
              <span class="text-xl font-bold text-indigo-600">₹{{ selectedSeva.price?.toLocaleString() }}</span>
            </div>
          </div>

          <div class="flex space-x-3 pt-4">
            <button
              type="button"
              @click="closeBookingModal"
              class="flex-1 py-2 px-4 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="bookingLoading"
              class="flex-1 py-2 px-4 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors"
            >
              {{ bookingLoading ? 'Booking...' : 'Confirm Booking' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

// Reactive data
const loading = ref(true)
const bookingLoading = ref(false)
const searchQuery = ref('')
const selectedCategory = ref('')
const selectedDate = ref('')
const priceRange = ref('')
const showBookingModal = ref(false)
const selectedSeva = ref({})

// Temple info
const templeInfo = ref({
  name: 'Sri Raghavendra Swamy Temple',
  location: 'Mantralayam, Karnataka'
})

// Booking form
const bookingForm = ref({
  date: '',
  timeSlot: '',
  notes: ''
})

// Available time slots
const availableTimeSlots = ref([
  '6:00 AM - 7:00 AM',
  '7:00 AM - 8:00 AM',
  '10:00 AM - 11:00 AM',
  '6:00 PM - 7:00 PM',
  '7:00 PM - 8:00 PM'
])

// Mock seva data
const sevas = ref([
  {
    id: 1,
    name: 'Abhisheka',
    description: 'Sacred bath ceremony for the deity with milk, honey, and other auspicious items',
    price: 2100,
    originalPrice: 2500,
    category: 'special',
    duration: '45 minutes',
    availableSlots: 3,
    nextAvailable: '2025-06-18',
    image: null
  },
  {
    id: 2,
    name: 'Archana',
    description: 'Chanting of names and offering flowers to the deity',
    price: 101,
    category: 'daily',
    duration: '15 minutes',
    availableSlots: 12,
    nextAvailable: '2025-06-18',
    image: null
  },
  {
    id: 3,
    name: 'Annadana',
    description: 'Serving free meals to devotees and the needy',
    price: 5000,
    category: 'donation',
    duration: 'Full day',
    availableSlots: 2,
    nextAvailable: '2025-06-20',
    image: null
  },
  {
    id: 4,
    name: 'Kalyanotsava',
    description: 'Divine marriage ceremony celebration',
    price: 11000,
    category: 'festival',
    duration: '2 hours',
    availableSlots: 1,
    nextAvailable: '2025-06-25',
    image: null
  },
  {
    id: 5,
    name: 'Deepotsava',
    description: 'Lighting of sacred lamps ceremony',
    price: 501,
    category: 'special',
    duration: '30 minutes',
    availableSlots: 8,
    nextAvailable: '2025-06-19',
    image: null
  },
  {
    id: 6,
    name: 'Nitya Pooja',
    description: 'Daily worship and offerings to the deity',
    price: 51,
    category: 'daily',
    duration: '20 minutes',
    availableSlots: 0, // Fully booked
    nextAvailable: '2025-06-21',
    image: null
  }
])

// Computed properties
const today = computed(() => {
  const date = new Date()
  return date.toISOString().split('T')[0]
})

const filteredSevas = computed(() => {
  let filtered = sevas.value

  if (searchQuery.value) {
    filtered = filtered.filter(seva =>
      seva.name.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      seva.description.toLowerCase().includes(searchQuery.value.toLowerCase())
    )
  }

  if (selectedCategory.value) {
    filtered = filtered.filter(seva => seva.category === selectedCategory.value)
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

const openBookingModal = (seva) => {
  selectedSeva.value = seva
  bookingForm.value = {
    date: '',
    timeSlot: '',
    notes: ''
  }
  showBookingModal.value = true
}

const closeBookingModal = () => {
  showBookingModal.value = false
  selectedSeva.value = {}
}

const bookSeva = async () => {
  bookingLoading.value = true
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1500))
    
    // Show success message (you would integrate with your toast system)
    alert(`Seva "${selectedSeva.value.name}" booked successfully for ${bookingForm.value.date} at ${bookingForm.value.timeSlot}`)
    
    closeBookingModal()
    
    // Optionally redirect to booking history or dashboard
    // router.push(`/entity/${route.params.id}/devotee/dashboard`)
    
  } catch (error) {
    console.error('Booking failed:', error)
    alert('Booking failed. Please try again.')
  } finally {
    bookingLoading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  try {
    // Simulate loading sevas from API
    await new Promise(resolve => setTimeout(resolve, 1000))
    loading.value = false
  } catch (error) {
    console.error('Failed to load sevas:', error)
    loading.value = false
  }
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
