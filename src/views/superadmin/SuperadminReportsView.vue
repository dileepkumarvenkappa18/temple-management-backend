<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Reports Management</h1>
            <p class="text-gray-600 mt-1">
              Select tenants to generate cross-organization reports
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">Super Admin</span>
            </div>
          </div>
        </div>
      </div>
    </div>
    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Tenants Selection Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-xl font-bold text-gray-900">Select Tenants</h3>
          <p class="text-gray-600 mt-1">Choose which tenants to include in your report</p>
        </div>
        <div class="p-6">
          <div class="mb-6">
            <div class="flex items-center justify-between mb-4">
              <div>
                <h4 class="text-lg font-medium text-gray-900">Tenant Selection</h4>
                <p class="text-sm text-gray-500">Select multiple tenants for combined reporting</p>
              </div>
              <div class="flex space-x-3">
                <button 
                  @click="selectAllTenants" 
                  class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
                >
                  Select All
                </button>
                <button 
                  @click="clearTenantSelection"
                  class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
                >
                  Clear
                </button>
              </div>
            </div>
            <!-- Filters -->
            <div class="flex flex-wrap items-center gap-4 mb-4">
              <div class="relative flex-1 min-w-[250px]">
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
              
              <div class="relative">
                <select 
                  v-model="statusFilter" 
                  class="border border-gray-300 rounded-lg py-2 px-4 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white"
                >
                  <option value="all">All Status</option>
                  <option value="approved">Approved</option>
                  <option value="pending">Pending</option>
                  <option value="rejected">Rejected</option>
                </select>
              </div>
            </div>
            <!-- Tenant List -->
            <div v-if="loading" class="py-10 flex justify-center">
              <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-500"></div>
            </div>
            <div v-else-if="filteredTenants.length === 0" class="text-center py-8">
              <p class="text-gray-500">No tenants found matching your filters</p>
            </div>
            <div v-else class="overflow-x-auto border border-gray-200 rounded-lg">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                      <div class="flex items-center">
                        <input
                          type="checkbox"
                          :checked="allSelected"
                          @change="toggleSelectAll"
                          :indeterminate="someSelected && !allSelected"
                          class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                        />
                      </div>
                    </th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Tenant Name</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Temple Name</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Location</th>
                    <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
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
                    <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ getTenantDisplayName(tenant) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ getTempleNameDisplay(tenant) }}</td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ getLocationDisplay(tenant) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span class="px-2 py-1 inline-flex text-xs leading-5 font-semibold rounded-full" 
                            :class="getStatusClass(tenant.status)">
                        {{ tenant.status }}
                      </span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
          <!-- Report Type Selection -->
          <div class="mt-8 border-t border-gray-200 pt-6">
            <h4 class="text-lg font-medium text-gray-900 mb-4">Select Report Type</h4>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
              <div 
                v-for="report in reportTypes" 
                :key="report.id"
                @click="selectReport(report.id)"
                :class="[
                  'border p-4 rounded-lg cursor-pointer transition-colors',
                  selectedReport === report.id 
                    ? 'border-indigo-500 bg-indigo-50' 
                    : 'border-gray-200 hover:border-indigo-300 hover:bg-indigo-50/50'
                ]"
              >
                <div class="flex items-start">
                  <div class="flex-shrink-0">
                    <div class="h-10 w-10 flex items-center justify-center rounded-lg" :class="report.bgColor">
                      <svg class="h-6 w-6" :class="report.iconColor" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="report.icon" />
                      </svg>
                    </div>
                  </div>
                  <div class="ml-4">
                    <h5 class="text-base font-medium text-gray-900">{{ report.name }}</h5>
                    <p class="mt-1 text-sm text-gray-500">{{ report.description }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <!-- Proceed Button -->
          <div class="mt-8 flex justify-end">
            <button 
              @click="proceedToReport"
              :disabled="!canProceed"
              class="px-6 py-3 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
            >
              Proceed to Report
            </button>
          </div>
        </div>
      </div>
      <!-- Selection Summary -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Selection Summary</h3>
          
          <div class="flex flex-wrap gap-2 mb-4">
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Selected Tenants:</span>
              {{ selectedTenants.length }}
            </div>
            
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Report Type:</span>
              {{ getReportName() }}
            </div>
          </div>
          
          <div v-if="selectedTenants.length > 0" class="mt-4">
            <h4 class="text-sm font-medium text-gray-700 mb-2">Selected Tenant Names:</h4>
            <div class="flex flex-wrap gap-2">
              <span 
                v-for="tenantId in selectedTenants" 
                :key="tenantId"
                class="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm bg-gray-100 text-gray-800"
              >
                {{ getTenantName(tenantId) }}
                <button 
                  @click="toggleSelect(tenantId)" 
                  class="ml-1.5 text-gray-500 hover:text-gray-700"
                >
                  <svg class="h-3.5 w-3.5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </span>
            </div>
          </div>
          
          <p class="mt-4 text-sm text-gray-600">
            Select tenants and a report type above, then click "Proceed to Report" to generate a combined report.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useSuperAdminStore } from '@/stores/superadmin';
