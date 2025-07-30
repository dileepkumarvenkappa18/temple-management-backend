<template>
  <Teleport to="body">
    <div
      v-if="modelValue"
      class="fixed inset-0 z-50 overflow-y-auto"
      aria-labelledby="modal-title"
      role="dialog"
      aria-modal="true"
    >
      <!-- Background overlay -->
      <div
        class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0"
      >
        <Transition
          enter-active-class="transition-opacity ease-out duration-300"
          enter-from-class="opacity-0"
          enter-to-class="opacity-100"
          leave-active-class="transition-opacity ease-in duration-200"
          leave-from-class="opacity-100"
          leave-to-class="opacity-0"
        >
          <div
            v-if="modelValue"
            class="fixed inset-0 bg-gray-500 bg-opacity-75 backdrop-blur-sm transition-opacity"
            @click="handleBackdropClick"
          ></div>
        </Transition>

        <!-- This element is to trick the browser into centering the modal contents. -->
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen" aria-hidden="true">&#8203;</span>

        <!-- Modal panel -->
        <Transition
          enter-active-class="transition-all ease-out duration-300"
          enter-from-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
          enter-to-class="opacity-100 translate-y-0 sm:scale-100"
          leave-active-class="transition-all ease-in duration-200"
          leave-from-class="opacity-100 translate-y-0 sm:scale-100"
          leave-to-class="opacity-0 translate-y-4 sm:translate-y-0 sm:scale-95"
        >
          <div
            v-if="modelValue"
            class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:w-full border border-gray-200"
            :class="sizeClasses"
          >
            <!-- Modal Header -->
            <div v-if="showHeader" class="bg-white px-6 py-4 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <div class="flex items-center">
                  <!-- Icon -->
                  <div v-if="icon" class="flex-shrink-0 mr-4">
                    <div class="mx-auto flex items-center justify-center h-10 w-10 rounded-full" :class="iconBgClass">
                      <component :is="icon" class="h-5 w-5" :class="iconClass" />
                    </div>
                  </div>
                  
                  <!-- Title -->
                  <div>
                    <h3 class="text-lg font-semibold text-gray-900" id="modal-title">
                      {{ title }}
                    </h3>
                    <p v-if="subtitle" class="text-sm text-gray-600 mt-1">
                      {{ subtitle }}
                    </p>
                  </div>
                </div>
                
                <!-- Close button -->
                <button
                  v-if="showCloseButton"
                  @click="handleClose"
                  class="rounded-md text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-colors duration-200"
                >
                  <span class="sr-only">Close</span>
                  <XMarkIcon class="h-5 w-5" />
                </button>
              </div>
            </div>

            <!-- Modal Body -->
            <div class="bg-white" :class="bodyClass">
              <div v-if="$slots.default" class="text-gray-700">
                <slot />
              </div>
              <div v-else-if="message" class="text-gray-700">
                <p>{{ message }}</p>
              </div>
            </div>

            <!-- Modal Footer -->
            <div v-if="showFooter" class="bg-gray-50 px-6 py-4 border-t border-gray-200 flex justify-end space-x-3">
              <slot name="footer">
                <!-- Cancel Button -->
                <button
                  v-if="showCancelButton"
                  @click="handleCancel"
                  type="button"
                  class="inline-flex justify-center rounded-md border border-gray-300 shadow-sm px-4 py-2 bg-white text-base font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200 sm:text-sm"
                  :disabled="loading"
                >
                  {{ cancelText }}
                </button>
                
                <!-- Confirm Button -->
                <button
                  v-if="showConfirmButton"
                  @click="handleConfirm"
                  type="button"
                  class="inline-flex items-center justify-center rounded-md shadow-sm px-4 py-2 text-base font-medium focus:outline-none focus:ring-2 focus:ring-offset-2 transition-all duration-200 sm:text-sm"
                  :class="confirmButtonClass"
                  :disabled="loading"
                >
                  <svg
                    v-if="loading"
                    class="animate-spin -ml-1 mr-2 h-4 w-4 text-white"
                    xmlns="http://www.w3.org/2000/svg"
                    fill="none"
                    viewBox="0 0 24 24"
                  >
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  {{ loading ? loadingText : confirmText }}
                </button>
              </slot>
            </div>
          </div>
        </Transition>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { computed, nextTick, onMounted, onUnmounted, toRefs, watch } from 'vue'
import { XMarkIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  modelValue: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  subtitle: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  icon: {
    type: [String, Object],
    default: null
  },
  variant: {
    type: String,
    default: 'default', // default, danger, warning, success
    validator: value => ['default', 'danger', 'warning', 'success'].includes(value)
  },
  size: {
    type: String,
    default: 'md', // sm, md, lg, xl, full
    validator: value => ['sm', 'md', 'lg', 'xl', 'full'].includes(value)
  },
  
  // Header options
  showHeader: {
    type: Boolean,
    default: true
  },
  showCloseButton: {
    type: Boolean,
    default: true
  },
  
  // Footer options
  showFooter: {
    type: Boolean,
    default: true
  },
  showCancelButton: {
    type: Boolean,
    default: true
  },
  showConfirmButton: {
    type: Boolean,
    default: true
  },
  
  // Button text
  cancelText: {
    type: String,
    default: 'Cancel'
  },
  confirmText: {
    type: String,
    default: 'Confirm'
  },
  loadingText: {
    type: String,
    default: 'Loading...'
  },
  
  // State
  loading: {
    type: Boolean,
    default: false
  },
  
  // Behavior
  closeOnBackdrop: {
    type: Boolean,
    default: true
  },
  closeOnEscape: {
    type: Boolean,
    default: true
  },
  persistent: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'confirm', 'cancel', 'close'])

// Computed properties
const sizeClasses = computed(() => {
  const sizes = {
    sm: 'sm:max-w-sm',
    md: 'sm:max-w-md',
    lg: 'sm:max-w-lg',
    xl: 'sm:max-w-2xl',
    full: 'sm:max-w-7xl'
  }
  return sizes[props.size]
})

const bodyClass = computed(() => {
  return props.showHeader || props.showFooter ? 'px-6 py-4' : 'p-6'
})

const iconBgClass = computed(() => {
  const classes = {
    default: 'bg-indigo-100',
    danger: 'bg-red-100',
    warning: 'bg-yellow-100',
    success: 'bg-green-100'
  }
  return classes[props.variant]
})

const iconClass = computed(() => {
  const classes = {
    default: 'text-indigo-600',
    danger: 'text-red-600',
    warning: 'text-yellow-600',
    success: 'text-green-600'
  }
  return classes[props.variant]
})

const confirmButtonClass = computed(() => {
  const classes = {
    default: 'border border-transparent bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500',
    danger: 'border border-transparent bg-red-600 text-white hover:bg-red-700 focus:ring-red-500',
    warning: 'border border-transparent bg-yellow-600 text-white hover:bg-yellow-700 focus:ring-yellow-500',
    success: 'border border-transparent bg-green-600 text-white hover:bg-green-700 focus:ring-green-500'
  }
  return `${classes[props.variant]} ${props.loading ? 'opacity-75 cursor-not-allowed' : ''}`
})

// Event handlers
const handleClose = () => {
  if (!props.persistent && !props.loading) {
    emit('update:modelValue', false)
    emit('close')
  }
}

const handleBackdropClick = () => {
  if (props.closeOnBackdrop) {
    handleClose()
  }
}

const handleCancel = () => {
  if (!props.loading) {
    emit('cancel')
    emit('update:modelValue', false)
  }
}

const handleConfirm = () => {
  if (!props.loading) {
    emit('confirm')
  }
}

const handleEscape = (event) => {
  if (event.key === 'Escape' && props.closeOnEscape && props.modelValue) {
    handleClose()
  }
}

// Lifecycle
onMounted(() => {
  document.addEventListener('keydown', handleEscape)
  
  // Prevent body scroll when modal is open
  if (props.modelValue) {
    document.body.style.overflow = 'hidden'
  }
})

onUnmounted(() => {
  document.removeEventListener('keydown', handleEscape)
  document.body.style.overflow = ''
})

// Watch for model value changes to handle body scroll
const { modelValue } = toRefs(props)
watch(modelValue, (newValue) => {
  nextTick(() => {
    if (newValue) {
      document.body.style.overflow = 'hidden'
    } else {
      document.body.style.overflow = ''
    }
  })
})
</script>

<style scoped>
/* Ensure modal appears above everything */
.z-50 {
  z-index: 50;
}

/* Backdrop blur for modern browsers */
.backdrop-blur-sm {
  backdrop-filter: blur(4px);
}
</style>