<template>
  <div class="profile-edit-container">
    <div class="max-w-4xl mx-auto py-8 px-4">
      <!-- Header -->
      <div class="mb-8">
        <h1 class="text-3xl font-bold text-gray-900 mb-2">Edit Profile</h1>
        <p class="text-gray-600">Update your personal and spiritual information</p>
      </div>

      <!-- Progress Bar -->
      <div class="mb-8">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-indigo-600">Profile Completion</span>
          <span class="text-sm text-gray-500">{{ completionPercentage }}%</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div 
            class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
            :style="{ width: completionPercentage + '%' }"
          ></div>
        </div>
      </div>

      <!-- Multi-step Form -->
      <div class="bg-white rounded-lg shadow-sm border">
        <!-- Step Navigation -->
        <div class="border-b border-gray-200 px-6 py-4">
          <nav class="flex space-x-8">
            <button
              v-for="(step, index) in steps"
              :key="step.id"
              @click="currentStep = index"
              :class="[
                'py-2 px-1 border-b-2 font-medium text-sm transition-colors',
                currentStep === index
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
              ]"
            >
              {{ step.title }}
            </button>
          </nav>
        </div>

        <!-- Step Content -->
        <div class="p-6">
          <!-- Personal Details -->
          <div v-if="currentStep === 0" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Personal Information</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Full Name</label>
                <input
                  v-model="profile.fullName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter your full name"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Date of Birth</label>
                <input
                  v-model="profile.dateOfBirth"
                  type="date"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Phone Number</label>
                <input
                  v-model="profile.phone"
                  type="tel"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="+91 9876543210"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Gender</label>
                <select
                  v-model="profile.gender"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="">Select Gender</option>
                  <option value="male">Male</option>
                  <option value="female">Female</option>
                  <option value="other">Other</option>
                </select>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
              <textarea
                v-model="profile.address"
                rows="3"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="Enter your complete address"
              ></textarea>
            </div>
          </div>

          <!-- Spiritual Information -->
          <div v-if="currentStep === 1" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Spiritual Information</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Gotra</label>
                <input
                  v-model="profile.gotra"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter your gotra"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Nakshatra</label>
                <input
                  v-model="profile.nakshatra"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter your nakshatra"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Rashi</label>
                <input
                  v-model="profile.rashi"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter your rashi"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Preferred Language</label>
                <select
                  v-model="profile.preferredLanguage"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                >
                  <option value="">Select Language</option>
                  <option value="kannada">Kannada</option>
                  <option value="hindi">Hindi</option>
                  <option value="sanskrit">Sanskrit</option>
                  <option value="tamil">Tamil</option>
                  <option value="telugu">Telugu</option>
                  <option value="english">English</option>
                </select>
              </div>
            </div>
          </div>

          <!-- Family Information -->
          <div v-if="currentStep === 2" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Family Information</h3>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Father's Name</label>
                <input
                  v-model="profile.fatherName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter father's name"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Mother's Name</label>
                <input
                  v-model="profile.motherName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter mother's name"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Spouse Name</label>
                <input
                  v-model="profile.spouseName"
                  type="text"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="Enter spouse name (if applicable)"
                />
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Number of Children</label>
                <input
                  v-model="profile.numberOfChildren"
                  type="number"
                  min="0"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                  placeholder="0"
                />
              </div>
            </div>
          </div>

          <!-- Preferences -->
          <div v-if="currentStep === 3" class="space-y-6">
            <h3 class="text-lg font-medium text-gray-900 mb-4">Seva Preferences</h3>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Preferred Sevas</label>
              <div class="grid grid-cols-2 md:grid-cols-3 gap-3">
                <label
                  v-for="seva in sevaTypes"
                  :key="seva.id"
                  class="flex items-center space-x-2 p-3 border border-gray-200 rounded-lg hover:bg-gray-50 cursor-pointer"
                >
                  <input
                    type="checkbox"
                    :value="seva.id"
                    v-model="profile.preferredSevas"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <span class="text-sm text-gray-700">{{ seva.name }}</span>
                </label>
              </div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">Notification Preferences</label>
              <div class="space-y-3">
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    v-model="profile.emailNotifications"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <span class="text-sm text-gray-700">Email notifications for events and sevas</span>
                </label>
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    v-model="profile.smsNotifications"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <span class="text-sm text-gray-700">SMS notifications for important updates</span>
                </label>
                <label class="flex items-center space-x-3">
                  <input
                    type="checkbox"
                    v-model="profile.whatsappNotifications"
                    class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  />
                  <span class="text-sm text-gray-700">WhatsApp notifications</span>
                </label>
              </div>
            </div>
          </div>
        </div>

        <!-- Navigation Buttons -->
        <div class="border-t border-gray-200 px-6 py-4 flex justify-between">
          <button
            v-if="currentStep > 0"
            @click="currentStep--"
            class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-md hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500"
          >
            Previous
          </button>
          <div v-else></div>
          
          <div class="flex space-x-3">
            <button
              @click="saveProfile"
              :disabled="isSaving"
              class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 disabled:opacity-50"
            >
              <span v-if="isSaving">Saving...</span>
              <span v-else>Save Changes</span>
            </button>
            
            <button
              v-if="currentStep < steps.length - 1"
              @click="currentStep++"
              class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const router = useRouter()
