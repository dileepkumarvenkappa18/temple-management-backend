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
 <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
  <!-- Seva Bookings Card -->
<router-link 
  :to="{ name: 'DevoteeMySevaBookings', params: { id: $route.params.id } }"
  class="block"
>
  <DashboardWidget 
    title="Seva Bookings"
    :value="sevaCount"
    icon="calendar"
    color="indigo"
    :subtitle="`${upcomingSevaCount} upcoming`"
  />
</router-link>

<!-- Total Events Card -->
<router-link 
  :to="{ name: 'DevoteeMyEvents', params: { id: $route.params.id } }"
  class="block"
>
  <DashboardWidget 
    title="Total Events"
    :value="eventCount"
    icon="calendar-days"
    color="emerald"
    :subtitle="`${upcomingEventCount} upcoming`"
  />
</router-link>

  <!-- NEW: Birthday Card -->
  <DashboardWidget 
    title="My Birthday"
    :value="formatBirthday()"
    icon="cake"
    color="purple"
    :subtitle="getBirthdaySubtitle()"
  />
<router-link 
  :to="{ name: 'DevoteeTempleSelection' }"
  class="block transition-transform hover:scale-105 focus:outline-none focus:ring-2 focus:ring-rose-500 focus:ring-offset-2 rounded-lg"
>
  <DashboardWidget 
    title="Temple Profile"
    :value="currentTemple ? 'Active' : 'Select Temple'"
    icon="building-library"
    color="rose"
    :subtitle="currentTemple ? currentTemple.name : 'Join a temple'"
    class="cursor-pointer hover:shadow-lg transition-shadow duration-200"
  />
</router-link>
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
                      {{ formatDate(donation.created_at || donation.date) }}
                    </p>
                  </div>
                </div>
                <div class="text-right">
                  <!-- <span 
                    :class="getDonationStatusClass(donation.status || 'completed')" 
                    class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium"
                  >
                    {{ donation.status || 'Completed' }}
                  </span> -->
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
import { ref, reactive, computed, onMounted, watch, nextTick } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import { useDevoteeStore } from '@/stores/devotee'
import { useTempleStore } from '@/stores/temple'
import { useSevaStore } from '@/stores/seva'
import { useDonationStore } from '@/stores/donation'
import { useTempleActivities } from '@/composables/useTempleActivities'
import devoteeService from '@/services/devotee.service'

import WelcomeBanner from '@/components/dashboard/WelcomeBanner.vue'
import DashboardWidget from '@/components/dashboard/DashboardWidget.vue'
import ActivityFeed from '@/components/dashboard/ActivityFeed.vue'
import BaseCard from '@/components/common/BaseCard.vue'
import BaseButton from '@/components/common/BaseButton.vue'
import BaseModal from '@/components/common/BaseModal.vue'

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
const donationStore = useDonationStore()
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

// FIX: Add data loading trackers
const dataLoadingStatus = reactive({
  profile: false,
  entity: false,
  sevas: false,
  events: false,
  donations: false
})

// FIX: Add a cache for seva names to ensure consistency
const sevaNameCache = reactive({})

// PERSISTENT COUNTER STATE
const dashboardCounts = reactive({
  sevaCount: 0,
  upcomingSevaCount: 0,
  eventCount: 0,
  upcomingEventCount: 0,
  countsInitialized: false
})

// Local persistent storage for dashboard data
const localData = reactive({
  recentSevas: [],
  recentDonations: [],
  upcomingEvents: []
})

// Compute the actual counts with locking mechanism
const sevaCount = computed(() => {
  // Return the locked value if set
  if (dashboardCounts.countsInitialized && dashboardCounts.sevaCount > 0) {
    return dashboardCounts.sevaCount
  }
  
  // Calculate from available data
  const count = localData.recentSevas.length || 
                sevaStore.recentSevas?.length ||
                mySevaBookings.value?.length || 
                devoteeStore.dashboardData?.summary?.totalSevaBookings || 
                0
                
  // Lock the value if it's greater than 0 and not already locked
  if (count > 0 && !dashboardCounts.countsInitialized) {
    dashboardCounts.sevaCount = count
  }
  
  return count
})

const upcomingSevaCount = computed(() => {
  // Return the locked value if set
  if (dashboardCounts.countsInitialized && dashboardCounts.upcomingSevaCount > 0) {
    return dashboardCounts.upcomingSevaCount
  }
  
  // Calculate from available data
  const sevasToCheck = localData.recentSevas.length > 0 ? 
                       localData.recentSevas : 
                       sevaStore.recentSevas || []
                       
  const pendingSevas = sevasToCheck.filter(
    s => (s.status?.toLowerCase() === 'pending' || 
         s.Status?.toLowerCase() === 'pending')
  )?.length || 0
  
  const count = pendingSevas || 
                devoteeStore.dashboardData?.summary?.upcomingSevas || 
                0
                
  // Lock the value if it's greater than 0 and not already locked
  if (count > 0 && !dashboardCounts.countsInitialized) {
    dashboardCounts.upcomingSevaCount = count
  }
  
  return count
})

const eventCount = computed(() => {
  // Return the locked value if set
  if (dashboardCounts.countsInitialized && dashboardCounts.eventCount > 0) {
    return dashboardCounts.eventCount
  }
  
  // Calculate from available data
  const count = localData.upcomingEvents.length > 0 ? 
                localData.upcomingEvents.length : 
                upcomingEvents.value?.length || 
                activityStats.totalEvents || 
                devoteeStore.dashboardData?.summary?.eventsAttended || 
                0
                
  // Lock the value if it's greater than 0 and not already locked
  if (count > 0 && !dashboardCounts.countsInitialized) {
    dashboardCounts.eventCount = count
  }
  
  return count
})

const upcomingEventCount = computed(() => {
  // Return the locked value if set
  if (dashboardCounts.countsInitialized && dashboardCounts.upcomingEventCount > 0) {
    return dashboardCounts.upcomingEventCount
  }
  
  // Calculate from available data
  const count = localData.upcomingEvents.length > 0 ? 
                localData.upcomingEvents.length : 
                upcomingEvents.value?.length || 
                activityStats.upcomingEvents || 
                devoteeStore.dashboardData?.summary?.upcomingEvents || 
                0
                
  // Lock the value if it's greater than 0 and not already locked
  if (count > 0 && !dashboardCounts.countsInitialized) {
    dashboardCounts.upcomingEventCount = count
  }
  
  return count
})

