<template>
  <div class="bg-white rounded-xl shadow-lg border border-indigo-100">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-indigo-100 bg-gradient-to-r from-indigo-50 to-white">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-semibold text-indigo-900">Donation History</h2>
          <p class="text-sm text-indigo-600 mt-1">Track all donation transactions</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <select 
            v-model="selectedFilter" 
            class="px-4 py-2 bg-white border border-indigo-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          >
            <option value="all">All Donations</option>
            <option value="today">Today</option>
            <option value="week">This Week</option>
            <option value="month">This Month</option>
            <option value="year">This Year</option>
          </select>
          <button 
            @click="exportToExcel"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
          >
            <i class="fas fa-download mr-2"></i>
            Export
          </button>
        </div>
      </div>
    </div>

    <!-- Search and Filters -->
    <div class="px-6 py-4 bg-gray-50 border-b border-indigo-100">
      <div class="flex flex-col lg:flex-row gap-4">
        <div class="flex-1">
          <div class="relative">
            <i class="fas fa-search absolute left-3 top-1/2 transform -translate-y-1/2 text-indigo-400"></i>
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="Search by donor name, transaction ID, or purpose..."
              class="w-full pl-10 pr-4 py-2 bg-white border border-indigo-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
          </div>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <select 
            v-model="selectedType" 
            class="px-4 py-2 bg-white border border-indigo-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          >
            <option value="">All Types</option>
            <option value="seva">Seva Donation</option>
            <option value="general">General Donation</option>
            <option value="festival">Festival Fund</option>
            <option value="maintenance">Maintenance</option>
          </select>
          <div class="flex gap-2">
            <input 
              v-model="minAmount" 
              type="number" 
              placeholder="Min ₹"
              class="w-24 px-3 py-2 bg-white border border-indigo-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
            <input 
              v-model="maxAmount" 
              type="number" 
              placeholder="Max ₹"
              class="w-24 px-3 py-2 bg-white border border-indigo-200 rounded-lg text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      <span class="ml-3 text-indigo-600">Loading donations...</span>
    </div>

    <!-- Table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-indigo-50">
          <tr>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              <button @click="sortBy('date')" class="flex items-center hover:text-indigo-700 transition-colors duration-200">
                Date
                <i :class="getSortIcon('date')" class="ml-1 text-xs"></i>
              </button>
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              Donor Details
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              <button @click="sortBy('amount')" class="flex items-center hover:text-indigo-700 transition-colors duration-200">
                Amount
                <i :class="getSortIcon('amount')" class="ml-1 text-xs"></i>
              </button>
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              Purpose
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              Payment
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-4 text-left text-xs font-semibold text-indigo-900 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-indigo-100">
          <tr 
            v-for="donation in filteredDonations" 
            :key="donation.id"
            class="hover:bg-indigo-50 transition-colors duration-200"
          >
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900 font-medium">{{ formatDate(donation.date) }}</div>
              <div class="text-xs text-gray-500">{{ formatTime(donation.date) }}</div>
            </td>
            <td class="px-6 py-4">
  <div class="flex items-center">
    <div class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center mr-3">
      <i class="fas fa-user text-indigo-600"></i>
    </div>
    <div>
      <div class="text-sm font-medium text-gray-900">{{ donation.donorName }}</div>
      <div class="text-xs text-gray-500">{{ donation.donorEmail }}</div>
      <div class="text-xs text-gray-500">{{ donation.donor.phone }}</div>
      <div class="text-xs text-indigo-600">ID: {{ donation.transactionId }}</div>
    </div>
  </div>
