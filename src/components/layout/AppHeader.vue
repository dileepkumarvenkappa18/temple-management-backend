<template>
  <header class="bg-white shadow-lg border-b border-gray-200 sticky top-0 z-50">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center h-16">
        <!-- Logo and Brand -->
        <div class="flex items-center space-x-4">
          <div class="flex-shrink-0 flex items-center">
            <!-- Temple Logo/Icon -->
            <router-link :to="getHomeLink()" class="flex items-center">
              <div class="w-10 h-10 bg-gradient-to-br from-indigo-500 to-indigo-700 rounded-xl flex items-center justify-center shadow-lg">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2L2 7v10c0 5.55 3.84 9.74 9 11 5.16-1.26 9-5.45 9-11V7l-10-5z"/>
                </svg>
              </div>
              <div class="ml-3 hidden sm:block">
                <h1 class="text-xl font-bold text-gray-900 hover:text-indigo-600 transition-colors">{{ appTitle }}</h1>
                <p v-if="currentTemple" class="text-sm text-gray-500 truncate max-w-48">{{ currentTemple }}</p>
              </div>
            </router-link>
          </div>
        </div>

        <!-- Center - Breadcrumb (Desktop only) -->
        <!-- <div class="hidden lg:block flex-1 max-w-md mx-8">
          <AppBreadcrumb />
        </div> -->

        <!-- Right Side - User Actions -->
        <div class="flex items-center space-x-4">
          
          <!-- <button 
            @click="toggleNotifications"
            class="relative p-2 text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 notifications-bell"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                d="M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11c0-3.07-1.64-5.64-4.5-6.32V4a1.5 1.5 0 00-3 0v.68C7.64 5.36 6 7.929 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9"/>
            </svg>

            
            <span 
              v-if="notificationCount > 0" 
              class="absolute -top-1 -right-1 bg-red-500 text-white text-xs rounded-full h-5 w-5 flex items-center justify-center font-medium shadow-lg"
            >
              {{ notificationCount > 9 ? '9+' : notificationCount }}
            </span>
          </button> -->

          <!-- Mobile Menu Toggle -->
          <button
            @click="toggleMobileSidebar"
            class="lg:hidden p-2 text-gray-500 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
            </svg>
          </button>

          <!-- Profile Dropdown -->
          <AppProfileDropdown 
            :user="currentUser"
            :role="currentRole"
            @logout="handleLogout"
          />
        </div>
      </div>
    </div>

    <!-- Mobile Breadcrumb -->
    <div class="lg:hidden border-t border-gray-200 px-4 py-2 bg-gray-50">
      <AppBreadcrumb />
    </div>

    <!-- Notifications Dropdown -->
    <div 
      v-if="showNotifications" 
      class="notifications-container absolute top-16 right-4 w-80 bg-white rounded-xl shadow-2xl border border-gray-200 z-50 max-h-96 overflow-hidden"
    >
      <div class="p-4 border-b border-gray-200">
        <div class="flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900">Notifications</h3>
          <button 
            @click="markAllAsRead"
            class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
          >
            Mark all read
          </button>
        </div>
      </div>
      
      <div class="max-h-64 overflow-y-auto">
        <div v-if="notifications.length === 0" class="p-6 text-center">
          <svg class="w-12 h-12 text-gray-300 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2-2v-5m16 0h-2M4 13h2m13-8V4a1 1 0 00-1-1H7a1 1 0 00-1 1v1m8 0V4.5"/>
          </svg>
          <p class="text-gray-500">No new notifications</p>
        </div>
        
        <div 
          v-for="notification in notifications" 
          :key="notification.id"
          class="p-4 border-b border-gray-100 hover:bg-gray-50 cursor-pointer transition-colors duration-150"
          @click="markAsRead(notification.id)"
        >
          <div class="flex items-start space-x-3">
            <div 
              class="flex-shrink-0 w-2 h-2 mt-2 rounded-full"
              :class="notification.read ? 'bg-gray-300' : 'bg-indigo-500'"
            ></div>
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900">{{ notification.title }}</p>
              <p class="text-sm text-gray-600 mt-1">{{ notification.message }}</p>
              <p class="text-xs text-gray-400 mt-2">{{ formatTime(notification.timestamp) }}</p>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Overlay for notifications -->
    <div 
      v-if="showNotifications" 
      @click="showNotifications = false"
      class="fixed inset-0 z-40"
    ></div>
  </header>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter, useRoute } from 'vue-router'
