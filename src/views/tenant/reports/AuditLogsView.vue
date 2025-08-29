<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6 flex justify-between items-center">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Audit Logs Report</h1>
          <p class="text-gray-600 mt-1">
            Monitor audit activities within your temple system
            <span v-if="currentEntityId && currentEntityId !== 'all'" class="text-indigo-600 font-medium">
              (Entity ID: {{ currentEntityId }})
            </span>
          </p>
        </div>
        <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200 flex items-center space-x-4">
          <span class="text-indigo-800 font-medium">{{ adminName }}</span>
          <span class="text-indigo-600 text-sm">({{ adminRole }})</span>
        </div>
      </div>
    </div>

    <!-- Filters and Download Card -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8 bg-white rounded-2xl shadow-sm border border-gray-200 mt-6">
      <h3 class="text-xl font-bold mb-4">Filter Audit Logs</h3>

      <!-- Temple Selection -->
      <div class="mb-4">
        <label class="block text-gray-700 font-medium mb-2">Select Temple</label>
        <div class="relative">
          <select v-model="currentEntityId" :disabled="loading || downloading || loadingEntities" class="w-full px-3 py-2 border border-gray-300 rounded-md">
            <option value="all">All Temples</option>
            <option v-for="entity in entities" :key="entity.id" :value="entity.id">{{ entity.name }}</option>
          </select>
          <div v-if="loadingEntities" class="absolute right-3 top-3">
            <svg class="animate-spin h-4 w-4 text-gray-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Action Filter -->
      <div class="mb-4">
        <label class="block text-gray-700 font-medium mb-2">Action Type</label>
        <select v-model="filters.action" @change="fetchPreview" :disabled="loading || downloading" class="w-full px-3 py-2 border border-gray-300 rounded-md">
          <option value="">All Actions</option>
          <!-- Auth Module Actions -->
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
          
          <!-- Notification Module Actions -->
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

      <!-- Status Filter -->
      <div class="mb-4">
        <label class="block text-gray-700 font-medium mb-2">Status</label>
        <div class="flex space-x-2">
          <button
            v-for="status in statusOptions"
            :key="status.value"
            @click="filters.status = status.value; fetchPreview()"
            :class="[
              'px-4 py-2 rounded-md text-sm font-medium',
              filters.status === status.value ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            ]"
            :disabled="loading || downloading"
          >
            {{ status.label }}
          </button>
        </div>
      </div>

      <!-- Date Range -->
      <div class="mb-4">
        <label class="block text-gray-700 font-medium mb-2">Date Range</label>
        <div class="flex space-x-2">
          <button
            v-for="filter in timeFilters"
            :key="filter.value"
            @click="setActiveFilter(filter.value)"
            :class="[
              'px-4 py-2 rounded-md text-sm font-medium',
              activeFilter === filter.value ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'
            ]"
            :disabled="loading || downloading"
          >
            {{ filter.label }}
          </button>
        </div>
      </div>

      <!-- Custom Date Range Inputs -->
      <div v-if="activeFilter === 'custom'" class="mb-4 grid grid-cols-1 md:grid-cols-2 gap-4">
        <div>
          <label class="block text-gray-700 font-medium mb-2">Start Date</label>
          <input
            type="date"
            v-model="filters.startDate"
            @change="fetchPreview"
            :disabled="loading || downloading"
            class="w-full px-3 py-2 border border-gray-300 rounded-md"
          />
        </div>
        <div>
          <label class="block text-gray-700 font-medium mb-2">End Date</label>
          <input
            type="date"
            v-model="filters.endDate"
            @change="fetchPreview"
            :disabled="loading || downloading"
            class="w-full px-3 py-2 border border-gray-300 rounded-md"
          />
        </div>
      </div>

      <!-- Format Selection and Download -->
      <div class="flex flex-wrap items-center justify-between pt-4 border-t border-gray-200">
        <div class="w-40">
          <label class="block text-gray-700 font-medium mb-2">Export Format</label>
          <select v-model="selectedFormat" :disabled="loading || downloading" class="w-full px-3 py-2 border border-gray-300 rounded-md">
            <option v-for="format in formats" :key="format.value" :value="format.value">{{ format.label }}</option>
          </select>
        </div>
        <button
          @click="downloadReport"
          :disabled="loading || downloading"
          class="mt-4 md:mt-0 inline-flex items-center px-5 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <svg
            v-if="downloading"
            class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
            fill="none"
            viewBox="0 0 24 24"
            stroke="currentColor"
          >
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="mr-2 h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          {{ downloading ? 'Downloading...' : 'Download Report' }}
        </button>
      </div>
    </div>

    <!-- Applied Filters Summary -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6 bg-white rounded-2xl shadow-sm border border-gray-200 p-6">
      <h3 class="text-lg font-medium mb-4">Applied Filters</h3>
      <div class="flex flex-wrap gap-3">
        <span class="inline-flex items-center bg-indigo-100 text-indigo-800 rounded-full px-3 py-1 text-sm font-semibold">
          Action: {{ getActionLabel(filters.action) || 'All' }}
        </span>
        <span class="inline-flex items-center bg-indigo-100 text-indigo-800 rounded-full px-3 py-1 text-sm font-semibold">
          Status: {{ filters.status || 'All' }}
        </span>
        <span class="inline-flex items-center bg-indigo-100 text-indigo-800 rounded-full px-3 py-1 text-sm font-semibold">
          Date Range: {{ activeFilterLabel }}
          <span v-if="activeFilter === 'custom'">
            ({{ filters.startDate }} - {{ filters.endDate }})
          </span>
        </span>
        <span class="inline-flex items-center bg-indigo-100 text-indigo-800 rounded-full px-3 py-1 text-sm font-semibold">
          Entity: {{ currentEntityLabel }}
        </span>
        <span class="inline-flex items-center bg-indigo-100 text-indigo-800 rounded-full px-3 py-1 text-sm font-semibold">
          Format: {{ selectedFormatLabel }}
        </span>
      </div>
    </div>

    <!-- Report Preview Table -->
    <div v-if="previewData.length" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6 bg-white rounded-lg shadow-sm border border-gray-200 overflow-auto">
      <div class="p-4 border-b border-gray-200">
        <h3 class="text-lg font-medium">Audit Logs Preview ({{ previewData.length }} records)</h3>
      </div>
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-50">
            <tr>
              <th v-for="column in previewColumns" :key="column.key" class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                {{ column.label }}
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="(row, idx) in previewData.slice(0, 100)" :key="idx" class="hover:bg-gray-50">
              <td v-for="column in previewColumns" :key="column.key" class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                {{ formatCellValue(row[column.key]) }}
              </td>
            </tr>
          </tbody>
        </table>
      </div>
      <div v-if="previewData.length > 100" class="p-4 border-t border-gray-200 bg-gray-50">
        <p class="text-sm text-gray-600 text-center">
          Showing first 100 records out of {{ previewData.length }}. Download the full report to see all records.
        </p>
      </div>
    </div>

    <!-- No Data State -->
    <div v-else-if="!loading && previewColumns.length === 0" class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 mt-6 bg-white rounded-lg shadow-sm border border-gray-200 p-12 text-center">
      <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
        <svg class="w-12 h-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No Audit Logs Found</h3>
      <p class="text-gray-600">No audit logs match the selected filters. Try adjusting your filter criteria.</p>
    </div>

    <!-- Loading state -->
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-30 flex items-center justify-center z-50">
      <div class="bg-white p-6 rounded-lg shadow-lg flex items-center space-x-3">
        <svg class="animate-spin h-6 w-6 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"></path>
        </svg>
        <span class="font-medium text-gray-900">Loading audit logs...</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import ReportsService from '@/services/reports.service'
