<template>
  <div class="bg-white rounded-xl shadow-sm border border-gray-200">
    <div class="p-6 border-b border-gray-200">
      <h3 class="text-lg font-semibold text-gray-900 font-roboto">User Role Management</h3>
      <p class="text-sm text-gray-600 mt-1">Create and manage user roles</p>
    </div>
    
    <div class="p-6">
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
        <!-- Create Role Form -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <h4 class="text-lg font-medium text-gray-900 mb-4">{{ isEditing ? 'Edit Role' : 'Create New Role' }}</h4>
          
          <form @submit.prevent="handleSaveRole" class="space-y-4">
            <!-- Role Name Field -->
            <div>
              <label for="roleName" class="block text-sm font-medium text-gray-700 mb-2">
                Role Name <span class="text-red-500">*</span>
              </label>
              <BaseInput
                id="roleName"
                v-model="roleForm.role_name"
                type="text"
                placeholder="Enter role name"
                :error="roleErrors.role_name"
                required
              />
            </div>

            <!-- Role Description Field -->
            <div>
              <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
                Description <span class="text-red-500">*</span>
              </label>
              <BaseInput
                id="description"
                v-model="roleForm.description"
                type="text"
                placeholder="Enter role description"
                :error="roleErrors.description"
                required
              />
            </div>

            <!-- Public Registration Toggle -->
            <div class="flex items-center justify-between py-2">
              <label for="canRegister" class="text-sm font-medium text-gray-700">
                Allow Public Registration
              </label>
              <div class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in">
                <input 
                  type="checkbox" 
                  id="canRegister" 
                  v-model="roleForm.can_register_publicly" 
                  class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
                />
                <label 
                  for="canRegister" 
                  class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
                ></label>
              </div>
            </div>

            <!-- Role Status Toggle (for editing) -->
            <div v-if="isEditing" class="flex items-center justify-between py-2">
    <label for="roleStatus" class="text-sm font-medium text-gray-700">
        Role Status
    </label>
    <div class="relative inline-block w-10 mr-2 align-middle select-none transition duration-200 ease-in">
        <input 
            type="checkbox" 
            id="roleStatus" 
            v-model="roleStatusComputed" 
            class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
        />
        <label for="roleStatus" 
            class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
        ></label>
    </div>
</div>

            <!-- Form Actions -->
            <div class="flex justify-end space-x-3 pt-3">
              <BaseButton 
                v-if="isEditing" 
                type="button" 
                variant="secondary" 
                @click="cancelEdit"
              >
                Cancel
              </BaseButton>
              <BaseButton 
                type="submit" 
                variant="primary" 
                :loading="isSaving"
              >
                {{ isEditing ? 'Update Role' : 'Create Role' }}
              </BaseButton>
            </div>
          </form>
        </div>

        <!-- Roles List -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-xl shadow-sm border border-gray-200">
            <div class="p-4 border-b border-gray-200 flex justify-between items-center">
              <h4 class="font-medium text-gray-900">Available Roles</h4>
              <span class="text-xs px-2 py-1 bg-indigo-100 text-indigo-800 rounded-full">
                {{ roles.length }} Roles
              </span>
            </div>
            
            <div class="overflow-hidden">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Role Name
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Description
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Public Reg.
                    </th>
                    <!-- New Status Column -->
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Status
                    </th>
                    <th scope="col" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                      Actions
                    </th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-if="isLoading" class="animate-pulse">
                    <td colspan="5" class="px-6 py-4">
                      <div class="h-4 bg-gray-200 rounded w-3/4"></div>
                    </td>
                  </tr>
                  <tr v-else-if="roles.length === 0">
                    <td colspan="5" class="px-6 py-4 text-center text-sm text-gray-500">
                      No roles found. Create your first role.
                    </td>
                  </tr>
                  <tr v-for="role in roles" :key="role.id" class="hover:bg-gray-50">
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                      {{ role.role_name }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ role.description }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      <span 
                        class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                        :class="role.can_register_publicly ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'"
                      >
                        {{ role.can_register_publicly ? 'Yes' : 'No' }}
                      </span>
                    </td>
                    <!-- Status Toggle Column -->
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
    <div class="relative inline-block w-10 mr-2 align-middle select-none">
        <input 
            type="checkbox" 
            :id="`toggle-${role.id}`" 
            :checked="role.status === 'active'" 
            @change="toggleRoleStatus(role)"
            class="toggle-checkbox absolute block w-6 h-6 rounded-full bg-white border-4 appearance-none cursor-pointer"
            :disabled="['superadmin', 'templeadmin', 'devotee', 'volunteer'].includes(role.role_name)"
        />
        <label 
            :for="`toggle-${role.id}`" 
            class="toggle-label block overflow-hidden h-6 rounded-full bg-gray-300 cursor-pointer"
        ></label>
    </div>
    <span class="text-xs ml-1">{{ role.status === 'active' ? 'Active' : 'Inactive' }}</span>
