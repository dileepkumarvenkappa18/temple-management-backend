<template>
  <div>
    <!-- Page Header -->
    <div class="bg-white p-6 rounded-xl shadow-sm mb-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">My Seva Bookings</h1>
          <p class="text-gray-600 mt-1">View and manage all your seva bookings</p>
        </div>
        <router-link 
          :to="`/entity/${route.params.id}/devotee/seva-booking`"
          class="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
          </svg>
          Book New Seva
        </router-link>
      </div>
    </div>

    <!-- Filters -->
    <div class="bg-white p-6 rounded-xl shadow-sm mb-6">
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="filters.status"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          >
            <option value="all">All Statuses</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
          <select
            v-model="filters.dateRange"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          >
            <option value="all">All Time</option>
            <option value="today">Today</option>
            <option value="week">This Week</option>
            <option value="month">This Month</option>
            <option value="year">This Year</option>
          </select>
        </div>
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <div class="relative">
            <input
              v-model="filters.search"
              type="text"
              placeholder="Search seva name..."
              class="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
            />
            <svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="bg-white p-6 rounded-xl shadow-sm flex justify-center">
      <div class="flex flex-col items-center">
        <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600"></div>
        <p class="mt-2 text-gray-600">Loading your seva bookings...</p>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-white p-6 rounded-xl shadow-sm flex justify-center">
      <div class="flex flex-col items-center text-center">
        <svg class="h-10 w-10 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"></path>
        </svg>
        <h3 class="mt-2 text-lg font-medium text-gray-900">Error loading seva bookings</h3>
        <p class="mt-1 text-sm text-gray-500">{{ error }}</p>
        <button 
          @click="fetchSevaBookings"
          class="mt-3 inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none"
        >
          Try Again
        </button>
      </div>
    </div>

    <!-- Seva Bookings List -->
    <div v-else-if="filteredBookings.length > 0" class="bg-white p-6 rounded-xl shadow-sm">
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Seva
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Date & Time
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Price
              </th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                Status
              </th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="booking in filteredBookings" :key="booking.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-10 w-10 bg-indigo-100 rounded-full flex items-center justify-center">
                    <svg class="h-5 w-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-medium text-gray-900">{{ booking.seva_name || (booking.seva && booking.seva.name) || 'Seva' }}</div>
                    <div class="text-xs text-gray-500">{{ booking.seva_type || (booking.seva && booking.seva.type) || '' }}</div>
                  </div>
                </div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ formatDate(booking.booking_time || booking.BookingTime) }}</div>
                <div class="text-xs text-gray-500">{{ formatTime(booking.booking_time || booking.BookingTime) }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">₹{{ booking.price || booking.amount || 0 }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span 
                  :class="getStatusClass(booking.status || booking.Status)"
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                >
                  {{ formatStatus(booking.status || booking.Status) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button
                  @click="viewBookingDetails(booking)"
                  class="text-indigo-600 hover:text-indigo-900 mr-3"
                >
                  View
                </button>
                <button
                  v-if="(booking.status || booking.Status || '').toLowerCase() === 'pending'"
                  @click="cancelBooking(booking)"
                  class="text-red-600 hover:text-red-900"
                >
                  Cancel
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="bg-white p-8 rounded-xl shadow-sm text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
      </svg>
      <h3 class="mt-2 text-lg font-medium text-gray-900">No seva bookings found</h3>
      <p class="mt-1 text-sm text-gray-500">You haven't booked any sevas yet.</p>
      <div class="mt-6">
        <router-link 
          :to="`/entity/${route.params.id}/devotee/seva-booking`"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none"
        >
          Book Your First Seva
        </router-link>
      </div>
    </div>

    <!-- Booking Details Modal -->
    <BaseModal 
      :show="showBookingDetailsModal" 
      @close="showBookingDetailsModal = false"
      title="Booking Details"
    >
      <div v-if="selectedBooking" class="p-4">
        <div class="space-y-4">
          <div class="flex justify-between">
            <h3 class="text-lg font-medium text-gray-900">{{ selectedBooking.seva_name || (selectedBooking.seva && selectedBooking.seva.name) || 'Seva' }}</h3>
            <span 
              :class="getStatusClass(selectedBooking.status || selectedBooking.Status)" 
              class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
            >
              {{ formatStatus(selectedBooking.status || selectedBooking.Status) }}
            </span>
          </div>
          
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div>
              <p class="text-gray-500">Booking Date</p>
              <p class="font-medium">{{ formatDate(selectedBooking.booking_time || selectedBooking.BookingTime) }}</p>
            </div>
            <div>
              <p class="text-gray-500">Booking Time</p>
              <p class="font-medium">{{ formatTime(selectedBooking.booking_time || selectedBooking.BookingTime) }}</p>
            </div>
            <div>
              <p class="text-gray-500">Amount</p>
              <p class="font-medium">₹{{ selectedBooking.price || selectedBooking.amount || 0 }}</p>
            </div>
            <div>
              <p class="text-gray-500">Type</p>
              <p class="font-medium">{{ selectedBooking.seva_type || (selectedBooking.seva && selectedBooking.seva.type) || '' }}</p>
            </div>
            <div class="col-span-2">
              <p class="text-gray-500">Description</p>
              <p class="font-medium">{{ selectedBooking.description || (selectedBooking.seva && selectedBooking.seva.description) || 'No description available' }}</p>
            </div>
            <div v-if="selectedBooking.special_requests || selectedBooking.SpecialRequests" class="col-span-2">
              <p class="text-gray-500">Special Requests</p>
              <p class="font-medium">{{ selectedBooking.special_requests || selectedBooking.SpecialRequests }}</p>
            </div>
          </div>
        </div>
        
        <div class="mt-6 flex justify-end space-x-3">
          <button 
            @click="showBookingDetailsModal = false"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            Close
          </button>
          <button 
            v-if="(selectedBooking.status || selectedBooking.Status || '').toLowerCase() === 'pending'"
            @click="cancelBookingFromModal"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none"
          >
            Cancel Booking
          </button>
        </div>
      </div>
    </BaseModal>
    
    <!-- Confirm Cancel Modal -->
    <BaseModal 
      :show="showCancelConfirmModal" 
      @close="showCancelConfirmModal = false"
      title="Cancel Seva Booking"
    >
      <div class="p-4">
        <p class="text-gray-700">Are you sure you want to cancel this seva booking? This action cannot be undone.</p>
        <div class="mt-6 flex justify-end space-x-3">
          <button 
            @click="showCancelConfirmModal = false"
            class="px-4 py-2 border border-gray-300 rounded-md text-sm font-medium text-gray-700 hover:bg-gray-50"
          >
            No, Keep Booking
          </button>
          <button 
            @click="confirmCancelBooking"
            class="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-red-600 hover:bg-red-700 focus:outline-none"
          >
            Yes, Cancel Booking
          </button>
        </div>
      </div>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useSevaStore } from '@/stores/seva'
import { useToast } from '@/composables/useToast'
import BaseModal from '@/components/common/BaseModal.vue'
import { sevaService } from '@/services/seva.service'

const route = useRoute()
const router = useRouter()
const sevaStore = useSevaStore()
const toast = useToast()

// State
const loading = ref(false)
const error = ref(null)
const filters = ref({
  status: 'all',
  dateRange: 'all',
  search: ''
})
const selectedBooking = ref(null)
const showBookingDetailsModal = ref(false)
const showCancelConfirmModal = ref(false)

// Format helpers
const formatDate = (dateTime) => {
  if (!dateTime) return ''
  
  try {
    const date = new Date(dateTime)
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return ''
  }
}

const formatTime = (dateTime) => {
  if (!dateTime) return ''
  
  try {
    const date = new Date(dateTime)
    return date.toLocaleTimeString('en-US', {
      hour: '2-digit',
      minute: '2-digit',
      hour12: true
    })
  } catch (error) {
    console.error('Error formatting time:', error)
    return ''
  }
}

const formatStatus = (status) => {
  if (!status) return 'Pending'
  return status.charAt(0).toUpperCase() + status.slice(1).toLowerCase()
}

// Status styles
const getStatusClass = (status) => {
  if (!status) return 'bg-gray-100 text-gray-800'
  
  const statusMap = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    completed: 'bg-blue-100 text-blue-800',
    cancelled: 'bg-red-100 text-red-800',
    rejected: 'bg-red-100 text-red-800'
  }
  
  return statusMap[status.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

// Computed
const filteredBookings = computed(() => {
  let result = sevaStore.recentSevas || []
  
  // Log the data for debugging
  console.log('Seva bookings data:', result);
  
  // Filter by status
  if (filters.value.status !== 'all') {
    result = result.filter(booking => {
      const bookingStatus = (booking.status || booking.Status || '').toLowerCase()
      return bookingStatus === filters.value.status.toLowerCase()
    })
  }
  
  // Filter by date range
  if (filters.value.dateRange !== 'all') {
    const now = new Date()
    const today = new Date(now.getFullYear(), now.getMonth(), now.getDate())
    const startOfWeek = new Date(now)
    startOfWeek.setDate(now.getDate() - now.getDay()) // Sunday
    startOfWeek.setHours(0, 0, 0, 0)
    const startOfMonth = new Date(now.getFullYear(), now.getMonth(), 1)
    const startOfYear = new Date(now.getFullYear(), 0, 1)
    
    result = result.filter(booking => {
      const bookingDate = new Date(booking.booking_time || booking.BookingTime || new Date())
      
      switch (filters.value.dateRange) {
        case 'today':
          return bookingDate >= today
        case 'week':
          return bookingDate >= startOfWeek
        case 'month':
          return bookingDate >= startOfMonth
        case 'year':
          return bookingDate >= startOfYear
        default:
          return true
      }
    })
  }
  
  // Filter by search
  if (filters.value.search) {
    const search = filters.value.search.toLowerCase()
    result = result.filter(booking => {
      const sevaName = (booking.seva_name || (booking.seva && booking.seva.name) || '').toLowerCase()
      const sevaType = (booking.seva_type || (booking.seva && booking.seva.type) || '').toLowerCase()
      const specialRequests = (booking.special_requests || booking.SpecialRequests || '').toLowerCase()
      
      return sevaName.includes(search) || 
             sevaType.includes(search) || 
             specialRequests.includes(search)
    })
  }
  
  return result
})

// Methods
const fetchSevaBookings = async () => {
  loading.value = true
  error.value = null
  
  try {
    await sevaStore.fetchRecentSevas()
    
    // After fetching, retrieve catalog/seva names
    if (sevaStore.fetchSevaCatalog) {
      await sevaStore.fetchSevaCatalog()
    }
    
    console.log("Fetched seva bookings:", sevaStore.recentSevas)
  } catch (err) {
    console.error('Error fetching seva bookings:', err)
    error.value = err.message || 'Failed to load seva bookings'
  } finally {
    loading.value = false
  }
}

const viewBookingDetails = (booking) => {
  selectedBooking.value = booking
  showBookingDetailsModal.value = true
}

const cancelBooking = (booking) => {
  selectedBooking.value = booking
  showCancelConfirmModal.value = true
}

const cancelBookingFromModal = () => {
  showBookingDetailsModal.value = false
  showCancelConfirmModal.value = true
}

const confirmCancelBooking = async () => {
  try {
    const bookingId = selectedBooking.value.id || selectedBooking.value.ID
    
    // Call API to cancel the booking
    await sevaService.updateBookingStatus(bookingId, 'cancelled')
    
    // Refresh the list
    await fetchSevaBookings()
    
    toast.success('Seva booking cancelled successfully')
    showCancelConfirmModal.value = false
    showBookingDetailsModal.value = false
  } catch (err) {
    console.error('Error cancelling booking:', err)
    toast.error('Failed to cancel booking: ' + (err.message || 'Unknown error'))
  }
}

// Call your API to get the real seva catalog data and match it with bookings
const fetchSevaNamesForBookings = async () => {
  if (sevaStore.fetchSevaCatalog) {
    try {
      await sevaStore.fetchSevaCatalog()
    } catch (error) {
      console.error("Error fetching seva catalog:", error)
    }
  }
}

// Lifecycle hooks
onMounted(() => {
  fetchSevaBookings()
  fetchSevaNamesForBookings()
  
  // Log store state for debugging
  console.log("Initial sevaStore state:", {
    recentSevas: sevaStore.recentSevas,
    sevaCatalog: sevaStore.sevaCatalog,
    sevaList: sevaStore.sevaList
  })
})
</script>