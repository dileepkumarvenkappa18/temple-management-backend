<!-- src/components/devotee/DevoteeCard.vue -->
<template>
  <div class="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden">
    <!-- Header with Avatar and Status -->
    <div class="p-6 pb-4">
      <div class="flex items-start justify-between">
        <div class="flex items-center space-x-4">
          <!-- Avatar -->
          <div class="relative">
            <div v-if="devotee.avatar" class="w-16 h-16 rounded-full overflow-hidden">
              <img :src="devotee.avatar" :alt="devotee.name" class="w-full h-full object-cover">
            </div>
            <div v-else class="w-16 h-16 rounded-full bg-indigo-100 flex items-center justify-center">
              <span class="text-indigo-600 font-semibold text-xl">
                {{ getInitials(devotee.name) }}
              </span>
            </div>
            <!-- Status Indicator -->
            <div class="absolute -bottom-1 -right-1">
              <div :class="statusClasses" class="w-5 h-5 rounded-full border-2 border-white"></div>
            </div>
          </div>
          
          <!-- Basic Info -->
          <div class="flex-1">
            <h3 class="text-lg font-semibold text-gray-900">{{ devotee.name }}</h3>
            <p class="text-sm text-gray-600">{{ devotee.email }}</p>
            <p class="text-sm text-gray-600">{{ devotee.phone }}</p>
          </div>
        </div>
        
        <!-- Actions Dropdown -->
        <div class="relative" v-if="showActions">
          <button
            @click="showDropdown = !showDropdown"
            class="p-2 text-gray-400 hover:text-gray-600 rounded-lg hover:bg-gray-50 transition-colors"
          >
            <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
              <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z"/>
            </svg>
          </button>
          
          <!-- Dropdown Menu -->
          <div v-if="showDropdown" class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 z-10">
            <div class="py-1">
              <button
                @click="$emit('view', devotee)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
              >
                View Profile
              </button>
              <button
                @click="$emit('edit', devotee)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Edit Details
              </button>
              <button
                @click="$emit('message', devotee)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-50 transition-colors"
              >
                Send Message
              </button>
              <hr class="my-1">
              <button
                @click="$emit('deactivate', devotee)"
                class="w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors"
              >
                {{ devotee.status === 'active' ? 'Deactivate' : 'Activate' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Profile Completion Progress -->
    <div class="px-6 pb-4" v-if="devotee.profileCompletion !== undefined">
      <div class="flex items-center justify-between mb-2">
        <span class="text-sm font-medium text-gray-700">Profile Completion</span>
        <span class="text-sm font-semibold text-indigo-600">{{ devotee.profileCompletion }}%</span>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div 
          class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
          :style="{ width: `${devotee.profileCompletion}%` }"
        ></div>
      </div>
    </div>
    
    <!-- Details Grid -->
    <div class="px-6 pb-4">
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
        <!-- Gotra -->
        <div v-if="devotee.gotra" class="flex items-center space-x-2">
          <span class="text-gray-500">Gotra:</span>
          <span class="text-gray-900 font-medium">{{ devotee.gotra }}</span>
        </div>
        
        <!-- Native Place -->
        <div v-if="devotee.nativePlace" class="flex items-center space-x-2">
          <span class="text-gray-500">Native Place:</span>
          <span class="text-gray-900 font-medium">{{ devotee.nativePlace }}</span>
        </div>
        
        <!-- Nakshatra -->
        <div v-if="devotee.nakshatra" class="flex items-center space-x-2">
          <span class="text-gray-500">Nakshatra:</span>
          <span class="text-gray-900 font-medium">{{ devotee.nakshatra }}</span>
        </div>
        
        <!-- Rashi -->
        <div v-if="devotee.rashi" class="flex items-center space-x-2">
          <span class="text-gray-500">Rashi:</span>
          <span class="text-gray-900 font-medium">{{ devotee.rashi }}</span>
        </div>
      </div>
    </div>
    
    <!-- Recent Activity -->
    <div class="px-6 pb-4" v-if="devotee.recentActivity && devotee.recentActivity.length > 0">
      <h4 class="text-sm font-medium text-gray-700 mb-2">Recent Activity</h4>
      <div class="space-y-1">
        <div 
          v-for="activity in devotee.recentActivity.slice(0, 2)" 
          :key="activity.id"
          class="flex items-center space-x-2 text-xs text-gray-600"
        >
          <div class="w-1.5 h-1.5 bg-indigo-500 rounded-full"></div>
          <span>{{ activity.description }}</span>
          <span class="text-gray-400">{{ formatDate(activity.date) }}</span>
        </div>
      </div>
    </div>
    
    <!-- Stats Footer -->
    <div class="bg-gray-50 px-6 py-4">
      <div class="grid grid-cols-3 gap-4 text-center">
        <div>
          <div class="text-lg font-semibold text-indigo-600">{{ devotee.totalSevas || 0 }}</div>
          <div class="text-xs text-gray-600">Sevas</div>
        </div>
        <div>
          <div class="text-lg font-semibold text-green-600">₹{{ formatAmount(devotee.totalDonations || 0) }}</div>
          <div class="text-xs text-gray-600">Donations</div>
        </div>
        <div>
          <div class="text-lg font-semibold text-purple-600">{{ devotee.eventsAttended || 0 }}</div>
          <div class="text-xs text-gray-600">Events</div>
        </div>
      </div>
    </div>
    
    <!-- Join Date -->
    <div class="px-6 py-3 bg-gray-50 border-t border-gray-100">
      <div class="flex items-center justify-between text-xs text-gray-600">
        <span>Joined on {{ formatDate(devotee.joinedDate) }}</span>
        <span v-if="devotee.lastActive" class="text-gray-500">
          Last active {{ formatRelativeTime(devotee.lastActive) }}
        </span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  devotee: {
    type: Object,
    required: true
  },
  showActions: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['view', 'edit', 'message', 'deactivate'])

const showDropdown = ref(false)

const statusClasses = computed(() => {
  switch (props.devotee.status) {
    case 'active':
      return 'bg-green-500'
    case 'inactive':
      return 'bg-gray-400'
    case 'pending':
      return 'bg-yellow-500'
    default:
      return 'bg-gray-400'
  }
})

// Helper functions
const getInitials = (name) => {
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-IN', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const formatRelativeTime = (date) => {
  const now = new Date()
  const past = new Date(date)
  const diffInSeconds = Math.floor((now - past) / 1000)
  
  if (diffInSeconds < 60) return 'just now'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`
  if (diffInSeconds < 2592000) return `${Math.floor(diffInSeconds / 86400)}d ago`
  
  return formatDate(date)
}

const formatAmount = (amount) => {
  if (amount >= 100000) {
    return (amount / 100000).toFixed(1) + 'L'
  } else if (amount >= 1000) {
    return (amount / 1000).toFixed(1) + 'K'
  }
  return amount.toString()
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>