<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Page Header -->
    <div class="bg-white border-b border-gray-200 px-6 py-4">
      <div class="flex items-center justify-between flex-wrap gap-y-2">
        <div>
          <h1 class="text-2xl font-bold text-indigo-700">Donation Management</h1>
          <p class="text-gray-600 mt-1">Track donations, generate receipts, and view financial reports</p>
        </div>
        <div class="flex flex-wrap items-center gap-3">
          <!-- Uncomment to enable export and add donation buttons as needed -->
          <!--
          <button
            @click="exportReport"
            class="bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 px-4 py-2 rounded-xl font-medium transition-all duration-200"
          >
            Export Report
          </button>
          <button
            @click="showAddDonationModal = true"
            class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-xl font-medium transition-all duration-200"
          >
            Add Donation
          </button>
          -->
        </div>
      </div>
    </div>

    <!-- Stats Section - Uncomment if needing stats component -->
    <!--
    <div class="px-6 pt--20 pb-2">
      <DonationStats />
    </div>
    -->

    <!-- Donation Ledger Section -->
    <div class="px-6 pb-10">
      <div class="bg-white rounded-2xl shadow-md border border-gray-200">
        <DonationLedger
          :donations="donations"
          :date-range="selectedDateRange"
          :search-query="searchQuery"
          @receipt-download="handleReceiptDownload"
        />
      </div>
    </div>

    <!-- Add Donation Modal -->
    <div
      v-if="showAddDonationModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      @click.self="showAddDonationModal = false"
    >
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-900">Add New Donation</h3>
          <button
            @click="showAddDonationModal = false"
            class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
            aria-label="Close Add Donation Modal"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-6">
          <DonationForm
            @submitted="handleDonationSubmitted"
            @cancelled="showAddDonationModal = false"
          />
        </div>
      </div>
    </div>

    <!-- Receipt Modal -->
    <div
      v-if="showReceiptModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4"
      @click.self="showReceiptModal = false"
    >
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-900">Donation Receipt</h3>
          <button
            @click="showReceiptModal = false"
            class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
            aria-label="Close Receipt Modal"
          >
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24" >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-6">
          <DonationReceipt :donation="selectedDonation" />
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div
      v-if="loading"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      aria-live="assertive" aria-busy="true"
    >
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <div class="animate-spin rounded-full h-6 w-6 border-2 border-indigo-600 border-t-transparent"></div>
        <span class="text-gray-700">Processing...</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, computed, watch } from 'vue'
import { useRoute } from 'vue-router'
import DonationStats from '@/components/donation/DonationStats.vue'
import DonationLedger from '@/components/donation/DonationLedger.vue'
import DonationForm from '@/components/donation/DonationForm.vue'
import DonationReceipt from '@/components/donation/DonationReceipt.vue'
import { useToast } from '@/composables/useToast'
import donationService from '@/services/donation.service'
import { useDonationStore } from '@/stores/donation'

export default {
  name: 'DonationManagement',
  components: {
    DonationStats,
    DonationLedger,
    DonationForm,
    DonationReceipt
  },
  setup() {
    const route = useRoute()
    const { showToast } = useToast()
    const donationStore = useDonationStore()

    const loading = ref(false)
    const showAddDonationModal = ref(false)
    const showReceiptModal = ref(false)
    const selectedDateRange = ref('month')
    const searchQuery = ref('')
    const selectedDonation = ref(null)
    const entityId = computed(() => route.params.id)
    
    // This is a lookup table for donor information based on user ID
    const userLookup = {
      1: { name: 'Super Admin', email: 'superadmin@example.com' },
      2: { name: 'Naresh V', email: 'naresh12@gmail.com' },
      3: { name: 'Dileep', email: 'dileep12@gmail.com' }
    }
    
    // Get donor information by user ID
    const getDonorInfo = (userId) => {
      return userLookup[userId] || { name: `User ${userId}`, email: '' }
    }
    
    // Called after successfully adding a donation
    const handleDonationSubmitted = (donation) => {
      showAddDonationModal.value = false
      showToast('Donation added successfully', 'success')
      refreshData()
    }

    // Load donations list from the service but PATCH the store data
    const loadDonations = async () => {
      try {
        loading.value = true
        console.log('Starting donation fetch...')

        // First, fetch donations using the store
        await donationStore.fetchDonations()
        
        // Get the donations from the store
        const originalDonations = donationStore.donations || []
        console.log('Original store donations:', originalDonations)
        
        // Create transformed donations with proper names
        if (originalDonations.length > 0) {
  donationStore.$patch({
    donations: originalDonations
  })
        }
        
      } catch (error) {
        console.error('Error loading donations:', error)
        showToast('Failed to load donations', 'error')
      } finally {
        loading.value = false
      }
    }

    // Load donation statistics for the given entity ID
    const loadStats = async () => {
      try {
        await donationStore.fetchDashboard()
      } catch (error) {
        console.error('Error loading donation stats:', error)
      }
    }

    // Handle receipt generation and show modal
    const handleReceiptDownload = async (donation) => {
      try {
        loading.value = true
        selectedDonation.value = donation
        await donationService.generateReceipt(donation.id)
        showReceiptModal.value = true
      } catch (error) {
        showToast('Failed to generate receipt', 'error')
      } finally {
        loading.value = false
      }
    }

    // Export report action (currently simulated)
    const exportReport = async () => {
      loading.value = true
      try {
        await new Promise((resolve) => setTimeout(resolve, 1500))
        showToast('Report exported successfully', 'success')
      } catch (error) {
        showToast('Failed to export report', 'error')
      } finally {
        loading.value = false
      }
    }

    // Reload donations and stats data
    const refreshData = () => {
      loadDonations()
      loadStats()
    }

    // Load initial data on component mount
    onMounted(() => {
      refreshData()
    })

    return {
      loading,
      showAddDonationModal,
      showReceiptModal,
      selectedDateRange,
      searchQuery,
      selectedDonation,
      donations: computed(() => donationStore.donations), // Use the store's data directly
      handleDonationSubmitted,
      handleReceiptDownload,
      exportReport,
      refreshData,
    }
  },
}
</script>