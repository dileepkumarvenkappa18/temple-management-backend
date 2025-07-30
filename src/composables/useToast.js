// src/composables/useToast.js
import { ref } from 'vue'

const toasts = ref([])
let toastId = 0

export const useToast = () => {
  // Toast types with indigo theme styling
  const TOAST_TYPES = {
    SUCCESS: {
      icon: '✅',
      bgClass: 'bg-green-50 border-green-200',
      textClass: 'text-green-800',
      iconClass: 'text-green-600'
    },
    ERROR: {
      icon: '❌',
      bgClass: 'bg-red-50 border-red-200',
      textClass: 'text-red-800',
      iconClass: 'text-red-600'
    },
    WARNING: {
      icon: '⚠️',
      bgClass: 'bg-yellow-50 border-yellow-200',
      textClass: 'text-yellow-800',
      iconClass: 'text-yellow-600'
    },
    INFO: {
      icon: 'ℹ️',
      bgClass: 'bg-indigo-50 border-indigo-200',
      textClass: 'text-indigo-800',
      iconClass: 'text-indigo-600'
    }
  }

  // Main toast function - using regular function instead of arrow function
  // so it has access to the arguments object
  function showToast(message, type = 'info', duration = 5000) {
    // Handle both formats:
    // 1. showToast('message', 'success', 5000)
    // 2. showToast({ message: 'Hello', type: 'success', duration: 5000 })
    
    let finalMessage, finalType, finalDuration;
    
    if (typeof message === 'object') {
      // New format with options object
      finalMessage = message.message;
      finalType = (message.type || 'info').toUpperCase();
      finalDuration = message.duration || 5000;
    } else {
      // Old format with separate arguments
      finalMessage = message;
      finalType = (type || 'info').toUpperCase();
      finalDuration = duration;
    }
    
    return addToast(finalMessage, finalType, finalDuration);
  }

  const addToast = (message, type = 'INFO', duration = 5000) => {
    const id = ++toastId
    const toast = {
      id,
      message,
      type,
      duration,
      ...TOAST_TYPES[type],
      timestamp: Date.now()
    }
    
    toasts.value.push(toast)
    
    // Auto remove toast after duration
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

  const clearAllToasts = () => {
    toasts.value = []
  }

  // Convenience methods for different toast types
  const success = (message, duration = 5000) => {
    return addToast(message, 'SUCCESS', duration)
  }

  const error = (message, duration = 7000) => {
    return addToast(message, 'ERROR', duration)
  }

  const warning = (message, duration = 6000) => {
    return addToast(message, 'WARNING', duration)
  }

  const info = (message, duration = 5000) => {
    return addToast(message, 'INFO', duration)
  }

  return {
    // State
    toasts,
    
    // Methods
    showToast,
    addToast,
    removeToast,
    clearAllToasts,
    
    // Convenience methods
    success,
    error,
    warning,
    info,
    
    // Constants
    TOAST_TYPES
  }
}

// Create and export a singleton instance
export default useToast();