<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <div class="bg-white border-b border-gray-200 px-4 py-6 sm:px-6">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Communication Center</h1>
          <p class="mt-1 text-sm text-gray-600">
            Send messages, manage templates, and track communication history
          </p>
        </div>
        
        <!-- Quick Stats -->
        <!-- <div class="flex space-x-4">
          <div class="text-center">
            <div class="text-2xl font-bold text-indigo-600">{{ stats.totalMessages }}</div>
            <div class="text-xs text-gray-500">Total Sent</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-green-600">{{ stats.deliveredMessages }}</div>
            <div class="text-xs text-gray-500">Delivered</div>
          </div>
          <div class="text-center">
            <div class="text-2xl font-bold text-orange-600">{{ stats.activeTemplates }}</div>
            <div class="text-xs text-gray-500">Templates</div>
          </div>
        </div> -->
      </div>
    </div>

    <!-- Navigation Tabs -->
    <div class="bg-white border-b border-gray-200">
      <nav class="flex space-x-8 px-4 sm:px-6">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'py-4 px-1 border-b-2 font-medium text-sm transition-colors duration-200',
            activeTab === tab.id
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          ]"
        >
          <div class="flex items-center space-x-2">
            <component :is="tab.icon" class="w-5 h-5" />
            <span>{{ tab.name }}</span>
            <span
              v-if="tab.count"
              class="bg-gray-100 text-gray-900 py-0.5 px-2 rounded-full text-xs"
            >
              {{ tab.count }}
            </span>
          </div>
        </button>
      </nav>
    </div>

    <!-- Content Area -->
    <div class="p-4 sm:p-6">
      <!-- Compose Message Tab -->
      <div v-if="activeTab === 'compose'" class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200">
          <div class="p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Compose New Message</h2>
            <MessageComposer 
              @message-sent="handleMessageSent"
              @save-template="handleSaveTemplate"
            />
          </div>
        </div>
      </div>

      <!-- Bulk Message Tab -->
      <div v-if="activeTab === 'bulk'" class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200">
          <div class="p-6">
            <h2 class="text-lg font-semibold text-gray-900 mb-4">Bulk Messaging</h2>
            <BulkMessage 
              @bulk-sent="handleBulkSent"
              @preview-requested="handlePreviewRequest"
            />
          </div>
        </div>
      </div>

      <!-- Templates Tab -->
      <div v-if="activeTab === 'templates'" class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200">
          <div class="p-6">
            <!-- <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-gray-900">Message Templates</h2>
              <button
                @click="showCreateTemplate = true"
                class="bg-indigo-600 text-white px-4 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200 flex items-center space-x-2"
              >
                <Plus class="w-4 h-4" />
                <span>New Template</span>
              </button>
            </div> -->
            <MessageTemplates 
              @template-selected="handleTemplateSelect"
              @template-edited="handleTemplateEdit"
              @template-deleted="handleTemplateDelete"
            />
          </div>
        </div>
      </div>

      <!-- History Tab -->
      <div v-if="activeTab === 'history'" class="space-y-6">
        <div class="bg-white rounded-xl shadow-sm border border-gray-200">
          <div class="p-6">
            <!-- <div class="flex items-center justify-between mb-4">
              <h2 class="text-lg font-semibold text-gray-900">Message History</h2>
              
              
              <div class="flex items-center space-x-3">
                <select 
                  v-model="historyFilter"
                  class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                >
                  <option value="all">All Messages</option>
                  <option value="sent">Sent</option>
                  <option value="delivered">Delivered</option>
                  <option value="failed">Failed</option>
                </select>
                
                <input
                  type="date"
                  v-model="dateFilter"
                  class="border border-gray-300 rounded-lg px-3 py-2 text-sm focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
                />
                
                <button
                  @click="exportHistory"
                  class="bg-gray-100 text-gray-700 px-4 py-2 rounded-lg hover:bg-gray-200 transition-colors duration-200 flex items-center space-x-2"
                >
                  <Download class="w-4 h-4" />
                  <span>Export</span>
                </button>
              </div>
            </div> -->
            
            <MessageHistory 
              :filter="historyFilter"
              :date-filter="dateFilter"
              @resend-message="handleResendMessage"
              @view-details="handleViewDetails"
            />
          </div>
        </div>
      </div>

      <!-- Analytics Tab -->
      <div v-if="activeTab === 'analytics'" class="space-y-6">
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Message Statistics -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Message Statistics</h3>
            <div class="space-y-4">
              <div class="flex items-center justify-between">
                <span class="text-gray-600">Messages Sent Today</span>
                <span class="font-semibold text-gray-900">{{ analytics.today }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600">This Week</span>
                <span class="font-semibold text-gray-900">{{ analytics.week }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600">This Month</span>
                <span class="font-semibold text-gray-900">{{ analytics.month }}</span>
              </div>
              <div class="flex items-center justify-between">
                <span class="text-gray-600">Delivery Rate</span>
                <span class="font-semibold text-green-600">{{ analytics.deliveryRate }}%</span>
              </div>
            </div>
          </div>

          <!-- Popular Templates -->
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Popular Templates</h3>
            <div class="space-y-3">
              <div 
                v-for="template in popularTemplates" 
                :key="template.id"
                class="flex items-center justify-between p-3 bg-gray-50 rounded-lg"
              >
                <div>
                  <div class="font-medium text-gray-900">{{ template.name }}</div>
                  <div class="text-sm text-gray-500">{{ template.category }}</div>
                </div>
                <div class="text-sm font-medium text-indigo-600">
                  {{ template.usage }} uses
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Recent Activity -->
        <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
          <h3 class="text-lg font-semibold text-gray-900 mb-4">Recent Activity</h3>
          <div class="space-y-4">
            <div 
              v-for="activity in recentActivity" 
              :key="activity.id"
              class="flex items-center space-x-4 p-4 bg-gray-50 rounded-lg"
            >
              <div 
                :class="[
                  'w-10 h-10 rounded-full flex items-center justify-center',
                  activity.type === 'sent' ? 'bg-blue-100 text-blue-600' :
                  activity.type === 'delivered' ? 'bg-green-100 text-green-600' :
                  'bg-red-100 text-red-600'
                ]"
              >
                <Mail class="w-5 h-5" />
              </div>
              <div class="flex-1">
                <div class="font-medium text-gray-900">{{ activity.message }}</div>
                <div class="text-sm text-gray-500">{{ activity.time }}</div>
              </div>
              <div 
                :class="[
                  'px-3 py-1 rounded-full text-xs font-medium',
                  activity.type === 'sent' ? 'bg-blue-100 text-blue-800' :
                  activity.type === 'delivered' ? 'bg-green-100 text-green-800' :
                  'bg-red-100 text-red-800'
                ]"
              >
                {{ activity.status }}
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Success Toast -->
    <div 
      v-if="showToast"
      class="fixed bottom-4 right-4 bg-green-500 text-white px-6 py-3 rounded-lg shadow-lg transition-all duration-300 z-50"
    >
      <div class="flex items-center space-x-2">
        <CheckCircle class="w-5 h-5" />
        <span>{{ toastMessage }}</span>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { 
  Mail, 
  Users, 
  FileText, 
  Clock, 
  BarChart3, 
  Plus, 
  Download,
  CheckCircle
} from 'lucide-vue-next'