import { useToast } from '@/composables/useToast';
import superAdminService from '@/services/superadmin.service';

// Router and stores
const router = useRouter();
const superAdminStore = useSuperAdminStore();
const toast = useToast();

// Reactive state
const loading = ref(false);
const searchQuery = ref('');
const statusFilter = ref('all');
const selectedTenants = ref([]);
const selectedReport = ref('');

// Report types (with "User Details" and "Approval Status" removed)
const reportTypes = [
  {
    id: 'temple-register',
    name: 'Temple Register Report',
    description: 'View temple registration data across all selected tenants',
    icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z',
    bgColor: 'bg-blue-100',
    iconColor: 'text-blue-600',
    route: '/superadmin/reports/temple-register'
  },
  {
    id: 'temple-activities',
    name: 'Temple Activities Report',
    description: 'Track events, seva, and other temple activities',
    icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2',
    bgColor: 'bg-green-100',
    iconColor: 'text-green-600',
    route: '/superadmin/reports/temple-activities'
  },
  {
    id: 'birthdays',
    name: 'Devotee Birthdays Report',
    description: 'View upcoming birthdays of devotees across temples',
    icon: 'M21 15.546c-.523 0-1.046.151-1.5.454a2.704 2.704 0 01-3 0 2.704 2.704 0 00-3 0 2.704 2.704 0 01-3 0 2.704 2.704 0 00-3 0 2.704 2.704 0 01-3 0 2.701 2.701 0 00-1.5-.454M9 6v2m3-2v2m3-2v2M9 3h.01M12 3h.01M15 3h.01M21 21v-7a2 2 0 00-2-2H5a2 2 0 00-2 2v7h18zm-3-9v-2a2 2 0 00-2-2H8a2 2 0 00-2 2v2h12z',
    bgColor: 'bg-purple-100',
    iconColor: 'text-purple-600',
    route: '/superadmin/reports/birthdays'
  }
];

// Temple data mapping for tenants
const tenantTempleMap = ref({});

// Computed properties
const filteredTenants = computed(() => {
  let tenants = superAdminStore.tenants;
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    tenants = tenants.filter(tenant => 
      tenant.name?.toLowerCase().includes(query) || 
      tenant.fullName?.toLowerCase().includes(query) ||
      getTempleNameDisplay(tenant).toLowerCase().includes(query) ||
      getLocationDisplay(tenant).toLowerCase().includes(query)
    );
  }
  if (statusFilter.value !== 'all') {
    tenants = tenants.filter(tenant => 
      tenant.status?.toLowerCase() === statusFilter.value.toLowerCase()
    );
  }
  return tenants;
});

const allSelected = computed(() => {
  return filteredTenants.value.length > 0 && 
         filteredTenants.value.every(tenant => selectedTenants.value.includes(tenant.id));
});
const someSelected = computed(() => {
  return selectedTenants.value.length > 0 && !allSelected.value;
});
const canProceed = computed(() => {
  return selectedTenants.value.length > 0 && selectedReport.value;
});

// Utility methods
const getTenantDisplayName = (tenant) => {
  return tenant.fullName || tenant.name || tenant.FullName || `Temple Admin #${tenant.id || tenant.ID}`;
};
const getTempleNameDisplay = (tenant) => {
  if (tenantTempleMap.value[tenant.id || tenant.ID]) {
    return tenantTempleMap.value[tenant.id || tenant.ID].name;
  }
  if (tenant.temple && tenant.temple.name) return tenant.temple.name;
  if (tenant.Temple && tenant.Temple.name) return tenant.Temple.name;
  if (tenant.Temple && tenant.Temple.Name) return tenant.Temple.Name;
  if (tenant.EntityID && tenant.EntityName) return tenant.EntityName;
  return getTenantDisplayName(tenant) + "'s Temple";
};
const getLocationDisplay = (tenant) => {
  if (tenantTempleMap.value[tenant.id || tenant.ID]) {
    const templeInfo = tenantTempleMap.value[tenant.id || tenant.ID];
    if (templeInfo.city && templeInfo.state) return `${templeInfo.city}, ${templeInfo.state}`;
  }
  if (tenant.temple) {
    if (tenant.temple.city && tenant.temple.state) return `${tenant.temple.city}, ${tenant.temple.state}`;
    if (tenant.temple.address) return tenant.temple.address;
  }
  if (tenant.Temple) {
    if (tenant.Temple.City && tenant.Temple.State) return `${tenant.Temple.City}, ${tenant.Temple.State}`;
    if (tenant.Temple.Address) return tenant.Temple.Address;
  }
  if (tenant.EntityCity && tenant.EntityState) return `${tenant.EntityCity}, ${tenant.EntityState}`;
  if (tenant.EntityAddress) return tenant.EntityAddress;
  return "Location not available";
};

