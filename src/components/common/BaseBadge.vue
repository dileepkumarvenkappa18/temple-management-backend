<template>
  <span
    :class="badgeClasses"
    class="inline-flex items-center font-medium transition-all duration-200"
  >
    <!-- Icon slot for leading icons -->
    <slot name="icon" v-if="$slots.icon">
      <span :class="iconClasses">
        <slot name="icon"></slot>
      </span>
    </slot>
    
    <!-- Badge text content -->
    <span class="truncate">{{ text }}</span>
    
    <!-- Dot indicator for status badges -->
    <span
      v-if="showDot"
      :class="dotClasses"
      class="ml-1.5 h-1.5 w-1.5 rounded-full"
    ></span>
    
    <!-- Close/dismiss button -->
    <button
      v-if="dismissible"
      @click="$emit('dismiss')"
      :class="closeButtonClasses"
      class="ml-1.5 inline-flex h-4 w-4 items-center justify-center rounded-full hover:bg-black/10 focus:outline-none focus:ring-2 focus:ring-offset-1 transition-colors duration-150"
      aria-label="Remove badge"
    >
      <svg class="h-2.5 w-2.5" fill="currentColor" viewBox="0 0 20 20">
        <path
          fill-rule="evenodd"
          d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
          clip-rule="evenodd"
        />
      </svg>
    </button>
  </span>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Badge text content
  text: {
    type: String,
    required: true
  },
  
  // Badge variant/type
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => [
      'primary', 'secondary', 'success', 'error', 'warning', 'info',
      'pending', 'approved', 'rejected', 'active', 'inactive', 'neutral'
    ].includes(value)
  },
  
  // Badge size
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
  },
  
  // Badge style
  style: {
    type: String,
    default: 'filled',
    validator: (value) => ['filled', 'outlined', 'subtle', 'minimal'].includes(value)
  },
  
  // Show status dot
  showDot: {
    type: Boolean,
    default: false
  },
  
  // Make badge dismissible
  dismissible: {
    type: Boolean,
    default: false
  },
  
  // Make badge rounded (pill style)
  rounded: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['dismiss'])

// Size classes mapping
const sizeClasses = computed(() => {
  const sizes = {
    xs: 'px-2 py-0.5 text-xs',
    sm: 'px-2.5 py-0.5 text-sm',
    md: 'px-3 py-1 text-sm',
    lg: 'px-3.5 py-1.5 text-base',
    xl: 'px-4 py-2 text-lg'
  }
  return sizes[props.size]
})

// Variant and style classes mapping
const variantClasses = computed(() => {
  const variants = {
    primary: {
      filled: 'bg-indigo-600 text-white shadow-sm',
      outlined: 'border-2 border-indigo-600 text-indigo-600 bg-white',
      subtle: 'bg-indigo-50 text-indigo-700 border border-indigo-200',
      minimal: 'text-indigo-600 bg-transparent'
    },
    secondary: {
      filled: 'bg-indigo-500 text-white shadow-sm',
      outlined: 'border-2 border-indigo-500 text-indigo-500 bg-white',
      subtle: 'bg-indigo-50 text-indigo-600 border border-indigo-200',
      minimal: 'text-indigo-500 bg-transparent'
    },
    success: {
      filled: 'bg-green-600 text-white shadow-sm',
      outlined: 'border-2 border-green-600 text-green-600 bg-white',
      subtle: 'bg-green-50 text-green-700 border border-green-200',
      minimal: 'text-green-600 bg-transparent'
    },
    error: {
      filled: 'bg-red-600 text-white shadow-sm',
      outlined: 'border-2 border-red-600 text-red-600 bg-white',
      subtle: 'bg-red-50 text-red-700 border border-red-200',
      minimal: 'text-red-600 bg-transparent'
    },
    warning: {
      filled: 'bg-amber-500 text-white shadow-sm',
      outlined: 'border-2 border-amber-500 text-amber-600 bg-white',
      subtle: 'bg-amber-50 text-amber-700 border border-amber-200',
      minimal: 'text-amber-600 bg-transparent'
    },
    info: {
      filled: 'bg-blue-600 text-white shadow-sm',
      outlined: 'border-2 border-blue-600 text-blue-600 bg-white',
      subtle: 'bg-blue-50 text-blue-700 border border-blue-200',
      minimal: 'text-blue-600 bg-transparent'
    },
    pending: {
      filled: 'bg-yellow-500 text-white shadow-sm',
      outlined: 'border-2 border-yellow-500 text-yellow-600 bg-white',
      subtle: 'bg-yellow-50 text-yellow-700 border border-yellow-200',
      minimal: 'text-yellow-600 bg-transparent'
    },
    approved: {
      filled: 'bg-green-600 text-white shadow-sm',
      outlined: 'border-2 border-green-600 text-green-600 bg-white',
      subtle: 'bg-green-50 text-green-700 border border-green-200',
      minimal: 'text-green-600 bg-transparent'
    },
    rejected: {
      filled: 'bg-red-600 text-white shadow-sm',
      outlined: 'border-2 border-red-600 text-red-600 bg-white',
      subtle: 'bg-red-50 text-red-700 border border-red-200',
      minimal: 'text-red-600 bg-transparent'
    },
    active: {
      filled: 'bg-green-600 text-white shadow-sm',
      outlined: 'border-2 border-green-600 text-green-600 bg-white',
      subtle: 'bg-green-50 text-green-700 border border-green-200',
      minimal: 'text-green-600 bg-transparent'
    },
    inactive: {
      filled: 'bg-gray-500 text-white shadow-sm',
      outlined: 'border-2 border-gray-500 text-gray-600 bg-white',
      subtle: 'bg-gray-50 text-gray-700 border border-gray-200',
      minimal: 'text-gray-600 bg-transparent'
    },
    neutral: {
      filled: 'bg-gray-600 text-white shadow-sm',
      outlined: 'border-2 border-gray-600 text-gray-600 bg-white',
      subtle: 'bg-gray-50 text-gray-700 border border-gray-200',
      minimal: 'text-gray-600 bg-transparent'
    }
  }
  
  return variants[props.variant]?.[props.style] || variants.primary.filled
})

