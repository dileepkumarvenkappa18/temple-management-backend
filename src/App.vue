<template>
  <div id="app" class="min-h-screen bg-gray-50">
    <!-- Global Loading Overlay - Now also used for auth initialization -->
    <Transition
      name="fade"
      enter-active-class="transition-opacity duration-300 ease-out"
      leave-active-class="transition-opacity duration-200 ease-in"
      enter-from-class="opacity-0"
      enter-to-class="opacity-100"
      leave-from-class="opacity-100"
      leave-to-class="opacity-0"
    >
      <div
        v-if="isGlobalLoading || isInitializing"
        class="fixed inset-0 z-50 flex items-center justify-center bg-white bg-opacity-90 backdrop-blur-sm"
      >
        <div class="flex flex-col items-center space-y-4">
          <div class="animate-spin rounded-full h-12 w-12 border-4 border-indigo-200 border-t-indigo-600"></div>
          <p class="text-indigo-600 font-medium">{{ isInitializing ? 'Initializing application...' : 'Loading...' }}</p>
        </div>
      </div>
    </Transition>

    <!-- Toast Notifications Container -->
    <Teleport to="body">
      <div class="fixed top-4 right-4 z-50 space-y-2">
        <TransitionGroup
          name="toast"
          enter-active-class="transition-all duration-300 ease-out"
          leave-active-class="transition-all duration-200 ease-in"
          enter-from-class="opacity-0 transform translate-x-full"
          enter-to-class="opacity-100 transform translate-x-0"
          leave-from-class="opacity-100 transform translate-x-0"
          leave-to-class="opacity-0 transform translate-x-full"
        >
          <div
            v-for="toast in toasts"
            :key="toast.id"
            :class="[
              'max-w-sm w-full bg-white shadow-lg rounded-2xl pointer-events-auto overflow-hidden border-l-4',
              getToastBorderColor(toast.type)
            ]"
          >
            <div class="p-4">
              <div class="flex items-start">
                <div class="flex-shrink-0">
                  <component
                    :is="getToastIcon(toast.type)"
                    :class="[
                      'h-6 w-6',
                      getToastIconColor(toast.type)
                    ]"
                  />
                </div>
                <div class="ml-3 w-0 flex-1">
                  <p class="text-sm font-medium text-gray-900">
                    {{ toast.title }}
                  </p>
                  <p v-if="toast.message" class="mt-1 text-sm text-gray-500">
                    {{ toast.message }}
                  </p>
                </div>
                <div class="ml-4 flex-shrink-0 flex">
                  <button
                    @click="removeToast(toast.id)"
                    class="bg-white rounded-md inline-flex text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors duration-200"
                  >
                    <span class="sr-only">Close</span>
                    <XMarkIcon class="h-5 w-5" />
                  </button>
                </div>
              </div>
            </div>
          </div>
        </TransitionGroup>
      </div>
    </Teleport>

    <!-- Main Router View - Only shown after initialization -->
    <template v-if="!isInitializing">
      <!-- THIS IS THE UPDATED SECTION TO HANDLE LAYOUTS -->
      <component :is="currentLayout">
        <router-view v-slot="{ Component, route }">
          <Transition
            :name="getPageTransition(route)"
            mode="out-in"
            enter-active-class="transition-all duration-300 ease-out"
            leave-active-class="transition-all duration-200 ease-in"
          >
            <component :is="Component" :key="route.path" />
          </Transition>
        </router-view>
      </component>
    </template>
    
    <!-- Initialization Error Fallback -->
    <div v-else-if="initializationFailed" class="min-h-screen flex flex-col items-center justify-center bg-gray-50">
      <div class="text-center max-w-md px-4">
        <h1 class="text-2xl font-bold text-indigo-600 mb-4">Something went wrong</h1>
        <p class="text-gray-600 mb-6">There was a problem starting the application. Please try again.</p>
        <div class="space-y-4">
          <button @click="refreshPage" class="w-full py-2 px-4 bg-indigo-600 text-white rounded hover:bg-indigo-700">
            Refresh Page
          </button>
          <button @click="resetAndLogout" class="w-full py-2 px-4 border border-gray-300 bg-white text-gray-700 rounded hover:bg-gray-50">
            Reset & Logout
          </button>
        </div>
      </div>
    </div>

    <!-- Global Modal Container -->
    <Teleport to="body">
      <Transition
        name="modal"
        enter-active-class="transition-opacity duration-300 ease-out"
        leave-active-class="transition-opacity duration-200 ease-in"
        enter-from-class="opacity-0"
        enter-to-class="opacity-100"
        leave-from-class="opacity-100"
        leave-to-class="opacity-0"
      >
        <div
          v-if="isModalOpen"
          class="fixed inset-0 z-40 overflow-y-auto"
          @click.self="closeModal"
        >
          <div class="flex min-h-screen items-center justify-center p-4">
            <div class="fixed inset-0 bg-black bg-opacity-25 backdrop-blur-sm transition-opacity"></div>
            
            <Transition
              name="modal-content"
              enter-active-class="transition-all duration-300 ease-out"
              leave-active-class="transition-all duration-200 ease-in"
              enter-from-class="opacity-0 transform scale-95"
              enter-to-class="opacity-100 transform scale-100"
              leave-from-class="opacity-100 transform scale-100"
              leave-to-class="opacity-0 transform scale-95"
            >
              <div
                v-if="isModalOpen"
                class="relative bg-white rounded-2xl shadow-2xl max-w-lg w-full mx-auto overflow-hidden"
              >
                <slot name="modal-content"></slot>
              </div>
            </Transition>
          </div>
        </div>
      </Transition>
    </Teleport>

    <!-- Offline Indicator -->
    <Transition
      name="slide-up"
      enter-active-class="transition-transform duration-300 ease-out"
      leave-active-class="transition-transform duration-200 ease-in"
      enter-from-class="transform translate-y-full"
      enter-to-class="transform translate-y-0"
      leave-from-class="transform translate-y-0"
      leave-to-class="transform translate-y-full"
    >
      <div
        v-if="!isOnline"
        class="fixed bottom-4 left-4 right-4 md:left-auto md:right-4 md:w-80 bg-red-600 text-white px-4 py-3 rounded-2xl shadow-lg flex items-center space-x-3 z-30"
      >
        <WifiIcon class="h-5 w-5 text-red-200" />
        <div class="flex-1">
          <p class="font-medium text-sm">No Internet Connection</p>
          <p class="text-xs text-red-200">Some features may be limited</p>
        </div>
      </div>
    </Transition>

    <!-- PWA Install Prompt -->
    <Transition
      name="slide-up"
      enter-active-class="transition-transform duration-300 ease-out"
      leave-active-class="transition-transform duration-200 ease-in"
      enter-from-class="transform translate-y-full"
      enter-to-class="transform translate-y-0"
      leave-from-class="transform translate-y-0"
      leave-to-class="transform translate-y-full"
    >
      <div
        v-if="showPWAPrompt"
        class="fixed bottom-4 left-4 right-4 md:left-auto md:right-4 md:w-96 bg-white border border-gray-200 rounded-2xl shadow-2xl p-4 z-30"
      >
        <div class="flex items-start space-x-3">
          <div class="flex-shrink-0">
            <div class="h-10 w-10 bg-indigo-100 rounded-xl flex items-center justify-center">
              <DevicePhoneMobileIcon class="h-6 w-6 text-indigo-600" />
            </div>
          </div>
          <div class="flex-1 min-w-0">
            <h3 class="text-sm font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
              Install Temple Manager
            </h3>
            <p class="text-xs text-gray-500 mt-1">
              Add to your home screen for quick access
            </p>
          </div>
        </div>
        <div class="mt-4 flex space-x-2">
          <button
            @click="installPWA"
            class="flex-1 bg-indigo-600 text-white text-xs font-medium py-2 px-3 rounded-xl hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-colors duration-200"
          >
            Install
          </button>
          <button
            @click="dismissPWAPrompt"
            class="flex-1 bg-gray-100 text-gray-700 text-xs font-medium py-2 px-3 rounded-xl hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors duration-200"
          >
            Not Now
          </button>
        </div>
      </div>
    </Transition>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed, markRaw, shallowRef } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import {
  CheckCircleIcon,
  ExclamationTriangleIcon,
  InformationCircleIcon,
  XCircleIcon,
  XMarkIcon,
  WifiIcon,
  DevicePhoneMobileIcon
} from '@heroicons/vue/24/outline'

