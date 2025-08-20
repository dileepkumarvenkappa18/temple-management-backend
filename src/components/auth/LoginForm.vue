<template>
  <div class="login-form">
    <!-- Alert for errors -->
    <div v-if="error" class="bg-red-100 border border-red-400 text-red-700 px-4 py-3 rounded relative mb-4">
      <span class="block sm:inline">{{ error }}</span>
      <button 
        @click="error = null" 
        class="absolute top-0 right-0 p-2"
        type="button"
      >
        <span class="text-red-700">&times;</span>
      </button>
    </div>
    
    <form @submit.prevent="handleLogin" class="space-y-6">
      <!-- Email Input -->
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700">Email</label>
        <div class="mt-1">
          <input
            id="email"
            v-model="form.email"
            type="email"
            required
            autocomplete="email"
            class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Enter your email"
            :disabled="isLoading"
          />
        </div>
      </div>
      
      <!-- Password Input -->
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700">Password</label>
        <div class="mt-1">
          <input
            id="password"
            v-model="form.password"
            type="password"
            required
            autocomplete="current-password"
            class="appearance-none block w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
            placeholder="Enter your password"
            :disabled="isLoading"
          />
        </div>
      </div>
      
      <!-- Remember Me Checkbox -->
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <input
            id="remember_me"
            v-model="rememberMe"
            type="checkbox"
            class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            :disabled="isLoading"
          />
          <label for="remember_me" class="ml-2 block text-sm text-gray-900">Remember me</label>
        </div>
        
        <div class="text-sm">
          <router-link to="/forgot-password" class="font-medium text-indigo-600 hover:text-indigo-500">
            Forgot your password?
          </router-link>
        </div>
      </div>
      
      <!-- Submit Button -->
      <div>
        <button
          type="submit"
          :disabled="isLoading"
          class="w-full flex justify-center py-2 px-4 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <template v-if="isLoading">
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          </template>
          <template v-else>
            Sign in
          </template>
        </button>
      </div>
      
      <!-- Register Link -->
      <div class="text-center">
        <router-link to="/register" class="font-medium text-indigo-600 hover:text-indigo-500">
          Don't have an account? Register
        </router-link>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import axios from 'axios'
import api from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth' // Added import for auth store

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore() // Added auth store

// Form state
const form = ref({
  email: '',
  password: ''
})
const rememberMe = ref(false)
const isLoading = ref(false)
const error = ref(null)

// Get redirect path from URL if available
const redirectPath = computed(() => route.query.redirect || '')

const handleLogin = async () => {
  if (!form.value.email || !form.value.password) {
    error.value = 'Please enter your email and password'
    return
  }
  
  isLoading.value = true
  error.value = null
  
  try {
    // Use auth store login instead of direct API call
    const result = await authStore.login({
      email: form.value.email,
      password: form.value.password
    })
    
    if (result.success) {
      // CRITICAL FIX: Check for standard_user and monitoring_user roles
      const userData = result.user;
      const userRole = (userData.role || '').toLowerCase();
      
      // Log the user role for debugging
      console.log('User logged in with role:', userRole);
      
      if (userRole === 'standard_user' || userRole === 'standarduser' || 
          userRole === 'monitoring_user' || userRole === 'monitoringuser') {
        console.log('CRITICAL FIX: Redirecting standard/monitoring user to tenant selection');
        
        // Force redirect to tenant selection
        window.location.href = `${window.location.origin}/tenant-selection`;
        return;
      }
      
      // Use router to navigate to the dashboard path after successful login for other roles
      console.log('Login successful, navigating to:', result.redirectPath)
      router.push(redirectPath.value || result.redirectPath)
    } else {
      error.value = result.error || 'Login failed'
    }
  } catch (err) {
    console.error('Login error:', err)
    
    // Handle different error scenarios
    if (err.response) {
      if (err.response.status === 401) {
        error.value = 'Invalid email or password'
      } else if (err.response.status === 403) {
        error.value = 'Your account is locked or inactive'
      } else if (err.response.status >= 500) {
        error.value = 'Server error. Please try again later'
      } else {
        error.value = err.response.data?.message || err.response.data?.error || 'Login failed'
      }
    } else if (err.request) {
      error.value = 'Cannot connect to the server. Please check your internet connection'
    } else {
      error.value = err.message || 'An unexpected error occurred'
    }
  } finally {
    isLoading.value = false
  }
}

// Helper function to extract auth data from response
const extractAuthData = (data) => {
  // Handle different response formats from different backend implementations
  return {
    token: data.token || data.accessToken || data.access_token || data.jwt,
    user: data.user || data.userData || data
  }
}

// Helper function to determine dashboard path based on role
const getDashboardPath = (user) => {
  // Default to basic dashboard
  let path = '/'
  
  // Get role - handle different formats
  const role = user.role || (user.roleId ? mapRoleIdToString(user.roleId) : null)
  const entityId = user.entityId || user.current_entity?.id
  
  // Route based on role
  if (role === 'superadmin' || role === 'super_admin') {
    path = '/superadmin/dashboard'
  } else if (role === 'templeadmin' || role === 'tenant') {
    path = '/tenant/dashboard'
  } else if (role === 'devotee') {
    path = entityId ? `/entity/${entityId}/devotee/dashboard` : '/devotee/temple-selection'
  } else if (role === 'volunteer') {
    path = entityId ? `/entity/${entityId}/volunteer/dashboard` : '/volunteer/temple-selection'
  } else if (role === 'standard_user' || role === 'standarduser' || 
             role === 'monitoring_user' || role === 'monitoringuser') {
    path = '/tenant-selection'
  }
  
  return path
}

// Helper function to map role IDs to string roles
const mapRoleIdToString = (roleId) => {
  const ROLE_MAP = {
    1: 'superadmin',
    2: 'templeadmin',
    3: 'devotee',
    4: 'volunteer'
  }
  return ROLE_MAP[roleId] || null
}
</script>

<style scoped>
.login-form {
  width: 100%;
}
</style>