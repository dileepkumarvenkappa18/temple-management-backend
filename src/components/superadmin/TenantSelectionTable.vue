<template>
  <div class="tenant-selection-container">
    <div class="mb-4 flex items-center justify-between">
      <h3 class="text-lg font-medium">Select Tenants to Assign</h3>
      <div class="relative w-64">
        <input
          v-model="searchQuery"
          placeholder="Search tenants..."
          type="search"
          class="w-full border border-gray-300 rounded-lg py-2 px-4 pl-10 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg class="h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
      </div>
    </div>

    <div v-if="loading" class="py-10 flex justify-center">
      <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
    </div>

    <div v-else-if="filteredTenants.length === 0" class="text-center py-8">
      <p class="text-gray-500">No tenants found</p>
    </div>

    <div v-else class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              <div class="flex items-center">
                <input
                  type="checkbox"
                  :checked="allSelected"
                  @change="toggleSelectAll"
                  :indeterminate="someSelected"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
              </div>
            </th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User ID</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tenant Name</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Temple Address</th>
            <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Temple Name</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="tenant in filteredTenants" :key="tenant.id" :class="{ 'bg-indigo-50': isSelected(tenant.id) }">
            <td class="px-6 py-4 whitespace-nowrap">
              <input
                type="checkbox"
                :checked="isSelected(tenant.id)"
                @change="toggleSelect(tenant.id)"
                class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ tenant.userId }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ tenant.name }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ tenant.temple.address }}</td>
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ tenant.temple.name }}</td>
          </tr>
        </tbody>
      </table>
    </div>

    <div class="mt-6 flex justify-end space-x-3">
      <button 
        @click="$emit('cancel')" 
        class="px-4 py-2 bg-gray-200 text-gray-800 rounded-lg hover:bg-gray-300 transition-colors"
      >
        Cancel
      </button>
      <button 
        @click="$emit('assign', selectedTenants)"
        :disabled="selectedTenants.length === 0"
        class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
      >
        Assign Selected ({{ selectedTenants.length }})
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue';
import superAdminService from '@/services/superadmin.service';

export default {
  name: 'TenantSelectionTable',
  props: {
    userId: {
      type: [String, Number],
      required: true
    }
  },
  emits: ['assign', 'cancel'],
  setup(props, { emit }) {
    const tenants = ref([]);
    const selectedTenants = ref([]);
    const loading = ref(true);
    const searchQuery = ref('');
    
    const fetchTenants = async () => {
      loading.value = true;
      try {
        const response = await superAdminService.getAvailableTenants(props.userId);
        if (response.success) {
          tenants.value = response.data;
        } else {
          console.error('Failed to fetch tenants:', response.message);
        }
      } catch (error) {
        console.error('Error in fetchTenants:', error);
      } finally {
        loading.value = false;
      }
    };
    
    onMounted(fetchTenants);
    
    const filteredTenants = computed(() => {
      if (!searchQuery.value) return tenants.value;
      
      const query = searchQuery.value.toLowerCase();
      return tenants.value.filter(tenant => 
        tenant.name.toLowerCase().includes(query) ||
        tenant.temple.name.toLowerCase().includes(query) ||
        tenant.temple.address.toLowerCase().includes(query) ||
        tenant.userId.toString().includes(query)
      );
    });
    
    const isSelected = (tenantId) => selectedTenants.value.includes(tenantId);
    
    const toggleSelect = (tenantId) => {
      if (isSelected(tenantId)) {
        selectedTenants.value = selectedTenants.value.filter(id => id !== tenantId);
      } else {
        selectedTenants.value.push(tenantId);
      }
    };
    
    const allSelected = computed(() => 
      filteredTenants.value.length > 0 && 
      filteredTenants.value.every(tenant => selectedTenants.value.includes(tenant.id))
    );
    
    const someSelected = computed(() => 
      !allSelected.value && 
      filteredTenants.value.some(tenant => selectedTenants.value.includes(tenant.id))
    );
    
    const toggleSelectAll = () => {
      if (allSelected.value) {
        selectedTenants.value = selectedTenants.value.filter(
          id => !filteredTenants.value.some(tenant => tenant.id === id)
        );
      } else {
        const newSelectedIds = filteredTenants.value
          .filter(tenant => !selectedTenants.value.includes(tenant.id))
          .map(tenant => tenant.id);
        
        selectedTenants.value = [...selectedTenants.value, ...newSelectedIds];
      }
    };
    
    return {
      tenants,
      filteredTenants,
      selectedTenants,
      loading,
      searchQuery,
      isSelected,
      toggleSelect,
      allSelected,
      someSelected,
      toggleSelectAll
    };
  }
}
</script>