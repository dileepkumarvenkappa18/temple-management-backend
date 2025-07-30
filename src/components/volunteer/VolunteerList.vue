<template>
  <div class="bg-white rounded-xl shadow-lg overflow-hidden">
    <!-- Header Section -->
    <div class="bg-gradient-to-r from-indigo-600 to-indigo-700 px-6 py-4">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-bold text-white font-roboto">Volunteer Management</h2>
          <p class="text-indigo-100 text-sm mt-1">Manage temple volunteers and their assignments</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <button 
            @click="$emit('add-volunteer')"
            class="bg-white text-indigo-600 px-4 py-2 rounded-lg font-medium hover:bg-indigo-50 transition-all duration-200 flex items-center gap-2 shadow-md"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
            </svg>
            Add Volunteer
          </button>
          <button 
            @click="exportVolunteers"
            class="bg-indigo-800 text-white px-4 py-2 rounded-lg font-medium hover:bg-indigo-900 transition-all duration-200 flex items-center gap-2 shadow-md"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            Export
          </button>
        </div>
      </div>
    </div>

    <!-- Search and Filter Section -->
    <div class="p-6 border-b border-gray-200 bg-gray-50">
      <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
        <!-- Search Input -->
        <div class="relative flex-1 max-w-md">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search volunteers..."
            class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          />
        </div>

        <!-- Filters -->
        <div class="flex flex-col sm:flex-row gap-3">
          <select 
            v-model="selectedStatus"
            class="px-3 py-2 border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          >
            <option value="">All Status</option>
            <option value="active">Active</option>
            <option value="inactive">Inactive</option>
          </select>
          
          <select 
            v-model="selectedRole"
            class="px-3 py-2 border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          >
            <option value="">All Roles</option>
            <option value="seva_coordinator">Seva Coordinator</option>
            <option value="event_manager">Event Manager</option>
            <option value="general_volunteer">General Volunteer</option>
          </select>
        </div>
      </div>
    </div>

    <!-- Statistics Cards -->
    <div class="p-6 bg-gray-50 border-b border-gray-200">
      <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
        <div class="bg-white rounded-lg p-4 shadow-md border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Total Volunteers</p>
              <p class="text-2xl font-bold text-gray-900">{{ volunteerStats.total }}</p>
            </div>
            <div class="p-2 bg-indigo-100 rounded-lg">
              <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg p-4 shadow-md border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Active</p>
              <p class="text-2xl font-bold text-green-600">{{ volunteerStats.active }}</p>
            </div>
            <div class="p-2 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg p-4 shadow-md border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">This Month</p>
              <p class="text-2xl font-bold text-blue-600">{{ volunteerStats.thisMonth }}</p>
            </div>
            <div class="p-2 bg-blue-100 rounded-lg">
              <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
              </svg>
            </div>
          </div>
        </div>
        
        <div class="bg-white rounded-lg p-4 shadow-md border border-gray-100">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600">Coordinators</p>
              <p class="text-2xl font-bold text-purple-600">{{ volunteerStats.coordinators }}</p>
            </div>
            <div class="p-2 bg-purple-100 rounded-lg">
              <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5.121 17.804A13.937 13.937 0 0112 16c2.5 0 4.847.655 6.879 1.804M15 10a3 3 0 11-6 0 3 3 0 016 0zm6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Table Section -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-50 border-b border-gray-200">
          <tr>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Volunteer
            </th>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Contact
            </th>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Role
            </th>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Joined
            </th>
            <th class="px-6 py-4 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr 
            v-for="volunteer in filteredVolunteers" 
            :key="volunteer.id"
            class="hover:bg-gray-50 transition-colors duration-200"
          >
            <!-- Volunteer Info -->
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="flex-shrink-0 h-10 w-10">
                  <div class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
                    <span class="font-medium text-indigo-600 text-sm">
                      {{ getInitials(volunteer.full_name) }}
                    </span>
                  </div>
                </div>
                <div class="ml-4">
                  <div class="text-sm font-medium text-gray-900">{{ volunteer.full_name }}</div>
                  <div class="text-sm text-gray-500">ID: {{ volunteer.volunteer_id }}</div>
                </div>
              </div>
            </td>

            <!-- Contact -->
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ volunteer.email }}</div>
              <div class="text-sm text-gray-500">{{ volunteer.phone || 'N/A' }}</div>
            </td>

            <!-- Role -->
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getRoleBadgeClass(volunteer.role)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                {{ formatRole(volunteer.role) }}
              </span>
            </td>

            <!-- Status -->
            <td class="px-6 py-4 whitespace-nowrap">
              <span :class="getStatusBadgeClass(volunteer.status)" class="inline-flex px-2 py-1 text-xs font-semibold rounded-full">
                {{ volunteer.status }}
              </span>
            </td>

            <!-- Joined Date -->
            <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
              {{ formatDate(volunteer.joined_date) }}
            </td>

            <!-- Actions -->
            <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
              <div class="flex items-center gap-2">
                <button 
                  @click="viewVolunteer(volunteer)"
                  class="text-indigo-600 hover:text-indigo-900 transition-colors duration-200 p-1 rounded hover:bg-indigo-50"
                  title="View Details"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                </button>
                <button 
                  @click="editVolunteer(volunteer)"
                  class="text-blue-600 hover:text-blue-900 transition-colors duration-200 p-1 rounded hover:bg-blue-50"
                  title="Edit"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                </button>
                <button 
                  @click="toggleStatus(volunteer)"
                  :class="volunteer.status === 'active' ? 'text-red-600 hover:text-red-900 hover:bg-red-50' : 'text-green-600 hover:text-green-900 hover:bg-green-50'"
                  class="transition-colors duration-200 p-1 rounded"
                  :title="volunteer.status === 'active' ? 'Deactivate' : 'Activate'"
                >
                  <svg v-if="volunteer.status === 'active'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="filteredVolunteers.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No volunteers found</h3>
      <p class="mt-1 text-sm text-gray-500">Try adjusting your search or filter criteria.</p>
    </div>

    <!-- Pagination -->
    <div v-if="totalPages > 1" class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
      <div class="flex items-center justify-between">
        <div class="flex-1 flex justify-between sm:hidden">
          <button 
            @click="previousPage"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            Previous
          </button>
          <button 
            @click="nextPage"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing <span class="font-medium">{{ startItem }}</span> to <span class="font-medium">{{ endItem }}</span> of <span class="font-medium">{{ totalItems }}</span> results
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px">
              <button 
                @click="previousPage"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              <button 
                v-for="page in visiblePages" 
                :key="page"
                @click="goToPage(page)"
                :class="page === currentPage ? 'bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
                class="relative inline-flex items-center px-4 py-2 border text-sm font-medium transition-all duration-200"
              >
                {{ page }}
              </button>
              <button 
                @click="nextPage"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Define props
