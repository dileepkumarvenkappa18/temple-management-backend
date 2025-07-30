<template>
  <div class="bg-white rounded-xl shadow-lg p-8">
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">
        Family Members
      </h2>
      <p class="text-gray-600">
        Add your family members and emergency contacts for better temple community engagement.
      </p>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-8">
      <!-- Spouse Details -->
      <div class="bg-gray-50 rounded-xl p-6">
        <div class="flex items-center mb-6">
          <HeartIcon class="w-6 h-6 text-indigo-600 mr-3" />
          <h3 class="text-lg font-semibold text-gray-900">Spouse Details</h3>
          <label class="ml-auto flex items-center cursor-pointer">
            <input
              v-model="hasSpouse"
              type="checkbox"
              class="sr-only"
            />
            <div class="relative">
              <div class="block bg-gray-300 w-10 h-6 rounded-full"></div>
              <div :class="[
                'absolute left-1 top-1 bg-white w-4 h-4 rounded-full transition-transform duration-200',
                hasSpouse ? 'transform translate-x-4 bg-indigo-600' : ''
              ]"></div>
            </div>
            <span class="ml-3 text-sm text-gray-600">{{ hasSpouse ? 'Married' : 'Single' }}</span>
          </label>
        </div>

        <div v-if="hasSpouse" class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <BaseInput
            v-model="formData.spouse.name"
            label="Spouse Name"
            placeholder="Enter spouse's full name"
            required
          />
          
          <BaseInput
            v-model="formData.spouse.email"
            type="email"
            label="Spouse Email"
            placeholder="spouse@example.com"
          />
          
          <BaseInput
            v-model="formData.spouse.phone"
            type="tel"
            label="Spouse Phone"
            placeholder="+91 00000 00000"
          />
          
          <BaseInput
            v-model="formData.spouse.dateOfBirth"
            type="date"
            label="Date of Birth"
          />
          
          <div>
            <label class="block text-sm font-semibold text-gray-900 mb-2">
              Gotra
            </label>
            <BaseSelect
              v-model="formData.spouse.gotra"
              :options="gotraOptions"
              placeholder="Select gotra"
            />
          </div>
          
          <BaseInput
            v-model="formData.spouse.nakshatra"
            label="Nakshatra"
            placeholder="Enter nakshatra"
          />
        </div>
      </div>

      <!-- Children Details -->
      <div class="bg-gray-50 rounded-xl p-6">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center">
            <UsersIcon class="w-6 h-6 text-indigo-600 mr-3" />
            <h3 class="text-lg font-semibold text-gray-900">Children</h3>
          </div>
          <button
            type="button"
            @click="addChild"
            class="flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors duration-200"
          >
            <PlusIcon class="w-4 h-4 mr-2" />
            Add Child
          </button>
        </div>

        <div v-if="formData.children.length === 0" class="text-center py-8 text-gray-500">
          <UserGroupIcon class="w-12 h-12 mx-auto mb-4 text-gray-300" />
          <p>No children added yet</p>
          <p class="text-sm">Click "Add Child" to include family members</p>
        </div>

        <div v-else class="space-y-6">
          <div 
            v-for="(child, index) in formData.children" 
            :key="index"
            class="bg-white rounded-lg p-6 border-2 border-gray-200"
          >
            <div class="flex items-center justify-between mb-4">
              <h4 class="font-semibold text-gray-900">Child {{ index + 1 }}</h4>
              <button
                type="button"
                @click="removeChild(index)"
                class="text-red-600 hover:text-red-800 p-1"
              >
                <TrashIcon class="w-5 h-5" />
              </button>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              <BaseInput
                v-model="child.name"
                label="Name"
                placeholder="Child's name"
                required
              />
              
              <BaseInput
                v-model="child.dateOfBirth"
                type="date"
                label="Date of Birth"
                required
              />
              
              <div>
                <label class="block text-sm font-semibold text-gray-900 mb-2">
                  Gender
                </label>
                <BaseSelect
                  v-model="child.gender"
                  :options="genderOptions"
                  placeholder="Select gender"
                />
              </div>
              
              <BaseInput
                v-model="child.education"
                label="Education/Class"
                placeholder="Current class or qualification"
              />
              
              <BaseInput
                v-model="child.interests"
                label="Interests"
                placeholder="Hobbies, sports, etc."
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Emergency Contacts -->
      <div class="bg-red-50 rounded-xl p-6">
        <div class="flex items-center justify-between mb-6">
          <div class="flex items-center">
            <PhoneIcon class="w-6 h-6 text-red-600 mr-3" />
            <h3 class="text-lg font-semibold text-gray-900">Emergency Contacts</h3>
          </div>
          <button
            type="button"
            @click="addEmergencyContact"
            class="flex items-center px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200"
          >
            <PlusIcon class="w-4 h-4 mr-2" />
            Add Contact
          </button>
        </div>

        <div v-if="formData.emergencyContacts.length === 0" class="text-center py-8 text-gray-500">
          <ExclamationTriangleIcon class="w-12 h-12 mx-auto mb-4 text-gray-300" />
          <p>No emergency contacts added</p>
          <p class="text-sm">Add at least one emergency contact for safety</p>
        </div>

        <div v-else class="space-y-6">
          <div 
            v-for="(contact, index) in formData.emergencyContacts" 
            :key="index"
            class="bg-white rounded-lg p-6 border-2 border-gray-200"
          >
            <div class="flex items-center justify-between mb-4">
              <h4 class="font-semibold text-gray-900">Emergency Contact {{ index + 1 }}</h4>
              <button
                type="button"
                @click="removeEmergencyContact(index)"
                class="text-red-600 hover:text-red-800 p-1"
              >
                <TrashIcon class="w-5 h-5" />
              </button>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <BaseInput
                v-model="contact.name"
                label="Name"
                placeholder="Contact name"
                required
              />
              
              <BaseInput
                v-model="contact.relationship"
                label="Relationship"
                placeholder="Brother, Sister, Friend"
                required
              />
              
              <BaseInput
                v-model="contact.phone"
                type="tel"
                label="Phone"
                placeholder="+91 00000 00000"
                required
              />
              
              <BaseInput
                v-model="contact.address"
                label="Address"
                placeholder="City, State"
              />
            </div>
          </div>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-4 pt-6 border-t border-gray-200">
        <!-- <button
          type="button"
          @click="$emit('previous')"
          class="flex items-center justify-center px-6 py-3 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 font-semibold transition-colors duration-200"
        >
          <ChevronLeftIcon class="w-5 h-5 mr-2" />
          Previous Step
        </button> -->
        
        <!-- <div class="flex gap-4 flex-1">
          <button
            type="button"
            @click="handleSkip"
            class="flex-1 px-6 py-3 text-indigo-600 font-semibold rounded-xl hover:bg-indigo-50 transition-colors duration-200"
          >
            Skip for Now
          </button>
          
          <button
            type="submit"
            :disabled="isSubmitting"
            class="flex-1 bg-indigo-600 text-white px-6 py-3 rounded-xl hover:bg-indigo-700 font-semibold transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
          >
            <BaseLoader v-if="isSubmitting" size="sm" color="white" class="mr-2" />
            Continue
          </button>
        </div> -->
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { 
  ChevronLeftIcon,
  HeartIcon,
  UsersIcon,
  UserGroupIcon,
  PhoneIcon,
  PlusIcon,
  TrashIcon,
  ExclamationTriangleIcon
} from '@heroicons/vue/24/outline'
import BaseInput from '../common/BaseInput.vue'
import BaseSelect from '../common/BaseSelect.vue'
import BaseLoader from '../common/BaseLoader.vue'

