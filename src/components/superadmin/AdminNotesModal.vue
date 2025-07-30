<template>
  <div
    v-if="isVisible"
    class="fixed inset-0 z-50 overflow-y-auto"
    @click.self="closeModal"
  >
    <!-- Backdrop -->
    <div class="fixed inset-0 bg-black bg-opacity-50 transition-opacity"></div>
    
    <!-- Modal -->
    <div class="flex min-h-full items-center justify-center p-4">
      <div
        class="relative w-full max-w-2xl transform rounded-2xl bg-white shadow-2xl transition-all"
        @click.stop
      >
        <!-- Header -->
        <div class="flex items-center justify-between border-b border-gray-200 px-6 py-4">
          <div class="flex items-center space-x-3">
            <div class="flex h-10 w-10 items-center justify-center rounded-xl bg-indigo-100">
              <svg class="h-5 w-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
              </svg>
            </div>
            <div>
              <h2 class="text-xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
                {{ modalTitle }}
              </h2>
              <p class="text-sm text-gray-600" style="font-family: sans-serif;">
                {{ modalSubtitle }}
              </p>
            </div>
          </div>
          
          <button
            @click="closeModal"
            class="rounded-xl p-2 text-gray-400 hover:bg-gray-100 hover:text-gray-600 transition-all duration-200"
          >
            <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>

        <!-- Body -->
        <div class="px-6 py-6">
          <!-- Temple/Entity Info -->
          <div class="mb-6 rounded-xl bg-gray-50 p-4">
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-2">
              <div>
                <h3 class="text-sm font-semibold text-gray-700" style="font-family: sans-serif;">Temple Name</h3>
                <p class="mt-1 text-sm text-gray-900">{{ entityData.name || 'N/A' }}</p>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-700" style="font-family: sans-serif;">Admin Name</h3>
                <p class="mt-1 text-sm text-gray-900">{{ entityData.adminName || 'N/A' }}</p>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-700" style="font-family: sans-serif;">Email</h3>
                <p class="mt-1 text-sm text-gray-900">{{ entityData.email || 'N/A' }}</p>
              </div>
              <div>
                <h3 class="text-sm font-semibold text-gray-700" style="font-family: sans-serif;">Status</h3>
                <span
                  :class="getStatusClass(entityData.status)"
                  class="mt-1 inline-flex rounded-lg px-2 py-1 text-xs font-medium"
                >
                  {{ entityData.status || 'N/A' }}
                </span>
              </div>
            </div>
          </div>

          <!-- Current Notes (if any) -->
          <div v-if="entityData.notes" class="mb-6">
            <h3 class="mb-3 text-sm font-semibold text-gray-700" style="font-family: sans-serif;">Previous Notes</h3>
            <div class="rounded-xl bg-yellow-50 border border-yellow-200 p-4">
              <p class="text-sm text-yellow-800">{{ entityData.notes }}</p>
            </div>
          </div>

          <!-- Notes Form -->
          <div class="space-y-4">
            <div>
              <label for="notes" class="block text-sm font-semibold text-gray-700 mb-2" style="font-family: sans-serif;">
                {{ actionType === 'approve' ? 'Approval Notes (Optional)' : 'Rejection Reason' }}
                <span v-if="actionType === 'reject'" class="text-red-500">*</span>
              </label>
              <textarea
                id="notes"
                v-model="notes"
                rows="5"
                :placeholder="getPlaceholder()"
                class="w-full rounded-xl border border-gray-300 px-4 py-3 text-sm text-gray-900 placeholder-gray-500 focus:border-indigo-500 focus:outline-none focus:ring-2 focus:ring-indigo-200 transition-all duration-200"
                :class="{ 'border-red-300 focus:border-red-500 focus:ring-red-200': actionType === 'reject' && !notes.trim() && showValidation }"
              ></textarea>
              <div v-if="actionType === 'reject' && !notes.trim() && showValidation" class="mt-2 text-sm text-red-600">
                Please provide a reason for rejection
              </div>
            </div>

            <!-- Character count -->
            <div class="flex justify-between text-xs text-gray-500">
              <span>{{ notes.length }}/500 characters</span>
              <span v-if="notes.length > 450" class="text-orange-600">
                {{ 500 - notes.length }} characters remaining
              </span>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="flex items-center justify-end space-x-3 border-t border-gray-200 px-6 py-4">
          <button
            @click="closeModal"
            class="rounded-xl border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-200 transition-all duration-200"
          >
            Cancel
          </button>
          
          <button
            @click="handleSubmit"
            :disabled="isSubmitting || (actionType === 'reject' && !notes.trim())"
            :class="getSubmitButtonClass()"
            class="rounded-xl px-6 py-2 text-sm font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center space-x-2"
          >
            <svg v-if="isSubmitting" class="h-4 w-4 animate-spin" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ getSubmitButtonText() }}</span>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

// Props
const props = defineProps({
  isVisible: {
    type: Boolean,
    default: false
  },
  entityData: {
    type: Object,
    default: () => ({})
  },
  actionType: {
    type: String,
    default: 'approve', // 'approve' or 'reject'
    validator: (value) => ['approve', 'reject'].includes(value)
  }
})

// Emits
const emit = defineEmits(['close', 'submit'])

// Reactive data
const notes = ref('')
const isSubmitting = ref(false)
const showValidation = ref(false)

// Computed properties
const modalTitle = computed(() => {
  return props.actionType === 'approve' ? 'Approve Temple' : 'Reject Temple'
})

const modalSubtitle = computed(() => {
  return props.actionType === 'approve' 
    ? 'Add optional notes for temple approval'
    : 'Please provide a reason for rejection'
})

// Methods
const closeModal = () => {
  if (!isSubmitting.value) {
    notes.value = ''
    showValidation.value = false
    emit('close')
  }
}

const handleSubmit = async () => {
  if (props.actionType === 'reject' && !notes.value.trim()) {
    showValidation.value = true
    return
  }

  if (notes.value.length > 500) {
    return
  }

  try {
    isSubmitting.value = true
    
    const payload = {
      entityId: props.entityData.id,
      action: props.actionType,
      notes: notes.value.trim()
    }

    emit('submit', payload)
    
    // Reset form
    notes.value = ''
    showValidation.value = false
    
  } catch (error) {
    console.error('Error submitting notes:', error)
  } finally {
    isSubmitting.value = false
  }
}

const getStatusClass = (status) => {
  const statusClasses = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'APPROVED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800'
  }
  return statusClasses[status] || 'bg-gray-100 text-gray-800'
}

const getSubmitButtonClass = () => {
  if (props.actionType === 'approve') {
    return 'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500'
  } else {
    return 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500'
  }
}

const getSubmitButtonText = () => {
  if (isSubmitting.value) {
    return props.actionType === 'approve' ? 'Approving...' : 'Rejecting...'
  }
  return props.actionType === 'approve' ? 'Approve Temple' : 'Reject Temple'
}

const getPlaceholder = () => {
  if (props.actionType === 'approve') {
    return 'Add any additional notes or instructions for the temple admin... (Optional)'
  } else {
    return 'Please provide a clear reason for rejection. This will be sent to the temple admin.'
  }
}

// Watch for character limit
watch(notes, (newValue) => {
  if (newValue.length > 500) {
    notes.value = newValue.substring(0, 500)
  }
})

// Reset validation when modal closes
watch(() => props.isVisible, (newValue) => {
  if (!newValue) {
    showValidation.value = false
  }
})
</script>