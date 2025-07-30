<!-- src/layouts/AuthLayout.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100 flex flex-col">
    <!-- Minimal Header -->
    <header class="bg-white/80 backdrop-blur-sm border-b border-indigo-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <router-link to="/" class="flex items-center space-x-3">
            <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2L2 7v10c0 5.55 3.84 9.739 9 11 5.16-1.261 9-5.45 9-11V7l-10-5z"/>
              </svg>
            </div>
            <span class="text-xl font-bold text-gray-900">TempleSaaS</span>
          </div>
          
          </router-link>

          <!-- Back to Home -->
          <router-link 
            to="/" 
            class="text-gray-600 hover:text-indigo-600 transition-colors font-medium flex items-center space-x-2"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
            </svg>
            <span>Back to Home</span>
          </router-link>
        </div>
      </div>
    </header>

    <!-- Main Content Area -->
    <main class="flex-1 flex items-center justify-center px-4 sm:px-6 lg:px-8 py-12">
      <div class="w-full max-w-md">
        <!-- Auth Card Container -->
        <div class="bg-white/90 backdrop-blur-sm rounded-2xl shadow-xl border border-white/20 p-8">
          <!-- Content from router-view (LoginView, RegisterView, etc.) -->
          <router-view />
        </div>

        <!-- Additional Links/Info -->
        <div class="mt-6 text-center">
          <p class="text-sm text-gray-600">
            Need help? 
            <a href="mailto:support@templesync.com" class="text-indigo-600 hover:text-indigo-700 font-medium transition-colors">
              Contact Support
            </a>
          </p>
        </div>
      </div>
    </main>

    <!-- Minimal Footer -->
    <footer class="bg-white/80 backdrop-blur-sm border-t border-indigo-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-4">
        <div class="flex flex-col sm:flex-row justify-between items-center space-y-2 sm:space-y-0">
          <p class="text-sm text-gray-600">
            &copy; 2025 TempleSync. All rights reserved.
          </p>
          <div class="flex items-center space-x-6">
            <router-link to="/terms" class="text-sm text-gray-600 hover:text-indigo-600 transition-colors">
              Terms
            </router-link>
            <router-link to="/privacy" class="text-sm text-gray-600 hover:text-indigo-600 transition-colors">
              Privacy
            </router-link>
          </div>
        </div>
      </div>
    </footer>

    <!-- Decorative Background Elements -->
    <div class="fixed inset-0 -z-10 overflow-hidden pointer-events-none">
      <div class="absolute -top-40 -right-40 w-80 h-80 bg-indigo-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse"></div>
      <div class="absolute -bottom-40 -left-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-xl opacity-20 animate-pulse" style="animation-delay: 2s;"></div>
      <div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-xl opacity-10 animate-pulse" style="animation-delay: 4s;"></div>
    </div>
  </div>
</template>

<script setup>
// No logic needed - this is a pure layout component
// The magic happens in the child routes (LoginView, RegisterView, etc.)
</script>