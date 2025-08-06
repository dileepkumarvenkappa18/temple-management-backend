<template>
    <!-- Welcome Banner -->
    <WelcomeBanner 
      :user="currentUser" 
      :temple="currentTemple"
      role="devotee"
      class="mb-8"
      v-if="currentUser && currentTemple"
    />

    <!-- Profile Completion Alert -->
    <div v-if="profileCompletionPercentage < 100" class="mb-6">
      <div class="bg-amber-50 border border-amber-200 rounded-xl p-4">
        <div class="flex items-center">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-amber-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M8.257 3.099c.765-1.36 2.722-1.36 3.486 0l5.58 9.92c.75 1.334-.213 2.98-1.742 2.98H4.42c-1.53 0-2.493-1.646-1.743-2.98l5.58-9.92zM11 13a1 1 0 11-2 0 1 1 0 012 0zm-1-8a1 1 0 00-1 1v3a1 1 0 002 0V6a1 1 0 00-1-1z" clip-rule="evenodd"/>
            </svg>
          </div>
          <div class="ml-3 flex-1">
            <h3 class="text-sm font-medium text-amber-800">Complete Your Profile</h3>
            <p class="text-sm text-amber-700 mt-1">
              Your profile is {{ profileCompletionPercentage }}% complete. 
                <router-link 
                :to="`/entity/${route.params.id}/devotee/profile/create`" 
                class="font-medium text-amber-800 hover:text-amber-900 underline"
              >
                Complete now
              </router-link>
            </p>
          </div>
          <div class="ml-4">
            <div class="w-16 h-2 bg-amber-200 rounded-full">
              <div 
                class="h-2 bg-amber-500 rounded-full transition-all duration-300"
                :style="{ width: `${profileCompletionPercentage}%` }"
              ></div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="flex items-center justify-center py-12">
      <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600"></div>
      <span class="ml-3 text-gray-600">Loading dashboard...</span>
    </div>

    <!-- Error State -->
    <div v-else-if="error" class="bg-white rounded-xl shadow-md p-8 text-center">
      <svg class="mx-auto h-12 w-12 text-red-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
      </svg>
      <h3 class="mt-2 text-lg font-medium text-gray-900">Unable to load dashboard</h3>
      <p class="mt-1 text-sm text-red-600">{{ error }}</p>
      <button 
        @click="loadDashboardData"
        class="mt-4 inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
      >
        Try Again
      </button>
    </div>

    <!-- Dashboard Content -->
    <div v-else>
  <!-- Dashboard Stats -->
  <div class="grid grid-cols-1 md:grid-cols-3 lg:grid-cols-3 gap-6 mb-8">
    <DashboardWidget 
      title="Seva Bookings"
      :value="stats.totalSevaBookings || mySevaBookings.length || 0"
      icon="calendar"
      color="indigo"
      :subtitle="`${stats.upcomingSevas || mySevaBookings.length || 0} upcoming`"
    />
    <DashboardWidget 
      title="Total Events"
      :value="activityStats.totalEvents || stats.eventsAttended || 0"
      icon="calendar-days"
      color="emerald"
      :subtitle="`${activityStats.upcomingEvents || stats.upcomingEvents || 0} upcoming`"
    />
    <DashboardWidget 
      title="Temple Profile"
      :value="currentTemple ? 'Active' : 'Select Temple'"
      icon="building-library"
      color="rose"
      :subtitle="currentTemple ? currentTemple.name : 'Join a temple'"
    />
  </div>

      <!-- Main Content Grid -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Left Column - Activities Feed -->
        <div class="lg:col-span-2 space-y-6">
          <!-- Recent Donations -->
          <BaseCard class="p-6">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-gray-900">Recent Donations</h3>
              <router-link 
                :to="`/entity/${route.params.id}/devotee/donations`"
                class="text-sm font-medium text-indigo-600 hover:text-indigo-700"
              >
                View All
              </router-link>
            </div>
            
            <div v-if="recentDonations && recentDonations.length > 0" class="space-y-3">
              <div 
                v-for="donation in recentDonations" 
                :key="donation.id"
                class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div class="flex items-center space-x-4">
                  <div class="flex-shrink-0">
                    <div class="w-10 h-10 bg-rose-100 rounded-full flex items-center justify-center">
                      <svg class="w-5 h-5 text-rose-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                      </svg>
                    </div>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">₹{{ donation.amount || 0 }}</p>
                    <p class="text-xs text-gray-500">
                      {{ formatDate(donation.donation_date || donation.date) }}
                    </p>
                  </div>
                </div>
                <div class="text-right">
                  <span 
                    :class="getDonationStatusClass(donation.status || 'completed')" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ donation.status || 'Completed' }}
                  </span>
                </div>
              </div>
            </div>
            
            <div v-else class="text-center py-8">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No recent donations</h3>
              <p class="mt-1 text-sm text-gray-500">Make your first donation to support the temple.</p>
              <div class="mt-6">
                <router-link 
                  :to="`/entity/${route.params.id}/devotee/donations`"
                  class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Donate Now
                </router-link>
              </div>
            </div>
          </BaseCard>

          <!-- My Seva Bookings -->
          <BaseCard class="p-6">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-gray-900">My Seva Bookings</h3>
              <router-link 
                :to="`/entity/${route.params.id}/devotee/my-seva-bookings`"
                class="text-sm font-medium text-indigo-600 hover:text-indigo-700"
              >
                View All
              </router-link>
            </div>
            
            <div v-if="loadingRecentSevas" class="text-center py-8">
              <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-indigo-600 mx-auto"></div>
              <span class="mt-2 text-sm text-gray-500 block">Loading seva bookings...</span>
            </div>
            
            <div v-else-if="recentSevaBookings && recentSevaBookings.length > 0" class="space-y-3">
              <div 
                v-for="booking in recentSevaBookings.slice(0, 3)" 
                :key="getId(booking)"
                class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div class="flex items-center space-x-4">
                  <div class="flex-shrink-0">
                    <div class="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
                      <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ getSevaName(booking) }}</p>
                    <p class="text-xs text-gray-500">
                      {{ formatBookingDate(booking) }}
                    </p>
                  </div>
                </div>
                <div class="text-right">
                  <span 
                    :class="getStatusClass(getStatus(booking))" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ (getStatus(booking)).charAt(0).toUpperCase() + (getStatus(booking)).slice(1) }}
                  </span>
                </div>
              </div>
            </div>

            <div v-else-if="mySevaBookings && mySevaBookings.length > 0" class="space-y-3">
              <div 
                v-for="booking in mySevaBookings.slice(0, 3)" 
                :key="booking.id"
                class="flex items-center justify-between p-4 bg-gray-50 rounded-lg"
              >
                <div class="flex items-center space-x-4">
                  <div class="flex-shrink-0">
                    <div class="w-10 h-10 bg-indigo-100 rounded-full flex items-center justify-center">
                      <svg class="w-5 h-5 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                      </svg>
                    </div>
                  </div>
                  <div>
                    <p class="text-sm font-medium text-gray-900">{{ booking.seva?.name || booking.name }}</p>
                    <p class="text-xs text-gray-500">
                      {{ formatDate(booking.booking_date) }} at {{ formatTime(booking.booking_time) }}
                    </p>
                  </div>
                </div>
                <div class="text-right">
                  <span 
                    :class="getStatusClass(booking.status || 'pending')" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ booking.status || 'Pending' }}
                  </span>
                </div>
              </div>
            </div>
            
            <div v-else class="text-center py-8">
              <svg class="mx-auto h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <h3 class="mt-2 text-sm font-medium text-gray-900">No upcoming sevas</h3>
              <p class="mt-1 text-sm text-gray-500">Book your first seva to get started.</p>
              <div class="mt-6">
                <router-link 
                  :to="`/entity/${route.params.id}/devotee/seva-booking`"
                  class="inline-flex items-center px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
                >
                  Book Seva
                </router-link>
              </div>
            </div>
          </BaseCard>

          <!-- Temple Info -->
          <BaseCard class="p-6" v-if="currentTemple">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-gray-900">Temple Information</h3>
              <span class="text-sm font-medium text-indigo-600">{{ currentTemple.status || 'Active' }}</span>
            </div>
            
            <div class="space-y-4">
              <div class="flex items-start">
                <div class="flex-shrink-0 w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                </div>
                <div class="ml-4">
                  <h4 class="text-sm font-medium text-gray-900">Address</h4>
                  <p class="text-sm text-gray-500 mt-1">
                    {{ currentTemple.address || 'No address available' }}<br>
                    {{ [currentTemple.city, currentTemple.state, currentTemple.pincode].filter(Boolean).join(', ') }}
                  </p>
                </div>
              </div>
              
              <div class="flex items-start">
                <div class="flex-shrink-0 w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                  </svg>
                </div>
                <div class="ml-4">
                  <h4 class="text-sm font-medium text-gray-900">Contact Information</h4>
                  <p class="text-sm text-gray-500 mt-1">
                    {{ currentTemple.phone || 'No phone available' }}<br>
                    {{ currentTemple.email || 'No email available' }}
                  </p>
                </div>
              </div>
              
              <div class="flex items-start">
                <div class="flex-shrink-0 w-12 h-12 bg-indigo-100 rounded-lg flex items-center justify-center">
                  <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                </div>
                <div class="ml-4">
                  <h4 class="text-sm font-medium text-gray-900">Temple Hours</h4>
                  <p class="text-sm text-gray-500 mt-1">
                    {{ currentTemple.operating_hours || 'Hours not specified' }}
                  </p>
                </div>
              </div>
            </div>
          </BaseCard>
        </div>

        <!-- Right Column - Events & Notifications -->
        <div class="space-y-6">
          <!-- Upcoming Events -->
          <BaseCard class="p-6">
            <div class="flex items-center justify-between mb-4">
              <h3 class="text-lg font-semibold text-gray-900">Upcoming Events</h3>
              <router-link 
                :to="`/entity/${route.params.id}/devotee/events`"
                class="text-sm font-medium text-indigo-600 hover:text-indigo-700"
              >
                View All
              </router-link>
            </div>
            
            <div v-if="upcomingEvents && upcomingEvents.length > 0" class="space-y-4">
              <div 
                v-for="event in upcomingEvents" 
                :key="event.id"
                class="border border-gray-200 rounded-lg p-4 hover:border-indigo-300 transition-colors"
              >
                <div class="flex items-start space-x-3">
                  <div class="flex-shrink-0">
                    <div class="w-12 h-12 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-lg flex items-center justify-center">
                      <span class="text-white font-semibold">{{ formatEventDate(event.event_date || event.date) }}</span>
                    </div>
                  </div>
                  <div class="flex-1 min-w-0">
                    <p class="text-sm font-medium text-gray-900 truncate">{{ event.title || event.name }}</p>
                    <p class="text-xs text-gray-500 mt-1">{{ formatDate(event.event_date || event.date) }}</p>
                    <div class="mt-2 flex items-center space-x-2">
                      <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-green-100 text-green-800">
                        {{ event.status || 'Upcoming' }}
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            </div>
            
            <div v-else class="text-center py-6">
              <svg class="mx-auto h-10 w-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
              </svg>
              <p class="mt-2 text-sm text-gray-500">No upcoming events</p>
            </div>
          </BaseCard>

          <!-- Notifications -->
          <BaseCard class="p-6" v-if="notifications && notifications.length > 0">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Notifications</h3>
            
            <div class="space-y-3">
              <div 
                v-for="notification in notifications" 
                :key="notification.id"
                class="flex items-start space-x-3 p-3 bg-gray-50 rounded-lg"
              >
                <div class="flex-shrink-0">
                  <div class="w-2 h-2 bg-indigo-500 rounded-full mt-2"></div>
                </div>
                <div class="flex-1 min-w-0">
                  <p class="text-sm text-gray-900">{{ notification.message }}</p>
                  <p class="text-xs text-gray-500 mt-1">{{ formatDate(notification.date) }}</p>
                </div>
              </div>
            </div>
          </BaseCard>

          <!-- Quick Actions -->
          <BaseCard class="p-6">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">Quick Actions</h3>
            <div class="space-y-3">
              <router-link 
                :to="`/entity/${route.params.id}/devotee/donations`"
                class="inline-flex items-center justify-center w-full px-4 py-2 border border-transparent shadow-sm text-sm font-medium rounded-md text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
                </svg>
                Donate Now
              </router-link>
              
              <router-link 
                :to="`/entity/${route.params.id}/devotee/seva-booking`"
                class="inline-flex items-center justify-center w-full px-4 py-2 border border-indigo-300 shadow-sm text-sm font-medium rounded-md text-indigo-700 bg-indigo-50 hover:bg-indigo-100 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                </svg>
                Book Seva
              </router-link>
              
              <router-link 
                :to="`/entity/${route.params.id}/devotee/profile/create`"
                class="inline-flex items-center justify-center w-full px-4 py-2 border border-gray-300 shadow-sm text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500"
              >
                <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"/>
                </svg>
                Edit Profile
              </router-link>
            </div>
          </BaseCard>

          <!-- Temple Information -->
          <BaseCard class="p-6" v-if="currentTemple">
            <h3 class="text-lg font-semibold text-gray-900 mb-4">About Temple</h3>
            <div class="space-y-4">
              <div class="flex items-center gap-2">
                <div class="text-sm font-medium text-gray-500">Type:</div>
                <div class="text-sm text-gray-900">{{ currentTemple.temple_type || 'Not specified' }}</div>
              </div>
              
              <div class="flex items-center gap-2">
                <div class="text-sm font-medium text-gray-500">Established:</div>
                <div class="text-sm text-gray-900">{{ currentTemple.established_year || 'Unknown' }}</div>
              </div>
              
              <div>
                <div class="text-sm font-medium text-gray-500 mb-1">Description:</div>
                <div class="text-sm text-gray-900">
                  {{ currentTemple.description || 'No description available' }}
                </div>
              </div>
              
              <div v-if="currentTemple.website">
                <a 
                  :href="currentTemple.website" 
                  target="_blank"
                  class="inline-flex items-center text-sm text-indigo-600 hover:text-indigo-800"
                >
                  <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                  </svg>
                  Visit Website
                </a>
              </div>
            </div>
          </BaseCard>
        </div>
      </div>
    </div>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDevoteeStore } from '@/stores/devotee'
