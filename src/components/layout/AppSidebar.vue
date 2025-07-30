<template>
  <!-- Root layout container -->
  <div class="flex min-h-screen">
    <!-- Fixed Sidebar -->
    <div class="fixed inset-y-0 left-0 z-20 hidden lg:block">
      <aside class="h-full w-70 bg-white border-r border-gray-200 overflow-y-auto">
        <!-- Navigation Menu -->
        <nav class="px-12 py-6">
          <div class="space-y-2">
            <!-- Role-based Navigation Items -->
            <template v-for="item in navigationItems" :key="item.name">
              <!-- Single Menu Item -->
              <router-link
                v-if="!item.children"
                :to="item.to"
                class="group flex items-center px-3 py-2.5 text-sm font-medium rounded-xl transition-all duration-200"
                :class="isActiveRoute(item.to) 
                  ? 'bg-indigo-50 text-indigo-700 border-r-2 border-indigo-600' 
                  : 'text-gray-700 hover:bg-gray-50 hover:text-indigo-600'"
                @click="$emit('navigate')"
              >
                <component 
                  :is="item.icon" 
                  class="mr-3 h-5 w-5 flex-shrink-0"
                  :class="isActiveRoute(item.to) ? 'text-indigo-600' : 'text-gray-400 group-hover:text-indigo-500'"
                />
                {{ item.name }}
                
                <!-- Badge for notifications/counts -->
                <span 
                  v-if="item.badge" 
                  class="ml-auto inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
                  :class="item.badgeClass || 'bg-indigo-100 text-indigo-800'"
                >
                  {{ item.badge }}
                </span>
              </router-link>

              <!-- Menu Group with Children -->
              <div v-else class="space-y-1">
                <button
                  @click="toggleGroup(item.name)"
                  class="group w-full flex items-center justify-between px-3 py-2.5 text-sm font-medium text-gray-700 hover:bg-gray-50 hover:text-indigo-600 rounded-xl transition-all duration-200"
                >
                  <div class="flex items-center">
                    <component 
                      :is="item.icon" 
                      class="mr-3 h-5 w-5 flex-shrink-0 text-gray-400 group-hover:text-indigo-500"
                    />
                    {{ item.name }}
                  </div>
                  <svg 
                    class="h-4 w-4 text-gray-400 transition-transform duration-200"
                    :class="{ 'rotate-90': openGroups.includes(item.name) }"
                    fill="none" 
                    stroke="currentColor" 
                    viewBox="0 0 24 24"
                  >
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                  </svg>
                </button>
                
                <!-- Submenu Items -->
                <div 
                  v-show="openGroups.includes(item.name)"
                  class="ml-8 space-y-1 border-l-2 border-gray-100 pl-4"
                >
                  <router-link
                    v-for="child in item.children"
                    :key="child.name"
                    :to="child.to"
                    class="block px-3 py-2 text-sm text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-all duration-200"
                    :class="isActiveRoute(child.to) ? 'text-indigo-700 bg-indigo-50 font-medium' : ''"
                    @click="$emit('navigate')"
                  >
                    {{ child.name }}
                  </router-link>
                </div>
              </div>
            </template>
          </div>

          <!-- Bottom Section -->
          <div class="mt-8 pt-6 border-t border-gray-200">
            <!-- Quick Stats (for Entity Admin) -->
            <div v-if="currentRole === 'entity_admin'" class="mb-6">
              <h4 class="text-xs font-semibold text-gray-400 uppercase tracking-wide mb-3">Quick Stats</h4>
              <div class="space-y-2">
                <div class="flex items-center justify-between p-2 bg-green-50 rounded-lg">
                  <span class="text-sm text-green-700">Devotees</span>
                  <span class="text-sm font-semibold text-green-800">{{ stats.devotees }}</span>
                </div>
                <div class="flex items-center justify-between p-2 bg-blue-50 rounded-lg">
                  <span class="text-sm text-blue-700">Today's Sevas</span>
                  <span class="text-sm font-semibold text-blue-800">{{ stats.sevas }}</span>
                </div>
                <div class="flex items-center justify-between p-2 bg-purple-50 rounded-lg">
                  <span class="text-sm text-purple-700">This Month</span>
                  <span class="text-sm font-semibold text-purple-800">₹{{ stats.donations }}</span>
                </div>
              </div>
            </div>

            <!-- Help & Support -->
            <router-link
              to="/support"
              class="group flex items-center px-3 py-2.5 text-sm font-medium text-gray-600 hover:text-indigo-600 hover:bg-indigo-50 rounded-xl transition-all duration-200"
            >
              <svg class="mr-3 h-5 w-5 text-gray-400 group-hover:text-indigo-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              Help & Support
            </router-link>
          </div>
        </nav>
      </aside>
    </div>

    <!-- Mobile Sidebar (when open) -->
    <div v-if="isOpen" class="lg:hidden">
      <div 
        class="fixed inset-0 bg-gray-600 bg-opacity-50 z-40"
        @click="$emit('close')"
      ></div>
      <div class="fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg">
        <div class="flex items-center justify-between h-16 px-6 border-b border-gray-200">
          <button
            @click="$emit('close')"
            class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors duration-200"
          >
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>

        <!-- Mobile Navigation (duplicate of desktop nav) -->
        <nav class="px-4 py-6 overflow-y-auto">
          <!-- Same navigation content as desktop sidebar -->
        </nav>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="flex-1 flex flex-col min-w-0 lg:pl-64">
      <!-- Main content goes here -->
      <slot></slot>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Props
