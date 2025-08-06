<template>
  <div class="min-h-screen bg-gray-50 py-8">
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8">
      <!-- Header -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-8 p-6">
        <div class="flex items-center justify-between mb-4">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Complete Your Profile</h1>
            <p class="text-gray-600 mt-1">Build your spiritual profile step by step</p>
          </div>
          <button
            @click="skipToEnd"
            class="text-indigo-600 hover:text-indigo-500 font-medium text-sm"
          >
            Skip All →
          </button>
        </div>
        
        <!-- Progress Indicator -->
        <div>
          <div class="flex items-center justify-between text-sm mb-2">
            <span class="font-medium text-indigo-600">
              {{ completedSteps.filter(s => s.completed).length }} of {{ totalSteps }} completed
            </span>
            <span class="text-gray-600">{{ completionPercentage }}% Complete</span>
          </div>
          <div class="w-full bg-gray-200 rounded-full h-2">
            <div 
              class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
              :style="{ width: `${completionPercentage}%` }"
            ></div>
          </div>
        </div>
      </div>

      <!-- Step Navigation -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-6">
        <div class="px-6 py-4 border-b border-gray-100">
          <nav class="flex space-x-8 overflow-x-auto">
            <button
              v-for="(step, index) in steps"
              :key="step.id"
              @click="goToStep(index)"
              :class="[
                'flex items-center space-x-2 py-2 px-3 rounded-lg text-sm font-medium whitespace-nowrap transition-colors',
                currentStep === index 
                  ? 'bg-indigo-100 text-indigo-700' 
                  : step.completed 
                    ? 'text-green-600 hover:bg-green-50' 
                    : 'text-gray-500 hover:bg-gray-50'
              ]"
            >
              <div :class="[
                'w-5 h-5 rounded-full flex items-center justify-center text-xs',
                currentStep === index 
                  ? 'bg-indigo-600 text-white' 
                  : step.completed 
                    ? 'bg-green-100 text-green-600' 
                    : 'bg-gray-200 text-gray-500'
              ]">
                <svg v-if="step.completed" xmlns="http://www.w3.org/2000/svg" class="w-3 h-3" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
                <span v-else>{{ index + 1 }}</span>
              </div>
              <span>{{ step.title }}</span>
            </button>
          </nav>
        </div>
      </div>

      <!-- Form Container -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200">
        <div class="p-6">
          <!-- Loading State -->
          <div v-if="loading" class="flex flex-col items-center justify-center py-12">
            <svg class="animate-spin h-10 w-10 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <p class="mt-4 text-gray-600">Loading...</p>
          </div>
          
          <!-- Error State -->
          <div v-else-if="error" class="bg-red-50 p-4 rounded-lg">
            <div class="flex">
              <svg class="h-5 w-5 text-red-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
              </svg>
              <div class="ml-3">
                <h3 class="text-sm font-medium text-red-800">An error occurred</h3>
                <p class="mt-2 text-sm text-red-700">{{ error }}</p>
                <div class="mt-4">
                  <button @click="retryLoading" class="text-sm font-medium text-red-700 hover:text-red-600">
                    Try again
                  </button>
                </div>
              </div>
            </div>
          </div>

          <!-- Step Content -->
          <div v-else class="min-h-[400px]">
            <!-- Personal Details -->
            <div v-if="currentStep === 0">
              <PersonalDetailsForm 
                :data="stepData.personal" 
                @update="updateStepData('personal', $event)"
                @update-data="updateStepData('personal', $event)"
                @next="nextStep"
              />
            </div>

            <!-- Spiritual Info -->
            <div v-else-if="currentStep === 1">
              <SpiritualInfoForm 
                :data="stepData.spiritual" 
                @update="updateStepData('spiritual', $event)"
                @update-data="updateStepData('spiritual', $event)"
                @next="nextStep"
              />
            </div>

            <!-- Family Lineage -->
            <div v-else-if="currentStep === 2">
              <FamilyLineageForm 
                :data="stepData.lineage" 
                @update="updateStepData('lineage', $event)"
                @update-data="updateStepData('lineage', $event)"
                @next="nextStep"
              />
            </div>

            <!-- Seva Preferences -->
            <div v-else-if="currentStep === 3">
              <SevaPreferencesForm 
                :data="stepData.preferences" 
                @update="updateStepData('preferences', $event)"
                @update-data="updateStepData('preferences', $event)"
                @next="nextStep"
              />
            </div>

            <!-- Family Members -->
            <div v-else-if="currentStep === 4">
              <FamilyMembersForm 
                :data="stepData.family" 
                @update="updateStepData('family', $event)"
                @update-data="updateStepData('family', $event)"
                @next="nextStep"
              />
            </div>

            <!-- Notes & Attachments -->
            <div v-else-if="currentStep === 5">
              <NotesAttachmentsForm 
                :data="stepData.notes" 
                @update="updateStepData('notes', $event)"
                @update-data="updateStepData('notes', $event)"
                @submit="completeProfile"
              />
            </div>
          </div>

          <!-- Navigation Buttons -->
          <div v-if="!loading && !error" class="flex items-center justify-between pt-6 border-t border-gray-100 mt-8">
            <button
              @click="previousStep"
              :disabled="currentStep === 0"
              class="flex items-center space-x-2 px-4 py-2 text-gray-600 hover:text-gray-900 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              <svg xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
              </svg>
              <span>Previous</span>
            </button>

            <div class="flex items-center space-x-3">
              <button
                @click="skipStep"
                class="px-4 py-2 text-gray-500 hover:text-gray-700 font-medium"
              >
                Skip This Step
              </button>
              
              <button
                @click="nextStep"
                :disabled="isSubmitting"
                class="flex items-center space-x-2 bg-indigo-600 hover:bg-indigo-700 disabled:bg-indigo-400 text-white px-6 py-2 rounded-lg font-medium transition-colors"
              >
                <span v-if="isSubmitting">Saving...</span>
                <span v-else-if="currentStep === totalSteps - 1">Complete Profile</span>
                <span v-else>Next Step</span>
                <svg v-if="!isSubmitting" xmlns="http://www.w3.org/2000/svg" class="w-4 h-4" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Help Section -->
      <div class="mt-6 bg-indigo-50 rounded-xl p-6">
        <div class="flex items-start space-x-3">
          <svg xmlns="http://www.w3.org/2000/svg" class="w-5 h-5 text-indigo-500 mt-0.5" viewBox="0 0 20 20" fill="currentColor">
            <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
          </svg>
          <div>
            <h3 class="text-sm font-medium text-indigo-900 mb-1">Need Help?</h3>
            <p class="text-sm text-indigo-700 mb-3">
              All steps are optional and can be completed later from your dashboard. 
              However, a complete profile helps temple administrators serve you better.
            </p>
            <button class="text-sm text-indigo-600 hover:text-indigo-500 font-medium">
              Contact Support →
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useDevoteeStore } from '@/stores/devotee'
import { useToast } from '@/composables/useToast'
import PersonalDetailsForm from '@/components/devotee/PersonalDetailsForm.vue'
import SpiritualInfoForm from '@/components/devotee/SpiritualInfoForm.vue'
import FamilyLineageForm from '@/components/devotee/FamilyLineageForm.vue'
import SevaPreferencesForm from '@/components/devotee/SevaPreferencesForm.vue'
import FamilyMembersForm from '@/components/devotee/FamilyMembersForm.vue'
import NotesAttachmentsForm from '@/components/devotee/NotesAttachmentsForm.vue'

