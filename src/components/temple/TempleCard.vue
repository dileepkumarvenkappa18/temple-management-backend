<template>
  <div class="bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden group">
    <!-- Temple Image/Header -->
    <div class="relative h-48 bg-gradient-to-br from-indigo-500 to-purple-600 overflow-hidden">
      <div class="absolute inset-0 bg-black bg-opacity-20"></div>
      
      <!-- Temple Icon -->
      <div class="absolute top-4 left-4 p-3 bg-white bg-opacity-20 backdrop-blur-sm rounded-xl">
        <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
        </svg>
      </div>

      <!-- Status Badge -->
      <div class="absolute top-4 right-4">
        <TempleApprovalStatus :status="temple.status" />
      </div>

      <!-- Temple Name -->
      <div class="absolute bottom-4 left-4 right-4">
        <h3 class="text-xl font-bold text-white mb-1 group-hover:text-yellow-200 transition-colors">
          {{ temple.name }}
        </h3>
        <p class="text-indigo-100 text-sm">
          {{ temple.presiding_deity || temple.main_deity }}
        </p>
      </div>
    </div>

    <!-- Card Content -->
    <div class="p-6">
      <!-- Description -->
      <p class="text-gray-600 text-sm mb-4 line-clamp-2" v-if="temple.description">
        {{ temple.description }}
      </p>

      <!-- Temple Details Grid -->
      <div class="space-y-3 mb-6">
        <!-- Category -->
        <div class="flex items-center gap-3">
          <div class="flex-shrink-0 w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
            <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z"/>
            </svg>
          </div>
          <div>
            <p class="text-sm font-medium text-gray-900">{{ temple.category || temple.temple_type }}</p>
            <p class="text-xs text-gray-500">Category</p>
          </div>
        </div>

        <!-- Location -->
        <div class="flex items-center gap-3">
          <div class="flex-shrink-0 w-8 h-8 bg-green-100 rounded-lg flex items-center justify-center">
            <svg class="w-4 h-4 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
          </div>
          <div>
            <p class="text-sm font-medium text-gray-900">{{ temple.city }}, {{ temple.state }}</p>
            <p class="text-xs text-gray-500">Location</p>
          </div>
        </div>

        <!-- Established Year -->
        <div class="flex items-center gap-3" v-if="temple.establishedYear || temple.established_year">
          <div class="flex-shrink-0 w-8 h-8 bg-amber-100 rounded-lg flex items-center justify-center">
            <svg class="w-4 h-4 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
          </div>
          <div>
            <p class="text-sm font-medium text-gray-900">{{ temple.establishedYear || temple.established_year }}</p>
            <p class="text-xs text-gray-500">Established</p>
          </div>
        </div>
      </div>

      <!-- Contact Information -->
      <div class="border-t border-gray-100 pt-4 mb-4">
        <div class="grid grid-cols-1 gap-2">
          <div class="flex items-center gap-2" v-if="temple.phone">
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
            </svg>
            <span class="text-sm text-gray-600">{{ temple.phone }}</span>
          </div>
          
          <div class="flex items-center gap-2" v-if="temple.email">
            <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207"/>
            </svg>
            <span class="text-sm text-gray-600">{{ temple.email }}</span>
          </div>
        </div>
      </div>

      <!-- Rejection Reason (if rejected) -->
      <div v-if="isRejected && temple.rejection_reason" 
           class="mb-4 p-3 bg-red-50 border border-red-200 rounded-lg">
        <div class="flex items-start gap-2">
          <svg class="w-4 h-4 text-red-500 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <div>
            <p class="text-sm font-medium text-red-800">Rejection Reason:</p>
              <p class="text-sm text-red-700">{{ temple.rejection_reason }}</p>
            </div>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row gap-2 pt-4 border-t border-gray-100">
          <button
            @click="$emit('view', temple)"
            class="flex-1 px-4 py-2 text-sm font-medium text-indigo-700 bg-indigo-50 border border-indigo-200 rounded-lg hover:bg-indigo-100 transition-colors flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
            </svg>
            View Details
          </button>

          <!-- Edit Button (only for non-approved temples) -->
          <button
            v-if="!isApprovedCaseInsensitive"
            @click="$emit('edit', temple)"
            class="flex-1 px-4 py-2 text-sm font-medium text-gray-700 bg-gray-50 border border-gray-200 rounded-lg hover:bg-gray-100 transition-colors flex items-center justify-center gap-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
            </svg>
            Edit
          </button>

          <!-- Primary Action Button -->
          <button
            v-if="isApprovedCaseInsensitive"
            @click="handlePrimaryAction"
            class="flex-1 px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-lg hover:bg-indigo-700 transition-colors flex items-center justify-center gap-2"
          >
            <svg v-if="isJoined" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z"/>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            </svg>
            <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"/>
            </svg>
            {{ isJoined ? 'Manage' : 'Join Temple' }}
          </button>

          <!-- Delete Button (only for non-approved temples) -->
          <button
            v-if="!isApprovedCaseInsensitive && showDeleteButton"
            @click="$emit('delete', temple)"
            class="px-3 py-2 text-sm font-medium text-red-700 bg-red-50 border border-red-200 rounded-lg hover:bg-red-100 transition-colors flex items-center justify-center"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M9 7h6m-6 0V4a1 1 0 011-1h4a1 1 0 011 1v3"/>
            </svg>
          </button>
        </div>

        <!-- Creation Date -->
        <div class="flex items-center justify-between mt-4 pt-3 border-t border-gray-100 text-xs text-gray-500">
          <span>Created: {{ formatDate(temple.created_at) }}</span>
          <span v-if="temple.approved_at">Approved: {{ formatDate(temple.approved_at) }}</span>
        </div>

        <!-- Debug Info - Only visible in development -->
        <div v-if="isDevelopment" class="mt-4 pt-3 border-t border-gray-200 text-xs text-gray-500">
          <details>
            <summary class="cursor-pointer font-semibold">Debug Info</summary>
            <div class="mt-2 p-2 bg-gray-50 rounded">
              <div><b>Temple ID:</b> {{ temple.id }}</div>
              <div><b>Raw Status:</b> "{{ temple.status }}"</div>
              <div><b>Status Type:</b> {{ typeof temple.status }}</div>
              <div><b>Is Approved (case-sensitive):</b> {{ isApproved }}</div>
              <div><b>Is Approved (case-insensitive):</b> {{ isApprovedCaseInsensitive }}</div>
              <div><b>Is Joined:</b> {{ isJoined }}</div>
              <div><b>Available Keys:</b> {{ Object.keys(temple).join(', ') }}</div>
              <button
                @click="forceApprovedStatus"
                class="mt-2 px-2 py-1 bg-green-100 text-green-800 rounded text-xs"
              >
                Force Status to "approved"
              </button>
            </div>
          </details>
        </div>
      </div>
    </div>