const emit = defineEmits(['next', 'previous', 'skip'])

const isSubmitting = ref(false)
const hasSpouse = ref(false)

const formData = reactive({
  spouse: {
    name: '',
    email: '',
    phone: '',
    dateOfBirth: '',
    gotra: '',
    nakshatra: ''
  },
  children: [],
  emergencyContacts: []
})

const gotraOptions = [
  { value: 'bharadvaja', label: 'Bharadvaja' },
  { value: 'kashyapa', label: 'Kashyapa' },
  { value: 'atri', label: 'Atri' },
  { value: 'vasishta', label: 'Vasishta' },
  { value: 'gautama', label: 'Gautama' },
  { value: 'jamadagni', label: 'Jamadagni' },
  { value: 'vishvamitra', label: 'Vishvamitra' }
]

const genderOptions = [
  { value: 'male', label: 'Male' },
  { value: 'female', label: 'Female' }
]

const addChild = () => {
  formData.children.push({
    name: '',
    dateOfBirth: '',
    gender: '',
    education: '',
    interests: ''
  })
}

const removeChild = (index) => {
  formData.children.splice(index, 1)
}

const addEmergencyContact = () => {
  formData.emergencyContacts.push({
    name: '',
    relationship: '',
    phone: '',
    address: ''
  })
}

const removeEmergencyContact = (index) => {
  formData.emergencyContacts.splice(index, 1)
}

const handleSubmit = async () => {
  isSubmitting.value = true
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    const submitData = {
      ...formData,
      hasSpouse: hasSpouse.value
    }
    
    emit('next', submitData)
  } catch (error) {
    console.error('Error saving family details:', error)
  } finally {
    isSubmitting.value = false
  }
}

const handleSkip = () => {
  emit('skip')
}
</script>