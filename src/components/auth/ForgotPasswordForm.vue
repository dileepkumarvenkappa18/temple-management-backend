<template>
  <BaseCard class="w-full max-w-md mx-auto">
    <div class="text-center mb-8">
      <h2 class="text-3xl font-bold text-gray-900 mb-2" style="font-family: 'Roboto', sans-serif;">
        Reset Password
      </h2>
      <p class="text-gray-600">
        {{ emailSent ? 
          'Check your email for reset instructions' : 
          'Enter your email to receive reset instructions' 
        }}
      </p>
    </div>

    <div v-if="emailSent" class="text-center space-y-6">
      <!-- Success State -->
      <div class="bg-green-50 border border-green-200 rounded-xl p-6">
        <div class="flex items-center justify-center w-12 h-12 mx-auto mb-4 bg-green-100 rounded-full">
          <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-green-900 mb-2" style="font-family: sans-serif;">
          Email Sent!
        </h3>
        <p class="text-green-700 text-sm mb-4">
          We've sent password reset instructions to <strong>{{ form.email }}</strong>
        </p>
        <p class="text-green-600 text-xs">
          Didn't receive the email? Check your spam folder or try again in a few minutes.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="space-y-3">
        <BaseButton
          @click="resetForm"
          variant="outline"
          size="md"
          class="w-full"
        >
          Try Different Email
        </BaseButton>
        
        <BaseButton
          @click="router.push('/login')"
          variant="primary"
          size="md"
          class="w-full"
        >
          Back to Login
        </BaseButton>
      </div>
    </div>

    <form v-else @submit.prevent="handleForgotPassword" class="space-y-6">
      <!-- Email Field -->
      <div>
        <label for="email" class="block text-sm font-medium text-gray-700 mb-2" style="font-family: sans-serif;">
          Email Address
        </label>
        <BaseInput
          id="email"
          v-model="form.email"
          type="email"
          placeholder="Enter your registered email"
          :error="errors.email"
          required
          autocomplete="email"
        />
        <p class="text-sm text-gray-500 mt-2">
          Enter the email address associated with your account
        </p>
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
        Send Reset Instructions
      </BaseButton>

      <!-- Back to Login Link -->
      <div class="text-center pt-4 border-t border-gray-200">
        <p class="text-sm text-gray-600">
          Remember your password?
          <router-link 
            to="/login" 
            class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
          >
            Back to Login
          </router-link>
        </p>
      </div>
    </form>

    <!-- Help Section -->
    <div class="mt-8 p-4 bg-gray-50 rounded-xl">
      <h4 class="text-sm font-medium text-gray-900 mb-2" style="font-family: sans-serif;">
        Need Help?
      </h4>
      <p class="text-xs text-gray-600 leading-relaxed">
        If you're having trouble accessing your account, please contact our support team. 
        We're here to help you get back to managing your temple activities.
      </p>
    </div>
  </BaseCard>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'

const router = useRouter()
const authStore = useAuthStore()
const { showToast } = useToast()

const isLoading = ref(false)
const emailSent = ref(false)
const form = ref({
  email: ''
})

const errors = ref({
  email: ''
})

const isFormValid = computed(() => {
  return form.value.email && !errors.value.email
})

const validateForm = () => {
  errors.value = {}
  
  // Email validation
  if (!form.value.email) {
    errors.value.email = 'Email is required'
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.value.email)) {
    errors.value.email = 'Please enter a valid email address'
  }
  
  return Object.keys(errors.value).length === 0
}

const handleForgotPassword = async () => {
  if (!validateForm()) return
  
  isLoading.value = true
  
  try {
    const result = await authStore.forgotPassword(form.value.email)
    
    if (result.success) {
      emailSent.value = true
      showToast('Password reset instructions sent to your email', 'success')
    } else {
      showToast(result.message || 'Failed to send reset instructions. Please try again.', 'error')
    }
  } catch (error) {
    console.error('Forgot password error:', error)
    showToast('An error occurred. Please try again later.', 'error')
  } finally {
    isLoading.value = false
  }
}

const resetForm = () => {
  emailSent.value = false
  form.value.email = ''
  errors.value = {}
}
</script>