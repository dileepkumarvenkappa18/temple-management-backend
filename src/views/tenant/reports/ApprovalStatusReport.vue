<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Back to SuperAdmin Reports button (when viewed from superadmin) -->
    <div v-if="fromSuperadmin" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4">
      <router-link 
        to="/superadmin/reports" 
        class="inline-flex items-center px-4 py-2 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors"
      >
        <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
        </svg>
        Back to SuperAdmin Reports
      </router-link>
    </div>
    
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Approval Status Report</h1>
            <p class="text-gray-600 mt-1">
              Track approval workflows and statuses across your organization
              <span v-if="fromSuperadmin && tenantIds.length > 1" class="text-indigo-600 font-medium">
                (Multiple Tenants Selected)
              </span>
              <span v-else-if="tenantId" class="text-indigo-600 font-medium">
                (Tenant ID: {{ tenantId }})
              </span>
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">{{ userStore.user?.name || 'Tenant User' }}</span>
              <span class="text-indigo-600 text-sm ml-2">{{ fromSuperadmin ? '(Super Admin)' : '(Tenant)' }}</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filter & Download Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-xl font-bold text-gray-900">Approval Status Report</h3>
          <p class="text-gray-600 mt-1">Configure filters and download your approval workflow data</p>
        </div>

        <div class="p-6">
          <!-- Filter Section -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <!-- Date Range Filter -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Date Range</label>
              <div class="flex flex-wrap gap-2">
                <button 
                  v-for="filter in timeFilters" 
                  :key="filter.value"
                  @click="setActiveFilter(filter.value)"
                  :disabled="reportsStore.downloadLoading"
                  class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="activeFilter === filter.value ? 
                    'bg-indigo-600 text-white' : 
                    'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  {{ filter.label }}
                </button>
              </div>
            </div>

            <!-- Approval Status Filter -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Approval Status</label>
              <div class="flex flex-wrap gap-2">
                <button 
                  v-for="status in statusFilters" 
                  :key="status.value"
                  @click="setActiveStatus(status.value)"
                  :disabled="reportsStore.downloadLoading"
                  class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="activeStatus === status.value ? 
                    'bg-indigo-600 text-white' : 
                    'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  {{ status.label }}
                </button>
              </div>
            </div>
          </div>

          <!-- Custom Date Range (shown only when custom date range is selected) -->
          <div v-if="activeFilter === 'custom'" class="mb-6 p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">Start Date</label>
                <input 
                  type="date" 
                  v-model="startDate"
                  :disabled="reportsStore.downloadLoading"
                  :max="endDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  :disabled="reportsStore.downloadLoading"
                  :min="startDate"
                  :max="today"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                />
              </div>
            </div>
          </div>

          <!-- Preview and Download Section -->
          <div class="border-t border-gray-200 pt-6">
            <div class="flex flex-col md:flex-row md:items-center md:justify-between mb-4">
              <div>
                <h4 class="text-lg font-medium text-gray-900">Report Actions</h4>
                <p class="text-sm text-gray-600">Preview or download your approval status report</p>
              </div>
              <div class="flex items-center space-x-3 mt-4 md:mt-0">
                <!-- Format Selection -->
                <div class="relative">
                  <select 
                    v-model="selectedFormat" 
                    :disabled="reportsStore.downloadLoading"
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                </div>

                <!-- Preview Button -->
                <button 
                  @click="previewReport"
                  :disabled="reportsStore.loading || !isFormValid"
                  class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  <svg v-if="reportsStore.loading" class="animate-spin -ml-1 mr-2 h-5 w-5 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                  {{ reportsStore.loading ? 'Loading...' : 'Preview' }}
                </button>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="reportsStore.downloadLoading || !isFormValid"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  <svg v-if="reportsStore.downloadLoading" class="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  {{ reportsStore.downloadLoading ? 'Downloading...' : 'Download' }}
                </button>
              </div>
            </div>
          </div>

          <!-- Error Display -->
          <div v-if="reportsStore.error" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Error</h3>
                <p class="mt-1 text-sm text-red-700">{{ reportsStore.error }}</p>
              </div>
              <div class="ml-auto pl-3">
                <button 
                  @click="reportsStore.clearError"
                  class="inline-flex text-red-400 hover:text-red-500"
                >
                  <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Preview Results Card -->
      <div v-if="reportsStore.reportPreview && reportsStore.hasReportData" class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="text-xl font-bold text-gray-900">Report Preview</h3>
              <p class="text-gray-600 mt-1">
                Showing {{ reportsStore.reportPreview.totalRecords }} approval records
              </p>
            </div>
            <div class="text-sm text-gray-500">
              Last updated: {{ new Date().toLocaleString() }}
            </div>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th v-for="column in reportsStore.reportPreview.columns" :key="column.key" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  {{ column.label }}
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="(item, index) in reportsStore.reportPreview.data.slice(0, 10)" :key="index" class="hover:bg-gray-50">
                <td v-for="column in reportsStore.reportPreview.columns" :key="column.key" class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ item[column.key] || '-' }}
                </td>
              </tr>
            </tbody>
          </table>
          
          <div v-if="reportsStore.reportPreview.data.length > 10" class="px-6 py-3 bg-gray-50 border-t border-gray-200">
            <p class="text-sm text-gray-600">
              Showing first 10 of {{ reportsStore.reportPreview.totalRecords }} records. Download the full report to see all data.
            </p>
          </div>
        </div>
      </div>

      <!-- Applied Filters Display -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Applied Filters</h3>
          
          <div class="flex flex-wrap gap-2">
            <!-- Tenant Filter (only in superadmin view with multiple tenants) -->
            <div v-if="fromSuperadmin && tenantIds.length > 1" class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Tenants:</span>
              {{ tenantIds.length }} selected
            </div>
            
            <!-- Date Range Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Period:</span>
              {{ getTimeFilterLabel(activeFilter) }}
              <span v-if="activeFilter === 'custom' && startDate && endDate">
                ({{ formatDate(startDate) }} - {{ formatDate(endDate) }})
              </span>
            </div>
            
            <!-- Status Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Status:</span>
              {{ getStatusFilterLabel(activeStatus) }}
            </div>
            
            <!-- Format -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Format:</span>
              {{ getFormatLabel(selectedFormat) }}
            </div>
          </div>
          
          <p class="mt-4 text-sm text-gray-600">
            Your report will include approval status data based on the filters above. Use Preview to see a sample of the data before downloading.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useReportsStore } from '@/stores/reports';
