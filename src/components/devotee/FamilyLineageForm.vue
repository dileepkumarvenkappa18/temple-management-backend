<template>
  <div class="bg-white rounded-xl shadow-lg p-6 md:p-8">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center gap-3 mb-4">
        <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
          </svg>
        </div>
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Family Lineage</h2>
          <p class="text-gray-600">Build your family tree and heritage</p>
        </div>
      </div>
      <div class="w-full bg-gray-200 rounded-full h-2">
        <div class="bg-indigo-600 h-2 rounded-full transition-all duration-300" :style="{ width: '60%' }"></div>
      </div>
      <p class="text-sm text-gray-500 mt-2">Step 3 of 6 - 60% Complete</p>
    </div>

    <form @submit.prevent="handleSubmit" class="space-y-8">
      <!-- Parents Section -->
      <div class="bg-indigo-50 rounded-xl p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
          </svg>
          Parents
        </h3>
        
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Father Details -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-800 border-b border-indigo-200 pb-2">Father's Details</h4>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Full Name</label>
              <input
                v-model="formData.father.name"
                type="text"
                placeholder="Enter father's full name"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Gotra</label>
              <input
                v-model="formData.father.gotra"
                type="text"
                placeholder="Enter father's gotra"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Native Place</label>
              <input
                v-model="formData.father.nativePlace"
                type="text"
                placeholder="Enter father's native place"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Veda Shaka</label>
              <select
                v-model="formData.father.vedaShaka"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
                <option value="">Select Veda Shaka</option>
                <option value="Rigveda">Rigveda</option>
                <option value="Yajurveda">Yajurveda</option>
                <option value="Samaveda">Samaveda</option>
                <option value="Atharvaveda">Atharvaveda</option>
              </select>
            </div>
          </div>

          <!-- Mother Details -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-800 border-b border-indigo-200 pb-2">Mother's Details</h4>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Full Name</label>
              <input
                v-model="formData.mother.name"
                type="text"
                placeholder="Enter mother's full name"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Maiden Gotra</label>
              <input
                v-model="formData.mother.maidenGotra"
                type="text"
                placeholder="Enter mother's maiden gotra"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Native Place</label>
              <input
                v-model="formData.mother.nativePlace"
                type="text"
                placeholder="Enter mother's native place"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Father's Name</label>
              <input
                v-model="formData.mother.fatherName"
                type="text"
                placeholder="Enter maternal grandfather's name"
                class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
            </div>
          </div>
        </div>
      </div>

      <!-- Grandparents Section -->
      <div class="bg-amber-50 rounded-xl p-6">
        <h3 class="text-lg font-semibold text-gray-900 mb-6 flex items-center gap-2">
          <svg class="w-5 h-5 text-amber-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"/>
          </svg>
          Grandparents (Optional)
        </h3>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <!-- Paternal Grandparents -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-800 border-b border-amber-200 pb-2">Paternal Grandparents</h4>
            
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Grandfather</label>
                <input
                  v-model="formData.paternalGrandfather"
                  type="text"
                  placeholder="Name"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Grandmother</label>
                <input
                  v-model="formData.paternalGrandmother"
                  type="text"
                  placeholder="Name"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors"
                >
              </div>
            </div>
          </div>

          <!-- Maternal Grandparents -->
          <div class="space-y-4">
            <h4 class="font-medium text-gray-800 border-b border-amber-200 pb-2">Maternal Grandparents</h4>
            
            <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Grandfather</label>
                <input
                  v-model="formData.maternalGrandfather"
                  type="text"
                  placeholder="Name"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors"
                >
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Grandmother</label>
                <input
                  v-model="formData.maternalGrandmother"
                  type="text"
                  placeholder="Name"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-amber-500 focus:border-amber-500 transition-colors"
                >
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- File upload section is commented out in original code -->

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-4 pt-6 border-t border-gray-200">
        <!-- Action buttons are commented out in original code -->
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive } from 'vue'
import { watch,onMounted } from 'vue';


// Updated emit to include 'update' event
const emit = defineEmits(['next', 'skip', 'update'])

const isSubmitting = ref(false)
const fileInput = ref(null)
const uploadedFiles = ref([])
const props = defineProps({
  data: {
    type: Object,
    default: () => ({})
  }
})

const formData = reactive({
  father: {
    name: '',
    gotra: '',
    nativePlace: '',
    vedaShaka: ''
  },
  mother: {
    name: '',
    maidenGotra: '',
    nativePlace: '',
    fatherName: ''
  },
  paternalGrandfather: '',
  paternalGrandmother: '',
  maternalGrandfather: '',
  maternalGrandmother: '',
  familyTreeFiles: []
})

// Keeping original file handling code
const handleFileUpload = (event) => {
  const files = Array.from(event.target.files)
  
  files.forEach(file => {
    if (file.size <= 10 * 1024 * 1024) { // 10MB limit
      uploadedFiles.value.push(file)
    }
  })
  
  formData.familyTreeFiles = uploadedFiles.value
  event.target.value = '' // Reset input
}

const removeFile = (index) => {
  uploadedFiles.value.splice(index, 1)
  formData.familyTreeFiles = uploadedFiles.value
}

// FIXED: Updated handleSubmit to emit original formData structure
const handleSubmit = async () => {
  isSubmitting.value = true
  
  try {
    // Emit the frontend data structure directly without transforming
    emit('update', formData)
    
    // Move to next step
    emit('next')
  } catch (error) {
    console.error('Error saving family lineage:', error)
  } finally {
    isSubmitting.value = false
  }
}

const skipStep = () => {
  emit('skip', { step: 3 })
}

// Add watch handler to emit updates on data change
watch(formData, () => {
  // Update parent component with data changes
  emit('update', formData)
}, { deep: true })

// In onMounted, initialize with existing data if available
onMounted(() => {
  if (props.data) {
    // Copy existing data to form
    formData.father = props.data.father || { name: '', gotra: '', nativePlace: '', vedaShaka: '' }
    formData.mother = props.data.mother || { name: '', maidenGotra: '', nativePlace: '', fatherName: '' }
    formData.paternalGrandfather = props.data.paternalGrandfather || ''
    formData.paternalGrandmother = props.data.paternalGrandmother || ''
    formData.maternalGrandfather = props.data.maternalGrandfather || ''
    formData.maternalGrandmother = props.data.maternalGrandmother || ''
  }
})
</script>