// Lock the counters once we have valid data to prevent flickering
const lockCounters = () => {
  if (dashboardCounts.countsInitialized) return
  
  // Only lock if we have at least some valid data
  if (sevaCount.value > 0 || eventCount.value > 0) {
    console.log('📊 Locking dashboard counters:')
    console.log(`- Seva count: ${sevaCount.value}`)
    console.log(`- Upcoming seva count: ${upcomingSevaCount.value}`)
    console.log(`- Event count: ${eventCount.value}`)
    console.log(`- Upcoming event count: ${upcomingEventCount.value}`)
    
    dashboardCounts.sevaCount = sevaCount.value 
    dashboardCounts.upcomingSevaCount = upcomingSevaCount.value
    dashboardCounts.eventCount = eventCount.value
    dashboardCounts.upcomingEventCount = upcomingEventCount.value
    dashboardCounts.countsInitialized = true
    
    // Save to localStorage for persistence
    try {
      localStorage.setItem('dashboard_counts', JSON.stringify({
        sevaCount: dashboardCounts.sevaCount,
        upcomingSevaCount: dashboardCounts.upcomingSevaCount,
        eventCount: dashboardCounts.eventCount,
        upcomingEventCount: dashboardCounts.upcomingEventCount,
        timestamp: Date.now()
      }))
    } catch (e) {
      console.warn('Could not save dashboard counts to localStorage', e)
    }
  }
}

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

// Get devotee profile
const devoteeProfile = computed(() => {
  return devoteeStore.profile || authStore.user || {}
})

// Get recent donations
const recentDonations = computed(() => {
  // Check current donation store first
  if (donationStore.recentDonations && donationStore.recentDonations.length > 0) {
    return donationStore.recentDonations
  }
  
  // Then check local data
  if (localData.recentDonations && localData.recentDonations.length > 0) {
    return localData.recentDonations
  }
  
  // Return empty array if nothing is available
  return []
})

// Get loading state for recent donations
const loadingRecentDonations = computed(() => {
  return donationStore.loadingRecent && !dataLoadingStatus.donations
})

// Use existing data from store if available
const stats = computed(() => ({
  totalSevaBookings: devoteeStore.totalSevaBookings || 0,
  upcomingSevas: devoteeStore.upcomingSevas?.length || 0,
  totalDonations: donationStore.donationStats?.total || devoteeStore.totalDonations || 0,
  donationCount: donationStore.donations?.length || devoteeStore.donations?.length || 0,
  eventsAttended: devoteeStore.eventsAttended || 0,
  upcomingEvents: devoteeStore.upcomingEvents?.length || 0
}))

const upcomingSevas = computed(() => {
  return devoteeStore.upcomingSevas?.slice(0, 3) || []
})

const notifications = computed(() => {
  return devoteeStore.notifications?.slice(0, 5) || []
})

// Get recent seva bookings with name caching
const recentSevaBookings = computed(() => {
  const sevas = localData.recentSevas.length > 0 ? 
                localData.recentSevas : 
                sevaStore.recentSevas || []
  
  // Ensure seva names are cached as we compute this
  if (sevas.length > 0) {
    sevas.forEach(seva => {
      const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
      if (sevaId) {
        // Pre-cache name if it exists in the seva
        if (seva.seva_name || seva.SevaName || seva.name || seva.Name) {
          sevaNameCache[sevaId] = seva.seva_name || seva.SevaName || seva.name || seva.Name
        }
      }
    })
  }
  
  return sevas
})

const loadingRecentSevas = computed(() => sevaStore.loadingRecentSevas && !dataLoadingStatus.sevas)

// Restore saved data from localStorage on mount
const restoreSavedData = () => {
  try {
    // Try to restore saved dashboard counts
    const savedCounts = localStorage.getItem('dashboard_counts')
    if (savedCounts) {
      const parsed = JSON.parse(savedCounts)
      
      // Only use if saved less than 1 hour ago
      const oneHourAgo = Date.now() - (60 * 60 * 1000)
      if (parsed.timestamp && parsed.timestamp > oneHourAgo) {
        console.log('📊 Restoring saved dashboard counts:', parsed)
        dashboardCounts.sevaCount = parsed.sevaCount || 0
        dashboardCounts.upcomingSevaCount = parsed.upcomingSevaCount || 0
        dashboardCounts.eventCount = parsed.eventCount || 0
        dashboardCounts.upcomingEventCount = parsed.upcomingEventCount || 0
        dashboardCounts.countsInitialized = true
      }
    }
    
    // Try to restore saved local data
    const savedLocalData = localStorage.getItem('dashboard_local_data')
    if (savedLocalData) {
      const parsed = JSON.parse(savedLocalData)
      
      // Only use if saved less than 1 hour ago
      const oneHourAgo = Date.now() - (60 * 60 * 1000)
      if (parsed.timestamp && parsed.timestamp > oneHourAgo) {
        console.log('🗃️ Restoring saved local data')
        if (parsed.recentSevas && Array.isArray(parsed.recentSevas)) {
          localData.recentSevas = parsed.recentSevas
        }
        if (parsed.recentDonations && Array.isArray(parsed.recentDonations)) {
          localData.recentDonations = parsed.recentDonations
        }
        if (parsed.upcomingEvents && Array.isArray(parsed.upcomingEvents)) {
          localData.upcomingEvents = parsed.upcomingEvents
          upcomingEvents.value = parsed.upcomingEvents
        }
      }
    }
    
    // Try to restore saved seva name cache
    const savedSevaNames = localStorage.getItem('dashboard_seva_names')
    if (savedSevaNames) {
      try {
        const parsed = JSON.parse(savedSevaNames)
        Object.assign(sevaNameCache, parsed)
        console.log('📋 Restored seva name cache with', Object.keys(parsed).length, 'entries')
      } catch (e) {
        console.warn('Error parsing saved seva names:', e)
      }
    }
  } catch (e) {
    console.warn('Error restoring saved data:', e)
  }
}

// FIX: Enhanced name caching system
// Helper function to save the seva name cache to localStorage
const saveSevaNameCache = () => {
  try {
    localStorage.setItem('dashboard_seva_names', JSON.stringify(sevaNameCache))
  } catch (e) {
    console.warn('Could not save seva name cache to localStorage', e)
  }
}

