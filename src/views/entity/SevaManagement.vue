<!-- 
  Key changes to the SevaManagement.vue file to display data from seva_bookings table:
  1. Updated loadSevas() function to use sevaService.getEntityBookings() instead of sevaStore.fetchSevas()
  2. Adjusted the UI and data display to handle the seva_bookings structure
  3. Maintained mapping of date to availability_schedule for the form
  4. Fixed field inconsistencies and added better error handling
-->

<template>
  <div class="min-h-screen bg-gray-50/90">
    <!-- Page Header -->
    <div class="bg-white border-b border-gray-200 shadow-sm">
      <div class="px-6 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-2xl font-bold text-gray-900">Seva Management</h1>
            <p class="mt-1 text-sm text-gray-500">
              Manage seva bookings, approvals, and scheduling
            </p>
          </div>
          <div class="flex gap-3">
            <button
              @click="showCreateForm = true"
              class="inline-flex items-center px-4 py-2.5 bg-indigo-600 text-white text-sm font-medium rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 shadow-sm"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
              </svg>
              Create Seva
            </button>
            <!-- <button
              @click="exportSevas"
              class="inline-flex items-center px-4 py-2.5 bg-white text-gray-700 text-sm font-medium rounded-lg border border-gray-300 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 shadow-sm"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              Export
            </button> -->
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Overview -->
    <div class="px-6 py-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-all duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-indigo-100 rounded-lg">
              <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Total Bookings</p>
              <p class="text-2xl font-bold text-gray-900">{{ sevaList.length }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-all duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-green-100 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Approved</p>
              <p class="text-2xl font-bold text-gray-900">
                {{ sevaList.filter(s => s.status === 'approved').length }}
              </p>
            </div>
          </div>
        </div>

        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-all duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-yellow-100 rounded-lg">
              <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Pending</p>
              <p class="text-2xl font-bold text-gray-900">
                {{ sevaList.filter(s => s.status === 'pending').length }}
              </p>
            </div>
          </div>
        </div>

        <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-all duration-200">
          <div class="flex items-center">
            <div class="p-3 bg-red-100 rounded-lg">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 14l2-2m0 0l2-2m-2 2l-2-2m2 2l2 2m7-2a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-500">Rejected</p>
              <p class="text-2xl font-bold text-gray-900">
                {{ sevaList.filter(s => s.status === 'rejected').length }}
              </p>
            </div>
          </div>
        </div>
      </div>

      <!-- Filters Section -->
      <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 mb-6 hover:shadow-md transition-all duration-200">
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Search Bookings</label>
            <div class="relative">
              <input
                v-model="searchFilter"
                type="text"
                placeholder="Search by name, devotee, or type..."
                class="w-full pl-10 pr-4 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
              />
              <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
                </svg>
              </div>
            </div>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Status</label>
            <select
              v-model="statusFilter"
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
              <option value="">All Status</option>
              <option value="pending">Pending</option>
              <option value="approved">Approved</option>
              <option value="rejected">Rejected</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Date Range</label>
            <select
              v-model="dateRangeFilter"
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
              <option value="">All Dates</option>
              <option value="today">Today</option>
              <option value="week">This Week</option>
              <option value="month">This Month</option>
            </select>
          </div>

          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Seva Type</label>
            <select
              v-model="typeFilter"
              class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
            >
              <option value="">All Types</option>
              <option value="daily">Daily Sevas</option>
              <option value="special">Special Sevas</option>
              <option value="festival">Festival Sevas</option>
              <option value="personal">Personal Sevas</option>
            </select>
          </div>
        </div>

        <div class="flex justify-between items-center mt-4">
          <div class="text-sm text-gray-500">
            Showing {{ filteredSevas.length }} of {{ sevaList.length }} bookings
          </div>
          <button
            @click="clearFilters"
            class="text-sm text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200"
          >
            Clear Filters
          </button>
        </div>
      </div>

      <!-- Loading indicator -->
      <div v-if="loading" class="flex justify-center items-center py-12">
        <svg class="animate-spin h-8 w-8 text-indigo-600" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      </div>

      <!-- Seva List Component -->
      <div v-else class="bg-white rounded-xl shadow-sm border border-gray-200 hover:shadow-md transition-all duration-200">
        <div class="bg-white rounded-xl shadow-md overflow-hidden">
          <!-- Table Section -->
          <div v-if="filteredSevas.length > 0" class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50">
                <tr>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Booking Details
                  </th>
                  <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">
                    Devotee
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
                  v-for="booking in filteredSevas"
                  :key="booking.id"
                  class="hover:bg-gray-50 transition-colors duration-150"
                >
                  <!-- Booking Details Column -->
                  <td class="px-6 py-4">
                    <div class="text-sm font-medium text-gray-900">{{ booking.seva?.name || 'Unnamed Seva' }}</div>
                    <div class="text-sm text-gray-500">{{ booking.seva?.type || 'Not Categorized' }}</div>
                    <div class="text-xs text-gray-400 mt-1" v-if="booking.notes">
                      {{ booking.notes }}
                    </div>
                  </td>

                  <!-- Devotee Column -->
                  <td class="px-6 py-4">
                    <div class="text-sm font-medium text-gray-900">{{ booking.devotee_name || 'Anonymous' }}</div>
                    <div class="text-xs text-gray-500" v-if="booking.devotee_contact">
                      {{ booking.devotee_contact }}
                    </div>
                  </td>

                  <!-- Date & Time Column -->
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div v-if="booking.scheduled_date" class="text-sm text-gray-900">{{ formatDate(booking.scheduled_date) }}</div>
                    <div v-else class="text-sm text-gray-900">No date specified</div>
                    <div class="text-sm text-gray-500">
                      {{ booking.start_time || '' }} {{ booking.start_time && booking.end_time ? '-' : '' }} {{ booking.end_time || '' }}
                    </div>
                  </td>

                  <!-- Amount Column -->
                  <td class="px-6 py-4 whitespace-nowrap">
                    <div class="text-sm font-medium text-gray-900">
                      ₹{{ (booking.amount || (booking.seva?.price || 0)).toLocaleString() }}
                    </div>
                  </td>

                  <!-- Status Column -->
                  <td class="px-6 py-4 whitespace-nowrap">
                    <span 
                      :class="getStatusClass(booking.status || 'pending')" 
                      class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                    >
                      <span :class="getStatusDotClass(booking.status || 'pending')" class="w-1.5 h-1.5 rounded-full mr-1.5"></span>
                      {{ (booking.status || 'pending').charAt(0).toUpperCase() + (booking.status || 'pending').slice(1) }}
                    </span>
                  </td>

                  <!-- Actions Column -->
                  <td class="px-6 py-4 whitespace-nowrap text-sm font-medium">
                    <div class="flex items-center gap-2">
                      <!-- View Button -->
                      <button
                        @click="viewSevaDetails(booking)"
                        class="text-indigo-600 hover:text-indigo-900 transition-colors duration-150"
                        title="View Details"
                      >
                        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                        </svg>
                      </button>

                      <!-- Approve/Reject for Pending -->
                      <template v-if="booking.status === 'pending'">
                        <button
                          @click="approveSeva(booking)"
                          class="text-green-600 hover:text-green-900 transition-colors duration-150"
                          title="Approve"
                        >
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                          </svg>
                        </button>
                        <button
                          @click="rejectSeva(booking)"
                          class="text-red-600 hover:text-red-900 transition-colors duration-150"
                          title="Reject"
                        >
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                          </svg>
                        </button>
                      </template>

                      <!-- Complete Button for Approved -->
                      <button
                        v-if="booking.status === 'approved'"
                        @click="completeSeva(booking)"
                        class="text-blue-600 hover:text-blue-900 transition-colors duration-150"
                        title="Mark Complete"
                      >
                        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
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
            <h3 class="mt-2 text-sm font-medium text-gray-900">No bookings found</h3>
            <p class="mt-1 text-sm text-gray-500">Get started by creating a new seva.</p>
            <div class="mt-6">
              <button
                @click="showCreateForm = true"
                class="inline-flex items-center px-4 py-2.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 shadow-sm"
              >
                <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Add New Seva
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create/Edit Seva Modal -->
    <div
      v-if="showCreateForm || editingSeva"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 backdrop-blur-sm"
    >
      <div class="bg-white rounded-xl shadow-xl max-w-2xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 class="text-xl font-bold text-gray-900">
            {{ editingSeva ? 'Edit Seva' : 'Create New Seva' }}
          </h2>
          <button
            @click="closeForm"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors duration-200"
          >
            <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <!-- Seva Form -->
        <div class="p-6">
          <div class="space-y-5">
            <!-- Basic Details Section -->
            <div class="pb-3 mb-3 border-b border-gray-200">
              <h3 class="text-md font-semibold text-gray-700 mb-3">Basic Details</h3>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Seva Name*</label>
                  <input
                    v-model="sevaForm.name"
                    type="text"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    placeholder="Enter seva name"
                  />
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Seva Type*</label>
                  <select
                    v-model="sevaForm.type"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                  >
                    <option value="" disabled>Select Seva Type</option>
                    <option value="daily">Daily Seva</option>
                    <option value="special">Special Seva</option>
                    <option value="festival">Festival Seva</option>
                    <option value="personal">Personal Seva</option>
                  </select>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Description*</label>
                  <textarea
                    v-model="sevaForm.description"
                    rows="3"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    placeholder="Add any additional details about this seva"
                  ></textarea>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Price (₹)*</label>
                  <input
                    v-model="sevaForm.price"
                    type="number"
                    min="0"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    placeholder="Enter amount in ₹"
                  />
                </div>
              </div>
            </div>

            <!-- Scheduling Section -->
            <div class="pb-3 mb-3 border-b border-gray-200">
              <h3 class="text-md font-semibold text-gray-700 mb-3">Scheduling</h3>
              <div class="space-y-4">
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Date*</label>
                  <input
                    v-model="sevaForm.availability_schedule"
                    type="date"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                  />
                </div>
                
                <!-- Time Section with Start and End Time -->
                <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">Start Time</label>
                    <input
                      v-model="sevaForm.startTime"
                      type="time"
                      class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    />
                  </div>
                  <div>
                    <label class="block text-sm font-medium text-gray-700 mb-1">End Time</label>
                    <input
                      v-model="sevaForm.endTime"
                      type="time"
                      class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    />
                  </div>
                </div>
                
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Duration (minutes)</label>
                  <input
                    v-model="sevaForm.duration"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    placeholder="e.g. 30, 60, 120"
                  />
                </div>
                
                <!-- Maximum Slots Per Day Field -->
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Maximum Bookings Per Day</label>
                  <input
                    v-model="sevaForm.max_bookings_per_day"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
                    placeholder="Enter maximum slots available per day"
                  />
                </div>
              </div>
            </div>
          </div>
          
          <!-- Form Actions -->
          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="closeForm"
              class="px-4 py-2.5 border border-gray-300 text-gray-700 rounded-lg hover:bg-gray-50 transition-all duration-200"
            >
              Cancel
            </button>
            <button
              type="button"
              @click="saveSeva"
              :disabled="formLoading"
              class="px-4 py-2.5 bg-indigo-600 text-white rounded-lg hover:bg-indigo-700 transition-all duration-200 shadow-sm"
            >
              <span v-if="formLoading" class="flex items-center">
                <svg class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Processing...
              </span>
              <span v-else>
                {{ editingSeva ? 'Update Seva' : 'Create Seva' }}
              </span>
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Seva Details Modal -->
    <div
      v-if="selectedSeva"
      class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50 backdrop-blur-sm"
    >
      <div class="bg-white rounded-xl shadow-xl max-w-3xl w-full max-h-[90vh] overflow-y-auto">
        <div class="flex items-center justify-between p-6 border-b border-gray-200">
          <h2 class="text-xl font-bold text-gray-900">Booking Details</h2>
          <button
            @click="selectedSeva = null"
            class="p-2 hover:bg-gray-100 rounded-lg transition-colors duration-200"
          >
            <svg class="w-5 h-5 text-gray-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>
        
        <div class="p-6">
          <div class="bg-gray-50 rounded-lg p-6 mb-6 border border-gray-100">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-gray-900">{{ selectedSeva.seva?.name || 'Unnamed Seva' }}</h3>
              <span 
                :class="getStatusClass(selectedSeva.status || 'pending')" 
                class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
              >
                <span :class="getStatusDotClass(selectedSeva.status || 'pending')" class="w-1.5 h-1.5 rounded-full mr-1.5"></span>
                {{ (selectedSeva.status || 'pending').charAt(0).toUpperCase() + (selectedSeva.status || 'pending').slice(1) }}
              </span>
            </div>
            
            <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-4">
              <div>
                <p class="text-sm font-medium text-gray-500">Seva Type</p>
                <p class="text-base text-gray-900">{{ selectedSeva.seva?.type || 'Not specified' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Devotee</p>
                <p class="text-base text-gray-900">{{ selectedSeva.devotee_name || 'Anonymous' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Contact</p>
                <p class="text-base text-gray-900">{{ selectedSeva.devotee_contact || 'Not provided' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Amount</p>
                <p class="text-base text-gray-900">₹{{ (selectedSeva.amount || (selectedSeva.seva?.price || 0)).toLocaleString() }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Date</p>
                <p class="text-base text-gray-900">{{ selectedSeva.scheduled_date ? formatDate(selectedSeva.scheduled_date) : 'Not specified' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Time</p>
                <p class="text-base text-gray-900">
                  {{ selectedSeva.start_time || 'Not specified' }} 
                  {{ selectedSeva.start_time && selectedSeva.end_time ? ' - ' : '' }} 
                  {{ selectedSeva.end_time || '' }}
                </p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Booking ID</p>
                <p class="text-base text-gray-900">{{ selectedSeva.id || 'Not available' }}</p>
              </div>
              <div>
                <p class="text-sm font-medium text-gray-500">Booking Date</p>
                <p class="text-base text-gray-900">{{ selectedSeva.created_at ? formatDate(selectedSeva.created_at) : 'Not available' }}</p>
              </div>
            </div>
            
            <div v-if="selectedSeva.notes" class="mb-4">
              <p class="text-sm font-medium text-gray-500">Notes</p>
              <p class="text-base text-gray-900">{{ selectedSeva.notes }}</p>
            </div>
          </div>
          
          <!-- Action buttons -->
          <div v-if="selectedSeva.status === 'pending'" class="flex justify-end space-x-3">
            <button
              @click="rejectSeva(selectedSeva)"
              class="px-4 py-2.5 bg-white border border-red-300 text-red-700 rounded-lg hover:bg-red-50 transition-all duration-200"
            >
              Reject Booking
            </button>
            <button
              @click="approveSeva(selectedSeva)"
              class="px-4 py-2.5 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-all duration-200 shadow-sm"
            >
              Approve Booking
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Toast Notifications -->
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="fixed top-4 right-4 z-50 transform transition-all duration-300"
      :class="toast.show ? 'translate-x-0 opacity-100' : 'translate-x-full opacity-0'"
    >
      <div
        class="px-4 py-3 rounded-lg shadow-lg border-l-4 max-w-sm"
        :class="{
          'bg-green-50 border-green-400 text-green-800': toast.type === 'success',
          'bg-red-50 border-red-400 text-red-800': toast.type === 'error',
          'bg-yellow-50 border-yellow-400 text-yellow-800': toast.type === 'warning',
          'bg-blue-50 border-blue-400 text-blue-800': toast.type === 'info'
        }"
      >
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg v-if="toast.type === 'success'" class="w-5 h-5 text-green-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
            </svg>
            <svg v-else-if="toast.type === 'error'" class="w-5 h-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm font-medium">{{ toast.message }}</p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useSevaStore } from '@/stores/seva'
import { sevaService } from '@/services/seva.service'

// Get entity ID from route
const route = useRoute()
const entityId = route.params.id

// Initialize stores
const sevaStore = useSevaStore()

// Local state
const loading = ref(false)
const formLoading = ref(false)
const showCreateForm = ref(false)
const editingSeva = ref(null)
const selectedSeva = ref(null)
const toasts = ref([])
const nextToastId = ref(1)
const sevaList = ref([])

// Filters - directly on component
const searchFilter = ref('')
const statusFilter = ref('')
const typeFilter = ref('')
const dateRangeFilter = ref('')

// Computed filtered sevas
const filteredSevas = computed(() => {
  let filtered = [...sevaList.value]
  
  if (searchFilter.value) {
    const search = searchFilter.value.toLowerCase()
    filtered = filtered.filter(booking => 
      (booking.seva?.name || '').toLowerCase().includes(search) ||
      (booking.notes || '').toLowerCase().includes(search) ||
      (booking.devotee_name || '').toLowerCase().includes(search) ||
      (booking.seva?.type || '').toLowerCase().includes(search)
    )
  }
  
  if (statusFilter.value) {
    filtered = filtered.filter(booking => booking.status === statusFilter.value)
  }
  
  if (typeFilter.value) {
    filtered = filtered.filter(booking => booking.seva?.type === typeFilter.value)
  }
  
  // Add date range filtering
  if (dateRangeFilter.value) {
    const today = new Date()
    today.setHours(0, 0, 0, 0)
    
    const getBookingDate = (booking) => {
      if (!booking.scheduled_date) return null
      return new Date(booking.scheduled_date)
    }
    
    if (dateRangeFilter.value === 'today') {
      filtered = filtered.filter(booking => {
        const bookingDate = getBookingDate(booking)
        if (!bookingDate) return false
        return bookingDate.getTime() === today.getTime()
      })
    } else if (dateRangeFilter.value === 'week') {
      const weekStart = new Date(today)
      weekStart.setDate(today.getDate() - today.getDay())
      
      const weekEnd = new Date(weekStart)
      weekEnd.setDate(weekStart.getDate() + 6)
      
      filtered = filtered.filter(booking => {
        const bookingDate = getBookingDate(booking)
        if (!bookingDate) return false
        return bookingDate >= weekStart && bookingDate <= weekEnd
      })
    } else if (dateRangeFilter.value === 'month') {
      const monthStart = new Date(today.getFullYear(), today.getMonth(), 1)
      const monthEnd = new Date(today.getFullYear(), today.getMonth() + 1, 0)
      
      filtered = filtered.filter(booking => {
        const bookingDate = getBookingDate(booking)
        if (!bookingDate) return false
        return bookingDate >= monthStart && bookingDate <= monthEnd
      })
    }
  }
  
  return filtered
})

// Form data for create/edit with new fields - UPDATED to use availability_schedule instead of date
const sevaForm = ref({
  name: '',
  type: '',
  description: '',
  price: 0,
  duration: 30,
  availability_schedule: '', // CHANGED from date to availability_schedule
  startTime: '',
  endTime: '',
  max_bookings_per_day: 10,
  is_active: true
})

// Toast function
const showToast = (message, type = 'info') => {
  const id = nextToastId.value++
  const toast = { id, message, type, show: true }
  toasts.value.push(toast)
  
  // Auto-hide after 3 seconds
  setTimeout(() => {
    const index = toasts.value.findIndex(t => t.id === id)
    if (index !== -1) {
      toasts.value[index].show = false
      setTimeout(() => {
        toasts.value = toasts.value.filter(t => t.id !== id)
      }, 300)
    }
  }, 3000)
  
  return id
}

// Helper functions
const getInitials = (name) => {
  if (!name) return 'NA'
  return name.split(' ').map(n => n[0]).join('').toUpperCase()
}

const formatDate = (dateString) => {
  if (!dateString) return 'No Date'
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    })
  } catch (error) {
    return dateString
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

// API Integration Functions - UPDATED to use sevaService.getEntityBookings
const loadSevas = async () => {
  console.log('Loading seva bookings for entity ID:', entityId)
  try {
    loading.value = true
    const result = await sevaService.getEntityBookings(entityId)

    // ✅ Ensure sevaList is always an array
    const bookings = result?.data
    sevaList.value = Array.isArray(bookings) ? bookings : []

    console.log('Seva bookings loaded:', sevaList.value)
  } catch (error) {
    console.error('Error loading seva bookings:', error)
    showToast('Failed to load seva bookings. Please try again.', 'error')
  } finally {
    loading.value = false
  }
}


// Action handlers
const approveSeva = async (booking) => {
  console.log('Approving booking:', booking.id)
  
  try {
    const result = await sevaService.updateBookingStatus(booking.id, 'approved')
    if (result.success) {
      await loadSevas() // Reload the list
      showToast(`Booking for "${booking.seva?.name}" approved successfully`, 'success')
      
      // Close detail view if open
      if (selectedSeva.value && selectedSeva.value.id === booking.id) {
        selectedSeva.value = null
      }
    } else {
      showToast(result.message || 'Failed to approve booking', 'error')
    }
  } catch (error) {
    console.error('Error approving booking:', error)
    showToast('Failed to approve booking. Please try again.', 'error')
  }
}

const rejectSeva = async (booking) => {
  console.log('Rejecting booking:', booking.id)
  
  try {
    const result = await sevaService.updateBookingStatus(booking.id, 'rejected')
    if (result.success) {
      await loadSevas() // Reload the list
      showToast(`Booking for "${booking.seva?.name}" rejected`, 'warning')
      
      // Close detail view if open
      if (selectedSeva.value && selectedSeva.value.id === booking.id) {
        selectedSeva.value = null
      }
    } else {
      showToast(result.message || 'Failed to reject booking', 'error')
    }
  } catch (error) {
    console.error('Error rejecting booking:', error)
    showToast('Failed to reject booking. Please try again.', 'error')
  }
}

const completeSeva = async (booking) => {
  console.log('Completing booking:', booking.id)
  
  try {
    const result = await sevaService.updateBookingStatus(booking.id, 'completed')
    if (result.success) {
      await loadSevas() // Reload the list
      showToast(`Booking for "${booking.seva?.name}" marked as completed`, 'success')
    } else {
      showToast(result.message || 'Failed to mark booking as completed', 'error')
    }
  } catch (error) {
    console.error('Error completing booking:', error)
    showToast('Failed to mark booking as completed. Please try again.', 'error')
  }
}

// Keep the seva creation functionality intact - UPDATED to use availability_schedule
const editSeva = (seva) => {
  console.log('Editing seva:', seva)
  
  // Populate the form with the seva data - using availability_schedule
  sevaForm.value = {
    name: seva.name,
    type: seva.type,
    description: seva.description || '',
    price: seva.price || 0,
    duration: seva.duration || 30,
    availability_schedule: seva.availability_schedule || '', // Using consistent field name
    startTime: seva.startTime || '',
    endTime: seva.endTime || '',
    max_bookings_per_day: seva.max_bookings_per_day || 10,
    is_active: seva.is_active !== undefined ? seva.is_active : true
  }
  
  // Store the original for updating
  editingSeva.value = seva
}

// UPDATED saveSeva function with improved validation and error handling
const saveSeva = async () => {
  // Validate form
  if (!sevaForm.value.name || !sevaForm.value.type || !sevaForm.value.description) {
    showToast('Please fill in all required fields', 'error')
    return
  }
  
  // Validate entity ID
  if (!entityId) {
    showToast('Missing entity ID. Please try again.', 'error')
    return
  }
  
  formLoading.value = true
  
  try {
    // Check if we have any existing sevas to inspect
    console.log("🔍 Checking for existing sevas:", sevaList.value.length > 0 ? "Found some" : "None available")
    
    // Include all scheduling fields in the seva_catalog structure
   const payload = {
  entity_id: parseInt(entityId),
  name: sevaForm.value.name,
  seva_type: sevaForm.value.type,
  description: sevaForm.value.description,
  price: parseFloat(sevaForm.value.price),
  duration: parseInt(sevaForm.value.duration),
  availability_schedule: sevaForm.value.availability_schedule,
  start_time: sevaForm.value.startTime,
  end_time: sevaForm.value.endTime,
  max_bookings_per_day: parseInt(sevaForm.value.max_bookings_per_day),
  status: "active",
  is_active: true
}
    
    console.log("🔍 COMPLETE PAYLOAD WITH SCHEDULING:", payload)
    
    if (editingSeva.value) {
      console.log("Updating seva ID:", editingSeva.value.id)
      
      const result = await sevaService.updateSeva(editingSeva.value.id, payload)
      if (result.success) {
        await loadSevas()
        showToast(`Seva "${sevaForm.value.name}" updated successfully`, 'success')
        closeForm()
      } else {
        showToast(result.message || 'Failed to update seva', 'error')
      }
    } else {
      console.log("🚀 Creating new seva with complete scheduling fields")
      
      try {
        console.log("Before sevaService.createSeva call with payload:", payload)
        const result = await sevaService.createSeva(payload)
        console.log("After sevaService.createSeva call, result:", result)
        
        // Detailed inspection of result
        if (result) {
          console.log("🔎 Result success:", result.success)
          console.log("🔎 Result message:", result.message)
          console.log("🔎 Result error:", result.error)
          console.log("🔎 Result errors detail:", result.errors)
          
          if (result.success) {
            await loadSevas()
            showToast(`Seva "${sevaForm.value.name}" created successfully`, 'success')
            closeForm()
          } else {
            // More specific error handling
            if (result.errors && Object.keys(result.errors).length > 0) {
              // Show the first validation error
              const firstErrorField = Object.keys(result.errors)[0]
              const errorMessage = result.errors[firstErrorField]
              showToast(`Validation error: ${firstErrorField} - ${errorMessage}`, 'error')
            } else {
              showToast(result.message || result.error || 'Failed to create seva', 'error')
            }
          }
        } else {
          showToast('Unexpected response from server', 'error')
        }
      } catch (serviceError) {
        console.error("❌ SERVICE ERROR:", serviceError)
        
        // Try to extract more details
        if (serviceError.response) {
          console.error("Response status:", serviceError.response.status)
          console.error("Response data:", serviceError.response.data)
          
          // Try to log the raw response
          try {
            console.error("Raw response:", JSON.stringify(serviceError.response))
          } catch (e) {
            console.error("Could not stringify response")
          }
          
          showToast(`API Error (${serviceError.response.status}): ${serviceError.response.data?.message || 'Server error'}`, 'error')
        } else {
          showToast(`Error: ${serviceError.message}`, 'error')
        }
      }
    }
  } catch (error) {
    console.error('❌ TOP LEVEL ERROR:', error)
    if (error.response) {
      console.error('Response data:', error.response.data)
      console.error('Response status:', error.response.status)
      showToast(`Error (${error.response.status}): ${error.response.data?.message || 'Failed to save seva'}`, 'error')
    } else {
      showToast('Failed to save seva. Please try again.', 'error')
    }
  } finally {
    formLoading.value = false
  }
}

const viewSevaDetails = (booking) => {
  console.log('Viewing booking details:', booking)
  selectedSeva.value = booking
}

// UPDATED to use availability_schedule
const closeForm = () => {
  showCreateForm.value = false
  editingSeva.value = null
  
  // Reset form with consistent field names
  sevaForm.value = {
    name: '',
    type: '',
    description: '',
    price: 0,
    duration: 30,
    availability_schedule: '',
    startTime: '',
    endTime: '',
    max_bookings_per_day: 10,
    is_active: true
  }
}

const clearFilters = () => {
  searchFilter.value = ''
  statusFilter.value = ''
  typeFilter.value = ''
  dateRangeFilter.value = ''
}

const exportSevas = () => {
  showToast('Export feature is not implemented yet', 'info')
}

// Lifecycle
onMounted(() => {
  console.log('SevaManagement component mounted')
  loadSevas()
})
</script>

<style scoped>
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>