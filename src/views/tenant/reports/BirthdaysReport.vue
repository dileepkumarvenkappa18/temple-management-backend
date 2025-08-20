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
            <h1 class="text-2xl font-bold text-gray-900">Devotee Reports</h1>
            <p class="text-gray-600 mt-1">
              Download devotee data for your temples
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
      <!-- Error Display -->
      <div v-if="reportsStore.error" class="mb-6 bg-red-50 border border-red-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error</h3>
            <div class="mt-2 text-sm text-red-700">
              {{ reportsStore.error }}
            </div>
            <div class="mt-4">
              <button 
                @click="reportsStore.clearError()"
                class="text-sm bg-red-100 text-red-800 px-3 py-1 rounded hover:bg-red-200"
              >
                Dismiss
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Report Type Selector -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Select Report Type</h3>
          <div class="flex space-x-4">
            <button 
              @click="activeReportType = 'birthdays'"
              :class="[
                'px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200',
                activeReportType === 'birthdays' 
                  ? 'bg-indigo-600 text-white' 
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              ]"
            >
              Devotee Birthdays
            </button>
            <button 
              @click="activeReportType = 'devotees'"
              :class="[
                'px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200',
                activeReportType === 'devotees' 
                  ? 'bg-indigo-600 text-white' 
                  : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
              ]"
            >
              Devotee List
            </button>
          </div>
        </div>
      </div>

      <!-- Filter & Download Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-xl font-bold text-gray-900">
            {{ activeReportType === 'birthdays' ? 'Devotee Birthdays' : 'Devotee List' }}
          </h3>
          <p class="text-gray-600 mt-1">
            Configure filters and download {{ activeReportType === 'birthdays' ? 'devotee birthday data' : 'devotee list data' }}
          </p>
        </div>

        <div class="p-6">
          <!-- Temple Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
            <div class="relative">
              <select 
                v-model="selectedTemple" 
                :disabled="reportsStore.loading || reportsStore.downloadLoading"
                class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
              >
                <option value="all">All Temples</option>
                <option v-for="temple in templeStore.temples" :key="temple.id" :value="temple.id">
                  {{ temple.name }}
                </option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700"></div>
            </div>
          </div>

          <!-- Filter Section based on report type -->
          <div v-if="activeReportType === 'birthdays'">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <!-- Birthday Date Range Filter -->
              <div>
                <label class="block text-gray-700 font-medium mb-2">Birthday Period</label>
                <div class="flex flex-wrap gap-2">
                  <button 
                    v-for="filter in timeFilters" 
                    :key="filter.value"
                    @click="setActiveFilter(filter.value)"
                    :disabled="reportsStore.loading || reportsStore.downloadLoading"
                    class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="activeFilter === filter.value ? 
                      'bg-indigo-600 text-white' : 
                      'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                  >
                    {{ filter.label }}
                  </button>
                </div>
              </div>
            </div>
          </div>

          <div v-else>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
              <!-- Join Date Range Filter -->
              <div>
                <label class="block text-gray-700 font-medium mb-2">Join Date Period</label>
                <div class="flex flex-wrap gap-2">
                  <button 
                    v-for="filter in timeFilters" 
                    :key="filter.value"
                    @click="setActiveFilter(filter.value)"
                    :disabled="reportsStore.loading || reportsStore.downloadLoading"
                    class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
                    :class="activeFilter === filter.value ? 
                      'bg-indigo-600 text-white' : 
                      'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                  >
                    {{ filter.label }}
                  </button>
                </div>
              </div>
              
              <!-- Devotee Status Filter -->
              <!-- <div>
                <label class="block text-gray-700 font-medium mb-2">Devotee Status</label>
                <div class="relative">
                  <select 
                    v-model="devoteeStatus" 
                    :disabled="reportsStore.loading || reportsStore.downloadLoading"
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                  >
                    <option value="all">All Statuses</option>
                    <option value="active">Active</option>
                    <option value="inactive">Inactive</option>
                    <option value="new">New Members</option>
                  </select>
                  <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700"></div>
                </div>
              </div> -->
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
                  :disabled="reportsStore.loading || reportsStore.downloadLoading"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  :disabled="reportsStore.loading || reportsStore.downloadLoading"
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
                    :disabled="reportsStore.loading || reportsStore.downloadLoading"
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 disabled:bg-gray-100 disabled:cursor-not-allowed"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                  <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700"></div>
                </div>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="reportsStore.loading || reportsStore.downloadLoading"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg v-if="reportsStore.downloadLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
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
        </div>
      </div>

      <!-- Current Applied Filters -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Applied Filters</h3>
          
          <div class="flex flex-wrap gap-2">
            <!-- Report Type Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Report Type:</span>
              {{ activeReportType === 'birthdays' ? 'Devotee Birthdays' : 'Devotee List' }}
            </div>
            
            <!-- Tenant Filter (only in superadmin view with multiple tenants) -->
            <div v-if="fromSuperadmin && tenantIds.length > 1" class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Tenants:</span>
              {{ tenantIds.length }} selected
            </div>
            
            <!-- Temple Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Temple:</span>
              {{ selectedTemple === 'all' ? 'All Temples' : getTempleName(selectedTemple) }}
            </div>
            
            <!-- Period Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">{{ activeReportType === 'birthdays' ? 'Birthday' : 'Join Date' }} Period:</span>
              {{ getTimeFilterLabel(activeFilter) }}
              <span v-if="activeFilter === 'custom'">
                ({{ formatDate(startDate) }} - {{ formatDate(endDate) }})
              </span>
            </div>

            <!-- Devotee Status Filter (only for Devotee List) -->
            <div v-if="activeReportType === 'devotees'" class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Status:</span>
              {{ getDevoteeStatusLabel(devoteeStatus) }}
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

      <!-- Loading State -->
      <div v-if="reportsStore.loading" class="mt-6 bg-white rounded-lg shadow-sm border border-gray-200 p-6">
        <div class="flex items-center justify-center">
          <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <span class="text-gray-600">Loading report data...</span>
        </div>
      </div>

      <!-- Report Preview (if available) -->
      <div v-if="reportsStore.hasReportData && reportsStore.reportPreview" class="mt-6 bg-white rounded-lg shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-lg font-medium text-gray-900">Report Preview</h3>
          <p class="text-sm text-gray-600 mt-1">
            Showing {{ reportsStore.reportPreview.totalRecords }} records
          </p>
        </div>
        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th v-for="column in reportsStore.reportPreview.columns" :key="column.key" 
                    class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                  {{ column.label }}
                </th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="(row, index) in reportsStore.reportPreview.data.slice(0, 10)" :key="index">
                <td v-for="column in reportsStore.reportPreview.columns" :key="column.key" 
                    class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ row[column.key] || '-' }}
                </td>
              </tr>
            </tbody>
          </table>
        </div>
        <div v-if="reportsStore.reportPreview.data.length > 10" class="p-4 border-t border-gray-200 text-center text-sm text-gray-600">
          Showing first 10 records. Download the full report to see all {{ reportsStore.reportPreview.totalRecords }} records.
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useTempleStore } from '@/stores/temple';
import { useAuthStore } from '@/stores/auth';
import { useReportsStore } from '@/stores/reports';
import { useToast } from '@/composables/useToast';

