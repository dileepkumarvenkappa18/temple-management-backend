<template>
  <div class="min-h-screen bg-gray-50 p-4 sm:p-6 lg:p-8">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-2">Donation Management</h1>
          <p class="text-gray-600">Track and manage all donations for your temple</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3 mt-4 sm:mt-0">
          <button @click="generateReport" class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-xl text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            Export Report
          </button>
          <button @click="newDonation" class="inline-flex items-center px-4 py-2 bg-indigo-600 border border-transparent rounded-xl text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            New Donation
          </button>
        </div>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div v-if="stats" class="grid gap-4 sm:gap-6 mb-6 sm:mb-8 grid-cols-1 sm:grid-cols-2 xl:grid-cols-4">
      <!-- Total Donations -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200 border-l-4 border-indigo-500">
        <div class="p-4 sm:p-6">
          <div class="flex items-center">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-600 mb-1">Total Donations</p>
              <p class="text-2xl sm:text-3xl font-bold text-gray-900">{{ formatCurrency(stats?.totalAmount || 0) }}</p>
              <div class="flex items-center mt-2">
                <span class="text-sm text-gray-500">from {{ stats?.total || 0 }} transactions</span>
              </div>
            </div>
            <div class="p-3 bg-indigo-100 rounded-lg flex-shrink-0 ml-3">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Average Donation -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200 border-l-4 border-green-500">
        <div class="p-4 sm:p-6">
          <div class="flex items-center">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-600 mb-1">Average Donation</p>
              <p class="text-2xl sm:text-3xl font-bold text-gray-900">{{ formatCurrency(stats?.averageAmount || 0) }}</p>
              <div class="flex items-center mt-2">
                <span class="text-sm text-gray-500">per transaction</span>
              </div>
            </div>
            <div class="p-3 bg-green-100 rounded-lg flex-shrink-0 ml-3">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- This Month -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200 border-l-4 border-purple-500">
        <div class="p-4 sm:p-6">
          <div class="flex items-center">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-600 mb-1">This Month</p>
              <p class="text-2xl sm:text-3xl font-bold text-gray-900">{{ formatCurrency(stats?.thisMonth || 0) }}</p>
              <div class="flex items-center mt-2">
                <svg class="h-4 w-4 text-purple-500 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                </svg>
                <span class="text-sm text-gray-500">in donations</span>
              </div>
            </div>
            <div class="p-3 bg-purple-100 rounded-lg flex-shrink-0 ml-3">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Total Donors -->
      <div class="bg-white rounded-xl shadow-md overflow-hidden hover:shadow-lg transition-shadow duration-200 border-l-4 border-yellow-500">
        <div class="p-4 sm:p-6">
          <div class="flex items-center">
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-600 mb-1">Total Donors</p>
              <p class="text-2xl sm:text-3xl font-bold text-gray-900">{{ stats?.totalDonors || 0 }}</p>
              <div class="flex items-center mt-2">
                <span class="text-sm text-gray-500">unique contributors</span>
              </div>
            </div>
            <div class="p-3 bg-yellow-100 rounded-lg flex-shrink-0 ml-3">
              <svg class="h-6 w-6 sm:h-7 sm:w-7 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="bg-white rounded-xl shadow-md p-6 text-center">
      <div class="animate-pulse">
        <div class="h-8 bg-gray-200 rounded w-1/4 mx-auto mb-4"></div>
        <div class="h-32 bg-gray-200 rounded mb-4"></div>
        <div class="h-4 bg-gray-200 rounded w-1/2 mx-auto"></div>
      </div>
      <p class="text-gray-500 mt-4">Loading donation statistics...</p>
    </div>

    <!-- Donation Status Summary -->
    <div v-if="stats" class="grid grid-cols-1 md:grid-cols-2 gap-4 sm:gap-6 mb-6 sm:mb-8">
      <!-- Completed vs Pending -->
      <div class="bg-white rounded-xl shadow-md p-4 sm:p-6 hover:shadow-lg transition-shadow duration-200">
        <h3 class="text-lg font-semibold text-gray-900 mb-4">Donation Status</h3>
        
        <div class="grid grid-cols-2 gap-4">
          <div class="bg-green-50 rounded-lg p-4 text-center">
            <div class="inline-flex items-center justify-center w-12 h-12 bg-green-100 rounded-full mb-3">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <h4 class="text-2xl font-bold text-gray-900">{{ stats?.completed || 0 }}</h4>
            <p class="text-sm text-gray-600 mt-1">Completed</p>
          </div>
          
          <div class="bg-yellow-50 rounded-lg p-4 text-center">
            <div class="inline-flex items-center justify-center w-12 h-12 bg-yellow-100 rounded-full mb-3">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <h4 class="text-2xl font-bold text-gray-900">{{ stats?.pending || 0 }}</h4>
            <p class="text-sm text-gray-600 mt-1">Pending</p>
          </div>
        </div>
      </div>

      <!-- Recent Donations -->
      <div class="bg-white rounded-xl shadow-md p-4 sm:p-6 hover:shadow-lg transition-shadow duration-200">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900">Recent Donations</h3>
          <button class="text-sm text-indigo-600 hover:text-indigo-800 font-medium">View All</button>
        </div>
        
        <div class="space-y-3" v-if="recentDonations && recentDonations.length">
          <div 
            v-for="donation in recentDonations.slice(0, 3)" 
            :key="donation.id || donation._id"
            class="flex items-center justify-between p-3 bg-gray-50 rounded-lg hover:bg-gray-100 transition-colors duration-200"
          >
            <div class="flex items-center flex-1 min-w-0">
              <div class="flex-1 min-w-0">
                <p class="text-sm font-medium text-gray-900 truncate">{{ donation.userName || 'Anonymous' }}</p>
                <p class="text-xs text-gray-500">{{ formatDate(donation.donatedAt || donation.date) }}</p>
              </div>
            </div>
            <div class="ml-4">
              <span class="text-sm font-semibold text-gray-900">{{ formatCurrency(donation.amount || 0) }}</span>
            </div>
          </div>
        </div>
        
        <div v-else class="text-center py-3">
          <p class="text-gray-500 text-sm">No recent donations found</p>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
