<template>
  <div class="w-full">
    <!-- Filters Section -->
    <BaseCard class="mb-4 overflow-hidden border border-gray-200 shadow-sm">
      <!-- Your existing filters here -->
    </BaseCard>

    <!-- Logs Table -->
    <BaseCard class="border border-gray-200 shadow-sm">
      <template #header>
        <div class="flex justify-between items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">Audit Logs</h3>
        </div>
      </template>
      
      <div class="overflow-x-auto">
        <!-- Basic debugging panel -->
        <div v-if="logs.length === 0 && !isLoading" class="p-3 bg-yellow-50 text-sm">
          <strong>No logs found. API Response information:</strong>
          <div class="mt-1">Total: {{ total }}, Page: {{ currentPage }}, Limit: {{ limit }}</div>
          <div v-if="error" class="mt-1 text-red-600">Error: {{ error }}</div>
        </div>
        
        <!-- Add columns prop to BaseTable -->
        <BaseTable :columns="tableColumns">
          <template #header>
            <BaseTable.HeaderCell>ID</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>User Name</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>Entity Name</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>Action</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>Status</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>IP Address</BaseTable.HeaderCell>
            <BaseTable.HeaderCell>Created At</BaseTable.HeaderCell>
          </template>
          
          <template #body>
            <BaseTable.Row 
              v-for="log in logs" 
              :key="log.id"
              @click="selectLog(log.id)"
              class="cursor-pointer hover:bg-indigo-50 transition-colors duration-150"
            >
              <BaseTable.Cell>{{ log.id }}</BaseTable.Cell>
              <BaseTable.Cell>{{ log.userName }}</BaseTable.Cell>
              <BaseTable.Cell>{{ log.entityName }}</BaseTable.Cell>
              <BaseTable.Cell>
                <span class="font-medium">{{ log.action }}</span>
              </BaseTable.Cell>
              <BaseTable.Cell>
                <BaseBadge
                  :variant="log.status === 'success' ? 'success' : 'error'"
                >
                  {{ log.status }}
                </BaseBadge>
              </BaseTable.Cell>
              <BaseTable.Cell>{{ log.ipAddress }}</BaseTable.Cell>
              <BaseTable.Cell>{{ log.formattedDate }}</BaseTable.Cell>
            </BaseTable.Row>
            
            <!-- Empty state -->
            <tr v-if="!isLoading && logs.length === 0">
              <td colspan="7" class="px-6 py-10 text-center text-gray-500">
                <div class="flex flex-col items-center justify-center">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-gray-400 mb-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                  </svg>
                  <p class="text-lg font-medium">No audit logs found</p>
                  <p class="text-sm mt-1">Try adjusting your filters or check back later</p>
                </div>
              </td>
            </tr>
            
            <!-- Loading state -->
            <tr v-if="isLoading">
              <td colspan="7" class="px-6 py-10 text-center">
                <BaseLoader class="mx-auto" />
              </td>
            </tr>
          </template>
        </BaseTable>
      </div>
      
      <!-- Pagination -->
      <template #footer>
        <div v-if="logs.length > 0 || total > 0" class="flex items-center justify-between border-t border-gray-200 px-4 py-3">
          <div class="flex flex-1 justify-between sm:hidden">
            <BaseButton
              variant="outline"
              size="sm"
              :disabled="currentPage <= 1"
              @click="setPage(currentPage - 1)"
            >
              Previous
            </BaseButton>
            <BaseButton
              variant="outline"
              size="sm"
              :disabled="currentPage >= totalPages"
              @click="setPage(currentPage + 1)"
            >
              Next
            </BaseButton>
          </div>
          <div class="hidden sm:flex sm:flex-1 sm:items-center sm:justify-between">
            <div>
              <p class="text-sm text-gray-700">
                Showing <span class="font-medium">{{ logs.length ? (currentPage - 1) * limit + 1 : 0 }}</span> to <span class="font-medium">{{ Math.min(currentPage * limit, total) }}</span> of <span class="font-medium">{{ total }}</span> results
              </p>
            </div>
            <div>
              <nav class="isolate inline-flex -space-x-px rounded-md shadow-sm" aria-label="Pagination">
                <button
                  @click="setPage(currentPage - 1)"
                  :disabled="currentPage <= 1"
                  class="relative inline-flex items-center rounded-l-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage <= 1 }"
                >
                  <span class="sr-only">Previous</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M12.79 5.23a.75.75 0 01-.02 1.06L8.832 10l3.938 3.71a.75.75 0 11-1.04 1.08l-4.5-4.25a.75.75 0 010-1.08l4.5-4.25a.75.75 0 011.06.02z" clip-rule="evenodd" />
                  </svg>
                </button>
                
                <template v-for="page in displayedPages" :key="page">
                  <span
                    v-if="page === '...'"
                    class="relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-700 ring-1 ring-inset ring-gray-300 focus:outline-offset-0"
                  >
                    ...
                  </span>
                  <button
                    v-else
                    @click="setPage(page)"
                    :class="[
                      page === currentPage
                        ? 'relative z-10 inline-flex items-center bg-indigo-600 px-4 py-2 text-sm font-semibold text-white focus:z-20 focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600'
                        : 'relative inline-flex items-center px-4 py-2 text-sm font-semibold text-gray-900 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0'
                    ]"
                  >
                    {{ page }}
                  </button>
                </template>
                
                <button
                  @click="setPage(currentPage + 1)"
                  :disabled="currentPage >= totalPages"
                  class="relative inline-flex items-center rounded-r-md px-2 py-2 text-gray-400 ring-1 ring-inset ring-gray-300 hover:bg-gray-50 focus:z-20 focus:outline-offset-0"
                  :class="{ 'opacity-50 cursor-not-allowed': currentPage >= totalPages }"
                >
                  <span class="sr-only">Next</span>
                  <svg class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor" aria-hidden="true">
                    <path fill-rule="evenodd" d="M7.21 14.77a.75.75 0 01.02-1.06L11.168 10 7.23 6.29a.75.75 0 111.04-1.08l4.5 4.25a.75.75 0 010 1.08l-4.5 4.25a.75.75 0 01-1.06-.02z" clip-rule="evenodd" />
                  </svg>
                </button>
              </nav>
            </div>
          </div>
        </div>
      </template>
    </BaseCard>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { useAuditLogStore } from '@/stores/auditlog';
