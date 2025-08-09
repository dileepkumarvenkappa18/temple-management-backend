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
      <!-- Status Alert -->
      <div v-if="statusMessage" :class="[
        'p-4 rounded-lg text-sm',
        isError ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
      ]">
        <div class="flex items-center justify-between">
          <span>{{ statusMessage }}</span>
          <button @click="statusMessage = ''" class="ml-2 text-gray-500 hover:text-gray-700">
            <svg class="w-4 h-4" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z" clip-rule="evenodd"/>
            </svg>
          </button>
        </div>
      </div>

      <!-- Debug Info (only in development) -->
      <div v-if="isDevelopment && showDebugInfo" class="p-4 bg-yellow-50 border border-yellow-200 rounded-lg text-sm">
        <div class="font-semibold text-yellow-800 mb-2">Debug Info:</div>
        <div class="space-y-1 text-yellow-700">
          <div>Sender ID: {{ senderID || 'Not set' }}</div>
          <div>Entity ID: {{ entityID || 'Not set' }}</div>
          <div>Auth Token: {{ hasAuthToken ? 'Present' : 'Missing' }}</div>
          <div>Recipient Type: {{ recipientType }}</div>
          <div>Selected Recipients: {{ selectedRecipients.length }}</div>
          <div>Delivery Methods: {{ deliveryMethods.join(', ') }}</div>
        </div>
        <button type="button" @click="showDebugInfo = false" class="text-yellow-600 underline mt-2">Hide Debug</button>
      </div>

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
            <button 
              v-if="isDevelopment"
              type="button"
              @click="showDebugInfo = !showDebugInfo"
              class="text-gray-500 hover:text-gray-700 font-medium"
            >
              {{ showDebugInfo ? 'Hide' : 'Show' }} Debug
            </button>
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

      <!-- Action Buttons -->
      <div class="flex flex-col sm:flex-row gap-3 pt-6 border-t border-gray-100">
        <button 
          type="submit"
          :disabled="!canSend || isSubmitting"
          class="flex-1 bg-indigo-600 hover:bg-indigo-700 disabled:bg-gray-300 disabled:cursor-not-allowed text-white font-semibold py-3 px-6 rounded-lg transition-all duration-200 flex items-center justify-center space-x-2"
        >
          <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
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
import { apiClient } from '@/plugins/axios'
import { useAuthStore } from '@/stores/auth'
import CommunicationService from '@/services/communication.service.js'
import { useRoute } from 'vue-router'