// Fetch tenants and temple details
const fetchTenants = async () => {
  loading.value = true;
  try {
    const response = await superAdminService.getAvailableTenants();
    if (response && response.success && response.data && response.data.length > 0) {
      superAdminStore.tenants = response.data;
      superAdminStore.tenants.forEach(tenant => {
        if (tenant.temple) {
          tenantTempleMap.value[tenant.id] = {
            name: tenant.temple.name,
            city: tenant.temple.city || (tenant.temple.address ? tenant.temple.address.split(',')[0]?.trim() : ''),
            state: tenant.temple.state || (tenant.temple.address ? tenant.temple.address.split(',')[1]?.trim() : '')
          };
        }
      });
    } else {
      const fallbackResponse = await superAdminStore.fetchTenantsForReports();
      if (fallbackResponse && fallbackResponse.success) {
        await fetchTempleDetailsForTenants();
      } else {
        throw new Error('Failed to fetch tenants');
      }
    }
  } catch (error) {
    toast.error('Failed to load tenants with complete details');
    try {
      await superAdminStore.fetchTenants();
      if (superAdminStore.tenants.length > 0) {
        await fetchTempleDetailsForTenants();
      }
    } catch (fallbackError) {
      toast.error('Using limited tenant data for demonstration');
    }
  } finally {
    loading.value = false;
  }
};

const fetchTempleDetailsForTenants = async () => {
  for (const tenant of superAdminStore.tenants) {
    if (!tenant.temple && !tenantTempleMap.value[tenant.id || tenant.ID]) {
      try {
        const tenantDetailsResponse = await superAdminService.getTenantDetails(tenant.id || tenant.ID);
        if (tenantDetailsResponse.success && tenantDetailsResponse.data) {
          const tenantData = tenantDetailsResponse.data;
          const templeInfo = tenantData.temple || tenantData.Temple || 
                           (tenantDetailsResponse.temples && tenantDetailsResponse.temples.length > 0 
                            ? tenantDetailsResponse.temples[0] : null);
          if (templeInfo) {
            tenantTempleMap.value[tenant.id || tenant.ID] = {
              name: templeInfo.name || templeInfo.Name || '',
              city: templeInfo.city || templeInfo.City || 
                    (templeInfo.address ? templeInfo.address.split(',')[0]?.trim() : '') ||
                    (templeInfo.Address ? templeInfo.Address.split(',')[0]?.trim() : ''),
              state: templeInfo.state || templeInfo.State || 
                     (templeInfo.address ? templeInfo.address.split(',')[1]?.trim() : '') ||
                     (templeInfo.Address ? templeInfo.Address.split(',')[1]?.trim() : '')
            };
            tenant.temple = {
              name: tenantTempleMap.value[tenant.id || tenant.ID].name,
              city: tenantTempleMap.value[tenant.id || tenant.ID].city,
              state: tenantTempleMap.value[tenant.id || tenant.ID].state
            };
          }
        }
      } catch {}
    }
  }
};

const isSelected = (tenantId) => selectedTenants.value.includes(tenantId);
const toggleSelect = (tenantId) => {
  const index = selectedTenants.value.indexOf(tenantId);
  if (index === -1) selectedTenants.value.push(tenantId);
  else selectedTenants.value.splice(index, 1);
};

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

const selectAllTenants = () => {
  selectedTenants.value = filteredTenants.value.map(tenant => tenant.id);
};
const clearTenantSelection = () => { selectedTenants.value = []; };
const selectReport = (reportId) => { selectedReport.value = reportId; };

const getStatusClass = (status) => {
  const statusMap = {
    'approved': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'rejected': 'bg-red-100 text-red-800',
    'active': 'bg-green-100 text-green-800'
  };
  return statusMap[status?.toLowerCase()] || 'bg-gray-100 text-gray-800';
};
const getTenantName = (tenantId) => {
  const tenant = superAdminStore.tenants.find(t => t.id === tenantId);
  return tenant ? getTenantDisplayName(tenant) : `Tenant #${tenantId}`;
};
const getReportName = () => {
  const report = reportTypes.find(r => r.id === selectedReport.value);
  return report ? report.name : 'None Selected';
};
const proceedToReport = () => {
  if (!canProceed.value) return;
  const report = reportTypes.find(r => r.id === selectedReport.value);
  if (!report) return;
  const tenantsParam = selectedTenants.value.join(',');
  router.push({
    path: report.route,
    query: {
      tenants: tenantsParam,
      from: 'superadmin'
    }
  });
};

// Lifecycle hook
onMounted(async () => {
  await fetchTenants();
});
</script>
