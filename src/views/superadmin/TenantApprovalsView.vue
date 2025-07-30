<template>
  <div class="space-y-6">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 font-roboto">
          Tenant Approvals
        </h2>
        <p class="text-sm text-gray-600 mt-1">
          Review and manage temple admin registration requests
        </p>
      </div>
      
      <!-- Stats Cards -->
      <div class="flex gap-3">
        <div class="bg-yellow-50 border border-yellow-200 rounded-xl px-4 py-3 text-center min-w-[100px]">
          <div class="text-2xl font-bold text-yellow-800">{{ superAdminStore.pendingCount }}</div>
          <div class="text-xs text-yellow-600 font-medium">Pending</div>
        </div>
        <div class="bg-green-50 border border-green-200 rounded-xl px-4 py-3 text-center min-w-[100px]">
          <div class="text-2xl font-bold text-green-800">{{ superAdminStore.activeCount }}</div>
          <div class="text-xs text-green-600 font-medium">Approved</div>
        </div>
        <div class="bg-red-50 border border-red-200 rounded-xl px-4 py-3 text-center min-w-[100px]">
          <div class="text-2xl font-bold text-red-800">{{ superAdminStore.rejectedCount }}</div>
          <div class="text-xs text-red-600 font-medium">Rejected</div>
        </div>
      </div>
    </div>

    <!-- Debug Button -->
    <div class="flex justify-end">
      <button 
        @click="debugTenantData" 
        class="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg text-sm flex items-center gap-2"
      >
        <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Debug API
      </button>
    </div>

    <!-- Filters Section -->
    <div class="bg-white rounded-xl shadow-md border border-gray-200 p-6">
      <div class="flex flex-col lg:flex-row gap-4 items-start lg:items-center">
        <div class="flex-1 grid grid-cols-1 md:grid-cols-3 gap-4">
          <!-- Search Input -->
          <div class="relative">
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search by name, email, or temple..."
              class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm"
            />
            <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>

          <!-- Status Filter -->
          <select
            v-model="statusFilter"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm bg-white"
          >
            <option value="">All Status</option>
            <option value="PENDING">Pending</option>
            <option value="APPROVED">Approved</option>
            <option value="REJECTED">Rejected</option>
          </select>

          <!-- Date Filter -->
          <input
            v-model="dateFilter"
            type="date"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 text-sm"
          />
        </div>

        <!-- Clear Filters Button -->
        <button
          @click="clearFilters"
          class="px-4 py-3 text-sm font-medium text-gray-600 bg-gray-100 hover:bg-gray-200 rounded-xl transition-all duration-200 whitespace-nowrap"
        >
          Clear Filters
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Tenant Applications List -->
    <div v-else-if="filteredTenants.length > 0" class="space-y-4">
      <div
        v-for="tenant in paginatedTenants"
        :key="tenant.id"
        class="bg-white rounded-xl shadow-md border border-gray-200 hover:shadow-lg transition-all duration-200"
      >
        <div class="p-6">
          <!-- Header Row -->
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-4">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center flex-shrink-0">
                <span class="text-indigo-600 font-semibold text-lg">
                  {{ tenant.fullName ? tenant.fullName.charAt(0).toUpperCase() : 'T' }}
                </span>
              </div>
              
              <!-- Basic Info -->
              <div class="flex-1 min-w-0">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">
                  {{ tenant.fullName || tenant.name || 'Unknown Tenant' }}
                </h3>
                <p class="text-sm text-gray-600">{{ tenant.email || 'No email provided' }}</p>
                <p class="text-sm text-gray-500">{{ tenant.phone || 'No phone provided' }}</p>
              </div>
            </div>

            <!-- Status & Date -->
            <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
              <span
                :class="getStatusBadgeClass(tenant.status)"
                class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium"
              >
                {{ tenant.status }}
              </span>
              <span class="text-xs text-gray-500">
                Applied: {{ formatDate(tenant.createdAt || tenant.submitted) }}
              </span>
            </div>
          </div>

          <!-- Temple Information -->
          <div v-if="tenant.temple" class="bg-gray-50 rounded-lg p-4 mb-4">
            <h4 class="text-sm font-semibold text-gray-900 mb-2">Temple Information</h4>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-3 text-sm">
              <div>
                <span class="text-gray-600">Name:</span>
                <span class="ml-2 text-gray-900">{{ tenant.temple.name || 'Pending' }}</span>
              </div>
              <div>
                <span class="text-gray-600">Type:</span>
                <span class="ml-2 text-gray-900">{{ tenant.temple.type || 'Hindu Temple' }}</span>
              </div>
              <div>
                <span class="text-gray-600">City:</span>
                <span class="ml-2 text-gray-900">{{ tenant.temple.city || tenant.location || 'Pending' }}</span>
              </div>
              <div>
                <span class="text-gray-600">State:</span>
                <span class="ml-2 text-gray-900">{{ tenant.temple.state || 'Pending' }}</span>
              </div>
            </div>
            <div class="mt-2">
              <span class="text-gray-600">Address:</span>
              <span class="ml-2 text-gray-900">{{ tenant.temple.address || tenant.address || 'Pending' }}</span>
            </div>
          </div>

          <!-- Documents Section -->
          <div v-if="tenant.documents && tenant.documents.length > 0" class="mb-4">
            <h4 class="text-sm font-semibold text-gray-900 mb-2">Uploaded Documents</h4>
            <div class="flex flex-wrap gap-2">
              <button
                v-for="doc in tenant.documents"
                :key="doc.id"
                @click="viewDocument(doc)"
                class="inline-flex items-center px-3 py-1 bg-indigo-50 text-indigo-700 rounded-lg text-xs hover:bg-indigo-100 transition-colors duration-200"
              >
                <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                {{ doc.name }}
              </button>
            </div>
          </div>

          <!-- Rejection Notes (if rejected) -->
          <div v-if="tenant.status === 'REJECTED' && tenant.rejectionNotes" class="mb-4">
            <div class="bg-red-50 border border-red-200 rounded-lg p-3">
              <h4 class="text-sm font-semibold text-red-800 mb-1">Rejection Notes</h4>
              <p class="text-sm text-red-700">{{ tenant.rejectionNotes }}</p>
            </div>
          </div>

          <!-- Action Buttons -->
          <div v-if="tenant.status === 'PENDING'" class="flex flex-col sm:flex-row gap-3">
            <button
              @click="openApprovalModal(tenant)"
              :disabled="isProcessing"
              class="flex-1 sm:flex-none px-6 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
              </svg>
              Approve
            </button>
            
            <button
              @click="openRejectionModal(tenant)"
              :disabled="isProcessing"
              class="flex-1 sm:flex-none px-6 py-2 bg-red-600 hover:bg-red-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
              Reject
            </button>
            
            <button
              @click="viewTenantDetails(tenant)"
              class="flex-1 sm:flex-none px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
              View Details
            </button>
          </div>

          <!-- Approved/Rejected Actions -->
          <div v-else class="flex flex-col sm:flex-row gap-3">
            <button
              @click="viewTenantDetails(tenant)"
              class="flex-1 sm:flex-none px-6 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
              View Details
            </button>
            
            <div class="text-sm text-gray-500 flex items-center">
              {{ tenant.status === 'APPROVED' ? 'Approved' : 'Rejected' }} on {{ formatDate(tenant.updatedAt) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
      </svg>
      <h3 class="mt-4 text-lg font-medium text-gray-900">No tenant applications found</h3>
      <p class="mt-2 text-sm text-gray-500">
        {{ searchQuery || statusFilter || dateFilter ? 'Try adjusting your filters' : 'No applications have been submitted yet' }}
      </p>
      <button 
        @click="refreshData" 
        class="mt-4 px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm"
      >
        Refresh Data
      </button>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="flex items-center justify-between bg-white rounded-xl shadow-md border border-gray-200 px-6 py-4">
      <div class="text-sm text-gray-700">
        Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredTenants.length) }} of {{ filteredTenants.length }} results
      </div>
      
      <div class="flex items-center gap-2">
        <button
          @click="currentPage--"
          :disabled="currentPage === 1"
          class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        >
          Previous
        </button>
        
        <span
          v-for="page in visiblePages"
          :key="page"
          @click="currentPage = page"
          :class="[
            'px-3 py-2 text-sm font-medium cursor-pointer rounded-lg transition-all duration-200',
            currentPage === page
              ? 'bg-indigo-600 text-white'
              : 'text-gray-500 bg-white border border-gray-300 hover:bg-gray-50'
          ]"
        >
          {{ page }}
        </span>
        
        <button
          @click="currentPage++"
          :disabled="currentPage === totalPages"
          class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        >
          Next
        </button>
      </div>
    </div>

    <!-- Approval Modal -->
    <div
      v-if="approvalModal.show"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeApprovalModal"></div>
        
        <div class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-6 pt-6 pb-4">
            <div class="flex items-center">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-green-100">
                <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
              </div>
              <div class="ml-4 text-left">
                <h3 class="text-lg leading-6 font-semibold text-gray-900 font-roboto" id="modal-title">
                  Approve Tenant Application
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500">
                    Are you sure you want to approve {{ approvalModal.tenant?.fullName || approvalModal.tenant?.name || 'this tenant' }}'s application{{ approvalModal.tenant?.temple ? ' for ' + approvalModal.tenant?.temple?.name : '' }}?
                  </p>
                </div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 px-6 py-4 flex flex-col sm:flex-row-reverse gap-3">
            <button
              @click="approveTenant"
              :disabled="isProcessing"
              class="w-full sm:w-auto px-6 py-2 bg-green-600 hover:bg-green-700 disabled:opacity-50 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
            >
              <div v-if="isProcessing" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              {{ isProcessing ? 'Processing...' : 'Approve' }}
            </button>
            <button
              @click="closeApprovalModal"
              :disabled="isProcessing"
              class="w-full sm:w-auto px-6 py-2 bg-white hover:bg-gray-50 text-gray-700 text-sm font-medium border border-gray-300 rounded-xl transition-all duration-200"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Rejection Modal -->
    <div
      v-if="rejectionModal.show"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeRejectionModal"></div>
        
        <div class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-6 pt-6 pb-4">
            <div class="flex items-start">
              <div class="mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
                <svg class="h-6 w-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                </svg>
              </div>
              <div class="ml-4 w-full text-left">
                <h3 class="text-lg leading-6 font-semibold text-gray-900 font-roboto" id="modal-title">
                  Reject Tenant Application
                </h3>
                <div class="mt-2">
                  <p class="text-sm text-gray-500 mb-4">
                    Please provide a reason for rejecting {{ rejectionModal.tenant?.fullName || rejectionModal.tenant?.name || 'this tenant' }}'s application:
                  </p>
                  <textarea
                    v-model="rejectionModal.notes"
                    rows="4"
                    placeholder="Enter rejection reason..."
                    class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-red-500 focus:border-red-500 text-sm resize-none"
                  ></textarea>
                </div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 px-6 py-4 flex flex-col sm:flex-row-reverse gap-3">
            <button
              @click="rejectTenant"
              :disabled="isProcessing || !rejectionModal.notes.trim()"
              class="w-full sm:w-auto px-6 py-2 bg-red-600 hover:bg-red-700 disabled:opacity-50 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
            >
              <div v-if="isProcessing" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
              {{ isProcessing ? 'Processing...' : 'Reject' }}
            </button>
            <button
              @click="closeRejectionModal"
              :disabled="isProcessing"
              class="w-full sm:w-auto px-6 py-2 bg-white hover:bg-gray-50 text-gray-700 text-sm font-medium border border-gray-300 rounded-xl transition-all duration-200"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from '@/composables/useToast'
