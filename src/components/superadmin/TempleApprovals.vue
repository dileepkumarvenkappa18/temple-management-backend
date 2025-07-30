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
                    <p class="text-sm text-gray-600">Admin Name</p>
                    <p class="font-medium text-gray-900">{{ application.adminName }}</p>
                  </div>
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
                <div>
                  <label class="text-sm text-gray-600">Full Name</label>
                  <p class="font-medium text-gray-900">{{ selectedApplication.adminName }}</p>
                </div>
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
import { apiClient } from '@/plugins/axios'

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
  return templeApplications.value.filter(app => app.status === filterStatus.value)
})

// Methods
const loadDashboardStats = async () => {
  try {
    const response = await apiClient.admin.getSystemStats()
    
    if (response.data) {
      // Map API response to our stats object
      summaryStats.value = {
        pendingApprovals: response.data.pendingTemples || 0,
        activeTemples: response.data.activeTemples || 0,
        totalUsers: response.data.totalUsers || 0,
        rejectedTemples: response.data.rejectedTemples || 0,
        newThisWeek: response.data.newTemplesThisWeek || 0
      }
    }
  } catch (error) {
    console.error('Failed to load stats:', error)
    toast.error('Failed to load dashboard statistics')
  }
}

const loadPendingTemples = async () => {
  try {
    isLoading.value = true
    const response = await apiClient.admin.getPendingTemples()
    
    if (response.data) {
      // Map API response to our applications array
      templeApplications.value = response.data.map(temple => ({
        id: temple.id,
        templeName: temple.name,
        adminName: temple.adminName,
        adminEmail: temple.adminEmail,
        phone: temple.phone,
        city: temple.city,
        state: temple.state,
        address: temple.address,
        status: temple.status,
        submittedAt: temple.createdAt,
        approvedAt: temple.approvedAt,
        rejectedAt: temple.rejectedAt,
        notes: temple.notes,
        registrationType: temple.registrationType,
        documents: temple.documents || []
      }))
    }
  } catch (error) {
    console.error('Failed to load pending temples:', error)
    toast.error('Failed to load pending temple applications')
  } finally {
    isLoading.value = false
  }
}

const loadRecentActivities = async () => {
  try {
    isLoadingActivities.value = true
    // This endpoint might need to be implemented on your backend
    // For now, we'll use the current activities as a placeholder
    
    // Simulating an API call for activities
    // In a real implementation, you would replace this with an actual API call:
    // const response = await apiClient.admin.getRecentActivities()
    
    // Placeholder data - replace with actual API implementation
    setTimeout(() => {
      recentActivities.value = [
        {
          id: 1,
          type: 'approval',
          description: 'Approved Hanuman Temple registration for Suresh Reddy',
          timestamp: new Date().toISOString()
        },
        {
          id: 2,
          type: 'rejection',
          description: 'Rejected Krishna Temple application due to incomplete documents',
          timestamp: new Date(Date.now() - 86400000).toISOString() // 1 day ago
        },
        {
          id: 3,
          type: 'new_application',
          description: 'New temple application received from Ganesh Mandir',
          timestamp: new Date(Date.now() - 172800000).toISOString() // 2 days ago
        }
      ]
      isLoadingActivities.value = false
    }, 1000)
  } catch (error) {
    console.error('Failed to load activities:', error)
    toast.error('Failed to load recent activities')
    isLoadingActivities.value = false
  }
}

const refreshData = async () => {
  try {
    isLoading.value = true
    await Promise.all([
      loadDashboardStats(),
      loadPendingTemples(),
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
    
    // Call API to approve temple
    await apiClient.admin.approveTemple(application.id, 'Approved by super admin')
    
    // Update local state
    const index = templeApplications.value.findIndex(app => app.id === application.id)
    if (index !== -1) {
      templeApplications.value[index].status = 'approved'
      templeApplications.value[index].approvedAt = new Date().toISOString()
    }
    
    // Update stats
    summaryStats.value.pendingApprovals--
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
    
    // Call API to reject temple
    await apiClient.admin.rejectTemple(application.id, rejectionNotes.value)
    
    // Update local state
    const index = templeApplications.value.findIndex(app => app.id === application.id)
    if (index !== -1) {
      templeApplications.value[index].status = 'rejected'
      templeApplications.value[index].rejectedAt = new Date().toISOString()
      templeApplications.value[index].notes = rejectionNotes.value
    }
    
    // Update stats
    summaryStats.value.pendingApprovals--
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

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  
  const date = new Date(dateString)
  return date.toLocaleDateString('en-IN', {
    year: 'numeric',
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

// Lifecycle
onMounted(async () => {
  // Load initial data
  await refreshData()
})
</script>