<template>
  <div class="w-full">
    <!-- Table Header with Search and Actions -->
    <!-- <div v-if="showHeader" class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mb-6"> -->
      <!-- Title -->
      <div>
        <h3 v-if="title" class="text-lg font-semibold text-gray-900">{{ title }}</h3>
        <p v-if="subtitle" class="text-sm text-gray-600 mt-1">{{ subtitle }}</p>
      </div>

      <!-- Actions Row -->
      <div class="flex flex-col sm:flex-row gap-3">
        <!-- Search -->
        <div v-if="searchable" class="relative">
          <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
            <svg class="h-5 w-5 text-gray-400" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
              <path fill-rule="evenodd" d="M8 4a4 4 0 100 8 4 4 0 000-8zM2 8a6 6 0 1110.89 3.476l4.817 4.817a1 1 0 01-1.414 1.414l-4.816-4.816A6 6 0 012 8z" clip-rule="evenodd" />
            </svg>
          </div>
          <!-- <input
            v-model="searchQuery"
            type="text"
            class="block w-full pl-10 pr-3 py-2 border border-gray-300 rounded-xl text-sm placeholder-gray-500 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            :placeholder="searchPlaceholder"
          /> -->
        </div>

        <!-- Action Buttons -->
        <div class="flex gap-2">
          <slot name="actions" />
        </div>
      </div>
    </div>

    <!-- Table Container -->
    <div class="bg-white shadow-md rounded-2xl border border-gray-200 overflow-hidden">
      <!-- Loading State -->
      <div v-if="loading" class="flex justify-center items-center py-12">
        <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
        <span class="ml-3 text-gray-600">Loading...</span>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredData.length === 0" class="text-center py-12">
        <div class="mx-auto h-12 w-12 text-gray-400 mb-4">
          <svg fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M20 13V6a2 2 0 00-2-2H6a2 2 0 00-2 2v7m16 0v5a2 2 0 01-2 2H6a2 2 0 01-2 2v-5m16 0h-2M4 13h2m0 0V9a2 2 0 012-2h2m0 0V6a2 2 0 012-2h2.586a1 1 0 01.707.293l2.414 2.414A1 1 0 0016 7.414V9" />
          </svg>
        </div>
        <h3 class="text-sm font-medium text-gray-900 mb-1">{{ emptyTitle }}</h3>
        <p class="text-sm text-gray-500">{{ emptyMessage }}</p>
      </div>

      <!-- Table -->
      <div v-else class="overflow-x-auto">
        <table class="min-w-full divide-y divide-gray-200">
          <!-- Table Head -->
          <thead class="bg-gray-50">
            <tr>
              <!-- Selection Column -->
              <th v-if="selectable" class="px-6 py-3 w-4">
                <input
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  :checked="isAllSelected"
                  :indeterminate="isIndeterminate"
                  @change="toggleSelectAll"
                />
              </th>

              <!-- Data Columns -->
              <th
                v-for="column in columns"
                :key="column.key"
                class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider"
                :class="column.headerClass || ''"
              >
                <div class="flex items-center gap-2">
                  {{ column.title }}
                  
                  <!-- Sort Icons -->
                  <div v-if="column.sortable" class="flex flex-col">
                    <button
                      class="text-gray-400 hover:text-gray-600"
                      @click="sort(column.key, 'asc')"
                    >
                      <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M14.707 12.707a1 1 0 01-1.414 0L10 9.414l-3.293 3.293a1 1 0 01-1.414-1.414l4-4a1 1 0 011.414 0l4 4a1 1 0 010 1.414z" clip-rule="evenodd" />
                      </svg>
                    </button>
                    <button
                      class="text-gray-400 hover:text-gray-600 -mt-1"
                      @click="sort(column.key, 'desc')"
                    >
                      <svg class="h-3 w-3" fill="currentColor" viewBox="0 0 20 20">
                        <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
                      </svg>
                    </button>
                  </div>
                </div>
              </th>

              <!-- Actions Column -->
              <th v-if="hasActions" class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">
                Actions
              </th>
            </tr>
          </thead>

          <!-- Table Body -->
          <tbody class="bg-white divide-y divide-gray-200">
            <tr
              v-for="(item, index) in paginatedData"
              :key="getRowKey(item, index)"
              class="hover:bg-gray-50 transition-colors duration-150"
              :class="{ 'bg-indigo-50': selectedRows.includes(getRowKey(item, index)) }"
            >
              <!-- Selection Column -->
              <td v-if="selectable" class="px-6 py-4 w-4">
                <input
                  type="checkbox"
                  class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300 rounded"
                  :checked="selectedRows.includes(getRowKey(item, index))"
                  @change="toggleRowSelection(getRowKey(item, index), item)"
                />
              </td>

              <!-- Data Columns -->
              <td
                v-for="column in columns"
                :key="column.key"
                class="px-6 py-4 text-sm"
                :class="column.cellClass || 'text-gray-900'"
              >
                <!-- Custom Cell Content -->
                <slot
                  v-if="$slots[`cell-${column.key}`]"
                  :name="`cell-${column.key}`"
                  :item="item"
                  :value="getNestedValue(item, column.key)"
                  :index="index"
                />
                
                <!-- Default Cell Content -->
                <span v-else>
                  {{ formatCellValue(getNestedValue(item, column.key), column) }}
                </span>
              </td>

              <!-- Actions Column -->
              <td v-if="hasActions" class="px-6 py-4 text-right text-sm">
                <slot name="row-actions" :item="item" :index="index" />
              </td>
            </tr>
          </tbody>
        </table>
      </div>

      <!-- Pagination -->
      <div v-if="showPagination && totalPages > 1" class="bg-white px-6 py-3 border-t border-gray-200">
        <div class="flex items-center justify-between">
          <div class="text-sm text-gray-700">
            Showing {{ startIndex + 1 }} to {{ endIndex }} of {{ totalItems }} results
          </div>
          
          <div class="flex items-center gap-2">
            <!-- Previous Button -->
            <button
              class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="currentPage === 1"
              @click="goToPage(currentPage - 1)"
            >
              Previous
            </button>

            <!-- Page Numbers -->
            <div class="flex gap-1">
              <button
                v-for="page in visiblePages"
                :key="page"
                class="px-3 py-2 text-sm font-medium rounded-lg transition-colors duration-150"
                :class="page === currentPage 
                  ? 'bg-indigo-600 text-white' 
                  : 'text-gray-700 bg-white border border-gray-300 hover:bg-gray-50'"
                @click="goToPage(page)"
              >
                {{ page }}
              </button>
            </div>

            <!-- Next Button -->
            <button
              class="px-3 py-2 text-sm font-medium text-gray-500 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              :disabled="currentPage === totalPages"
              @click="goToPage(currentPage + 1)"
            >
              Next
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Selection Summary -->
    <div v-if="selectable && selectedRows.length > 0" class="mt-4 p-4 bg-indigo-50 rounded-xl border border-indigo-200">
      <div class="flex items-center justify-between">
        <span class="text-sm text-indigo-700">
          {{ selectedRows.length }} item{{ selectedRows.length > 1 ? 's' : '' }} selected
        </span>
        <button
          class="text-sm text-indigo-600 hover:text-indigo-800 font-medium"
          @click="clearSelection"
        >
          Clear selection
        </button>
      </div>
    </div>
  <!-- </div> -->
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  data: {
    type: Array,
    default: () => []
  },
  columns: {
    type: Array,
    required: true
  },
  loading: {
    type: Boolean,
    default: false
  },
  searchable: {
    type: Boolean,
    default: true
  },
  searchPlaceholder: {
    type: String,
    default: 'Search...'
  },
  selectable: {
    type: Boolean,
    default: false
  },
  showHeader: {
    type: Boolean,
    default: true
  },
  title: {
    type: String,
    default: ''
  },
  subtitle: {
    type: String,
    default: ''
  },
  emptyTitle: {
    type: String,
    default: 'No data found'
  },
  emptyMessage: {
    type: String,
    default: 'Get started by adding your first item'
  },
  itemsPerPage: {
    type: Number,
    default: 10
  },
  showPagination: {
    type: Boolean,
    default: true
  },
  rowKey: {
    type: String,
    default: 'id'
  }
})

