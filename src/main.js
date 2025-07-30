import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'

// Composables
import { useAuth } from './composables/useAuth'
import { useToast } from './composables/useToast'

// Global CSS with Tailwind v4 and Indigo Theme
import './assets/main.css'

// Create Vue app instance
const app = createApp(App)

// Create Pinia store instance
const pinia = createPinia()

// Install plugins
app.use(pinia)
app.use(router)

// Global error handler
app.config.errorHandler = (error, instance, info) => {
  console.error('Global error:', error)
  console.error('Component instance:', instance)
  console.error('Error info:', info)
  
  // Toast notification for production errors
  if (process.env.NODE_ENV === 'production') {
    // Use a safer way to access toast without potential circular dependency
    try {
      const { showToast } = useToast()
      showToast('An unexpected error occurred. Please try again.', 'error')
    } catch (toastError) {
      console.error('Failed to show error toast:', toastError)
    }
  }
}

// Performance monitoring
app.config.performance = process.env.NODE_ENV === 'development'

// Global constants
app.config.globalProperties.$constants = {
  // App metadata
  APP_NAME: 'Temple SaaS Platform',
  APP_VERSION: '1.0.0',
  
  // Roles
  ROLES: {
    SUPER_ADMIN: 'superadmin',
    TENANT: 'tenant',
    DEVOTEE: 'devotee',
    VOLUNTEER: 'volunteer'
  },
  
  // Status types
  STATUS: {
    PENDING: 'PENDING',
    APPROVED: 'APPROVED',
    REJECTED: 'REJECTED',
    ACTIVE: 'ACTIVE',
    INACTIVE: 'INACTIVE'
  },
  
  // Temple approval statuses
  TEMPLE_STATUS: {
    PENDING: 'PENDING',
    APPROVED: 'APPROVED',
    REJECTED: 'REJECTED',
    UNDER_REVIEW: 'UNDER_REVIEW'
  },
  
  // Seva types
  SEVA_TYPES: {
    DAILY: 'daily',
    WEEKLY: 'weekly',
    MONTHLY: 'monthly',
    FESTIVAL: 'festival',
    SPECIAL: 'special'
  },
  
  // Donation types
  DONATION_TYPES: {
    GENERAL: 'general',
    SEVA: 'seva',
    FESTIVAL: 'festival',
    DEVELOPMENT: 'development',
    FOOD: 'food',
    MAINTENANCE: 'maintenance'
  },
  
  // Communication channels
  COMMUNICATION_CHANNELS: {
    EMAIL: 'email',
    SMS: 'sms',
    WHATSAPP: 'whatsapp',
    IN_APP: 'in_app'
  },
  
  // File upload limits
  FILE_LIMITS: {
    MAX_SIZE: 5 * 1024 * 1024, // 5MB
    ALLOWED_TYPES: [
      'image/jpeg', 
      'image/png', 
      'image/jpg', 
      'application/pdf', 
      'application/msword', 
      'application/vnd.openxmlformats-officedocument.wordprocessingml.document'
    ]
  },
  
  // Theme colors (Indigo-based)
  THEME: {
    PRIMARY: 'indigo-600',
    SECONDARY: 'indigo-500',
    ACCENT: 'indigo-800',
    SUCCESS: 'green-600',
    ERROR: 'red-600',
    WARNING: 'yellow-600',
    INFO: 'blue-600',
    GRAY: 'gray-600'
  },
  
  // Responsive breakpoints
  BREAKPOINTS: {
    SM: '640px',
    MD: '768px',
    LG: '1024px',
    XL: '1280px',
    '2XL': '1536px'
  }
}

// Provide global composables - Initialize after app is created but before mount
const authComposable = useAuth()
const toastComposable = useToast()

app.provide('auth', authComposable)
app.provide('toast', toastComposable)

// Global properties for easy access (safer initialization)
app.config.globalProperties.$auth = authComposable
app.config.globalProperties.$toast = toastComposable

