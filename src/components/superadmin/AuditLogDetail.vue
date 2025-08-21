<template>
  <div>
    <!-- Use wrapper div to handle props properly -->
    <BaseModal v-if="!!selectedLog" @close="closeDrawer">
      <template #header>
        <div class="flex items-center justify-between bg-gray-50 px-4 py-3 border-b border-gray-200 rounded-t-lg">
          <h3 class="text-lg font-bold leading-6 text-gray-900">
            Audit Log Details
          </h3>
          <BaseBadge
            v-if="selectedLog"
            :variant="selectedLog.status === 'success' ? 'success' : 'error'"
            class="ml-2"
          >
            {{ selectedLog.status }}
          </BaseBadge>
        </div>
      </template>
      
      <div class="p-6 space-y-6">
        <!-- Loading state -->
        <div v-if="isLoading" class="flex justify-center items-center py-16">
          <BaseLoader />
        </div>

        <!-- Error state -->
        <div v-else-if="error" class="flex justify-center items-center py-8">
          <BaseAlert variant="error">
            {{ error }}
          </BaseAlert>
        </div>

        <!-- Content -->
        <div v-else-if="selectedLog" class="space-y-6">
          <!-- Basic info section -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-x-6 gap-y-4 p-5 bg-gray-50 rounded-lg border border-gray-200">
            <div>
              <label class="block text-sm font-medium text-gray-500">ID</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ selectedLog.id }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Action</label>
              <div class="mt-1 text-sm font-semibold text-indigo-700">{{ selectedLog.action }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">User</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ selectedLog.userName }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Entity</label>
              <div class="mt-1 text-sm text-gray-900 font-semibold">{{ selectedLog.entityName || '-' }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">IP Address</label>
              <div class="mt-1 text-sm text-gray-900">{{ selectedLog.ipAddress }}</div>
            </div>
            
            <div>
              <label class="block text-sm font-medium text-gray-500">Date & Time</label>
              <div class="mt-1 text-sm text-gray-900">{{ selectedLog.formattedDate }}</div>
            </div>
          </div>
          
          <!-- Details -->
          <div class="pt-4">
            <label class="block text-sm font-medium text-gray-500 mb-3">Details</label>
            
            <div v-if="!selectedLog.details" class="text-sm text-gray-500 italic bg-gray-50 p-4 rounded-md border border-gray-200">
              No additional details available
            </div>
            
            <!-- JSON Viewer -->
            <div 
              v-else 
              class="bg-gray-50 p-4 rounded-md border border-gray-200 font-mono text-sm overflow-auto max-h-96 shadow-inner"
            >
              <pre class="json-formatter">{{ formattedDetails }}</pre>
            </div>
          </div>
          
          <!-- Request/Response Info (if available) -->
          <div v-if="selectedLog.request || selectedLog.response" class="border-t border-gray-200 pt-4 grid grid-cols-1 gap-4">
            <div v-if="selectedLog.request">
              <label class="block text-sm font-medium text-gray-500 mb-2">Request</label>
              <div class="bg-gray-50 p-4 rounded-md border border-gray-200 font-mono text-sm overflow-auto max-h-48 shadow-inner">
                <pre>{{ formatJSON(selectedLog.request) }}</pre>
              </div>
            </div>
            
            <div v-if="selectedLog.response">
              <label class="block text-sm font-medium text-gray-500 mb-2">Response</label>
              <div class="bg-gray-50 p-4 rounded-md border border-gray-200 font-mono text-sm overflow-auto max-h-48 shadow-inner">
                <pre>{{ formatJSON(selectedLog.response) }}</pre>
              </div>
            </div>
          </div>
        </div>
        
        <div v-else class="flex justify-center items-center py-8 text-gray-500">
          No log selected
        </div>
      </div>
      
      <template #footer>
        <div class="flex justify-end gap-3 bg-gray-50 px-4 py-3 border-t border-gray-200 rounded-b-lg">
          <BaseButton 
            variant="outline" 
            @click="closeDrawer"
            class="w-full sm:w-auto"
          >
            Close
          </BaseButton>
        </div>
      </template>
    </BaseModal>
  </div>
</template>

<script setup>
import { computed, watch } from 'vue';
import { useAuditLogStore } from '@/stores/auditlog';
import BaseAlert from '@/components/common/BaseAlert.vue';
import BaseBadge from '@/components/common/BaseBadge.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseModal from '@/components/common/BaseModal.vue';
import BaseLoader from '@/components/common/BaseLoader.vue';

const props = defineProps({
  logId: {
    type: [String, Number],
    default: null
  }
});

const emit = defineEmits(['close']);

const auditLogStore = useAuditLogStore();

const selectedLog = computed(() => auditLogStore.selectedLog);
const isLoading = computed(() => auditLogStore.isDetailLoading);
const error = computed(() => auditLogStore.error);

// Format JSON details for better readability
const formattedDetails = computed(() => {
  if (!selectedLog.value || !selectedLog.value.details) return '';
  
  // If details is already a string, format it as JSON
  if (typeof selectedLog.value.details === 'string') {
    try {
      return formatJSON(JSON.parse(selectedLog.value.details));
    } catch (e) {
      return selectedLog.value.details;
    }
  }
  
  // If details is an object, format it directly
  return formatJSON(selectedLog.value.details);
});

// Helper function to format JSON with proper indentation and colors
function formatJSON(obj) {
  return JSON.stringify(obj, null, 2);
}

// Watch for changes in logId prop
watch(() => props.logId, (newId) => {
  if (newId) {
    auditLogStore.fetchAuditLogDetails(newId);
  }
}, { immediate: true });

// Close the drawer and clear selected log
function closeDrawer() {
  auditLogStore.clearSelectedLog();
  emit('close');
}
</script>

<style scoped>
.json-formatter {
  white-space: pre-wrap;
  word-break: break-word;
}

/* Syntax highlighting colors */
.json-formatter .string { color: #22c55e; }
.json-formatter .number { color: #3b82f6; }
.json-formatter .boolean { color: #f59e0b; }
.json-formatter .null { color: #ef4444; }
.json-formatter .key { color: #8b5cf6; }

/* Custom scrollbar for JSON viewer */
.max-h-96::-webkit-scrollbar,
.max-h-48::-webkit-scrollbar {
  width: 6px;
  height: 6px;
}

.max-h-96::-webkit-scrollbar-track,
.max-h-48::-webkit-scrollbar-track {
  background: #f1f5f9;
  border-radius: 3px;
}

.max-h-96::-webkit-scrollbar-thumb,
.max-h-48::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 3px;
}

.max-h-96::-webkit-scrollbar-thumb:hover,
.max-h-48::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}
</style>