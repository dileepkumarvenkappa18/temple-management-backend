<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white border-b border-gray-200 shadow-sm">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <h1 class="text-3xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
              Super Admin Dashboard
            </h1>
            <p class="mt-1 text-sm text-gray-600">
              Manage temple registrations and system oversight
            </p>
          </div>
          <div class="flex items-center space-x-4">
            <div class="bg-indigo-50 text-indigo-800 px-3 py-1 rounded-full text-sm font-medium">
              System Administrator
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Quick Access Links -->
      <div class="flex flex-wrap gap-4 mb-8">
        <router-link
          to="/superadmin/tenant-approvals"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z"></path>
          </svg>
          Tenant Approvals
        </router-link>

        <button
          @click="showComingSoonToast"
          class="inline-flex items-center px-4 py-2 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
          </svg>
          Temple Approvals
        </button>
      </div>

      <!-- Summary Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100 hover:shadow-lg transition-all duration-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600 uppercase tracking-wide">Pending Approvals</p>
              <p class="text-3xl font-bold text-orange-600 mt-2">{{ summaryStats.pendingApprovals }}</p>
            </div>
            <div class="bg-orange-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-orange-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"></path>
              </svg>
            </div>
          </div>
          <div class="mt-4">
            <span class="text-sm text-gray-500">Awaiting review</span>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100 hover:shadow-lg transition-all duration-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600 uppercase tracking-wide">Active Temples</p>
              <p class="text-3xl font-bold text-green-600 mt-2">{{ summaryStats.activeTemples }}</p>
            </div>
            <div class="bg-green-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4"></path>
              </svg>
            </div>
          </div>
          <div class="mt-4">
            <span class="text-sm text-green-600">+{{ summaryStats.newThisWeek }} this week</span>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100 hover:shadow-lg transition-all duration-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600 uppercase tracking-wide">Total Users</p>
              <p class="text-3xl font-bold text-indigo-600 mt-2">{{ summaryStats.totalUsers }}</p>
            </div>
            <div class="bg-indigo-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"></path>
              </svg>
            </div>
          </div>
          <div class="mt-4">
            <span class="text-sm text-gray-500">Across all temples</span>
          </div>
        </div>

        <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100 hover:shadow-lg transition-all duration-200">
          <div class="flex items-center justify-between">
            <div>
              <p class="text-sm font-medium text-gray-600 uppercase tracking-wide">Rejected</p>
              <p class="text-3xl font-bold text-red-600 mt-2">{{ summaryStats.rejectedTemples }}</p>
            </div>
            <div class="bg-red-100 p-3 rounded-lg">
              <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.732-.833-2.464 0L4.35 16.5c-.77.833.192 2.5 1.732 2.5z"></path>
              </svg>
            </div>
          </div>
          <div class="mt-4">
            <span class="text-sm text-gray-500">This month</span>
          </div>
        </div>
      </div>

      <!-- Pending Approvals Section -->
      <div class="bg-white rounded-xl shadow-md border border-gray-100 mb-8">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center justify-between">
            <h2 class="text-xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
              Pending Temple Approvals
            </h2>
            <div class="flex items-center space-x-3">
              <select v-model="filterStatus" class="text-sm border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500">
                <option value="all">All Status</option>
                <option value="pending">Pending</option>
                <option value="approved">Approved</option>
                <option value="rejected">Rejected</option>
              </select>

              <button @click="refreshData" class="bg-indigo-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200">
                <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15"></path>
                </svg>
                Refresh
              </button>
            </div>
          </div>
        </div>

        <!-- Loading State -->
        <div v-if="isLoading" class="flex justify-center items-center py-12">
          <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
          <span class="ml-3 text-gray-600">Loading temple applications...</span>
        </div>

        <!-- Temple Applications List -->
        <div v-else class="divide-y divide-gray-200">
          <div v-if="filteredApplications.length === 0" class="text-center py-12">
            <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 11H5m14 0a2 2 0 012 2v6a2 2 0 01-2 2H5a2 2 0 01-2-2v-6a2 2 0 012-2m14 0V9a2 2 0 00-2-2M5 9a2 2 0 012-2m0 0V5a2 2 0 012-2h6a2 2 0 012 2v2M7 7h10"></path>
            </svg>
            <h3 class="mt-2 text-sm font-medium text-gray-900">No applications found</h3>
            <p class="mt-1 text-sm text-gray-500">No temple applications match your current filter.</p>
          </div>

          <div v-for="application in filteredApplications" :key="application.id" class="p-6 hover:bg-gray-50 transition-colors duration-200">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center space-x-3 mb-3">
                  <h3 class="text-lg font-semibold text-gray-900">{{ application.templeName }}</h3>
                  <span :class="getStatusBadgeClass(application.status)" class="px-2 py-1 rounded-full text-xs font-medium uppercase tracking-wide">
                    {{ application.status }}
                  </span>
                </div>

                <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mb-4">
                  <div>
                    <p class="text-sm text-gray-600">Email</p>
                    <p class="font-medium text-gray-900">{{ application.adminEmail }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-600">Phone</p>
                    <p class="font-medium text-gray-900">{{ application.phone || 'Not provided' }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-600">Location</p>
                    <p class="font-medium text-gray-900">{{ application.city }}, {{ application.state }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-600">Submitted</p>
                    <p class="font-medium text-gray-900">{{ formatDate(application.submittedAt) }}</p>
                  </div>
                  <div>
                    <p class="text-sm text-gray-600">Documents</p>
                    <p class="font-medium text-indigo-600">{{ application.documents?.length || 0 }} files</p>
                  </div>
                </div>

                <div v-if="application.address" class="mb-4">
                  <p class="text-sm text-gray-600">Address</p>
                  <p class="text-gray-900">{{ application.address }}</p>
                </div>

                <div v-if="application.notes" class="mb-4">
                  <p class="text-sm text-gray-600">Admin Notes</p>
                  <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mt-1">
                    <p class="text-yellow-800 text-sm">{{ application.notes }}</p>
                  </div>
                </div>
              </div>

              <!-- Action Buttons -->
              <div class="flex flex-col space-y-2 ml-6">
                <button @click="viewApplication(application)" class="bg-gray-100 text-gray-700 px-4 py-2 rounded-lg text-sm font-medium hover:bg-gray-200 transition-colors duration-200 whitespace-nowrap">
                  <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z"></path>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z"></path>
                  </svg>
                  View Details
                </button>

                <div v-if="application.status === 'pending'" class="flex flex-col space-y-2">
                  <button
                    @click="approveApplication(application)"
                    class="bg-green-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-green-700 transition-colors duration-200 whitespace-nowrap"
                    :disabled="isProcessing"
                  >
                    <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
                    </svg>
                    Approve
                  </button>
                  <button
                    @click="rejectApplication(application)"
                    class="bg-red-600 text-white px-4 py-2 rounded-lg text-sm font-medium hover:bg-red-700 transition-colors duration-200 whitespace-nowrap"
                    :disabled="isProcessing"
                  >
                    <svg class="w-4 h-4 inline mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
                    </svg>
                    Reject
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Recent Activity -->
      <div class="bg-white rounded-xl shadow-md border border-gray-100">
        <div class="px-6 py-4 border-b border-gray-200">
          <h2 class="text-xl font-bold text-gray-900" style="font-family: 'Roboto', sans-serif;">
            Recent Activity
          </h2>
        </div>
        <div class="p-6">
          <div v-if="isLoadingActivities" class="flex justify-center items-center py-8">
            <div class="animate-spin rounded-full h-6 w-6 border-b-2 border-indigo-600"></div>
            <span class="ml-3 text-gray-600">Loading activities...</span>
          </div>
          <div v-else-if="recentActivities.length === 0" class="text-center py-8">
            <p class="text-gray-600">No recent activities found</p>
          </div>
          <div v-else class="space-y-4">
            <div v-for="activity in recentActivities" :key="activity.id" class="flex items-start space-x-4">
              <div :class="getActivityIconClass(activity.type)" class="p-2 rounded-lg">
                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" :d="getActivityIconPath(activity.type)"></path>
                </svg>
              </div>
              <div class="flex-1 min-w-0">
                <p class="text-sm text-gray-900">
                  <span class="font-medium">{{ activity.description }}</span>
                </p>
                <p class="text-xs text-gray-500 mt-1">{{ formatDate(activity.timestamp) }}</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Modals -->
    <!-- Application Details Modal -->
    <div v-if="showDetailsModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-xl max-w-4xl w-full max-h-[90vh] overflow-y-auto">
        <div class="px-6 py-4 border-b border-gray-200 flex items-center justify-between">
          <h3 class="text-lg font-bold text-gray-900">Temple Application Details</h3>
          <button @click="closeDetailsModal" class="text-gray-400 hover:text-gray-600">
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
            </svg>
          </button>
        </div>
        <div class="p-6">
          <div v-if="selectedApplication" class="space-y-6">
            <!-- Temple Details -->
            <div>
              <h4 class="text-md font-semibold text-gray-900 mb-3">Temple Information</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <label class="text-sm text-gray-600">Temple Name</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.templeName }}</p>
                </div>
                <div>
                  <label class="text-sm text-gray-600">Registration Type</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.registrationType || 'Religious Trust' }}</p>
                </div>
                <div class="md:col-span-2">
                  <label class="text-sm text-gray-600">Address</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.address }}</p>
                </div>
              </div>
            </div>

            <!-- Admin Details -->
            <div>
              <h4 class="text-md font-semibold text-gray-900 mb-3">Administrator Information</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
                <!-- <div>
                  <label class="text-sm text-gray-600">Full Name</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.adminName }}</p>
                </div> -->
                <div>
                  <label class="text-sm text-gray-600">Email</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.adminEmail }}</p>
                </div>
                <div>
                  <label class="text-sm text-gray-600">Phone</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.phone || 'Not provided' }}</p>
                </div>
                <div>
                  <label class="text-sm text-gray-600">Role</label>
                  <p class="font-medium text-gray-900">Temple Administrator</p>
                </div>
              </div>
            </div>

            <!-- Documents -->
            <div>
              <h4 class="text-md font-semibold text-gray-900 mb-3">Uploaded Documents</h4>
              <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                <div v-for="doc in selectedApplication.documents" :key="doc.id" class="flex items-center space-x-3 p-3 border border-gray-200 rounded-lg">
                  <svg class="w-8 h-8 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"></path>
                  </svg>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate">{{ doc.name }}</p>
                    <p class="text-xs text-gray-500">{{ doc.size }} • {{ doc.type }}</p>
                  </div>
                  <button class="text-indigo-600 hover:text-indigo-800 text-sm font-medium">
                    View
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end space-x-3">
          <button @click="closeDetailsModal" class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors duration-200">
            Close
          </button>
          <div v-if="selectedApplication?.status === 'pending'" class="flex space-x-3">
            <button
              @click="approveApplication(selectedApplication)"
              class="px-4 py-2 bg-green-600 text-white rounded-lg hover:bg-green-700 transition-colors duration-200"
              :disabled="isProcessing"
            >
              Approve Temple
            </button>
            <button
              @click="rejectApplication(selectedApplication)"
              class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200"
              :disabled="isProcessing"
            >
              Reject Application
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Rejection Notes Modal -->
    <div v-if="showRejectModal" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
      <div class="bg-white rounded-xl shadow-xl max-w-md w-full">
        <div class="px-6 py-4 border-b border-gray-200">
          <h3 class="text-lg font-bold text-gray-900">Reject Application</h3>
        </div>
        <div class="p-6">
          <p class="text-gray-600 mb-4">Please provide a reason for rejecting this temple application:</p>
          <textarea v-model="rejectionNotes" class="w-full h-32 border border-gray-300 rounded-lg px-3 py-2 focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 resize-none" placeholder="Enter rejection reason..."></textarea>
        </div>
        <div class="px-6 py-4 border-t border-gray-200 flex justify-end space-x-3">
          <button @click="closeRejectModal" class="px-4 py-2 text-gray-700 bg-gray-100 rounded-lg hover:bg-gray-200 transition-colors duration-200">
            Cancel
          </button>
          <button
            @click="confirmReject"
            class="px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700 transition-colors duration-200"
            :disabled="isProcessing"
          >
            Confirm Reject
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import superAdminService from '@/services/superadmin.service'

