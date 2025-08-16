<template>
  <BaseCard class="w-full max-w-md mx-auto">
    <div class="text-center mb-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-2">Create Account</h2>
      <p class="text-gray-600">Join our temple management platform</p>
    </div>

    <!-- Registration Success Alert -->
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
          Full Name <span class="text-red-500">*</span>
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
          Email Address <span class="text-red-500">*</span>
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
          Password <span class="text-red-500">*</span>
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
        <label class="block text-sm font-medium text-gray-700 mb-2">
          I want to join as 
        </label>
        
        <!-- Simple Role Cards -->
        <div class="grid grid-cols-3 gap-4">
          <div
            v-for="role in roleOptions"
            :key="role.value"
            :class="[
              'border-2 rounded-lg p-4 text-center cursor-pointer transition-all',
              form.role === role.value 
                ? 'border-indigo-600 bg-indigo-50' 
                : 'border-gray-200 hover:border-indigo-300'
            ]"
            @click="selectRole(role.value)"
          >
            <div class="flex flex-col items-center">
              <!-- Avatar -->
              <div :class="[
                'w-16 h-16 rounded-full bg-gray-300 flex items-center justify-center mb-3',
                form.role === role.value ? 'bg-indigo-100' : ''
              ]">
                <svg class="w-8 h-8 text-gray-600" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                </svg>
              </div>
              
              <!-- Role Name -->
              <span class="font-medium">{{ role.label }}</span>
            </div>
          </div>
        </div>
        
        <!-- Role Error -->
        <div v-if="errors.role" class="mt-1 text-sm text-red-600">
          {{ errors.role }}
        </div>
      </div>
      
      <!-- Temple Details Form (shown when Temple Admin is selected) -->
      <div v-if="form.role === 'tenant'" class="bg-gray-50 border border-gray-200 p-4 rounded-lg">
        <h3 class="text-sm font-medium text-gray-700 mb-3">Temple Details</h3>
        
        <!-- Temple Name -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Temple Name <span class="text-red-500">*</span>
          </label>
          <input
            v-model="templeDetails.name"
            type="text"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="Enter temple name"
          />
          <div v-if="errors.templeName" class="mt-1 text-sm text-red-600">
            {{ errors.templeName }}
          </div>
        </div>
        
        <!-- Temple Place -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Temple Place <span class="text-red-500">*</span>
          </label>
          <input
            v-model="templeDetails.place"
            type="text"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="Enter temple location/city"
          />
          <div v-if="errors.templePlace" class="mt-1 text-sm text-red-600">
            {{ errors.templePlace }}
          </div>
        </div>
        
        <!-- Temple Address -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Temple Address <span class="text-red-500">*</span>
          </label>
          <textarea
            v-model="templeDetails.address"
            rows="2"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="Enter full temple address"
          ></textarea>
          <div v-if="errors.templeAddress" class="mt-1 text-sm text-red-600">
            {{ errors.templeAddress }}
          </div>
        </div>
        
        <!-- Temple Phone -->
        <div class="mb-4">
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Temple Phone <span class="text-red-500">*</span>
          </label>
          <input
            v-model="templeDetails.phoneNumber"
            type="tel"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="Enter temple contact number"
          />
          <div v-if="errors.templePhoneNo" class="mt-1 text-sm text-red-600">
            {{ errors.templePhoneNo }}
          </div>
        </div>
        
        <!-- Temple Description -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">
            Temple Description
          </label>
          <textarea
            v-model="templeDetails.description"
            rows="3"
            class="w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500 sm:text-sm"
            placeholder="Briefly describe your temple (optional)"
          ></textarea>
        </div>
      </div>

      <!-- Terms and Privacy -->
      <div class="flex items-start">
        <input
          id="terms"
          v-model="form.acceptTerms"
          type="checkbox"
          required
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
          <span class="text-red-500">*</span>
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

    <!-- Registration Success Modal -->
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
import { apiClient } from '@/plugins/axios'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import PasswordStrengthMeter from '@/components/auth/PasswordStrengthMeter.vue'

const router = useRouter()
const authStore = useAuthStore()
const { success, error: showError } = useToast()

// State
const isLoading = ref(false)
const showSuccessModal = ref(false)
const registrationSuccess = ref(false)
const needsApproval = ref(false)
const errors = ref({})

// Fixed role options for reliability
const roleOptions = [
  { value: 'tenant', label: 'Temple Admin' },
  { value: 'devotee', label: 'Devotee' },
  { value: 'volunteer', label: 'Volunteer' }
]

// Form data
const form = ref({
  fullName: '',
  email: '',
  password: '',
  phone: '',
  role: '',
  acceptTerms: false
})

// Temple details (separate for clarity)
const templeDetails = ref({
  name: '',
  place: '',
  address: '',
  phoneNumber: '',
  description: ''
})

