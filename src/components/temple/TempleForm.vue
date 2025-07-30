<template>
  <div class="max-w-4xl mx-auto p-6">
    <div class="bg-white rounded-2xl shadow-lg overflow-hidden">
      <!-- Header -->
      <div class="bg-gradient-to-r from-indigo-600 to-indigo-700 px-8 py-6">
        <h2 class="text-2xl font-bold text-white">
          {{ isEditing ? 'Edit Temple' : 'Create New Temple' }}
        </h2>
        <p class="text-indigo-100 mt-2">
          {{ isEditing ? 'Update temple information' : 'Fill in the details to register your temple' }}
        </p>
      </div>

      <form @submit.prevent="handleSubmit" class="p-8 space-y-8">
        <!-- Basic Information -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              Basic Information
            </h3>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Temple Name *
              </label>
              <input
                v-model="form.name"
                type="text"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Enter temple name"
              />
            </div>

            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Description
              </label>
              <textarea
                v-model="form.description"
                rows="3"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
                placeholder="Brief description of the temple"
              ></textarea>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Category *
              </label>
              <select
                v-model="form.category"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
                <option value="">Select Category</option>
                <option value="Vishnu Temple">Vishnu Temple</option>
                <option value="Shiva Temple">Shiva Temple</option>
                <option value="Devi Temple">Devi Temple</option>
                <option value="Hanuman Temple">Hanuman Temple</option>
                <option value="Ganesha Temple">Ganesha Temple</option>
                <option value="Spiritual Center">Spiritual Center</option>
                <option value="Other">Other</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Presiding Deity *
              </label>
              <input
                v-model="form.presiding_deity"
                type="text"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Main deity of the temple"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Established Year *
              </label>
              <input
                v-model="form.establishedYear"
                type="number"
                required
                min="1800"
                :max="new Date().getFullYear()"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="e.g., 1950"
              />
            </div>
          </div>
        </div>

        <!-- Address Information -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
              </svg>
              Address Information
            </h3>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Street Address *
              </label>
              <input
                v-model="form.address.street"
                type="text"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="House/Plot number, Street name"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                City *
              </label>
              <input
                v-model="form.address.city"
                type="text"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Enter city"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                State *
              </label>
              <select
                v-model="form.address.state"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
                <option value="">Select State</option>
                <option value="Andhra Pradesh">Andhra Pradesh</option>
                <option value="Karnataka">Karnataka</option>
                <option value="Tamil Nadu">Tamil Nadu</option>
                <option value="Kerala">Kerala</option>
                <option value="Telangana">Telangana</option>
                <option value="Maharashtra">Maharashtra</option>
                <option value="Gujarat">Gujarat</option>
                <option value="Rajasthan">Rajasthan</option>
                <option value="Uttar Pradesh">Uttar Pradesh</option>
                <option value="West Bengal">West Bengal</option>
                <!-- Add more states as needed -->
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                District *
              </label>
              <input
                v-model="form.address.district"
                type="text"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Enter district"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Pincode *
              </label>
              <input
                v-model="form.address.pincode"
                type="text"
                required
                pattern="[0-9]{6}"
                maxlength="6"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="6-digit pincode"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Country
              </label>
              <input
                v-model="form.address.country"
                type="text"
                readonly
                class="w-full px-4 py-3 border border-gray-300 rounded-xl bg-gray-50 text-gray-500"
                value="India"
              />
            </div>
          </div>
        </div>

        <!-- Contact Information -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"/>
              </svg>
              Contact Information
            </h3>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Phone Number *
              </label>
              <input
                v-model="form.contact.phone"
                type="tel"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="+91-XXXXXXXXXX"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Email Address *
              </label>
              <input
                v-model="form.contact.email"
                type="email"
                required
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="temple@example.com"
              />
            </div>

            <div class="md:col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Website (Optional)
              </label>
              <input
                v-model="form.contact.website"
                type="url"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="https://yourtemple.com"
              />
            </div>
          </div>
        </div>

        <!-- Temple Timings -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
              </svg>
              Temple Timings
            </h3>
          </div>

          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Morning Open
              </label>
              <input
                v-model="form.timings.morning_open"
                type="time"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Morning Close
              </label>
              <input
                v-model="form.timings.morning_close"
                type="time"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Evening Open
              </label>
              <input
                v-model="form.timings.evening_open"
                type="time"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">
                Evening Close
              </label>
              <input
                v-model="form.timings.evening_close"
                type="time"
                class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              />
            </div>
          </div>
        </div>

        <!-- Facilities -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
              </svg>
              Facilities Available
            </h3>
          </div>

          <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
            <label v-for="facility in facilityOptions" :key="facility" class="flex items-center space-x-3 cursor-pointer">
              <input
                type="checkbox"
                :value="facility"
                v-model="form.facilities"
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              />
              <span class="text-sm text-gray-700">{{ facility }}</span>
            </label>
          </div>
        </div>

        <!-- Document Upload Section -->
        <div class="space-y-6">
          <div class="border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
              <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 21h10a2 2 0 002-2V9.414a1 1 0 00-.293-.707l-5.414-5.414A1 1 0 0012.586 3H7a2 2 0 00-2 2v14a2 2 0 002 2z"/>
              </svg>
              Document Upload
            </h3>
            <p class="text-sm text-gray-600 mt-2">
              Upload required documents for temple verification
            </p>
          </div>

          <div class="border-2 border-dashed border-gray-300 rounded-xl p-8 text-center hover:border-indigo-400 transition-colors">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 48 48">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"/>
            </svg>
            <div class="mt-4">
              <label class="cursor-pointer">
                <span class="text-sm font-medium text-indigo-600 hover:text-indigo-500">
                  Click to upload documents
                </span>
                <input type="file" multiple accept=".pdf,.jpg,.jpeg,.png" class="hidden" @change="handleFileUpload" />
              </label>
              <p class="text-xs text-gray-500 mt-1">
                PDF, JPG, PNG up to 10MB each
              </p>
            </div>
          </div>

          <!-- Uploaded Files List -->
          <div v-if="uploadedFiles.length > 0" class="space-y-2">
            <h4 class="text-sm font-medium text-gray-700">Uploaded Documents:</h4>
            <div class="space-y-2">
              <div v-for="(file, index) in uploadedFiles" :key="index" 
                   class="flex items-center justify-between p-3 bg-gray-50 rounded-lg">
                <div class="flex items-center space-x-3">
                  <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <span class="text-sm text-gray-700">{{ file.name }}</span>
                </div>
                <button @click="removeFile(index)" type="button" 
                        class="text-red-500 hover:text-red-700 transition-colors">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M9 7h6m-6 0V4a1 1 0 011-1h4a1 1 0 011 1v3"/>
                  </svg>
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Form Actions -->
        <div class="flex flex-col sm:flex-row justify-end gap-4 pt-6 border-t border-gray-200">
          <button
            type="button"
            @click="$emit('cancel')"
            class="px-6 py-3 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-xl hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all"
          >
            Cancel
          </button>
          <button
            type="submit"
            :disabled="loading"
            class="px-8 py-3 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-xl hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all flex items-center gap-2"
          >
            <svg v-if="loading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ loading ? 'Saving...' : (isEditing ? 'Update Temple' : 'Create Temple') }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useTempleStore } from '@/stores/temple'
