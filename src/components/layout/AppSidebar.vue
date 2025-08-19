<template>
  <aside class="fixed left-0 top-0 bottom-0 w-72 bg-white border-r border-gray-200 overflow-y-auto shadow-md z-40" style="margin-top: 48px; display: block !important;">
    <!-- Logo Area with INCREASED top padding -->
    <div class="p-4 pt-8 border-b border-gray-200"> <!-- Added pt-8 for extra top padding -->
      <h3 class="text-lg font-semibold text-indigo-600">Temple Management</h3>
      <p v-if="actualRole" class="text-sm text-gray-500">{{ actualRole }}</p>
    </div>
    
    <!-- Navigation Menu -->
    <nav class="px-4 py-4">
      <!-- TENANT Navigation -->
      <div v-if="actualRole === 'tenant'" class="space-y-1">
        <router-link to="/tenant/dashboard" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/tenant/dashboard') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/tenant/dashboard') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </router-link>

        <router-link to="/tenant/dashboard" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/tenant/entities') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/tenant/dashboard') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          My Temples
        </router-link>

        <router-link to="/tenant/entities/create" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/tenant/entities/create') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/tenant/entities/create') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create Temple
        </router-link>

        <!-- NEW: Reports Dashboard with Subsections that Link to Separate Pages -->
        <div class="mt-3">
          <!-- Reports Dashboard Main Button -->
          <div 
            class="flex items-center justify-between px-3 py-2 text-sm font-medium rounded-md bg-indigo-50 text-indigo-700 cursor-pointer"
            @click="toggleReportsDashboard">
            <div class="flex items-center">
              <svg class="mr-3 h-5 w-5 text-indigo-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 17v-2m3 2v-4m3 4v-6m2 10H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span>Reports Dashboard</span>
            </div>
            <!-- Replaced arrow SVG with text indicator -->
            <span class="text-xs font-bold text-indigo-600">
              {{ isReportsDashboardActive ? '−' : '+' }}
            </span>
          </div>

          <!-- Reports Dashboard Subsections -->
          <transition name="slide-down">
            <div v-if="isReportsDashboardActive" class="mt-1 ml-6 space-y-1 bg-indigo-50 p-2 rounded-md">
              <!-- Temple Register -->
              <router-link 
                to="/tenant/reports/temple-register" 
                class="px-3 py-2 text-sm rounded-md cursor-pointer flex items-center"
                :class="isActiveRoute('/tenant/reports/temple-register') ? 'text-indigo-700 font-medium' : 'text-gray-700 hover:bg-indigo-100'"
              >
                <svg class="mr-2 h-4 w-4" :class="isActiveRoute('/tenant/reports/temple-register') ? 'text-indigo-600' : 'text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
                </svg>
                <span>Temple Register Report</span>
              </router-link>

              <!-- Temple Activities -->
              <router-link 
                to="/tenant/reports/temple-activities" 
                class="px-3 py-2 text-sm rounded-md cursor-pointer flex items-center"
                :class="isActiveRoute('/tenant/reports/temple-activities') ? 'text-indigo-700 font-medium' : 'text-gray-700 hover:bg-indigo-100'"
              >
                <svg class="mr-2 h-4 w-4" :class="isActiveRoute('/tenant/reports/temple-activities') ? 'text-indigo-600' : 'text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                </svg>
                <span>Temple Activities Report</span>
              </router-link>

              <!-- Birthdays -->
              <router-link 
                to="/tenant/reports/birthdays" 
                class="px-3 py-2 text-sm rounded-md cursor-pointer flex items-center"
                :class="isActiveRoute('/tenant/reports/birthdays') ? 'text-indigo-700 font-medium' : 'text-gray-700 hover:bg-indigo-100'"
              >
                <svg class="mr-2 h-4 w-4" :class="isActiveRoute('/tenant/reports/birthdays') ? 'text-indigo-600' : 'text-gray-500'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 15.546c-.523 0-1.046.151-1.5.454a2.704 2.704 0 01-3 0 2.704 2.704 0 00-3 0 2.704 2.704 0 01-3 0 2.704 2.704 0 00-3 0 2.704 2.704 0 01-3 0 2.701 2.701 0 00-1.5-.454M9 6v2m3-2v2m3-2v2M9 3h.01M12 3h.01M15 3h.01M21 21v-7a2 2 0 00-2-2H5a2 2 0 00-2 2v7h18zm-3-9v-2a2 2 0 00-2-2H8a2 2 0 00-2 2v2h12z" />
                </svg>
                <span>Birthdays Report</span>
              </router-link>
            </div>
          </transition>
        </div>
      </div>

      <!-- ENTITY ADMIN Navigation -->
      <div v-else-if="actualRole === 'entity_admin'" class="space-y-1">
        <router-link :to="`/entity/${entityId}/dashboard`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/dashboard`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/dashboard`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </router-link>

        <router-link :to="`/entity/${entityId}/devotees`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotees`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotees`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          Devotees
        </router-link>

        <router-link :to="`/entity/${entityId}/sevas`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/sevas`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/sevas`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          Seva Management
        </router-link>

        <router-link :to="`/entity/${entityId}/donations`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/donations`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/donations`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          Donations
        </router-link>

        <router-link :to="`/entity/${entityId}/events`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/events`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/events`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          Events & Festivals
        </router-link>

        <router-link :to="`/entity/${entityId}/communication`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/communication`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/communication`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 10h.01M12 10h.01M16 10h.01M9 16H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-5l-5 5v-5z" />
          </svg>
          Communication
        </router-link>
      </div>

      <!-- DEVOTEE Navigation -->
      <div v-else-if="actualRole === 'devotee'" class="space-y-1">
        <router-link :to="`/entity/${entityId}/devotee/dashboard`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotee/dashboard`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotee/dashboard`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </router-link>

        <router-link :to="`/entity/${entityId}/devotee/seva-booking`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotee/seva-booking`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotee/seva-booking`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4.318 6.318a4.5 4.5 0 000 6.364L12 20.364l7.682-7.682a4.5 4.5 0 00-6.364-6.364L12 7.636l-1.318-1.318a4.5 4.5 0 00-6.364 0z" />
          </svg>
          Book Seva
        </router-link>

        <router-link :to="`/entity/${entityId}/devotee/donations`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotee/donations`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotee/donations`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
          </svg>
          My Donations
        </router-link>

        <router-link :to="`/entity/${entityId}/devotee/events`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotee/events`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotee/events`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          Temple Events
        </router-link>

        <router-link :to="`/entity/${entityId}/devotee/profile/edit`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/devotee/profile`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/devotee/profile`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          My Profile
        </router-link>
      </div>

      <!-- VOLUNTEER Navigation -->
      <div v-else-if="actualRole === 'volunteer'" class="space-y-1">
        <router-link :to="`/entity/${entityId}/volunteer/dashboard`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/volunteer/dashboard`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/volunteer/dashboard`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </router-link>

        <router-link :to="`/entity/${entityId}/volunteer/assignments`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/volunteer/assignments`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/volunteer/assignments`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          My Assignments
        </router-link>

        <router-link :to="`/entity/${entityId}/volunteer/schedule`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/volunteer/schedule`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/volunteer/schedule`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          My Schedule
        </router-link>

        <router-link :to="`/entity/${entityId}/volunteer/events`" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute(`/entity/${entityId}/volunteer/events`) ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute(`/entity/${entityId}/volunteer/events`) ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
          </svg>
          Temple Events
        </router-link>
      </div>

      <!-- SUPERADMIN Navigation -->
      <div v-else-if="actualRole === 'superadmin'" class="space-y-1">
        <router-link to="/superadmin/dashboard" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/dashboard') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/dashboard') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Dashboard
        </router-link>

        <router-link to="/superadmin/tenant-approvals" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/tenant-approvals') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/tenant-approvals') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
          </svg>
          Manage Tenants
        </router-link>

        <router-link to="/superadmin/dashboard" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/temples') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/dashboard') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
          Manage Temples
        </router-link>

        <!-- New User Management Link -->
        <router-link to="/superadmin/user-management" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/user-management') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/user-management') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197M13 7a4 4 0 11-8 0 4 4 0 018 0z" />
          </svg>
          User Management
        </router-link>

        <!-- New User Role Management Link -->
        <router-link to="/superadmin/role-management" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/role-management') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/role-management') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z" />
          </svg>
          User Role Management
        </router-link>
        
        <!-- Reset Password Link -->
        <router-link to="/superadmin/reset-password" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/reset-password') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/reset-password') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z" />
          </svg>
          Reset Password
        </router-link>
        
        <!-- NEW: Audit Logs Link -->
        <router-link to="/superadmin/audit-logs" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/audit-logs') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/audit-logs') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
          </svg>
          Audit Logs
        </router-link>
      </div>

      <!-- Default message when no role is set -->
      <div v-else class="text-center py-6">
        <div class="mx-auto w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mb-4">
          <svg class="w-8 h-8 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
          </svg>
        </div>
        <p class="text-sm text-gray-500">Menu will appear once your role is assigned.</p>
        <p class="text-xs text-gray-400 mt-2">Current path: {{ route.path }}</p>
        <p v-if="user" class="text-xs text-gray-400 mt-1">User role: {{ user.role }}</p>
      </div>
    </nav>
    
    <!-- Help & Support Section -->
    <div class="mt-6 px-4 py-4 border-t border-gray-200">
      <router-link to="/support" class="flex items-center px-3 py-2 text-sm font-medium text-gray-700 rounded-md hover:bg-gray-50">
        <svg class="mr-3 h-5 w-5 text-gray-400" fill="none" viewBox="0 0 24 24" stroke="currentColor">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Help & Support
      </router-link>
    </div>
  </aside>