// Composables
const route = useRoute();
const router = useRouter();
const templeStore = useTempleStore();
const userStore = useAuthStore();
const reportsStore = useReportsStore();
const { showToast } = useToast();

// Reactive state
const activeReportType = ref('birthdays'); // 'birthdays' or 'devotees'
const selectedTemple = ref('all');
const activeFilter = ref('monthly');
const selectedFormat = ref('pdf');
const startDate = ref(new Date().toISOString().split('T')[0]);
const endDate = ref(new Date(new Date().setDate(new Date().getDate() + 30)).toISOString().split('T')[0]);
const devoteeStatus = ref('all'); // For Devotee List report

// Filter options
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

const devoteeStatusOptions = [
  { label: 'All Statuses', value: 'all' },
  { label: 'Active', value: 'active' },
  { label: 'Inactive', value: 'inactive' },
  { label: 'New Members', value: 'new' },
];

// Computed properties
const tenantId = computed(() => {
  return route.params.tenantId || userStore.user?.id || localStorage.getItem('current_tenant_id');
});

// Check for tenants parameter from superadmin
const fromSuperadmin = computed(() => route.query.from === 'superadmin');
const tenantIds = computed(() => {
  if (route.query.tenants) {
    return route.query.tenants.split(',');
  }
  return [tenantId.value]; // Default to current tenant only
});

