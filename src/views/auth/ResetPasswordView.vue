<template>
  <div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Reset Your Password</h2>
      <p class="mt-2 text-center text-sm text-gray-600">Set a new password for your account</p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <BaseCard class="w-full max-w-md mx-auto">
        <div v-if="tokenValidating" class="text-center py-4">
          <div class="animate-spin rounded-full h-10 w-10 border-b-2 border-indigo-600 mx-auto"></div>
          <p class="mt-2 text-gray-600">Validating your reset token...</p>
        </div>

        <div v-else-if="error" class="text-center py-4">
          <div class="rounded-full bg-red-100 p-2 mx-auto w-12 h-12 flex items-center justify-center mb-4">
            <svg class="w-6 h-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-red-800">Invalid or Expired Link</h3>
          <p class="mt-2 text-sm text-gray-600">
            {{ error || 'The password reset link is invalid or has expired. Please request a new password reset link.' }}
          </p>
          <BaseButton 
            @click="router.push('/forgot-password')" 
            variant="primary"
            size="lg"
            class="mt-4 w-full"
          >
            Request New Link
          </BaseButton>
        </div>

        <form v-else-if="!success" @submit.prevent="handleResetPassword" class="space-y-6">
          <div>
            <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
              New Password
            </label>
            <BaseInput
              id="password"
              v-model="form.password"
              type="password"
              placeholder="Enter your new password"
              :error="errors.password"
              required
            />
          </div>

          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-2">
              Confirm New Password
            </label>
            <BaseInput
              id="confirmPassword"
              v-model="form.confirmPassword"
              type="password"
              placeholder="Confirm your new password"
              :error="errors.confirmPassword"
              required
            />
          </div>

          <!-- Password Strength Meter -->
          <div v-if="form.password">
            <div class="mt-2">
              <div class="h-2 bg-gray-200 rounded-full overflow-hidden">
                <div 
                  class="h-full rounded-full transition-all duration-300"
                  :class="{
                    'bg-red-500': passwordStrength.score < 3,
                    'bg-yellow-500': passwordStrength.score === 3,
                    'bg-green-500': passwordStrength.score > 3
                  }"
                  :style="{ width: `${passwordStrength.score * 20}%` }"
                ></div>
              </div>
              <div class="mt-2 grid grid-cols-1 gap-y-1 text-xs">
                <div :class="passwordStrength.criteria.minLength ? 'text-green-600' : 'text-gray-500'">
                  ✓ At least 8 characters
                </div>
                <div :class="passwordStrength.criteria.hasUppercase ? 'text-green-600' : 'text-gray-500'">
                  ✓ At least 1 uppercase letter
                </div>
                <div :class="passwordStrength.criteria.hasLowercase ? 'text-green-600' : 'text-gray-500'">
                  ✓ At least 1 lowercase letter
                </div>
                <div :class="passwordStrength.criteria.hasNumber ? 'text-green-600' : 'text-gray-500'">
                  ✓ At least 1 number
                </div>
                <div :class="passwordStrength.criteria.hasSpecialChar ? 'text-green-600' : 'text-gray-500'">
                  ✓ At least 1 special character
                </div>
              </div>
            </div>
          </div>

          <BaseButton
            type="submit"
            variant="primary"
            size="lg"
            :loading="isLoading"
            :disabled="!isFormValid"
            class="w-full"
          >
            Reset Password
          </BaseButton>
        </form>

        <div v-else class="text-center space-y-6">
          <!-- Success State -->
          <div class="bg-green-50 border border-green-200 rounded-xl p-6">
            <div class="flex items-center justify-center w-12 h-12 mx-auto mb-4 bg-green-100 rounded-full">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
            </div>
            <h3 class="text-lg font-semibold text-green-900 mb-2">Password Reset Successful!</h3>
            <p class="text-green-700 text-sm mb-4">
              Your password has been reset successfully. You can now log in with your new password.
            </p>
          </div>

          <!-- Action Button -->
          <BaseButton
            @click="router.push('/login')"
            variant="primary"
            size="lg"
            class="w-full"
          >
            Go to Login
          </BaseButton>
        </div>
      </BaseCard>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useToast } from '@/composables/useToast';