import { useTempleStore } from '@/stores/temple'
import { useSevaStore } from '@/stores/seva'
import { useTempleActivities } from '@/composables/useTempleActivities'
import WelcomeBanner from '@/components/dashboard/WelcomeBanner.vue'
import DashboardWidget from '@/components/dashboard/DashboardWidget.vue'
import ActivityFeed from '@/components/dashboard/ActivityFeed.vue'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'
// Donation form is not needed until donation backend is completed
// import DonationForm from '@/components/donation/DonationForm.vue'

// Define props to handle the 'id' attribute that's being passed
const props = defineProps({
  id: {
    type: [String, Number],
    default: null
  }
})

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()
const devoteeStore = useDevoteeStore()
const templeStore = useTempleStore()
const sevaStore = useSevaStore()
const { 
  loading: activitiesLoading, 
  error: activitiesError, 
  activities, 
  stats: activityStats, 
  fetchAllActivities, 
  fetchUpcomingEvents, 
  fetchMySevaBookings,
  fetchEntityDetails 
} = useTempleActivities()

// Reactive data
const loading = ref(true)
const error = ref(null)
const mySevaBookings = ref([])
const upcomingEvents = ref([])
const entityDetails = ref(null)
const recentDonations = ref([])

// Computed properties
const currentUser = computed(() => authStore.user)
const currentTemple = computed(() => {
  if (entityDetails.value) {
    return entityDetails.value
  }
  return templeStore.temples?.find(t => t.id === parseInt(route.params.id)) || templeStore.currentTemple
})