const { showToast } = useToast()

// Current step
const currentStep = ref(0)
const isSaving = ref(false)

// Steps configuration
const steps = ref([
  { id: 'personal', title: 'Personal' },
  { id: 'spiritual', title: 'Spiritual' },
  { id: 'family', title: 'Family' },
  { id: 'preferences', title: 'Preferences' }
])

// Profile data
const profile = ref({
  fullName: '',
  dateOfBirth: '',
  phone: '',
  gender: '',
  address: '',
  gotra: '',
  nakshatra: '',
  rashi: '',
  preferredLanguage: '',
  fatherName: '',
  motherName: '',
  spouseName: '',
  numberOfChildren: 0,
  preferredSevas: [],
  emailNotifications: true,
  smsNotifications: true,
  whatsappNotifications: false
})

// Seva types (mock data)
const sevaTypes = ref([
  { id: 1, name: 'Abhisheka' },
  { id: 2, name: 'Archana' },
  { id: 3, name: 'Annadana' },
  { id: 4, name: 'Pradakshina' },
  { id: 5, name: 'Bhajan' },
  { id: 6, name: 'Cleaning' }
])

// Computed properties
const completionPercentage = computed(() => {
  const fields = [
    profile.value.fullName,
    profile.value.dateOfBirth,
    profile.value.phone,
    profile.value.gender,
    profile.value.address,
    profile.value.gotra,
    profile.value.preferredLanguage
  ]
  
  const filledFields = fields.filter(field => field && field.toString().trim() !== '').length
  return Math.round((filledFields / fields.length) * 100)
})

// Methods
const loadProfile = async () => {
  try {
    // Mock API call - replace with actual API
    const mockProfile = {
      fullName: 'John Devotee',
      dateOfBirth: '1990-01-01',
      phone: '+91 9876543210',
      gender: 'male',
      address: '123 Temple Street, Bengaluru',
      gotra: 'Bharadvaja',
      nakshatra: 'Rohini',
      rashi: 'Vrishabha',
      preferredLanguage: 'kannada',
      fatherName: 'Father Name',
      motherName: 'Mother Name',
      spouseName: '',
      numberOfChildren: 0,
      preferredSevas: [1, 2],
      emailNotifications: true,
      smsNotifications: true,
      whatsappNotifications: false
    }
    
    Object.assign(profile.value, mockProfile)
  } catch (error) {
    console.error('Error loading profile:', error)
    showToast('Error loading profile data', 'error')
  }
}

const saveProfile = async () => {
  try {
    isSaving.value = true
    
    // Mock API call - replace with actual API
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    showToast('Profile updated successfully!', 'success')
    
    // Redirect back to dashboard
    router.push(`/entity/${route.params.id}/devotee/dashboard`)
  } catch (error) {
    console.error('Error saving profile:', error)
    showToast('Error saving profile', 'error')
  } finally {
    isSaving.value = false
  }
}

// Lifecycle
onMounted(() => {
  loadProfile()
})
</script>

<style scoped>
.profile-edit-container {
  min-height: calc(100vh - 64px);
  background-color: #f9fafb;
}
</style>