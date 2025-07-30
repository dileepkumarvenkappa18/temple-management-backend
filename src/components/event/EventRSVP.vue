<template>
  <div class="bg-white rounded-xl shadow-md p-6">
    <!-- Header -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between mb-6">
      <div>
        <h3 class="text-xl font-semibold text-gray-900 mb-1">Event RSVPs</h3>
        <p class="text-sm text-gray-600">Manage event registrations and attendance</p>
      </div>
      <div class="mt-4 sm:mt-0 flex items-center space-x-3">
        <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-indigo-100 text-indigo-800">
          {{ totalRSVPs }} Total RSVPs
        </span>
        <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-green-100 text-green-800">
          {{ confirmedCount }} Confirmed
        </span>
      </div>
    </div>

    <!-- Filters and Search -->
    <div class="mb-6 grid grid-cols-1 md:grid-cols-3 gap-4">
      <!-- Search -->
      <div class="relative">
        <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
          <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
        </div>
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by name, phone, email..."
          class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm"
        />
      </div>

      <!-- Status Filter -->
      <select
        v-model="selectedStatus"
        class="block w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm"
      >
        <option value="all">All Status</option>
        <option value="confirmed">Confirmed</option>
        <option value="pending">Pending</option>
        <option value="declined">Declined</option>
        <option value="attended">Attended</option>
      </select>

      <!-- Bulk Actions -->
      <div class="flex space-x-2">
        <select
          v-model="selectedAction"
          class="flex-1 px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-transparent text-sm"
          :disabled="selectedRSVPs.length === 0"
        >
          <option value="">Bulk Actions</option>
          <option value="confirm">Mark as Confirmed</option>
          <option value="decline">Mark as Declined</option>
          <option value="attend">Mark as Attended</option>
          <option value="send_reminder">Send Reminder</option>
        </select>
        <button
          @click="executeBulkAction"
          :disabled="!selectedAction || selectedRSVPs.length === 0"
          class="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          Apply
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      <span class="ml-3 text-gray-600">Loading RSVPs...</span>
    </div>

    <!-- RSVP List -->
    <div v-else-if="filteredRSVPs.length > 0" class="space-y-4">
      <!-- Select All -->
      <div class="flex items-center space-x-3 pb-3 border-b border-gray-200">
        <input
          type="checkbox"
          :checked="isAllSelected"
          @change="toggleSelectAll"
          class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
        />
        <label class="text-sm font-medium text-gray-700">
          Select All ({{ filteredRSVPs.length }} items)
        </label>
      </div>

      <!-- RSVP Cards -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
        <div
          v-for="rsvp in filteredRSVPs"
          :key="rsvp.id"
          class="border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow duration-200"
          :class="{'ring-2 ring-indigo-500 border-transparent': selectedRSVPs.includes(rsvp.id)}"
        >
          <div class="flex items-start justify-between">
            <div class="flex items-start space-x-3">
              <input
                type="checkbox"
                :value="rsvp.id"
                v-model="selectedRSVPs"
                class="mt-1 h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
              />
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-2 mb-1">
                  <h4 class="text-sm font-medium text-gray-900 truncate">{{ rsvp.devoteeName }}</h4>
                  <span
                    class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium"
                    :class="getStatusBadgeClass(rsvp.status)"
                  >
                    {{ rsvp.status.charAt(0).toUpperCase() + rsvp.status.slice(1) }}
                  </span>
                </div>
                <div class="space-y-1">
                  <p class="text-xs text-gray-600 flex items-center">
                    <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                    </svg>
                    {{ rsvp.devoteePhone }}
                  </p>
                  <p class="text-xs text-gray-600 flex items-center">
                    <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 12a4 4 0 10-8 0 4 4 0 008 0zm0 0v1.5a2.5 2.5 0 005 0V12a9 9 0 10-9 9m4.5-1.206a8.959 8.959 0 01-4.5 1.207" />
                    </svg>
                    {{ rsvp.devoteeEmail }}
                  </p>
                  <p class="text-xs text-gray-500">
                    RSVP Date: {{ formatDate(rsvp.rsvpDate) }}
                  </p>
                  <p v-if="rsvp.familyMembers > 0" class="text-xs text-gray-500">
                    Family Members: +{{ rsvp.familyMembers }}
                  </p>
                  <p v-if="rsvp.notes" class="text-xs text-gray-600 mt-2 italic">
                    "{{ rsvp.notes }}"
                  </p>
                </div>
              </div>
            </div>
            
            <!-- Action Dropdown -->
            <div class="relative" v-if="showActions">
              <button
                @click="toggleActionMenu(rsvp.id)"
                class="p-1 text-gray-400 hover:text-gray-600 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 rounded"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M10 6a2 2 0 110-4 2 2 0 010 4zM10 12a2 2 0 110-4 2 2 0 010 4zM10 18a2 2 0 110-4 2 2 0 010 4z" />
                </svg>
              </button>
              
              <div
                v-if="activeActionMenu === rsvp.id"
                class="absolute right-0 mt-2 w-48 bg-white rounded-md shadow-lg z-10 border border-gray-200"
              >
                <div class="py-1">
                  <button
                    @click="updateRSVPStatus(rsvp.id, 'confirmed')"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Mark as Confirmed
                  </button>
                  <button
                    @click="updateRSVPStatus(rsvp.id, 'declined')"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Mark as Declined
                  </button>
                  <button
                    @click="updateRSVPStatus(rsvp.id, 'attended')"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Mark as Attended
                  </button>
                  <hr class="my-1">
                  <button
                    @click="sendReminder(rsvp.id)"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    Send Reminder
                  </button>
                  <button
                    @click="viewDevoteeProfile(rsvp.devoteeId)"
                    class="block w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-gray-100"
                  >
                    View Profile
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
      </svg>
      <h3 class="mt-4 text-sm font-medium text-gray-900">No RSVPs found</h3>
      <p class="mt-2 text-sm text-gray-500">
        {{ searchQuery || selectedStatus !== 'all' ? 'Try adjusting your filters.' : 'No one has RSVP\'d for this event yet.' }}
      </p>
    </div>

    <!-- Summary Stats -->
    <div v-if="rsvpList.length > 0" class="mt-6 grid grid-cols-2 md:grid-cols-4 gap-4 pt-6 border-t border-gray-200">
      <div class="text-center">
        <div class="text-2xl font-bold text-indigo-600">{{ totalRSVPs }}</div>
        <div class="text-sm text-gray-600">Total RSVPs</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-bold text-green-600">{{ confirmedCount }}</div>
        <div class="text-sm text-gray-600">Confirmed</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-bold text-yellow-600">{{ pendingCount }}</div>
        <div class="text-sm text-gray-600">Pending</div>
      </div>
      <div class="text-center">
        <div class="text-2xl font-bold text-purple-600">{{ attendedCount }}</div>
        <div class="text-sm text-gray-600">Attended</div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted, onUnmounted } from 'vue'

