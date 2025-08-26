<template>
  <div class="min-h-screen bg-gray-50">
   

    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Audit Logs Report</h1>
            <p class="text-gray-600 mt-1">
              Track and monitor all activities within your temple system
              <span v-if="currentEntityId" class="text-indigo-600 font-medium">
                (Tenant ID: {{ currentEntityId }})
              </span>
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">{{ adminName }}</span>
              <span class="text-indigo-600 text-sm ml-2">({{ adminRole }})</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
        <span class="text-gray-900 font-medium">Loading audit logs...</span>
      </div>
    </div>

    <!-- Error Alert -->
    <div v-if="exportStatus.show && exportStatus.type === 'error'" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4">
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-red-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-red-800">Error</h3>
            <p class="mt-1 text-sm text-red-700">{{ exportStatus.message }}</p>
          </div>
          <div class="ml-auto pl-3">
            <button @click="exportStatus.show = false" class="text-red-400 hover:text-red-500">
              <span class="sr-only">Dismiss</span>
              <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Alert -->
    <div v-if="exportStatus.show && exportStatus.type === 'success'" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-4">
      <div class="bg-green-50 border border-green-200 rounded-md p-4">
        <div class="flex">
          <svg class="h-5 w-5 text-green-400" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
          </svg>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-green-800">Success</h3>
            <p class="mt-1 text-sm text-green-700">{{ exportStatus.message }}</p>
          </div>
          <div class="ml-auto pl-3">
            <button @click="exportStatus.show = false" class="text-green-400 hover:text-green-500">
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
          <h3 class="text-xl font-bold text-gray-900">Audit Logs</h3>
          <p class="text-gray-600 mt-1">Configure filters and download your audit log data</p>
        </div>

        <div class="p-6">
          <!-- Action Type Selection (Dropdown with grouped options excluding Superadmin) -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Action Type</label>
            <select
              v-model="filters.action"
              @change="applyFilters"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Actions</option>
              <!-- Authentication Module Actions -->
              <optgroup label="Authentication">
                <option value="REGISTRATION_SUCCESS">Registration Success</option>
                <option value="REGISTRATION_FAILED">Registration Failed</option>
                <option value="TEMPLEADMIN_REGISTRATION_SUCCESS">Temple Admin Registration</option>
                <option value="REGISTRATION_BLOCKED">Registration Blocked</option>
                <option value="LOGIN_SUCCESS">Login Success</option>
                <option value="LOGIN_FAILED">Login Failed</option>
                <option value="LOGOUT">Logout</option>
                <option value="PASSWORD_RESET_REQUESTED">Password Reset Requested</option>
                <option value="PASSWORD_RESET_SUCCESS">Password Reset Success</option>
                <option value="PASSWORD_RESET_FAILED">Password Reset Failed</option>
              </optgroup>
              
              <!-- Entity Module Actions -->
              <optgroup label="Entity">
                <option value="TEMPLE_CREATED">Temple Created</option>
                <option value="TEMPLE_UPDATED">Temple Updated</option>
                <option value="TEMPLE_CREATE_FAILED">Temple Create Failed</option>
                <option value="TEMPLE_UPDATE_FAILED">Temple Update Failed</option>
              </optgroup>
              
              <!-- Events Module Actions -->
              <optgroup label="Events">
                <option value="EVENT_CREATED">Event Created</option>
                <option value="EVENT_UPDATED">Event Updated</option>
                <option value="EVENT_DELETED">Event Deleted</option>
              </optgroup>
              
              <!-- Seva Module Actions -->
              <optgroup label="Seva">
                <option value="SEVA_CREATED">Seva Created</option>
                <option value="SEVA_UPDATED">Seva Updated</option>
                <option value="SEVA_BOOKED">Seva Booked</option>
                <option value="SEVA_BOOKING_APPROVED">Seva Booking Approved</option>
                <option value="SEVA_BOOKING_REJECTED">Seva Booking Rejected</option>
              </optgroup>
              
              <!-- Donations Module Actions -->
              <optgroup label="Donations">
                <option value="DONATION_INITIATED">Donation Initiated</option>
                <option value="DONATION_SUCCESS">Donation Success</option>
                <option value="DONATION_FAILED">Donation Failed</option>
                <option value="DONATION_VERIFICATION_FAILED">Donation Verification Failed</option>
              </optgroup>
              
              <!-- Notifications Module Actions -->
              <optgroup label="Notifications">
                <option value="TEMPLATE_CREATED">Template Created</option>
                <option value="TEMPLATE_UPDATED">Template Updated</option>
                <option value="TEMPLATE_DELETED">Template Deleted</option>
                <option value="EMAIL_SENT">Email Sent</option>
                <option value="SMS_SENT">SMS Sent</option>
                <option value="WHATSAPP_SENT">WhatsApp Sent</option>
              </optgroup>
              
              <!-- User Profile Module Actions -->
              <optgroup label="User Profile">
                <option value="PROFILE_CREATED">Profile Created</option>
                <option value="PROFILE_UPDATED">Profile Updated</option>
                <option value="DEVOTEE_JOINED_TEMPLE">Devotee Joined Temple</option>
                <option value="VOLUNTEER_JOINED_TEMPLE">Volunteer Joined Temple</option>
              </optgroup>
              
              <!-- Reports Module Actions -->
              <optgroup label="Reports">
                <option value="DEVOTEE_BIRTHDAYS_REPORT_VIEWED">Birthday Report Viewed</option>
                <option value="DEVOTEE_BIRTHDAYS_REPORT_DOWNLOADED">Birthday Report Downloaded</option>
                <option value="TEMPLE_REGISTER_REPORT_VIEWED">Temple Register Report Viewed</option>
                <option value="TEMPLE_REGISTER_REPORT_DOWNLOADED">Temple Register Report Downloaded</option>
                <option value="TEMPLE_ACTIVITIES_REPORT_VIEWED">Activities Report Viewed</option>
                <option value="TEMPLE_ACTIVITIES_REPORT_DOWNLOADED">Activities Report Downloaded</option>
              </optgroup>
            </select>
          </div>

          <!-- Status Selection -->
          <div class="mb-6">
            <label class="block text-gray-700 font-medium mb-2">Status</label>
            <div class="flex flex-wrap gap-2">
              <button 
                v-for="status in statusOptions" 
                :key="status.value"
                @click="setStatus(status.value)"
                class="px-4 py-2 rounded-md text-sm font-medium transition-colors duration-200"
                :class="filters.status === status.value ? 
                  'bg-indigo-600 text-white' : 
                  'bg-gray-100 text-gray-700 hover:bg-gray-200'"
              >
                {{ status.label }}
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
                  v-model="filters.startDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  @change="handleDateChange"
                />
              </div>
              <div>
                <label class="block text-gray-700 text-sm font-medium mb-2">End Date</label>
                <input 
                  type="date" 
                  v-model="filters.endDate"
                  class="w-full py-2 px-3 border border-gray-300 bg-white rounded-md shadow-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500"
                  @change="handleDateChange"
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
                </div>
                <!-- Download Button -->
                <button 
                  @click="downloadReport"
                  :disabled="loading || downloadLoading"
                  class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg v-if="downloadLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
                  </svg>
                  {{ downloadLoading ? 'Downloading...' : 'Download' }}
                </button>

                <!-- Refresh Button -->
                <button 
                  @click="fetchAuditLogs"
                  :disabled="loading"
                  class="inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md shadow-sm text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed"
                >
                  <svg class="mr-2 h-5 w-5" :class="{ 'animate-spin': loading }" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                  </svg>
                  {{ loading ? 'Refreshing...' : 'Refresh' }}
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Current Applied Filters -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden mb-8">
        <div class="p-6">
          <h3 class="text-lg font-medium text-gray-900 mb-4">Applied Filters</h3>
          <div class="flex flex-wrap gap-2">
            <!-- Action Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Action:</span>
              {{ getActionTypeLabel(filters.action) }}
            </div>
            <!-- Status Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Status:</span>
              {{ getStatusLabel(filters.status) }}
            </div>
            <!-- Date Range Filter -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Period:</span>
              {{ getTimeFilterLabel(activeFilter) }}
              <span v-if="activeFilter === 'custom'">
                ({{ formatDateShort(filters.startDate) }} - {{ formatDateShort(filters.endDate) }})
              </span>
            </div>
            <!-- Format -->
            <div class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-indigo-100 text-indigo-800">
              <span class="font-medium mr-1">Format:</span>
              {{ getFormatLabel(selectedFormat) }}
            </div>
            <!-- Clear All Filters -->
            <button 
              @click="clearAllFilters"
              class="inline-flex items-center px-3 py-1.5 rounded-full text-sm bg-red-100 text-red-800 hover:bg-red-200 transition-colors"
            >
              <svg class="w-4 h-4 mr-1" fill="currentColor" viewBox="0 0 20 20">
                <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
              </svg>
              Clear All
            </button>
          </div>
          <p class="mt-4 text-sm text-gray-600">
            Your report will include audit log data based on the filters above. Click Download to generate and download the report.
          </p>
        </div>
      </div>

      <!-- Detail Modal (unchanged, omitted here for brevity) -->
      <div v-if="selectedLogData" class="fixed inset-0 z-50 overflow-y-auto">
        <!-- Modal content unchanged as before -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { nextTick } from 'vue';
