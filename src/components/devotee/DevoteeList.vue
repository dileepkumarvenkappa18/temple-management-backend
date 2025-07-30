<!-- src/components/devotee/DevoteeList.vue -->
<template>
  <div class="bg-white rounded-xl shadow-lg border border-gray-100">
    <!-- Header with Title and Actions -->
    <div class="px-6 py-4 border-b border-gray-100 bg-gradient-to-r from-indigo-50 to-white">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-semibold text-gray-900">Devotee Management</h2>
          <p class="text-sm text-gray-600 mt-1">Manage and track devotee information</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <!-- <button class="inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors duration-200 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
            </svg>
            Add Devotee
          </button> -->
          <!-- <button class="inline-flex items-center px-4 py-2 text-sm font-medium text-indigo-600 bg-white hover:bg-indigo-50 border border-indigo-200 rounded-lg transition-colors duration-200 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2">
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
            </svg>
            Export
          </button> -->
        </div>
      </div>
    </div>

    <!-- Search and Filters -->
    <!-- <div class="px-6 py-4 bg-gray-50 border-b border-gray-100">
      <DevoteeSearch 
        v-model:search="searchQuery"
        v-model:filters="activeFilters"
        @search="handleSearch"
        @filter="handleFilter"
        @reset="resetFilters"
      />
    </div> -->

    <!-- Table Container with Loading State -->
    <div class="relative overflow-hidden">
      <!-- Loading Overlay -->
      <div v-if="loading" class="absolute inset-0 bg-white bg-opacity-75 flex items-center justify-center z-10">
        <div class="flex items-center space-x-2">
          <div class="animate-spin rounded-full h-6 w-6 border-2 border-indigo-600 border-t-transparent"></div>
          <span class="text-sm text-gray-600">Loading devotees...</span>
        </div>
      </div>

      <!-- Table -->
      <div class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <thead class="bg-gray-100">
            <tr>
              <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                <div class="flex items-center space-x-1">
                  <span>Devotee</span>
                  <button @click="sortBy('name')" class="text-gray-400 hover:text-gray-600">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
                    </svg>
                  </button>
                </div>
              </th>
              <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Contact Info
              </th>
              <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Spiritual Info
              </th>
              <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                <div class="flex items-center space-x-1">
                  <span>Profile Status</span>
                  <button @click="sortBy('profileCompletion')" class="text-gray-400 hover:text-gray-600">
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"/>
                    </svg>
                  </button>
                </div>
              </th>
              <th scope="col" class="px-6 py-4 text-left text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Status
              </th>
              <th scope="col" class="px-6 py-4 text-center text-xs font-semibold text-gray-700 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="devotee in paginatedDevotees" :key="devotee.id" class="hover:bg-gray-50 transition-colors duration-150">
              <!-- Devotee Info -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex items-center">
                  <div class="flex-shrink-0 h-12 w-12">
                    <img class="h-12 w-12 rounded-full border-2 border-gray-200 object-cover" 
                         :src="devotee.photo || '/placeholder-avatar.png'" 
                         :alt="devotee.name">
                  </div>
                  <div class="ml-4">
                    <div class="text-sm font-semibold text-gray-900">{{ devotee.name }}</div>
                    <div class="text-sm text-gray-600">ID: {{ devotee.devoteeId }}</div>
                    <div class="text-xs text-gray-500">{{ devotee.nativePlace }}</div>
                  </div>
                </div>
              </td>

              <!-- Contact Info -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ devotee.email }}</div>
                <div class="text-sm text-gray-600">{{ devotee.phone }}</div>
                <div class="text-xs text-gray-500">{{ devotee.address?.city }}</div>
              </td>

              <!-- Spiritual Info -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="text-sm text-gray-900">{{ devotee.gotra }}</div>
                <div class="text-sm text-gray-600">{{ devotee.nakshatra }}</div>
                <div class="text-xs text-gray-500">{{ devotee.rashi }}</div>
              </td>

              <!-- Profile Completion -->
              <td class="px-6 py-4 whitespace-nowrap">
                <div class="flex flex-col space-y-2">
                  <div class="flex items-center justify-between">
                    <span class="text-sm font-medium text-gray-700">{{ devotee.profileCompletion }}%</span>
                    <span class="text-xs" :class="getCompletionStatusClass(devotee.profileCompletion)">
                      {{ getCompletionStatus(devotee.profileCompletion) }}
                    </span>
                  </div>
                  <div class="w-full bg-gray-200 rounded-full h-2">
                    <div class="h-2 rounded-full transition-all duration-300" 
                         :class="getCompletionBarClass(devotee.profileCompletion)"
                         :style="{ width: `${devotee.profileCompletion}%` }">
                    </div>
                  </div>
                </div>
              </td>

              <!-- Status -->
              <td class="px-6 py-4 whitespace-nowrap">
                <span class="inline-flex items-center px-2.5 py-1 rounded-full text-xs font-medium"
                      :class="getStatusClass(devotee.status)">
                  <span class="w-1.5 h-1.5 rounded-full mr-1.5" :class="getStatusDotClass(devotee.status)"></span>
                  {{ devotee.status }}
                </span>
              </td>

              <!-- Actions -->
              <td class="px-6 py-4 whitespace-nowrap text-center">
                <div class="flex items-center justify-center space-x-2">
                  <button @click="viewDevotee(devotee)" 
                          class="text-indigo-600 hover:text-indigo-800 transition-colors duration-200"
                          title="View Details">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
                    </svg>
                  </button>
                  <button @click="editDevotee(devotee)" 
                          class="text-gray-600 hover:text-gray-800 transition-colors duration-200"
                          title="Edit">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"/>
                    </svg>
                  </button>
                  <button @click="toggleDevoteeStatus(devotee)" 
                          :class="devotee.status === 'Active' ? 'text-red-600 hover:text-red-800' : 'text-green-600 hover:text-green-800'"
                          class="transition-colors duration-200"
                          :title="devotee.status === 'Active' ? 'Deactivate' : 'Activate'">
                    <svg v-if="devotee.status === 'Active'" class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M18.364 18.364A9 9 0 005.636 5.636m12.728 12.728L5.636 5.636m12.728 12.728L18 12M6 12l12.728 6.364"/>
                    </svg>
                    <svg v-else class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Empty State -->
      <div v-if="!loading && filteredDevotees.length === 0" class="text-center py-12">
        <svg class="mx-auto h-16 w-16 text-gray-300" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
        </svg>
        <h3 class="mt-4 text-lg font-medium text-gray-900">No devotees found</h3>
        <p class="mt-2 text-gray-600">{{ searchQuery ? 'Try adjusting your search criteria' : 'Get started by adding your first devotee' }}</p>
        <button v-if="!searchQuery" class="mt-4 inline-flex items-center px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors duration-200">
          <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
          </svg>
          Add First Devotee
        </button>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="filteredDevotees.length > 0" class="px-6 py-4 bg-gray-50 border-t border-gray-100">
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-600">
          Showing {{ startIndex + 1 }} to {{ Math.min(endIndex, filteredDevotees.length) }} of {{ filteredDevotees.length }} devotees
        </div>
        <div class="flex items-center space-x-2">
          <button @click="previousPage" 
                  :disabled="currentPage === 1"
                  :class="currentPage === 1 ? 'text-gray-400 cursor-not-allowed' : 'text-indigo-600 hover:text-indigo-800'"
                  class="p-2 rounded-lg transition-colors duration-200">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
            </svg>
          </button>
          
          <div class="flex items-center space-x-1">
            <button v-for="page in visiblePages" :key="page"
                    @click="goToPage(page)"
                    :class="page === currentPage ? 'bg-indigo-600 text-white' : 'text-gray-600 hover:bg-gray-100'"
                    class="px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-200">
              {{ page }}
            </button>
          </div>
          
          <button @click="nextPage" 
                  :disabled="currentPage === totalPages"
                  :class="currentPage === totalPages ? 'text-gray-400 cursor-not-allowed' : 'text-indigo-600 hover:text-indigo-800'"
                  class="p-2 rounded-lg transition-colors duration-200">
            <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
            </svg>
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useDevoteeStore } from '@/stores/devotee'
import DevoteeSearch from './DevoteeSearch.vue'

