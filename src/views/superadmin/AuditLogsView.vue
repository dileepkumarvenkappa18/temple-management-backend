<template>
  <div class="w-full max-w-full">
    <!-- Header -->
    <div class="mb-6">
      <AppBreadcrumb :items="breadcrumbItems" />
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Audit Logs</h1>
          <p class="mt-1 text-sm text-gray-600">
            Track and monitor all system activities across your platform
          </p>
        </div>
        <div class="mt-4 sm:mt-0">
          <BaseButton
            variant="primary"
            @click="fetchAuditLogs"
            class="w-full sm:w-auto"
          >
            <template #icon>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </template>
            Refresh
          </BaseButton>
        </div>
      </div>
    </div>

    <!-- Filters Section -->
    <BaseCard class="border border-gray-200 shadow-sm mb-6">
      <template #header>
        <div class="flex justify-between items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
          <h3 class="text-lg font-semibold text-gray-900">Filters</h3>
          <BaseButton
            variant="outline"
            size="sm"
            @click="clearFilters"
          >
            Clear Filters
          </BaseButton>
        </div>
      </template>
      
      <div class="p-6">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <!-- Action Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Action</label>
            <select
              v-model="filters.action"
              @change="applyFilters"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Actions</option>
              <option value="login">Login</option>
              <option value="logout">Logout</option>
              <option value="create">Create</option>
              <option value="update">Update</option>
              <option value="delete">Delete</option>
              <option value="view">View</option>
            </select>
          </div>

          <!-- Status Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              v-model="filters.status"
              @change="applyFilters"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Status</option>
              <option value="success">Success</option>
              <option value="failure">Failure</option>
            </select>
          </div>

          <!-- User Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">User</label>
            <input
              v-model="filters.user"
              @input="debounceFilter"
              type="text"
              placeholder="Search by user..."
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>

          <!-- Date Range Filter -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
            <select
              v-model="filters.dateRange"
              @change="applyFilters"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">All Time</option>
              <option value="today">Today</option>
              <option value="yesterday">Yesterday</option>
              <option value="last7days">Last 7 Days</option>
              <option value="last30days">Last 30 Days</option>
              <option value="last90days">Last 90 Days</option>
            </select>
          </div>
        </div>

        <!-- Custom Date Range -->
        <div v-if="filters.dateRange === 'custom'" class="mt-4 grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">From Date</label>
            <input
              v-model="filters.startDate"
              @change="applyFilters"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">To Date</label>
            <input
              v-model="filters.endDate"
              @change="applyFilters"
              type="date"
              class="w-full px-3 py-2 border border-gray-300 rounded-md shadow-sm focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
        </div>

        <!-- Active Filters Display -->
        <div v-if="hasActiveFilters" class="mt-4">
          <div class="flex flex-wrap gap-2">
            <span class="text-sm font-medium text-gray-700">Active filters:</span>
            <span v-if="filters.action" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800">
              Action: {{ filters.action }}
              <button @click="clearFilter('action')" class="ml-1 text-indigo-600 hover:text-indigo-500">
                <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </button>
            </span>
            <span v-if="filters.status" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800">
              Status: {{ filters.status }}
              <button @click="clearFilter('status')" class="ml-1 text-green-600 hover:text-green-500">
                <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </button>
            </span>
            <span v-if="filters.user" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
              User: {{ filters.user }}
              <button @click="clearFilter('user')" class="ml-1 text-blue-600 hover:text-blue-500">
                <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </button>
            </span>
            <span v-if="filters.dateRange" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
              Date: {{ formatDateRangeLabel(filters.dateRange) }}
              <button @click="clearFilter('dateRange')" class="ml-1 text-purple-600 hover:text-purple-500">
                <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </button>
            </span>
          </div>
        </div>
      </div>
    </BaseCard>
    
    <!-- Audit Logs Table -->
    <BaseCard class="border border-gray-200 shadow-sm mb-4">
      <template #header>
        <div class="flex justify-between items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">
            Audit Logs
            <span v-if="filteredLogs.length !== logs.length" class="text-sm font-normal text-gray-500">
              ({{ filteredLogs.length }} of {{ logs.length }} records)
            </span>
          </h3>
        </div>
      </template>
      
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">ID</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User Name</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Entity</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Action</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">IP Address</th>
              <th scope="col" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Created At</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-if="isLoading">
              <td colspan="7" class="px-6 py-10 text-center">
                <div class="flex justify-center">
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Loading...
                </div>
              </td>
            </tr>
            <tr v-else-if="!filteredLogs.length">
              <td colspan="7" class="px-6 py-10 text-center text-gray-500">
                <div class="flex flex-col items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <p class="text-lg font-medium">
                    {{ hasActiveFilters ? 'No audit logs match your filters' : 'No audit logs found' }}
                  </p>
                  <p class="text-sm mt-1">
                    {{ hasActiveFilters ? 'Try adjusting your filters' : 'Try adjusting your filters or check back later' }}
                  </p>
                </div>
              </td>
            </tr>
            <tr 
              v-else
              v-for="log in filteredLogs" 
              :key="log.id"
              @click="selectLog(log.id)"
              class="cursor-pointer hover:bg-indigo-50 transition-colors duration-150"
            >
              <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">{{ log.id }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ getUserName(log) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ getEntityName(log) }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.action }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                <span 
                  :class="log.status === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
                  class="px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
                >
                  {{ log.status }}
                </span>
              </td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ log.ip_address }}</td>
              <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">{{ formatDate(log.created_at) }}</td>
            </tr>
          </tbody>
        </table>
      </div>
    </BaseCard>
    
    <!-- Debug Information -->
    <!-- <BaseCard class="border border-gray-200 shadow-sm mb-4" v-if="showDebug">
      <template #header>
        <div class="flex justify-between items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">Debug Information</h3>
          <BaseButton
            variant="outline"
            size="sm"
            @click="showDebug = false"
          >
            Hide
          </BaseButton>
        </div>
      </template>
      
      <div class="p-4">
        <h4 class="text-md font-semibold mb-2">Raw API Response:</h4>
        <pre class="bg-gray-100 p-4 rounded-lg overflow-auto max-h-96 text-xs">{{ JSON.stringify(apiResponse, null, 2) }}</pre>
        
        <h4 class="text-md font-semibold mt-4 mb-2">Parsed Logs:</h4>
        <pre class="bg-gray-100 p-4 rounded-lg overflow-auto max-h-96 text-xs">{{ JSON.stringify(logs, null, 2) }}</pre>
        
        <h4 class="text-md font-semibold mt-4 mb-2">Filtered Logs:</h4>
        <pre class="bg-gray-100 p-4 rounded-lg overflow-auto max-h-96 text-xs">{{ JSON.stringify(filteredLogs, null, 2) }}</pre>
        
        <h4 class="text-md font-semibold mt-4 mb-2">API Endpoint:</h4>
        <div class="flex space-x-2">
          <BaseButton
            variant="outline"
            size="sm"
            @click="tryEndpoint('/api/v1/auditlogs')"
          >
            Try /api/v1/auditlogs
          </BaseButton>
          <BaseButton
            variant="outline"
            size="sm"
            @click="tryEndpoint('/v1/auditlogs')"
          >
            Try /v1/auditlogs
          </BaseButton>
          <BaseButton
            variant="outline"
            size="sm"
            @click="tryEndpoint('/api/auditlogs')"
          >
            Try /api/auditlogs
          </BaseButton>
          <BaseButton
            variant="outline"
            size="sm"
            @click="tryEndpoint('/auditlogs')"
          >
            Try /auditlogs
          </BaseButton>
        </div>
      </div>
    </BaseCard> -->
    
    <!-- Detail Modal (simplified) -->
    <div v-if="selectedLogData">
      <BaseModal @close="closeDetailModal">
        <template #header>
          <div class="flex items-center justify-between bg-gray-50 px-4 py-3 border-b border-gray-200 rounded-t-lg">
            <h3 class="text-lg font-bold leading-6 text-gray-900">Audit Log Details</h3>
            <span 
              :class="selectedLogData.status === 'success' ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'" 
              class="ml-2 px-2 inline-flex text-xs leading-5 font-semibold rounded-full"
            >
              {{ selectedLogData.status }}
            </span>
          </div>
        </template>
        
        <div class="p-6 space-y-6">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4 p-5 bg-gray-50 rounded-lg border border-gray-200">
            <div>
              <label class="block text-sm font-medium text-gray-500">ID</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ selectedLogData.id }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Action</label>
              <div class="mt-1 text-sm font-semibold text-indigo-700">{{ selectedLogData.action }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">User</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ getUserName(selectedLogData) }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Entity</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ getEntityName(selectedLogData) }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">IP Address</label>
              <div class="mt-1 text-sm text-gray-900">{{ selectedLogData.ip_address }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Date & Time</label>
              <div class="mt-1 text-sm text-gray-900">{{ formatDate(selectedLogData.created_at) }}</div>
            </div>
          </div>
          
          <div class="pt-4">
            <label class="block text-sm font-medium text-gray-500 mb-3">Details</label>
            
            <div v-if="!selectedLogData.details" class="text-sm text-gray-500 italic bg-gray-50 p-4 rounded-md border border-gray-200">
              No additional details available
            </div>
            
            <div v-else class="bg-gray-50 p-4 rounded-md border border-gray-200 font-mono text-sm overflow-auto max-h-96 shadow-inner">
              <pre>{{ JSON.stringify(parseDetails(selectedLogData.details), null, 2) }}</pre>
            </div>
          </div>
        </div>
        
        <template #footer>
          <div class="flex justify-end gap-3 bg-gray-50 px-4 py-3 border-t border-gray-200 rounded-b-lg">
            <BaseButton 
              variant="outline" 
              @click="closeDetailModal"
              class="w-full sm:w-auto"
            >
              Close
            </BaseButton>
          </div>
        </template>
      </BaseModal>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { nextTick } from 'vue';
import { api } from '@/services/api'; // Import your API service
import AppBreadcrumb from '@/components/layout/AppBreadcrumb.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseCard from '@/components/common/BaseCard.vue';
import BaseModal from '@/components/common/BaseModal.vue';

// State
const logs = ref([]);
const apiResponse = ref(null);
const isLoading = ref(false);
const showDebug = ref(true);
const selectedLogId = ref(null);
const selectedLogData = ref(null);
const currentEndpoint = ref('/api/v1/auditlogs');
const debounceTimer = ref(null);

// Filter state
const filters = ref({
  action: '',
  status: '',
  user: '',
  dateRange: '',
  startDate: '',
  endDate: ''
});

// Computed values
const breadcrumbItems = computed(() => [
  { name: 'Dashboard', path: '/superadmin' },
  { name: 'Audit Logs', path: '/superadmin/audit-logs' }
]);

const hasActiveFilters = computed(() => {
  return filters.value.action || 
         filters.value.status || 
         filters.value.user || 
         filters.value.dateRange;
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

  // Filter by user
  if (filters.value.user) {
    filtered = filtered.filter(log => {
      const userName = getUserName(log).toLowerCase();
      return userName.includes(filters.value.user.toLowerCase());
    });
  }

  // Filter by date range
  if (filters.value.dateRange) {
    const now = new Date();
    let startDate, endDate;

    switch (filters.value.dateRange) {
      case 'today':
        startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate() + 1);
        break;
      case 'yesterday':
        startDate = new Date(now.getFullYear(), now.getMonth(), now.getDate() - 1);
        endDate = new Date(now.getFullYear(), now.getMonth(), now.getDate());
        break;
      case 'last7days':
        startDate = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000);
        endDate = now;
        break;
      case 'last30days':
        startDate = new Date(now.getTime() - 30 * 24 * 60 * 60 * 1000);
        endDate = now;
        break;
      case 'last90days':
        startDate = new Date(now.getTime() - 90 * 24 * 60 * 60 * 1000);
        endDate = now;
        break;
      case 'custom':
        if (filters.value.startDate) {
          startDate = new Date(filters.value.startDate);
        }
        if (filters.value.endDate) {
          endDate = new Date(filters.value.endDate);
          endDate.setDate(endDate.getDate() + 1); // Include the end date
        }
        break;
    }

    if (startDate || endDate) {
      filtered = filtered.filter(log => {
        const logDate = new Date(log.created_at);
        if (startDate && logDate < startDate) return false;
        if (endDate && logDate >= endDate) return false;
        return true;
      });
    }
  }

  return filtered;
});

