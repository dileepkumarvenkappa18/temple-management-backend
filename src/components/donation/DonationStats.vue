<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
    <!-- Total Amount -->
    <div class="bg-gradient-to-r from-blue-500 to-blue-600 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-blue-200 text-sm font-medium">Total Amount</p>
          <p class="text-3xl font-bold">₹{{ formatCurrency(stats.totalAmount || 0) }}</p>
          <p class="text-blue-200 text-xs mt-1">{{ stats.totalCount || 0 }} donations</p>
        </div>
        <div class="bg-white bg-opacity-20 rounded-lg p-3">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3-.895 3-2-1.343-2-3-2z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M12 14v8" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Average Donation -->
    <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-green-200 text-sm font-medium">Average Donation</p>
          <p class="text-3xl font-bold">₹{{ formatCurrency(stats.averageAmount || 0) }}</p>
          <p class="text-green-200 text-xs mt-1">per transaction</p>
        </div>
        <div class="bg-white bg-opacity-20 rounded-lg p-3">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2-2V7a2 2 0 012-2h2a2 2 0 002 2v2a2 2 0 002 2h2a2 2 0 002-2V7a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 00-2 2h-2a2 2 0 00-2 2v6" />
          </svg>
        </div>
      </div>
    </div>

    <!-- This Month -->
    <div class="bg-gradient-to-r from-purple-500 to-purple-600 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-purple-200 text-sm font-medium">This Month</p>
          <p class="text-3xl font-bold">₹{{ formatCurrency(stats.thisMonth || 0) }}</p>
          <p class="text-purple-200 text-xs mt-1">current month</p>
        </div>
        <div class="bg-white bg-opacity-20 rounded-lg p-3">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2-2z" />
          </svg>
        </div>
      </div>
    </div>

    <!-- Total Donors -->
    <div class="bg-gradient-to-r from-orange-500 to-orange-600 rounded-xl p-6 text-white">
      <div class="flex items-center justify-between">
        <div>
          <p class="text-orange-200 text-sm font-medium">Total Donors</p>
          <p class="text-3xl font-bold">{{ stats.totalDonors || 0 }}</p>
          <p class="text-orange-200 text-xs mt-1">unique devotees</p>
        </div>
        <div class="bg-white bg-opacity-20 rounded-lg p-3">
          <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" 
                  d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-1.657-1.343-3-3-3H9c-1.657 0-3 1.343-3 3v2m4-10a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
        </div>
      </div>
    </div>
  </div>

  <!-- Status Overview -->
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
    <!-- Successful -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <div class="w-10 h-10 bg-green-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7" />
            </svg>
          </div>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-gray-600">Successful</p>
          <p class="text-2xl font-bold text-gray-900">{{ stats.completedCount || 0 }}</p>
        </div>
      </div>
    </div>

    <!-- Pending -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <div class="w-10 h-10 bg-yellow-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </div>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-gray-600">Pending</p>
          <p class="text-2xl font-bold text-gray-900">{{ stats.pendingCount || 0 }}</p>
        </div>
      </div>
    </div>

    <!-- Failed -->
    <div class="bg-white rounded-lg shadow-sm border border-gray-200 p-4">
      <div class="flex items-center">
        <div class="flex-shrink-0">
          <div class="w-10 h-10 bg-red-100 rounded-full flex items-center justify-center">
            <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </div>
        </div>
        <div class="ml-4">
          <p class="text-sm font-medium text-gray-600">Failed</p>
          <p class="text-2xl font-bold text-gray-900">{{ stats.failedCount || 0 }}</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useDonationStore } from '@/stores/donation'

const donationStore = useDonationStore()

const stats = computed(() => donationStore.donationStats)

const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount || 0)
}
</script>