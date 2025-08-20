//TempleActivities.vue

<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Temple Activities Report</h1>
            <p class="text-gray-600 mt-1">
              Download event, seva, and booking activity data for your temples
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

    <!-- Loading Overlay -->
    <div v-if="reportsStore.loading" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
        <span class="text-gray-900 font-medium">Loading report data...</span>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="reportsStore.error" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error</h3>
            <p class="mt-1 text-sm text-red-700">{{ reportsStore.error }}</p>
          </div>
          <div class="ml-auto pl-3">
            <button @click="reportsStore.clearError" class="text-red-400 hover:text-red-500">
              <span class="sr-only">Dismiss</span>
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Filter & Download Card -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6 border-b border-gray-200">
          <h3 class="text-xl font-bold text-gray-900">Temple Activities</h3>
          <p class="text-gray-600 mt-1">Configure filters and download your temple activity data</p>
        </div>

        <div class="p-6">
          <!-- Temple Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
            <div class="relative">
              <select 
                v-model="selectedTemple" 
                class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                @change="handleTempleChange"
              >
                <option value="all">All Temples</option>
                <option v-for="temple in templeStore.temples" :key="temple.id" :value="temple.id">
                  {{ temple.name }}
                </option>
              </select>
            </div>
          </div>

          <!-- Activity Type Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Activity Type</label>
            <div class="flex flex-wrap gap-2">
              <button 
                v-for="type in activityTypes" 
                :key="type.value"
                @click="setActivityType(type.value)"
                class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                :class="activityType === type.value ? 
                  'bg-indigo-600 text-white' : 
                  'bg-gray-100 text-gray-700 hover:bg-gray-200'"
              >
                {{ type.label }}
              </button>
            </div>
          </div>

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
                  class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                  :class="activeFilter === filter.value ? 
                    'bg-indigo-600 text-white' : 
                    'bg-gray-100 text-gray-700 hover:bg-gray-200'"
                >
                  {{ filter.label }}
                </button>
              </div>
            </div>
          </div>

          <!-- Custom Date Range -->
          <div v-if="activeFilter === 'custom'" class="mb-6 p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">Start Date</label>
                <input 
                  type="date" 
                  v-model="startDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  @change="handleDateChange"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  @change="handleDateChange"
                />
              </div>
            </div>
          </div>

          <!-- Preview Section -->
          <!-- <div class="mb-6">
            <div class="flex items-center justify-between mb-4">
              <h4 class="text-lg font-medium text-gray-900">Report Preview</h4>
              <button 
                @click="loadReportPreview"
                :disabled="reportsStore.loading"
                class="inline-flex items-center px-3 py-1 border border-transparent text-sm font-medium rounded-md text-indigo-600 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50"
              >
                <svg class="mr-1 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Refresh Preview
              </button>
            </div>
            
            Preview Data
            <div v-if="reportsStore.hasReportData" class="bg-gray-50 border border-gray-200 rounded-lg p-4">
              <div class="flex items-center justify-between mb-3">
                <span class="text-sm font-medium text-gray-900">
                  {{ reportsStore.reportSummary?.totalRecords || 0 }} records found
                </span>
                <span class="text-xs text-gray-600">
                  {{ getActivityTypeLabel(activityType) }} | {{ getTimeFilterLabel(activeFilter) }}
                </span>
              </div>
              
              Preview Table
              <div class="overflow-x-auto">
                <table class="min-w-full divide-y divide-gray-200">
                  <thead class="bg-gray-100">
                    <tr>
                      <th v-for="column in reportsStore.reportPreview?.columns" :key="column.key" 
                          class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                        {{ column.label }}
                      </th>
                    </tr>
                  </thead>
                  <tbody class="bg-white divide-y divide-gray-200">
                    <tr v-for="(row, index) in previewRows" :key="index">
                      <td v-for="column in reportsStore.reportPreview?.columns" :key="column.key"
                          class="px-3 py-2 whitespace-nowrap text-xs text-gray-900">
                        {{ row[column.key] || '-' }}
                      </td>
                    </tr>
                  </tbody>
                </table>
              </div>
              
              <p v-if="reportsStore.reportPreview?.data.length > 5" class="text-xs text-gray-500 mt-2 text-center">
                Showing first 5 rows. Download full report to see all {{ reportsStore.reportSummary?.totalRecords }} records.
              </p>
            </div>
            
            <div v-else-if="!reportsStore.loading" class="bg-gray-50 border border-gray-200 rounded-lg p-8 text-center">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <p class="mt-2 text-sm text-gray-600">No data available. Click "Refresh Preview" to load data.</p>
            </div>
          </div> -->

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
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                </div>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="reportsStore.downloadLoading"
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
            <!-- Temple Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Temple:</span>
              {{ selectedTemple === 'all' ? 'All Temples' : getTempleName(selectedTemple) }}
            </div>
            
            <!-- Activity Type Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Activity:</span>
              {{ getActivityTypeLabel(activityType) }}
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
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute } from 'vue-router';
import { useTempleStore } from '@/stores/temple';
import { useAuthStore } from '@/stores/auth';
import { useReportsStore } from '@/stores/reports';
import { useToast } from '@/composables/useToast';

