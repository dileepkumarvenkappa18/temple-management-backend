<template>
  <div class="bg-white rounded-xl shadow-md overflow-hidden">
    <div class="p-6 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900">Recent Temple Activities</h3>
    </div>
    
    <!-- Loading State -->
    <div v-if="loading" class="p-6 flex justify-center">
      <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
    </div>
    
    <!-- Error State -->
    <div v-else-if="error" class="p-6 text-center">
      <svg class="mx-auto h-12 w-12 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <p class="mt-2 text-sm text-red-600">{{ error }}</p>
      <button 
        @click="loadActivities"
        class="mt-4 px-4 py-2 bg-indigo-600 text-white rounded-md text-sm hover:bg-indigo-700 transition-colors duration-150"
      >
        Try Again
      </button>
    </div>
    
    <!-- Empty State -->
    <div v-else-if="!hasActivities" class="p-6 text-center">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2.586a1 1 0 00-.707.293l-2.414 2.414a1 1 0 01-.707.293h-3.172a1 1 0 01-.707-.293l-2.414-2.414A1 1 0 006.586 13H4" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No activities yet</h3>
      <p class="mt-1 text-sm text-gray-500">Temple activities will appear here once available.</p>
    </div>
    
    <!-- Activities List -->
    <div v-else class="divide-y divide-gray-200">
      <div 
        v-for="(activity, index) in sortedActivities" 
        :key="index"
        class="p-4 hover:bg-gray-50 transition-colors"
      >
        <div class="flex items-start space-x-3">
          <!-- Activity Icon -->
          <div class="flex-shrink-0">
            <div 
              :class="getActivityIconClass(activity.type)"
              class="w-10 h-10 rounded-full flex items-center justify-center"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path 
                  v-if="activity.type === 'event'" 
                  stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" 
                />
                <path 
                  v-else-if="activity.type === 'seva'" 
                  stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" 
                />
                <path 
                  v-else
                  stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9" 
                />
              </svg>
            </div>
          </div>
          
          <!-- Activity Content -->
          <div class="flex-1 min-w-0">
            <p class="text-sm font-medium text-gray-900">
              {{ activity.title }}
            </p>
            <p class="text-xs text-gray-500 mt-1">{{ formatDate(activity.date) }}</p>
            <p v-if="activity.description" class="text-xs text-gray-600 mt-1 line-clamp-2">
              {{ activity.description }}
            </p>
          </div>
          
          <!-- Activity Status -->
          <div class="ml-2">
            <span 
              v-if="activity.status" 
              :class="getStatusClass(activity.status)"
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
            >
              {{ formatStatus(activity.status) }}
            </span>
          </div>
        </div>
      </div>
    </div>
    
    <!-- View All Link -->
    <div v-if="hasActivities" class="p-4 bg-gray-50 text-center">
      <button 
        @click="$emit('view-all')"
        class="text-sm font-medium text-indigo-600 hover:text-indigo-800 transition-colors"
      >
        View All Activities
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, defineProps, defineEmits, watch } from 'vue'
import { useTempleActivities } from '@/composables/useTempleActivities'

const props = defineProps({
  entityId: {
    type: [String, Number],
    required: true
  },
  limit: {
    type: Number,
    default: 10
  },
  autoLoad: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['view-all', 'loaded', 'error'])

// Setup composable
const {
  loading,
  error,
  activities,
  fetchAllActivities
} = useTempleActivities()

// Computed properties
const hasActivities = computed(() => {
  return sortedActivities.value.length > 0
})

const sortedActivities = computed(() => {
  // Combine all activities
  const allActivities = [
    // Map events from backend format to our display format
    ...activities.events.map(event => ({
      ...event,
      type: 'event',
      title: event.title || event.name,
      date: event.event_date,
      description: event.description,
      status: event.status
    })),
    // Map sevas from backend format to our display format
    ...activities.sevas.map(seva => ({
      ...seva,
      type: 'seva',
      title: seva.name,
      date: seva.created_at, // Use created_at if no specific date
      description: seva.description,
      status: seva.active ? 'active' : 'inactive'
    }))
  ]
  
  // Sort by date (newest first)
  const sorted = allActivities.sort((a, b) => {
    return new Date(b.date) - new Date(a.date)
  })
  
  // Limit the number of activities
  return sorted.slice(0, props.limit)
})

// Methods
const loadActivities = async () => {
  try {
    await fetchAllActivities(props.entityId)
    emit('loaded', sortedActivities.value)
  } catch (err) {
    emit('error', err.message)
  }
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  
  try {
    const date = new Date(dateString)
    
    // Check if date is valid
    if (isNaN(date.getTime())) return 'Invalid date'
    
    // Format date
    return date.toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short',
      year: 'numeric'
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Invalid date'
  }
}

const formatStatus = (status) => {
  if (!status) return ''
  
  // Capitalize first letter
  return status.charAt(0).toUpperCase() + status.slice(1)
}

const getActivityIconClass = (type) => {
  const classes = {
    event: 'bg-indigo-100 text-indigo-600',
    seva: 'bg-green-100 text-green-600',
    announcement: 'bg-yellow-100 text-yellow-600'
  }
  
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getStatusClass = (status) => {
  const statusMap = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    completed: 'bg-blue-100 text-blue-800',
    cancelled: 'bg-red-100 text-red-800',
    confirmed: 'bg-green-100 text-green-800',
    upcoming: 'bg-indigo-100 text-indigo-800',
    active: 'bg-green-100 text-green-800',
    inactive: 'bg-gray-100 text-gray-800'
  }
  
  // Handle case insensitive
  const normalizedStatus = (status || '').toLowerCase()
  return statusMap[normalizedStatus] || 'bg-gray-100 text-gray-800'
}

// Watch for changes to entityId
watch(() => props.entityId, (newId) => {
  if (newId) {
    loadActivities()
  }
})

// Lifecycle hooks
onMounted(() => {
  if (props.autoLoad && props.entityId) {
    loadActivities()
  }
})
</script>