import { useToast } from '@/composables/useToast';

// Composables
const route = useRoute();
const router = useRouter();
const userStore = useAuthStore();
const reportsStore = useReportsStore();
const { showToast } = useToast();

// Reactive state
const activeFilter = ref('monthly');
const activeStatus = ref('all');
const selectedFormat = ref('pdf');
const startDate = ref('');
const endDate = ref('');

// Initialize dates
const initializeDates = () => {
  const today = new Date();
  const monthAgo = new Date();
  monthAgo.setDate(today.getDate() - 30);
  
  endDate.value = today.toISOString().split('T')[0];
  startDate.value = monthAgo.toISOString().split('T')[0];
};

// Filter options
const timeFilters = [
  { label: 'Weekly', value: 'weekly' },
  { label: 'Monthly', value: 'monthly' },
  { label: 'Yearly', value: 'yearly' },
  { label: 'Custom Range', value: 'custom' },
];

const statusFilters = [
  { label: 'All Status', value: 'all' },
  { label: 'Approved', value: 'approved' },
  { label: 'Pending', value: 'pending' },
  { label: 'Rejected', value: 'rejected' },
];

const formats = [
  { label: 'PDF', value: 'pdf' },
  { label: 'CSV', value: 'csv' },
  { label: 'Excel', value: 'excel' },
];

// Computed properties
const tenantId = computed(() => {
  return route.params.tenantId || userStore.user?.id || localStorage.getItem('current_tenant_id');
});

