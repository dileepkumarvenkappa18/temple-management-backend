<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Donation Ledger</h2>
        <p class="text-gray-600">Complete financial record</p>
      </div>
      <div class="flex flex-col sm:flex-row gap-3"></div>
    </div>

    <!-- Filters -->
    <div class="grid grid-cols-1 md:grid-cols-4 gap-4 mb-6 p-4 bg-gray-50 rounded-lg">
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
        <select
          v-model="filters.dateRange"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        >
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
        <select
          v-model="filters.donationType"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        >
          <option value="all">All Types</option>
          <option value="general">General</option>
          <option value="seva">Seva</option>
          <option value="festival">Festival</option>
          <option value="construction">Construction</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Payment Method</label>
        <select
          v-model="filters.paymentMethod"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        >
          <option value="all">All Methods</option>
          <option value="cash">Cash</option>
          <option value="online">Online</option>
          <option value="cheque">Cheque</option>
          <option value="dd">Demand Draft</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Amount Range</label>
        <select
          v-model="filters.amountRange"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        >
          <option value="all">All Amounts</option>
          <option value="0-500">₹0 - ₹500</option>
          <option value="500-1000">₹500 - ₹1,000</option>
          <option value="1000-5000">₹1,000 - ₹5,000</option>
          <option value="5000+">₹5,000+</option>
        </select>
      </div>
    </div>

    <!-- Summary Cards -->
    <div v-if="!isLoading" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      <div class="bg-gradient-to-r from-indigo-500 to-indigo-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-indigo-200 text-sm">Total Donations</p>
            <p class="text-2xl font-bold text-white">₹{{ formatCurrency(donationStats?.totalAmount || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3-0.895 3-2-1.343-2-3-2z M12 8v-1 m0 1v8"
              />
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-r from-green-500 to-green-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-200 text-sm">Total Donors</p>
            <p class="text-2xl font-bold text-white">{{ donationStats?.totalDonors || 0 }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2m-2-8v-1m0 1v8"
              />
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-r from-amber-500 to-amber-600 rounded-xl p-4 text-white">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-amber-200 text-sm">This Month</p>
            <p class="text-2xl font-bold text-white">₹{{ formatCurrency(donationStats?.thisMonth || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path
                stroke-linecap="round"
                stroke-linejoin="round"
                stroke-width="2"
                d="M8 7V3m8 4V3m-9 8h10 M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5z"
              />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      <div v-for="i in 3" :key="i" class="animate-pulse bg-gray-200 h-16 rounded"></div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"
            />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by donor name, transaction ID, or reference"
          class="w-full px-4 py-3 pl-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
      </div>
    </div>

    <!-- Donations Table -->
    <div class="overflow-x-auto">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date & Reference</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donor Details</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donation Type</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Payment Method</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
          </tr>
        </thead>
        <tbody v-if="!isLoading && paginatedDonations.length" class="bg-white divide-y divide-gray-200">
          <tr
            v-for="donation in paginatedDonations"
            :key="donation.id"
            class="hover:bg-gray-50 transition-colors duration-200"
          >
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ formatDate(donation.donatedAt || donation.date) }}</div>
              <div class="text-sm text-gray-500">#{{ donation.transactionId || donation.referenceID || 'N/A' }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="flex items-center justify-center w-10 h-10 rounded-full bg-indigo-100">
                  <span class="text-indigo-600 font-medium text-sm">{{ getInitials(getDonorName(donation)) }}</span>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">{{ getDonorName(donation) }}</div>
                  <div class="text-sm text-gray-500">{{ getDonorEmail(donation) }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="getDonationTypeClass(donation.donationType || donation.type || 'General')"
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
              >
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
              <span
                :class="getStatusClass(donation.status || 'Pending')"
                class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium"
              >
                {{ donation.status || 'Pending' }}
              </span>
            </td>
          </tr>
        </tbody>
        <tbody v-else-if="isLoading" class="bg-white divide-y divide-gray-200">
          <tr v-for="i in 5" :key="i" class="animate-pulse">
            <td colspan="6" class="px-6 py-4">
              <div class="h-16 bg-gray-200 rounded"></div>
            </td>
          </tr>
        </tbody>
        <tbody v-else>
          <tr>
            <td colspan="6" class="px-6 py-10 text-center text-gray-500">
              No donations found matching your filters
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between mt-6 px-4 py-3 bg-gray-50 rounded-lg">
      <div class="text-sm text-gray-700">
        Showing
        {{ paginatedDonations.length ? (currentPage - 1) * itemsPerPage + 1 : 0 }} to
        {{ Math.min(currentPage * itemsPerPage, totalItems) }} of {{ totalItems }} entries
      </div>
      <div class="flex space-x-2">
        <button
          :disabled="currentPage === 1"
          @click="handlePageChange(currentPage - 1)"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Previous
        </button>
        <span class="px-3 py-1 rounded-md bg-indigo-100 text-indigo-600 text-sm font-medium">{{ currentPage }}</span>
        <button
          :disabled="currentPage === totalPages"
          @click="handlePageChange(currentPage + 1)"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
        >
          Next
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useDonationStore } from '@/stores/donation'

const donationStore = useDonationStore()

// Reactive states
const isLoading = computed(() => donationStore.loading)
const donationStats = computed(() => donationStore.donationStats)
const donations = computed(() => donationStore.donations || [])

const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(10)

const filters = ref({
  dateRange: 'all',
  donationType: 'all',
  paymentMethod: 'all',
  amountRange: 'all',
})

// Helper functions
function safeString(val, fallback) {
  return val && val.trim() !== '' ? val : fallback
}
function getDonorName(donation) {
  return safeString(donation.userName || donation.donorName, 'Anonymous')
}
function getDonorEmail(donation) {
  return safeString(donation.userEmail || donation.donorEmail, 'N/A')
}
function getInitials(name) {
  if (!name || name.trim() === '') return 'NA'
  return name
    .trim()
    .split(/\s+/)
    .map((n) => n[0])
    .join('')
    .toUpperCase()
}
function formatCurrency(amount) {
  return new Intl.NumberFormat('en-IN').format(amount || 0)
}
function formatDate(dateString) {
  if (!dateString) return 'N/A'
  return new Date(dateString).toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
  })
}
function getDonationTypeClass(type) {
  switch ((type || '').toLowerCase()) {
    case 'general':
      return 'bg-gray-100 text-gray-800'
    case 'seva':
      return 'bg-green-100 text-green-800'
    case 'festival':
      return 'bg-yellow-100 text-yellow-800'
    case 'construction':
      return 'bg-blue-100 text-blue-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}
function getStatusClass(status) {
  switch ((status || '').toLowerCase()) {
    case 'pending':
      return 'bg-yellow-100 text-yellow-800'
    case 'completed':
      return 'bg-green-100 text-green-800'
    case 'failed':
      return 'bg-red-100 text-red-800'
    case 'cancelled':
      return 'bg-gray-100 text-gray-600'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// Watch filters and update donationStore filters + fetch data
watch(
  filters,
  (newFilters) => {
    const mappedFilters = {
      page: currentPage.value,
      limit: itemsPerPage.value,
      dateRange: newFilters.dateRange || 'all',
      type: newFilters.donationType !== 'all' ? newFilters.donationType : '',
      method: newFilters.paymentMethod !== 'all' ? newFilters.paymentMethod : '',
      search: searchQuery.value.trim() || '',
    }

    // Handle amount range
    switch (newFilters.amountRange) {
      case '0-500':
        mappedFilters.min = 0
        mappedFilters.max = 500
        break
      case '500-1000':
        mappedFilters.min = 500
        mappedFilters.max = 1000
        break
      case '1000-5000':
        mappedFilters.min = 1000
        mappedFilters.max = 5000
        break
      case '5000+':
        mappedFilters.min = 5000
        mappedFilters.max = ''
        break
      default:
        mappedFilters.min = ''
        mappedFilters.max = ''
    }

    donationStore.setFilters(mappedFilters)
    currentPage.value = 1
    fetchDonations()
  },
  { deep: true }
)

// Watch searchQuery to reset page and fetch
watch(searchQuery, () => {
  currentPage.value = 1
  fetchDonations()
})

// Client-side filter for search only (other filters handled backend or in store)
const filteredDonations = computed(() => {
  if (!donations.value || !Array.isArray(donations.value)) return []
  if (!searchQuery.value.trim()) return donations.value
  const query = searchQuery.value.toLowerCase()
  return donations.value.filter((donation) => {
    const name = getDonorName(donation).toLowerCase()
    const transactionId = (donation.transactionId || donation.referenceID || '').toLowerCase()
    const email = getDonorEmail(donation).toLowerCase()
    return name.includes(query) || transactionId.includes(query) || email.includes(query)
  })
})

// Pagination slicing on filtered list
const paginatedDonations = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredDonations.value.slice(start, end)
})

// Total count for pagination UI
const totalItems = computed(() => filteredDonations.value.length)
const totalPages = computed(() => Math.max(1, Math.ceil(totalItems.value / itemsPerPage.value)))

// Handle changing page
const handlePageChange = (page) => {
  if (page < 1 || page > totalPages.value) return
  currentPage.value = page
  donationStore.setFilters({ page })
  fetchDonations()
}

// Fetch donations from store
const fetchDonations = async () => {
  try {
    await donationStore.fetchDonations()
  } catch (error) {
    console.error('Error fetching donations:', error)
  }
}

// Fetch dashboard from store
const fetchDashboard = async () => {
  try {
    await donationStore.fetchDashboard()
  } catch (error) {
    console.error('Error fetching donation dashboard:', error)
  }
}

// Initial load
onMounted(async () => {
  await fetchDashboard()
  await fetchDonations()
})
</script>