import { useToast } from '@/composables/useToast'

// State vars
const entities = ref([])
const loadingEntities = ref(false)

const currentEntityId = ref('all')
const filters = ref({
  action: '',
  status: '',
  startDate: '',
  endDate: ''
})
const activeFilter = ref('weekly')
const loading = ref(false)
const downloading = ref(false)
const selectedFormat = ref('pdf')
const previewData = ref([])
const previewColumns = ref([])

const adminName = 'Temple Admin'
const adminRole = 'Administrator'
const { showToast } = useToast()

const actionTypes = [
  // Auth Module Actions
  { value: 'REGISTRATION_SUCCESS', label: 'Registration Success', group: 'Authentication' },
  { value: 'REGISTRATION_FAILED', label: 'Registration Failed', group: 'Authentication' },
  { value: 'TEMPLEADMIN_REGISTRATION_SUCCESS', label: 'Temple Admin Registration', group: 'Authentication' },
  { value: 'REGISTRATION_BLOCKED', label: 'Registration Blocked', group: 'Authentication' },
  { value: 'LOGIN_SUCCESS', label: 'Login Success', group: 'Authentication' },
  { value: 'LOGIN_FAILED', label: 'Login Failed', group: 'Authentication' },
  { value: 'LOGOUT', label: 'Logout', group: 'Authentication' },
  { value: 'PASSWORD_RESET_REQUESTED', label: 'Password Reset Requested', group: 'Authentication' },
  { value: 'PASSWORD_RESET_SUCCESS', label: 'Password Reset Success', group: 'Authentication' },
  { value: 'PASSWORD_RESET_FAILED', label: 'Password Reset Failed', group: 'Authentication' },

  // Entity Module Actions
  { value: 'TEMPLE_CREATED', label: 'Temple Created', group: 'Entity' },
  { value: 'TEMPLE_UPDATED', label: 'Temple Updated', group: 'Entity' },
  { value: 'TEMPLE_CREATE_FAILED', label: 'Temple Create Failed', group: 'Entity' },
  { value: 'TEMPLE_UPDATE_FAILED', label: 'Temple Update Failed', group: 'Entity' },

  // Events Module Actions
  { value: 'EVENT_CREATED', label: 'Event Created', group: 'Events' },
  { value: 'EVENT_UPDATED', label: 'Event Updated', group: 'Events' },
  { value: 'EVENT_DELETED', label: 'Event Deleted', group: 'Events' },

  // Seva Module Actions
  { value: 'SEVA_CREATED', label: 'Seva Created', group: 'Seva' },
  { value: 'SEVA_UPDATED', label: 'Seva Updated', group: 'Seva' },
  { value: 'SEVA_BOOKED', label: 'Seva Booked', group: 'Seva' },
  { value: 'SEVA_BOOKING_APPROVED', label: 'Seva Booking Approved', group: 'Seva' },
  { value: 'SEVA_BOOKING_REJECTED', label: 'Seva Booking Rejected', group: 'Seva' },

  // Donations Module Actions
  { value: 'DONATION_INITIATED', label: 'Donation Initiated', group: 'Donations' },
  { value: 'DONATION_SUCCESS', label: 'Donation Success', group: 'Donations' },
  { value: 'DONATION_FAILED', label: 'Donation Failed', group: 'Donations' },
  { value: 'DONATION_VERIFICATION_FAILED', label: 'Donation Verification Failed', group: 'Donations' },

  // Notification Module Actions
  { value: 'TEMPLATE_CREATED', label: 'Template Created', group: 'Notifications' },
  { value: 'TEMPLATE_UPDATED', label: 'Template Updated', group: 'Notifications' },
  { value: 'TEMPLATE_DELETED', label: 'Template Deleted', group: 'Notifications' },
  { value: 'EMAIL_SENT', label: 'Email Sent', group: 'Notifications' },
  { value: 'SMS_SENT', label: 'SMS Sent', group: 'Notifications' },
  { value: 'WHATSAPP_SENT', label: 'WhatsApp Sent', group: 'Notifications' },

  // User Profile Module Actions
  { value: 'PROFILE_CREATED', label: 'Profile Created', group: 'User Profile' },
  { value: 'PROFILE_UPDATED', label: 'Profile Updated', group: 'User Profile' },
  { value: 'DEVOTEE_JOINED_TEMPLE', label: 'Devotee Joined Temple', group: 'User Profile' },
  { value: 'VOLUNTEER_JOINED_TEMPLE', label: 'Volunteer Joined Temple', group: 'User Profile' },

  // Reports Module Actions
  { value: 'DEVOTEE_BIRTHDAYS_REPORT_VIEWED', label: 'Birthday Report Viewed', group: 'Reports' },
  { value: 'DEVOTEE_BIRTHDAYS_REPORT_DOWNLOADED', label: 'Birthday Report Downloaded', group: 'Reports' },
  { value: 'TEMPLE_REGISTER_REPORT_VIEWED', label: 'Temple Register Report Viewed', group: 'Reports' },
  { value: 'TEMPLE_REGISTER_REPORT_DOWNLOADED', label: 'Temple Register Report Downloaded', group: 'Reports' },
  { value: 'TEMPLE_ACTIVITIES_REPORT_VIEWED', label: 'Activities Report Viewed', group: 'Reports' },
  { value: 'TEMPLE_ACTIVITIES_REPORT_DOWNLOADED', label: 'Activities Report Downloaded', group: 'Reports' }
]