const props = defineProps({
  volunteers: {
    type: Array,
    default: () => []
  },
  loading: {
    type: Boolean,
    default: false
  }
})

// Define emits
const emit = defineEmits(['add-volunteer', 'view-volunteer', 'edit-volunteer', 'toggle-status'])

// Reactive data
const searchQuery = ref('')
const selectedStatus = ref('')
const selectedRole = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(10)

// Mock data for demonstration
const mockVolunteers = ref([
  {
    id: 1,
    volunteer_id: 'VOL001',
    full_name: 'Rajesh Kumar',
    email: 'rajesh@temple.com',
    phone: '+91 98765 43210',
    role: 'seva_coordinator',
    status: 'active',
    joined_date: '2024-01-15'
  },
  {
    id: 2,
    volunteer_id: 'VOL002',
    full_name: 'Priya Sharma',
    email: 'priya@temple.com',
    phone: '+91 87654 32109',
    role: 'event_manager',
    status: 'active',
    joined_date: '2024-02-20'
  },
  {
    id: 3,
    volunteer_id: 'VOL003',
    full_name: 'Amit Patel',
    email: 'amit@temple.com',
    phone: null,
    role: 'general_volunteer',
    status: 'inactive',
    joined_date: '2024-03-10'
  }
])

// Computed properties
const allVolunteers = computed(() => {
  return props.volunteers.length > 0 ? props.volunteers : mockVolunteers.value
})

