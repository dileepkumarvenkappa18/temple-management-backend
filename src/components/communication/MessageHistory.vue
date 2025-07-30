<template>
  <div class="bg-white rounded-xl shadow-md">
    <!-- Header -->
    <div class="px-6 py-4 border-b border-gray-200">
      <div class="flex items-center justify-between">
        <div>
          <h3 class="text-lg font-semibold text-gray-900">Message History</h3>
          <p class="text-sm text-gray-600">View all sent messages and their delivery status</p>
        </div>
        <div class="flex items-center space-x-3">
          <!-- Export Button -->
          <button 
            @click="exportHistory"
            class="px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-colors duration-200 flex items-center space-x-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
            </svg>
            <span>Export</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Filters -->
    <div class="px-6 py-4 border-b border-gray-200 bg-gray-50">
      <div class="grid grid-cols-1 md:grid-cols-4 gap-4">
        <!-- Date Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
          <select v-model="filters.dateRange" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
            <option value="all">All Time</option>
            <option value="today">Today</option>
            <option value="week">This Week</option>
            <option value="month">This Month</option>
            <option value="custom">Custom Range</option>
          </select>
        </div>

        <!-- Message Type -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Message Type</label>
          <select v-model="filters.messageType" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
            <option value="all">All Types</option>
            <option value="bulk">Bulk Message</option>
            <option value="event">Event Notification</option>
            <option value="seva">Seva Reminder</option>
            <option value="donation">Donation Receipt</option>
            <option value="announcement">Announcement</option>
          </select>
        </div>

        <!-- Delivery Status -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Delivery Status</label>
          <select v-model="filters.status" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
            <option value="all">All Status</option>
            <option value="sent">Sent</option>
            <option value="delivered">Delivered</option>
            <option value="failed">Failed</option>
            <option value="pending">Pending</option>
          </select>
        </div>

        <!-- Channel -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Channel</label>
          <select v-model="filters.channel" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
            <option value="all">All Channels</option>
            <option value="email">Email</option>
            <option value="sms">SMS</option>
            <option value="whatsapp">WhatsApp</option>
            <option value="push">Push Notification</option>
          </select>
        </div>
      </div>

      <!-- Search -->
      <div class="mt-4">
        <div class="relative">
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search by message content, recipient name, or subject..."
            class="w-full pl-10 pr-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
          >
          <svg class="absolute left-3 top-2.5 h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
          </svg>
        </div>
      </div>
    </div>

    <!-- Message History Table -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Message</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Recipients</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Channel</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Sent Date</th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="message in filteredMessages" :key="message.id" class="hover:bg-gray-50 transition-colors duration-150">
            <!-- Message Content -->
            <td class="px-6 py-4">
              <div class="max-w-xs">
                <p class="text-sm font-medium text-gray-900 truncate">{{ message.subject || message.title }}</p>
                <p class="text-sm text-gray-600 truncate">{{ message.content }}</p>
                <div class="flex items-center mt-1 space-x-2">
                  <span :class="getMessageTypeClass(message.type)" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium">
                    {{ getMessageTypeLabel(message.type) }}
                  </span>
                </div>
              </div>
            </td>

            <!-- Recipients -->
            <td class="px-6 py-4">
              <div class="text-sm text-gray-900">
                <div class="flex items-center space-x-2">
                  <span class="font-medium">{{ message.recipientCount }}</span>
                  <span class="text-gray-600">recipients</span>
                </div>
                <div v-if="message.recipientGroups" class="text-xs text-gray-500 mt-1">
                  {{ message.recipientGroups.join(', ') }}
                </div>
              </div>
            </td>

            <!-- Channel -->
            <td class="px-6 py-4">
              <div class="flex items-center space-x-2">
                <div :class="getChannelIconClass(message.channel)" class="w-8 h-8 rounded-full flex items-center justify-center">
                  <component :is="getChannelIcon(message.channel)" class="w-4 h-4" />
                </div>
                <span class="text-sm text-gray-900 capitalize">{{ message.channel }}</span>
              </div>
            </td>

            <!-- Status -->
            <td class="px-6 py-4">
              <div class="space-y-1">
                <span :class="getStatusClass(message.deliveryStatus)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                  {{ getStatusLabel(message.deliveryStatus) }}
                </span>
                <div class="text-xs text-gray-500">
                  {{ message.deliveredCount }}/{{ message.recipientCount }} delivered
                </div>
              </div>
            </td>

            <!-- Sent Date -->
            <td class="px-6 py-4 text-sm text-gray-900">
              <div>{{ formatDate(message.sentAt) }}</div>
              <div class="text-xs text-gray-500">{{ formatTime(message.sentAt) }}</div>
            </td>

            <!-- Actions -->
            <td class="px-6 py-4">
              <div class="flex items-center space-x-2">
                <!-- View Details -->
                <button
                  @click="viewMessageDetails(message)"
                  class="text-indigo-600 hover:text-indigo-800 transition-colors duration-200"
                  title="View Details"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                </button>

                <!-- Resend (if failed) -->
                <button
                  v-if="message.deliveryStatus === 'failed'"
                  @click="resendMessage(message)"
                  class="text-green-600 hover:text-green-800 transition-colors duration-200"
                  title="Resend Message"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                  </svg>
                </button>

                <!-- Download Report -->
                <button
                  @click="downloadMessageReport(message)"
                  class="text-gray-600 hover:text-gray-800 transition-colors duration-200"
                  title="Download Report"
                >
                  <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="filteredMessages.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-3.582 8-8 8a8.959 8.959 0 01-4.906-1.456A3.975 3.975 0 000 16c0-1.105.895-2 2-2h1.09a8.969 8.969 0 000-4H2c-1.105 0-2-.895-2-2 0-1.105.895-2 2-2h1.09a8.969 8.969 0 000-4H2c-1.105 0-2-.895-2-2s.895-2 2-2h1.09C5.456 2.59 8.646 0 12.01 0 17.424 0 21.818 3.582 21.818 8c0 1.105-.895 2-2 2h-1.09a8.969 8.969 0 000 4h1.09c1.105 0 2 .895 2 2z"></path>
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No messages found</h3>
      <p class="mt-1 text-sm text-gray-500">No messages match your current filters.</p>
    </div>

    <!-- Pagination -->
    <div v-if="filteredMessages.length > 0" class="px-6 py-4 border-t border-gray-200 bg-gray-50">
      <div class="flex items-center justify-between">
        <div class="text-sm text-gray-700">
          Showing {{ (currentPage - 1) * itemsPerPage + 1 }} to {{ Math.min(currentPage * itemsPerPage, filteredMessages.length) }} of {{ filteredMessages.length }} messages
        </div>
        <div class="flex items-center space-x-2">
          <button
            @click="currentPage = Math.max(1, currentPage - 1)"
            :disabled="currentPage === 1"
            class="px-3 py-1 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span class="text-sm text-gray-700">{{ currentPage }} of {{ totalPages }}</span>
          <button
            @click="currentPage = Math.min(totalPages, currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="px-3 py-1 border border-gray-300 rounded-md text-sm text-gray-700 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>

  <!-- Message Details Modal -->
  <div v-if="selectedMessage" class="fixed inset-0 z-50 overflow-y-auto" aria-labelledby="modal-title" role="dialog" aria-modal="true">
    <div class="flex items-end justify-center min-h-screen pt-4 px-4 pb-20 text-center sm:block sm:p-0">
      <div class="fixed inset-0 bg-gray-500 bg-opacity-75 transition-opacity" @click="closeMessageDetails"></div>
      <div class="inline-block align-bottom bg-white rounded-lg text-left overflow-hidden shadow-xl transform transition-all sm:my-8 sm:align-middle sm:max-w-4xl sm:w-full">
        <div class="bg-white px-6 py-4">
          <div class="flex items-center justify-between border-b border-gray-200 pb-4">
            <h3 class="text-lg font-semibold text-gray-900">Message Details</h3>
            <button @click="closeMessageDetails" class="text-gray-400 hover:text-gray-600">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
              </svg>
            </button>
          </div>
          
          <div class="mt-6 grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Message Info -->
            <div>
              <h4 class="text-sm font-medium text-gray-900 mb-3">Message Information</h4>
              <dl class="space-y-3">
                <div>
                  <dt class="text-sm text-gray-500">Subject/Title</dt>
                  <dd class="text-sm text-gray-900 font-medium">{{ selectedMessage.subject || selectedMessage.title }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Content</dt>
                  <dd class="text-sm text-gray-900">{{ selectedMessage.content }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Type</dt>
                  <dd>
                    <span :class="getMessageTypeClass(selectedMessage.type)" class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium">
                      {{ getMessageTypeLabel(selectedMessage.type) }}
                    </span>
                  </dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Channel</dt>
                  <dd class="text-sm text-gray-900 capitalize">{{ selectedMessage.channel }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Sent Date</dt>
                  <dd class="text-sm text-gray-900">{{ formatFullDate(selectedMessage.sentAt) }}</dd>
                </div>
              </dl>
            </div>

            <!-- Delivery Stats -->
            <div>
              <h4 class="text-sm font-medium text-gray-900 mb-3">Delivery Statistics</h4>
              <dl class="space-y-3">
                <div>
                  <dt class="text-sm text-gray-500">Total Recipients</dt>
                  <dd class="text-sm text-gray-900 font-medium">{{ selectedMessage.recipientCount }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Delivered</dt>
                  <dd class="text-sm text-green-600 font-medium">{{ selectedMessage.deliveredCount }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Failed</dt>
                  <dd class="text-sm text-red-600 font-medium">{{ selectedMessage.failedCount || 0 }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Pending</dt>
                  <dd class="text-sm text-yellow-600 font-medium">{{ selectedMessage.pendingCount || 0 }}</dd>
                </div>
                <div>
                  <dt class="text-sm text-gray-500">Delivery Rate</dt>
                  <dd class="text-sm text-gray-900 font-medium">{{ ((selectedMessage.deliveredCount / selectedMessage.recipientCount) * 100).toFixed(1) }}%</dd>
                </div>
              </dl>
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
  const messages = ref([])
  const loading = ref(false)
  const searchQuery = ref('')
  const selectedMessage = ref(null)
  const currentPage = ref(1)
  const itemsPerPage = ref(10)

  // Filters
  const filters = ref({
    dateRange: 'all',
    messageType: 'all',
    status: 'all',
    channel: 'all'
  })

  defineProps({
  filter: String,
  dateFilter: String
})

defineEmits(['resendMessage', 'viewDetails'])


  // Computed properties
  const filteredMessages = computed(() => {
    let filtered = messages.value

    // Apply search filter
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase()
      filtered = filtered.filter(message => 
        message.subject?.toLowerCase().includes(query) ||
        message.title?.toLowerCase().includes(query) ||
        message.content.toLowerCase().includes(query)
      )
    }

    // Apply other filters
    if (filters.value.messageType !== 'all') {
      filtered = filtered.filter(message => message.type === filters.value.messageType)
    }

    if (filters.value.status !== 'all') {
      filtered = filtered.filter(message => message.deliveryStatus === filters.value.status)
    }

    if (filters.value.channel !== 'all') {
      filtered = filtered.filter(message => message.channel === filters.value.channel)
    }

    return filtered
  })

  const totalPages = computed(() => Math.ceil(filteredMessages.value.length / itemsPerPage.value))

  // Methods
  const getMessageTypeClass = (type) => {
    const classes = {
      'bulk': 'bg-blue-100 text-blue-800',
      'event': 'bg-purple-100 text-purple-800',
      'seva': 'bg-orange-100 text-orange-800',
      'donation': 'bg-green-100 text-green-800',
      'announcement': 'bg-indigo-100 text-indigo-800'
    }
    return classes[type] || 'bg-gray-100 text-gray-800'
  }

  const getMessageTypeLabel = (type) => {
    const labels = {
      'bulk': 'Bulk Message',
      'event': 'Event Notification',
      'seva': 'Seva Reminder',
      'donation': 'Donation Receipt',
      'announcement': 'Announcement'
    }
    return labels[type] || 'Message'
  }

  const getChannelIconClass = (channel) => {
    const classes = {
      'email': 'bg-blue-100 text-blue-600',
      'sms': 'bg-green-100 text-green-600',
      'whatsapp': 'bg-green-100 text-green-600',
      'push': 'bg-purple-100 text-purple-600'
    }
    return classes[channel] || 'bg-gray-100 text-gray-600'
  }

  const getChannelIcon = (channel) => {
    // Return appropriate icon component based on channel
    return 'svg' // placeholder
  }

  const getStatusClass = (status) => {
    const classes = {
      'sent': 'bg-blue-100 text-blue-800',
      'delivered': 'bg-green-100 text-green-800',
      'failed': 'bg-red-100 text-red-800',
      'pending': 'bg-yellow-100 text-yellow-800'
    }
    return classes[status] || 'bg-gray-100 text-gray-800'
  }

  const getStatusLabel = (status) => {
    const labels = {
      'sent': 'Sent',
      'delivered': 'Delivered',
      'failed': 'Failed',
      'pending': 'Pending'
    }
    return labels[status] || 'Unknown'
  }

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString()
  }

  const formatTime = (dateString) => {
    return new Date(dateString).toLocaleTimeString()
  }

  const formatFullDate = (dateString) => {
    return new Date(dateString).toLocaleString()
  }

  const viewMessageDetails = (message) => {
    selectedMessage.value = message
  }

  const closeMessageDetails = () => {
    selectedMessage.value = null
  }

  const resendMessage = (message) => {
    // Implement resend logic
    console.log('Resending message:', message.id)
  }

  const downloadMessageReport = (message) => {
    // Implement download logic
    console.log('Downloading report for message:', message.id)
  }

  const exportHistory = () => {
    // Implement export logic
    console.log('Exporting message history')
  }

  // Load mock data
  onMounted(() => {
    loading.value = true
    
    // Mock message history data
    messages.value = [
      {
        id: 1,
        subject: 'Ganesha Chaturthi Celebration',
        content: 'Join us for the grand Ganesha Chaturthi celebration on August 31st...',
        type: 'event',
        channel: 'email',
        deliveryStatus: 'delivered',
        recipientCount: 245,
        deliveredCount: 242,
        failedCount: 3,
        pendingCount: 0,
        recipientGroups: ['All Devotees', 'Volunteers'],
        sentAt: '2024-08-25T10:30:00Z',
        sentBy: 'Temple Admin'
      },
      {
        id: 2,
        title: 'Weekly Seva Reminder',
        content: 'This is a reminder for your upcoming seva booking...',
        type: 'seva',
        channel: 'sms',
        deliveryStatus: 'delivered',
        recipientCount: 89,
        deliveredCount: 87,
        failedCount: 2,
        pendingCount: 0,
        recipientGroups: ['Seva Participants'],
        sentAt: '2024-08-24T08:00:00Z',
        sentBy: 'System'
      },
      {
        id: 3,
        subject: 'Donation Receipt',
        content: 'Thank you for your generous donation of ₹5,000...',
        type: 'donation',
        channel: 'email',
        deliveryStatus: 'delivered',
        recipientCount: 1,
        deliveredCount: 1,
        failedCount: 0,
        pendingCount: 0,
        recipientGroups: ['Individual'],
        sentAt: '2024-08-23T14:15:00Z',
        sentBy: 'System'
      },
      {
        id: 4,
        title: 'Temple Announcement',
        content: 'Important announcement regarding temple timings...',
        type: 'announcement',
        channel: 'whatsapp',
        deliveryStatus: 'failed',
        recipientCount: 156,
        deliveredCount: 0,
        failedCount: 156,
        pendingCount: 0,
        recipientGroups: ['All Devotees'],
        sentAt: '2024-08-22T16:45:00Z',
        sentBy: 'Temple Admin'
      },
      {
        id: 5,
        subject: 'Monthly Newsletter',
        content: 'Read about this month\'s temple activities and upcoming events...',
        type: 'bulk',
        channel: 'email',
        deliveryStatus: 'pending',
        recipientCount: 324,
        deliveredCount: 298,
        failedCount: 12,
        pendingCount: 14,
        recipientGroups: ['All Devotees', 'Volunteers', 'Donors'],
        sentAt: '2024-08-21T09:00:00Z',
        sentBy: 'Temple Admin'
      }
    ]
    
    loading.value = false
  })
  </script>