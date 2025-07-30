<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Page Header -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <!-- <div class="py-6">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Devotee Management</h1>
              <p class="mt-1 text-sm text-gray-600">Manage your temple's registered devotees</p>
            </div>
            <div class="flex items-center space-x-3">
              <button 
                @click="exportData" 
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                Export
              </button>
              <button 
                @click="showAddDevoteeModal = true"
                class="inline-flex items-center px-4 py-2 bg-indigo-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Add Devotee
              </button>
            </div>
          </div>
        </div> -->
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 pt-2">
      <!-- Stats Overview -->
      <div class="mb-0">
        <DevoteeStats />
      </div>
      
      <!-- Search and Filters -->
      <!-- <div class="bg-white rounded-xl shadow-sm border border-gray-200 mb-6">
        <div class="p-6">
          <DevoteeSearch 
            @search="handleSearch"
            @filter="handleFilter"
            :loading="loading"
          />
        </div>
      </div> -->

      <!-- Devotee List -->
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 -mt-64">
        <div class="p-6">
          <div class="flex items-center justify-between mb-6">
            <div>
              <h2 class="text-lg font-semibold text-gray-900">All Devotees</h2>
              <p class="text-sm text-gray-600">{{ filteredDevotees.length }} devotees found</p>
            </div>
            <div class="flex items-center space-x-3">
              <!-- View Toggle -->
              <div class="flex bg-gray-100 rounded-lg p-1">
                <button
                  @click="viewMode = 'table'"
                  :class="[
                    'px-3 py-1 text-sm font-medium rounded-md transition-all duration-200',
                    viewMode === 'table' 
                      ? 'bg-white text-indigo-600 shadow-sm' 
                      : 'text-gray-500 hover:text-gray-700'
                  ]"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16"></path>
                  </svg>
                </button>
                <!-- <button
                  @click="viewMode = 'grid'"
                  :class="[
                    'px-3 py-1 text-sm font-medium rounded-md transition-all duration-200',
                    viewMode === 'grid' 
                      ? 'bg-white text-indigo-600 shadow-sm' 
                      : 'text-gray-500 hover:text-gray-700'
                  ]"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z"></path>
                  </svg>
                </button> -->
              </div>
            </div>
          </div>

          <!-- Table View -->
          <div v-if="viewMode === 'table'">
            <DevoteeList 
              :devotees="filteredDevotees"
              :loading="loading"
              @view="viewDevotee"
              @edit="editDevotee"
              @deactivate="deactivateDevotee"
            />
          </div>

          <!-- Grid View -->
          <div v-else class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 xl:grid-cols-4 gap-6">
            <DevoteeCard
              v-for="devotee in filteredDevotees"
              :key="devotee.id"
              :devotee="devotee"
              @view="viewDevotee"
              @edit="editDevotee"
              @deactivate="deactivateDevotee"
            />
          </div>

          <!-- Empty State -->
          <div v-if="filteredDevotees.length === 0 && !loading" class="text-center py-8">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No devotees found</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by adding your first devotee.</p>
            <div class="mt-6">
              <!-- <button 
                @click="showAddDevoteeModal = true"
                class="inline-flex items-center px-4 py-2 bg-indigo-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Add First Devotee
              </button> -->
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Add Devotee Modal -->
    <BaseModal 
      v-if="showAddDevoteeModal" 
      @close="showAddDevoteeModal = false"
      title="Add New Devotee"
      size="lg"
    >
      <div class="space-y-6">
        <div class="text-sm text-gray-600">
          Send an invitation to a devotee to join your temple. They will receive an email with instructions to complete their profile.
        </div>
        
        <form @submit.prevent="addDevotee" class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <BaseInput
              v-model="newDevotee.fullName"
              label="Full Name"
              placeholder="Enter devotee's full name"
              required
            />
            <BaseInput
              v-model="newDevotee.email"
              label="Email"
              type="email"
              placeholder="Enter email address"
              required
            />
          </div>
          
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <BaseInput
              v-model="newDevotee.phone"
              label="Phone Number"
              placeholder="Enter phone number"
            />
            <BaseSelect
              v-model="newDevotee.gotra"
              label="Gotra (Optional)"
              placeholder="Select gotra"
              :options="gotraOptions"
            />
          </div>

          <div class="flex items-center">
            <input
              id="send-invitation"
              v-model="newDevotee.sendInvitation"
              type="checkbox"
              class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
            >
            <label for="send-invitation" class="ml-2 block text-sm text-gray-900">
              Send invitation email to devotee
            </label>
          </div>

          <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
            <button
              type="button"
              @click="showAddDevoteeModal = false"
              class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
            >
              Cancel
            </button>
            <button
              type="submit"
              :disabled="addingDevotee"
              class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 border border-transparent rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            >
              <span v-if="addingDevotee" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Adding...
              </span>
              <span v-else>Add Devotee</span>
            </button>
          </div>
        </form>
      </div>
    </BaseModal>

    <!-- Loading Overlay -->
    <div v-if="loading" class="fixed inset-0 bg-black bg-opacity-25 flex items-center justify-center z-50">
      <div class="bg-white rounded-lg p-6 flex items-center space-x-3">
        <BaseLoader size="sm" />
        <span class="text-gray-700">Loading devotees...</span>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useToast } from '@/composables/useToast'
