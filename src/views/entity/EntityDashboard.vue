<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header with Temple Info -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex items-center justify-between h-16">
          <div class="flex items-center space-x-4">
            <div class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
              <svg class="h-6 w-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-2 0H3m2-16l9-9 9 9"></path>
              </svg>
            </div>
            <div>
              <h1 class="text-xl font-semibold text-gray-900">{{ temple ? temple.name : 'Loading...' }}</h1>
              <p class="text-sm text-gray-500">{{ temple ? `${temple.city}, ${temple.state}` : 'Loading location...' }}</p>
            </div>
          </div>
          <div class="flex items-center space-x-3">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              Active
            </span>
            <button class="relative p-2 text-gray-400 hover:text-gray-500 transition-colors duration-200">
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19H3m6 0V5l-6 5m0 9h6"></path>
              </svg>
              <span v-if="notifications.length > 0" class="absolute -top-1 -right-1 h-4 w-4 bg-red-500 rounded-full flex items-center justify-center text-xs text-white">{{ notifications.length }}</span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Banner -->
      <div class="bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-2xl p-6 mb-8">
        <div class="flex items-center justify-between">
          <div>
            <h2 class="text-2xl font-bold text-white mb-2">Welcome back, {{ user ? user.fullName : 'Admin' }}!</h2>
            <p class="text-indigo-100">Here's what's happening at your temple today</p>
          </div>
          <div class="hidden md:block">
            <svg class="h-16 w-16 text-indigo-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-2 0H3m2-16l9-9 9 9"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Stats Grid -->
      <div v-if="loading" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div v-for="i in 4" :key="i" class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="animate-pulse flex items-center">
            <div class="h-12 w-12 bg-gray-200 rounded-lg mr-4"></div>
            <div class="space-y-2 flex-1">
              <div class="h-4 bg-gray-200 rounded w-1/3"></div>
              <div class="h-6 bg-gray-200 rounded w-1/4"></div>
              <div class="h-3 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- Registered Devotees -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-12 w-12 bg-blue-100 rounded-lg flex items-center justify-center">
                <svg class="h-6 w-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"></path>
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Registered Devotees</p>
              <p class="text-2xl font-bold text-gray-900">{{ dashboardData.devotees?.total || 0 }}</p>
              <p class="text-xs text-green-600">+{{ dashboardData.devotees?.newThisMonth || 0 }} this month</p>
            </div>
          </div>
        </div>

        <!-- Seva Bookings Today -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-12 w-12 bg-orange-100 rounded-lg flex items-center justify-center">
                <svg class="h-6 w-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Seva Bookings</p>
              <div class="flex items-baseline space-x-2">
                <p class="text-2xl font-bold text-gray-900">{{ dashboardData.sevas?.today || 0 }}</p>
                <span class="text-sm text-gray-500">today</span>
              </div>
              <p class="text-xs text-blue-600">{{ dashboardData.sevas?.thisMonth || 0 }} this month</p>
            </div>
          </div>
        </div>

        <!-- Donations This Month -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-12 w-12 bg-green-100 rounded-lg flex items-center justify-center">
                <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Month Donations</p>
              <p class="text-2xl font-bold text-gray-900">₹{{ formatCurrency(dashboardData.donations?.thisMonth || 0) }}</p>
              <p class="text-xs text-green-600">+{{ dashboardData.donations?.growth || 0 }}% from last month</p>
            </div>
          </div>
        </div>

        <!-- Upcoming Events -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-12 w-12 bg-purple-100 rounded-lg flex items-center justify-center">
                <svg class="h-6 w-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Upcoming Events</p>
              <p class="text-2xl font-bold text-gray-900">{{ dashboardData.events?.upcoming || 0 }}</p>
              <p class="text-xs text-purple-600">{{ dashboardData.events?.thisWeek || 0 }} this week</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column: Recent Activity & Quick Actions -->
        <div class="lg:col-span-2 space-y-8">
          <!-- Recent Notifications -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200">
            <!-- <div class="p-6 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900">Recent Notifications</h3>
                <router-link 
                  to="/notifications" 
                  class="text-sm text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200">
                  View All
                </router-link>
              </div>
            </div> -->
            <!-- <div class="p-6">
              <div v-if="notifications.length === 0" class="text-center py-8">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z M9 19H3m6 0V5l-6 5m0 9h6"></path>
                </svg>
                <p class="mt-4 text-sm text-gray-500">No new notifications</p>
              </div>
              <div v-else class="space-y-4">
                <div v-for="notification in notifications.slice(0, 5)" :key="notification.id" 
                     class="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors duration-200">
                  <div class="flex-shrink-0">
                    <div :class="getNotificationIconClass(notification.type)" class="h-8 w-8 rounded-full flex items-center justify-center">
                      <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getNotificationIconPath(notification.type)"></path>
                      </svg>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
                    <p class="text-sm text-gray-500">{{ notification.message }}</p>
                    <p class="text-xs text-gray-400 mt-1">{{ formatTimeAgo(notification.createdAt) }}</p>
                  </div>
                </div>
              </div>
            </div> -->
          </div>

          <!-- Quick Actions -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200">
            <div class="p-6 border-b border-gray-200">
              <h3 class="text-lg font-semibold text-gray-900">Quick Actions</h3>
            </div>
            <div class="p-6">
              <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                <router-link :to="`/entity/${entityId}/sevas`" 
                           class="flex items-center p-4 border border-gray-200 rounded-xl hover:border-indigo-300 hover:shadow-md transition-all duration-200 group">
                  <div class="h-10 w-10 bg-orange-100 rounded-lg flex items-center justify-center group-hover:bg-orange-200 transition-colors duration-200">
                    <svg class="h-5 w-5 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Manage Sevas</p>
                    <p class="text-xs text-gray-500">Create and approve sevas</p>
                  </div>
                </router-link>

                <router-link :to="`/entity/${entityId}/events`" 
                           class="flex items-center p-4 border border-gray-200 rounded-xl hover:border-indigo-300 hover:shadow-md transition-all duration-200 group">
                  <div class="h-10 w-10 bg-purple-100 rounded-lg flex items-center justify-center group-hover:bg-purple-200 transition-colors duration-200">
                    <svg class="h-5 w-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Create Event</p>
                    <p class="text-xs text-gray-500">Schedule temple events</p>
                  </div>
                </router-link>

                <router-link :to="`/entity/${entityId}/devotees`" 
                           class="flex items-center p-4 border border-gray-200 rounded-xl hover:border-indigo-300 hover:shadow-md transition-all duration-200 group">
                  <div class="h-10 w-10 bg-blue-100 rounded-lg flex items-center justify-center group-hover:bg-blue-200 transition-colors duration-200">
                    <svg class="h-5 w-5 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z"></path>
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">View Devotees</p>
                    <p class="text-xs text-gray-500">Manage devotee profiles</p>
                  </div>
                </router-link>

                <router-link :to="`/entity/${entityId}/communication`" 
                           class="flex items-center p-4 border border-gray-200 rounded-xl hover:border-indigo-300 hover:shadow-md transition-all duration-200 group">
                  <div class="h-10 w-10 bg-green-100 rounded-lg flex items-center justify-center group-hover:bg-green-200 transition-colors duration-200">
                    <svg class="h-5 w-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
                    </svg>
                  </div>
                  <div class="ml-3">
                    <p class="text-sm font-medium text-gray-900">Send Message</p>
                    <p class="text-xs text-gray-500">Communicate with devotees</p>
                  </div>
                </router-link>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column: Upcoming Events & Top Donors -->
        <div class="space-y-8">
          <!-- Upcoming Events -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200">
            <div class="p-6 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900">Upcoming Events</h3>
                <router-link 
                  :to="`/entity/${entityId}/events`" 
                  class="text-sm text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200">
                  View All
                </router-link>
              </div>
            </div>
            <div class="p-6">
              <div v-if="upcomingEvents.length === 0" class="text-center py-8">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
                <p class="mt-4 text-sm text-gray-500">No upcoming events</p>
                <router-link :to="`/entity/${entityId}/events`" 
                           class="mt-2 inline-flex items-center text-sm text-indigo-600 hover:text-indigo-700 font-medium">
                  Create Event
                </router-link>
              </div>
              <div v-else class="space-y-4">
                <div v-for="event in upcomingEvents.slice(0, 3)" :key="event.id" 
                     class="flex items-start space-x-3 p-3 rounded-lg hover:bg-gray-50 transition-colors duration-200">
                  <div class="flex-shrink-0">
                    <div class="h-10 w-10 bg-purple-100 rounded-lg flex items-center justify-center">
                      <span class="text-xs font-semibold text-purple-600">{{ formatEventDate(event.date) }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900">{{ event.title }}</p>
                    <p class="text-xs text-gray-500">{{ formatFullDate(event.date) }}</p>
                    <div class="flex items-center mt-1">
                      <span class="text-xs text-green-600">{{ event.rsvpCount }} RSVPs</span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Top Donors -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200">
            <div class="p-6 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900">Top Donors</h3>
                <router-link 
                  :to="`/entity/${entityId}/donations`" 
                  class="text-sm text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200">
                  View All
                </router-link>
              </div>
            </div>
            <div class="p-6">
              <div v-if="topDonors.length === 0" class="text-center py-8">
                <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
                </svg>
                <p class="mt-4 text-sm text-gray-500">No donations yet</p>
              </div>
              <div v-else class="space-y-4">
                <div v-for="(donor, index) in topDonors.slice(0, 5)" :key="donor.id" 
                     class="flex items-center space-x-3">
                  <div class="flex-shrink-0">
                    <div class="h-8 w-8 bg-gradient-to-br from-green-400 to-green-600 rounded-full flex items-center justify-center">
                      <span class="text-xs font-bold text-white">{{ index + 1 }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900">{{ donor.name }}</p>
                    <p class="text-xs text-gray-500">{{ donor.totalDonations }} donations</p>
                  </div>
                  <div class="text-sm font-semibold text-green-600">
                    ₹{{ formatCurrency(donor.amount) }}
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useTempleStore } from '@/stores/temple'
import { useToast } from '@/composables/useToast'
import api from '@/plugins/axios'
import { apiClient } from '@/plugins/axios'

const route = useRoute()
const entityId = computed(() => route.params.id)
const authStore = useAuthStore()
const templeStore = useTempleStore()
const toast = useToast()

// User info
const user = computed(() => authStore.user)

// Temple data
const temple = ref(null)
const loading = ref(true)
const dashboardData = ref({
  devotees: { total: 0, newThisMonth: 0 },
  sevas: { today: 0, thisMonth: 0 },
  donations: { thisMonth: 0, growth: 0 },
  events: { upcoming: 0, thisWeek: 0 }
})

// Notifications
const notifications = ref([])

// Events and donors - will be loaded from API
const upcomingEvents = ref([])
const topDonors = ref([])

// Load entity data
const loadEntityData = async () => {
  try {
    loading.value = true
    console.log('Loading entity data for ID:', entityId.value)
    
    // Try to get temple from store first
    const storedTemple = templeStore.getTempleById(entityId.value)
    
    if (storedTemple) {
      console.log('Found temple in store:', storedTemple)
      temple.value = storedTemple
    } else {
      // Fetch temple details from API
      console.log('Temple not in store, fetching from API')
      try {
        const response = await api.get(`/v1/entities/${entityId.value}`)
        
        if (response.data) {
          temple.value = response.data
        }
      } catch (err) {
        console.error('Failed to fetch temple details:', err)
        toast.error('Failed to load temple details')
      }
    }
    
    // Inside loadEntityData function in EntityDashboard.vue
// Fixed dashboard data fetching code:

// Fetch dashboard data from the new API endpoint
try {
  // Remove the /api prefix from the URL since it's likely already configured in Axios
  // or is being handled by the proxy configuration
  const dashboardResponse = await api.get(`/v1/entities/dashboard-summary?entity_id=${entityId.value}`)
  
  if (dashboardResponse && dashboardResponse.data) {
    // Map the backend response to the structure expected by the frontend
    dashboardData.value = {
      devotees: {
        total: dashboardResponse.data.registered_devotees.total || 0,
        newThisMonth: dashboardResponse.data.registered_devotees.this_month || 0
      },
      sevas: {
        today: dashboardResponse.data.seva_bookings.today || 0,
        thisMonth: dashboardResponse.data.seva_bookings.this_month || 0
      },
      donations: {
        thisMonth: dashboardResponse.data.month_donations.amount || 0,
        growth: dashboardResponse.data.month_donations.percent_change || 0
      },
      events: {
        upcoming: dashboardResponse.data.upcoming_events.total || 0,
        thisWeek: dashboardResponse.data.upcoming_events.this_week || 0
      }
    }
    
    console.log('Dashboard data loaded successfully:', dashboardData.value)
  }
} catch (err) {
  console.error('Failed to fetch dashboard data:', err)
  toast.error('Failed to load dashboard statistics')
  // Keep the default values if API fails
}
    
    // Fetch upcoming events
    try {
      const eventsResponse = await api.get(`/v1/events/upcoming?entity_id=${entityId.value}&limit=3`)
      
      if (eventsResponse && eventsResponse.data) {
        upcomingEvents.value = eventsResponse.data.map(event => ({
          ...event,
          date: new Date(event.event_date || event.date || Date.now())
        }))
      }
    } catch (err) {
      console.error('Failed to fetch upcoming events:', err)
      upcomingEvents.value = []
    }
    
    // Fetch top donors
    try {
      const donorsResponse = await api.get(`/v1/donations/top?entity_id=${entityId.value}&limit=5`)
      
      if (donorsResponse && donorsResponse.data) {
        topDonors.value = donorsResponse.data
      }
    } catch (err) {
      console.error('Failed to fetch top donors:', err)
      topDonors.value = []
    }
    
    // Fetch notifications
    try {
      const notificationsResponse = await api.get(`/v1/notifications?entity_id=${entityId.value}&limit=5`)
      
      if (notificationsResponse && notificationsResponse.data) {
        notifications.value = notificationsResponse.data.map(notification => ({
          ...notification,
          createdAt: new Date(notification.created_at || notification.timestamp || Date.now())
        }))
      }
    } catch (err) {
      console.error('Failed to fetch notifications:', err)
      notifications.value = []
    }
    
  } catch (err) {
    console.error('Error loading entity data:', err)
    toast.error('Failed to load entity data')
  } finally {
    loading.value = false
  }
}

// Helper functions
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount)
}

const formatTimeAgo = (date) => {
  if (!date) return ''
  const now = new Date()
  const diff = now - date
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (minutes < 60) {
    return `${minutes}m ago`
  } else if (hours < 24) {
    return `${hours}h ago`
  } else {
    return `${days}d ago`
  }
}

const formatEventDate = (date) => {
  if (!date) return ''
  return date.getDate().toString()
}

const formatFullDate = (date) => {
  if (!date) return ''
  return date.toLocaleDateString('en-IN', {
    weekday: 'short',
    month: 'short',
    day: 'numeric'
  })
}

const getNotificationIconClass = (type) => {
  const classes = {
    seva: 'bg-orange-100 text-orange-600',
    donation: 'bg-green-100 text-green-600',
    event: 'bg-purple-100 text-purple-600',
    devotee: 'bg-blue-100 text-blue-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getNotificationIconPath = (type) => {
  const paths = {
    seva: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
    donation: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1',
    event: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
    devotee: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.25 2.25 0 11-4.5 0 2.25 2.25 0 014.5 0z'
  }
  return paths[type] || 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
}

onMounted(() => {
  // Load dashboard data
  console.log('EntityDashboard mounted for entity:', entityId.value)
  loadEntityData()
})
</script>
