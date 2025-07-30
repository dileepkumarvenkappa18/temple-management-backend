// Approval functions
const approveTenant = (id) => {
  console.log(`Approving tenant ${id}`)
  // Handle tenant approval logic
}

const rejectTenant = (id) => {
  console.log(`Rejecting tenant ${id}`)
  // Handle tenant rejection logic
}

const approveTemple = (id) => {
  console.log(`Approving temple ${id}`)
  // Handle temple approval logic
}

const rejectTemple = (id) => {
  console.log(`Rejecting temple ${id}`)
  // Handle temple rejection logic
}<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-50">
    <!-- Header -->
    <header class="bg-white shadow-lg border-b border-indigo-100 sticky top-0 z-40">
      <div class="px-6 py-4">
        <div class="flex items-center justify-between">
          <!-- Logo & Title -->
          <div class="flex items-center space-x-4">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
              </svg>
            </div>
            <div>
              <h1 class="text-xl font-bold text-gray-900 font-roboto">Temple SaaS</h1>
              <p class="text-sm text-indigo-600 font-medium font-sans">Super Admin Portal</p>
            </div>
          </div>

          <!-- Header Actions -->
          <div class="flex items-center space-x-4">
            <button class="p-2 text-gray-600 hover:text-gray-900 hover:bg-indigo-50 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 17h5l-5 5-5-5h5v-5a7.07 7.07 0 017-7 7.07 7.07 0 017 7v5z" />
              </svg>
            </button>
            <button class="p-2 text-gray-600 hover:text-gray-900 hover:bg-indigo-50 rounded-lg transition-colors">
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10.325 4.317c.426-1.756 2.924-1.756 3.35 0a1.724 1.724 0 002.573 1.066c1.543-.94 3.31.826 2.37 2.37a1.724 1.724 0 001.065 2.572c1.756.426 1.756 2.924 0 3.35a1.724 1.724 0 00-1.066 2.573c.94 1.543-.826 3.31-2.37 2.37a1.724 1.724 0 00-2.572 1.065c-.426 1.756-2.924 1.756-3.35 0a1.724 1.724 0 00-2.573-1.066c-1.543.94-3.31-.826-2.37-2.37a1.724 1.724 0 00-1.065-2.572c-1.756-.426-1.756-2.924 0-3.35a1.724 1.724 0 001.066-2.573c-.94-1.543.826-3.31 2.37-2.37.996.608 2.296.07 2.572-1.065z" />
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
              </svg>
            </button>

            <!-- Profile Dropdown -->
            <div class="relative" ref="profileDropdownRef">
              <button 
                @click="toggleProfileDropdown"
                class="flex items-center space-x-3 p-2 rounded-xl hover:bg-indigo-50 transition-all duration-200"
              >
                <div class="w-8 h-8 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-lg flex items-center justify-center">
                  <span class="text-white font-semibold text-sm">SA</span>
                </div>
                <div class="text-left">
                  <p class="text-sm font-semibold text-gray-900 font-sans">Super Admin</p>
                  <p class="text-xs text-gray-500 font-sans">admin@templesaas.com</p>
                </div>
                <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                </svg>
              </button>

              <!-- Dropdown Menu -->
              <div v-if="showProfileDropdown" 
                   class="absolute right-0 mt-2 w-48 bg-white rounded-xl shadow-lg border border-gray-100 py-2 z-50">
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors font-sans">
                  Profile Settings
                </a>
                <a href="#" class="block px-4 py-2 text-sm text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-colors font-sans">
                  System Logs
                </a>
                <hr class="my-2 border-gray-100">
                <button @click="logout" class="block w-full text-left px-4 py-2 text-sm text-red-600 hover:bg-red-50 transition-colors font-sans">
                  Logout
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <div class="flex">
      <!-- Sidebar -->
      <aside class="w-64 bg-white shadow-lg min-h-screen border-r border-indigo-100">
        <nav class="p-6">
          <div class="space-y-2">
            <!-- Dashboard -->
            <router-link 
              to="/superadmin/dashboard"
              class="flex items-center space-x-3 p-3 rounded-xl text-gray-700 hover:bg-indigo-50 hover:text-indigo-600 transition-all duration-200 group"
              :class="{ 'bg-indigo-50 text-indigo-600 border border-indigo-200': $route.path === '/superadmin/dashboard' }"
            >
              <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
              </svg>
              <span class="font-medium font-sans">Dashboard</span>
            </router-link>
          </div>

          <!-- Divider -->
          <hr class="my-6 border-gray-200">

          <!-- Quick Stats -->
          <div class="space-y-4">
            <h3 class="text-xs font-semibold text-gray-500 uppercase tracking-wider font-sans">Quick Overview</h3>
            
            <div class="bg-gradient-to-br from-indigo-50 to-indigo-100 p-4 rounded-xl">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-2xl font-bold text-indigo-600 font-sans">24</p>
                  <p class="text-sm text-indigo-600 font-sans">Active Temples</p>
                </div>
                <div class="w-10 h-10 bg-indigo-600 rounded-lg flex items-center justify-center">
                  <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                </div>
              </div>
            </div>

            <div class="bg-gradient-to-br from-green-50 to-green-100 p-4 rounded-xl">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-2xl font-bold text-green-600 font-sans">1,247</p>
                  <p class="text-sm text-green-600 font-sans">Total Users</p>
                </div>
                <div class="w-10 h-10 bg-green-600 rounded-lg flex items-center justify-center">
                  <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z" />
                  </svg>
                </div>
              </div>
            </div>
          </div>
        </nav>
      </aside>

      <!-- Main Content Area -->
      <main class="flex-1 p-8">
        <!-- Dashboard Content (when route is dashboard) -->
        <div v-if="$route.path === '/superadmin/dashboard' || !$route.path.includes('/superadmin/')">
          <div class="mb-6">
            <h2 class="text-2xl font-bold text-gray-900 font-roboto">Super Admin Dashboard</h2>
            <p class="text-gray-600 mt-1 font-sans">Manage tenant and temple approvals</p>
          </div>

          <!-- Stats Cards -->
          <div class="grid grid-cols-1 md:grid-cols-3 gap-6 mb-8">
            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-gray-600 font-sans">Pending Tenants</p>
                  <p class="text-2xl font-bold text-gray-900 font-sans">12</p>
                </div>
                <div class="p-3 bg-yellow-100 rounded-full">
                  <svg class="w-6 h-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                  </svg>
                </div>
              </div>
            </div>

            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-gray-600 font-sans">Pending Temples</p>
                  <p class="text-2xl font-bold text-gray-900 font-sans">8</p>
                </div>
                <div class="p-3 bg-blue-100 rounded-full">
                  <svg class="w-6 h-6 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                  </svg>
                </div>
              </div>
            </div>

            <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
              <div class="flex items-center justify-between">
                <div>
                  <p class="text-sm font-medium text-gray-600 font-sans">Total Active</p>
                  <p class="text-2xl font-bold text-gray-900 font-sans">156</p>
                </div>
                <div class="p-3 bg-green-100 rounded-full">
                  <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                  </svg>
                </div>
              </div>
            </div>
          </div>

          <!-- Pending Approvals -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
            <!-- Tenant Approvals -->
            <div class="bg-white rounded-xl shadow-sm border border-gray-200">
              <div class="p-6 border-b border-gray-200">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">Pending Tenant Approvals</h3>
              </div>
              <div class="p-6">
                <div class="space-y-4">
                  <div v-for="item in [1, 2, 3]" :key="item" class="flex items-center justify-between p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                    <div>
                      <h4 class="font-medium text-gray-900 font-sans">Temple Admin {{ item }}</h4>
                      <p class="text-sm text-gray-600 font-sans">admin{{ item }}@temple.com</p>
                      <p class="text-xs text-gray-500 font-sans">Applied 2 days ago</p>
                    </div>
                    <div class="flex space-x-2">
                      <button @click="approveTenant(item)" class="px-3 py-1 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors font-sans">
                        Approve
                      </button>
                      <button @click="rejectTenant(item)" class="px-3 py-1 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700 transition-colors font-sans">
                        Reject
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>

            <!-- Temple Approvals -->
            <div class="bg-white rounded-xl shadow-sm border border-gray-200">
              <div class="p-6 border-b border-gray-200">
                <h3 class="text-lg font-semibold text-gray-900 font-roboto">Pending Temple Approvals</h3>
              </div>
              <div class="p-6">
                <div class="space-y-4">
                  <div v-for="item in [1, 2]" :key="item" class="flex flex-col p-4 border border-gray-200 rounded-lg hover:bg-gray-50 transition-colors">
                    <div class="flex items-center justify-between mb-2">
                      <h4 class="font-medium text-gray-900 font-sans">Sri Venkateswara Temple {{ item }}</h4>
                      <div class="flex space-x-2">
                        <button @click="approveTemple(item)" class="px-3 py-1 bg-green-600 text-white text-sm rounded-lg hover:bg-green-700 transition-colors font-sans">
                          Approve
                        </button>
                        <button @click="rejectTemple(item)" class="px-3 py-1 bg-red-600 text-white text-sm rounded-lg hover:bg-red-700 transition-colors font-sans">
                          Reject
                        </button>
                      </div>
                    </div>
                    <p class="text-sm text-gray-600 font-sans">Bangalore, Karnataka</p>
                    <p class="text-xs text-gray-500 font-sans">Submitted 1 day ago</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- Router View for other pages -->
        <router-view v-else />
      </main>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const showProfileDropdown = ref(false)
const profileDropdownRef = ref(null)

const logout = () => {
  // Handle logout logic
  router.push('/login')
}

const toggleProfileDropdown = () => {
  showProfileDropdown.value = !showProfileDropdown.value
}

// Close dropdown when clicking outside
const handleClickOutside = (event) => {
  if (profileDropdownRef.value && !profileDropdownRef.value.contains(event.target)) {
    showProfileDropdown.value = false
  }
}

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style>
@import url('https://fonts.googleapis.com/css2?family=Roboto:wght@400;700&display=swap');

.font-roboto {
  font-family: 'Roboto', sans-serif;
  font-weight: 700;
}

.font-sans {
  font-family: sans-serif;
}
</style>