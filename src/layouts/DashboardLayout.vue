<!-- src/layouts/DashboardLayout.vue -->
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <AppHeader 
      :user="user" 
      :sidebar-open="sidebarOpen"
      @toggle-sidebar="sidebarOpen = !sidebarOpen"
      @logout="handleLogout"
    />

    <!-- Main Layout -->
    <div class="flex">
      <!-- Sidebar -->
      <AppSidebar 
  :isOpen="sidebarOpen"
  :user="user"
  :navigation="navigation"
  @close="sidebarOpen = false"
/>


      <!-- Main Content -->
      <div class="flex-1 flex flex-col min-w-0">
        <!-- Breadcrumb -->
        <!-- <AppBreadcrumb :items="breadcrumbItems" /> -->

        <!-- Page Content -->
        <main class="flex-1 p-4 lg:p-6 xl:p-8">
          <router-view />
        </main>
      </div>
    </div>

    <!-- Global Modals & Toasts -->
    <BaseToast />
    <BaseModal />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import AppBreadcrumb from '@/components/layout/AppBreadcrumb.vue'
import BaseToast from '@/components/common/BaseToast.vue'
import BaseModal from '@/components/common/BaseModal.vue'

const router = useRouter()
const authStore = useAuthStore()

// Sidebar state
const sidebarOpen = ref(false)

// User data from auth store
const user = computed(() => authStore.user)

// Navigation items based on user role
const navigation = computed(() => {
  const baseNavigation = []
  
  if (!user.value) return baseNavigation

  // Different navigation for different roles
  switch (user.value.role) {
    case 'TENANT':
      return [
        {
          name: 'Dashboard',
          icon: 'home',
          route: '/tenant/dashboard',
          current: router.currentRoute.value.name === 'TenantDashboard'
        },
        {
          name: 'My Temples',
          icon: 'building',
          route: '/tenant/temples',
          current: router.currentRoute.value.path.includes('/tenant/temples')
        },
        {
          name: 'Create Temple',
          icon: 'plus',
          route: '/tenant/temples/create',
          current: router.currentRoute.value.name === 'CreateTemple'
        },
        {
          name: 'Profile',
          icon: 'user',
          route: '/tenant/profile',
          current: router.currentRoute.value.name === 'TenantProfile'
        }
      ]

    case 'DEVOTEE':
      return [
        {
          name: 'Dashboard',
          icon: 'home',
          route: `/entity/${user.value.entityId}/devotee/dashboard`,
          current: router.currentRoute.value.name === 'DevoteeDashboard'
        },
        {
          name: 'My Profile',
          icon: 'user',
          route: `/entity/${user.value.entityId}/devotee/profile`,
          current: router.currentRoute.value.path.includes('/profile')
        },
        {
          name: 'Book Seva',
          icon: 'calendar',
          route: `/entity/${user.value.entityId}/devotee/seva-booking`,
          current: router.currentRoute.value.path.includes('/seva-booking')
        },
        {
          name: 'My Donations',
          icon: 'gift',
          route: `/entity/${user.value.entityId}/devotee/donations`,
          current: router.currentRoute.value.path.includes('/donations')
        },
        {
          name: 'Events',
          icon: 'star',
          route: `/entity/${user.value.entityId}/devotee/events`,
          current: router.currentRoute.value.path.includes('/events')
        }
      ]

    case 'VOLUNTEER':
      return [
        {
          name: 'Dashboard',
          icon: 'home',
          route: `/entity/${user.value.entityId}/volunteer/dashboard`,
          current: router.currentRoute.value.name === 'VolunteerDashboard'
        },
        {
          name: 'My Assignments',
          icon: 'clipboard',
          route: `/entity/${user.value.entityId}/volunteer/assignments`,
          current: router.currentRoute.value.path.includes('/assignments')
        },
        {
          name: 'Schedule',
          icon: 'calendar',
          route: `/entity/${user.value.entityId}/volunteer/schedule`,
          current: router.currentRoute.value.path.includes('/schedule')
        },
        {
          name: 'Profile',
          icon: 'user',
          route: `/entity/${user.value.entityId}/volunteer/profile`,
          current: router.currentRoute.value.path.includes('/profile')
        }
      ]

    case 'ENTITY_ADMIN':
      return [
        {
          name: 'Dashboard',
          icon: 'home',
          route: `/entity/${user.value.entityId}/dashboard`,
          current: router.currentRoute.value.name === 'EntityDashboard'
        },
        {
          name: 'Devotees',
          icon: 'users',
          route: `/entity/${user.value.entityId}/devotees`,
          current: router.currentRoute.value.path.includes('/devotees')
        },
        {
          name: 'Seva Management',
          icon: 'calendar',
          route: `/entity/${user.value.entityId}/sevas`,
          current: router.currentRoute.value.path.includes('/sevas')
        },
        {
          name: 'Donations',
          icon: 'gift',
          route: `/entity/${user.value.entityId}/donations`,
          current: router.currentRoute.value.path.includes('/donations')
        },
        {
          name: 'Events & Festivals',
          icon: 'star',
          route: `/entity/${user.value.entityId}/events`,
          current: router.currentRoute.value.path.includes('/events')
        },
        {
          name: 'Volunteers',
          icon: 'user-group',
          route: `/entity/${user.value.entityId}/volunteers`,
          current: router.currentRoute.value.path.includes('/volunteers')
        },
        {
          name: 'Communication',
          icon: 'mail',
          route: `/entity/${user.value.entityId}/communication`,
          current: router.currentRoute.value.path.includes('/communication')
        }
      ]

    default:
      return baseNavigation
  }
})

// Breadcrumb items based on current route
const breadcrumbItems = computed(() => {
  const route = router.currentRoute.value
  const items = []

  // Build breadcrumb based on route path
  if (route.path.includes('/tenant')) {
    items.push({ name: 'Tenant Dashboard', href: '/tenant/dashboard' })
    if (route.path.includes('/temples/create')) {
      items.push({ name: 'Create Temple' })
    } else if (route.path.includes('/temples')) {
      items.push({ name: 'My Temples' })
    }
  } else if (route.path.includes('/entity')) {
    const entityId = route.params.id
    items.push({ name: 'Dashboard', href: `/entity/${entityId}/dashboard` })
    
    if (route.path.includes('/devotees')) {
      items.push({ name: 'Devotees' })
    } else if (route.path.includes('/sevas')) {
      items.push({ name: 'Seva Management' })
    } else if (route.path.includes('/donations')) {
      items.push({ name: 'Donations' })
    } else if (route.path.includes('/events')) {
      items.push({ name: 'Events & Festivals' })
    } else if (route.path.includes('/volunteers')) {
      items.push({ name: 'Volunteers' })
    } else if (route.path.includes('/communication')) {
      items.push({ name: 'Communication' })
    }
  }

  return items
})

// Logout handler
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

// Initialize on mount
// onMounted(() => {
//   // Ensure user is authenticated
//   if (!user.value) {
//     router.push('/login')
//   }
// })
</script>