// Import layouts
import AuthLayout from '@/layouts/AuthLayout.vue'
import PublicLayout from '@/layouts/PublicLayout.vue'
import DashboardLayout from '@/layouts/DashboardLayout.vue'
import SuperAdminLayout from '@/layouts/SuperAdminLayout.vue'

// Router setup
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

// Global loading state
const isGlobalLoading = ref(false)
const isInitializing = ref(true)
const initializationFailed = ref(false)

// Online/Offline detection
const isOnline = ref(navigator.onLine)

// Toast notifications
const toasts = ref([])
let toastIdCounter = 0

// Modal state
const isModalOpen = ref(false)

// PWA Installation
const showPWAPrompt = ref(false)
let deferredPrompt = null

// Layout management - NEW CODE
const layouts = {
  'AuthLayout': markRaw(AuthLayout),
  'PublicLayout': markRaw(PublicLayout),
  'DashboardLayout': markRaw(DashboardLayout),
  'SuperAdminLayout': markRaw(SuperAdminLayout)
}

const currentLayout = computed(() => {
  const layoutName = route.meta.layout || 'PublicLayout'
  console.log('Current layout:', layoutName, 'for route:', route.path)
  return layouts[layoutName] || layouts['PublicLayout']
})

// Auth-related computed properties
const isAuthenticated = computed(() => authStore.isAuthenticated)
const isPublicRoute = computed(() => !route.meta?.requiresAuth)
const isAuthRoute = computed(() => 
  ['/login', '/register', '/forgot-password', '/reset-password'].includes(route.path) || 
  route.path.startsWith('/auth/')
)

