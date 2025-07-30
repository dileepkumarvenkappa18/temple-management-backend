<template>
  <div class="bg-white rounded-xl shadow-lg p-6 transition-all duration-200">
    <!-- Calendar Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center space-x-4">
        <h2 class="text-2xl font-bold text-gray-900">Events Calendar</h2>
        <div class="flex items-center space-x-2">
          <span class="text-sm text-gray-600">{{ currentMonthYear }}</span>
        </div>
      </div>
      
      <!-- Calendar Navigation -->
      <div class="flex items-center space-x-3">
        <button
          @click="previousMonth"
          class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"></path>
          </svg>
        </button>
        
        <button
          @click="goToToday"
          class="px-4 py-2 text-sm font-medium text-indigo-600 bg-indigo-50 hover:bg-indigo-100 rounded-lg transition-colors duration-200"
        >
          Today
        </button>
        
        <button
          @click="nextMonth"
          class="p-2 text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"></path>
          </svg>
        </button>
      </div>
    </div>

    <!-- Calendar View Toggle -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center space-x-2">
        <button
          v-for="view in calendarViews"
          :key="view.value"
          @click="currentView = view.value"
          :class="[
            'px-4 py-2 text-sm font-medium rounded-lg transition-colors duration-200',
            currentView === view.value
              ? 'bg-indigo-600 text-white'
              : 'text-gray-600 hover:bg-gray-100'
          ]"
        >
          {{ view.label }}
        </button>
      </div>
      
      <!-- Event Legend -->
      <div class="flex items-center space-x-4">
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-indigo-500 rounded-full"></div>
          <span class="text-sm text-gray-600">Festival</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-emerald-500 rounded-full"></div>
          <span class="text-sm text-gray-600">Special Seva</span>
        </div>
        <div class="flex items-center space-x-2">
          <div class="w-3 h-3 bg-amber-500 rounded-full"></div>
          <span class="text-sm text-gray-600">Community Event</span>
        </div>
      </div>
    </div>

    <!-- Month View Calendar -->
    <div v-if="currentView === 'month'" class="grid grid-cols-7 gap-1">
      <!-- Day Headers -->
      <div
        v-for="day in dayHeaders"
        :key="day"
        class="p-3 text-center text-sm font-semibold text-gray-700 bg-gray-50 rounded-lg"
      >
        {{ day }}
      </div>
      
      <!-- Calendar Days -->
      <div
        v-for="day in calendarDays"
        :key="`${day.date}-${day.month}`"
        :class="[
          'min-h-[120px] p-2 border border-gray-100 rounded-lg transition-all duration-200 hover:shadow-md',
          day.isCurrentMonth ? 'bg-white' : 'bg-gray-50',
          day.isToday ? 'ring-2 ring-indigo-500 bg-indigo-50' : '',
          day.hasEvents ? 'cursor-pointer' : ''
        ]"
        @click="day.hasEvents && selectDay(day)"
      >
        <!-- Day Number -->
        <div class="flex items-start justify-between mb-2">
          <span
            :class="[
              'text-sm font-medium',
              day.isCurrentMonth ? 'text-gray-900' : 'text-gray-400',
              day.isToday ? 'text-indigo-600 font-bold' : ''
            ]"
          >
            {{ day.date }}
          </span>
          
          <!-- Events Count Badge -->
          <span
            v-if="day.events.length > 0"
            class="inline-flex items-center justify-center w-5 h-5 text-xs font-bold text-white bg-indigo-600 rounded-full"
          >
            {{ day.events.length }}
          </span>
        </div>
        
        <!-- Event Items -->
        <div class="space-y-1">
          <div
            v-for="(event, index) in day.events.slice(0, 2)"
            :key="event.id"
            :class="[
              'text-xs p-1 rounded truncate cursor-pointer transition-colors duration-200',
              getEventClass(event.type)
            ]"
            @click.stop="selectEvent(event)"
            :title="event.title"
          >
            {{ event.title }}
          </div>
          
          <!-- More Events Indicator -->
          <div
            v-if="day.events.length > 2"
            class="text-xs text-gray-500 cursor-pointer hover:text-gray-700"
            @click.stop="selectDay(day)"
          >
            +{{ day.events.length - 2 }} more
          </div>
        </div>
      </div>
    </div>

    <!-- Week View -->
    <div v-else-if="currentView === 'week'" class="space-y-4">
      <!-- Week Header -->
      <div class="grid grid-cols-8 gap-4">
        <div class="text-sm font-semibold text-gray-700"></div>
        <div
          v-for="day in weekDays"
          :key="day.date"
          class="text-center"
        >
          <div class="text-sm font-semibold text-gray-700">{{ day.dayName }}</div>
          <div
            :class="[
              'text-lg font-bold mt-1',
              day.isToday ? 'text-indigo-600' : 'text-gray-900'
            ]"
          >
            {{ day.date }}
          </div>
        </div>
      </div>
      
      <!-- Time Slots -->
      <div class="grid grid-cols-8 gap-4" style="grid-template-rows: repeat(12, minmax(60px, auto));">
        <!-- Time Labels -->
        <div v-for="hour in timeSlots" :key="hour" class="text-sm text-gray-600 text-right pr-2">
          {{ hour }}
        </div>
        
        <!-- Week Events -->
        <div
          v-for="day in weekDays"
          :key="`week-${day.date}`"
          class="border-l border-gray-200 relative"
        >
          <div
            v-for="event in day.events"
            :key="event.id"
            :class="[
              'absolute left-1 right-1 p-1 rounded text-xs cursor-pointer transition-all duration-200 hover:shadow-md',
              getEventClass(event.type)
            ]"
            :style="getEventPosition(event)"
            @click="selectEvent(event)"
          >
            <div class="font-medium truncate">{{ event.title }}</div>
            <div class="text-xs opacity-75">{{ event.time }}</div>
          </div>
        </div>
      </div>
    </div>

    <!-- Day View -->
    <div v-else-if="currentView === 'day'" class="space-y-4">
      <!-- Day Header -->
      <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
        <div>
          <h3 class="text-lg font-bold text-gray-900">{{ selectedDay.dayName }}</h3>
          <p class="text-sm text-gray-600">{{ selectedDay.fullDate }}</p>
        </div>
        <div class="text-sm text-gray-600">
          {{ selectedDay.events.length }} event{{ selectedDay.events.length !== 1 ? 's' : '' }}
        </div>
      </div>
      
      <!-- Day Events -->
      <div class="space-y-3">
        <div
          v-for="event in selectedDay.events"
          :key="event.id"
          :class="[
            'p-4 rounded-lg border-l-4 cursor-pointer transition-all duration-200 hover:shadow-md',
            getEventBorderClass(event.type)
          ]"
          @click="selectEvent(event)"
        >
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <h4 class="font-semibold text-gray-900">{{ event.title }}</h4>
              <p class="text-sm text-gray-600 mt-1">{{ event.description }}</p>
              <div class="flex items-center space-x-4 mt-2 text-sm text-gray-500">
                <span>{{ event.time }}</span>
                <span>{{ event.location }}</span>
                <span>{{ event.rsvpCount }} RSVPs</span>
              </div>
            </div>
            <div :class="['px-2 py-1 rounded-full text-xs font-medium', getEventClass(event.type)]">
              {{ event.type }}
            </div>
          </div>
        </div>
        
        <!-- No Events State -->
        <div
          v-if="selectedDay.events.length === 0"
          class="text-center py-12 text-gray-500"
        >
          <svg class="w-12 h-12 mx-auto mb-4 text-gray-300" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"></path>
          </svg>
          <p>No events scheduled for this day</p>
        </div>
      </div>
    </div>

    <!-- Event Details Modal -->
    <div
      v-if="selectedEvent"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50"
      @click="closeEventModal"
    >
      <div
        class="bg-white rounded-xl shadow-2xl max-w-lg w-full mx-4 transform transition-all duration-200"
        @click.stop
      >
        <div class="p-6">
          <div class="flex items-start justify-between mb-4">
            <div>
              <h3 class="text-xl font-bold text-gray-900">{{ selectedEvent.title }}</h3>
              <p class="text-sm text-gray-600 mt-1">{{ selectedEvent.date }} at {{ selectedEvent.time }}</p>
            </div>
            <button
              @click="closeEventModal"
              class="text-gray-400 hover:text-gray-600 transition-colors duration-200"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
          
          <div class="space-y-4">
            <div>
              <p class="text-gray-700">{{ selectedEvent.description }}</p>
            </div>
            
            <div class="flex items-center space-x-6 text-sm text-gray-600">
              <div class="flex items-center space-x-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"></path>
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"></path>
                </svg>
                <span>{{ selectedEvent.location }}</span>
              </div>
              
              <div class="flex items-center space-x-2">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z"></path>
                </svg>
                <span>{{ selectedEvent.rsvpCount }} RSVPs</span>
              </div>
            </div>
            
            <div class="flex items-center justify-between pt-4 border-t border-gray-200">
              <div :class="['px-3 py-1 rounded-full text-sm font-medium', getEventClass(selectedEvent.type)]">
                {{ selectedEvent.type }}
              </div>
              
              <div class="flex space-x-3">
                <button
                  @click="editEvent(selectedEvent)"
                  class="px-4 py-2 text-sm text-indigo-600 hover:bg-indigo-50 rounded-lg transition-colors duration-200"
                >
                  Edit Event
                </button>
                <button
                  @click="viewRSVPs(selectedEvent)"
                  class="px-4 py-2 text-sm bg-indigo-600 text-white hover:bg-indigo-700 rounded-lg transition-colors duration-200"
                >
                  View RSVPs
                </button>
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