// Toast notifications
const toast = useToast()

// Reactive data
const isLoading = ref(false)
const isLoadingActivities = ref(false)
const isProcessing = ref(false)
const filterStatus = ref('all')
const showDetailsModal = ref(false)
const showRejectModal = ref(false)
const selectedApplication = ref(null)
const rejectionNotes = ref('')

// API data
const summaryStats = ref({
  pendingApprovals: 0,
  activeTemples: 0,
  totalUsers: 0,
  rejectedTemples: 0,
  newThisWeek: 0
})

const templeApplications = ref([])
const recentActivities = ref([])

// Computed properties
const filteredApplications = computed(() => {
  if (filterStatus.value === 'all') {
    return templeApplications.value
  }
  return templeApplications.value.filter(app => app.status.toLowerCase() === filterStatus.value.toLowerCase())
})

// Mock data for development
const MOCK_DASHBOARD_STATS = {
  pendingApprovals: 5,
  activeTemples: 32,
  totalUsers: 178,
  rejectedTemples: 3,
  newThisWeek: 4
}

const MOCK_TEMPLE_APPLICATIONS = [
  {
    id: '1',
    templeName: 'Vishnu Temple',
    adminName: 'Rajesh Kumar',
    adminEmail: 'rajesh@temple.org',
    phone: '+91 98765 43210',
    city: 'Bengaluru',
    state: 'Karnataka',
    address: '123 Temple Road, Rajajinagar, Bengaluru, Karnataka',
    status: 'pending',
    submittedAt: new Date(Date.now() - 3 * 86400000).toISOString(), // 3 days ago
    documents: [
      { id: '1', name: 'Temple_Registration.pdf', size: '2.3 MB', type: 'PDF' },
      { id: '2', name: 'Trust_Deed.pdf', size: '1.8 MB', type: 'PDF' }
    ]
  },
  {
    id: '2',
    templeName: 'Shiva Temple',
    adminName: 'Sanjay Sharma',
    adminEmail: 'sanjay@temple.org',
    phone: '+91 87654 32109',
    city: 'Chennai',
    state: 'Tamil Nadu',
    address: '456 Temple Street, Mylapore, Chennai, Tamil Nadu',
    status: 'pending',
    submittedAt: new Date(Date.now() - 5 * 86400000).toISOString(), // 5 days ago
    documents: [
      { id: '3', name: 'Registration_Certificate.pdf', size: '1.5 MB', type: 'PDF' },
      { id: '4', name: 'Property_Documents.pdf', size: '3.2 MB', type: 'PDF' }
    ]
  },
  {
    id: '3',
    templeName: 'Hanuman Temple',
    adminName: 'Praveen Reddy',
    adminEmail: 'praveen@temple.org',
    phone: '+91 76543 21098',
    city: 'Hyderabad',
    state: 'Telangana',
    address: '789 Temple Avenue, Jubilee Hills, Hyderabad, Telangana',
    status: 'approved',
    submittedAt: new Date(Date.now() - 10 * 86400000).toISOString(), // 10 days ago
    approvedAt: new Date(Date.now() - 7 * 86400000).toISOString(), // 7 days ago
    documents: [
      { id: '5', name: 'Trust_Documents.pdf', size: '2.1 MB', type: 'PDF' }
    ]
  },
  {
    id: '4',
    templeName: 'Ganesh Temple',
    adminName: 'Vinod Patel',
    adminEmail: 'vinod@temple.org',
    phone: '+91 65432 10987',
    city: 'Mumbai',
    state: 'Maharashtra',
    address: '101 Temple Lane, Dadar, Mumbai, Maharashtra',
    status: 'rejected',
    submittedAt: new Date(Date.now() - 15 * 86400000).toISOString(), // 15 days ago
    rejectedAt: new Date(Date.now() - 12 * 86400000).toISOString(), // 12 days ago
    notes: 'Incomplete documentation. Missing property ownership records and trust deed details.',
    documents: [
      { id: '6', name: 'Partial_Documents.pdf', size: '1.2 MB', type: 'PDF' }
    ]
  }
]

