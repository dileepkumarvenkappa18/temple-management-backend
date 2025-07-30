// src/services/auth.service.js
import axios from 'axios'

const API_BASE_URL = process.env.VUE_APP_API_BASE_URL || 'http://localhost:8080/api'

class AuthService {
  constructor() {
    this.token = localStorage.getItem('token')
    this.user = JSON.parse(localStorage.getItem('user') || 'null')
    
    // Setup axios interceptor for auth token
    this.setupAxiosInterceptor()
  }

  setupAxiosInterceptor() {
    axios.interceptors.request.use(
      (config) => {
        if (this.token) {
          config.headers.Authorization = `Bearer ${this.token}`
        }
        return config
      },
      (error) => {
        return Promise.reject(error)
      }
    )

    // Response interceptor to handle 401 errors
    axios.interceptors.response.use(
      (response) => response,
      (error) => {
        if (error.response?.status === 401) {
          this.logout()
          window.location.href = '/login'
        }
        return Promise.reject(error)
      }
    )
  }

  /**
   * Register a new user (Tenant, Devotee, or Volunteer only)
   * Super Admin is not registerable via public flow
   */
  async register(userData) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/register`, {
        fullName: userData.fullName,
        email: userData.email,
        password: userData.password,
        phone: userData.phone,
        role: userData.role // tenant, devotee, volunteer
      })

      // Handle different registration flows based on role
      if (userData.role === 'tenant') {
        // Tenant needs approval - show pending message
        return {
          success: true,
          message: 'Registration received. You\'ll be notified after approval.',
          needsApproval: true,
          user: response.data.user
        }
      } else {
        // End users can immediately proceed
        return {
          success: true,
          message: 'Welcome! You can now join a temple and complete your profile.',
          needsApproval: false,
          user: response.data.user,
          token: response.data.token
        }
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Login user - handles all roles including Super Admin
   */
  async login(credentials) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/login`, {
        email: credentials.email,
        password: credentials.password
      })

      const { user, token } = response.data

      // Store auth data
      this.token = token
      this.user = user
      localStorage.setItem('token', token)
      localStorage.setItem('user', JSON.stringify(user))

      // Return role-based routing information
      return {
        success: true,
        user,
        token,
        redirectPath: this.getRedirectPath(user)
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Get redirect path based on user role and status
   */
  getRedirectPath(user) {
    switch (user.role) {
      case 'superadmin':
        return '/superadmin/dashboard'
      
      case 'tenant':
        // Check if tenant is approved
        if (user.status === 'pending') {
          return '/tenant/pending'
        }
        return '/tenant/dashboard'
      
      case 'devotee':
        // Check if devotee has joined a temple
        if (!user.entityId) {
          return '/devotee/temple-selection'
        }
        return `/entity/${user.entityId}/devotee/dashboard`
      
      case 'volunteer':
        // Check if volunteer has joined a temple
        if (!user.entityId) {
          return '/volunteer/temple-selection'
        }
        return `/entity/${user.entityId}/volunteer/dashboard`
      
      default:
        return '/dashboard'
    }
  }

  /**
   * Logout user
   */
  async logout() {
    try {
      if (this.token) {
        await axios.post(`${API_BASE_URL}/auth/logout`)
      }
    } catch (error) {
      console.error('Logout error:', error)
    } finally {
      // Clear local storage and state
      this.token = null
      this.user = null
      localStorage.removeItem('token')
      localStorage.removeItem('user')
    }
  }

  /**
   * Refresh token
   */
  async refreshToken() {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/refresh`)
      const { token } = response.data

      this.token = token
      localStorage.setItem('token', token)

      return { success: true, token }
    } catch (error) {
      this.logout()
      throw this.handleError(error)
    }
  }

  /**
   * Forgot password
   */
  async forgotPassword(email) {
    try {
      await axios.post(`${API_BASE_URL}/auth/forgot-password`, { email })
      return {
        success: true,
        message: 'Password reset link sent to your email'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Reset password
   */
  async resetPassword(token, newPassword) {
    try {
      await axios.post(`${API_BASE_URL}/auth/reset-password`, {
        token,
        password: newPassword
      })
      return {
        success: true,
        message: 'Password reset successfully'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Change password (authenticated user)
   */
  async changePassword(currentPassword, newPassword) {
    try {
      await axios.post(`${API_BASE_URL}/auth/change-password`, {
        currentPassword,
        newPassword
      })
      return {
        success: true,
        message: 'Password changed successfully'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Get current user profile
   */
  async getCurrentUser() {
    try {
      const response = await axios.get(`${API_BASE_URL}/auth/me`)
      const user = response.data.user

      // Update stored user data
      this.user = user
      localStorage.setItem('user', JSON.stringify(user))

      return { success: true, user }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Update user profile
   */
  async updateProfile(profileData) {
    try {
      const response = await axios.put(`${API_BASE_URL}/auth/profile`, profileData)
      const user = response.data.user

      // Update stored user data
      this.user = user
      localStorage.setItem('user', JSON.stringify(user))

      return { success: true, user }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Verify email
   */
  async verifyEmail(token) {
    try {
      await axios.post(`${API_BASE_URL}/auth/verify-email`, { token })
      return {
        success: true,
        message: 'Email verified successfully'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Resend email verification
   */
  async resendEmailVerification() {
    try {
      await axios.post(`${API_BASE_URL}/auth/resend-verification`)
      return {
        success: true,
        message: 'Verification email sent'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Join temple (for devotees and volunteers)
   */
  async joinTemple(templeId) {
    try {
      const response = await axios.post(`${API_BASE_URL}/auth/join-temple`, {
        templeId
      })

      const user = response.data.user
      
      // Update stored user data with new entityId
      this.user = user
      localStorage.setItem('user', JSON.stringify(user))

      return {
        success: true,
        user,
        message: 'Successfully joined temple!'
      }
    } catch (error) {
      throw this.handleError(error)
    }
  }

  /**
   * Check if user is authenticated
   */
  isAuthenticated() {
    return !!(this.token && this.user)
  }

  /**
   * Check if user has specific role
   */
  hasRole(role) {
    return this.user?.role === role
  }

  /**
   * Check if user has any of the specified roles
   */
  hasAnyRole(roles) {
    return roles.includes(this.user?.role)
  }

  /**
   * Check if user has specific permission
   */
  hasPermission(permission) {
    return this.user?.permissions?.includes(permission)
  }

  /**
   * Get user's entity ID (temple ID for devotees/volunteers)
   */
  getEntityId() {
    return this.user?.entityId
  }

  /**
   * Check if tenant is approved
   */
  isTenantApproved() {
    return this.user?.role === 'tenant' && this.user?.status === 'approved'
  }

  /**
   * Check if user needs to complete profile
   */
  needsProfileCompletion() {
    return this.user?.role === 'devotee' && !this.user?.profileCompleted
  }

  /**
   * Check if user needs to join temple
   */
  needsTempleSelection() {
    const rolesNeedingTemple = ['devotee', 'volunteer']
    return rolesNeedingTemple.includes(this.user?.role) && !this.user?.entityId
  }

  /**
   * Get current user data
   */
  getCurrentUserData() {
    return this.user
  }

  /**
   * Get auth token
   */
  getToken() {
    return this.token
  }

  /**
   * Handle API errors consistently
   */
  handleError(error) {
    if (error.response?.data?.message) {
      return new Error(error.response.data.message)
    } else if (error.message) {
      return new Error(error.message)
    } else {
      return new Error('An unexpected error occurred')
    }
  }

  /**
   * Validate password strength
   */
  validatePasswordStrength(password) {
    const criteria = {
      minLength: password.length >= 8,
      hasUppercase: /[A-Z]/.test(password),
      hasLowercase: /[a-z]/.test(password),
      hasNumber: /\d/.test(password),
      hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password)
    }

    const score = Object.values(criteria).filter(Boolean).length
    
    return {
      criteria,
      score,
      strength: score < 3 ? 'weak' : score < 5 ? 'medium' : 'strong'
    }
  }

  /**
   * Format user display name
   */
  getUserDisplayName() {
    if (!this.user) return ''
    return this.user.fullName || this.user.email || 'User'
  }

  /**
   * Get role display name
   */
  getRoleDisplayName() {
    const roleMap = {
      superadmin: 'Super Admin',
      tenant: 'Temple Admin',
      devotee: 'Devotee',
      volunteer: 'Volunteer'
    }
    return roleMap[this.user?.role] || 'User'
  }
}

// Create and export singleton instance
const authService = new AuthService()
export default authService