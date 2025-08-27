<template>
  <div class="tenant-selection-page">
    <!-- Header with user info -->
    <div class="bg-indigo-50 rounded-2xl shadow-md p-6 mb-6">
      <div class="flex flex-col md:flex-row md:items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-indigo-800 mb-2">Welcome, {{ userInfo.name }}</h1>
          <div class="flex flex-wrap gap-3 items-center">
            <span class="bg-indigo-600 text-white text-sm px-3 py-1 rounded-xl">
              {{ userInfo.role }}
            </span>
            <span class="text-indigo-700">{{ userInfo.email }}</span>
          </div>
        </div>
        
        <div class="mt-4 md:mt-0" v-if="isSuperAdmin">
          <p class="text-indigo-600 font-semibold">
            Super Admin Access
          </p>
        </div>
      </div>
    </div>

    <!-- Main content area -->
    <div class="bg-white rounded-2xl shadow-md p-6">
      <h2 class="text-xl font-bold text-indigo-700 mb-4">
        {{ isSuperAdmin ? 'Tenant Management' : 'Select a Tenant' }}
      </h2>
      
      <p class="text-gray-600 mb-6">
        {{ getSelectionInstructions }}
      </p>

      <div v-if="isSuperAdmin && !showTenantList" class="flex justify-center">
        <button 
          @click="loadTenants" 
          class="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-6 rounded-xl shadow-md transition-all duration-200 flex items-center"
        >
          <span>Show Tenants</span>
          <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v3.586L7.707 9.293a1 1 0 00-1.414 1.414l3 3a1 1 0 001.414 0l3-3a1 1 0 00-1.414-1.414L11 10.586V7z" clip-rule="evenodd" />
          </svg>
        </button>
      </div>

      <!-- Tenant list -->
      <div v-if="showTenantList">
        <!-- Search and filter -->
        <div class="mb-6">
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search tenants..."
              class="w-full pl-10 pr-4 py-2 border border-indigo-200 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            />
            <div class="absolute left-3 top-2.5 text-indigo-400">
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
              </svg>
            </div>
          </div>
        </div>

        <!-- Loading -->
        <div v-if="loading" class="flex justify-center py-12">
          <div class="animate-spin rounded-full h-12 w-12 border-b-2 border-indigo-600"></div>
        </div>

        <!-- Tenant list -->
        <div v-else-if="filteredTenants.length > 0" class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="tenant in filteredTenants" 
            :key="tenant.id"
            @click="selectTenant(tenant.id)"
            class="bg-white border rounded-xl overflow-hidden shadow-md hover:shadow-lg transition-all duration-200 cursor-pointer"
            :class="{'border-indigo-500 ring-2 ring-indigo-500': selectedTenantId === tenant.id, 'border-gray-200': selectedTenantId !== tenant.id}"
          >
            <div class="relative h-36 bg-indigo-100 overflow-hidden">
              <img 
                v-if="tenant.imageUrl" 
                :src="tenant.imageUrl" 
                :alt="tenant.name" 
                class="w-full h-full object-cover"
              />
              <div v-else class="w-full h-full flex items-center justify-center">
                <span class="text-indigo-800 text-5xl opacity-30">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                </span>
              </div>
              <div class="absolute top-2 right-2">
                <span 
                  class="px-2 py-1 text-xs font-semibold rounded-lg"
                  :class="{
                    'bg-green-100 text-green-800': tenant.status === 'active' || tenant.status === 'approved',
                    'bg-yellow-100 text-yellow-800': tenant.status === 'pending',
                    'bg-red-100 text-red-800': tenant.status === 'inactive'
                  }"
                >
                  {{ tenant.status }}
                </span>
              </div>
            </div>
            
            <div class="p-4">
              <h3 class="font-bold text-lg text-indigo-900 mb-1">{{ tenant.name }}</h3>
              <p class="text-gray-600 text-sm mb-3">{{ tenant.location || tenant.email }}</p>
              
              <div class="flex items-center text-gray-500 text-sm mb-3">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
                </svg>
                <span>{{ tenant.templesCount || 0 }} temples</span>
              </div>
              
              <button 
                @click.stop="selectTenant(tenant.id)"
                class="w-full py-2 px-4 flex justify-center items-center rounded-lg text-sm font-medium transition-all duration-200"
                :class="selectedTenantId === tenant.id ? 
                  'bg-indigo-600 text-white hover:bg-indigo-700' : 
                  'border border-indigo-500 text-indigo-600 hover:bg-indigo-50'"
              >
                {{ selectedTenantId === tenant.id ? 'Selected' : 'Select' }}
              </button>
            </div>
          </div>
        </div>

        <!-- Empty state -->
        <div v-else class="text-center py-8">
          <div class="text-indigo-300 mb-4">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-16 w-16 mx-auto" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M8 16l2.879-2.879m0 0a3 3 0 104.243-4.242 3 3 0 00-4.243 4.242zM21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
          <h3 class="text-lg font-medium text-gray-800 mb-1">No tenants found</h3>
          <p class="text-gray-500">Try adjusting your search or filters</p>
        </div>

        <!-- Proceed button -->
        <div v-if="selectedTenantId && filteredTenants.length > 0" class="mt-8 flex justify-center">
          <button 
            @click="proceedToTenantDashboard" 
            class="bg-indigo-600 hover:bg-indigo-700 text-white font-semibold py-3 px-8 rounded-xl shadow-md transition-all duration-200 flex items-center"
            id="proceed-button"
          >
            <span>Proceed to Dashboard</span>
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 ml-2" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10.293 5.293a1 1 0 011.414 0l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414-1.414L12.586 11H5a1 1 0 110-2h7.586l-2.293-2.293a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import axios from 'axios';
import { useToast } from '@/composables/useToast';
import tenantService from '@/services/tenant.service';