// Helper functions to handle database field names
const getId = (booking) => booking?.ID || booking?.id || 0
const getSevaId = (booking) => booking?.SevaID || booking?.seva_id || 0
const getStatus = (booking) => booking?.Status || booking?.status || 'pending'
const getBookingTime = (booking) => booking?.BookingTime || booking?.booking_time || null

// FIX: Enhanced getSevaName with persistent caching 
const getSevaName = (booking) => {
  if (!booking) return 'Unknown Seva'
  
  // Extract the sevaId first
  const sevaId = getSevaId(booking)
  
  // First check our local cache for a previously resolved name
  if (sevaId && sevaNameCache[sevaId]) {
    return sevaNameCache[sevaId]
  }
  
  // Check for direct seva_name property first (both cases)
  if (booking.SevaName || booking.seva_name) {
    const name = booking.SevaName || booking.seva_name
    if (sevaId) {
      sevaNameCache[sevaId] = name
      saveSevaNameCache() // Save the cache after adding a new name
    }
    return name
  }
  
  // Check if the booking has seva object with name
  if (booking.seva && booking.seva.name) {
    const name = booking.seva.name
    if (sevaId) {
      sevaNameCache[sevaId] = name
      saveSevaNameCache()
    }
    return name
  }
  
  // Check for seva catalog data in sevaStore
  if (sevaId && sevaStore.sevaCatalog && Array.isArray(sevaStore.sevaCatalog)) {
    const catalogItem = sevaStore.sevaCatalog.find(item => 
      item.id === sevaId || item.ID === sevaId
    )
    if (catalogItem && (catalogItem.name || catalogItem.Name)) {
      const name = catalogItem.name || catalogItem.Name
      sevaNameCache[sevaId] = name
      saveSevaNameCache()
      return name
    }
  }
  
  // Check for seva in seva list
  if (sevaId && sevaStore.sevaList && Array.isArray(sevaStore.sevaList)) {
    const seva = sevaStore.sevaList.find(s => 
      s.id === sevaId || s.ID === sevaId
    )
    if (seva && (seva.name || seva.Name)) {
      const name = seva.name || seva.Name
      sevaNameCache[sevaId] = name
      saveSevaNameCache()
      return name
    }
  }
  
  // Check if there's a name field directly on the booking
  if (booking.name || booking.Name) {
    const name = booking.name || booking.Name
    if (sevaId) {
      sevaNameCache[sevaId] = name
      saveSevaNameCache()
    }
    return name
  }
  
  // Try to lookup in global catalog or seva entries
  if (sevaStore.getSevaNameById && typeof sevaStore.getSevaNameById === 'function') {
    try {
      const name = sevaStore.getSevaNameById(sevaId)
      if (name && sevaId) {
        sevaNameCache[sevaId] = name
        saveSevaNameCache()
        return name
      }
    } catch (e) {
      console.warn(`Error getting seva name for ID ${sevaId}:`, e)
    }
  }
  
  // Return generic name with ID as a last resort
  const genericName = `Seva ${sevaId || booking.id || 'Unknown'}`
  
  // Cache even the generic name to be consistent
  if (sevaId) {
    sevaNameCache[sevaId] = genericName
    saveSevaNameCache()
  }
  
  return genericName
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

// Method to format birthday - Shows actual date
const formatBirthday = () => {
  const profile = devoteeProfile.value
  
  // Look in the correct location: profile.personal.dateOfBirth
  const dob = profile.personal?.dateOfBirth || 
             profile.dob || 
             profile.date_of_birth || 
             profile.birthday
  
  if (!dob) {
    return 'Not Set'
  }
  
  try {
    // Parse the database format (2003-12-18T00:00:00Z)
    const date = new Date(dob)
    
    // Check if date is valid
    if (isNaN(date.getTime())) {
      return 'Invalid Date'
    }
    
    // Return in a nice format like "18 Dec" (without year to save space)
    return date.toLocaleDateString('en-IN', {
      day: 'numeric',
      month: 'short'
    })
  } catch (error) {
    return 'Invalid Date'
  }
}

// Method to get birthday subtitle with countdown
const getBirthdaySubtitle = () => {
  const profile = devoteeProfile.value
  
  // Look in the correct location: profile.personal.dateOfBirth  
  const dob = profile.personal?.dateOfBirth || 
             profile.dob || 
             profile.date_of_birth || 
             profile.birthday
  
  if (!dob) {
    return 'Add your birthday in profile'
  }
  
  try {
    // Parse the database format
    const birthDate = new Date(dob)
    
    // Check if date is valid
    if (isNaN(birthDate.getTime())) {
      return 'Add your birthday in profile'
    }
    
    const today = new Date()
    const currentYear = today.getFullYear()
    
    // Get this year's birthday
    const thisYearBirthday = new Date(currentYear, birthDate.getMonth(), birthDate.getDate())
    
    // If birthday already passed this year, check next year
    if (thisYearBirthday < today) {
      thisYearBirthday.setFullYear(currentYear + 1)
    }
    
    // Calculate days until birthday
    const diffTime = thisYearBirthday - today
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    
    // Calculate age (current year - birth year, adjust if birthday hasn't happened yet)
    let age = currentYear - birthDate.getFullYear()
    if (thisYearBirthday > today) {
      // Birthday hasn't happened this year yet, so they're still the previous age
      age = age  // Keep current calculation
    } else {
      // Birthday has passed this year, so add 1 to age for next year's birthday
      age = age + 1
    }
    
    if (diffDays === 0) {
      return `🎉 Happy ${age}th Birthday!`
    } else if (diffDays === 1) {
      return `🎂 Tomorrow! Turning ${age}`
    } else if (diffDays <= 7) {
      return `🎈 In ${diffDays} days (${age}th)`
    } else if (diffDays <= 30) {
      return `In ${diffDays} days • ${age} years`
    } else if (diffDays <= 100) {
      // Show months for longer periods
      const months = Math.floor(diffDays / 30)
      const remainingDays = diffDays % 30
      if (remainingDays === 0) {
        return `In ${months} month${months === 1 ? '' : 's'}`
      } else {
        return `In ${months}m ${remainingDays}d • Turning ${age}`
      }
    } else {
      // For very far dates, show year and age
      const nextYear = thisYearBirthday.getFullYear()
      return `${birthDate.toLocaleDateString('en-IN', { month: 'short', day: 'numeric' })}, ${nextYear} • Age ${age}`
    }
  } catch (error) {
    return 'Add your birthday in profile'
  }
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
    success: 'bg-green-100 text-green-800',
    cancelled: 'bg-red-100 text-red-800',
    failed: 'bg-red-100 text-red-800'
  }
  
  // Handle case insensitive
  const normalizedStatus = (status || '').toLowerCase()
  return statusMap[normalizedStatus] || 'bg-gray-100 text-gray-800'
}

