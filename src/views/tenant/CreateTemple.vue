<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200 rounded-2xl">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <nav class="flex items-center space-x-2 text-sm text-gray-500 mb-2">
              <router-link to="/tenant/dashboard" class="hover:text-indigo-600">Dashboard</router-link>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
              <span class="text-gray-900">Create Temple</span>
            </nav>
            <h1 class="text-2xl font-bold text-gray-900">Register New Temple</h1>
            <p class="text-gray-600 mt-1">Fill out the form below to register your temple for approval</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Progress Indicator -->
      <div class="mb-8">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-indigo-600">Step {{ currentStep }} of 3</span>
          <span class="text-sm text-gray-500">{{ Math.round((currentStep / 3) * 100) }}% Complete</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div 
            class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
            :style="{ width: `${(currentStep / 3) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- Form Container -->
      <div class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <form @submit.prevent="handleSubmit">
          <!-- Step 1: Basic Information -->
          <div v-show="currentStep === 1" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Temple Basic Information</h2>
              <p class="text-gray-600">Provide the essential details about your temple</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Temple Name -->
              <div class="md:col-span-2">
                <label for="templeName" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Name *
                </label>
                <input
                  id="templeName"
                  v-model="form.name"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter temple name"
                />
              </div>

              <!-- Deity -->
              <div>
                <label for="deity" class="block text-sm font-medium text-gray-700 mb-2">
                  Main Deity *
                </label>
                <input
                  id="deity"
                  v-model="form.main_deity"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="e.g., Lord Vishnu, Goddess Lakshmi"
                />
              </div>

              <!-- Temple Type -->
              <div>
                <label for="templeType" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Type *
                </label>
                <select
                  id="templeType"
                  v-model="form.temple_type"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                >
                  <option value="">Select temple type</option>
                  <option value="traditional">Traditional</option>
                  <option value="modern">Modern</option>
                  <option value="heritage">Heritage</option>
                  <option value="community">Community</option>
                </select>
              </div>

              <!-- Established Year -->
              <div>
                <label for="establishedYear" class="block text-sm font-medium text-gray-700 mb-2">
                  Established Year *
                </label>
                <input
                  id="establishedYear"
                  v-model="form.established_year"
                  type="number"
                  min="1800"
                  :max="new Date().getFullYear()"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="YYYY"
                />
              </div>

              <!-- Phone -->
              <div>
                <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
                  Contact Phone *
                </label>
                <input
                  id="phone"
                  v-model="form.phone"
                  type="tel"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="+91 XXXXX XXXXX"
                />
              </div>

              <!-- Email -->
              <div class="md:col-span-2">
                <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Email *
                </label>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="temple@example.com"
                />
              </div>

              <!-- Description -->
              <div class="md:col-span-2">
                <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Description
                </label>
                <textarea
                  id="description"
                  v-model="form.description"
                  rows="4"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Brief description of the temple, its history, and significance..."
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Step 2: Address Information -->
          <div v-show="currentStep === 2" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Temple Address</h2>
              <p class="text-gray-600">Provide the complete address details</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Street Address -->
              <div class="md:col-span-2">
                <label for="street_address" class="block text-sm font-medium text-gray-700 mb-2">
                  Street Address *
                </label>
                <input
                  id="street_address"
                  v-model="form.street_address"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter complete street address"
                />
              </div>

              <!-- City -->
              <div>
                <label for="city" class="block text-sm font-medium text-gray-700 mb-2">
                  City *
                </label>
                <input
                  id="city"
                  v-model="form.city"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter city name"
                />
              </div>

              <!-- State -->
              <div>
                <label for="state" class="block text-sm font-medium text-gray-700 mb-2">
                  State *
                </label>
                <select
                  id="state"
                  v-model="form.state"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                >
                  <option value="">Select state</option>
                  <option value="Karnataka">Karnataka</option>
                  <option value="Tamil Nadu">Tamil Nadu</option>
                  <option value="Andhra Pradesh">Andhra Pradesh</option>
                  <option value="Kerala">Kerala</option>
                  <option value="Maharashtra">Maharashtra</option>
                  <option value="Gujarat">Gujarat</option>
                  <option value="Rajasthan">Rajasthan</option>
                  <option value="Uttar Pradesh">Uttar Pradesh</option>
                  <option value="West Bengal">West Bengal</option>
                  <option value="Odisha">Odisha</option>
                  <!-- Add more states as needed -->
                </select>
              </div>

              <!-- District -->
              <div>
                <label for="district" class="block text-sm font-medium text-gray-700 mb-2">
                  District *
                </label>
                <input
                  id="district"
                  v-model="form.district"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter district name"
                />
              </div>

              <!-- PIN Code -->
              <div>
                <label for="pincode" class="block text-sm font-medium text-gray-700 mb-2">
                  PIN Code *
                </label>
                <input
                  id="pincode"
                  v-model="form.pincode"
                  type="text"
                  pattern="[0-9]{6}"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter 6-digit PIN code"
                />
              </div>

              <!-- Landmark -->
              <div class="md:col-span-2">
                <label for="landmark" class="block text-sm font-medium text-gray-700 mb-2">
                  Landmark (Optional)
                </label>
                <input
                  id="landmark"
                  v-model="form.landmark"
                  type="text"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Nearby landmark for easy identification"
                />
              </div>

              <!-- Google Maps Link -->
              <div class="md:col-span-2">
                <label for="map_link" class="block text-sm font-medium text-gray-700 mb-2">
                  Google Maps Link (Optional)
                </label>
                <input
                  id="map_link"
                  v-model="form.map_link"
                  type="url"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="https://maps.google.com/..."
                />
              </div>
            </div>
          </div>

          <!-- Step 3: Documents Upload -->
          <div v-show="currentStep === 3" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Upload Documents</h2>
              <p class="text-gray-600">Upload required documents for temple verification</p>
            </div>

            <div class="space-y-6">
              <!-- Registration Certificate -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Registration Certificate *</h3>
                  <p class="text-gray-600 mb-4">Upload temple registration certificate</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200">
                      <span>Choose File</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'registration')" />
                    </label>
                  </div>
                  <p v-if="form.documents.registration" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.registration.name }}
                  </p>
                </div>
              </div>

              <!-- Trust Deed -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Trust Deed *</h3>
                  <p class="text-gray-600 mb-4">Upload temple trust deed document</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200">
                      <span>Choose File</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'trustDeed')" />
                    </label>
                  </div>
                  <p v-if="form.documents.trustDeed" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.trustDeed.name }}
                  </p>
                </div>
              </div>

              <!-- Property Documents -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Property Documents</h3>
                  <p class="text-gray-600 mb-4">Upload property ownership documents (Optional)</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-gray-600 text-white px-6 py-2 rounded-lg hover:bg-gray-700 transition-colors duration-200">
                      <span>Choose File</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'property')" />
                    </label>
                  </div>
                  <p v-if="form.documents.property" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.property.name }}
                  </p>
                </div>
              </div>

              <!-- Additional Documents -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Additional Documents</h3>
                  <p class="text-gray-600 mb-4">Upload any additional supporting documents (Optional)</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-gray-600 text-white px-6 py-2 rounded-lg hover:bg-gray-700 transition-colors duration-200">
                      <span>Choose Files</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" multiple @change="handleFileUpload($event, 'additional')" />
                    </label>
                  </div>
                  <div v-if="form.documents.additional && form.documents.additional.length > 0" class="mt-2">
                    <p v-for="(file, index) in form.documents.additional" :key="index" class="text-sm text-green-600">
                      ✓ {{ file.name }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Terms and Conditions -->
            <div class="mt-8 p-4 bg-gray-50 rounded-lg">
              <label class="flex items-start">
                <input
                  v-model="form.acceptTerms"
                  type="checkbox"
                  class="mt-1 mr-3 h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  required
                />
                <span class="text-sm text-gray-700">
                  I agree to the 
                  <router-link to="/terms" class="text-indigo-600 hover:text-indigo-800 underline">Terms and Conditions</router-link>
                  and 
                  <router-link to="/privacy" class="text-indigo-600 hover:text-indigo-800 underline">Privacy Policy</router-link>.
                  I confirm that all information provided is accurate and complete.
                </span>
              </label>
            </div>
          </div>

          <!-- Navigation Buttons -->
          <div class="px-8 py-6 bg-gray-50 border-t border-gray-200 flex justify-between">
            <button
              v-if="currentStep > 1"
              type="button"
              @click="previousStep"
              class="inline-flex items-center px-6 py-3 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors duration-200"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
              Previous
            </button>
            <div v-else></div>

            <div class="flex space-x-3">
              <router-link
                to="/tenant/dashboard"
                class="inline-flex items-center px-6 py-3 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors duration-200"
              >
                Cancel
              </router-link>

              <button
                v-if="currentStep < 3"
                type="button"
                @click="nextStep"
                class="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200"
              >
                Next
                <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </button>

              <button
                v-else
                type="submit"
                :disabled="isSubmitting"
                class="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
              >
                <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ isSubmitting ? 'Submitting...' : 'Submit Temple Registration' }}
              </button>
            </div>
          </div>
        </form>
      </div>

      <!-- Help Section -->
      <div class="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
        <div class="flex items-start">
          <svg class="w-6 h-6 text-blue-600 mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <div>
            <h3 class="text-sm font-medium text-blue-800 mb-1">Need Help?</h3>
            <p class="text-sm text-blue-700">
              If you need assistance with temple registration, please contact our support team at 
              <a href="mailto:support@templesaas.com" class="underline">support@templesaas.com</a> 
              or call us at <a href="tel:+91-800-123-4567" class="underline">+91-800-123-4567</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useTempleStore } from '@/stores/temple'
