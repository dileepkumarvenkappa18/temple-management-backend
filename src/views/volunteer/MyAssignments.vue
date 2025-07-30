<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-6">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900 font-['Roboto']">My Assignments</h1>
              <p class="mt-1 text-sm text-gray-600 font-sans">
                View and manage your volunteering assignments and schedules
              </p>
            </div>
            <div class="mt-4 sm:mt-0 flex flex-col sm:flex-row gap-3">
              <!-- <button
                @click="showAvailabilityModal = true"
                class="inline-flex items-center px-4 py-2 border border-indigo-300 text-sm font-medium rounded-lg text-indigo-700 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"/>
                </svg>
                Update Availability
              </button> -->
              <button
                @click="refreshAssignments"
                :disabled="loading"
                class="inline-flex items-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
              >
                <svg 
                  class="w-4 h-4 mr-2" 
                  :class="{ 'animate-spin': loading }" 
                  fill="none" 
                  stroke="currentColor" 
                  viewBox="0 0 24 24"
                >
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"/>
                </svg>
                Refresh
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Cards -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white overflow-hidden shadow-md rounded-xl">
          <div class="p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-10 h-10 bg-indigo-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <dt class="text-sm font-medium text-gray-500 font-sans">Total Assignments</dt>
                <dd class="text-2xl font-bold text-gray-900">{{ stats.totalAssignments }}</dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow-md rounded-xl">
          <div class="p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-10 h-10 bg-green-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <dt class="text-sm font-medium text-gray-500 font-sans">Completed</dt>
                <dd class="text-2xl font-bold text-gray-900">{{ stats.completed }}</dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow-md rounded-xl">
          <div class="p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-10 h-10 bg-yellow-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <dt class="text-sm font-medium text-gray-500 font-sans">Pending</dt>
                <dd class="text-2xl font-bold text-gray-900">{{ stats.pending }}</dd>
              </div>
            </div>
          </div>
        </div>

        <div class="bg-white overflow-hidden shadow-md rounded-xl">
          <div class="p-6">
            <div class="flex items-center">
              <div class="flex-shrink-0">
                <div class="w-10 h-10 bg-blue-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                  </svg>
                </div>
              </div>
              <div class="ml-4">
                <dt class="text-sm font-medium text-gray-500 font-sans">This Week</dt>
                <dd class="text-2xl font-bold text-gray-900">{{ stats.thisWeek }}</dd>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Filter Tabs -->
      <div class="bg-white shadow-md rounded-xl mb-6">
        <div class="border-b border-gray-200">
          <nav class="-mb-px flex space-x-8 px-6" aria-label="Tabs">
            <button
              v-for="tab in filterTabs"
              :key="tab.id"
              @click="activeTab = tab.id"
              :class="[
                activeTab === tab.id
                  ? 'border-indigo-500 text-indigo-600'
                  : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300',
                'whitespace-nowrap py-4 px-1 border-b-2 font-medium text-sm transition-all duration-200'
              ]"
            >
              {{ tab.name }}
              <span
                v-if="tab.count"
                :class="[
                  activeTab === tab.id
                    ? 'bg-indigo-100 text-indigo-600'
                    : 'bg-gray-100 text-gray-900',
                  'ml-2 py-0.5 px-2.5 rounded-full text-xs font-medium'
                ]"
              >
                {{ tab.count }}
              </span>
            </button>
          </nav>
        </div>
      </div>

      <!-- Assignments List -->
      <div class="bg-white shadow-md rounded-xl">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
            <h3 class="text-lg font-semibold text-gray-900 font-['Roboto']">
              {{ getTabName(activeTab) }} Assignments
            </h3>
            <div class="mt-3 sm:mt-0 flex items-center space-x-3">
              <div class="relative">
                <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search assignments..."
                  class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg leading-5 bg-white placeholder-gray-500 focus:outline-none focus:placeholder-gray-400 focus:ring-1 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
                >
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"/>
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="divide-y divide-gray-200">
          <div v-if="loading" class="p-8 text-center">
            <div class="inline-flex items-center px-4 py-2 font-semibold leading-6 text-sm shadow rounded-md text-indigo-600 bg-indigo-50">
              <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24">
                <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              Loading assignments...
            </div>
          </div>

          <div v-else-if="filteredAssignments.length === 0" class="p-8 text-center">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No assignments found</h3>
            <p class="mt-1 text-sm text-gray-500">
              {{ searchQuery ? 'Try adjusting your search criteria.' : 'No assignments match the current filter.' }}
            </p>
          </div>

          <div
            v-else
            v-for="assignment in filteredAssignments"
            :key="assignment.id"
            class="p-6 hover:bg-gray-50 transition-colors duration-200"
          >
            <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between">
              <div class="flex-1 min-w-0">
                <div class="flex items-center space-x-3">
                  <div
                    :class="[
                      'flex-shrink-0 w-3 h-3 rounded-full',
                      getStatusColor(assignment.status)
                    ]"
                  ></div>
                  <h4 class="text-lg font-semibold text-gray-900 font-['Roboto'] truncate">
                    {{ assignment.title }}
                  </h4>
                  <span
                    :class="[
                      'inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium',
                      getStatusBadgeColor(assignment.status)
                    ]"
                  >
                    {{ assignment.status }}
                  </span>
                </div>
                
                <div class="mt-2 flex flex-col sm:flex-row sm:items-center text-sm text-gray-500 space-y-1 sm:space-y-0 sm:space-x-6">
                  <div class="flex items-center">
                    <svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    {{ formatDate(assignment.date) }}
                  </div>
                  <div class="flex items-center">
                    <svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                    </svg>
                    {{ assignment.startTime }} - {{ assignment.endTime }}
                  </div>
                  <div class="flex items-center">
                    <svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                    </svg>
                    {{ assignment.location }}
                  </div>
                </div>
                
                <p class="mt-2 text-sm text-gray-600 font-sans">
                  {{ assignment.description }}
                </p>
                
                <div v-if="assignment.coordinator" class="mt-2 flex items-center text-sm text-gray-500">
                  <svg class="flex-shrink-0 mr-1.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                  </svg>
                  Coordinator: {{ assignment.coordinator }}
                </div>
              </div>
              
              <div class="mt-4 lg:mt-0 lg:ml-6 flex flex-col sm:flex-row gap-2">
                <button
                  v-if="assignment.status === 'pending'"
                  @click="acceptAssignment(assignment.id)"
                  class="inline-flex items-center px-3 py-2 border border-transparent text-sm leading-4 font-medium rounded-lg text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 transition-all duration-200"
                >
                  Accept
                </button>
                
                <button
                  v-if="assignment.status === 'pending'"
                  @click="declineAssignment(assignment.id)"
                  class="inline-flex items-center px-3 py-2 border border-red-300 text-sm leading-4 font-medium rounded-lg text-red-700 bg-red-50 hover:bg-red-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500 transition-all duration-200"
                >
                  Decline
                </button>
                
                <button
                  @click="viewAssignmentDetails(assignment)"
                  class="inline-flex items-center px-3 py-2 border border-gray-300 text-sm leading-4 font-medium rounded-lg text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200"
                >
                  View Details
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Availability Modal -->
    <div v-if="showAvailabilityModal" class="fixed inset-0 z-50 overflow-y-auto">
      <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
        <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="showAvailabilityModal = false"></div>
        
        <span class="hidden sm:inline-block sm:align-middle sm:h-screen">&#8203;</span>
        
        <div class="inline-block align-bottom bg-white rounded-xl text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-lg sm:w-full">
          <div class="bg-white px-6 pt-6 pb-4">
            <div class="flex items-center justify-between">
              <h3 class="text-lg leading-6 font-bold text-gray-900 font-['Roboto']">
                Update Availability
              </h3>
              <button
                @click="showAvailabilityModal = false"
                class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
              >
                <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
                </svg>
              </button>
            </div>
            <p class="mt-2 text-sm text-gray-600 font-sans">
              Set your availability preferences for volunteer assignments.
            </p>
          </div>
          
          <div class="px-6 pb-6">
            <div class="space-y-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 font-sans mb-2">
                  Preferred Days
                </label>
                <div class="grid grid-cols-2 gap-2">
                  <label v-for="day in weekDays" :key="day" class="flex items-center">
                    <input
                      v-model="availability.days"
                      :value="day"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    >
                    <span class="ml-2 text-sm text-gray-700 font-sans">{{ day }}</span>
                  </label>
                </div>
              </div>
              
              <div>
                <label class="block text-sm font-medium text-gray-700 font-sans mb-2">
                  Preferred Time Slots
                </label>
                <div class="space-y-2">
                  <label v-for="slot in timeSlots" :key="slot" class="flex items-center">
                    <input
                      v-model="availability.timeSlots"
                      :value="slot"
                      type="checkbox"
                      class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                    >
                    <span class="ml-2 text-sm text-gray-700 font-sans">{{ slot }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>
          
          <div class="bg-gray-50 px-6 py-4 flex flex-col sm:flex-row-reverse gap-3">
            <button
              @click="saveAvailability"
              class="w-full inline-flex justify-center px-4 py-2 border border-transparent text-sm font-medium rounded-lg text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 sm:w-auto"
            >
              Save Changes
            </button>
            <button
              @click="showAvailabilityModal = false"
              class="w-full inline-flex justify-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-lg text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 sm:w-auto"
            >
              Cancel
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Reactive data
const loading = ref(false)
const searchQuery = ref('')
const activeTab = ref('all')
const showAvailabilityModal = ref(false)

// Stats data
const stats = ref({
  totalAssignments: 12,
  completed: 8,
  pending: 3,
  thisWeek: 2
})

// Filter tabs configuration
const filterTabs = ref([
  { id: 'all', name: 'All', count: 12 },
  { id: 'pending', name: 'Pending', count: 3 },
  { id: 'accepted', name: 'Accepted', count: 6 },
  { id: 'completed', name: 'Completed', count: 8 },
  { id: 'upcoming', name: 'Upcoming', count: 4 }
])

// Availability data
const availability = ref({
  days: ['Monday', 'Wednesday', 'Friday'],
  timeSlots: ['Morning (6AM-12PM)', 'Evening (6PM-10PM)']
})

const weekDays = ['Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', 'Sunday']
const timeSlots = [
  'Early Morning (4AM-8AM)',
  'Morning (6AM-12PM)', 
  'Afternoon (12PM-6PM)',
  'Evening (6PM-10PM)',
  'Night (10PM-12AM)'
]

// Mock assignments data
const assignments = ref([
  {
    id: 1,
    title: 'Temple Cleaning - Main Hall',
    description: 'Assist with daily cleaning and maintenance of the main prayer hall including floor mopping, decoration arrangement, and general upkeep.',
    date: '2025-06-22',
    startTime: '6:00 AM',
    endTime: '8:00 AM',
    location: 'Main Prayer Hall',
    status: 'pending',
    coordinator: 'Radhika Sharma'
  },
  {
    id: 2,
    title: 'Festival Decoration Setup',
    description: 'Help set up decorations for the upcoming Janmashtami celebration. Tasks include hanging garlands, arranging flowers, and setting up stage decorations.',
    date: '2025-06-23',
    startTime: '4:00 PM',
    endTime: '8:00 PM',
    location: 'Temple Courtyard',
    status: 'accepted',
    coordinator: 'Priya Nair'
  },
  {
    id: 3,
    title: 'Prasadam Distribution',
    description: 'Assist in preparing and distributing prasadam to devotees during evening aarti. Includes serving, crowd management, and cleanup.',
    date: '2025-06-21',
    startTime: '6:00 PM',
    endTime: '8:30 PM',
    location: 'Distribution Counter',
    status: 'completed',
    coordinator: 'Amit Kumar'
  },
  {
    id: 4,
    title: 'Library Organization',
    description: 'Help organize religious books and scriptures in the temple library. Tasks include cataloging, shelving, and maintaining reading area.',
    date: '2025-06-24',
    startTime: '10:00 AM',
    endTime: '1:00 PM',
    location: 'Temple Library',
    status: 'accepted',
    coordinator: 'Dr. Sunita Mehta'
  },
  {
    id: 5,
    title: 'Youth Program Coordination',
    description: 'Assist with organizing and conducting youth spiritual education program. Help with activities, storytelling, and managing children.',
    date: '2025-06-25',
    startTime: '3:00 PM',
    endTime: '5:00 PM',
    location: 'Community Hall',
    status: 'pending',
    coordinator: 'Rajesh Gupta'
  }
])

// Computed properties
const filteredAssignments = computed(() => {
  let filtered = assignments.value

  // Filter by active tab
  if (activeTab.value !== 'all') {
    if (activeTab.value === 'upcoming') {
      const today = new Date().toISOString().split('T')[0]
      filtered = filtered.filter(assignment => 
        assignment.date >= today && 
        (assignment.status === 'accepted' || assignment.status === 'pending')
      )
    } else {
      filtered = filtered.filter(assignment => assignment.status === activeTab.value)
    }
  }

  // Filter by search query
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(assignment =>
      assignment.title.toLowerCase().includes(query) ||
      assignment.description.toLowerCase().includes(query) ||
      assignment.location.toLowerCase().includes(query) ||
      assignment.coordinator.toLowerCase().includes(query)
    )
  }

  return filtered
})

