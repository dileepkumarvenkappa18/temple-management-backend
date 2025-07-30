<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
    <!-- Pending Tenants Card -->
    <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-200 border border-gray-100">
      <div class="flex items-center justify-between mb-4">
        <div class="bg-amber-100 p-3 rounded-xl">
          <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <span class="bg-amber-100 text-amber-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
          Pending
        </span>
      </div>
      <div class="space-y-2">
        <h3 class="text-2xl font-bold text-gray-900 font-roboto">{{ stats.pendingTenants }}</h3>
        <p class="text-sm text-gray-600 font-sans">Tenant Approvals</p>
        <div class="flex items-center text-xs text-amber-600">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
          </svg>
          Requires action
        </div>
      </div>
    </div>

    <!-- Active Temples Card -->
    <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-200 border border-gray-100">
      <div class="flex items-center justify-between mb-4">
        <div class="bg-indigo-100 p-3 rounded-xl">
          <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
          </svg>
        </div>
        <span class="bg-green-100 text-green-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
          Active
        </span>
      </div>
      <div class="space-y-2">
        <h3 class="text-2xl font-bold text-gray-900 font-roboto">{{ stats.activeTemples }}</h3>
        <p class="text-sm text-gray-600 font-sans">Active Temples</p>
        <div class="flex items-center text-xs text-green-600">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
          </svg>
          {{ stats.newTemplesThisMonth }}+ this month
        </div>
      </div>
    </div>

    <!-- Total Users Card -->
    <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-200 border border-gray-100">
      <div class="flex items-center justify-between mb-4">
        <div class="bg-purple-100 p-3 rounded-xl">
          <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"/>
          </svg>
        </div>
        <span class="bg-purple-100 text-purple-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
          {{ formatGrowth(stats.userGrowthPercent) }}%
        </span>
      </div>
      <div class="space-y-2">
        <h3 class="text-2xl font-bold text-gray-900 font-roboto">{{ formatNumber(stats.totalUsers) }}</h3>
        <p class="text-sm text-gray-600 font-sans">Total Users</p>
        <div class="flex items-center text-xs text-purple-600">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M3.293 9.707a1 1 0 010-1.414l6-6a1 1 0 011.414 0l6 6a1 1 0 01-1.414 1.414L11 5.414V17a1 1 0 11-2 0V5.414L4.707 9.707a1 1 0 01-1.414 0z" clip-rule="evenodd"/>
          </svg>
          +{{ stats.newUsersThisWeek }} this week
        </div>
      </div>
    </div>

    <!-- System Health Card -->
    <div class="bg-white rounded-2xl shadow-lg p-6 hover:shadow-xl transition-all duration-200 border border-gray-100">
      <div class="flex items-center justify-between mb-4">
        <div class="bg-emerald-100 p-3 rounded-xl">
          <svg class="w-6 h-6 text-emerald-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div class="flex items-center">
          <div class="w-2 h-2 bg-emerald-500 rounded-full mr-2 animate-pulse"></div>
          <span class="bg-emerald-100 text-emerald-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
            Healthy
          </span>
        </div>
      </div>
      <div class="space-y-2">
        <h3 class="text-2xl font-bold text-gray-900 font-roboto">{{ stats.systemUptime }}%</h3>
        <p class="text-sm text-gray-600 font-sans">System Uptime</p>
        <div class="flex items-center text-xs text-emerald-600">
          <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-8.293l-3-3a1 1 0 00-1.414 0l-3 3a1 1 0 001.414 1.414L9 9.414V13a1 1 0 102 0V9.414l1.293 1.293a1 1 0 001.414-1.414z" clip-rule="evenodd"/>
          </svg>
          {{ stats.avgResponseTime }}ms avg response
        </div>
      </div>
    </div>
  </div>

  <!-- Quick Stats Row -->
  <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-6">
    <div class="bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-xl p-4 text-white">
      <div class="text-center">
        <h4 class="text-xl font-bold font-roboto">{{ stats.pendingTemples }}</h4>
        <p class="text-indigo-100 text-sm font-sans">Pending Temples</p>
      </div>
    </div>
    
    <div class="bg-gradient-to-r from-cyan-500 to-cyan-600 rounded-xl p-4 text-white">
      <div class="text-center">
        <h4 class="text-xl font-bold font-roboto">{{ formatNumber(stats.totalDonations) }}</h4>
        <p class="text-cyan-100 text-sm font-sans">Total Donations</p>
      </div>
    </div>
    
    <div class="bg-gradient-to-r from-pink-500 to-pink-600 rounded-xl p-4 text-white">
      <div class="text-center">
        <h4 class="text-xl font-bold font-roboto">{{ stats.activeSevas }}</h4>
        <p class="text-pink-100 text-sm font-sans">Active Sevas</p>
      </div>
    </div>
    
    <div class="bg-gradient-to-r from-violet-500 to-violet-600 rounded-xl p-4 text-white">
      <div class="text-center">
        <h4 class="text-xl font-bold font-roboto">{{ stats.upcomingEvents }}</h4>
        <p class="text-violet-100 text-sm font-sans">Upcoming Events</p>
      </div>
    </div>
  </div>

  <!-- Recent Activity Summary -->
  <div class="bg-white rounded-2xl shadow-lg p-6 border border-gray-100">
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-bold text-gray-900 font-roboto">Recent Platform Activity</h3>
      <button class="text-indigo-600 hover:text-indigo-800 text-sm font-medium font-sans transition-colors duration-200">
        View All
      </button>
    </div>
    
    <div class="space-y-4">
      <div v-for="activity in recentActivities" :key="activity.id" 
           class="flex items-center justify-between p-3 bg-gray-50 rounded-xl hover:bg-gray-100 transition-colors duration-200">
        <div class="flex items-center space-x-3">
          <div :class="getActivityIconClass(activity.type)" class="p-2 rounded-lg">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path v-if="activity.type === 'approval'" fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd"/>
              <path v-else-if="activity.type === 'registration'" fill-rule="evenodd" d="M10 9a3 3 0 100-6 3 3 0 000 6zm-7 9a7 7 0 1114 0H3z" clip-rule="evenodd"/>
              <path v-else fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
            </svg>
          </div>
          <div>
            <p class="text-sm font-medium text-gray-900 font-sans">{{ activity.title }}</p>
            <p class="text-xs text-gray-500 font-sans">{{ activity.description }}</p>
          </div>
        </div>
        <span class="text-xs text-gray-400 font-sans">{{ formatTimeAgo(activity.timestamp) }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Props (if needed to pass data from parent)
const props = defineProps({
  refreshTrigger: {
    type: Number,
    default: 0
  }
})

// Reactive data
const stats = ref({
  pendingTenants: 12,
  activeTemples: 45,
  totalUsers: 2847,
  systemUptime: 99.8,
  newTemplesThisMonth: 8,
  userGrowthPercent: 15.2,
  newUsersThisWeek: 127,
  avgResponseTime: 145,
  pendingTemples: 7,
  totalDonations: 1250000,
  activeSevas: 234,
  upcomingEvents: 18
})

const recentActivities = ref([
  {
    id: 1,
    type: 'approval',
    title: 'Temple Approved',
    description: 'Sri Venkateshwara Temple - Bangalore approved',
    timestamp: new Date(Date.now() - 1000 * 60 * 15) // 15 minutes ago
  },
  {
    id: 2,
    type: 'registration',
    title: 'New Tenant Registration',
    description: 'Admin from Shirdi Sai Temple registered',
    timestamp: new Date(Date.now() - 1000 * 60 * 45) // 45 minutes ago
  },
  {
    id: 3,
    type: 'approval',
    title: 'Temple Approved',
    description: 'Hanuman Temple - Chennai approved',
    timestamp: new Date(Date.now() - 1000 * 60 * 120) // 2 hours ago
  },
  {
    id: 4,
    type: 'info',
    title: 'System Maintenance',
    description: 'Scheduled maintenance completed successfully',
    timestamp: new Date(Date.now() - 1000 * 60 * 180) // 3 hours ago
  }
])

const loading = ref(false)

// Computed properties
const formattedStats = computed(() => {
  return {
    ...stats.value,
    totalDonations: formatCurrency(stats.value.totalDonations)
  }
})

// Methods
const formatNumber = (num) => {
  if (num >= 1000000) {
    return (num / 1000000).toFixed(1) + 'M'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'K'
  }
  return num.toString()
}

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-IN', {
    style: 'currency',
    currency: 'INR',
    maximumFractionDigits: 0
  }).format(amount)
}