const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  },
  user: {
    type: Object,
    default: null
  },
  navigation: {
    type: Array,
    default: () => []
  }
})

// Emits
defineEmits(['close', 'navigate'])

// Composables
const route = useRoute()
const authStore = useAuthStore()

// Reactive data
const openGroups = ref(['Management']) // Default open groups
const stats = ref({
  devotees: 1247,
  sevas: 8,
  donations: '45,000'
})

// Computed properties
const currentRole = computed(() => authStore.user?.role)

// Define SVG components as strings for dynamic rendering
const DashboardIcon = 'svg'
const UsersIcon = 'svg'
const HeartIcon = 'svg'
const CurrencyIcon = 'svg'
const CalendarIcon = 'svg'
const ChatIcon = 'svg'
const UserGroupIcon = 'svg'
const CogIcon = 'svg'
const TemplateIcon = 'svg'
const BuildingIcon = 'svg'

// Navigation items based on role
const navigationItems = computed(() => {
  const baseItems = []
  
  switch (currentRole.value) {
    case 'tenant':
      return [
        {
          name: 'Dashboard',
          to: '/tenant/dashboard',
          icon: DashboardIcon
        },
        {
          name: 'My Temples',
          to: '/tenant/entities',
          icon: BuildingIcon,
          badge: '3'
        },
        {
          name: 'Create Temple',
          to: '/tenant/entities/create',
          icon: TemplateIcon
        },
        {
          name: 'Settings',
          to: '/tenant/settings',
          icon: CogIcon
        }
      ]

    case 'entity_admin':
      return [
        {
          name: 'Dashboard',
          to: `/entity/${route.params.entityId || 1}/dashboard`,
          icon: DashboardIcon
        },
        {
          name: 'Management',
          icon: UsersIcon,
          children: [
            {
              name: 'Devotees',
              to: `/entity/${route.params.entityId || 1}/devotees`
            },
            {
              name: 'Volunteers',
              to: `/entity/${route.params.entityId || 1}/volunteers`
            }
          ]
        },
        {
          name: 'Seva Management',
          to: `/entity/${route.params.entityId || 1}/sevas`,
          icon: HeartIcon,
          badge: '5',
          badgeClass: 'bg-orange-100 text-orange-800'
        },
        {
          name: 'Donations',
          to: `/entity/${route.params.entityId || 1}/donations`,
          icon: CurrencyIcon
        },
        {
          name: 'Events & Festivals',
          to: `/entity/${route.params.entityId || 1}/events`,
          icon: CalendarIcon,
          badge: '2',
          badgeClass: 'bg-purple-100 text-purple-800'
        },
        {
          name: 'Communication',
          to: `/entity/${route.params.entityId || 1}/communication`,
          icon: ChatIcon
        },
        {
          name: 'Settings',
          to: `/entity/${route.params.entityId || 1}/settings`,
          icon: CogIcon
        }
      ]

    case 'devotee':
      return [
        {
          name: 'Dashboard',
          to: `/entity/${route.params.entityId || 1}/devotee/dashboard`,
          icon: DashboardIcon
        },
        {
          name: 'Book Seva',
          to: `/entity/${route.params.entityId || 1}/devotee/seva-booking`,
          icon: HeartIcon
        },
        {
          name: 'My Donations',
          to: `/entity/${route.params.entityId || 1}/devotee/donations`,
          icon: CurrencyIcon
        },
        {
          name: 'Temple Events',
          to: `/entity/${route.params.entityId || 1}/devotee/events`,
          icon: CalendarIcon
        },
        {
          name: 'My Profile',
          to: `/entity/${route.params.entityId || 1}/devotee/profile`,
          icon: UsersIcon
        }
      ]

    case 'volunteer':
      return [
        {
          name: 'Dashboard',
          to: `/entity/${route.params.entityId || 1}/volunteer/dashboard`,
          icon: DashboardIcon
        },
        {
          name: 'My Assignments',
          to: `/entity/${route.params.entityId || 1}/volunteer/assignments`,
          icon: HeartIcon,
          badge: '3',
          badgeClass: 'bg-green-100 text-green-800'
        },
        {
          name: 'My Schedule',
          to: `/entity/${route.params.entityId || 1}/volunteer/schedule`,
          icon: CalendarIcon
        },
        {
          name: 'Temple Events',
          to: `/entity/${route.params.entityId || 1}/volunteer/events`,
          icon: CalendarIcon
        }
      ]

    default:
      return []
  }
})

// Methods
const isActiveRoute = (routePath) => {
  return route.path === routePath || route.path.startsWith(routePath + '/')
}

const toggleGroup = (groupName) => {
  const index = openGroups.value.indexOf(groupName)
  if (index > -1) {
    openGroups.value.splice(index, 1)
  } else {
    openGroups.value.push(groupName)
  }
}
</script>

<style scoped>
/* Custom scrollbar for webkit browsers */
nav::-webkit-scrollbar {
  width: 4px;
}

nav::-webkit-scrollbar-track {
  background: #f1f5f9;
}

nav::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

nav::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>