// Computed
const isFormValid = computed(() => {
  const hasBasicFields = form.value.fullName && 
                        form.value.email && 
                        form.value.password && 
                        form.value.phone && 
                        form.value.role &&
                        form.value.acceptTerms;
                        
  // Temple admin requires temple details
  const hasTempleDetails = form.value.role !== 'tenant' || 
                          (templeDetails.value.name && 
                           templeDetails.value.place && 
                           templeDetails.value.address && 
                           templeDetails.value.phoneNumber);
                           
  return hasBasicFields && hasTempleDetails && Object.keys(errors.value).length === 0;
})

// Methods
const selectRole = (role) => {
  console.log('Role selected:', role)
  form.value.role = role
  
  // Clear temple-related errors when changing roles
  if (role !== 'tenant') {
    Object.keys(errors.value).forEach(key => {
      if (key.startsWith('temple')) {
        delete errors.value[key]
      }
    })
  }
}

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
  }
  
  // Password validation
  if (!form.value.password) {
    errors.value.password = 'Password is required'
  } else if (form.value.password.length < 8) {
    errors.value.password = 'Password must be at least 8 characters'
  }
  
  // Phone validation
  if (!form.value.phone) {
    errors.value.phone = 'Phone number is required'
  } else if (!/^[\+]?[1-9][\d]{0,15}$/.test(form.value.phone.replace(/\s/g, ''))) {
    errors.value.phone = 'Please enter a valid phone number'
  }
  
  // Role validation
  if (!form.value.role) {
    errors.value.role = 'Please select your role'
  }
  
  // Temple details validation (only for temple admin)
  if (form.value.role === 'tenant') {
    if (!templeDetails.value.name) {
      errors.value.templeName = 'Temple name is required'
    }
    
    if (!templeDetails.value.place) {
      errors.value.templePlace = 'Temple place is required'
    }
    
    if (!templeDetails.value.address) {
      errors.value.templeAddress = 'Temple address is required'
    }
    
    if (!templeDetails.value.phoneNumber) {
      errors.value.templePhoneNo = 'Temple phone number is required'
    }
  }
  
  // Terms validation
  if (!form.value.acceptTerms) {
    errors.value.acceptTerms = 'You must accept the terms and privacy policy'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleRegister = async () => {
  console.log('Register form submitted:', form.value)
  
  if (!validateForm()) {
    console.log('Form validation failed. Errors:', errors.value)
    return
  }
  
  isLoading.value = true
  
  try {
    // Map frontend role to backend role
    const roleMapping = {
      'tenant': 'templeadmin',
      'devotee': 'devotee', 
      'volunteer': 'volunteer'
    }

    // Prepare registration data
    const registrationData = {
      fullName: form.value.fullName,
      email: form.value.email,
      password: form.value.password,
      phone: form.value.phone,
      role: roleMapping[form.value.role] || form.value.role
    }
    
    // Add temple details for temple admin role
    if (form.value.role === 'tenant') {
      registrationData.templeName = templeDetails.value.name
      registrationData.templePlace = templeDetails.value.place
      registrationData.templeAddress = templeDetails.value.address
      registrationData.templePhoneNo = templeDetails.value.phoneNumber
      registrationData.templeDescription = templeDetails.value.description || ''
    }
    
    console.log('Sending registration data to API:', registrationData)
    
    // Send registration request
    const response = await apiClient.auth.register(registrationData)
    console.log('Registration successful:', response)
    
    // Set success state
    needsApproval.value = form.value.role === 'tenant'
    showSuccessModal.value = true
    registrationSuccess.value = true
    
    // Show success message
    success(needsApproval.value 
      ? 'Your temple admin account has been created! You\'ll be notified after approval.' 
      : 'Your account has been created successfully!')
    
    // Store registration result in auth store if available
    if (authStore && typeof authStore.setRegistrationResult === 'function') {
      authStore.setRegistrationResult({
        success: true,
        needsApproval: needsApproval.value,
        message: needsApproval.value 
          ? 'Your temple admin account has been created. You\'ll be notified after approval.'
          : 'Your account has been created successfully.'
      })
    }
    
    // Clear form
    resetForm()
    
    // Redirect after delay
    setTimeout(() => {
      goToLogin()
    }, 3000)
  } catch (err) {
    console.error('Registration error:', err)
    
    // Handle API errors
    if (err.response?.data?.error) {
      showError(err.response.data.error)
      
      // Map backend errors to form fields
      if (err.response.data.errors) {
        const backendErrors = err.response.data.errors
        
        Object.keys(backendErrors).forEach(field => {
          errors.value[field] = backendErrors[field]
        })
      }
    } else {
      showError('Registration failed. Please try again.')
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
  
  templeDetails.value = {
    name: '',
    place: '',
    address: '',
    phoneNumber: '',
    description: ''
  }
  
  errors.value = {}
}

const goToLogin = () => {
  showSuccessModal.value = false
  
  nextTick(() => {
    router.push('/login')
  })
}
</script>