const formatGrowth = (percent) => {
  return percent > 0 ? `+${percent}` : percent.toString()
}

const formatTimeAgo = (timestamp) => {
  const now = new Date()
  const diff = now - timestamp
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(diff / (1000 * 60 * 60))
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days > 0) {
    return `${days}d ago`
  } else if (hours > 0) {
    return `${hours}h ago`
  } else if (minutes > 0) {
    return `${minutes}m ago`
  } else {
    return 'Just now'
  }
}

const getActivityIconClass = (type) => {
  const baseClass = 'text-white'
  switch (type) {
    case 'approval':
      return `bg-green-500 ${baseClass}`
    case 'registration':
      return `bg-indigo-500 ${baseClass}`
    case 'info':
      return `bg-blue-500 ${baseClass}`
    default:
      return `bg-gray-500 ${baseClass}`
  }
}

const fetchStats = async () => {
  loading.value = true
  try {
    // Mock API call - Replace with actual API
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // This would be replaced with actual API call
    // const response = await api.get('/superadmin/stats')
    // stats.value = response.data
    
    console.log('Stats refreshed')
  } catch (error) {
    console.error('Error fetching stats:', error)
  } finally {
    loading.value = false
  }
}

const refreshData = async () => {
  await fetchStats()
}

// Lifecycle
onMounted(() => {
  fetchStats()
})

// Watch for refresh trigger
watch(() => props.refreshTrigger, () => {
  refreshData()
})

// Expose methods to parent component
defineExpose({
  refreshData
})
</script>

<style scoped>
/* Custom styles for enhanced visual appeal */
.font-roboto {
  font-family: 'Roboto', system-ui, -apple-system, sans-serif;
  font-weight: 700;
}

.font-sans {
  font-family: system-ui, -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .grid-cols-1 {
    gap: 1rem;
  }
  
  .text-2xl {
    font-size: 1.5rem;
  }
  
  .p-6 {
    padding: 1rem;
  }
}

/* Enhanced hover effects */
@media (hover: hover) {
  .hover\:shadow-xl:hover {
    box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 0 10px 10px -5px rgba(0, 0, 0, 0.04);
  }
}

/* Loading state styles */
.loading {
  opacity: 0.6;
  pointer-events: none;
}

/* Animation for the pulse effect */
@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.5;
  }
}

.animate-pulse {
  animation: pulse 2s cubic-bezier(0.4, 0, 0.6, 1) infinite;
}
</style>