const profileCompletionPercentage = computed(() => {
  return devoteeStore.completionPercentage || 0
})

// Use existing data from store if available
const stats = computed(() => ({
  totalSevaBookings: devoteeStore.totalSevaBookings || 0,
  upcomingSevas: devoteeStore.upcomingSevas?.length || 0,
  totalDonations: devoteeStore.totalDonations || 0,
  donationCount: devoteeStore.donations?.length || 0,
  eventsAttended: devoteeStore.eventsAttended || 0,
  upcomingEvents: devoteeStore.upcomingEvents?.length || 0
}))

const upcomingSevas = computed(() => {
  return devoteeStore.upcomingSevas?.slice(0, 3) || []
})

const notifications = computed(() => {
  return devoteeStore.notifications?.slice(0, 5) || []
})

// For seva store
const recentSevaBookings = computed(() => {
  return sevaStore.recentSevas
})

const loadingRecentSevas = computed(() => sevaStore.loadingRecentSevas)

// Helper functions to handle database field names
const getId = (booking) => booking?.ID || booking?.id || 0
const getSevaId = (booking) => booking?.SevaID || booking?.seva_id || 0
const getStatus = (booking) => booking?.Status || booking?.status || 'pending'
const getBookingTime = (booking) => booking?.BookingTime || booking?.booking_time || null