import { useToast } from '@/composables/useToast'
import { validators } from '@/utils/validators'

// Props & Emits
const props = defineProps({
  temple: {
    type: Object,
    default: null
  }
})

const emit = defineEmits(['submit', 'cancel'])

// Store
const templeStore = useTempleStore()
const { showToast } = useToast()

// Reactive data
const loading = computed(() => templeStore.loading)
const error = computed(() => templeStore.error)
const uploadedFiles = ref([])

const facilityOptions = [
  'Parking Available',
  'Wheelchair Accessible',
  'Audio System',
  'Photography Allowed',
  'Prasadam Counter',
  'Donation Counter',
  'Rest Rooms',
  'Drinking Water',
  'Garden/Park',
  'Guest House',
  'Library',
  'Community Hall'
]

// Form data
const form = ref({
  name: '',
  description: '',
  address: {
    street: '',
    city: '',
    state: '',
    district: '',
    pincode: '',
    country: 'India'
  },
  contact: {
    phone: '',
    email: '',
    website: ''
  },
  category: '',
  establishedYear: '',
  presiding_deity: '',
  timings: {
    morning_open: '',
    morning_close: '',
    evening_open: '',
    evening_close: ''
  },
  facilities: [],
  documents: []
})

// Computed
const isEditing = computed(() => !!props.temple)

