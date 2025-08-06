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
          
          <!-- Fixed: Using computed property for date conversion -->
          <div>
            <label class="block text-sm font-semibold text-gray-900 mb-2">
              Date of Birth
            </label>
            <input
              :value="formatDateForInput(formData.spouse.dateOfBirth)"
              @input="updateSpouseDOB($event.target.value)"
              type="date"
              class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors duration-200"
            />
          </div>
          
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
              
              <!-- Fixed: Using computed property for date conversion -->
              <div>
                <label class="block text-sm font-semibold text-gray-900 mb-2">
                  Date of Birth *
                </label>
                <input
                  :value="formatDateForInput(child.dateOfBirth)"
                  @input="updateChildDOB(index, $event.target.value)"
                  type="date"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent transition-colors duration-200"
                />
              </div>
              
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
        <!-- Action buttons are commented out in the original code -->
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref, watch } from 'vue'
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

// Define props to receive initial data
const props = defineProps({
  data: {
    type: Object,
    default: () => ({
      spouse: {
        name: '',
        email: '',
        phone: '',
        dateOfBirth: '',
        gotra: '',
        nakshatra: ''
      },
      children: [],
      emergencyContacts: [],
      hasSpouse: false
    })
  }
})

// Updated to include 'update' event
const emit = defineEmits(['next', 'previous', 'skip', 'update', 'updateData'])

const isSubmitting = ref(false)
const hasSpouse = ref(props.data.hasSpouse || false)

// Initialize formData with props data
const formData = reactive({
  spouse: { ...props.data.spouse },
  children: props.data.children ? [...props.data.children] : [],
  emergencyContacts: props.data.emergencyContacts ? [...props.data.emergencyContacts] : []
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

// Date formatting utilities
const formatDateForInput = (dateValue) => {
  if (!dateValue) return ''
  
  try {
    // Handle ISO string format (from backend)
    if (typeof dateValue === 'string' && dateValue.includes('T')) {
      const date = new Date(dateValue)
      if (isNaN(date.getTime())) return ''
      return date.toISOString().split('T')[0]
    }
    
    // Handle Date object
    if (dateValue instanceof Date) {
      if (isNaN(dateValue.getTime())) return ''
      return dateValue.toISOString().split('T')[0]
    }
    
    // Handle yyyy-mm-dd format (already correct)
    if (typeof dateValue === 'string' && /^\d{4}-\d{2}-\d{2}$/.test(dateValue)) {
      return dateValue
    }
    
    return ''
  } catch (error) {
    console.warn('Date formatting error:', error)
    return ''
  }
}

const formatDateForBackend = (dateString) => {
  if (!dateString) return ''
  
  try {
    // Convert yyyy-mm-dd to ISO string
    const date = new Date(dateString + 'T00:00:00Z')
    if (isNaN(date.getTime())) return ''
    return date.toISOString()
  } catch (error) {
    console.warn('Date conversion error:', error)
    return ''
  }
}

// Date update handlers
const updateSpouseDOB = (dateString) => {
  formData.spouse.dateOfBirth = formatDateForBackend(dateString)
  emitUpdate()
}

const updateChildDOB = (index, dateString) => {
  if (formData.children[index]) {
    formData.children[index].dateOfBirth = formatDateForBackend(dateString)
    emitUpdate()
  }
}

const addChild = () => {
  formData.children.push({
    name: '',
    dateOfBirth: '',
    gender: '',
    education: '',
    interests: ''
  })
  // Emit update immediately when adding child
  emitUpdate()
}

const removeChild = (index) => {
  formData.children.splice(index, 1)
  // Emit update immediately when removing child
  emitUpdate()
}

const addEmergencyContact = () => {
  formData.emergencyContacts.push({
    name: '',
    relationship: '',
    phone: '',
    address: ''
  })
  // Emit update immediately when adding contact
  emitUpdate()
}

const removeEmergencyContact = (index) => {
  formData.emergencyContacts.splice(index, 1)
  // Emit update immediately when removing contact
  emitUpdate()
}

// Helper function to emit updates
const emitUpdate = () => {
  const familyData = {
    ...formData,
    hasSpouse: hasSpouse.value
  }
  emit('update', familyData)
  emit('updateData', familyData)
}

// Updated handleSubmit to emit both update and next events with frontend format data
const handleSubmit = async () => {
  isSubmitting.value = true
  
  try {
    // Create the data object with hasSpouse flag
    const familyData = {
      ...formData,
      hasSpouse: hasSpouse.value
    }
    
    // Emit update with the frontend format data
    emit('update', familyData)
    emit('updateData', familyData)
    
    // Also emit next to proceed to the next step
    emit('next')
  } catch (error) {
    console.error('Error saving family details:', error)
  } finally {
    isSubmitting.value = false
  }
}

const handleSkip = () => {
  emit('skip')
}

// Watch for changes in formData and hasSpouse (excluding date fields to avoid conflicts)
watch([() => formData.spouse.name, () => formData.spouse.email, () => formData.spouse.phone, 
       () => formData.spouse.gotra, () => formData.spouse.nakshatra, hasSpouse], () => {
  emitUpdate()
}, { deep: true })

// Watch for children changes (excluding date fields)
watch(() => formData.children.map(child => ({
  name: child.name,
  gender: child.gender,
  education: child.education,
  interests: child.interests
})), () => {
  emitUpdate()
}, { deep: true })

// Watch for emergency contacts changes
watch(() => formData.emergencyContacts, () => {
  emitUpdate()
}, { deep: true })

// Watch for changes in hasSpouse specifically
watch(hasSpouse, () => {
  emitUpdate()
})
</script>