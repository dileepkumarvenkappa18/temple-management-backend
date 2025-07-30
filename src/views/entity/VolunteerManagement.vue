<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Page Header -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="py-6">
          <div class="flex items-center justify-between">
            <div>
              <h1 class="text-2xl font-bold text-gray-900">Volunteer Management</h1>
              <p class="mt-1 text-sm text-gray-600">
                Manage volunteers, assignments, and schedules for {{ templeName }}
              </p>
            </div>
            <div class="flex space-x-3">
              <button
                @click="showBulkAssignModal = true"
                class="inline-flex items-center px-4 py-2 border border-indigo-300 rounded-lg text-sm font-medium text-indigo-700 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4"></path>
                </svg>
                Bulk Assignment
              </button>
              <button
                @click="exportVolunteers"
                class="inline-flex items-center px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                </svg>
                Export
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <div class="grid grid-cols-1 lg:grid-cols-4 gap-8">
        <!-- Stats Cards -->
        <div class="lg:col-span-4">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
                    </svg>
                  </div>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-600">Total Volunteers</p>
                  <p class="text-2xl font-bold text-gray-900">{{ volunteerStats.total }}</p>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-green-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-600">Active Today</p>
                  <p class="text-2xl font-bold text-gray-900">{{ volunteerStats.activeToday }}</p>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-yellow-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                    </svg>
                  </div>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-600">Assignments</p>
                  <p class="text-2xl font-bold text-gray-900">{{ volunteerStats.assignments }}</p>
                </div>
              </div>
            </div>

            <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
              <div class="flex items-center">
                <div class="flex-shrink-0">
                  <div class="w-8 h-8 bg-purple-100 rounded-lg flex items-center justify-center">
                    <svg class="w-5 h-5 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6"></path>
                    </svg>
                  </div>
                </div>
                <div class="ml-4">
                  <p class="text-sm font-medium text-gray-600">Hours This Month</p>
                  <p class="text-2xl font-bold text-gray-900">{{ volunteerStats.hoursThisMonth }}</p>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Filters and Search -->
        <div class="lg:col-span-4">
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6 mb-6">
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <!-- Search -->
              <div class="relative">
                <div class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none">
                  <svg class="h-5 w-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z"></path>
                  </svg>
                </div>
                <input
                  v-model="searchQuery"
                  type="text"
                  placeholder="Search volunteers..."
                  class="block w-full pl-10 pr-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition-colors duration-200"
                />
              </div>

              <!-- Status Filter -->
              <select
                v-model="statusFilter"
                class="block w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition-colors duration-200"
              >
                <option value="">All Status</option>
                <option value="active">Active</option>
                <option value="inactive">Inactive</option>
                <option value="on-leave">On Leave</option>
              </select>

              <!-- Department Filter -->
              <select
                v-model="departmentFilter"
                class="block w-full px-3 py-2.5 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 sm:text-sm transition-colors duration-200"
              >
                <option value="">All Departments</option>
                <option value="seva">Seva Management</option>
                <option value="events">Events</option>
                <option value="maintenance">Maintenance</option>
                <option value="security">Security</option>
                <option value="reception">Reception</option>
              </select>

              <!-- Clear Filters -->
              <button
                @click="clearFilters"
                class="inline-flex items-center justify-center px-4 py-2.5 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-colors duration-200"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                Clear
              </button>
            </div>
          </div>
        </div>

        <!-- Volunteers List -->
        <div class="lg:col-span-4">
          <div class="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
            <div class="px-6 py-4 border-b border-gray-200">
              <h3 class="text-lg font-medium text-gray-900">Volunteers ({{ filteredVolunteers.length }})</h3>
            </div>

            <!-- Loading State -->
            <div v-if="loading" class="p-8 text-center">
              <div class="inline-flex items-center px-4 py-2 font-semibold leading-6 text-sm shadow rounded-md text-indigo-500 bg-white">
                <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-indigo-500" xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                Loading volunteers...
              </div>
            </div>

            <!-- Empty State -->
            <div v-else-if="filteredVolunteers.length === 0" class="p-8 text-center">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No volunteers found</h3>
              <p class="mt-1 text-sm text-gray-500">
                {{ searchQuery || statusFilter || departmentFilter ? 'Try adjusting your search or filters.' : 'No volunteers have joined this temple yet.' }}
              </p>
            </div>

            <!-- Volunteers Table -->
            <div v-else class="overflow-x-auto">
              <table class="min-w-full divide-y divide-gray-200">
                <thead class="bg-gray-50">
                  <tr>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Volunteer</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Contact</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Department</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Status</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Assignments</th>
                    <th class="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Last Active</th>
                    <th class="px-6 py-3 text-right text-xs font-medium text-gray-500 uppercase tracking-wider">Actions</th>
                  </tr>
                </thead>
                <tbody class="bg-white divide-y divide-gray-200">
                  <tr v-for="volunteer in paginatedVolunteers" :key="volunteer.id" class="hover:bg-gray-50 transition-colors duration-150">
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="flex items-center">
                        <div class="flex-shrink-0 h-10 w-10">
                          <img v-if="volunteer.avatar" :src="volunteer.avatar" :alt="volunteer.name" class="h-10 w-10 rounded-full object-cover" />
                          <div v-else class="h-10 w-10 rounded-full bg-indigo-100 flex items-center justify-center">
                            <span class="text-sm font-medium text-indigo-600">{{ volunteer.name.charAt(0) }}</span>
                          </div>
                        </div>
                        <div class="ml-4">
                          <div class="text-sm font-medium text-gray-900">{{ volunteer.name }}</div>
                          <div class="text-sm text-gray-500">{{ volunteer.role }}</div>
                        </div>
                      </div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <div class="text-sm text-gray-900">{{ volunteer.email }}</div>
                      <div class="text-sm text-gray-500">{{ volunteer.phone }}</div>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                        {{ volunteer.department }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap">
                      <span :class="getStatusBadgeClass(volunteer.status)" class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium">
                        {{ volunteer.status }}
                      </span>
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-900">
                      {{ volunteer.activeAssignments }}/{{ volunteer.totalAssignments }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-sm text-gray-500">
                      {{ formatDate(volunteer.lastActive) }}
                    </td>
                    <td class="px-6 py-4 whitespace-nowrap text-right text-sm font-medium">
                      <div class="flex items-center justify-end space-x-2">
                        <button
                          @click="viewVolunteer(volunteer)"
                          class="text-indigo-600 hover:text-indigo-900 transition-colors duration-150"
                          title="View Details"
                        >
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                          </svg>
                        </button>
                        <button
                          @click="editAssignments(volunteer)"
                          class="text-blue-600 hover:text-blue-900 transition-colors duration-150"
                          title="Edit Assignments"
                        >
                          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2"></path>
                          </svg>
                        </button>
                        <button
                          @click="toggleVolunteerStatus(volunteer)"
                          :class="volunteer.status === 'active' ? 'text-yellow-600 hover:text-yellow-900' : 'text-green-600 hover:text-green-900'"
                          class="transition-colors duration-150"
                          :title="volunteer.status === 'active' ? 'Deactivate' : 'Activate'"
                        >
                          <svg v-if="volunteer.status === 'active'" class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                          </svg>
                          <svg v-else class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1m4 0h1m-6 4h.01M19 10a9 9 0 11-18 0 9 9 0 0118 0z"></path>
                          </svg>
                        </button>
                      </div>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>

            <!-- Pagination -->
            <div v-if="filteredVolunteers.length > itemsPerPage" class="bg-white px-4 py-3 border-t border-gray-200 sm:px-6">
              <div class="flex items-center justify-between">
                <div class="flex-1 flex justify-between sm:hidden">
                  <button
                    @click="previousPage"
                    :disabled="currentPage === 1"
                    class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Previous
                  </button>
                  <button
                    @click="nextPage"
                    :disabled="currentPage === totalPages"
                    class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                  >
                    Next
                  </button>
                </div>
                <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
                  <div>
                    <p class="text-sm text-gray-700">
                      Showing <span class="font-medium">{{ startIndex + 1 }}</span> to <span class="font-medium">{{ endIndex }}</span> of
                      <span class="font-medium">{{ filteredVolunteers.length }}</span> results
                    </p>
                  </div>
                  <div>
                    <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
                      <button
                        @click="previousPage"
                        :disabled="currentPage === 1"
                        class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                        </svg>
                      </button>
                      <button
                        v-for="page in visiblePages"
                        :key="page"
                        @click="goToPage(page)"
                        :class="page === currentPage ? 'bg-indigo-50 border-indigo-500 text-indigo-600' : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'"
                        class="relative inline-flex items-center px-4 py-2 border text-sm font-medium"
                      >
                        {{ page }}
                      </button>
                      <button
                        @click="nextPage"
                        :disabled="currentPage === totalPages"
                        class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
                      >
                        <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                          <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                        </svg>
                      </button>
                    </nav>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Bulk Assignment Modal -->
    <div v-if="showBulkAssignModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-xl max-w-md w-full p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium text-gray-900">Bulk Assignment</h3>
          <button @click="showBulkAssignModal = false" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <div class="space-y-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Select Department</label>
            <select v-model="bulkAssignment.department" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
              <option value="">Choose department</option>
              <option value="seva">Seva Management</option>
              <option value="events">Events</option>
              <option value="maintenance">Maintenance</option>
              <option value="security">Security</option>
              <option value="reception">Reception</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Assignment Type</label>
            <select v-model="bulkAssignment.type" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
              <option value="">Choose type</option>
              <option value="daily">Daily Task</option>
              <option value="weekly">Weekly Task</option>
              <option value="event">Event Assignment</option>
              <option value="special">Special Duty</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-2">Description</label>
            <textarea v-model="bulkAssignment.description" rows="3" class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500" placeholder="Assignment details..."></textarea>
          </div>
        </div>
        <div class="flex justify-end space-x-3 mt-6">
          <button @click="showBulkAssignModal = false" class="px-4 py-2 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Cancel
          </button>
          <button @click="processBulkAssignment" class="px-4 py-2 bg-indigo-600 border border-transparent rounded-lg text-sm font-medium text-white hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
            Assign
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useVolunteerStore } from '@/stores/volunteer'
import { useToast } from '@/composables/useToast'

const route = useRoute()
const volunteerStore = useVolunteerStore()
const { showToast } = useToast()

// Reactive state
const loading = ref(false)
const searchQuery = ref('')
const statusFilter = ref('')
const departmentFilter = ref('')
const currentPage = ref(1)
const itemsPerPage = ref(10)
const showBulkAssignModal = ref(false)

// Bulk assignment data
const bulkAssignment = ref({
  department: '',
  type: '',
  description: ''
})

// Mock data - replace with actual API calls
const volunteers = ref([
  {
    id: 1,
    name: 'Rajesh Kumar',
    email: 'rajesh@email.com',
    phone: '+91 98765 43210',
    role: 'Senior Volunteer',
    department: 'seva',
    status: 'active',
    activeAssignments: 3,
    totalAssignments: 5,
    lastActive: '2024-06-17T10:30:00Z',
    avatar: null
  },
  {
    id: 2,
    name: 'Priya Sharma',
    email: 'priya@email.com',
    phone: '+91 98765 43211',
    role: 'Event Coordinator',
    department: 'events',
    status: 'active',
    activeAssignments: 2,
    totalAssignments: 3,
    lastActive: '2024-06-16T14:20:00Z',
    avatar: null
  },
  {
    id: 3,
    name: 'Amit Patel',
    email: 'amit@email.com',
    phone: '+91 98765 43212',
    role: 'Security Volunteer',
    department: 'security',
    status: 'inactive',
    activeAssignments: 0,
    totalAssignments: 2,
    lastActive: '2024-06-10T09:15:00Z',
    avatar: null
  },
  {
    id: 4,
    name: 'Sunita Devi',
    email: 'sunita@email.com',
    phone: '+91 98765 43213',
    role: 'Reception Volunteer',
    department: 'reception',
    status: 'on-leave',
    activeAssignments: 1,
    totalAssignments: 4,
    lastActive: '2024-06-15T16:45:00Z',
    avatar: null
  },
  {
    id: 5,
    name: 'Vikram Singh',
    email: 'vikram@email.com',
    phone: '+91 98765 43214',
    role: 'Maintenance Volunteer',
    department: 'maintenance',
    status: 'active',
    activeAssignments: 2,
    totalAssignments: 2,
    lastActive: '2024-06-17T08:30:00Z',
    avatar: null
  }
])

const templeName = ref('Sri Lakshmi Temple')

// Computed properties
const volunteerStats = computed(() => ({
  total: volunteers.value.length,
  activeToday: volunteers.value.filter(v => v.status === 'active').length,
  assignments: volunteers.value.reduce((sum, v) => sum + v.activeAssignments, 0),
  hoursThisMonth: 248 // Mock data
}))

const filteredVolunteers = computed(() => {
  let filtered = volunteers.value

  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase()
    filtered = filtered.filter(volunteer =>
      volunteer.name.toLowerCase().includes(query) ||
      volunteer.email.toLowerCase().includes(query) ||
      volunteer.phone.includes(query)
    )
  }

  if (statusFilter.value) {
    filtered = filtered.filter(volunteer => volunteer.status === statusFilter.value)
  }

  if (departmentFilter.value) {
    filtered = filtered.filter(volunteer => volunteer.department === departmentFilter.value)
  }

  return filtered
})

