import { ref, computed, watch } from 'vue'
import { apiClient } from '@/plugins/axios'
import { useToast } from './useToast'
import { USER_ROLES, USER_STATUS } from '@/utils/constants'
import rolesConfig from '@/data/constants/roles.json'

// Global auth state (shared across all components)
const user = ref(null)
const token = ref(localStorage.getItem('auth_token'))
const isLoading = ref(false)
const isInitialized = ref(false)
const assignedTenantId = ref(localStorage.getItem('assigned_tenant_id'))

export function useAuth() {
  // Computed properties
  const isAuthenticated = computed(() => !!token.value && !!user.value)
  const userRole = computed(() => user.value?.role || null)
  const userStatus = computed(() => user.value?.status || null)
  const currentEntity = computed(() => user.value?.current_entity || null)

  // Role-based computed properties
  const isSuperAdmin = computed(() => userRole.value === USER_ROLES.SUPER_ADMIN)
  const isTenant = computed(() => userRole.value === USER_ROLES.TENANT)
  const isDevotee = computed(() => userRole.value === USER_ROLES.DEVOTEE)
  const isVolunteer = computed(() => userRole.value === USER_ROLES.VOLUNTEER)
  const isEndUser = computed(() => isDevotee.value || isVolunteer.value)
  
  // Special role checks
  const isStandardUser = computed(() => {
    const role = (userRole.value || '').toLowerCase();
    return role === 'standard_user' || role === 'standarduser';
  })
  
  const isMonitoringUser = computed(() => {
    const role = (userRole.value || '').toLowerCase();
    return role === 'monitoring_user' || role === 'monitoringuser';
  })

  // Status checks
  const isPending = computed(() => userStatus.value === USER_STATUS.PENDING)
  const isApproved = computed(() => userStatus.value === USER_STATUS.APPROVED)
  const isRejected = computed(() => userStatus.value === USER_STATUS.REJECTED)

  // Function to parse JWT token
  const parseJwt = (token) => {
    try {
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

  // Get role configuration
  const getRoleConfig = (role) => {
    return rolesConfig.roles[role] || null
  }

  // Get dashboard path based on role, entity, and assigned tenant ID
  // UPDATED: Changed redirection for monitoring_user and standard_user to tenant dashboard
  const getDashboardPath = () => {
    if (!user.value) return '/login'
    
    const roleConfig = getRoleConfig(userRole.value)
    if (!roleConfig) return '/login'
    
    // Special handling for monitoring_user and standard_user
    if (isStandardUser.value || isMonitoringUser.value) {
      // If assigned tenant ID is available, redirect to that tenant's dashboard
      if (assignedTenantId.value) {
        console.log(`✅ Redirecting ${userRole.value} to assigned tenant: ${assignedTenantId.value}`);
        // UPDATED: Redirect to tenant dashboard instead of entity dashboard
        return `/tenant/${assignedTenantId.value}/dashboard`;
      }
      
      // Otherwise, redirect to tenant selection
      return '/tenant-selection';
    }

    // For devotee role, always redirect to temple selection first
    if (isDevotee.value) {
      return '/devotee/temple-selection';
    }

    // Replace :id with actual entity ID for other end users
    if (isVolunteer.value && currentEntity.value) {
      return roleConfig.dashboard.replace(':id', currentEntity.value.id)
    }
    
    return roleConfig.dashboard
  }

  // Check if user has specific permission
  const hasPermission = (permission) => {
    if (!user.value) return false
    
    const roleConfig = getRoleConfig(userRole.value)
    return roleConfig?.permissions?.includes(permission) || false
  }

  // Check if user can access a specific route
  const canAccessRoute = (routeName, requiredRole = null, requiredPermissions = []) => {
    if (!isAuthenticated.value) return false
    
    // Check role requirement
    if (requiredRole && userRole.value !== requiredRole) return false
    
    // Check permission requirements
    if (requiredPermissions.length > 0) {
      return requiredPermissions.every(permission => hasPermission(permission))
    }
    
    return true
  }

  // Authentication actions with token parsing
  const login = async (credentials) => {
    const { showToast } = useToast()
    
    try {
      isLoading.value = true
      const response = await apiClient.auth.login(credentials)
      
      const { user: userData, token: authToken } = response.data
      
      // Store authentication data
      token.value = authToken
      user.value = userData
      localStorage.setItem('auth_token', authToken)
      localStorage.setItem('user_data', JSON.stringify(userData))
      
      // Extract assigned_tenant_id from token if present
      try {
        const tokenData = parseJwt(authToken);
        if (tokenData.assigned_tenant_id) {
          assignedTenantId.value = tokenData.assigned_tenant_id;
          localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
          console.log('Found assigned tenant ID in token:', tokenData.assigned_tenant_id);
          
          // Set tenant header for API calls
          apiClient.setTenantId(tokenData.assigned_tenant_id);
        }
      } catch (e) {
        console.error('Failed to parse token for assigned tenant ID:', e);
      }
      
      showToast('Login successful!', 'success')
      
      // Get dashboard path with assigned tenant ID awareness
      const dashboardPath = getDashboardPath()
      
      return { 
        success: true, 
        user: userData,
        redirectPath: dashboardPath
      }
    } catch (error) {
      console.error('Login failed:', error)
      return { success: false, error: error.response?.data?.message || 'Login failed' }
    } finally {
      isLoading.value = false
    }
  }

  const register = async (userData) => {
    const { showToast } = useToast()
    
    try {
      isLoading.value = true
      const response = await apiClient.auth.register(userData)
      
      const { user: newUser, message } = response.data
      
      // Handle different registration flows based on role
      const roleConfig = rolesConfig.registrationFlow[userData.role]
      
      if (roleConfig.approvalRequired) {
        // For TENANT - redirect to pending approval page
        showToast(message || 'Registration successful! You will be notified after approval.', 'success')
        return { 
          success: true, 
          user: newUser,
          redirectPath: '/auth/pending-approval'
        }
      } else {
        // For DEVOTEE/VOLUNTEER - proceed to temple selection
        showToast(message || 'Registration successful! Please join a temple to continue.', 'success')
        
        // Auto-login for end users
        if (response.data.token) {
          token.value = response.data.token
          user.value = newUser
          localStorage.setItem('auth_token', response.data.token)
          localStorage.setItem('user_data', JSON.stringify(newUser))
          
          // Check for assigned tenant ID in token
          try {
            const tokenData = parseJwt(response.data.token);
            if (tokenData.assigned_tenant_id) {
              assignedTenantId.value = tokenData.assigned_tenant_id;
              localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
              console.log('Found assigned tenant ID in token:', tokenData.assigned_tenant_id);
              
              // Set tenant header for API calls
              apiClient.setTenantId(tokenData.assigned_tenant_id);
            }
          } catch (e) {
            console.error('Failed to parse token for assigned tenant ID:', e);
          }
          
          // Redirect based on role
          return { 
            success: true, 
            user: newUser,
            redirectPath: roleConfig.redirectAfterRegistration
          }
        } else {
          return { 
            success: true, 
            user: newUser,
            redirectPath: '/login'
          }
        }
      }
    } catch (error) {
      console.error('Registration failed:', error)
      return { success: false, error: error.response?.data?.message || 'Registration failed' }
    } finally {
      isLoading.value = false
    }
  }

  const logout = async () => {
    const { showToast } = useToast()
    
    try {
      isLoading.value = true
      
      // Call logout endpoint if user is authenticated
      if (token.value) {
        await apiClient.auth.logout()
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // Clear local state regardless of API call success
      token.value = null
      user.value = null
      assignedTenantId.value = null
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_data')
      localStorage.removeItem('current_tenant_id')
      localStorage.removeItem('assigned_tenant_id')
      
      isLoading.value = false
      
      showToast('Logged out successfully', 'success')
      return { success: true, redirectPath: '/login' }
    }
  }

  const forgotPassword = async (email) => {
    const { showToast } = useToast()
    
    try {
      isLoading.value = true
      await apiClient.auth.forgotPassword(email)
      
      showToast('Password reset link sent to your email', 'success')
      return { success: true }
    } catch (error) {
      console.error('Forgot password failed:', error)
      return { success: false, error: error.response?.data?.message || 'Failed to send reset link' }
    } finally {
      isLoading.value = false
    }
  }

  const resetPassword = async (token, password) => {
    const { showToast } = useToast()
      
    try {
      isLoading.value = true
      await apiClient.auth.resetPassword(token, password)
      
      showToast('Password reset successful! Please login with your new password.', 'success')
      return { success: true, redirectPath: '/login' }
    } catch (error) {
      console.error('Reset password failed:', error)
      return { success: false, error: error.response?.data?.message || 'Password reset failed' }
    } finally {
      isLoading.value = false
    }
  }

  // Initialize auth state from localStorage with token parsing
  const initializeAuth = async () => {
    try {
      const storedToken = localStorage.getItem('auth_token')
      const storedUser = localStorage.getItem('user_data')
      
      if (storedToken && storedUser) {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
        
        // Extract assigned_tenant_id from token if present
        try {
          const tokenData = parseJwt(storedToken);
          if (tokenData.assigned_tenant_id) {
            assignedTenantId.value = tokenData.assigned_tenant_id;
            localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
            console.log('Found assigned tenant ID in token:', tokenData.assigned_tenant_id);
            
            // Set tenant header for API calls
            apiClient.setTenantId(tokenData.assigned_tenant_id);
          }
        } catch (e) {
          console.error('Failed to parse token for assigned tenant ID:', e);
        }
        
        // Verify token is still valid
        try {
          const response = await apiClient.auth.refreshToken()
          if (response.data.token) {
            token.value = response.data.token
            localStorage.setItem('auth_token', response.data.token)
            
            // Re-check assigned tenant ID in new token
            try {
              const tokenData = parseJwt(response.data.token);
              if (tokenData.assigned_tenant_id) {
                assignedTenantId.value = tokenData.assigned_tenant_id;
                localStorage.setItem('assigned_tenant_id', tokenData.assigned_tenant_id);
                console.log('Found assigned tenant ID in refreshed token:', tokenData.assigned_tenant_id);
                
                // Set tenant header for API calls
                apiClient.setTenantId(tokenData.assigned_tenant_id);
              }
            } catch (e) {
              console.error('Failed to parse refreshed token for assigned tenant ID:', e);
            }
          }
        } catch (error) {
          // Token is invalid, clear auth state
          console.warn('Token refresh failed, clearing auth state')
          await logout()
        }
      }
    } catch (error) {
      console.error('Auth initialization failed:', error)
    } finally {
      isInitialized.value = true
    }
  }

  // Update user profile
  const updateProfile = async (profileData) => {
    const { showToast } = useToast()
      
    try {
      isLoading.value = true
      
      // Call appropriate API based on role
      let response
      if (isTenant.value) {
        response = await apiClient.tenant.updateProfile(profileData)
      } else if (isDevotee.value) {
        response = await apiClient.devotee.updateProfile(profileData)
      } else {
        throw new Error('Profile update not supported for this role')
      }
      
      // Update local user data
      user.value = { ...user.value, ...response.data.user }
      localStorage.setItem('user_data', JSON.stringify(user.value))
      
      showToast('Profile updated successfully!', 'success')
      return { success: true, user: user.value }
    } catch (error) {
      console.error('Profile update failed:', error)
      return { success: false, error: error.response?.data?.message || 'Profile update failed' }
    } finally {
      isLoading.value = false
    }
  }

  // Join temple (for end users)
  const joinTemple = async (templeId) => {
    const { showToast } = useToast()
      
    try {
      isLoading.value = true
      
      let response
      if (isDevotee.value) {
        response = await apiClient.devotee.joinTemple(templeId)
      } else if (isVolunteer.value) {
        response = await apiClient.volunteer.joinTemple(templeId)
      } else {
        throw new Error('Only devotees and volunteers can join temples')
      }
      
      // Update user's current entity with the entity_id from the response
      const entityId = response.data.entity_id || templeId
      const updatedUser = { 
        ...user.value, 
        current_entity: { 
          id: entityId,
          ...response.data.temple 
        }
      }
      
      user.value = updatedUser
      localStorage.setItem('user_data', JSON.stringify(updatedUser))
      
      showToast(`Successfully joined ${response.data.temple?.name || 'temple'}!`, 'success')
      
      // Get dashboard path with proper entity ID
      const dashboardPath = getDashboardPath()
      
      return { 
        success: true, 
        temple: response.data.temple,
        entityId: entityId,
        redirectPath: dashboardPath
      }
    } catch (error) {
      console.error('Join temple failed:', error)
      return { success: false, error: error.response?.data?.message || 'Failed to join temple' }
    } finally {
      isLoading.value = false
    }
  }

  // Watch for token changes to update axios headers
  watch(token, (newToken) => {
    if (newToken) {
      // Token is set in axios interceptor
    } else {
      // Token cleared, user logged out
    }
  }, { immediate: true })

  return {
    // State
    user: computed(() => user.value),
    token: computed(() => token.value),
    isLoading: computed(() => isLoading.value),
    isInitialized: computed(() => isInitialized.value),
    assignedTenantId: computed(() => assignedTenantId.value),
    
    // Computed properties
    isAuthenticated,
    userRole,
    userStatus,
    currentEntity,
    isSuperAdmin,
    isTenant,
    isDevotee,
    isVolunteer,
    isEndUser,
    isStandardUser,
    isMonitoringUser,
    isPending,
    isApproved,
    isRejected,
    
    // Methods
    login,
    register,
    logout,
    forgotPassword,
    resetPassword,
    initializeAuth,
    updateProfile,
    joinTemple,
    hasPermission,
    canAccessRoute,
    getDashboardPath,
    getRoleConfig,
    parseJwt
  }
}