<template>
  <div class="space-y-6">
    <!-- Debug Info (remove in production) -->
    <div class="bg-gray-100 p-4 rounded-lg mb-4 text-xs font-mono overflow-auto max-h-40" v-if="debugMode">
      <div class="mb-2 font-bold">Debug Information:</div>
      <div v-if="selectedTenant">Selected Tenant: ID={{ selectedTenant.id || selectedTenant.ID }}, Name={{ getTenantName(selectedTenant) }}</div>
      <div>Tenant Count: {{ allTenants.length }}</div>
      <div class="mt-2 flex gap-2">
        <button @click="debugTenantData" class="px-3 py-1 bg-gray-200 rounded text-xs">Run API Debug</button>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="flex flex-wrap gap-3">
      <div class="bg-yellow-50 border border-yellow-200 rounded-xl px-4 py-3 text-center min-w-[100px] flex-1">
        <div class="text-2xl font-bold text-yellow-800">{{ pendingCount }}</div>
        <div class="text-xs text-yellow-600 font-medium">Pending</div>
      </div>
      <div class="bg-green-50 border border-green-200 rounded-xl px-4 py-3 text-center min-w-[100px] flex-1">
        <div class="text-2xl font-bold text-green-800">{{ approvedCount }}</div>
        <div class="text-xs text-green-600 font-medium">Approved</div>
      </div>
      <div class="bg-red-50 border border-red-200 rounded-xl px-4 py-3 text-center min-w-[100px] flex-1">
        <div class="text-2xl font-bold text-red-800">{{ rejectedCount }}</div>
        <div class="text-xs text-red-600 font-medium">Rejected</div>
      </div>
    </div>

    <!-- Filter & Refresh -->
    <div class="flex flex-col sm:flex-row justify-between gap-3">
      <!-- Filters -->
      <div class="flex gap-3">
        <select 
          v-model="statusFilter" 
          class="px-3 py-2 bg-white border border-gray-300 rounded-lg text-sm w-full sm:w-auto"
          @change="applyFilters"
        >
          <option value="">View All</option>
          <option value="pending">Pending</option>
          <option value="approved">Approved</option>
          <option value="rejected">Rejected</option>
        </select>
        
        <!-- Debug Mode Toggle -->
        <button 
          @click="debugMode = !debugMode" 
          class="px-3 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg text-sm"
        >
          {{ debugMode ? 'Hide Debug' : 'Debug Mode' }}
        </button>
      </div>
      
      <!-- Refresh Button -->
      <button 
        @click="loadTenants" 
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm"
      >
        Refresh Data
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Tenant Applications List -->
    <div v-if="!loading && filteredTenants.length > 0" class="space-y-4">
      <div
        v-for="(tenant, idx) in paginatedTenants"
        :key="tenant.id || tenant.ID || idx"
        class="bg-white rounded-xl shadow-md border border-gray-200 hover:shadow-lg transition-all duration-200"
      >
        <div class="p-6">
          <!-- Tenant ID Debug (remove in production) -->
          <div v-if="debugMode" class="bg-gray-100 p-2 mb-3 rounded text-xs">
            <span class="font-bold">ID:</span> {{ tenant.id || tenant.ID || 'Not available' }} |
            <span class="font-bold">Status:</span> {{ tenant.status || tenant.Status || 'Not available' }}
          </div>
          
          <!-- Header Row -->
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-4">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center flex-shrink-0">
                <span class="text-indigo-600 font-semibold text-lg">
                  {{ getTenantInitial(tenant) }}
                </span>
              </div>
              
              <!-- Basic Info - Only name and email -->
              <div class="flex-1 min-w-0">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">
                  {{ getTenantName(tenant) }}
                </h3>
                <p class="text-sm text-gray-600">{{ getTenantEmail(tenant) }}</p>
              </div>
            </div>

            <!-- Status & Date -->
            <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
              <span
                :class="getStatusBadgeClass(tenant.status || tenant.Status)"
                class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium"
              >
                {{ tenant.status || tenant.Status || 'pending' }}
              </span>
              <span class="text-xs text-gray-500">
                Applied: {{ formatDate(tenant.created_at || tenant.CreatedAt) }}
              </span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div v-if="isStatusPending(tenant)" class="flex flex-col sm:flex-row gap-3">
            <button
              @click="handleApprove(tenant)"
              class="flex-1 sm:flex-none px-6 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
              :disabled="isProcessing"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
              {{ isProcessing ? 'Processing...' : 'Approve' }}
            </button>
            
            <button
              @click="handleRejectClick(tenant)"
              class="flex-1 sm:flex-none px-6 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
              :disabled="isProcessing"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
              {{ isProcessing ? 'Processing...' : 'Reject' }}
            </button>
            
            <!-- Direct API call debug button (remove in production) -->
            <button 
              v-if="debugMode"
              @click="testApprovalApi(tenant)"
              class="flex-1 sm:flex-none px-6 py-2 bg-gray-600 text-white text-sm font-medium rounded-xl flex items-center justify-center gap-2"
            >
              Test API Call
            </button>
          </div>
        </div>
      </div>

      <!-- Pagination Controls -->
      <div class="flex flex-col sm:flex-row items-center justify-between bg-white p-4 rounded-xl shadow-sm border border-gray-200">
        <div class="text-sm text-gray-700 mb-3 sm:mb-0">
          Showing <span class="font-medium">{{ paginationStart }}</span> to 
          <span class="font-medium">{{ paginationEnd }}</span> of 
          <span class="font-medium">{{ filteredTenants.length }}</span> results
        </div>
        
        <div class="flex space-x-2">
          <button 
            @click="currentPage--" 
            :disabled="currentPage === 1"
            class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          
          <div class="flex space-x-1">
            <button 
              v-for="page in displayedPageNumbers" 
              :key="page"
              @click="goToPage(page)"
              :class="[
                'inline-flex items-center px-3 py-2 border text-sm font-medium rounded-md',
                currentPage === page 
                  ? 'border-indigo-500 bg-indigo-50 text-indigo-600'
                  : 'border-gray-300 bg-white text-gray-700 hover:bg-gray-50'
              ]"
            >
              {{ page }}
            </button>
          </div>
          
          <button 
            @click="currentPage++" 
            :disabled="currentPage === totalPages"
            class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm font-medium rounded-md bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!loading" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
      </svg>
      <h3 class="mt-4 text-lg font-medium text-gray-900">No tenant applications found</h3>
      <p class="mt-2 text-sm text-gray-500">
        {{ statusFilter ? `No tenants with "${statusFilter}" status found` : 'Try refreshing or checking back later' }}
      </p>
    </div>

    <!-- Rejection Modal -->
    <div v-if="showRejectModal" class="fixed inset-0 z-50 overflow-y-auto bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-xl p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">Reject Tenant</h3>
        
        <!-- Debug info -->
        <div v-if="debugMode" class="bg-gray-100 p-2 mb-3 rounded text-xs">
          <div><span class="font-bold">Selected Tenant ID:</span> {{ selectedTenant?.id || selectedTenant?.ID }}</div>
          <div><span class="font-bold">Selected Tenant Name:</span> {{ getTenantName(selectedTenant) }}</div>
        </div>
        
        <p class="mb-4">Please provide a reason for rejecting <span class="font-medium">{{ selectedTenant ? getTenantName(selectedTenant) : '' }}</span>:</p>
        <textarea 
          v-model="rejectReason" 
          class="w-full border border-gray-300 rounded-lg p-3 mb-4 min-h-[100px] focus:ring-2 focus:ring-indigo-500 focus:border-transparent"
          placeholder="Enter rejection reason..."
        ></textarea>
        <div class="flex justify-end gap-3">
          <button 
            @click="showRejectModal = false"
            class="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg"
            :disabled="isProcessing"
          >
            Cancel
          </button>
          <button 
            @click="confirmReject"
            :disabled="!rejectReason.trim() || isProcessing"
            class="px-4 py-2 bg-red-600 hover:bg-red-700 text-white rounded-lg disabled:opacity-50"
          >
            {{ isProcessing ? 'Processing...' : 'Reject' }}
          </button>
        </div>
      </div>
    </div>

    <!-- API Debug Modal -->
    <div v-if="showApiDebugModal" class="fixed inset-0 z-50 overflow-y-auto bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-xl p-6 w-full max-w-3xl mx-4 max-h-[80vh] overflow-auto">
        <h3 class="text-lg font-semibold mb-4">API Debug Information</h3>
        <div class="bg-gray-100 p-4 rounded font-mono text-xs overflow-auto whitespace-pre">
          {{ apiDebugInfo }}
        </div>
        <div class="mt-4 flex justify-end">
          <button 
            @click="showApiDebugModal = false"
            class="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from '@/composables/useToast'