import ReportsService from '@/services/reports.service';

// State
const logs = ref([]);
const apiResponse = ref(null);
const loading = ref(false);
const downloadLoading = ref(false);
const selectedLogId = ref(null);
const selectedLogData = ref(null);
const activeFilter = ref('monthly');
const selectedFormat = ref('pdf');

// Entity ID - replace with your actual logic
const currentEntityId = ref('all');
const fromReports = ref(true);
//
const adminName = ref('Temple Admin');
const adminRole = ref('Administrator');
// Export status
const exportStatus = ref({
  show: false,
  type: 'info', // 'success', 'error', 'info'
  message: ''
});

// Pagination state
const currentPage = ref(1);
const pageSize = ref(10);
const totalItems = ref(0);
const totalPages = ref(1);

// Filter state
const filters = ref({
  action: '',
  status: '',
  startDate: '',
  endDate: ''
});

// Summary state
const summary = ref({
  totalActivities: 0,
  activeUsers: 0,
  failedActions: 0,
  recentActions: 0
});

// Filter options
const actionTypes = [
  { label: 'All Actions', value: '' },
  { label: 'Authentication', value: 'LOGIN_SUCCESS' },
  { label: 'User Management', value: 'USER_CREATED' },
  { label: 'Temple Updates', value: 'TEMPLE_UPDATED' },
  { label: 'Events', value: 'EVENT_CREATED' },
  { label: 'Seva', value: 'SEVA_CREATED' },
  { label: 'Donations', value: 'DONATION_SUCCESS' },
  { label: 'Reports', value: 'REPORT_VIEWED' }
];

