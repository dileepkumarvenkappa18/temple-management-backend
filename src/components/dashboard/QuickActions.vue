<template>
  <div class="bg-white rounded-xl shadow-md p-6 border border-gray-100">
    <h3 class="text-lg font-semibold text-gray-900 mb-6 flex items-center">
      <svg class="w-5 h-5 text-indigo-600 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 10V3L4 14h7v7l9-11h-7z"></path>
      </svg>
      Quick Actions
    </h3>

    <!-- Tenant Quick Actions -->
    <div v-if="userRole === 'tenant'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <QuickActionButton
        title="Create Temple"
        description="Add a new temple"
        icon="temple"
        color="indigo"
        @click="handleAction('create-temple')"
      />
      <QuickActionButton
        title="View Applications"
        description="Check temple status"
        icon="document"
        color="blue"
        @click="handleAction('view-applications')"
      />
      <QuickActionButton
        title="Upload Documents"
        description="Submit required docs"
        icon="upload"
        color="green"
        @click="handleAction('upload-documents')"
      />
      <QuickActionButton
        title="Profile Settings"
        description="Update your profile"
        icon="user"
        color="purple"
        @click="handleAction('profile-settings')"
      />
      <QuickActionButton
        title="Help & Support"
        description="Get assistance"
        icon="help"
        color="orange"
        @click="handleAction('help-support')"
      />
    </div>

    <!-- Entity Admin Quick Actions -->
    <div v-else-if="userRole === 'entity'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
      <QuickActionButton
        title="Add Devotee"
        description="Register new devotee"
        icon="user-add"
        color="indigo"
        @click="handleAction('add-devotee')"
      />
      <QuickActionButton
        title="Create Seva"
        description="Add new seva service"
        icon="seva"
        color="purple"
        @click="handleAction('create-seva')"
      />
      <QuickActionButton
        title="New Event"
        description="Schedule temple event"
        icon="calendar"
        color="orange"
        @click="handleAction('create-event')"
      />
      <QuickActionButton
        title="Send Notification"
        description="Notify devotees"
        icon="notification"
        color="blue"
        @click="handleAction('send-notification')"
      />
      <QuickActionButton
        title="Generate Report"
        description="Export data reports"
        icon="report"
        color="green"
        @click="handleAction('generate-report')"
      />
      <QuickActionButton
        title="Manage Volunteers"
        description="View assignments"
        icon="volunteer"
        color="cyan"
        @click="handleAction('manage-volunteers')"
      />
      <QuickActionButton
        title="Donation Receipt"
        description="Generate receipts"
        icon="receipt"
        color="emerald"
        @click="handleAction('donation-receipt')"
      />
      <QuickActionButton
        title="Temple Settings"
        description="Configure temple"
        icon="settings"
        color="gray"
        @click="handleAction('temple-settings')"
      />
    </div>

    <!-- Devotee Quick Actions -->
    <div v-else-if="userRole === 'devotee'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <QuickActionButton
        title="Book Seva"
        description="Schedule seva booking"
        icon="seva"
        color="purple"
        @click="handleAction('book-seva')"
      />
      <QuickActionButton
        title="Make Donation"
        description="Contribute to temple"
        icon="donation"
        color="green"
        @click="handleAction('make-donation')"
      />
      <QuickActionButton
        title="View Events"
        description="Upcoming events"
        icon="calendar"
        color="orange"
        @click="handleAction('view-events')"
      />
      <QuickActionButton
        title="Update Profile"
        description="Complete your profile"
        icon="profile"
        color="indigo"
        :badge="profileCompletion < 100 ? `${profileCompletion}%` : null"
        @click="handleAction('update-profile')"
      />
      <QuickActionButton
        title="Donation History"
        description="View past donations"
        icon="history"
        color="blue"
        @click="handleAction('donation-history')"
      />
      <QuickActionButton
        title="Family Members"
        description="Manage family info"
        icon="family"
        color="pink"
        @click="handleAction('family-members')"
      />
    </div>

    <!-- Volunteer Quick Actions -->
    <div v-else-if="userRole === 'volunteer'" class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-4">
      <QuickActionButton
        title="My Assignments"
        description="View seva assignments"
        icon="assignment"
        color="indigo"
        :badge="pendingAssignments > 0 ? pendingAssignments.toString() : null"
        @click="handleAction('my-assignments')"
      />
      <QuickActionButton
        title="My Schedule"
        description="Volunteer calendar"
        icon="schedule"
        color="blue"
        @click="handleAction('my-schedule')"
      />
      <QuickActionButton
        title="Update Availability"
        description="Set your availability"
        icon="availability"
        color="green"
        @click="handleAction('update-availability')"
      />
      <QuickActionButton
        title="Report Activity"
        description="Log seva completion"
        icon="report"
        color="purple"
        @click="handleAction('report-activity')"
      />
      <QuickActionButton
        title="Temple Events"
        description="Upcoming events"
        icon="calendar"
        color="orange"
        @click="handleAction('temple-events')"
      />
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'
import { useAuth } from '@/composables/useAuth'
import { useRouter } from 'vue-router'
import QuickActionButton from './QuickActionButton.vue'