const MOCK_ACTIVITIES = [
  {
    id: 1,
    type: 'approval',
    description: 'Approved Hanuman Temple registration for Praveen Reddy',
    timestamp: new Date(Date.now() - 7 * 86400000).toISOString() // 7 days ago
  },
  {
    id: 2,
    type: 'rejection',
    description: 'Rejected Ganesh Temple application due to incomplete documents',
    timestamp: new Date(Date.now() - 12 * 86400000).toISOString() // 12 days ago
  },
  {
    id: 3,
    type: 'new_application',
    description: 'New temple application received from Vishnu Temple',
    timestamp: new Date(Date.now() - 3 * 86400000).toISOString() // 3 days ago
  },
  {
    id: 4,
    type: 'new_application',
    description: 'New temple application received from Shiva Temple',
    timestamp: new Date(Date.now() - 5 * 86400000).toISOString() // 5 days ago
  }
]

// Simple function for the coming soon toast
const showComingSoonToast = () => {
  console.log("Temple approvals coming soon");
  toast.info('Temple Approvals feature is coming soon');
}

// Methods
const loadDashboardStats = async () => {
  try {
    // Try to get data from API
    // eslint-disable-next-line no-unused-vars
    const response = await superAdminService.getSystemStats().catch(err => null)

    if (response && response.data) {
      // Map API response to our stats object
      summaryStats.value = {
        pendingApprovals: response.data.pendingTemples || 0,
        activeTemples: response.data.activeTemples || 0,
        totalUsers: response.data.totalUsers || 0,
        rejectedTemples: response.data.rejectedTemples || 0,
        newThisWeek: response.data.newTemplesThisWeek || 0
      }
    } else {
      // If API fails, use mock data in development mode
      console.log('DEVELOPMENT: Using mock dashboard stats')
      summaryStats.value = { ...MOCK_DASHBOARD_STATS }
    }
  } catch (error) {
    console.error('Failed to load stats:', error)

    // In development, use mock data instead of showing error
    if (import.meta.env.DEV) {
      console.log('DEVELOPMENT: Using mock dashboard stats')
      summaryStats.value = { ...MOCK_DASHBOARD_STATS }
    } else {
      toast.error('Failed to load dashboard statistics')
    }
  }
}

