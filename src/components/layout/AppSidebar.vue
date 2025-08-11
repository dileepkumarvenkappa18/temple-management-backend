<template>
  <aside class="fixed left-0 top-0 bottom-0 w-64 bg-white border-r border-gray-200 overflow-y-auto shadow-md z-40" style="margin-top: 48px; display: block !important;">
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

        <router-link to="/superadmin/tenant-approvals" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/tenants') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
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

        <!-- <router-link to="/superadmin/settings" class="flex items-center px-3 py-2 text-sm font-medium rounded-md" :class="isActiveRoute('/superadmin/settings') ? 'bg-indigo-100 text-indigo-700' : 'text-gray-700 hover:bg-gray-50'">
          <svg class="mr-3 h-5 w-5" :class="isActiveRoute('/superadmin/settings') ? 'text-indigo-500' : 'text-gray-400'" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
          </svg>
          System Settings
        </router-link> -->
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
import { computed } from 'vue';
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

// Methods
const isActiveRoute = (path) => {
  return route.path.startsWith(path);
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
</style>