// Methods
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files)
  files.forEach(file => {
    if (file.size <= 10 * 1024 * 1024) { // 10MB limit
      uploadedFiles.value.push(file)
    } else {
      showToast('File size exceeds 10MB limit', 'error')
    }
  })
}

const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
}

const validateForm = () => {
  // Required fields validation
  if (!form.value.name.trim()) {
    showToast('Temple name is required', 'error')
    return false
  }
  
  if (!form.value.category) {
    showToast('Temple category/type is required', 'error')
    return false
  }
  
  if (!form.value.establishedYear) {
    showToast('Established year is required', 'error')
    return false
  }
  
  if (!form.value.address.state) {
    showToast('State is required', 'error')
    return false
  }
  
  if (!form.value.address.district) {
    showToast('District is required', 'error')
    return false
  }
  
  if (!form.value.contact.phone) {
    showToast('Phone number is required', 'error')
    return false
  }
  
  return true
}

const handleSubmit = async () => {
  try {
    if (!validateForm()) {
      return
    }
    
    // Prepare form data for API
    const formData = new FormData()
    
    // Convert form data to a format the backend API expects
    const apiData = {
      name: form.value.name.trim(),
      description: form.value.description,
      email: form.value.contact.email,
      phone: form.value.contact.phone,
      addressLine1: form.value.address.street,
      city: form.value.address.city,
      state: form.value.address.state,
      district: form.value.address.district,
      pincode: form.value.address.pincode,
      country: 'India',
      templeType: form.value.category, // Ensure this matches the backend field name
      contactPerson: '',
      // Convert establishedYear to a proper number (not a date string)
      establishedYear: form.value.establishedYear ? parseInt(form.value.establishedYear) : null,
      main_deity: form.value.presiding_deity || null
    }
    
    // Important: make sure these field names match exactly what the backend expects
    
    // Add each field as a separate entry in FormData
    Object.keys(apiData).forEach(key => {
      // Skip null values but include empty strings and zeros
      if (apiData[key] !== null && apiData[key] !== undefined) {
        formData.append(key, apiData[key])
      }
    })
    
    // Add documents
    uploadedFiles.value.forEach((file, index) => {
      if (index === 0) {
        formData.append('registration', file)
      } else if (index === 1) {
        formData.append('trustDeed', file)
      } else if (index === 2) {
        formData.append('property', file)
      } else {
        formData.append(`additional_${index-3}`, file)
      }
    })
    
    // Submit to API
    let result
    
    if (isEditing.value && props.temple) {
      result = await templeStore.updateTemple(props.temple.id, formData)
    } else {
      result = await templeStore.createTemple(formData)
    }
    
    // Emit success event
    emit('submit', result)
    showToast('Temple saved successfully', 'success')
    
  } catch (error) {
    console.error('Error submitting form:', error)
    showToast(error.message || 'Failed to save temple', 'error')
  }
}

// Initialize form with temple data if editing
onMounted(() => {
  if (props.temple) {
    form.value = {
      name: props.temple.name || '',
      description: props.temple.description || '',
      address: {
        street: props.temple.addressLine1 || '',
        city: props.temple.city || '',
        state: props.temple.state || '',
        district: props.temple.district || '',
        pincode: props.temple.pincode || '',
        country: props.temple.country || 'India'
      },
      contact: {
        phone: props.temple.phone || '',
        email: props.temple.email || '',
        website: ''  // Not in API but in UI
      },
      category: props.temple.templeType || '',
      establishedYear: props.temple.establishedYear || '',
      presiding_deity: props.temple.main_deity || '',
      timings: {
        morning_open: '',
        morning_close: '',
        evening_open: '',
        evening_close: ''
      },
      facilities: props.temple.facilities || [],
      documents: []
    }
  }
})

// Expose methods for parent component
defineExpose({
  resetForm: () => {
    form.value = {
      name: '',
      description: '',
      address: {
        street: '',
        city: '',
        state: '',
        district: '',
        pincode: '',
        country: 'India'
      },
      contact: {
        phone: '',
        email: '',
        website: ''
      },
      category: '',
      establishedYear: '',
      presiding_deity: '',
      timings: {
        morning_open: '',
        morning_close: '',
        evening_open: '',
        evening_close: ''
      },
      facilities: [],
      documents: []
    }
    uploadedFiles.value = []
  }
})
</script>