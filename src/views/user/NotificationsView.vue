<template>
  <div class="notifications-page min-h-screen bg-gradient-to-br from-indigo-50 via-white to-purple-50 py-8 px-4 sm:px-6 lg:px-8">
    <div class="max-w-4xl mx-auto">
      <!-- Header Section -->
      <div class="mb-8">
        <div class="text-center sm:text-left">
          <h1 class="text-4xl font-bold bg-gradient-to-r from-indigo-800 to-purple-700 bg-clip-text text-transparent mb-3">
            Notifications
          </h1>
          <p class="text-gray-600 text-lg font-medium">Manage all your notifications in one place</p>
        </div>
      </div>

      <!-- Stats Cards -->
      <div class="grid grid-cols-1 sm:grid-cols-3 gap-6 mb-8">
        <div class="bg-white rounded-2xl shadow-lg border border-indigo-100 p-6 transform hover:scale-105 transition-all duration-300">
          <div class="flex items-center">
            <div class="p-3 rounded-xl bg-indigo-100">
              <Bell class="h-6 w-6 text-indigo-600" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total</p>
              <p class="text-2xl font-bold text-gray-900">{{ notifications.length }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-lg border border-amber-100 p-6 transform hover:scale-105 transition-all duration-300">
          <div class="flex items-center">
            <div class="p-3 rounded-xl bg-amber-100">
              <AlertCircle class="h-6 w-6 text-amber-600" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Unread</p>
              <p class="text-2xl font-bold text-gray-900">{{ unreadCount }}</p>
            </div>
          </div>
        </div>

        <div class="bg-white rounded-2xl shadow-lg border border-green-100 p-6 transform hover:scale-105 transition-all duration-300">
          <div class="flex items-center">
            <div class="p-3 rounded-xl bg-green-100">
              <CheckCircle class="h-6 w-6 text-green-600" />
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Read</p>
              <p class="text-2xl font-bold text-gray-900">{{ notifications.length - unreadCount }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Notification Actions -->
      <div class="flex flex-col sm:flex-row justify-between items-start sm:items-center gap-4 mb-8">
        <div class="flex flex-wrap gap-3">
          <BaseButton
            v-if="unreadCount > 0"
            @click="markAllAsRead"
            variant="secondary"
            size="sm"
            class="inline-flex items-center bg-indigo-100 text-indigo-700 hover:bg-indigo-200 border-0 rounded-xl px-6 py-2 font-semibold transition-all duration-200 transform hover:scale-105"
          >
            <CheckCircle class="h-4 w-4 mr-2" />
            Mark all as read
          </BaseButton>
          <BaseButton
            v-if="notifications.length > 0"
            @click="clearRead"
            variant="secondary"
            size="sm"
            class="inline-flex items-center bg-amber-100 text-amber-700 hover:bg-amber-200 border-0 rounded-xl px-6 py-2 font-semibold transition-all duration-200 transform hover:scale-105"
          >
            <Trash2 class="h-4 w-4 mr-2" />
            Clear read notifications
          </BaseButton>
        </div>
        <div class="w-full sm:w-auto">
          <BaseSelect
            v-model="currentFilter"
            :options="filterOptions"
            placeholder="Filter notifications"
            class="min-w-[200px] rounded-xl border-indigo-200 focus:border-indigo-500 focus:ring-indigo-500"
          />
        </div>
      </div>

      <!-- Custom Enhanced Loading State -->
      <div v-if="isLoading" class="flex justify-center items-center py-16">
        <div class="relative">
          <div class="w-16 h-16 border-4 border-indigo-200 border-t-indigo-600 rounded-full animate-spin"></div>
          <div class="absolute inset-0 flex items-center justify-center">
            <Bell class="h-6 w-6 text-indigo-600" />
          </div>
        </div>
      </div>

      <!-- Empty State -->
      <div v-else-if="filteredNotifications.length === 0" class="bg-white rounded-2xl shadow-xl border border-gray-100 p-12 text-center my-8 transform hover:scale-105 transition-all duration-300">
        <div class="flex justify-center mb-6">
          <div class="relative">
            <div class="h-20 w-20 rounded-full bg-gradient-to-br from-indigo-100 to-purple-100 flex items-center justify-center">
              <Bell class="h-10 w-10 text-indigo-600" />
            </div>
            <div class="absolute -top-1 -right-1 h-6 w-6 rounded-full bg-gradient-to-r from-indigo-500 to-purple-500 flex items-center justify-center">
              <Check class="h-3 w-3 text-white" />
            </div>
          </div>
        </div>
        <h3 class="text-2xl font-bold text-gray-900 mb-3">No notifications</h3>
        <p class="text-gray-500 text-lg max-w-md mx-auto leading-relaxed">
          {{ emptyStateMessage }}
        </p>
      </div>

      <!-- Notification List -->
      <div v-else class="space-y-6">
        <div
          v-for="notification in filteredNotifications"
          :key="notification.id"
          :class="[
            'group bg-white rounded-2xl shadow-lg overflow-hidden transition-all duration-300 transform hover:scale-102 hover:shadow-xl',
            notification.isRead
              ? 'border border-gray-100 opacity-80 hover:opacity-100'
              : 'border border-indigo-100 border-l-4 border-l-indigo-500 shadow-indigo-100'
          ]"
        >
          <!-- notification content goes here -->
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useNotificationStore } from '@/stores/notification'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseLoader from '@/components/common/BaseLoader.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'

import {
  Bell,
  Check,
  X,
  ArrowRight,
  AlertCircle,
  CheckCircle,
  Trash2,
  Clock
} from 'lucide-vue-next'
import * as lucideIcons from 'lucide-vue-next'

const notificationStore = useNotificationStore()
const {
  notifications,
  unreadCount,
  isLoading,
  markAsRead,
  markAllAsRead,
  removeNotification,
  clearRead,
  fetchNotifications,
  NOTIFICATION_TYPES,
  CATEGORIES
} = notificationStore

const currentFilter = ref('all')
const filterOptions = [
  { value: 'all', label: 'All Notifications' },
  { value: 'unread', label: 'Unread Only' },
  { value: 'temple', label: 'Temple Updates' },
  { value: 'seva', label: 'Seva Bookings' },
  { value: 'donation', label: 'Donations' },
  { value: 'event', label: 'Events' },
  { value: 'system', label: 'System Updates' }
]

const filteredNotifications = computed(() => {
  if (currentFilter.value === 'all') return notifications
  if (currentFilter.value === 'unread') return notifications.filter(n => !n.isRead)
  return notifications.filter(n => n.category === currentFilter.value)
})

const emptyStateMessage = computed(() => {
  if (currentFilter.value === 'all') return 'You have no notifications at the moment...'
  if (currentFilter.value === 'unread') return 'All caught up! You have no unread notifications.'
  const map = {
    temple: 'temple updates',
    seva: 'seva booking',
    donation: 'donation',
    event: 'event',
    system: 'system'
  }
  return `You have no ${map[currentFilter.value] || currentFilter.value} notifications.`
})

onMounted(() => fetchNotifications('current-user-id'))
</script>