const totalPages = computed(() => Math.ceil(filteredVolunteers.value.length / itemsPerPage.value))

const paginatedVolunteers = computed(() => {
  const start = (currentPage.value - 1) * itemsPerPage.value
  const end = start + itemsPerPage.value
  return filteredVolunteers.value.slice(start, end)
})

const startIndex = computed(() => (currentPage.value - 1) * itemsPerPage.value)
const endIndex = computed(() => Math.min(startIndex.value + itemsPerPage.value, filteredVolunteers.value.length))

const visiblePages = computed(() => {
  const pages = []
  const start = Math.max(1, currentPage.value - 2)
  const end = Math.min(totalPages.value, currentPage.value + 2)
  
  for (let i = start; i <= end; i++) {
    pages.push(i)
  }
  
  return pages
})

// Methods
const getStatusBadgeClass = (status) => {
  switch (status) {
    case 'active':
      return 'bg-green-100 text-green-800'
    case 'inactive':
      return 'bg-red-100 text-red-800'
    case 'on-leave':
      return 'bg-yellow-100 text-yellow-800'
    default:
      return 'bg-gray-100 text-gray-800'
  }
}

const formatDate = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffTime = Math.abs(now - date)
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  
  if (diffDays === 1) return 'Today'
  if (diffDays === 2) return 'Yesterday'
  if (diffDays <= 7) return `${diffDays - 1} days ago`
  
  return date.toLocaleDateString('en-IN', {
    day: '2-digit',
    month: 'short',
    year: 'numeric'
  })
}