// FIX: Enhanced data processing for sevas to extract and cache names
const processSevaData = (sevas) => {
  if (!sevas || !Array.isArray(sevas) || sevas.length === 0) return sevas
  
  return sevas.map(seva => {
    const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
    
    // Extract seva name from any available property
    let sevaName = null
    
    // Try to get name from seva object first
    if (seva.seva && (seva.seva.name || seva.seva.Name)) {
      sevaName = seva.seva.name || seva.seva.Name
    } 
    // Then direct properties
    else if (seva.seva_name || seva.SevaName || seva.name || seva.Name) {
      sevaName = seva.seva_name || seva.SevaName || seva.name || seva.Name
    }
    
    // Cache the name if we found one and have an ID
    if (sevaName && sevaId) {
      sevaNameCache[sevaId] = sevaName
    }
    
    // Return the original seva object
    return seva
  })
}

// FIX: Add a specialized function to enrich seva objects with their names
const enrichSevaWithNames = async (sevas) => {
  if (!sevas || !Array.isArray(sevas) || sevas.length === 0) return sevas

  // Make a deep copy to avoid reactivity issues
  const enrichedSevas = JSON.parse(JSON.stringify(sevas))
  
  console.log('🔄 Enriching', enrichedSevas.length, 'sevas with names')
  
  // First attempt: Try to map names from the seva catalog directly
  if (sevaStore.sevaCatalog && Array.isArray(sevaStore.sevaCatalog)) {
    const catalogMap = {}
    sevaStore.sevaCatalog.forEach(catalogItem => {
      const id = catalogItem.id || catalogItem.ID
      if (id && (catalogItem.name || catalogItem.Name)) {
        catalogMap[id] = catalogItem.name || catalogItem.Name
      }
    })
    
    // Apply catalog names to sevas
    enrichedSevas.forEach(seva => {
      const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
      if (sevaId && catalogMap[sevaId]) {
        // Store the name in multiple properties to ensure it's used
        seva.seva_name = catalogMap[sevaId]
        seva.SevaName = catalogMap[sevaId]
        if (!seva.name) seva.name = catalogMap[sevaId]
        if (!seva.Name) seva.Name = catalogMap[sevaId]
        
        // If seva has a seva object, ensure it has the name
        if (seva.seva && typeof seva.seva === 'object') {
          seva.seva.name = catalogMap[sevaId]
        }
        
        // Cache the name
        sevaNameCache[sevaId] = catalogMap[sevaId]
      }
    })
  }
  
  // Second attempt: Try to use the cached seva names
  enrichedSevas.forEach(seva => {
    const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
    if (sevaId && sevaNameCache[sevaId]) {
      // Apply the cached name to all possible fields
      seva.seva_name = sevaNameCache[sevaId]
      seva.SevaName = sevaNameCache[sevaId]
      if (!seva.name) seva.name = sevaNameCache[sevaId]
      if (!seva.Name) seva.Name = sevaNameCache[sevaId]
      
      // If seva has a seva object, ensure it has the name
      if (seva.seva && typeof seva.seva === 'object') {
        seva.seva.name = sevaNameCache[sevaId]
      }
    }
  })
  
  // Third attempt: If all else fails, try to look up names from the server
  // This would be specific to your API structure, but here's a placeholder
  if (sevaStore.getSevaDetails && typeof sevaStore.getSevaDetails === 'function') {
    // Find sevas that still don't have names
    const sevasNeedingNames = enrichedSevas.filter(seva => {
      const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
      return sevaId && 
        !seva.seva_name && 
        !seva.SevaName && 
        !seva.name && 
        !seva.Name && 
        !sevaNameCache[sevaId]
    })
    
    // If we have sevas needing names, fetch them one by one
    if (sevasNeedingNames.length > 0) {
      console.log('🔄 Fetching names for', sevasNeedingNames.length, 'sevas')
      
      for (const seva of sevasNeedingNames) {
        const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
        if (!sevaId) continue
        
        try {
          // This would be your API call to get seva details by ID
          const sevaDetails = await sevaStore.getSevaDetails(sevaId)
          if (sevaDetails && (sevaDetails.name || sevaDetails.Name)) {
            const name = sevaDetails.name || sevaDetails.Name
            
            // Find the original seva in our enriched array
            const targetSeva = enrichedSevas.find(s => {
              const id = s.seva_id || s.SevaID || s.id || s.ID
              return id === sevaId
            })
            
            if (targetSeva) {
              // Apply the name to all possible fields
              targetSeva.seva_name = name
              targetSeva.SevaName = name
              if (!targetSeva.name) targetSeva.name = name
              if (!targetSeva.Name) targetSeva.Name = name
              
              // If seva has a seva object, ensure it has the name
              if (targetSeva.seva && typeof targetSeva.seva === 'object') {
                targetSeva.seva.name = name
              }
              
              // Cache the name
              sevaNameCache[sevaId] = name
            }
          }
        } catch (err) {
          console.warn(`⚠️ Error fetching details for seva ${sevaId}:`, err)
        }
      }
      
      // Save the updated cache
      saveSevaNameCache()
    }
  }
  
  // Fourth attempt: Hardcode default names for common seva types
  // This is a fallback if all else fails
  const defaultSevaNames = {
    '1': 'Archana Seva',
    '2': 'Abhishekam',
    '3': 'Puja Seva',
    '4': 'Homam'
  }
  
  enrichedSevas.forEach(seva => {
    const sevaId = seva.seva_id || seva.SevaID || seva.id || seva.ID
    
    // Only apply default names if we still don't have a name
    if (sevaId && 
        !seva.seva_name && 
        !seva.SevaName && 
        !seva.name && 
        !seva.Name && 
        defaultSevaNames[sevaId]) {
      
      // Apply the default name
      seva.seva_name = defaultSevaNames[sevaId]
      seva.SevaName = defaultSevaNames[sevaId]
      seva.name = defaultSevaNames[sevaId]
      seva.Name = defaultSevaNames[sevaId]
      
      // If seva has a seva object, ensure it has the name
      if (seva.seva && typeof seva.seva === 'object') {
        seva.seva.name = defaultSevaNames[sevaId]
      }
      
      // Cache the name
      sevaNameCache[sevaId] = defaultSevaNames[sevaId]
    }
  })
  
  // Save the updated cache after all attempts
  saveSevaNameCache()
  
  return enrichedSevas
}