const statusOptions = [
  { label: 'All Status', value: '' },
  { label: 'Success', value: 'success' },
  { label: 'Failed', value: 'failure' }
];

const timeFilters = [
  { label: 'Daily', value: 'daily' },
  { label: 'Weekly', value: 'weekly' },
  { label: 'Monthly', value: 'monthly' },
  
  { label: 'Custom Range', value: 'custom' }
];

const formats = [
  { label: 'PDF', value: 'pdf' },
  { label: 'CSV', value: 'csv' },
  { label: 'Excel', value: 'excel' }
];

// Computed values
const hasActiveFilters = computed(() => {
  return filters.value.action || 
         filters.value.status || 
         filters.value.startDate ||
         filters.value.endDate ||
         activeFilter.value !== 'monthly';
});

const filteredLogs = computed(() => {
  let filtered = [...logs.value];

  // Filter by action
  if (filters.value.action) {
    filtered = filtered.filter(log => 
      log.action && log.action.toLowerCase().includes(filters.value.action.toLowerCase())
    );
  }

  // Filter by status
  if (filters.value.status) {
    filtered = filtered.filter(log => 
      log.status && log.status.toLowerCase() === filters.value.status.toLowerCase()
    );
  }

  // Filter by date range
  if (filters.value.startDate || filters.value.endDate) {
    filtered = filtered.filter(log => {
      const logDate = new Date(log.created_at);
      
      if (filters.value.startDate) {
        const startDate = new Date(filters.value.startDate);
        if (logDate < startDate) return false;
      }
      
      if (filters.value.endDate) {
        const endDate = new Date(filters.value.endDate);
        endDate.setDate(endDate.getDate() + 1);
        if (logDate >= endDate) return false;
      }
      
      return true;
    });
  }

  totalItems.value = filtered.length;
  totalPages.value = Math.ceil(totalItems.value / pageSize.value);
  
  return filtered;
});

