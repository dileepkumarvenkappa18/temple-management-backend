<template>
  <div class="min-h-screen bg-gray-50 flex flex-col font-sans">
    <!-- Header -->
    <AppHeader
      :user="user"
      :sidebar-open="sidebarOpen"
      @toggle-sidebar="toggleSidebar"
      @logout="handleLogout"
      class="z-30 transition-all duration-300 ease-in-out"
    />
    
    <!-- Main Layout with Fixed Header Height -->
    <div class="flex flex-1 pt-0">
      <!-- Sidebar -->
      <AppSidebar
        :isOpen="sidebarOpen"
        :user="user"
        class="transition-all duration-300 ease-in-out shadow-lg"
      />
      
      <!-- Main Content -->
      <div 
        :class="[
          'flex-1 flex flex-col min-w-0 transition-all duration-300 ease-in-out',
          sidebarOpen ? 'ml-72 md:ml-72 sm:ml-0' : 'ml-0'
        ]"
      >
        <!-- Page Content -->
        <main class="flex-1 p-4 lg:p-6 xl:p-8 pt-4 relative">
          <!-- Content Container with enhanced styling -->
          <div class="transition-all duration-200 ease-in-out">
            <router-view v-slot="{ Component }">
              <transition 
                name="fade" 
                mode="out-in" 
                appear
              >
                <component :is="Component" />
              </transition>
            </router-view>
          </div>
        </main>
      </div>
    </div>
    
    <!-- Global Modals & Toasts -->
    <BaseToast />
    <BaseModal />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, provide } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import BaseToast from '@/components/common/BaseToast.vue'
import BaseModal from '@/components/common/BaseModal.vue'

const router = useRouter()
const authStore = useAuthStore()

// Sidebar state - responsive by default
const sidebarOpen = ref(window.innerWidth > 768)

// Make sidebar state available to child components
provide('sidebarOpen', sidebarOpen)

// Toggle sidebar with animation
const toggleSidebar = () => {
  sidebarOpen.value = !sidebarOpen.value
}

// User data from auth store
const user = computed(() => authStore.user)

// Logout handler with improved error handling
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}

// Handle responsive behavior
onMounted(() => {
  const handleResize = () => {
    if (window.innerWidth < 768) {
      sidebarOpen.value = false
    } else {
      sidebarOpen.value = true
    }
  }
  
  window.addEventListener('resize', handleResize)
  
  // Clean up event listener
  return () => {
    window.removeEventListener('resize', handleResize)
  }
})
</script>

<style>
/* Add Roboto font for headings */
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;500;700&display=swap');

h1, h2, h3, h4, h5, h6 {
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
}

/* Transition for page changes */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease-in-out;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* Indigo theme color variables */
:root {
  --color-primary: #4f46e5; /* indigo-600 */
  --color-secondary: #6366f1; /* indigo-500 */
  --color-accent: #3730a3; /* indigo-800 */
  --color-success-bg: #dcfce7; /* green-100 */
  --color-success-text: #166534; /* green-800 */
  --color-error-bg: #fee2e2; /* red-100 */
  --color-error-text: #991b1b; /* red-800 */
  --color-warning-bg: #fef3c7; /* amber-100 */
  --color-warning-text: #92400e; /* amber-800 */
  --color-info-bg: #dbeafe; /* blue-100 */
  --color-info-text: #1e40af; /* blue-800 */
}

/* Custom scrollbar for better UX */
::-webkit-scrollbar {
  width: 8px;
  height: 8px;
}

::-webkit-scrollbar-track {
  background: #f1f1f1;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb {
  background: #d1d5db;
  border-radius: 10px;
}

::-webkit-scrollbar-thumb:hover {
  background: #6366f1;
}

/* Ensure proper box sizing everywhere */
*, *::before, *::after {
  box-sizing: border-box;
}

/* Responsive adjustments */
@media (max-width: 768px) {
  .ml-72 {
    margin-left: 0;
  }
}
</style>