const devoteeStore = useDevoteeStore()

// Reactive data
const loading = ref(false)
const searchQuery = ref('')
const activeFilters = ref({})
const currentPage = ref(1)
const itemsPerPage = ref(10)
const sortField = ref('name')
const sortDirection = ref('asc')

// Computed properties
const filteredDevotees = computed(() => {
  let devotees = devoteeStore.devotees

  // Apply search
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    devotees = devotees.filter(devotee => 
      devotee.name.toLowerCase().includes(query) ||
      devotee.email.toLowerCase().includes(query) ||
      devotee.gotra?.toLowerCase().includes(query) ||
      devotee.nativePlace?.toLowerCase().includes(query)
    )
  }

  // Apply filters
  if (activeFilters.value.status) {
    devotees = devotees.filter(devotee => devotee.status === activeFilters.value.status)
  }
  if (activeFilters.value.gotra) {
    devotees = devotees.filter(devotee => devotee.gotra === activeFilters.value.gotra)
  }
  if (activeFilters.value.completionRange) {
    const [min, max] = activeFilters.value.completionRange
    devotees = devotees.filter(devotee => 
      devotee.profileCompletion >= min && devotee.profileCompletion <= max
    )
  }

  // Apply sorting
  devotees.sort((a, b) => {
    let aVal = a[sortField.value]
    let bVal = b[sortField.value]
    
    if (typeof aVal === 'string') {
      aVal = aVal.toLowerCase()
      bVal = bVal.toLowerCase()
    }
    
    if (sortDirection.value === 'asc') {
      return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
    } else {
      return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
    }
  })

  return devotees
})