</template>

<script setup>
import { computed, ref } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

// Props
const props = defineProps({
  isOpen: {
    type: Boolean,
    default: true
  },
  user: {
    type: Object,
    default: null
  }
});

// Composables
const route = useRoute();
const authStore = useAuthStore();

// Extract entity ID from route
const entityId = computed(() => {
  return route.params.id || route.params.entityId || '1';
});

// Detect role from various sources for reliability
const actualRole = computed(() => {
  // First try to get from auth store
  if (authStore.user && authStore.user.role) {
    return authStore.user.role.toLowerCase();
  }
  
  // Then try from props
  if (props.user && props.user.role) {
    return props.user.role.toLowerCase();
  }
  
  // Try to infer from route path
  const path = route.path;
  if (path.includes('/tenant/')) {
    return 'tenant';
  } else if (path.includes('/entity/') && path.includes('/devotee/')) {
    return 'devotee';
  } else if (path.includes('/entity/') && path.includes('/volunteer/')) {
    return 'volunteer';
  } else if (path.includes('/entity/') && !path.includes('/devotee/') && !path.includes('/volunteer/')) {
    return 'entity_admin';
  } else if (path.includes('/superadmin/')) {
    return 'superadmin';
  }
  
  // Default to empty if can't determine
  return '';
});

// Reports Dashboard state
const isReportsDashboardActive = ref(false);

// Methods
const isActiveRoute = (path) => {
  return route.path.startsWith(path);
};

const toggleReportsDashboard = () => {
  isReportsDashboardActive.value = !isReportsDashboardActive.value;
};
</script>

<style scoped>
/* Custom scrollbar for webkit browsers */
aside::-webkit-scrollbar {
  width: 4px;
}

aside::-webkit-scrollbar-track {
  background: #f1f5f9;
}

aside::-webkit-scrollbar-thumb {
  background: #cbd5e1;
  border-radius: 2px;
}

aside::-webkit-scrollbar-thumb:hover {
  background: #94a3b8;
}

/* Transitions */
.slide-down-enter-active,
.slide-down-leave-active {
  transition: all 0.3s ease-out;
}

.slide-down-enter-from,
.slide-down-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>