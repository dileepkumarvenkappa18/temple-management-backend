<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Superadmin Reports</h1>
            <p class="text-gray-600 mt-1">
              Download temple reports across all tenants
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">{{ userStore.user?.name || 'System Administrator' }}</span>
              <span class="text-indigo-600 text-sm ml-2">(Superadmin)</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Error Display -->
      <div v-if="errorMessage" class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error</h3>
            <div class="mt-2 text-sm text-red-700">
              {{ errorMessage }}
            </div>
            <div class="mt-4">
              <button 
                @click="clearError"
                class="text-sm bg-red-100 text-red-800 px-3 py-1 rounded hover:bg-red-200"
              >
                Dismiss
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter & Download Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-xl font-bold text-gray-900">Generate Reports</h3>
          <p class="text-gray-600 mt-1">Select a tenant, temple, and report type</p>
        </div>

        <div class="p-6">
          <!-- Tenant and Temple Selection - Two-step process -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <!-- Step 1: Tenant Selection -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Select Tenant</label>
              <div class="relative">
                <select 
                  v-model="selectedTenantId" 
                  @change="handleTenantChange"
                  :disabled="isLoading"
                  class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                >
                  <option value="">Select a tenant</option>
                  <option v-for="tenant in tenants" :key="tenant.id" :value="tenant.id">
                    {{ tenant.name }}
                  </option>
                </select>
                <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                  <!-- Dropdown indicator -->
                </div>
              </div>
              <p v-if="loadingTenants" class="mt-1 text-sm text-gray-500">Loading tenants...</p>
            </div>

            <!-- Step 2: Temple Selection (Only enabled after tenant is selected) -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
              <div class="relative">
                <select 
                  v-model="selectedTempleId" 
                  :disabled="!selectedTenantId || isLoading || loadingTemples"
                  class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                >
                  <option value="">Select a temple</option>
                  <option value="all" v-if="temples.length > 0">All Temples</option>
                  <option v-for="temple in temples" :key="temple.id" :value="temple.id">
                    {{ temple.name }}
                  </option>
                </select>
                <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                  <!-- Dropdown indicator -->
                </div>
              </div>
              <p v-if="loadingTemples" class="mt-1 text-sm text-gray-500">Loading temples...</p>
              <p v-else-if="selectedTenantId && temples.length === 0 && !loadingTemples" class="mt-1 text-sm text-red-500">
                No temples found for this tenant
              </p>
            </div>
          </div>

          <!-- Report Type Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Report Type</label>
            <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
              <button 
                v-for="type in reportTypes" 
                :key="type.value"
                @click="selectReportType(type.value)"
                :disabled="!selectedTempleId || isLoading"
                :class="[
                  'px-4 py-3 rounded-lg text-sm font-medium transition-colors duration-200 flex items-center',
                  selectedReportType === type.value 
                    ? 'bg-indigo-600 text-white' 
                    : 'bg-gray-100 text-gray-700 hover:bg-gray-200',
                  (!selectedTempleId || isLoading) ? 'opacity-50 cursor-not-allowed' : ''
                ]"
              >
                <span class="mr-2">
                  <svg class="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor" v-html="type.icon"></svg>
                </span>
                {{ type.label }}
              </button>
            </div>
          </div>

          <!-- Date Range Selection -->
          <div v-if="selectedReportType" class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Date Range</label>
            <div class="flex flex-wrap gap-2">
              <button 
                v-for="filter in timeFilters" 
                :key="filter.value"
                @click="setActiveFilter(filter.value)"
                :disabled="isLoading"
                class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                :class="activeFilter === filter.value ? 
                  'bg-indigo-600 text-white' : 
                  'bg-gray-100 text-gray-700 hover:bg-gray-200'"
              >
                {{ filter.label }}
              </button>
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
                  :disabled="isLoading"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  :disabled="isLoading"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                />
              </div>
            </div>
          </div>

          <!-- Download Section -->
          <div class="border-t border-gray-200 pt-6">
            <div class="flex flex-col md:flex-row md:items-center md:justify-between">
              <div class="mb-4 md:mb-0">
                <h4 class="text-lg font-medium text-gray-900">Download Report</h4>
                <p class="text-sm text-gray-600">Select a format and click download</p>
              </div>
              <div class="flex items-center space-x-3">
                <!-- Format Selection -->
                <div class="relative">
                  <select 
                    v-model="selectedFormat" 
                    :disabled="isLoading"
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                </div>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="!canDownload || isLoading"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  {{ isLoading ? 'Downloading...' : 'Download' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Current Applied Filters -->
      <div v-if="selectedTenantId && selectedTempleId && selectedReportType" class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Applied Filters</h3>
          
          <div class="flex flex-wrap gap-2">
            <!-- Tenant Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Tenant:</span>
              {{ getTenantName(selectedTenantId) }}
            </div>
            
            <!-- Temple Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Temple:</span>
              {{ selectedTempleId === 'all' ? 'All Temples' : getTempleName(selectedTempleId) }}
            </div>
            
            <!-- Report Type Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Report:</span>
              {{ getReportTypeLabel(selectedReportType) }}
            </div>
            
            <!-- Date Range Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Period:</span>
              {{ getTimeFilterLabel(activeFilter) }}
              <span v-if="activeFilter === 'custom'">
                ({{ formatDate(startDate) }} - {{ formatDate(endDate) }})
              </span>
            </div>
            
            <!-- Format -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Format:</span>
              {{ getFormatLabel(selectedFormat) }}
            </div>
          </div>
          
          <p class="mt-4 text-sm text-gray-600">
            Your report will include data based on the filters above. Click Download to generate and download the report.
          </p>
        </div>
      </div>
    </div>

    <!-- Download Format Modal -->
    <DownloadFormatModal
      v-if="showDownloadModal"
      :show="showDownloadModal"
      :title="'Report'"
      v-model="selectedFormat"
      @close="showDownloadModal = false"
      @download="handleModalDownload"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useToast } from '@/composables/useToast';
import superadminService from '@/services/superadmin.service';
import ReportsService from '@/services/reports.service';
import DownloadFormatModal from '@/components/common/DownloadFormatModal.vue';

// Store and composables
const userStore = useAuthStore();
const { showToast } = useToast();

// Reactive state
const tenants = ref([]);
const temples = ref([]);
const selectedTenantId = ref('');
const selectedTempleId = ref('');
const selectedReportType = ref('');
const activeFilter = ref('monthly');
const selectedFormat = ref('pdf');
const startDate = ref('');
const endDate = ref('');
const isLoading = ref(false);
const loadingTenants = ref(false);
const loadingTemples = ref(false);
const errorMessage = ref('');
const showDownloadModal = ref(false);

// Initialize dates
const initializeDates = () => {
  const today = new Date();
  const monthAgo = new Date();
  monthAgo.setDate(today.getDate() - 30);
  
  endDate.value = today.toISOString().split('T')[0];
  startDate.value = monthAgo.toISOString().split('T')[0];
};

// Filter options
const reportTypes = [
  { 
    label: 'Temple Register', 
    value: 'temple-register',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />'
  },
  { 
    label: 'Temple Activities', 
    value: 'temple-activities',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />'
  },
  { 
    label: 'Devotee Birthdays', 
    value: 'devotee-birthdays',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />'
  },
  { 
    label: 'Donation Summary', 
    value: 'donations',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2z" />'
  },
  { 
    label: 'Seva Bookings', 
    value: 'sevas',
    icon: '<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 10h18M7 15h1m4 0h1m-7 4h12a3 3 0 003-3V8a3 3 0 00-3-3H6a3 3 0 00-3 3v8a3 3 0 003 3z" />'
  }
];

const timeFilters = [
  { label: 'Weekly', value: 'weekly' },
  { label: 'Monthly', value: 'monthly' },
  { label: 'Yearly', value: 'yearly' },
  { label: 'Custom Range', value: 'custom' },
];

const formats = [
  { label: 'PDF', value: 'pdf' },
  { label: 'CSV', value: 'csv' },
  { label: 'Excel', value: 'excel' },
];

// Computed properties
const canDownload = computed(() => {
  return selectedTenantId.value && 
         selectedTempleId.value && 
         selectedReportType.value && 
         selectedFormat.value &&
         (activeFilter.value !== 'custom' || (startDate.value && endDate.value));
});

// Methods
const loadTenants = async () => {
  try {
    loadingTenants.value = true;
    clearError();
    
    const response = await superadminService.getAllTenants();
    
    if (response.success && response.data) {
      tenants.value = response.data.map(tenant => ({
        id: tenant.ID || tenant.id,
        name: tenant.FullName || tenant.name || tenant.Name || `Tenant ${tenant.ID || tenant.id}`,
        email: tenant.Email || tenant.email,
        status: tenant.Status || tenant.status
      }));
    } else {
      errorMessage.value = 'Failed to load tenants. Please try again.';
    }
  } catch (error) {
    console.error('Error loading tenants:', error);
    errorMessage.value = error.message || 'Failed to load tenants. Please try again.';
  } finally {
    loadingTenants.value = false;
  }
};

const handleTenantChange = async () => {
  selectedTempleId.value = '';
  temples.value = [];
  
  if (!selectedTenantId.value) return;
  
  try {
    loadingTemples.value = true;
    clearError();
    
    const response = await superadminService.getTemplesByTenant(selectedTenantId.value);
    
    if (response.success && response.data) {
      temples.value = response.data.map(temple => ({
        id: temple.ID || temple.id,
        name: temple.Name || temple.name,
        status: temple.Status || temple.status
      }));
    } else {
      errorMessage.value = 'Failed to load temples. Please try again.';
    }
  } catch (error) {
    console.error('Error loading temples:', error);
    errorMessage.value = error.message || 'Failed to load temples. Please try again.';
  } finally {
    loadingTemples.value = false;
  }
};

const selectReportType = (type) => {
  selectedReportType.value = type;
  // Depending on report type, we might want to set specific default filters
  switch (type) {
    case 'devotee-birthdays':
      activeFilter.value = 'monthly';
      break;
    case 'temple-activities':
      activeFilter.value = 'monthly';
      break;
    case 'temple-register':
      activeFilter.value = 'yearly';
      break;
    default:
      activeFilter.value = 'monthly';
  }
  
  // Update date range based on the active filter
  updateDateRange();
};

const setActiveFilter = (filter) => {
  activeFilter.value = filter;
  updateDateRange();
};

const updateDateRange = () => {
  // Set appropriate date range based on filter
  const today = new Date();
  
  if (activeFilter.value === 'weekly') {
    // Past week
    const weekAgo = new Date();
    weekAgo.setDate(today.getDate() - 7);
    startDate.value = weekAgo.toISOString().split('T')[0];
    endDate.value = today.toISOString().split('T')[0];
  } else if (activeFilter.value === 'monthly') {
    // Past month
    const monthAgo = new Date();
    monthAgo.setDate(today.getDate() - 30);
    startDate.value = monthAgo.toISOString().split('T')[0];
    endDate.value = today.toISOString().split('T')[0];
  } else if (activeFilter.value === 'yearly') {
    // Current year
    const currentYear = today.getFullYear();
    startDate.value = new Date(currentYear, 0, 1).toISOString().split('T')[0]; // January 1st
    endDate.value = today.toISOString().split('T')[0];
  }
  // For custom, we leave the dates as they are
};

const clearError = () => {
  errorMessage.value = '';
};

const getTenantName = (tenantId) => {
  const tenant = tenants.value.find(t => t.id.toString() === tenantId.toString());
  return tenant ? tenant.name : 'Unknown Tenant';
};

const getTempleName = (templeId) => {
  if (templeId === 'all') return 'All Temples';
  const temple = temples.value.find(t => t.id.toString() === templeId.toString());
  return temple ? temple.name : 'Unknown Temple';
};

const getReportTypeLabel = (type) => {
  const found = reportTypes.find(t => t.value === type);
  return found ? found.label : 'Unknown Report';
};

const getTimeFilterLabel = (filter) => {
  const found = timeFilters.find(f => f.value === filter);
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
  const params = {
    entityId: selectedTempleId.value,
    format: selectedFormat.value,
    dateRange: activeFilter.value,
    tenantId: selectedTenantId.value // Add tenantId for superadmin reports
  };

  // Add custom date range if needed
  if (activeFilter.value === 'custom') {
    params.startDate = startDate.value;
    params.endDate = endDate.value;
  }

  return params;
};

const downloadReport = async () => {
  if (!canDownload.value || isLoading.value) return;

  try {
    isLoading.value = true;
    clearError();
    
    const params = buildReportParams();
    
    // Add specific parameters based on report type
    let result;
    
    switch (selectedReportType.value) {
      case 'temple-register':
        result = await ReportsService.downloadTempleRegisteredReport(params);
        break;
      case 'temple-activities':
        result = await ReportsService.downloadActivitiesReport({
          ...params, 
          type: 'events' // Default to events, but could be configurable
        });
        break;
      case 'devotee-birthdays':
        result = await ReportsService.downloadDevoteeBirthdaysReport(params);
        break;
      case 'donations':
        result = await ReportsService.downloadActivitiesReport({
          ...params, 
          type: 'donations'
        });
        break;
      case 'sevas':
        result = await ReportsService.downloadActivitiesReport({
          ...params, 
          type: 'sevas'
        });
        break;
      default:
        throw new Error('Invalid report type selected');
    }
    
    if (result.success) {
      showToast(`Report downloaded successfully in ${getFormatLabel(selectedFormat.value)} format`, 'success');
    } else {
      throw new Error(result.message || 'Download failed');
    }
    
  } catch (error) {
    console.error('Error downloading report:', error);
    errorMessage.value = error.message || 'Failed to download report. Please try again.';
    showToast(error.message || 'Failed to download report', 'error');
  } finally {
    isLoading.value = false;
  }
};

const handleModalDownload = (data) => {
  selectedFormat.value = data.format;
  downloadReport();
};

// Lifecycle hooks
onMounted(async () => {
  initializeDates();
  await loadTenants();
});
</script>