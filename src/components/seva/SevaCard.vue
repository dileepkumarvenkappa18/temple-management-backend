<template>
  <div class="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 border border-gray-100 overflow-hidden">
    <!-- Seva Header -->
    <div class="p-6 border-b border-gray-100">
      <div class="flex items-start justify-between">
        <div class="flex-1">
          <h3 class="text-lg font-semibold text-gray-900 mb-2">{{ seva.name }}</h3>
          <p class="text-sm text-gray-600 mb-3">{{ seva.description }}</p>
          
          <!-- Seva Details Grid -->
          <div class="grid grid-cols-2 gap-4 text-sm">
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
              <span class="text-gray-700">{{ formatDate(seva.date) }}</span>
            </div>
            
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
              <span class="text-gray-700">{{ seva.time }}</span>
            </div>
            
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
              <span class="text-gray-700">₹{{ formatAmount(seva.amount) }}</span>
            </div>
            
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
              </svg>
              <span class="text-gray-700">{{ seva.devotee_name || 'Available' }}</span>
            </div>
          </div>
        </div>
        
        <!-- Status Badge -->
        <div class="ml-4">
          <span :class="getStatusBadgeClass(seva.status)" class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium">
            {{ seva.status }}
          </span>
        </div>
      </div>
    </div>
    
    <!-- Action Buttons -->
    <div class="px-6 py-4 bg-gray-50 flex items-center justify-between">
      <div class="flex items-center space-x-3 text-sm text-gray-600">
        <span>Booked by: <span class="font-medium">{{ seva.booked_by || 'Not booked' }}</span></span>
        <span v-if="seva.phone" class="text-indigo-600">{{ seva.phone }}</span>
      </div>
      
      <div class="flex items-center space-x-2">
        <button 
          @click="$emit('view', seva)"
          class="px-3 py-1.5 text-sm font-medium text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
        >
          View Details
        </button>
        
        <button 
          v-if="canEdit"
          @click="$emit('edit', seva)"
          class="px-3 py-1.5 text-sm font-medium text-gray-600 hover:text-gray-700 hover:bg-gray-100 rounded-lg transition-colors duration-200"
        >
          Edit
        </button>
        
        <button 
          v-if="seva.status === 'PENDING' && canApprove"
          @click="$emit('approve', seva)"
          class="px-3 py-1.5 text-sm font-medium text-green-600 hover:text-green-700 hover:bg-green-50 rounded-lg transition-colors duration-200"
        >
          Approve
        </button>
        
        <button 
          v-if="seva.status === 'PENDING' && canApprove"
          @click="$emit('reject', seva)"
          class="px-3 py-1.5 text-sm font-medium text-red-600 hover:text-red-700 hover:bg-red-50 rounded-lg transition-colors duration-200"
        >
          Reject
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  seva: {
    type: Object,
    required: true
  },
  canEdit: {
    type: Boolean,
    default: false
  },
  canApprove: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['view', 'edit', 'approve', 'reject'])

const getStatusBadgeClass = (status) => {
  const classes = {
    'PENDING': 'bg-yellow-100 text-yellow-800',
    'APPROVED': 'bg-green-100 text-green-800',
    'REJECTED': 'bg-red-100 text-red-800',
    'COMPLETED': 'bg-blue-100 text-blue-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (date) => {
  if (!date) return 'Not scheduled'
  return new Date(date).toLocaleDateString('en-IN', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const formatAmount = (amount) => {
  if (!amount) return '0'
  return new Intl.NumberFormat('en-IN').format(amount)
}
</script>