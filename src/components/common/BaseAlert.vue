<template>
  <Transition
    enter-active-class="transition-all duration-300 ease-out"
    enter-from-class="opacity-0 transform -translate-y-2 scale-95"
    enter-to-class="opacity-100 transform translate-y-0 scale-100"
    leave-active-class="transition-all duration-200 ease-in"
    leave-from-class="opacity-100 transform translate-y-0 scale-100"
    leave-to-class="opacity-0 transform -translate-y-2 scale-95"
  >
    <div
      v-if="visible"
      :class="alertClasses"
      class="relative rounded-xl shadow-lg border transition-all duration-200 hover:shadow-xl"
      role="alert"
      :aria-live="type === 'error' ? 'assertive' : 'polite'"
    >
      <!-- Icon Section -->
      <div class="flex items-start gap-3 p-4">
        <div :class="iconWrapperClasses" class="flex-shrink-0 rounded-full p-1">
          <component 
            :is="iconComponent" 
            :class="iconClasses"
            class="w-5 h-5"
          />
        </div>

        <!-- Content Section -->
        <div class="flex-1 min-w-0">
          <!-- Title -->
          <h4 
            v-if="title" 
            :class="titleClasses"
            class="font-bold text-sm font-roboto leading-5 mb-1"
          >
            {{ title }}
          </h4>

          <!-- Message -->
          <div :class="messageClasses" class="text-sm font-sans leading-5">
            <slot>
              {{ message }}
            </slot>
          </div>

          <!-- Actions Slot -->
          <div v-if="$slots.actions || showActions" class="mt-3 flex gap-2">
            <slot name="actions">
              <button
                v-if="actionText"
                @click="handleAction"
                :class="actionButtonClasses"
                class="text-xs font-medium px-3 py-1.5 rounded-lg transition-all duration-200 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2"
              >
                {{ actionText }}
              </button>
            </slot>
          </div>
        </div>

        <!-- Close Button -->
        <button
          v-if="dismissible"
          @click="dismiss"
          :class="closeButtonClasses"
          class="flex-shrink-0 rounded-lg p-1.5 transition-all duration-200 hover:shadow-md focus:outline-none focus:ring-2 focus:ring-offset-2"
          aria-label="Dismiss alert"
        >
          <XMarkIcon class="w-4 h-4" />
        </button>
      </div>

      <!-- Progress Bar (for auto-dismiss) -->
      <div
        v-if="autoDismiss && autoDismissTime"
        class="absolute bottom-0 left-0 right-0 h-1 bg-black/10 rounded-b-xl overflow-hidden"
      >
        <div
          class="h-full bg-current transition-all ease-linear"
          :style="{ width: `${progressWidth}%`, transitionDuration: `${autoDismissTime}ms` }"
        ></div>
      </div>
    </div>
  </Transition>
</template>

<script setup>
import { computed, ref, onMounted, onUnmounted } from 'vue'
import { 
  CheckCircleIcon, 
  ExclamationTriangleIcon, 
  XCircleIcon, 
  InformationCircleIcon,
  XMarkIcon 
} from '@heroicons/vue/24/outline'

const props = defineProps({
  type: {
    type: String,
    default: 'info',
    validator: (value) => ['success', 'error', 'warning', 'info'].includes(value)
  },
  title: {
    type: String,
    default: ''
  },
  message: {
    type: String,
    default: ''
  },
  dismissible: {
    type: Boolean,
    default: true
  },
  autoDismiss: {
    type: Boolean,
    default: false
  },
  autoDismissTime: {
    type: Number,
    default: 5000
  },
  actionText: {
    type: String,
    default: ''
  },
  showActions: {
    type: Boolean,
    default: false
  },
  persistent: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['dismiss', 'action'])

const visible = ref(true)
const progressWidth = ref(0)
let dismissTimer = null

// Icon mapping
const iconComponent = computed(() => {
  const icons = {
    success: CheckCircleIcon,
    error: XCircleIcon,
    warning: ExclamationTriangleIcon,
    info: InformationCircleIcon
  }
  return icons[props.type]
})

// Dynamic classes based on type
const alertClasses = computed(() => {
  const baseClasses = 'w-full max-w-md mx-auto sm:max-w-lg md:max-w-xl lg:max-w-2xl'
  const typeClasses = {
    success: 'bg-green-50 border-green-200 text-green-900',
    error: 'bg-red-50 border-red-200 text-red-900',
    warning: 'bg-yellow-50 border-yellow-200 text-yellow-900',
    info: 'bg-blue-50 border-blue-200 text-blue-900'
  }
  return `${baseClasses} ${typeClasses[props.type]}`
})

const iconWrapperClasses = computed(() => {
  const classes = {
    success: 'bg-green-100 text-green-600',
    error: 'bg-red-100 text-red-600',
    warning: 'bg-yellow-100 text-yellow-600',
    info: 'bg-blue-100 text-blue-600'
  }
  return classes[props.type]
})

const iconClasses = computed(() => {
  const classes = {
    success: 'text-green-600',
    error: 'text-red-600',
    warning: 'text-yellow-600',
    info: 'text-blue-600'
  }
  return classes[props.type]
})

const titleClasses = computed(() => {
  const classes = {
    success: 'text-green-900',
    error: 'text-red-900',
    warning: 'text-yellow-900',
    info: 'text-blue-900'
  }
  return classes[props.type]
})

const messageClasses = computed(() => {
  const classes = {
    success: 'text-green-800',
    error: 'text-red-800',
    warning: 'text-yellow-800',
    info: 'text-blue-800'
  }
  return classes[props.type]
})

const actionButtonClasses = computed(() => {
  const classes = {
    success: 'bg-green-100 text-green-800 hover:bg-green-200 focus:ring-green-500',
    error: 'bg-red-100 text-red-800 hover:bg-red-200 focus:ring-red-500',
    warning: 'bg-yellow-100 text-yellow-800 hover:bg-yellow-200 focus:ring-yellow-500',
    info: 'bg-blue-100 text-blue-800 hover:bg-blue-200 focus:ring-blue-500'
  }
  return classes[props.type]
})

const closeButtonClasses = computed(() => {
  const classes = {
    success: 'text-green-500 hover:bg-green-100 focus:ring-green-500',
    error: 'text-red-500 hover:bg-red-100 focus:ring-red-500',
    warning: 'text-yellow-500 hover:bg-yellow-100 focus:ring-yellow-500',
    info: 'text-blue-500 hover:bg-blue-100 focus:ring-blue-500'
  }
  return classes[props.type]
})

const dismiss = () => {
  if (props.persistent) return
  visible.value = false
  if (dismissTimer) {
    clearTimeout(dismissTimer)
  }
  setTimeout(() => {
    emit('dismiss')
  }, 300) // Wait for exit animation
}

const handleAction = () => {
  emit('action')
}

onMounted(() => {
  if (props.autoDismiss && props.autoDismissTime && !props.persistent) {
    // Start progress bar animation
    setTimeout(() => {
      progressWidth.value = 100
    }, 100)
    
    // Auto dismiss after specified time
    dismissTimer = setTimeout(() => {
      dismiss()
    }, props.autoDismissTime)
  }
})

onUnmounted(() => {
  if (dismissTimer) {
    clearTimeout(dismissTimer)
  }
})

defineExpose({
  dismiss
})
</script>

<style scoped>
/* Additional responsive utilities */
@media (max-width: 640px) {
  .font-roboto {
    font-size: 0.875rem;
  }
}

@media (min-width: 768px) {
  .font-roboto {
    font-size: 0.875rem;
  }
}
</style>