const clearFilters = () => {
  searchQuery.value = ''
  statusFilter.value = ''
  departmentFilter.value = ''
  currentPage.value = 1
}

const previousPage = () => {
  if (currentPage.value > 1) {
    currentPage.value--
  }
}

const nextPage = () => {
  if (currentPage.value < totalPages.value) {
    currentPage.value++
  }
}

const goToPage = (page) => {
  currentPage.value = page
}

const viewVolunteer = (volunteer) => {
  // Navigate to volunteer detail view or show modal
  showToast(`Viewing details for ${volunteer.name}`, 'info')
}

const editAssignments = (volunteer) => {
  // Navigate to assignment editing view or show modal
  showToast(`Editing assignments for ${volunteer.name}`, 'info')
}

const toggleVolunteerStatus = (volunteer) => {
  const newStatus = volunteer.status === 'active' ? 'inactive' : 'active'
  
  // Update volunteer status
  const index = volunteers.value.findIndex(v => v.id === volunteer.id)
  if (index !== -1) {
    volunteers.value[index].status = newStatus
    showToast(`${volunteer.name} has been ${newStatus === 'active' ? 'activated' : 'deactivated'}`, 'success')
  }
}

const processBulkAssignment = () => {
  if (!bulkAssignment.value.department || !bulkAssignment.value.type) {
    showToast('Please fill in all required fields', 'error')
    return
  }

  // Process bulk assignment logic here
  showToast('Bulk assignment processed successfully', 'success')
  showBulkAssignModal.value = false
  
  // Reset form
  bulkAssignment.value = {
    department: '',
    type: '',
    description: ''
  }
}

const exportVolunteers = () => {
  // Export functionality
  showToast('Exporting volunteer data...', 'info')
  
  // Mock export process
  setTimeout(() => {
    showToast('Volunteer data exported successfully', 'success')
  }, 2000)
}

// Lifecycle
onMounted(async () => {
  loading.value = true
  
  try {
    // Load volunteers data
    // await volunteerStore.fetchVolunteers(route.params.id)
    
    // Mock loading delay
    await new Promise(resolve => setTimeout(resolve, 1000))
  } catch (error) {
    showToast('Failed to load volunteers', 'error')
  } finally {
    loading.value = false
  }
})

// Watch for route changes
// watch(() => route.params.id, (newId) => {
//   if (newId) {
//     loadVolunteers()
//   }
// })
</script>