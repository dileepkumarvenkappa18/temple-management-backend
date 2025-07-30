<template>
  <BaseCard class="w-full max-w-md mx-auto">
    <div class="text-center mb-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">Create Account</h2>
      <p class="text-gray-600">Join our temple management platform</p>
    </div>

    <!-- Registration Success Alert (Non-modal version) -->
    <div v-if="registrationSuccess" class="mb-6 bg-green-50 border border-green-200 rounded-lg p-4">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <svg class="h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
        </div>
        <div class="ml-3">
          <h3 class="text-sm font-medium text-green-800">Registration Successful!</h3>
          <div class="mt-1 text-sm text-green-700">
            <p v-if="needsApproval">
              Your temple admin account has been created. You'll receive an email once your account is approved by our team.
            </p>
            <p v-else>
              Your account has been created successfully. You can now login to access your dashboard.
            </p>
          </div>
          <div class="mt-3">
            <BaseButton @click="goToLogin" variant="success" size="sm">
              Go to Login
            </BaseButton>
          </div>
        </div>
      </div>
    </div>

    <form v-if="!registrationSuccess" @submit.prevent="handleRegister" class="space-y-6">
      <!-- Full Name Field -->
      <div>
        <label for="fullName" class="block text-sm font-medium text-gray-700 mb-2">
          Full Name
        </label>
        <BaseInput
          id="fullName"
          v-model="form.fullName"
          type="text"
          placeholder="Enter your full name"
          :error="errors.fullName"
          required
          autocomplete="name"
        />
      </div>

      <!-- Email Field -->
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
          Email Address
        </label>
        <BaseInput
          id="email"
          v-model="form.email"
          type="email"
          placeholder="Enter your email"
          :error="errors.email"
          required
          autocomplete="email"
        />
      </div>

      <!-- Password Field -->
      <div>
        <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
          Password
        </label>
        <BaseInput
          id="password"
          v-model="form.password"
          type="password"
          placeholder="Create a strong password"
          :error="errors.password"
          required
          autocomplete="new-password"
        />
        <PasswordStrengthMeter 
          :password="form.password" 
          class="mt-2"
        />
      </div>

      <!-- Phone Field (Required) -->
      <div>
        <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
          Phone Number <span class="text-red-500">*</span>
        </label>
        <BaseInput
          id="phone"
          v-model="form.phone"
          type="tel"
          placeholder="Enter your phone number"
          :error="errors.phone"
          required
          autocomplete="tel"
        />
      </div>

      <!-- Role Selector -->
      <div>
        <label for="role" class="block text-sm font-medium text-gray-700 mb-2">
          I want to join as
        </label>
        <RoleSelector
          v-model="form.role"
          :error="errors.role"
          required
        />
      </div>

      <!-- Terms and Privacy -->
      <div class="flex items-start">
        <input
          id="terms"
          v-model="form.acceptTerms"
          type="checkbox"
          class="mt-1 h-4 w-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500 focus:ring-2"
        />
        <label for="terms" class="ml-3 text-sm text-gray-700">
          I agree to the 
          <router-link to="/terms" class="text-indigo-600 hover:text-indigo-800 font-medium">
            Terms of Service
          </router-link>
          and 
          <router-link to="/privacy" class="text-indigo-600 hover:text-indigo-800 font-medium">
            Privacy Policy
          </router-link>
        </label>
      </div>
      <div v-if="errors.acceptTerms" class="text-red-600 text-sm">
        {{ errors.acceptTerms }}
      </div>

      <!-- Submit Button -->
      <BaseButton
        type="submit"
        variant="primary"
        size="lg"
        :loading="isLoading"
        :disabled="!isFormValid"
        class="w-full"
      >
        Create Account
      </BaseButton>

      <!-- Login Link -->
      <div class="text-center pt-4 border-t border-gray-200">
        <p class="text-sm text-gray-600">
          Already have an account?
          <router-link 
            to="/login" 
            class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
          >
            Sign In
          </router-link>
        </p>
      </div>
    </form>

    <!-- Registration Success Modal (alternative approach) -->
    <BaseModal v-if="showSuccessModal" @close="goToLogin">
      <template #header>
        <h3 class="text-lg font-medium text-gray-900">Registration Successful!</h3>
      </template>
      <template #default>
        <div class="py-4">
          <p v-if="needsApproval" class="text-gray-600">
            Your temple admin account has been created. You'll receive an email once your account is approved by our team.
          </p>
          <p v-else class="text-gray-600">
            Your account has been created successfully. You can now login to access your dashboard.
          </p>
        </div>
      </template>
      <template #footer>
        <BaseButton @click="goToLogin" variant="primary" class="w-full">
          Go to Login
        </BaseButton>
      </template>
    </BaseModal>
  </BaseCard>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import { apiClient } from '@/plugins/axios' // Import apiClient for direct API calls
