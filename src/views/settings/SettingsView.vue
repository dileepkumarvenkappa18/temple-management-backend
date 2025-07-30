<template>
  <div class="settings-container">
    <h1 class="text-2xl font-bold mb-6">Account Settings</h1>
    
    <BaseCard class="mb-6">
      <template #header>
        <h2 class="text-xl font-semibold">Profile Information</h2>
      </template>
      <div class="p-4">
        <form @submit.prevent="updateProfile">
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <BaseInput
              v-model="profile.fullName"
              label="Full Name"
              placeholder="Enter your full name"
              required
            />
            <BaseInput
              v-model="profile.email"
              label="Email Address"
              type="email"
              placeholder="Enter your email"
              disabled
            />
            <BaseInput
              v-model="profile.phone"
              label="Phone Number"
              placeholder="Enter your phone number"
            />
            <BaseSelect
              v-model="profile.language"
              label="Preferred Language"
              :options="languageOptions"
            />
          </div>
          <div class="mt-4">
            <BaseButton type="submit" variant="primary">
              Save Changes
            </BaseButton>
          </div>
        </form>
      </div>
    </BaseCard>

    <BaseCard class="mb-6">
      <template #header>
        <h2 class="text-xl font-semibold">Change Password</h2>
      </template>
      <div class="p-4">
        <form @submit.prevent="changePassword">
          <div class="grid grid-cols-1 gap-4">
            <BaseInput
              v-model="passwordForm.currentPassword"
              label="Current Password"
              type="password"
              placeholder="Enter your current password"
              required
            />
            <BaseInput
              v-model="passwordForm.newPassword"
              label="New Password"
              type="password"
              placeholder="Enter your new password"
              required
            />
            <BaseInput
              v-model="passwordForm.confirmPassword"
              label="Confirm New Password"
              type="password"
              placeholder="Confirm your new password"
              required
            />
          </div>
          <div class="mt-4">
            <BaseButton type="submit" variant="primary">
              Update Password
            </BaseButton>
          </div>
        </form>
      </div>
    </BaseCard>

    <BaseCard>
      <template #header>
        <h2 class="text-xl font-semibold">Notification Settings</h2>
      </template>
      <div class="p-4">
        <div class="space-y-4">
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-medium">Email Notifications</h3>
              <p class="text-gray-600 text-sm">Receive updates and alerts via email</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="notifications.email" class="sr-only peer">
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-medium">SMS Notifications</h3>
              <p class="text-gray-600 text-sm">Receive important alerts via SMS</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="notifications.sms" class="sr-only peer">
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
            </label>
          </div>
          
          <div class="flex items-center justify-between">
            <div>
              <h3 class="font-medium">Event Reminders</h3>
              <p class="text-gray-600 text-sm">Receive reminders about upcoming events</p>
            </div>
            <label class="relative inline-flex items-center cursor-pointer">
              <input type="checkbox" v-model="notifications.events" class="sr-only peer">
              <div class="w-11 h-6 bg-gray-200 peer-focus:outline-none peer-focus:ring-4 peer-focus:ring-indigo-300 rounded-full peer peer-checked:after:translate-x-full peer-checked:after:border-white after:content-[''] after:absolute after:top-[2px] after:left-[2px] after:bg-white after:border-gray-300 after:border after:rounded-full after:h-5 after:w-5 after:transition-all peer-checked:bg-indigo-600"></div>
            </label>
          </div>
        </div>
        <div class="mt-4">
          <BaseButton @click="saveNotificationSettings" variant="primary">
            Save Notification Settings
          </BaseButton>
        </div>
      </div>
    </BaseCard>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useToast } from '@/composables/useToast'
import BaseInput from '@/components/common/BaseInput.vue'
import BaseSelect from '@/components/common/BaseSelect.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseCard from '@/components/common/BaseCard.vue'

const { showToast } = useToast()

// Profile form data
const profile = ref({
  fullName: '',
  email: '',
  phone: '',
  language: 'en'
})

// Password form data
const passwordForm = ref({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
})

// Notification settings
const notifications = ref({
  email: true,
  sms: false,
  events: true
})

// Language options
const languageOptions = [
  { value: 'en', label: 'English' },
  { value: 'hi', label: 'Hindi' },
  { value: 'ta', label: 'Tamil' },
  { value: 'te', label: 'Telugu' },
  { value: 'mr', label: 'Marathi' }
]

// Load user data
onMounted(async () => {
  try {
    // This would normally fetch from an API
    // For demonstration, using mock data
    profile.value = {
      fullName: 'Arjun Sharma',
      email: 'arjun.sharma@example.com',
      phone: '+91 98765 43210',
      language: 'en'
    }
  } catch (error) {
    showToast('Error loading profile data', 'error')
  }
})

// Update profile information
const updateProfile = async () => {
  try {
    // This would normally call an API
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    showToast('Profile updated successfully', 'success')
  } catch (error) {
    showToast('Failed to update profile', 'error')
  }
}

// Change password
const changePassword = async () => {
  try {
    // Validate passwords match
    if (passwordForm.value.newPassword !== passwordForm.value.confirmPassword) {
      showToast('New passwords do not match', 'error')
      return
    }
    
    // This would normally call an API
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    
    // Reset form
    passwordForm.value = {
      currentPassword: '',
      newPassword: '',
      confirmPassword: ''
    }
    
    showToast('Password changed successfully', 'success')
  } catch (error) {
    showToast('Failed to change password', 'error')
  }
}

// Save notification settings
const saveNotificationSettings = async () => {
  try {
    // This would normally call an API
    // Simulate API call
    await new Promise(resolve => setTimeout(resolve, 1000))
    showToast('Notification settings saved', 'success')
  } catch (error) {
    showToast('Failed to save notification settings', 'error')
  }
}
</script>

<style scoped>
.settings-container {
  max-width: 900px;
  margin: 0 auto;
  padding: 20px;
}
</style>