const statusOptions = [
  { value: '', label: 'All' },
  { value: 'success', label: 'Success' },
  { value: 'failure', label: 'Failure' },
]

const timeFilters = [
  { value: 'daily', label: 'Daily' },
  { value: 'weekly', label: 'Weekly' },
  { value: 'monthly', label: 'Monthly' },
  { value: 'custom', label: 'Custom Range' },
]

const formats = [
  { value: 'pdf', label: 'PDF' },
  { value: 'csv', label: 'CSV' },
  { value: 'excel', label: 'Excel' },
]

// Fetch temples/entities from API
async function fetchEntities() {
  loadingEntities.value = true
  try {
    // Try multiple possible API endpoints/methods
    let response
    
    // Option 1: Try ReportsService method (adjust method name as needed)
    try {
      response = await ReportsService.getTemples()
    } catch (e1) {
      console.log('ReportsService.getTemples() failed, trying alternatives...')
      
      // Option 2: Try different method names
      try {
        response = await ReportsService.getEntities()
      } catch (e2) {
        // Option 3: Try another common method name
        response = await ReportsService.getTenantEntities()
      }
    }
    
    console.log('Raw API Response:', response)
    
    // Handle different response structures
    let temples = []
    
    if (response?.data) {
      // Try different nested structures
      temples = response.data.data || response.data.temples || response.data.entities || response.data || []
    } else if (response?.temples) {
      temples = response.temples
    } else if (response?.entities) {
      temples = response.entities
    } else if (Array.isArray(response)) {
      temples = response
    }
    
    console.log('Extracted temples:', temples)
    
    // Map the response to the format needed for the dropdown
    if (Array.isArray(temples) && temples.length > 0) {
      entities.value = temples.map(temple => {
        console.log('Processing temple:', temple)
        return {
          id: temple.id || temple._id || temple.temple_id || temple.entityId,
          name: temple.name || temple.temple_name || temple.title || temple.entity_name || `Temple ${temple.id}`
        }
      })
      
      console.log('Mapped entities:', entities.value)
      showToast(`${entities.value.length} temples loaded successfully`, 'success')
    } else {
      console.warn('No temples found in response')
      entities.value = []
      showToast('No temples found', 'warning')
    }
    
  } catch (error) {
    console.error('Error fetching temples:', error)
    showToast('Failed to load temples: ' + error.message, 'error')
    
    // Fallback - you can add some default entities for testing
    entities.value = [
      { id: '1', name: 'Default Temple 1' },
      { id: '2', name: 'Default Temple 2' }
    ]
  } finally {
    loadingEntities.value = false
  }
}

