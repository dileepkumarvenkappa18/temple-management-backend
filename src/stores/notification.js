import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notification', () => {
  // State
  const notifications = ref([])
  const unreadCount = ref(0)
  const isLoading = ref(false)
  const lastFetch = ref(null)

  // Notification types based on the temple SaaS platform
  const NOTIFICATION_TYPES = {
    // Temple Admin notifications
    TEMPLE_APPROVED: 'temple_approved',
    TEMPLE_REJECTED: 'temple_rejected',
    TEMPLE_PENDING: 'temple_pending',
    NEW_DEVOTEE_REGISTRATION: 'new_devotee_registration',
    NEW_SEVA_BOOKING: 'new_seva_booking',
    NEW_DONATION: 'new_donation',
    EVENT_RSVP: 'event_rsvp',
    
    // Devotee notifications
    PROFILE_INCOMPLETE: 'profile_incomplete',
    SEVA_CONFIRMATION: 'seva_confirmation',
    SEVA_REMINDER: 'seva_reminder',
    DONATION_RECEIPT: 'donation_receipt',
    EVENT_INVITATION: 'event_invitation',
    EVENT_REMINDER: 'event_reminder',
    TEMPLE_ANNOUNCEMENT: 'temple_announcement',
    
    // Volunteer notifications
    ASSIGNMENT_CREATED: 'assignment_created',
    ASSIGNMENT_UPDATED: 'assignment_updated',
    SCHEDULE_CHANGED: 'schedule_changed',
    
    // Super Admin notifications
    NEW_TENANT_REGISTRATION: 'new_tenant_registration',
    TEMPLE_APPROVAL_REQUIRED: 'temple_approval_required',
    
    // System notifications
    SYSTEM_MAINTENANCE: 'system_maintenance',
    SYSTEM_UPDATE: 'system_update',
    WELCOME: 'welcome',
    SUCCESS: 'success',
    ERROR: 'error',
    WARNING: 'warning',
    INFO: 'info'
  }

  // Priority levels
  const PRIORITY_LEVELS = {
    LOW: 'low',
    MEDIUM: 'medium',
    HIGH: 'high',
    URGENT: 'urgent'
  }

  // Notification categories
  const CATEGORIES = {
    TEMPLE: 'temple',
    SEVA: 'seva',
    DONATION: 'donation',
    EVENT: 'event',
    PROFILE: 'profile',
    SYSTEM: 'system',
    APPROVAL: 'approval',
    COMMUNICATION: 'communication'
  }

  // Getters
  const unreadNotifications = computed(() => 
    notifications.value.filter(n => !n.isRead)
  )

  const recentNotifications = computed(() => 
    notifications.value
      .filter(n => !n.isRead)
      .sort((a, b) => new Date(b.createdAt) - new Date(a.createdAt))
      .slice(0, 5)
  )

  const notificationsByCategory = computed(() => {
    const grouped = {}
    notifications.value.forEach(notification => {
      const category = notification.category || 'general'
      if (!grouped[category]) {
        grouped[category] = []
      }
      grouped[category].push(notification)
    })
    return grouped
  })

  const highPriorityNotifications = computed(() =>
    notifications.value.filter(n => 
      n.priority === PRIORITY_LEVELS.HIGH || n.priority === PRIORITY_LEVELS.URGENT
    )
  )

  // Actions
  const addNotification = (notification) => {
    const newNotification = {
      id: Date.now() + Math.random(),
      title: notification.title || 'Notification',
      message: notification.message || '',
      type: notification.type || NOTIFICATION_TYPES.INFO,
      category: notification.category || CATEGORIES.SYSTEM,
      priority: notification.priority || PRIORITY_LEVELS.MEDIUM,
      isRead: false,
      isPersistent: notification.isPersistent || false,
      autoClose: notification.autoClose !== false, // Default to true
      duration: notification.duration || 5000,
      createdAt: new Date().toISOString(),
      readAt: null,
      userId: notification.userId || null,
      entityId: notification.entityId || null,
      actionUrl: notification.actionUrl || null,
      actionText: notification.actionText || null,
      metadata: notification.metadata || {},
      icon: notification.icon || getDefaultIcon(notification.type),
      color: notification.color || getDefaultColor(notification.type)
    }

    notifications.value.unshift(newNotification)
    if (!newNotification.isRead) {
      unreadCount.value++
    }

    return newNotification.id
  }

  const markAsRead = (notificationId) => {
    const notification = notifications.value.find(n => n.id === notificationId)
    if (notification && !notification.isRead) {
      notification.isRead = true
      notification.readAt = new Date().toISOString()
      unreadCount.value = Math.max(0, unreadCount.value - 1)
    }
  }

  const markAllAsRead = () => {
    notifications.value.forEach(notification => {
      if (!notification.isRead) {
        notification.isRead = true
        notification.readAt = new Date().toISOString()
      }
    })
    unreadCount.value = 0
  }

  const removeNotification = (notificationId) => {
    const index = notifications.value.findIndex(n => n.id === notificationId)
    if (index > -1) {
      const notification = notifications.value[index]
      if (!notification.isRead) {
        unreadCount.value = Math.max(0, unreadCount.value - 1)
      }
      notifications.value.splice(index, 1)
    }
  }

  const clearAll = () => {
    notifications.value = []
    unreadCount.value = 0
  }

  const clearRead = () => {
    notifications.value = notifications.value.filter(n => !n.isRead)
  }

  // Toast notifications (temporary)
  const showToast = (message, type = 'info', options = {}) => {
    return addNotification({
      title: options.title || getToastTitle(type),
      message,
      type: type,
      category: CATEGORIES.SYSTEM,
      priority: options.priority || PRIORITY_LEVELS.MEDIUM,
      autoClose: options.autoClose !== false,
      duration: options.duration || getToastDuration(type),
      isPersistent: false,
      ...options
    })
  }

  const showSuccess = (message, options = {}) => 
    showToast(message, NOTIFICATION_TYPES.SUCCESS, options)

  const showError = (message, options = {}) => 
    showToast(message, NOTIFICATION_TYPES.ERROR, { 
      duration: 8000, 
      priority: PRIORITY_LEVELS.HIGH,
      ...options 
    })

  const showWarning = (message, options = {}) => 
    showToast(message, NOTIFICATION_TYPES.WARNING, options)

  const showInfo = (message, options = {}) => 
    showToast(message, NOTIFICATION_TYPES.INFO, options)

  // Role-specific notification helpers
  const showTempleApprovalNotification = (templeData, isApproved) => {
    return addNotification({
      title: isApproved ? 'Temple Approved!' : 'Temple Rejected',
      message: isApproved 
        ? `Your temple "${templeData.name}" has been approved. You can now start managing your temple.`
        : `Your temple "${templeData.name}" was rejected. Please check the notes and resubmit.`,
      type: isApproved ? NOTIFICATION_TYPES.TEMPLE_APPROVED : NOTIFICATION_TYPES.TEMPLE_REJECTED,
      category: CATEGORIES.TEMPLE,
      priority: PRIORITY_LEVELS.HIGH,
      isPersistent: true,
      autoClose: false,
      actionUrl: isApproved ? '/tenant/dashboard' : '/tenant/entities/create',
      actionText: isApproved ? 'View Dashboard' : 'Resubmit Temple',
      metadata: { templeId: templeData.id }
    })
  }

  const showSevaBookingNotification = (sevaData, devoteeData) => {
    return addNotification({
      title: 'New Seva Booking',
      message: `${devoteeData.name} has booked ${sevaData.name} seva.`,
      type: NOTIFICATION_TYPES.NEW_SEVA_BOOKING,
      category: CATEGORIES.SEVA,
      priority: PRIORITY_LEVELS.MEDIUM,
      isPersistent: true,
      actionUrl: `/entity/${sevaData.entityId}/sevas`,
      actionText: 'View Seva',
      metadata: { sevaId: sevaData.id, devoteeId: devoteeData.id }
    })
  }

  const showDonationNotification = (donationData) => {
    return addNotification({
      title: 'Donation Received',
      message: `New donation of ₹${donationData.amount} received from ${donationData.donorName}.`,
      type: NOTIFICATION_TYPES.NEW_DONATION,
      category: CATEGORIES.DONATION,
      priority: PRIORITY_LEVELS.MEDIUM,
      isPersistent: true,
      actionUrl: `/entity/${donationData.entityId}/donations`,
      actionText: 'View Donation',
      metadata: { donationId: donationData.id }
    })
  }

  const showProfileIncompleteReminder = (devoteeData) => {
    return addNotification({
      title: 'Complete Your Profile',
      message: 'Please complete your profile to enjoy full temple services.',
      type: NOTIFICATION_TYPES.PROFILE_INCOMPLETE,
      category: CATEGORIES.PROFILE,
      priority: PRIORITY_LEVELS.MEDIUM,
      isPersistent: true,
      actionUrl: '/profile/complete',
      actionText: 'Complete Profile',
      metadata: { devoteeId: devoteeData.id }
    })
  }

  const showEventInvitation = (eventData) => {
    return addNotification({
      title: 'Event Invitation',
      message: `You're invited to ${eventData.name} on ${new Date(eventData.date).toLocaleDateString()}.`,
      type: NOTIFICATION_TYPES.EVENT_INVITATION,
      category: CATEGORIES.EVENT,
      priority: PRIORITY_LEVELS.MEDIUM,
      isPersistent: true,
      actionUrl: `/entity/${eventData.entityId}/events/${eventData.id}`,
      actionText: 'View Event',
      metadata: { eventId: eventData.id }
    })
  }

  // Utility functions
  const getDefaultIcon = (type) => {
    const iconMap = {
      [NOTIFICATION_TYPES.SUCCESS]: 'CheckCircle',
      [NOTIFICATION_TYPES.ERROR]: 'XCircle',
      [NOTIFICATION_TYPES.WARNING]: 'AlertTriangle',
      [NOTIFICATION_TYPES.INFO]: 'Info',
      [NOTIFICATION_TYPES.TEMPLE_APPROVED]: 'CheckCircle',
      [NOTIFICATION_TYPES.TEMPLE_REJECTED]: 'XCircle',
      [NOTIFICATION_TYPES.NEW_SEVA_BOOKING]: 'Calendar',
      [NOTIFICATION_TYPES.NEW_DONATION]: 'CreditCard',
      [NOTIFICATION_TYPES.EVENT_INVITATION]: 'Calendar',
      [NOTIFICATION_TYPES.PROFILE_INCOMPLETE]: 'User',
      [NOTIFICATION_TYPES.ASSIGNMENT_CREATED]: 'UserCheck',
      [NOTIFICATION_TYPES.SYSTEM_MAINTENANCE]: 'Settings'
    }
    return iconMap[type] || 'Bell'
  }

  const getDefaultColor = (type) => {
    const colorMap = {
      [NOTIFICATION_TYPES.SUCCESS]: 'text-green-600 bg-green-50',
      [NOTIFICATION_TYPES.ERROR]: 'text-red-600 bg-red-50',
      [NOTIFICATION_TYPES.WARNING]: 'text-yellow-600 bg-yellow-50',
      [NOTIFICATION_TYPES.INFO]: 'text-indigo-600 bg-indigo-50',
      [NOTIFICATION_TYPES.TEMPLE_APPROVED]: 'text-green-600 bg-green-50',
      [NOTIFICATION_TYPES.TEMPLE_REJECTED]: 'text-red-600 bg-red-50',
      [NOTIFICATION_TYPES.NEW_SEVA_BOOKING]: 'text-indigo-600 bg-indigo-50',
      [NOTIFICATION_TYPES.NEW_DONATION]: 'text-purple-600 bg-purple-50',
      [NOTIFICATION_TYPES.EVENT_INVITATION]: 'text-blue-600 bg-blue-50'
    }
    return colorMap[type] || 'text-indigo-600 bg-indigo-50'
  }

  const getToastTitle = (type) => {
    const titleMap = {
      [NOTIFICATION_TYPES.SUCCESS]: 'Success',
      [NOTIFICATION_TYPES.ERROR]: 'Error',
      [NOTIFICATION_TYPES.WARNING]: 'Warning',
      [NOTIFICATION_TYPES.INFO]: 'Information'
    }
    return titleMap[type] || 'Notification'
  }

  const getToastDuration = (type) => {
    const durationMap = {
      [NOTIFICATION_TYPES.SUCCESS]: 4000,
      [NOTIFICATION_TYPES.ERROR]: 8000,
      [NOTIFICATION_TYPES.WARNING]: 6000,
      [NOTIFICATION_TYPES.INFO]: 5000
    }
    return durationMap[type] || 5000
  }

  // Fetch notifications from API (mock implementation)
  const fetchNotifications = async (userId, entityId = null) => {
    isLoading.value = true
    try {
      // Mock API call - replace with actual API
      const mockNotifications = [
        {
          title: 'Welcome to Temple Management',
          message: 'Complete your profile to get started.',
          type: NOTIFICATION_TYPES.WELCOME,
          category: CATEGORIES.SYSTEM,
          priority: PRIORITY_LEVELS.MEDIUM,
          isPersistent: true,
          actionUrl: '/profile/complete',
          actionText: 'Complete Profile'
        }
      ]

      // Clear existing notifications and add new ones
      notifications.value = []
      unreadCount.value = 0
      
      mockNotifications.forEach(notification => {
        addNotification(notification)
      })

      lastFetch.value = new Date().toISOString()
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
      showError('Failed to load notifications')
    } finally {
      isLoading.value = false
    }
  }

  // Auto-cleanup old notifications
  const cleanupOldNotifications = (daysOld = 30) => {
    const cutoffDate = new Date()
    cutoffDate.setDate(cutoffDate.getDate() - daysOld)
    
    const initialCount = notifications.value.length
    notifications.value = notifications.value.filter(notification => {
      const notificationDate = new Date(notification.createdAt)
      return notificationDate > cutoffDate || notification.isPersistent
    })
    
    // Recalculate unread count
    unreadCount.value = notifications.value.filter(n => !n.isRead).length
    
    const removedCount = initialCount - notifications.value.length
    if (removedCount > 0) {
      console.log(`Cleaned up ${removedCount} old notifications`)
    }
  }

  return {
    // State
    notifications,
    unreadCount,
    isLoading,
    lastFetch,
    
    // Constants
    NOTIFICATION_TYPES,
    PRIORITY_LEVELS,
    CATEGORIES,
    
    // Getters
    unreadNotifications,
    recentNotifications,
    notificationsByCategory,
    highPriorityNotifications,
    
    // Actions
    addNotification,
    markAsRead,
    markAllAsRead,
    removeNotification,
    clearAll,
    clearRead,
    
    // Toast methods
    showToast,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    
    // Role-specific notifications
    showTempleApprovalNotification,
    showSevaBookingNotification,
    showDonationNotification,
    showProfileIncompleteReminder,
    showEventInvitation,
    
    // API methods
    fetchNotifications,
    cleanupOldNotifications
  }
})