// Helper function to get Seva name from a booking
const getSevaName = (booking) => {
  if (!booking) return 'Unknown Seva'
  
  // Check for direct seva_name property first (both cases)
  if (booking.SevaName || booking.seva_name) {
    return booking.SevaName || booking.seva_name
  }
  
  // Check for seva catalog data in sevaStore
  const sevaId = getSevaId(booking)
  if (sevaStore.sevaCatalog && Array.isArray(sevaStore.sevaCatalog)) {
    const catalogItem = sevaStore.sevaCatalog.find(item => item.id === sevaId)
    if (catalogItem && catalogItem.name) {
      return catalogItem.name
    }
  }
  
  // Check for seva in seva list
  if (sevaStore.sevaList && Array.isArray(sevaStore.sevaList)) {
    const seva = sevaStore.sevaList.find(s => s.id === sevaId)
    if (seva && seva.name) {
      return seva.name
    }
  }
  
  // Try to lookup in global catalog or seva entries
  if (sevaStore.getSevaNameById && typeof sevaStore.getSevaNameById === 'function') {
    const name = sevaStore.getSevaNameById(sevaId)
    if (name) return name
  }
  
  // Return generic name with ID as a last resort
  return `Seva ${sevaId}`
}

// Format the booking date
const formatBookingDate = (booking) => {
  const dateTime = getBookingTime(booking)
  if (!dateTime) return ''
  
  try {
    const date = new Date(dateTime)
    return date.toLocaleDateString('en-IN', {
      day: '2-digit',
      month: 'short',
      year: 'numeric'
    })
  } catch (error) {
    console.error('Error formatting booking date:', error)
    return ''
  }
}

