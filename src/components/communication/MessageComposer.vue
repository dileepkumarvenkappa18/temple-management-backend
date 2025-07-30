<template>
  <div class="bg-white rounded-2xl shadow-lg border border-gray-100 overflow-hidden">
    <!-- Header -->
    <div class="bg-gradient-to-r from-indigo-600 to-indigo-700 px-6 py-4">
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div class="bg-white/20 rounded-lg p-2">
            <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
            </svg>
          </div>
          <div>
            <h2 class="text-xl font-bold text-white">Compose Message</h2>
            <p class="text-indigo-100 text-sm">Send notifications to devotees and volunteers</p>
          </div>
        </div>
        <button 
          v-if="isDraft"
          @click="clearDraft"
          class="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg text-sm font-medium transition-all duration-200"
        >
          Clear Draft
        </button>
      </div>
    </div>

    <form @submit.prevent="sendMessage" class="p-6 space-y-6">
      <!-- Recipients Section -->
      <div class="space-y-4">
        <div class="flex items-center justify-between">
          <label class="text-sm font-semibold text-gray-900">Recipients</label>
          <span class="text-xs text-gray-500">{{ selectedRecipients.length }} selected</span>
        </div>
        
        <!-- Recipient Type Selector -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-3">
          <label class="relative cursor-pointer">
            <input 
              type="radio" 
              value="all" 
              v-model="recipientType" 
              class="sr-only"
              @change="updateRecipients"
            >
            <div class="border-2 rounded-xl p-4 text-center transition-all duration-200" 
                 :class="recipientType === 'all' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-8 h-8 mx-auto mb-2 rounded-full flex items-center justify-center"
                   :class="recipientType === 'all' ? 'bg-indigo-500 text-white' : 'bg-gray-100 text-gray-400'">
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="recipientType === 'all' ? 'text-indigo-700' : 'text-gray-600'">
                All Members
              </span>
            </div>
          </label>

          <label class="relative cursor-pointer">
            <input 
              type="radio" 
              value="devotees" 
              v-model="recipientType" 
              class="sr-only"
              @change="updateRecipients"
            >
            <div class="border-2 rounded-xl p-4 text-center transition-all duration-200" 
                 :class="recipientType === 'devotees' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-8 h-8 mx-auto mb-2 rounded-full flex items-center justify-center"
                   :class="recipientType === 'devotees' ? 'bg-indigo-500 text-white' : 'bg-gray-100 text-gray-400'">
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="recipientType === 'devotees' ? 'text-indigo-700' : 'text-gray-600'">
                Devotees Only
              </span>
            </div>
          </label>

          <label class="relative cursor-pointer">
            <input 
              type="radio" 
              value="volunteers" 
              v-model="recipientType" 
              class="sr-only"
              @change="updateRecipients"
            >
            <div class="border-2 rounded-xl p-4 text-center transition-all duration-200" 
                 :class="recipientType === 'volunteers' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-8 h-8 mx-auto mb-2 rounded-full flex items-center justify-center"
                   :class="recipientType === 'volunteers' ? 'bg-indigo-500 text-white' : 'bg-gray-100 text-gray-400'">
                <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
                  <path d="M13 6a3 3 0 11-6 0 3 3 0 016 0zM18 8a2 2 0 11-4 0 2 2 0 014 0zM14 15a4 4 0 00-8 0v3h8v-3z"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="recipientType === 'volunteers' ? 'text-indigo-700' : 'text-gray-600'">
                Volunteers Only
              </span>
            </div>
          </label>
        </div>

        <!-- Custom Selection -->
        <div v-if="recipientType === 'custom'" class="space-y-3">
          <div class="flex items-center justify-between">
            <label class="text-sm font-medium text-gray-700">Select Recipients</label>
            <button 
              type="button"
              @click="showRecipientModal = true"
              class="text-indigo-600 hover:text-indigo-700 text-sm font-medium"
            >
              Browse & Select
            </button>
          </div>
          <div class="flex flex-wrap gap-2">
            <span 
              v-for="recipient in selectedRecipients" 
              :key="recipient.id"
              class="inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-indigo-100 text-indigo-800"
            >
              {{ recipient.name }}
              <button 
                type="button"
                @click="removeRecipient(recipient.id)"
                class="ml-2 text-indigo-600 hover:text-indigo-800"
              >
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
                </svg>
              </button>
            </span>
          </div>
        </div>
      </div>

      <!-- Message Type -->
      <div class="space-y-3">
        <label class="text-sm font-semibold text-gray-900">Message Type</label>
        <div class="grid grid-cols-1 sm:grid-cols-3 gap-3">
          <label class="relative cursor-pointer">
            <input type="radio" value="announcement" v-model="messageType" class="sr-only">
            <div class="border-2 rounded-lg p-3 flex items-center space-x-3 transition-all duration-200"
                 :class="messageType === 'announcement' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-5 h-5 rounded-full flex items-center justify-center"
                   :class="messageType === 'announcement' ? 'bg-indigo-500 text-white' : 'bg-gray-200'">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M18 10a8 8 0 11-16 0 8 8 0 0116 0zm-7-4a1 1 0 11-2 0 1 1 0 012 0zM9 9a1 1 0 000 2v3a1 1 0 001 1h1a1 1 0 100-2v-3a1 1 0 00-1-1H9z" clip-rule="evenodd"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="messageType === 'announcement' ? 'text-indigo-700' : 'text-gray-600'">
                Announcement
              </span>
            </div>
          </label>

          <label class="relative cursor-pointer">
            <input type="radio" value="reminder" v-model="messageType" class="sr-only">
            <div class="border-2 rounded-lg p-3 flex items-center space-x-3 transition-all duration-200"
                 :class="messageType === 'reminder' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-5 h-5 rounded-full flex items-center justify-center"
                   :class="messageType === 'reminder' ? 'bg-indigo-500 text-white' : 'bg-gray-200'">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm1-12a1 1 0 10-2 0v4a1 1 0 00.293.707l2.828 2.829a1 1 0 101.415-1.415L11 9.586V6z" clip-rule="evenodd"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="messageType === 'reminder' ? 'text-indigo-700' : 'text-gray-600'">
                Reminder
              </span>
            </div>
          </label>

          <label class="relative cursor-pointer">
            <input type="radio" value="urgent" v-model="messageType" class="sr-only">
            <div class="border-2 rounded-lg p-3 flex items-center space-x-3 transition-all duration-200"
                 :class="messageType === 'urgent' ? 'border-red-500 bg-red-50' : 'border-gray-200 hover:border-gray-300'">
              <div class="w-5 h-5 rounded-full flex items-center justify-center"
                   :class="messageType === 'urgent' ? 'bg-red-500 text-white' : 'bg-gray-200'">
                <svg class="w-3 h-3" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
                </svg>
              </div>
              <span class="font-medium text-sm" :class="messageType === 'urgent' ? 'text-red-700' : 'text-gray-600'">
                Urgent
              </span>
            </div>
          </label>
        </div>
      </div>

      <!-- Subject -->
      <div class="space-y-2">
        <label class="text-sm font-semibold text-gray-900">Subject</label>
        <input 
          type="text" 
          v-model="subject"
          placeholder="Enter message subject..."
          class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
          required
        >
      </div>

      <!-- Message Content -->
      <div class="space-y-2">
        <label class="text-sm font-semibold text-gray-900">Message</label>
        <textarea 
          v-model="message"
          placeholder="Type your message here..."
          rows="6"
          class="w-full px-4 py-3 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200 resize-none"
          required
        ></textarea>
        <div class="flex items-center justify-between text-sm">
          <div class="flex items-center space-x-4">
            <button 
              type="button"
              @click="insertTemplate"
              class="text-indigo-600 hover:text-indigo-700 font-medium"
            >
              Insert Template
            </button>
            <!-- <button 
              type="button"
              @click="saveDraft"
              class="text-gray-600 hover:text-gray-700 font-medium"
            >
              Save Draft
            </button> -->
          </div>
          <span class="text-gray-500">{{ message.length }}/1000 characters</span>
        </div>
      </div>

      <!-- Delivery Options -->
      <div class="space-y-3">
        <label class="text-sm font-semibold text-gray-900">Delivery Method</label>
        <div class="space-y-2">
          <label class="flex items-center space-x-3">
            <input 
              type="checkbox" 
              v-model="deliveryMethods" 
              value="sms"
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            >
            <span class="text-sm text-gray-700">SMS</span>
          </label>
          <label class="flex items-center space-x-3">
            <input 
              type="checkbox" 
              v-model="deliveryMethods" 
              value="email"
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            >
            <span class="text-sm text-gray-700">Email</span>
          </label>
          <label class="flex items-center space-x-3">
            <input 
              type="checkbox" 
              v-model="deliveryMethods" 
              value="whatsapp"
              class="w-4 h-4 text-indigo-600 border-gray-300 rounded focus:ring-indigo-500"
            >
            <span class="text-sm text-gray-700">WhatsApp</span>
          </label>
        </div>
      </div>

      <!-- Schedule Options -->
      <!-- <div class="space-y-3">
        <label class="text-sm font-semibold text-gray-900">Schedule</label>
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <label class="flex items-center space-x-3">
            <input 
              type="radio" 
              value="now" 
              v-model="scheduleType"
              class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
            >
            <span class="text-sm text-gray-700">Send Now</span>
          </label>
          <label class="flex items-center space-x-3">
            <input 
              type="radio" 
              value="later" 
              v-model="scheduleType"
              class="w-4 h-4 text-indigo-600 border-gray-300 focus:ring-indigo-500"
            >
            <span class="text-sm text-gray-700">Schedule for Later</span>
          </label>
        </div>
        
        <div v-if="scheduleType === 'later'" class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Date</label>
            <input 
              type="date" 
              v-model="scheduledDate"
              class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
          </div>
          <div>
            <label class="block text-xs font-medium text-gray-600 mb-1">Time</label>
            <input 
              type="time" 
              v-model="scheduledTime"
              class="w-full px-3 py-2 border border-gray-200 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            >
          </div>
        </div>
      </div> -->

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100">
        <button 
          type="button"
          @click="previewMessage"
          class="flex-1 bg-gray-100 hover:bg-gray-200 text-gray-700 font-semibold py-3 px-6 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"/>
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"/>
          </svg>
          <span>Preview</span>
        </button>
        
        <button 
          type="submit"
          :disabled="!canSend"
          class="flex-1 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 19l9 2-9-18-9 18 9-2zm0 0v-8"/>
          </svg>
          <span>{{ scheduleType === 'now' ? 'Send Message' : 'Schedule Message' }}</span>
        </button>
      </div>
    </form>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  name: 'MessageComposer',
  setup() {
    const recipientType = ref('all')
    const selectedRecipients = ref([])
    const messageType = ref('announcement')
    const subject = ref('')
    const message = ref('')
    const deliveryMethods = ref(['email'])
    const scheduleType = ref('now')
    const scheduledDate = ref('')
    const scheduledTime = ref('')
    const isDraft = ref(false)
    const showRecipientModal = ref(false)

    const canSend = computed(() => {
      return subject.value.trim() && 
             message.value.trim() && 
             deliveryMethods.value.length > 0 &&
             selectedRecipients.value.length > 0
    })

    const updateRecipients = () => {
      // Mock logic - in real app, this would fetch from API
      if (recipientType.value === 'all') {
        selectedRecipients.value = [
          { id: 1, name: 'All Devotees', type: 'devotee' },
          { id: 2, name: 'All Volunteers', type: 'volunteer' }
        ]
      } else if (recipientType.value === 'devotees') {
        selectedRecipients.value = [
          { id: 1, name: 'All Devotees', type: 'devotee' }
        ]
      } else if (recipientType.value === 'volunteers') {
        selectedRecipients.value = [
          { id: 2, name: 'All Volunteers', type: 'volunteer' }
        ]
      }
    }

    const removeRecipient = (id) => {
      selectedRecipients.value = selectedRecipients.value.filter(r => r.id !== id)
    }

    const insertTemplate = () => {
      // This would open a template selector modal
      console.log('Opening template selector...')
    }

    const saveDraft = () => {
      isDraft.value = true
      console.log('Saving draft...')
      // Save to localStorage or API
    }

    const clearDraft = () => {
      subject.value = ''
      message.value = ''
      isDraft.value = false
    }

    const previewMessage = () => {
      console.log('Opening preview modal...')
      // Show preview modal
    }

    const sendMessage = () => {
      console.log('Sending message...', {
        recipientType: recipientType.value,
        recipients: selectedRecipients.value,
        messageType: messageType.value,
        subject: subject.value,
        message: message.value,
        deliveryMethods: deliveryMethods.value,
        scheduleType: scheduleType.value,
        scheduledDate: scheduledDate.value,
        scheduledTime: scheduledTime.value
      })
    }

    onMounted(() => {
      updateRecipients()
    })

    return {
      recipientType,
      selectedRecipients,
      messageType,
      subject,
      message,
      deliveryMethods,
      scheduleType,
      scheduledDate,
      scheduledTime,
      isDraft,
      showRecipientModal,
      canSend,
      updateRecipients,
      removeRecipient,
      insertTemplate,
      saveDraft,
      clearDraft,
      previewMessage,
      sendMessage
    }
  }
}
</script>