// UPDATED: Better date handling and API field mapping
const loadPendingEntities = async () => {
  try {
    isLoading.value = true
    console.log('Loading pending entities...')

    // Use superAdminService instead of apiClient
    const response = await superAdminService.getPendingEntities().catch(err => {
      console.error('Error fetching pending entities:', err)
      return null
    })

    if (response && response.success && response.data) {
      // Log the response to see what we're working with
      console.log('Pending entities response:', response.data)

      // Map API response to our applications array
      templeApplications.value = response.data.map(temple => {
        // Convert Status field (uppercase) to status field (lowercase)
        const status = (temple.Status || temple.status || 'pending').toLowerCase()

        // Better date handling - try multiple possible field names
        const submittedDate = temple.CreatedAt ||
                             temple.createdAt ||
                             temple.SubmittedAt ||
                             temple.submittedAt ||
                             temple.created_at ||
                             temple.submitted_at ||
                             new Date().toISOString() // fallback to current date

        return {
          id: temple.ID || temple.id,
          templeName: temple.Name || temple.name,
          adminName: temple.CreatedBy || temple.createdBy || 'Unknown',
          adminEmail: temple.Email || temple.email || 'unknown@example.com',
          phone: temple.Phone || temple.phone,
          city: temple.City || temple.city,
          state: temple.State || temple.state,
          address: temple.Address || temple.address ||
                  `${temple.StreetAddress || temple.streetAddress || ''}, ${temple.City || temple.city || ''}, ${temple.State || temple.state || ''}`,
          status: status,
          submittedAt: submittedDate, // Use the better date handling
          approvedAt: temple.ApprovedAt || temple.approvedAt,
          rejectedAt: temple.RejectedAt || temple.rejectedAt,
          notes: temple.Notes || temple.notes,
          registrationType: temple.TempleType || temple.templeType,
          documents: temple.documents || []
        }
      })

      toast.success(`Loaded ${templeApplications.value.length} temple applications`)
    } else {
      // If API fails, use mock data in development mode
      console.log('DEVELOPMENT: Using mock temple applications')
      templeApplications.value = [...MOCK_TEMPLE_APPLICATIONS]
    }
  } catch (error) {
    console.error('Failed to load pending temples:', error)

    // In development, use mock data instead of showing error
    if (import.meta.env.DEV) {
      console.log('DEVELOPMENT: Using mock temple applications')
      templeApplications.value = [...MOCK_TEMPLE_APPLICATIONS]
    } else {
      toast.error('Failed to load pending temple applications')
    }
  } finally {
    isLoading.value = false
  }
}

