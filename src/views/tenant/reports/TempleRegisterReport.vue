//TempleRegisterReport.vue
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Temple Register Report</h1>
            <p class="text-gray-600 mt-1">
              Download registration data for your temples
              <span v-if="tenantId" class="text-indigo-600 font-medium"> (Tenant ID: {{ tenantId }})</span>
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">{{ userStore.user?.name || 'Tenant User' }}</span>
              <span class="text-indigo-600 text-sm ml-2">(Tenant)</span>
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
          <h3 class="text-xl font-bold text-gray-900">Temple Register</h3>
          <p class="text-gray-600 mt-1">Configure filters and download your temple registration data</p>
        </div>

        <div class="p-6">
          <!-- Temple Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
            <div class="relative">
              <select 
                v-model="selectedTemple" 
                @change="onTempleChange"
                class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="all">All Temples</option>
                <option v-for="temple in templeStore.temples" :key="temple.id" :value="temple.id">
                  {{ temple.name }}
                </option>
              </select>
            </div>
          </div>

          <!-- Filter Section -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <!-- Registration Date Filter -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Registration Date</label>
              <div class="flex flex-wrap gap-2">
                <button 
                  v-for="filter in timeFilters" 
                  :key="filter.value"
                  @click="setActiveFilter(filter.value)"
                  :disabled="isDownloading"
                  class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                  :class="activeFilter === filter.value ? 
                    'bg-indigo-600 text-white' : 
                    'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  {{ filter.label }}
                </button>
              </div>
            </div>

            <!-- Temple Status Filter -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Temple Status</label>
              <div class="flex flex-wrap gap-2">
                <button 
                  v-for="status in statusFilters" 
                  :key="status.value"
                  @click="setActiveStatus(status.value)"
                  :disabled="isDownloading"
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
                  :disabled="isDownloading"
                  :max="endDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  :disabled="isDownloading"
                  :min="startDate"
                  :max="today"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
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
                    :disabled="isDownloading"
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                </div>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="isDownloading || !isFormValid"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
                >
                  <svg v-if="isDownloading" class="animate-spin -ml-1 mr-2 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  {{ isDownloading ? 'Downloading...' : 'Download' }}
                </button>
              </div>
            </div>
          </div>

          <!-- Error Display -->
          <div v-if="errorMessage" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-md">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">Error</h3>
                <p class="mt-1 text-sm text-red-700">{{ errorMessage }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Current Applied Filters -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Applied Filters</h3>
          
          <div class="flex flex-wrap gap-2">
            <!-- Temple Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Temple:</span>
              {{ selectedTemple === 'all' ? 'All Temples' : getTempleName(selectedTemple) }}
            </div>
            
            <!-- Registration Date Filter -->
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
            Your report will include data based on the filters above. Click Download to generate and download the report.
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useTempleStore } from '@/stores/temple';
import { useAuthStore } from '@/stores/auth';
import { useToast } from '@/composables/useToast';
import ReportsService from '@/services/reports.service';

// Composables
const route = useRoute();
const templeStore = useTempleStore();
const userStore = useAuthStore();
const { showToast } = useToast();

// Reactive state
const selectedTemple = ref('all');
const activeFilter = ref('weekly');
const activeStatus = ref('all');
const selectedFormat = ref('pdf');
const startDate = ref('');
const endDate = ref('');
const isDownloading = ref(false);
const errorMessage = ref('');

// Initialize dates
const initializeDates = () => {
  const today = new Date();
  const weekAgo = new Date();
  weekAgo.setDate(today.getDate() - 7);
  
  endDate.value = today.toISOString().split('T')[0];
  startDate.value = weekAgo.toISOString().split('T')[0];
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

// Computed
const tenantId = computed(() => {
  return route.params.tenantId || userStore.user?.id || localStorage.getItem('current_tenant_id');
});

const today = computed(() => {
  return new Date().toISOString().split('T')[0];
});

const entityId = computed(() => {
  return selectedTemple.value === 'all' ? 'all' : selectedTemple.value;
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
  clearError();
  
  // Set appropriate date range based on filter
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
  clearError();
};

const onTempleChange = () => {
  clearError();
};

const clearError = () => {
  errorMessage.value = '';
};

const getTempleName = (templeId) => {
  if (templeId === 'all') return 'All Temples';
  const temple = templeStore.temples.find(t => t.id.toString() === templeId.toString());
  return temple ? temple.name : 'Unknown Temple';
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
  const params = {
    entityId: entityId.value,
    dateRange: activeFilter.value,
    format: selectedFormat.value
  };

  // Add status filter if not 'all'
  if (activeStatus.value !== 'all') {
    params.status = activeStatus.value;
  }

  // Add custom date range
  if (activeFilter.value === 'custom') {
    params.startDate = startDate.value;
    params.endDate = endDate.value;
  }

  return params;
};

const downloadReport = async () => {
  if (isDownloading.value || !isFormValid.value) return;

  clearError();
  isDownloading.value = true;

  try {
    const params = buildReportParams();
    
    // Validate parameters using the service validation
    const validation = ReportsService.validateReportParams({
      ...params,
      type: 'temple-registered'
    });

    if (!validation.isValid) {
      throw new Error(validation.errors.join(', '));
    }

    // Console logging (matching your existing format)
    console.log('Downloading report with the following parameters:');
    console.log('- Temple:', selectedTemple.value === 'all' ? 'All Temples' : getTempleName(selectedTemple.value));
    console.log('- Time filter:', getTimeFilterLabel(activeFilter.value));
    console.log('- Date range:', formatDate(startDate.value), 'to', formatDate(endDate.value));
    console.log('- Status:', getStatusFilterLabel(activeStatus.value));
    console.log('- Format:', getFormatLabel(selectedFormat.value));

    // Call the service method
    const result = await ReportsService.downloadTempleRegisteredReport(params);
    
    // Show success message
    showToast(`Report downloaded successfully: ${result.filename}`, 'success');
    
  } catch (error) {
    console.error('Download failed:', error);
    errorMessage.value = error.message || 'Failed to download report. Please try again.';
    showToast(`Download failed: ${error.message}`, 'error');
  } finally {
    isDownloading.value = false;
  }
};

// Lifecycle hooks
onMounted(async () => {
  // Initialize default dates
  initializeDates();
  
  // Fetch temples if not already loaded
  if (templeStore.temples.length === 0) {
    try {
      await templeStore.fetchTemples(tenantId.value);
    } catch (error) {
      console.error('Error loading temple data:', error);
      showToast('Failed to load temple data. Please try again.', 'error');
      errorMessage.value = 'Failed to load temple data. Some features may not work correctly.';
    }
  }
});
</script>
