<!-- src/components/temple/TempleApprovalStatus.vue -->
<template>
  <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full text-sm font-medium transition-all duration-200"
       :class="statusConfig.classes">
    <!-- Status Icon -->
    <component :is="statusConfig.icon" class="w-4 h-4" />
    
    <!-- Status Text -->
    <span>{{ statusConfig.text }}</span>
    
    <!-- Animated Pulse for Pending -->
    <div v-if="normalizedStatus === 'PENDING'" class="w-2 h-2 bg-current rounded-full animate-pulse opacity-60"></div>
  </div>
</template>

<script setup>
import { computed, h } from 'vue'

// Props
const props = defineProps({
  status: {
    type: String,
    required: true,
    validator: (value) => {
      // Case-insensitive validation
      const normalizedValue = String(value || '').toUpperCase();
      return ['PENDING', 'APPROVED', 'REJECTED'].includes(normalizedValue);
    }
  },
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  }
})

// Normalized status for consistent comparison
const normalizedStatus = computed(() => {
  return String(props.status || '').toUpperCase();
})

// Status configuration
const statusConfig = computed(() => {
  const configs = {
    'PENDING': {
      text: 'Pending Approval',
      classes: 'bg-amber-100 text-amber-800 border border-amber-200',
      icon: () => h('svg', {
        fill: 'none',
        stroke: 'currentColor',
        viewBox: '0 0 24 24'
      }, [
        h('path', {
          'stroke-linecap': 'round',
          'stroke-linejoin': 'round',
          'stroke-width': '2',
          d: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z'
        })
      ])
    },
    'APPROVED': {
      text: 'Approved',
      classes: 'bg-green-100 text-green-800 border border-green-200',
      icon: () => h('svg', {
        fill: 'none',
        stroke: 'currentColor',
        viewBox: '0 0 24 24'
      }, [
        h('path', {
          'stroke-linecap': 'round',
          'stroke-linejoin': 'round',
          'stroke-width': '2',
          d: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
        })
      ])
    },
    'REJECTED': {
      text: 'Rejected',
      classes: 'bg-red-100 text-red-800 border border-red-200',
      icon: () => h('svg', {
        fill: 'none',
        stroke: 'currentColor',
        viewBox: '0 0 24 24'
      }, [
        h('path', {
          'stroke-linecap': 'round',
          'stroke-linejoin': 'round',
          'stroke-width': '2',
          d: 'M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z'
        })
      ])
    }
  }
  
  // Use normalizedStatus for lookup instead of raw props.status
  return configs[normalizedStatus.value] || configs['PENDING']
})
</script>