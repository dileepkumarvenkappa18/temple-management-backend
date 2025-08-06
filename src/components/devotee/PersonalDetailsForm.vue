<template>
  <div class="bg-white rounded-2xl shadow-md p-6">
    <!-- Header -->
    <div class="mb-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-2">Personal Details</h2>
      <p class="text-gray-600">Let us know more about you to personalize your temple experience</p>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Full Name -->
      <div>
        <label for="fullName" class="block text-sm font-medium text-gray-700 mb-2">
          Full Name <span class="text-red-500">*</span>
        </label>
        <input
          id="fullName"
          v-model="formData.fullName"
          type="text"
          required
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          placeholder="Enter your full name"
        />
        <p v-if="errors.fullName" class="mt-1 text-sm text-red-600">{{ errors.fullName }}</p>
      </div>

      <!-- Date of Birth & Gender -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label for="dateOfBirth" class="block text-sm font-medium text-gray-700 mb-2">
            Date of Birth <span class="text-red-500">*</span>
          </label>
          <input
            id="dateOfBirth"
            v-model="formData.dateOfBirth"
            type="date"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          />
          <p v-if="errors.dateOfBirth" class="mt-1 text-sm text-red-600">{{ errors.dateOfBirth }}</p>
        </div>

        <div>
          <label for="gender" class="block text-sm font-medium text-gray-700 mb-2">
            Gender <span class="text-red-500">*</span>
          </label>
          <select
            id="gender"
            v-model="formData.gender"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          >
            <option value="">Select Gender</option>
            <option value="male">Male</option>
            <option value="female">Female</option>
            <option value="other">Other</option>
          </select>
          <p v-if="errors.gender" class="mt-1 text-sm text-red-600">{{ errors.gender }}</p>
        </div>
      </div>

      <!-- Contact Information -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Contact fields commented out in original code -->
      </div>

      <!-- Address Section -->
      <div class="border-t pt-6">
        <h3 class="text-lg font-medium text-gray-900 mb-4">Address Information</h3>
        
        <!-- Street Address -->
        <div class="mb-6">
          <label for="street" class="block text-sm font-medium text-gray-700 mb-2">
            Street Address <span class="text-red-500">*</span>
          </label>
          <textarea
            id="street"
            v-model="formData.address.street"
            rows="2"
            required
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200 resize-none"
            placeholder="Enter your full address"
          ></textarea>
          <p v-if="errors['address.street']" class="mt-1 text-sm text-red-600">{{ errors['address.street'] }}</p>
        </div>

        <!-- City, State, Pincode -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-6">
          <div>
            <label for="city" class="block text-sm font-medium text-gray-700 mb-2">
              City <span class="text-red-500">*</span>
            </label>
            <input
              id="city"
              v-model="formData.address.city"
              type="text"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
              placeholder="City"
            />
            <p v-if="errors['address.city']" class="mt-1 text-sm text-red-600">{{ errors['address.city'] }}</p>
          </div>

          <div>
            <label for="state" class="block text-sm font-medium text-gray-700 mb-2">
              State <span class="text-red-500">*</span>
            </label>
            <select
              id="state"
              v-model="formData.address.state"
              required
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
            >
              <option value="">Select State</option>
              <option v-for="state in indianStates" :key="state" :value="state">{{ state }}</option>
            </select>
            <p v-if="errors['address.state']" class="mt-1 text-sm text-red-600">{{ errors['address.state'] }}</p>
          </div>

          <div>
            <label for="pincode" class="block text-sm font-medium text-gray-700 mb-2">
              Pincode
            </label>
            <input
              id="pincode"
              v-model="formData.address.pincode"
              type="text"
              pattern="[0-9]{6}"
              class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
              placeholder="000000"
              maxlength="6"
            />
            <p v-if="errors['address.pincode']" class="mt-1 text-sm text-red-600">{{ errors['address.pincode'] }}</p>
          </div>
        </div>

        <!-- Country -->
        <div class="mt-6">
          <label for="country" class="block text-sm font-medium text-gray-700 mb-2">
            Country
          </label>
          <input
            id="country"
            v-model="formData.address.country"
            type="text"
            readonly
            class="w-full px-4 py-3 border border-gray-300 rounded-xl bg-gray-50 text-gray-500"
          />
        </div>
      </div>

      <!-- Form Actions commented out in original code -->
    </form>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import { useDevoteeStore } from '@/stores/devotee'