const emit = defineEmits(['selection-change', 'row-click', 'sort'])

// Component state
const searchQuery = ref('')
const selectedRows = ref([])
const currentPage = ref(1)
const sortKey = ref('')
const sortOrder = ref('asc')

// Computed properties
const filteredData = computed(() => {
  let result = [...props.data]
  
  // Apply search filter
  if (searchQuery.value && props.searchable) {
    const query = searchQuery.value.toLowerCase()
    result = result.filter(item => {
      return props.columns.some(column => {
        const value = getNestedValue(item, column.key)
        return String(value).toLowerCase().includes(query)
      })
    })
  }
  
  // Apply sorting
  if (sortKey.value) {
    result.sort((a, b) => {
      const aVal = getNestedValue(a, sortKey.value)
      const bVal = getNestedValue(b, sortKey.value)
      
      if (sortOrder.value === 'asc') {
        return aVal < bVal ? -1 : aVal > bVal ? 1 : 0
      } else {
        return aVal > bVal ? -1 : aVal < bVal ? 1 : 0
      }
    })
  }
  
  return result
})

const totalItems = computed(() => filteredData.value.length)
const totalPages = computed(() => Math.ceil(totalItems.value / props.itemsPerPage))

const paginatedData = computed(() => {
  if (!props.showPagination) return filteredData.value
  
  const start = (currentPage.value - 1) * props.itemsPerPage
  const end = start + props.itemsPerPage
  return filteredData.value.slice(start, end)
})

