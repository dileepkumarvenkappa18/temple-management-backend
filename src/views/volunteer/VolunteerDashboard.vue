<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Welcome Banner -->
    <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 text-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold font-roboto">Welcome back, {{ volunteer.fullName }}!</h1>
            <p class="mt-2 text-indigo-100 text-lg">{{ templeInfo.name }} - Volunteer Dashboard</p>
          </div>
          <div class="hidden md:flex items-center space-x-4">
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 text-center">
              <div class="text-2xl font-bold">{{ stats.assignedSevas }}</div>
              <div class="text-sm text-indigo-200">Assigned Sevas</div>
            </div>
            <div class="bg-white/10 backdrop-blur-sm rounded-xl p-4 text-center">
              <div class="text-2xl font-bold">{{ stats.upcomingEvents }}</div>
              <div class="text-sm text-indigo-200">Upcoming Events</div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Mobile Stats -->
      <div class="grid grid-cols-2 gap-4 md:hidden mb-6">
        <div class="bg-white rounded-xl shadow-md p-4 text-center border-l-4 border-indigo-500">
          <div class="text-2xl font-bold text-gray-900">{{ stats.assignedSevas }}</div>
          <div class="text-sm text-gray-600">Assigned Sevas</div>
        </div>
        <div class="bg-white rounded-xl shadow-md p-4 text-center border-l-4 border-indigo-500">
          <div class="text-2xl font-bold text-gray-900">{{ stats.upcomingEvents }}</div>
          <div class="text-sm text-gray-600">Upcoming Events</div>
        </div>
      </div>

      <!-- Quick Actions -->
      <div class="mb-8">
        <h2 class="text-xl font-bold text-gray-900 mb-4 font-roboto">Quick Actions</h2>
        <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
          <button 
            @click="viewAssignments"
            class="bg-white hover:bg-indigo-50 border border-gray-200 rounded-xl p-6 text-left transition-all duration-200 hover:shadow-lg hover:border-indigo-300"
          >
            <div class="flex items-center">
              <div class="bg-indigo-100 rounded-lg p-3">
                <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                </svg>
              </div>
              <div class="ml-4">
                <h3 class="font-semibold text-gray-900">My Assignments</h3>
                <p class="text-sm text-gray-600">View seva assignments</p>
              </div>
            </div>
          </button>

          <button 
            @click="viewSchedule"
            class="bg-white hover:bg-indigo-50 border border-gray-200 rounded-xl p-6 text-left transition-all duration-200 hover:shadow-lg hover:border-indigo-300"
          >
            <div class="flex items-center">
              <div class="bg-indigo-100 rounded-lg p-3">
                <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
              </div>
              <div class="ml-4">
                <h3 class="font-semibold text-gray-900">My Schedule</h3>
                <p class="text-sm text-gray-600">Check upcoming tasks</p>
              </div>
            </div>
          </button>

          <button 
            @click="viewNotifications"
            class="bg-white hover:bg-indigo-50 border border-gray-200 rounded-xl p-6 text-left transition-all duration-200 hover:shadow-lg hover:border-indigo-300 relative"
          >
            <div class="flex items-center">
              <div class="bg-indigo-100 rounded-lg p-3">
                <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z"/>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 19H6.5a2.5 2.5 0 010-5H9"/>
                </svg>
              </div>
              <div class="ml-4">
                <h3 class="font-semibold text-gray-900">Notifications</h3>
                <p class="text-sm text-gray-600">Temple updates</p>
              </div>
            </div>
            <span v-if="unreadNotifications > 0" class="absolute -top-2 -right-2 bg-red-500 text-white text-xs rounded-full w-6 h-6 flex items-center justify-center">
              {{ unreadNotifications }}
            </span>
          </button>

          <button 
            @click="viewTempleDetails"
            class="bg-white hover:bg-indigo-50 border border-gray-200 rounded-xl p-6 text-left transition-all duration-200 hover:shadow-lg hover:border-indigo-300"
          >
            <div class="flex items-center">
              <div class="bg-indigo-100 rounded-lg p-3">
                <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"/>
                </svg>
              </div>
              <div class="ml-4">
                <h3 class="font-semibold text-gray-900">Temple Info</h3>
                <p class="text-sm text-gray-600">View temple details</p>
              </div>
            </div>
          </button>
        </div>
      </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column - Assigned Sevas -->
        <div class="lg:col-span-2">
          <div class="bg-white rounded-xl shadow-md border border-gray-200">
            <div class="px-6 py-4 border-b border-gray-200">
              <h3 class="text-lg font-semibold text-gray-900 font-roboto">My Assigned Sevas</h3>
            </div>
            <div class="p-6">
              <div v-if="assignedSevas.length === 0" class="text-center py-8">
                <div class="text-gray-400 mb-4">
                  <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"/>
                  </svg>
                </div>
                <p class="text-gray-600">No sevas assigned yet</p>
                <p class="text-sm text-gray-500 mt-2">Check back later for new assignments</p>
              </div>
              <div v-else class="space-y-4">
                <div 
                  v-for="seva in assignedSevas" 
                  :key="seva.id"
                  class="border border-gray-200 rounded-lg p-4 hover:bg-gray-50 transition-colors duration-200"
                >
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <h4 class="font-semibold text-gray-900">{{ seva.name }}</h4>
                      <p class="text-sm text-gray-600 mt-1">{{ seva.description }}</p>
                      <div class="flex items-center mt-3 space-x-4 text-sm">
                        <div class="flex items-center text-gray-500">
                          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                          </svg>
                          {{ formatDate(seva.scheduledDate) }}
                        </div>
                        <div class="flex items-center text-gray-500">
                          <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                          </svg>
                          {{ seva.timeSlot }}
                        </div>
                      </div>
                    </div>
                    <div class="ml-4">
                      <span 
                        :class="getSevaStatusClass(seva.status)"
                        class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                      >
                        {{ seva.status }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Right Column - Upcoming Events & Notifications -->
        <div class="space-y-6">
          <!-- Upcoming Events -->
          <div class="bg-white rounded-xl shadow-md border border-gray-200">
            <div class="px-6 py-4 border-b border-gray-200">
              <h3 class="text-lg font-semibold text-gray-900 font-roboto">Upcoming Events</h3>
            </div>
            <div class="p-6">
              <div v-if="upcomingEvents.length === 0" class="text-center py-6">
                <div class="text-gray-400 mb-3">
                  <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                  </svg>
                </div>
                <p class="text-sm text-gray-600">No upcoming events</p>
              </div>
              <div v-else class="space-y-3">
                <div 
                  v-for="event in upcomingEvents" 
                  :key="event.id"
                  class="border-l-4 border-indigo-500 bg-indigo-50 p-4 rounded-r-lg"
                >
                  <h4 class="font-semibold text-gray-900 text-sm">{{ event.name }}</h4>
                  <div class="flex items-center mt-2 text-xs text-gray-600">
                    <svg class="w-3 h-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    {{ formatDate(event.date) }}
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Recent Notifications
          <div class="bg-white rounded-xl shadow-md border border-gray-200">
            <div class="px-6 py-4 border-b border-gray-200">
              <div class="flex items-center justify-between">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">Recent Notifications</h3>
                <span v-if="unreadNotifications > 0" class="bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded-full">
                  {{ unreadNotifications }} new
                </span>
              </div>
            </div>
            <div class="p-6">
              <div v-if="notifications.length === 0" class="text-center py-6">
                <div class="text-gray-400 mb-3">
                  <svg class="w-12 h-12 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z"/>
                  </svg>
                </div>
                <p class="text-sm text-gray-600">No notifications</p>
              </div>
              <div v-else class="space-y-3">
                <div 
                  v-for="notification in notifications.slice(0, 5)" 
                  :key="notification.id"
                  :class="[
                    'p-3 rounded-lg border',
                    notification.isRead ? 'bg-gray-50 border-gray-200' : 'bg-blue-50 border-blue-200'
                  ]"
                >
                  <p class="text-sm text-gray-900">{{ notification.message }}</p>
                  <p class="text-xs text-gray-500 mt-1">{{ formatTimeAgo(notification.createdAt) }}</p>
                </div>
              </div>
            </div>
          </div> -->

          <!-- Temple Information -->
          <!-- <div class="bg-white rounded-xl shadow-md border border-gray-200">
            <div class="px-6 py-4 border-b border-gray-200">
              <h3 class="text-lg font-semibold text-gray-900 font-roboto">Temple Information</h3>
            </div>
            <div class="p-6">
              <div class="space-y-3">
                <div>
                  <p class="text-sm font-medium text-gray-900">{{ templeInfo.name }}</p>
                  <p class="text-sm text-gray-600">{{ templeInfo.address }}</p>
                </div>
                <div>
                  <p class="text-sm text-gray-600">
                    <span class="font-medium">Contact:</span> {{ templeInfo.phone }}
                  </p>
                </div>
                <div>
                  <p class="text-sm text-gray-600">
                    <span class="font-medium">Email:</span> {{ templeInfo.email }}
                  </p>
                </div>
              </div>
            </div>
          </div> -->
        </div>
      </div>
    </div>

    <!-- Notifications Modal -->
    <div 
      v-if="showNotificationsModal" 
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeNotificationsModal"
    >
      <div 
        class="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden"
        @click.stop
      >
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-xl font-semibold text-gray-900 font-roboto">All Notifications</h3>
          <button 
            @click="closeNotificationsModal"
            class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        <div class="p-6 overflow-y-auto max-h-[60vh]">
          <div v-if="notifications.length === 0" class="text-center py-8">
            <div class="text-gray-400 mb-4">
              <svg class="w-16 h-16 mx-auto" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5v-5z"/>
              </svg>
            </div>
            <p class="text-gray-600">No notifications found</p>
          </div>
          <div v-else class="space-y-4">
            <div 
              v-for="notification in notifications" 
              :key="notification.id"
              :class="[
                'p-4 rounded-lg border transition-colors duration-200',
                notification.isRead ? 'bg-gray-50 border-gray-200' : 'bg-blue-50 border-blue-200'
              ]"
            >
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <p class="text-gray-900">{{ notification.message }}</p>
                  <p class="text-sm text-gray-500 mt-2">{{ formatFullDate(notification.createdAt) }}</p>
                </div>
                <div class="ml-4">
                  <span 
                    v-if="!notification.isRead"
                    class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800"
                  >
                    New
                  </span>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 bg-gray-50">
          <button 
            @click="markAllAsRead"
            class="w-full bg-indigo-600 text-white py-2 px-4 rounded-lg hover:bg-indigo-700 transition-colors duration-200 font-medium"
          >
            Mark All as Read
          </button>
        </div>
      </div>
    </div>

    <!-- Temple Details Modal -->
    <div 
      v-if="showTempleModal" 
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeTempleModal"
    >
      <div 
        class="bg-white rounded-xl shadow-xl max-w-2xl w-full mx-4 max-h-[80vh] overflow-hidden"
        @click.stop
      >
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-xl font-semibold text-gray-900 font-roboto">Temple Details</h3>
          <button 
            @click="closeTempleModal"
            class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"/>
            </svg>
          </button>
        </div>
        <div class="p-6 overflow-y-auto max-h-[60vh]">
          <div class="space-y-6">
            <!-- Basic Information -->
            <div>
              <h4 class="text-lg font-semibold text-gray-900 mb-3">Basic Information</h4>
              <div class="bg-gray-50 rounded-lg p-4 space-y-3">
                <div>
                  <label class="text-sm font-medium text-gray-700">Temple Name</label>
                  <p class="text-gray-900">{{ templeInfo.name }}</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-700">Address</label>
                  <p class="text-gray-900">{{ templeInfo.address }}</p>
                </div>
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium text-gray-700">Phone</label>
                    <p class="text-gray-900">{{ templeInfo.phone }}</p>
                  </div>
                  <div>
                    <label class="text-sm font-medium text-gray-700">Email</label>
                    <p class="text-gray-900">{{ templeInfo.email }}</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Temple Timings -->
            <div>
              <h4 class="text-lg font-semibold text-gray-900 mb-3">Temple Timings</h4>
              <div class="bg-gray-50 rounded-lg p-4">
                <div class="grid grid-cols-1 sm:grid-cols-2 gap-4">
                  <div>
                    <label class="text-sm font-medium text-gray-700">Morning</label>
                    <p class="text-gray-900">5:30 AM - 12:00 PM</p>
                  </div>
                  <div>
                    <label class="text-sm font-medium text-gray-700">Evening</label>
                    <p class="text-gray-900">4:00 PM - 9:00 PM</p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Daily Schedule -->
            <div>
              <h4 class="text-lg font-semibold text-gray-900 mb-3">Daily Schedule</h4>
              <div class="bg-gray-50 rounded-lg p-4 space-y-2">
                <div class="flex justify-between">
                  <span class="text-gray-700">Suprabhatam</span>
                  <span class="text-gray-900">5:30 AM</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-700">Abhishekam</span>
                  <span class="text-gray-900">6:00 AM</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-700">Archana</span>
                  <span class="text-gray-900">7:00 AM</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-700">Afternoon Aarti</span>
                  <span class="text-gray-900">12:00 PM</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-700">Evening Aarti</span>
                  <span class="text-gray-900">7:00 PM</span>
                </div>
                <div class="flex justify-between">
                  <span class="text-gray-700">Night Aarti</span>
                  <span class="text-gray-900">8:30 PM</span>
                </div>
              </div>
            </div>

            <!-- Management -->
            <div>
              <h4 class="text-lg font-semibold text-gray-900 mb-3">Temple Management</h4>
              <div class="bg-gray-50 rounded-lg p-4 space-y-3">
                <div>
                  <label class="text-sm font-medium text-gray-700">Head Priest</label>
                  <p class="text-gray-900">Sri Venkateswara Sharma</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-700">Temple Administrator</label>
                  <p class="text-gray-900">Sri Rajesh Kumar</p>
                </div>
                <div>
                  <label class="text-sm font-medium text-gray-700">Volunteer Coordinator</label>
                  <p class="text-gray-900">Smt. Lakshmi Devi</p>
                  <p class="text-sm text-gray-600">coordinator@srivenkateswara.temple</p>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'

