<template>
    <!-- Breadcrumb Navigation -->
    <!-- <AppBreadcrumb 
      :items="[
        { label: 'Dashboard', to: `/entity/${$route.params.id}/devotee/dashboard` },
        { label: 'My Events', to: `/entity/${$route.params.id}/devotee/events` }
      ]" 
    /> -->

    <div class="space-y-6">
      <!-- Welcome Banner -->
      <WelcomeBanner 
        :user="currentUser"
        title="My Temple Events"
        :subtitle="`Explore upcoming events and festivals at ${templeName}`"
      />

      <!-- Event Filters & Actions -->
      <BaseCard class="p-6">
        <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
          <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center flex-1">
            <!-- Search Events -->
            <div class="relative flex-1 min-w-0">
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <Search class="h-5 w-5 text-gray-400" />
              </div>
              <input
                v-model="searchQuery"
                type="text"
                placeholder="Search events..."
                class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              />
            </div>

            <!-- Filter by Status -->
            <select
              v-model="selectedFilter"
              class="px-4 py-2 border border-gray-300 rounded-xl focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white transition-all duration-200"
            >
              <option value="all">All Events</option>
              <option value="upcoming">Upcoming</option>
              <option value="attending">My RSVPs</option>
              <option value="past">Past Events</option>
            </select>
          </div>

          <!-- View Toggle -->
          <div class="flex bg-gray-100 rounded-xl p-1">
            <button
              @click="viewMode = 'grid'"
              :class="[
                'px-3 py-1 rounded-lg text-sm font-medium transition-all duration-200',
                viewMode === 'grid' 
                  ? 'bg-white text-indigo-600 shadow-sm' 
                  : 'text-gray-600 hover:text-gray-800'
              ]"
            >
              <Grid3x3 class="h-4 w-4" />
            </button>
            <button
              @click="viewMode = 'list'"
              :class="[
                'px-3 py-1 rounded-lg text-sm font-medium transition-all duration-200',
                viewMode === 'list' 
                  ? 'bg-white text-indigo-600 shadow-sm' 
                  : 'text-gray-600 hover:text-gray-800'
              ]"
            >
              <List class="h-4 w-4" />
            </button>
          </div>
        </div>
      </BaseCard>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <BaseLoader />
      </div>

      <!-- Events Display -->
      <div v-else-if="filteredEvents.length > 0">
        <!-- Grid View -->
        <div v-if="viewMode === 'grid'" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          <EventCard
            v-for="event in filteredEvents"
            :key="event.id"
            :event="event"
            :show-rsvp="true"
            :user-rsvp="getUserRSVP(event.id)"
            @rsvp="handleRSVP"
            @view-details="openEventModal"
            class="transform hover:scale-105 transition-all duration-200"
          />
        </div>

        <!-- List View -->
        <div v-else class="space-y-4">
          <BaseCard
            v-for="event in filteredEvents"
            :key="event.id"
            class="p-6 hover:shadow-lg transition-all duration-200"
          >
            <div class="flex flex-col lg:flex-row gap-6">
              <!-- Event Image -->
              <div class="flex-shrink-0">
                <img
                  :src="event.image || '/api/placeholder/120/120'"
                  :alt="event.title"
                  class="w-full lg:w-32 h-48 lg:h-32 object-cover rounded-xl"
                />
              </div>

              <!-- Event Details -->
              <div class="flex-1 min-w-0">
                <div class="flex flex-col sm:flex-row gap-4 justify-between">
                  <div class="flex-1">
                    <h3 class="text-xl font-bold text-gray-900 mb-2">{{ event.title }}</h3>
                    <p class="text-gray-600 mb-3 line-clamp-2">{{ event.description }}</p>
                    
                    <div class="flex flex-wrap gap-4 text-sm text-gray-500">
                      <div class="flex items-center gap-1">
                        <Calendar class="h-4 w-4" />
                        {{ formatDate(event.date) }}
                      </div>
                      <div class="flex items-center gap-1">
                        <Clock class="h-4 w-4" />
                        {{ event.time }}
                      </div>
                      <div class="flex items-center gap-1">
                        <MapPin class="h-4 w-4" />
                        {{ event.venue || 'Temple Premises' }}
                      </div>
                    </div>
                  </div>

                  <!-- RSVP Actions -->
                  <div class="flex flex-col gap-2 items-start sm:items-end">
                    <div class="flex items-center gap-2">
                      <span :class="getEventStatusClass(event)">
                        {{ getEventStatus(event) }}
                      </span>
                    </div>
                    
                    <div class="flex gap-2">
                      <BaseButton
                        variant="outline"
                        size="sm"
                        @click="openEventModal(event)"
                      >
                        View Details
                      </BaseButton>
                      
                      <BaseButton
                        v-if="canRSVP(event)"
                        :variant="getUserRSVP(event.id) ? 'danger' : 'primary'"
                        size="sm"
                        @click="handleRSVP(event.id, !getUserRSVP(event.id))"
                      >
                        {{ getUserRSVP(event.id) ? 'Cancel RSVP' : 'RSVP' }}
                      </BaseButton>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </BaseCard>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-12">
        <BaseCard class="p-8">
          <Calendar class="h-16 w-16 text-gray-300 mx-auto mb-4" />
          <h3 class="text-xl font-semibold text-gray-900 mb-2">No Events Found</h3>
          <p class="text-gray-600 mb-6">
            {{ getEmptyStateMessage() }}
          </p>
          <BaseButton
            variant="outline"
            @click="resetFilters"
          >
            View All Events
          </BaseButton>
        </BaseCard>
      </div>

      <!-- My RSVPs Summary -->
      <DashboardWidget
        title="My Event Participation"
        class="mt-8"
      >
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-4">
          <div class="bg-indigo-50 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-indigo-600">{{ rsvpStats.attending }}</div>
            <div class="text-sm text-indigo-600">Attending</div>
          </div>
          <div class="bg-green-50 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-green-600">{{ rsvpStats.attended }}</div>
            <div class="text-sm text-green-600">Attended</div>
          </div>
          <div class="bg-blue-50 rounded-xl p-4 text-center">
            <div class="text-2xl font-bold text-blue-600">{{ rsvpStats.upcoming }}</div>
            <div class="text-sm text-blue-600">Upcoming</div>
          </div>
        </div>
      </DashboardWidget>
    </div>

    <!-- Event Details Modal -->
    <BaseModal
      v-if="selectedEvent"
      :show="showEventModal"
      @close="closeEventModal"
      size="lg"
    >
      <div class="p-6">
        <!-- Event Header -->
        <div class="flex justify-between items-start mb-6">
          <div>
            <h2 class="text-2xl font-bold text-gray-900 mb-2">{{ selectedEvent.title }}</h2>
            <div class="flex flex-wrap gap-4 text-sm text-gray-500">
              <div class="flex items-center gap-1">
                <Calendar class="h-4 w-4" />
                {{ formatDate(selectedEvent.date) }}
              </div>
              <div class="flex items-center gap-1">
                <Clock class="h-4 w-4" />
                {{ selectedEvent.time }}
              </div>
              <div class="flex items-center gap-1">
                <MapPin class="h-4 w-4" />
                {{ selectedEvent.venue || 'Temple Premises' }}
              </div>
            </div>
          </div>
          <button
            @click="closeEventModal"
            class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
          >
            <X class="h-6 w-6" />
          </button>
        </div>

        <!-- Event Image -->
        <div v-if="selectedEvent.image" class="mb-6">
          <img
            :src="selectedEvent.image"
            :alt="selectedEvent.title"
            class="w-full h-64 object-cover rounded-xl"
          />
        </div>

        <!-- Event Description -->
        <div class="mb-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-3">About This Event</h3>
          <p class="text-gray-700 leading-relaxed">{{ selectedEvent.description }}</p>
        </div>

        <!-- Attached Sevas -->
        <div v-if="selectedEvent.sevas && selectedEvent.sevas.length > 0" class="mb-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-3">Available Sevas</h3>
          <div class="grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div
              v-for="seva in selectedEvent.sevas"
              :key="seva.id"
              class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
            >
              <div>
                <div class="font-medium text-gray-900">{{ seva.name }}</div>
                <div class="text-sm text-gray-600">₹{{ seva.amount }}</div>
              </div>
              <BaseButton size="sm" variant="outline">
                Book Now
              </BaseButton>
            </div>
          </div>
        </div>

        <!-- RSVP Section -->
        <div class="flex flex-col sm:flex-row gap-4 justify-between items-start sm:items-center pt-6 border-t border-gray-200">
          <div>
            <div class="text-sm text-gray-600 mb-1">RSVP Status</div>
            <span :class="getRSVPStatusClass(selectedEvent.id)">
              {{ getUserRSVP(selectedEvent.id) ? 'Attending' : 'Not Attending' }}
            </span>
          </div>
          
          <div class="flex gap-3">
            <BaseButton
              variant="outline"
              @click="closeEventModal"
            >
              Close
            </BaseButton>
            <BaseButton
              v-if="canRSVP(selectedEvent)"
              :variant="getUserRSVP(selectedEvent.id) ? 'danger' : 'primary'"
              @click="handleRSVP(selectedEvent.id, !getUserRSVP(selectedEvent.id))"
            >
              {{ getUserRSVP(selectedEvent.id) ? 'Cancel RSVP' : 'RSVP Now' }}
            </BaseButton>
          </div>
        </div>
      </div>
    </BaseModal>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
