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
  const assignedTenantId = ref(localStorage.getItem('assigned_tenant_id') || null)
  
  // Getters
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  
  // Complete role mapping including ALL backend roleIds
  const userRole = computed(() => {
  if (!user.value) return null

  // Handle both string roles and numeric role_id
  if (typeof user.value.role === 'string') {
    const roleStr = user.value.role.toLowerCase().trim()

    // Normalize role string variations
    if (roleStr === 'monitoringuser' || roleStr === 'monitoring_user') return 'monitoring_user'
    if (roleStr === 'standarduser' || roleStr === 'standard_user') return 'standard_user'
    if (roleStr === 'superadmin' || roleStr === 'super_admin') return 'superadmin'
    if (roleStr === 'templeadmin' || roleStr === 'tenant') return 'templeadmin'
    return roleStr
  }

  // COMPLETE role mapping for ALL backend roleIds
  const ROLE_MAP = {
    1: 'superadmin',
    2: 'templeadmin',
    3: 'devotee',
    4: 'volunteer',
    5: 'standard_user',
    6: 'monitoring_user'
  }

  const mappedRole = user.value.roleId ? ROLE_MAP[user.value.roleId] : null

  console.log('🎭 Role mapping:', {
    roleId: user.value.roleId,
    mappedRole,
    userObject: user.value
  })

  return mappedRole
})
  
  // Parse JWT token to extract claims
  const parseJwt = (token) => {
    try {
      // JWT tokens are base64 encoded and have three parts separated by dots
      const base64Url = token.split('.')[1];
      const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
      const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
      }).join(''));

      return JSON.parse(jsonPayload);
    } catch (e) {
      console.error('Error parsing JWT token:', e);
      return {};
    }
  };
  
  // Dashboard path logic with assigned tenant ID support
  // UPDATED: Changed redirection path for monitoring_user and standard_user to tenant dashboard
  const getDashboardPath = (userRole) => {
    const role = (userRole || '').toLowerCase().trim()
    
    console.log('🎯 Getting dashboard path for role:', role)
    
    // Handle monitoring_user and standard_user with assigned tenant ID
    if (role === 'standard_user' || role === 'standarduser' || 
        role === 'monitoring_user' || role === 'monitoringuser') {
      // Check for assigned tenant ID
      if (assignedTenantId.value) {
        console.log(`✅ ${role} redirecting to assigned tenant: ${assignedTenantId.value}`)
        // UPDATED: Redirect to tenant dashboard instead of entity dashboard
        return `/tenant/${assignedTenantId.value}/dashboard`
      }
      
      // Fallback to tenant selection if no assigned tenant
      console.log(`✅ ${role} with no assigned tenant -> tenant selection`)
      return '/tenant-selection'
    }
    
    if (role === 'superadmin' || role === 'super_admin') {
      console.log('✅ Super admin -> /superadmin/dashboard')
      return '/superadmin/dashboard'
    }
    
    // Handle other roles
    const entityId = user.value?.entityId || user.value?.current_entity?.id
    
    if (role === 'templeadmin' || role === 'tenant') {
      const tenantId = user.value?.id || currentTenantId.value
      return tenantId ? `/tenant/${tenantId}/dashboard` : '/tenant/dashboard'
    } else if (role === 'devotee') {
      return entityId ? `/entity/${entityId}/devotee/dashboard` : '/devotee/temple-selection'
    } else if (role === 'volunteer') {
      return entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
    }
    
    console.log('⚠️ Unknown role, defaulting to home')
    return '/'
  }
  
  // Dashboard route getter
  const dashboardRoute = computed(() => {
    if (!user.value) return '/login'
    return getDashboardPath(userRole.value)
  })
  
  // Role-specific getters
  const isTenant = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'tenant' || role === 'templeadmin';
  })
  
  const isDevotee = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'devotee';
  })
  
  const isVolunteer = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'volunteer';
  })
  
  const isSuperAdmin = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'superadmin' || role === 'super_admin';
  })
  
  const isStandardUser = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'standard_user' || role === 'standarduser';
  })
  
  const isMonitoringUser = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    return role === 'monitoring_user' || role === 'monitoringuser';
  })
  
  const isEndUser = computed(() => isDevotee.value || isVolunteer.value)
  
  // Check if user needs tenant selection
  const needsTenantSelection = computed(() => {
    const role = userRole.value?.toLowerCase() || '';
    // Only need tenant selection if standard/monitoring user AND no assigned tenant
    return (role === 'standard_user' || role === 'standarduser' || 
           role === 'monitoring_user' || role === 'monitoringuser') && 
           !assignedTenantId.value;
  })
  
  // Refresh the auth token
  const refreshToken = async () => {
    if (!token.value) return false;
    
    try {
      console.log('Refreshing auth token...');
      
      // Set the token in headers for the refresh request
      api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
      
      // Update token timestamp to avoid validation loops
      localStorage.setItem('token_last_refreshed', new Date().getTime());
      localStorage.setItem('auth_token', token.value);
      
      return true;
    } catch (err) {
      console.error('Token refresh failed:', err);
      return false;
    }
  }


