<!-- src/layouts/DashboardLayout.vue -->
<template>
  <div class="min-h-screen bg-gray-50 flex flex-col">
    <!-- Header -->
    <AppHeader 
      :user="user" 
      :sidebar-open="sidebarOpen"
      @toggle-sidebar="sidebarOpen = !sidebarOpen"
      @logout="handleLogout"
    />

    <!-- Main Layout with Fixed Header Height -->
    <div class="flex flex-1 pt-0"> <!-- Removed excess padding -->
      <!-- Sidebar -->
      <AppSidebar 
        :isOpen="sidebarOpen"
        :user="user"
      />

      <!-- Main Content -->
      <div class="flex-1 flex flex-col min-w-0 ml-72"> <!-- Extended from ml-64 to ml-72 (288px to match sidebar width) -->
        <!-- Page Content -->
        <main class="flex-1 p-4 lg:p-6 xl:p-8 pt-0"> <!-- Removed mt-16 and added pt-0 -->
          <router-view />
        </main>
      </div>
    </div>

    <!-- Global Modals & Toasts -->
    <BaseToast />
    <BaseModal />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import AppHeader from '@/components/layout/AppHeader.vue'
import AppSidebar from '@/components/layout/AppSidebar.vue'
import BaseToast from '@/components/common/BaseToast.vue'
import BaseModal from '@/components/common/BaseModal.vue'

const router = useRouter()
const authStore = useAuthStore()

// Sidebar state - always visible
const sidebarOpen = ref(true)

// User data from auth store
const user = computed(() => authStore.user)

// Logout handler
const handleLogout = async () => {
  try {
    await authStore.logout()
    router.push('/login')
  } catch (error) {
    console.error('Logout failed:', error)
  }
}
</script>