import { useToast } from '@/composables/useToast'
import { 
  Calendar, 
  Clock, 
  MapPin, 
  Search, 
  Grid3x3, 
  List, 
  X,
  Users
} from 'lucide-vue-next'

// Layout Components
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import AppBreadcrumb from '@/components/layout/AppBreadcrumb.vue'

// Base Components
import BaseCard from '@/components/common/BaseCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseLoader from '@/components/common/BaseLoader.vue'
import BaseModal from '@/components/common/BaseModal.vue'

// Dashboard Components
import WelcomeBanner from '@/components/dashboard/WelcomeBanner.vue'
import DashboardWidget from '@/components/dashboard/DashboardWidget.vue'

// Event Components
import EventCard from '@/components/event/EventCard.vue'

const route = useRoute()
const { currentUser } = useAuth()
const { showToast } = useToast()

// Reactive state
const loading = ref(true)
const searchQuery = ref('')
const selectedFilter = ref('all')
const viewMode = ref('grid')
const showEventModal = ref(false)
const selectedEvent = ref(null)

// Mock data - replace with API calls
const templeName = ref('Sri Venkateswara Temple')
const events = ref([])
const userRSVPs = ref(new Set())

// Mock events data
const mockEvents = [
  {
    id: 1,
    title: 'Maha Shivaratri Celebration',
    description: 'Grand celebration of Lord Shiva with special poojas, abhishekam, and cultural programs throughout the night.',
    date: '2024-03-08',
    time: '6:00 PM - 6:00 AM',
    venue: 'Main Temple Hall',
    image: '/api/placeholder/400/250',
    status: 'upcoming',
    attendees: 150,
    maxAttendees: 200,
    sevas: [
      { id: 1, name: 'Rudrabhishek', amount: 251 },
      { id: 2, name: 'Maha Aarti', amount: 101 }
    ]
  },
  {
    id: 2,
    title: 'Holi Festival',
    description: 'Colorful celebration of Holi with traditional colors, sweets, and community gathering.',
    date: '2024-03-25',
    time: '10:00 AM - 2:00 PM',
    venue: 'Temple Courtyard',
    image: '/api/placeholder/400/250',
    status: 'upcoming',
    attendees: 80,
    maxAttendees: 100,
    sevas: [
      { id: 3, name: 'Holika Dahan', amount: 151 }
    ]
  },
  {
    id: 3,
    title: 'Ram Navami Celebration',
    description: 'Celebrating the birth of Lord Rama with special kirtans, bhajans, and prasadam distribution.',
    date: '2024-04-17',
    time: '9:00 AM - 12:00 PM',
    venue: 'Main Temple',
    status: 'upcoming',
    attendees: 95,
    maxAttendees: 120
  },
  {
    id: 4,
    title: 'Diwali Celebration',
    description: 'Festival of lights celebrated with diyas, rangoli, and special prayers.',
    date: '2023-11-12',
    time: '6:00 PM - 10:00 PM',
    venue: 'Temple Complex',
    status: 'past',
    attendees: 200,
    maxAttendees: 200
  }
]

