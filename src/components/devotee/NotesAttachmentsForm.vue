<template>
  <div class="bg-white rounded-xl shadow-lg p-6 border border-indigo-100">
    <!-- Step Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h3 class="text-xl font-semibold text-gray-900">Notes & Attachments</h3>
        <p class="text-sm text-gray-600 mt-1">Add photos, health information, and personal notes</p>
      </div>
      <div class="text-right">
        <span class="text-xs font-medium text-indigo-600 bg-indigo-50 px-2 py-1 rounded-full">
          Step 6 of 6
        </span>
      </div>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-8">
      <!-- Photo Upload Section -->
      <!-- <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16l4.586-4.586a2 2 0 012.828 0L16 16m-2-2l1.586-1.586a2 2 0 012.828 0L20 14m-6-6h.01M6 20h12a2 2 0 002-2V6a2 2 0 00-2-2H6a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          <h4 class="text-lg font-medium text-gray-900">Profile Photos</h4>
        </div>
        
        
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div class="space-y-3">
            <label class="block text-sm font-medium text-gray-700">Profile Photo</label>
            <div class="relative">
              <div v-if="!formData.profilePhoto" 
                   class="w-full h-40 border-2 border-dashed border-indigo-300 rounded-lg flex flex-col items-center justify-center hover:border-indigo-400 transition-colors cursor-pointer"
                   @click="triggerFileInput('profile')"
                   @dragover.prevent
                   @drop.prevent="handleFileDrop($event, 'profile')">
                <svg class="w-8 h-8 text-indigo-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                <p class="text-sm text-gray-600">Click to upload or drag and drop</p>
                <p class="text-xs text-gray-500">PNG, JPG up to 5MB</p>
              </div>
              
              
              <div v-else class="relative">
                <img :src="formData.profilePhoto" alt="Profile" class="w-full h-40 object-cover rounded-lg">
                <button type="button" 
                        @click="removePhoto('profile')"
                        class="absolute top-2 right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
            <input ref="profileFileInput" type="file" class="hidden" accept="image/*" @change="handleFileSelect($event, 'profile')">
          </div>

          
          <div class="space-y-3">
            <label class="block text-sm font-medium text-gray-700">Family Photo</label>
            <div class="relative">
              <div v-if="!formData.familyPhoto" 
                   class="w-full h-40 border-2 border-dashed border-indigo-300 rounded-lg flex flex-col items-center justify-center hover:border-indigo-400 transition-colors cursor-pointer"
                   @click="triggerFileInput('family')"
                   @dragover.prevent
                   @drop.prevent="handleFileDrop($event, 'family')">
                <svg class="w-8 h-8 text-indigo-400 mb-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
                </svg>
                <p class="text-sm text-gray-600">Click to upload or drag and drop</p>
                <p class="text-xs text-gray-500">PNG, JPG up to 5MB</p>
              </div>
              
              <div v-else class="relative">
                <img :src="formData.familyPhoto" alt="Family" class="w-full h-40 object-cover rounded-lg">
                <button type="button" 
                        @click="removePhoto('family')"
                        class="absolute top-2 right-2 bg-red-500 text-white rounded-full p-1 hover:bg-red-600 transition-colors">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                  </svg>
                </button>
              </div>
            </div>
            <input ref="familyFileInput" type="file" class="hidden" accept="image/*" @change="handleFileSelect($event, 'family')">
          </div>
        </div>
      </div> -->

      <!-- Health Information -->
      <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          <h4 class="text-lg font-medium text-gray-900">Health Information</h4>
        </div>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <!-- Allergies -->
          <div>
            <label for="allergies" class="block text-sm font-medium text-gray-700 mb-2">
              Allergies/Medical Conditions
            </label>
            <textarea
              id="allergies"
              v-model="formData.allergies"
              rows="4"
              placeholder="List any allergies or medical conditions..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
            ></textarea>
          </div>

          <!-- Dietary Restrictions -->
          <div>
            <label for="dietaryRestrictions" class="block text-sm font-medium text-gray-700 mb-2">
              Dietary Restrictions
            </label>
            <textarea
              id="dietaryRestrictions"
              v-model="formData.dietaryRestrictions"
              rows="4"
              placeholder="Vegetarian, vegan, specific food restrictions..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- Sankalpa & Personal Notes -->
      <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
          </svg>
          <h4 class="text-lg font-medium text-gray-900">Sankalpa & Personal Notes</h4>
        </div>

        <div class="space-y-6">
          <!-- Sankalpa -->
          <div>
            <label for="sankalpa" class="block text-sm font-medium text-gray-700 mb-2">
              Personal Sankalpa/Intentions
            </label>
            <textarea
              id="sankalpa"
              v-model="formData.sankalpa"
              rows="4"
              placeholder="Share your spiritual intentions, goals, or personal sankalpa..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
            ></textarea>
          </div>

          <!-- Additional Notes -->
          <div>
            <label for="additionalNotes" class="block text-sm font-medium text-gray-700 mb-2">
              Additional Notes
            </label>
            <textarea
              id="additionalNotes"
              v-model="formData.additionalNotes"
              rows="4"
              placeholder="Any other information you'd like to share with the temple..."
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
            ></textarea>
          </div>
        </div>
      </div>

      <!-- Document Attachments -->
      <!-- <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13" />
          </svg>
          <h4 class="text-lg font-medium text-gray-900">Document Attachments</h4>
        </div>

        <div class="border-2 border-dashed border-indigo-300 rounded-lg p-6 hover:border-indigo-400 transition-colors">
          <div class="text-center">
            <svg class="mx-auto w-12 h-12 text-indigo-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
            </svg>
            <p class="text-sm text-gray-600 mb-2">
              Upload supporting documents (Birth certificate, ID proof, etc.)
            </p>
            <button type="button" 
                    @click="triggerFileInput('documents')"
                    class="bg-indigo-50 text-indigo-600 px-4 py-2 rounded-lg hover:bg-indigo-100 transition-colors text-sm font-medium">
              Choose Files
            </button>
            <p class="text-xs text-gray-500 mt-2">PDF, JPG, PNG up to 10MB each</p>
          </div>
          
          
          <div v-if="formData.documents.length > 0" class="mt-4 space-y-2">
            <div v-for="(doc, index) in formData.documents" :key="index" 
                 class="flex items-center justify-between bg-gray-50 p-3 rounded-lg">
              <div class="flex items-center space-x-3">
                <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span class="text-sm text-gray-700">{{ doc.name }}</span>
              </div>
              <button type="button" 
                      @click="removeDocument(index)"
                      class="text-red-500 hover:text-red-700 transition-colors">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
                </svg>
              </button>
            </div>
          </div>
        </div>
        <input ref="documentsFileInput" type="file" class="hidden" multiple accept=".pdf,.jpg,.jpeg,.png" @change="handleDocumentSelect">
      </div> -->

      <!-- Privacy Preferences -->
      <!-- <div class="space-y-4">
        <div class="flex items-center space-x-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
          <h4 class="text-lg font-medium text-gray-900">Privacy Preferences</h4>
        </div>

        <div class="space-y-4 bg-gray-50 p-4 rounded-lg">
          <div class="flex items-center justify-between">
            <div>
              <label class="text-sm font-medium text-gray-700">Share profile with other devotees</label>
              <p class="text-xs text-gray-600">Allow other devotees to see your profile information</p>
            </div>
            <button type="button" 
                    @click="formData.shareProfile = !formData.shareProfile"
                    :class="[
                      'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                      formData.shareProfile ? 'bg-indigo-600' : 'bg-gray-300'
                    ]">
              <span :class="[
                'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                formData.shareProfile ? 'translate-x-6' : 'translate-x-1'
              ]" />
            </button>
          </div>

          <div class="flex items-center justify-between">
            <div>
              <label class="text-sm font-medium text-gray-700">Receive event notifications</label>
              <p class="text-xs text-gray-600">Get notified about temple events and announcements</p>
            </div>
            <button type="button" 
                    @click="formData.receiveNotifications = !formData.receiveNotifications"
                    :class="[
                      'relative inline-flex h-6 w-11 items-center rounded-full transition-colors',
                      formData.receiveNotifications ? 'bg-indigo-600' : 'bg-gray-300'
                    ]">
              <span :class="[
                'inline-block h-4 w-4 transform rounded-full bg-white transition-transform',
                formData.receiveNotifications ? 'translate-x-6' : 'translate-x-1'
              ]" />
            </button>
          </div>
        </div>
      </div> -->

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row justify-between items-center pt-6 border-t border-gray-200 space-y-4 sm:space-y-0">
        <!-- <button type="button" 
                @click="$emit('previous')"
                class="w-full sm:w-auto px-6 py-3 text-indigo-600 bg-white border border-indigo-300 rounded-lg hover:bg-indigo-50 transition-colors font-medium">
          <div class="flex items-center justify-center space-x-2">
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            <span>Previous Step</span>
          </div>
        </button> -->

        <div class="flex space-x-4">
          <!-- <button type="button" 
                  @click="$emit('skip')"
                  class="px-6 py-3 text-gray-600 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors font-medium">
            Skip for Now
          </button> -->
          
          <!-- <button type="submit" 
                  :disabled="isSubmitting"
                  class="px-8 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium">
            <div class="flex items-center space-x-2">
              <svg v-if="isSubmitting" class="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              <span>{{ isSubmitting ? 'Saving...' : 'Complete Profile' }}</span>
            </div>
          </button> -->
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { useToast } from '@/composables/useToast'
import { watch, onMounted } from 'vue';