export default {
  name: 'MessageComposer',
  setup() {
    const authStore = useAuthStore()
    const route = useRoute()
    
    // Reactive state
    const isSubmitting = ref(false)
    const statusMessage = ref('')
    const isError = ref(false)
    const showDebugInfo = ref(false)
    
    // Computed values
    const senderID = computed(() => authStore.user?.id)
    const entityID = computed(() => 
      authStore.user?.entity_id ?? 
      authStore.user?.entityId ?? 
      authStore.user?.current_entity?.id ?? 
      localStorage.getItem('current_entity_id') ??
      null
    )
    
    const hasAuthToken = computed(() => !!localStorage.getItem('auth_token'))
    const isDevelopment = computed(() => import.meta.env.DEV)

    // Form data
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
      const isCustom = recipientType.value === 'custom'
      const hasRecipients = isCustom ? selectedRecipients.value.length > 0 : true

      return subject.value.trim() && 
             message.value.trim() && 
             deliveryMethods.value.length > 0 &&
             hasRecipients &&
             !isSubmitting.value
    })

    const showStatus = (msg, error = false) => {
      statusMessage.value = msg
      isError.value = error
      
      // Auto-clear after 5 seconds
      setTimeout(() => {
        statusMessage.value = ''
      }, 5000)
    }

    const updateRecipients = () => {
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
      console.log('Opening template selector...')
    }

    const saveDraft = () => {
      isDraft.value = true
      console.log('Saving draft...')
    }

    const clearDraft = () => {
      subject.value = ''
      message.value = ''
      isDraft.value = false
    }

    const sendMessage = async () => {
      console.log('📤 Starting message send process...')
      
      // Validation checks
      if (!senderID.value || !entityID.value) {
        const errorMsg = `Missing required data - Sender ID: ${senderID.value}, Entity ID: ${entityID.value}`
        console.error('❌ Validation failed:', errorMsg)
        showStatus(errorMsg, true)
        return
      }

      if (!hasAuthToken.value) {
        const errorMsg = 'No authentication token found'
        console.error('❌ Auth check failed:', errorMsg)
        showStatus(errorMsg, true)
        return
      }

      if (!subject.value.trim() || !message.value.trim()) {
        const errorMsg = 'Subject and message are required'
        console.error('❌ Content validation failed:', errorMsg)
        showStatus(errorMsg, true)
        return
      }

      if (!deliveryMethods.value.length) {
        const errorMsg = 'Please select at least one delivery method'
        console.error('❌ Delivery method validation failed:', errorMsg)
        showStatus(errorMsg, true)
        return
      }

      isSubmitting.value = true
      
      try {
        const isCustom = recipientType.value === 'custom'
        let successCount = 0
        let errorCount = 0
        let lastError = ''

        showStatus('Sending messages...', false)

        for (const channel of deliveryMethods.value) {
          console.log(`📡 Preparing to send via ${channel}...`)
          
          // Create payload matching expected backend format
          const payload = {
            channel: channel,
            subject: subject.value.trim(),
            body: message.value.trim(),
            template_id: null,
            // For custom recipients, send specific emails/phones
            // For audience-based, let backend handle the recipient list
            recipients: isCustom ? selectedRecipients.value.map(r => {
              // Assuming recipients have email/phone fields
              if (channel === 'email') return r.email
              if (channel === 'sms') return r.phone
              if (channel === 'whatsapp') return r.phone
              return r.id
            }).filter(Boolean) : [],
            // Send audience type if not custom
            audience: !isCustom ? recipientType.value : null,
            // Additional metadata that might be useful
            metadata: {
              sender_id: senderID.value,
              entity_id: entityID.value,
              message_type: messageType.value,
              scheduled_type: scheduleType.value
            }
          }

          console.log(`📋 Payload for ${channel}:`, JSON.stringify(payload, null, 2))

          try {
            const result = await CommunicationService.sendDirectNotification(payload)
            
            if (result.success) {
              successCount++
              console.log(`✅ Successfully sent via ${channel}`)
            } else {
              errorCount++
              lastError = result.error || 'Unknown error'
              console.error(`❌ Failed to send via ${channel}:`, result.error)
            }
          } catch (err) {
            errorCount++
            lastError = err.message || 'Network error'
            console.error(`💥 Exception sending via ${channel}:`, err)
          }
        }

        // Show final result
        if (successCount > 0 && errorCount === 0) {
          showStatus(`✅ Message sent successfully via ${successCount} channel(s)!`, false)
          // Clear form on complete success
          subject.value = ''
          message.value = ''
          selectedRecipients.value = []
          updateRecipients() // Reset recipients
        } else if (successCount > 0 && errorCount > 0) {
          showStatus(`⚠️ Partially sent: ${successCount} successful, ${errorCount} failed. Last error: ${lastError}`, true)
        } else {
          showStatus(`❌ Failed to send message: ${lastError}`, true)
        }

      } catch (error) {
        console.error('💥 Unexpected error in sendMessage:', error)
        showStatus(`❌ Unexpected error: ${error.message}`, true)
      } finally {
        isSubmitting.value = false
      }
    }

    onMounted(() => {
      console.log('🔧 MessageComposer mounted')
      console.log('Auth state:', {
        user: authStore.user,
        senderID: senderID.value,
        entityID: entityID.value,
        hasToken: hasAuthToken.value
      })
      
      updateRecipients()
      
      // Pre-fill if coming from a template
      const subjectFromQuery = route.query.subject
      const contentFromQuery = route.query.content

      if (subjectFromQuery) subject.value = subjectFromQuery
      if (contentFromQuery) message.value = contentFromQuery
    })

    return {
      // State
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
      isSubmitting,
      statusMessage,
      isError,
      showDebugInfo,
      
      // Computed
      canSend,
      senderID,
      entityID,
      hasAuthToken,
      isDevelopment,
      
      // Methods
      updateRecipients,
      removeRecipient,
      insertTemplate,
      saveDraft,
      clearDraft,
      sendMessage
    }
  }
}
</script>