// src/stores/auth.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import api from '@/plugins/axios'
import { USER_ROLES } from '@/utils/constants'
import router from '@/router'

export const useAuthStore = defineStore('auth', () => {
  // State
  const userDataRaw = localStorage.getItem('user_data')
  const user = ref(userDataRaw && userDataRaw !== 'undefined' ? JSON.parse(userDataRaw) : null)
  const token = ref(localStorage.getItem('auth_token') || null)
  const loading = ref(false)
  const error = ref(null)
  const currentTenantId = ref(localStorage.getItem('current_tenant_id') || null)
  
  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => {
    if (!user.value) return null
    
    // Handle both string roles and numeric role_id
    if (typeof user.value.role === 'string') {
      return user.value.role
    }
    
    // Map role_id to string role if needed
    const ROLE_MAP = {
      1: 'superadmin',
      2: 'templeadmin',
      3: 'devotee',
      4: 'volunteer'
    }
    
    return user.value.roleId ? ROLE_MAP[user.value.roleId] : null
  })
  
  // Dashboard route getter
  const dashboardRoute = computed(() => {
    if (!user.value) return '/login'
    
    const role = userRole.value
    console.log('Determining dashboard route for role:', role)
    
    const entityId = user.value.entityId || user.value.current_entity?.id
    
    if (role === 'superadmin' || role === 'super_admin') {
      return '/superadmin/dashboard'
    } else if (role === 'templeadmin' || role === 'tenant') {
      // Use direct path to tenant dashboard with ID
      const tenantId = user.value.id || currentTenantId.value
      if (tenantId) {
        return `/tenant/${tenantId}/dashboard`
      }
      return '/tenant/dashboard'
    } else if (role === 'devotee') {
      // MODIFIED: Always return temple selection for devotees
      return '/devotee/temple-selection'
    } else if (role === 'volunteer') {
      return entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
    }
    
    return '/'
  })
  
  // Role-specific getters
  const isTenant = computed(() => userRole.value === USER_ROLES.TENANT || userRole.value === 'templeadmin' || userRole.value === 'tenant')
  const isDevotee = computed(() => userRole.value === USER_ROLES.DEVOTEE || userRole.value === 'devotee')
  const isVolunteer = computed(() => userRole.value === USER_ROLES.VOLUNTEER || userRole.value === 'volunteer')
  const isSuperAdmin = computed(() => userRole.value === USER_ROLES.SUPER_ADMIN || userRole.value === 'superadmin' || userRole.value === 'super_admin')
  const isEndUser = computed(() => isDevotee.value || isVolunteer.value)
  
  // Clear all browser storage completely
  const clearAllStorage = () => {
    console.log('Clearing all browser storage')
    localStorage.clear()
    sessionStorage.clear()
    
    // Try to clear cookies related to the app
    document.cookie.split(";").forEach(function(c) {
      document.cookie = c.replace(/^ +/, "").replace(/=.*/, "=;expires=" + new Date().toUTCString() + ";path=/");
    });
    
    console.log('All storage cleared')
  }
  
  // Reset all app state
  const resetAppState = () => {
    console.log('Resetting all application state')
    
    // Clear browser storage and cache
    clearAllStorage()
    
    // Reset API headers
    delete api.defaults.headers.common['Authorization']
    delete api.defaults.headers.common['X-Tenant-ID']
    
    console.log('Application state reset complete')
  }
  
  // Initialize auth state
  const initialize = () => {
  const storedToken = localStorage.getItem('auth_token')
  if (storedToken) {
    token.value = storedToken

    const storedUser = localStorage.getItem('user_data')
    if (storedUser && storedUser !== 'undefined') {
      try {
        user.value = JSON.parse(storedUser)
        // ❌ REMOVE this: isAuthenticated.value = true
        api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
        console.log('Auth initialized with stored user:', user.value)
        console.log('Set Authorization header with token')
      } catch (e) {
        console.error('Failed to parse stored user_data:', e)
      }
    } else {
      console.warn('Token exists but no valid user data found')
    }
  }
}

  
  // Login action - only uses real backend
  const login = async (credentials) => {
    loading.value = true
    error.value = null
    
    try {
      console.log('Starting login request')
      
      // Reset all app state before login to ensure clean isolation
      resetAppState()
      
      console.log('Performing backend login')
      
      const response = await api.post('/v1/auth/login', credentials)
      
      console.log('Login response received:', response.data)
      
      // Extract data from response based on backend structure
      const data = response.data
      const accessToken = data.accessToken || data.token || data.access_token || data.jwt
      const userData = data.user || data.userData || data
      
      if (!accessToken) {
        throw new Error('No authentication token received from server')
      }
      
      // Store auth data
      token.value = accessToken
      user.value = userData
      localStorage.setItem('auth_token', accessToken)
      // Clear old token key
      localStorage.setItem('user_data', JSON.stringify(userData))
      
      // Set axios default header for authentication
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
      console.log('Set Authorization header with token')
      
      // For tenant users, store and set tenant ID
      if (userData.roleId === 2 || userData.role === 'templeadmin' || userData.role === 'tenant') {
        const tenantId = userData.id
        currentTenantId.value = tenantId
        localStorage.setItem('current_tenant_id', tenantId)
        
        // Set tenant header for API calls
        api.defaults.headers.common['X-Tenant-ID'] = tenantId
        console.log('Set X-Tenant-ID header:', tenantId)
      }
      
      // Determine where to redirect based on role
      let redirectPath = '/'
      if (userData.roleId === 1 || userData.role === 'superadmin' || userData.role === 'super_admin') {
        redirectPath = '/superadmin/dashboard'
      } else if (userData.roleId === 2 || userData.role === 'templeadmin' || userData.role === 'tenant') {
        const tenantId = userData.id
        redirectPath = `/tenant/${tenantId}/dashboard`
        console.log(`Setting tenant-specific redirect path: ${redirectPath}`)
      } else if (userData.roleId === 3 || userData.role === 'devotee') {
        // MODIFIED: Always redirect devotees to temple selection regardless of entityId
        redirectPath = '/devotee/temple-selection'
      } else if (userData.roleId === 4 || userData.role === 'volunteer') {
        redirectPath = userData.entityId 
          ? `/entity/${userData.entityId}/volunteer/dashboard` 
          : '/volunteer/temple-selection'
      }
      
      console.log('Login successful, will redirect to:', redirectPath)
      
      // DIRECT APPROACH: Force navigation for tenant
      if (userData.roleId === 2 || userData.role === 'templeadmin' || userData.role === 'tenant') {
        console.log('DIRECT NAVIGATION: Forcing hard redirect for tenant dashboard')
        
        // Store the redirect info first
        const result = {
          success: true,
          user: userData,
          redirectPath
        }
        
        // Give the browser time to process storage changes
        setTimeout(() => {
          // Force a hard navigation to bypass Vue Router
          window.location.href = window.location.origin + redirectPath
        }, 500)
        
        return result
      }
      
      return {
        success: true,
        user: userData,
        redirectPath
      }
    } catch (err) {
      console.error('Login error:', err)
      error.value = err.response?.data?.message || err.response?.data?.error || err.message || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
      console.log('Login process completed, loading set to false')
    }
  }
  
  // Get dashboard path based on role and entity
  const getDashboardPath = () => {
    if (!user.value) return '/login'
    
    const role = userRole.value
    const entityId = user.value.entityId || user.value.current_entity?.id
    
    if (role === 'superadmin' || role === 'super_admin') {
      return '/superadmin/dashboard'
    } else if (role === 'templeadmin' || role === 'tenant') {
      const tenantId = user.value.id || currentTenantId.value
      return `/tenant/${tenantId}/dashboard`
    } else if (role === 'devotee') {
      // MODIFIED: Always return temple selection for devotees
      return '/devotee/temple-selection'
    } else if (role === 'volunteer') {
      return entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
    }
    
    return '/'
  }
  
  // Logout action
  const logout = async () => {
    try {
      if (token.value) {
        try {
          await api.post('/v1/auth/logout')
        } catch (logoutErr) {
          console.warn('Logout endpoint error:', logoutErr)
        }
      }
    } catch (err) {
      console.error('Logout error:', err)
    } finally {
      // Reset all application state
      resetAppState()
      
      // Clear auth state variables
      user.value = null
      token.value = null
      currentTenantId.value = null
      
      // Force a complete page reload to clear any Vue Router state
      window.location.href = window.location.origin + '/login'
    }
  }
  
  // Register new user
  const register = async (userData) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.post('/v1/auth/register', userData)
      
      const data = response.data
      
      // Handle different registration flows
      const redirectPath = userData.role === 'templeadmin' 
        ? '/auth/pending-approval'
        : userData.role === 'devotee'
          ? '/devotee/temple-selection'
          : '/volunteer/temple-selection'
          
      return {
        success: true,
        message: data.message || 'Registration successful!',
        redirectPath
      }
    } catch (err) {
      console.error('Registration failed:', err)
      error.value = err.response?.data?.message || err.message || 'Registration failed'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }
  
  // Join temple (for devotees/volunteers)
  const joinTemple = async (templeId) => {
    loading.value = true
    error.value = null
    
    try {
      // Different endpoints based on role
      const endpoint = isDevotee.value 
        ? '/v1/devotee/join-temple' 
        : '/v1/volunteer/join-temple'
        
      const response = await api.post(endpoint, { templeId })
      
      const data = response.data
      
      // Update user data with temple
      const temple = data.temple || data.entity || data
      const updatedUser = {
        ...user.value,
        entityId: templeId,
        current_entity: temple
      }
      
      user.value = updatedUser
      localStorage.setItem('user_data', JSON.stringify(updatedUser))

      localStorage.setItem('current_entity_id', templeId)
      
      // Get redirect path
      const redirectPath = getDashboardPath()
      
      return {
        success: true,
        message: `Successfully joined temple!`,
        redirectPath
      }
    } catch (err) {
      console.error('Failed to join temple:', err)
      error.value = err.response?.data?.message || err.message || 'Failed to join temple'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }
  
  // Debug function to check if we can see the tenant data
  const debugTenantData = () => {
    console.log('------- DEBUG TENANT DATA -------')
    console.log('Current user:', user.value)
    console.log('Current tenant ID:', currentTenantId.value)
    console.log('Is tenant?', isTenant.value)
    console.log('Dashboard route:', dashboardRoute.value)
    console.log('Auth headers:', api.defaults.headers.common)
    console.log('Local storage:', {
      auth_token: localStorage.getItem('auth_token'),
      user_data: localStorage.getItem('user_data'),
      current_tenant_id: localStorage.getItem('current_tenant_id')
    })
    console.log('--------------------------------')
  }
  
  return {
    // State
    user,
    token,
    loading,
    error,
    currentTenantId,
    
    // Getters
    isAuthenticated,
    userRole,
    isTenant,
    isDevotee,
    isVolunteer,
    isSuperAdmin,
    isEndUser,
    dashboardRoute,
    
    // Actions
    initialize,
    login,
    logout,
    register,
    joinTemple,
    getDashboardPath,
    resetAppState,
    clearAllStorage,
    debugTenantData,
    
    // For backwards compatibility
    initializeAuth: initialize,
    clearError: () => { error.value = null }
  }
})
