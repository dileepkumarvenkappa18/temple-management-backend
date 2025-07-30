<template>
  <Teleport to="body">
    <div class="fixed top-4 right-4 z-50 space-y-2">
      <Transition
        v-for="toast in toasts"
        :key="toast.id"
        enter-active-class="transform transition-all duration-300 ease-out"
        enter-from-class="translate-x-full opacity-0"
        enter-to-class="translate-x-0 opacity-100"
        leave-active-class="transform transition-all duration-200 ease-in"
        leave-from-class="translate-x-0 opacity-100"
        leave-to-class="translate-x-full opacity-0"
      >
        <div
          class="max-w-sm w-full shadow-lg rounded-xl pointer-events-auto overflow-hidden border-l-4"
          :class="getToastClasses(toast.type)"
        >
          <div class="p-4">
            <div class="flex items-center">
              <!-- Icon -->
              <div class="flex-shrink-0">
                <component 
                  :is="getIcon(toast.type)" 
                  :class="getIconClasses(toast.type)"
                  class="h-5 w-5"
                />
              </div>
              
              <!-- Content -->
              <div class="ml-3 flex-1">
                <p class="text-sm font-semibold" :class="getTitleClasses(toast.type)">
                  {{ toast.title }}
                </p>
                <p v-if="toast.message" class="mt-1 text-sm" :class="getMessageClasses(toast.type)">
                  {{ toast.message }}
                </p>
              </div>
              
              <!-- Close Button -->
              <div class="ml-4 flex-shrink-0 flex">
                <button
                  @click="removeToast(toast.id)"
                  class="inline-flex rounded-md focus:outline-none focus:ring-2 focus:ring-offset-2 transition-colors duration-200"
                  :class="getCloseButtonClasses(toast.type)"
                >
                  <span class="sr-only">Close</span>
                  <XMarkIcon class="h-4 w-4" />
                </button>
              </div>
            </div>
            
            <!-- Action Buttons -->
            <div v-if="toast.actions && toast.actions.length" class="mt-3 flex space-x-2">
              <button
                v-for="action in toast.actions"
                :key="action.label"
                @click="handleAction(action, toast)"
                class="text-xs font-medium px-3 py-1 rounded-md transition-colors duration-200"
                :class="getActionButtonClasses(toast.type)"
              >
                {{ action.label }}
              </button>
            </div>
          </div>
          
          <!-- Progress Bar -->
          <div 
            v-if="toast.autoHide && toast.duration"
            class="h-1 transition-all duration-100 ease-linear"
            :class="getProgressClasses(toast.type)"
            :style="{ width: `${getProgress(toast)}%` }"
          ></div>
        </div>
      </Transition>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { 
  CheckCircleIcon, 
  XCircleIcon, 
  ExclamationTriangleIcon, 
  InformationCircleIcon,
  XMarkIcon 
} from '@heroicons/vue/24/solid'

// Toast store
const toasts = ref([])
let toastId = 0

// Toast methods
const addToast = (toast) => {
  const id = ++toastId
  const newToast = {
    id,
    type: 'success',
    autoHide: true,
    duration: 5000,
    startTime: Date.now(),
    ...toast
  }
  
  toasts.value.push(newToast)
  
  if (newToast.autoHide) {
    setTimeout(() => {
      removeToast(id)
    }, newToast.duration)
  }
  
  return id
}

const removeToast = (id) => {
  const index = toasts.value.findIndex(toast => toast.id === id)
  if (index > -1) {
    toasts.value.splice(index, 1)
  }
}

const clearAll = () => {
  toasts.value = []
}

// Toast type helpers
const getToastClasses = (type) => {
  const classes = {
    success: 'bg-green-50 border-green-400',
    error: 'bg-red-50 border-red-400',
    warning: 'bg-yellow-50 border-yellow-400',
    info: 'bg-blue-50 border-blue-400'
  }
  return classes[type] || classes.info
}

const getIcon = (type) => {
  const icons = {
    success: CheckCircleIcon,
    error: XCircleIcon,
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon
  }
  return icons[type] || icons.info
}

const getIconClasses = (type) => {
  const classes = {
    success: 'text-green-600',
    error: 'text-red-600',
    warning: 'text-yellow-600',
    info: 'text-blue-600'
  }
  return classes[type] || classes.info
}

const getTitleClasses = (type) => {
  const classes = {
    success: 'text-green-900',
    error: 'text-red-900',
    warning: 'text-yellow-900',
    info: 'text-blue-900'
  }
  return classes[type] || classes.info
}

const getMessageClasses = (type) => {
  const classes = {
    success: 'text-green-700',
    error: 'text-red-700',
    warning: 'text-yellow-700',
    info: 'text-blue-700'
  }
  return classes[type] || classes.info
}

const getCloseButtonClasses = (type) => {
  const classes = {
    success: 'text-green-400 hover:text-green-600 focus:ring-green-500',
    error: 'text-red-400 hover:text-red-600 focus:ring-red-500',
    warning: 'text-yellow-400 hover:text-yellow-600 focus:ring-yellow-500',
    info: 'text-blue-400 hover:text-blue-600 focus:ring-blue-500'
  }
  return classes[type] || classes.info
}

const getActionButtonClasses = (type) => {
  const classes = {
    success: 'bg-green-100 text-green-800 hover:bg-green-200',
    error: 'bg-red-100 text-red-800 hover:bg-red-200',
    warning: 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200',
    info: 'bg-blue-100 text-blue-800 hover:bg-blue-200'
  }
  return classes[type] || classes.info
}

const getProgressClasses = (type) => {
  const classes = {
    success: 'bg-green-400',
    error: 'bg-red-400',
    warning: 'bg-yellow-400',
    info: 'bg-blue-400'
  }
  return classes[type] || classes.info
}

const getProgress = (toast) => {
  if (!toast.autoHide || !toast.duration) return 100
  
  const elapsed = Date.now() - toast.startTime
  const progress = Math.max(0, 100 - (elapsed / toast.duration) * 100)
  return progress
}

const handleAction = (action, toast) => {
  if (action.handler) {
    action.handler()
  }
  if (action.closeOnClick !== false) {
    removeToast(toast.id)
  }
}

// Convenience methods
const success = (title, message, options = {}) => {
  return addToast({ type: 'success', title, message, ...options })
}

const error = (title, message, options = {}) => {
  return addToast({ type: 'error', title, message, ...options })
}

const warning = (title, message, options = {}) => {
  return addToast({ type: 'warning', title, message, ...options })
}

const info = (title, message, options = {}) => {
  return addToast({ type: 'info', title, message, ...options })
}

// Expose methods
defineExpose({
  addToast,
  removeToast,
  clearAll,
  success,
  error,
  warning,
  info
})

// Global toast instance
let toastInstance = null

onMounted(() => {
  // Make this the global toast instance
  toastInstance = {
    success,
    error,
    warning,
    info,
    addToast,
    removeToast,
    clearAll
  }
  
  // Attach to window for global access (optional)
  window.$toast = toastInstance
})

onUnmounted(() => {
  if (window.$toast === toastInstance) {
    delete window.$toast
  }
})
</script>

<style scoped>
/* Custom animations for better UX */
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
</style>