// Props
const props = defineProps({
  events: {
    type: Array,
    default: () => []
  },
  entityId: {
    type: String,
    required: true
  }
})

// Emits
const emit = defineEmits(['edit-event', 'view-rsvps', 'select-day'])

// Reactive data
const currentDate = ref(new Date())
const currentView = ref('month')
const selectedEvent = ref(null)
const selectedDay = ref(null)

// Calendar views
const calendarViews = [
  { label: 'Month', value: 'month' },
  { label: 'Week', value: 'week' },
  { label: 'Day', value: 'day' }
]

// Day headers
const dayHeaders = ['Sun', 'Mon', 'Tue', 'Wed', 'Thu', 'Fri', 'Sat']

// Time slots for week/day view
const timeSlots = [
  '6:00 AM', '7:00 AM', '8:00 AM', '9:00 AM', '10:00 AM', '11:00 AM',
  '12:00 PM', '1:00 PM', '2:00 PM', '3:00 PM', '4:00 PM', '5:00 PM'
]

// Mock events data
const mockEvents = [
  {
    id: 1,
    title: 'Diwali Celebration',
    description: 'Grand Diwali festival with special poojas and cultural programs',
    date: '2024-11-01',
    time: '6:00 PM',
    location: 'Main Temple Hall',
    type: 'Festival',
    rsvpCount: 45
  },
  {
    id: 2,
    title: 'Rudra Abhishek',
    description: 'Monthly Rudra Abhishek seva',
    date: '2024-11-05',
    time: '7:00 AM',
    location: 'Shiva Temple',
    type: 'Special Seva',
    rsvpCount: 12
  },
  {
    id: 3,
    title: 'Community Feast',
    description: 'Free community meal for all devotees',
    date: '2024-11-10',
    time: '12:00 PM',
    location: 'Community Hall',
    type: 'Community Event',
    rsvpCount: 78
  }
]