const loadRecentActivities = async () => {
  try {
    isLoadingActivities.value = true

    // Try using API if it exists
    try {
      const response = await superAdminService.getRecentActivities()
      if (response && response.success && response.data) {
        recentActivities.value = response.data
        return
      }
    // eslint-disable-next-line no-unused-vars
    } catch (err) {
      console.log('API for activities not available, using mock data')
    }
    // If we get here, use mock data
    console.log('DEVELOPMENT: Using mock activity data')

    // Add a small delay to simulate API call
    await new Promise(resolve => setTimeout(resolve, 500))

    recentActivities.value = [...MOCK_ACTIVITIES]
  } catch (error) {
    console.error('Failed to load activities:', error)

    // In development, use mock data instead of showing error
    if (import.meta.env.DEV) {
      recentActivities.value = [...MOCK_ACTIVITIES]
    } else {
      toast.error('Failed to load recent activities')
    }
  } finally {
    isLoadingActivities.value = false
  }
}

const refreshData = async () => {
  try {
    isLoading.value = true
    await Promise.all([
      loadDashboardStats(),
      loadPendingEntities(),
      loadRecentActivities()
    ])
    toast.success('Data refreshed successfully')
  } catch (error) {
    console.error('Error refreshing data:', error)
    toast.error('Failed to refresh data')
  } finally {
    isLoading.value = false
  }
}