// Composables
const route = useRoute();
const templeStore = useTempleStore();
const userStore = useAuthStore();
const reportsStore = useReportsStore();
const { showToast } = useToast();

// Reactive state
const selectedTemple = ref('all');
const activityType = ref('events');
const activeFilter = ref('monthly');
const selectedFormat = ref('pdf');
const startDate = ref(new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0]);
const endDate = ref(new Date().toISOString().split('T')[0]);

// Filter options
const activityTypes = [
  { label: 'Events Report', value: 'events' },
  { label: 'Seva Reports', value: 'sevas' },
  { label: 'Seva Bookings', value: 'bookings' },
  { label: 'Donations Reports', value: 'donations' }, // Added new report type
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

// Computed
const tenantId = computed(() => {
  return route.params.tenantId || userStore.user?.id || localStorage.getItem('current_tenant_id');
});

const previewRows = computed(() => {
  if (!reportsStore.reportPreview?.data) return []
  // Show only first 5 rows for preview
  return reportsStore.formatReportData(
    reportsStore.reportPreview.data.slice(0, 5), 
    activityType.value
  )
});

// Methods
const setActivityType = async (type) => {
  activityType.value = type;
  await loadReportPreview();
};

const setActiveFilter = (filter) => {
  activeFilter.value = filter;
  
  // Set appropriate date range based on filter
  if (filter !== 'custom') {
    const dateRange = reportsStore.getDefaultDateRange(filter);
    startDate.value = dateRange.startDate;
    endDate.value = dateRange.endDate;
  }
  
  // Auto-refresh preview after filter change
  loadReportPreview();
};

const handleTempleChange = () => {
  loadReportPreview();
};

const handleDateChange = () => {
  if (activeFilter.value === 'custom') {
    loadReportPreview();
  }
};

const getTempleName = (templeId) => {
  if (templeId === 'all') return 'All Temples';
  const temple = templeStore.temples.find(t => t.id.toString() === templeId.toString());
  return temple ? temple.name : 'Unknown Temple';
};

const getActivityTypeLabel = (type) => {
  const found = activityTypes.find(t => t.value === type);
  return found ? found.label : 'Unknown';
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

const loadReportPreview = async () => {
  try {
    const params = reportsStore.buildReportParams({
      selectedTemple: selectedTemple.value,
      activityType: activityType.value,
      activeFilter: activeFilter.value,
      selectedFormat: '', // No format for preview
      startDate: startDate.value,
      endDate: endDate.value
    });

    await reportsStore.getReportPreview(params);
  } catch (error) {
    console.error('Error loading report preview:', error);
    showToast('Failed to load report preview. Please try again.', 'error');
  }
};

const downloadReport = async () => {
  try {
    const params = reportsStore.buildReportParams({
      selectedTemple: selectedTemple.value,
      activityType: activityType.value,
      activeFilter: activeFilter.value,
      selectedFormat: selectedFormat.value,
      startDate: startDate.value,
      endDate: endDate.value
    });

    const result = await reportsStore.downloadActivitiesReport(params);
    
    if (result.success) {
      showToast(`${getActivityTypeLabel(activityType.value)} downloaded successfully as ${getFormatLabel(selectedFormat.value)}`, 'success');
    }
  } catch (error) {
    console.error('Error downloading report:', error);
    showToast('Failed to download report. Please try again.', 'error');
  }
};

// Lifecycle hooks
onMounted(async () => {
  // Fetch temples if not already loaded
  if (templeStore.temples.length === 0) {
    try {
      await templeStore.fetchTemples(tenantId.value);
    } catch (error) {
      console.error('Error loading temple data:', error);
      showToast('Failed to load temple data. Please try again.', 'error');
    }
  }

  // Load initial preview
  await loadReportPreview();
});

// Cleanup on unmount
onMounted(() => {
  return () => {
    reportsStore.clearReportData();
  };
});
</script>