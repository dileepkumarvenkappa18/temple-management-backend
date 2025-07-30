<template>
  <div 
    v-if="show"
    class="fixed inset-0 z-50 overflow-y-auto"
    aria-labelledby="modal-title"
    role="dialog"
    aria-modal="true"
  >
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <!-- Background overlay -->
      <div 
        class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" 
        @click="$emit('close')"
      ></div>
      
      <!-- Modal panel -->
      <div class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
        <div class="bg-white px-6 pt-6 pb-4">
          <div class="flex items-start">
            <!-- Icon -->
            <div 
              :class="[
                'mx-auto flex-shrink-0 flex items-center justify-center h-12 w-12 rounded-full',
                isApproval ? 'bg-green-100' : 'bg-red-100'
              ]"
            >
              <svg 
                v-if="isApproval" 
                class="h-6 w-6 text-green-600" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path 
                  stroke-linecap="round" 
                  stroke-linejoin="round" 
                  stroke-width="2" 
                  d="M5 13l4 4L19 7"
                ></path>
              </svg>
              <svg 
                v-else 
                class="h-6 w-6 text-red-600" 
                fill="none" 
                stroke="currentColor" 
                viewBox="0 0 24 24"
              >
                <path 
                  stroke-linecap="round" 
                  stroke-linejoin="round" 
                  stroke-width="2" 
                  d="M6 18L18 6M6 6l12 12"
                ></path>
              </svg>
            </div>

            <!-- Content -->
            <div class="ml-4 w-full text-left">
              <h3 class="text-lg leading-6 font-semibold text-gray-900 font-roboto" id="modal-title">
                {{ isApproval ? 'Approve Tenant Application' : 'Reject Tenant Application' }}
              </h3>
              
              <div class="mt-2">
                <p class="text-sm text-gray-500 mb-4">
                  {{ 
                    isApproval 
                      ? `Are you sure you want to approve ${tenant?.fullName}'s application?` 
                      : `Please provide a reason for rejecting ${tenant?.fullName}'s application:`
                  }}
                </p>

                <!-- Notes field (optional for approval, required for rejection) -->
                <textarea
                  v-model="notes"
                  rows="4"
                  :placeholder="isApproval ? 'Add optional approval notes...' : 'Enter rejection reason (required)...'"
                  class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 text-sm resize-none"
                  :class="isApproval ? 'focus:ring-green-500 focus:border-green-500' : 'focus:ring-red-500 focus:border-red-500'"
                ></textarea>

                <!-- Notification toggle -->
                <div class="mt-4 flex items-center">
                  <input
                    id="sendNotification"
                    v-model="sendNotification"
                    type="checkbox"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <label for="sendNotification" class="ml-2 block text-sm text-gray-700">
                    Send email notification
                  </label>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Footer/Buttons -->
        <div class="bg-gray-50 px-6 py-4 flex flex-col sm:flex-row-reverse gap-3">
          <button
            @click="submitAction"
            :disabled="isProcessing || (!isApproval && !notes.trim())"
            :class="[
              'w-full sm:w-auto px-6 py-2 text-white text-sm font-medium rounded-xl transition-all duration-200 flex items-center justify-center gap-2 disabled:opacity-50',
              isApproval 
                ? 'bg-green-600 hover:bg-green-700' 
                : 'bg-red-600 hover:bg-red-700'
            ]"
          >
            <div v-if="isProcessing" class="animate-spin rounded-full h-4 w-4 border-b-2 border-white"></div>
            {{ isProcessing ? 'Processing...' : (isApproval ? 'Approve' : 'Reject') }}
          </button>
          
          <button
            @click="$emit('close')"
            :disabled="isProcessing"
            class="w-full sm:w-auto px-6 py-2 bg-white hover:bg-gray-50 text-gray-700 text-sm font-medium border border-gray-300 rounded-xl transition-all duration-200"
          >
            Cancel
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, defineProps, defineEmits } from 'vue';
import { useToast } from '@/composables/useToast';
import superAdminService from '@/services/superadmin.service';

const props = defineProps({
  show: {
    type: Boolean,
    required: true
  },
  tenant: {
    type: Object,
    default: () => ({})
  },
  action: {
    type: String,
    default: 'approve',
    validator: (value) => ['approve', 'reject'].includes(value)
  }
});

const emit = defineEmits(['close', 'tenant-updated']);

// Composables
const toast = useToast();

// Reactive state
const notes = ref('');
const sendNotification = ref(true);
const isProcessing = ref(false);

// Computed properties
const isApproval = computed(() => props.action === 'approve');

// Methods
const submitAction = async () => {
  // Validation for rejection
  if (!isApproval.value && !notes.value.trim()) {
    toast.error('Please provide a reason for rejection');
    return;
  }

  isProcessing.value = true;

  try {
    let response;

    if (isApproval.value) {
      // Call approval service
      response = await superAdminService.approveTenant(props.tenant.id, {
        notes: notes.value,
        sendEmail: sendNotification.value
      });
    } else {
      // Call rejection service
      response = await superAdminService.rejectTenant(props.tenant.id, {
        notes: notes.value,
        sendEmail: sendNotification.value
      });
    }

    if (response.success) {
      toast.success(response.message || `Tenant ${isApproval.value ? 'approved' : 'rejected'} successfully`);
      emit('tenant-updated', { 
        id: props.tenant.id, 
        status: isApproval.value ? 'APPROVED' : 'REJECTED',
        notes: notes.value
      });
      emit('close');
    } else {
      toast.error(response.message || `Failed to ${isApproval.value ? 'approve' : 'reject'} tenant`);
    }
  } catch (error) {
    console.error(`Error ${isApproval.value ? 'approving' : 'rejecting'} tenant:`, error);
    toast.error(`An error occurred while ${isApproval.value ? 'approving' : 'rejecting'} the tenant`);
  } finally {
    isProcessing.value = false;
  }
};
</script>

<style scoped>
/* Custom font classes */
.font-roboto {
  font-family: 'Roboto', sans-serif;
}

/* Modal backdrop blur effect */
.fixed.inset-0.bg-gray-500 {
  backdrop-filter: blur(4px);
}

/* Enhanced focus states for accessibility */
button:focus,
textarea:focus,
input:focus {
  outline: 2px solid #6366f1;
  outline-offset: 2px;
}

/* Button loading state */
button:disabled {
  cursor: not-allowed;
}

/* Animation for processing spinner */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}
</style>