// Fix #1: Import watch from Vue
import { ref, computed, onMounted, watch } from 'vue'
import { useDonationStore } from '@/stores/donation'

export default {
  name: 'DonationStats',
  
  setup() {
    const donationStore = useDonationStore()
    
    // Create computed properties for stats
    const stats = computed(() => donationStore.donationStats)
    const recentDonations = computed(() => donationStore.recentDonations)
    
    // Fetch dashboard data on component mount
    onMounted(async () => {
      try {
        await donationStore.fetchDashboard()
        await donationStore.fetchDonations()
      } catch (error) {
        console.error('Error loading donation statistics:', error)
      }
    })
    
    // Watch for filter changes
    watch(() => donationStore.filters, async (newFilters) => {
      try {
        await donationStore.fetchDashboard()
      } catch (error) {
        console.error('Error updating donation statistics:', error)
      }
    }, { deep: true })
    
    // Utility functions
    const formatCurrency = (amount) => {
      return new Intl.NumberFormat('en-IN', {
        style: 'currency',
        currency: 'INR',
        maximumFractionDigits: 0
      }).format(amount || 0)
    }
    
    const formatDate = (dateString) => {
      if (!dateString) return '-'
      return new Date(dateString).toLocaleDateString('en-IN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric'
      })
    }
    
    const truncate = (str, length = 20) => {
      if (!str) return ''
      return str.length > length ? str.substring(0, length) + '...' : str
    }
    
    const generateReport = () => {
      console.log('Generating donation report...')
      // Implementation for report generation
    }
    
    const newDonation = () => {
      console.log('Creating new donation...')
      // Implementation for new donation
    }
    
    return {
      stats,
      recentDonations,
      formatCurrency,
      formatDate,
      truncate,
      generateReport,
      newDonation
    }
  }
}
</script>