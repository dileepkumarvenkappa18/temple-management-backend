<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900 mb-2">Donation Ledger</h2>
        <p class="text-gray-600">Complete financial record</p>
      </div>
      <div class="flex flex-col sm:flex-row gap-3">
        <button
          @click="refreshData"
          :disabled="isLoading"
          class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 transition-colors duration-200"
        >
          <svg v-if="isLoading" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {{ isLoading ? 'Refreshing...' : 'Refresh' }}
        </button>
      </div>
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
          <option value="annadanam">Annadanam</option>
          <option value="education">Education</option>
          <option value="maintenance">Maintenance</option>
          <option value="event">Event</option>
        </select>
      </div>

      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">Payment Status</label>
        <select
          v-model="filters.status"
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        >
          <option value="all">All Status</option>
          <option value="success">Success</option>
          <option value="pending">Pending</option>
          <option value="failed">Failed</option>
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
            <p class="text-indigo-200 text-sm">Total Amount</p>
            <p class="text-2xl font-bold text-white">₹{{ formatCurrency(donationStats?.totalAmount || 0) }}</p>
          </div>
          <div class="bg-white bg-opacity-20 rounded-lg p-2">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3-.895 3-2-1.343-2-3-2z" />
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8V7m0 1v8m0-8c1.657 0 3 .895 3 2s-1.343 2-3 2-3-.895-3-2 1.343-2 3-2z" />
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
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-3.313-2.687-6-6-6s-6 2.687-6 6v2m2-10a4 4 0 118 0 4 4 0 01-8 0z" />
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
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 002 2z" />
            </svg>
          </div>
        </div>
      </div>
    </div>

    <div v-else class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
      <div v-for="i in 3" :key="i" class="animate-pulse bg-gray-200 h-20 rounded-xl"></div>
    </div>

    <!-- Search Bar -->
    <div class="mb-6">
      <div class="relative">
        <div class="absolute inset-y-0 left-0 flex items-center pl-3 pointer-events-none">
          <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by donor name, email, transaction ID..."
          class="w-full px-4 py-3 pl-10 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
        />
        <div v-if="searchQuery" class="absolute inset-y-0 right-0 flex items-center pr-3">
          <button @click="clearSearch" class="text-gray-400 hover:text-gray-600">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
      </div>
    </div>

    <!-- Donations Table -->
    <div class="overflow-x-auto bg-white border border-gray-200 rounded-lg">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date & ID</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Donor Info</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Method</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
          </tr>
        </thead>
        <tbody v-if="!isLoading && currentDonations.length" class="bg-white divide-y divide-gray-200">
          <tr
            v-for="donation in currentDonations"
            :key="donation.id"
            class="hover:bg-gray-50 transition-colors duration-200"
          >
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">{{ formatDate(donation.donatedAt || donation.created_at) }}</div>
              <div class="text-sm text-gray-500 font-mono">#{{ getTransactionId(donation) }}</div>
            </td>
            <td class="px-6 py-4">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div class="flex items-center justify-center w-10 h-10 rounded-full bg-gradient-to-r from-indigo-400 to-indigo-600 text-white">
                    <span class="text-sm font-medium">{{ getInitials(getDonorName(donation)) }}</span>
                  </div>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">{{ getDonorName(donation) }}</div>
                  <div class="text-sm text-gray-500">{{ getDonorEmail(donation) }}</div>
                </div>
              </div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="getDonationTypeClass(donation.donationType || donation.type)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize"
              >
                {{ formatDonationType(donation.donationType || donation.type) }}
              </span>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900 capitalize">{{ formatPaymentMethod(donation.method || donation.paymentMethod) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-lg font-bold text-gray-900">₹{{ formatCurrency(donation.amount || 0) }}</div>
            </td>
            <td class="px-6 py-4 whitespace-nowrap">
              <span
                :class="getStatusClass(donation.status)"
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium capitalize"
              >
                <span class="w-1.5 h-1.5 rounded-full mr-1.5" :class="getStatusDotClass(donation.status)"></span>
                {{ formatStatus(donation.status) }}
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
            <td colspan="6" class="px-6 py-10 text-center">
              <div class="text-gray-400 text-lg mb-2">
                <svg class="mx-auto h-12 w-12" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
              </div>
              <p class="text-gray-500 font-medium">No donations found</p>
              <p class="text-gray-400 text-sm">Try adjusting your filters or search terms</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Pagination -->
    <div class="flex items-center justify-between mt-6 px-4 py-3 bg-gray-50 rounded-lg">
      <div class="text-sm text-gray-700">
        Showing {{ getDisplayRange() }} of {{ totalItems }} donations
      </div>
      <div class="flex items-center space-x-2">
        <button
          :disabled="currentPage === 1 || isLoading"
          @click="handlePageChange(currentPage - 1)"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
        >
          Previous
        </button>
        <div class="flex space-x-1">
          <button
            v-for="page in visiblePages"
            :key="page"
            :class="page === currentPage ? 'bg-indigo-600 text-white' : 'bg-white text-gray-700 hover:bg-gray-50'"
            @click="handlePageChange(page)"
            class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium transition-colors"
          >
            {{ page }}
          </button>
        </div>
        <button
          :disabled="currentPage === totalPages || isLoading"
          @click="handlePageChange(currentPage + 1)"
          class="px-3 py-1 border border-gray-300 rounded-md text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
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
const totalItems = computed(() => donationStore.totalItems || 0)
const totalPages = computed(() => donationStore.totalPages || 1)