// Toast management
const addToast = (type, title, message = '', duration = 5000) => {
  const id = ++toastIdCounter
  const toast = { id, type, title, message }
  
  toasts.value.push(toast)
  
  if (duration > 0) {
    setTimeout(() => {
      removeToast(id)
    }, duration)
  }
  
  return id
}

const removeToast = (id) => {
  const index = toasts.value.findIndex(toast => toast.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

// Toast styling helpers
const getToastBorderColor = (type) => {
  const colors = {
    success: 'border-green-500',
    error: 'border-red-500',
    warning: 'border-yellow-500',
    info: 'border-indigo-500'
  }
  return colors[type] || colors.info
}

const getToastIconColor = (type) => {
  const colors = {
    success: 'text-green-500',
    error: 'text-red-500',
    warning: 'text-yellow-500',
    info: 'text-indigo-500'
  }
  return colors[type] || colors.info
}

const getToastIcon = (type) => {
  const icons = {
    success: CheckCircleIcon,
    error: XCircleIcon,
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon
  }
  return icons[type] || icons.info
}

// Page transition logic
const getPageTransition = (route) => {
  // Different transitions for different route types
  if (route.path.includes('/auth/')) return 'slide-left'
  if (route.path.includes('/dashboard/')) return 'fade'
  if (route.path.includes('/profile/')) return 'slide-right'
  return 'fade'
}

// Modal management
const openModal = () => {
  isModalOpen.value = true
  document.body.style.overflow = 'hidden'
}

const closeModal = () => {
  isModalOpen.value = false
  document.body.style.overflow = 'auto'
}

// PWA Installation
const handleBeforeInstallPrompt = (e) => {
  e.preventDefault()
  deferredPrompt = e
  showPWAPrompt.value = true
}

const installPWA = async () => {
  if (!deferredPrompt) return
  
  deferredPrompt.prompt()
  const { outcome } = await deferredPrompt.userChoice
  
  if (outcome === 'accepted') {
    addToast('success', 'App Installed!', 'Temple Manager has been added to your home screen')
  }
  
  deferredPrompt = null
  showPWAPrompt.value = false
}

const dismissPWAPrompt = () => {
  showPWAPrompt.value = false
  deferredPrompt = null
}

// Network status handlers
const handleOnline = () => {
  isOnline.value = true
  addToast('success', 'Back Online', 'Internet connection restored')
}

const handleOffline = () => {
  isOnline.value = false
}

// Global error handler
const handleGlobalError = (error) => {
  console.error('Global error:', error)
  addToast('error', 'Something went wrong', 'Please try again or contact support')
}

// Authentication initialization - MODIFIED for role-based redirection
const initializeAuth = async () => {
  console.log('App mounting, initializing auth...')
  try {
    isInitializing.value = true
    
    // Initialize auth state from localStorage
    const isAuth = await authStore.initialize()
    console.log('Auth initialized, authenticated:', isAuth)
    
    // Check if we need to redirect
    if (isAuth) {
      // If user is trying to access login/register while logged in, redirect to appropriate dashboard
      if (isAuthRoute.value) {
        // Get user role
        const userRole = authStore.userRole
        console.log('Authenticated user with role:', userRole)
        
        // Determine dashboard path based on role
        let dashboardPath
        
        if (userRole === 'superadmin' || userRole === 'super_admin') {
          dashboardPath = '/superadmin/dashboard'
        } else if (userRole === 'templeadmin' || userRole === 'tenant') {
          dashboardPath = '/tenant/dashboard'
        } else if (userRole === 'devotee') {
          const entityId = authStore.user?.entityId || authStore.user?.current_entity?.id
          dashboardPath = entityId ? `/entity/${entityId}/devotee/dashboard` : '/devotee/temple-selection'
        } else if (userRole === 'volunteer') {
          const entityId = authStore.user?.entityId || authStore.user?.current_entity?.id
          dashboardPath = entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
        } else {
          dashboardPath = '/'
        }
        
        console.log('Redirecting authenticated user to dashboard:', dashboardPath)
        router.replace(dashboardPath)
      }
    } else {
      // If user is trying to access protected route while not logged in
      if (route.meta?.requiresAuth) {
        console.log('Redirecting unauthenticated user from protected route to login')
        router.replace('/login')
      }
    }
    
    // Show welcome toast for authenticated users
    if (isAuth && !isAuthRoute.value) {
      const userName = authStore.user?.fullName || authStore.user?.name || 'User'
      addToast('success', `Welcome back, ${userName}!`)
    }
    
  } catch (error) {
    console.error('Failed to initialize auth:', error)
    initializationFailed.value = true
    addToast('error', 'Authentication Error', 'Failed to initialize application')
  } finally {
    // Always finish initialization, even if it fails
    isInitializing.value = false
  }
}

// Helper methods for error recovery
const refreshPage = () => {
  window.location.reload()
}

const resetAndLogout = async () => {
  try {
    await authStore.logout()
    addToast('info', 'Session reset', 'Please log in again')
    router.replace('/login')
  } catch (error) {
    console.error('Error during reset:', error)
    // Force logout and redirect
    localStorage.clear()
    window.location.href = '/login'
  }
}

// Route change handler for loading states
const handleRouteChange = () => {
  isGlobalLoading.value = true
  
  // Simulate loading delay (replace with actual route loading logic)
  setTimeout(() => {
    isGlobalLoading.value = false
  }, 500)
}

// Lifecycle hooks
onMounted(async () => {
  // Initialize authentication first
  await initializeAuth()
  
  // Add event listeners
  window.addEventListener('online', handleOnline)
  window.addEventListener('offline', handleOffline)
  window.addEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.addEventListener('error', handleGlobalError)
  
  // Add global styles
  document.documentElement.style.setProperty('--color-primary', '#4f46e5') // indigo-600
  document.documentElement.style.setProperty('--color-primary-dark', '#3730a3') // indigo-800
  document.documentElement.style.setProperty('--color-primary-light', '#6366f1') // indigo-500
  
  // Set up router loading
  router.beforeEach((to, from, next) => {
    if (to.path !== from.path) {
      handleRouteChange()
    }
    next()
  })
  
  // Global keyboard shortcuts
  document.addEventListener('keydown', (e) => {
    // Escape key closes modals
    if (e.key === 'Escape' && isModalOpen.value) {
      closeModal()
    }
  })
})

onUnmounted(() => {
  // Clean up event listeners
  window.removeEventListener('online', handleOnline)
  window.removeEventListener('offline', handleOffline)
  window.removeEventListener('beforeinstallprompt', handleBeforeInstallPrompt)
  window.removeEventListener('error', handleGlobalError)
  
  // Reset body overflow
  document.body.style.overflow = 'auto'
})

// Expose methods globally (for use in other components)
window.showToast = addToast
window.showModal = openModal
window.hideModal = closeModal
</script>

<style scoped>
/* Page Transitions */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

.slide-left-enter-active,
.slide-left-leave-active {
  transition: all 0.3s ease;
}

.slide-left-enter-from {
  opacity: 0;
  transform: translateX(-30px);
}

.slide-left-leave-to {
  opacity: 0;
  transform: translateX(30px);
}

.slide-right-enter-active,
.slide-right-leave-active {
  transition: all 0.3s ease;
}

.slide-right-enter-from {
  opacity: 0;
  transform: translateX(30px);
}

.slide-right-leave-to {
  opacity: 0;
  transform: translateX(-30px);
}

/* Toast Transitions */
.toast-move,
.toast-enter-active,
.toast-leave-active {
  transition: all 0.3s ease;
}

.toast-enter-from {
  opacity: 0;
  transform: translateX(100%);
}

.toast-leave-to {
  opacity: 0;
  transform: translateX(100%);
}

/* Modal Transitions */
.modal-enter-active,
.modal-leave-active {
  transition: opacity 0.3s ease;
}

.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}

.modal-content-enter-active,
.modal-content-leave-active {
  transition: all 0.3s ease;
}

.modal-content-enter-from,
.modal-content-leave-to {
  opacity: 0;
  transform: scale(0.95);
}

/* Slide Up Transitions */
.slide-up-enter-active,
.slide-up-leave-active {
  transition: transform 0.3s ease;
}

.slide-up-enter-from,
.slide-up-leave-to {
  transform: translateY(100%);
}

/* Custom Scrollbar */
:deep(*::-webkit-scrollbar) {
  width: 8px;
  height: 8px;
}

:deep(*::-webkit-scrollbar-track) {
  background: #f1f5f9;
  border-radius: 4px;
}

:deep(*::-webkit-scrollbar-thumb) {
  background: #cbd5e1;
  border-radius: 4px;
}

:deep(*::-webkit-scrollbar-thumb:hover) {
  background: #94a3b8;
}

/* Focus styles */
:deep(*:focus) {
  outline: none;
}

:deep(*:focus-visible) {
  outline: 2px solid #4f46e5;
  outline-offset: 2px;
}

/* Selection styles */
:deep(::selection) {
  background-color: #4f46e5;
  color: white;
}

/* Animations */
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

/* Responsive adjustments */
@media (max-width: 640px) {
  .toast-container {
    left: 1rem;
    right: 1rem;
  }
}

/* Print styles */
@media print {
  .fixed {
    display: none !important;
  }
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .bg-indigo-600 {
    background-color: #000 !important;
  }
  
  .text-indigo-600 {
    color: #000 !important;
  }
  
  .border-indigo-600 {
    border-color: #000 !important;
  }
}

/* Reduced motion support */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}
</style>