<!-- src/components/devotee/SpiritualInfoForm.vue -->
<template>
  <div class="bg-white rounded-2xl shadow-md p-6">
    <!-- Header -->
    <div class="mb-6">
      <h2 class="text-xl font-semibold text-gray-900 mb-2">Spiritual Information</h2>
      <p class="text-gray-600">Help us understand your spiritual background and traditions</p>
    </div>

    <!-- Form -->
    <form @submit.prevent="handleSubmit" class="space-y-6">
      <!-- Gotra -->
      <div>
        <label for="gotra" class="block text-sm font-medium text-gray-700 mb-2">
          Gotra
          <span class="text-gray-400 text-xs ml-1">(Your family lineage)</span>
        </label>
        <input
          id="gotra"
          v-model="formData.gotra"
          type="text"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          placeholder="e.g., Bharadwaja, Kashyapa, Vasishtha"
          list="gotra-suggestions"
        />
        <datalist id="gotra-suggestions">
          <option v-for="gotra in popularGotras" :key="gotra" :value="gotra">{{ gotra }}</option>
        </datalist>
        <p v-if="errors.gotra" class="mt-1 text-sm text-red-600">{{ errors.gotra }}</p>
      </div>

      <!-- Nakshatra & Rashi -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label for="nakshatra" class="block text-sm font-medium text-gray-700 mb-2">
            Nakshatra
            <span class="text-gray-400 text-xs ml-1">(Birth star)</span>
          </label>
          <select
            id="nakshatra"
            v-model="formData.nakshatra"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          >
            <option value="">Select Nakshatra</option>
            <option v-for="nakshatra in nakshatras" :key="nakshatra" :value="nakshatra">
              {{ nakshatra }}
            </option>
          </select>
          <p v-if="errors.nakshatra" class="mt-1 text-sm text-red-600">{{ errors.nakshatra }}</p>
        </div>

        <div>
          <label for="rashi" class="block text-sm font-medium text-gray-700 mb-2">
            Rashi
            <span class="text-gray-400 text-xs ml-1">(Zodiac sign)</span>
          </label>
          <select
            id="rashi"
            v-model="formData.rashi"
            class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
          >
            <option value="">Select Rashi</option>
            <option v-for="rashi in rashis" :key="rashi" :value="rashi">
              {{ rashi }}
            </option>
          </select>
          <p v-if="errors.rashi" class="mt-1 text-sm text-red-600">{{ errors.rashi }}</p>
        </div>
      </div>

      <!-- Lagna -->
      <div>
        <label for="lagna" class="block text-sm font-medium text-gray-700 mb-2">
          Lagna (Ascendant)
          <span class="text-gray-400 text-xs ml-1">(Rising sign at birth)</span>
        </label>
        <select
          id="lagna"
          v-model="formData.lagna"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
        >
          <option value="">Select Lagna</option>
          <option v-for="lagna in rashis" :key="lagna" :value="lagna">
            {{ lagna }}
          </option>
        </select>
        <p v-if="errors.lagna" class="mt-1 text-sm text-red-600">{{ errors.lagna }}</p>
      </div>

      <!-- Veda Shaka -->
      <div>
        <label for="vedaShaka" class="block text-sm font-medium text-gray-700 mb-2">
          Veda Shaka
          <span class="text-gray-400 text-xs ml-1">(Vedic branch/tradition)</span>
        </label>
        <select
          id="vedaShaka"
          v-model="formData.vedaShaka"
          class="w-full px-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
        >
          <option value="">Select Veda Shaka</option>
          <option v-for="shaka in vedaShakas" :key="shaka" :value="shaka">
            {{ shaka }}
          </option>
        </select>
        <p v-if="errors.vedaShaka" class="mt-1 text-sm text-red-600">{{ errors.vedaShaka }}</p>
      </div>

      <!-- Additional Information -->
      <div class="bg-indigo-50 rounded-xl p-4 border border-indigo-200">
        <div class="flex items-start">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-indigo-400 mt-0.5" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-indigo-800">Optional Information</h3>
            <p class="text-sm text-indigo-700 mt-1">
              This information helps us provide personalized spiritual guidance and appropriate rituals. 
              All fields are optional and can be filled later.
            </p>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="flex flex-col sm:flex-row gap-4 pt-6 border-t">
        <!-- <button
          type="button"
          @click="$emit('previous')"
          class="px-6 py-3 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 transition-colors duration-200 font-medium order-3 sm:order-1"
        >
          Previous
        </button>

        <button
          type="button"
          @click="$emit('skip')"
          class="px-6 py-3 border border-gray-300 text-gray-700 rounded-xl hover:bg-gray-50 transition-colors duration-200 font-medium order-2 sm:order-2"
        >
          Skip for Now
        </button> -->
        
        <!-- <div class="flex gap-4 order-1 sm:order-3 sm:ml-auto">
          <button
            type="button"
            @click="saveDraft"
            :disabled="!hasChanges"
            class="px-6 py-3 border border-indigo-300 text-indigo-700 rounded-xl hover:bg-indigo-50 transition-colors duration-200 font-medium disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Save Draft
          </button>
          
          <button
            type="submit"
            class="px-6 py-3 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 transition-colors duration-200 font-medium"
          >
            {{ hasAnyData ? 'Save & Continue' : 'Skip & Continue' }}
          </button>
        </div> -->
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { useDevoteeStore } from '@/stores/devotee'
import { storeToRefs } from 'pinia'

