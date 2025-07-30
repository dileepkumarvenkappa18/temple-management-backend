<template>
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Your Temples</h2>
        <p class="text-gray-600 mt-1">Manage and track your temple applications</p>
      </div>
      <div class="flex space-x-2">
        <button
          @click="showDebugPanel = !showDebugPanel"
          class="px-3 py-1 bg-amber-600 text-white rounded text-sm"
        >
          Debug Panel
        </button>
        <router-link
          to="/tenant/entities/create"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition-all duration-200 shadow-md hover:shadow-lg"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Temple
        </router-link>
      </div>
    </div>

    <!-- Temple List -->
    <div v-for="temple in temples" :key="temple.id" class="bg-white rounded-lg shadow mb-4 overflow-hidden">
      <!-- Temple Header -->
      <div class="p-4">
        <div class="flex justify-between items-center mb-2">
          <h3 class="text-xl font-bold text-gray-900">{{ temple.name }}</h3>
          <div>
            <span 
              v-if="isApproved(temple)" 
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-green-100 text-green-800"
            >
              Approved
            </span>
            <span 
              v-else-if="isRejected(temple)" 
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-red-100 text-red-800"
            >
              Rejected
            </span>
            <span 
              v-else 
              class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-yellow-100 text-yellow-800"
            >
              Pending
            </span>
          </div>
        </div>
        
        <!-- Address -->
        <div class="text-sm text-gray-600 mb-3">
          {{ getAddressString(temple) }}
        </div>
        
        <!-- Created Date -->
        <div class="text-xs text-gray-500">
          <div class="flex items-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
            Created: {{ formatDate(temple.created_at) }}
          </div>
        </div>
      </div>
      
      <!-- Action Buttons -->
      <div class="flex border-t border-gray-200">
        <button
          @click="handleViewTemple(temple)"
          class="flex-1 py-3 text-gray-700 hover:bg-gray-50 transition-colors"
        >
          <span class="flex items-center justify-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
            </svg>
            View
          </span>
        </button>
        
        <button
          @click="handleEditTemple(temple)"
          class="flex-1 py-3 text-gray-700 hover:bg-gray-50 transition-colors border-l border-gray-200"
        >
          <span class="flex items-center justify-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
            </svg>
            Edit
          </span>
        </button>

        <!-- Manage Button -->
        <button
          v-if="isApproved(temple)"
          @click="handleManageTemple(temple)"
          class="flex-1 py-3 text-indigo-700 font-medium hover:bg-indigo-50 transition-colors border-l border-gray-200"
        >
          <span class="flex items-center justify-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
            </svg>
            Manage
          </span>
        </button>
      </div>
      
      <!-- Tenant Approval Required Message -->
      <div class="p-2 text-xs text-center text-gray-500 italic border-t border-gray-200">
        {{ isApproved(temple) ? '' : 'Tenant approval required' }}
      </div>
    </div>

    <!-- Debug Panel -->
    <TempleDebug 
      v-if="showDebugPanel"
      :temples="temples"
      @close="showDebugPanel = false"
      @update="forceUpdate++"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'
import { useTempleStore } from '@/stores/temple'
import TempleDebug from './TempleDebug.vue'

const router = useRouter()
const { showToast } = useToast()
const templeStore = useTempleStore()

// Reactive data
const temples = ref([])
const loading = ref(true)
const showDebugPanel = ref(false)
const forceUpdate = ref(0)

// Helper functions for status checks - already using case-insensitive comparison
const isApproved = (temple) => {
  const status = (temple.status || '').toString().toLowerCase()
  return status === 'approved'
}

const isRejected = (temple) => {
  const status = (temple.status || '').toString().toLowerCase()
  return status === 'rejected'
}

// Methods
const loadTemples = async () => {
  try {
    loading.value = true
    
    // Use the temple store to fetch real data
    await templeStore.fetchTemples()
    temples.value = templeStore.temples
    
    console.log('Loaded temples:', JSON.stringify(temples.value, null, 2))
    
    // If using mock data for development
    if (temples.value.length === 0) {
      temples.value = [
        {
          id: 1,
          name: 'Sri Venkateswara Temple',
          main_deity: 'Lord Venkateswara',
          temple_type: 'Vaishnavite',
          status: 'approved',
          created_at: '2023-01-15',
          approved_at: '2023-01-20',
          city: 'Tirupati',
          state: 'Andhra Pradesh',
          district: 'Chittoor',
          phone: '9876543210',
          email: 'info@tirumala.org',
          street_address: '2nd street, Aditya nagar'
        },
        {
          id: 2,
          name: 'Jagannath Temple',
          main_deity: 'Lord Jagannath',
          temple_type: 'Vaishnavite',
          status: 'pending',
          created_at: '2023-02-20',
          city: 'Puri',
          state: 'Odisha',
          district: 'Puri',
          phone: '9876543211',
          email: 'info@jagannath.org',
          street_address: 'Temple Street'
        }
      ]
    }
  } catch (error) {
    console.error('Error loading temples:', error)
    showToast('Failed to load temples', 'error')
  } finally {
    loading.value = false
  }
}

const getAddressString = (temple) => {
  return `${temple.street_address || ''}, ${temple.city || ''}, ${temple.state || ''}, ${temple.district || ''} ${temple.pincode || ''}`
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  } catch (e) {
    return dateString
  }
}

const handleViewTemple = (temple) => {
  router.push(`/tenant/entities/${temple.id}`)
}

const handleEditTemple = (temple) => {
  router.push(`/tenant/entities/${temple.id}/edit`)
}

const handleManageTemple = (temple) => {
  // Since we're already checking with v-if, we can be confident the temple is approved
  router.push(`/entity/${temple.id}/dashboard`)
}

// Watch for updates to refresh the data
watch(forceUpdate, () => {
  console.log('Force update triggered')
  loadTemples() // Also reload data when debug panel triggers an update
})

// Lifecycle
onMounted(() => {
  loadTemples()
})
</script>