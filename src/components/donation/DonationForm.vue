<template>
  <div class="max-w-2xl mx-auto p-6">
    <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 px-6 py-4">
        <h2 class="text-xl font-semibold text-white flex items-center gap-2">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1" />
          </svg>
          Make a Donation
        </h2>
        <p class="text-indigo-100 text-sm mt-1">Support the temple with your generous contribution</p>
      </div>

      <!-- Form Content -->
      <form @submit.prevent="submitDonation" class="p-6 space-y-6">
        <!-- Donation Type -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Donation Type <span class="text-red-500">*</span>
          </label>
          <select 
            v-model="form.donationType" 
            :class="[
              'w-full px-4 py-3 rounded-xl border transition-all duration-200',
              'focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500',
              errors.donationType ? 'border-red-300 bg-red-50' : 'border-gray-300 hover:border-gray-400'
            ]"
            required
          >
            <option value="">Select donation type</option>
            <option value="general">General Donation</option>
            <option value="seva">Seva Donation</option>
            <option value="festival">Festival Donation</option>
            <option value="construction">Construction Fund</option>
            <option value="annadanam">Annadanam</option>
            <option value="education">Education Fund</option>
            <option value="maintenance">Maintenance</option>
          </select>
          <p v-if="errors.donationType" class="text-red-600 text-sm">{{ errors.donationType }}</p>
        </div>

        <!-- Amount -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Donation Amount <span class="text-red-500">*</span>
          </label>
          <div class="relative">
            <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-500 font-medium">₹</span>
            <input
              v-model="form.amount"
              type="number"
              min="1"
              step="0.01"
              placeholder="Enter amount"
              :class="[
                'w-full pl-8 pr-4 py-3 rounded-xl border transition-all duration-200',
                'focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500',
                errors.amount ? 'border-red-300 bg-red-50' : 'border-gray-300 hover:border-gray-400'
              ]"
              required
            />
          </div>
          <p v-if="errors.amount" class="text-red-600 text-sm">{{ errors.amount }}</p>
          
          <!-- Quick Amount Buttons -->
          <div class="flex flex-wrap gap-2 mt-3">
            <button
              v-for="quickAmount in quickAmounts"
              :key="quickAmount"
              type="button"
              @click="form.amount = quickAmount"
              class="px-4 py-2 text-sm rounded-lg border border-indigo-200 text-indigo-700 hover:bg-indigo-50 hover:border-indigo-300 transition-all duration-200"
            >
              ₹{{ quickAmount }}
            </button>
          </div>
        </div>

        <!-- Donor Information -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-900">
              Donor Name <span class="text-red-500">*</span>
            </label>
            <input
              v-model="form.donorName"
              type="text"
              placeholder="Enter donor name"
              :class="[
                'w-full px-4 py-3 rounded-xl border transition-all duration-200',
                'focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500',
                errors.donorName ? 'border-red-300 bg-red-50' : 'border-gray-300 hover:border-gray-400'
              ]"
              required
            />
            <p v-if="errors.donorName" class="text-red-600 text-sm">{{ errors.donorName }}</p>
          </div>

          <div class="space-y-2">
            <label class="block text-sm font-medium text-gray-900">
              Phone Number
            </label>
            <input
              v-model="form.phone"
              type="tel"
              placeholder="Enter phone number"
              class="w-full px-4 py-3 rounded-xl border border-gray-300 hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            />
          </div>
        </div>

        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Email Address
          </label>
          <input
            v-model="form.email"
            type="email"
            placeholder="Enter email address"
            class="w-full px-4 py-3 rounded-xl border border-gray-300 hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          />
        </div>

        <!-- Address -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Address
          </label>
          <textarea
            v-model="form.address"
            rows="3"
            placeholder="Enter donor address"
            class="w-full px-4 py-3 rounded-xl border border-gray-300 hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 resize-none"
          ></textarea>
        </div>

        <!-- Payment Method -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Payment Method <span class="text-red-500">*</span>
          </label>
          <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
            <button
              v-for="method in paymentMethods"
              :key="method.value"
              type="button"
              @click="form.paymentMethod = method.value"
              :class="[
                'p-4 rounded-xl border-2 transition-all duration-200 text-center',
                form.paymentMethod === method.value
                  ? 'border-indigo-500 bg-indigo-50 text-indigo-700'
                  : 'border-gray-200 hover:border-gray-300 text-gray-700'
              ]"
            >
              <div class="flex items-center justify-center mb-2">
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="method.icon" />
                </svg>
              </div>
              <span class="text-sm font-medium">{{ method.label }}</span>
            </button>
          </div>
        </div>

        <!-- Purpose/Notes -->
        <div class="space-y-2">
          <label class="block text-sm font-medium text-gray-900">
            Purpose/Notes
          </label>
          <textarea
            v-model="form.purpose"
            rows="3"
            placeholder="Enter purpose or special notes for this donation"
            class="w-full px-4 py-3 rounded-xl border border-gray-300 hover:border-gray-400 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 resize-none"
          ></textarea>
        </div>

        <!-- Receipt Options -->
        <div class="space-y-3">
          <div class="flex items-center space-x-3">
            <input
              id="needReceipt"
              v-model="form.needReceipt"
              type="checkbox"
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            />
            <label for="needReceipt" class="text-sm text-gray-700">
              I need a donation receipt
            </label>
          </div>
          
          <div class="flex items-center space-x-3">
            <input
              id="anonymous"
              v-model="form.anonymous"
              type="checkbox"
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            />
            <label for="anonymous" class="text-sm text-gray-700">
              Make this donation anonymous
            </label>
          </div>
        </div>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row gap-3 pt-4">
          <button
            type="button"
            @click="resetForm"
            class="flex-1 px-6 py-3 text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-xl font-medium transition-all duration-200"
          >
            Reset Form
          </button>
          <button
            type="submit"
            :disabled="isSubmitting"
            :class="[
              'flex-1 px-6 py-3 rounded-xl font-medium transition-all duration-200',
              'flex items-center justify-center gap-2',
              isSubmitting
                ? 'bg-gray-400 cursor-not-allowed text-white'
                : 'bg-gradient-to-r from-indigo-600 to-indigo-700 hover:from-indigo-700 hover:to-indigo-800 text-white shadow-lg hover:shadow-xl'
            ]"
          >
            <svg v-if="isSubmitting" class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>{{ isSubmitting ? 'Processing...' : 'Proceed to Payment' }}</span>
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script>
import { ref, reactive } from 'vue'

