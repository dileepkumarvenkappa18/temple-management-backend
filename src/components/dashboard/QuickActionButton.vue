<template>
  <button
    @click="$emit('click')"
    class="group relative overflow-hidden rounded-xl border border-gray-200 bg-white p-4 hover:border-gray-300 hover:shadow-lg transition-all duration-300 focus:outline-none focus:ring-2 focus:ring-offset-2 text-left"
    :class="focusRingClass"
  >
    <!-- Background gradient on hover -->
    <div 
      class="absolute inset-0 opacity-0 group-hover:opacity-5 transition-opacity duration-300"
      :class="backgroundClass"
    ></div>
    
    <!-- Badge (if provided) -->
    <div v-if="badge" class="absolute -top-1 -right-1 z-10">
      <span 
        class="inline-flex items-center justify-center px-2 py-1 text-xs font-bold text-white rounded-full"
        :class="badgeClass"
      >
        {{ badge }}
      </span>
    </div>

    <!-- Content -->
    <div class="relative z-10">
      <!-- Icon -->
      <div class="flex items-center justify-center w-10 h-10 rounded-lg mb-3" :class="iconBackgroundClass">
        <svg class="w-5 h-5" :class="iconClass" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="iconPath"></path>
        </svg>
      </div>

      <!-- Text -->
      <div>
        <h4 class="font-semibold text-gray-900 text-sm mb-1 group-hover:text-gray-800 transition-colors">
          {{ title }}
        </h4>
        <p class="text-xs text-gray-600 group-hover:text-gray-700 transition-colors">
          {{ description }}
        </p>
      </div>
    </div>

    <!-- Hover arrow -->
    <div class="absolute top-4 right-4 opacity-0 group-hover:opacity-100 transition-all duration-300 transform translate-x-1 group-hover:translate-x-0">
      <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
      </svg>
    </div>
  </button>
</template>

<script setup>
import { computed } from 'vue'

// Props
const props = defineProps({
  title: {
    type: String,
    required: true
  },
  description: {
    type: String,
    required: true
  },
  icon: {
    type: String,
    required: true
  },
  color: {
    type: String,
    default: 'indigo',
    validator: (value) => [
      'indigo', 'blue', 'purple', 'green', 'orange', 'red', 
      'pink', 'cyan', 'emerald', 'gray'
    ].includes(value)
  },
  badge: {
    type: String,
    default: null
  }
})

// Emits
defineEmits(['click'])

// Color classes mapping
const colorClasses = {
  indigo: {
    icon: 'text-indigo-600',
    iconBg: 'bg-indigo-100',
    background: 'bg-indigo-500',
    focus: 'focus:ring-indigo-500',
    badge: 'bg-indigo-500'
  },
  blue: {
    icon: 'text-blue-600',
    iconBg: 'bg-blue-100',
    background: 'bg-blue-500',
    focus: 'focus:ring-blue-500',
    badge: 'bg-blue-500'
  },
  purple: {
    icon: 'text-purple-600',
    iconBg: 'bg-purple-100',
    background: 'bg-purple-500',
    focus: 'focus:ring-purple-500',
    badge: 'bg-purple-500'
  },
  green: {
    icon: 'text-green-600',
    iconBg: 'bg-green-100',
    background: 'bg-green-500',
    focus: 'focus:ring-green-500',
    badge: 'bg-green-500'
  },
  orange: {
    icon: 'text-orange-600',
    iconBg: 'bg-orange-100',
    background: 'bg-orange-500',
    focus: 'focus:ring-orange-500',
    badge: 'bg-orange-500'
  },
  red: {
    icon: 'text-red-600',
    iconBg: 'bg-red-100',
    background: 'bg-red-500',
    focus: 'focus:ring-red-500',
    badge: 'bg-red-500'
  },
  pink: {
    icon: 'text-pink-600',
    iconBg: 'bg-pink-100',
    background: 'bg-pink-500',
    focus: 'focus:ring-pink-500',
    badge: 'bg-pink-500'
  },
  cyan: {
    icon: 'text-cyan-600',
    iconBg: 'bg-cyan-100',
    background: 'bg-cyan-500',
    focus: 'focus:ring-cyan-500',
    badge: 'bg-cyan-500'
  },
  emerald: {
    icon: 'text-emerald-600',
    iconBg: 'bg-emerald-100',
    background: 'bg-emerald-500',
    focus: 'focus:ring-emerald-500',
    badge: 'bg-emerald-500'
  },
  gray: {
    icon: 'text-gray-600',
    iconBg: 'bg-gray-100',
    background: 'bg-gray-500',
    focus: 'focus:ring-gray-500',
    badge: 'bg-gray-500'
  }
}

