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

  // Status checks
  const isPending = computed(() => userStatus.value === USER_STATUS.PENDING)
  const isApproved = computed(() => userStatus.value === USER_STATUS.APPROVED)
  const isRejected = computed(() => userStatus.value === USER_STATUS.REJECTED)

  // Get role configuration
  const getRoleConfig = (role) => {
    return rolesConfig.roles[role] || null
  }

  // Get dashboard path based on role and entity
  const getDashboardPath = () => {
    if (!user.value) return '/login'
    
    const roleConfig = getRoleConfig(userRole.value)
    if (!roleConfig) return '/login'

    // Replace :id with actual entity ID for end users
    if (isEndUser.value && currentEntity.value) {
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

  // Authentication actions
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
      
      showToast('Login successful!', 'success')
      
      // Redirect to appropriate dashboard
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
      localStorage.removeItem('auth_token')
      localStorage.removeItem('user_data')
      localStorage.removeItem('current_tenant_id')
      
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

  // Initialize auth state from localStorage
  const initializeAuth = async () => {
    try {
      const storedToken = localStorage.getItem('auth_token')
      const storedUser = localStorage.getItem('user_data')
      
      if (storedToken && storedUser) {
        token.value = storedToken
        user.value = JSON.parse(storedUser)
        
        // Verify token is still valid
        try {
          const response = await apiClient.auth.refreshToken()
          if (response.data.token) {
            token.value = response.data.token
            localStorage.setItem('auth_token', response.data.token)
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
    getRoleConfig
  }
}