import { useToast } from '@/composables/useToast'
import { validators } from '@/utils/validators'
import api from '@/plugins/axios'

const router = useRouter()
const templeStore = useTempleStore()
const toast = useToast()

// Form state
const currentStep = ref(1)
const isSubmitting = ref(false)
const error = ref(null)

// Form data with snake_case field names to EXACTLY match backend
const form = reactive({
  name: '',
  main_deity: '',
  temple_type: '',
  established_year: '',
  phone: '',
  email: '',
  description: '',
  street_address: '',
  city: '',
  state: '',
  district: '',
  pincode: '',
  landmark: '',
  map_link: '',
  documents: {
    registration: null,
    trustDeed: null,
    property: null,
    additional: []
  },
  acceptTerms: false
})

// Methods
const nextStep = () => {
  if (validateCurrentStep()) {
    currentStep.value++
  }
}

const previousStep = () => {
  currentStep.value--
}

const validateCurrentStep = () => {
  if (currentStep.value === 1) {
    // Check for all required fields in step 1
    const requiredFields = ['name', 'main_deity', 'temple_type', 'phone', 'email', 'established_year']
    const missingFields = requiredFields.filter(field => !form[field])
    
    if (missingFields.length > 0) {
      toast.error(`Please fill in all required fields: ${missingFields.join(', ')}`)
      return false
    }
    
    // Validate establishedYear as a number
    if (form.established_year && isNaN(parseInt(form.established_year))) {
      toast.error('Established Year must be a valid number')
      return false
    }
    
    // Validate phone number format
    if (form.phone && !/^[0-9]{10}$/.test(form.phone.replace(/[^0-9]/g, ''))) {
      toast.error('Please enter a valid 10-digit phone number')
      return false
    }
    
    // Validate email format
    if (form.email && !validators.email(form.email)) {
      toast.error('Please enter a valid email address')
      return false
    }
  } else if (currentStep.value === 2) {
    const requiredFields = ['street_address', 'city', 'state', 'district', 'pincode']
    const missingFields = requiredFields.filter(field => !form[field])
    
    if (missingFields.length > 0) {
      toast.error(`Please fill in all required address fields: ${missingFields.join(', ')}`)
      return false
    }
    
    // Specifically validate street address is not empty
    if (!form.street_address || form.street_address.trim() === '') {
      toast.error('Street address is required')
      return false
    }
    
    // Validate PIN code format
    if (form.pincode && !/^\d{6}$/.test(form.pincode)) {
      toast.error('PIN Code must be a 6-digit number')
      return false
    }
  }
  
  return true
}