const fromSuperadmin = computed(() => route.query.from === 'superadmin');
const tenantIds = computed(() => {
  if (route.query.tenants) {
    return route.query.tenants.split(',');
  }
  return [tenantId.value];
});

const today = computed(() => {
  return new Date().toISOString().split('T')[0];
});

const isFormValid = computed(() => {
  if (activeFilter.value === 'custom') {
    return startDate.value && endDate.value && new Date(startDate.value) <= new Date(endDate.value);
  }
  return true;
});

// Methods
const setActiveFilter = (filter) => {
  activeFilter.value = filter;
  reportsStore.clearError();
  
  const today = new Date();
  const start = new Date();
  
  if (filter === 'weekly') {
    start.setDate(today.getDate() - 7);
  } else if (filter === 'monthly') {
    start.setDate(today.getDate() - 30);
  } else if (filter === 'yearly') {
    start.setDate(today.getDate() - 365);
  }
  
  if (filter !== 'custom') {
    startDate.value = start.toISOString().split('T')[0];
    endDate.value = today.toISOString().split('T')[0];
  }
};

const setActiveStatus = (status) => {
  activeStatus.value = status;
  reportsStore.clearError();
};

const getTimeFilterLabel = (filter) => {
  const found = timeFilters.find(f => f.value === filter);
  return found ? found.label : 'Unknown';
};

const getStatusFilterLabel = (status) => {
  const found = statusFilters.find(s => s.value === status);
  return found ? found.label : 'Unknown';
};

const getFormatLabel = (format) => {
  const found = formats.find(f => f.value === format);
  return found ? found.label : 'Unknown';
};

const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
};

const buildReportParams = () => {
  const baseParams = {
    dateRange: activeFilter.value,
    startDate: startDate.value,
    endDate: endDate.value,
  };

  // Add status filter if not 'all'
  if (activeStatus.value !== 'all') {
    baseParams.status = activeStatus.value;
  }

  // Handle multi-tenant vs single tenant
  if (fromSuperadmin.value && tenantIds.value.length > 1) {
    baseParams.entityIds = tenantIds.value;
    baseParams.isSuperAdmin = true;
  } else {
    baseParams.entityId = tenantId.value;
  }

  return baseParams;
};

const previewReport = async () => {
  if (!isFormValid.value) return;

  try {
    const params = buildReportParams();
    
    console.log('Previewing Approval Status report with parameters:', {
      tenants: fromSuperadmin.value && tenantIds.value.length > 1 ? tenantIds.value.length : 1,
      dateRange: getTimeFilterLabel(activeFilter.value),
      status: getStatusFilterLabel(activeStatus.value),
      period: `${formatDate(startDate.value)} - ${formatDate(endDate.value)}`
    });

    await reportsStore.getApprovalStatusPreview(params);
    
    if (reportsStore.hasReportData) {
      showToast(`Preview loaded: ${reportsStore.reportPreview.totalRecords} records found`, 'success');
    } else {
      showToast('No data found for the selected criteria', 'info');
    }
    
  } catch (error) {
    console.error('Preview failed:', error);
    showToast(`Preview failed: ${error.message}`, 'error');
  }
};

const downloadReport = async () => {
  if (!isFormValid.value) return;

  try {
    const params = {
      ...buildReportParams(),
      format: selectedFormat.value
    };
    
    console.log('Downloading Approval Status report:', {
      tenants: fromSuperadmin.value && tenantIds.value.length > 1 ? tenantIds.value.length : 1,
      dateRange: getTimeFilterLabel(activeFilter.value),
      status: getStatusFilterLabel(activeStatus.value),
      format: getFormatLabel(selectedFormat.value),
      period: `${formatDate(startDate.value)} - ${formatDate(endDate.value)}`
    });

    const result = await reportsStore.downloadApprovalStatusReport(params);
    showToast(`Report downloaded successfully: ${result.filename}`, 'success');
    
  } catch (error) {
    console.error('Download failed:', error);
    showToast(`Download failed: ${error.message}`, 'error');
  }
};

// Lifecycle hooks
onMounted(() => {
  initializeDates();
  reportsStore.clearReportData();
});
</script>