const searchQuery = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(20)

const filters = ref({
  dateRange: 'all',
  donationType: 'all',
  status: 'all',
  amountRange: 'all',
})

// Helper functions for robust data mapping
function safeString(val, fallback = '') {
  return typeof val === 'string' && val.trim() !== '' ? val.trim() : fallback;
}

// Donor name mapping – checks every possible location!
function getDonorName(donation) {
  const name =
    safeString(donation.fullName) ||
    (donation.donor && safeString(donation.donor.name)) ||
    safeString(donation.userName) ||
    safeString(donation.donorName) ||
    safeString(donation.user_name) ||
    safeString(donation.name);
  return name || 'Anonymous';
}

function getDonorEmail(donation) {
  const email =
    safeString(donation.emailAddress) ||
    (donation.donor && safeString(donation.donor.email)) ||
    safeString(donation.userEmail) ||
    safeString(donation.donorEmail) ||
    safeString(donation.user_email) ||
    safeString(donation.email);
  return email || 'N/A';
}

function getTransactionId(donation) {
  return safeString(donation.transactionId) ||
    safeString(donation.order_id) ||
    safeString(donation.payment_id) ||
    safeString(donation.paymentId) ||
    `DON-${donation.id}` ||
    'N/A';
}

function getInitials(name) {
  if (!name || name === 'Anonymous' || name === 'N/A') return 'AN';
  return name.trim().split(/\s+/).map(n => n[0]).join('').toUpperCase().substring(0, 2);
}

function formatCurrency(amount) {
  return new Intl.NumberFormat('en-IN').format(Math.floor(amount || 0));
}

function formatDate(dateString) {
  if (!dateString) return 'N/A';
  try {
    return new Date(dateString).toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  } catch {
    return 'Invalid Date';
  }
}

function formatDonationType(type) {
  if (!type) return 'General';
  return type.charAt(0).toUpperCase() + type.slice(1).toLowerCase();
}

function formatPaymentMethod(method) {
  if (!method) return 'Online';
  // Convert backend method names to user-friendly names
  const methodMap = {
    'upi': 'UPI',
    'card': 'Card',
    'netbanking': 'Net Banking',
    'wallet': 'Wallet',
    'pending': 'Pending'
  };
  return methodMap[method.toLowerCase()] || method;
}

function formatStatus(status) {
  if (!status) return 'Pending';
  const normalizedStatus = status.toUpperCase();
  const statusMap = {
    'SUCCESS': 'Success',
    'PENDING': 'Pending',
    'FAILED': 'Failed',
    'CANCELLED': 'Cancelled'
  };
  return statusMap[normalizedStatus] || status.charAt(0).toUpperCase() + status.slice(1).toLowerCase();
}