// Router
const router = useRouter()

// Modal states
const showNotificationsModal = ref(false)
const showTempleModal = ref(false)

// Reactive data
const volunteer = ref({
  fullName: 'Ravi Kumar',
  email: 'ravi.volunteer@temple.com',
  phone: '+91 98765 43210'
})

const templeInfo = ref({
  name: 'Sri Venkateswara Temple',
  address: 'Temple Street, Bangalore, Karnataka 560001',
  phone: '+91 80 2345 6789',
  email: 'info@srivenkateswara.temple'
})

const assignedSevas = ref([
  {
    id: 1,
    name: 'Abhishekam',
    description: 'Assist with morning abhishekam ritual',
    scheduledDate: '2025-06-25',
    timeSlot: '6:00 AM - 8:00 AM',
    status: 'Scheduled'
  },
  {
    id: 2,
    name: 'Decoration',
    description: 'Flower decoration for evening aarti',
    scheduledDate: '2025-06-27',
    timeSlot: '5:00 PM - 6:00 PM',
    status: 'Pending'
  },
  {
    id: 3,
    name: 'Prasadam Distribution',
    description: 'Help distribute prasadam after evening prayers',
    scheduledDate: '2025-06-23',
    timeSlot: '7:00 PM - 8:00 PM',
    status: 'Completed'
  }
])

