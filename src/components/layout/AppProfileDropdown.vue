<template>
  <div class="relative">
    <!-- Profile Button -->
    <button
      @click="toggleDropdown"
      class="flex items-center space-x-3 p-2 text-sm rounded-xl hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
      :class="{ 'bg-gray-50': isOpen }"
    >
      <!-- User Avatar -->
      <div class="relative">
        <img
          v-if="user?.avatar"
          :src="user.avatar"
          :alt="user.name"
          class="w-8 h-8 rounded-full object-cover ring-2 ring-white shadow-sm"
        >
        <div
          v-else
          class="w-8 h-8 rounded-full bg-gradient-to-br from-indigo-500 to-indigo-700 flex items-center justify-center shadow-sm"
        >
          <span class="text-white text-sm font-semibold">
            {{ getInitials(user?.name) }}
          </span>
        </div>
        
        <!-- Online Status Indicator -->
        <div class="absolute -bottom-0.5 -right-0.5 w-3 h-3 bg-green-400 rounded-full ring-2 ring-white"></div>
      </div>

      <!-- User Info (Hidden on mobile) -->
      <div class="hidden sm:block text-left">
        <p class="text-gray-900 font-medium">{{ user?.name || 'User' }}</p>
        <p class="text-gray-500 text-xs capitalize">{{ getRoleDisplayName(role) }}</p>
      </div>

      <!-- Dropdown Arrow -->
      <svg 
        class="w-4 h-4 text-gray-400 transition-transform duration-200"
        :class="{ 'rotate-180': isOpen }"
        fill="none" 
        stroke="currentColor" 
        viewBox="0 0 24 24"
      >
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"/>
      </svg>
    </button>

    <!-- Dropdown Menu -->
    <div
      v-if="isOpen"
      class="absolute right-0 mt-2 w-72 bg-white rounded-xl shadow-2xl border border-gray-200 z-50 overflow-hidden"
    >
      <!-- User Info Header -->
      <div class="p-4 bg-gradient-to-r from-indigo-50 to-purple-50 border-b border-gray-200">
        <div class="flex items-center space-x-3">
          <!-- Larger Avatar -->
          <div class="relative">
            <img
              v-if="user?.avatar"
              :src="user.avatar"
              :alt="user.name"
              class="w-12 h-12 rounded-full object-cover ring-3 ring-white shadow-lg"
            >
            <div
              v-else
              class="w-12 h-12 rounded-full bg-gradient-to-br from-indigo-500 to-indigo-700 flex items-center justify-center shadow-lg"
            >
              <span class="text-white text-lg font-semibold">
                {{ getInitials(user?.name) }}
              </span>
            </div>
            <div class="absolute -bottom-1 -right-1 w-4 h-4 bg-green-400 rounded-full ring-2 ring-white"></div>
          </div>
          
          <div class="flex-1 min-w-0">
            <h3 class="text-lg font-semibold text-gray-900 truncate">{{ user?.name || 'User' }}</h3>
            <p class="text-sm text-gray-600">{{ user?.email || 'user@example.com' }}</p>
            <div class="flex items-center mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800 capitalize">
                {{ getRoleDisplayName(role) }}
              </span>
            </div>
          </div>
        </div>
      </div>

      <!-- Menu Items -->
      <div class="py-2">
        <!-- Role-specific Profile Link -->
        <router-link
          :to="getProfileLink()"
          @click="closeDropdown"
          class="flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600 transition-colors duration-200"
        >
          <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
          My Profile
        </router-link>

        <!-- Account Settings -->
        <router-link
          :to="getSettingsLink()"
          @click="closeDropdown"
          class="flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600 transition-colors duration-200"
        >
          <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
          </svg>
          Account Settings
        </router-link>

        <!-- Notifications -->
        <button
          @click="toggleNotifications"
          class="w-full flex items-center justify-between px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600 transition-colors duration-200"
        >
          <div class="flex items-center">
            <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2z"/>
            </svg>
            Notifications
          </div>
          <div class="flex items-center space-x-2">
            <span class="text-xs text-gray-500">3 new</span>
            <div class="relative">
              <input
                type="checkbox"
                :checked="notificationsEnabled"
                @change="toggleNotifications"
                class="sr-only"
              >
              <div class="w-9 h-5 bg-gray-200 rounded-full shadow-inner transition-colors duration-200"
                   :class="notificationsEnabled ? 'bg-indigo-600' : 'bg-gray-200'">
                <div class="w-4 h-4 bg-white rounded-full shadow transform transition-transform duration-200"
                     :class="notificationsEnabled ? 'translate-x-4' : 'translate-x-0.5'"
                     style="margin-top: 2px;">
                </div>
              </div>
            </div>
          </div>
        </button>

        <div class="border-t border-gray-200 my-2"></div>

        <!-- Help & Support -->
        <router-link
          to="/support"
          @click="closeDropdown"
          class="flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600 transition-colors duration-200"
        >
          <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          Help & Support
        </router-link>

        <!-- Keyboard Shortcuts -->
        <!-- <button
          @click="showKeyboardShortcuts = true; closeDropdown()"
          class="w-full flex items-center px-4 py-3 text-sm text-gray-700 hover:bg-gray-50 hover:text-indigo-600 transition-colors duration-200"
        >
          <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          Keyboard Shortcuts
          <span class="ml-auto text-xs text-gray-400">⌘K</span>
        </button> -->

        <div class="border-t border-gray-200 my-2"></div>

        <!-- Logout -->
        <button
          @click="handleLogout"
          class="w-full flex items-center px-4 py-3 text-sm text-red-600 hover:bg-red-50 hover:text-red-700 transition-colors duration-200"
        >
          <svg class="mr-3 h-5 w-5 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1"/>
          </svg>
          Sign Out
        </button>
      </div>
    </div>

    <!-- Overlay -->
    <div 
      v-if="isOpen" 
      @click="closeDropdown"
      class="fixed inset-0 z-40"
    ></div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'