// Icon paths mapping
const iconPaths = {
  // Temple & Management Icons
  temple: 'M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4',
  users: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z',
  user: 'M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z',
  userAdd: 'M18 9v3m0 0v3m0-3h3m-3 0h-3m-2-5a4 4 0 11-8 0 4 4 0 018 0zM3 20a6 6 0 0112 0v1H3v-1z',
  
  // Financial Icons
  currency: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  wallet: 'M3 6l3 1m0 0l-3 9a5.002 5.002 0 006.001 0M6 7l3 9M6 7l6-2m6 2l3-1m-3 1l-3 9a5.002 5.002 0 006.001 0M18 7l3 9m-3-9l-6-2m0-2v2m0 16V5m0 16l-3-9',
  chart: 'M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z',
  receipt: 'M9 12h6m-6 4h6m2 5l-8-5H3v-2h6V5h2l8 5-2 2h-6v5z',
  
  // Events & Calendar Icons
  calendar: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
  calendarAdd: 'M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v4m0-4H4m2 0V9a2 2 0 012-2h8a2 2 0 012 2v3m0 0a2 2 0 104 0m-4 0a2 2 0 014 0m0 0v4a2 2 0 01-2 2h-3',
  star: 'M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z',
  
  // Communication Icons
  bell: 'M15 17h5l-1.405-1.405A2.032 2.032 0 0118 14.158V11a6.002 6.002 0 00-4-5.659V5a2 2 0 10-4 0v.341C7.67 6.165 6 8.388 6 11v3.159c0 .538-.214 1.055-.595 1.436L4 17h5m6 0v1a3 3 0 11-6 0v-1m6 0H9',
  mail: 'M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z',
  phone: 'M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z',
  chatAlt: 'M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z',
  
  // Management Icons
  cog: 'M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z M15 12a3 3 0 11-6 0 3 3 0 016 0z',
  clipboard: 'M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h.01M9 16h.01',
  document: 'M9 12h6m-6 4h6m2 5l-8-8H5a2 2 0 01-2-2V6a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1M9 12h6m-6 4h6m2 5l-8-8H5a2 2 0 01-2-2V6a2 2 0 012-2h4l2 2h4a2 2 0 012 2v1',
  folder: 'M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z',
  
  // Action Icons
  plus: 'M12 6v6m0 0v6m0-6h6m-6 0H6',
  edit: 'M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z',
  eye: 'M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z',
  download: 'M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4',
  upload: 'M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12',
  
  // Status Icons
  check: 'M5 13l4 4L19 7',
  x: 'M6 18L18 6M6 6l12 12',
  exclamation: 'M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z',
  
  // Special Temple Actions
  puja: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
  heart: 'M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z',
  gift: 'M12 8v13m0-13V6a2 2 0 112 2h-2zm0 0V5.5A2.5 2.5 0 109.5 8H12zm-7 4h14M5 12a2 2 0 110-4h14a2 2 0 110 4M5 12v7a2 2 0 002 2h10a2 2 0 002-2v-7'
}

// Computed properties for dynamic classes
const iconPath = computed(() => iconPaths[props.icon] || iconPaths.document)
const iconClass = computed(() => colorClasses[props.color]?.icon || colorClasses.indigo.icon)
const iconBackgroundClass = computed(() => colorClasses[props.color]?.iconBg || colorClasses.indigo.iconBg)
const backgroundClass = computed(() => colorClasses[props.color]?.background || colorClasses.indigo.background)
const focusRingClass = computed(() => colorClasses[props.color]?.focus || colorClasses.indigo.focus)
const badgeClass = computed(() => colorClasses[props.color]?.badge || colorClasses.indigo.badge)
</script>