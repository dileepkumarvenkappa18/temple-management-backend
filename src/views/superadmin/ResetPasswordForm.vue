<template>
  <div class="container mx-auto px-4 py-8">
    <div class="flex items-center mb-6">
      <button @click="goBack" class="mr-4 p-2 rounded-full hover:bg-gray-100">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
        </svg>
      </button>
      <h1 class="text-2xl font-bold text-gray-900">Reset User Password</h1>
    </div>
    
    <div class="bg-white rounded-lg shadow-md p-6 max-w-2xl mx-auto">
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-700"></div>
      </div>
      
      <div v-else-if="loadError" class="text-center py-8">
        <div class="text-red-500 text-lg mb-4">{{ loadError }}</div>
        <button 
          @click="goBack" 
          class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
        >
          Go Back
        </button>
      </div>
      
      <div v-else>
        <div class="mb-6 p-4 bg-gray-50 rounded-md">
          <h2 class="text-lg font-medium text-gray-800 mb-2">User Information</h2>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <p class="text-sm text-gray-500">Name</p>
              <p class="font-medium">{{ user?.full_name || user?.fullName }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Email</p>
              <p class="font-medium">{{ user?.email }}</p>
            </div>
            <div>
              <p class="text-sm text-gray-500">Role</p>
              <p class="font-medium">{{ getRoleName(user) }}</p>
            </div>
          </div>
        </div>
        
        <form @submit.prevent="resetPassword" class="space-y-6">
          <div>
            <label for="newPassword" class="block text-sm font-medium text-gray-700 mb-1">New Password</label>
            <input 
              type="password" 
              id="newPassword"
              v-model="newPassword" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{'border-red-500': passwordError}"
              required
              minlength="8"
            />
            <div class="mt-2">
              <div class="text-xs text-gray-500 mb-2">Password must contain:</div>
              <ul class="text-xs space-y-1 ml-5 list-disc">
                <li :class="newPassword.length >= 8 ? 'text-green-600' : 'text-gray-500'">
                  At least 8 characters
                </li>
                <li :class="/[A-Z]/.test(newPassword) ? 'text-green-600' : 'text-gray-500'">
                  At least one uppercase letter
                </li>
                <li :class="/[a-z]/.test(newPassword) ? 'text-green-600' : 'text-gray-500'">
                  At least one lowercase letter
                </li>
                <li :class="/\d/.test(newPassword) ? 'text-green-600' : 'text-gray-500'">
                  At least one number
                </li>
              </ul>
            </div>
          </div>
          
          <div>
            <label for="confirmPassword" class="block text-sm font-medium text-gray-700 mb-1">Confirm Password</label>
            <input 
              type="password" 
              id="confirmPassword"
              v-model="confirmPassword" 
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
              :class="{'border-red-500': passwordError}"
              required
            />
            <p v-if="passwordError" class="mt-1 text-sm text-red-600">{{ passwordError }}</p>
          </div>
          
          <div class="pt-2">
            <button 
              type="submit" 
              class="w-full px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
              :disabled="!isFormValid || resetLoading"
            >
              <span v-if="resetLoading">Updating Password...</span>
              <span v-else>Update Password</span>
            </button>
          </div>
        </form>
      </div>
    </div>
    
    <!-- Success Toast -->
    <div 
      v-if="showToast" 
      class="fixed bottom-4 right-4 bg-green-500 text-white px-6 py-3 rounded-md shadow-lg z-50 flex items-center"
    >
      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
      </svg>
      {{ toastMessage }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import superAdminService from '@/services/superadmin.service';

const router = useRouter();
const route = useRoute();
const userId = route.params.userId;

// State
const user = ref(null);
const loading = ref(true);
const loadError = ref('');
const newPassword = ref('');
const confirmPassword = ref('');
const passwordError = ref('');
const resetLoading = ref(false);
const showToast = ref(false);
const toastMessage = ref('');

// Fetch user details
const fetchUserDetails = async () => {
  loading.value = true;
  loadError.value = '';
  
  try {
    // First try to get user by ID
    const response = await superAdminService.getUserById(userId);
    
    if (response.success && response.data) {
      user.value = response.data;
    } else {
      loadError.value = 'Unable to load user details. User may not exist.';
    }
  } catch (error) {
    console.error('Error fetching user details:', error);
    loadError.value = 'Error loading user details. Please try again.';
  } finally {
    loading.value = false;
  }
};

// Validate password
const validatePassword = () => {
  if (newPassword.value.length < 8) {
    passwordError.value = 'Password must be at least 8 characters long';
    return false;
  }
  
  if (!/[A-Z]/.test(newPassword.value)) {
    passwordError.value = 'Password must contain at least one uppercase letter';
    return false;
  }
  
  if (!/[a-z]/.test(newPassword.value)) {
    passwordError.value = 'Password must contain at least one lowercase letter';
    return false;
  }
  
  if (!/\d/.test(newPassword.value)) {
    passwordError.value = 'Password must contain at least one number';
    return false;
  }
  
  if (newPassword.value !== confirmPassword.value) {
    passwordError.value = 'Passwords do not match';
    return false;
  }
  
  passwordError.value = '';
  return true;
};

// Computed property to check if form is valid
const isFormValid = computed(() => {
  return newPassword.value && 
         confirmPassword.value && 
         newPassword.value === confirmPassword.value &&
         newPassword.value.length >= 8 &&
         /[A-Z]/.test(newPassword.value) &&
         /[a-z]/.test(newPassword.value) &&
         /\d/.test(newPassword.value);
});

// Reset password function
const resetPassword = async () => {
  if (!validatePassword()) return;
  
  try {
    resetLoading.value = true;
    
    const result = await superAdminService.resetUserPassword(
      userId, 
      newPassword.value
    );
    
    if (result.success) {
      // Show success toast
      toastMessage.value = 'Password updated successfully! An email has been sent to the user.';
      showToast.value = true;
      
      setTimeout(() => {
        showToast.value = false;
        // Navigate back to the password reset list page
        router.push({ name: 'ResetPassword' });
      }, 3000);
    } else {
      // Display error
      passwordError.value = result.message || 'Failed to reset password. Please try again.';
    }
  } catch (error) {
    console.error('Error resetting password:', error);
    passwordError.value = 'An unexpected error occurred. Please try again.';
  } finally {
    resetLoading.value = false;
  }
};

// Go back function
const goBack = () => {
  router.push({ name: 'ResetPassword' });
};

// Get role name from user object with different possible structures
const getRoleName = (user) => {
  if (!user) return 'User';
  if (user.role?.role_name) return user.role.role_name;
  if (user.role?.roleName) return user.role.roleName;
  if (user.roleName) return user.roleName;
  if (typeof user.role === 'string') return user.role;
  return 'User';
};

// Watch for password changes to validate in real-time
watch([newPassword, confirmPassword], () => {
  if (newPassword.value && confirmPassword.value) {
    validatePassword();
  } else {
    passwordError.value = '';
  }
});

// Fetch user data on component mount
onMounted(() => {
  if (userId) {
    fetchUserDetails();
  } else {
    loadError.value = 'No user ID provided. Please go back and try again.';
    loading.value = false;
  }
});
</script>