// Rounded classes
const roundedClasses = computed(() => {
  if (!props.rounded) return 'rounded-md'
  
  const roundedSizes = {
    xs: 'rounded-full',
    sm: 'rounded-full',
    md: 'rounded-full',
    lg: 'rounded-2xl',
    xl: 'rounded-2xl'
  }
  
  return roundedSizes[props.size]
})

// Icon classes for proper spacing
const iconClasses = computed(() => {
  const iconSizes = {
    xs: 'mr-1 h-3 w-3',
    sm: 'mr-1.5 h-3.5 w-3.5',
    md: 'mr-1.5 h-4 w-4',
    lg: 'mr-2 h-5 w-5',
    xl: 'mr-2 h-6 w-6'
  }
  
  return iconSizes[props.size]
})

// Dot indicator classes
const dotClasses = computed(() => {
  const dotColors = {
    primary: 'bg-indigo-400',
    secondary: 'bg-indigo-300',
    success: 'bg-green-400',
    error: 'bg-red-400',
    warning: 'bg-amber-400',
    info: 'bg-blue-400',
    pending: 'bg-yellow-400',
    approved: 'bg-green-400',
    rejected: 'bg-red-400',
    active: 'bg-green-400',
    inactive: 'bg-gray-400',
    neutral: 'bg-gray-400'
  }
  
  return dotColors[props.variant] || dotColors.primary
})

// Close button classes
const closeButtonClasses = computed(() => {
  const closeColors = {
    primary: 'text-indigo-200 hover:text-indigo-100 focus:ring-indigo-500',
    secondary: 'text-indigo-300 hover:text-indigo-200 focus:ring-indigo-400',
    success: 'text-green-200 hover:text-green-100 focus:ring-green-500',
    error: 'text-red-200 hover:text-red-100 focus:ring-red-500',
    warning: 'text-amber-200 hover:text-amber-100 focus:ring-amber-500',
    info: 'text-blue-200 hover:text-blue-100 focus:ring-blue-500',
    pending: 'text-yellow-200 hover:text-yellow-100 focus:ring-yellow-500',
    approved: 'text-green-200 hover:text-green-100 focus:ring-green-500',
    rejected: 'text-red-200 hover:text-red-100 focus:ring-red-500',
    active: 'text-green-200 hover:text-green-100 focus:ring-green-500',
    inactive: 'text-gray-300 hover:text-gray-200 focus:ring-gray-500',
    neutral: 'text-gray-300 hover:text-gray-200 focus:ring-gray-500'
  }
  
  // For outlined and subtle styles, use darker colors
  if (props.style === 'outlined' || props.style === 'subtle') {
    const outlinedColors = {
      primary: 'text-indigo-500 hover:text-indigo-600 focus:ring-indigo-500',
      secondary: 'text-indigo-400 hover:text-indigo-500 focus:ring-indigo-400',
      success: 'text-green-500 hover:text-green-600 focus:ring-green-500',
      error: 'text-red-500 hover:text-red-600 focus:ring-red-500',
      warning: 'text-amber-500 hover:text-amber-600 focus:ring-amber-500',
      info: 'text-blue-500 hover:text-blue-600 focus:ring-blue-500'
    }
    return outlinedColors[props.variant] || outlinedColors.primary
  }
  
  return closeColors[props.variant] || closeColors.primary
})

// Final badge classes
const badgeClasses = computed(() => {
  return [
    sizeClasses.value,
    variantClasses.value,
    roundedClasses.value,
    'select-none',
    'max-w-full',
    // Hover effects for interactive badges
    props.dismissible ? 'hover:shadow-md' : ''
  ].filter(Boolean).join(' ')
})
</script>