const router = useRouter()
const devoteeStore = useDevoteeStore()
const toast = useToast()

// State
const loading = ref(true)
const error = ref(null)
const currentStep = ref(0)
const isSubmitting = ref(false)
// Update the stepData initialization to include empty objects for all sections
const stepData = ref({
  personal: {},
  spiritual: {},
  lineage: {
    father: { name: '', gotra: '', nativePlace: '', vedaShaka: '' },
    mother: { name: '', maidenGotra: '', nativePlace: '', fatherName: '' },
    paternalGrandfather: '',
    paternalGrandmother: '',
    maternalGrandfather: '',
    maternalGrandmother: ''
  },
  preferences: {
    favoriteSevas: [],
    donationInterests: [],
    preferredFrequency: '',
    specialInterests: ''
  },
  family: {
    spouse: { name: '', email: '', phone: '', dateOfBirth: '', gotra: '', nakshatra: '' },
    children: [],
    emergencyContacts: []
  },
  notes: {
    healthNotes: '',
    allergies: '',
    dietaryRestrictions: '',
    sankalpa: '',
    additionalNotes: '',
    documents: [],
    profilePhoto: null,
    familyPhoto: null,
    shareProfile: true,
    receiveNotifications: true
  }
})

// Steps configuration
const steps = ref([
  { 
    id: 'personal', 
    title: 'Personal Details', 
    component: 'PersonalDetailsForm',
    completed: false 
  },
  { 
    id: 'spiritual', 
    title: 'Spiritual Info', 
    component: 'SpiritualInfoForm',
    completed: false 
  },
  { 
    id: 'lineage', 
    title: 'Family Lineage', 
    component: 'FamilyLineageForm',
    completed: false 
  },
  { 
    id: 'preferences', 
    title: 'Seva Preferences', 
    component: 'SevaPreferencesForm',
    completed: false 
  },
  { 
    id: 'family', 
    title: 'Family Members', 
    component: 'FamilyMembersForm',
    completed: false 
  },
  { 
    id: 'notes', 
    title: 'Notes & Attachments', 
    component: 'NotesAttachmentsForm',
    completed: false 
  }
])

// Computed
const totalSteps = computed(() => steps.value.length)

const completionPercentage = computed(() => {
  const completedStepsCount = steps.value.filter(step => step.completed).length
  return Math.round((completedStepsCount / totalSteps.value) * 100)
})

const completedSteps = computed(() => steps.value)

