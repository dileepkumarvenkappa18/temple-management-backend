<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Page Header -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-3">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-xl font-bold text-gray-900">Devotee Management</h1>
              <p class="mt-1 text-sm text-gray-600">Manage your temple's registered devotees</p>
            </div>
            <div class="flex items-center space-x-2">
              <button @click="refreshData" :disabled="loading" class="inline-flex items-center px-3 py-1.5 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 disabled:opacity-50">
                <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                Refresh
              </button>
              <button @click="exportData" class="inline-flex items-center px-3 py-1.5 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200">
                <svg class="w-4 h-4 mr-1.5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                Export
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
      <!-- Stats Overview -->
      <div class="mb-4">
        <DevoteeStats :stats="devoteeStore.devoteeStats" :loading="loading" />
      </div>

      <!-- Search and Filters -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-100 mb-4">
        <div class="p-3">
          <div class="flex items-center space-x-3">
            <div class="flex-1">
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                  </svg>
                </div>
                <input v-model="searchQuery" type="text" placeholder="Search devotees by name, email, or phone..." class="block w-full pl-9 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm" />
              </div>
            </div>
            <div class="flex items-center space-x-2">
              <select v-model="statusFilter" class="block w-32 px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option value="">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
              </select>
              <select v-model="itemsPerPage" @change="resetPagination" class="block w-20 px-3 py-2 border border-gray-300 rounded-lg focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm">
                <option :value="10">10</option>
                <option :value="25">25</option>
                <option :value="50">50</option>
                <option :value="100">100</option>
              </select>
            </div>
          </div>
        </div>
      </div>

      <!-- Devotee List -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200">
        <div class="p-3">
          <div class="flex items-center justify-between mb-3">
            <div>
              <h2 class="text-lg font-semibold text-gray-900">All Devotees</h2>
              <p class="text-sm text-gray-600">
                Showing {{ paginationInfo.start }}-{{ paginationInfo.end }} of {{ paginationInfo.total }} devotees
              </p>
            </div>
            <div class="flex items-center space-x-2">
              <!-- View Toggle -->
              <div class="flex bg-gray-100 rounded-lg p-0.5">
                <button @click="viewMode = 'table'" :class="['px-2 py-1 text-sm font-medium rounded-md transition-all duration-200', viewMode === 'table' ? 'bg-white text-indigo-600 shadow-sm' : 'text-gray-500 hover:text-gray-700']">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
                  </svg>
                </button>
              </div>
            </div>
          </div>

          <!-- Table View -->
          <div v-if="viewMode === 'table'" class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200">
              <thead class="bg-gray-50">
                <tr>
                  <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">User ID</th>
                  <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Full Name</th>
                  <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Email</th>
                  <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Phone</th>
                  <th scope="col" class="px-3 py-2 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                  <th scope="col" class="px-3 py-2 text-center text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                </tr>
              </thead>
              <tbody class="bg-white divide-y divide-gray-200">
                <tr v-for="devotee in paginatedDevotees" :key="getDevoteeId(devotee)" class="hover:bg-gray-50">
                  <td class="px-3 py-2 whitespace-nowrap text-sm font-medium text-gray-900">#{{ getDevoteeId(devotee) }}</td>
                  <td class="px-3 py-2 whitespace-nowrap">
                    <div class="flex items-center">
                      <div class="flex-shrink-0 h-7 w-7">
                        <div class="h-7 w-7 rounded-full bg-indigo-100 flex items-center justify-center">
                          <span class="text-xs font-medium text-indigo-600">{{ getDevoteeName(devotee).charAt(0).toUpperCase() }}</span>
                        </div>
                      </div>
                      <div class="ml-2">
                        <div class="text-sm font-medium text-gray-900">{{ getDevoteeName(devotee) }}</div>
                      </div>
                    </div>
                  </td>
                  <td class="px-3 py-2 whitespace-nowrap text-sm text-gray-500">{{ devotee.email || 'Not provided' }}</td>
                  <td class="px-3 py-2 whitespace-nowrap text-sm text-gray-500">{{ devotee.phone || 'Not provided' }}</td>
                  <td class="px-3 py-2 whitespace-nowrap">
                    <span :class="getStatusClass(devotee.status)" class="inline-flex px-2 py-0.5 text-xs font-semibold rounded-full">{{ devotee.status || 'active' }}</span>
                  </td>
                  <td class="px-3 py-2 whitespace-nowrap text-center text-sm font-medium">
                    <div class="flex items-center justify-center">
                      <button
                        @click="handleToggleStatus(devotee)"
                        :disabled="statusUpdateLoading === getDevoteeId(devotee)"
                        :class="getToggleButtonClass(devotee.status)"
                        class="px-2 py-0.5 rounded text-xs transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
                        :title="getStatusToggleText(devotee.status)">
                        <span v-if="statusUpdateLoading === getDevoteeId(devotee)" class="flex items-center">
                          <svg class="animate-spin -ml-1 mr-1 h-3 w-3 text-current" fill="none" viewBox="0 0 24 24">
                            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                          </svg>
                          Updating...
                        </span>
                        <span v-else>{{ getStatusToggleText(devotee.status) }}</span>
                      </button>
                    </div>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>

          <!-- Grid View -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-4">
            <DevoteeCard
              v-for="devotee in paginatedDevotees"
              :key="getDevoteeId(devotee)"
              :devotee="devotee"
              :loading="statusUpdateLoading === getDevoteeId(devotee)"
              @toggle-status="handleToggleStatus"
            />
          </div>

          <!-- Pagination Controls -->
          <div v-if="totalPages > 1" class="flex items-center justify-between border-t border-gray-200 bg-white px-3 py-2 mt-3">
            <div class="flex flex-1 justify-between sm:hidden">
              <button
                @click="goToPage(currentPage - 1)"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Previous
              </button>
              <button
                @click="goToPage(currentPage + 1)"
                :disabled="currentPage === totalPages"
                class="relative ml-3 inline-flex items-center rounded-md border border-gray-300 bg-white px-3 py-1.5 text-sm font-medium text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                Next
              </button>
            </div>
            <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
              <div>
                <p class="text-sm text-gray-700">
                  Showing
                  <span class="font-medium">{{ paginationInfo.start }}</span>
                  to
                  <span class="font-medium">{{ paginationInfo.end }}</span>
                  of
                  <span class="font-medium">{{ paginationInfo.total }}</span>
                  results
                </p>
              </div>
              <div>
                <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                  <!-- Previous button -->
                  <button
                    @click="goToPage(currentPage - 1)"
                    :disabled="currentPage === 1"
                    class="relative inline-flex items-center rounded-l-md px-2 py-1.5 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">Previous</span>
                    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                      <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                    </svg>
                  </button>

                  <!-- Page numbers -->
                  <template v-for="page in visiblePages" :key="page">
                    <button
                      v-if="page !== '...'"
                      @click="goToPage(page)"
                      :class="[
                        'relative inline-flex items-center px-3 py-1.5 text-sm font-semibold focus:z-20 focus:outline-offset-0',
                        page === currentPage
                          ? 'z-10 bg-indigo-600 text-white focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
                          : 'text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50'
                      ]"
                    >
                      {{ page }}
                    </button>
                    <span
                      v-else
                      class="relative inline-flex items-center px-3 py-1.5 text-sm font-semibold text-gray-700 ring-1 ring-inset ring-gray-300"
                    >
                      ...
                    </span>
                  </template>

                  <!-- Next button -->
                  <button
                    @click="goToPage(currentPage + 1)"
                    :disabled="currentPage === totalPages"
                    class="relative inline-flex items-center rounded-r-md px-2 py-1.5 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    <span class="sr-only">Next</span>
                    <svg class="h-4 w-4" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                      <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
                    </svg>
                  </button>
                </nav>
              </div>
            </div>
          </div>

          <!-- Empty State -->
          <div v-if="paginatedDevotees.length === 0 && !loading" class="text-center py-6">
            <svg class="mx-auto h-10 w-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No devotees found</h3>
            <p class="mt-1 text-sm text-gray-500">
              {{ searchQuery || statusFilter ? 'Try adjusting your search criteria or status filter.' : 'No devotees have joined this temple yet.' }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading Overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-25 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-4 flex items-center space-x-3">
        <BaseLoader size="sm" />
        <span class="text-gray-700">Loading devotees...</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useToast } from '@/composables/useToast'
