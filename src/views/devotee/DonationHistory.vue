<!-- src/views/devotee/DonationHistory.vue -->
<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Page Header -->
    <div class="bg-white border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-6">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">My Donations</h1>
              <p class="mt-1 text-sm text-gray-600">Track your donation history and receipts</p>
            </div>
            <button 
              @click="showDonationForm = true"
              class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200 flex items-center gap-2"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
              </svg>
              Donate Now
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="p-3 bg-indigo-100 rounded-lg">
              <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Donated</p>
              <p class="text-2xl font-bold text-gray-900">₹{{ totalDonated.toLocaleString() }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="p-3 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Donations</p>
              <p class="text-2xl font-bold text-gray-900">{{ donationHistory.length }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <div class="flex items-center">
            <div class="p-3 bg-amber-100 rounded-lg">
              <svg class="w-6 h-6 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3a4 4 0 118 0v4m-4 8a4 4 0 11-8 0v-4m4-8a4 4 0 118 0v4"></path>
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">This Month</p>
              <p class="text-2xl font-bold text-gray-900">₹{{ thisMonthDonated.toLocaleString() }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filters -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
            <select 
              v-model="filters.dateRange" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="all">All Time</option>
              <option value="this-month">This Month</option>
              <option value="last-month">Last Month</option>
              <option value="this-year">This Year</option>
            </select>
          </div>
          
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Donation Type</label>
            <select 
              v-model="filters.donationType" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="all">All Types</option>
              <option value="general">General Donation</option>
              <option value="seva">Seva Donation</option>
              <option value="festival">Festival Donation</option>
              <option value="construction">Construction Fund</option>
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
              <option value="501-1000">₹501 - ₹1,000</option>
              <option value="1001-5000">₹1,001 - ₹5,000</option>
              <option value="5000+">₹5,000+</option>
            </select>
          </div>

          <div class="flex items-end">
            <button 
              @click="resetFilters"
              class="w-full px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors duration-200"
            >
              Reset Filters
            </button>
          </div>
        </div>
      </div>

      <!-- Donation History Table -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-lg font-semibold text-gray-900">Donation History</h3>
          </div>
        </div>

        <div class="overflow-x-auto">
          <table class="min-w-full divide-y divide-gray-200">
            <thead class="bg-gray-50">
              <tr>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Date</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Amount</th>
                <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
              </tr>
            </thead>
            <tbody class="bg-white divide-y divide-gray-200">
              <tr v-for="donation in filteredDonations" :key="donation.id" class="hover:bg-gray-50">
                <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                  {{ formatDate(donation.created_at || donation.donated_at || donation.donation_date) }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getDonationTypeClass(donation.type || donation.donation_type)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ donation.type || donation.donation_type }}
                  </span>
                </td>
                <td class="px-6 py-4 whitespace-nowrap text-sm font-medium text-gray-900">
                  ₹{{ donation.amount.toLocaleString() }}
                </td>
                <td class="px-6 py-4 whitespace-nowrap">
                  <span :class="getStatusClass(donation.status)" class="px-2 py-1 text-xs font-medium rounded-full">
                    {{ donation.status }}
                  </span>
                </td>
              </tr>
            </tbody>
          </table>
        </div>

        <!-- Empty State -->
        <div v-if="filteredDonations.length === 0" class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1"></path>
          </svg>
          <h3 class="mt-2 text-sm font-medium text-gray-900">No donations found</h3>
          <p class="mt-1 text-sm text-gray-500">No donations match your current filters.</p>
        </div>
      </div>
    </div>

    <!-- Donation Form Modal -->
    <div v-if="showDonationForm" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl max-w-md w-full p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-semibold text-gray-900">Make a Donation</h3>
          <button @click="showDonationForm = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        
        <form @submit.prevent="submitDonation" class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Donation Type</label>
            <select 
              v-model="donationForm.type" 
              required
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
              <option value="">Select Type</option>
              <option value="general">General Donation</option>
              <option value="seva">Seva Donation</option>
              <option value="festival">Festival Donation</option>
              <option value="event">Event Donation</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Amount (₹)</label>
            <input 
              type="number" 
              v-model="donationForm.amount" 
              required
              min="1"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="Enter amount"
            >
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Purpose</label>
            <textarea 
              v-model="donationForm.purpose" 
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              placeholder="Purpose of donation (optional)"
            ></textarea>
          </div>

          <div class="flex gap-3 pt-4">
            <button 
              type="button" 
              @click="showDonationForm = false"
              class="flex-1 px-4 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors duration-200"
            >
              Cancel
            </button>
            <button 
              type="submit"
              :disabled="submittingDonation"
              class="flex-1 bg-indigo-600 hover:bg-indigo-700 disabled:opacity-50 text-white px-4 py-2 rounded-lg font-medium transition-colors duration-200"
            >
              {{ submittingDonation ? 'Processing...' : 'Proceed to Payment' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { donationService } from '@/services/donation.service'

const loadRazorpayScript = () => {
  return new Promise((resolve) => {
    const script = document.createElement("script");
    script.src = "https://checkout.razorpay.com/v1/checkout.js";
    script.onload = () => resolve(true);
    script.onerror = () => resolve(false);
    document.body.appendChild(script);
  });
};

// Reactive data
const showDonationForm = ref(false)
const donationHistory = ref([])
const loading = ref(true)
const submittingDonation = ref(false)

// Filters
const filters = ref({
  dateRange: 'all',
  donationType: 'all',
  amountRange: 'all'
})

// Donation form
const donationForm = ref({
  type: '',
  amount: '',
  purpose: ''
})

// Computed properties
const totalDonated = computed(() => {
  return donationHistory.value.reduce((total, donation) => total + donation.amount, 0)
})

const thisMonthDonated = computed(() => {
  const thisMonth = new Date().getMonth()
  const thisYear = new Date().getFullYear()
  
  return donationHistory.value
    .filter(donation => {
      const donationDate = new Date(donation.date || donation.donated_at || donation.donation_date)
      return donationDate.getMonth() === thisMonth && donationDate.getFullYear() === thisYear
    })
    .reduce((total, donation) => total + donation.amount, 0)
})

const filteredDonations = computed(() => {
  let filtered = [...donationHistory.value]

  // Date range filter
  if (filters.value.dateRange !== 'all') {
    const now = new Date()
    filtered = filtered.filter(donation => {
      const donationDate = new Date(donation.date || donation.donated_at || donation.donation_date)
      
      switch (filters.value.dateRange) {
        case 'this-month':
          return donationDate.getMonth() === now.getMonth() && donationDate.getFullYear() === now.getFullYear()
        case 'last-month':
          const lastMonth = new Date(now.getFullYear(), now.getMonth() - 1, 1)
          return donationDate.getMonth() === lastMonth.getMonth() && donationDate.getFullYear() === lastMonth.getFullYear()
        case 'this-year':
          return donationDate.getFullYear() === now.getFullYear()
        default:
          return true
      }
    })
  }

  // Donation type filter
  if (filters.value.donationType !== 'all') {
    filtered = filtered.filter(donation => 
      (donation.type || donation.donation_type) === filters.value.donationType
    )
  }

  // Amount range filter
  if (filters.value.amountRange !== 'all') {
    filtered = filtered.filter(donation => {
      const amount = donation.amount
      switch (filters.value.amountRange) {
        case '0-500':
          return amount >= 0 && amount <= 500
        case '501-1000':
          return amount >= 501 && amount <= 1000
        case '1001-5000':
          return amount >= 1001 && amount <= 5000
        case '5000+':
          return amount > 5000
        default:
          return true
      }
    })
  }

  return filtered.sort((a, b) => {
    const dateA = new Date(a.date || a.donated_at || a.donation_date)
    const dateB = new Date(b.date || b.donated_at || b.donation_date)
    return dateB - dateA
  })
})

// Debug function to analyze response structure
const debugDonationResponse = (response) => {
  console.log('=== DONATION RESPONSE DEBUG ===')
  console.log('Full response object:', response)
  console.log('Response type:', typeof response)
  console.log('Response keys:', Object.keys(response || {}))
  
  // Check for nested data
  if (response.data) {
    console.log('response.data:', response.data)
    console.log('response.data keys:', Object.keys(response.data || {}))
  }
  
  // Check for different field name variations
  const possibleOrderIds = [
    response.order_id,
    response.OrderID,
    response.orderId,
    response.data?.order_id,
    response.data?.OrderID,
    response.data?.orderId
  ].filter(Boolean)
  
  const possibleKeys = [
    response.razorpay_key,
    response.RazorpayKey,
    response.razorpayKey,
    response.key,
    response.data?.razorpay_key,
    response.data?.RazorpayKey,
    response.data?.razorpayKey,
    response.data?.key
  ].filter(Boolean)
  
  console.log('Found possible order_id values:', possibleOrderIds)
  console.log('Found possible razorpay_key values:', possibleKeys)
  
  console.log('JSON stringified response:', JSON.stringify(response, null, 2))
  console.log('=== END DEBUG ===')
}

// Comprehensive response parsing - handles multiple backend response formats
const parseRazorpayResponse = (response) => {
  // The response might be nested in different ways
  const possibleDataSources = [
    response,
    response.data,
    response.result,
    response.payload,
    response.order,
    response.razorpay
  ].filter(source => source && typeof source === 'object')
  
  let order_id, razorpay_key, amount
  
  for (const source of possibleDataSources) {
    // Try different field name variations for order_id
    order_id = order_id || 
      source.order_id || 
      source.OrderID || 
      source.orderId || 
      source.order_id || 
      source.razorpay_order_id ||
      source.RazorpayOrderID ||
      source.orderID ||
      source.id
    
    // Try different field name variations for razorpay_key
    razorpay_key = razorpay_key || 
      source.razorpay_key || 
      source.RazorpayKey || 
      source.razorpayKey || 
      source.key || 
      source.Key ||
      source.api_key ||
      source.ApiKey ||
      source.razorpay_api_key ||
      source.RazorpayApiKey
    
    // Try different field name variations for amount
    amount = amount || 
      source.amount || 
      source.Amount || 
      source.total || 
      source.Total ||
      source.payment_amount ||
      source.PaymentAmount
  }
  
  // Fallback for amount
  amount = amount || Number(donationForm.value.amount)
  
  return { order_id, razorpay_key, amount }
}

// Methods
const formatDate = (date) => {
  if (!date) return 'N/A'
  return new Date(date).toLocaleDateString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const getDonationTypeClass = (type) => {
  const classes = {
    'general': 'bg-blue-100 text-blue-800',
    'seva': 'bg-green-100 text-green-800',
    'festival': 'bg-amber-100 text-amber-800',
    'construction': 'bg-purple-100 text-purple-800',
    'event': 'bg-pink-100 text-pink-800'
  }
  return classes[type] || 'bg-gray-100 text-gray-800'
}

const getStatusClass = (status) => {
  const statusLower = (status || '').toLowerCase()
  const classes = {
    'completed': 'bg-green-100 text-green-800',
    'success': 'bg-green-100 text-green-800',
    'pending': 'bg-amber-100 text-amber-800',
    'processing': 'bg-amber-100 text-amber-800',
    'failed': 'bg-red-100 text-red-800',
    'error': 'bg-red-100 text-red-800'
  }
  return classes[statusLower] || 'bg-gray-100 text-gray-800'
}

const resetFilters = () => {
  filters.value = {
    dateRange: 'all',
    donationType: 'all',
    amountRange: 'all'
  }
}

const submitDonation = async () => {
  if (submittingDonation.value) return
  
  try {
    submittingDonation.value = true

    // Validate form
    if (!donationForm.value.type || !donationForm.value.amount) {
      alert('Please fill in all required fields')
      return
    }

    // Load Razorpay script
    const scriptLoaded = await loadRazorpayScript();
    if (!scriptLoaded) {
      alert("Razorpay SDK failed to load. Please check your internet connection.");
      return;
    }

    // Create donation using service layer
    console.log('Creating donation with data:', {
      amount: Number(donationForm.value.amount),
      donationType: donationForm.value.type,
      note: donationForm.value.purpose,
    })

    const response = await donationService.createDonation({
      amount: Number(donationForm.value.amount),
      donationType: donationForm.value.type,
      note: donationForm.value.purpose,
    })

    // Debug the response structure
    debugDonationResponse(response)

    // Parse response using comprehensive handler
    const { order_id, razorpay_key, amount } = parseRazorpayResponse(response)

    console.log('Extracted values:', { order_id, razorpay_key, amount })

    if (!order_id || !razorpay_key) {
      console.error('Missing required fields in response:', response)
      
      // Show more helpful error message
      const missingFields = []
      if (!order_id) missingFields.push('order_id/OrderID/orderId')
      if (!razorpay_key) missingFields.push('razorpay_key/RazorpayKey/key')
      
      const errorMsg = `Server response missing required fields: ${missingFields.join(', ')}. 
      
Available response fields: ${Object.keys(response || {}).join(', ')}
${response.data ? 'Data fields: ' + Object.keys(response.data).join(', ') : ''}

Please check backend implementation or contact support.`
      
      alert(errorMsg)
      throw new Error(`Server response missing required fields: ${missingFields.join(', ')}`)
    }

    const options = {
      key: razorpay_key,
      amount: amount * 100, // Razorpay expects paise
      currency: "INR",
      name: "Temple Donation",
      description: "Thank you for your generosity!",
      order_id: order_id,
      handler: async function (response) {
        try {
          console.log('Payment successful, verifying...', response)
          
          await donationService.verifyDonation({
            razorpay_payment_id: response.razorpay_payment_id,
            razorpay_order_id: response.razorpay_order_id,
            razorpay_signature: response.razorpay_signature,
          });

          alert("🎉 Donation successful! Thank you for your generosity.");
          showDonationForm.value = false;
          
          // Reset form
          donationForm.value = {
            type: '',
            amount: '',
            purpose: ''
          }
          
          // Refresh donation history
          await fetchDonationHistory();
        } catch (verifyError) {
          console.error("❌ Verification failed:", verifyError);
          alert("Payment verification failed. Please contact support if amount was deducted.");
        }
      },
      prefill: {
        name: localStorage.getItem("user_name") || "Donor",
        email: localStorage.getItem("user_email") || "donor@example.com",
      },
      theme: {
        color: "#6366f1",
      },
      modal: {
        ondismiss: function() {
          console.log('Payment modal dismissed')
        }
      }
    };

    console.log('Opening Razorpay with options:', options)
    const rzp = new window.Razorpay(options);
    rzp.open();
  } catch (error) {
    console.error("❌ Donation Flow Error:", error);
    
    const errorMessage = error.response?.data?.error || 
                        error.response?.data?.message || 
                        error.message || 
                        'Something went wrong while processing your donation'
    
    alert("Error: " + errorMessage);
  } finally {
    submittingDonation.value = false
  }
};

const fetchDonationHistory = async () => {
  try {
    loading.value = true

    console.log('Fetching donation history...')
    const response = await donationService.getMyDonations()
    
    console.log('Donation history response:', response)
    donationHistory.value = Array.isArray(response) ? response : []
  } catch (error) {
    console.error('Error fetching donation history:', error)
    donationHistory.value = []
    
    // Don't show alert for 404 errors (endpoint might not be implemented)
    if (error.response?.status !== 404) {
      alert('Failed to load donation history: ' + (error.response?.data?.error || error.message))
    }
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(() => {
  fetchDonationHistory()
})
</script>