// Methods
const setActiveFilter = (filter) => {
  activeFilter.value = filter;
  
  // Set appropriate date range based on filter
  const today = new Date();
  
  if (filter === 'weekly') {
    // Next 7 days
    startDate.value = new Date().toISOString().split('T')[0];
    const weekEnd = new Date();
    weekEnd.setDate(weekEnd.getDate() + 7);
    endDate.value = weekEnd.toISOString().split('T')[0];
  } else if (filter === 'monthly') {
    // Next 30 days
    startDate.value = new Date().toISOString().split('T')[0];
    const monthEnd = new Date();
    monthEnd.setDate(monthEnd.getDate() + 30);
    endDate.value = monthEnd.toISOString().split('T')[0];
  } else if (filter === 'yearly') {
    // Current year
    const currentYear = today.getFullYear();
    startDate.value = new Date(currentYear, 0, 1).toISOString().split('T')[0]; // January 1st
    endDate.value = new Date(currentYear, 11, 31).toISOString().split('T')[0]; // December 31st
  }
  
  // For custom, we leave the dates as they are
  
  // Automatically fetch preview when filter changes
  fetchPreview();
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

const getFormatLabel = (format) => {
  const found = formats.find(f => f.value === format);
  return found ? found.label : 'Unknown';
};

const getDevoteeStatusLabel = (status) => {
  const found = devoteeStatusOptions.find(s => s.value === status);
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
  // When in superadmin view with multiple tenants, use entityIds instead of entityId
  const params = fromSuperadmin.value && tenantIds.value.length > 1
    ? { 
        entityIds: tenantIds.value, 
        dateRange: activeFilter.value,
        startDate: startDate.value,
        endDate: endDate.value,
        format: selectedFormat.value
      }
    : {
        entityId: selectedTemple.value === 'all' ? 'all' : selectedTemple.value.toString(),
        dateRange: activeFilter.value,
        startDate: startDate.value,
        endDate: endDate.value,
        format: selectedFormat.value
      };
  
  // Add devotee status for Devotee List report
  if (activeReportType.value === 'devotees') {
    params.status = devoteeStatus.value;
  }
  
  return params;
};

const fetchPreview = async () => {
  try {
    const params = buildReportParams();
    delete params.format; // Don't include format for preview
    
    if (activeReportType.value === 'birthdays') {
      await reportsStore.getDevoteeBirthdaysPreview(params);
    } else {
      await reportsStore.getDevoteeListPreview(params);
    }
  } catch (error) {
    console.error('Error fetching preview:', error);
    // Error is already handled by the store
  }
};

const downloadReport = async () => {
  try {
    // Clear any previous errors
    reportsStore.clearError();
    
    // Validate required fields
    if (activeFilter.value === 'custom' && (!startDate.value || !endDate.value)) {
      showToast('Please select both start and end dates for custom range', 'error');
      return;
    }
    
    if (new Date(startDate.value) > new Date(endDate.value)) {
      showToast('Start date must be before end date', 'error');
      return;
    }

    const params = buildReportParams();
    
    console.log(`Downloading ${activeReportType.value} report with parameters:`, params);
    
    let result;
    if (activeReportType.value === 'birthdays') {
      result = await reportsStore.downloadDevoteeBirthdaysReport(params);
    } else {
      result = await reportsStore.downloadDevoteeListReport(params);
    }
    
    showToast(
      `${activeReportType.value === 'birthdays' ? 'Devotee Birthdays' : 'Devotee List'} Report downloaded successfully in ${getFormatLabel(selectedFormat.value)} format`, 
      'success'
    );
    
    console.log('Download completed:', result);
    
  } catch (error) {
    console.error('Error downloading report:', error);
    showToast(error.message || 'Failed to download report. Please try again.', 'error');
  }
};

// Watch for report type changes
watch(activeReportType, () => {
  // Clear any previous report data
  reportsStore.clearReportData();
  
  // Fetch new preview based on selected report type
  fetchPreview();
});

// Watch for filter changes
watch([selectedTemple, devoteeStatus], () => {
  fetchPreview();
});

// Lifecycle hooks
onMounted(async () => {
  // Clear any previous report data
  reportsStore.clearReportData();
  
  // Fetch temples if not already loaded
  if (templeStore.temples.length === 0) {
    try {
      await templeStore.fetchTemples(tenantId.value);
    } catch (error) {
      console.error('Error loading temple data:', error);
      showToast('Failed to load temple data. Please try again.', 'error');
    }
  }
  
  // Fetch initial preview
  await fetchPreview();
});
</script>