import superAdminService from '@/services/superadmin.service'
import api from '@/services/api'

export default {
  name: 'TenantApprovals',
  emits: ['updated'],
  setup(props, { emit }) {
    // Data
    const loading = ref(true)
    const tenants = ref([])
    const allTenants = ref([]) // Store all tenants for "View All" functionality
    const toast = useToast()
    const isProcessing = ref(false)
    
    // Debug mode
    const debugMode = ref(false)
    const showApiDebugModal = ref(false)
    const apiDebugInfo = ref('')
    
    // Simple rejection modal
    const showRejectModal = ref(false)
    const selectedTenant = ref(null)
    const rejectReason = ref('')
    
    // Filters
    const statusFilter = ref('')
    
    // Pagination
    const currentPage = ref(1)
    const itemsPerPage = ref(5)
    
    // Computed properties
    const pendingCount = computed(() => 
      allTenants.value.filter(t => isStatusPending(t)).length
    )
    
    const approvedCount = computed(() => 
      allTenants.value.filter(t => isStatusApproved(t)).length
    )
    
    const rejectedCount = computed(() => 
      allTenants.value.filter(t => isStatusRejected(t)).length
    )
    
    // Status check helpers - Updated to check both uppercase and lowercase properties
    const isStatusPending = (tenant) => {
      if (!tenant) return false;
      const status = (tenant.status || tenant.Status || '').toLowerCase();
      return status === 'pending' || status === 'pending approval' || status === '';
    }
    
    const isStatusApproved = (tenant) => {
      if (!tenant) return false;
      const status = (tenant.status || tenant.Status || '').toLowerCase();
      return status === 'active' || status === 'approved';
    }
    
    const isStatusRejected = (tenant) => {
      if (!tenant) return false;
      const status = (tenant.status || tenant.Status || '').toLowerCase();
      return status === 'rejected' || status === 'declined';
    }
    
    // Filtered tenants based on status
    const filteredTenants = computed(() => {
      let sourceData = statusFilter.value === '' ? allTenants.value : tenants.value;
      
      if (!statusFilter.value) {
        return allTenants.value; // Show all tenants when "View All" is selected
      }
      
      switch (statusFilter.value.toLowerCase()) {
        case 'pending':
          return allTenants.value.filter(tenant => isStatusPending(tenant));
        case 'approved':
          return allTenants.value.filter(tenant => isStatusApproved(tenant));
        case 'rejected':
          return allTenants.value.filter(tenant => isStatusRejected(tenant));
        default:
          return allTenants.value;
      }
    })
    
    // Paginated tenants
    const paginatedTenants = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value
      const end = start + itemsPerPage.value
      return filteredTenants.value.slice(start, end)
    })
    
    // Pagination helpers
    const totalPages = computed(() => {
      return Math.ceil(filteredTenants.value.length / itemsPerPage.value) || 1
    })
    
    const paginationStart = computed(() => {
      return filteredTenants.value.length === 0 
        ? 0 
        : (currentPage.value - 1) * itemsPerPage.value + 1
    })
    
    const paginationEnd = computed(() => {
      return Math.min(currentPage.value * itemsPerPage.value, filteredTenants.value.length)
    })
    
    const displayedPageNumbers = computed(() => {
      const maxDisplayed = 5
      const pages = []
      
      if (totalPages.value <= maxDisplayed) {
        for (let i = 1; i <= totalPages.value; i++) {
          pages.push(i)
        }
      } else {
        pages.push(1)
        
        let startPage = Math.max(2, currentPage.value - 1)
        let endPage = Math.min(totalPages.value - 1, currentPage.value + 1)
        
        if (currentPage.value <= 2) {
          endPage = 3
        }
        
        if (currentPage.value >= totalPages.value - 1) {
          startPage = totalPages.value - 2
        }
        
        if (startPage > 2) {
          pages.push('...')
        }
        
        for (let i = startPage; i <= endPage; i++) {
          pages.push(i)
        }
        
        if (endPage < totalPages.value - 1) {
          pages.push('...')
        }
        
        if (totalPages.value > 1) {
          pages.push(totalPages.value)
        }
      }
      
      return pages
    })
    
    // Helper methods for displaying tenant information
    const getTenantName = (tenant) => {
      if (!tenant) return 'Unknown Tenant';
      
      // Check both uppercase and lowercase properties
      if (tenant.full_name) return tenant.full_name;
      if (tenant.FullName) return tenant.FullName;
      if (tenant.fullName) return tenant.fullName;
      if (tenant.name) return tenant.name;
      if (tenant.Name) return tenant.Name;
      
      // Check email as fallback
      if (tenant.email) return tenant.email.split('@')[0];
      if (tenant.Email) return tenant.Email.split('@')[0];
      
      // Return ID if available
      return (tenant.id || tenant.ID) ? `User ${tenant.id || tenant.ID}` : 'Unknown Tenant';
    }
    
    const getTenantEmail = (tenant) => {
      if (!tenant) return 'No email provided';
      
      if (tenant.email) return tenant.email;
      if (tenant.Email) return tenant.Email;
      if (tenant.email_address) return tenant.email_address;
      if (tenant.EmailAddress) return tenant.EmailAddress;
      
      return 'No email provided';
    }
    
    const getTenantInitial = (tenant) => {
      if (!tenant) return 'T';
      
      const name = getTenantName(tenant);
      if (name && name !== 'Unknown Tenant' && name.length > 0) {
        return name.charAt(0).toUpperCase();
      }
      
      const email = getTenantEmail(tenant);
      if (email && email !== 'No email provided' && email.length > 0) {
        return email.charAt(0).toUpperCase();
      }
      
      return 'T';
    }
    
    // Methods
    const getStatusBadgeClass = (status) => {
      const statusLower = status?.toLowerCase() || 'pending'
      const classes = {
        'pending': 'bg-yellow-100 text-yellow-800 border border-yellow-200',
        'active': 'bg-green-100 text-green-800 border border-green-200',
        'approved': 'bg-green-100 text-green-800 border border-green-200',
        'rejected': 'bg-red-100 text-red-800 border border-red-200',
        'declined': 'bg-red-100 text-red-800 border border-red-200'
      }
      return classes[statusLower] || 'bg-gray-100 text-gray-800 border border-gray-200'
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      try {
        const date = new Date(dateString)
        return date.toLocaleDateString('en-IN', {
          year: 'numeric',
          month: 'short',
          day: 'numeric'
        })
      } catch {
        return 'N/A'
      }
    }
    
    const goToPage = (page) => {
      if (page !== '...') {
        currentPage.value = page
      }
    }
    
    // Updated filter function to properly apply filters
    const applyFilters = () => {
      // Reset to first page when filter changes
      currentPage.value = 1
      
      // If "View All" is selected, load all tenants
      if (statusFilter.value === '') {
        loadAllTenants()
      }
    }
    
    // Load pending tenants (default view)
    const loadTenants = async () => {
      loading.value = true
      
      try {
        console.log('Loading tenant data...')
        
        try {
          // First try the correct API endpoint
          const response = await fetch('/api/v1/superadmin/tenants?status=pending', {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
            }
          })
          
          if (!response.ok) {
            throw new Error(`API returned ${response.status}: ${response.statusText}`)
          }
          
          const data = await response.json()
          console.log('API response:', data)
          
          if (data && data.data) {
            tenants.value = data.data
            allTenants.value = data.data
            console.log(`Loaded ${tenants.value.length} tenant applications from API`)
            toast.success(`Loaded ${tenants.value.length} tenant applications`)
          } else if (data) {
            tenants.value = data
            allTenants.value = data
            console.log(`Loaded ${tenants.value.length} tenant applications from API`)
            toast.success(`Loaded ${tenants.value.length} tenant applications`)
          } else {
            console.warn('API returned unexpected format:', data)
            toast.error('Could not load tenant data: unexpected API response format')
          }
        } catch (apiError) {
          console.error('Error calling API directly:', apiError)
          
          // Try using the service
          try {
            const serviceResponse = await superAdminService.getPendingTenants()
            console.log('Service response:', serviceResponse)
            
            if (serviceResponse && serviceResponse.success && serviceResponse.data) {
              // Check different possible structures from the service
              let tenantsData = []
              
              if (serviceResponse.data.pending_tenants) {
                tenantsData = serviceResponse.data.pending_tenants
                console.log('Found pending_tenants in service response')
              } else if (Array.isArray(serviceResponse.data)) {
                tenantsData = serviceResponse.data
                console.log('Found array in service response')
              } else {
                console.warn('Unexpected service response format:', serviceResponse)
                tenantsData = []
              }
              
              tenants.value = tenantsData
              allTenants.value = tenantsData
              
              if (tenantsData.length > 0) {
                toast.success(`Loaded ${tenantsData.length} tenant applications`)
              } else {
                toast.info('No pending tenant applications found')
              }
            } else {
              console.warn('Service call failed or returned unexpected format')
              toast.error('Could not load tenant data from API')
            }
          } catch (serviceError) {
            console.error('Error calling service:', serviceError)
            toast.error('Could not load tenant data')
          }
        }
      } catch (error) {
        console.error('Error in loadTenants:', error)
        toast.error('Error loading tenant data')
      } finally {
        loading.value = false
      }
    }
    
    // Load all tenants (for "View All" option)
    const loadAllTenants = async () => {
      loading.value = true
      
      try {
        const response = await superAdminService.getAllTenants()
        
        if (response && response.success && Array.isArray(response.data)) {
          allTenants.value = response.data
          toast.success(`Loaded ${response.data.length} total tenants`)
        } else {
          // Fallback to pending tenants if getAllTenants fails
          await loadTenants()
        }
      } catch (error) {
        console.error('Error loading all tenants:', error)
        // Fallback to pending tenants if getAllTenants fails
        await loadTenants()
        toast.warning('Showing pending tenants only')
      } finally {
        loading.value = false
      }
    }
    
    // Debug methods
    const debugTenantData = async () => {
      apiDebugInfo.value = 'Loading API debug information...'
      showApiDebugModal.value = true
      
      try {
        const debugInfo = []
        debugInfo.push('==== TENANT DATA DEBUG ====\n')
        
        // 1. Test direct API call to verify backend
        debugInfo.push('Testing direct API call to /api/v1/superadmin/tenants?status=pending...')
        try {
          const response = await fetch('/api/v1/superadmin/tenants?status=pending', {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
            }
          })
          
          if (!response.ok) {
            debugInfo.push(`API Error: ${response.status} ${response.statusText}`)
            const errorText = await response.text()
            debugInfo.push(`Error details: ${errorText}`)
          } else {
            const data = await response.json()
            debugInfo.push(`API Response Status: ${response.status}`)
            debugInfo.push(`API Response Data: ${JSON.stringify(data, null, 2)}`)
            
            // Try to load the data if not already loaded
            if (allTenants.value.length === 0 && data && data.data) {
              tenants.value = data.data
              allTenants.value = data.data
              debugInfo.push(`\nAUTOMATICALLY LOADED ${data.data.length} TENANTS FROM API RESPONSE`)
            } else if (allTenants.value.length === 0 && Array.isArray(data)) {
              tenants.value = data
              allTenants.value = data
              debugInfo.push(`\nAUTOMATICALLY LOADED ${data.length} TENANTS FROM API RESPONSE`)
            }
          }
        } catch (error) {
          debugInfo.push(`API Error: ${error.message}`)
        }
        
        // 2. Check service methods
        debugInfo.push('\n==== SERVICE METHODS DEBUG ====\n')
        debugInfo.push(`superAdminService.baseURL: ${superAdminService.baseURL}`)
        
        // 3. Display tenant data
        debugInfo.push('\n==== TENANT DATA ====\n')
        debugInfo.push(`Total Tenants: ${allTenants.value.length}`)
        if (allTenants.value.length > 0) {
          const firstTenant = allTenants.value[0]
          debugInfo.push(`First Tenant Data: ${JSON.stringify(firstTenant, null, 2)}`)
          debugInfo.push(`First Tenant ID: ${firstTenant.id || firstTenant.ID || 'Not found'}`)
          
          // List all properties that could be used as ID
          const potentialIds = ['ID', 'id', 'user_id', 'userId', '_id']
          for (const idField of potentialIds) {
            debugInfo.push(`${idField}: ${firstTenant[idField]}`)
          }
        }
        
        apiDebugInfo.value = debugInfo.join('\n')
      } catch (error) {
        apiDebugInfo.value = `Error running debug: ${error.message}`
      }
    }
    
    // Direct API test for approval
    const testApprovalApi = async (tenant) => {
      if (!tenant) {
        toast.error('No tenant selected for testing')
        return
      }
      
      apiDebugInfo.value = `Testing API call for tenant ${tenant.id || tenant.ID || 'unknown'}...`
      showApiDebugModal.value = true
      
      try {
        const debugInfo = []
        debugInfo.push(`Testing API approval for tenant ${JSON.stringify(tenant, null, 2)}\n`)
        
        // Extract tenant ID with fallbacks
        const tenantId = tenant.id || tenant.ID || tenant.user_id || tenant.userId
        
        if (!tenantId) {
          debugInfo.push('ERROR: No tenant ID found in tenant object!')
          apiDebugInfo.value = debugInfo.join('\n')
          return
        }
        
        debugInfo.push(`Using tenant ID: ${tenantId}`)
        
        // Test direct API call
        debugInfo.push('\nTesting direct PATCH call to approve tenant...')
        try {
          const token = localStorage.getItem('auth_token')
          debugInfo.push(`Auth token available: ${!!token}`)
          
          const response = await fetch(`/api/v1/superadmin/tenants/${tenantId}/approval`, {
            method: 'PATCH',
            headers: {
              'Content-Type': 'application/json',
              'Authorization': `Bearer ${token}`
            },
            body: JSON.stringify({
              status: 'APPROVED'
            })
          })
          
          debugInfo.push(`API Response Status: ${response.status} ${response.statusText}`)
          
          if (response.ok) {
            const responseData = await response.json()
            debugInfo.push(`API Response Data: ${JSON.stringify(responseData, null, 2)}`)
            debugInfo.push('\nSUCCESS: Direct API call worked! The backend endpoint is functional.')
            toast.success('Test API call successful!')
            
            // Refresh the tenant list to show the updated status
            await loadTenants()
          } else {
            const errorText = await response.text()
            debugInfo.push(`Error response: ${errorText}`)
          }
        } catch (error) {
          debugInfo.push(`API Error: ${error.message}`)
        }
        
        apiDebugInfo.value = debugInfo.join('\n')
      } catch (error) {
        apiDebugInfo.value = `Error running API test: ${error.message}`
      }
    }
    
    // Updated to use id instead of ID
    const handleApprove = async (tenant) => {
      if (!tenant) {
        toast.error('Cannot approve: Missing tenant');
        return;
      }
      
      // Use lowercase id first, fall back to uppercase ID
      const tenantId = tenant.id || tenant.ID;
      
      if (!tenantId) {
        toast.error('Cannot approve: Missing tenant ID');
        return;
      }
      
      isProcessing.value = true;
      console.log(`Attempting to approve tenant ${tenantId} with PATCH to /api/v1/superadmin/tenants/${tenantId}/approval`);
      
      try {
        // Make a direct API call to the approval endpoint
        const response = await fetch(`/api/v1/superadmin/tenants/${tenantId}/approval`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
          },
          body: JSON.stringify({
            status: "APPROVED"
          })
        });
        
        if (response.ok) {
          const data = await response.json();
          console.log('Approval successful:', data);
          toast.success(`Tenant ${getTenantName(tenant)} approved successfully`);
          
          // Refresh the tenant list
          loadTenants();
          emit('updated');
        } else {
          console.error('Approval failed with status:', response.status);
          const errorData = await response.text();
          console.error('Error response:', errorData);
          toast.error(`Failed to approve tenant: ${response.statusText}`);
        }
      } catch (error) {
        console.error('Error in approval process:', error);
        toast.error('Error approving tenant: ' + (error.message || 'Unknown error'));
      } finally {
        isProcessing.value = false;
      }
    };
    
    // Updated to use id instead of ID
    const handleRejectClick = (tenant) => {
      if (!tenant) {
        toast.error('Cannot reject: Missing tenant');
        return;
      }
      
      // Use lowercase id first, fall back to uppercase ID
      const tenantId = tenant.id || tenant.ID;
      
      if (!tenantId) {
        toast.error('Cannot reject: Missing tenant ID');
        return;
      }
      
      selectedTenant.value = tenant;
      rejectReason.value = '';
      showRejectModal.value = true;
    };
    
    // Updated to use id instead of ID
    const confirmReject = async () => {
      const tenant = selectedTenant.value;
      
      if (!tenant) {
        toast.error('Cannot reject: Missing tenant');
        showRejectModal.value = false;
        return;
      }
      
      // Use lowercase id first, fall back to uppercase ID
      const tenantId = tenant.id || tenant.ID;
      
      if (!tenantId) {
        toast.error('Cannot reject: Missing tenant ID');
        showRejectModal.value = false;
        return;
      }
      
      if (!rejectReason.value.trim()) {
        toast.warning('Please provide a rejection reason');
        return;
      }
      
      isProcessing.value = true;
      
      try {
        // Make a direct API call to the rejection endpoint
        const response = await fetch(`/api/v1/superadmin/tenants/${tenantId}/approval`, {
          method: 'PATCH',
          headers: {
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
          },
          body: JSON.stringify({
            status: "REJECTED",
            reason: rejectReason.value
          })
        });
        
        if (response.ok) {
          const data = await response.json();
          console.log('Rejection successful:', data);
          toast.success(`Tenant ${getTenantName(tenant)} rejected successfully`);
          showRejectModal.value = false;
          
          // Refresh the tenant list
          loadTenants();
          emit('updated');
        } else {
          console.error('Rejection failed with status:', response.status);
          const errorData = await response.text();
          console.error('Error response:', errorData);
          toast.error(`Failed to reject tenant: ${response.statusText}`);
        }
      } catch (error) {
        console.error('Error in rejection process:', error);
        toast.error('Error rejecting tenant: ' + (error.message || 'Unknown error'));
      } finally {
        isProcessing.value = false;
      }
    };
    
    // Load data on mount
    onMounted(() => {
      loadTenants()
    })
    
    return {
      loading,
      tenants,
      allTenants,
      pendingCount,
      approvedCount,
      rejectedCount,
      statusFilter,
      currentPage,
      itemsPerPage,
      filteredTenants,
      paginatedTenants,
      totalPages,
      paginationStart,
      paginationEnd,
      displayedPageNumbers,
      getStatusBadgeClass,
      formatDate,
      loadTenants,
      loadAllTenants,
      handleApprove,
      handleRejectClick,
      getTenantName,
      getTenantEmail,
      getTenantInitial,
      goToPage,
      isStatusPending,
      isStatusApproved,
      isStatusRejected,
      applyFilters,
      isProcessing,
      showRejectModal,
      selectedTenant,
      rejectReason,
      confirmReject,
      // Debug & Mock data
      debugMode,
      debugTenantData,
      testApprovalApi,
      showApiDebugModal,
      apiDebugInfo
    }
  }
}
</script>