const handleFileUpload = (event, type) => {
  const files = Array.from(event.target.files)
  
  if (type === 'additional') {
    form.documents.additional = files
  } else {
    form.documents[type] = files[0]
  }
}

const handleSubmit = async () => {
  if (!validateCurrentStep()) return

  try {
    isSubmitting.value = true

    // CRITICAL CHANGE: We'll now directly create the payload and call the API without going through the temple service
    // This will help us ensure exactly what data is being sent
    
    console.log('🔍 FORM DATA AT SUBMISSION:')
    console.log('Street address value:', form.street_address)
    console.log('Main deity value:', form.main_deity)
    console.log('Temple type value:', form.temple_type)

    // Create the exact payload that matches the backend fields
    const payload = {
      name: form.name || "",
      main_deity: form.main_deity || "",
      temple_type: form.temple_type || "",
      established_year: parseInt(form.established_year || 0),
      phone: form.phone ? form.phone.replace(/[^0-9]/g, '') : "",
      email: form.email || "",
      description: form.description || '',
      street_address: form.street_address || '', // CRITICAL field we're having issues with
      city: form.city || "",
      district: form.district || "",
      state: form.state || "",
      pincode: form.pincode || "",
      landmark: form.landmark || '',
      map_link: form.map_link || '',
      accepted_terms: true,
      status: 'pending'
    }

    console.log('📦 Final direct payload:', payload)
    console.log('🧪 Checking street_address once more:', payload.street_address)

    // Make a direct API call to create the temple
    const response = await api.post('/v1/entities', payload)
    
    console.log('✅ Temple created successfully:', response.data)

    toast.success('Temple registration submitted successfully!')
    router.push('/tenant/dashboard')
  } catch (err) {
    error.value = err.response?.data?.error || 'Failed to submit temple registration. Please try again.'
    toast.error(error.value)
    console.error('Temple registration failed:', err)
    console.error('Error details:', err.response?.data || err.message)
  } finally {
    isSubmitting.value = false
  }
}
</script>