<template>
  <div class="bg-white rounded-xl shadow-lg p-6">
    <!-- Header -->
    <div class="flex items-center justify-between mb-6">
      <div>
        <h2 class="text-2xl font-bold text-gray-900">Message Templates</h2>
        <p class="text-gray-600 mt-1">Create and manage reusable message templates</p>
      </div>
      <button
        @click="showCreateModal = true"
        class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-all duration-200 flex items-center gap-2"
      >
        <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Create Template
      </button>
    </div>

    <!-- Search and Filter -->
    <div class="flex flex-col md:flex-row gap-4 mb-6">
      <div class="flex-1">
        <div class="relative">
          <svg class="absolute left-3 top-1/2 transform -translate-y-1/2 w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
          </svg>
          <input
            v-model="searchQuery"
            type="text"
            placeholder="Search templates..."
            class="pl-10 pr-4 py-2 w-full border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
          />
        </div>
      </div>
      <select
        v-model="selectedCategory"
        class="px-4 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
      >
        <option value="">All Categories</option>
        <option value="seva">Seva Notifications</option>
        <option value="event">Event Announcements</option>
        <option value="donation">Donation Receipts</option>
        <option value="general">General Messages</option>
        <option value="festival">Festival Greetings</option>
      </select>
    </div>

    <!-- Templates Grid -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6" v-if="filteredTemplates.length > 0">
      <div
        v-for="template in filteredTemplates"
        :key="template.id"
        class="bg-gray-50 rounded-lg p-4 hover:shadow-md transition-all duration-200 border border-gray-200"
      >
        <!-- Template Header -->
        <div class="flex items-start justify-between mb-3">
          <div class="flex-1">
            <h3 class="font-semibold text-gray-900 mb-1">{{ template.name }}</h3>
            <span :class="getCategoryBadgeClass(template.category)" class="px-2 py-1 rounded-full text-xs font-medium">
              {{ template.category }}
            </span>
          </div>
          <div class="flex items-center gap-2">
            <button
              @click="editTemplate(template)"
              class="text-gray-400 hover:text-indigo-600 transition-colors"
              title="Edit Template"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
              </svg>
            </button>
            <button
              @click="deleteTemplate(template.id)"
              class="text-gray-400 hover:text-red-600 transition-colors"
              title="Delete Template"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
              </svg>
            </button>
          </div>
        </div>

        <!-- Template Preview -->
        <div class="mb-3">
          <p class="text-sm text-gray-600 line-clamp-3">{{ template.content }}</p>
        </div>

        <!-- Template Info -->
        <div class="flex items-center justify-between text-xs text-gray-500">
          <span>{{ template.usageCount }} times used</span>
          <span>{{ formatDate(template.lastUsed) }}</span>
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-2 mt-3">
          <button
            @click="previewTemplate(template)"
            class="flex-1 bg-white hover:bg-gray-50 text-gray-700 border border-gray-300 px-3 py-2 rounded-lg text-sm font-medium transition-colors"
          >
            Preview
          </button>
          <button
            @click="useTemplate(template)"
            class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white px-3 py-2 rounded-lg text-sm font-medium transition-colors"
          >
            Use Template
          </button>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No templates found</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by creating your first message template.</p>
      <div class="mt-6">
        <button
          @click="showCreateModal = true"
          class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
        >
          Create Template
        </button>
      </div>
    </div>

    <!-- Create/Edit Modal -->
    <div v-if="showCreateModal || editingTemplate" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="p-6">
          <div class="flex items-center justify-between mb-6">
            <h3 class="text-lg font-semibold text-gray-900">
              {{ editingTemplate ? 'Edit Template' : 'Create New Template' }}
            </h3>
            <button
              @click="closeModal"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="saveTemplate" class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Template Name</label>
              <input
                v-model="templateForm.name"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Enter template name"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Category</label>
              <select
                v-model="templateForm.category"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
              >
                <option value="">Select Category</option>
                <option value="seva">Seva Notifications</option>
                <option value="event">Event Announcements</option>
                <option value="donation">Donation Receipts</option>
                <option value="general">General Messages</option>
                <option value="festival">Festival Greetings</option>
              </select>
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Subject</label>
              <input
                v-model="templateForm.subject"
                type="text"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors"
                placeholder="Enter message subject"
              />
            </div>

            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Message Content</label>
              <textarea
                v-model="templateForm.content"
                rows="6"
                required
                class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors resize-none"
                placeholder="Enter your message content..."
              ></textarea>
            </div>

            <!-- <div class="bg-gray-50 rounded-lg p-4"> -->
              <!-- <h4 class="text-sm font-medium text-gray-700 mb-2">Available Variables</h4> -->
              <!-- <div class="flex flex-wrap gap-2"> -->
                <!-- <span
                  v-for="variable in availableVariables"
                  :key="variable"
                  @click="insertVariable(variable)"
                  class="bg-indigo-100 text-indigo-800 px-2 py-1 rounded text-xs cursor-pointer hover:bg-indigo-200 transition-colors"
                >
                  {{ variable }}
                </span> -->
              <!-- </div> -->
            <!-- </div> -->

            <div class="flex justify-end gap-3 pt-4">
              <button
                type="button"
                @click="closeModal"
                class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
              >
                Cancel
              </button>
              <button
                type="submit"
                class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
              >
                {{ editingTemplate ? 'Update Template' : 'Create Template' }}
              </button>
            </div>
          </form>
        </div>
      </div>
    </div>

    <!-- Preview Modal -->
    <div v-if="previewingTemplate" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl max-w-lg w-full">
        <div class="p-6">
          <div class="flex items-center justify-between mb-4">
            <h3 class="text-lg font-semibold text-gray-900">Template Preview</h3>
            <button
              @click="previewingTemplate = null"
              class="text-gray-400 hover:text-gray-600 transition-colors"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <div class="space-y-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Subject:</label>
              <p class="text-gray-900">{{ previewingTemplate.subject }}</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Content:</label>
              <div class="bg-gray-50 rounded-lg p-4">
                <p class="text-gray-900 whitespace-pre-wrap">{{ previewingTemplate.content }}</p>
              </div>
            </div>
          </div>

          <div class="flex justify-end gap-3 mt-6">
            <button
              @click="previewingTemplate = null"
              class="px-4 py-2 text-gray-700 border border-gray-300 rounded-lg hover:bg-gray-50 transition-colors"
            >
              Close
            </button>
            <button
              @click="useTemplate(previewingTemplate)"
              class="px-4 py-2 bg-indigo-600 hover:bg-indigo-700 text-white rounded-lg transition-colors"
            >
              Use This Template
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import communicationService from '@/services/communication.service'
import { useAuthStore } from '@/stores/auth' // If using Pinia
import { useRouter } from 'vue-router'
// Reactive data
const templates = ref([])
const searchQuery = ref('')
const selectedCategory = ref('')
const showCreateModal = ref(false)
const editingTemplate = ref(null)
const previewingTemplate = ref(null)
const authStore = useAuthStore()
const entityId = authStore.user?.entityId
const router = useRouter()