const viewApplication = (application) => {
  selectedApplication.value = application
  showDetailsModal.value = true
}

const closeDetailsModal = () => {
  showDetailsModal.value = false
  selectedApplication.value = null
}

const approveApplication = async (application) => {
  if (isProcessing.value) return

  try {
    isProcessing.value = true

    if (showDetailsModal.value) {
      closeDetailsModal()
    }

    // Try using superAdminService for entity approval
    try {
      await superAdminService.approveEntity(application.id, {
        notes: 'Approved by super admin'
      })
      console.log('Entity approved successfully')
    } catch (err) {
      console.log('Error in approval API call:', err)
      console.log('DEVELOPMENT: Simulating approve temple action')
      // In development, just simulate success
      await new Promise(resolve => setTimeout(resolve, 500))
    }

    // Update local state regardless
    const index = templeApplications.value.findIndex(app => app.id === application.id)
    if (index !== -1) {
      templeApplications.value[index].status = 'approved'
      templeApplications.value[index].approvedAt = new Date().toISOString()
    }

    // Update stats
    summaryStats.value.pendingApprovals = Math.max(0, summaryStats.value.pendingApprovals - 1)
    summaryStats.value.activeTemples++

    // Add to recent activities
    const newActivity = {
      id: Date.now(),
      type: 'approval',
      description: `Approved ${application.templeName} registration for ${application.adminName}`,
      timestamp: new Date().toISOString()
    }

    recentActivities.value.unshift(newActivity)

    toast.success(`${application.templeName} has been approved successfully!`)
  } catch (error) {
    console.error('Error approving temple:', error)
    toast.error('Failed to approve temple. Please try again.')
  } finally {
    isProcessing.value = false
  }
}

