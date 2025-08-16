<template>
  <div class="container mx-auto px-4 py-8">
    <h1 class="text-2xl font-bold text-gray-900 mb-6">Reset User Password</h1>
    
    <!-- Search Section -->
    <div class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-lg font-medium text-gray-700 mb-4">Search User by Email</h2>
      <div class="flex flex-col md:flex-row md:space-x-4 space-y-4 md:space-y-0">
        <div class="flex-1">
          <input 
            type="email" 
            v-model="searchEmail" 
            placeholder="Enter user email address" 
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:ring-indigo-500 focus:border-indigo-500"
            @keyup.enter="searchUser"
          />
        </div>
        <button 
          @click="searchUser" 
          class="px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
          :disabled="loading"
        >
          <span v-if="loading">Searching...</span>
          <span v-else>Search</span>
        </button>
      </div>
    </div>
    
    <!-- User Results Section (shown only when searchPerformed is true) -->
    <div v-if="searchPerformed" class="bg-white rounded-lg shadow-md p-6 mb-8">
      <h2 class="text-lg font-medium text-gray-700 mb-4">Search Results</h2>
      
      <div v-if="loading" class="flex justify-center py-8">
        <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-700"></div>
      </div>
      
      <div v-else-if="userFound" class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">NAME</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">EMAIL</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ROLE</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ACTION</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ foundUser.full_name || foundUser.fullName }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-500">{{ foundUser.email }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-green-100 text-green-800">
                  {{ getRoleName(foundUser) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button 
                  @click="navigateToResetForm(foundUser)" 
                  class="text-indigo-600 hover:text-indigo-900 bg-indigo-50 px-3 py-1 rounded-md"
                >
                  Reset Password
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      
      <div v-else class="text-center py-8">
        <div class="text-gray-500">No user found with that email address.</div>
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
      Password updated successfully!
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import superAdminService from '@/services/superadmin.service';

const router = useRouter();

// State for search functionality
const searchEmail = ref('');
const searchPerformed = ref(false);
const loading = ref(false);
const userFound = ref(false);
const foundUser = ref(null);

// State for toast notification
const showToast = ref(false);

// Search user by email
const searchUser = async () => {
  if (!searchEmail.value) return;
  
  searchPerformed.value = true;
  loading.value = true;
  userFound.value = false;
  
  try {
    const response = await superAdminService.searchUserByEmail(searchEmail.value);
    console.log('Search response:', response);
    
    // Check if response has data and it's not an empty object
    if (response.success && response.data && Object.keys(response.data).length > 0) {
      userFound.value = true;
      foundUser.value = response.data;
      console.log('Found user:', foundUser.value);
    } else {
      userFound.value = false;
      foundUser.value = null;
    }
  } catch (error) {
    console.error('Error searching for user:', error);
    userFound.value = false;
    foundUser.value = null;
  } finally {
    loading.value = false;
  }
};

// Get role name from user object with different possible structures
const getRoleName = (user) => {
  if (user.role?.role_name) return user.role.role_name;
  if (user.role?.roleName) return user.role.roleName;
  if (user.roleName) return user.roleName;
  if (typeof user.role === 'string') return user.role;
  return 'User';
};

// Navigate to the reset password form page
const navigateToResetForm = (user) => {
  // Get user ID from various possible formats
  const userId = user.id || user.ID || user.user_id;
  // Navigate to the reset password form page with the user ID
  router.push({
    name: 'ResetPasswordForm',
    params: { userId }
  });
};
</script>