<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 text-center">
      <!-- Animated 404 Icon -->
      <div class="relative">
        <div class="absolute inset-0 bg-indigo-200 rounded-full opacity-20 animate-pulse"></div>
        <div class="relative bg-white rounded-full p-8 shadow-lg border border-indigo-100">
          <div class="flex items-center justify-center">
            <!-- Temple Icon -->
            <svg class="w-24 h-24 text-indigo-600 animate-bounce" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="1.5" 
                    d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
            </svg>
          </div>
        </div>
      </div>

      <!-- Error Code -->
      <div class="space-y-2">
        <h1 class="text-8xl font-black text-indigo-600 tracking-tight">
          404
        </h1>
        <div class="h-1 w-32 bg-gradient-to-r from-indigo-500 to-indigo-700 rounded-full mx-auto"></div>
      </div>

      <!-- Error Message -->
      <div class="space-y-4">
        <h2 class="text-2xl font-bold text-gray-900 sm:text-3xl">
          Page Not Found
        </h2>
        <p class="text-gray-600 text-lg leading-relaxed">
          The page you're looking for doesn't exist or has been moved. 
          Let's get you back to your spiritual journey.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="space-y-4 pt-6">
        <button 
          @click="goHome"
          class="w-full flex items-center justify-center px-8 py-3 border border-transparent text-base font-medium rounded-xl text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transform transition-all duration-200 hover:scale-105 shadow-lg hover:shadow-xl"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
          </svg>
          Go to Home
        </button>

        <button 
          @click="goBack"
          class="w-full flex items-center justify-center px-8 py-3 border-2 border-indigo-200 text-base font-medium rounded-xl text-indigo-700 bg-white hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transform transition-all duration-200 hover:scale-105"
        >
          <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18" />
          </svg>
          Go Back
        </button>
      </div>

      <!-- Help Text -->
      <div class="pt-8 border-t border-indigo-100">
        <p class="text-sm text-gray-500">
          Still having trouble? 
          <a href="#" class="font-medium text-indigo-600 hover:text-indigo-500 transition-colors duration-200">
            Contact support
          </a>
        </p>
      </div>

      <!-- Decorative Elements -->
      <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
        <!-- Floating Circles -->
        <div class="absolute top-20 left-10 w-4 h-4 bg-indigo-300 rounded-full opacity-20 animate-pulse"></div>
        <div class="absolute top-40 right-16 w-6 h-6 bg-indigo-400 rounded-full opacity-15 animate-pulse" style="animation-delay: 1s;"></div>
        <div class="absolute bottom-32 left-20 w-3 h-3 bg-indigo-500 rounded-full opacity-25 animate-pulse" style="animation-delay: 2s;"></div>
        <div class="absolute bottom-20 right-8 w-5 h-5 bg-indigo-300 rounded-full opacity-20 animate-pulse" style="animation-delay: 0.5s;"></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const router = useRouter()

const goHome = () => {
  router.push('/')
}

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}
</script>

<style scoped>
/* Custom animations */
@keyframes bounce {
  0%, 20%, 53%, 80%, 100% {
    transform: translate3d(0,0,0);
  }
  40%, 43% {
    transform: translate3d(0, -8px, 0);
  }
  70% {
    transform: translate3d(0, -4px, 0);
  }
  90% {
    transform: translate3d(0, -2px, 0);
  }
}

@keyframes pulse {
  0% {
    opacity: 0.4;
    transform: scale(1);
  }
  50% {
    opacity: 0.8;
    transform: scale(1.05);
  }
  100% {
    opacity: 0.4;
    transform: scale(1);
  }
}

.animate-bounce {
  animation: bounce 2s infinite;
}

.animate-pulse {
  animation: pulse 3s infinite;
}

/* Responsive text scaling */
@media (max-width: 640px) {
  h1 {
    font-size: 4rem;
  }
}

@media (max-width: 480px) {
  h1 {
    font-size: 3rem;
  }
}
</style>