const upcomingEvents = ref([
  {
    id: 1,
    name: 'Janmashtami Celebration',
    date: '2025-08-16',
    description: 'Grand celebration of Lord Krishna\'s birthday'
  },
  {
    id: 2,
    name: 'Monthly Satsang',
    date: '2025-06-30',
    description: 'Community gathering and spiritual discourse'
  }
])

const notifications = ref([
  {
    id: 1,
    message: 'New seva assignment: Abhishekam on June 25th',
    isRead: false,
    createdAt: '2025-06-20T10:30:00Z'
  },
  {
    id: 2,
    message: 'Janmashtami celebration planning meeting on June 22nd',
    isRead: false,
    createdAt: '2025-06-19T14:15:00Z'
  },
  {
    id: 3,
    message: 'Thank you for your service during last weekend\'s event',
    isRead: true,
    createdAt: '2025-06-18T09:00:00Z'
  },
  {
    id: 4,
    message: 'Reminder: Monthly volunteer meeting tomorrow at 6 PM',
    isRead: false,
    createdAt: '2025-06-17T16:45:00Z'
  },
  {
    id: 5,
    message: 'New temple guidelines have been updated',
    isRead: true,
    createdAt: '2025-06-15T11:20:00Z'
  }
])

// Computed properties
const stats = computed(() => ({
  assignedSevas: assignedSevas.value.filter(seva => seva.status !== 'Completed').length,
  upcomingEvents: upcomingEvents.value.length
}))