const paginatedLogs = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value;
  const end = start + pageSize.value;
  return filteredLogs.value.slice(start, end);
});

const paginationRange = computed(() => {
  const range = [];
  const showPages = 5;
  
  let start = Math.max(1, currentPage.value - Math.floor(showPages / 2));
  let end = Math.min(totalPages.value, start + showPages - 1);
  
  if (end - start + 1 < showPages) {
    start = Math.max(1, end - showPages + 1);
  }
  
  if (start > 1) {
    range.push(1);
    if (start > 2) range.push('...');
  }
  
  for (let i = start; i <= end; i++) {
    range.push(i);
  }
  
  if (end < totalPages.value) {
    if (end < totalPages.value - 1) range.push('...');
    range.push(totalPages.value);
  }
  
  return range;
});

// Methods
function getDefaultDateRange(filter) {
  const today = new Date();
  let startDate, endDate;
  
  switch (filter) {
    case 'daily':
      startDate = new Date(today);
      endDate = new Date(today);
      break;
    case 'weekly':
      startDate = new Date(today.getTime() - (7 * 24 * 60 * 60 * 1000));
      endDate = new Date(today);
      break;
    case 'monthly':
      startDate = new Date(today.getTime() - (30 * 24 * 60 * 60 * 1000));
      endDate = new Date(today);
      break;
    
    default:
      startDate = new Date(today.getTime() - (30 * 24 * 60 * 60 * 1000));
      endDate = new Date(today);
  }
  
  return {
    startDate: startDate.toISOString().split('T'),
    endDate: endDate.toISOString().split('T')
  };
}

function setActionType(type) {
  filters.value.action = type;
  currentPage.value = 1;
  applyFilters();
}

function setStatus(status) {
  filters.value.status = status;
  currentPage.value = 1;
  applyFilters();
}

function setActiveFilter(filter) {
  activeFilter.value = filter;
  
  if (filter !== 'custom') {
    const dateRange = getDefaultDateRange(filter);
    filters.value.startDate = dateRange.startDate;
    filters.value.endDate = dateRange.endDate;
  }
  
  currentPage.value = 1;
  applyFilters();
}

function handleDateChange() {
  if (activeFilter.value === 'custom') {
    applyFilters();
  }
}

function clearAllFilters() {
  filters.value = {
    action: '',
    status: '',
    startDate: '',
    endDate: ''
  };
  activeFilter.value = 'monthly';
  
  const dateRange = getDefaultDateRange('monthly');
  filters.value.startDate = dateRange.startDate;
  filters.value.endDate = dateRange.endDate;
  
  currentPage.value = 1;
  applyFilters();
}

function getActionTypeLabel(action) {
  if (!action) return 'All Actions';
  const found = actionTypes.find(t => t.value === action);
  return found ? found.label : action;
}