import { useSuperAdminStore } from '@/stores/superadmin'
import superAdminService from '@/services/superadmin.service'

export default {
  name: 'TenantApprovalView',
  setup() {
    // Toast notifications
    const toast = useToast()
    
    // Use the SuperAdmin store
    const superAdminStore = useSuperAdminStore()
    
    // Reactive data
    const isProcessing = ref(false)
    const searchQuery = ref('')
    const statusFilter = ref('PENDING') // Default to pending
    const dateFilter = ref('')
    const currentPage = ref(1)
    const itemsPerPage = ref(10)

    // Modal states
    const approvalModal = ref({
      show: false,
      tenant: null
    })

    const rejectionModal = ref({
      show: false,
      tenant: null,
      notes: ''
    })

    // Computed properties
    const isLoading = computed(() => superAdminStore.loadingTenants)
    
    const filteredTenants = computed(() => {
      let filtered = superAdminStore.tenants

      // Search filter
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(tenant => 
          (tenant.fullName || tenant.name || '').toLowerCase().includes(query) ||
          (tenant.email || '').toLowerCase().includes(query) ||
          (tenant.temple?.name || '').toLowerCase().includes(query)
        )
      }

      // Status filter
      if (statusFilter.value) {
        filtered = filtered.filter(tenant => tenant.status === statusFilter.value)
      }

      // Date filter
      if (dateFilter.value) {
        filtered = filtered.filter(tenant => {
          if (!tenant.createdAt && !tenant.submitted) return false;
          const date = tenant.createdAt || tenant.submitted;
          const tenantDate = new Date(date).toISOString().split('T')[0]
          return tenantDate === dateFilter.value
        })
      }

      return filtered
    })

    const paginatedTenants = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value
      const end = start + itemsPerPage.value
      return filteredTenants.value.slice(start, end)
    })

    const totalPages = computed(() => {
      return Math.ceil(filteredTenants.value.length / itemsPerPage.value)
    })

    const visiblePages = computed(() => {
      const pages = []
      const total = totalPages.value
      const current = currentPage.value
      
      if (total <= 7) {
        for (let i = 1; i <= total; i++) {
          pages.push(i)
        }
      } else {
        if (current <= 4) {
          for (let i = 1; i <= 5; i++) {
            pages.push(i)
          }
          pages.push('...', total)
        } else if (current >= total - 3) {
          pages.push(1, '...')
          for (let i = total - 4; i <= total; i++) {
            pages.push(i)
          }
        } else {
          pages.push(1, '...')
          for (let i = current - 1; i <= current + 1; i++) {
            pages.push(i)
          }
          pages.push('...', total)
        }
      }
      
      return pages.filter(page => page !== '...' || pages.indexOf(page) === pages.lastIndexOf(page))
    })

    // Methods
    const getStatusBadgeClass = (status) => {
      const classes = {
        'PENDING': 'bg-yellow-100 text-yellow-800 border border-yellow-200',
        'APPROVED': 'bg-green-100 text-green-800 border border-green-200',
        'REJECTED': 'bg-red-100 text-red-800 border border-red-200',
        'pending': 'bg-yellow-100 text-yellow-800 border border-yellow-200',
        'approved': 'bg-green-100 text-green-800 border border-green-200',
        'rejected': 'bg-red-100 text-red-800 border border-red-200'
      }
      return classes[status] || 'bg-gray-100 text-gray-800 border border-gray-200'
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'N/A'
      const date = new Date(dateString)
      return date.toLocaleDateString('en-IN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }

    const clearFilters = () => {
      searchQuery.value = ''
      statusFilter.value = ''
      dateFilter.value = ''
      currentPage.value = 1
    }

    const refreshData = () => {
      superAdminStore.fetchTenants()
      superAdminStore.fetchStats()
      toast.info('Refreshing data...')
    }

    const openApprovalModal = (tenant) => {
      approvalModal.value = {
        show: true,
        tenant
      }
    }

    const closeApprovalModal = () => {
      approvalModal.value = {
        show: false,
        tenant: null
      }
    }

    const openRejectionModal = (tenant) => {
      rejectionModal.value = {
        show: true,
        tenant,
        notes: ''
      }
    }

    const closeRejectionModal = () => {
      rejectionModal.value = {
        show: false,
        tenant: null,
        notes: ''
      }
    }

    const approveTenant = async () => {
      if (!approvalModal.value.tenant) return
      
      isProcessing.value = true
      
      try {
        // Call store action to approve tenant
        const result = await superAdminStore.approveTenant(approvalModal.value.tenant.id, {
          notes: 'Approved by Super Admin'
        })
        
        if (result.success) {
          toast.success(`Tenant ${approvalModal.value.tenant.fullName || approvalModal.value.tenant.name} approved successfully`)
        } else {
          toast.error(result.error || 'Failed to approve tenant')
        }
        
        closeApprovalModal()
      } catch (error) {
        console.error('Error approving tenant:', error)
        toast.error('Failed to approve tenant. Please try again.')
      } finally {
        isProcessing.value = false
      }
    }

    const rejectTenant = async () => {
      if (!rejectionModal.value.tenant || !rejectionModal.value.notes.trim()) return
      
      isProcessing.value = true
      
      try {
        // Call store action to reject tenant
        const result = await superAdminStore.rejectTenant(rejectionModal.value.tenant.id, {
          notes: rejectionModal.value.notes
        })
        
        if (result.success) {
          toast.success(`Tenant ${rejectionModal.value.tenant.fullName || rejectionModal.value.tenant.name} rejected successfully`)
        } else {
          toast.error(result.error || 'Failed to reject tenant')
        }
        
        closeRejectionModal()
      } catch (error) {
        console.error('Error rejecting tenant:', error)
        toast.error('Failed to reject tenant. Please try again.')
      } finally {
        isProcessing.value = false
      }
    }

    const viewDocument = (document) => {
      // Open document in new tab if URL is available
      if (document.url) {
        window.open(document.url, '_blank')
      } else {
        toast.info('Document preview not available')
      }
    }

    const viewTenantDetails = (tenant) => {
      // This could open a modal or navigate to detailed view
      toast.info(`Viewing details for ${tenant.fullName || tenant.name}`)
      
      // In a real implementation, you might navigate to a details page
      // router.push(`/superadmin/tenants/${tenant.id}`)
    }

    const debugTenantData = async () => {
      try {
        console.log('==== DEBUGGING API REQUESTS ====')
        
        // 1. Log the API endpoint being called
        console.log('API endpoint:', `${window.location.origin}/api/v1/superadmin/tenants/pending`)
        
        // 2. Make a direct fetch request to bypass the service
        const response = await fetch('/api/v1/superadmin/tenants/pending')
        const responseText = await response.text()
        
        console.log('Raw API Response:', {
          status: response.status,
          statusText: response.statusText,
          headers: Object.fromEntries([...response.headers.entries()]),
          responseLength: responseText.length
        })
        
        // 3. Try to parse as JSON
        try {
          const jsonData = JSON.parse(responseText)
          console.log('API Response Data:', jsonData)
          
          // Check if the response contains tenant data
          if (Array.isArray(jsonData)) {
            console.log('Response contains array with', jsonData.length, 'items')
            if (jsonData.length > 0) {
              console.log('First item:', jsonData[0])
            } else {
              console.log('Array is empty - no tenant data')
            }
          } else if (jsonData && Array.isArray(jsonData.tenants)) {
            console.log('Response contains tenant array with', jsonData.tenants.length, 'items')
            if (jsonData.tenants.length > 0) {
              console.log('First item:', jsonData.tenants[0])
            } else {
              console.log('Tenant array is empty - no tenant data')
            }
          } else {
            console.log('Response data structure does not contain expected tenant array')
            console.log('Response data keys:', Object.keys(jsonData))
          }
        } catch (e) {
          console.log('Response is not valid JSON:', e)
        }
        
        // 4. Check localStorage for authentication token
        const token = localStorage.getItem('auth_token')
        console.log('Auth token exists:', !!token)
        if (token) {
          console.log('Token prefix:', token.substring(0, 15) + '...')
        }
        
        // 5. Test the service directly
        console.log('Testing superAdminService.getPendingTenants() directly...')
        const serviceResponse = await superAdminService.getPendingTenants()
        console.log('Service response:', serviceResponse)
        
        // 6. Refresh data from the store
        console.log('Refreshing data from the store...')
        superAdminStore.fetchTenants()
        
        toast.info('API debug completed - check browser console')
      } catch (error) {
        console.error('Debug API call failed:', error)
        toast.error('API debug failed: ' + error.message)
      }
    }

    // Watchers
    watch([searchQuery, statusFilter, dateFilter], () => {
      currentPage.value = 1
    })

    // Lifecycle
    onMounted(() => {
      console.log('TenantApproval component mounted')
      // Initialize the store data
      superAdminStore.initialize()
    })

    return {
      // Store
      superAdminStore,
      
      // Reactive data
      isLoading,
      isProcessing,
      searchQuery,
      statusFilter,
      dateFilter,
      currentPage,
      itemsPerPage,
      approvalModal,
      rejectionModal,
      
      // Computed properties
      filteredTenants,
      paginatedTenants,
      totalPages,
      visiblePages,
      
      // Methods
      getStatusBadgeClass,
      formatDate,
      clearFilters,
      refreshData,
      openApprovalModal,
      closeApprovalModal,
      openRejectionModal,
      closeRejectionModal,
      approveTenant,
      rejectTenant,
      viewDocument,
      viewTenantDetails,
      debugTenantData
    }
  }
}
</script>