export default {
  name: 'EventRSVP',
  props: {
    eventId: {
      type: String,
      required: true
    },
    showActions: {
      type: Boolean,
      default: true
    }
  },
  emits: ['rsvp-updated', 'reminder-sent'],
  setup(props, { emit }) {
    const rsvpList = ref([])
    const isLoading = ref(false)
    const searchQuery = ref('')
    const selectedStatus = ref('all')
    const selectedAction = ref('')
    const selectedRSVPs = ref([])
    const activeActionMenu = ref(null)

    // Mock data - replace with actual API call
    const mockRSVPData = [
      {
        id: '1',
        devoteeId: 'dev_001',
        devoteeName: 'Rajesh Kumar',
        devoteePhone: '+91 98765 43210',
        devoteeEmail: 'rajesh@example.com',
        status: 'confirmed',
        rsvpDate: '2024-01-15T10:30:00Z',
        familyMembers: 3,
        notes: 'Coming with family for the special pooja'
      },
      {
        id: '2',
        devoteeId: 'dev_002',
        devoteeName: 'Priya Sharma',
        devoteePhone: '+91 98765 43211',
        devoteeEmail: 'priya@example.com',
        status: 'pending',
        rsvpDate: '2024-01-16T14:20:00Z',
        familyMembers: 0,
        notes: ''
      },
      {
        id: '3',
        devoteeId: 'dev_003',
        devoteeName: 'Amit Patel',
        devoteePhone: '+91 98765 43212',
        devoteeEmail: 'amit@example.com',
        status: 'attended',
        rsvpDate: '2024-01-14T09:15:00Z',
        familyMembers: 1,
        notes: 'Will help with prasadam distribution'
      },
      {
        id: '4',
        devoteeId: 'dev_004',
        devoteeName: 'Sunita Reddy',
        devoteePhone: '+91 98765 43213',
        devoteeEmail: 'sunita@example.com',
        status: 'declined',
        rsvpDate: '2024-01-17T16:45:00Z',
        familyMembers: 0,
        notes: 'Unable to attend due to travel'
      }
    ]

    // Computed properties
    const filteredRSVPs = computed(() => {
      let filtered = rsvpList.value

      // Filter by search query
      if (searchQuery.value) {
        const query = searchQuery.value.toLowerCase()
        filtered = filtered.filter(rsvp =>
          rsvp.devoteeName.toLowerCase().includes(query) ||
          rsvp.devoteePhone.includes(query) ||
          rsvp.devoteeEmail.toLowerCase().includes(query)
        )
      }

      // Filter by status
      if (selectedStatus.value !== 'all') {
        filtered = filtered.filter(rsvp => rsvp.status === selectedStatus.value)
      }

      return filtered
    })

    const totalRSVPs = computed(() => rsvpList.value.length)
    const confirmedCount = computed(() => rsvpList.value.filter(r => r.status === 'confirmed').length)
    const pendingCount = computed(() => rsvpList.value.filter(r => r.status === 'pending').length)
    const attendedCount = computed(() => rsvpList.value.filter(r => r.status === 'attended').length)
    
    const isAllSelected = computed({
      get: () => filteredRSVPs.value.length > 0 && selectedRSVPs.value.length === filteredRSVPs.value.length,
      set: (value) => {
        selectedRSVPs.value = value ? filteredRSVPs.value.map(r => r.id) : []
      }
    })

    // Methods
    const loadRSVPs = async () => {
      isLoading.value = true
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1000))
        rsvpList.value = mockRSVPData
      } catch (error) {
        console.error('Failed to load RSVPs:', error)
      } finally {
        isLoading.value = false
      }
    }

    const toggleSelectAll = () => {
      isAllSelected.value = !isAllSelected.value
    }

    const toggleActionMenu = (rsvpId) => {
      activeActionMenu.value = activeActionMenu.value === rsvpId ? null : rsvpId
    }

    const updateRSVPStatus = async (rsvpId, newStatus) => {
      try {
        const rsvp = rsvpList.value.find(r => r.id === rsvpId)
        if (rsvp) {
          rsvp.status = newStatus
          emit('rsvp-updated', { rsvpId, status: newStatus })
          
          // Show success message based on status
          const messages = {
            confirmed: 'RSVP marked as confirmed',
            declined: 'RSVP marked as declined',
            attended: 'Attendance marked successfully'
          }
          
          // You would show a toast notification here
          console.log(messages[newStatus] || 'RSVP updated successfully')
        }
      } catch (error) {
        console.error('Failed to update RSVP status:', error)
      }
      activeActionMenu.value = null
    }

    const sendReminder = async (rsvpId) => {
      try {
        // Simulate API call to send reminder
        await new Promise(resolve => setTimeout(resolve, 500))
        emit('reminder-sent', { rsvpId })
        console.log('Reminder sent successfully')
      } catch (error) {
        console.error('Failed to send reminder:', error)
      }
      activeActionMenu.value = null
    }

    const viewDevoteeProfile = (devoteeId) => {
      // Navigate to devotee profile or open modal
      console.log('View profile for devotee:', devoteeId)
      activeActionMenu.value = null
    }

    const executeBulkAction = async () => {
      if (!selectedAction.value || selectedRSVPs.value.length === 0) return

      try {
        switch (selectedAction.value) {
          case 'confirm':
            selectedRSVPs.value.forEach(rsvpId => {
              const rsvp = rsvpList.value.find(r => r.id === rsvpId)
              if (rsvp) rsvp.status = 'confirmed'
            })
            console.log(`${selectedRSVPs.value.length} RSVPs marked as confirmed`)
            break
          case 'decline':
            selectedRSVPs.value.forEach(rsvpId => {
              const rsvp = rsvpList.value.find(r => r.id === rsvpId)
              if (rsvp) rsvp.status = 'declined'
            })
            console.log(`${selectedRSVPs.value.length} RSVPs marked as declined`)
            break
          case 'attend':
            selectedRSVPs.value.forEach(rsvpId => {
              const rsvp = rsvpList.value.find(r => r.id === rsvpId)
              if (rsvp) rsvp.status = 'attended'
            })
            console.log(`${selectedRSVPs.value.length} RSVPs marked as attended`)
            break
          case 'send_reminder':
            console.log(`Reminder sent to ${selectedRSVPs.value.length} devotees`)
            break
        }
        
        selectedRSVPs.value = []
        selectedAction.value = ''
      } catch (error) {
        console.error('Failed to execute bulk action:', error)
      }
    }

    const getStatusBadgeClass = (status) => {
      const classes = {
        confirmed: 'bg-green-100 text-green-800',
        pending: 'bg-yellow-100 text-yellow-800',
        declined: 'bg-red-100 text-red-800',
        attended: 'bg-purple-100 text-purple-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    }

    const formatDate = (dateString) => {
      return new Date(dateString).toLocaleDateString('en-IN', {
        year: 'numeric',
        month: 'short',
        day: 'numeric',
        hour: '2-digit',
        minute: '2-digit'
      })
    }

    // Close action menu when clicking outside
    const handleClickOutside = (event) => {
      if (!event.target.closest('.relative')) {
        activeActionMenu.value = null
      }
    }

    onMounted(() => {
      loadRSVPs()
      document.addEventListener('click', handleClickOutside)
    })

    onUnmounted(() => {
      document.removeEventListener('click', handleClickOutside)
    })

    return {
      rsvpList,
      isLoading,
      searchQuery,
      selectedStatus,
      selectedAction,
      selectedRSVPs,
      activeActionMenu,
      filteredRSVPs,
      totalRSVPs,
      confirmedCount,
      pendingCount,
      attendedCount,
      isAllSelected,
      toggleSelectAll,
      toggleActionMenu,
      updateRSVPStatus,
      sendReminder,
      viewDevoteeProfile,
      executeBulkAction,
      getStatusBadgeClass,
      formatDate
    }
  }
}
</script>