// Props with defaults to prevent undefined errors
const props = defineProps({
  data: {
    type: Object,
    default: () => ({
      gotra: '',
      nakshatra: '',
      rashi: '',
      lagna: '',
      vedaShaka: ''
    })
  }
})

// Emits
const emit = defineEmits(['skip', 'next', 'previous', 'save-draft', 'update-data'])

// Store
const devoteeStore = useDevoteeStore()

// Form Data
const formData = ref({
  gotra: '',
  nakshatra: '',
  rashi: '',
  lagna: '',
  vedaShaka: ''
})

// Validation
const errors = ref({})

// Data Lists
const popularGotras = [
  'Bharadwaja', 'Kashyapa', 'Vasishtha', 'Vishwamitra', 'Jamadagni', 'Atri', 'Gautama',
  'Agastya', 'Kaundinya', 'Kaushika', 'Bhargava', 'Shandilya', 'Angirasa', 'Harita',
  'Mudgala', 'Parasara', 'Vatsya', 'Garga', 'Shaktri', 'Pulastya'
]

const nakshatras = [
  'Ashwini', 'Bharani', 'Krittika', 'Rohini', 'Mrigashirsha', 'Ardra', 'Punarvasu',
  'Pushya', 'Ashlesha', 'Magha', 'Purva Phalguni', 'Uttara Phalguni', 'Hasta',
  'Chitra', 'Swati', 'Vishakha', 'Anuradha', 'Jyeshtha', 'Mula', 'Purva Ashadha',
  'Uttara Ashadha', 'Shravana', 'Dhanishta', 'Shatabhisha', 'Purva Bhadrapada',
  'Uttara Bhadrapada', 'Revati'
]

const rashis = [
  'Mesha (Aries)', 'Vrishabha (Taurus)', 'Mithuna (Gemini)', 'Karka (Cancer)',
  'Simha (Leo)', 'Kanya (Virgo)', 'Tula (Libra)', 'Vrishchika (Scorpio)',
  'Dhanu (Sagittarius)', 'Makara (Capricorn)', 'Kumbha (Aquarius)', 'Meena (Pisces)'
]

const vedaShakas = [
  'Rig Veda - Shankhayana', 'Rig Veda - Kaushitaki', 'Sama Veda - Kauthuma',
  'Sama Veda - Jaiminiya', 'Yajur Veda - Madhyandina', 'Yajur Veda - Kanva',
  'Yajur Veda - Taittiriya', 'Yajur Veda - Maitrayani', 'Atharva Veda - Shaunaka',
  'Atharva Veda - Paippalada'
]

// Computed
const hasAnyData = computed(() => {
  return formData.value.gotra || 
         formData.value.nakshatra || 
         formData.value.rashi || 
         formData.value.lagna || 
         formData.value.vedaShaka
})

const hasChanges = computed(() => {
  // Check if props.data exists first to avoid undefined errors
  if (!props.data) return false
  
  return formData.value.gotra !== (props.data.gotra || '') ||
         formData.value.nakshatra !== (props.data.nakshatra || '') ||
         formData.value.rashi !== (props.data.rashi || '') ||
         formData.value.lagna !== (props.data.lagna || '') ||
         formData.value.vedaShaka !== (props.data.vedaShaka || '')
})

// Methods
function validateForm() {
  errors.value = {}
  // No required fields for spiritual info - all optional
  return true
}

function handleSubmit() {
  if (validateForm()) {
    emit('update-data', formData.value)
    emit('next')
  }
}

function saveDraft() {
  emit('update-data', formData.value)
  emit('save-draft')
}

// Load existing data
onMounted(() => {
  // Initialize form with data from props or empty values if not available
  if (props.data) {
    formData.value = {
      gotra: props.data.gotra || '',
      nakshatra: props.data.nakshatra || '',
      rashi: props.data.rashi || '',
      lagna: props.data.lagna || '',
      vedaShaka: props.data.vedaShaka || ''
    }
  }
})

// Watch for validation on input change
watch(formData, () => {
  if (Object.keys(errors.value).length > 0) {
    validateForm()
  }
  
  // Update parent component with data changes
  emit('update-data', formData.value)
}, { deep: true })
</script>