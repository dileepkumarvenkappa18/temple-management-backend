<template>
  <div class="bg-white rounded-2xl shadow-lg p-6">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div class="flex items-center space-x-3">
        <div class="p-2 bg-indigo-100 rounded-lg">
          <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z"></path>
          </svg>
        </div>
        <div>
          <h2 class="text-xl font-semibold text-gray-900">Bulk Message</h2>
          <p class="text-sm text-gray-600">Send messages to multiple devotees at once</p>
        </div>
      </div>
      <div class="flex items-center space-x-2">
        <span class="text-sm text-gray-500">Selected: {{ selectedCount }}</span>
        <div class="w-2 h-2 bg-indigo-500 rounded-full"></div>
      </div>
    </div>

    <!-- Recipient Selection -->
    <div class="mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Select Recipients</h3>
      
      <!-- Filter Options -->
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-4">
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Group</label>
          <select v-model="selectedGroup" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200">
            <option value="">All Devotees</option>
            <option value="active">Active Devotees</option>
            <option value="volunteers">Volunteers</option>
            <!-- <option value="donors">Regular Donors</option>
            <option value="seva_participants">Seva Participants</option> -->
          </select>
        </div>
        
        <!-- <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Location</label>
          <select v-model="selectedLocation" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200">
            <option value="">All Locations</option>
            <option value="local">Local Devotees</option>
            <option value="outstation">Outstation Devotees</option>
            <option value="international">International Devotees</option>
          </select>
        </div> -->
        
        <!-- <div>
          <label class="block text-sm font-medium text-gray-700 mb-2">Filter by Age Group</label>
          <select v-model="selectedAgeGroup" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200">
            <option value="">All Ages</option>
            <option value="youth">Youth (18-35)</option>
            <option value="adult">Adult (36-60)</option>
            <option value="senior">Senior (60+)</option>
          </select>
        </div> -->
      </div>

      <!-- Quick Select Buttons -->
      <div class="flex flex-wrap gap-2 mb-4">
        <button 
          @click="selectAll"
          class="px-4 py-2 bg-indigo-100 text-indigo-700 rounded-lg hover:bg-indigo-200 transition-all duration-200 text-sm font-medium"
        >
          Select All
        </button>
        <button 
          @click="clearSelection"
          class="px-4 py-2 bg-gray-100 text-gray-700 rounded-lg hover:bg-gray-200 transition-all duration-200 text-sm font-medium"
        >
          Clear Selection
        </button>
        <button 
          @click="selectActiveOnly"
          class="px-4 py-2 bg-green-100 text-green-700 rounded-lg hover:bg-green-200 transition-all duration-200 text-sm font-medium"
        >
          Active Only
        </button>
      </div>

      <!-- Devotee List -->
      <div class="border border-gray-200 rounded-lg max-h-64 overflow-y-auto">
        <div class="p-3 bg-gray-50 border-b border-gray-200 flex items-center space-x-3">
          <input 
            type="checkbox" 
            v-model="selectAllChecked"
            @change="toggleSelectAll"
            class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
          >
          <span class="text-sm font-medium text-gray-700">Select All ({{ filteredDevotees.length }})</span>
        </div>
        
        <div class="divide-y divide-gray-200">
          <div 
            v-for="devotee in filteredDevotees" 
            :key="devotee.id"
            class="p-3 hover:bg-gray-50 transition-colors duration-150"
          >
            <label class="flex items-center space-x-3 cursor-pointer">
              <input 
                type="checkbox" 
                v-model="selectedDevotees"
                :value="devotee.id"
                class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
              >
              <div class="flex items-center space-x-3 flex-1">
                <div class="w-8 h-8 bg-indigo-100 rounded-full flex items-center justify-center">
                  <span class="text-sm font-medium text-indigo-600">{{ devotee.name.charAt(0) }}</span>
                </div>
                <div class="flex-1">
                  <p class="text-sm font-medium text-gray-900">{{ devotee.name }}</p>
                  <p class="text-xs text-gray-500">{{ devotee.phone }} • {{ devotee.location }}</p>
                </div>
                <div class="flex items-center space-x-2">
                  <span 
                    :class="[
                      'px-2 py-1 rounded-full text-xs font-medium',
                      devotee.status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'
                    ]"
                  >
                    {{ devotee.status }}
                  </span>
                </div>
              </div>
            </label>
          </div>
        </div>
      </div>
    </div>

    <!-- Message Composition -->
    <div class="mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Compose Message</h3>
      
      <!-- Message Type Selection -->
      <div class="flex space-x-4 mb-4">
        <label class="flex items-center space-x-2 cursor-pointer">
          <input 
            type="radio" 
            v-model="messageType" 
            value="sms"
            class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
          >
          <span class="text-sm font-medium text-gray-700">SMS</span>
        </label>
        <label class="flex items-center space-x-2 cursor-pointer">
          <input 
            type="radio" 
            v-model="messageType" 
            value="whatsapp"
            class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
          >
          <span class="text-sm font-medium text-gray-700">WhatsApp</span>
        </label>
        <label class="flex items-center space-x-2 cursor-pointer">
          <input 
            type="radio" 
            v-model="messageType" 
            value="email"
            class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
          >
          <span class="text-sm font-medium text-gray-700">Email</span>
        </label>
      </div>

      <!-- Template Selection -->
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Use Template (Optional)</label>
        <select v-model="selectedTemplate" @change="loadTemplate" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200">
          <option value="">Select a template...</option>
          <option value="festival_greeting">Festival Greeting</option>
          <option value="seva_reminder">Seva Reminder</option>
          <option value="donation_thank_you">Donation Thank You</option>
          <option value="event_announcement">Event Announcement</option>
          <option value="general_announcement">General Announcement</option>
        </select>
      </div>

      <!-- Subject (for Email) -->
      <div v-if="messageType === 'email'" class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Subject</label>
        <input 
          v-model="messageSubject" 
          type="text" 
          placeholder="Enter email subject..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
        >
      </div>

      <!-- Message Content -->
      <div class="mb-4">
        <label class="block text-sm font-medium text-gray-700 mb-2">Message Content</label>
        <textarea 
          v-model="messageContent" 
          rows="6" 
          placeholder="Type your message here..."
          class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 resize-none"
        ></textarea>
        <div class="flex justify-between items-center mt-2">
          <p class="text-xs text-gray-500">{{ messageContent.length }}/{{ maxLength }} characters</p>
          <div class="flex space-x-2">
            <!-- <button 
              @click="insertVariable('name')"
              class="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs hover:bg-gray-200 transition-colors duration-150"
            >
              Insert {name}
            </button> -->
            <!-- <button 
              @click="insertVariable('temple')"
              class="px-2 py-1 bg-gray-100 text-gray-700 rounded text-xs hover:bg-gray-200 transition-colors duration-150"
            >
              Insert {temple}
            </button> -->
          </div>
        </div>
      </div>

      <!-- Schedule Options -->
      <div class="mb-4">
        <!-- <label class="flex items-center space-x-2 cursor-pointer mb-3">
          <input 
            type="checkbox" 
            v-model="scheduleMessage"
            class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
          >
          <span class="text-sm font-medium text-gray-700">Schedule Message</span>
        </label> -->
        
        <div v-if="scheduleMessage" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date</label>
            <input 
              v-model="scheduleDate" 
              type="date" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Time</label>
            <input 
              v-model="scheduleTime" 
              type="time" 
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
          </div>
        </div>
      </div>
    </div>

    <!-- Preview Section -->
    <div v-if="messageContent" class="mb-6">
      <h3 class="text-lg font-medium text-gray-900 mb-4">Message Preview</h3>
      <div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
        <div class="flex items-center space-x-2 mb-3">
          <div class="w-2 h-2 bg-indigo-500 rounded-full"></div>
          <span class="text-sm font-medium text-gray-700">{{ messageType.toUpperCase() }} Message</span>
        </div>
        <div v-if="messageType === 'email'" class="mb-2">
          <p class="text-sm font-medium text-gray-900">Subject: {{ messageSubject || 'No Subject' }}</p>
        </div>
        <p class="text-sm text-gray-800 whitespace-pre-wrap">{{ previewMessage }}</p>
      </div>
    </div>

    <!-- Action Buttons -->
    <div class="flex flex-col sm:flex-row gap-3 sm:justify-end">
      <!-- <button 
        @click="saveDraft"
        class="px-6 py-2 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-all duration-200 font-medium"
      >
        Save as Draft
      </button> -->
      <!-- <button 
        @click="sendPreview"
        class="px-6 py-2 bg-indigo-100 text-indigo-700 rounded-lg hover:bg-indigo-200 transition-all duration-200 font-medium"
      >
        Send Preview
      </button> -->
      <button 
        @click="sendMessage"
        :disabled="!canSendMessage"
        :class="[
          'px-6 py-2 rounded-lg font-medium transition-all duration-200',
          canSendMessage 
            ? 'bg-indigo-600 text-white hover:bg-indigo-700 shadow-md hover:shadow-lg' 
            : 'bg-gray-300 text-gray-500 cursor-not-allowed'
        ]"
      >
        {{ scheduleMessage ? 'Schedule Message' : 'Send Message' }}
      </button>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  name: 'BulkMessage',
  setup() {
    // Reactive data
    const selectedGroup = ref('')
    const selectedLocation = ref('')
    const selectedAgeGroup = ref('')
    const selectedDevotees = ref([])
    const selectAllChecked = ref(false)
    const messageType = ref('sms')
    const selectedTemplate = ref('')
    const messageSubject = ref('')
    const messageContent = ref('')
    const scheduleMessage = ref(false)
    const scheduleDate = ref('')
    const scheduleTime = ref('')

    // Mock data
    const devotees = ref([
      { id: 1, name: 'Rajesh Kumar', phone: '+91 98765 43210', location: 'Bangalore', status: 'active', group: 'active', ageGroup: 'adult' },
      { id: 2, name: 'Priya Sharma', phone: '+91 98765 43211', location: 'Mumbai', status: 'active', group: 'volunteers', ageGroup: 'youth' },
      { id: 3, name: 'Amit Patel', phone: '+91 98765 43212', location: 'Delhi', status: 'inactive', group: 'donors', ageGroup: 'senior' },
      { id: 4, name: 'Sunita Devi', phone: '+91 98765 43213', location: 'Chennai', status: 'active', group: 'seva_participants', ageGroup: 'adult' },
      { id: 5, name: 'Vikram Singh', phone: '+91 98765 43214', location: 'Pune', status: 'active', group: 'active', ageGroup: 'youth' },
    ])

    const templates = ref({
      'festival_greeting': {
        subject: 'Festival Greetings from {temple}',
        content: 'Dear {name},\n\nWarm festival greetings from {temple}! May this auspicious occasion bring joy, peace, and prosperity to you and your family.\n\nWith divine blessings,\n{temple} Team'
      },
      'seva_reminder': {
        subject: 'Seva Reminder - {temple}',
        content: 'Dear {name},\n\nThis is a gentle reminder about your upcoming seva at {temple}. Please arrive 15 minutes early for preparation.\n\nThank you for your devotion.\n\nWith regards,\n{temple} Team'
      },
      'donation_thank_you': {
        subject: 'Thank You for Your Donation - {temple}',
        content: 'Dear {name},\n\nThank you for your generous donation to {temple}. Your contribution helps us serve the community better.\n\nMay the divine bless you abundantly.\n\nWith gratitude,\n{temple} Team'
      }
    })

    // Computed properties
    const filteredDevotees = computed(() => {
      return devotees.value.filter(devotee => {
        const groupMatch = !selectedGroup.value || devotee.group === selectedGroup.value
        const locationMatch = !selectedLocation.value || devotee.location.toLowerCase().includes(selectedLocation.value.toLowerCase())
        const ageMatch = !selectedAgeGroup.value || devotee.ageGroup === selectedAgeGroup.value
        return groupMatch && locationMatch && ageMatch
      })
    })

    const selectedCount = computed(() => selectedDevotees.value.length)

    const maxLength = computed(() => {
      return messageType.value === 'sms' ? 160 : messageType.value === 'whatsapp' ? 4096 : 10000
    })

    const previewMessage = computed(() => {
      return messageContent.value
        .replace(/{name}/g, 'Devotee Name')
        .replace(/{temple}/g, 'Temple Name')
    })

    const canSendMessage = computed(() => {
      return selectedDevotees.value.length > 0 && 
             messageContent.value.trim() !== '' && 
             (!scheduleMessage.value || (scheduleDate.value && scheduleTime.value))
    })

    // Methods
    const selectAll = () => {
      selectedDevotees.value = filteredDevotees.value.map(d => d.id)
      selectAllChecked.value = true
    }

    const clearSelection = () => {
      selectedDevotees.value = []
      selectAllChecked.value = false
    }

    const selectActiveOnly = () => {
      selectedDevotees.value = filteredDevotees.value
        .filter(d => d.status === 'active')
        .map(d => d.id)
    }

    const toggleSelectAll = () => {
      if (selectAllChecked.value) {
        selectAll()
      } else {
        clearSelection()
      }
    }

    const loadTemplate = () => {
      if (selectedTemplate.value && templates.value[selectedTemplate.value]) {
        const template = templates.value[selectedTemplate.value]
        if (messageType.value === 'email') {
          messageSubject.value = template.subject
        }
        messageContent.value = template.content
      }
    }

    const insertVariable = (variable) => {
      const textarea = document.querySelector('textarea')
      const start = textarea.selectionStart
      const end = textarea.selectionEnd
      const text = messageContent.value
      messageContent.value = text.substring(0, start) + `{${variable}}` + text.substring(end)
    }

    const saveDraft = () => {
      // Save draft logic
      console.log('Saving draft...')
      // Show success toast
    }

    const sendPreview = () => {
      // Send preview logic
      console.log('Sending preview...')
      // Show success toast
    }

    const sendMessage = () => {
      // Send message logic
      console.log('Sending message to:', selectedDevotees.value)
      console.log('Message:', messageContent.value)
      console.log('Type:', messageType.value)
      if (scheduleMessage.value) {
        console.log('Scheduled for:', scheduleDate.value, scheduleTime.value)
      }
      // Show success toast
    }

    return {
      selectedGroup,
      selectedLocation,
      selectedAgeGroup,
      selectedDevotees,
      selectAllChecked,
      messageType,
      selectedTemplate,
      messageSubject,
      messageContent,
      scheduleMessage,
      scheduleDate,
      scheduleTime,
      filteredDevotees,
      selectedCount,
      maxLength,
      previewMessage,
      canSendMessage,
      selectAll,
      clearSelection,
      selectActiveOnly,
      toggleSelectAll,
      loadTemplate,
      insertVariable,
      saveDraft,
      sendPreview,
      sendMessage
    }
  }
}
</script>