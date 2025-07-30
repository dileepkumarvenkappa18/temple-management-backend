<template>
  <div class="space-y-4">
    <label class="block text-sm font-medium text-gray-700 mb-2">
      Select Your Role
    </label>
    
    <div class="grid gap-4 md:grid-cols-3">
      <!-- Tenant Role -->
      <div
        :class="[
          'relative border-2 rounded-xl p-6 cursor-pointer transition-all duration-200',
          selectedRole === 'tenant' 
            ? 'border-indigo-600 bg-indigo-50 shadow-md' 
            : 'border-gray-200 hover:border-indigo-300 hover:bg-gray-50'
        ]"
        @click="selectRole('tenant')"
      >
        <div class="flex flex-col items-center text-center space-y-3">
          <!-- Icon -->
          <div :class="[
            'w-12 h-12 rounded-full flex items-center justify-center transition-colors duration-200',
            selectedRole === 'tenant' ? 'bg-indigo-600' : 'bg-gray-400'
          ]">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-2 8h2V9h-2v8z" />
            </svg>
          </div>
          
          <!-- Title -->
          <h3 :class="[
            'font-semibold text-lg transition-colors duration-200',
            selectedRole === 'tenant' ? 'text-indigo-600' : 'text-gray-900'
          ]">
            Temple Admin
          </h3>
          
          <!-- Description -->
          <p class="text-sm text-gray-600 leading-relaxed">
            Manage temple operations, devotees, sevas, donations, and events
          </p>
          
          <!-- Features List -->
          <ul class="text-xs text-gray-500 space-y-1 mt-2">
            <li>• Create & manage temples</li>
            <li>• Approve seva bookings</li>
            <li>• Track donations</li>
            <li>• Organize events</li>
          </ul>
        </div>
        
        <!-- Selection Indicator -->
        <div v-if="selectedRole === 'tenant'" 
             class="absolute top-3 right-3 w-6 h-6 bg-indigo-600 rounded-full flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
      </div>

      <!-- Devotee Role -->
      <div
        :class="[
          'relative border-2 rounded-xl p-6 cursor-pointer transition-all duration-200',
          selectedRole === 'devotee' 
            ? 'border-indigo-600 bg-indigo-50 shadow-md' 
            : 'border-gray-200 hover:border-indigo-300 hover:bg-gray-50'
        ]"
        @click="selectRole('devotee')"
      >
        <div class="flex flex-col items-center text-center space-y-3">
          <!-- Icon -->
          <div :class="[
            'w-12 h-12 rounded-full flex items-center justify-center transition-colors duration-200',
            selectedRole === 'devotee' ? 'bg-indigo-600' : 'bg-gray-400'
          ]">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
            </svg>
          </div>
          
          <!-- Title -->
          <h3 :class="[
            'font-semibold text-lg transition-colors duration-200',
            selectedRole === 'devotee' ? 'text-indigo-600' : 'text-gray-900'
          ]">
            Devotee
          </h3>
          
          <!-- Description -->
          <p class="text-sm text-gray-600 leading-relaxed">
            Join temples, book sevas, make donations, and participate in events
          </p>
          
          <!-- Features List -->
          <ul class="text-xs text-gray-500 space-y-1 mt-2">
            <li>• Join multiple temples</li>
            <li>• Book seva slots</li>
            <li>• Make donations</li>
            <li>• View events & festivals</li>
          </ul>
        </div>
        
        <!-- Selection Indicator -->
        <div v-if="selectedRole === 'devotee'" 
             class="absolute top-3 right-3 w-6 h-6 bg-indigo-600 rounded-full flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
      </div>

      <!-- Volunteer Role -->
      <div
        :class="[
          'relative border-2 rounded-xl p-6 cursor-pointer transition-all duration-200',
          selectedRole === 'volunteer' 
            ? 'border-indigo-600 bg-indigo-50 shadow-md' 
            : 'border-gray-200 hover:border-indigo-300 hover:bg-gray-50'
        ]"
        @click="selectRole('volunteer')"
      >
        <div class="flex flex-col items-center text-center space-y-3">
          <!-- Icon -->
          <div :class="[
            'w-12 h-12 rounded-full flex items-center justify-center transition-colors duration-200',
            selectedRole === 'volunteer' ? 'bg-indigo-600' : 'bg-gray-400'
          ]">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
          
          <!-- Title -->
          <h3 :class="[
            'font-semibold text-lg transition-colors duration-200',
            selectedRole === 'volunteer' ? 'text-indigo-600' : 'text-gray-900'
          ]">
            Volunteer
          </h3>
          
          <!-- Description -->
          <p class="text-sm text-gray-600 leading-relaxed">
            Assist temples with seva activities and event management
          </p>
          
          <!-- Features List -->
          <ul class="text-xs text-gray-500 space-y-1 mt-2">
            <li>• Join temple teams</li>
            <li>• Get seva assignments</li>
            <li>• Help with events</li>
            <li>• View schedules</li>
          </ul>
        </div>
        
        <!-- Selection Indicator -->
        <div v-if="selectedRole === 'volunteer'" 
             class="absolute top-3 right-3 w-6 h-6 bg-indigo-600 rounded-full flex items-center justify-center">
          <svg class="w-4 h-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Error Message -->
    <div v-if="error" class="bg-red-100 border border-red-200 text-red-800 px-4 py-3 rounded-lg text-sm">
      {{ error }}
    </div>

    <!-- Role Description -->
    <div v-if="selectedRole" class="bg-indigo-50 border border-indigo-200 rounded-lg p-4">
      <p class="text-sm text-indigo-800">
        <span class="font-medium">{{ roleDescriptions[selectedRole].title }}:</span>
        {{ roleDescriptions[selectedRole].description }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'

const props = defineProps({
  modelValue: {
    type: String,
    default: ''
  },
  error: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue'])

const selectedRole = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
})

const roleDescriptions = {
  tenant: {
    title: 'Temple Admin',
    description: 'You will be able to create and manage temples, approve seva bookings, track donations, and organize events. Your registration will be pending approval.'
  },
  devotee: {
    title: 'Devotee',
    description: 'You can join temples, book sevas, make donations, and participate in temple events. You can start using the platform immediately after registration.'
  },
  volunteer: {
    title: 'Volunteer',
    description: 'You can join temple volunteer teams, get seva assignments, help with events, and view your schedules. You can start volunteering immediately after registration.'
  }
}

const selectRole = (role) => {
  selectedRole.value = role
}
</script>