function getStatusLabel(status) {
  if (!status) return 'All Status';
  const found = statusOptions.find(s => s.value === status);
  return found ? found.label : status;
}

function getTimeFilterLabel(filter) {
  const found = timeFilters.find(f => f.value === filter);
  return found ? found.label : filter;
}

function getFormatLabel(format) {
  const found = formats.find(f => f.value === format);
  return found ? found.label : format;
}

function formatDate(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleString('en-US', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit',
    second: '2-digit',
    hour12: true
  });
}

function formatDateShort(dateString) {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  });
}

function parseDetails(details) {
  if (!details) return {};
  if (typeof details === 'object') return details;
  
  try {
    return JSON.parse(details);
  } catch (e) {
    console.error('Failed to parse details:', e);
    return { raw: details };
  }
}

function getUserName(log) {
  const details = parseDetails(log.details);
  return details.target_user_name || details.user_name || `User ${log.user_id || 'Unknown'}`;
}

function getEntityName(log) {
  const details = parseDetails(log.details);
  return details.entity_name || details.temple_name || 'Temple';
}

// leave calculateSummary unchanged in logic
function calculateSummary() {
  if (!logs.value.length) return;
  
  summary.value.totalActivities = logs.value.length;
  summary.value.activeUsers = new Set(logs.value.map(log => getUserName(log))).size;
  summary.value.failedActions = logs.value.filter(log => log.status === 'failure').length;
  
  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);
  summary.value.recentActions = logs.value.filter(log => 
    new Date(log.created_at) > yesterday
  ).length;
}

async function fetchAuditLogs() {
  loading.value = true;
  try {
    console.log('Fetching audit logs using ReportsService...');
    const params = {
      entityId: currentEntityId.value,
      dateRange: 'custom',
      startDate: filters.value.startDate,
      endDate: filters.value.endDate,
      actionType: filters.value.action,
    };

    Object.keys(params).forEach(key => {
      if (!params[key]) {
        delete params[key];
      }
    });

    console.log('API parameters:', params);
    const response = await ReportsService.getAuditLogsReport(params);
    console.log('ReportsService response:', response);

    apiResponse.value = response.data;

    let responseData = response.data;
    if (responseData && responseData.data) {
      responseData = responseData.data;
    }

    if (Array.isArray(responseData)) {
      logs.value = responseData;
    } else if (responseData && responseData.audit_logs) {
      logs.value = responseData.audit_logs;
    } else if (responseData && responseData.logs) {
      logs.value = responseData.logs;
    } else {
      console.warn('Unexpected response format:', responseData);
      logs.value = [];
    } 
    calculateSummary();
    await nextTick();
    console.log('Audit logs loaded successfully:', logs.value.length, 'records');
  } catch (error) {
    console.error('Error fetching audit logs:', error);
    exportStatus.value = {
      show: true,
      type: 'error',
      message: 'Failed to load audit logs. Please check your connection and try again.'
    };
    // Mock data for demonstration
    logs.value = [
      {
        id: 1,
        user_id: 101,
        action: 'LOGIN_SUCCESS',
        status: 'success',
        ip_address: '192.168.1.100',
        created_at: '2025-01-15T10:30:00Z',
        details: JSON.stringify({
          user_name: 'Temple Admin',
          entity_name: 'Sri Venkateswara Temple',
          device: 'Desktop'
        })
      },
      {
        id: 2,
        user_id: 102,
        action: 'SEVA_CREATED',
        status: 'success',
        ip_address: '192.168.1.101',
        created_at: '2025-01-15T11:45:00Z',
        details: JSON.stringify({
          user_name: 'Seva Coordinator',
          entity_name: 'Sri Venkateswara Temple',
          seva_name: 'Abhisheka Seva',
          seva_price: '₹500'
        })
      },
      {
        id: 3,
        user_id: 103,
        action: 'DONATION_SUCCESS',
        status: 'success',
        ip_address: '192.168.1.102',
        created_at: '2025-01-15T12:15:00Z',
        details: JSON.stringify({
          user_name: 'Devotee Kumar',
          entity_name: 'Sri Venkateswara Temple',
          donation_amount: '₹1000',
          payment_method: 'UPI'
        })
      },
      {
        id: 4,
        user_id: 104,
        action: 'EVENT_CREATED',
        status: 'success',
        ip_address: '192.168.1.103',
        created_at: '2025-01-15T13:20:00Z',
        details: JSON.stringify({
          user_name: 'Event Manager',
          entity_name: 'Sri Venkateswara Temple',
          event_name: 'Brahmotsavam Festival',
          event_date: '2025-02-15'
        })
      },
      {
        id: 5,
        user_id: 105,
        action: 'LOGIN_FAILED',
        status: 'failure',
        ip_address: '192.168.1.104',
        created_at: '2025-01-15T14:25:00Z',
        details: JSON.stringify({
          user_name: 'Unknown User',
          entity_name: 'Sri Venkateswara Temple',
          error: 'Invalid credentials'
        })
      }
    ];
    calculateSummary();
  } finally {
    loading.value = false;
    console.log('Fetch completed, logs:', logs.value.length);
  }
}