// Form data
const templateForm = ref({
  name: '',
  category: '',
  subject: '',
  content: ''
})

// Available variables for templates
const availableVariables = ref([
  '{devotee_name}',
  '{temple_name}',
  '{seva_name}',
  '{event_name}',
  '{date}',
  '{time}',
  '{amount}',
  '{receipt_number}'
])

// Computed properties
const filteredTemplates = computed(() => {
  let filtered = templates.value || []

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()

    filtered = filtered.filter(template => {
      const name = template.name?.toLowerCase() || ''
      const content = template.content?.toLowerCase() || ''
      return name.includes(query) || content.includes(query)
    })
  }

  if (selectedCategory.value) {
    filtered = filtered.filter(template => template.category === selectedCategory.value)
  }

  return filtered
})


// Methods
const getCategoryBadgeClass = (category) => {
  const classes = {
    'seva': 'bg-blue-100 text-blue-800',
    'event': 'bg-green-100 text-green-800',
    'donation': 'bg-yellow-100 text-yellow-800',
    'general': 'bg-gray-100 text-gray-800',
    'festival': 'bg-purple-100 text-purple-800'
  }
  return classes[category] || 'bg-gray-100 text-gray-800'
}

const formatDate = (dateString) => {
  if (!dateString) return 'Never used'
  return new Date(dateString).toLocaleDateString()
}

const editTemplate = (template) => {
  editingTemplate.value = template
  templateForm.value = { ...template }
  showCreateModal.value = false
}

const deleteTemplate = (templateId) => {
  if (confirm('Are you sure you want to delete this template?')) {
    templates.value = templates.value.filter(t => t.id !== templateId)
  }
}

const previewTemplate = (template) => {
  previewingTemplate.value = {
    ...template,
    content: template.content || template.body // 👈 fallback for correct rendering
  }
}


const useTemplate = (template) => {
  // Navigate with template ID or subject/content in query
  router.push({
    name: 'MessageComposer',
    query: {
      subject: template.subject,
      content: template.content || template.body
    }
  })
  previewingTemplate.value = null
  showCreateModal.value = false
  editingTemplate.value = null
}
const insertVariable = (variable) => {
  const textarea = document.querySelector('textarea')
  if (textarea) {
    const start = textarea.selectionStart
    const end = textarea.selectionEnd
    const content = templateForm.value.content
    templateForm.value.content = content.substring(0, start) + variable + content.substring(end)
  }
}

const saveTemplate = async () => {
  const payload = {
    name: templateForm.value.name,
    subject: templateForm.value.subject,
    body: templateForm.value.content,
    category: templateForm.value.category,
    userId: authStore.user?.id,
    entityId: authStore.user?.entity_id
  }

  try {
    let result
    if (editingTemplate.value) {
      // ⚠️ Update API not yet connected — optional future step
      console.warn('Edit functionality is not wired to backend yet')
    } else {
      result = await communicationService.createTemplate(payload)
    }

    if (result?.success) {
      templates.value.push(result.data)
      closeModal()
    } else {
      alert(result?.error || 'Something went wrong while saving the template.')
    }
  } catch (err) {
    console.error('Error saving template:', err)
    alert('Unexpected error occurred.')
  }
}

const closeModal = () => {
  showCreateModal.value = false
  editingTemplate.value = null
  templateForm.value = {
    name: '',
    category: '',
    subject: '',
    content: ''
  }
}

// Load templates on mount
onMounted(async () => {
  const result = await communicationService.getTemplates()

  if (result.success) {
    templates.value = result.data
  } else {
    console.error("Failed to load templates:", result.error)
  }
})

</script>
