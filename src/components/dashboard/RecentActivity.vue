<template>
  <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100">
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-semibold text-gray-900 flex items-center">
        <svg class="w-5 h-5 text-indigo-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        Recent Activity
      </h3>
      <button class="text-sm text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200">
        View All
      </button>
    </div>

    <div class="space-y-4">
      <!-- Loading State -->
      <div v-if="isLoading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="animate-pulse">
          <div class="flex items-start space-x-3">
            <div class="w-8 h-8 bg-gray-200 rounded-full"></div>
            <div class="flex-1 space-y-2">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-3 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Activity Items -->
      <div v-else-if="activities.length > 0" class="space-y-4">
        <div 
          v-for="activity in activities" 
          :key="activity.id"
          class="flex items-start space-x-3 p-3 hover:bg-gray-50 rounded-lg transition-colors duration-200"
        >
          <!-- Activity Icon -->
          <div class="flex-shrink-0">
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center"
              :class="getActivityIconClass(activity.type)"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path 
                  stroke-linecap="round" 
                  stroke-linejoin="round" 
                  stroke-width="2" 
                  :d="getActivityIconPath(activity.type)"
                ></path>
              </svg>
            </div>
          </div>

          <!-- Activity Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-900 font-medium">{{ activity.title }}</p>
              <span class="text-xs text-gray-500">{{ formatTimeAgo(activity.timestamp) }}</span>
            </div>
            <p class="text-sm text-gray-600 mt-1">{{ activity.description }}</p>
            
            <!-- Activity Metadata -->
            <div v-if="activity.metadata" class="flex items-center mt-2 space-x-4">
              <span 
                v-if="activity.metadata.user" 
                class="text-xs text-gray-500 flex items-center"
              >
                <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
                {{ activity.metadata.user }}
              </span>
              <span 
                v-if="activity.metadata.amount" 
                class="text-xs font-medium"
                :class="activity.type === 'donation' ? 'text-green-600' : 'text-gray-600'"
              >
                ₹{{ formatAmount(activity.metadata.amount) }}
              </span>
            </div>
          </div>

          <!-- Activity Action -->
          <div v-if="activity.actionable" class="flex-shrink-0">
            <button 
              class="text-xs text-indigo-600 hover:text-indigo-800 font-medium px-2 py-1 rounded hover:bg-indigo-50 transition-all duration-200"
              @click="handleActivityAction(activity)"
            >
              {{ activity.actionText || 'View' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-8">
        <svg class="w-12 h-12 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <p class="text-gray-500 text-sm">No recent activity to show</p>
        <p class="text-gray-400 text-xs mt-1">Activities will appear here as they happen</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useRoute, useRouter } from 'vue-router'

// Props
const props = defineProps({
  limit: {
    type: Number,
    default: 10
  },
  entityId: {
    type: String,
    default: null
  }
})

// Composables
const { user } = useAuth()
const route = useRoute()
const router = useRouter()

// Reactive data
const isLoading = ref(false)
const activities = ref([])

// Mock activity data based on user role
const mockActivities = computed(() => {
  const baseActivities = [
    {
      id: 1,
      type: 'devotee_registered',
      title: 'New Devotee Registration',
      description: 'Ramesh Kumar completed profile registration',
      timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2 hours ago
      metadata: { user: 'Ramesh Kumar' },
      actionable: true,
      actionText: 'View Profile'
    },
    {
      id: 2,
      type: 'seva_booking',
      title: 'Seva Booking',
      description: 'Ganesh Chaturthi Pooja booked for tomorrow',
      timestamp: new Date(Date.now() - 4 * 60 * 60 * 1000), // 4 hours ago
      metadata: { user: 'Priya Sharma' },
      actionable: true,
      actionText: 'View Booking'
    },
    {
      id: 3,
      type: 'donation',
      title: 'Donation Received',
      description: 'Temple development fund contribution',
      timestamp: new Date(Date.now() - 6 * 60 * 60 * 1000), // 6 hours ago
      metadata: { user: 'Anonymous Devotee', amount: 5000 },
      actionable: true,
      actionText: 'View Receipt'
    },
    {
      id: 4,
      type: 'event_created',
      title: 'New Event Created',
      description: 'Diwali Celebration 2025 event published',
      timestamp: new Date(Date.now() - 8 * 60 * 60 * 1000), // 8 hours ago
      metadata: { user: 'Temple Admin' },
      actionable: true,
      actionText: 'View Event'
    },
    {
      id: 5,
      type: 'volunteer_assigned',
      title: 'Volunteer Assignment',
      description: 'Kitchen seva volunteer assigned for weekend',
      timestamp: new Date(Date.now() - 12 * 60 * 60 * 1000), // 12 hours ago
      metadata: { user: 'Suresh Patel' },
      actionable: true,
      actionText: 'View Assignment'
    }
  ]

  // Filter activities based on user role
  if (user.value?.role === 'devotee') {
    return baseActivities.filter(activity => 
      ['seva_booking', 'donation', 'event_created'].includes(activity.type)
    ).slice(0, 3)
  } else if (user.value?.role === 'volunteer') {
    return baseActivities.filter(activity => 
      ['volunteer_assigned', 'event_created', 'seva_booking'].includes(activity.type)
    ).slice(0, 3)
  }

  return baseActivities
})

// Methods
const getActivityIconClass = (type) => {
  const iconClasses = {
    devotee_registered: 'bg-blue-100 text-blue-600',
    seva_booking: 'bg-purple-100 text-purple-600',
    donation: 'bg-green-100 text-green-600',
    event_created: 'bg-orange-100 text-orange-600',
    volunteer_assigned: 'bg-indigo-100 text-indigo-600',
    default: 'bg-gray-100 text-gray-600'
  }
  return iconClasses[type] || iconClasses.default
}

const getActivityIconPath = (type) => {
  const iconPaths = {
    devotee_registered: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
    seva_booking: 'M8 7V3a1 1 0 011-1h6a1 1 0 011 1v4m-6 0h6m-6 0V3m6 4v4m-6-4v4m0 0H5a1 1 0 00-1 1v4a1 1 0 001 1h14a1 1 0 001-1v-4a1 1 0 00-1-1h-2',
    donation: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1',
    event_created: 'M8 7V3a1 1 0 011-1h6a1 1 0 011 1v4m-6 0h6m-6 0V3m6 4v4m-6-4v4',
    volunteer_assigned: 'M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z',
    default: 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
  }
  return iconPaths[type] || iconPaths.default
}

const formatTimeAgo = (timestamp) => {
  const now = new Date()
  const diff = now - new Date(timestamp)
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  return `${days}d ago`
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount)
}

const handleActivityAction = (activity) => {
  // Get the temple ID from the route or props
  const entityId = props.entityId || route.params.id || route.params.entityId || '1'
  
  // Navigate based on activity type
  if (activity.type === 'donation') {
    router.push({
      name: 'DevoteeDonationHistory',
      params: { id: entityId }
    })
  } else if (activity.type === 'seva_booking') {
    router.push({
      name: 'SevaBooking',
      params: { id: entityId }
    })
  } else if (activity.type === 'event_created') {
    router.push({
      name: 'DevoteeMyEvents',
      params: { id: entityId }
    })
  } else {
    console.log('Activity action clicked:', activity)
  }
}

const loadActivities = async () => {
  isLoading.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    activities.value = mockActivities.value.slice(0, props.limit)
  } catch (error) {
    console.error('Failed to load activities:', error)
  } finally {
    isLoading.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadActivities()
})
</script>