// Status styles
const getStatusClass = (status) => {
  const classes = {
    pending: 'bg-yellow-100 text-yellow-800',
    approved: 'bg-green-100 text-green-800',
    rejected: 'bg-red-100 text-red-800',
    completed: 'bg-blue-100 text-blue-800',
    cancelled: 'bg-red-100 text-red-800'
  }
  return classes[status.toLowerCase()] || 'bg-gray-100 text-gray-800'
}

// Methods
const formatDate = (date) => {
  if (!date) return 'N/A'
  try {
    return new Intl.DateTimeFormat('en-IN', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    }).format(new Date(date))
  } catch (error) {
    return 'Invalid Date'
  }
}

const formatTime = (time) => {
  if (!time) return 'N/A'
  // For handling both string time formats and Date objects
  try {
    if (typeof time === 'string' && time.includes(':')) {
      // Handle "HH:MM" format
      return time
    } else {
      // Handle date object or ISO string
      return new Date(time).toLocaleTimeString('en-IN', {
        hour: '2-digit',
        minute: '2-digit'
      })
    }
  } catch (error) {
    return 'Invalid Time'
  }
}

const formatEventDate = (date) => {
  if (!date) return '?'
  try {
    return new Intl.DateTimeFormat('en-IN', {
      day: 'numeric'
    }).format(new Date(date))
  } catch (error) {
    return '?'
  }
}

