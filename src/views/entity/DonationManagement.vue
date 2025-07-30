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
         
          <!-- <button
            @click="exportReport"
            class="bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 px-4 py-2 rounded-xl font-medium transition-all duration-200"
          >
            Export Report
          </button> -->
           <!-- <button
            @click="showAddDonationModal = true"
            class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-xl font-medium transition-all duration-200"
          >
            Add Donation
          </button> -->
        </div>
      </div>
    </div>

    <!-- Stats -->
    <!-- <div class="px-6 pt--20 pb-2">
      <DonationStats />
    </div> -->

    <!-- Donation Ledger Section -->
    <div class="px-6 pb-10">
      <div class="bg-white rounded-2xl shadow-md border border-gray-200">
        <!-- <div class="p-6"> -->
          <DonationLedger
            :date-range="selectedDateRange"
            :search-query="searchQuery"
            @receipt-download="handleReceiptDownload"
          />
        </div>
      </div>
    </div>

    <!-- Add Donation Modal -->
    <div v-if="showAddDonationModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="showAddDonationModal = false">
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-900">Add New Donation</h3>
          <button @click="showAddDonationModal = false" class="text-gray-400 hover:text-gray-600 transition-colors duration-200">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-6">
          <DonationForm @submitted="handleDonationSubmitted" @cancelled="showAddDonationModal = false" />
        </div>
      </div>
    </div>

    <!-- Receipt Modal -->
    <div v-if="showReceiptModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4" @click.self="showReceiptModal = false">
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h3 class="text-xl font-semibold text-gray-900">Donation Receipt</h3>
          <button @click="showReceiptModal = false" class="text-gray-400 hover:text-gray-600 transition-colors duration-200">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        <div class="p-6">
          <DonationReceipt :donation="selectedDonation" />
        </div>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <div class="animate-spin rounded-full h-6 w-6 border-2 border-indigo-600 border-t-transparent"></div>
        <span class="text-gray-700">Processing...</span>
      </div>
    </div>
  <!-- </div> -->
</template>

<script>
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import DonationStats from '@/components/donation/DonationStats.vue'
import DonationLedger from '@/components/donation/DonationLedger.vue'
import DonationForm from '@/components/donation/DonationForm.vue'
import DonationReceipt from '@/components/donation/DonationReceipt.vue'
import { useToast } from '@/composables/useToast'
import {useDonationStore} from '@/stores/donation.js'


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

    const loading = ref(false)
    const showAddDonationModal = ref(false)
    const showReceiptModal = ref(false)
    const selectedDateRange = ref('month')
    const searchQuery = ref('')
    const selectedDonation = ref(null)

    const handleDonationSubmitted = (donation) => {
      showAddDonationModal.value = false
      showToast('Donation added successfully', 'success')
      refreshData()
    }

    const handleReceiptDownload = (donation) => {
      selectedDonation.value = donation
      showReceiptModal.value = true
    }

    const exportReport = async () => {
      loading.value = true
      try {
        await new Promise(resolve => setTimeout(resolve, 1500))
        showToast('Report exported successfully', 'success')
      } catch {
        showToast('Failed to export report', 'error')
      } finally {
        loading.value = false
      }
    }

    const refreshData = () => {
      // You can add data refetch logic here
    }

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
      handleDonationSubmitted,
      handleReceiptDownload,
      exportReport,
      refreshData
    }
  }
}
</script>