// Props
const props = defineProps({
  user: {
    type: Object,
    default: () => ({})
  },
  role: {
    type: String,
    default: 'user'
  }
})

// Emits
const emit = defineEmits(['logout'])

// Composables
const route = useRoute()

// Reactive data
const isOpen = ref(false)
const notificationsEnabled = ref(true)
const showKeyboardShortcuts = ref(false)

// Methods
const toggleDropdown = () => {
  isOpen.value = !isOpen.value
}

const closeDropdown = () => {
  isOpen.value = false
}

const getInitials = (name) => {
  if (!name) return 'U'
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const getRoleDisplayName = (role) => {
  const roleNames = {
    'tenant': 'Temple Admin',
    'entity_admin': 'Temple Manager',
    'devotee': 'Devotee',
    'volunteer': 'Volunteer',
    'super_admin': 'Super Admin'
  }
  return roleNames[role] || 'User'
}

const getProfileLink = () => {
  const entityId = route.params.entityId
  switch (props.role) {
    case 'tenant':
      return '/tenant/profile'
    case 'entity_admin':
      return `/entity/${entityId}/profile`
    case 'devotee':
      return `/entity/${entityId}/devotee/profile`
    case 'volunteer':
      return `/entity/${entityId}/volunteer/profile`
    default:
      return '/profile'
  }
}

const getSettingsLink = () => {
  const entityId = route.params.entityId
  switch (props.role) {
    case 'tenant':
      return '/tenant/settings'
    case 'entity_admin':
      return `/entity/${entityId}/settings`
    case 'devotee':
      return `/entity/${entityId}/devotee/settings`
    case 'volunteer':
      return `/entity/${entityId}/volunteer/settings`
    default:
      return '/settings'
  }
}

const toggleNotifications = () => {
  notificationsEnabled.value = !notificationsEnabled.value
}

const handleLogout = () => {
  closeDropdown()
  emit('logout')
}

// Click outside handler
const handleClickOutside = (event) => {
  if (isOpen.value && !event.target.closest('.relative')) {
    closeDropdown()
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>