import { useDevoteeStore } from '@/stores/devotee'
import { storeToRefs } from 'pinia'

// Components
import DevoteeCard from '@/components/devotee/DevoteeCard.vue'
import DevoteeStats from '@/components/devotee/DevoteeStats.vue'
import BaseLoader from '@/components/common/BaseLoader.vue'

export default {
  name: 'DevoteeManagement',
  components: {
    DevoteeCard,
    DevoteeStats,
    BaseLoader
  },
  setup() {
    const route = useRoute()
    const router = useRouter()
    const toast = useToast()
    const devoteeStore = useDevoteeStore()
    const { devotees } = storeToRefs(devoteeStore)

    // Reactive data
    const loading = ref(false)
    const statusUpdateLoading = ref(null)
    const viewMode = ref('table')
    const searchQuery = ref('')
    const statusFilter = ref('')
    const activeFilters = ref({})

    // Pagination state
    const currentPage = ref(1)
    const itemsPerPage = ref(25)

    // Get entity ID safely
    const getEntityId = () => {
      const entityId = route.params.id || route.params.entityId
      if (!entityId) {
        console.error('❌ No entity ID found in route params')
        toast.error('Invalid entity ID')
        return null
      }
      return entityId
    }

    // Computed
    const filteredDevotees = computed(() => {
      let devotees = devoteeStore.devotees || []

      console.log('🔍 Filtering devotees:', {
        totalDevotees: devotees.length,
        searchQuery: searchQuery.value,
        statusFilter: statusFilter.value
      })

      // Apply search
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        devotees = devotees.filter(devotee =>
          getDevoteeName(devotee)?.toLowerCase().includes(query) ||
          devotee.email?.toLowerCase().includes(query) ||
          devotee.phone?.toLowerCase().includes(query) ||
          devotee.gotra?.toLowerCase().includes(query) ||
          devotee.nativePlace?.toLowerCase().includes(query)
        )
      }

      // Apply status filter
      if (statusFilter.value) {
        devotees = devotees.filter(devotee =>
          (devotee.status || 'active') === statusFilter.value
        )
      }

      // Apply other filters
      if (activeFilters.value.gotra) {
        devotees = devotees.filter(devotee => devotee.gotra === activeFilters.value.gotra)
      }

      if (activeFilters.value.profileCompletion) {
        const completionRange = activeFilters.value.profileCompletion
        devotees = devotees.filter(devotee => {
          const completion = devotee.profileCompletion || 0
          switch (completionRange) {
            case 'incomplete': return completion < 50
            case 'partial': return completion >= 50 && completion < 100
            case 'complete': return completion === 100
            default: return true
          }
        })
      }

      console.log('✅ Filtered devotees result:', devotees.length)
      return devotees
    })

    // Pagination computed properties
    const totalPages = computed(() => {
      return Math.ceil(filteredDevotees.value.length / itemsPerPage.value)
    })

    const paginatedDevotees = computed(() => {
      const start = (currentPage.value - 1) * itemsPerPage.value
      const end = start + itemsPerPage.value
      return filteredDevotees.value.slice(start, end)
    })

    const paginationInfo = computed(() => {
      const start = filteredDevotees.value.length === 0 ? 0 : (currentPage.value - 1) * itemsPerPage.value + 1
      const end = Math.min(currentPage.value * itemsPerPage.value, filteredDevotees.value.length)
      const total = filteredDevotees.value.length

      return { start, end, total }
    })

    const visiblePages = computed(() => {
      const total = totalPages.value
      const current = currentPage.value
      const delta = 2 // Number of pages to show on each side of current page

      if (total <= 7) {
        // Show all pages if total is small
        return Array.from({ length: total }, (_, i) => i + 1)
      }

      const pages = []

      // Always show first page
      pages.push(1)

      // Add ellipsis if needed
      if (current - delta > 2) {
        pages.push('...')
      }

      // Add pages around current
      const start = Math.max(2, current - delta)
      const end = Math.min(total - 1, current + delta)

      for (let i = start; i <= end; i++) {
        pages.push(i)
      }

      // Add ellipsis if needed
      if (current + delta < total - 1) {
        pages.push('...')
      }

      // Always show last page if more than 1 page
      if (total > 1) {
        pages.push(total)
      }

      return pages
    })

    // Watch for filter changes and reset pagination
    watch([searchQuery, statusFilter], () => {
      currentPage.value = 1
    })

    // Helper methods
    const getDevoteeId = (devotee) => {
      return devotee.user_id || devotee.id || 'unknown'
    }

    const getDevoteeName = (devotee) => {
      return devotee.full_name || devotee.fullName || 'Unknown'
    }

    const getStatusClass = (status) => {
      switch (status || 'active') {
        case 'active':
          return 'bg-green-100 text-green-800'
        case 'inactive':
          return 'bg-red-100 text-red-800'
        case 'pending':
          return 'bg-yellow-100 text-yellow-800'
        default:
          return 'bg-gray-100 text-gray-800'
      }
    }

    const getToggleButtonClass = (status) => {
      return (status || 'active') === 'active'
        ? 'text-red-600 hover:text-red-900 hover:bg-red-50 border border-red-300'
        : 'text-green-600 hover:text-green-900 hover:bg-green-50 border border-green-300'
    }

    const getStatusToggleText = (status) => {
      return (status || 'active') === 'active' ? 'Deactivate' : 'Activate'
    }

    // Pagination methods
    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        currentPage.value = page
      }
    }

    const resetPagination = () => {
      currentPage.value = 1
    }

    // Data loading methods
    const loadDevotees = async () => {
      const entityId = getEntityId()
      if (!entityId) return

      loading.value = true
      try {
        console.log('🔄 Loading devotees for entity:', entityId)
        await devoteeStore.fetchDevoteesByEntity(entityId)
        console.log('✅ Devotees loaded:', devoteeStore.devotees.length)
      } catch (error) {
        console.error('❌ Error loading devotees:', error)
        toast.error('Failed to load devotees')
      } finally {
        loading.value = false
      }
    }

    const loadDevoteeStats = async () => {
      const entityId = getEntityId()
      if (!entityId) return

      try {
        console.log('🔄 Loading devotee stats for entity:', entityId)
        await devoteeStore.fetchDevoteeStats(entityId)
        console.log('✅ Devotee stats loaded:', devoteeStore.devoteeStats)
      } catch (error) {
        console.error('❌ Error loading devotee stats:', error)
        toast.error('Failed to load devotee statistics')
      }
    }

    const refreshData = async () => {
      await Promise.all([
        loadDevotees(),
        loadDevoteeStats()
      ])
    }

    // Status toggle handler - FIXED
    const handleToggleStatus = async (devotee) => {
      const entityId = getEntityId()
      if (!entityId) return

      const devoteeId = getDevoteeId(devotee)
      const currentStatus = devotee.status || 'active'
      const newStatus = currentStatus === 'active' ? 'inactive' : 'active'

      console.log('🔄 Toggling devotee status:', {
        devoteeId,
        currentStatus,
        newStatus,
        entityId
      })

      // Validate inputs
      if (!devoteeId || devoteeId === 'unknown') {
        toast.error('Invalid devotee ID')
        return
      }

      statusUpdateLoading.value = devoteeId

      try {
        // Use the correct store method
        await devoteeStore.updateDevoteeStatus(entityId, devoteeId, newStatus)

        // Refresh stats after status change
        await loadDevoteeStats()

        console.log('✅ Status updated successfully')
        toast.success(`Devotee ${newStatus === 'active' ? 'activated' : 'deactivated'} successfully`)

      } catch (error) {
        console.error('❌ Error updating status:', error)
        
        // More detailed error handling
        const errorMessage = error.response?.data?.error || 
                            error.response?.data?.message || 
                            error.message || 
                            'Failed to update devotee status'
        
        toast.error(errorMessage)
      } finally {
        statusUpdateLoading.value = null
      }
    }

    // Export functionality
    const exportData = () => {
      try {
        // Export all filtered data, not just paginated data
        const csvContent = [
          ['User ID', 'Full Name', 'Email', 'Phone', 'Status', 'Gotra', 'Created Date'],
          ...filteredDevotees.value.map(d => [
            getDevoteeId(d),
            getDevoteeName(d),
            d.email || '',
            d.phone || '',
            d.status || 'active',
            d.gotra || '',
            d.created_at ? new Date(d.created_at).toLocaleDateString() : ''
          ])
        ].map(row => row.join(',')).join('\n')

        // Create and download file
        const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
        const url = window.URL.createObjectURL(blob)
        const a = document.createElement('a')
        a.href = url
        a.download = `devotees_${new Date().toISOString().split('T')[0]}.csv`
        document.body.appendChild(a)
        a.click()
        document.body.removeChild(a)
        window.URL.revokeObjectURL(url)

        toast.success('Devotee data exported successfully')
      } catch (error) {
        console.error('❌ Export error:', error)
        toast.error('Failed to export data')
      }
    }

    // Error handling for missing route params
    const validateRouteParams = () => {
      const entityId = getEntityId()
      if (!entityId) {
        console.error('❌ No entity ID in route, redirecting...')
        router.push('/dashboard')
        return false
      }
      return true
    }

    // Lifecycle
    onMounted(async () => {
      console.log('🚀 DevoteeManagement mounted')
      console.log('Route params:', route.params)
      
      if (!validateRouteParams()) {
        return
      }

      await refreshData()
    })

    return {
      // Store
      devoteeStore,

      // Reactive data
      loading,
      statusUpdateLoading,
      viewMode,
      searchQuery,
      statusFilter,

      // Pagination
      currentPage,
      itemsPerPage,
      totalPages,
      paginatedDevotees,
      paginationInfo,
      visiblePages,

      // Computed
      filteredDevotees,

      // Helper methods
      getDevoteeId,
      getDevoteeName,
      getStatusClass,
      getToggleButtonClass,
      getStatusToggleText,

      // Action methods
      handleToggleStatus,
      exportData,
      refreshData,
      goToPage,
      resetPagination
    }
  }
}
</script>