// Methods
const getTabName = (tabId) => {
  const tab = filterTabs.value.find(t => t.id === tabId)
  return tab ? tab.name : 'All'
}

const getStatusColor = (status) => {
  const colors = {
    pending: 'bg-yellow-400',
    accepted: 'bg-blue-400',
    completed: 'bg-green-400',
    declined: 'bg-red-400'
  }
  return colors[status] || 'bg-gray-400'
}

const getStatusBadgeColor = (status) => {
  const colors = {
    pending: 'bg-yellow-100 text-yellow-800',
    accepted: 'bg-blue-100 text-blue-800', 
    completed: 'bg-green-100 text-green-800',
    declined: 'bg-red-100 text-red-800'
  }
  return colors[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    weekday: 'short', 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric'
  })
}

const refreshAssignments = async () => {
  loading.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    // In real app, fetch fresh data from API
    console.log('Refreshing assignments...')
  } catch (error) {
    console.error('Error refreshing assignments:', error)
  } finally {
    loading.value = false
  }
}

const acceptAssignment = async (assignmentId) => {
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Update assignment status
    const assignment = assignments.value.find(a => a.id === assignmentId)
    if (assignment) {
      assignment.status = 'accepted'
    }
    
    // Update stats
    stats.value.pending -= 1
    
    // Show success message
    console.log('Assignment accepted successfully!')
  } catch (error) {
    console.error('Error accepting assignment:', error)
  }
}

const declineAssignment = async (assignmentId) => {
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))
    
    // Update assignment status
    const assignment = assignments.value.find(a => a.id === assignmentId)
    if (assignment) {
      assignment.status = 'declined'
    }
    
    // Update stats
    stats.value.pending -= 1
    
    // Show success message
    console.log('Assignment declined.')
  } catch (error) {
    console.error('Error declining assignment:', error)
  }
}

const viewAssignmentDetails = (assignment) => {
  // In real app, navigate to assignment details page or open modal
  console.log('Viewing assignment details:', assignment)
  alert(`Assignment Details:\n\nTitle: ${assignment.title}\nDate: ${formatDate(assignment.date)}\nTime: ${assignment.startTime} - ${assignment.endTime}\nLocation: ${assignment.location}\nCoordinator: ${assignment.coordinator}\n\nDescription: ${assignment.description}`)
}

const saveAvailability = async () => {
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // In real app, save availability to backend
    console.log('Saving availability:', availability.value)
    
    showAvailabilityModal.value = false
    
    // Show success message
    console.log('Availability updated successfully!')
  } catch (error) {
    console.error('Error saving availability:', error)
  }
}

// Lifecycle
onMounted(() => {
  // Load initial data
  refreshAssignments()
});
</script>