const emit = defineEmits(['previous', 'skip', 'submit', 'update', 'update-data'])

const { showToast } = useToast()

const isSubmitting = ref(false)
const profileFileInput = ref(null)
const familyFileInput = ref(null)
const documentsFileInput = ref(null)
const props = defineProps({
  data: {
    type: Object,
    default: () => ({})
  }
})

const formData = reactive({
  // profilePhoto: null,
  // familyPhoto: null,
  allergies: '',
  dietaryRestrictions: '',
  sankalpa: '',
  additionalNotes: '',
  // documents: [],
  shareProfile: true,
  receiveNotifications: true
})

// File handling methods
// const triggerFileInput = (type) => {
//   if (type === 'profile') {
//     profileFileInput.value.click()
//   } else if (type === 'family') {
//     familyFileInput.value.click()
//   } else if (type === 'documents') {
//     documentsFileInput.value.click()
//   }
// }

const handleFileSelect = (event, type) => {
  const file = event.target.files[0]
  if (file) {
    if (file.size > 5 * 1024 * 1024) {
      showToast('File size must be less than 5MB', 'error')
      return
    }
    
    const reader = new FileReader()
    reader.onload = (e) => {
      if (type === 'profile') {
        formData.profilePhoto = e.target.result
      } else if (type === 'family') {
        formData.familyPhoto = e.target.result
      }
    }
    reader.readAsDataURL(file)
  }
}