</td>

            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-lg font-bold text-green-600">₹{{ formatAmount(donation.amount) }}</div>
              <div v-if="donation.taxDeduction" class="text-xs text-indigo-600">
                <i class="fas fa-receipt mr-1"></i>80G Eligible
              </div>
            </td>
            <td class="px-6 py-4">
              <div class="text-sm text-gray-900 font-medium">{{ donation.purpose.title }}</div>
              <span :class="getTypeColor(donation.purpose.type)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium mt-1">
                {{ donation.purpose.type }}
              </span>
              <div v-if="donation.notes" class="text-xs text-gray-500 mt-1">{{ donation.notes }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ donation.paymentMethod }}</div>
              <div class="text-xs text-gray-500">{{ donation.paymentId }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getStatusColor(donation.status)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                <span :class="getStatusDot(donation.status)" class="w-1.5 h-1.5 rounded-full mr-1.5"></span>
                {{ donation.status }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex items-center space-x-2">
                <button 
                  @click="viewReceipt(donation)"
                  class="text-indigo-600 hover:text-indigo-900 transition-colors duration-200"
                  title="View Receipt"
                >
                  <i class="fas fa-eye"></i>
                </button>
                <button 
                  @click="downloadReceipt(donation)"
                  class="text-green-600 hover:text-green-900 transition-colors duration-200"
                  title="Download Receipt"
                >
                  <i class="fas fa-download"></i>
                </button>
                <button 
                  v-if="donation.status === 'pending'"
                  @click="resendReceipt(donation)"
                  class="text-yellow-600 hover:text-yellow-900 transition-colors duration-200"
                  title="Resend Receipt"
                >
                  <i class="fas fa-paper-plane"></i>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="!loading && filteredDonations.length === 0" class="text-center py-12">
      <div class="w-24 h-24 mx-auto bg-indigo-100 rounded-full flex items-center justify-center mb-4">
        <i class="fas fa-search text-3xl text-indigo-400"></i>
      </div>
      <h3 class="text-lg font-medium text-gray-900 mb-2">No donations found</h3>
      <p class="text-gray-500">Try adjusting your search criteria or filters.</p>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="px-6 py-4 bg-gray-50 border-t border-indigo-100">
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-700">
          Showing {{ ((currentPage - 1) * itemsPerPage) + 1 }} to {{ Math.min(currentPage * itemsPerPage, totalItems) }} of {{ totalItems }} results
        </div>
        <div class="flex items-center space-x-2">
          <button 
            @click="currentPage--" 
            :disabled="currentPage === 1"
            class="px-3 py-2 text-sm bg-white border border-indigo-200 rounded-lg hover:bg-indigo-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            Previous
          </button>
          <span class="px-3 py-2 text-sm bg-indigo-600 text-white rounded-lg">{{ currentPage }}</span>
          <button 
            @click="currentPage++" 
            :disabled="currentPage === totalPages"
            class="px-3 py-2 text-sm bg-white border border-indigo-200 rounded-lg hover:bg-indigo-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import donationService from '@/services/donation.service'

// Props
const props = defineProps({
  entityId: {
    type: String,
    required: true
  }
})

// Reactive data
const loading = ref(true)
const donations = ref([])
const searchQuery = ref('')
const selectedFilter = ref('all')
const selectedType = ref('')
const minAmount = ref('')
const maxAmount = ref('')
const sortField = ref('date')
const sortDirection = ref('desc')
const currentPage = ref(1)
const itemsPerPage = ref(25)

// Computed properties
const filteredDonations = computed(() => {
  let filtered = [...donations.value]

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(donation => 
      donation.donor.name.toLowerCase().includes(query) ||
      donation.transactionId.toLowerCase().includes(query) ||
      donation.purpose.title.toLowerCase().includes(query)
    )
  }

  // Type filter
  if (selectedType.value) {
    filtered = filtered.filter(donation => donation.purpose.type === selectedType.value)
  }

  // Amount filter
  if (minAmount.value) {
    filtered = filtered.filter(donation => donation.amount >= parseFloat(minAmount.value))
  }
  if (maxAmount.value) {
    filtered = filtered.filter(donation => donation.amount <= parseFloat(maxAmount.value))
  }

  // Date filter
  const now = new Date()
  if (selectedFilter.value !== 'all') {
    filtered = filtered.filter(donation => {
      const donationDate = new Date(donation.date)
      switch (selectedFilter.value) {
        case 'today':
          return donationDate.toDateString() === now.toDateString()
        case 'week':
          const weekAgo = new Date(now.getTime() - 7 * 24 * 60 * 60 * 1000)
          return donationDate >= weekAgo
        case 'month':
          return donationDate.getMonth() === now.getMonth() && donationDate.getFullYear() === now.getFullYear()
        case 'year':
          return donationDate.getFullYear() === now.getFullYear()
        default:
          return true
      }
    })
  }

  // Sort
  filtered.sort((a, b) => {
    let aVal = a[sortField.value]
    let bVal = b[sortField.value]
    
    if (sortField.value === 'amount') {
      aVal = parseFloat(aVal)
      bVal = parseFloat(bVal)
    } else if (sortField.value === 'date') {
      aVal = new Date(aVal)
      bVal = new Date(bVal)
    }
    
    if (sortDirection.value === 'asc') {
      return aVal > bVal ? 1 : -1
    } else {
      return aVal < bVal ? 1 : -1
    }
  })

  // Pagination
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filtered.slice(start, end)
})

const totalItems = computed(() => donations.value.length)
const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))

// Methods
const sortBy = (field) => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'desc'
  }
}

const getSortIcon = (field) => {
  if (sortField.value !== field) return 'fas fa-sort text-gray-400'
  return sortDirection.value === 'asc' ? 'fas fa-sort-up text-indigo-600' : 'fas fa-sort-down text-indigo-600'
}

const getStatusColor = (status) => {
  const colors = {
    completed: 'bg-green-100 text-green-800',
    pending: 'bg-yellow-100 text-yellow-800',
    failed: 'bg-red-100 text-red-800',
    refunded: 'bg-gray-100 text-gray-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const getStatusDot = (status) => {
  const colors = {
    completed: 'bg-green-400',
    pending: 'bg-yellow-400',
    failed: 'bg-red-400',
    refunded: 'bg-gray-400'
  }
  return colors[status] || 'bg-gray-400'
}

const getTypeColor = (type) => {
  const colors = {
    seva: 'bg-indigo-100 text-indigo-800',
    general: 'bg-blue-100 text-blue-800',
    festival: 'bg-purple-100 text-purple-800',
    maintenance: 'bg-orange-100 text-orange-800'
  }
  return colors[type] || 'bg-gray-100 text-gray-800'
}

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const formatTime = (date) => {
  return new Date(date).toLocaleTimeString('en-IN', {
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount)
}

const viewReceipt = (donation) => {
  // Open receipt in modal or new tab
  console.log('View receipt:', donation.id)
}

const downloadReceipt = (donation) => {
  // Generate and download PDF receipt
  console.log('Download receipt:', donation.id)
}

const resendReceipt = (donation) => {
  // Resend receipt email
  console.log('Resend receipt:', donation.id)
}

const exportToExcel = () => {
  // Export filtered data to Excel
  console.log('Export to Excel')
}

const loadDonations = async () => {
  try {
    loading.value = true

    const params = {
      entity_id: props.entityId,   // 👈 assuming entityId is passed in
      page: 1,
      limit: 100
    }

    const response = await donationService.getDonations(params)

    // Check structure
    console.log('Fetched donations:', response)

    if (response?.data) {
      donations.value = response.data
    } else {
      donations.value = []
    }
  } catch (error) {
    console.error('Error loading donations:', error)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadDonations()
})
</script>