// FIX: Improved watchers for store data
// Store watcher to track changes to store data and extract names
watch(() => sevaStore.recentSevas, (newVal) => {
  if (newVal && Array.isArray(newVal) && newVal.length > 0) {
    console.log('💾 Seva store updated with', newVal.length, 'sevas')
    
    // Process sevas to extract names
    const processedSevas = processSevaData(newVal)
    
    // Make a deep copy to prevent reference issues
    localData.recentSevas = JSON.parse(JSON.stringify(processedSevas))
    
    // Save to localStorage
    try {
      localStorage.setItem('dashboard_local_data', JSON.stringify({
        ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
        recentSevas: localData.recentSevas,
        timestamp: Date.now()
      }))
      
      // Also save the updated name cache
      saveSevaNameCache()
    } catch (e) {
      console.warn('Could not save sevas to localStorage', e)
    }
    
    // Mark seva data as loaded
    dataLoadingStatus.sevas = true
    
    // Lock counters after a short delay to ensure computed properties have updated
    setTimeout(lockCounters, 500)
  }
}, { deep: true, immediate: true })

// Watch catalog for name extraction
watch(() => sevaStore.sevaCatalog, (newVal) => {
  if (newVal && Array.isArray(newVal) && newVal.length > 0) {
    console.log('💾 Seva catalog updated with', newVal.length, 'entries')
    
    // Extract names from catalog
    newVal.forEach(seva => {
      const sevaId = seva.id || seva.ID
      const sevaName = seva.name || seva.Name
      
      if (sevaId && sevaName) {
        sevaNameCache[sevaId] = sevaName
      }
    })
    
    // Save updated name cache
    saveSevaNameCache()
    
    // Process any existing sevas with new catalog information
    if (localData.recentSevas.length > 0) {
      localData.recentSevas = processSevaData(localData.recentSevas)
    }
  }
}, { deep: true, immediate: true })

// FIX: Improved watcher for donations
watch(() => donationStore.recentDonations, (newVal) => {
  if (newVal && Array.isArray(newVal) && newVal.length > 0) {
    console.log('💾 Donation store updated with', newVal.length, 'donations')
    // Make a deep copy to prevent reference issues
    localData.recentDonations = JSON.parse(JSON.stringify(newVal))
    
    // Log and debug output to track data flow
    console.log('✅ Donations saved to localData:', localData.recentDonations.length)
    
    // Save to localStorage
    try {
      localStorage.setItem('dashboard_local_data', JSON.stringify({
        ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
        recentDonations: localData.recentDonations,
        timestamp: Date.now()
      }))
    } catch (e) {
      console.warn('Could not save donations to localStorage', e)
    }
    
    // Force a UI update by making Vue see the change
    const tempDonations = [...localData.recentDonations]
    localData.recentDonations = []
    nextTick(() => {
      localData.recentDonations = tempDonations
    })
    
    // Mark donations as loaded
    dataLoadingStatus.donations = true
  }
}, { deep: true, immediate: true })

// FIX: Improved watcher for events that updates the local event data
watch(() => upcomingEvents.value, (newVal) => {
  if (newVal && Array.isArray(newVal) && newVal.length > 0) {
    console.log('💾 Events updated with', newVal.length, 'events')
    // Update local storage with events
    localData.upcomingEvents = JSON.parse(JSON.stringify(newVal))
    
    // Save to localStorage
    try {
      localStorage.setItem('dashboard_local_data', JSON.stringify({
        ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
        upcomingEvents: localData.upcomingEvents,
        timestamp: Date.now()
      }))
    } catch (e) {
      console.warn('Could not save events to localStorage', e)
    }
    
    // Mark events as loaded
    dataLoadingStatus.events = true
  }
}, { deep: true, immediate: true })

