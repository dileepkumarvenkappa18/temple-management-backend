<template>
  <div class="bg-white rounded-xl shadow-lg">
    <!-- Form Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <h2 class="text-xl font-bold text-gray-900">
            {{ isEdit ? 'Edit Event' : 'Create New Event' }}
          </h2>
          <p class="text-sm text-gray-600 mt-1">
            {{ isEdit ? 'Update event details and notify devotees' : 'Plan a new temple event or festival' }}
          </p>
        </div>
        <div v-if="isEdit" class="flex items-center space-x-2">
          <span :class="statusBadgeClasses" class="px-3 py-1 rounded-full text-xs font-medium">
            {{ form.status }}
          </span>
        </div>
      </div>
    </div>

    <!-- Form Content -->
    <form @submit.prevent="handleSubmit" class="p-6 space-y-6">
      <!-- Basic Information -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <!-- Event Title -->
        <div class="lg:col-span-2">
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Event Title <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.title"
            type="text"
            placeholder="e.g., Maha Shivaratri Celebration"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            :class="{ 'border-red-300 focus:ring-red-500 focus:border-red-500': errors.title }"
          />
          <p v-if="errors.title" class="mt-1 text-sm text-red-600">{{ errors.title }}</p>
        </div>

        <!-- Event Type -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Event Type <span class="text-red-500">*</span>
          </label>
          <select
            v-model="form.type"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            :class="{ 'border-red-300 focus:ring-red-500 focus:border-red-500': errors.type }"
          >
            <option value="">Select Event Type</option>
            <option value="festival">Festival</option>
            <option value="ceremony">Ceremony</option>
            <option value="spiritual">Spiritual Discourse</option>
            <option value="cultural">Cultural Program</option>
            <option value="community">Community Gathering</option>
            <option value="special">Special Occasion</option>
          </select>
          <p v-if="errors.type" class="mt-1 text-sm text-red-600">{{ errors.type }}</p>
        </div>

        <!-- Event Status -->
        <div v-if="isEdit">
          <!-- <label class="block text-sm font-medium text-gray-700 mb-2">
            Status
          </label> -->
          <!-- <select
            v-model="form.status"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
          >
            <option value="upcoming">Upcoming</option>
            <option value="ongoing">Ongoing</option>
            <option value="completed">Completed</option>
            <option value="cancelled">Cancelled</option>
          </select> -->
        </div>
      </div>

      <!-- Date and Time -->
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Event Date <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.date"
            type="date"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            :class="{ 'border-red-300 focus:ring-red-500 focus:border-red-500': errors.date }"
          />
          <p v-if="errors.date" class="mt-1 text-sm text-red-600">{{ errors.date }}</p>
        </div>

        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">
            Event Time <span class="text-red-500">*</span>
          </label>
          <input
            v-model="form.time"
            type="time"
            class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
            :class="{ 'border-red-300 focus:ring-red-500 focus:border-red-500': errors.time }"
          />
          <p v-if="errors.time" class="mt-1 text-sm text-red-600">{{ errors.time }}</p>
        </div>
      </div>

      <!-- Location -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Location
        </label>
        <input
          v-model="form.location"
          type="text"
          placeholder="e.g., Main Temple Hall, Garden Area"
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
        />
        <p class="mt-1 text-sm text-gray-500">Leave empty to default to "Temple Premises"</p>
      </div>

      <!-- Description -->
      <div>
        <label class="block text-sm font-medium text-gray-700 mb-2">
          Event Description <span class="text-red-500">*</span>
        </label>
        <textarea
          v-model="form.description"
          rows="4"
          placeholder="Describe the event, its significance, schedule, and what devotees can expect..."
          class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
          :class="{ 'border-red-300 focus:ring-red-500 focus:border-red-500': errors.description }"
        ></textarea>
        <p v-if="errors.description" class="mt-1 text-sm text-red-600">{{ errors.description }}</p>
      </div>

      <!-- Form Actions -->
      <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200">
        <button
          type="button"
          @click="handleCancel"
          class="px-6 py-3 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-colors font-medium"
        >
          Cancel
        </button>
        <button
          type="submit"
          :disabled="isSubmitting"
          class="px-6 py-3 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors font-medium flex items-center"
        >
          <svg v-if="isSubmitting" xmlns="http://www.w3.org/2000/svg" class="animate-spin h-4 w-4 mr-2" viewBox="0 0 24 24">
            <path fill="none" stroke="currentColor" stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 22c5.523 0 10-4.477 10-10h-2c0 4.418-3.582 8-8 8s-8-3.582-8-8 3.582-8 8-8v2M18.59 13.51l2.12-2.12-2.12-2.12M5.41 8.59l-2.12 2.12 2.12 2.12" />
          </svg>
          {{ isEdit ? 'Update Event' : 'Create Event' }}
        </button>
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useEventStore } from '@/stores/event'
import { useToast } from '@/composables/useToast'

