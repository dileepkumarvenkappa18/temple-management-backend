<template>
  <div 
    :class="[
      'bg-white rounded-xl shadow-md border border-gray-100 hover:shadow-lg transition-all duration-200',
      paddingClass,
      heightClass
    ]"
  >
    <!-- Widget Header -->
    <div v-if="title || $slots.header" class="flex items-center justify-between mb-6">
      <div class="flex items-center">
        <div 
          v-if="icon"
          :class="[
            'flex items-center justify-center w-10 h-10 rounded-lg mr-3',
            iconBg
          ]"
        >
          <!-- Profile Icon -->
          <svg v-if="icon === 'profile' || icon === 'user'" :class="['w-5 h-5', iconColor]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2"/>
            <circle cx="12" cy="7" r="4"/>
          </svg>
          
          <!-- Bookings/Calendar Icon -->
          <svg v-else-if="icon === 'bookings' || icon === 'calendar'" :class="['w-5 h-5', iconColor]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="3" y="4" width="18" height="18" rx="2" ry="2"/>
            <line x1="16" y1="2" x2="16" y2="6"/>
            <line x1="8" y1="2" x2="8" y2="6"/>
            <line x1="3" y1="10" x2="21" y2="10"/>
            <path d="M8 14h.01"/>
            <path d="M12 14h.01"/>
            <path d="M16 14h.01"/>
            <path d="M8 18h.01"/>
            <path d="M12 18h.01"/>
            <path d="M16 18h.01"/>
          </svg>
          
          <!-- Donations/Heart Icon -->
          <svg v-else-if="icon === 'donations' || icon === 'heart'" :class="['w-5 h-5', iconColor]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M20.84 4.61a5.5 5.5 0 0 0-7.78 0L12 5.67l-1.06-1.06a5.5 5.5 0 0 0-7.78 7.78l1.06 1.06L12 21.23l7.78-7.78 1.06-1.06a5.5 5.5 0 0 0 0-7.78z"/>
          </svg>
          
          <!-- Events/Lightning Icon -->
          <svg v-else-if="icon === 'events' || icon === 'lightning'" :class="['w-5 h-5', iconColor]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <path d="M13 2L3 14h9l-1 8 10-12h-9l1-8z"/>
          </svg>
          
          <!-- Default Icon (if no match) -->
          <svg v-else :class="['w-5 h-5', iconColor]" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="12" r="10"/>
            <path d="M12 6v6l4 2"/>
          </svg>
        </div>
        <div>
          <h3 v-if="title" class="text-lg font-semibold text-gray-900">{{ title }}</h3>
          <p v-if="subtitle" class="text-sm text-gray-500 mt-1">{{ subtitle }}</p>
        </div>
      </div>
      
      <!-- Header Actions -->
      <div v-if="$slots.header" class="flex items-center space-x-2">
        <slot name="header" />
      </div>
      
      <!-- Default Actions -->
      <div v-else-if="showActions" class="flex items-center space-x-2">
        <button
          v-if="refreshable"
          @click="$emit('refresh')"
          class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
          :disabled="loading"
        >
          <svg 
            :class="['w-4 h-4', loading ? 'animate-spin' : '']"
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24"
          >
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
        </button>
        
        <button
          v-if="expandable"
          @click="$emit('expand')"
          class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 8V4m0 0h4M4 4l5 5m11-1V4m0 0h-4m4 0l-5 5M4 16v4m0 0h4m-4 0l5-5m11 5l-5-5m5 5v-4m0 4h-4" />
          </svg>
        </button>
        
        <div v-if="hasMenu" class="relative">
          <button
            @click="showMenu = !showMenu"
            class="p-2 text-gray-400 hover:text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01M12 6a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2zm0 7a1 1 0 110-2 1 1 0 010 2z" />
            </svg>
          </button>
          
          <!-- Menu Dropdown -->
          <div 
            v-if="showMenu"
            v-click-outside="() => showMenu = false"
            class="absolute right-0 mt-2 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-2 z-50"
          >
            <slot name="menu" />
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="flex items-center space-x-3">
        <svg class="animate-spin w-6 h-6 text-indigo-600" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <span class="text-gray-500">{{ loadingText }}</span>
      </div>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="flex items-center justify-center py-12">
      <div class="text-center">
        <svg class="w-12 h-12 text-red-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <p class="text-gray-500 mb-4">{{ error }}</p>
        <button
          @click="$emit('retry')"
          class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
        >
          Try Again
        </button>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="empty" class="flex items-center justify-center py-12">
      <div class="text-center">
        <svg class="w-12 h-12 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
        </svg>
        <p class="text-gray-500 mb-2">{{ emptyText }}</p>
        <p class="text-sm text-gray-400">{{ emptySubtext }}</p>
      </div>
    </div>

    <!-- Widget Content -->
    <div v-else>
      <slot />
    </div>

    <!-- Widget Footer -->
    <div v-if="$slots.footer" class="mt-6 pt-4 border-t border-gray-100">
      <slot name="footer" />
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  title: String,
  subtitle: String,
  icon: [String, Object],
  iconBg: {
    type: String,
    default: 'bg-indigo-100'
  },
  iconColor: {
    type: String,
    default: 'text-indigo-600'
  },
  loading: {
    type: Boolean,
    default: false
  },
  loadingText: {
    type: String,
    default: 'Loading...'
  },
  error: String,
  empty: {
    type: Boolean,
    default: false
  },
  emptyText: {
    type: String,
    default: 'No data available'
  },
  emptySubtext: {
    type: String,
    default: 'Data will appear here when available'
  },
  refreshable: {
    type: Boolean,
    default: false
  },
  expandable: {
    type: Boolean,
    default: false
  },
  hasMenu: {
    type: Boolean,
    default: false
  },
  padding: {
    type: String,
    default: 'normal',
    validator: (value) => ['none', 'sm', 'normal', 'lg'].includes(value)
  },
  height: {
    type: String,
    default: 'auto',
    validator: (value) => ['auto', 'sm', 'md', 'lg', 'full'].includes(value)
  }
})

const emit = defineEmits(['refresh', 'expand', 'retry'])

const showMenu = ref(false)

const paddingClass = computed(() => {
  const paddingMap = {
    none: 'p-0',
    sm: 'p-4',
    normal: 'p-6',
    lg: 'p-8'
  }
  return paddingMap[props.padding]
})

const heightClass = computed(() => {
  const heightMap = {
    auto: '',
    sm: 'h-64',
    md: 'h-80',
    lg: 'h-96',
    full: 'h-full'
  }
  return heightMap[props.height]
})

const showActions = computed(() => {
  return props.refreshable || props.expandable || props.hasMenu
})

// Click outside directive
const vClickOutside = {
  mounted(el, binding) {
    el.clickOutsideEvent = (event) => {
      if (!(el === event.target || el.contains(event.target))) {
        binding.value()
      }
    }
    document.addEventListener('click', el.clickOutsideEvent)
  },
  unmounted(el) {
    document.removeEventListener('click', el.clickOutsideEvent)
  }
}
</script>