const handleFileDrop = (event, type) => {
  const file = event.dataTransfer.files[0]
  if (file) {
    const mockEvent = { target: { files: [file] } }
    handleFileSelect(mockEvent, type)
  }
}

// const removePhoto = (type) => {
//   if (type === 'profile') {
//     formData.profilePhoto = null
//     profileFileInput.value.value = ''
//   } else if (type === 'family') {
//     formData.familyPhoto = null
//     familyFileInput.value.value = ''
//   }
// }

// const handleDocumentSelect = (event) => {
//   const files = Array.from(event.target.files)
  
//   files.forEach(file => {
//     if (file.size > 10 * 1024 * 1024) {
//       showToast(`${file.name} exceeds 10MB limit`, 'error')
//       return
//     }
    
//     formData.documents.push({
//       name: file.name,
//       size: file.size,
//       file: file
//     })
//   })
  
//   if (files.length > 0) {
//     showToast(`${files.length} document(s) added successfully`, 'success')
//   }
// }

// const removeDocument = (index) => {
//   formData.documents.splice(index, 1)
//   showToast('Document removed', 'success')
// }

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    showToast('Profile completed successfully!', 'success')
    emit('submit', formData)
    
  } catch (error) {
    showToast('Error completing profile. Please try again.', 'error')
  } finally {
    isSubmitting.value = false
  }
}

// Add watch handler
watch(formData, () => {
  // Update parent component with data changes
  emit('update', formData)
}, { deep: true })

// In onMounted, initialize with existing data
onMounted(() => {
  if (props.data) {
    formData.healthNotes = props.data.healthNotes || ''
    formData.allergies = props.data.allergies || ''
    formData.dietaryRestrictions = props.data.dietaryRestrictions || ''
    formData.sankalpa = props.data.sankalpa || ''
    formData.additionalNotes = props.data.additionalNotes || ''
    // Don't copy over files/photos as they need special handling
  }
})
</script>