import AppBreadcrumb from './AppBreadcrumb.vue'
import AppProfileDropdown from './AppProfileDropdown.vue'

// Props
const props = defineProps({
  showSidebar: {
    type: Boolean,
    default: true
  }
})

// Emits
const emit = defineEmits(['toggle-sidebar'])

// Composables
const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()

// Reactive data
const showNotifications = ref(false)
const notifications = ref([
  {
    id: 1,
    title: 'New Seva Booking',
    message: 'Rajesh Kumar booked Abhishekam for tomorrow',
    timestamp: new Date(Date.now() - 30 * 60 * 1000), // 30 minutes ago
    read: false,
    type: 'seva'
  },
  {
    id: 2,
    title: 'Donation Received',
    message: 'Priya Sharma donated ₹5,000 for Annadanam',
    timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000), // 2 hours ago
    read: false,
    type: 'donation'
  },
  {
    id: 3,
    title: 'Event RSVP',
    message: '15 devotees confirmed for Diwali celebration',
    timestamp: new Date(Date.now() - 24 * 60 * 60 * 1000), // 1 day ago
    read: true,
    type: 'event'
  }
])

// Computed properties
const currentUser = computed(() => authStore.user)
const currentRole = computed(() => authStore.user?.role)
const notificationCount = computed(() => notifications.value.filter(n => !n.read).length)
const appTitle = computed(() => {
  switch (currentRole.value) {
    case 'tenant':
      return 'Temple Admin'
    case 'entity_admin':
      return 'Devotee Connect'
    case 'devotee':
      return 'Devotee Dashboard'
    case 'volunteer':
      return 'Volunteer Portal'
    default:
      return 'Devotee Connect'
  }
})

const currentTemple = computed(() => {
  // This would come from a temple store in real implementation
  if (currentRole.value === 'entity_admin' || currentRole.value === 'devotee' || currentRole.value === 'volunteer') {
    return 'Sri Venkateshwara Temple, Bangalore'
  }
  return null
})

// Methods
const getHomeLink = () => {
  const entityId = route.params.entityId
  
  switch (currentRole.value) {
    case 'tenant':
      return '/tenant/dashboard'
    case 'entity_admin':
      return `/entity/${entityId}/dashboard`
    case 'devotee':
      return `/entity/${entityId}/devotee/dashboard`
    case 'volunteer':
      return `/entity/${entityId}/volunteer/dashboard`
    default:
      return '/'
  }
}

const toggleNotifications = (event) => {
  // Stop propagation to prevent the click outside handler from closing immediately
  event.stopPropagation()
  showNotifications.value = !showNotifications.value
}

const toggleMobileSidebar = () => {
  emit('toggle-sidebar')
}

const markAsRead = (notificationId) => {
  const notification = notifications.value.find(n => n.id === notificationId)
  if (notification) {
    notification.read = true
  }
}

const markAllAsRead = () => {
  notifications.value.forEach(notification => {
    notification.read = true
  })
  showNotifications.value = false  // Close the notifications panel
}

const formatTime = (timestamp) => {
  const now = new Date()
  const diff = now - timestamp
  
  if (diff < 60 * 1000) {
    return 'Just now'
  } else if (diff < 60 * 60 * 1000) {
    const minutes = Math.floor(diff / (60 * 1000))
    return `${minutes}m ago`
  } else if (diff < 24 * 60 * 60 * 1000) {
    const hours = Math.floor(diff / (60 * 60 * 1000))
    return `${hours}h ago`
  } else {
    const days = Math.floor(diff / (24 * 60 * 60 * 1000))
    return `${days}d ago`
  }
}

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}

// Click outside handler
const handleClickOutside = (event) => {
  // Check if the click is outside both the notifications container and the bell button
  const isClickOnBell = event.target.closest('.notifications-bell')
  const isClickInContainer = event.target.closest('.notifications-container')
  
  if (showNotifications.value && !isClickInContainer && !isClickOnBell) {
    showNotifications.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>