<style scoped>
/* Loading spinner animation */
@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.animate-spin {
  animation: spin 1s linear infinite;
}

/* Improved button states */
.disabled\:opacity-50:disabled {
  opacity: 0.5;
}

.disabled\:cursor-not-allowed:disabled {
  cursor: not-allowed;
}

/* Hover effects */
.hover\:bg-gray-50:hover {
  background-color: #f9fafb;
}

.hover\:bg-red-50:hover {
  background-color: #fef2f2;
}

.hover\:bg-green-50:hover {
  background-color: #f0fdf4;
}

.hover\:bg-indigo-50:hover {
  background-color: #eef2ff;
}

/* Focus states for accessibility */
.focus\:ring-2:focus {
  --tw-ring-width: 2px;
}

.focus\:ring-indigo-500:focus {
  --tw-ring-color: #6366f1;
}

.focus\:ring-offset-2:focus {
  --tw-ring-offset-width: 2px;
}

/* Transition effects */
.transition-all {
  transition-property: all;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

.transition-colors {
  transition-property: color, background-color, border-color, text-decoration-color, fill, stroke;
  transition-timing-function: cubic-bezier(0.4, 0, 0.2, 1);
  transition-duration: 150ms;
}

.duration-200 {
  transition-duration: 200ms;
}
</style>