const authStore = useAuthStore();
const router = useRouter();
const { showToast } = useToast();

const userInfo = ref({
  id: authStore.user?.id || 1,
  name: authStore.user?.fullName || authStore.user?.name || 'User',
  email: authStore.user?.email || 'user@example.com',
  role: authStore.userRole || authStore.user?.role || 'standard_user',
});

const loading = ref(false);
const tenants = ref([]);
const selectedTenantId = ref(null);
const searchQuery = ref('');

const isSuperAdmin = computed(() => {
  const role = userInfo.value.role?.toLowerCase() || '';
  return role === 'superadmin' || role === 'super_admin';
});

const isMonitoringUser = computed(() => {
  const role = userInfo.value.role?.toLowerCase() || '';
  return role === 'monitoringuser' || role === 'monitoring_user';
});

const showTenantList = ref(!isSuperAdmin.value);

const getSelectionInstructions = computed(() => {
  if (isSuperAdmin.value) {
    return 'As a Super Admin, you can access and manage all tenant dashboards.';
  } else if (isMonitoringUser.value) {
    return 'Select a tenant to view their dashboard. As a Monitoring User, you have limited access to view data.';
  } else {
    return 'Select a tenant to access their dashboard and management features.';
  }
});

const filteredTenants = computed(() => {
  if (!searchQuery.value) return tenants.value;
  
  const query = searchQuery.value.toLowerCase();
  return tenants.value.filter(tenant => 
    tenant.name?.toLowerCase().includes(query) || 
    tenant.location?.toLowerCase().includes(query) ||
    tenant.email?.toLowerCase().includes(query)
  );
});

const loadTenants = async () => {
  loading.value = true;
  showTenantList.value = true;
  try {
    const tenantsData = await tenantService.getTenantsForSelection();
    if (Array.isArray(tenantsData) && tenantsData.length > 0) {
      tenants.value = tenantsData;
    } else {
      tenants.value = [];
      showToast('No tenants available for selection', 'warning');
    }
  } catch (error) {
    console.error('Failed to load tenants:', error);
    showToast('Using development data while API is being set up', 'info');
  } finally {
    loading.value = false;
  }
};

const selectTenant = (tenantId) => {
  selectedTenantId.value = Number(tenantId);
};

const proceedToTenantDashboard = () => {
  if (!selectedTenantId.value) {
    showToast('Please select a tenant to proceed', 'warning');
    return;
  }

  try {
    localStorage.setItem('selected_tenant_id', selectedTenantId.value);
    localStorage.setItem('current_tenant_id', selectedTenantId.value);
    localStorage.setItem('current_entity_id', selectedTenantId.value);

    const token = localStorage.getItem('auth_token');
    if (token) {
      axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
      axios.defaults.headers.common['X-Tenant-ID'] = selectedTenantId.value;
    }

    const selectedTenant = tenants.value.find(t => Number(t.id) === Number(selectedTenantId.value));
    if (selectedTenant) {
      localStorage.setItem('selected_tenant_name', selectedTenant.name);
    }

    // ✅ Redirect to tenant dashboard
    const redirectPath = `/tenant/${selectedTenantId.value}/dashboard`;
    window.location.href = redirectPath;
  } catch (error) {
    console.error('Navigation error:', error);
    showToast('Failed to navigate to dashboard. Please try again.', 'error');
  }
};

const directNavigate = (tenantId) => {
  localStorage.setItem('selected_tenant_id', tenantId);
  localStorage.setItem('current_tenant_id', tenantId);
  localStorage.setItem('current_entity_id', tenantId);

  const selectedTenant = tenants.value.find(t => Number(t.id) === Number(tenantId));
  if (selectedTenant) {
    localStorage.setItem('selected_tenant_name', selectedTenant.name);
  }

  // ✅ Redirect to tenant dashboard
  const path = `/tenant/${tenantId}/dashboard`;
  window.location.href = path;
};

onMounted(() => {
  if (!isSuperAdmin.value) {
    loadTenants();
  }
  const savedTenantId = localStorage.getItem('selected_tenant_id');
  if (savedTenantId) {
    selectedTenantId.value = Number(savedTenantId);
  }
  const token = localStorage.getItem('auth_token');
  if (token) {
    axios.defaults.headers.common['Authorization'] = `Bearer ${token}`;
  }
});
</script>

<style scoped>
.tenant-selection-page {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

@media (max-width: 768px) {
  .tenant-selection-page {
    padding: 1rem;
  }
}
</style>
