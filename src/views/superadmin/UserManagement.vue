<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-200">
    <div class="p-6 border-b border-gray-200 flex justify-between items-center">
      <div>
        <h3 class="text-lg font-semibold text-gray-900 font-roboto">User Management</h3>
        <p class="text-sm text-gray-600 mt-1">Create and manage system users</p>
      </div>
      <!-- Regular button instead of BaseButton to ensure it works -->
      <button 
        @click="openCreateModal" 
        class="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors"
      >
        <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
        </svg>
        Create User
      </button>
    </div>
    
    <div class="p-6">
      <!-- Users List Table -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Full Name</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Phone</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Role</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="isUsersLoading" class="animate-pulse">
              <td colspan="6" class="px-6 py-4">
                <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              </td>
            </tr>
            <tr v-else-if="users.length === 0">
              <td colspan="6" class="px-6 py-4 text-center text-sm text-gray-500">
                No users found. Create your first user by clicking the 'Create User' button.
              </td>
            </tr>
            <tr v-for="user in users" :key="user.id" class="hover:bg-gray-50">
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm font-medium text-gray-900">{{ user.full_name }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-500">{{ user.email }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-500">{{ user.phone }}</div>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full bg-indigo-100 text-indigo-800">
                  {{ getRoleDisplay(user.role) }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="relative inline-block w-10 mr-2 align-middle select-none">
                  <input 
                    type="checkbox" 
                    :id="`toggle-${user.id}`" 
                    :checked="user.status === 'active'" 
                    @change="toggleUserStatus(user)"
                    class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                  />
                  <label 
                    :for="`toggle-${user.id}`" 
                    class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                  ></label>
                </div>
                <span class="text-xs ml-1">{{ user.status === 'active' ? 'Active' : 'Inactive' }}</span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                <button 
                  @click="editUser(user)" 
                  class="text-indigo-600 hover:text-indigo-900 mr-3"
                >
                  Edit
                </button>
                <user-assign-button :user="user" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>

    <!-- Direct Modal Implementation (use this if BaseModal doesn't work) -->
    <div v-if="showCreateModal" class="fixed inset-0 z-50 overflow-auto bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-lg shadow-xl max-w-3xl w-full mx-4 max-h-[90vh] overflow-y-auto">
        <!-- Modal Header -->
        <div class="px-6 py-4 border-b border-gray-200 flex justify-between items-center">
          <h3 class="text-lg font-medium text-gray-900">{{ isEditing ? 'Edit User' : 'Create New User' }}</h3>
          <button @click="cancelUserForm" class="text-gray-400 hover:text-gray-500">
            <svg class="h-6 w-6" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Modal Body -->
        <div class="px-6 py-4">
          <!-- Success Alert -->
          <div v-if="registrationSuccess" class="mb-6 bg-green-50 border border-green-200 rounded-lg p-4">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-green-500" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-green-800">User {{ isEditing ? 'Updated' : 'Created' }} Successfully!</h3>
                <div class="mt-1 text-sm text-green-700">
                  <p>
                    User account has been {{ isEditing ? 'updated' : 'created' }} successfully.
                  </p>
                </div>
              </div>
            </div>
          </div>

          <form v-if="!registrationSuccess" @submit.prevent="handleSubmitUser" class="space-y-6">
            <!-- Full Name Field -->
            <div>
              <label for="fullName" class="block text-sm font-medium text-gray-700 mb-2">
                Full Name <span class="text-red-500">*</span>
              </label>
              <input
                id="fullName"
                v-model="form.fullName"
                type="text"
                placeholder="Enter full name"
                required
                autocomplete="name"
                class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                :class="{ 'border-red-500': errors.fullName }"
              />
              <div v-if="errors.fullName" class="mt-1 text-sm text-red-600">
                {{ errors.fullName }}
              </div>
            </div>

            <!-- Email Field -->
            <div>
              <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
                Email Address <span class="text-red-500">*</span>
              </label>
              <input
                id="email"
                v-model="form.email"
                type="email"
                placeholder="Enter email"
                required
                autocomplete="email"
                class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                :class="{ 'border-red-500': errors.email }"
              />
              <div v-if="errors.email" class="mt-1 text-sm text-red-600">
                {{ errors.email }}
              </div>
            </div>

            <!-- Password Field (only for new users) -->
            <div v-if="!isEditing">
              <label for="password" class="block text-sm font-medium text-gray-700 mb-2">
                Password <span class="text-red-500">*</span>
              </label>
              <input
                id="password"
                v-model="form.password"
                type="password"
                placeholder="Create a strong password"
                required
                autocomplete="new-password"
                class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                :class="{ 'border-red-500': errors.password }"
              />
              <div v-if="errors.password" class="mt-1 text-sm text-red-600">
                {{ errors.password }}
              </div>
              <PasswordStrengthMeter 
                :password="form.password" 
                class="mt-2"
              />
            </div>

            <!-- Phone Field -->
            <div>
              <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
                Phone Number <span class="text-red-500">*</span>
              </label>
              <input
                id="phone"
                v-model="form.phone"
                type="tel"
                placeholder="Enter phone number"
                required
                autocomplete="tel"
                class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                :class="{ 'border-red-500': errors.phone }"
              />
              <div v-if="errors.phone" class="mt-1 text-sm text-red-600">
                {{ errors.phone }}
              </div>
            </div>

            <!-- Role Dropdown Selector -->
            <div>
              <label for="role" class="block text-sm font-medium text-gray-700 mb-2">
                User Role <span class="text-red-500">*</span>
              </label>
              <select
                id="role"
                v-model="form.role"
                class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                required
                :disabled="isRolesLoading"
              >
                <option value="" disabled>{{ isRolesLoading ? 'Loading roles...' : 'Select a role' }}</option>
                <option 
                  v-for="role in roles" 
                  :key="role.id" 
                  :value="role.role_name"
                >
                  {{ role.role_name }}
                </option>
              </select>
              <div v-if="roles.length === 0 && !isRolesLoading" class="mt-2 text-sm text-amber-600">
                No roles available. Please check role configuration.
              </div>
              <div v-if="errors.role" class="mt-2 text-sm text-red-600">
                {{ errors.role }}
              </div>
            </div>

            <!-- Temple details section (conditional - only for temple admin role) -->
            <div v-if="form.role === 'templeadmin'" class="space-y-4 p-4 bg-gray-50 rounded-lg border border-gray-200">
              <h4 class="font-medium text-gray-900">Temple Details</h4>
              
              <!-- Temple Name -->
              <div>
                <label for="templeName" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Name <span class="text-red-500">*</span>
                </label>
                <input
                  id="templeName"
                  v-model="form.templeDetails.name"
                  type="text"
                  placeholder="Enter temple name"
                  required
                  class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                  :class="{ 'border-red-500': errors.templeName }"
                />
                <div v-if="errors.templeName" class="mt-1 text-sm text-red-600">
                  {{ errors.templeName }}
                </div>
              </div>
              
              <!-- Temple Place -->
              <div>
                <label for="templePlace" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Location <span class="text-red-500">*</span>
                </label>
                <input
                  id="templePlace"
                  v-model="form.templeDetails.place"
                  type="text"
                  placeholder="Enter city, state"
                  required
                  class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                  :class="{ 'border-red-500': errors.templePlace }"
                />
                <div v-if="errors.templePlace" class="mt-1 text-sm text-red-600">
                  {{ errors.templePlace }}
                </div>
              </div>
              
              <!-- Temple Address -->
              <div>
                <label for="templeAddress" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Address <span class="text-red-500">*</span>
                </label>
                <input
                  id="templeAddress"
                  v-model="form.templeDetails.address"
                  type="text"
                  placeholder="Enter complete address"
                  required
                  class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                  :class="{ 'border-red-500': errors.templeAddress }"
                />
                <div v-if="errors.templeAddress" class="mt-1 text-sm text-red-600">
                  {{ errors.templeAddress }}
                </div>
              </div>
              
              <!-- Temple Phone -->
              <div>
                <label for="templePhoneNumber" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Phone <span class="text-red-500">*</span>
                </label>
                <input
                  id="templePhoneNumber"
                  v-model="form.templeDetails.phoneNumber"
                  type="tel"
                  placeholder="Enter temple contact number"
                  required
                  class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm p-2.5"
                  :class="{ 'border-red-500': errors.templePhoneNumber }"
                />
                <div v-if="errors.templePhoneNumber" class="mt-1 text-sm text-red-600">
                  {{ errors.templePhoneNumber }}
                </div>
              </div>
              
              <!-- Temple Description -->
              <div>
                <label for="templeDescription" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Description <span class="text-red-500">*</span>
                </label>
                <textarea
                  id="templeDescription"
                  v-model="form.templeDetails.description"
                  rows="3"
                  class="block w-full border border-gray-300 rounded-lg shadow-sm focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition-colors duration-200 p-3 leading-relaxed"
                  placeholder="Brief description about the temple"
                  required
                ></textarea>
                <div v-if="errors.templeDescription" class="mt-1 text-sm text-red-600">
                  {{ errors.templeDescription }}
                </div>
              </div>
            </div>

            <!-- Status (only for editing) -->
            <div v-if="isEditing" class="flex items-center space-x-2">
              <label for="status" class="text-sm font-medium text-gray-700">User Status:</label>
              <div class="relative inline-block w-10 mr-2 align-middle select-none">
                <input 
                  type="checkbox" 
                  id="status" 
                  v-model="form.isActive" 
                  class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label 
                  for="status" 
                  class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                ></label>
              </div>
              <span class="text-sm text-gray-500">{{ form.isActive ? 'Active' : 'Inactive' }}</span>
            </div>

            <!-- Form Actions -->
            <div class="flex justify-end space-x-3 pt-3">
              <button 
                type="button" 
                class="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors"
                @click="cancelUserForm"
              >
                Cancel
              </button>
              <button 
                type="submit" 
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors flex items-center"
                :disabled="isLoading"
              >
                <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ isEditing ? 'Update User' : 'Create User' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import { useSuperAdminStore } from '@/stores/superadmin'
import PasswordStrengthMeter from '@/components/auth/PasswordStrengthMeter.vue'
import UserAssignButton from '@/components/superadmin/UserAssignButton.vue'

const { success, error } = useToast()
const superAdminStore = useSuperAdminStore()

// State variables - using store reactive state
const isLoading = computed(() => superAdminStore.loadingUserAction)
const isRolesLoading = computed(() => superAdminStore.loadingUserRoles)
const isUsersLoading = computed(() => superAdminStore.loadingUsers)
const registrationSuccess = ref(false)
const showCreateModal = ref(false)
const isEditing = ref(false)
const editingUserId = ref(null)

// Reactive data from store
const roles = computed(() => superAdminStore.userRoles)
const users = computed(() => superAdminStore.users)

// Form state
const form = ref({
  fullName: '',
  email: '',
  password: '',
  phone: '',
  role: '',
  isActive: true,
  templeDetails: {
    name: '',
    place: '',
    address: '',
    phoneNumber: '',
    description: ''
  }
})

const errors = ref({})

// Explicitly defined function to open the modal
const openCreateModal = () => {
  console.log('Opening create modal')
  showCreateModal.value = true
}

// Fetch all data when component mounts
onMounted(async () => {
  await superAdminStore.refreshUserData()
})

// Helper function to get role display name
const getRoleDisplay = (role) => {
  if (typeof role === 'object' && role.description) {
    return role.description
  }
  if (typeof role === 'object' && role.role_name) {
    return role.role_name
  }
  // Fallback for string roles
  const roleObj = roles.value.find(r => r.role_name === role)
  return roleObj ? roleObj.description || roleObj.role_name : role
}

// Form validation
const isFormValid = computed(() => {
  // Basic validation
  if (!form.value.fullName || !form.value.email || (!isEditing.value && !form.value.password) || !form.value.phone || !form.value.role) {
    return false
  }
  
  // Temple admin specific validation
  if (form.value.role === 'templeadmin') {
    if (!form.value.templeDetails.name || !form.value.templeDetails.place || 
        !form.value.templeDetails.address || !form.value.templeDetails.phoneNumber || 
        !form.value.templeDetails.description) {
      return false
    }
  }
  
  // No errors
  return Object.keys(errors.value).length === 0
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
  }
  
  // Password validation (only for new users)
  if (!isEditing.value) {
    if (!form.value.password) {
      errors.value.password = 'Password is required'
    } else if (form.value.password.length < 6) {
      errors.value.password = 'Password must be at least 6 characters'
    }
  }
  
  // Phone validation
  if (!form.value.phone) {
    errors.value.phone = 'Phone number is required'
  } else if (!/^[\+]?[1-9][\d]{0,15}$/.test(form.value.phone.replace(/\s/g, ''))) {
    errors.value.phone = 'Please enter a valid phone number'
  }
  
  // Role validation
  if (!form.value.role) {
    errors.value.role = 'Please select a user role'
  }
  
  // Validate temple details if role is templeadmin
  if (form.value.role === 'templeadmin') {
    if (!form.value.templeDetails.name) {
      errors.value.templeName = 'Temple name is required'
    }
    
    if (!form.value.templeDetails.place) {
      errors.value.templePlace = 'Temple location is required'
    }
    
    if (!form.value.templeDetails.address) {
      errors.value.templeAddress = 'Temple address is required'
    }
    
    if (!form.value.templeDetails.phoneNumber) {
      errors.value.templePhoneNumber = 'Temple phone number is required'
    }
    
    if (!form.value.templeDetails.description?.trim()) {
      errors.value.templeDescription = 'Temple description is required'
    }
  }
  
  return Object.keys(errors.value).length === 0
}

// Handle form submission (create or update user)
const handleSubmitUser = async () => {
  if (!validateForm()) return
  
  try {
    // Prepare user data according to backend structure
    const userData = {
      fullName: form.value.fullName,
      email: form.value.email,
      phone: form.value.phone,
      role: form.value.role
    }
    
    // Add password for new users
    if (!isEditing.value) {
      userData.password = form.value.password
    }
    
    // Add temple details if applicable
    if (form.value.role === 'templeadmin') {
      userData.templeName = form.value.templeDetails.name
      userData.templePlace = form.value.templeDetails.place
      userData.templeAddress = form.value.templeDetails.address
      userData.templePhoneNo = form.value.templeDetails.phoneNumber
      userData.templeDescription = form.value.templeDetails.description
    }
    
    console.log(`${isEditing.value ? 'Updating' : 'Creating'} user with data:`, userData)
    
    let result
    
    if (isEditing.value) {
      // Update existing user
      result = await superAdminStore.updateUser(editingUserId.value, userData)
    } else {
      // Create new user
      result = await superAdminStore.createUser(userData)
    }
    
    if (result.success) {
      // Set success state
      registrationSuccess.value = true
      
      // Show success notification
      success(result.message || `User ${isEditing.value ? 'updated' : 'created'} successfully!`)
      
      // Auto-close modal after success
      setTimeout(() => {
        resetForm()
        showCreateModal.value = false
        registrationSuccess.value = false
      }, 2000)
    } else {
      error(result.error || `Failed to ${isEditing.value ? 'update' : 'create'} user`)
    }
    
  } catch (apiError) {
    console.error(`Error ${isEditing.value ? 'updating' : 'creating'} user:`, apiError)
    error(`Failed to ${isEditing.value ? 'update' : 'create'} user. Please try again.`)
  }
}

// Toggle user active status
const toggleUserStatus = async (user) => {
  try {
    const newStatus = user.status === 'active' ? 'inactive' : 'active'
    
    const result = await superAdminStore.updateUserStatus(user.id, newStatus)
    
    if (result.success) {
      success(result.message || `User status updated to ${newStatus}`)
    } else {
      error(result.error || 'Failed to update user status')
    }
  } catch (err) {
    console.error('Error toggling user status:', err)
    error('Failed to update user status')
  }
}

// Edit user
const editUser = (user) => {
  isEditing.value = true
  editingUserId.value = user.id
  
  // Clone user data to form
  form.value = {
    fullName: user.full_name,
    email: user.email,
    phone: user.phone,
    role: user.role.role_name,
    isActive: user.status === 'active',
    templeDetails: {
      name: user.temple_details?.temple_name || '',
      place: user.temple_details?.temple_place || '',
      address: user.temple_details?.temple_address || '',
      phoneNumber: user.temple_details?.temple_phone_no || '',
      description: user.temple_details?.temple_description || ''
    }
  }
  
  // Show modal
  console.log('Opening edit modal for user:', user.id)
  showCreateModal.value = true
}

// Cancel form
const cancelUserForm = () => {
  console.log('Canceling form, closing modal')
  showCreateModal.value = false
  setTimeout(() => {
    resetForm()
    isEditing.value = false
    editingUserId.value = null
    registrationSuccess.value = false
  }, 300)
}

// Reset form
const resetForm = () => {
  form.value = {
    fullName: '',
    email: '',
    password: '',
    phone: '',
    role: '',
    isActive: true,
    templeDetails: {
      name: '',
      place: '',
      address: '',
      phoneNumber: '',
      description: ''
    }
  }
  errors.value = {}
}
</script>

<style scoped>
.toggle-checkbox:checked {
  right: 0;
  border-color: #4f46e5;
}
.toggle-checkbox:checked + .toggle-label {
  background-color: #4f46e5;
}

/* Modal styles */
.fixed {
  position: fixed;
}
.inset-0 {
  top: 0;
  right: 0;
  bottom: 0;
  left: 0;
}
.z-50 {
  z-index: 50;
}
</style>