// Computed properties
const currentMonthYear = computed(() => {
  return currentDate.value.toLocaleDateString('en-US', { 
    month: 'long', 
    year: 'numeric' 
  })
})

const calendarDays = computed(() => {
  const year = currentDate.value.getFullYear()
  const month = currentDate.value.getMonth()
  
  // Get first day of month and number of days
  const firstDay = new Date(year, month, 1)
  const lastDay = new Date(year, month + 1, 0)
  const daysInMonth = lastDay.getDate()
  const startingDayOfWeek = firstDay.getDay()
  
  const days = []
  const today = new Date()
  
  // Add previous month's trailing days
  const prevMonth = new Date(year, month - 1, 0)
  for (let i = startingDayOfWeek - 1; i >= 0; i--) {
    const date = prevMonth.getDate() - i
    days.push({
      date,
      month: month - 1,
      isCurrentMonth: false,
      isToday: false,
      events: [],
      hasEvents: false
    })
  }
  
  // Add current month's days
  for (let date = 1; date <= daysInMonth; date++) {
    const dayDate = new Date(year, month, date)
    const dayEvents = getEventsForDate(dayDate)
    
    days.push({
      date,
      month,
      isCurrentMonth: true,
      isToday: isSameDay(dayDate, today),
      events: dayEvents,
      hasEvents: dayEvents.length > 0,
      fullDate: dayDate
    })
  }
  
  // Add next month's leading days
  const remainingDays = 42 - days.length // 6 rows × 7 days
  for (let date = 1; date <= remainingDays; date++) {
    days.push({
      date,
      month: month + 1,
      isCurrentMonth: false,
      isToday: false,
      events: [],
      hasEvents: false
    })
  }
  
  return days
})

