<template>
  <div class="assign-tenants-view">
    <div class="mb-4">
      <div class="flex items-center">
        <a @click="goBack" class="cursor-pointer text-indigo-600 hover:text-indigo-800 mr-2">
          <svg class="w-5 h-5 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
          </svg>
        </a>
        <div>
          <h2 class="text-xl font-bold">Assign Tenants to User</h2>
          <!-- <div class="text-sm text-gray-500">
            <span>Dashboard</span>
            <span class="mx-2">/</span>
            <span>User Management</span>
            <span class="mx-2">/</span>
            <span>Assign Tenants</span>
          </div> -->
        </div>
      </div>
    </div>

    <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-6">
      <div class="p-6 border-b border-gray-200 flex justify-between items-center">
        <h2 class="text-lg font-semibold">Assign Tenants to User</h2>
        <button 
          @click="goBack" 
          class="px-3 py-1.5 border border-gray-300 rounded-lg text-gray-700 hover:bg-gray-50"
        >
          Back
        </button>
      </div>

      <div class="p-6">
        <div v-if="loading" class="py-10 flex justify-center">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
        </div>

        <div v-else-if="error" class="bg-red-50 border border-red-200 rounded p-4 mb-6">
          <p class="text-red-700">{{ error }}</p>
          <button @click="fetchUserDetails" class="text-red-600 mt-2 underline">
            Retry
          </button>
        </div>

        <div v-else>
          <div class="bg-indigo-50 border border-indigo-100 rounded p-4 mb-6">
            <h3 class="font-medium text-indigo-800 mb-2">User Details</h3>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div>
                <p class="text-gray-500 text-sm">Name</p>
                <p class="font-medium">{{ user.fullName || user.full_name }}</p>
              </div>
              <div>
                <p class="text-gray-500 text-sm">Email</p>
                <p class="font-medium">{{ user.email }}</p>
              </div>
              <div>
                <p class="text-gray-500 text-sm">Role</p>
                <span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full bg-indigo-100 text-indigo-800">
                  {{ userRole }}
                </span>
              </div>
            </div>
          </div>

          <tenant-selection-table 
            :userId="userId" 
            @assign="assignTenants" 
            @cancel="goBack" 
          />
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import superAdminService from '@/services/superadmin.service';
import TenantSelectionTable from '@/components/superadmin/TenantSelectionTable.vue';

export default {
  name: 'AssignTenantsView',
  components: {
    TenantSelectionTable
  },
  setup() {
    const router = useRouter();
    const route = useRoute();
    const userId = ref(route.params.userId);
    const user = ref({});
    const loading = ref(true);
    const error = ref('');

    const fetchUserDetails = async () => {
      loading.value = true;
      error.value = '';
      
      try {
        const response = await superAdminService.getUserById(userId.value);
        if (response.success) {
          user.value = response.data;
          
          // Validate user role
          const role = userRole.value.toLowerCase();
          if (!['standarduser', 'standard user', 'monitoringuser', 'monitoring user'].includes(role)) {
            error.value = 'Tenant assignment is only available for Standard User and Monitoring User roles.';
            setTimeout(() => {
              router.push({ name: 'SuperadminUserManagement' });
            }, 3000);
          }
        } else {
          error.value = response.message || 'Failed to load user details';
        }
      } catch (err) {
        error.value = 'An error occurred while loading user details. Please try again.';
        console.error(err);
      } finally {
        loading.value = false;
      }
    };

    const assignTenants = async (selectedTenantIds) => {
      try {
        const response = await superAdminService.assignTenantsToUser(userId.value, selectedTenantIds);
        if (response.success) {
          // Show success message using toast if available
          if (window.$toast) {
            window.$toast.success('Tenants assigned successfully');
          } else {
            alert('Tenants assigned successfully');
          }
          router.push({ name: 'SuperadminUserManagement' });
        } else {
          if (window.$toast) {
            window.$toast.error(response.message || 'Failed to assign tenants');
          } else {
            alert(response.message || 'Failed to assign tenants');
          }
        }
      } catch (err) {
        if (window.$toast) {
          window.$toast.error('An error occurred while assigning tenants');
        } else {
          alert('An error occurred while assigning tenants');
        }
        console.error(err);
      }
    };

    const goBack = () => {
      router.push({ name: 'SuperadminUserManagement' });
    };

    const userRole = computed(() => {
      if (!user.value) return '';
      
      if (typeof user.value.role === 'string') {
        return user.value.role;
      } else if (typeof user.value.role === 'object') {
        return user.value.role.role_name || user.value.role.roleName || '';
      }
      return '';
    });

    onMounted(fetchUserDetails);

    return {
      userId,
      user,
      loading,
      error,
      assignTenants,
      goBack,
      fetchUserDetails,
      userRole
    };
  }
}
</script>