// Import communication components
import MessageComposer from '@/components/communication/MessageComposer.vue'
import MessageTemplates from '@/components/communication/MessageTemplates.vue'
import BulkMessage from '@/components/communication/BulkMessage.vue'
import MessageHistory from '@/components/communication/MessageHistory.vue'

const route = useRoute()
const templeId = route.params.id

// Reactive state
const activeTab = ref('compose')
const historyFilter = ref('all')
const dateFilter = ref('')
const showCreateTemplate = ref(false)
const showToast = ref(false)
const toastMessage = ref('')

// Stats data
const stats = reactive({
  totalMessages: 1247,
  deliveredMessages: 1195,
  activeTemplates: 12
})

// Analytics data
const analytics = reactive({
  today: 23,
  week: 156,
  month: 678,
  deliveryRate: 95.8
})

// Tab configuration
const tabs = computed(() => [
  { 
    id: 'compose', 
    name: 'Compose', 
    icon: Mail,
    count: null
  },
  // { 
  //   id: 'bulk', 
  //   name: 'Bulk Message', 
  //   icon: Users,
  //   count: null
  // },
  { 
    id: 'templates', 
    name: 'Templates', 
    icon: FileText,
    count: stats.activeTemplates
  },
  // { 
  //   id: 'history', 
  //   name: 'History', 
  //   icon: Clock,
  //   count: stats.totalMessages
  // },
  // { 
  //   id: 'analytics', 
  //   name: 'Analytics', 
  //   icon: BarChart3,
  //   count: null
  // }
])