import BaseCard from '@/components/common/BaseCard.vue';
import BaseInput from '@/components/common/BaseInput.vue';
import BaseButton from '@/components/common/BaseButton.vue';

const route = useRoute();
const router = useRouter();
const authStore = useAuthStore();
const { showToast } = useToast();

// Get token from route query with ref for reactivity
const token = ref(route.query.token || '');

// State
const tokenValidating = ref(true); // Start with validating state
const error = ref('');
const success = ref(false);
const isLoading = ref(false);
const form = ref({
  password: '',
  confirmPassword: ''
});
const errors = ref({
  password: '',
  confirmPassword: ''
});

// Validate password strength
const passwordStrength = computed(() => {
  const password = form.value.password;
  const criteria = {
    minLength: password.length >= 8,
    hasUppercase: /[A-Z]/.test(password),
    hasLowercase: /[a-z]/.test(password),
    hasNumber: /\d/.test(password),
    hasSpecialChar: /[!@#$%^&*(),.?":{}|<>]/.test(password)
  };
  
  const score = Object.values(criteria).filter(Boolean).length;
  
  return {
    criteria,
    score,
    strength: score < 3 ? 'weak' : score < 5 ? 'medium' : 'strong'
  };
});

// Check if form is valid
const isFormValid = computed(() => {
  return (
    form.value.password && 
    form.value.confirmPassword && 
    form.value.password === form.value.confirmPassword &&
    passwordStrength.value.score >= 3
  );
});

// Watch for route query changes
watch(() => route.query.token, (newToken) => {
  token.value = newToken || '';
  validateToken();
});

// Validate token
const validateToken = async () => {
  console.log('Validating token:', token.value);
  console.log('Current route path:', route.path);
  tokenValidating.value = true;
  
  try {
    if (!token.value) {
      error.value = 'No reset token provided. Please check your email link.';
      console.error('No token provided in URL');
      return false;
    }
    
    // For development: just validate token exists
    setTimeout(() => {
      tokenValidating.value = false;
    }, 1000);
    
    return true;
  } catch (err) {
    console.error('Token validation error:', err);
    error.value = 'Error validating token';
    return false;
  } finally {
    // This will be called if there's no actual async operation
    // For development, we'll keep the artificial delay
  }
};

// Reset password handler
const handleResetPassword = async () => {
  // Validate form
  errors.value = {};
  
  if (!form.value.password) {
    errors.value.password = 'Password is required';
  } else if (form.value.password.length < 8) {
    errors.value.password = 'Password must be at least 8 characters';
  }
  
  if (!form.value.confirmPassword) {
    errors.value.confirmPassword = 'Please confirm your password';
  } else if (form.value.password !== form.value.confirmPassword) {
    errors.value.confirmPassword = 'Passwords do not match';
  }
  
  if (Object.keys(errors.value).length > 0) return;
  
  // Submit form
  isLoading.value = true;
  
  try {
    console.log('Submitting reset password with token:', token.value);
    const result = await authStore.resetPassword(token.value, form.value.password);
    console.log('Password reset result:', result);
    
    if (result?.success) {
      success.value = true;
      showToast('Password has been reset successfully', 'success');
    } else {
      error.value = result?.error || 'Failed to reset password';
      showToast(error.value, 'error');
    }
  } catch (err) {
    console.error('Password reset error:', err);
    error.value = 'An error occurred. Please try again later.';
    showToast(error.value, 'error');
  } finally {
    isLoading.value = false;
  }
};

// Single onMounted hook with debugging
onMounted(() => {
  console.log('ResetPasswordView mounted!');
  console.log('Current route:', route.path);
  console.log('Route query params:', route.query);
  console.log('Token from query:', token.value);
  
  // Add a delay to ensure the token is available from the route
  setTimeout(() => {
    validateToken();
  }, 100);
});
</script>