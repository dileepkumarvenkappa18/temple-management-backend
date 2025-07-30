<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Donation Ledger</h2>
        <p class="text-gray-600">Complete financial record of all donations</p>
      </div>
      
      <!-- Export Actions -->
      <div class="flex flex-col sm:flex-row gap-3 mt-4 sm:mt-0">
        <!-- <button class="inline-flex items-center px-4 py-2 bg-indigo-100 text-indigo-700 rounded-lg hover:bg-indigo-200 transition-all duration-200">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Export PDF
        </button>
        <button class="inline-flex items-center px-4 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition-all duration-200">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a4 4 0 01-4-4V5a4 4 0 014-4h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a4 4 0 01-4 4z" />
          </svg>
          Export Excel
        </button> -->
      </div>
    </div>

    <!-- Filters -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
        <select v-model="filters.dateRange" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
          <option value="all">All Time</option>
          <option value="today">Today</option>
          <option value="week">This Week</option>
          <option value="month">This Month</option>
          <option value="year">This Year</option>
          <option value="custom">Custom Range</option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Donation Type</label>
        <select v-model="filters.donationType" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
          <option value="all">All Types</option>
          <option value="general">General</option>
          <option value="seva">Seva</option>
          <option value="festival">Festival</option>
          <option value="construction">Construction</option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Payment Method</label>
        <select v-model="filters.paymentMethod" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
          <option value="all">All Methods</option>
          <option value="cash">Cash</option>
          <option value="online">Online</option>
          <option value="cheque">Cheque</option>
          <option value="dd">Demand Draft</option>
        </select>
      </div>
      
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Amount Range</label>
        <select v-model="filters.amountRange" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
          <option value="all">All Amounts</option>
          <option value="0-500">₹0 - ₹500</option>
          <option value="500-1000">₹500 - ₹1,000</option>
          <option value="1000-5000">₹1,000 - ₹5,000</option>
          <option value="5000+">₹5,000+</option>
        </select>
      </div>
    </div>

    <!-- Summary Cards -->
    <div v-if="!isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div class="bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-indigo-200 text-sm">Total Donations</p>
            <p class="text-2xl font-bold">₹{{ formatCurrency(donationStats?.totalAmount || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-200 text-sm">Total Donors</p>
            <p class="text-2xl font-bold">{{ donationStats?.totalDonors || 0 }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-gradient-to-r from-purple-500 to-purple-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-purple-200 text-sm">Average Donation</p>
            <p class="text-2xl font-bold">₹{{ formatCurrency(donationStats?.averageAmount || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
        </div>
      </div>
      
      <div class="bg-gradient-to-r from-amber-500 to-amber-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-amber-200 text-sm">This Month</p>
            <p class="text-2xl font-bold">₹{{ formatCurrency(donationStats?.thisMonth || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>
    
    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
      <div v-for="i in 4" :key="i" class="bg-gray-100 rounded-xl p-4 animate-pulse">
        <div class="h-16 bg-gray-200 rounded"></div>
      </div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by donor name, transaction ID, or reference..."
          class="w-full pl-10 pr-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
    </div>

    <!-- Ledger Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date & Reference
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Donor Details
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Donation Type
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Payment Method
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody v-if="!isLoading && filteredDonations.length" class="bg-white divide-y divide-gray-200">
          <tr v-for="donation in paginatedDonations" :key="donation.id || donation._id" class="hover:bg-gray-50 transition-colors duration-200">
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900 font-medium">{{ formatDate(donation.donatedAt || donation.date) }}</div>
              <div class="text-sm text-gray-500">#{{ donation.transactionId || donation.referenceID || 'N/A' }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
                  <span class="text-indigo-600 font-medium text-sm">{{ getInitials(donation.userName || donation.donorName || 'Anonymous') }}</span>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">{{ donation.userName || donation.donorName || 'Anonymous' }}</div>
                  <div class="text-sm text-gray-500">{{ donation.userEmail || donation.donorEmail || 'N/A' }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getDonationTypeClass(donation.donationType || donation.type || 'General')">
                {{ donation.donationType || donation.type || 'General' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="text-sm text-gray-900">{{ donation.paymentMethod || 'Online' }}</div>
                <div v-if="donation.referenceNumber || donation.referenceID" class="text-xs text-gray-500 ml-2">
                  Ref: {{ donation.referenceNumber || donation.referenceID }}
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-lg font-bold text-gray-900">₹{{ formatCurrency(donation.amount || 0) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                :class="getStatusClass(donation.status || 'Pending')">
                {{ donation.status || 'Pending' }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex items-center space-x-2">
                <button class="text-indigo-600 hover:text-indigo-900 p-1 rounded">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                  </svg>
                </button>
                <button @click="generateReceipt(donation.id || donation._id)" class="text-green-600 hover:text-green-900 p-1 rounded">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                </button>
                <button class="text-blue-600 hover:text-blue-900 p-1 rounded">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
        <tbody v-else-if="isLoading" class="bg-white divide-y divide-gray-200">
          <tr v-for="i in 5" :key="i" class="animate-pulse">
            <td colspan="7" class="px-6 py-4">
              <div class="h-10 bg-gray-200 rounded"></div>
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr>
            <td colspan="7" class="px-6 py-10 text-center text-gray-500">
              No donations found matching your filters
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between mt-6 px-4 py-3 bg-gray-50 rounded-lg">
      <div class="flex items-center text-sm text-gray-700">
        Showing {{ paginatedDonations.length ? (currentPage - 1) * itemsPerPage + 1 : 0 }} to {{ Math.min(currentPage * itemsPerPage, totalItems) }} of {{ totalItems }} entries
      </div>
      <div class="flex items-center space-x-2">
        <button
          @click="currentPage > 1 && currentPage--"
          :disabled="currentPage === 1"
          class="px-3 py-1 rounded-md text-sm font-medium border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>
        <span class="px-3 py-1 text-sm font-medium text-indigo-600 bg-indigo-100 rounded-md">
          {{ currentPage }}
        </span>
        <button
          @click="currentPage < totalPages && currentPage++"
          :disabled="currentPage === totalPages"
          class="px-3 py-1 rounded-md text-sm font-medium border border-gray-300 bg-white text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useDonationStore } from '@/stores/donation'

// Initialize store
const donationStore = useDonationStore()
const isLoading = computed(() => donationStore.loading)

// Reactive data
const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(10)

const filters = ref({
  dateRange: 'all',
  donationType: 'all',
  paymentMethod: 'all',
  amountRange: 'all'
})

// Get data from store instead of mock
const donationStats = computed(() => donationStore.donationStats)
const donations = computed(() => donationStore.donations || [])

// Watch for filter changes
watch(() => filters.value, (newFilters) => {
  donationStore.setFilters(newFilters)
  currentPage.value = 1 // Reset pagination when filters change
}, { deep: true })

// Computed properties
const filteredDonations = computed(() => {
  let filtered = donations.value || []

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(donation => {
      const donorName = donation.userName || donation.donorName || ''
      const transactionId = donation.transactionId || donation.referenceID || ''
      const email = donation.userEmail || donation.donorEmail || ''
      
      return donorName.toLowerCase().includes(query) ||
        transactionId.toLowerCase().includes(query) ||
        email.toLowerCase().includes(query)
    })
  }

  // Apply filters based on your filter.value if needed
  // This can be extended as needed

  return filtered
})

const paginatedDonations = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredDonations.value.slice(start, end)
})

const totalItems = computed(() => filteredDonations.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalItems.value / itemsPerPage.value)))

// Methods
const formatCurrency = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount || 0)
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const getInitials = (name) => {
  if (!name) return 'NA'
  return name.split(' ').map(n => n[0]).join('').toUpperCase()
}

const getDonationTypeClass = (type) => {
  const classes = {
    'General': 'bg-blue-100 text-blue-800',
    'general': 'bg-blue-100 text-blue-800',
    'Seva': 'bg-purple-100 text-purple-800',
    'seva': 'bg-purple-100 text-purple-800',
    'Festival': 'bg-orange-100 text-orange-800',
    'festival': 'bg-orange-100 text-orange-800',
    'Construction': 'bg-green-100 text-green-800',
    'construction': 'bg-green-100 text-green-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getStatusClass = (status) => {
  const lowercaseStatus = (status || '').toLowerCase()
  const classes = {
    'completed': 'bg-green-100 text-green-800',
    'success': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'failed': 'bg-red-100 text-red-800'
  }
  return classes[lowercaseStatus] || 'bg-gray-100 text-gray-800'
}

const generateReceipt = (donationId) => {
  if (!donationId) return
  donationStore.generateReceipt(donationId)
}

onMounted(async () => {
  // Load data from API
  try {
    await donationStore.fetchDashboard()
    await donationStore.fetchDonations()
  } catch (error) {
    console.error('Error loading donation data:', error)
  }
})
</script>