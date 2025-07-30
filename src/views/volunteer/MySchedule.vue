<template>
  <div class="min-h-screen bg-gray-50 p-4 sm:p-6 lg:p-8">
    <!-- Page Header -->
    <div class="mb-8">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h1 class="text-3xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">My Schedule</h1>
          <p class="mt-2 text-gray-600">View your upcoming seva assignments and events</p>
        </div>
        <div class="flex flex-col sm:flex-row gap-3">
          <select 
            v-model="selectedMonth" 
            @change="filterByMonth"
            class="px-4 py-2 border border-gray-300 rounded-xl bg-white text-gray-700 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          >
            <option v-for="month in months" :key="month.value" :value="month.value">
              {{ month.label }}
            </option>
          </select>
          <button
            @click="refreshSchedule"
            :disabled="isLoading"
            class="px-6 py-2 bg-indigo-600 text-white rounded-xl hover:bg-indigo-700 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 shadow-md"
          >
            <svg v-if="isLoading" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white inline" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            {{ isLoading ? 'Loading...' : 'Refresh' }}
          </button>
        </div>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
      <div class="bg-white rounded-2xl shadow-md p-6 border border-gray-200">
        <div class="flex items-center">
          <div class="p-3 bg-indigo-100 rounded-xl">
            <svg class="h-6 w-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">This Month</p>
            <p class="text-2xl font-bold text-gray-900">{{ stats.thisMonth }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-md p-6 border border-gray-200">
        <div class="flex items-center">
          <div class="p-3 bg-green-100 rounded-xl">
            <svg class="h-6 w-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Completed</p>
            <p class="text-2xl font-bold text-gray-900">{{ stats.completed }}</p>
          </div>
        </div>
      </div>

      <div class="bg-white rounded-2xl shadow-md p-6 border border-gray-200">
        <div class="flex items-center">
          <div class="p-3 bg-orange-100 rounded-xl">
            <svg class="h-6 w-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
            </svg>
          </div>
          <div class="ml-4">
            <p class="text-sm font-medium text-gray-600">Upcoming</p>
            <p class="text-2xl font-bold text-gray-900">{{ stats.upcoming }}</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Calendar View Toggle -->
    <div class="bg-white rounded-2xl shadow-md border border-gray-200 mb-8">
      <div class="p-6 border-b border-gray-200">
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
          <h2 class="text-xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">Schedule View</h2>
          <div class="flex bg-gray-100 rounded-xl p-1">
            <button
              @click="viewMode = 'list'"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                viewMode === 'list' 
                  ? 'bg-white text-indigo-600 shadow-sm' 
                  : 'text-gray-600 hover:text-gray-900'
              ]"
            >
              List View
            </button>
            <button
              @click="viewMode = 'calendar'"
              :class="[
                'px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200',
                viewMode === 'calendar' 
                  ? 'bg-white text-indigo-600 shadow-sm' 
                  : 'text-gray-600 hover:text-gray-900'
              ]"
            >
              Calendar View
            </button>
          </div>
        </div>
      </div>

      <!-- List View -->
      <div v-if="viewMode === 'list'" class="p-6">
        <div v-if="filteredSchedule.length === 0" class="text-center py-12">
          <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <h3 class="mt-4 text-lg font-medium text-gray-900">No scheduled items</h3>
          <p class="mt-2 text-gray-600">You don't have any assignments for the selected period.</p>
        </div>

        <div v-else class="space-y-4">
          <div
            v-for="item in filteredSchedule"
            :key="item.id"
            class="bg-gray-50 rounded-xl p-6 border border-gray-200 hover:shadow-md transition-all duration-200"
          >
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-2">
                  <span
                    :class="[
                      'px-2 py-1 rounded-lg text-xs font-medium',
                      item.type === 'seva' 
                        ? 'bg-indigo-100 text-indigo-800' 
                        : 'bg-purple-100 text-purple-800'
                    ]"
                  >
                    {{ item.type === 'seva' ? 'Seva' : 'Event' }}
                  </span>
                  <span
                    :class="[
                      'px-2 py-1 rounded-lg text-xs font-medium',
                      getStatusClass(item.status)
                    ]"
                  >
                    {{ item.status }}
                  </span>
                </div>
                <h3 class="text-lg font-semibold text-gray-900 mb-1">{{ item.title }}</h3>
                <p class="text-gray-600 mb-2">{{ item.description }}</p>
                <div class="flex flex-col sm:flex-row sm:items-center gap-2 text-sm text-gray-500">
                  <div class="flex items-center gap-1">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
                    </svg>
                    {{ formatDate(item.date) }}
                  </div>
                  <div class="flex items-center gap-1">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                    {{ item.time }}
                  </div>
                  <div class="flex items-center gap-1">
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    </svg>
                    {{ item.location }}
                  </div>
                </div>
              </div>
              <div class="flex flex-col sm:flex-row gap-2">
                <button
                  v-if="item.status === 'upcoming'"
                  @click="markAsCompleted(item.id)"
                  class="px-4 py-2 bg-green-600 text-white rounded-xl hover:bg-green-700 focus:ring-2 focus:ring-green-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
                >
                  Mark Complete
                </button>
                <button
                  @click="viewDetails(item)"
                  class="px-4 py-2 bg-indigo-100 text-indigo-700 rounded-xl hover:bg-indigo-200 focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 text-sm font-medium"
                >
                  View Details
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Calendar View -->
      <div v-else class="p-6">
        <div class="bg-gray-50 rounded-xl p-8 text-center">
          <svg class="mx-auto h-12 w-12 text-gray-400 mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Calendar View</h3>
          <p class="text-gray-600">Calendar view will be implemented in the next phase</p>
        </div>
      </div>
    </div>

    <!-- Details Modal -->
    <div
      v-if="showDetailsModal"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50"
      @click="closeDetailsModal"
    >
      <div
        class="bg-white rounded-2xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto"
        @click.stop
      >
        <div class="p-6 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h3 class="text-xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
              {{ selectedItem?.title }}
            </h3>
            <button
              @click="closeDetailsModal"
              class="p-2 hover:bg-gray-100 rounded-xl transition-all duration-200"
            >
              <svg class="h-5 w-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
        </div>
        <div class="p-6">
          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <p class="text-gray-900">{{ selectedItem?.description }}</p>
            </div>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Date</label>
                <p class="text-gray-900">{{ formatDate(selectedItem?.date) }}</p>
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-1">Time</label>
                <p class="text-gray-900">{{ selectedItem?.time }}</p>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Location</label>
              <p class="text-gray-900">{{ selectedItem?.location }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
              <span
                :class="[
                  'px-3 py-1 rounded-lg text-sm font-medium',
                  getStatusClass(selectedItem?.status)
                ]"
              >
                {{ selectedItem?.status }}
              </span>
            </div>
            <div v-if="selectedItem?.notes">
              <label class="block text-sm font-medium text-gray-700 mb-1">Notes</label>
              <p class="text-gray-900">{{ selectedItem.notes }}</p>
            </div>
          </div>
        </div>
        <div class="p-6 border-t border-gray-200 flex justify-end">
          <button
            @click="closeDetailsModal"
            class="px-6 py-2 bg-gray-100 text-gray-700 rounded-xl hover:bg-gray-200 focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-all duration-200 font-medium"
          >
            Close
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Reactive data
const isLoading = ref(false)
const selectedMonth = ref(new Date().getMonth())
const viewMode = ref('list')
const showDetailsModal = ref(false)
const selectedItem = ref(null)

