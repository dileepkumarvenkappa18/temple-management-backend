<!-- src/components/temple/TempleDetailsModal.vue -->
<template>
  <Teleport to="body">
    <div v-if="isOpen" class="fixed inset-0 z-[1000] overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
      <!-- Backdrop overlay -->
      <div class="fixed inset-0 bg-black bg-opacity-50" @click="$emit('close')"></div>
      
      <!-- Modal container -->
      <div class="flex min-h-screen items-center justify-center p-4">
        <!-- Modal content -->
        <div class="relative max-h-[90vh] w-full max-w-md overflow-y-auto rounded-lg bg-white shadow-xl">
          <!-- Modal header -->
          <div class="border-b border-gray-200 bg-white px-6 py-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg font-semibold text-gray-900">Temple Details</h3>
              <button 
                type="button" 
                class="rounded-md text-gray-400 hover:text-gray-500" 
                @click="$emit('close')"
              >
                <span class="sr-only">Close</span>
                <svg class="h-6 w-6" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
            <p class="mt-1 text-sm text-gray-600">Please provide the following details about your temple</p>
          </div>
          
          <!-- Modal body -->
          <div class="bg-white px-6 py-4">
            <form @submit.prevent="handleSubmit" class="space-y-4">
              <!-- Temple Name -->
              <div>
                <label for="temple-name" class="block text-sm font-medium text-gray-700">
                  Temple Name <span class="text-red-500">*</span>
                </label>
                <input
                  id="temple-name"
                  v-model="templeData.name"
                  type="text"
                  placeholder="Enter temple name"
                  required
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                  :class="{'border-red-500': errors.name}"
                />
                <p v-if="errors.name" class="mt-1 text-sm text-red-600">{{ errors.name }}</p>
              </div>
              
              <!-- Temple Place -->
              <div>
                <label for="temple-place" class="block text-sm font-medium text-gray-700">
                  Temple Place <span class="text-red-500">*</span>
                </label>
                <input
                  id="temple-place"
                  v-model="templeData.place"
                  type="text"
                  placeholder="Enter temple location"
                  required
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                  :class="{'border-red-500': errors.place}"
                />
                <p v-if="errors.place" class="mt-1 text-sm text-red-600">{{ errors.place }}</p>
              </div>
              
              <!-- Temple Address -->
              <div>
                <label for="temple-address" class="block text-sm font-medium text-gray-700">
                  Temple Address <span class="text-red-500">*</span>
                </label>
                <textarea
                  id="temple-address"
                  v-model="templeData.address"
                  rows="3"
                  placeholder="Enter complete temple address"
                  required
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                  :class="{'border-red-500': errors.address}"
                ></textarea>
                <p v-if="errors.address" class="mt-1 text-sm text-red-600">{{ errors.address }}</p>
              </div>
              
              <!-- Temple Phone Number -->
              <div>
                <label for="temple-phone" class="block text-sm font-medium text-gray-700">
                  Temple Phone Number <span class="text-red-500">*</span>
                </label>
                <input
                  id="temple-phone"
                  v-model="templeData.phoneNumber"
                  type="tel"
                  placeholder="Enter temple contact number"
                  required
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                  :class="{'border-red-500': errors.phoneNumber}"
                />
                <p v-if="errors.phoneNumber" class="mt-1 text-sm text-red-600">{{ errors.phoneNumber }}</p>
              </div>
              
              <!-- Temple Description (Now Required) -->
              <div>
                <label for="temple-description" class="block text-sm font-medium text-gray-700">
                  Description <span class="text-red-500">*</span>
                </label>
                <textarea
                  id="temple-description"
                  v-model="templeData.description"
                  rows="4"
                  placeholder="Provide additional details about your temple (history, services, etc.)"
                  required
                  class="mt-1 block w-full rounded-md border border-gray-300 px-3 py-2 shadow-sm focus:border-indigo-500 focus:outline-none focus:ring-indigo-500 sm:text-sm"
                  :class="{'border-red-500': errors.description}"
                ></textarea>
                <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
              </div>
            </form>
          </div>
          
          <!-- Modal footer -->
          <div class="border-t border-gray-200 bg-gray-50 px-6 py-4">
            <div class="flex justify-end space-x-3">
              <button 
                type="button" 
                class="inline-flex justify-center rounded-md border border-gray-300 bg-white px-4 py-2 text-sm font-medium text-gray-700 shadow-sm hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                @click="$emit('close')"
              >
                Cancel
              </button>
              <button 
                type="button" 
                class="inline-flex justify-center rounded-md border border-transparent bg-indigo-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
                @click="handleSubmit"
                :disabled="isSubmitting"
              >
                <span v-if="isSubmitting" class="mr-2">
                  <svg class="h-4 w-4 animate-spin text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                </span>
                Save Temple Details
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<script setup>
import { ref, reactive, watch } from 'vue';

const props = defineProps({
  isOpen: {
    type: Boolean,
    default: false
  }
});

const emit = defineEmits(['close', 'submit']);

const isSubmitting = ref(false);

// Temple data state
const templeData = reactive({
  name: '',
  place: '',
  address: '',
  phoneNumber: '',
  description: ''
});

// Form validation errors
const errors = reactive({
  name: '',
  place: '',
  address: '',
  phoneNumber: '',
  description: '' // Added description error field
});

// Validate form
const validateForm = () => {
  let isValid = true;
  
  // Reset errors
  Object.keys(errors).forEach(key => {
    errors[key] = '';
  });
  
  // Validate required fields
  if (!templeData.name.trim()) {
    errors.name = 'Temple name is required';
    isValid = false;
  }
  
  if (!templeData.place.trim()) {
    errors.place = 'Temple place is required';
    isValid = false;
  }
  
  if (!templeData.address.trim()) {
    errors.address = 'Temple address is required';
    isValid = false;
  }
  
  if (!templeData.phoneNumber.trim()) {
    errors.phoneNumber = 'Temple phone number is required';
    isValid = false;
  } else if (!/^\d{10}$/.test(templeData.phoneNumber.replace(/\D/g, ''))) {
    errors.phoneNumber = 'Please enter a valid 10-digit phone number';
    isValid = false;
  }
  
  // Validate description (now required)
  if (!templeData.description.trim()) {
    errors.description = 'Temple description is required';
    isValid = false;
  }
  
  return isValid;
};

// Handle form submission
const handleSubmit = async () => {
  if (!validateForm()) {
    return;
  }
  
  try {
    isSubmitting.value = true;
    
    // Emit the data to the parent component
    emit('submit', { ...templeData });
    
    // Reset form after successful submission
    resetForm();
  } catch (error) {
    console.error('Error saving temple details:', error);
  } finally {
    isSubmitting.value = false;
  }
};

// Reset form values
const resetForm = () => {
  Object.keys(templeData).forEach(key => {
    templeData[key] = '';
  });
  
  Object.keys(errors).forEach(key => {
    errors[key] = '';
  });
};

// Reset form when modal opens
watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    resetForm();
  }
});

// Prevent body scrolling when modal is open
watch(() => props.isOpen, (newVal) => {
  if (newVal) {
    document.body.style.overflow = 'hidden';
  } else {
    document.body.style.overflow = '';
  }
});
</script>