<style scoped>
/* Custom font classes */
.font-roboto {
  font-family: 'Roboto', sans-serif;
}

/* Custom scrollbar for better UX */
::-webkit-scrollbar {
  width: 6px;
}

::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 6px;
}

::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 6px;
}

::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Smooth transitions for all interactive elements */
* {
  transition: all 0.2s ease-in-out;
}

/* Enhanced focus states for accessibility */
button:focus,
input:focus,
select:focus,
textarea:focus {
  outline: 2px solid #6366f1;
  outline-offset: 2px;
}

/* Loading animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Card hover effects */
.bg-white:hover {
  transform: translateY(-2px);
}

/* Button loading state */
button:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

/* Modal backdrop blur effect */
.fixed.inset-0.bg-gray-500 {
  backdrop-filter: blur(4px);
}

/* Status badge pulse animation for pending items */
.bg-yellow-100 {
  animation: pulse 2s infinite;
}

@keyframes pulse {
  0%, 100% {
    opacity: 1;
  }
  50% {
    opacity: 0.8;
  }
}

/* Responsive grid adjustments */
@media (max-width: 640px) {
  .grid-cols-1 {
    gap: 1rem;
  }
}

/* Enhanced spacing for mobile */
@media (max-width: 768px) {
  .p-6 {
    padding: 1rem;
  }
  
  .px-6 {
    padding-left: 1rem;
    padding-right: 1rem;
  }
}

/* Print styles */
@media print {
  .no-print,
  button,
  .fixed {
    display: none !important;
  }
  
  .bg-white {
    box-shadow: none !important;
  }
}
</style>