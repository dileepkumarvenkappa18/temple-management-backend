<template>
  <div class="bg-white rounded-xl shadow-md p-6 mb-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Search Devotees</h3>
    
    <!-- Search Form -->
    <form @submit.prevent="handleSearch" class="space-y-4">
      <!-- Search Input -->
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <Search class="h-5 w-5 text-gray-400" />
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by name, email, or phone..."
          class="block w-full pl-10 pr-3 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          @input="debouncedSearch"
        />
      </div>

      <!-- Advanced Filters Toggle -->
      <div class="flex items-center justify-between">
        <button
          type="button"
          @click="showAdvancedFilters = !showAdvancedFilters"
          class="inline-flex items-center px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 border border-indigo-200 rounded-lg hover:bg-indigo-100 transition-all duration-200"
        >
          <Filter class="h-4 w-4 mr-2" />
          Advanced Filters
          <ChevronDown 
            :class="['h-4 w-4 ml-2 transition-transform duration-200', showAdvancedFilters ? 'rotate-180' : '']" 
          />
        </button>

        <!-- Clear Filters -->
        <button
          v-if="hasActiveFilters"
          type="button"
          @click="clearFilters"
          class="inline-flex items-center px-3 py-2 text-sm font-medium text-red-600 hover:text-red-700 transition-colors duration-200"
        >
          <X class="h-4 w-4 mr-1" />
          Clear All
        </button>
      </div>

      <!-- Advanced Filters Panel -->
      <transition
        enter-active-class="transition-all duration-300 ease-out"
        enter-from-class="opacity-0 max-h-0"
        enter-to-class="opacity-100 max-h-96"
        leave-active-class="transition-all duration-300 ease-in"
        leave-from-class="opacity-100 max-h-96"
        leave-to-class="opacity-0 max-h-0"
      >
        <div v-if="showAdvancedFilters" class="bg-gray-50 rounded-lg p-4 space-y-4 overflow-hidden">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            <!-- Gotra Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Gotra</label>
              <select
                v-model="filters.gotra"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              >
                <option value="">All Gotras</option>
                <option v-for="gotra in gotraOptions" :key="gotra" :value="gotra">
                  {{ gotra }}
                </option>
              </select>
            </div>

            <!-- Native Place Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Native Place</label>
              <select
                v-model="filters.nativePlace"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              >
                <option value="">All Places</option>
                <option v-for="place in nativePlaceOptions" :key="place" :value="place">
                  {{ place }}
                </option>
              </select>
            </div>

            <!-- Status Filter -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
              <select
                v-model="filters.status"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              >
                <option value="">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="pending">Pending</option>
              </select>
            </div>

            <!-- Registration Date Range -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Registered From</label>
              <input
                v-model="filters.dateFrom"
                type="date"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Registered To</label>
              <input
                v-model="filters.dateTo"
                type="date"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              />
            </div>

            <!-- Profile Completion -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Profile Completion</label>
              <select
                v-model="filters.profileCompletion"
                class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              >
                <option value="">Any Completion</option>
                <option value="incomplete">Incomplete (< 100%)</option>
                <option value="complete">Complete (100%)</option>
                <option value="partial">Partial (50-99%)</option>
                <option value="minimal">Minimal (< 50%)</option>
              </select>
            </div>
          </div>

          <!-- Search Button -->
          <div class="flex justify-end pt-2">
            <button
              type="submit"
              class="inline-flex items-center px-6 py-2 bg-indigo-600 text-white font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
            >
              <Search class="h-4 w-4 mr-2" />
              Apply Filters
            </button>
          </div>
        </div>
      </transition>
    </form>

    <!-- Search Results Summary -->
    <div v-if="searchResults.length > 0" class="mt-4 p-3 bg-indigo-50 rounded-lg">
      <p class="text-sm text-indigo-800">
        <span class="font-medium">{{ searchResults.length }}</span> 
        {{ searchResults.length === 1 ? 'devotee' : 'devotees' }} found
        <span v-if="searchQuery"> for "{{ searchQuery }}"</span>
      </p>
    </div>

    <!-- No Results -->
    <div v-else-if="hasSearched && searchResults.length === 0" class="mt-4 p-4 bg-yellow-50 rounded-lg">
      <div class="flex items-center">
        <AlertTriangle class="h-5 w-5 text-yellow-400 mr-2" />
        <p class="text-sm text-yellow-800">
          No devotees found matching your search criteria. Try adjusting your filters.
        </p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useDebounceFn } from '@vueuse/core'
import {
  Search,
  Filter,
  ChevronDown,
  X,
  AlertTriangle
} from 'lucide-vue-next'

// Props
const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  }
})

// Emits
const emit = defineEmits(['update:modelValue', 'search'])

// Reactive data
const searchQuery = ref('')
const showAdvancedFilters = ref(false)
const hasSearched = ref(false)

const filters = ref({
  gotra: '',
  nativePlace: '',
  status: '',
  dateFrom: '',
  dateTo: '',
  profileCompletion: ''
})

// Filter options (these would typically come from API)
const gotraOptions = ref([
  'Bharadwaja', 'Kashyapa', 'Atri', 'Vishwamitra', 'Jamadagni', 
  'Vasishta', 'Gautama', 'Agastya', 'Kaundinya', 'Sandilya'
])

const nativePlaceOptions = ref([
  'Bangalore', 'Mysore', 'Mangalore', 'Hubli', 'Belgaum', 
  'Davangere', 'Shimoga', 'Tumkur', 'Gulbarga', 'Bijapur'
])

const searchResults = ref([])

// Computed
const hasActiveFilters = computed(() => {
  return searchQuery.value || 
         filters.value.gotra || 
         filters.value.nativePlace || 
         filters.value.status ||
         filters.value.dateFrom ||
         filters.value.dateTo ||
         filters.value.profileCompletion
})

// Methods
const handleSearch = () => {
  hasSearched.value = true
  const searchCriteria = {
    query: searchQuery.value,
    ...filters.value
  }
  
  // Emit search event
  emit('search', searchCriteria)
  
  // Mock search results (replace with actual API call)
  mockSearch(searchCriteria)
}

const mockSearch = (criteria) => {
  // Simulate API call delay
  setTimeout(() => {
    const mockResults = [
      {
        id: 1,
        name: 'Ravi Kumar',
        email: 'ravi@example.com',
        gotra: 'Bharadwaja',
        nativePlace: 'Bangalore',
        status: 'active'
      },
      {
        id: 2,
        name: 'Priya Sharma',
        email: 'priya@example.com',
        gotra: 'Kashyapa',
        nativePlace: 'Mysore',
        status: 'active'
      }
    ]
    searchResults.value = mockResults
    emit('update:modelValue', mockResults)
  }, 300)
}

const clearFilters = () => {
  searchQuery.value = ''
  filters.value = {
    gotra: '',
    nativePlace: '',
    status: '',
    dateFrom: '',
    dateTo: '',
    profileCompletion: ''
  }
  searchResults.value = []
  hasSearched.value = false
  emit('update:modelValue', [])
}

// Debounced search for real-time filtering
const debouncedSearch = useDebounceFn(() => {
  if (searchQuery.value.length >= 2) {
    handleSearch()
  } else if (searchQuery.value.length === 0) {
    searchResults.value = []
    hasSearched.value = false
    emit('update:modelValue', [])
  }
}, 500)

// Lifecycle
onMounted(() => {
  // Load initial filter options if needed
})
</script>