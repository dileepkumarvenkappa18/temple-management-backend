<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100">
    <!-- Navigation Header -->
    <nav class="bg-white/80 backdrop-blur-md border-b border-indigo-100 sticky top-0 z-50">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <!-- Logo -->
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2L2 7v10c0 5.55 3.84 9.739 9 11 5.16-1.261 9-5.45 9-11V7l-10-5z"/>
              </svg>
            </div>
            <span class="text-xl font-bold text-gray-900">Devotee Connect</span>
          </div>
          
          <!-- Navigation Links -->
          <div class="hidden md:flex items-center space-x-8">
            <a href="#features" class="text-gray-600 hover:text-indigo-600 transition-colors duration-200">Features</a>
            <a href="#about" class="text-gray-600 hover:text-indigo-600 transition-colors duration-200">About</a>
            <a href="#contact" class="text-gray-600 hover:text-indigo-600 transition-colors duration-200">Contact</a>
          </div>

          <!-- Auth Buttons (Desktop only) -->
          <div class="hidden md:flex items-center space-x-4">
            <router-link 
              to="/login" 
              class="text-indigo-600 hover:text-indigo-700 font-medium transition-colors duration-200"
            >
              Login
            </router-link>
            <router-link 
              to="/register" 
              class="bg-indigo-600 hover:bg-indigo-700 text-white px-6 py-2 rounded-xl font-medium transition-all duration-200 hover:shadow-lg hover:shadow-indigo-500/25"
            >
              Register
            </router-link>
          </div>
          
          <!-- Mobile Menu Button -->
          <button 
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden p-2 rounded-lg hover:bg-gray-100 transition-colors duration-200"
          >
            <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16"/>
            </svg>
          </button>
        </div>
        
        <!-- Mobile Menu -->
        <transition name="fade">
          <div
            v-if="mobileMenuOpen"
            class="md:hidden mt-4 bg-white shadow-lg rounded-xl border border-indigo-100 px-6 py-6 space-y-4 transition-all duration-300"
          >
            <a href="#features" class="block text-gray-700 font-medium hover:text-indigo-600 transition-colors duration-200">Features</a>
            <a href="#about" class="block text-gray-700 font-medium hover:text-indigo-600 transition-colors duration-200">About</a>
            <a href="#contact" class="block text-gray-700 font-medium hover:text-indigo-600 transition-colors duration-200">Contact</a>
            <div class="pt-4 border-t border-gray-200">
              <router-link 
                to="/login" 
                class="block text-indigo-600 hover:text-indigo-700 font-medium transition duration-200 mb-2"
              >
                Login
              </router-link>
              <router-link 
                to="/register" 
                class="block bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-xl font-medium transition-all duration-200 text-center"
              >
                Register
              </router-link>
            </div>
          </div>
        </transition>
      </div>
    </nav>

    <!-- Hero Section -->
    <section class="relative py-20 px-4 sm:px-6 lg:px-8">
      <div class="max-w-7xl mx-auto">
        <div class="text-center">
          <!-- Hero Badge -->
          <div class="inline-flex items-center px-4 py-2 bg-indigo-100 rounded-full text-indigo-800 text-sm font-medium mb-8">
            <svg class="w-4 h-4 mr-2" fill="currentColor" viewBox="0 0 24 24">
              <path d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"/>
            </svg>
            Trusted by 500+ Temples Worldwide
          </div>
          
          <!-- Main Heading -->
          <h1 class="text-4xl md:text-6xl font-bold text-gray-900 mb-6 leading-tight">
            Manage your temple,
            <span class="text-indigo-600">devotees, sevas</span>
            <br class="hidden sm:block">
            & donations with ease
          </h1>
          
          <!-- Subtitle -->
          <p class="text-xl text-gray-600 mb-10 max-w-3xl mx-auto leading-relaxed">
            Streamline your temple operations with our comprehensive SaaS platform. 
            From devotee management to seva bookings, donations, and events - all in one place.
          </p>
          
          <!-- CTA Buttons -->
          <div class="flex flex-col sm:flex-row gap-4 justify-center items-center">
            <router-link 
              to="/register" 
              class="bg-indigo-600 hover:bg-indigo-700 text-white px-8 py-4 rounded-xl text-lg font-semibold transition-all duration-200 hover:shadow-xl hover:shadow-indigo-500/25 hover:-translate-y-1"
            >
              Get Started
            </router-link>
            <router-link 
              to="/login" 
              class="border-2 border-indigo-600 text-indigo-600 hover:bg-indigo-600 hover:text-white px-8 py-4 rounded-xl text-lg font-semibold transition-all duration-200 hover:shadow-lg"
            >
              Already User?
            </router-link>
          </div>
        </div>
      </div>
      
      <!-- Floating Elements -->
      <div class="absolute top-20 left-10 w-20 h-20 bg-indigo-200 rounded-full opacity-50 animate-pulse"></div>
      <div class="absolute bottom-20 right-10 w-16 h-16 bg-indigo-300 rounded-full opacity-60 animate-pulse delay-1000"></div>
    </section>

    <!-- Features Section -->
    <section id="features" class="py-20 bg-white">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="text-center mb-16">
          <h2 class="text-3xl md:text-4xl font-bold text-gray-900 mb-4">
            Everything you need to manage your temple
          </h2>
          <p class="text-xl text-gray-600 max-w-2xl mx-auto">
            Comprehensive tools designed specifically for temple administration and devotee engagement
          </p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-8">
          <!-- Feature Cards -->
          <div 
            v-for="feature in features" 
            :key="feature.id"
            class="bg-gradient-to-br from-white to-indigo-50 p-8 rounded-2xl border border-indigo-100 hover:shadow-xl hover:shadow-indigo-500/10 transition-all duration-300 hover:-translate-y-2"
          >
            <div class="w-14 h-14 bg-indigo-600 rounded-xl flex items-center justify-center mb-6">
              <svg class="w-7 h-7 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path :d="feature.icon"/>
              </svg>
            </div>
            <h3 class="text-xl font-bold text-gray-900 mb-3">{{ feature.title }}</h3>
            <p class="text-gray-600 leading-relaxed">{{ feature.description }}</p>
          </div>
        </div>
      </div>
    </section>

    <!-- Stats Section -->
    <section class="py-20 bg-indigo-600">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-2 md:grid-cols-4 gap-8 text-center">
          <div v-for="stat in stats" :key="stat.label" class="text-white">
            <div class="text-3xl md:text-4xl font-bold mb-2">{{ stat.value }}</div>
            <div class="text-indigo-200">{{ stat.label }}</div>
          </div>
        </div>
      </div>
    </section>

    <!-- Footer -->
    <footer class="bg-gray-900 text-white py-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="grid grid-cols-1 md:grid-cols-4 gap-8">
          <!-- Company Info -->
          <div class="md:col-span-2">
            <div class="flex items-center space-x-3 mb-4">
              <div class="w-10 h-10 bg-indigo-600 rounded-xl flex items-center justify-center">
                <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M12 2L2 7v10c0 5.55 3.84 9.739 9 11 5.16-1.261 9-5.45 9-11V7l-10-5z"/>
                </svg>
              </div>
              <span class="text-xl font-bold">Devotee Connect</span>
            </div>
            <p class="text-gray-400 mb-4 max-w-md">
              Empowering temples worldwide with modern technology for better devotee management and spiritual community building.
            </p>
          </div>
          
          <!-- Quick Links -->
          <div>
            <h4 class="font-semibold mb-4">Quick Links</h4>
            <ul class="space-y-2 text-gray-400">
              <li><a href="#features" class="hover:text-white transition-colors duration-200">Features</a></li>
              <li><a href="#about" class="hover:text-white transition-colors duration-200">About</a></li>
              <li><a href="#contact" class="hover:text-white transition-colors duration-200">Contact</a></li>
            </ul>
          </div>
          
          <!-- Legal -->
          <div>
            <h4 class="font-semibold mb-4">Legal</h4>
            <ul class="space-y-2 text-gray-400">
              <li><router-link to="/terms" class="hover:text-white transition-colors duration-200">Terms of Service</router-link></li>
              <li><router-link to="/privacy" class="hover:text-white transition-colors duration-200">Privacy Policy</router-link></li>
            </ul>
          </div>
        </div>
        
        <div class="border-t border-gray-800 mt-8 pt-8 text-center text-gray-400">
          <p>&copy; 2025 Devotee Connect. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'