import BaseTable from '@/components/common/BaseTable.vue';
import BaseCard from '@/components/common/BaseCard.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseBadge from '@/components/common/BaseBadge.vue';
import BaseLoader from '@/components/common/BaseLoader.vue';

const auditLogStore = useAuditLogStore();

// Table columns definition to fix the missing columns prop error
const tableColumns = [
  { key: 'id', label: 'ID' },
  { key: 'userName', label: 'User Name' },
  { key: 'entityName', label: 'Entity Name' },
  { key: 'action', label: 'Action' },
  { key: 'status', label: 'Status' },
  { key: 'ipAddress', label: 'IP Address' },
  { key: 'createdAt', label: 'Created At' }
];

// Expose emit events for parent component
const emit = defineEmits(['select-log']);

// Computed values from store
const logs = computed(() => {
  console.log("AuditLogsList - formatted logs:", auditLogStore.formattedLogs);
  return auditLogStore.formattedLogs;
});
const isLoading = computed(() => auditLogStore.isLoading);
const error = computed(() => auditLogStore.error);
const currentPage = computed(() => auditLogStore.currentPage);
const totalPages = computed(() => auditLogStore.totalPages);
const total = computed(() => auditLogStore.total);
const limit = computed(() => auditLogStore.limit);

// Calculate which pages to display in pagination
const displayedPages = computed(() => {
  const total = totalPages.value;
  const current = currentPage.value;
  
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1);
  }
  
  if (current <= 3) {
    return [1, 2, 3, 4, '...', total];
  }
  
  if (current >= total - 2) {
    return [1, '...', total - 3, total - 2, total - 1, total];
  }
  
  return [1, '...', current - 1, current, current + 1, '...', total];
});

// Methods
function selectLog(id) {
  emit('select-log', id);
}

function setPage(page) {
  if (page < 1 || page > totalPages.value) return;
  auditLogStore.setPage(page);
}
</script>