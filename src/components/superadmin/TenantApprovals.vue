<template>
  <div class="space-y-6">
    <!-- Heading Section -->
    <div class="border-b border-gray-200 pb-4">
      <h1 class="text-2xl font-bold text-gray-900 font-roboto">Temple Approvals</h1>
      <p class="mt-1 text-sm text-gray-500">Review and manage temple registration requests</p>
    </div>

    <!-- Debug Info (remove in production) -->
    <div class="bg-gray-100 p-4 rounded-lg mb-4 text-xs font-mono overflow-auto max-h-40" v-if="debugMode">
      <div class="mb-2 font-bold">Debug Information:</div>
      <div v-if="selectedEntity">Selected Temple: ID={{ selectedEntity.id || selectedEntity.ID }}, Name={{ getEntityName(selectedEntity) }}</div>
      <div>Temple Count: {{ Array.isArray(allEntities) ? allEntities.length : 0 }}</div>
      <div class="mt-2 flex gap-2">
        <button @click="debugEntityData" class="px-3 py-1 bg-gray-200 rounded text-xs">Run API Debug</button>
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
        @click="loadEntities" 
        class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg text-sm"
      >
        Refresh Data
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex justify-center items-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
    </div>

    <!-- Temple Applications List -->
    <div v-if="!loading && Array.isArray(filteredEntities) && filteredEntities.length > 0" class="space-y-4">
      <div
        v-for="(entity, idx) in paginatedEntities"
        :key="entity.id || entity.ID || idx"
        class="bg-white rounded-xl shadow-md border border-gray-200 hover:shadow-lg transition-all duration-200"
      >
        <div class="p-6">
          <!-- Entity ID Debug (remove in production) -->
          <div v-if="debugMode" class="bg-gray-100 p-2 mb-3 rounded text-xs">
            <span class="font-bold">ID:</span> {{ entity.id || entity.ID || 'Not available' }} |
            <span class="font-bold">Status:</span> {{ entity.status || entity.Status || 'Not available' }}
          </div>
          
          <!-- Header Row -->
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4 mb-4">
            <div class="flex items-start gap-4">
              <!-- Avatar -->
              <div class="h-12 w-12 rounded-full bg-indigo-100 flex items-center justify-center flex-shrink-0">
                <span class="text-indigo-600 font-semibold text-lg">
                  {{ getEntityInitial(entity) }}
                </span>
              </div>
              
              <!-- Basic Info - Temple name and main deity -->
              <div class="flex-1 min-w-0">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">
                  {{ getEntityName(entity) }}
                </h3>
                <p class="text-sm text-gray-600">{{ getMainDeity(entity) }}</p>
              </div>
            </div>

            <!-- Status & Date -->
            <div class="flex flex-col sm:flex-row items-start sm:items-center gap-3">
              <span
                :class="getStatusBadgeClass(entity.status || entity.Status)"
                class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium"
              >
                {{ entity.status || entity.Status || 'pending' }}
              </span>
              <span class="text-xs text-gray-500">
                Created: {{ formatDate(entity.created_at || entity.CreatedAt) }}
              </span>
            </div>
          </div>

          <!-- Action Buttons -->
          <div class="flex flex-col sm:flex-row gap-3">
            <!-- View Details Button (Always visible) -->
            <button
              @click="handleViewDetails(entity)"
              class="flex-1 sm:flex-none px-6 py-2 bg-indigo-100 hover:bg-indigo-200 text-indigo-700 text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
            >
              <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
              </svg>
              View Details
            </button>
            
            <!-- Approval/Rejection Buttons (Only for pending temples) -->
            <div v-if="isStatusPending(entity)" class="flex flex-col sm:flex-row gap-3 flex-1">
              <button
                @click="handleApprove(entity)"
                class="flex-1 sm:flex-none px-6 py-2 bg-green-600 hover:bg-green-700 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2"
                :disabled="isProcessing"
              >
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                </svg>
                {{ isProcessing ? 'Processing...' : 'Approve' }}
              </button>
              
              <button
                @click="handleRejectClick(entity)"
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
                @click="testApprovalApi(entity)"
                class="flex-1 sm:flex-none px-6 py-2 bg-gray-600 text-white text-sm font-medium rounded-xl flex items-center justify-center gap-2"
              >
                Test API Call
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination Controls -->
      <div class="flex flex-col sm:flex-row items-center justify-between bg-white p-4 rounded-xl shadow-sm border border-gray-200">
        <div class="text-sm text-gray-700 mb-3 sm:mb-0">
          Showing <span class="font-medium">{{ paginationStart }}</span> to 
          <span class="font-medium">{{ paginationEnd }}</span> of 
          <span class="font-medium">{{ Array.isArray(filteredEntities) ? filteredEntities.length : 0 }}</span> results
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
      <h3 class="mt-4 text-lg font-medium text-gray-900">No temple applications found</h3>
      <p class="mt-2 text-sm text-gray-500">
        {{ statusFilter ? `No temples with "${statusFilter}" status found` : 'Try refreshing or checking back later' }}
      </p>
    </div>

    <!-- Rejection Modal -->
    <div v-if="showRejectModal" class="fixed inset-0 z-50 overflow-y-auto bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-xl p-6 w-full max-w-md mx-4">
        <h3 class="text-lg font-semibold mb-4">Reject Temple</h3>
        
        <!-- Debug info -->
        <div v-if="debugMode" class="bg-gray-100 p-2 mb-3 rounded text-xs">
          <div><span class="font-bold">Selected Temple ID:</span> {{ selectedEntity?.id || selectedEntity?.ID }}</div>
          <div><span class="font-bold">Selected Temple Name:</span> {{ getEntityName(selectedEntity) }}</div>
        </div>
        
        <p class="mb-4">Please provide a reason for rejecting <span class="font-medium">{{ selectedEntity ? getEntityName(selectedEntity) : '' }}</span>:</p>
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

    <!-- Temple Details Modal -->
    <div v-if="showDetailsModal" class="fixed inset-0 z-50 overflow-y-auto bg-black bg-opacity-50 flex items-center justify-center">
      <div class="bg-white rounded-xl p-6 w-full max-w-2xl mx-4 max-h-[90vh] overflow-auto">
        <div class="flex justify-between items-center border-b border-gray-200 pb-4 mb-5">
          <h3 class="text-xl font-bold text-gray-900">Temple Details</h3>
          <button 
            @click="showDetailsModal = false"
            class="text-gray-400 hover:text-gray-600"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div v-if="selectedEntity" class="space-y-5">
          <!-- Basic Info -->
          <div>
            <div class="flex items-center gap-4 mb-4">
              <div class="h-16 w-16 rounded-full bg-indigo-100 flex items-center justify-center flex-shrink-0">
                <span class="text-indigo-600 font-semibold text-2xl">
                  {{ getEntityInitial(selectedEntity) }}
                </span>
              </div>
              <div>
                <h4 class="text-xl font-semibold text-gray-900">
                  {{ getEntityName(selectedEntity) }}
                </h4>
                <p class="text-gray-600">{{ getMainDeity(selectedEntity) }}</p>
              </div>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <!-- First Column - Temple Info -->
              <div class="space-y-4">
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Status</div>
                  <div class="text-base text-gray-900 mt-1">
                    <span :class="getStatusBadgeClass(selectedEntity.status || selectedEntity.Status)" class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium">
                      {{ selectedEntity.status || selectedEntity.Status || 'pending' }}
                    </span>
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Temple Type</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.temple_type || selectedEntity.TempleType || 'Not specified' }}
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Established Year</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.established_year || selectedEntity.EstablishedYear || 'Not specified' }}
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Created By</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.created_by || selectedEntity.CreatedBy || 'Unknown' }}
                  </div>
                </div>
              </div>
              
              <!-- Second Column - Contact Info -->
              <div class="space-y-4">
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Phone</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.phone || selectedEntity.Phone || 'Not provided' }}
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Email</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.email || selectedEntity.Email || 'Not provided' }}
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Address</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ getFullAddress(selectedEntity) }}
                  </div>
                </div>
                
                <div class="border-b border-gray-100 pb-3">
                  <div class="text-sm font-medium text-gray-500">Entity ID</div>
                  <div class="text-base text-gray-900 mt-1">
                    {{ selectedEntity.id || selectedEntity.ID || 'Not available' }}
                  </div>
                </div>
              </div>
            </div>
            
            <!-- Description Section -->
            <div class="mt-5 border-t border-gray-100 pt-4">
              <div class="text-sm font-medium text-gray-500 mb-2">Temple Description</div>
              <div class="text-base text-gray-900 bg-gray-50 p-4 rounded-lg">
                {{ selectedEntity.description || selectedEntity.Description || 'No description provided.' }}
              </div>
            </div>

            <!-- Debug Section -->
            <div v-if="debugMode" class="mt-4 p-3 bg-yellow-50 border border-yellow-200 rounded">
              <button 
                @click="debugEntityDetails(selectedEntity)"
                class="px-3 py-1 bg-yellow-200 hover:bg-yellow-300 text-yellow-800 rounded text-sm"
              >
                Debug Temple Details
              </button>
            </div>
          </div>
        </div>
        
        <!-- Actions -->
        <div class="flex justify-end gap-3 mt-6 pt-4 border-t border-gray-200">
          <button 
            @click="showDetailsModal = false"
            class="px-4 py-2 bg-gray-200 hover:bg-gray-300 text-gray-800 rounded-lg"
          >
            Close
          </button>
          
          <div v-if="isStatusPending(selectedEntity)">
            <button
              @click="handleApproveFromDetails"
              class="px-4 py-2 bg-green-600 hover:bg-green-700 text-white rounded-lg disabled:opacity-50"
              :disabled="isProcessing"
            >
              {{ isProcessing ? 'Processing...' : 'Approve' }}
            </button>
          </div>
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
  name: 'TempleApprovals',
  emits: ['updated'],
  setup(props, { emit }) {
    // Data
    const loading = ref(true)
    const entities = ref([])
    const allEntities = ref([]) // Store all entities for "View All" functionality
    const toast = useToast()
    const isProcessing = ref(false)
    
    // Debug mode
    const debugMode = ref(false)
    const showApiDebugModal = ref(false)
    const apiDebugInfo = ref('')
    
    // Simple rejection modal
    const showRejectModal = ref(false)
    const selectedEntity = ref(null)
    const rejectReason = ref('')
    
    // Details modal
    const showDetailsModal = ref(false)
    
    // Filters
    const statusFilter = ref('')
    
    // Pagination
    const currentPage = ref(1)
    const itemsPerPage = ref(5)
    
    // Computed properties
    const pendingCount = computed(() => 
      Array.isArray(allEntities.value) ? allEntities.value.filter(e => isStatusPending(e)).length : 0
    )
    
    const approvedCount = computed(() => 
      Array.isArray(allEntities.value) ? allEntities.value.filter(e => isStatusApproved(e)).length : 0
    )
    
    const rejectedCount = computed(() => 
      Array.isArray(allEntities.value) ? allEntities.value.filter(e => isStatusRejected(e)).length : 0
    )
    
    // Status check helpers
    const isStatusPending = (entity) => {
      if (!entity) return false;
      const status = (entity.status || entity.Status || '').toLowerCase();
      return status === 'pending' || status === 'pending approval' || status === '';
    }
    
    const isStatusApproved = (entity) => {
      if (!entity) return false;
      const status = (entity.status || entity.Status || '').toLowerCase();
      return status === 'active' || status === 'approved';
    }
    
    const isStatusRejected = (entity) => {
      if (!entity) return false;
      const status = (entity.status || entity.Status || '').toLowerCase();
      return status === 'rejected' || status === 'declined';
    }
    
    // Filtered entities based on status
    const filteredEntities = computed(() => {
      // Ensure we always work with arrays
      const allEntitiesArray = Array.isArray(allEntities.value) ? allEntities.value : [];
      
      if (!statusFilter.value) {
        return allEntitiesArray; // Show all entities when "View All" is selected
      }
      
      switch (statusFilter.value.toLowerCase()) {
        case 'pending':
          return allEntitiesArray.filter(entity => isStatusPending(entity));
        case 'approved':
          return allEntitiesArray.filter(entity => isStatusApproved(entity));
        case 'rejected':
          return allEntitiesArray.filter(entity => isStatusRejected(entity));
        default:
          return allEntitiesArray;
      }
    })
    
    // Paginated entities
    const paginatedEntities = computed(() => {
      const filtered = Array.isArray(filteredEntities.value) ? filteredEntities.value : [];
      const start = (currentPage.value - 1) * itemsPerPage.value
      const end = start + itemsPerPage.value
      return filtered.slice(start, end)
    })
    
    // Pagination helpers
    const totalPages = computed(() => {
      const filtered = Array.isArray(filteredEntities.value) ? filteredEntities.value : [];
      return Math.ceil(filtered.length / itemsPerPage.value) || 1
    })
    
    const paginationStart = computed(() => {
      const filtered = Array.isArray(filteredEntities.value) ? filteredEntities.value : [];
      return filtered.length === 0 
        ? 0 
        : (currentPage.value - 1) * itemsPerPage.value + 1
    })
    
    const paginationEnd = computed(() => {
      const filtered = Array.isArray(filteredEntities.value) ? filteredEntities.value : [];
      return Math.min(currentPage.value * itemsPerPage.value, filtered.length)
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
    
    // Helper methods for displaying entity information
    const getEntityName = (entity) => {
      if (!entity) return 'Unknown Temple';
      
      if (entity.name) return entity.name;
      if (entity.Name) return entity.Name;
      
      return (entity.id || entity.ID) ? `Temple ${entity.id || entity.ID}` : 'Unknown Temple';
    }
    
    const getMainDeity = (entity) => {
      if (!entity) return 'No deity specified';
      
      if (entity.main_deity) return entity.main_deity;
      if (entity.MainDeity) return entity.MainDeity;
      
      return 'No deity specified';
    }
    
    const getFullAddress = (entity) => {
      if (!entity) return 'No address provided';
      
      const addressParts = [];
      
      // Street address
      if (entity.street_address || entity.StreetAddress) {
        addressParts.push(entity.street_address || entity.StreetAddress);
      }
      
      // City
      if (entity.city || entity.City) {
        addressParts.push(entity.city || entity.City);
      }
      
      // District
      if (entity.district || entity.District) {
        addressParts.push(entity.district || entity.District);
      }
      
      // State
      if (entity.state || entity.State) {
        addressParts.push(entity.state || entity.State);
      }
      
      // Pincode
      if (entity.pincode || entity.Pincode) {
        addressParts.push(entity.pincode || entity.Pincode);
      }
      
      return addressParts.length > 0 ? addressParts.join(', ') : 'No address provided';
    }
    
    const getEntityInitial = (entity) => {
      if (!entity) return 'T';
      
      const name = getEntityName(entity);
      if (name && name !== 'Unknown Temple' && name.length > 0) {
        return name.charAt(0).toUpperCase();
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
    
    // Apply filters
    const applyFilters = () => {
      currentPage.value = 1
      
      if (statusFilter.value === '') {
        loadAllEntities()
      }
    }
    
    // Handle View Details
    const handleViewDetails = (entity) => {
      selectedEntity.value = entity;
      console.log('Opening details for entity:', entity);
      showDetailsModal.value = true;
    }
    
    // Handle Approve from Details modal
    const handleApproveFromDetails = () => {
      handleApprove(selectedEntity.value);
      showDetailsModal.value = false;
    }
    
    // Debug entity details
    const debugEntityDetails = (entity) => {
      if (!entity) return;
      
      console.log('=== ENTITY DEBUG ===');
      console.log('Full entity object:', entity);
      
      // Show all available keys
      console.log('Available keys on entity:', Object.keys(entity));
      console.log('==================');
    };
    
    // Load pending entities (default view)
    const loadEntities = async () => {
      loading.value = true
      
      try {
        console.log('Loading temple entities data...')
        
        try {
          // First try the correct API endpoint
          const response = await fetch('/api/v1/superadmin/entities?status=pending', {
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
            // Ensure data.data is an array
            const entitiesData = Array.isArray(data.data) ? data.data : [];
            entities.value = entitiesData
            allEntities.value = entitiesData
            console.log(`Loaded ${entitiesData.length} temple entities from API`)
            toast.success(`Loaded ${entitiesData.length} temple entities`)
          } else if (Array.isArray(data)) {
            entities.value = data
            allEntities.value = data
            console.log(`Loaded ${data.length} temple entities from API`)
            toast.success(`Loaded ${data.length} temple entities`)
          } else {
            console.warn('API returned unexpected format:', data)
            // Set empty arrays as fallback
            entities.value = []
            allEntities.value = []
            toast.error('Could not load temple data: unexpected API response format')
          }
        } catch (apiError) {
          console.error('Error calling API directly:', apiError)
          
          // Try using the service if available
          try {
            // Try to get a function for pending entities from the superAdminService
            let serviceFunction = superAdminService.getPendingEntities;
            
            // If that doesn't exist, try alternative function names
            if (!serviceFunction) serviceFunction = superAdminService.getEntitiesWithFilters;
            if (!serviceFunction) serviceFunction = superAdminService.getEntities;
            
            if (serviceFunction) {
              const serviceResponse = await serviceFunction('pending');
              console.log('Service response:', serviceResponse)
              
              if (serviceResponse && serviceResponse.success && serviceResponse.data) {
                // Check different possible structures from the service
                let entitiesData = []
                
                if (serviceResponse.data.pending_entities) {
                  entitiesData = Array.isArray(serviceResponse.data.pending_entities) ? serviceResponse.data.pending_entities : [];
                } else if (Array.isArray(serviceResponse.data)) {
                  entitiesData = serviceResponse.data
                } else {
                  console.warn('Unexpected service response format:', serviceResponse)
                  entitiesData = []
                }
                
                entities.value = entitiesData
                allEntities.value = entitiesData
                
                if (entitiesData.length > 0) {
                  toast.success(`Loaded ${entitiesData.length} temple entities`)
                } else {
                  toast.info('No pending temple entities found')
                }
              } else {
                console.warn('Service call failed or returned unexpected format')
                entities.value = []
                allEntities.value = []
                toast.error('Could not load temple data from API')
              }
            } else {
              // If no service function is available, try a direct API call with axios
              const axiosResponse = await api.get('/v1/entities/by-creator');
              console.log('Axios response:', axiosResponse);
              
              if (axiosResponse && axiosResponse.data) {
                let entitiesData = Array.isArray(axiosResponse.data) ? axiosResponse.data : [];
                entities.value = entitiesData;
                allEntities.value = entitiesData;
                
                if (entitiesData.length > 0) {
                  toast.success(`Loaded ${entitiesData.length} temple entities`);
                } else {
                  toast.info('No pending temple entities found');
                }
              }
            }
          } catch (serviceError) {
            console.error('Error calling service:', serviceError)
            entities.value = []
            allEntities.value = []
            toast.error('Could not load temple data')
          }
        }
      } catch (error) {
        console.error('Error in loadEntities:', error)
        entities.value = []
        allEntities.value = []
        toast.error('Error loading temple data')
      } finally {
        loading.value = false
      }
    }
    
    // Load all entities (for "View All" option)
    const loadAllEntities = async () => {
      loading.value = true
      
      try {
        // Try to call a service method if it exists
        let serviceFunction = superAdminService.getAllEntities;
        
        if (serviceFunction) {
          const response = await serviceFunction();
          
          if (response && response.success && Array.isArray(response.data)) {
            allEntities.value = response.data
            toast.success(`Loaded ${response.data.length} total temples`)
          } else {
            // Fallback to direct API call
            await loadEntityTypesDirectly();
          }
        } else {
          // Try direct API call
          await loadEntityTypesDirectly();
        }
      } catch (error) {
        console.error('Error loading all entities:', error)
        // Fallback to pending entities
        await loadEntities()
        toast.warning('Showing pending temples only')
      } finally {
        loading.value = false
      }
    }
    
    // Helper to load all entity types directly via API
    const loadEntityTypesDirectly = async () => {
      try {
        const statuses = ['pending', 'approved', 'rejected'];
        let allData = [];
        
        // Fetch each status type and combine
        for (const status of statuses) {
          const response = await fetch(`/api/v1/superadmin/entities?status=${status}`, {
            headers: {
              'Authorization': `Bearer ${localStorage.getItem('auth_token')}`
            }
          });
          
          if (response.ok) {
            const data = await response.json();
            if (data && data.data && Array.isArray(data.data)) {
              allData = [...allData, ...data.data];
            }
          }
        }
        
        if (allData.length > 0) {
          allEntities.value = allData;
          toast.success(`Loaded ${allData.length} total temples`);
        } else {
          // If direct API calls fail, fall back to pending entities
          await loadEntities();
          toast.warning('Showing pending temples only');
        }
      } catch (error) {
        console.error('Error in loadEntityTypesDirectly:', error);
        // Fall back to pending entities
        await loadEntities();
      }
    }
    
    // Debug methods
    const debugEntityData = async () => {
      apiDebugInfo.value = 'Loading API debug information...'
      showApiDebugModal.value = true
      
      try {
        const debugInfo = []
        debugInfo.push('==== ENTITY DATA DEBUG ====\n')
        
        // 1. Test direct API call to verify backend
        debugInfo.push('Testing direct API call to /api/v1/superadmin/entities?status=pending...')
        try {
          const response = await fetch('/api/v1/superadmin/entities?status=pending', {
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
            if ((!Array.isArray(allEntities.value) || allEntities.value.length === 0) && data && data.data) {
              const entitiesData = Array.isArray(data.data) ? data.data : [];
              entities.value = entitiesData
              allEntities.value = entitiesData
              debugInfo.push(`\nAUTOMATICALLY LOADED ${entitiesData.length} ENTITIES FROM API RESPONSE`)
            }
          }
        } catch (error) {
          debugInfo.push(`API Error: ${error.message}`)
        }
        
        // 2. Test the "by-creator" API endpoint
        debugInfo.push('\n==== ENTITY BY CREATOR API DEBUG ====\n')
        try {
          const response = await fetch('/api/v1/entities/by-creator', {
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
            if ((!Array.isArray(allEntities.value) || allEntities.value.length === 0) && Array.isArray(data)) {
              entities.value = data
              allEntities.value = data
              debugInfo.push(`\nAUTOMATICALLY LOADED ${data.length} ENTITIES FROM BY-CREATOR API RESPONSE`)
            }
          }
        } catch (error) {
          debugInfo.push(`API Error: ${error.message}`)
        }
        
        // 3. Display entity data
        debugInfo.push('\n==== ENTITY DATA ====\n')
        debugInfo.push(`Total Entities: ${Array.isArray(allEntities.value) ? allEntities.value.length : 0}`)
        if (Array.isArray(allEntities.value) && allEntities.value.length > 0) {
          const firstEntity = allEntities.value[0]
          debugInfo.push(`First Entity Data: ${JSON.stringify(firstEntity, null, 2)}`)
          debugInfo.push(`First Entity ID: ${firstEntity.id || firstEntity.ID || 'Not found'}`)
        }
        
        apiDebugInfo.value = debugInfo.join('\n')
      } catch (error) {
        apiDebugInfo.value = `Error running debug: ${error.message}`
      }
    }
    
    // Direct API test for approval
    const testApprovalApi = async (entity) => {
      if (!entity) {
        toast.error('No entity selected for testing')
        return
      }
      
      apiDebugInfo.value = `Testing API call for entity ${entity.id || entity.ID || 'unknown'}...`
      showApiDebugModal.value = true
      
      try {
        const debugInfo = []
        debugInfo.push(`Testing API approval for entity ${JSON.stringify(entity, null, 2)}\n`)
        
        // Extract entity ID with fallbacks
        const entityId = entity.id || entity.ID
        
        if (!entityId) {
          debugInfo.push('ERROR: No entity ID found in entity object!')
          apiDebugInfo.value = debugInfo.join('\n')
          return
        }
        
        debugInfo.push(`Using entity ID: ${entityId}`)
        
        // Test direct API call
        debugInfo.push('\nTesting direct PATCH call to approve entity...')
        try {
          const token = localStorage.getItem('auth_token')
          debugInfo.push(`Auth token available: ${!!token}`)
          
          const response = await fetch(`/api/v1/superadmin/entities/${entityId}/approval`, {
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
            
            // Refresh the entity list
            loadEntities()
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
    
    // Handle approval
    const handleApprove = async (entity) => {
      if (!entity) {
        toast.error('Cannot approve: Missing entity');
        return;
      }
      
      const entityId = entity.id || entity.ID;
      
      if (!entityId) {
        toast.error('Cannot approve: Missing entity ID');
        return;
      }
      
      isProcessing.value = true;
      console.log(`Attempting to approve entity ${entityId}`);
      
      try {
        // Make a direct API call to the approval endpoint
        const response = await fetch(`/api/v1/superadmin/entities/${entityId}/approval`, {
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
          toast.success(`Temple ${getEntityName(entity)} approved successfully`);
          
          // Refresh the entity list
          loadEntities();
          emit('updated');
        } else {
          console.error('Approval failed with status:', response.status);
          const errorData = await response.text();
          console.error('Error response:', errorData);
          toast.error(`Failed to approve temple: ${response.statusText}`);
        }
      } catch (error) {
        console.error('Error in approval process:', error);
        toast.error('Error approving temple: ' + (error.message || 'Unknown error'));
      } finally {
        isProcessing.value = false;
      }
    };
    
    // Handle reject click
    const handleRejectClick = (entity) => {
      if (!entity) {
        toast.error('Cannot reject: Missing entity');
        return;
      }
      
      const entityId = entity.id || entity.ID;
      
      if (!entityId) {
        toast.error('Cannot reject: Missing entity ID');
        return;
      }
      
      selectedEntity.value = entity;
      rejectReason.value = '';
      showRejectModal.value = true;
    };
    
    // Confirm rejection
    const confirmReject = async () => {
      const entity = selectedEntity.value;
      
      if (!entity) {
        toast.error('Cannot reject: Missing entity');
        showRejectModal.value = false;
        return;
      }
      
      const entityId = entity.id || entity.ID;
      
      if (!entityId) {
        toast.error('Cannot reject: Missing entity ID');
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
        const response = await fetch(`/api/v1/superadmin/entities/${entityId}/approval`, {
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
          toast.success(`Temple ${getEntityName(entity)} rejected successfully`);
          showRejectModal.value = false;
          
          // Refresh the entity list
          loadEntities();
          emit('updated');
        } else {
          console.error('Rejection failed with status:', response.status);
          const errorData = await response.text();
          console.error('Error response:', errorData);
          toast.error(`Failed to reject temple: ${response.statusText}`);
        }
      } catch (error) {
        console.error('Error in rejection process:', error);
        toast.error('Error rejecting temple: ' + (error.message || 'Unknown error'));
      } finally {
        isProcessing.value = false;
      }
    };
    
    // Load data on mount
    onMounted(() => {
      loadEntities()
    })
    
    return {
      loading,
      entities,
      allEntities,
      pendingCount,
      approvedCount,
      rejectedCount,
      statusFilter,
      currentPage,
      itemsPerPage,
      filteredEntities,
      paginatedEntities,
      totalPages,
      paginationStart,
      paginationEnd,
      displayedPageNumbers,
      getStatusBadgeClass,
      formatDate,
      loadEntities,
      loadAllEntities,
      handleApprove,
      handleRejectClick,
      getEntityName,
      getMainDeity,
      getFullAddress,
      getEntityInitial,
      goToPage,
      isStatusPending,
      isStatusApproved,
      isStatusRejected,
      applyFilters,
      isProcessing,
      showRejectModal,
      selectedEntity,
      rejectReason,
      confirmReject,
      // Details modal
      showDetailsModal,
      handleViewDetails,
      handleApproveFromDetails,
      debugEntityDetails,
      // Debug methods
      debugMode,
      debugEntityData,
      testApprovalApi,
      showApiDebugModal,
      apiDebugInfo
    }
  }
}
</script>