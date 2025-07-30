<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- Header with actions -->
    <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Events & Festivals</h2>
        <p class="text-gray-600 mt-1">Manage temple events and track RSVPs</p>
      </div>
    </div>

    <!-- Filters -->
    <div class="flex flex-col sm:flex-row gap-4 mb-6">
      <div class="flex-1">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search events..."
          class="w-full px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
        />
      </div>
      <select
        v-model="selectedStatus"
        class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
      >
        <option value="">All Status</option>
        <option value="upcoming">Upcoming</option>
        <option value="ongoing">Ongoing</option>
        <option value="completed">Completed</option>
        <option value="cancelled">Cancelled</option>
      </select>
      <select
        v-model="selectedCategory"
        class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
      >
        <option value="">All Categories</option>
        <option value="festival">Festival</option>
        <option value="ceremony">Ceremony</option>
        <option value="spiritual">Spiritual Program</option>
        <option value="cultural">Cultural Event</option>
        <option value="community">Community Service</option>
      </select>
    </div>

    <!-- Events List -->
    <div v-if="loading" class="flex justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <div v-else-if="filteredEvents.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
      </svg>
      <h3 class="mt-2 text-lg font-medium text-gray-900">No events found</h3>
      <p class="mt-1 text-gray-500">Get started by creating your first event.</p>
    </div>

    <div v-else class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
      <div
        v-for="event in filteredEvents"
        :key="event.id"
        class="bg-white border border-gray-200 rounded-xl shadow-sm hover:shadow-md transition-all duration-200 overflow-hidden"
      >
        <!-- Event Image -->
        <div class="relative h-48 bg-gradient-to-br from-indigo-400 to-indigo-600">
          <img
            v-if="event.image"
            :src="event.image"
            :alt="event.title"
            class="w-full h-full object-cover"
          />
          <div class="absolute inset-0 bg-black bg-opacity-20"></div>
          <div class="absolute top-4 right-4">
            <span
              :class="getStatusBadgeClass(getEventStatus(event))"
              class="px-3 py-1 rounded-full text-xs font-medium"
            >
              {{ getEventStatus(event) }}
            </span>
          </div>
        </div>

        <!-- Event Details -->
        <div class="p-6">
          <div class="flex items-start justify-between mb-3">
            <h3 class="text-lg font-semibold text-gray-900 line-clamp-2">
              {{ event.title }}
            </h3>
            <div class="flex gap-2 ml-2">
              <button
                @click="$emit('edit-event', event)"
                class="p-1 text-gray-400 hover:text-indigo-600 transition-colors duration-200"
                title="Edit Event"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                </svg>
              </button>
              <button
                @click="$emit('delete-event', event)"
                class="p-1 text-gray-400 hover:text-red-600 transition-colors duration-200"
                title="Delete Event"
              >
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                </svg>
              </button>
            </div>
          </div>

          <p class="text-gray-600 text-sm mb-4 line-clamp-2">
            {{ event.description }}
          </p>

          <!-- Event Meta Info -->
          <div class="space-y-2 mb-4">
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
              </svg>
              {{ formatDate(event.eventDate) }}
            </div>
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              {{ formatTime(event.eventDate) }}
            </div>
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
              {{ event.location }}
            </div>
          </div>

          <!-- Category & RSVP Stats -->
          <div class="flex items-center justify-between">
            <span
              :class="getCategoryBadgeClass(event.type)"
              class="px-2 py-1 rounded-md text-xs font-medium"
            >
              {{ formatEventType(event.type) }}
            </span>
            <div class="flex items-center text-sm text-gray-600">
              <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z" />
              </svg>
              {{ event.currentAttendees || 0 }} RSVPs
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex gap-2 mt-4">
            <button
              @click="$emit('view-rsvp', event)"
              class="flex-1 bg-indigo-50 text-indigo-700 hover:bg-indigo-100 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            >
              Manage RSVPs
            </button>
            <button
              @click="$emit('view-details', event)"
              class="flex-1 bg-gray-50 text-gray-700 hover:bg-gray-100 px-3 py-2 rounded-lg text-sm font-medium transition-all duration-200"
            >
              View Details
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex justify-center mt-8">
      <nav class="flex items-center gap-2">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          Previous
        </button>
        <span
          v-for="page in visiblePages"
          :key="page"
          @click="currentPage = page"
          :class="[
            'px-3 py-2 text-sm font-medium cursor-pointer transition-all duration-200',
            currentPage === page
              ? 'bg-indigo-600 text-white rounded-lg'
              : 'text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg'
          ]"
        >
          {{ page }}
        </span>
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="px-3 py-2 text-sm font-medium text-gray-600 hover:text-indigo-600 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          Next
        </button>
      </nav>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useEventStore } from '@/stores/event'

