<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Devotee Birthdays Report</h1>
            <p class="text-gray-600 mt-1">
              Download devotee birthday data for your temples
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
          <h3 class="text-xl font-bold text-gray-900">Devotee Birthdays</h3>
          <p class="text-gray-600 mt-1">Configure filters and download devotee birthday data</p>
        </div>

        <div class="p-6">
          <!-- Temple Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
            <div class="relative">
              <select 
                v-model="selectedTemple" 
                class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
              >
                <option value="all">All Temples</option>
                <option v-for="temple in templeStore.temples" :key="temple.id" :value="temple.id">
                  {{ temple.name }}
                </option>
              </select>
              <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                <!-- <span class="text-xs">▼</span> -->
              </div>
            </div>
          </div>

          <!-- Filter Section -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-6">
            <!-- Birthday Date Range Filter -->
            <div>
              <label class="block text-gray-700 font-medium mb-2">Birthday Period</label>
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

          <!-- Custom Date Range (shown only when custom date range is selected) -->
          <div v-if="activeFilter === 'custom'" class="mb-6 p-4 bg-gray-50 border border-gray-200 rounded-lg">
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">Start Date</label>
                <input 
                  type="date" 
                  v-model="startDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="endDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
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
                    class="block w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  >
                    <option v-for="format in formats" :key="format.value" :value="format.value">
                      {{ format.label }}
                    </option>
                  </select>
                  <div class="pointer-events-none absolute inset-y-0 right-0 flex items-center px-2 text-gray-700">
                    <!-- <span class="text-xs">▼</span> -->
                  </div>
                </div>

                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  <svg class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  Download
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
            
            <!-- Birthday Period Filter -->
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
import { ref, computed, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import { useTempleStore } from '@/stores/temple';
import { useAuthStore } from '@/stores/auth';
import { useToast } from '@/composables/useToast';

// Composables
const route = useRoute();
const templeStore = useTempleStore();
const userStore = useAuthStore();
const { showToast } = useToast();

// Reactive state
const selectedTemple = ref('all');
const activeFilter = ref('monthly');
const ageGroup = ref('all');
const selectedFormat = ref('pdf');
const startDate = ref(new Date(new Date().setDate(new Date().getDate() - 30)).toISOString().split('T')[0]);
const endDate = ref(new Date().toISOString().split('T')[0]);
const includeContactInfo = ref(true);
const includeFamilyMembers = ref(false);

// Filter options
const timeFilters = [
  { label: 'Weekly', value: 'weekly' },
  { label: 'Monthly', value: 'monthly' },
  { label: 'Yearly', value: 'yearly' },
  { label: 'Custom Range', value: 'custom' },
];

const ageGroups = [
  { label: 'All Ages', value: 'all' },
  { label: 'Children (0-12)', value: 'children' },
  { label: 'Teens (13-19)', value: 'teens' },
  { label: 'Adults (20-59)', value: 'adults' },
  { label: 'Seniors (60+)', value: 'seniors' },
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

// Methods
const setActiveFilter = (filter) => {
  activeFilter.value = filter;
  
  // Set appropriate date range based on filter
  const today = new Date();
  
  if (filter === 'weekly') {
    // Next 7 days
    startDate.value = new Date().toISOString().split('T')[0];
    endDate.value = new Date(today.setDate(today.getDate() + 7)).toISOString().split('T')[0];
  } else if (filter === 'monthly') {
    // Next 30 days
    startDate.value = new Date().toISOString().split('T')[0];
    endDate.value = new Date(today.setDate(today.getDate() + 30)).toISOString().split('T')[0];
  } else if (filter === 'yearly') {
    // Current year
    const currentYear = today.getFullYear();
    startDate.value = new Date(currentYear, 0, 1).toISOString().split('T')[0]; // January 1st
    endDate.value = new Date(currentYear, 11, 31).toISOString().split('T')[0]; // December 31st
  }
  
  // For custom, we leave the dates as they are
};

const setAgeGroup = (group) => {
  ageGroup.value = group;
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

const getAgeGroupLabel = (group) => {
  const found = ageGroups.find(g => g.value === group);
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

const downloadReport = () => {
  // Simulate downloading a report
  console.log('Downloading birthdays report with the following parameters:');
  console.log('- Temple:', selectedTemple.value === 'all' ? 'All Temples' : getTempleName(selectedTemple.value));
  console.log('- Time filter:', getTimeFilterLabel(activeFilter.value));
  console.log('- Date range:', formatDate(startDate.value), 'to', formatDate(endDate.value));
  console.log('- Age Group:', getAgeGroupLabel(ageGroup.value));
  console.log('- Include Contact Info:', includeContactInfo.value);
  console.log('- Include Family Members:', includeFamilyMembers.value);
  console.log('- Format:', getFormatLabel(selectedFormat.value));
  
  // In a real implementation, you would make an API call here
  showToast(`Birthdays Report downloaded in ${getFormatLabel(selectedFormat.value)} format`, 'success');
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
});
</script>