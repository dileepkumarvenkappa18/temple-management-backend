<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Tenant Dashboard</h1>
            <p class="text-gray-600 mt-1">
              Manage your temple registrations and applications
              <span v-if="tenantId" class="text-indigo-600 font-medium"> (Tenant ID: {{ tenantId }})</span>
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 px-4 py-2 rounded-lg border border-indigo-200">
              <span class="text-indigo-800 font-medium">{{ userStore.user?.name }}</span>
              <span class="text-indigo-600 text-sm ml-2">(Tenant)</span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Welcome Banner with Approval Status -->
      <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 rounded-2xl shadow-lg p-8 mb-8 text-white">
        <div class="flex flex-col md:flex-row md:items-center md:justify-between">
          <div class="flex-1">
            <h2 class="text-3xl font-bold mb-2">Welcome back!</h2>
            <p class="text-indigo-100 text-lg mb-4">Manage your temple applications and track their approval status</p>
            
            <!-- Account Status
            <div class="flex items-center space-x-3">
              <div class="flex items-center">
                <div :class="[
                  'w-3 h-3 rounded-full mr-2',
                  isUserApproved ? 'bg-green-400' : 
                  isUserRejected ? 'bg-red-400' : 'bg-yellow-400'
                ]"></div>
                <span class="text-indigo-100">Account Status: </span>
                <span class="font-semibold ml-1">{{ userStore.user?.status || 'pending' }}</span>
              </div>
            </div>
             -->
          </div>
          
          <div class="mt-6 md:mt-0">
            <router-link 
              to="/tenant/entities/create"
              class="inline-flex items-center bg-white text-indigo-600 px-6 py-3 rounded-xl font-semibold hover:bg-indigo-50 transition-all duration-200 shadow-lg hover:shadow-xl transform hover:-translate-y-1"
            >
              <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
              </svg>
              Create New Temple
            </router-link>
          </div>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-indigo-100 rounded-xl">
              <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H3m2 0h4M9 7h6m-6 4h6m-6 4h6"/>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Total Temples</p>
              <p class="text-2xl font-bold text-gray-900">{{ templeStore.temples.length }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-green-100 rounded-xl">
              <svg class="w-8 h-8 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Approved</p>
              <p class="text-2xl font-bold text-gray-900">{{ approvedTemplesCount }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-sm border border-gray-200 p-6 hover:shadow-md transition-shadow duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-amber-100 rounded-xl">
              <svg class="w-8 h-8 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm text-gray-600">Pending</p>
              <p class="text-2xl font-bold text-gray-900">{{ pendingTemplesCount }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Temple List Section -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="p-6 border-b border-gray-200">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h3 class="text-xl font-bold text-gray-900">Your Temples</h3>
              <p class="text-gray-600 mt-1">Manage and track your temple applications</p>
            </div>
            <div class="mt-4 sm:mt-0">
              <router-link 
                to="/tenant/entities/create"
                class="inline-flex items-center bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200 font-medium"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                </svg>
                Add Temple
              </router-link>
            </div>
          </div>
        </div>

        <!-- Temples List -->
        <div v-if="templeStore.temples.length > 0" class="divide-y divide-gray-200">
          <div 
            v-for="temple in templeStore.temples" 
            :key="temple.id"
            class="p-6 hover:bg-gray-50 transition-colors duration-150"
          >
            <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between">
              <!-- Temple Info -->
              <div class="flex-1">
                <div class="flex items-center mb-2">
                  <h4 class="text-lg font-semibold text-gray-900 mr-3">{{ temple.name }}</h4>
                  <TempleApprovalStatus :status="temple.status" />
                </div>
                <p class="text-gray-600 mb-2">{{ formatAddress(temple.address) }}</p>
                <div class="flex flex-wrap items-center text-sm text-gray-500 gap-4">
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    Created: {{ formatDate(temple.createdAt) }}
                  </span>
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                    {{ temple.city }}, {{ temple.state }}
                  </span>
                </div>
              </div>

              <!-- Actions -->
              <div class="mt-4 lg:mt-0 flex items-center space-x-3">
                <button
                  @click="viewTemple(temple)"
                  class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors duration-200"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                  </svg>
                  View
                </button>
                
                <button
                  @click="editTemple(temple)"
                  class="inline-flex items-center px-4 py-2 border border-indigo-300 rounded-lg text-sm font-medium text-indigo-700 bg-indigo-50 hover:bg-indigo-100 transition-colors duration-200"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                  </svg>
                  Edit
                </button>

                <!-- Manage Button - UPDATED TO ONLY CHECK TEMPLE STATUS -->
                <button
                  v-if="isTempleApproved(temple)"
                  @click="manageTemple(temple)"
                  class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200"
                >
                  <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6V4m0 2a2 2 0 100 4m0-4a2 2 0 110 4m-6 8a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4m6 6v10m6-2a2 2 0 100-4m0 4a2 2 0 100 4m0-4v2m0-6V4"/>
                  </svg>
                  Manage
                </button>
                <!-- Removed conditional showing "Tenant approval required" -->
              </div>
            </div>

            <!-- Rejection Notes (if rejected) -->
            <div v-if="isTempleRejected(temple) && temple.adminNotes" class="mt-4 p-4 bg-red-50 border border-red-200 rounded-lg">
              <div class="flex items-start">
                <svg class="w-5 h-5 text-red-400 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <div>
                  <h5 class="text-sm font-medium text-red-800 mb-1">Admin Notes:</h5>
                  <p class="text-sm text-red-700">{{ temple.adminNotes }}</p>
                  <button
                    @click="reUploadDocuments(temple)"
                    class="mt-2 text-sm text-red-600 hover:text-red-800 font-medium underline"
                  >
                    Re-upload Documents
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Empty State -->
        <div v-else class="p-12 text-center">
          <svg class="w-16 h-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H3m2 0h4M9 7h6m-6 4h6m-6 4h6"/>
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No temples registered yet</h3>
          <p class="text-gray-600 mb-6">Get started by creating your first temple registration</p>
          <router-link 
            to="/tenant/entities/create"
            class="inline-flex items-center bg-indigo-600 text-white px-6 py-3 rounded-lg hover:bg-indigo-700 transition-colors duration-200 font-medium"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Create Your First Temple
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, onMounted, watch, ref } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useTempleStore } from '@/stores/temple'
import { useAuthStore } from '@/stores/auth'
import { useToast } from '@/composables/useToast'
import TempleApprovalStatus from '@/components/temple/TempleApprovalStatus.vue'

const router = useRouter()
const route = useRoute()
const templeStore = useTempleStore()
const userStore = useAuthStore()
const { showToast } = useToast()

// Status check helper functions
const isTempleApproved = (temple) => {
  const status = (temple.status || '').toString().toLowerCase()
  return status === 'approved'
}

const isTempleRejected = (temple) => {
  const status = (temple.status || '').toString().toLowerCase()
  return status === 'rejected'
}

const isTemplePending = (temple) => {
  const status = (temple.status || '').toString().toLowerCase()
  return status === 'pending'
}

// User status helper computed properties
const isUserApproved = computed(() => {
  const status = (userStore.user?.status || '').toString().toLowerCase()
  return status === 'approved'
})

const isUserRejected = computed(() => {
  const status = (userStore.user?.status || '').toString().toLowerCase()
  return status === 'rejected'
})

const isUserPending = computed(() => {
  const status = (userStore.user?.status || '').toString().toLowerCase()
  return status === 'pending' || !status
})

// Get the tenant ID from route params or auth store
const tenantId = computed(() => {
  // First try to get from route params
  if (route.params.tenantId) {
    console.log('Using tenantId from route params:', route.params.tenantId)
    return route.params.tenantId
  }
  
  // Fallback to auth store
  const authTenantId = userStore.user?.id || localStorage.getItem('current_tenant_id')
  console.log('Using tenantId from auth store:', authTenantId)
  return authTenantId
})

// Reset temple data to ensure clean state
const resetTempleData = () => {
  // Clear any cached temple data
  localStorage.removeItem('temple_data')
  localStorage.removeItem('temple_cache')
  // Optional: Clear temples array if your store supports it
  if (typeof templeStore.clearTempleData === 'function') {
    templeStore.clearTempleData()
  }
}

// Computed properties
const approvedTemplesCount = computed(() => {
  return templeStore.temples.filter(temple => isTempleApproved(temple)).length
})

const pendingTemplesCount = computed(() => {
  return templeStore.temples.filter(temple => isTemplePending(temple)).length
})

// Methods
const formatDate = (dateString) => {
  return new Date(dateString).toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const viewTemple = (temple) => {
  // Navigate to temple view/details page
  router.push(`/tenant/entities/${temple.id}`)
}

const editTemple = (temple) => {
  // Navigate to temple edit page
  router.push(`/tenant/entities/${temple.id}/edit`)
}

const manageTemple = (temple) => {
  // Navigate to temple management dashboard
  router.push(`/entity/${temple.id}/dashboard`)
}

const reUploadDocuments = (temple) => {
  // Navigate to document re-upload page
  router.push(`/tenant/entities/${temple.id}/documents`)
}

// Watch for changes to tenantId and reload data if it changes
watch(tenantId, async (newValue, oldValue) => {
  if (newValue && newValue !== oldValue) {
    console.log(`TenantId changed from ${oldValue} to ${newValue}, reloading data...`)
    await loadTempleData()
  }
})

// Load temple data for the current tenant
const loadTempleData = async () => {
  if (!tenantId.value) {
    console.warn('No tenant ID available, cannot load temples')
    return
  }
  
  try {
    // Reset temple data to ensure clean state
    resetTempleData()
    
    // Load temples for current tenant with the tenantId parameter
    console.log('Loading temples for tenant ID:', tenantId.value)
    await templeStore.fetchTemples(tenantId.value)
  } catch (error) {
    console.error('Error loading temple data:', error)
    showToast('Failed to load temple data. Please try again.', 'error')
  }
}

const formatAddress = (address) => {
  try {
    const addr = typeof address === 'string' ? JSON.parse(address) : address
    if (!addr) return ''
    return `${addr.street || ''}, ${addr.city || ''}, ${addr.state || ''}, ${addr.pincode || ''}, ${addr.country || ''}`
  } catch {
    return address || ''
  }
}

// Lifecycle
onMounted(async () => {
  console.log('TenantDashboard mounted with tenantId:', tenantId.value)
  
  // Load temples for current tenant
  await loadTempleData()
  
  // Check if tenant is pending approval
  if (isUserPending.value) {
    // Show notification about pending status
    showToast('Your account is pending approval from the administrator.', 'info')
  } else if (isUserRejected.value) {
    showToast('Your account has been rejected. Please contact support.', 'error')
  }
})
</script>