const totalPages = computed(() => Math.ceil(filteredDevotees.value.length / itemsPerPage.value))
const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage.value)
const endIndex = computed(() => startIndex.value + itemsPerPage.value)

const paginatedDevotees = computed(() => 
  filteredDevotees.value.slice(startIndex.value, endIndex.value)
)

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  return pages
})

// Methods
const handleSearch = (query) => {
  searchQuery.value = query
  currentPage.value = 1
}

const handleFilter = (filters) => {
  activeFilters.value = filters
  currentPage.value = 1
}

const resetFilters = () => {
  searchQuery.value = ''
  activeFilters.value = {}
  currentPage.value = 1
}

const sortBy = (field) => {
  if (sortField.value === field) {
    sortDirection.value = sortDirection.value === 'asc' ? 'desc' : 'asc'
  } else {
    sortField.value = field
    sortDirection.value = 'asc'
  }
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const goToPage = (page) => {
  currentPage.value = page
}

const getCompletionStatus = (completion) => {
  if (completion >= 90) return 'Complete'
  if (completion >= 70) return 'Almost Done'
  if (completion >= 40) return 'In Progress'
  return 'Getting Started'
}

const getCompletionStatusClass = (completion) => {
  if (completion >= 90) return 'text-green-600'
  if (completion >= 70) return 'text-blue-600'
  if (completion >= 40) return 'text-yellow-600'
  return 'text-gray-600'
}

const getCompletionBarClass = (completion) => {
  if (completion >= 90) return 'bg-green-500'
  if (completion >= 70) return 'bg-blue-500'
  if (completion >= 40) return 'bg-yellow-500'
  return 'bg-gray-400'
}

const getStatusClass = (status) => {
  const classes = {
    'Active': 'bg-green-100 text-green-800',
    'Inactive': 'bg-gray-100 text-gray-800',
    'Pending': 'bg-yellow-100 text-yellow-800',
    'Suspended': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getStatusDotClass = (status) => {
  const classes = {
    'Active': 'bg-green-400',
    'Inactive': 'bg-gray-400',
    'Pending': 'bg-yellow-400',
    'Suspended': 'bg-red-400'
  }
  return classes[status] || 'bg-gray-400'
}

const viewDevotee = (devotee) => {
  // Navigate to devotee detail view
  console.log('View devotee:', devotee)
}

const editDevotee = (devotee) => {
  // Navigate to edit form
  console.log('Edit devotee:', devotee)
}

const toggleDevoteeStatus = (devotee) => {
  const newStatus = devotee.status === 'Active' ? 'Inactive' : 'Active'
  devoteeStore.updateDevoteeStatus(devotee.id, newStatus)
}

// Lifecycle
onMounted(async () => {
  loading.value = true
  try {
    await devoteeStore.fetchDevotees()
  } finally {
    loading.value = false
  }
})
</script>