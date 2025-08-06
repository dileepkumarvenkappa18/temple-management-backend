<template>
  <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100">
    <div class="flex items-center justify-between mb-6">
      <h3 class="text-lg font-semibold text-gray-900 flex items-center">
        <svg class="w-5 h-5 text-indigo-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        Recent Donations
      </h3>
      <router-link 
        :to="{ name: 'DonationManagement', params: { id: entityId } }"
        class="text-sm text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200"
      >
        View All
      </router-link>
    </div>

    <div class="space-y-4">
      <!-- Loading State -->
      <div v-if="isLoading" class="space-y-3">
        <div v-for="i in 5" :key="i" class="animate-pulse">
          <div class="flex items-start space-x-3">
            <div class="w-8 h-8 bg-gray-200 rounded-full"></div>
            <div class="flex-1 space-y-2">
              <div class="h-4 bg-gray-200 rounded w-3/4"></div>
              <div class="h-3 bg-gray-200 rounded w-1/2"></div>
            </div>
          </div>
        </div>
      </div>

      <!-- Donation Items -->
      <div v-else-if="donations.length > 0" class="space-y-4">
        <div 
          v-for="donation in donations" 
          :key="donation.id"
          class="flex items-start space-x-3 p-3 hover:bg-gray-50 rounded-lg transition-colors duration-200"
        >
          <!-- Donation Icon -->
          <div class="flex-shrink-0">
            <div class="w-8 h-8 rounded-full flex items-center justify-center bg-green-100 text-green-600">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
          </div>

          <!-- Donation Content -->
          <div class="flex-1 min-w-0">
            <div class="flex items-center justify-between">
              <p class="text-sm text-gray-900 font-medium">
                {{ donation.user?.name || 'Anonymous' }} donated
              </p>
              <span class="text-xs text-gray-500">{{ formatTimeAgo(donation.donatedAt) }}</span>
            </div>
            <p class="text-sm text-gray-600 mt-1">
              {{ donation.donationType || 'General Donation' }}
            </p>
            
            <!-- Donation Metadata -->
            <div class="flex items-center mt-2 space-x-4">
              <span class="text-xs font-medium text-green-600">
                ₹{{ formatAmount(donation.amount) }}
              </span>
              <span class="text-xs text-gray-500">
                {{ donation.method || 'Online' }}
              </span>
            </div>
          </div>

          <!-- View Action -->
          <div class="flex-shrink-0">
            <button 
              class="text-xs text-indigo-600 hover:text-indigo-800 font-medium px-2 py-1 rounded hover:bg-indigo-50 transition-all duration-200"
              @click="viewDonation(donation.id)"
            >
              View
            </button>
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else class="text-center py-8">
        <svg class="w-12 h-12 text-gray-300 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <p class="text-gray-500 text-sm">No recent donations</p>
        <p class="text-gray-400 text-xs mt-1">Donations will appear here as they are made</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import donationService from '@/services/donation.service'
import { useAuth } from '@/composables/useAuth'

const route = useRoute()
const { user } = useAuth()
const entityId = computed(() => user.value?.entityId || route.params.id)

// Reactive data
const isLoading = ref(false)
const donations = ref([])

// Methods
const formatTimeAgo = (timestamp) => {
  if (!timestamp) return 'Just now'
  const now = new Date()
  const diff = now - new Date(timestamp)
  const minutes = Math.floor(diff / (1000 * 60))
  const hours = Math.floor(minutes / 60)
  const days = Math.floor(hours / 24)

  if (minutes < 60) return `${minutes}m ago`
  if (hours < 24) return `${hours}h ago`
  return `${days}d ago`
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount || 0)
}

const viewDonation = (donationId) => {
  // Navigate to donation details or show modal
  console.log('View donation:', donationId)
}

const fetchRecentDonations = async () => {
  try {
    isLoading.value = true
    
    if (!entityId.value) {
      console.error('No entity ID available')
      return
    }

    const response = await donationService.getRecentDonations(entityId.value, 5)
    donations.value = response.data || []
  } catch (error) {
    // Error handling is done in the service
  } finally {
    isLoading.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchRecentDonations()
})
</script>