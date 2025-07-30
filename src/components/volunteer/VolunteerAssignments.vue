<template>
  <div class="bg-white rounded-2xl shadow-lg p-6 space-y-6">
    <!-- Header Section -->
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
      <div>
        <h2 class="text-2xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
          My Assignments
        </h2>
        <p class="text-gray-600 mt-1">Manage your volunteer tasks and responsibilities</p>
      </div>
      
      <!-- Status Filter -->
      <div class="flex items-center gap-3">
        <select 
          v-model="selectedStatus" 
          class="px-4 py-2 border border-gray-300 rounded-xl focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm font-medium transition-all duration-200"
        >
          <option value="all">All Status</option>
          <option value="pending">Pending</option>
          <option value="in-progress">In Progress</option>
          <option value="completed">Completed</option>
          <option value="overdue">Overdue</option>
        </select>
        
        <button
          @click="refreshAssignments"
          class="px-4 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 flex items-center gap-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
          </svg>
          Refresh
        </button>
      </div>
    </div>

    <!-- Summary Cards -->
    <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <div class="bg-gradient-to-br from-indigo-50 to-indigo-100 p-4 rounded-xl border border-indigo-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-indigo-600 text-sm font-medium">Total Assignments</p>
            <p class="text-2xl font-bold text-indigo-800">{{ totalAssignments }}</p>
          </div>
          <div class="p-2 bg-indigo-600 rounded-lg">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-yellow-50 to-yellow-100 p-4 rounded-xl border border-yellow-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-yellow-600 text-sm font-medium">Pending</p>
            <p class="text-2xl font-bold text-yellow-800">{{ pendingCount }}</p>
          </div>
          <div class="p-2 bg-yellow-500 rounded-lg">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-blue-50 to-blue-100 p-4 rounded-xl border border-blue-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-blue-600 text-sm font-medium">In Progress</p>
            <p class="text-2xl font-bold text-blue-800">{{ inProgressCount }}</p>
          </div>
          <div class="p-2 bg-blue-600 rounded-lg">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
            </svg>
          </div>
        </div>
      </div>

      <div class="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-xl border border-green-200">
        <div class="flex items-center justify-between">
          <div>
            <p class="text-green-600 text-sm font-medium">Completed</p>
            <p class="text-2xl font-bold text-green-800">{{ completedCount }}</p>
          </div>
          <div class="p-2 bg-green-600 rounded-lg">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
        </div>
      </div>
    </div>

    <!-- Assignments List -->
    <div class="space-y-4">
      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      </div>

      <div v-else-if="filteredAssignments.length === 0" class="text-center py-12">
        <div class="mx-auto w-24 h-24 bg-gray-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-12 h-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">No assignments found</h3>
        <p class="text-gray-600">{{ selectedStatus === 'all' ? 'You don\'t have any assignments yet.' : `No ${selectedStatus} assignments found.` }}</p>
      </div>

      <div v-else class="space-y-4">
        <div 
          v-for="assignment in filteredAssignments" 
          :key="assignment.id"
          class="bg-white border border-gray-200 rounded-xl p-6 hover:shadow-md transition-all duration-200"
        >
          <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
            <!-- Assignment Details -->
            <div class="flex-1 space-y-3">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="text-lg font-semibold text-gray-900 mb-1">{{ assignment.title }}</h3>
                  <p class="text-gray-600 text-sm">{{ assignment.description }}</p>
                </div>
                <span 
                  :class="getStatusClass(assignment.status)"
                  class="px-3 py-1 rounded-full text-xs font-medium ml-4"
                >
                  {{ formatStatus(assignment.status) }}
                </span>
              </div>

              <!-- Assignment Meta -->
              <div class="flex flex-wrap items-center gap-4 text-sm text-gray-600">
                <div class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                  </svg>
                  <span>Due: {{ formatDate(assignment.dueDate) }}</span>
                </div>
                
                <div class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  </svg>
                  <span>{{ assignment.location }}</span>
                </div>

                <div class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
                  </svg>
                  <span>{{ assignment.assignedBy }}</span>
                </div>

                <div v-if="assignment.priority" class="flex items-center gap-1">
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v3m0 0v3m0-3h3m-3 0H9m12 0a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                  <span :class="getPriorityClass(assignment.priority)">
                    {{ assignment.priority }} Priority
                  </span>
                </div>
              </div>
            </div>

            <!-- Action Buttons -->
            <div class="flex flex-col sm:flex-row gap-2 min-w-fit">
              <button
                v-if="assignment.status === 'pending'"
                @click="startAssignment(assignment.id)"
                class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
              >
                Start Task
              </button>
              
              <button
                v-if="assignment.status === 'in-progress'"
                @click="completeAssignment(assignment.id)"
                class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
              >
                Mark Complete
              </button>

              <button
                @click="viewAssignmentDetails(assignment.id)"
                class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
              >
                View Details
              </button>
            </div>
          </div>

          <!-- Progress Bar (for in-progress assignments) -->
          <div v-if="assignment.status === 'in-progress' && assignment.progress !== undefined" class="mt-4">
            <div class="flex items-center justify-between mb-2">
              <span class="text-sm font-medium text-gray-700">Progress</span>
              <span class="text-sm text-gray-600">{{ assignment.progress }}%</span>
            </div>
            <div class="w-full bg-gray-200 rounded-full h-2">
              <div 
                class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
                :style="{ width: assignment.progress + '%' }"
              ></div>
            </div>
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
const selectedStatus = ref('all')
const assignments = ref([])