// Methods
const goToStep = (stepIndex) => {
  if (stepIndex >= 0 && stepIndex < totalSteps.value) {
    currentStep.value = stepIndex
  }
}

// Add to nextStep function
const nextStep = async () => {
  isSubmitting.value = true
  
  try {
    // Save current step data to store
    const stepId = steps.value[currentStep.value].id
    console.log(`⏭️ Moving to next step, saving data for ${stepId}:`, stepData.value[stepId])
    await devoteeStore.updateProfileStep(stepId, stepData.value[stepId])
    
    // Mark current step as completed
    steps.value[currentStep.value].completed = true
    
    if (currentStep.value < totalSteps.value - 1) {
      // Go to next step
      currentStep.value++
    } else {
      // Complete profile creation
      await completeProfile()
    }
  } catch (error) {
    toast.error('Failed to save step data')
    console.error('Error saving step:', error)
  } finally {
    isSubmitting.value = false
  }
}

const previousStep = () => {
  if (currentStep.value > 0) {
    currentStep.value--
  }
}

const skipStep = () => {
  if (currentStep.value < totalSteps.value - 1) {
    currentStep.value++
  } else {
    skipToEnd()
  }
}

const skipToEnd = async () => {
  try {
    isSubmitting.value = true
    
    // Mark user has been through the profile creation process
    const result = await devoteeStore.completeProfile()
    
    if (result && result.redirectPath) {
      router.push(result.redirectPath)
    } else {
      // Fallback redirect
      const entityId = devoteeStore.getCurrentEntityId()
      if (entityId) {
        router.push(`/entity/${entityId}/devotee/dashboard`)
      } else {
        router.push('/devotee/dashboard')
      }
    }
  } catch (error) {
    toast.error('Failed to complete profile')
    console.error('Error skipping to end:', error)
  } finally {
    isSubmitting.value = false
  }
}

const updateStepData = (stepId, data) => {
  console.log(`💾 UPDATING STEP DATA for ${stepId}:`, data);
  stepData.value[stepId] = data;
}

const completeProfile = async () => {
  isSubmitting.value = true
  
  try {
    // Log current data state before completing
    console.log("📊 Current stepData before completion:", JSON.stringify(stepData.value, null, 2));
    
    // Build the complete profile object from all steps (for logging purposes)
    const completeProfileData = {
      personal: stepData.value.personal || {},
      spiritual: stepData.value.spiritual || {},
      lineage: stepData.value.lineage || {},
      preferences: stepData.value.preferences || {},
      family: stepData.value.family || {},
      notes: stepData.value.notes || {}
    }
    
    // Log the data that will be sent to the store
    console.log(`📤 COMPLETE PROFILE DATA GOING TO STORE:`, JSON.stringify(completeProfileData, null, 2));
    
    // IMPORTANT: Skip re-updating each step since we've already updated them individually
    // This was causing the issue - each step was being updated when moving through the form
    // and then again here, which could lead to data inconsistencies
    
    // Instead, directly call the store method to save everything to backend
    const result = await devoteeStore.completeProfile()
    
    console.log("✅ Profile completion result:", result);
    
    if (result && result.success) {
      // Use the redirect path from the result
      if (result.redirectPath) {
        router.push(result.redirectPath)
      } else {
        // Fallback redirect
        const entityId = devoteeStore.getCurrentEntityId()
        if (entityId) {
          router.push(`/entity/${entityId}/devotee/dashboard`)
        } else {
          router.push('/devotee/dashboard')
        }
      }
    }
  } catch (err) {
    toast.error('Failed to complete profile')
    console.error('Error completing profile:', err)
  } finally {
    isSubmitting.value = false
  }
}

const retryLoading = async () => {
  error.value = null
  loading.value = true
  try {
    await loadProfileData()
  } catch (err) {
    error.value = 'Failed to load profile data. Please try again.'
  } finally {
    loading.value = false
  }
}

// Load existing profile data
const loadProfileData = async () => {
  try {
    loading.value = true
    error.value = null
    
    // Load profile data from API
    const data = await devoteeStore.loadProfileData()
    
    if (data) {
      // Populate step data from API response
      stepData.value = {
        personal: data.personal || {},
        spiritual: data.spiritual || {},
        lineage: data.lineage || {},
        preferences: data.preferences || {},
        family: data.family || {},
        notes: data.notes || {}
      }
      
      // Update completed steps based on data
      steps.value.forEach((step, index) => {
        const stepId = step.id
        const hasData = stepData.value[stepId] && Object.keys(stepData.value[stepId]).length > 0
        steps.value[index].completed = hasData
      })
    }
  } catch (err) {
    // Only set error for serious issues, not 404s (which are normal for new users)
    if (err?.response?.status !== 404) {
      error.value = 'Failed to load profile data. Please try again.'
    }
    console.error('Error loading profile data:', err)
  } finally {
    loading.value = false
  }
}

// Lifecycle
onMounted(async () => {
  await loadProfileData()
})
</script>