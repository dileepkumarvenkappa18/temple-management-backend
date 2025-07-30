<template>
  <div class="group bg-white rounded-2xl shadow-lg hover:shadow-xl transition-all duration-300 overflow-hidden border border-gray-100 hover:border-indigo-200">
    <!-- Header Section -->
    <div class="bg-gradient-to-r from-indigo-50 to-indigo-100 px-6 py-4 border-b border-indigo-100">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-4">
          <!-- Avatar -->
          <div class="relative">
            <div class="w-12 h-12 bg-indigo-600 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-md">
              {{ getInitials(volunteer.fullName) }}
            </div>
            <div class="absolute -bottom-1 -right-1">
              <div :class="[
                'w-4 h-4 rounded-full border-2 border-white',
                volunteer.isActive ? 'bg-green-500' : 'bg-gray-400'
              ]"></div>
            </div>
          </div>
          
          <!-- Name and Role -->
          <div>
            <h3 class="font-bold text-gray-900 text-lg leading-tight">{{ volunteer.fullName }}</h3>
            <p class="text-indigo-600 font-medium text-sm">{{ volunteer.volunteerType || 'General Volunteer' }}</p>
          </div>
        </div>

        <!-- Actions Dropdown -->
        <div class="relative">
          <button 
            @click="toggleDropdown"
            class="p-2 hover:bg-indigo-200 rounded-lg transition-colors duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500"
            :class="{ 'bg-indigo-200': showDropdown }"
          >
            <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 5v.01M12 12v.01M12 19v.01"></path>
            </svg>
          </button>
          
          <!-- Dropdown Menu -->
          <div v-if="showDropdown" class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 z-10">
            <div class="py-2">
              <button 
                @click="$emit('view', volunteer)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors duration-200"
              >
                <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                </svg>
                View Details
              </button>
              <button 
                @click="$emit('edit', volunteer)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors duration-200"
              >
                <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                </svg>
                Edit Volunteer
              </button>
              <button 
                @click="$emit('assign-seva', volunteer)"
                class="w-full text-left px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors duration-200"
              >
                <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6"></path>
                </svg>
                Assign Seva
              </button>
              <hr class="my-1 border-gray-100">
              <button 
                @click="toggleStatus"
                :class="[
                  'w-full text-left px-4 py-2 text-sm transition-colors duration-200',
                  volunteer.isActive 
                    ? 'text-red-600 hover:bg-red-50' 
                    : 'text-green-600 hover:bg-green-50'
                ]"
              >
                <svg v-if="volunteer.isActive" class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                </svg>
                <svg v-else class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6.4-6.4A7.966 7.966 0 0112 2c4.418 0 8 3.582 8 8s-3.582 8-8 8-8-3.582-8-8c0-1.527.432-2.954 1.182-4.182"></path>
                </svg>
                {{ volunteer.isActive ? 'Deactivate' : 'Activate' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Content Section -->
    <div class="p-6 space-y-4">
      <!-- Contact Information -->
      <div class="space-y-3">
        <div class="flex items-center space-x-3 text-sm">
          <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
            <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
            </svg>
          </div>
          <span class="text-gray-700">{{ volunteer.email }}</span>
        </div>
        
        <div v-if="volunteer.phone" class="flex items-center space-x-3 text-sm">
          <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
            <svg class="w-4 h-4 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
            </svg>
          </div>
          <span class="text-gray-700">{{ volunteer.phone }}</span>
        </div>
      </div>

      <!-- Skills/Specializations -->
      <div v-if="volunteer.skills && volunteer.skills.length > 0" class="space-y-2">
        <h4 class="font-medium text-gray-900">Skills & Specializations</h4>
        <div class="flex flex-wrap gap-2">
          <span 
            v-for="skill in volunteer.skills" 
            :key="skill"
            class="px-3 py-1 bg-indigo-100 text-indigo-700 text-xs font-medium rounded-full"
          >
            {{ skill }}
          </span>
        </div>
      </div>

      <!-- Availability -->
      <div v-if="volunteer.availability" class="space-y-2">
        <h4 class="font-medium text-gray-900">Availability</h4>
        <div class="flex flex-wrap gap-2">
          <span 
            v-for="day in volunteer.availability" 
            :key="day"
            class="px-3 py-1 bg-green-100 text-green-700 text-xs font-medium rounded-full"
          >
            {{ day }}
          </span>
        </div>
      </div>

      <!-- Statistics Row -->
      <div class="grid grid-cols-3 gap-4 pt-4 border-t border-gray-100">
        <div class="text-center">
          <div class="text-2xl font-bold text-indigo-600">{{ volunteer.assignedSevas || 0 }}</div>
          <div class="text-xs text-gray-500 font-medium">Assigned Sevas</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-green-600">{{ volunteer.completedTasks || 0 }}</div>
          <div class="text-xs text-gray-500 font-medium">Completed</div>
        </div>
        <div class="text-center">
          <div class="text-2xl font-bold text-orange-600">{{ volunteer.hoursContributed || 0 }}</div>
          <div class="text-xs text-gray-500 font-medium">Hours</div>
        </div>
      </div>
    </div>

    <!-- Footer Actions -->
    <div class="px-6 py-4 bg-gray-50 border-t border-gray-100">
      <div class="flex items-center justify-between">
        <!-- Join Date -->
        <div class="text-sm text-gray-500">
          <span class="font-medium">Joined:</span> {{ formatDate(volunteer.joinedDate) }}
        </div>
        
        <!-- Quick Actions -->
        <div class="flex space-x-2">
          <button 
            @click="$emit('contact', volunteer)"
            class="px-4 py-2 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all duration-200 shadow-sm hover:shadow-md"
          >
            Contact
          </button>
          <button 
            @click="$emit('view-schedule', volunteer)"
            class="px-4 py-2 bg-white text-indigo-600 text-sm font-medium rounded-lg border border-indigo-200 hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 transition-all duration-200"
          >
            Schedule
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// Props
const props = defineProps({
  volunteer: {
    type: Object,
    required: true,
    default: () => ({
      id: '',
      fullName: '',
      email: '',
      phone: '',
      volunteerType: 'General Volunteer',
      isActive: true,
      skills: [],
      availability: [],
      assignedSevas: 0,
      completedTasks: 0,
      hoursContributed: 0,
      joinedDate: new Date().toISOString()
    })
  }
})

// Emits
const emit = defineEmits([
  'view', 
  'edit', 
  'assign-seva', 
  'contact', 
  'view-schedule', 
  'toggle-status'
])

// Reactive data
const showDropdown = ref(false)

// Methods
const getInitials = (name) => {
  if (!name) return 'V'
  return name
    .split(' ')
    .map(word => word.charAt(0))
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleDateString('en-US', { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric' 
  })
}

const toggleDropdown = () => {
  showDropdown.value = !showDropdown.value
}

const toggleStatus = () => {
  emit('toggle-status', {
    ...props.volunteer,
    isActive: !props.volunteer.isActive
  })
  showDropdown.value = false
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (!event.target.closest('.relative')) {
    showDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
/* Custom scrollbar for any overflow content */
.scrollbar-thin::-webkit-scrollbar {
  width: 4px;
}

.scrollbar-thin::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 2px;
}

.scrollbar-thin::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

.scrollbar-thin::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Smooth animations */
.group:hover .group-hover\:scale-105 {
  transform: scale(1.05);
}

/* Focus states for accessibility */
button:focus {
  outline: none;
}

/* Custom gradient animation */
.bg-gradient-to-r {
  background-size: 200% 200%;
  animation: gradient-shift 3s ease infinite;
}

@keyframes gradient-shift {
  0% {
    background-position: 0% 50%;
  }
  50% {
    background-position: 100% 50%;
  }
  100% {
    background-position: 0% 50%;
  }
}
</style>