<template>
  <div 
    :class="cardClasses"
    @click="handleClick"
  >
    <!-- Header Section -->
    <div v-if="hasHeader" class="border-b border-gray-100 pb-4 mb-4">
      <slot name="header">
        <div v-if="title" class="flex items-center justify-between">
          <h3 class="text-lg font-semibold text-gray-900">{{ title }}</h3>
          <slot name="headerActions"></slot>
        </div>
      </slot>
    </div>

    <!-- Main Content -->
    <div :class="contentClass">
      <slot></slot>
    </div>

    <!-- Footer Section -->
    <div v-if="hasFooter" class="border-t border-gray-100 pt-4 mt-4">
      <slot name="footer"></slot>
    </div>
  </div>
</template>

<script setup>
import { computed, useSlots } from 'vue'

const props = defineProps({
  title: {
    type: String,
    default: ''
  },
  elevated: {
    type: Boolean,
    default: false
  },
  hoverable: {
    type: Boolean,
    default: false
  },
  clickable: {
    type: Boolean,
    default: false
  },
  padding: {
    type: String,
    default: 'default',
    validator: (value) => ['none', 'sm', 'default', 'lg'].includes(value)
  },
  contentClass: {
    type: String,
    default: ''
  },
  customClass: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['click'])
const slots = useSlots()

// Simple static computeds to prevent any reactivity issues
const hasHeader = computed(() => !!(slots.header || props.title))
const hasFooter = computed(() => !!slots.footer)

// Pre-built class strings to avoid dynamic computation
const baseClasses = 'bg-white rounded-xl shadow-md border border-gray-200 transition-all duration-200'

const paddingClasses = {
  none: '',
  sm: 'p-4',
  default: 'p-6',
  lg: 'p-8'
}

// Static class computation
const cardClasses = computed(() => {
  let classes = baseClasses
  
  // Add padding
  const paddingClass = paddingClasses[props.padding] || paddingClasses.default
  if (paddingClass) {
    classes += ` ${paddingClass}`
  }
  
  // Add conditional classes
  if (props.hoverable) {
    classes += ' hover:shadow-lg hover:border-indigo-200'
  }
  
  if (props.elevated) {
    classes += ' shadow-lg border-indigo-200'
  }
  
  if (props.clickable) {
    classes += ' cursor-pointer'
  }
  
  // Add custom class
  if (props.customClass) {
    classes += ` ${props.customClass}`
  }
  
  return classes
})

const handleClick = (event) => {
  if (props.clickable) {
    emit('click', event)
  }
}
</script>