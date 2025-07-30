<template>
  <div class="bg-white">
    <!-- Receipt Header -->
    <div class="flex items-center justify-between mb-8">
      <div class="flex items-center space-x-4">
        <div class="w-16 h-16 bg-indigo-100 rounded-xl flex items-center justify-center">
          <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
        </div>
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Donation Receipt</h2>
          <p class="text-sm text-gray-600">Tax-deductible donation receipt</p>
        </div>
      </div>
      <div class="flex space-x-3">
        <button 
          @click="downloadPDF"
          :disabled="isGenerating"
          class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
          </svg>
          {{ isGenerating ? 'Generating...' : 'Download PDF' }}
        </button>
        <button 
          @click="printReceipt"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
        >
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z"/>
          </svg>
          Print Receipt
        </button>
      </div>
    </div>

    <!-- Receipt Content -->
    <div ref="receiptContent" class="max-w-2xl mx-auto bg-white border border-gray-200 rounded-2xl shadow-lg overflow-hidden">
      <!-- Temple Header -->
      <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 px-8 py-6 text-white">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold">{{ receipt.templeName }}</h1>
            <p class="text-indigo-100 mt-1">{{ receipt.templeAddress }}</p>
            <p class="text-indigo-100">Phone: {{ receipt.templePhone }} | Email: {{ receipt.templeEmail }}</p>
          </div>
          <div class="text-right">
            <div class="bg-white/20 rounded-lg px-4 py-2">
              <p class="text-sm font-medium">Receipt #</p>
              <p class="text-xl font-bold">{{ receipt.receiptNumber }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Receipt Details -->
      <div class="px-8 py-6">
        <!-- Donation Information -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6 mb-8">
          <div class="space-y-4">
            <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">Donor Information</h3>
            <div class="space-y-2">
              <div>
                <p class="text-sm text-gray-600">Full Name</p>
                <p class="font-medium text-gray-900">{{ receipt.donorName }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-600">Address</p>
                <p class="font-medium text-gray-900">{{ receipt.donorAddress }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-600">Phone</p>
                <p class="font-medium text-gray-900">{{ receipt.donorPhone }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-600">Email</p>
                <p class="font-medium text-gray-900">{{ receipt.donorEmail }}</p>
              </div>
              <div v-if="receipt.panNumber">
                <p class="text-sm text-gray-600">PAN Number</p>
                <p class="font-medium text-gray-900">{{ receipt.panNumber }}</p>
              </div>
            </div>
          </div>

          <div class="space-y-4">
            <h3 class="text-lg font-semibold text-gray-900 border-b border-gray-200 pb-2">Donation Details</h3>
            <div class="space-y-2">
              <div>
                <p class="text-sm text-gray-600">Date of Donation</p>
                <p class="font-medium text-gray-900">{{ formatDate(receipt.donationDate) }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-600">Payment Method</p>
                <p class="font-medium text-gray-900">{{ receipt.paymentMethod }}</p>
              </div>
              <div v-if="receipt.transactionId">
                <p class="text-sm text-gray-600">Transaction ID</p>
                <p class="font-medium text-gray-900">{{ receipt.transactionId }}</p>
              </div>
              <div>
                <p class="text-sm text-gray-600">Purpose</p>
                <p class="font-medium text-gray-900">{{ receipt.purpose }}</p>
              </div>
              <div v-if="receipt.sevaName">
                <p class="text-sm text-gray-600">Seva Name</p>
                <p class="font-medium text-gray-900">{{ receipt.sevaName }}</p>
              </div>
            </div>
          </div>
        </div>

        <!-- Amount Details -->
        <div class="bg-gray-50 rounded-xl p-6 mb-8">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Amount Details</h3>
          <div class="space-y-3">
            <div class="flex justify-between items-center py-2 border-b border-gray-200">
              <span class="text-gray-600">Donation Amount</span>
              <span class="font-medium text-gray-900">₹{{ formatAmount(receipt.amount) }}</span>
            </div>
            <div v-if="receipt.processingFee" class="flex justify-between items-center py-2 border-b border-gray-200">
              <span class="text-gray-600">Processing Fee</span>
              <span class="font-medium text-gray-900">₹{{ formatAmount(receipt.processingFee) }}</span>
            </div>
            <div class="flex justify-between items-center py-3 bg-indigo-50 rounded-lg px-4">
              <span class="text-lg font-semibold text-indigo-900">Total Amount</span>
              <span class="text-2xl font-bold text-indigo-900">₹{{ formatAmount(receipt.totalAmount) }}</span>
            </div>
          </div>
          
          <!-- Amount in Words -->
          <div class="mt-4 pt-4 border-t border-gray-200">
            <p class="text-sm text-gray-600">Amount in Words:</p>
            <p class="font-medium text-gray-900 capitalize">{{ amountInWords }}</p>
          </div>
        </div>

        <!-- Tax Information -->
        <div class="bg-green-50 border border-green-200 rounded-xl p-6 mb-8">
          <div class="flex items-start space-x-3">
            <div class="flex-shrink-0">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
            </div>
            <div class="flex-1">
              <h4 class="text-lg font-semibold text-green-900 mb-2">Tax Deduction Information</h4>
              <p class="text-sm text-green-800 mb-2">
                This donation is eligible for tax deduction under Section 80G of Income Tax Act, 1961.
              </p>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
                <div>
                  <p class="text-green-700 font-medium">80G Registration Number:</p>
                  <p class="text-green-900">{{ receipt.taxRegistrationNumber }}</p>
                </div>
                <div>
                  <p class="text-green-700 font-medium">Exemption Percentage:</p>
                  <p class="text-green-900">{{ receipt.taxExemptionPercentage }}% deduction available</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="text-center space-y-4 pt-6 border-t border-gray-200">
          <p class="text-sm text-gray-600">
            This is a computer-generated receipt and does not require a signature.
          </p>
          <div class="flex items-center justify-center space-x-2 text-xs text-gray-500">
            <span>Generated on {{ formatDate(new Date()) }}</span>
            <span>•</span>
            <span>Valid for tax purposes</span>
          </div>
          
          <!-- QR Code for verification (placeholder) -->
          <div class="flex justify-center mt-6">
            <div class="w-20 h-20 bg-gray-200 rounded-lg flex items-center justify-center">
              <svg class="w-8 h-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z"/>
              </svg>
            </div>
          </div>
          <p class="text-xs text-gray-500">Scan for receipt verification</p>
        </div>
      </div>
    </div>


  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const props = defineProps({
  donation: {
    type: Object,
    required: true
  }
})

const receiptContent = ref(null)
const isGenerating = ref(false)

// Mock receipt data - in real app, this would come from API
const receipt = ref({
  receiptNumber: 'RCT-2024-001234',
  templeName: 'Sri Venkateswara Temple',
  templeAddress: '123 Temple Street, Bengaluru, Karnataka 560001',
  templePhone: '+91 80 1234 5678',
  templeEmail: 'info@srivenkateswaratemple.org',
  donorName: 'Rajesh Kumar',
  donorAddress: '456 MG Road, Bengaluru, Karnataka 560025',
  donorPhone: '+91 98765 43210',
  donorEmail: 'rajesh.kumar@email.com',
  panNumber: 'ABCDE1234F',
  donationDate: '2024-01-15',
  paymentMethod: 'Online Transfer',
  transactionId: 'TXN123456789',
  purpose: 'General Donation',
  sevaName: 'Abhishekam',
  amount: 5000,
  processingFee: 50,
  totalAmount: 5050,
  taxRegistrationNumber: '12A-0001234-56789',
  taxExemptionPercentage: 50
})

const amountInWords = computed(() => {
  // Simple number to words conversion (you'd use a library in production)
  const amount = receipt.value.totalAmount
  if (amount === 5050) return 'Five Thousand Fifty Only'
  return 'Amount in words conversion needed'
})

const formatDate = (date) => {
  return new Date(date).toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const formatAmount = (amount) => {
  return new Intl.NumberFormat('en-IN').format(amount)
}

const downloadPDF = async () => {
  isGenerating.value = true
  try {
    // In a real app, you'd use a library like jsPDF or html2pdf
    // For now, we'll simulate the download
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Create a mock download
    const link = document.createElement('a')
    link.href = '#'
    link.download = `donation-receipt-${receipt.value.receiptNumber}.pdf`
    link.click()
    
    // Show success message
    // You'd integrate with your toast system here
    console.log('PDF download initiated')
  } catch (error) {
    console.error('Error generating PDF:', error)
  } finally {
    isGenerating.value = false
  }
}

const printReceipt = () => {
  window.print()
}

onMounted(() => {
  // Initialize with prop data if provided
  if (props.donation) {
    receipt.value = { ...receipt.value, ...props.donation }
  }
})
</script>

    <!-- Print Styles -->
    <style scoped>
    @media print {
      body * {
        visibility: hidden;
      }
      .receipt-content, .receipt-content * {
        visibility: visible;
      }
      .receipt-content {
        position: absolute;
        left: 0;
        top: 0;
        width: 100%;
      }
      .no-print {
        display: none !important;
      }
    }
    </style>