const filteredVolunteers = computed(() => {
  let filtered = allVolunteers.value

  // Search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(volunteer => 
      volunteer.full_name.toLowerCase().includes(query) ||
      volunteer.email.toLowerCase().includes(query) ||
      volunteer.volunteer_id.toLowerCase().includes(query)
    )
  }

  // Status filter
  if (selectedStatus.value) {
    filtered = filtered.filter(volunteer => volunteer.status === selectedStatus.value)
  }

  // Role filter
  if (selectedRole.value) {
    filtered = filtered.filter(volunteer => volunteer.role === selectedRole.value)
  }

  // Pagination
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filtered.slice(start, end)
})

const volunteerStats = computed(() => {
  const total = allVolunteers.value.length
  const active = allVolunteers.value.filter(v => v.status === 'active').length
  const thisMonth = allVolunteers.value.filter(v => {
    const joinedDate = new Date(v.joined_date)
    const now = new Date()
    return joinedDate.getMonth() === now.getMonth() && joinedDate.getFullYear() === now.getFullYear()
  }).length
  const coordinators = allVolunteers.value.filter(v => v.role === 'seva_coordinator').length

  return { total, active, thisMonth, coordinators }
})

const totalItems = computed(() => {
  let filtered = allVolunteers.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(volunteer => 
      volunteer.full_name.toLowerCase().includes(query) ||
      volunteer.email.toLowerCase().includes(query) ||
      volunteer.volunteer_id.toLowerCase().includes(query)
    )
  }

  if (selectedStatus.value) {
    filtered = filtered.filter(volunteer => volunteer.status === selectedStatus.value)
  }

  if (selectedRole.value) {
    filtered = filtered.filter(volunteer => volunteer.role === selectedRole.value)
  }

  return filtered.length
})

const totalPages = computed(() => Math.ceil(totalItems.value / itemsPerPage.value))

const startItem = computed(() => (currentPage.value - 1) * itemsPerPage.value + 1)
const endItem = computed(() => Math.min(currentPage.value * itemsPerPage.value, totalItems.value))

const visiblePages = computed(() => {
  const pages = []
  const maxVisible = 5
  let start = Math.max(1, currentPage.value - Math.floor(maxVisible / 2))
  let end = Math.min(totalPages.value, start + maxVisible - 1)
  
  if (end - start + 1 < maxVisible) {
    start = Math.max(1, end - maxVisible + 1)
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// Methods
const getInitials = (name) => {
  return name.split(' ').map(n => n[0]).join('').toUpperCase()
}

const getRoleBadgeClass = (role) => {
  const classes = {
    seva_coordinator: 'bg-purple-100 text-purple-800',
    event_manager: 'bg-blue-100 text-blue-800',
    general_volunteer: 'bg-gray-100 text-gray-800'
  }
  return classes[role] || 'bg-gray-100 text-gray-800'
}

const getStatusBadgeClass = (status) => {
  return status === 'active' 
    ? 'bg-green-100 text-green-800' 
    : 'bg-red-100 text-red-800'
}

const formatRole = (role) => {
  const roleMap = {
    seva_coordinator: 'Seva Coordinator',
    event_manager: 'Event Manager',
    general_volunteer: 'General Volunteer'
  }
  return roleMap[role] || role
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    year: 'numeric',
    month: 'short',
    day: 'numeric'
  })
}

const viewVolunteer = (volunteer) => {
  emit('view-volunteer', volunteer)
}

const editVolunteer = (volunteer) => {
  emit('edit-volunteer', volunteer)
}

const toggleStatus = (volunteer) => {
  emit('toggle-status', volunteer)
}

const exportVolunteers = () => {
  // Export functionality implementation
  console.log('Exporting volunteers...')
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

// Lifecycle
onMounted(() => {
  // Initialize component
})
</script>

<style scoped>
.font-roboto {
  font-family: 'Roboto', sans-serif;
}
</style>