const getDonationStatusClass = (status) => {
  const statusMap = {
    pending: 'bg-yellow-100 text-yellow-800',
    processing: 'bg-blue-100 text-blue-800',
    completed: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
    failed: 'bg-red-100 text-red-800'
  }
  
  // Handle case insensitive
  const normalizedStatus = (status || '').toLowerCase()
  return statusMap[normalizedStatus] || 'bg-gray-100 text-gray-800'
}

// Load dashboard data from API
const loadDashboardData = async () => {
  const entityId = route.params.id
  if (!entityId) {
    error.value = 'Entity ID is missing'
    loading.value = false
    return
  }

  loading.value = true
  error.value = null

  try {
    // Set the entity ID for API calls
    localStorage.setItem('current_entity_id', entityId)
    localStorage.setItem('current_tenant_id', entityId) // ✅ Enables proper X-Tenant-ID header in axios

    // Make parallel API requests
    const promises = []

    // Load profile from Pinia store if available
    if (devoteeStore.loadProfileData) {
      promises.push(devoteeStore.loadProfileData())
    }

    // Load temple data
    promises.push(fetchEntityDetails(entityId)
      .then(data => {
        entityDetails.value = data
      })
      .catch(err => {
        console.error('Error loading entity details:', err)
        // If we can't fetch entity details, try loading from store
        if (templeStore.fetchTemples) {
          return templeStore.fetchTemples()
        } else if (templeStore.loadTemple) {
          return templeStore.loadTemple(entityId)
        }
      })
    )

    // Load upcoming events
    promises.push(fetchUpcomingEvents(entityId)
      .then(data => {
        if (Array.isArray(data)) {
          upcomingEvents.value = data.slice(0, 3)
        } else {
          console.warn('Expected array but got:', data)
          upcomingEvents.value = []
        }
      })
    )

    // Load seva bookings
    promises.push(fetchMySevaBookings()
      .then(data => {
        if (Array.isArray(data)) {
          mySevaBookings.value = data.slice(0, 3)
        } else {
          console.warn('Expected array but got:', data)
          mySevaBookings.value = []
        }
      })
    )

    // Load donations data
    promises.push(
      // If we have a fetchDonations method, use it
      devoteeStore.fetchDonations 
        ? devoteeStore.fetchDonations()
            .then(data => {
              if (Array.isArray(data)) {
                recentDonations.value = data.slice(0, 5)
              }
            })
        // Otherwise try to use store data
        : Promise.resolve().then(() => {
            if (devoteeStore.donations && Array.isArray(devoteeStore.donations)) {
              recentDonations.value = devoteeStore.donations.slice(0, 5)
            }
          })
    )
    
    // Load recent sevas
    promises.push(
      sevaStore.fetchRecentSevas()
        .catch(err => {
          console.error('Error loading recent sevas:', err)
        })
    )
    
    // Load seva catalog to get actual seva names
    if (sevaStore.fetchSevaCatalog) {
      promises.push(
        sevaStore.fetchSevaCatalog()
          .then(() => {
            console.log('Seva catalog loaded:', sevaStore.sevaCatalog)
          })
          .catch(err => {
            console.error('Error loading seva catalog:', err)
          })
      )
    }

    // Wait for all promises to resolve
    await Promise.all(promises)
    
    // Debug the loaded data
    console.log('Recent sevas:', sevaStore.recentSevas)
    console.log('Seva catalog:', sevaStore.sevaCatalog)
    
  } catch (err) {
    console.error('Error loading dashboard data:', err)
    error.value = err.message || 'Failed to load dashboard data'
  } finally {
    loading.value = false
  }
}

// Handle activity feed events
const handleActivitiesLoaded = (activities) => {
  // Update dashboard with loaded activities
  console.log('Activities loaded:', activities.length)
}

const handleActivitiesError = (errorMessage) => {
  console.error('Error loading activities:', errorMessage)
  // Don't set global error to avoid hiding entire dashboard
}

const handleViewAllActivities = () => {
  // Navigate to activities page or open modal
  router.push(`/entity/${route.params.id}/devotee/activities`)
}

// Lifecycle
onMounted(async () => {
  await loadDashboardData()
})
</script>