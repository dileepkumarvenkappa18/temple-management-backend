<template>
  <div class="bg-white rounded-xl shadow-md hover:shadow-lg transition-all duration-200 overflow-hidden group">
    <!-- Event Image -->
    <div class="relative h-48 bg-gradient-to-br from-indigo-500 to-indigo-600 overflow-hidden">
      <img 
        v-if="event.image" 
        :src="event.image" 
        :alt="event.title"
        class="w-full h-full object-cover group-hover:scale-105 transition-transform duration-300"
      />
      <div v-else class="flex items-center justify-center h-full">
        <Calendar class="h-16 w-16 text-white/70" />
      </div>
      
      <!-- Status Badge -->
      <div class="absolute top-3 right-3">
        <span 
          :class="statusClasses"
          class="px-3 py-1 rounded-full text-xs font-medium backdrop-blur-sm"
        >
          {{ event.status }}
        </span>
      </div>
      
      <!-- Date Badge -->
      <div class="absolute bottom-3 left-3 bg-white/90 backdrop-blur-sm rounded-lg px-3 py-2">
        <div class="text-xs text-gray-600 font-medium">{{ formatDate(event.date) }}</div>
        <div class="text-sm font-bold text-gray-900">{{ formatTime(event.date) }}</div>
      </div>
    </div>

    <!-- Event Content -->
    <div class="p-6">
      <!-- Title and Type -->
      <div class="flex items-start justify-between mb-3">
        <h3 class="text-lg font-bold text-gray-900 group-hover:text-indigo-600 transition-colors line-clamp-2">
          {{ event.title }}
        </h3>
        <span class="ml-2 px-2 py-1 bg-indigo-100 text-indigo-800 text-xs rounded-md font-medium whitespace-nowrap">
          {{ event.type }}
        </span>
      </div>

      <!-- Description -->
      <p class="text-gray-600 text-sm mb-4 line-clamp-3">
        {{ event.description }}
      </p>

      <!-- Event Stats -->
      <div class="grid grid-cols-2 gap-4 mb-4">
        <div class="flex items-center text-sm text-gray-500">
          <Users class="h-4 w-4 mr-2 text-indigo-500" />
          <span>{{ event.rsvpCount || 0 }} RSVPs</span>
        </div>
        <div class="flex items-center text-sm text-gray-500">
          <MapPin class="h-4 w-4 mr-2 text-indigo-500" />
          <span class="truncate">{{ event.location || 'Temple Premises' }}</span>
        </div>
      </div>

      <!-- Attached Sevas -->
      <div v-if="event.attachedSevas && event.attachedSevas.length > 0" class="mb-4">
        <h4 class="text-sm font-medium text-gray-700 mb-2 flex items-center">
          <Star class="h-4 w-4 mr-1 text-indigo-500" />
          Associated Sevas
        </h4>
        <div class="flex flex-wrap gap-2">
          <span 
            v-for="seva in event.attachedSevas.slice(0, 3)" 
            :key="seva.id"
            class="px-2 py-1 bg-amber-100 text-amber-800 text-xs rounded-md"
          >
            {{ seva.name }}
          </span>
          <span 
            v-if="event.attachedSevas.length > 3"
            class="px-2 py-1 bg-gray-100 text-gray-600 text-xs rounded-md"
          >
            +{{ event.attachedSevas.length - 3 }} more
          </span>
        </div>
      </div>

      <!-- Action Buttons -->
      <div class="flex items-center justify-between pt-4 border-t border-gray-100">
        <div class="flex space-x-2">
          <button
            @click="handleView"
            class="flex items-center px-3 py-2 text-sm font-medium text-indigo-600 hover:text-indigo-700 hover:bg-indigo-50 rounded-lg transition-colors"
          >
            <Eye class="h-4 w-4 mr-1" />
            View
          </button>
          
          <button
            v-if="canEdit"
            @click="handleEdit"
            class="flex items-center px-3 py-2 text-sm font-medium text-gray-600 hover:text-gray-700 hover:bg-gray-50 rounded-lg transition-colors"
          >
            <Edit3 class="h-4 w-4 mr-1" />
            Edit
          </button>
        </div>

        <!-- RSVP Status -->
        <div v-if="showRSVP" class="flex items-center">
          <button
            v-if="!userRSVPStatus"
            @click="handleRSVP"
            class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-sm font-medium rounded-lg transition-colors"
          >
            RSVP Now
          </button>
          <span
            v-else
            class="px-3 py-2 bg-green-100 text-green-800 text-sm font-medium rounded-lg flex items-center"
          >
            <Check class="h-4 w-4 mr-1" />
            Going
          </span>
        </div>
      </div>
    </div>

    <!-- Hover Overlay for Quick Actions -->
    <div class="absolute inset-0 bg-indigo-600/0 group-hover:bg-indigo-600/5 transition-colors pointer-events-none" />
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { Calendar, Users, MapPin, Star, Eye, Edit3, Check } from 'lucide-vue-next'

const props = defineProps({
  event: {
    type: Object,
    required: true
  },
  canEdit: {
    type: Boolean,
    default: false
  },
  showRSVP: {
    type: Boolean,
    default: false
  },
  userRSVPStatus: {
    type: String,
    default: null // 'going', 'maybe', 'not_going'
  }
})

const emit = defineEmits(['view', 'edit', 'rsvp'])

const statusClasses = computed(() => {
  const status = props.event.status?.toLowerCase()
  switch (status) {
    case 'upcoming':
      return 'bg-blue-100/80 text-blue-800'
    case 'ongoing':
      return 'bg-green-100/80 text-green-800'
    case 'completed':
      return 'bg-gray-100/80 text-gray-700'
    case 'cancelled':
      return 'bg-red-100/80 text-red-800'
    default:
      return 'bg-gray-100/80 text-gray-700'
  }
})

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', { 
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const formatTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('en-IN', { 
    hour: '2-digit',
    minute: '2-digit',
    hour12: true
  })
}

const handleView = () => {
  emit('view', props.event)
}

const handleEdit = () => {
  emit('edit', props.event)
}

const handleRSVP = () => {
  emit('rsvp', props.event)
}
</script>

<style scoped>
.line-clamp-2 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  line-clamp: 2; /* optional future support */
}

.line-clamp-3 {
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 3;
  overflow: hidden;
  line-clamp: 3; /* optional future support */
}
</style>
