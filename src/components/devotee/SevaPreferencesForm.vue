<template>
  <div class="bg-white rounded-xl shadow-lg p-8">
    <!-- Header -->
    <div class="mb-8">
      <h2 class="text-2xl font-bold text-gray-900 mb-2">
        Seva Preferences
      </h2>
      <p class="text-gray-600">
        Select your preferred sevas and donation interests to personalize your temple experience.
      </p>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-8">
      <!-- Favorite Sevas -->
      <div>
        <label class="block text-sm font-semibold text-gray-900 mb-4">
          Favorite Sevas <span class="text-gray-500 font-normal">(Select multiple)</span>
        </label>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          <div 
            v-for="seva in sevaTypes" 
            :key="seva.id"
            class="relative"
          >
            <input
              :id="`seva-${seva.id}`"
              v-model="formData.favoriteSevas"
              :value="seva.id"
              type="checkbox"
              class="sr-only peer"
            />
            <label
              :for="`seva-${seva.id}`"
              class="flex items-start p-4 border-2 border-gray-200 rounded-xl cursor-pointer transition-all duration-200 hover:border-indigo-300 peer-checked:border-indigo-600 peer-checked:bg-indigo-50"
            >
              <div class="flex-1">
                <div class="flex items-center mb-2">
                  <component 
                    :is="seva.icon" 
                    class="w-5 h-5 text-indigo-600 mr-2 flex-shrink-0"
                  />
                  <h3 class="font-semibold text-gray-900 text-sm">
                    {{ seva.name }}
                  </h3>
                </div>
                <p class="text-xs text-gray-600">
                  {{ seva.description }}
                </p>
                <div class="mt-2 flex items-center text-xs text-indigo-600">
                  <span class="bg-indigo-100 px-2 py-1 rounded-full">
                    ₹{{ seva.minAmount }}+
                  </span>
                </div>
              </div>
              <div class="ml-3 flex items-center">
                <div class="w-5 h-5 border-2 border-gray-300 rounded peer-checked:border-indigo-600 peer-checked:bg-indigo-600 peer-checked:text-white flex items-center justify-center">
                  <svg class="w-3 h-3 hidden peer-checked:block" fill="currentColor" viewBox="0 0 20 20">
                    <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                  </svg>
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Donation Interests -->
      <div>
        <label class="block text-sm font-semibold text-gray-900 mb-4">
          Donation Interests <span class="text-gray-500 font-normal">(Select categories you're interested in)</span>
        </label>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div 
            v-for="category in donationCategories" 
            :key="category.id"
            class="relative"
          >
            <input
              :id="`donation-${category.id}`"
              v-model="formData.donationInterests"
              :value="category.id"
              type="checkbox"
              class="sr-only peer"
            />
            <label
              :for="`donation-${category.id}`"
              class="flex items-center p-4 border-2 border-gray-200 rounded-xl cursor-pointer transition-all duration-200 hover:border-indigo-300 peer-checked:border-indigo-600 peer-checked:bg-indigo-50"
            >
              <component 
                :is="category.icon" 
                class="w-6 h-6 text-indigo-600 mr-4 flex-shrink-0"
              />
              <div class="flex-1">
                <h3 class="font-semibold text-gray-900 text-sm mb-1">
                  {{ category.name }}
                </h3>
                <p class="text-xs text-gray-600">
                  {{ category.description }}
                </p>
              </div>
              <div class="w-5 h-5 border-2 border-gray-300 rounded peer-checked:border-indigo-600 peer-checked:bg-indigo-600 peer-checked:text-white flex items-center justify-center">
                <svg class="w-3 h-3 hidden peer-checked:block" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Preferred Frequency -->
      <div>
        <label class="block text-sm font-semibold text-gray-900 mb-4">
          How often would you like to participate in temple activities?
        </label>
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div 
            v-for="frequency in frequencies" 
            :key="frequency.value"
            class="relative"
          >
            <input
              :id="`frequency-${frequency.value}`"
              v-model="formData.preferredFrequency"
              :value="frequency.value"
              type="radio"
              class="sr-only peer"
            />
            <label
              :for="`frequency-${frequency.value}`"
              class="block p-4 border-2 border-gray-200 rounded-xl cursor-pointer transition-all duration-200 hover:border-indigo-300 peer-checked:border-indigo-600 peer-checked:bg-indigo-50 text-center"
            >
              <div class="font-semibold text-gray-900 text-sm mb-1">
                {{ frequency.label }}
              </div>
              <div class="text-xs text-gray-600">
                {{ frequency.description }}
              </div>
            </label>
          </div>
        </div>
      </div>

      <!-- Special Interests -->
      <div>
        <label for="specialInterests" class="block text-sm font-semibold text-gray-900 mb-2">
          Special Interests or Notes
        </label>
        <textarea
          id="specialInterests"
          v-model="formData.specialInterests"
          rows="4"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200 resize-none"
          placeholder="Any specific festivals, rituals, or temple activities you're particularly interested in..."
        ></textarea>
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
        
        <div class="flex gap-4 flex-1">
          <!-- <button
            type="button"
            @click="handleSkip"
            class="flex-1 px-6 py-3 text-indigo-600 font-semibold rounded-xl hover:bg-indigo-50 transition-colors duration-200"
          >
            Skip for Now
          </button> -->
          
          <!-- <button
            type="submit"
            :disabled="isSubmitting"
            class="flex-1 bg-indigo-600 text-white px-6 py-3 rounded-xl hover:bg-indigo-700 font-semibold transition-colors duration-200 disabled:opacity-50 disabled:cursor-not-allowed flex items-center justify-center"
          >
            <BaseLoader v-if="isSubmitting" size="sm" color="white" class="mr-2" />
            Continue
          </button> -->
        </div>
      </div>
    </form>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { ChevronLeftIcon } from '@heroicons/vue/24/outline'
import { 
  HandRaisedIcon, 
  SparklesIcon, 
  SunIcon, 
  HeartIcon,
  HomeIcon,
  GiftIcon,
  AcademicCapIcon,
  BuildingLibraryIcon 
} from '@heroicons/vue/24/outline'
import BaseLoader from '../common/BaseLoader.vue'

const emit = defineEmits(['next', 'previous', 'skip'])

const isSubmitting = ref(false)

const formData = reactive({
  favoriteSevas: [],
  donationInterests: [],
  preferredFrequency: '',
  specialInterests: ''
})

const sevaTypes = [
  {
    id: 'abhisheka',
    name: 'Abhisheka',
    description: 'Sacred bathing ritual of the deity',
    icon: SparklesIcon,
    minAmount: 51
  },
  {
    id: 'arti',
    name: 'Arti',
    description: 'Divine light ceremony',
    icon: SunIcon,
    minAmount: 21
  },
  {
    id: 'annadana',
    name: 'Annadana',
    description: 'Food donation service',
    icon: HeartIcon,
    minAmount: 101
  },
  {
    id: 'archana',
    name: 'Archana',
    description: 'Personal prayer offering',
    icon: HandRaisedIcon,
    minAmount: 11
  },
  {
    id: 'kalyanam',
    name: 'Kalyanam',
    description: 'Divine marriage ceremony',
    icon: GiftIcon,
    minAmount: 501
  },
  {
    id: 'homam',
    name: 'Homam',
    description: 'Sacred fire ritual',
    icon: SunIcon,
    minAmount: 251
  }
]

const donationCategories = [
  {
    id: 'temple-maintenance',
    name: 'Temple Maintenance',
    description: 'Upkeep and renovation of temple premises',
    icon: HomeIcon
  },
  {
    id: 'annadana',
    name: 'Annadana Program',
    description: 'Free food distribution to devotees',
    icon: HeartIcon
  },
  {
    id: 'festivals',
    name: 'Festival Celebrations',
    description: 'Special occasions and religious festivals',
    icon: SparklesIcon
  },
  {
    id: 'education',
    name: 'Religious Education',
    description: 'Vedic studies and spiritual learning programs',
    icon: AcademicCapIcon
  },
  {
    id: 'construction',
    name: 'Temple Construction',
    description: 'New building projects and expansions',
    icon: BuildingLibraryIcon
  },
  {
    id: 'general',
    name: 'General Donations',
    description: 'Flexible contributions for temple needs',
    icon: GiftIcon
  }
]

const frequencies = [
  {
    value: 'daily',
    label: 'Daily',
    description: 'Every day'
  },
  {
    value: 'weekly',
    label: 'Weekly',
    description: 'Once a week'
  },
  {
    value: 'monthly',
    label: 'Monthly',
    description: 'Once a month'
  },
  {
    value: 'festivals',
    label: 'Festivals Only',
    description: 'Special occasions'
  }
]

const handleSubmit = async () => {
  isSubmitting.value = true
  
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    emit('next', formData)
  } catch (error) {
    console.error('Error saving seva preferences:', error)
  } finally {
    isSubmitting.value = false
  }
}

const handleSkip = () => {
  emit('skip')
}
</script>