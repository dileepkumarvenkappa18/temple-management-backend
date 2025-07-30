<template>
  <div class="w-full">
    <!-- Search Input -->
    <div class="relative mb-6">
      <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
        <SearchIcon class="h-5 w-5 text-gray-400" />
      </div>
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search temples by name, city, or state..."
        class="block w-full pl-10 pr-4 py-3 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-gray-900 placeholder-gray-500 transition-all duration-200"
        @input="handleSearch"
      />
    </div>

    <!-- Filters -->
    <div class="flex flex-wrap gap-4 mb-6">
      <!-- State Filter -->
      <div class="min-w-0 flex-1 sm:flex-none sm:min-w-[180px]">
        <select
          v-model="filters.state"
          class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-gray-900 text-sm transition-all duration-200"
          @change="handleFilterChange"
        >
          <option value="">All States</option>
          <option v-for="state in states" :key="state" :value="state">{{ state }}</option>
        </select>
      </div>

      <!-- District Filter -->
      <div class="min-w-0 flex-1 sm:flex-none sm:min-w-[180px]">
        <select
          v-model="filters.district"
          class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-gray-900 text-sm transition-all duration-200"
          @change="handleFilterChange"
          :disabled="!filters.state"
        >
          <option value="">All Districts</option>
          <option v-for="district in availableDistricts" :key="district" :value="district">{{ district }}</option>
        </select>
      </div>

      <!-- Category Filter -->
      <div class="min-w-0 flex-1 sm:flex-none sm:min-w-[180px]">
        <select
          v-model="filters.category"
          class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 bg-white text-gray-900 text-sm transition-all duration-200"
          @change="handleFilterChange"
        >
          <option value="">All Categories</option>
          <option v-for="category in categories" :key="category" :value="category">{{ category }}</option>
        </select>
      </div>

      <!-- Clear Filters Button -->
      <button
        v-if="hasActiveFilters"
        @click="clearFilters"
        class="px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 border border-indigo-200 rounded-lg hover:bg-indigo-100 transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
      >
        Clear Filters
      </button>
    </div>

    <!-- Results Count -->
    <div class="mb-4 text-sm text-gray-600">
      <span v-if="loading">Searching...</span>
      <span v-else-if="temples.length > 0">
        {{ temples.length }} temple{{ temples.length !== 1 ? 's' : '' }} found
      </span>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch, onMounted } from 'vue'
import { SearchIcon } from 'lucide-vue-next'

// Props
const props = defineProps({
  temples: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Emits
const emit = defineEmits(['search', 'filter-change'])

// Reactive data
const searchQuery = ref('')
const filters = ref({
  state: '',
  district: '',
  category: ''
})

// Mock data for filters
const states = ref([
  'Andhra Pradesh', 'Arunachal Pradesh', 'Assam', 'Bihar', 'Chhattisgarh',
  'Goa', 'Gujarat', 'Haryana', 'Himachal Pradesh', 'Jharkhand',
  'Karnataka', 'Kerala', 'Madhya Pradesh', 'Maharashtra', 'Manipur',
  'Meghalaya', 'Mizoram', 'Nagaland', 'Odisha', 'Punjab',
  'Rajasthan', 'Sikkim', 'Tamil Nadu', 'Telangana', 'Tripura',
  'Uttar Pradesh', 'Uttarakhand', 'West Bengal'
])

const categories = ref([
  'Vishnu Temple', 'Shiva Temple', 'Devi Temple', 'Ganesha Temple',
  'Hanuman Temple', 'Murugan Temple', 'Ayyappa Temple', 'Jain Temple'
])

// Mock districts data (simplified)
const districtsByState = ref({
  'Karnataka': ['Bangalore Urban', 'Bangalore Rural', 'Mysore', 'Mangalore', 'Hubli-Dharwad', 'Bellary', 'Tumkur'],
  'Tamil Nadu': ['Chennai', 'Coimbatore', 'Madurai', 'Tiruchirappalli', 'Salem', 'Tirunelveli', 'Erode'],
  'Maharashtra': ['Mumbai', 'Pune', 'Nagpur', 'Nashik', 'Aurangabad', 'Solapur', 'Kolhapur'],
  'Gujarat': ['Ahmedabad', 'Surat', 'Vadodara', 'Rajkot', 'Bhavnagar', 'Jamnagar', 'Junagadh'],
  'Rajasthan': ['Jaipur', 'Jodhpur', 'Kota', 'Bikaner', 'Ajmer', 'Udaipur', 'Bharatpur'],
  // Add more states and districts as needed
})

// Computed properties
const availableDistricts = computed(() => {
  if (!filters.value.state) return []
  return districtsByState.value[filters.value.state] || []
})

const hasActiveFilters = computed(() => {
  return searchQuery.value || filters.value.state || filters.value.district || filters.value.category
})

// Methods
const handleSearch = () => {
  emitSearchData()
}

const handleFilterChange = () => {
  // Reset district if state changes
  if (!filters.value.state) {
    filters.value.district = ''
  }
  emitSearchData()
}

const clearFilters = () => {
  searchQuery.value = ''
  filters.value = {
    state: '',
    district: '',
    category: ''
  }
  emitSearchData()
}

const emitSearchData = () => {
  const searchData = {
    query: searchQuery.value.trim(),
    filters: { ...filters.value }
  }
  
  emit('search', searchData)
  emit('filter-change', searchData)
}

// Watchers
watch(() => filters.value.state, (newState) => {
  if (!newState) {
    filters.value.district = ''
  }
})

// Lifecycle
onMounted(() => {
  // Initial search to load all temples
  emitSearchData()
})
</script>