const unreadNotifications = computed(() => 
  notifications.value.filter(n => !n.isRead).length
)

// Methods
const viewAssignments = () => {
  router.push('/entity/1/volunteer/assignments')
}

const viewSchedule = () => {
  router.push('/entity/1/volunteer/schedule')
}

const viewNotifications = () => {
  showNotificationsModal.value = true
}

const closeNotificationsModal = () => {
  showNotificationsModal.value = false
}

const viewTempleDetails = () => {
  showTempleModal.value = true
}

const closeTempleModal = () => {
  showTempleModal.value = false
}

const markAllAsRead = () => {
  notifications.value.forEach(n => n.isRead = true)
}

const getSevaStatusClass = (status) => {
  const classes = {
    'Scheduled': 'bg-blue-100 text-blue-800',
    'Pending': 'bg-yellow-100 text-yellow-800',
    'Completed': 'bg-green-100 text-green-800',
    'Cancelled': 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', {
    day: 'numeric',
    month: 'short',
    year: 'numeric'
  })
}

const formatFullDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatTimeAgo = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 1) return 'Yesterday'
  if (diffDays < 7) return `${diffDays} days ago`
  return formatDate(dateString)
}

// Lifecycle
onMounted(() => {
  // Could fetch volunteer data, assignments, etc.
  console.log('Volunteer dashboard mounted')
})
</script>

<style scoped>
.font-roboto {
  font-family: 'Roboto', sans-serif;
}
</style>