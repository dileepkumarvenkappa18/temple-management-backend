<template>
  <div class="seva-form">
    <!-- Success State -->
    <div v-if="submitted" class="card bg-white shadow-md rounded-lg p-6 text-center">
      <div class="mb-4 p-6 bg-green-100 text-green-800 rounded-md flex flex-col items-center">
        <svg class="h-16 w-16 text-green-500 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
        </svg>
        <h2 class="text-xl font-semibold mb-2">Seva Created Successfully!</h2>
        <p class="text-green-700 mb-4">Your seva has been created and is now available in the system.</p>
      </div>
    </div>

    <!-- Form State -->
    <div v-else class="card bg-white shadow-md rounded-lg p-6">
      <h2 class="text-xl font-semibold text-indigo-700 mb-6">
        {{ isEditing ? 'Edit Seva' : 'Create New Seva' }}
      </h2>
      
      <!-- Error message -->
      <div v-if="apiError" class="mb-4 p-3 bg-red-100 text-red-800 rounded-md flex items-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" viewBox="0 0 20 20" fill="currentColor">
          <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7 4a1 1 0 11-2 0 1 1 0 012 0zm-1-9a1 1 0 00-1 1v4a1 1 0 102 0V6a1 1 0 00-1-1z" clip-rule="evenodd" />
        </svg>
        <span>{{ apiError }}</span>
      </div>
      
      <div class="space-y-4">
        <!-- Basic Information -->
        <div class="form-group">
          <label for="name" class="block text-sm font-medium text-gray-700">Seva Name*</label>
          <input
            id="name"
            v-model="form.name"
            type="text"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            :class="{'border-red-500': errors.name}"
          />
          <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
        </div>
        
        <div class="form-group">
          <label for="type" class="block text-sm font-medium text-gray-700">Seva Type*</label>
          <select
            id="type"
            v-model="form.type"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            :class="{'border-red-500': errors.type}"
          >
            <option value="" disabled>Select Seva Type</option>
            <option value="daily">Daily Seva</option>
            <option value="special">Special Seva</option>
            <option value="festival">Festival Seva</option>
            <option value="personal">Personal Seva</option>
          </select>
          <p v-if="errors.type" class="mt-1 text-sm text-red-600">{{ errors.type }}</p>
        </div>
        
        <div class="form-group">
          <label for="description" class="block text-sm font-medium text-gray-700">Description*</label>
          <textarea
            id="description"
            v-model="form.description"
            rows="3"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            :class="{'border-red-500': errors.description}"
          ></textarea>
          <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div class="form-group">
            <label for="price" class="block text-sm font-medium text-gray-700">Price (₹)*</label>
            <input
              id="price"
              v-model="form.price"
              type="number"
              min="0"
              step="0.01"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              :class="{'border-red-500': errors.price}"
            />
            <p v-if="errors.price" class="mt-1 text-sm text-red-600">{{ errors.price }}</p>
          </div>
          
          <div class="form-group">
            <label for="duration" class="block text-sm font-medium text-gray-700">Duration (minutes)</label>
            <input
              id="duration"
              v-model="form.duration"
              type="number"
              min="5"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            />
          </div>
        </div>
        
        <div class="form-group">
          <label class="block text-sm font-medium text-gray-700">Availability</label>
          <div class="mt-2">
            <p class="text-sm text-gray-500 mb-2">Select days when this seva is available:</p>
            <div class="flex flex-wrap gap-2">
              <label 
                v-for="day in weekDays" 
                :key="day.value"
                class="inline-flex items-center px-3 py-1 rounded-full text-sm"
                :class="form.availableDays.includes(day.value) ? 'bg-indigo-100 text-indigo-800' : 'bg-gray-100 text-gray-700'"
              >
                <input 
                  type="checkbox"
                  :value="day.value"
                  v-model="form.availableDays"
                  class="mr-1.5 h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                {{ day.label }}
              </label>
            </div>
          </div>
        </div>
        
        <div class="form-group">
          <label class="block text-sm font-medium text-gray-700">Time Slots</label>
          <div class="mt-2 space-y-3">
            <div 
              v-for="(slot, index) in form.timeSlots" 
              :key="index"
              class="flex items-center space-x-2"
            >
              <div class="w-1/3">
                <input
                  v-model="slot.start"
                  type="time"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                />
              </div>
              <span class="text-gray-500">to</span>
              <div class="w-1/3">
                <input
                  v-model="slot.end"
                  type="time"
                  class="block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
                />
              </div>
              <button 
                type="button"
                @click="removeTimeSlot(index)"
                class="p-1 text-red-600 hover:text-red-800"
                :disabled="form.timeSlots.length === 1"
              >
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM7 9a1 1 0 000 2h6a1 1 0 100-2H7z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
            
            <button 
              type="button"
              @click="addTimeSlot"
              class="inline-flex items-center px-3 py-1.5 border border-transparent text-sm font-medium rounded-md text-indigo-700 bg-indigo-100 hover:bg-indigo-200 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="h-4 w-4 mr-1" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-11a1 1 0 10-2 0v2H7a1 1 0 100 2h2v2a1 1 0 102 0v-2h2a1 1 0 100-2h-2V7z" clip-rule="evenodd" />
              </svg>
              Add Time Slot
            </button>
          </div>
        </div>

        <div class="form-group">
          <label for="max_bookings" class="block text-sm font-medium text-gray-700">Maximum Bookings Per Day</label>
          <input
            id="max_bookings"
            v-model="form.max_bookings_per_day"
            type="number"
            min="1"
            class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>
        
        <!-- Toggle for advanced options -->
        <div class="pt-2">
          <button 
            type="button" 
            @click="showAdvanced = !showAdvanced"
            class="text-indigo-600 hover:text-indigo-800 text-sm font-medium flex items-center"
          >
            {{ showAdvanced ? 'Hide' : 'Show' }} Advanced Options
            <svg 
              xmlns="http://www.w3.org/2000/svg" 
              class="h-4 w-4 ml-1" 
              :class="showAdvanced ? 'transform rotate-180' : ''"
              viewBox="0 0 20 20" 
              fill="currentColor"
            >
              <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
            </svg>
          </button>
        </div>
        
        <!-- Advanced options -->
        <div v-if="showAdvanced" class="space-y-4 pt-2 border-t border-gray-200 mt-2">
          <div class="form-group">
            <label for="requirements" class="block text-sm font-medium text-gray-700">Special Requirements</label>
            <textarea
              id="requirements"
              v-model="form.requirements"
              rows="2"
              placeholder="Any special items required for this seva"
              class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
            ></textarea>
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div class="form-group">
              <label for="advanceBookingDays" class="block text-sm font-medium text-gray-700">Advance Booking Days</label>
              <input
                id="advanceBookingDays"
                v-model="form.advanceBookingDays"
                type="number"
                min="0"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              />
            </div>
            
            <div class="form-group">
              <label for="cancellationHours" class="block text-sm font-medium text-gray-700">Cancellation Hours</label>
              <input
                id="cancellationHours"
                v-model="form.cancellationHours"
                type="number"
                min="0"
                class="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-indigo-500 focus:ring-indigo-500"
              />
            </div>
          </div>
          
          <div class="space-y-3">
            <div class="form-group">
              <div class="flex items-center">
                <input
                  id="isRecurring"
                  v-model="form.isRecurring"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="isRecurring" class="ml-2 block text-sm font-medium text-gray-700">
                  Allow Recurring Bookings
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <div class="flex items-center">
                <input
                  id="requiresApproval"
                  v-model="form.requiresApproval"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="requiresApproval" class="ml-2 block text-sm font-medium text-gray-700">
                  Requires Approval
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <div class="flex items-center">
                <input
                  id="allowMultipleBookings"
                  v-model="form.allowMultipleBookings"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="allowMultipleBookings" class="ml-2 block text-sm font-medium text-gray-700">
                  Allow Multiple Bookings per User
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <div class="flex items-center">
                <input
                  id="sendReminders"
                  v-model="form.sendReminders"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="sendReminders" class="ml-2 block text-sm font-medium text-gray-700">
                  Send Reminders
                </label>
              </div>
            </div>
            
            <div class="form-group">
              <div class="flex items-center">
                <input
                  id="is_active"
                  v-model="form.is_active"
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                />
                <label for="is_active" class="ml-2 block text-sm font-medium text-gray-700">
                  Active
                </label>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Development mode indicator -->
      <div v-if="isDev" class="mt-4 p-2 bg-yellow-50 text-yellow-800 text-xs rounded-md">
        Development Mode: API calls will be simulated
      </div>
      
      <!-- Form actions -->
      <div class="mt-6 flex justify-end space-x-3">
        <button
          type="button"
          @click="handleCancel"
          class="inline-flex items-center px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleSubmit"
          :disabled="loading || formSubmitted"
          class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-md shadow-sm text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
        >
          <span v-if="loading" class="mr-2">
            <svg class="animate-spin h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
          </span>
          {{ isEditing ? 'Update Seva' : 'Create Seva' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import { useSevaStore } from '@/stores/seva'
import { useRouter } from 'vue-router'

// Props
const props = defineProps({
  sevaId: {
    type: [String, Number],
    default: null
  },
  initialData: {
    type: Object,
    default: () => ({})
  },
  entityId: {
    type: [String, Number],
    required: true
  }
})

// Emits
const emit = defineEmits(['submit', 'cancel', 'seva-created'])

// Setup
const router = useRouter()
const { showToast } = useToast()
const sevaStore = useSevaStore()
const loading = ref(false)
const showAdvanced = ref(false)
const errors = ref({})
const formSubmitted = ref(false)
const apiError = ref(null)
const submitted = ref(false) // New state to track successful submission

// Check if in development mode
const isDev = computed(() => {
  return import.meta.env.MODE === 'development' || 
         window.location.hostname === 'localhost' || 
         window.location.hostname === '127.0.0.1' ||
         window.location.port === '5173'
})

// Form data
const form = reactive({
  name: '',
  type: '',
  description: '',
  price: null,
  duration: 30,
  max_bookings_per_day: 10,
  availability_schedule: '{}',
  is_active: true,
  
  // Frontend only fields (will be transformed before API call)
  availableDays: [],
  timeSlots: [{ start: '09:00', end: '10:00' }],
  requirements: '',
  advanceBookingDays: 7,
  cancellationHours: 24,
  isRecurring: false,
  requiresApproval: true,
  allowMultipleBookings: false,
  sendReminders: true
})

// Week days options
const weekDays = [
  { value: 'monday', label: 'Monday' },
  { value: 'tuesday', label: 'Tuesday' },
  { value: 'wednesday', label: 'Wednesday' },
  { value: 'thursday', label: 'Thursday' },
  { value: 'friday', label: 'Friday' },
  { value: 'saturday', label: 'Saturday' },
  { value: 'sunday', label: 'Sunday' }
]

// Computed
const isEditing = computed(() => props.sevaId !== null)

// Methods
const addTimeSlot = () => {
  form.timeSlots.push({ start: '09:00', end: '10:00' })
}

const removeTimeSlot = (index) => {
  form.timeSlots.splice(index, 1)
}

const validateForm = () => {
  errors.value = {}
  
  if (!form.name.trim()) {
    errors.value.name = 'Seva name is required'
  }
  
  if (!form.type) {
    errors.value.type = 'Seva type is required'
  }
  
  if (!form.description.trim()) {
    errors.value.description = 'Description is required'
  }
  
  if (form.price === null || form.price < 0) {
    errors.value.price = 'Valid price is required'
  }
  
  return Object.keys(errors.value).length === 0
}

const prepareFormData = () => {
  // Convert time slots and available days to JSON schedule
  const schedule = {}
  
  form.availableDays.forEach(day => {
    schedule[day] = form.timeSlots.map(slot => [slot.start, slot.end])
  })
  
  // Return API-compatible data
  return {
    name: form.name,
    description: form.description,
    price: parseFloat(form.price),
    duration: parseInt(form.duration),
    max_bookings_per_day: parseInt(form.max_bookings_per_day),
    availability_schedule: JSON.stringify(schedule),
    is_active: form.is_active,
    entity_id: props.entityId
  }
}

// Mock function to simulate API response in development mode
const mockApiCall = async (payload) => {
  // Simulate network delay
  await new Promise(resolve => setTimeout(resolve, 800))
  
  // Generate a random ID for the new seva
  const mockId = Math.floor(Math.random() * 1000) + 1
  
  // Format the time in AM/PM format
  const formatTime = (timeStr) => {
    if (!timeStr) return '09:00 AM'
    const [hours, minutes] = timeStr.split(':')
    const hour = parseInt(hours)
    const ampm = hour >= 12 ? 'PM' : 'AM'
    const hour12 = hour % 12 || 12
    return `${hour12}:${minutes} ${ampm}`
  }
  
  // Create a mock seva object that matches what the SevaList component expects
  return {
    success: true,
    message: isEditing.value ? 'Seva updated successfully' : 'Seva created successfully',
    data: {
      id: isEditing.value ? props.sevaId : mockId,
      name: payload.name,
      description: payload.description,
      type: form.type || 'special', // Match the type to the list's expected values
      devotee: {
        name: 'Dev User',
        phone: '+91 98765 43210'
      },
      date: new Date().toISOString().split('T')[0], // Today's date in YYYY-MM-DD format
      time: formatTime(form.timeSlots[0]?.start), // Format the time
      amount: parseFloat(payload.price),
      status: form.requiresApproval ? 'pending' : 'approved',
      notes: form.requirements || payload.description,
      price: parseFloat(payload.price),
      duration: parseInt(payload.duration),
      max_bookings_per_day: parseInt(payload.max_bookings_per_day),
      availability_schedule: payload.availability_schedule,
      is_active: payload.is_active,
      entity_id: payload.entity_id,
      created_at: new Date().toISOString(),
      updated_at: new Date().toISOString()
    }
  }
}

const handleSubmit = async () => {
  if (!validateForm()) {
    showToast('Please fix the errors above', 'error')
    return
  }

  apiError.value = null
  formSubmitted.value = true
  loading.value = true // Set loading manually
  
  try {
    const payload = prepareFormData()
    let result

    // Use mock API call in development mode to bypass API errors
    if (isDev.value) {
      console.log('Development mode: Using mock API response')
      result = await mockApiCall(payload)
    } else {
      // Real API call for production
      if (isEditing.value) {
        result = await sevaStore.updateSeva(props.sevaId, payload)
      } else {
        result = await sevaStore.createSeva(payload)
      }
    }
    
    if (result && result.success) {
      // Show success message using toast notification
      showToast(isEditing.value ? 'Seva updated successfully!' : 'Seva created successfully!', 'success')
      
      // Emit events for parent components
      emit('submit', result.data)
      emit('seva-created', result.data)
      
      // For development mode, manually add the seva to the store
      if (isDev.value) {
        sevaStore.addSevaToStore(result.data)
      } else {
        // Force refresh the seva list in production
        await sevaStore.fetchSevas(props.entityId)
      }
      
      // Set submitted state to true to hide the form and show success message
      submitted.value = true
    } else {
      formSubmitted.value = false
      apiError.value = result?.message || 'An error occurred while saving the seva'
      showToast(apiError.value, 'error')
    }
  } catch (error) {
    console.error('Error saving seva:', error)
    formSubmitted.value = false
    apiError.value = 'Failed to save seva. Please try again.'
    showToast(apiError.value, 'error')
  } finally {
    loading.value = false // Always reset loading state
  }
}

const handleCancel = () => {
  emit('cancel')
  // Navigate back to the sevas list page
  router.push(`/entity/${props.entityId}/sevas`)
}

const loadSevaData = async () => {
  if (isEditing.value && props.initialData) {
    // Convert backend data to form format
    try {
      // Load basic fields
      form.name = props.initialData.name
      form.description = props.initialData.description
      form.price = props.initialData.price
      form.duration = props.initialData.duration
      form.max_bookings_per_day = props.initialData.max_bookings_per_day
      form.is_active = props.initialData.is_active
      
      // Parse availability schedule
      if (props.initialData.availability_schedule) {
        const schedule = JSON.parse(props.initialData.availability_schedule)
        form.availableDays = Object.keys(schedule)
        
        if (form.availableDays.length > 0) {
          // Extract time slots from the first available day
          const firstDaySlots = schedule[form.availableDays[0]]
          if (Array.isArray(firstDaySlots) && firstDaySlots.length > 0) {
            form.timeSlots = firstDaySlots.map(slot => ({
              start: slot[0],
              end: slot[1]
            }))
          }
        }
      }
    } catch (error) {
      console.error('Error parsing seva data:', error)
      showToast('Error loading seva data', 'error')
    }
  }
}

// Lifecycle
onMounted(() => {
  loadSevaData()
})
</script>