const props = defineProps({
  event: {
    type: Object,
    default: null
  },
  isEdit: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['submit', 'cancel'])

const eventStore = useEventStore()
const toast = useToast()

const isSubmitting = ref(false)
const errors = reactive({})

const form = reactive({
  title: '',
  type: '',
  status: 'upcoming',
  date: '',
  time: '',
  location: '',
  description: '',
  isActive: true
})

const statusBadgeClasses = computed(() => {
  const status = form.status?.toLowerCase()
  switch (status) {
    case 'upcoming':
      return 'bg-blue-100 text-blue-800'
    case 'ongoing':
      return 'bg-green-100 text-green-800'
    case 'completed':
      return 'bg-gray-100 text-gray-700'
    case 'cancelled':
      return 'bg-red-100 text-red-800'
    default:
      return 'bg-gray-100 text-gray-700'
  }
})

// Initialize form with event data if editing
onMounted(() => {
  if (props.event) {
    // Convert event data to the form structure
    let eventDate = '';
    let eventTime = '';
    
    // FIX 1: Better event type handling
    let eventType = '';
    if (props.event.event_type) {
      eventType = props.event.event_type;
    } else if (props.event.type) {
      eventType = props.event.type;
    } else if (props.event.eventType) {
      eventType = props.event.eventType;
    }
    
    // Make sure event type matches one of our dropdown options
    const validTypes = ['festival', 'ceremony', 'spiritual', 'cultural', 'community', 'special'];
    if (!validTypes.includes(eventType.toLowerCase())) {
      // If it's not a valid type, we'll set it to an empty string (which will show "Select Event Type")
      eventType = '';
    }
    
    // FIX 2: Better time handling
    if (props.event.eventDate) {
      try {
        const date = new Date(props.event.eventDate);
        if (!isNaN(date.getTime())) {
          eventDate = date.toISOString().split('T')[0];
          
          // Get hours and minutes directly from the date object
          const hours = String(date.getHours()).padStart(2, '0');
          const minutes = String(date.getMinutes()).padStart(2, '0');
          eventTime = `${hours}:${minutes}`;
        }
      } catch (e) {
        console.warn('Invalid eventDate format:', props.event.eventDate, e);
      }
    }

    // Additional fallbacks for date and time
    if (!eventDate && props.event.event_date) {
      eventDate = props.event.event_date;
    }

    // FIX 3: More specific time handling
    if (!eventTime && props.event.event_time) {
      // Try to handle event_time which might be a string or a Date object
      if (typeof props.event.event_time === 'string') {
        // If it's a string like "15:30", use it directly
        if (/^\d{2}:\d{2}$/.test(props.event.event_time)) {
          eventTime = props.event.event_time;
        } else {
          // If it's an ISO string or other date format
          try {
            const timeDate = new Date(props.event.event_time);
            if (!isNaN(timeDate.getTime())) {
              const hours = String(timeDate.getHours()).padStart(2, '0');
              const minutes = String(timeDate.getMinutes()).padStart(2, '0');
              eventTime = `${hours}:${minutes}`;
            }
          } catch (e) {
            console.warn('Could not parse event_time as date:', props.event.event_time);
          }
        }
      } else {
        // It might be a Date object directly
        try {
          const timeDate = new Date(props.event.event_time);
          if (!isNaN(timeDate.getTime())) {
            const hours = String(timeDate.getHours()).padStart(2, '0');
            const minutes = String(timeDate.getMinutes()).padStart(2, '0');
            eventTime = `${hours}:${minutes}`;
          }
        } catch (e) {
          console.warn('Invalid event_time object:', props.event.event_time);
        }
      }
    }
    
    // Final fallback for time
    if (!eventTime) {
      eventTime = '18:00'; // Default to 6 PM if no time found
    }
    
    const isActive = props.event.isActive !== undefined ? props.event.isActive : 
                     props.event.is_active !== undefined ? props.event.is_active : true;
    
    let status = props.event.status || 'upcoming';
    
    // If status is "cancelled", we need to update isActive
    if (status === 'cancelled') {
      status = 'upcoming'; // Default to upcoming, but mark as inactive
    }
    
    Object.assign(form, {
      title: props.event.title || '',
      type: eventType.toLowerCase(), // FIX: Ensure lowercase to match options
      status: status,
      date: eventDate,
      time: eventTime,
      location: props.event.location || '',
      description: props.event.description || '',
      isActive: isActive
    })
  } else {
    // Set default date to tomorrow for new events
    const tomorrow = new Date();
    tomorrow.setDate(tomorrow.getDate() + 1);
    form.date = tomorrow.toISOString().split('T')[0];
    form.time = '18:00'; // Default time 6 PM
  }
})

const validateForm = () => {
  // Reset errors
  Object.keys(errors).forEach(key => delete errors[key])
  
  let isValid = true

  if (!form.title.trim()) {
    errors.title = 'Event title is required'
    isValid = false
  }

  if (!form.type) {
    errors.type = 'Event type is required'
    isValid = false
  }

  if (!form.date) {
    errors.date = 'Event date is required'
    isValid = false
  } else {
    // Validate date is not in the past
    const selectedDate = new Date(form.date);
    const today = new Date();
    today.setHours(0, 0, 0, 0);
    
    if (selectedDate < today && !props.isEdit) {
      errors.date = 'Event date cannot be in the past'
      isValid = false
    }
  }

  if (!form.time) {
    errors.time = 'Event time is required'
    isValid = false
  }

  if (!form.description.trim()) {
    errors.description = 'Event description is required'
    isValid = false
  } else if (form.description.length < 10) {
    errors.description = 'Description is too short (minimum 10 characters)'
    isValid = false
  }

  return isValid
}

const handleSubmit = async () => {
  if (!validateForm()) {
    toast.error('Please fix the form errors before submitting')
    return
  }

  isSubmitting.value = true
  
  try {
    // Create proper request object that matches backend expectations
    const eventData = {
      title: form.title.trim(),
      description: form.description.trim(),
      event_type: form.type,
      event_date: form.date,
      event_time: form.time,
      location: form.location || 'Temple Premises',
      is_active: form.status !== 'cancelled' && form.isActive
    }
    
    let response;
    if (props.isEdit && props.event?.id) {
      response = await eventStore.updateEvent(props.event.id, eventData);
    } else {
      response = await eventStore.createEvent(eventData);
    }
    
    emit('submit', response)
  } catch (error) {
    console.error('Error submitting form:', error)
    toast.error('Failed to save event. Please try again.')
  } finally {
    isSubmitting.value = false
  }
}

const handleCancel = () => {
  emit('cancel')
}
</script>