// Development tools
if (process.env.NODE_ENV === 'development') {
  app.config.globalProperties.$log = console.log
  app.config.globalProperties.$debug = true
  
  // Dev helper functions
  app.config.globalProperties.$devHelpers = {
    // Quick role switcher for development
    switchRole: (role) => {
      try {
        const { setUser } = authComposable
        setUser({
          id: 1,
          email: 'dev@temple.com',
          fullName: 'Dev User',
          role: role,
          isActive: true
        })
      } catch (error) {
        console.error('Error switching role:', error)
      }
    },
    
    // Quick temple context setter
    setTempleContext: (templeId, templeName) => {
      try {
        localStorage.setItem('currentTempleId', templeId)
        localStorage.setItem('currentTempleName', templeName)
      } catch (error) {
        console.error('Error setting temple context:', error)
      }
    },
    
    // Clear all local storage
    clearStorage: () => {
      try {
        localStorage.clear()
        sessionStorage.clear()
        console.log('Storage cleared successfully')
      } catch (error) {
        console.error('Error clearing storage:', error)
      }
    },
    
    // ADDED: Login debugging tools
    // ============================
    
    // Analyze user data in local storage
    analyzeUser: () => {
      try {
        const userData = localStorage.getItem('user_data')
        
        if (!userData) {
          console.log('❌ No user data found in localStorage')
          return null
        }
        
        const user = JSON.parse(userData)
        
        console.group('🔍 User Data Analysis')
        console.log('User object:', user)
        
        // Check for role properties
        const roleProperties = ['role', 'roleId', 'role_id', 'userType', 'user_type', 'user_role', 'type']
        
        console.log('Role-related properties:')
        roleProperties.forEach(prop => {
          if (user[prop] !== undefined) {
            console.log(`- ${prop}: ${user[prop]} (type: ${typeof user[prop]})`)
          }
        })
        
        // Check for entity/temple ID
        const entityProps = ['entityId', 'templeId', 'entity_id', 'temple_id']
        console.log('Entity-related properties:')
        entityProps.forEach(prop => {
          if (user[prop] !== undefined) {
            console.log(`- ${prop}: ${user[prop]}`)
          }
        })
        
        // Current entity object
        if (user.current_entity) {
          console.log('Current entity:', user.current_entity)
        }
        
        console.groupEnd()
        return user
      } catch (error) {
        console.error('Error analyzing user data:', error)
        return null
      }
    },
    
    // Force a specific role for testing
    forceRole: (role) => {
      try {
        const userData = localStorage.getItem('user_data')
        
        if (!userData) {
          console.log('❌ No user data found in localStorage')
          return false
        }
        
        const user = JSON.parse(userData)
        
        // Backup original role
        const originalRole = user.role
        
        // Update role
        user.role = role
        
        // Save back to localStorage
        localStorage.setItem('user_data', JSON.stringify(user))
        
        console.log(`✅ Role changed from "${originalRole}" to "${role}"`)
        console.log('Refresh the page for changes to take effect')
        
        return true
      } catch (error) {
        console.error('Error modifying user role:', error)
        return false
      }
    },
    
    // Get dashboard path for current user
    getDashboardPath: () => {
      try {
        const userData = localStorage.getItem('user_data')
        
        if (!userData) {
          console.log('❌ No user data found in localStorage')
          return null
        }
        
        const user = JSON.parse(userData)
        
        // Simple role normalization
        let role = user.role
        
        // Map common role variations
        const roleMap = {
          'templeadmin': 'tenant',
          'temple_admin': 'tenant',
          'entity_admin': 'tenant',
          'super_admin': 'superadmin',
          'admin': 'superadmin',
          'bhakta': 'devotee',
          'user': 'devotee',
          'seva_volunteer': 'volunteer'
        }
        
        // Apply mapping if it exists
        role = roleMap[role] || role
        
        // Get entity ID
        const entityId = user.entityId || (user.current_entity ? user.current_entity.id : null)
        
        // Determine dashboard path
        let path
        
        switch (role) {
          case 'superadmin':
            path = '/superadmin/dashboard'
            break
          case 'tenant':
            path = '/tenant/dashboard'
            break
          case 'devotee':
            path = entityId ? `/entity/${entityId}/devotee/dashboard` : '/devotee/temple-selection'
            break
          case 'volunteer':
            path = entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
            break
          default:
            path = '/'
        }
        
        console.log(`Dashboard path for role "${role}": ${path}`)
        return path
      } catch (error) {
        console.error('Error determining dashboard path:', error)
        return null
      }
    },
    
    // Show login debug help
    showLoginHelp: () => {
      console.group('🛠️ Login Debug Tools')
      console.log('$devHelpers.analyzeUser() - Analyze current user data in localStorage')
      console.log('$devHelpers.forceRole("role") - Override current user role (e.g., "tenant", "superadmin")')
      console.log('$devHelpers.getDashboardPath() - Get dashboard path for current user')
      console.log('$devHelpers.clearStorage() - Clear all authentication data')
      console.groupEnd()
    }
  }
  
  // Auto-show debug help on startup
  console.log('🔧 Development tools loaded! Access debugging tools via app.$devHelpers')
  console.log('Try app.$devHelpers.showLoginHelp() for login debugging commands')
  
  // Make tools available on window for console access
  window.$devHelpers = app.config.globalProperties.$devHelpers
}