// Composables
const { user } = useAuth()
const router = useRouter()

// Props (for customization)
const props = defineProps({
  entityId: {
    type: String,
    default: null
  }
})

// Computed properties
const userRole = computed(() => user.value?.role)

// Mock data for badges/counts
const profileCompletion = computed(() => {
  // This would come from user profile data
  return user.value?.profileCompletion || 65
})

const pendingAssignments = computed(() => {
  // This would come from volunteer assignments
  return 3
})

// Methods
const handleAction = (actionType) => {
  console.log('Quick action clicked:', actionType)
  
  // Route mapping based on action type
  const routeMap = {
    // Tenant actions
    'create-temple': '/tenant/entities/create',
    'view-applications': '/tenant/entities',
    'upload-documents': '/tenant/entities/upload',
    'profile-settings': '/tenant/profile',
    'help-support': '/help',

    // Entity admin actions
    'add-devotee': `/entity/${props.entityId}/devotees/create`,
    'create-seva': `/entity/${props.entityId}/sevas/create`,
    'create-event': `/entity/${props.entityId}/events/create`,
    'send-notification': `/entity/${props.entityId}/communication`,
    'generate-report': `/entity/${props.entityId}/reports`,
    'manage-volunteers': `/entity/${props.entityId}/volunteers`,
    'donation-receipt': `/entity/${props.entityId}/donations/receipts`,
    'temple-settings': `/entity/${props.entityId}/settings`,

    // Devotee actions
    'book-seva': `/entity/${props.entityId}/devotee/seva-booking`,
    'make-donation': `/entity/${props.entityId}/devotee/donate`,
    'view-events': `/entity/${props.entityId}/devotee/events`,
    'update-profile': `/entity/${props.entityId}/devotee/profile`,
    'donation-history': `/entity/${props.entityId}/devotee/donations`,
    'family-members': `/entity/${props.entityId}/devotee/family`,

    // Volunteer actions
    'my-assignments': `/entity/${props.entityId}/volunteer/assignments`,
    'my-schedule': `/entity/${props.entityId}/volunteer/schedule`,
    'update-availability': `/entity/${props.entityId}/volunteer/availability`,
    'report-activity': `/entity/${props.entityId}/volunteer/report`,
    'temple-events': `/entity/${props.entityId}/volunteer/events`
  }

  const route = routeMap[actionType]
  if (route) {
    router.push(route)
  } else {
    // Handle custom actions or modals
    handleCustomAction(actionType)
  }
}

const handleCustomAction = (actionType) => {
  // Handle actions that don't require navigation
  // Like opening modals, triggering downloads, etc.
  switch (actionType) {
    case 'generate-report':
      // Open report generation modal
      break
    case 'send-notification':
      // Open notification composer modal
      break
    default:
      console.log('Custom action not implemented:', actionType)
  }
}
</script>