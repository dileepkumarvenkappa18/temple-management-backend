<template>
  <div v-if="show" class="flex items-center justify-center" :class="containerClass">
    <!-- Spinner Loader -->
    <div v-if="type === 'spinner'" class="relative">
      <div 
        class="animate-spin rounded-full border-2 border-indigo-200"
        :class="[
          spinnerSize,
          `border-t-${color}-600`
        ]"
      ></div>
    </div>

    <!-- Dots Loader -->
    <div v-else-if="type === 'dots'" class="flex space-x-2">
      <div 
        v-for="i in 3" 
        :key="i"
        class="animate-bounce rounded-full"
        :class="[
          dotSize,
          `bg-${color}-600`,
          { 'animation-delay-150': i === 2, 'animation-delay-300': i === 3 }
        ]"
        :style="{ animationDelay: `${(i - 1) * 150}ms` }"
      ></div>
    </div>

    <!-- Pulse Loader -->
    <div v-else-if="type === 'pulse'" class="flex space-x-1">
      <div 
        v-for="i in 3" 
        :key="i"
        class="animate-pulse rounded"
        :class="[
          pulseSize,
          `bg-${color}-600`
        ]"
        :style="{ animationDelay: `${(i - 1) * 200}ms` }"
      ></div>
    </div>

    <!-- Loading Text -->
    <div v-if="text" class="ml-3">
      <p :class="`text-${textSize} text-${color}-700 font-medium`">
        {{ text }}
      </p>
      <p v-if="subtext" :class="`text-${subtextSize} text-gray-500 mt-1`">
        {{ subtext }}
      </p>
    </div>
  </div>

  <!-- Overlay for full-screen loaders -->
  <div 
    v-if="overlay && show" 
    class="fixed inset-0 bg-white bg-opacity-75 backdrop-blur-sm flex items-center justify-center z-50"
    :class="{ 'bg-gray-900 bg-opacity-50': overlayDark }"
  >
    <div class="bg-white rounded-xl shadow-lg p-8 max-w-sm w-full mx-4">
      <div class="text-center">
        <div class="flex justify-center mb-4">
          <div class="animate-spin rounded-full h-12 w-12 border-2 border-indigo-200 border-t-indigo-600"></div>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ overlayTitle || 'Loading...' }}</h3>
        <p class="text-gray-600">{{ overlayText || 'Please wait while we process your request.' }}</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  // Basic props
  show: {
    type: Boolean,
    default: true
  },
  type: {
    type: String,
    default: 'spinner', // spinner, dots, pulse
    validator: value => ['spinner', 'dots', 'pulse'].includes(value)
  },
  size: {
    type: String,
    default: 'md', // sm, md, lg, xl
    validator: value => ['sm', 'md', 'lg', 'xl'].includes(value)
  },
  color: {
    type: String,
    default: 'indigo'
  },
  
  // Text props
  text: {
    type: String,
    default: ''
  },
  subtext: {
    type: String,
    default: ''
  },
  
  // Layout props
  center: {
    type: Boolean,
    default: false
  },
  fullHeight: {
    type: Boolean,
    default: false
  },
  
  // Overlay props
  overlay: {
    type: Boolean,
    default: false
  },
  overlayDark: {
    type: Boolean,
    default: false
  },
  overlayTitle: {
    type: String,
    default: ''
  },
  overlayText: {
    type: String,
    default: ''
  }
})

// Computed size classes
const spinnerSize = computed(() => {
  const sizes = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-8 w-8',
    xl: 'h-12 w-12'
  }
  return sizes[props.size]
})

const dotSize = computed(() => {
  const sizes = {
    sm: 'h-2 w-2',
    md: 'h-3 w-3',
    lg: 'h-4 w-4',
    xl: 'h-5 w-5'
  }
  return sizes[props.size]
})

const pulseSize = computed(() => {
  const sizes = {
    sm: 'h-2 w-1',
    md: 'h-4 w-1',
    lg: 'h-6 w-1',
    xl: 'h-8 w-2'
  }
  return sizes[props.size]
})

const textSize = computed(() => {
  const sizes = {
    sm: 'sm',
    md: 'base',
    lg: 'lg',
    xl: 'xl'
  }
  return sizes[props.size]
})

const subtextSize = computed(() => {
  const sizes = {
    sm: 'xs',
    md: 'sm',
    lg: 'base',
    xl: 'lg'
  }
  return sizes[props.size]
})

const containerClass = computed(() => {
  let classes = []
  
  if (props.center) classes.push('justify-center')
  if (props.fullHeight) classes.push('min-h-screen')
  
  return classes.join(' ')
})
</script>

<style scoped>
@keyframes bounce {
  0%, 80%, 100% { transform: translateY(0); }
  40% { transform: translateY(-10px); }
}

.animate-bounce {
  animation: bounce 1.5s infinite;
}

.animation-delay-150 {
  animation-delay: 150ms;
}

.animation-delay-300 {
  animation-delay: 300ms;
}
</style>