// Global directive for permissions
app.directive('permission', {
  mounted(el, binding) {
    try {
      const { hasPermission } = authComposable
      const permission = binding.value
      
      if (!hasPermission || !hasPermission(permission)) {
        el.style.display = 'none'
      }
    } catch (error) {
      console.error('Permission directive error:', error)
      // Fail safely - hide element if permission check fails
      el.style.display = 'none'
    }
  },
  updated(el, binding) {
    try {
      const { hasPermission } = authComposable
      const permission = binding.value
      
      if (!hasPermission || !hasPermission(permission)) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    } catch (error) {
      console.error('Permission directive error:', error)
      el.style.display = 'none'
    }
  }
})

// Global directive for role-based visibility
app.directive('role', {
  mounted(el, binding) {
    try {
      const { user } = authComposable
      const allowedRoles = Array.isArray(binding.value) ? binding.value : [binding.value]
      
      if (!user?.value || !allowedRoles.includes(user.value.role)) {
        el.style.display = 'none'
      }
    } catch (error) {
      console.error('Role directive error:', error)
      el.style.display = 'none'
    }
  },
  updated(el, binding) {
    try {
      const { user } = authComposable
      const allowedRoles = Array.isArray(binding.value) ? binding.value : [binding.value]
      
      if (!user?.value || !allowedRoles.includes(user.value.role)) {
        el.style.display = 'none'
      } else {
        el.style.display = ''
      }
    } catch (error) {
      console.error('Role directive error:', error)
      el.style.display = 'none'
    }
  }
})

// Global directive for loading states
app.directive('loading', {
  mounted(el, binding) {
    updateLoadingState(el, binding.value)
  },
  updated(el, binding) {
    updateLoadingState(el, binding.value)
  }
})

// Helper function for loading directive
function updateLoadingState(el, isLoading) {
  try {
    if (isLoading) {
      el.classList.add('opacity-50', 'pointer-events-none', 'cursor-not-allowed')
      
      // Add loading spinner if not already present
      if (!el.querySelector('.loading-spinner')) {
        const spinner = document.createElement('div')
        spinner.className = 'loading-spinner absolute inset-0 flex items-center justify-center bg-white bg-opacity-75'
        spinner.innerHTML = `
          <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
        `
        
        // Ensure parent has relative positioning
        const currentPosition = window.getComputedStyle(el).position
        if (currentPosition === 'static') {
          el.style.position = 'relative'
        }
        
        el.appendChild(spinner)
      }
    } else {
      el.classList.remove('opacity-50', 'pointer-events-none', 'cursor-not-allowed')
      const spinner = el.querySelector('.loading-spinner')
      if (spinner) {
        spinner.remove()
      }
    }
  } catch (error) {
    console.error('Loading directive error:', error)
  }
}