const props = defineProps({
  templeId: {
    type: String,
    required: true
  }
})

const emit = defineEmits([
  'create-event',
  'edit-event',
  'delete-event',
  'view-details',
  'view-rsvp'
])

const eventStore = useEventStore()

// Reactive data
const loading = ref(false)
const searchQuery = ref('')
const selectedStatus = ref('')
const selectedCategory = ref('')
const currentPage = ref(1)
const itemsPerPage = 9

// Computed properties
const filteredEvents = computed(() => {
  let events = eventStore.events

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    events = events.filter(event =>
      event.title.toLowerCase().includes(query) ||
      event.description.toLowerCase().includes(query) ||
      event.location.toLowerCase().includes(query)
    )
  }

  // Filter by status
  if (selectedStatus.value) {
    events = events.filter(event => getEventStatus(event).toLowerCase() === selectedStatus.value.toLowerCase())
  }

  // Filter by category
  if (selectedCategory.value) {
    events = events.filter(event => 
      (event.type || event.eventType)?.toLowerCase() === selectedCategory.value.toLowerCase()
    )
  }

  // Pagination
  const start = (currentPage.value - 1) * itemsPerPage
  const end = start + itemsPerPage
  return events.slice(start, end)
})

const totalPages = computed(() => {
  return Math.ceil(eventStore.events.length / itemsPerPage)
})

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, start + 4)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
const getStatusBadgeClass = (status) => {
  const classes = {
    'upcoming': 'bg-blue-100 text-blue-800',
    'ongoing': 'bg-green-100 text-green-800',
    'completed': 'bg-gray-100 text-gray-800',
    'cancelled': 'bg-red-100 text-red-800'
  }
  return classes[status.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

const getEventStatus = (event) => {
  if (!event || !event.eventDate) return 'Unknown'
  
  const now = new Date()
  const eventDate = new Date(event.eventDate)
  
  // If event has endDate, use it; otherwise assume it ends on the same day
  const endDate = event.endDate ? new Date(event.endDate) : new Date(eventDate)
  endDate.setHours(23, 59, 59) // End of the event day
  
  if (!event.isActive) return 'Cancelled'
  if (now < eventDate) return 'Upcoming'
  if (now >= eventDate && now <= endDate) return 'Ongoing'
  return 'Completed'
}

const getCategoryBadgeClass = (category) => {
  const classes = {
    'festival': 'bg-orange-100 text-orange-800',
    'ceremony': 'bg-purple-100 text-purple-800',
    'spiritual': 'bg-indigo-100 text-indigo-800',
    'cultural': 'bg-pink-100 text-pink-800',
    'community': 'bg-green-100 text-green-800'
  }
  return classes[category?.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

const formatEventType = (type) => {
  if (!type) return 'Other'
  
  // Convert snake_case or kebab-case to Title Case
  return type
    .replace(/_/g, ' ')
    .replace(/-/g, ' ')
    .replace(/\w\S*/g, (txt) => {
      return txt.charAt(0).toUpperCase() + txt.substr(1).toLowerCase()
    })
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('en-US', {
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  })
}

const loadEvents = async () => {
  loading.value = true
  try {
    await eventStore.fetchEvents()
  } catch (error) {
    console.error('Failed to load events:', error)
  } finally {
    loading.value = false
  }
}

// Reset page when filters change
watch([searchQuery, selectedStatus, selectedCategory], () => {
  currentPage.value = 1
})

// Lifecycle
onMounted(() => {
  loadEvents()
})
</script>