const checkAuth = () => {
  const token = localStorage.getItem('auth_token');
  
  if (token) {
    isAuthenticated.value = true;
    
    // Try to restore user data if available
    const storedUserId = localStorage.getItem('user_id');
    const storedUserRole = localStorage.getItem('user_role');
    
    if (!user.value && storedUserId) {
      // If user data is missing but we have an ID, create minimal user object
      user.value = {
        id: storedUserId,
        role: storedUserRole || 'unknown'
      };
    }
    
    // Set authorization header for API calls
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
    console.log('Auth check: Authenticated with token');
    return true;
  } else {
    isAuthenticated.value = false;
    user.value = null;
    console.log('Auth check: Not authenticated (no token)');
    return false;
  }
};

  const normalizedRole = computed(() => {
  // If user has a role object with roleName property
  if (user.value?.role?.roleName) {
    return user.value.role.roleName.toLowerCase();
  }
  // If user has a direct role property
  else if (user.value?.role) {
    return user.value.role.toLowerCase();
  }
  // If the store tracks userRole directly
  else if (userRole.value) {
    return userRole.value.toLowerCase();
  }
  return '';
});


// Check if user has a specific role
const hasRole = (roleToCheck) => {
  if (!isAuthenticated.value) return false;
  
  const currentRole = normalizedRole.value;
  if (!currentRole) return false;
  
  // Check if the role matches (normalizing the input)
  return currentRole === roleToCheck.toLowerCase() ||
    currentRole === roleToCheck.toLowerCase().replace('_', '');
};
  
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
      assignedTenantId.value = null
      
      // Force a complete page reload to clear any Vue Router state
      window.location.href = window.location.origin + '/login'
    }
  }
  
  // Initialize auth state
  const initialize = async () => {
    console.log('Initializing auth store...');
    
    const storedToken = localStorage.getItem('auth_token')
    if (storedToken) {
      token.value = storedToken

      // Extract assigned tenant ID from token
      try {
        const tokenData = parseJwt(storedToken);
        if (tokenData.assigned_tenant_id) {
          assignedTenantId.value = tokenData.assigned_tenant_id;
          localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
          console.log('✅ Found assigned tenant ID in token:', tokenData.assigned_tenant_id);
          
          // Set tenant header for assigned_tenant_id
          api.defaults.headers.common['X-Tenant-ID'] = tokenData.assigned_tenant_id;
          console.log('✅ Set X-Tenant-ID header from token:', tokenData.assigned_tenant_id);
        }
      } catch (e) {
        console.error('Failed to parse token for assigned tenant ID:', e);
      }

      const storedUser = localStorage.getItem('user_data')
      if (storedUser && storedUser !== 'undefined') {
        try {
          user.value = JSON.parse(storedUser)
          api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`
          console.log('Auth initialized with stored user:', user.value)
          console.log('User role computed as:', userRole.value) 
          console.log('Set Authorization header with token')
          
          // ADDED: Refresh token timestamp
          localStorage.setItem('token_last_refreshed', new Date().getTime());
          
          // Set tenant header if applicable (fallback to stored value)
          if (assignedTenantId.value) {
            // Already set from token parsing above
          } else {
            const tenantId = localStorage.getItem('current_tenant_id');
            if (tenantId) {
              api.defaults.headers.common['X-Tenant-ID'] = tenantId;
              currentTenantId.value = tenantId;
              console.log('Set X-Tenant-ID header from localStorage:', tenantId);
            }
          }
          
          return true;
        } catch (e) {
          console.error('Failed to parse stored user_data:', e)
          return false;
        }
      } else {
        console.warn('Token exists but no valid user data found')
        return false;
      }
    }
    
    return false;
  }

  // Verify token is still valid
  const verifyToken = async () => {
    if (!token.value) return false;
    
    try {
      // Set the token in headers
      api.defaults.headers.common['Authorization'] = `Bearer ${token.value}`;
      
      // Get current timestamp
      const now = new Date().getTime();
      const lastRefreshed = parseInt(localStorage.getItem('token_last_refreshed') || '0');
      
      // Only verify if it's been more than 10 seconds since last check
      if (now - lastRefreshed > 10000) {
        // Update last refreshed time
        localStorage.setItem('token_last_refreshed', now.toString());
      }
      
      return true;
    } catch (err) {
      console.error('Token verification failed:', err);
      return false;
    }
  }

  // Login action - Clean and consistent with assigned tenant ID extraction
  const login = async (credentials) => {
    loading.value = true
    error.value = null
    
    try {
      console.log('🚀 Starting login request')
      
      // Reset all app state before login to ensure clean isolation
      resetAppState()
      
      console.log('📡 Performing backend login')
      
      const response = await api.post('/v1/auth/login', credentials)
      
      console.log('✅ Login response received:', response.data)
      
      // Extract data from response based on backend structure
      const data = response.data
      const accessToken = data.accessToken || data.token || data.access_token || data.jwt
      const userData = data.user || data.userData || data
      
      if (!accessToken) {
        throw new Error('No authentication token received from server')
      }
      
      console.log('👤 Raw user data from backend:', userData)
      
      // Store auth data
      token.value = accessToken
      user.value = userData
      localStorage.setItem('auth_token', accessToken)
      localStorage.setItem('user_data', JSON.stringify(userData))
      
      // Extract assigned tenant ID from token
      try {
        const tokenData = parseJwt(accessToken);
        console.log('🔑 Token payload:', tokenData);
        
        if (tokenData.assigned_tenant_id) {
          assignedTenantId.value = tokenData.assigned_tenant_id;
          localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
          console.log('✅ Found assigned tenant ID in token:', tokenData.assigned_tenant_id);
          
          // Set tenant header for API calls when using assigned tenant ID
          api.defaults.headers.common['X-Tenant-ID'] = tokenData.assigned_tenant_id;
          console.log('✅ Set X-Tenant-ID header from token:', tokenData.assigned_tenant_id);
        }
      } catch (e) {
        console.error('Failed to parse token for assigned tenant ID:', e);
      }
      
      // Debug role computation after user is set
      console.log('🎭 Computed user role after login:', userRole.value)
      console.log('🔍 Role computed from roleId:', userData.roleId)
      
      // Set token refresh timestamp
      localStorage.setItem('token_last_refreshed', new Date().getTime());
      
      // Set axios default header for authentication
      api.defaults.headers.common['Authorization'] = `Bearer ${accessToken}`
      console.log('🔑 Set Authorization header with token')
      
      // For tenant users, store and set tenant ID (only if not already set from token)
      if (!assignedTenantId.value && (userData.roleId === 2 || userData.role === 'templeadmin' || userData.role === 'tenant')) {
        const tenantId = userData.id
        currentTenantId.value = tenantId
        localStorage.setItem('current_tenant_id', tenantId)
        
        // Set tenant header for API calls
        api.defaults.headers.common['X-Tenant-ID'] = tenantId
        console.log('🏛️ Set X-Tenant-ID header from user data:', tenantId)
      }
      
      // Use the getDashboardPath function with assigned tenant ID awareness
      const redirectPath = getDashboardPath(userRole.value)
      
      console.log('🎯 Login successful, redirecting to:', redirectPath)
      console.log('👤 User role:', userRole.value)
      console.log('📋 Full user object:', userData)
      
      return {
        success: true,
        user: userData,
        redirectPath
      }
    } catch (err) {
      console.error('❌ Login error:', err)
      error.value = err.response?.data?.message || err.response?.data?.error || err.message || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
      console.log('✅ Login process completed, loading set to false')
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

  const forgotPassword = async (email) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.post('/v1/auth/forgot-password', { email })
      return {
        success: true,
        message: 'Password reset instructions sent to your email'
      }
    } catch (err) {
      console.error('Forgot password error:', err)
      error.value = err.response?.data?.message || err.message || 'Failed to send reset instructions'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  const resetPassword = async (token, newPassword) => {
    loading.value = true
    error.value = null
    
    try {
      const response = await api.post('/v1/auth/reset-password', { 
        token, 
        newPassword
      })
      
      return {
        success: true,
        message: 'Password has been reset successfully'
      }
    } catch (err) {
      console.error('Reset password error:', err)
      error.value = err.response?.data?.message || err.message || 'Failed to reset password'
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
      const redirectPath = getDashboardPath(userRole.value)
      
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
  
  // Select tenant (for admins, standard users, and monitoring users)
  const selectTenant = async (tenantId) => {
    if (!tenantId) {
      console.error('No tenant ID provided');
      return { success: false, error: 'No tenant ID provided' };
    }
    
    try {
      console.log('🏛️ Selecting tenant ID:', tenantId);
      
      // Store the tenant ID
      localStorage.setItem('selected_tenant_id', tenantId);
      localStorage.setItem('current_tenant_id', tenantId);
      localStorage.setItem('current_entity_id', tenantId);
      
      // Update store state
      currentTenantId.value = tenantId;
      
      // Set headers for API calls
      api.defaults.headers.common['X-Tenant-ID'] = tenantId;
      
      // Determine redirect path based on role
      const role = (userRole.value || '').toLowerCase();
      let redirectPath;
      
      if (role === 'superadmin' || role === 'super_admin' || 
          role === 'standard_user' || role === 'standarduser' || 
          role === 'monitoring_user' || role === 'monitoringuser') {
        redirectPath = `/entity/${tenantId}/dashboard`;
      } else {
        redirectPath = `/tenant/${tenantId}/dashboard`;
      }
      
      return {
        success: true,
        redirectPath
      };
    } catch (err) {
      console.error('Error selecting tenant:', err);
      return { 
        success: false, 
        error: 'Failed to select tenant' 
      };
    }
  }
  
  // Debug function to check if we can see the tenant data
  const debugTenantData = () => {
    console.log('------- DEBUG TENANT DATA -------')
    console.log('Current user:', user.value)
    console.log('Current tenant ID:', currentTenantId.value)
    console.log('Assigned tenant ID:', assignedTenantId.value)
    console.log('Is tenant?', isTenant.value)
    console.log('Dashboard route:', dashboardRoute.value)
    console.log('Auth headers:', api.defaults.headers.common)
    console.log('Local storage:', {
      auth_token: localStorage.getItem('auth_token'),
      user_data: localStorage.getItem('user_data'),
      current_tenant_id: localStorage.getItem('current_tenant_id'),
      selected_tenant_id: localStorage.getItem('selected_tenant_id'),
      current_entity_id: localStorage.getItem('current_entity_id'),
      assigned_tenant_id: localStorage.getItem('assigned_tenant_id')
    })
    console.log('Token payload:', parseJwt(localStorage.getItem('auth_token')))
    console.log('--------------------------------')
  }
  
  return {
    // State
    user,
    token,
    loading,
    error,
    currentTenantId,
    assignedTenantId,
    normalizedRole,
    hasRole,
    checkAuth,
    
    // Getters
    isAuthenticated,
    userRole,
    isTenant,
    isDevotee,
    isVolunteer,
    isSuperAdmin,
    isStandardUser,
    isMonitoringUser,
    isEndUser,
    dashboardRoute,
    needsTenantSelection,
    
    // Actions
    initialize,
    login,
    logout,
    register,
    joinTemple,
    selectTenant,
    getDashboardPath,
    resetAppState,
    clearAllStorage,
    debugTenantData,
    forgotPassword,
    resetPassword,
    refreshToken,
    verifyToken,
    parseJwt,
    
    // For backwards compatibility
    initializeAuth: initialize,
    clearError: () => { error.value = null }
  }
})