<template>
  <div :class="bannerClasses">
    <!-- Background Pattern -->
    <div class="absolute inset-0 bg-gradient-to-r opacity-10">
      <svg class="absolute inset-0 w-full h-full" viewBox="0 0 100 100" preserveAspectRatio="none">
        <defs>
          <pattern id="welcome-pattern" x="0" y="0" width="20" height="20" patternUnits="userSpaceOnUse">
            <circle cx="2" cy="2" r="1" fill="currentColor" opacity="0.5"/>
          </pattern>
        </defs>
        <rect width="100%" height="100%" fill="url(#welcome-pattern)"/>
      </svg>
    </div>

    <div class="relative flex items-center justify-between">
      <!-- Welcome Content -->
      <div class="flex items-center space-x-4">
        <!-- Role Icon -->
        <div :class="iconBgClasses">
          <component :is="roleIcon" :class="iconClasses" />
        </div>

        <!-- Welcome Text -->
        <div>
          <h1 class="text-2xl md:text-3xl font-bold text-white mb-1">
            {{ welcomeMessage }}
          </h1>
          <p class="text-base md:text-lg text-white/90">
            {{ subMessage }}
          </p>
          
          <!-- Additional Info -->
          <div v-if="showAdditionalInfo" class="flex flex-wrap items-center gap-4 mt-3">
            <div v-if="entityName" class="flex items-center text-white/80 text-sm">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
              {{ entityName }}
            </div>
            
            <div v-if="lastLogin" class="flex items-center text-white/80 text-sm">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              Last login: {{ formatDate(lastLogin) }}
            </div>
            
            <div v-if="profileCompletion !== null" class="flex items-center text-white/80 text-sm">
              <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
              </svg>
              Profile {{ profileCompletion }}% complete
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div v-if="showActions" class="hidden md:flex items-center space-x-3">
        <slot name="actions">
          <!-- Default actions based on role -->
          <button
            v-if="userRole === 'tenant'"
            @click="$emit('create-temple')"
            :class="actionButtonClasses"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Create Temple
          </button>
          
          <!-- <button
            v-if="userRole === 'devotee' && profileCompletion < 100"
            @click="$emit('complete-profile')"
            :class="actionButtonClasses"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
            Complete Profile
          </button> -->
          
          <button
            v-if="userRole === 'entity_admin'"
            @click="$emit('quick-action')"
            :class="actionButtonClasses"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z" />
            </svg>
            Quick Actions
          </button>
        </slot>
      </div>
    </div>

    <!-- Mobile Actions -->
    <div v-if="showActions" class="md:hidden mt-4 pt-4 border-t border-white/20">
      <div class="flex flex-wrap gap-2">
        <slot name="mobile-actions">
          <button
            v-if="userRole === 'tenant'"
            @click="$emit('create-temple')"
            :class="mobileActionButtonClasses"
          >
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Create Temple
          </button>
        </slot>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  userRole: {
    type: String,
    required: true,
    validator: (value) => ['tenant', 'entity_admin', 'devotee', 'volunteer', 'super_admin'].includes(value)
  },
  userName: {
    type: String,
    required: true
  },
  entityName: String,
  lastLogin: [String, Date],
  profileCompletion: {
    type: Number,
    default: null
  },
  showActions: {
    type: Boolean,
    default: true
  },
  variant: {
    type: String,
    default: 'gradient',
    validator: (value) => ['gradient', 'solid', 'minimal'].includes(value)
  }
})

const emit = defineEmits(['create-temple', 'complete-profile', 'quick-action'])

const roleConfig = computed(() => {
  const configs = {
    tenant: {
      title: 'Temple Administrator',
      icon: 'BuildingOfficeIcon',
      gradient: 'from-indigo-600 to-purple-600',
      solid: 'bg-indigo-600'
    },
    entity_admin: {
      title: 'Temple Manager',
      icon: 'CogIcon',
      gradient: 'from-blue-600 to-indigo-600',
      solid: 'bg-blue-600'
    },
    devotee: {
      title: 'Devotee',
      icon: 'UserIcon',
      gradient: 'from-green-600 to-blue-600',
      solid: 'bg-green-600'
    },
    volunteer: {
      title: 'Volunteer',
      icon: 'HandRaisedIcon',
      gradient: 'from-orange-600 to-red-600',
      solid: 'bg-orange-600'
    },
    super_admin: {
      title: 'System Administrator',
      icon: 'ShieldCheckIcon',
      gradient: 'from-red-600 to-pink-600',
      solid: 'bg-red-600'
    }
  }
  return configs[props.userRole] || configs.devotee
})

const welcomeMessage = computed(() => {
  const timeOfDay = new Date().getHours()
  let greeting = 'Hello'
  
  if (timeOfDay < 12) greeting = 'Good Morning'
  else if (timeOfDay < 17) greeting = 'Good Afternoon'
  else greeting = 'Good Evening'
  
  return `${greeting}, ${props.userName}!`
})

const subMessage = computed(() => {
  const messages = {
    tenant: 'Manage your temples and oversee operations',
    entity_admin: props.entityName ? `Welcome back to ${props.entityName}` : 'Manage your temple operations',
    devotee: 'Continue your spiritual journey',
    volunteer: 'Thank you for your service',
    super_admin: 'System overview and administration'
  }
  return messages[props.userRole] || 'Welcome to your dashboard'
})

const roleIcon = computed(() => {
  // Using simple div with emoji for now since we don't have Heroicons imported
  // In a real app, you'd import and use the actual icons
  const icons = {
    tenant: 'div',
    entity_admin: 'div', 
    devotee: 'div',
    volunteer: 'div',
    super_admin: 'div'
  }
  return icons[props.userRole]
})

const showAdditionalInfo = computed(() => {
  return props.entityName || props.lastLogin || props.profileCompletion !== null
})

const bannerClasses = computed(() => {
  const base = 'relative p-6 md:p-8 rounded-xl overflow-hidden'
  
  if (props.variant === 'gradient') {
    return `${base} bg-gradient-to-r ${roleConfig.value.gradient}`
  } else if (props.variant === 'solid') {
    return `${base} ${roleConfig.value.solid}`
  } else {
    return `${base} bg-white border-2 border-gray-200`
  }
})

const iconBgClasses = computed(() => {
  if (props.variant === 'minimal') {
    return 'flex items-center justify-center w-16 h-16 bg-indigo-100 rounded-xl'
  }
  return 'flex items-center justify-center w-16 h-16 bg-white/20 backdrop-blur-sm rounded-xl'
})

const iconClasses = computed(() => {
  if (props.variant === 'minimal') {
    return 'w-8 h-8 text-indigo-600'
  }
  return 'w-8 h-8 text-white'
})

const actionButtonClasses = computed(() => {
  if (props.variant === 'minimal') {
    return 'inline-flex items-center px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors duration-200'
  }
  return 'inline-flex items-center px-4 py-2 bg-white/20 backdrop-blur-sm text-white text-sm font-medium rounded-lg hover:bg-white/30 transition-all duration-200'
})

const mobileActionButtonClasses = computed(() => {
  if (props.variant === 'minimal') {
    return 'inline-flex items-center px-3 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 transition-colors duration-200'
  }
  return 'inline-flex items-center px-3 py-2 bg-white/20 backdrop-blur-sm text-white text-sm font-medium rounded-lg hover:bg-white/30 transition-all duration-200'
})

const formatDate = (date) => {
  if (!date) return ''
  const d = new Date(date)
  return d.toLocaleDateString('en-US', { 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}
</script>