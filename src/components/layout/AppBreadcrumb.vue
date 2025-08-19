<template>
  <!-- <nav class="flex items-center space-x-2 text-sm" aria-label="Breadcrumb">
    <ol class="flex items-center space-x-2">
      <li v-for="(item, index) in finalBreadcrumbItems" :key="index" class="flex items-center">
        Home/Dashboard Icon for first item
        <div v-if="index === 0" class="flex items-center">
          <svg class="w-4 h-4 text-gray-400 mr-2" fill="currentColor" viewBox="0 0 24 24">
            <path d="M10 20v-6h4v6h5v-8h3L12 3 2 12h3v8z"/>
          </svg>
        </div>
        
        Separator
        <svg 
          v-if="index > 0" 
          class="w-4 h-4 text-gray-300 mx-2" 
          fill="currentColor" 
          viewBox="0 0 24 24"
        >
          <path d="M8.59 16.59L13.17 12 8.59 7.41 10 6l6 6-6 6-1.41-1.41z"/>
        </svg>

        Breadcrumb Link
        <router-link
          v-if="!item.current && item.to"
          :to="item.to"
          class="text-gray-500 hover:text-indigo-600 transition-colors duration-200 font-medium"
        >
          {{ item.name }}
        </router-link>

        Current Page (Non-clickable)
        <span
          v-else
          class="font-semibold"
          :class="item.current ? 'text-indigo-700' : 'text-gray-500'"
        >
          {{ item.name }}
        </span>

        Badge/Counter
        <span 
          v-if="item.badge" 
          class="ml-2 inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gray-100 text-gray-800"
        >
          {{ item.badge }}
        </span>
      </li>
    </ol>
  </nav> -->
</template>

<script setup>
import { computed, defineProps } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

// Props
const props = defineProps({
  items: {
    type: Array,
    default: () => []
  }
})

// Composables
const route = useRoute()
const authStore = useAuthStore()

// Computed properties
const currentRole = computed(() => authStore.user?.role)

const autoBreadcrumbItems = computed(() => {
  const items = []
  const path = route.path
  const params = route.params
  
  // Helper function to get role-specific dashboard
  const getDashboardItem = () => {
    switch (currentRole.value) {
      case 'tenant':
        return { name: 'Dashboard', to: '/tenant/dashboard' }
      case 'entity_admin':
        return { name: 'Dashboard', to: `/entity/${params.entityId}/dashboard` }
      case 'devotee':
        return { name: 'Dashboard', to: `/entity/${params.entityId}/devotee/dashboard` }
      case 'volunteer':
        return { name: 'Dashboard', to: `/entity/${params.entityId}/volunteer/dashboard` }
      default:
        return { name: 'Dashboard', to: '/' }
    }
  }

  // Build breadcrumb based on current route
  if (path === '/' || path === '/landing') {
    items.push({ name: 'Home', current: true })
  } else if (path.includes('/login') || path.includes('/register')) {
    items.push({ name: 'Home', to: '/' })
    items.push({ name: path.includes('/login') ? 'Login' : 'Register', current: true })
  } else if (path.includes('/tenant')) {
    items.push(getDashboardItem())
    
    if (path.includes('/entities/create')) {
      items.push({ name: 'My Temples', to: '/tenant/entities' })
      items.push({ name: 'Create Temple', current: true })
    } else if (path.includes('/entities') && !path.includes('/create') && !path.includes('/edit')) {
      // Check if we're on a temple details page
      if (path.match(/\/tenant\/entities\/\d+$/)) {
        items.push({ name: 'My Temples', to: '/tenant/entities' })
        items.push({ name: 'Temple Details', current: true })
      } else {
        items.push({ name: 'My Temples', current: true, badge: '3' })
      }
    } else if (path.includes('/entities') && path.includes('/edit')) {
      items.push({ name: 'My Temples', to: '/tenant/entities' })
      items.push({ name: 'Edit Temple', current: true })
    } else if (path.includes('/settings')) {
      items.push({ name: 'Settings', current: true })
    } else {
      items[0].current = true
    }
  } else if (path.includes('/entity') && currentRole.value === 'entity_admin') {
    items.push(getDashboardItem())
    
    if (path.includes('/devotees')) {
      items.push({ name: 'Management', to: `/entity/${params.entityId}/devotees` })
      items.push({ name: 'Devotees', current: true, badge: '1,247' })
    } else if (path.includes('/volunteers')) {
      items.push({ name: 'Management', to: `/entity/${params.entityId}/volunteers` })
      items.push({ name: 'Volunteers', current: true, badge: '89' })
    } else if (path.includes('/sevas')) {
      items.push({ name: 'Seva Management', current: true, badge: '5 pending' })
    } else if (path.includes('/donations')) {
      items.push({ name: 'Donations', current: true })
    } else if (path.includes('/events')) {
      items.push({ name: 'Events & Festivals', current: true, badge: '2 upcoming' })
    } else if (path.includes('/communication')) {
      items.push({ name: 'Communication', current: true })
    } else if (path.includes('/settings')) {
      items.push({ name: 'Settings', current: true })
    } else {
      items[0].current = true
    }
  } else if (path.includes('/entity') && currentRole.value === 'devotee') {
    items.push(getDashboardItem())
    
    if (path.includes('/seva-booking')) {
      items.push({ name: 'Book Seva', current: true })
    } else if (path.includes('/donations')) {
      items.push({ name: 'My Donations', current: true })
    } else if (path.includes('/events')) {
      items.push({ name: 'Temple Events', current: true })
    } else if (path.includes('/profile')) {
      items.push({ name: 'My Profile', current: true })
      
      if (path.includes('/create')) {
        items[items.length - 1].current = false
        items[items.length - 1].to = `/entity/${params.entityId}/devotee/profile`
        items.push({ name: 'Complete Profile', current: true })
      }
    } else {
      items[0].current = true
    }
  } else if (path.includes('/entity') && currentRole.value === 'volunteer') {
    items.push(getDashboardItem())
    
    if (path.includes('/assignments')) {
      items.push({ name: 'My Assignments', current: true, badge: '3 active' })
    } else if (path.includes('/schedule')) {
      items.push({ name: 'My Schedule', current: true })
    } else if (path.includes('/events')) {
      items.push({ name: 'Temple Events', current: true })
    } else {
      items[0].current = true
    }
  } else if (path.includes('/devotee/temple-selection')) {
    items.push({ name: 'Home', to: '/' })
    items.push({ name: 'Select Temple', current: true })
  } else if (path.includes('/volunteer/temple-selection')) {
    items.push({ name: 'Home', to: '/' })
    items.push({ name: 'Select Temple', current: true })
  } else {
    // Default fallback
    items.push({ name: 'Home', to: '/' })
    items.push({ name: 'Page', current: true })
  }

  return items
})

// Use passed items if available, otherwise generate automatically
const finalBreadcrumbItems = computed(() => {
  if (props.items && props.items.length > 0) {
    // Convert from the format used in TempleDetails to the format used in AppBreadcrumb
    return props.items.map((item, index) => ({
      name: item.name,
      to: item.path || undefined, // Use path as 'to' if present
      current: index === props.items.length - 1 // Last item is current
    }))
  }
  return autoBreadcrumbItems.value
})
</script>