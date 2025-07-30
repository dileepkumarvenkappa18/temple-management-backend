<!-- src/components/common/BaseButton.vue -->
<template>
  <button
    :type="type"
    :disabled="disabled || loading"
    :class="buttonClasses"
    @click="handleClick"
  >
    <!-- Loading spinner -->
    <svg
      v-if="loading"
      class="animate-spin -ml-1 mr-2 h-4 w-4"
      xmlns="http://www.w3.org/2000/svg"
      fill="none"
      viewBox="0 0 24 24"
    >
      <circle
        class="opacity-25"
        cx="12"
        cy="12"
        r="10"
        stroke="currentColor"
        stroke-width="4"
      ></circle>
      <path
        class="opacity-75"
        fill="currentColor"
        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
      ></path>
    </svg>

    <!-- Icon (left) -->
    <component
      v-if="iconLeft && !loading"
      :is="iconLeft"
      :class="iconSizeClass"
    />

    <!-- Button text -->
    <span v-if="$slots.default || text">
      <slot>{{ text }}</slot>
    </span>

    <!-- Icon (right) -->
    <component
      v-if="iconRight && !loading"
      :is="iconRight"
      :class="iconSizeClass"
    />
  </button>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Button types
  variant: {
    type: String,
    default: 'primary',
    validator: (value) => ['primary', 'secondary', 'outline', 'ghost', 'danger', 'success', 'warning'].includes(value)
  },
  
  // Button sizes
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['xs', 'sm', 'md', 'lg', 'xl'].includes(value)
  },
  
  // HTML button type
  type: {
    type: String,
    default: 'button',
    validator: (value) => ['button', 'submit', 'reset'].includes(value)
  },
  
  // Button text
  text: {
    type: String,
    default: ''
  },
  
  // States
  disabled: {
    type: Boolean,
    default: false
  },
  
  loading: {
    type: Boolean,
    default: false
  },
  
  // Icons
  iconLeft: {
    type: [String, Object],
    default: null
  },
  
  iconRight: {
    type: [String, Object],
    default: null
  },
  
  // Full width
  fullWidth: {
    type: Boolean,
    default: false
  },
  
  // Rounded
  rounded: {
    type: String,
    default: 'md',
    validator: (value) => ['none', 'sm', 'md', 'lg', 'xl', 'full'].includes(value)
  }
})

const emit = defineEmits(['click'])

// Button variant classes
const variantClasses = computed(() => {
  const variants = {
    primary: 'bg-indigo-600 text-white hover:bg-indigo-700 focus:ring-indigo-500 border-transparent',
    secondary: 'bg-indigo-100 text-indigo-700 hover:bg-indigo-200 focus:ring-indigo-500 border-transparent',
    outline: 'bg-transparent text-indigo-600 border-indigo-300 hover:bg-indigo-50 focus:ring-indigo-500',
    ghost: 'bg-transparent text-gray-700 hover:bg-gray-100 focus:ring-gray-500 border-transparent',
    danger: 'bg-red-600 text-white hover:bg-red-700 focus:ring-red-500 border-transparent',
    success: 'bg-green-600 text-white hover:bg-green-700 focus:ring-green-500 border-transparent',
    warning: 'bg-yellow-500 text-white hover:bg-yellow-600 focus:ring-yellow-400 border-transparent'
  }
  return variants[props.variant] || variants.primary
})

// Button size classes
const sizeClasses = computed(() => {
  const sizes = {
    xs: 'px-2.5 py-1.5 text-xs',
    sm: 'px-3 py-2 text-sm',
    md: 'px-4 py-2.5 text-sm',
    lg: 'px-6 py-3 text-base',
    xl: 'px-8 py-4 text-lg'
  }
  return sizes[props.size] || sizes.md
})

// Icon size classes
const iconSizeClass = computed(() => {
  const iconSizes = {
    xs: 'w-3 h-3',
    sm: 'w-4 h-4',
    md: 'w-4 h-4',
    lg: 'w-5 h-5',
    xl: 'w-6 h-6'
  }
  return iconSizes[props.size] || iconSizes.md
})

// Rounded classes
const roundedClasses = computed(() => {
  const rounded = {
    none: 'rounded-none',
    sm: 'rounded-sm',
    md: 'rounded-md',
    lg: 'rounded-lg',
    xl: 'rounded-xl',
    full: 'rounded-full'
  }
  return rounded[props.rounded] || rounded.md
})

// Combined button classes
const buttonClasses = computed(() => {
  const baseClasses = [
    'inline-flex items-center justify-center',
    'font-medium',
    'border',
    'transition-all duration-200',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    roundedClasses.value,
    variantClasses.value,
    sizeClasses.value
  ]

  // Add full width class
  if (props.fullWidth) {
    baseClasses.push('w-full')
  }

  // Add disabled styles
  if (props.disabled || props.loading) {
    baseClasses.push('opacity-50 cursor-not-allowed')
  }

  // Add spacing for icons
  if (props.iconLeft && (props.text || props.$slots.default)) {
    baseClasses.push('space-x-2')
  }
  if (props.iconRight && (props.text || props.$slots.default)) {
    baseClasses.push('space-x-2')
  }

  return baseClasses.join(' ')
})

const handleClick = (event) => {
  if (!props.disabled && !props.loading) {
    emit('click', event)
  }
}
</script>