async function downloadReport() {
  downloadLoading.value = true;
  try {
    exportStatus.value = {
      show: true,
      type: 'info',
      message: `Preparing ${selectedFormat.value.toUpperCase()} export...`
    };

    if (!filteredLogs.value.length) {
      exportStatus.value = {
        show: true,
        type: 'error',
        message: 'No data available to export. Please adjust your filters.'
      };
      return;
    }

    const downloadParams = {
      entityId: currentEntityId.value,
      format: selectedFormat.value,
      dateRange: 'custom',
      startDate: filters.value.startDate,
      endDate: filters.value.endDate,
      actionType: filters.value.action,
    };

    Object.keys(downloadParams).forEach(key => {
      if (!downloadParams[key]) {
        delete downloadParams[key];
      }
    });

    console.log('Downloading with parameters:', downloadParams);

    const result = await ReportsService.downloadAuditLogsReport(downloadParams);

    exportStatus.value = {
      show: true,
      type: 'success',
      message: `${selectedFormat.value.toUpperCase()} report "${result.filename}" downloaded successfully!`
    };

    setTimeout(() => {
      exportStatus.value.show = false;
    }, 3000);

  } catch (error) {
    console.error('Error exporting report:', error);
    exportStatus.value = {
      show: true,
      type: 'error',
      message: `Failed to export ${selectedFormat.value.toUpperCase()} report: ${error.message || 'Please try again.'}`
    };
  } finally {
    downloadLoading.value = false;
  }
}

function applyFilters() {
  console.log('Filters applied:', filters.value);
  currentPage.value = 1;
  fetchAuditLogs();
}

function selectLog(id) {
  selectedLogId.value = id;
  selectedLogData.value = logs.value.find(log => log.id === id);
}

function closeDetailModal() {
  selectedLogId.value = null;
  selectedLogData.value = null;
}

// Pagination methods
function goToPage(page) {
  currentPage.value = page;
}

function goToPreviousPage() {
  if (currentPage.value > 1) {
    currentPage.value--;
  }
}

function goToNextPage() {
  if (currentPage.value < totalPages.value) {
    currentPage.value++;
  }
}

// Initialize
onMounted(() => {
  console.log('Component mounted, initializing...');
  // Set default date range (monthly)
  const dateRange = getDefaultDateRange('monthly');
  filters.value.startDate = dateRange.startDate;
  filters.value.endDate = dateRange.endDate;

  console.log('Default date range set:', filters.value.startDate, 'to', filters.value.endDate);
  console.log('Using entity ID:', currentEntityId.value);

  fetchAuditLogs();
});

// Watch for filter changes
watch(filters, () => {
  currentPage.value = 1;
}, { deep: true });
</script>