// Global mixin for common functionality
app.mixin({
  methods: {
    // Format currency
    formatCurrency(amount, currency = 'INR') {
      try {
        return new Intl.NumberFormat('en-IN', {
          style: 'currency',
          currency: currency,
          minimumFractionDigits: 0
        }).format(amount)
      } catch (error) {
        console.error('Currency formatting error:', error)
        return `₹${amount}`
      }
    },
    
    // Format date
    formatDate(date, format = 'short') {
      try {
        const options = {
          short: { year: 'numeric', month: 'short', day: 'numeric' },
          long: { year: 'numeric', month: 'long', day: 'numeric', weekday: 'long' },
          time: { hour: '2-digit', minute: '2-digit' },
          datetime: { year: 'numeric', month: 'short', day: 'numeric', hour: '2-digit', minute: '2-digit' }
        }
        
        return new Intl.DateTimeFormat('en-IN', options[format] || options.short).format(new Date(date))
      } catch (error) {
        console.error('Date formatting error:', error)
        return new Date(date).toLocaleDateString()
      }
    },
    
    // Truncate text
    truncate(text, length = 50) {
      if (!text) return ''
      return text.length > length ? text.substring(0, length) + '...' : text
    },
    
    // Capitalize first letter
    capitalize(text) {
      if (!text) return ''
      return text.charAt(0).toUpperCase() + text.slice(1)
    },
    
    // Generate avatar from name
    generateAvatar(name) {
      if (!name) return ''
      const initials = name.split(' ')
        .map(word => word.charAt(0))
        .join('')
        .toUpperCase()
      return initials.substring(0, 2)
    },
    
    // Get status color
    getStatusColor(status) {
      const colors = {
        PENDING: 'bg-yellow-100 text-yellow-800',
        APPROVED: 'bg-green-100 text-green-800',
        REJECTED: 'bg-red-100 text-red-800',
        ACTIVE: 'bg-green-100 text-green-800',
        INACTIVE: 'bg-gray-100 text-gray-800',
        UNDER_REVIEW: 'bg-blue-100 text-blue-800'
      }
      return colors[status] || 'bg-gray-100 text-gray-800'
    },
    
    // Safe navigation helper
    safeGet(obj, path, defaultValue = null) {
      try {
        return path.split('.').reduce((current, key) => current?.[key], obj) ?? defaultValue
      } catch (error) {
        return defaultValue
      }
    }
  }
})

// Mount the app with error handling
try {
  app.mount('#app')
  
  // Log startup in development
  if (process.env.NODE_ENV === 'development') {
    console.log(`🏛️ Temple SaaS Platform v${app.config.globalProperties.$constants.APP_VERSION} started`)
    console.log('🚀 Development mode enabled')
    console.log('🎨 Indigo theme loaded')
    console.log('📱 Responsive design active')
  }
} catch (error) {
  console.error('Failed to mount app:', error)
  
  // Fallback error display
  const appElement = document.getElementById('app')
  if (appElement) {
    appElement.innerHTML = `
      <div class="min-h-screen flex items-center justify-center bg-gray-50">
        <div class="text-center">
          <h1 class="text-2xl font-bold text-red-600 mb-4">Application Failed to Start</h1>
          <p class="text-gray-600 mb-4">Please refresh the page or contact support.</p>
          <button onclick="window.location.reload()" class="bg-indigo-600 text-white px-4 py-2 rounded">
            Refresh Page
          </button>
        </div>
      </div>
    `
  }
}

// Export app instance for testing
export default app