// Mock data - replace with actual API calls
const mockAssignments = [
  {
    id: 1,
    title: 'Temple Cleaning - Main Hall',
    description: 'Deep cleaning of the main worship hall including floors, walls, and decorative items',
    status: 'pending',
    dueDate: '2025-06-25',
    location: 'Main Hall',
    assignedBy: 'Temple Admin',
    priority: 'high',
    progress: 0
  },
  {
    id: 2,
    title: 'Event Setup - Janmashtami Celebration',
    description: 'Setup decorations, sound system, and seating arrangements for the festival',
    status: 'in-progress',
    dueDate: '2025-06-22',
    location: 'Temple Grounds',
    assignedBy: 'Event Coordinator',
    priority: 'high',
    progress: 65
  },
  {
    id: 3,
    title: 'Devotee Registration Support',
    description: 'Assist new devotees with registration process and profile completion',
    status: 'completed',
    dueDate: '2025-06-20',
    location: 'Reception Desk',
    assignedBy: 'Registration Team',
    priority: 'medium',
    progress: 100
  },
  {
    id: 4,
    title: 'Seva Coordination - Morning Aarti',
    description: 'Coordinate with devotees for morning aarti seva scheduling',
    status: 'overdue',
    dueDate: '2025-06-18',
    location: 'Sanctum',
    assignedBy: 'Seva Coordinator',
    priority: 'high',
    progress: 30
  },
  {
    id: 5,
    title: 'Library Organization',
    description: 'Organize and catalog new spiritual books and scriptures',
    status: 'pending',
    dueDate: '2025-06-28',
    location: 'Temple Library',
    assignedBy: 'Librarian',
    priority: 'low',
    progress: 0
  }
]

// Computed properties
const filteredAssignments = computed(() => {
  if (selectedStatus.value === 'all') {
    return assignments.value
  }
  return assignments.value.filter(assignment => assignment.status === selectedStatus.value)
})

const totalAssignments = computed(() => assignments.value.length)
const pendingCount = computed(() => assignments.value.filter(a => a.status === 'pending').length)
const inProgressCount = computed(() => assignments.value.filter(a => a.status === 'in-progress').length)
const completedCount = computed(() => assignments.value.filter(a => a.status === 'completed').length)

// Methods
const loadAssignments = async () => {
  loading.value = true
  try {
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    assignments.value = mockAssignments
  } catch (error) {
    console.error('Error loading assignments:', error)
  } finally {
    loading.value = false
  }
}

const refreshAssignments = () => {
  loadAssignments()
}

const startAssignment = (assignmentId) => {
  const assignment = assignments.value.find(a => a.id === assignmentId)
  if (assignment) {
    assignment.status = 'in-progress'
    assignment.progress = 0
  }
}

const completeAssignment = (assignmentId) => {
  const assignment = assignments.value.find(a => a.id === assignmentId)
  if (assignment) {
    assignment.status = 'completed'
    assignment.progress = 100
  }
}

const viewAssignmentDetails = (assignmentId) => {
  // Navigate to assignment details or open modal
  console.log('View assignment details:', assignmentId)
}

const getStatusClass = (status) => {
  const classes = {
    'pending': 'bg-yellow-100 text-yellow-800',
    'in-progress': 'bg-blue-100 text-blue-800',
    'completed': 'bg-green-100 text-green-800',
    'overdue': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getPriorityClass = (priority) => {
  const classes = {
    'high': 'text-red-600 font-medium',
    'medium': 'text-yellow-600 font-medium',
    'low': 'text-green-600 font-medium'
  }
  return classes[priority] || 'text-gray-600'
}

const formatStatus = (status) => {
  return status.split('-').map(word => 
    word.charAt(0).toUpperCase() + word.slice(1)
  ).join(' ')
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

// Lifecycle
onMounted(() => {
  loadAssignments()
})
</script>