</template>

<script setup>
import { computed, ref, onMounted } from 'vue'
import TempleApprovalStatus from './TempleApprovalStatus.vue'
import { useAuthStore } from '@/stores/auth'

// Props
const props = defineProps({
  temple: {
    type: Object,
    required: true
  },
  showDeleteButton: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['view', 'edit', 'manage', 'delete', 'join'])

// Stores
const authStore = useAuthStore()
const membershipStore = useTempleMembershipStore()

// Development mode detection
const isDevelopment = ref(false)
onMounted(() => {
  // Check if we're in development mode
  isDevelopment.value = 
    window.location.hostname === 'localhost' || 
    window.location.hostname === '127.0.0.1' ||
    process.env.NODE_ENV === 'development'
  
  // Log temple data for debugging
  if (isDevelopment.value) {
    console.log('TempleCard mounted:', props.temple)
    console.log('Temple status:', props.temple.status)
    console.log('Status data type:', typeof props.temple.status)
  }

  // Check membership status when component mounts
  checkMembershipStatus()
})

// Check if temple is already joined by the devotee
const checkMembershipStatus = async () => {
  if (authStore.isDevotee && props.temple?.id) {
    await membershipStore.checkMembership(props.temple.id)
  }
}

// Computed properties
const isApproved = computed(() => {
  const status = (props.temple.status || '').toString()
  return status === 'APPROVED'
})

// Case-insensitive check - this is what we'll use throughout the component
const isApprovedCaseInsensitive = computed(() => {
  const status = (props.temple.status || '').toString().toLowerCase()
  return status === 'approved'
})

const isRejected = computed(() => {
  const status = (props.temple.status || '').toString().toLowerCase()
  return status === 'rejected'
})

const isJoined = computed(() => {
  return membershipStore.membershipStatus[props.temple.id] || false
})

// Methods
const formatDate = (dateString) => {
  if (!dateString) return ''
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const handlePrimaryAction = () => {
  if (isJoined.value) {
    emit('manage', props.temple)
  } else {
    emit('join', props.temple)
  }
}

// Debug method to force approval status
const forceApprovedStatus = () => {
  if (props.temple) {
    props.temple.status = 'approved'
    console.log('Temple status forced to approved')
  }
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;

  /* Future-proof: */
  line-clamp: 2;
}
</style>