// Popular templates data
const popularTemplates = ref([
  {
    id: 1,
    name: 'Seva Reminder',
    category: 'Seva',
    usage: 89
  },
  {
    id: 2,
    name: 'Festival Invitation',
    category: 'Event',
    usage: 67
  },
  {
    id: 3,
    name: 'Donation Receipt',
    category: 'Finance',
    usage: 45
  },
  {
    id: 4,
    name: 'Welcome Message',
    category: 'General',
    usage: 34
  }
])

// Recent activity data
const recentActivity = ref([
  {
    id: 1,
    message: 'Bulk message sent to 45 devotees about upcoming Diwali celebration',
    time: '2 hours ago',
    type: 'sent',
    status: 'Sent'
  },
  {
    id: 2,
    message: 'Seva reminder delivered to all registered volunteers',
    time: '4 hours ago',
    type: 'delivered',
    status: 'Delivered'
  },
  {
    id: 3,
    message: 'Welcome message sent to new devotee registrations',
    time: '6 hours ago',
    type: 'delivered',
    status: 'Delivered'
  },
  {
    id: 4,
    message: 'Event notification failed to deliver to 3 recipients',
    time: '1 day ago',
    type: 'failed',
    status: 'Failed'
  }
])

// Event handlers
const handleMessageSent = (messageData) => {
  stats.totalMessages++
  stats.deliveredMessages++
  showToastMessage('Message sent successfully!')
}

const handleBulkSent = (bulkData) => {
  stats.totalMessages += bulkData.recipientCount
  stats.deliveredMessages += bulkData.deliveredCount
  showToastMessage(`Bulk message sent to ${bulkData.recipientCount} recipients!`)
}

const handleSaveTemplate = (templateData) => {
  stats.activeTemplates++
  showToastMessage('Template saved successfully!')
}

const handleTemplateSelect = (template) => {
  // Switch to compose tab with selected template
  activeTab.value = 'compose'
}

const handleTemplateEdit = (template) => {
  showToastMessage('Template updated successfully!')
}

const handleTemplateDelete = (templateId) => {
  stats.activeTemplates--
  showToastMessage('Template deleted successfully!')
}

const handlePreviewRequest = (previewData) => {
  // Handle bulk message preview
  console.log('Preview requested:', previewData)
}

const handleResendMessage = (messageId) => {
  showToastMessage('Message resent successfully!')
}

const handleViewDetails = (messageId) => {
  // Handle viewing message details
  console.log('View details for message:', messageId)
}

const exportHistory = () => {
  // Handle history export
  showToastMessage('Export started! Check your downloads.')
}

const showToastMessage = (message) => {
  toastMessage.value = message
  showToast.value = true
  setTimeout(() => {
    showToast.value = false
  }, 3000)
}

onMounted(() => {
  // Load initial data
  console.log('Communication Center loaded for temple:', templeId)
})
</script>