function setActiveFilter(filter) {
  activeFilter.value = filter

  const today = new Date()
  let startDate, endDate

  if (filter === 'daily') {
    startDate = endDate = today
  } else if (filter === 'weekly') {
    startDate = new Date(today)
    startDate.setDate(startDate.getDate() - 7)
    endDate = today
  } else if (filter === 'monthly') {
    startDate = new Date(today)
    startDate.setMonth(startDate.getMonth() - 1)
    endDate = today
  }

  if (filter !== 'custom') {
    filters.value.startDate = startDate.toISOString().split('T')[0]
    filters.value.endDate = endDate.toISOString().split('T')[0]
  }

  fetchPreview()
}

const currentEntityLabel = computed(() => {
  if (currentEntityId.value === 'all') return 'All Entities'
  const e = entities.value.find(ent => ent.id === currentEntityId.value)
  return e ? e.name : 'Unknown Entity'
})

const activeFilterLabel = computed(() => {
  const found = timeFilters.find(f => f.value === activeFilter.value)
  return found ? found.label : activeFilter.value
})

const selectedFormatLabel = computed(() => {
  const found = formats.find(f => f.value === selectedFormat.value)
  return found ? found.label : selectedFormat.value
})

// Helper function to get action label
function getActionLabel(actionValue) {
  if (!actionValue) return null
  const action = actionTypes.find(a => a.value === actionValue)
  return action ? action.label : actionValue
}