// Initialize user RSVPs (mock data)
const mockUserRSVPs = new Set([1, 3]) // User has RSVP'd to events 1 and 3

// Computed properties
const filteredEvents = computed(() => {
  let filtered = events.value

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(event => 
      event.title.toLowerCase().includes(query) ||
      event.description.toLowerCase().includes(query)
    )
  }

  // Filter by status
  switch (selectedFilter.value) {
    case 'upcoming':
      filtered = filtered.filter(event => event.status === 'upcoming')
      break
    case 'attending':
      filtered = filtered.filter(event => userRSVPs.value.has(event.id))
      break
    case 'past':
      filtered = filtered.filter(event => event.status === 'past')
      break
  }

  return filtered
})

const rsvpStats = computed(() => {
  const attending = events.value.filter(event => 
    userRSVPs.value.has(event.id) && event.status === 'upcoming'
  ).length
  
  const attended = events.value.filter(event => 
    userRSVPs.value.has(event.id) && event.status === 'past'
  ).length
  
  const upcoming = events.value.filter(event => 
    event.status === 'upcoming'
  ).length

  return { attending, attended, upcoming }
})

// Methods
const loadEvents = async () => {
  try {
    loading.value = true
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    events.value = mockEvents
    userRSVPs.value = mockUserRSVPs
  } catch (error) {
    showToast('Failed to load events', 'error')
  } finally {
    loading.value = false
  }
}