export default {
  name: 'DonationForm',
  emits: ['donation-submitted'],
  setup(props, { emit }) {
    const isSubmitting = ref(false)
    
    const form = reactive({
      donationType: '',
      amount: '',
      donorName: '',
      phone: '',
      email: '',
      address: '',
      paymentMethod: 'online',
      purpose: '',
      needReceipt: true,
      anonymous: false
    })

    const errors = reactive({})

    const quickAmounts = [100, 500, 1000, 2500, 5000, 10000]

    const paymentMethods = [
      {
        value: 'online',
        label: 'Online Payment',
        icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z'
      },
      {
        value: 'cash',
        label: 'Cash',
        icon: 'M17 9V7a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2m2 4h10a2 2 0 002-2v-6a2 2 0 00-2-2H9a2 2 0 00-2 2v6a2 2 0 002 2zm7-5a2 2 0 11-4 0 2 2 0 014 0z'
      },
      {
        value: 'cheque',
        label: 'Cheque',
        icon: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
      }
    ]

    const validateForm = () => {
      // Reset errors
      Object.keys(errors).forEach(key => delete errors[key])

      if (!form.donationType) {
        errors.donationType = 'Please select a donation type'
      }

      if (!form.amount || form.amount <= 0) {
        errors.amount = 'Please enter a valid amount'
      }

      if (!form.donorName.trim()) {
        errors.donorName = 'Please enter donor name'
      }

      return Object.keys(errors).length === 0
    }

    const submitDonation = async () => {
      if (!validateForm()) return

      isSubmitting.value = true

      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 2000))
        
        const donationData = {
          ...form,
          id: Date.now(),
          date: new Date().toISOString(),
          status: 'pending'
        }

        emit('donation-submitted', donationData)
        
        // Success feedback could be added here
        console.log('Donation submitted:', donationData)
        
      } catch (error) {
        console.error('Error submitting donation:', error)
      } finally {
        isSubmitting.value = false
      }
    }

    const resetForm = () => {
      Object.keys(form).forEach(key => {
        if (typeof form[key] === 'boolean') {
          form[key] = key === 'needReceipt' ? true : false
        } else {
          form[key] = ''
        }
      })
      form.paymentMethod = 'online'
      Object.keys(errors).forEach(key => delete errors[key])
    }

    return {
      form,
      errors,
      isSubmitting,
      quickAmounts,
      paymentMethods,
      submitDonation,
      resetForm
    }
  }
}
</script>