// Mock schedule data
const schedule = ref([
  {
    id: 1,
    type: 'seva',
    title: 'Morning Prayer Service',
    description: 'Assist with morning prayers and temple setup',
    date: '2025-06-22',
    time: '06:00 AM - 08:00 AM',
    location: 'Main Temple Hall',
    status: 'upcoming',
    notes: 'Please arrive 15 minutes early for briefing'
  },
  {
    id: 2,
    type: 'event',
    title: 'Ganesha Chaturthi Celebration',
    description: 'Help with event setup and crowd management',
    date: '2025-06-25',
    time: '10:00 AM - 06:00 PM',
    location: 'Temple Grounds',
    status: 'upcoming',
    notes: 'Lunch will be provided'
  },
  {
    id: 3,
    type: 'seva',
    title: 'Temple Cleaning Drive',
    description: 'Weekly temple cleaning and maintenance',
    date: '2025-06-18',
    time: '07:00 AM - 09:00 AM',
    location: 'Entire Temple Complex',
    status: 'completed'
  },
  {
    id: 4,
    type: 'seva',
    title: 'Food Distribution',
    description: 'Distribute prasad to devotees',
    date: '2025-06-28',
    time: '12:00 PM - 02:00 PM',
    location: 'Temple Kitchen',
    status: 'upcoming'
  },
  {
    id: 5,
    type: 'event',
    title: 'Spiritual Discourse',
    description: 'Help with audio/video setup and seating arrangement',
    date: '2025-06-30',
    time: '05:00 PM - 08:00 PM',
    location: 'Discourse Hall',
    status: 'upcoming'
  }
])

// Months for filter
const months = [
  { value: 0, label: 'January' },
  { value: 1, label: 'February' },
  { value: 2, label: 'March' },
  { value: 3, label: 'April' },
  { value: 4, label: 'May' },
  { value: 5, label: 'June' },
  { value: 6, label: 'July' },
  { value: 7, label: 'August' },
  { value: 8, label: 'September' },
  { value: 9, label: 'October' },
  { value: 10, label: 'November' },
  { value: 11, label: 'December' }
]

// Computed properties
const filteredSchedule = computed(() => {
  return schedule.value.filter(item => {
    const itemDate = new Date(item.date)
    return itemDate.getMonth() === selectedMonth.value
  })
})

const stats = computed(() => {
  const thisMonth = filteredSchedule.value.length
  const completed = schedule.value.filter(item => item.status === 'completed').length
  const upcoming = schedule.value.filter(item => item.status === 'upcoming').length
  
  return { thisMonth, completed, upcoming }
})

// Methods
const refreshSchedule = async () => {
  isLoading.value = true
  // Simulate API call
  await new Promise(resolve => setTimeout(resolve, 1000))
  isLoading.value = false
}

const filterByMonth = () => {
  // Filter is already handled by computed property
}

const markAsCompleted = (id) => {
  const item = schedule.value.find(item => item.id === id)
  if (item) {
    item.status = 'completed'
  }
}

const viewDetails = (item) => {
  selectedItem.value = item
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedItem.value = null
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', {
    weekday: 'long',
    year: 'numeric',
    month: 'long',
    day: 'numeric'
  })
}

const getStatusClass = (status) => {
  switch (status) {
    case 'completed':
      return 'bg-green-100 text-green-800'
    case 'upcoming':
      return 'bg-orange-100 text-orange-800'
    case 'cancelled':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

// Lifecycle
onMounted(() => {
  // Initialize component
})
</script>