// FIX: Enhanced and sequentially load data with seva catalog prioritization and proper Promise handling
const loadDashboardData = async () => {
  const entityId = route.params.id
  if (!entityId) {
    error.value = 'Entity ID is missing'
    loading.value = false
    return
  }

  // First restore any saved data
  restoreSavedData()

  // If we already have saved counts, show them immediately
  if (dashboardCounts.countsInitialized) {
    loading.value = false
  } else {
    loading.value = true
  }
  
  error.value = null

  try {
    // Set the entity ID for API calls
    localStorage.setItem('current_entity_id', entityId)
    localStorage.setItem('current_tenant_id', entityId)
    
    // Array to hold all loading promises
    const loadingPromises = []
    
    // FIX: Load seva catalog first as it's critical for name resolution
    // We'll wait for this to complete before proceeding with other data
    console.log('🔄 Loading seva catalog first...')
    try {
      if (sevaStore.fetchSevaCatalog) {
        await sevaStore.fetchSevaCatalog()
        console.log('✅ Seva catalog loaded with', sevaStore.sevaCatalog?.length || 0, 'items')
        
        // Extract names from catalog
        if (sevaStore.sevaCatalog && Array.isArray(sevaStore.sevaCatalog)) {
          sevaStore.sevaCatalog.forEach(seva => {
            const sevaId = seva.id || seva.ID
            const sevaName = seva.name || seva.Name
            
            if (sevaId && sevaName) {
              sevaNameCache[sevaId] = sevaName
            }
          })
          
          // Save the name cache
          saveSevaNameCache()
        }
      }
    } catch (err) {
      console.warn('⚠️ Error loading seva catalog:', err)
    }
    
    // Attempt to load any additional seva related data 
    try {
      // Try to load any backup or additional seva data from the store
      if (sevaStore.fetchSevaList && typeof sevaStore.fetchSevaList === 'function') {
        await sevaStore.fetchSevaList()
        
        // Extract names from seva list if available
        if (sevaStore.sevaList && Array.isArray(sevaStore.sevaList)) {
          sevaStore.sevaList.forEach(seva => {
            const sevaId = seva.id || seva.ID
            const sevaName = seva.name || seva.Name
            
            if (sevaId && sevaName) {
              sevaNameCache[sevaId] = sevaName
              console.log(`✓ Cached seva name: ID ${sevaId} = "${sevaName}"`)
            }
          })
          saveSevaNameCache()
        }
      }
    } catch (err) {
      console.warn('⚠️ Error loading seva list:', err)
    }
    
    // Load profile data next (for birthday info)
    console.log('🔄 Loading profile data...')
    const profilePromise = (async () => {
      try {
        if (devoteeStore.loadProfileData) {
          await devoteeStore.loadProfileData()
        } else {
          const profileResponse = await devoteeService.getProfile()
          if (profileResponse && profileResponse.data) {
            devoteeStore.profile = profileResponse.data
          }
        }
        dataLoadingStatus.profile = true
      } catch (err) {
        console.warn('⚠️ Could not load profile:', err)
        dataLoadingStatus.profile = true
      }
    })()
    
    loadingPromises.push(profilePromise)
    
    // Load entity details
    console.log('🔄 Loading entity details...')
    const entityPromise = (async () => {
      try {
        const entityData = await fetchEntityDetails(entityId)
        if (entityData) {
          entityDetails.value = entityData
        }
        dataLoadingStatus.entity = true
      } catch (err) {
        console.warn('⚠️ Could not load entity details:', err)
        // Try loading from temple store as fallback
        try {
          if (templeStore.fetchTemples) {
            await templeStore.fetchTemples()
          }
        } catch (e) {
          console.warn('⚠️ Could not load temples as fallback:', e)
        }
        dataLoadingStatus.entity = true
      }
    })()
    
    loadingPromises.push(entityPromise)
    
    // If we already have counters, we can exit loading state
    if (dashboardCounts.countsInitialized) {
      loading.value = false
    }

    // FIX: Load critical data with proper Promise handling
    
    // SEVA DATA
    console.log('🔄 Loading seva bookings...')
    const sevaPromise = (async () => {
      try {
        // Check if we already have data
        if (localData.recentSevas.length === 0) {
          if (sevaStore.fetchRecentSevas) {
            await sevaStore.fetchRecentSevas()
            
            // Wait a bit for the data to propagate to the store
            await new Promise(resolve => setTimeout(resolve, 100))
          }
          
          // If store doesn't have data, try direct method
          if (!sevaStore.recentSevas || sevaStore.recentSevas.length === 0) {
            const sevaData = await fetchMySevaBookings()
            if (Array.isArray(sevaData) && sevaData.length > 0) {
              // IMPORTANT: Apply name resolution to sevas before setting to view
              const enrichedSevas = await enrichSevaWithNames(sevaData)
              mySevaBookings.value = enrichedSevas.slice(0, 3)
              
              // Store in local data
              localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas.slice(0, 3)))
              
              // Save to localStorage
              try {
                localStorage.setItem('dashboard_local_data', JSON.stringify({
                  ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                  recentSevas: localData.recentSevas,
                  timestamp: Date.now()
                }))
              } catch (e) {
                console.warn('Could not save sevas to localStorage', e)
              }
            }
          } else {
            // IMPORTANT: Apply name resolution to store sevas
            const enrichedSevas = await enrichSevaWithNames(sevaStore.recentSevas)
            if (enrichedSevas && enrichedSevas.length > 0) {
              // Update sevaStore data with enriched names if possible
              if (Array.isArray(sevaStore.recentSevas)) {
                // This may or may not work depending on reactivity
                sevaStore.recentSevas = enrichedSevas
              }
              
              // Store in local data
              localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas))
              
              // Save to localStorage
              try {
                localStorage.setItem('dashboard_local_data', JSON.stringify({
                  ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                  recentSevas: localData.recentSevas,
                  timestamp: Date.now()
                }))
              } catch (e) {
                console.warn('Could not save sevas to localStorage', e)
              }
            }
          }
        } else {
          // We have data in localData.recentSevas, let's make sure names are resolved
          const enrichedSevas = await enrichSevaWithNames(localData.recentSevas)
          localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas))
        }
        dataLoadingStatus.sevas = true
      } catch (err) {
        console.warn('⚠️ Error loading seva data:', err)
        dataLoadingStatus.sevas = true
      }
    })()
    
    loadingPromises.push(sevaPromise)
    
    // EVENTS DATA
    console.log('🔄 Loading event data...')
    const eventPromise = (async () => {
      try {
        // Check if we already have data
        if (localData.upcomingEvents.length === 0) {
          const eventData = await fetchUpcomingEvents(entityId)
          if (Array.isArray(eventData) && eventData.length > 0) {
            upcomingEvents.value = eventData.slice(0, 3)
            localData.upcomingEvents = eventData.slice(0, 3)
            
            // Save to localStorage
            try {
              localStorage.setItem('dashboard_local_data', JSON.stringify({
                ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                upcomingEvents: localData.upcomingEvents,
                timestamp: Date.now()
              }))
            } catch (e) {
              console.warn('Could not save events to localStorage', e)
            }
          }
        }
        dataLoadingStatus.events = true
      } catch (err) {
        console.warn('⚠️ Error loading event data:', err)
        dataLoadingStatus.events = true
      }
    })()
    
    loadingPromises.push(eventPromise)
    
    // DONATION DATA - Load donations first before other data
    console.log('🔄 Loading donation data...')
    const donationPromise = (async () => {
      try {
        // Check if we already have data
        if (localData.recentDonations.length === 0) {
          if (donationStore.fetchRecentDonations) {
            // Use a more reliable approach with retry mechanism
            try {
              await donationStore.fetchRecentDonations()
              
              // Wait a bit for the data to propagate to the store
              await new Promise(resolve => setTimeout(resolve, 200))
              
              // If no donations were loaded, try again
              if (!donationStore.recentDonations || donationStore.recentDonations.length === 0) {
                console.log('⚠️ No donations loaded on first attempt, retrying...')
                await donationStore.fetchRecentDonations()
                
                // Wait again
                await new Promise(resolve => setTimeout(resolve, 200))
              }
              
              // After loading, immediately copy to local data
              if (donationStore.recentDonations && donationStore.recentDonations.length > 0) {
                localData.recentDonations = JSON.parse(JSON.stringify(donationStore.recentDonations))
                console.log('✅ Successfully loaded', localData.recentDonations.length, 'donations')
                
                // Save to localStorage
                try {
                  localStorage.setItem('dashboard_local_data', JSON.stringify({
                    ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                    recentDonations: localData.recentDonations,
                    timestamp: Date.now()
                  }))
                } catch (e) {
                  console.warn('Could not save donations to localStorage', e)
                }
              } else {
                console.warn('⚠️ No donations available after two attempts')
              }
            } catch (fetchErr) {
              console.warn('⚠️ Error in fetchRecentDonations:', fetchErr)
              
              // Try alternative loading methods if available
              if (donationStore.getDonations && typeof donationStore.getDonations === 'function') {
                try {
                  const donationsResult = await donationStore.getDonations()
                  if (donationsResult && Array.isArray(donationsResult) && donationsResult.length > 0) {
                    localData.recentDonations = JSON.parse(JSON.stringify(donationsResult))
                    console.log('✅ Loaded donations via alternative method')
                  }
                } catch (altErr) {
                  console.warn('⚠️ Alternative donation loading also failed:', altErr)
                }
              }
            }
          }
        }
        dataLoadingStatus.donations = true
      } catch (err) {
        console.warn('⚠️ Error loading donation data:', err)
        dataLoadingStatus.donations = true
      }
    })()
    
    // Add this promise first to prioritize donations loading
    loadingPromises.unshift(donationPromise)
    
    // Load dashboard data if available (for supplementary data)
    if (devoteeStore.fetchDashboardData) {
      const dashboardPromise = (async () => {
        try {
          await devoteeStore.fetchDashboardData(entityId)
        } catch (err) {
          console.warn('⚠️ Error loading dashboard data:', err)
        }
      })()
      
      loadingPromises.push(dashboardPromise)
    }
    
    // Wait for all critical data to load
    Promise.all(loadingPromises)
      .then(() => {
        console.log('✅ All data loaded successfully')
        loading.value = false
        lockCounters()
      })
      .catch((err) => {
        console.error('❌ Error loading dashboard data:', err)
        loading.value = false
        lockCounters()
      })
    
  } catch (err) {
    console.error('❌ Error in dashboard data loading:', err)
    error.value = err.message || 'Failed to load dashboard data'
    loading.value = false
  }
}