// Methods
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

function formatDateRangeLabel(range) {
  const labels = {
    'today': 'Today',
    'yesterday': 'Yesterday',
    'last7days': 'Last 7 Days',
    'last30days': 'Last 30 Days',
    'last90days': 'Last 90 Days',
    'custom': 'Custom Range'
  };
  return labels[range] || range;
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
  return details.entity_name || (log.entity_id ? `Entity ${log.entity_id}` : '-');
}

async function fetchAuditLogs() {
  isLoading.value = true;
  try {
    const response = await api.get(currentEndpoint.value);
    console.log('API Response:', response.data); // Check this in the console
    apiResponse.value = response.data;
    logs.value = response.data.data || response.data || []; // Adjust based on actual structure
    await nextTick(); // Ensure reactive update
  } catch (error) {
    console.error('Error fetching audit logs:', error);
    logs.value = [];
  } finally {
    isLoading.value = false;
    console.log('Fetch completed, logs:', logs.value);
  }
}

async function tryEndpoint(endpoint) {
  currentEndpoint.value = endpoint;
  console.log('Trying endpoint:', endpoint);
  await fetchAuditLogs();
}

async function fetchLogDetails(id) {
  try {
    const response = await api.get(`${currentEndpoint.value}/${id}`);
    console.log('Log details response:', response.data);
    
    if (response.data && response.data.data) {
      return response.data.data;
    } else if (response.data) {
      return response.data;
    }
    return null;
  } catch (error) {
    console.error('Error fetching log details:', error);
    return null;
  }
}

function selectLog(id) {
  selectedLogId.value = id;
  selectedLogData.value = logs.value.find(log => log.id === id);
  
  // Optionally fetch more details if needed
  /*
  fetchLogDetails(id).then(detailedLog => {
    if (detailedLog) {
      selectedLogData.value = detailedLog;
    }
  });
  */
}

function closeDetailModal() {
  selectedLogId.value = null;
  selectedLogData.value = null;
}

// Filter methods
function applyFilters() {
  // Filters are applied automatically through computed property
  console.log('Filters applied:', filters.value);
}

function debounceFilter() {
  if (debounceTimer.value) {
    clearTimeout(debounceTimer.value);
  }
  debounceTimer.value = setTimeout(() => {
    applyFilters();
  }, 300);
}

function clearFilters() {
  filters.value = {
    action: '',
    status: '',
    user: '',
    dateRange: '',
    startDate: '',
    endDate: ''
  };
  applyFilters();
}

function clearFilter(filterKey) {
  filters.value[filterKey] = '';
  if (filterKey === 'dateRange') {
    filters.value.startDate = '';
    filters.value.endDate = '';
  }
  applyFilters();
}

// Initialize
onMounted(() => {
  fetchAuditLogs();
});
</script>