const handleRSVP = async (eventId, isRSVP) => {
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    
    if (isRSVP) {
      userRSVPs.value.add(eventId)
      showToast('RSVP confirmed successfully!', 'success')
    } else {
      userRSVPs.value.delete(eventId)
      showToast('RSVP cancelled successfully!', 'success')
    }
  } catch (error) {
    showToast('Failed to update RSVP', 'error')
  }
}

const openEventModal = (event) => {
  selectedEvent.value = event
  showEventModal.value = true
}

const closeEventModal = () => {
  showEventModal.value = false
  selectedEvent.value = null
}

const getUserRSVP = (eventId) => {
  return userRSVPs.value.has(eventId)
}

const canRSVP = (event) => {
  return event.status === 'upcoming' && 
         (!event.maxAttendees || event.attendees < event.maxAttendees)
}

const getEventStatus = (event) => {
  if (event.status === 'past') return 'Completed'
  if (event.maxAttendees && event.attendees >= event.maxAttendees) return 'Full'
  return 'Open'
}

const getEventStatusClass = (event) => {
  const status = getEventStatus(event)
  const baseClasses = 'px-2 py-1 rounded-lg text-xs font-medium'
  
  switch (status) {
    case 'Completed':
      return `${baseClasses} bg-gray-100 text-gray-800`
    case 'Full':
      return `${baseClasses} bg-red-100 text-red-800`
    default:
      return `${baseClasses} bg-green-100 text-green-800`
  }
}

const getRSVPStatusClass = (eventId) => {
  const baseClasses = 'px-3 py-1 rounded-lg text-sm font-medium'
  return getUserRSVP(eventId)
    ? `${baseClasses} bg-green-100 text-green-800`
    : `${baseClasses} bg-gray-100 text-gray-800`
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  })
}

const getEmptyStateMessage = () => {
  switch (selectedFilter.value) {
    case 'upcoming':
      return 'No upcoming events found. Check back later for new events.'
    case 'attending':
      return 'You haven\'t RSVP\'d to any events yet. Browse upcoming events to participate.'
    case 'past':
      return 'No past events to display.'
    default:
      return 'No events match your search criteria.'
  }
}

const resetFilters = () => {
  searchQuery.value = ''
  selectedFilter.value = 'all'
}

// Lifecycle
onMounted(() => {
  loadEvents()
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