// FIX: Check if donations and events failed to load and try again
const checkAndReloadMissingData = () => {
  const needsReload = {}
  
  if (!dataLoadingStatus.donations && localData.recentDonations.length === 0 && 
      (!donationStore.recentDonations || donationStore.recentDonations.length === 0)) {
    console.log('🔄 Donations data missing, scheduling reload')
    needsReload.donations = true
  }
  
  if (!dataLoadingStatus.events && localData.upcomingEvents.length === 0 && 
      upcomingEvents.value.length === 0) {
    console.log('🔄 Events data missing, scheduling reload')
    needsReload.events = true
  }
  
  if (!dataLoadingStatus.sevas || localData.recentSevas.length === 0 || 
      mySevaBookings.value.length === 0 || 
      (mySevaBookings.value.length > 0 && !getSevaName(mySevaBookings.value[0]).includes('Archana'))) {
    console.log('🔄 Seva names missing or incorrect, scheduling reload')
    needsReload.sevas = true
  }
  
  if (Object.keys(needsReload).length > 0) {
    console.log('🔄 Reloading missing data:', needsReload)
    
    if (needsReload.donations) {
      setTimeout(async () => {
        try {
          if (donationStore.fetchRecentDonations) {
            console.log('🔄 Reloading donations from checkAndReloadMissingData...')
            await donationStore.fetchRecentDonations()
            
            // Wait a bit for the store to update
            await new Promise(resolve => setTimeout(resolve, 200))
            
            // Copy data to local storage
            if (donationStore.recentDonations && donationStore.recentDonations.length > 0) {
              localData.recentDonations = JSON.parse(JSON.stringify(donationStore.recentDonations))
              console.log('✅ Donations reloaded successfully:', localData.recentDonations.length)
              
              // Force reactive update
              const tempDonations = [...localData.recentDonations]
              localData.recentDonations = []
              setTimeout(() => {
                localData.recentDonations = tempDonations
              }, 10)
              
              // Save to localStorage
              try {
                localStorage.setItem('dashboard_local_data', JSON.stringify({
                  ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                  recentDonations: localData.recentDonations,
                  timestamp: Date.now()
                }))
              } catch (e) {
                console.warn('Could not save donations to localStorage', e)
              }
            } else {
              console.warn('⚠️ No donations found after reload attempt')
              
              // Try one more approach - direct API call if available
              if (donationStore.getDonationsByApi && typeof donationStore.getDonationsByApi === 'function') {
                try {
                  const donationsApi = await donationStore.getDonationsByApi()
                  if (donationsApi && Array.isArray(donationsApi) && donationsApi.length > 0) {
                    localData.recentDonations = JSON.parse(JSON.stringify(donationsApi))
                    console.log('✅ Retrieved donations via direct API')
                  }
                } catch (apiErr) {
                  console.warn('⚠️ Direct API donation retrieval failed:', apiErr)
                }
              }
            }
          }
        } catch (e) {
          console.warn('⚠️ Failed to reload donations:', e)
        }
      }, 500)
    }
    
    if (needsReload.events) {
      setTimeout(async () => {
        try {
          const entityId = route.params.id
          const eventData = await fetchUpcomingEvents(entityId)
          if (Array.isArray(eventData) && eventData.length > 0) {
            console.log('✅ Events reloaded successfully:', eventData.length)
            upcomingEvents.value = eventData.slice(0, 3)
            localData.upcomingEvents = eventData.slice(0, 3)
            
            // Save to localStorage
            try {
              localStorage.setItem('dashboard_local_data', JSON.stringify({
                ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
                upcomingEvents: localData.upcomingEvents,
                timestamp: Date.now()
              }))
            } catch (e) {
              console.warn('Could not save events to localStorage', e)
            }
          }
        } catch (e) {
          console.warn('⚠️ Failed to reload events:', e)
        }
      }, 700)
    }
    
    if (needsReload.sevas) {
      setTimeout(async () => {
        try {
          // First try to reload seva catalog
          if (sevaStore.fetchSevaCatalog) {
            await sevaStore.fetchSevaCatalog()
          }
          
          // Then reload the seva bookings
          if (sevaStore.fetchRecentSevas) {
            await sevaStore.fetchRecentSevas()
            
            // Wait a bit for the data to propagate
            await new Promise(resolve => setTimeout(resolve, 100))
            
            if (sevaStore.recentSevas && sevaStore.recentSevas.length > 0) {
              // Apply name resolution
              const enrichedSevas = await enrichSevaWithNames(sevaStore.recentSevas)
              
              // Update the local data and view
              mySevaBookings.value = enrichedSevas.slice(0, 3)
              localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas))
              
              console.log('✅ Sevas reloaded and enriched successfully')
              
              // Force a UI update by clearing and resetting the data
              const tempSevas = [...mySevaBookings.value]
              mySevaBookings.value = []
              setTimeout(() => {
                mySevaBookings.value = tempSevas
              }, 10)
            }
          } else {
            // Direct method
            const sevaData = await fetchMySevaBookings()
            if (Array.isArray(sevaData) && sevaData.length > 0) {
              // Apply name resolution
              const enrichedSevas = await enrichSevaWithNames(sevaData)
              
              // Update the local data and view
              mySevaBookings.value = enrichedSevas.slice(0, 3)
              localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas))
              
              console.log('✅ Sevas reloaded directly and enriched successfully')
            }
          }
          
          // Save the updated data
          try {
            localStorage.setItem('dashboard_local_data', JSON.stringify({
              ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
              recentSevas: localData.recentSevas,
              timestamp: Date.now()
            }))
          } catch (e) {
            console.warn('Could not save sevas to localStorage', e)
          }
        } catch (e) {
          console.warn('⚠️ Failed to reload sevas:', e)
        }
      }, 1000)
    }
  }
}