function getDonationTypeClass(type) {
  const typeClasses = {
    'general': 'bg-gray-100 text-gray-800',
    'seva': 'bg-green-100 text-green-800',
    'festival': 'bg-yellow-100 text-yellow-800',
    'construction': 'bg-blue-100 text-blue-800',
    'annadanam': 'bg-orange-100 text-orange-800',
    'education': 'bg-purple-100 text-purple-800',
    'maintenance': 'bg-red-100 text-red-800',
    'event': 'bg-pink-100 text-pink-800',
  };
  return typeClasses[(type || 'general').toLowerCase()] || typeClasses.general;
}

function getStatusClass(status) {
  if (!status) return 'bg-yellow-100 text-yellow-800';
  const normalizedStatus = status.toLowerCase();
  const statusClasses = {
    'success': 'bg-green-100 text-green-800',
    'pending': 'bg-yellow-100 text-yellow-800',
    'failed': 'bg-red-100 text-red-800',
    'cancelled': 'bg-gray-100 text-gray-600',
  };
  return statusClasses[normalizedStatus] || statusClasses.pending;
}

function getStatusDotClass(status) {
  if (!status) return 'bg-yellow-400';
  const normalizedStatus = status.toLowerCase();
  const dotClasses = {
    'success': 'bg-green-400',
    'pending': 'bg-yellow-400',
    'failed': 'bg-red-400',
    'cancelled': 'bg-gray-400',
  };
  return dotClasses[normalizedStatus] || dotClasses.pending;
}

function getDisplayRange() {
  if (!donations.value.length) return '0-0';
  const start = (currentPage.value - 1) * itemsPerPage.value + 1;
  const end = Math.min(currentPage.value * itemsPerPage.value, totalItems.value);
  return `${start}-${end}`;
}

function clearSearch() {
  searchQuery.value = '';
}

// Current donations from store (server-side filtered and paginated)
const currentDonations = computed(() => donations.value);

// Pagination controls
const visiblePages = computed(() => {
  const pages = [];
  const maxVisible = 5;
  let startPage = Math.max(1, currentPage.value - Math.floor(maxVisible / 2));
  let endPage = Math.min(totalPages.value, startPage + maxVisible - 1);

  if (endPage - startPage + 1 < maxVisible) {
    startPage = Math.max(1, endPage - maxVisible + 1);
  }
  for (let i = startPage; i <= endPage; i++) {
    pages.push(i);
  }
  return pages;
});

const handlePageChange = (page) => {
  if (page < 1 || page > totalPages.value || page === currentPage.value) return;
  currentPage.value = page;
  fetchDonations();
};

// Build filters for API call
const buildApiFilters = () => {
  const apiFilters = {
    page: currentPage.value,
    limit: itemsPerPage.value,
    search: searchQuery.value.trim(),
  };
  if (filters.value.dateRange !== 'all') {
    apiFilters.dateRange = filters.value.dateRange;
  }
  if (filters.value.donationType !== 'all') {
    apiFilters.type = filters.value.donationType;
  }
  if (filters.value.status !== 'all') {
    apiFilters.status = filters.value.status;
  }
  if (filters.value.amountRange !== 'all') {
    const [min, max] = filters.value.amountRange.split('-').map(v => v.replace(/[₹+,]/g, ''));
    if (min) apiFilters.min = parseInt(min);
    if (max && max !== '') apiFilters.max = parseInt(max);
  }
  return apiFilters;
};

let searchTimeout;
watch([filters, searchQuery], () => {
  clearTimeout(searchTimeout);
  searchTimeout = setTimeout(() => {
    currentPage.value = 1;
    fetchDonations();
  }, 300);
}, { deep: true });

// Fetch data functions
const fetchDonations = async () => {
  try {
    const apiFilters = buildApiFilters();
    donationStore.setFilters(apiFilters);
    await donationStore.fetchDonations();
  } catch (error) {
    console.error('Error fetching donations:', error);
  }
};

const fetchDashboard = async () => {
  try {
    await donationStore.fetchDashboard();
  } catch (error) {
    console.error('Error fetching donation dashboard:', error);
  }
};

const refreshData = async () => {
  await Promise.all([fetchDashboard(), fetchDonations()]);
};

onMounted(async () => {
  await refreshData();
});
</script>