<template>
  <div class="bg-white rounded-xl shadow-md overflow-hidden">
    <!-- Header Section -->
    <div class="bg-gradient-to-r from-indigo-600 to-indigo-700 px-6 py-4">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div>
          <h2 class="text-xl font-semibold text-white">Seva Management</h2>
          <p class="text-indigo-100 text-sm mt-1">Manage and approve seva bookings</p>
        </div>
      </div>
    </div>

    <!-- Filters Section -->
    <div class="p-6 border-b border-gray-200 bg-gray-50">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Date Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Filter by Date</label>
          <input
            v-model="filters.date"
            type="date"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          />
        </div>

        <!-- Type Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Seva Type</label>
          <select
            v-model="filters.type"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          >
            <option value="">All Types</option>
            <option value="daily">Daily Sevas</option>
            <option value="special">Special Sevas</option>
            <option value="festival">Festival Sevas</option>
            <option value="personal">Personal Sevas</option>
          </select>
        </div>

        <!-- Status Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Status</label>
          <select
            v-model="filters.status"
            class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
          >
            <option value="">All Status</option>
            <option value="pending">Pending</option>
            <option value="approved">Approved</option>
            <option value="rejected">Rejected</option>
            <option value="completed">Completed</option>
          </select>
        </div>

        <!-- Search -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <div class="relative">
            <input
              v-model="filters.search"
              type="text"
              placeholder="Search devotee or seva..."
              class="w-full pl-10 pr-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 text-sm"
            />
            <svg class="absolute left-3 top-2.5 h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
            </svg>
          </div>
        </div>
      </div>

      <!-- Quick Stats -->
      <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mt-4">
        <div class="bg-white p-3 rounded-lg border border-gray-200">
          <div class="text-2xl font-bold text-indigo-600">{{ stats.total }}</div>
          <div class="text-xs text-gray-600">Total Sevas</div>
        </div>
        <div class="bg-white p-3 rounded-lg border border-gray-200">
          <div class="text-2xl font-bold text-yellow-600">{{ stats.pending }}</div>
          <div class="text-xs text-gray-600">Pending</div>
        </div>
        <div class="bg-white p-3 rounded-lg border border-gray-200">
          <div class="text-2xl font-bold text-green-600">{{ stats.approved }}</div>
          <div class="text-xs text-gray-600">Approved</div>
        </div>
        <div class="bg-white p-3 rounded-lg border border-gray-200">
          <div class="text-2xl font-bold text-blue-600">{{ stats.completed }}</div>
          <div class="text-xs text-gray-600">Completed</div>
        </div>
      </div>
    </div>

    <!-- Table Section -->
    <div class="overflow-x-auto">
      <table class="w-full">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Devotee
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Seva Details
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Date & Time
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Amount
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Status
            </th>
            <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
              Actions
            </th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr
            v-for="seva in filteredSevas"
            :key="seva.id"
            class="hover:bg-gray-50 transition-colors duration-150"
          >
            <!-- Devotee Column -->
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="flex items-center">
                <div class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
                  <span class="text-sm font-medium text-indigo-700">
                    {{ getInitials(seva.devotee?.name) || getInitials(seva.name) || 'NA' }}
                  </span>
                </div>  
                <div class="ml-3">
                  <div class="text-sm font-medium text-gray-900">{{ seva.devotee?.name || 'Not Assigned' }}</div>
                  <div class="text-sm text-gray-500">{{ seva.devotee?.phone || 'No Contact' }}</div>
                </div>
              </div>
            </td>

            <!-- Seva Details Column -->
            <td class="px-6 py-4">
              <div class="text-sm font-medium text-gray-900">{{ seva.name || 'Unnamed Seva' }}</div>
              <div class="text-sm text-gray-500">{{ seva.type || 'Not Categorized' }}</div>
              <div class="text-xs text-gray-400 mt-1" v-if="seva.notes || seva.description">
                {{ seva.notes || seva.description }}
              </div>
            </td>

            <!-- Date & Time Column -->
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm text-gray-900">{{ seva.date ? formatDate(seva.date) : 'No Date' }}</div>
              <div class="text-sm text-gray-500">{{ seva.time || 'No Time' }}</div>
            </td>

            <!-- Amount Column -->
            <td class="px-6 py-4 whitespace-nowrap">
              <div class="text-sm font-medium text-gray-900">
                ₹{{ (seva.amount || seva.price || 0).toLocaleString() }}
              </div>
            </td>

            <!-- Status Column -->
            <td class="px-6 py-4 whitespace-nowrap">
              <span 
                :class="getStatusClass(seva.status || 'pending')" 
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
              >
                <span :class="getStatusDotClass(seva.status || 'pending')" class="w-1.5 h-1.5 rounded-full mr-1.5"></span>
                {{ (seva.status || 'pending').charAt(0).toUpperCase() + (seva.status || 'pending').slice(1) }}
              </span>
            </td>

            <!-- Actions Column -->
            <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
              <div class="flex items-center gap-2">
                <!-- View Button -->
                <button
                  @click="$emit('view-seva', seva)"
                  class="text-indigo-600 hover:text-indigo-900 transition-colors duration-150"
                  title="View Details"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                </button>

                <!-- Edit Button -->
                <button
                  @click="$emit('edit-seva', seva)"
                  class="text-yellow-600 hover:text-yellow-900 transition-colors duration-150"
                  title="Edit"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
                  </svg>
                </button>

                <!-- Approve/Reject for Pending -->
                <div v-if="(seva.status || 'pending') === 'pending'" class="flex items-center gap-1">
                  <button
                    @click="$emit('approve-seva', seva)"
                    class="text-green-600 hover:text-green-900 transition-colors duration-150"
                    title="Approve"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                  </button>
                  <button
                    @click="$emit('reject-seva', seva)"
                    class="text-red-600 hover:text-red-900 transition-colors duration-150"
                    title="Reject"
                  >
                    <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                  </button>
                </div>

                <!-- Complete Button for Approved -->
                <button
                  v-if="(seva.status || 'pending') === 'approved'"
                  @click="$emit('complete-seva', seva)"
                  class="text-blue-600 hover:text-blue-900 transition-colors duration-150"
                  title="Mark Complete"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                  </svg>
                </button>

                <!-- Delete Button -->
                <button
                  @click="$emit('delete-seva', seva)"
                  class="text-red-600 hover:text-red-900 transition-colors duration-150"
                  title="Delete"
                >
                  <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16"></path>
                  </svg>
                </button>
              </div>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Empty State -->
    <div v-if="filteredSevas.length === 0" class="text-center py-12">
      <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
      </svg>
      <h3 class="mt-2 text-sm font-medium text-gray-900">No sevas found</h3>
      <p class="mt-1 text-sm text-gray-500">Get started by creating a new seva.</p>
      <div class="mt-6">
        <button
          @click="$emit('create-seva')"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200"
        >
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
          </svg>
          Add New Seva
        </button>
      </div>
    </div>

    <!-- Pagination -->
    <div v-if="filteredSevas.length > 0" class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
      <div class="flex items-center justify-between">
        <div class="flex items-center">
          <p class="text-sm text-gray-700">
            Showing
            <span class="font-medium">{{ (currentPage - 1) * itemsPerPage + 1 }}</span>
            to
            <span class="font-medium">{{ Math.min(currentPage * itemsPerPage, filteredSevas.length) }}</span>
            of
            <span class="font-medium">{{ filteredSevas.length }}</span>
            results
          </p>
        </div>
        <div class="flex items-center gap-2">
          <button
            @click="currentPage--"
            :disabled="currentPage === 1"
            class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <span class="px-3 py-1 text-sm bg-indigo-600 text-white rounded-md">{{ currentPage }}</span>
          <button
            @click="currentPage++"
            :disabled="currentPage * itemsPerPage >= filteredSevas.length"
            class="px-3 py-1 text-sm border border-gray-300 rounded-md hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  name: 'SevaList',
  emits: [
    'create-seva',
    'view-seva',
    'edit-seva',
    'approve-seva',
    'reject-seva',
    'complete-seva',
    'delete-seva'
  ],
  props: {
    sevas: {
      type: Array,
      default: () => []
    },
    loading: {
      type: Boolean,
      default: false
    }
  },
  setup(props) {
    const filters = ref({
      date: '',
      type: '',
      status: '',
      search: ''
    })

    const currentPage = ref(1)
    const itemsPerPage = ref(10)

    // Mock data if no sevas provided
    const mockSevas = ref([
      {
        id: 1,
        name: 'Abhishekam',
        type: 'special',
        devotee: {
          name: 'Rajesh Kumar',
          phone: '+91 98765 43210'
        },
        date: '2024-12-25',
        time: '06:00 AM',
        amount: 1100,
        status: 'pending',
        notes: 'Special request for Lord Ganesha'
      },
      {
        id: 2,
        name: 'Archana',
        type: 'daily',
        devotee: {
          name: 'Priya Sharma',
          phone: '+91 87654 32109'
        },
        date: '2024-12-24',
        time: '07:30 AM',
        amount: 51,
        status: 'approved',
        notes: ''
      },
      {
        id: 3,
        name: 'Satyanarayana Pooja',
        type: 'festival',
        devotee: {
          name: 'Venkatesh Rao',
          phone: '+91 76543 21098'
        },
        date: '2024-12-26',
        time: '10:00 AM',
        amount: 2100,
        status: 'completed',
        notes: 'Full moon day special'
      },
      {
        id: 4,
        name: 'Kalyanam',
        type: 'personal',
        devotee: {
          name: 'Lakshmi Devi',
          phone: '+91 65432 10987'
        },
        date: '2024-12-27',
        time: '11:30 AM',
        amount: 5100,
        status: 'rejected',
        notes: 'Anniversary celebration'
      }
    ])

    const sevasList = computed(() => props.sevas.length > 0 ? props.sevas : mockSevas.value)

    const filteredSevas = computed(() => {
      let filtered = sevasList.value || []

      if (filters.value.date) {
        filtered = filtered.filter(seva => seva.date === filters.value.date)
      }

      if (filters.value.type) {
        filtered = filtered.filter(seva => seva.type === filters.value.type)
      }

      if (filters.value.status) {
        filtered = filtered.filter(seva => seva.status === filters.value.status)
      }

      if (filters.value.search) {
        const search = filters.value.search.toLowerCase()
        filtered = filtered.filter(seva => 
          (seva.devotee?.name || '').toLowerCase().includes(search) ||
          (seva.name || '').toLowerCase().includes(search) ||
          (seva.type || '').toLowerCase().includes(search)
        )
      }

      return filtered
    })

    const stats = computed(() => {
      const list = sevasList.value || []
      const total = list.length
      const pending = list.filter(s => (s.status || 'pending') === 'pending').length
      const approved = list.filter(s => s.status === 'approved').length
      const completed = list.filter(s => s.status === 'completed').length

      return { total, pending, approved, completed }
    })

    const getInitials = (name) => {
      if (!name || typeof name !== 'string') return ''; 
      return name.split(' ').map(n => n[0]).join('').toUpperCase();
    }

    const formatDate = (dateString) => {
      if (!dateString) return 'No Date';
      try {
        const date = new Date(dateString);
        if (isNaN(date.getTime())) return 'Invalid Date';
        
        return date.toLocaleDateString('en-IN', {
          day: '2-digit',
          month: 'short',
          year: 'numeric'
        });
      } catch (error) {
        console.error('Error formatting date:', error);
        return 'Error';
      }
    }

    const getStatusClass = (status) => {
      const classes = {
        pending: 'bg-yellow-100 text-yellow-800',
        approved: 'bg-green-100 text-green-800',
        rejected: 'bg-red-100 text-red-800',
        completed: 'bg-blue-100 text-blue-800'
      }
      return classes[status] || 'bg-gray-100 text-gray-800'
    }

    const getStatusDotClass = (status) => {
      const classes = {
        pending: 'bg-yellow-400',
        approved: 'bg-green-400',
        rejected: 'bg-red-400',
        completed: 'bg-blue-400'
      }
      return classes[status] || 'bg-gray-400'
    }

    return {
      filters,
      currentPage,
      itemsPerPage,
      filteredSevas,
      stats,
      getInitials,
      formatDate,
      getStatusClass,
      getStatusDotClass
    }
  }
}
</script>