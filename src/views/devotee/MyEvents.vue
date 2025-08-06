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
        :userName="currentUser?.fullName || ''"
        :userRole="currentUser?.role || 'devotee'"
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
              <option value="past">Past Events</option>
            </select>
          </div>
        </div>
      </BaseCard>

      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center py-12">
        <BaseLoader />
      </div>

      <!-- Events Display - List View Only -->
      <div v-else-if="filteredEvents.length > 0" class="space-y-4">
        <BaseCard
          v-for="event in filteredEvents"
          :key="event.id"
          class="p-6 hover:shadow-lg transition-all duration-200 cursor-pointer"
          @click="openEventModal(event)"
        >
          <div class="flex flex-col">
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
                      {{ formatTime(event.time) }}
                    </div>
                    <div class="flex items-center gap-1">
                      <MapPin class="h-4 w-4" />
                      {{ event.venue || 'Temple Premises' }}
                    </div>
                  </div>
                </div>

                <!-- Status Badge -->
                <div class="flex flex-col gap-2 items-start sm:items-end">
                  <div class="flex items-center gap-2">
                    <span :class="getEventStatusClass(event)">
                      {{ getEventStatus(event) }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </BaseCard>
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
                {{ formatTime(selectedEvent.time) }}
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
  List, 
  X,
  Users
} from 'lucide-vue-next'
import { apiClient } from '@/plugins/axios'
import eventService from '@/services/event.service'

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
const viewMode = ref('list') // Default to list view and keep it fixed
const showEventModal = ref(false)
const selectedEvent = ref(null)

// Temple data
const templeName = ref('Sri Venkateswara Temple')
const events = ref([])
const userRSVPs = ref(new Set())

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
    case 'all':
    default:
      filtered = filtered.filter(e => e.status) // or just return filtered
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
    loading.value = true;
    
    // Use eventService instead of direct apiClient calls
    const data = await eventService.getEvents();
    console.log("🔍 Events Response:", data);
    
    const today = new Date();

    // Map API response to our UI model
    events.value = (Array.isArray(data) ? data : []).map(event => {
      const eventDate = new Date(event.event_date);
      let status = 'upcoming';
      if (!isNaN(eventDate)) {
        status = eventDate < today ? 'past' : 'upcoming';
      }

      return {
        id: event.id,
        title: event.title,
        description: event.description,
        date: event.event_date,
        time: event.event_time || null,
        venue: event.location || 'Temple Premises',
        status,
        attendees: event.rsvp_count || 0,
        eventType: event.event_type,
        maxAttendees: 200 // Default placeholder
      };
    });
    
    try {
      // Skip RSVP fetching for now as it's not working correctly
      // We'll implement this later when we have the correct API structure
      
      /* 
      // This doesn't work:
      const rsvpResponse = await apiClient.get('/api/v1/event-rsvps/my');
      console.log("🔍 User RSVPs:", rsvpResponse);
      
      if (rsvpResponse && rsvpResponse.data && Array.isArray(rsvpResponse.data)) {
        userRSVPs.value = new Set(
          rsvpResponse.data.map(rsvp => rsvp.event_id)
        );
      }
      */
    } catch (rsvpError) {
      console.error('Error fetching RSVPs:', rsvpError);
    }

    console.log('✅ Mapped Events:', events.value);
    console.log('✅ User RSVPs:', userRSVPs.value);

  } catch (error) {
    console.error('❌ Error fetching events:', error);
    showToast('Failed to load events', 'error');
  } finally {
    loading.value = false;
  }
};

const handleRSVP = async (eventId, isRSVP) => {
  try {
    if (isRSVP) {
      // Try to use eventService if available for RSVP
      if (typeof eventService.registerForEvent === 'function') {
        await eventService.registerForEvent(eventId);
        userRSVPs.value.add(eventId);
        showToast('RSVP confirmed successfully!', 'success');
      } else {
        // Fall back to showing error
        showToast('RSVP functionality not available', 'error');
      }
    } else {
      // There's no cancel RSVP endpoint in your API, so we'll show an error
      showToast('RSVP cancellation is not supported yet', 'error');
      return;
    }
  } catch (error) {
    console.error('Error updating RSVP:', error);
    showToast('Failed to update RSVP', 'error');
  }
}

const openEventModal = (event) => {
  selectedEvent.value = event;
  showEventModal.value = true;
}

const closeEventModal = () => {
  showEventModal.value = false;
  selectedEvent.value = null;
}

const getUserRSVP = (eventId) => {
  return userRSVPs.value.has(eventId);
}

const canRSVP = (event) => {
  return event.status === 'upcoming' && 
         (!event.maxAttendees || event.attendees < event.maxAttendees);
}

const getEventStatus = (event) => {
  if (event.status === 'past') return 'Completed';
  if (event.maxAttendees && event.attendees >= event.maxAttendees) return 'Full';
  return 'Open';
}

const getEventStatusClass = (event) => {
  const status = getEventStatus(event);
  const baseClasses = 'px-2 py-1 rounded-lg text-xs font-medium';
  
  switch (status) {
    case 'Completed':
      return `${baseClasses} bg-gray-100 text-gray-800`;
    case 'Full':
      return `${baseClasses} bg-red-100 text-red-800`;
    default:
      return `${baseClasses} bg-green-100 text-green-800`;
  }
}

const getRSVPStatusClass = (eventId) => {
  const baseClasses = 'px-3 py-1 rounded-lg text-sm font-medium';
  return getUserRSVP(eventId)
    ? `${baseClasses} bg-green-100 text-green-800`
    : `${baseClasses} bg-gray-100 text-gray-800`;
}

const formatDate = (dateString) => {
  if (!dateString) return 'No date available';
  
  const date = new Date(dateString);
  if (isNaN(date.getTime())) return dateString;
  
  return date.toLocaleDateString('en-IN', { 
    weekday: 'long', 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric' 
  });
}

const formatTime = (timeString) => {
  if (!timeString) return 'Time not specified';
  
  // Handle ISO-formatted datetime strings
  if (timeString.includes('T')) {
    const date = new Date(timeString);
    if (!isNaN(date.getTime())) {
      return date.toLocaleTimeString('en-IN', {
        hour: '2-digit',
        minute: '2-digit',
        hour12: true
      });
    }
  }
  
  // Handle time-only strings (HH:MM or HH:MM:SS)
  if (timeString.includes(':')) {
    const timeParts = timeString.split(':');
    if (timeParts.length >= 2) {
      const hours = parseInt(timeParts[0], 10);
      const minutes = parseInt(timeParts[1], 10);
      
      if (!isNaN(hours) && !isNaN(minutes)) {
        const date = new Date();
        date.setHours(hours, minutes, 0);
        
        return date.toLocaleTimeString('en-IN', {
          hour: '2-digit',
          minute: '2-digit',
          hour12: true
        });
      }
    }
  }
  
  // Return the original if we can't parse it
  return timeString;
}

const getEmptyStateMessage = () => {
  switch (selectedFilter.value) {
    case 'upcoming':
      return 'No upcoming events found. Check back later for new events.';
    case 'attending':
      return 'You haven\'t RSVP\'d to any events yet. Browse upcoming events to participate.';
    case 'past':
      return 'No past events to display.';
    default:
      return 'No events match your search criteria.';
  }
}

const resetFilters = () => {
  searchQuery.value = '';
  selectedFilter.value = 'all';
}

// Lifecycle
onMounted(() => {
  loadEvents();
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