// Helper function to format cell values
function formatCellValue(value) {
  if (value === null || value === undefined) return '-'
  if (typeof value === 'boolean') return value ? 'Yes' : 'No'
  if (typeof value === 'object') return JSON.stringify(value)
  // Format dates if they look like ISO strings
  if (typeof value === 'string' && /^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}/.test(value)) {
    return new Date(value).toLocaleString()
  }
  return String(value)
}

async function fetchPreview() {
  loading.value = true
  try {
    const params = {
      entityId: currentEntityId.value,
      actionType: filters.value.action || undefined,
      dateRange: activeFilter.value,
      startDate: activeFilter.value === 'custom' ? filters.value.startDate : undefined,
      endDate: activeFilter.value === 'custom' ? filters.value.endDate : undefined,
      status: filters.value.status || undefined,
    }

    const response = await ReportsService.getAuditLogsReport(params)
    let rows = response.data?.data || response.data || []
    // Fallback if data wrapped differently
    if (rows.audit_logs) rows = rows.audit_logs

    previewData.value = rows

    // Set columns from one of the rows or fallback columns
    if (rows.length > 0) {
      previewColumns.value = Object.keys(rows[0]).map(key => ({
        key,
        label: key.replace(/_/g, ' ').toUpperCase()
      }))
    } else {
      previewColumns.value = []
    }
  } catch (error) {
    console.error('Error fetching audit logs:', error)
    showToast('Failed to load audit logs preview', 'error')
    previewData.value = []
    previewColumns.value = []
  } finally {
    loading.value = false
  }
}

async function downloadReport() {
  downloading.value = true
  try {
    const params = {
      entityId: currentEntityId.value,
      dateRange: activeFilter.value,
      startDate: activeFilter.value === 'custom' ? filters.value.startDate : undefined,
      endDate: activeFilter.value === 'custom' ? filters.value.endDate : undefined,
      format: selectedFormat.value,
      actionType: filters.value.action || undefined,
      status: filters.value.status || undefined,
    }

    await ReportsService.downloadAuditLogsReport(params)
    showToast('Audit logs report downloaded successfully', 'success')
  } catch (error) {
    console.error('Error downloading report:', error)
    showToast('Failed to download report', 'error')
  } finally {
    downloading.value = false
  }
}

onMounted(async () => {
  // Debug: Check what methods are available in ReportsService
  console.log('Available ReportsService methods:', Object.getOwnPropertyNames(ReportsService))
  console.log('ReportsService object:', ReportsService)
  
  await fetchEntities() // Fetch temples on component mount
  setActiveFilter('weekly')
})

watch([() => filters.value.action, () => filters.value.status, () => currentEntityId.value], () => {
  fetchPreview()
})
</script>