const weekDays = computed(() => {
  const days = []
  const startOfWeek = getStartOfWeek(currentDate.value)
  
  for (let i = 0; i < 7; i++) {
    const date = new Date(startOfWeek)
    date.setDate(startOfWeek.getDate() + i)
    
    days.push({
      date: date.getDate(),
      dayName: date.toLocaleDateString('en-US', { weekday: 'short' }),
      fullDate: date,
      isToday: isSameDay(date, new Date()),
      events: getEventsForDate(date)
    })
  }
  
  return days
})

// Methods
const getEventsForDate = (date) => {
  const dateString = date.toISOString().split('T')[0]
  return mockEvents.filter(event => event.date === dateString)
}

const isSameDay = (date1, date2) => {
  return date1.toDateString() === date2.toDateString()
}

const getStartOfWeek = (date) => {
  const start = new Date(date)
  const day = start.getDay()
  const diff = start.getDate() - day
  return new Date(start.setDate(diff))
}

const getEventClass = (type) => {
  const classes = {
    'Festival': 'bg-indigo-100 text-indigo-800 hover:bg-indigo-200',
    'Special Seva': 'bg-emerald-100 text-emerald-800 hover:bg-emerald-200',
    'Community Event': 'bg-amber-100 text-amber-800 hover:bg-amber-200'
  }
  return classes[type] || 'bg-gray-100 text-gray-800 hover:bg-gray-200'
}

const getEventBorderClass = (type) => {
  const classes = {
    'Festival': 'border-indigo-500 bg-indigo-50',
    'Special Seva': 'border-emerald-500 bg-emerald-50',
    'Community Event': 'border-amber-500 bg-amber-50'
  }
  return classes[type] || 'border-gray-500 bg-gray-50'
}

const getEventPosition = (event) => {
  // Simple positioning based on time
  const hour = parseInt(event.time.split(':')[0])
  const top = (hour - 6) * 60 // Assuming 6 AM start
  return {
    top: `${top}px`,
    height: '50px'
  }
}

const previousMonth = () => {
  const newDate = new Date(currentDate.value)
  newDate.setMonth(newDate.getMonth() - 1)
  currentDate.value = newDate
}

const nextMonth = () => {
  const newDate = new Date(currentDate.value)
  newDate.setMonth(newDate.getMonth() + 1)
  currentDate.value = newDate
}

const goToToday = () => {
  currentDate.value = new Date()
}

const selectDay = (day) => {
  selectedDay.value = {
    ...day,
    dayName: day.fullDate.toLocaleDateString('en-US', { weekday: 'long' }),
    fullDate: day.fullDate.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  }
  currentView.value = 'day'
  emit('select-day', day)
}

const selectEvent = (event) => {
  selectedEvent.value = event
}

const closeEventModal = () => {
  selectedEvent.value = null
}

const editEvent = (event) => {
  emit('edit-event', event)
  closeEventModal()
}

const viewRSVPs = (event) => {
  emit('view-rsvps', event)
  closeEventModal()
}

onMounted(() => {
  // Initialize with today if no specific date selected
  if (currentView.value === 'day' && !selectedDay.value) {
    const today = new Date()
    selectDay({
      date: today.getDate(),
      fullDate: today,
      events: getEventsForDate(today)
    })
  }
})
</script>