const mobileMenuOpen = ref(false)

const features = [
  {
    id: 1,
    title: 'Devotee Management',
    description: 'Comprehensive devotee profiles with family lineage, spiritual information, and preferences management.',
    icon: 'M12 4.354a4 4 0 110 5.292M15 21H3v-1a6 6 0 0112 0v1zm0 0h6v-1a6 6 0 00-9-5.197m13.5-9a2.5 2.5 0 11-5 0 2.5 2.5 0 015 0z'
  },
  {
    id: 2,
    title: 'Seva Bookings',
    description: 'Easy online seva booking system with calendar integration and automated confirmations.',
    icon: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z'
  },
  {
    id: 3,
    title: 'Donation Tracking',
    description: 'Secure donation processing with automated receipts and comprehensive financial reporting.',
    icon: 'M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1'
  },
  {
    id: 4,
    title: 'Event Management',
    description: 'Organize festivals and events with RSVP tracking and automated devotee notifications.',
    icon: 'M11 3.055A9.001 9.001 0 1020.945 13H11V3.055z M20.488 9H15V3.512A9.025 9.025 0 0120.488 9z'
  },
  {
    id: 5,
    title: 'Communication Hub',
    description: 'Bulk messaging system with SMS, WhatsApp, and email integration for devotee outreach.',
    icon: 'M7 8h10M7 12h4m1 8l-4-4H5a2 2 0 01-2-2V6a2 2 0 012-2h14a2 2 0 012 2v8a2 2 0 01-2 2h-3l-4 4z'
  },
  {
    id: 6,
    title: 'Multi-Role Access',
    description: 'Role-based permissions for temple admins, volunteers, and devotees with secure access control.',
    icon: 'M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.03 9-11.622 0-1.042-.133-2.052-.382-3.016z'
  }
]

const stats = [
  { value: '500+', label: 'Active Temples' },
  { value: '50K+', label: 'Devotees Managed' },
  { value: '1M+', label: 'Sevas Booked' },
  { value: '99.9%', label: 'Uptime' }
]

// Close mobile menu when clicking outside or on links
const closeMobileMenu = () => {
  mobileMenuOpen.value = false
}

// Add click outside listener
import { onMounted, onUnmounted } from 'vue'

onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})

const handleClickOutside = (event) => {
  if (mobileMenuOpen.value && !event.target.closest('nav')) {
    closeMobileMenu()
  }
}
</script>

<style scoped>
.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s ease;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>