// Props
const props = defineProps({
  data: {
    type: Object,
    default: () => ({})
  }
})

// Emits
const emit = defineEmits(['update', 'next'])

const toast = useToast()
const devoteeStore = useDevoteeStore()

// Form Data
const formData = ref({
  fullName: '',
  dateOfBirth: '',
  gender: '',
  // phone: '',
  // email: '',
  address: {
    street: '',
    city: '',
    state: '',
    pincode: '',
    country: 'India'
  }
})

// Validation
const errors = ref({})

// Indian States
const indianStates = [
  'Andhra Pradesh', 'Arunachal Pradesh', 'Assam', 'Bihar', 'Chhattisgarh',
  'Goa', 'Gujarat', 'Haryana', 'Himachal Pradesh', 'Jharkhand', 'Karnataka',
  'Kerala', 'Madhya Pradesh', 'Maharashtra', 'Manipur', 'Meghalaya',
  'Mizoram', 'Nagaland', 'Odisha', 'Punjab', 'Rajasthan', 'Sikkim',
  'Tamil Nadu', 'Telangana', 'Tripura', 'Uttar Pradesh', 'Uttarakhand',
  'West Bengal', 'Delhi', 'Jammu and Kashmir', 'Ladakh'
]

// Methods
function validateForm() {
  errors.value = {}
  
  if (!formData.value.fullName.trim()) {
    errors.value.fullName = 'Full name is required'
  }
  
  if (!formData.value.dateOfBirth) {
    errors.value.dateOfBirth = 'Date of birth is required'
  } else {
    const birthDate = new Date(formData.value.dateOfBirth)
    const today = new Date()
    const age = today.getFullYear() - birthDate.getFullYear()
    if (age < 5 || age > 120) {
      errors.value.dateOfBirth = 'Please enter a valid date of birth'
    }
  }
  
  if (!formData.value.gender) {
    errors.value.gender = 'Gender is required'
  }
  
  // if (formData.value.phone && !/^[+]?[0-9\s-()]{10,15}$/.test(formData.value.phone)) {
  //   errors.value.phone = 'Please enter a valid phone number'
  // }
  
  // if (formData.value.email && !/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.value.email)) {
  //   errors.value.email = 'Please enter a valid email address'
  // }
  
  if (!formData.value.address.street.trim()) {
    errors.value['address.street'] = 'Street address is required'
  }
  
  if (!formData.value.address.city.trim()) {
    errors.value['address.city'] = 'City is required'
  }
  
  if (!formData.value.address.state) {
    errors.value['address.state'] = 'State is required'
  }
  
  if (formData.value.address.pincode && !/^[0-9]{6}$/.test(formData.value.address.pincode)) {
    errors.value['address.pincode'] = 'Pincode must be 6 digits'
  }
  
  return Object.keys(errors.value).length === 0
}

// THIS IS THE FIXED FUNCTION - Don't transform data before emitting it
function handleSubmit() {
  if (validateForm()) {
    // Emit the original form data without transforming it
    emit('update', formData.value)
    
    // Move to next step
    emit('next')
  } else {
    toast.error('Please fix the errors before continuing')
  }
}

// Load existing data
onMounted(() => {
  if (props.data) {
    // Handle both nested format and flat API format
    formData.value = {
      fullName: props.data.full_name || props.data.fullName || '',
      dateOfBirth: props.data.dob || props.data.dateOfBirth || '',
      gender: props.data.gender || '',
      // phone: props.data.phone || '',
      // email: props.data.email || '',
      address: {
        street: props.data.street_address || props.data.address?.street || '',
        city: props.data.city || props.data.address?.city || '',
        state: props.data.state || props.data.address?.state || '',
        pincode: props.data.pincode || props.data.address?.pincode || '',
        country: props.data.country || props.data.address?.country || 'India'
      }
    }
  }
})

// Watch for validation on input change
// Add this watch handler
watch(formData, () => {
  if (Object.keys(errors.value).length > 0) {
    validateForm()
  }
  
  // Update parent component with data changes
  emit('update', formData.value)
}, { deep: true })
</script>