</td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div class="flex justify-end space-x-2">
                        <button 
                          @click="editRole(role)" 
                          class="text-indigo-600 hover:text-indigo-900 focus:outline-none"
                        >
                          Edit
                        </button>
                        <button 
                          @click="confirmDeleteRole(role.id)" 
                          class="text-red-600 hover:text-red-900 focus:outline-none"
                          :disabled="['superadmin', 'templeadmin', 'devotee', 'volunteer'].includes(role.role_name)"
                        >
                          Delete
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Delete Confirmation Modal -->
    <BaseModal v-if="showDeleteModal" @close="showDeleteModal = false">
      <template #header>
        <h3 class="text-lg font-medium text-gray-900">Confirm Delete</h3>
      </template>
      <template #default>
        <div class="py-4">
          <p class="text-gray-600">Are you sure you want to delete this role? This action cannot be undone.</p>
        </div>
      </template>
      <template #footer>
        <div class="flex justify-end space-x-3">
          <BaseButton @click="showDeleteModal = false" variant="secondary">
            Cancel
          </BaseButton>
          <BaseButton @click="deleteRole" variant="danger" :loading="isDeleting">
            Delete
          </BaseButton>
        </div>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import superAdminService from '@/services/superadmin.service.js'

const { success, error } = useToast()

// State
const roles = ref([])
const isLoading = ref(false)
const isSaving = ref(false)
const isDeleting = ref(false)
const isEditing = ref(false)
const editingRoleId = ref(null)
const showDeleteModal = ref(false)
const roleToDeleteId = ref(null)

// Form state
// ...
const roleForm = ref({
    role_name: '',
    description: '',
    can_register_publicly: true,
    status: 'active' // 🎯 FIX: Changed 'is_active' to 'status' with a string value
});
// ...

const roleErrors = ref({})

// Fetch roles on component mount
onMounted(async () => {
  await fetchRoles()
})

// File: RoleManagement.vue
const fetchRoles = async () => {
    try {
        const response = await superAdminService.getRoles();
        console.log('Service: Roles response:', response);

        if (response.success && Array.isArray(response.data)) {
            // 🎯 FIX: Directly use the data returned from the backend
            // The backend is now returning the 'status' field, so no mapping is needed.
            roles.value = response.data;
        } else {
            console.error('API Error fetching roles:', response.message);
            // You can add a toast or other user-facing error message here.
            error(response.message || 'Failed to fetch roles.');
        }
    } catch (err) {
        console.error('Error fetching roles:', err);
        error('A network error occurred. Please try again.');
    }
};

const validateRoleForm = () => {
  roleErrors.value = {}
  
  if (!roleForm.value.role_name.trim()) {
    roleErrors.value.role_name = 'Role name is required'
  }
  
  if (!roleForm.value.description.trim()) {
    roleErrors.value.description = 'Description is required'
  }
  
  return Object.keys(roleErrors.value).length === 0
}

const handleSaveRole = async () => {
  if (!validateRoleForm()) return
  
  isSaving.value = true
  
  try {
    if (isEditing.value) {
      // 🔄 REPLACED: apiClient.put(...)
      const response = await superAdminService.updateRole(editingRoleId.value, roleForm.value)
      if (response.success) {
        success('Role updated successfully')
      } else {
        throw new Error(response.message)
      }
    } else {
      // 🔄 REPLACED: apiClient.post(...)
      const response = await superAdminService.createRole(roleForm.value)
      if (response.success) {
        success('Role created successfully')
      } else {
        throw new Error(response.message)
      }
    }
    
    resetForm()
    await fetchRoles()
  } catch (err) {
    console.error('Error saving role:', err)
    error(err.message || 'Failed to save role')
  } finally {
    isSaving.value = false
  }
}

// In your RoleManagement.vue or similar component
const toggleRoleStatus = async (role) => {
    try {
        // Get the current status from the role object
        const currentStatus = role.status;
        // Determine the new status
        const newStatus = currentStatus === 'active' ? 'inactive' : 'active';

        // Send the new status to the backend in the correct format
        const response = await superAdminService.updateRole(role.id, { status: newStatus });

        if (response.success) {
            // Update the local role object with the new status
            role.status = newStatus;
            success(`Role status updated to ${newStatus}`);
        } else {
            throw new Error(response.message);
        }
    } catch (err) {
        console.error('Error toggling role status:', err);
        error('Failed to update role status');
    }
};
const editRole = (role) => {
  isEditing.value = true;
  editingRoleId.value = role.id;
  
  // Clone the role data to the form, using the new 'status' field
  roleForm.value = {
    role_name: role.role_name,
    description: role.description,
    can_register_publicly: role.can_register_publicly,
    status: role.status // 🎯 FIX: Use the 'status' field directly
  };
};

const cancelEdit = () => {
  isEditing.value = false
  editingRoleId.value = null
  resetForm()
}

const confirmDeleteRole = (id) => {
  roleToDeleteId.value = id
  showDeleteModal.value = true
}

const deleteRole = async () => {
  if (!roleToDeleteId.value) return
  
  isDeleting.value = true
  
  try {
    // 🔄 REPLACED: apiClient.delete(...)
    const response = await superAdminService.deleteRole(roleToDeleteId.value)
    if (response.success) {
      success('Role deleted successfully')
    } else {
      throw new Error(response.message)
    }

    showDeleteModal.value = false
    roleToDeleteId.value = null
    await fetchRoles()
  } catch (err) {
    console.error('Error deleting role:', err)
    error(err.message || 'Failed to delete role')
  } finally {
    isDeleting.value = false
  }
}


const resetForm = () => {
  roleForm.value = {
    role_name: '',
    description: '',
    can_register_publicly: true,
    is_active: true // Reset is_active field
  }
  roleErrors.value = {}
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
</style>
