<template>
  <div class="w-full">
    <!-- Filters Section -->
    <BaseCard class="mb-4 overflow-hidden border border-gray-200 shadow-sm">
      <template #header>
        <div class="flex items-center justify-between px-4 py-0">
          <h3 class="text-base font-bold text-gray-900">Filters</h3>
          <BaseButton
            variant="outline"
            size="sm"
            @click="resetFilters"
            class="text-indigo-600 border-indigo-600 hover:bg-indigo-50"
          >
            Reset
          </BaseButton>
        </div>
      </template>
      
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-3 p-3">
        <!-- Action Filter -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">
            Action
          </label>
          <div class="relative">
            <div class="border border-gray-300 rounded-md shadow-sm">
              <input
                type="text"
                v-model="actionSearchTerm"
                class="block w-full pl-3 pr-10 py-1.5 text-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
                placeholder="Search actions..."
                @focus="openActionDropdown"
                @blur="handleDropdownBlur"
                ref="actionInput"
              />
              <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                </svg>
              </div>
            </div>
            
            <!-- Action Dropdown using Teleport to render at document body level -->
            <Teleport to="body">
              <div 
                v-if="showActionDropdown" 
                class="bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm border border-gray-200"
                ref="actionDropdown"
                :style="dropdownStyle"
              >
                <!-- Group by module -->
                <template v-for="(actions, module) in filteredGroupedActions" :key="module">
                  <!-- Don't show Auth section -->
                  <template v-if="module !== 'Auth' && actions.length > 0">
                    <div class="px-3 py-2 text-xs font-semibold text-gray-500 bg-gray-50 sticky top-0 border-t border-b border-gray-200 first:border-t-0">
                      {{ module }}
                    </div>
                    <div 
                      v-for="action in actions" 
                      :key="action"
                      @click="setActionFilter(action)"
                      @mousedown.prevent
                      class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-indigo-50"
                      :class="{'bg-indigo-50 text-indigo-700': action === filters.action}"
                    >
                      <span class="block truncate">{{ action }}</span>
                      <span v-if="action === filters.action" class="absolute inset-y-0 right-0 flex items-center pr-4">
                        <svg class="h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                          <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                        </svg>
                      </span>
                    </div>
                  </template>
                </template>
                
                <div v-if="Object.entries(filteredGroupedActions).filter(([module]) => module !== 'Auth').every(([_, actions]) => actions.length === 0)" class="px-3 py-4 text-sm text-gray-500 text-center">
                  No actions found
                </div>
              </div>
            </Teleport>
          </div>
        </div>
        
        <!-- Status Filter -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">
            Status
          </label>
          <div class="relative">
            <div class="border border-gray-300 rounded-md shadow-sm">
              <div 
                class="block w-full pl-3 pr-10 py-1.5 text-sm focus:outline-none sm:text-sm rounded-md cursor-pointer"
                @click="toggleStatusDropdown"
                ref="statusInput"
              >
                <span v-if="filters.status" class="block truncate">
                  {{ statusOptions.find(option => option.value === filters.status)?.label || 'Select status' }}
                </span>
                <span v-else class="block truncate text-gray-500">
                  Select status
                </span>
                <span class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
                  <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-gray-400" viewBox="0 0 20 20" fill="currentColor">
                    <path fill-rule="evenodd" d="M10 3a1 1 0 01.707.293l3 3a1 1 0 01-1.414 1.414L10 5.414 7.707 7.707a1 1 0 01-1.414-1.414l3-3A1 1 0 0110 3zm-3.707 9.293a1 1 0 011.414 0L10 14.586l2.293-2.293a1 1 0 011.414 1.414l-3 3a1 1 0 01-1.414 0l-3-3a1 1 0 010-1.414z" clip-rule="evenodd" />
                  </svg>
                </span>
              </div>
            </div>
            
            <!-- Status Dropdown using Teleport -->
            <Teleport to="body">
              <div 
                v-if="showStatusDropdown" 
                class="bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm border border-gray-200"
                ref="statusDropdown"
                :style="statusDropdownStyle"
              >
                <div 
                  v-for="option in statusOptions" 
                  :key="option.value"
                  @click="setStatusFilter(option.value)"
                  @mousedown.prevent
                  class="cursor-pointer select-none relative py-2 pl-3 pr-9 hover:bg-indigo-50"
                  :class="{'bg-indigo-50 text-indigo-700': option.value === filters.status}"
                >
                  <span class="block truncate">{{ option.label }}</span>
                  <span v-if="option.value === filters.status" class="absolute inset-y-0 right-0 flex items-center pr-4">
                    <svg class="h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                      <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                    </svg>
                  </span>
                </div>
              </div>
            </Teleport>
          </div>
        </div>
        
        <!-- Date Range - From -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">
            From Date
          </label>
          <div class="border border-gray-300 rounded-md shadow-sm">
            <input
              type="date"
              v-model="filters.from_date"
              class="block w-full pl-3 pr-10 py-1.5 text-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
              @change="applyFilters"
            />
          </div>
        </div>
        
        <!-- Date Range - To -->
        <div>
          <label class="block text-xs font-medium text-gray-700 mb-1">
            To Date
          </label>
          <div class="border border-gray-300 rounded-md shadow-sm">
            <input
              type="date"
              v-model="filters.to_date"
              class="block w-full pl-3 pr-10 py-1.5 text-sm focus:outline-none focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm rounded-md"
              @change="applyFilters"
            />
          </div>
        </div>
      </div>
    </BaseCard>

    <!-- Logs Table -->
    <BaseCard class="border border-gray-200 shadow-sm">
      <template #header>
        <div class="flex justify-between items-center px-4 py-3 bg-gray-50 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">Audit Logs</h3>
          <!-- <div class="flex items-center space-x-2">
            <span class="text-sm text-gray-600">
              Showing {{ logs.length ? (currentPage - 1) * limit + 1 : 0 }}-{{ Math.min(currentPage * limit, total) }} of {{ total }}
            </span>
            <BaseSelect
              v-model="limit"
              :options="limitOptions"
              @change="setLimit"
              class="w-20"
            />
          </div> -->
        </div>
      </template>
      
      <div class="overflow-x-auto">
        <BaseTable>
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
        <div class="flex items-center justify-between border-t border-gray-200 px-4 py-3">
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
import { ref, computed, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { useAuditLogStore } from '@/stores/auditlog';
import BaseTable from '@/components/common/BaseTable.vue';
import BaseCard from '@/components/common/BaseCard.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseBadge from '@/components/common/BaseBadge.vue';
import BaseLoader from '@/components/common/BaseLoader.vue';
import BaseSelect from '@/components/common/BaseSelect.vue';

const auditLogStore = useAuditLogStore();

// Expose the logs list to parent component
const props = defineProps({
  showDetailInDrawer: {
    type: Boolean,
    default: true
  }
});

const emit = defineEmits(['select-log']);

// Local state
const showActionDropdown = ref(false);
const showStatusDropdown = ref(false);
const actionSearchTerm = ref('');
const actionDropdown = ref(null);
const actionInput = ref(null);
const statusInput = ref(null);
const statusDropdown = ref(null);

// Dropdown positioning
const dropdownStyle = ref({
  position: 'fixed',
  zIndex: 9999,
  width: '0px',
  left: '0px',
  top: '0px'
});

const statusDropdownStyle = ref({
  position: 'fixed',
  zIndex: 9999,
  width: '0px',
  left: '0px',
  top: '0px'
});

function updateDropdownPosition() {
  if (!actionInput.value) return;
  
  const rect = actionInput.value.getBoundingClientRect();
  dropdownStyle.value = {
    position: 'fixed',
    zIndex: 9999,
    width: `${rect.width}px`,
    left: `${rect.left}px`,
    top: `${rect.bottom + 2}px` // Add a small gap
  };
}

function updateStatusDropdownPosition() {
  if (!statusInput.value) return;
  
  const rect = statusInput.value.getBoundingClientRect();
  statusDropdownStyle.value = {
    position: 'fixed',
    zIndex: 9999,
    width: `${rect.width}px`,
    left: `${rect.left}px`,
    top: `${rect.bottom + 2}px` // Add a small gap
  };
}

function openActionDropdown() {
  showActionDropdown.value = true;
  nextTick(() => {
    updateDropdownPosition();
  });
}

function toggleStatusDropdown() {
  showStatusDropdown.value = !showStatusDropdown.value;
  if (showStatusDropdown.value) {
    nextTick(() => {
      updateStatusDropdownPosition();
    });
  }
}

// Computed values from store
const logs = computed(() => auditLogStore.formattedLogs);
const isLoading = computed(() => auditLogStore.isLoading);
const error = computed(() => auditLogStore.error);
const currentPage = computed(() => auditLogStore.currentPage);
const totalPages = computed(() => auditLogStore.totalPages);
const total = computed(() => auditLogStore.total);
const limit = computed({
  get: () => auditLogStore.limit,
  set: (value) => auditLogStore.setLimit(Number(value))
});
const filters = computed({
  get: () => auditLogStore.filters,
  set: (value) => auditLogStore.setFilters(value)
});
const groupedActions = computed(() => auditLogStore.groupedActions);

// Dropdown options
const statusOptions = [
  { value: '', label: 'All' },
  { value: 'success', label: 'Success' },
  { value: 'failure', label: 'Failure' }
];

const limitOptions = [
  { value: 10, label: '10' },
  { value: 25, label: '25' },
  { value: 50, label: '50' },
  { value: 100, label: '100' }
];

// Filter action groups based on search term
const filteredGroupedActions = computed(() => {
  if (!actionSearchTerm.value) {
    return groupedActions.value;
  }
  
  const search = actionSearchTerm.value.toLowerCase();
  const result = {};
  
  Object.entries(groupedActions.value).forEach(([module, actions]) => {
    const filteredActions = actions.filter(action => 
      action.toLowerCase().includes(search)
    );
    
    if (filteredActions.length > 0) {
      result[module] = filteredActions;
    }
  });
  
  return result;
});

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
function applyFilters() {
  auditLogStore.setFilters({ ...filters.value });
}

function resetFilters() {
  actionSearchTerm.value = '';
  auditLogStore.resetFilters();
}

function setActionFilter(action) {
  filters.value = { ...filters.value, action };
  actionSearchTerm.value = action;
  showActionDropdown.value = false;
  applyFilters();
}

function handleDropdownBlur(event) {
  // Small delay to allow for click events on dropdown items
  setTimeout(() => {
    showActionDropdown.value = false;
  }, 200);
}

function setPage(page) {
  if (page < 1 || page > totalPages.value) return;
  auditLogStore.setPage(page);
}

function setLimit(newLimit) {
  auditLogStore.setLimit(Number(newLimit));
}

function selectLog(id) {
  emit('select-log', id);
}

// Event listeners
function handleWindowResize() {
  if (showActionDropdown.value) {
    updateDropdownPosition();
  }
  if (showStatusDropdown.value) {
    updateStatusDropdownPosition();
  }
}

function handleWindowScroll() {
  if (showActionDropdown.value) {
    updateDropdownPosition();
  }
  if (showStatusDropdown.value) {
    updateStatusDropdownPosition();
  }
}

function handleClickOutside(event) {
  // Handle action dropdown
  if (actionDropdown.value && !actionDropdown.value.contains(event.target) && 
      actionInput.value && !actionInput.value.contains(event.target)) {
    showActionDropdown.value = false;
  }
  
  // Handle status dropdown
  if (statusDropdown.value && !statusDropdown.value.contains(event.target) && 
      statusInput.value && !statusInput.value.contains(event.target)) {
    showStatusDropdown.value = false;
  }
}

function setStatusFilter(value) {
  filters.value = { ...filters.value, status: value };
  showStatusDropdown.value = false;
  applyFilters();
}

// Initial data load and event listeners
onMounted(() => {
  auditLogStore.fetchAuditLogs();
  
  // Add event listeners
  document.addEventListener('mousedown', handleClickOutside);
  window.addEventListener('resize', handleWindowResize);
  window.addEventListener('scroll', handleWindowScroll);
  
  // Clean up event listener when component unmounts
  watch(() => actionSearchTerm.value, (newValue) => {
    if (!newValue) {
      // Clear action filter when search term is empty
      if (filters.value.action) {
        filters.value = { ...filters.value, action: '' };
        applyFilters();
      }
    }
  });
});

onUnmounted(() => {
  document.removeEventListener('mousedown', handleClickOutside);
  window.removeEventListener('resize', handleWindowResize);
  window.removeEventListener('scroll', handleWindowScroll);
});
</script>

<style scoped>
/* Add subtle animation for hover effects */
.hover\:bg-indigo-50 {
  transition: background-color 0.15s ease-in-out;
}
</style>