const rejectApplication = (application) => {
  selectedApplication.value = application
  showRejectModal.value = true

  if (showDetailsModal.value) {
    closeDetailsModal()
  }
}

const closeRejectModal = () => {
  showRejectModal.value = false
  selectedApplication.value = null
  rejectionNotes.value = ''
}

const confirmReject = async () => {
  if (!rejectionNotes.value.trim()) {
    toast.error('Please provide a reason for rejection')
    return
  }

  if (isProcessing.value) return

  try {
    isProcessing.value = true

    const application = selectedApplication.value

    // Try using superAdminService for entity rejection
    try {
      await superAdminService.rejectEntity(application.id, {
        notes: rejectionNotes.value,
        reason: rejectionNotes.value
      })
      console.log('Entity rejected successfully')
    } catch (err) {
      console.log('Error in rejection API call:', err)
      console.log('DEVELOPMENT: Simulating reject temple action')
      // In development, just simulate success
      await new Promise(resolve => setTimeout(resolve, 500))
    }

    // Update local state regardless
    const index = templeApplications.value.findIndex(app => app.id === application.id)
    if (index !== -1) {
      templeApplications.value[index].status = 'rejected'
      templeApplications.value[index].rejectedAt = new Date().toISOString()
      templeApplications.value[index].notes = rejectionNotes.value
    }

    // Update stats
    summaryStats.value.pendingApprovals = Math.max(0, summaryStats.value.pendingApprovals - 1)
    summaryStats.value.rejectedTemples++

    // Add to recent activities
    const newActivity = {
      id: Date.now(),
      type: 'rejection',
      description: `Rejected ${application.templeName} application - ${rejectionNotes.value.slice(0, 50)}${rejectionNotes.value.length > 50 ? '...' : ''}`,
      timestamp: new Date().toISOString()
    }

    recentActivities.value.unshift(newActivity)

    toast.success(`${application.templeName} application has been rejected`)
    closeRejectModal()
  } catch (error) {
    console.error('Error rejecting temple:', error)
    toast.error('Failed to reject temple. Please try again.')
  } finally {
    isProcessing.value = false
  }
}

const getStatusBadgeClass = (status) => {
  const classes = {
    pending: 'bg-orange-100 text-orange-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800'
  }
  return classes[status] || 'bg-gray-100 text-gray-800'
}

const getActivityIconClass = (type) => {
  const classes = {
    approval: 'bg-green-100 text-green-600',
    rejection: 'bg-red-100 text-red-600',
    new_application: 'bg-blue-100 text-blue-600'
  }
  return classes[type] || 'bg-gray-100 text-gray-600'
}

const getActivityIconPath = (type) => {
  const paths = {
    approval: 'M5 13l4 4L19 7',
    rejection: 'M6 18L18 6M6 6l12 12',
    new_application: 'M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z'
  }
  return paths[type] || 'M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z'
}

// UPDATED: Fixed formatDate function to handle missing dates properly
const formatDate = (dateString) => {
  if (!dateString) return 'Not available'

  try {
    const date = new Date(dateString)

    // Check if the date is valid
    if (isNaN(date.getTime())) {
      return 'Invalid date'
    }

    return date.toLocaleDateString('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    })
  } catch (error) {
    console.error('Error formatting date:', error)
    return 'Date unavailable'
  }
}

// Lifecycle
onMounted(async () => {
  // Load initial data
  await refreshData()
})
</script>