import BaseCard from '@/components/common/BaseCard.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import RoleSelector from '@/components/auth/RoleSelector.vue'
import PasswordStrengthMeter from '@/components/auth/PasswordStrengthMeter.vue'

const router = useRouter()
const authStore = useAuthStore()
const { success, error } = useToast()

const isLoading = ref(false)
const showSuccessModal = ref(false)
const needsApproval = ref(false)
// New flag to show in-page success message
const registrationSuccess = ref(false)

const form = ref({
  fullName: '',
  email: '',
  password: '',
  phone: '',
  role: '',
  acceptTerms: false
})

const errors = ref({})

const isFormValid = computed(() => {
  return form.value.fullName && 
         form.value.email && 
         form.value.password && 
         form.value.phone && // Phone is now required
         form.value.role &&
         form.value.acceptTerms &&
         Object.keys(errors.value).length === 0
})

const validateForm = () => {
  errors.value = {}
  
  // Full name validation
  if (!form.value.fullName.trim()) {
    errors.value.fullName = 'Full name is required'
  } else if (form.value.fullName.trim().length < 2) {
    errors.value.fullName = 'Full name must be at least 2 characters'
  }
  
  // Email validation
  if (!form.value.email) {
    errors.value.email = 'Email is required'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    errors.value.email = 'Please enter a valid email address'
  } else if (!form.value.email.endsWith('@gmail.com')) {
    errors.value.email = 'Only @gmail.com email addresses are allowed'
  }
  
  // Password validation
  if (!form.value.password) {
    errors.value.password = 'Password is required'
  } else if (form.value.password.length < 8) {
    errors.value.password = 'Password must be at least 8 characters'
  }
  
  // Phone validation (required)
  if (!form.value.phone) {
    errors.value.phone = 'Phone number is required'
  } else if (!/^[\+]?[1-9][\d]{0,15}$/.test(form.value.phone.replace(/\s/g, ''))) {
    errors.value.phone = 'Please enter a valid phone number'
  }
  
  // Role validation
  if (!form.value.role) {
    errors.value.role = 'Please select your role'
  }
  
  // Terms validation
  if (!form.value.acceptTerms) {
    errors.value.acceptTerms = 'You must accept the terms and privacy policy'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleRegister = async () => {
  if (!validateForm()) return
  
  isLoading.value = true
  
  try {
    // Map frontend role to backend role if needed
    // The backend expects 'templeadmin' for tenant role
    const roleMapping = {
      'tenant': 'templeadmin',
      'devotee': 'devotee', 
      'volunteer': 'volunteer'
    }

    // Attempt to register with the API
    try {
      // Use apiClient directly to match backend API structure
      const response = await apiClient.auth.register({
        fullName: form.value.fullName,
        email: form.value.email,
        password: form.value.password,
        phone: form.value.phone, // Phone is required by backend
        role: roleMapping[form.value.role] || form.value.role
        // Removed status field as it's not expected by the backend
      })
      
      console.log('Registration successful:', response)
      
      // Set success state
      needsApproval.value = form.value.role === 'tenant'
      
      // Show both success indicators
      showSuccessModal.value = true
      registrationSuccess.value = true
      
      // Show a toast notification
      success(needsApproval.value 
        ? 'Your temple admin account has been created! You\'ll be notified after approval.' 
        : 'Your account has been created successfully!')
      
      // Store registration result in auth store if needed
      if (authStore && typeof authStore.setRegistrationResult === 'function') {
        authStore.setRegistrationResult({
          success: true,
          needsApproval: needsApproval.value,
          message: needsApproval.value 
            ? 'Your temple admin account has been created. You\'ll be notified after approval.'
            : 'Your account has been created successfully.'
        })
      }
      
      // Clear the form
      resetForm()
      
      // Automatically redirect after a short delay (3 seconds)
      setTimeout(() => {
        goToLogin()
      }, 3000)
    } catch (apiError) {
      console.error('API Error during registration:', apiError)
      throw apiError // Re-throw to be caught by the outer catch block
    }
    
  } catch (apiError) {
    console.error('Registration error:', apiError)
    
    // Handle specific backend validation errors
    if (apiError.response?.data?.error) {
      error(apiError.response.data.error)
      
      // Map backend errors to form fields if available
      if (apiError.response.data.errors) {
        const backendErrors = apiError.response.data.errors
        
        // Map backend error fields to form fields
        Object.keys(backendErrors).forEach(field => {
          if (field in errors.value) {
            errors.value[field] = backendErrors[field]
          }
        })
      }
    } else {
      error('An error occurred during registration. Please try again.')
    }
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  form.value = {
    fullName: '',
    email: '',
    password: '',
    phone: '',
    role: '',
    acceptTerms: false
  }
  errors.value = {}
}

// Single function for redirection to login
const goToLogin = () => {
  showSuccessModal.value = false
  
  // Always redirect to login page
  nextTick(() => {
    router.push('/login')
  })
}
</script>