// Lifecycle
onMounted(async () => {
  await loadDashboardData()
  
  // FIX: Check specifically for donation data immediately
  if (localData.recentDonations.length === 0 && 
      (!donationStore.recentDonations || donationStore.recentDonations.length === 0)) {
    console.log('🔍 Donation data not loaded after initial load, forcing reload...')
    
    try {
      // Direct and forceful donation loading
      if (donationStore.fetchRecentDonations) {
        console.log('🔄 Forcing donation reload...')
        await donationStore.fetchRecentDonations()
        
        // Wait a bit to ensure the store is updated
        await new Promise(resolve => setTimeout(resolve, 300))
        
        // Check if we got data and copy it to local data
        if (donationStore.recentDonations && donationStore.recentDonations.length > 0) {
          localData.recentDonations = JSON.parse(JSON.stringify(donationStore.recentDonations))
          console.log('✅ Successfully forced donation reload with', localData.recentDonations.length, 'items')
          
          // Force UI refresh by creating a temporary copy
          const tempDonations = [...localData.recentDonations]
          localData.recentDonations = []
          setTimeout(() => {
            localData.recentDonations = tempDonations
          }, 10)
          
          // Save to localStorage
          try {
            localStorage.setItem('dashboard_local_data', JSON.stringify({
              ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
              recentDonations: localData.recentDonations,
              timestamp: Date.now()
            }))
          } catch (e) {
            console.warn('Could not save donations to localStorage', e)
          }
        } else {
          console.warn('⚠️ Still no donations after forced reload')
        }
      }
    } catch (err) {
      console.warn('⚠️ Error in forced donation reload:', err)
    }
  }
  
  // FIX: Set up a check for missing data and reload if needed
  setTimeout(() => {
    checkAndReloadMissingData()
  }, 1500)
  
  // FIX: Add a second check for seva names specifically
  setTimeout(async () => {
    // Check if we have sevas but with generic names
    if (mySevaBookings.value.length > 0) {
      const firstSeva = mySevaBookings.value[0]
      const sevaName = getSevaName(firstSeva)
      
      // If the name looks like a generic "Seva 1" instead of "Archana Seva"
      if (sevaName.includes("Seva ") && !sevaName.includes("Archana")) {
        console.log('⚠️ Found generic seva names, applying enrichment...')
        
        // Apply name enrichment
        const enrichedSevas = await enrichSevaWithNames(mySevaBookings.value)
        
        // Update the display data
        mySevaBookings.value = enrichedSevas
        
        // Also update local storage
        localData.recentSevas = JSON.parse(JSON.stringify(enrichedSevas))
        
        // Save to localStorage
        try {
          localStorage.setItem('dashboard_local_data', JSON.stringify({
            ...JSON.parse(localStorage.getItem('dashboard_local_data') || '{}'),
            recentSevas: localData.recentSevas,
            timestamp: Date.now()
          }))
        } catch (e) {
          console.warn('Could not save enriched sevas to localStorage', e)
        }
        
        console.log('✅ Applied seva name enrichment on second pass')
      }
    }
  }, 2500)
  
  // Add a specific check for donations again after a delay
  setTimeout(async () => {
    if (localData.recentDonations.length === 0 && 
        (!donationStore.recentDonations || donationStore.recentDonations.length === 0)) {
      console.log('🔍 Donations still missing after delay, attempting final reload...')
      
      try {
        // Try one last time with a different approach
        if (donationStore.fetchRecentDonations) {
          await donationStore.fetchRecentDonations()
          
          if (donationStore.recentDonations && donationStore.recentDonations.length > 0) {
            // If successful, update local data
            localData.recentDonations = JSON.parse(JSON.stringify(donationStore.recentDonations))
            console.log('✅ Final donation reload successful')
            
            // Force reactive update
            const tempDonations = [...localData.recentDonations]
            localData.recentDonations = []
            setTimeout(() => {
              localData.recentDonations = tempDonations
            }, 10)
          }
        }
      } catch (err) {
        console.warn('⚠️ Final donation reload failed:', err)
      }
    }
  }, 3500)
  
  // Set up a delayed check to lock counters if they contain data
  setTimeout(() => {
    lockCounters()
  }, 4000)
})
</script>