const startIndex = computed(() => (currentPage.value - 1) * props.itemsPerPage)
const endIndex = computed(() => Math.min(startIndex.value + props.itemsPerPage, totalItems.value))

const visiblePages = computed(() => {
  const pages = []
  const current = currentPage.value
  const total = totalPages.value
  
  // Show max 5 page numbers
  let start = Math.max(1, current - 2)
  let end = Math.min(total, current + 2)
  
  // Adjust if we're at the beginning or end
  if (end - start < 4) {
    if (start === 1) {
      end = Math.min(total, start + 4)
    } else {
      start = Math.max(1, end - 4)
    }
  }
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

const isAllSelected = computed(() => {
  return paginatedData.value.length > 0 && 
         paginatedData.value.every(item => selectedRows.value.includes(getRowKey(item)))
})

const isIndeterminate = computed(() => {
  const selectedInPage = paginatedData.value.filter(item => 
    selectedRows.value.includes(getRowKey(item))
  ).length
  return selectedInPage > 0 && selectedInPage < paginatedData.value.length
})

const hasActions = computed(() => !!props.$slots['row-actions'])

// Methods
const getNestedValue = (obj, path) => {
  return path.split('.').reduce((o, p) => o?.[p], obj) ?? ''
}

const getRowKey = (item, index) => {
  return item[props.rowKey] || index
}

const formatCellValue = (value, column) => {
  if (column.formatter && typeof column.formatter === 'function') {
    return column.formatter(value)
  }
  return value
}

const sort = (key, order) => {
  sortKey.value = key
  sortOrder.value = order
  emit('sort', { key, order })
}

const toggleSelectAll = () => {
  if (isAllSelected.value) {
    // Deselect all in current page
    paginatedData.value.forEach(item => {
      const key = getRowKey(item)
      const index = selectedRows.value.indexOf(key)
      if (index > -1) {
        selectedRows.value.splice(index, 1)
      }
    })
  } else {
    // Select all in current page
    paginatedData.value.forEach(item => {
      const key = getRowKey(item)
      if (!selectedRows.value.includes(key)) {
        selectedRows.value.push(key)
      }
    })
  }
  
  emit('selection-change', selectedRows.value)
}

const toggleRowSelection = (key, item) => {
  const index = selectedRows.value.indexOf(key)
  if (index > -1) {
    selectedRows.value.splice(index, 1)
  } else {
    selectedRows.value.push(key)
  }
  
  emit('selection-change', selectedRows.value)
}

const clearSelection = () => {
  selectedRows.value = []
  emit('selection-change', [])
}

const goToPage = (page) => {
  if (page >= 1 && page <= totalPages.value) {
    currentPage.value = page
  }
}

// Watchers
watch(() => props.data, () => {
  currentPage.value = 1
  selectedRows.value = []
})

watch(searchQuery, () => {
  currentPage.value = 1
})
</script>