import { useDevoteeStore } from '@/stores/devotee'

// Components
import DevoteeCard from '@/components/devotee/DevoteeCard.vue'
import DevoteeList from '@/components/devotee/DevoteeList.vue'
import DevoteeSearch from '@/components/devotee/DevoteeSearch.vue'
import DevoteeStats from '@/components/devotee/DevoteeStats.vue'
import BaseModal from '@/components/common/BaseModal.vue'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import BaseLoader from '@/components/common/BaseLoader.vue'

export default {
  name: 'DevoteeManagement',
  components: {
    DevoteeCard,
    DevoteeList,
    DevoteeSearch,
    DevoteeStats,
    BaseModal,
    BaseInput,
    BaseSelect,
    BaseLoader
  },
  setup() {
    const route = useRoute()
    const { showToast } = useToast()
    const devoteeStore = useDevoteeStore()

    // Reactive data
    const loading = ref(false)
    const viewMode = ref('table')
    const showAddDevoteeModal = ref(false)
    const addingDevotee = ref(false)
    const searchQuery = ref('')
    const activeFilters = ref({})

    const newDevotee = ref({
      fullName: '',
      email: '',
      phone: '',
      gotra: '',
      sendInvitation: true
    })

    const gotraOptions = ref([
      { value: 'bharadwaj', label: 'Bharadwaj' },
      { value: 'kashyap', label: 'Kashyap' },
      { value: 'gautam', label: 'Gautam' },
      { value: 'atri', label: 'Atri' },
      { value: 'jamadagni', label: 'Jamadagni' },
      { value: 'vasishtha', label: 'Vasishtha' },
      { value: 'vishwamitra', label: 'Vishwamitra' }
    ])

    // Computed
    const filteredDevotees = computed(() => {
      let devotees = devoteeStore.devotees || []

      // Apply search
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        devotees = devotees.filter(devotee => 
          devotee.fullName?.toLowerCase().includes(query) ||
          devotee.email?.toLowerCase().includes(query) ||
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

      return devotees
    })

    // Methods
    const loadDevotees = async () => {
      loading.value = true
      try {
        await devoteeStore.fetchDevotees(route.params.id)
      } catch (error) {
        showToast('Failed to load devotees', 'error')
      } finally {
        loading.value = false
      }
    }

    const handleSearch = (query) => {
      searchQuery.value = query
    }

    const handleFilter = (filters) => {
      activeFilters.value = filters
    }

    const viewDevotee = (devotee) => {
      // Navigate to devotee profile view
      showToast(`Viewing ${devotee.fullName}'s profile`, 'info')
    }

    const editDevotee = (devotee) => {
      // Open edit modal or navigate to edit page
      showToast(`Editing ${devotee.fullName}'s profile`, 'info')
    }

    const deactivateDevotee = async (devotee) => {
      try {
        await devoteeStore.deactivateDevotee(devotee.id)
        showToast(`${devotee.fullName} has been deactivated`, 'success')
      } catch (error) {
        showToast('Failed to deactivate devotee', 'error')
      }
    }

    const addDevotee = async () => {
      addingDevotee.value = true
      try {
        await devoteeStore.addDevotee({
          ...newDevotee.value,
          templeId: route.params.id
        })
        
        showAddDevoteeModal.value = false
        newDevotee.value = {
          fullName: '',
          email: '',
          phone: '',
          gotra: '',
          sendInvitation: true
        }
        
        showToast('Devotee added successfully', 'success')
        if (newDevotee.value.sendInvitation) {
          showToast('Invitation email sent', 'info')
        }
      } catch (error) {
        showToast('Failed to add devotee', 'error')
      } finally {
        addingDevotee.value = false
      }
    }

    const exportData = () => {
      // Export functionality
      showToast('Exporting devotee data...', 'info')
    }

    // Lifecycle
    onMounted(() => {
      loadDevotees()
    })

    return {
      loading,
      viewMode,
      showAddDevoteeModal,
      addingDevotee,
      newDevotee,
      gotraOptions,
      filteredDevotees,
      handleSearch,
      handleFilter,
      viewDevotee,
      editDevotee,
      deactivateDevotee,
      addDevotee,
      exportData
    }
  }
}
</script>