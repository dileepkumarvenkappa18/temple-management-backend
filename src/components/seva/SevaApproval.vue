<template>
  <div class="bg-white rounded-xl shadow-md border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="bg-gradient-to-r from-indigo-50 to-indigo-100 px-6 py-4 border-b border-indigo-200">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="p-2 bg-white rounded-lg shadow-sm">
            <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <div>
            <h3 class="text-lg font-semibold text-gray-900">Seva Approval</h3>
            <p class="text-sm text-gray-600">Review and approve seva bookings</p>
          </div>
        </div>
        <div class="flex items-center space-x-2">
          <span class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
            {{ pendingSevas.length }} Pending Review
          </span>
        </div>
      </div>
    </div>

    <!-- Approval List -->
    <div class="divide-y divide-gray-100">
      <div v-if="pendingSevas.length === 0" class="px-6 py-12 text-center">
        <div class="w-16 h-16 mx-auto mb-4 bg-gray-100 rounded-full flex items-center justify-center">
          <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <h3 class="text-lg font-medium text-gray-900 mb-2">No Pending Approvals</h3>
        <p class="text-gray-500">All seva bookings have been reviewed.</p>
      </div>

      <div 
        v-for="seva in pendingSevas" 
        :key="seva.id"
        class="px-6 py-5 hover:bg-gray-50 transition-colors duration-200"
      >
        <div class="flex items-start justify-between space-x-4">
          <!-- Seva Details -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center space-x-3 mb-3">
              <div class="w-12 h-12 rounded-lg bg-gradient-to-br from-indigo-100 to-indigo-200 flex items-center justify-center">
                <span class="text-indigo-700 font-semibold text-sm">{{ seva.devotee.name.charAt(0) }}</span>
              </div>
              <div class="flex-1 min-w-0">
                <h4 class="text-base font-semibold text-gray-900 truncate">{{ seva.sevaType }}</h4>
                <p class="text-sm text-gray-600">{{ seva.devotee.name }} • {{ seva.devotee.phone }}</p>
              </div>
              <div class="flex flex-col items-end space-y-1">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800">
                  Pending
                </span>
                <span class="text-xs text-gray-500">{{ formatDate(seva.requestedDate) }}</span>
              </div>
            </div>

            <!-- Seva Information -->
            <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
              <div class="bg-gray-50 rounded-lg p-3">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Seva Date</label>
                <p class="text-sm font-semibold text-gray-900 mt-1">{{ formatDate(seva.sevaDate) }}</p>
              </div>
              <div class="bg-gray-50 rounded-lg p-3">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Amount</label>
                <p class="text-sm font-semibold text-gray-900 mt-1">₹{{ seva.amount.toLocaleString() }}</p>
              </div>
              <div class="bg-gray-50 rounded-lg p-3">
                <label class="text-xs font-medium text-gray-500 uppercase tracking-wide">Time Slot</label>
                <p class="text-sm font-semibold text-gray-900 mt-1">{{ seva.timeSlot }}</p>
              </div>
            </div>

            <!-- Special Requests -->
            <div v-if="seva.specialRequests" class="mb-4">
              <label class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2 block">Special Requests</label>
              <div class="bg-blue-50 rounded-lg p-3 border border-blue-200">
                <p class="text-sm text-blue-900">{{ seva.specialRequests }}</p>
              </div>
            </div>

            <!-- Admin Notes Input -->
            <div class="mb-4">
              <label class="text-xs font-medium text-gray-500 uppercase tracking-wide mb-2 block">Admin Notes</label>
              <textarea
                v-model="adminNotes[seva.id]"
                rows="2"
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200 resize-none"
                placeholder="Add notes for approval/rejection..."
              ></textarea>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex items-center justify-end space-x-3 pt-4 border-t border-gray-100">
          <button
            @click="rejectSeva(seva)"
            :disabled="processing[seva.id]"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 hover:border-red-300 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            <svg v-if="processing[seva.id] === 'rejecting'" class="animate-spin -ml-1 mr-2 h-4 w-4 text-red-700" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
            {{ processing[seva.id] === 'rejecting' ? 'Rejecting...' : 'Reject' }}
          </button>
          
          <button
            @click="approveSeva(seva)"
            :disabled="processing[seva.id]"
            class="inline-flex items-center px-6 py-2 text-sm font-medium text-white bg-gradient-to-r from-indigo-600 to-indigo-700 rounded-lg hover:from-indigo-700 hover:to-indigo-800 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-sm"
          >
            <svg v-if="processing[seva.id] === 'approving'" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg v-else class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
            </svg>
            {{ processing[seva.id] === 'approving' ? 'Approving...' : 'Approve' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Bulk Actions Footer -->
    <div v-if="selectedSevas.length > 0" class="bg-gray-50 px-6 py-4 border-t border-gray-200">
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-600">
          {{ selectedSevas.length }} seva(s) selected
        </div>
        <div class="flex items-center space-x-3">
          <button
            @click="bulkReject"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-red-700 bg-white border border-red-300 rounded-lg hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-red-500 focus:ring-offset-2 transition-all duration-200"
          >
            Bulk Reject
          </button>
          <button
            @click="bulkApprove"
            class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
          >
            Bulk Approve
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useSevaStore } from '@/stores/seva'
import { useToast } from '@/composables/useToast'

const sevaStore = useSevaStore()
const { showToast } = useToast()

// Reactive state
const adminNotes = reactive({})
const processing = reactive({})
const selectedSevas = ref([])

// Computed properties
const pendingSevas = computed(() => 
  sevaStore.sevas.filter(seva => seva.status === 'PENDING')
)

// Methods
const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-IN', {
    weekday: 'short',
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const approveSeva = async (seva) => {
  processing[seva.id] = 'approving'
  
  try {
    await sevaStore.updateSevaStatus(seva.id, {
      status: 'APPROVED',
      adminNotes: adminNotes[seva.id] || '',
      approvedAt: new Date().toISOString(),
      approvedBy: 'current-admin-id' // This should come from auth store
    })
    
    showToast('Seva booking approved successfully!', 'success')
    
    // Send notification to devotee
    await sevaStore.notifyDevotee(seva.devotee.id, {
      type: 'SEVA_APPROVED',
      message: `Your seva booking for ${seva.sevaType} has been approved.`,
      sevaId: seva.id
    })
    
  } catch (error) {
    showToast('Failed to approve seva booking', 'error')
    console.error('Error approving seva:', error)
  } finally {
    processing[seva.id] = null
  }
}

const rejectSeva = async (seva) => {
  processing[seva.id] = 'rejecting'
  
  try {
    await sevaStore.updateSevaStatus(seva.id, {
      status: 'REJECTED',
      adminNotes: adminNotes[seva.id] || '',
      rejectedAt: new Date().toISOString(),
      rejectedBy: 'current-admin-id' // This should come from auth store
    })
    
    showToast('Seva booking rejected', 'info')
    
    // Send notification to devotee
    await sevaStore.notifyDevotee(seva.devotee.id, {
      type: 'SEVA_REJECTED',
      message: `Your seva booking for ${seva.sevaType} has been rejected. ${adminNotes[seva.id] ? 'Reason: ' + adminNotes[seva.id] : ''}`,
      sevaId: seva.id
    })
    
  } catch (error) {
    showToast('Failed to reject seva booking', 'error')
    console.error('Error rejecting seva:', error)
  } finally {
    processing[seva.id] = null
  }
}

const bulkApprove = async () => {
  try {
    const promises = selectedSevas.value.map(sevaId => {
      const seva = pendingSevas.value.find(s => s.id === sevaId)
      return sevaStore.updateSevaStatus(sevaId, {
        status: 'APPROVED',
        adminNotes: adminNotes[sevaId] || '',
        approvedAt: new Date().toISOString(),
        approvedBy: 'current-admin-id'
      })
    })
    
    await Promise.all(promises)
    showToast(`${selectedSevas.value.length} seva bookings approved!`, 'success')
    selectedSevas.value = []
    
  } catch (error) {
    showToast('Failed to approve selected bookings', 'error')
    console.error('Error in bulk approve:', error)
  }
}

const bulkReject = async () => {
  try {
    const promises = selectedSevas.value.map(sevaId => {
      return sevaStore.updateSevaStatus(sevaId, {
        status: 'REJECTED',
        adminNotes: adminNotes[sevaId] || '',
        rejectedAt: new Date().toISOString(),
        rejectedBy: 'current-admin-id'
      })
    })
    
    await Promise.all(promises)
    showToast(`${selectedSevas.value.length} seva bookings rejected`, 'info')
    selectedSevas.value = []
    
  } catch (error) {
    showToast('Failed to reject selected bookings', 'error')
    console.error('Error in bulk reject:', error)
  }
}

// Lifecycle
onMounted(() => {
  sevaStore.fetchPendingSevas()
})
</script>