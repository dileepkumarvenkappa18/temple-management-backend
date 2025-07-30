<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-md w-full space-y-8 text-center">
      <!-- Error Icon -->
      <div class="mx-auto flex items-center justify-center h-32 w-32 rounded-full bg-gradient-to-br from-red-100 to-red-200 shadow-lg">
        <svg class="h-16 w-16 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor" stroke-width="2">
          <path stroke-linecap="round" stroke-linejoin="round" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z" />
        </svg>
      </div>

      <!-- Error Code -->
      <div class="space-y-2">
        <h1 class="text-8xl font-bold text-indigo-600 tracking-tight">
          401
        </h1>
        <div class="h-1 w-20 bg-gradient-to-r from-indigo-500 to-indigo-700 mx-auto rounded-full"></div>
      </div>

      <!-- Error Message -->
      <div class="space-y-4">
        <h2 class="text-3xl font-bold text-gray-900 font-roboto">
          Access Denied
        </h2>
        <p class="text-lg text-gray-600 leading-relaxed font-sans">
          Sorry, you don't have permission to access this page. Please contact your administrator or try logging in with appropriate credentials.
        </p>
      </div>

      <!-- Action Buttons -->
      <div class="space-y-4 pt-6">
        <div class="flex flex-col sm:flex-row gap-4 justify-center">
          <!-- Go Back Button -->
          <button
            @click="goBack"
            class="inline-flex items-center justify-center px-6 py-3 border border-transparent text-base font-medium rounded-xl text-white bg-indigo-600 hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
            </svg>
            Go Back
          </button>

          <!-- Login Button -->
          <button
            @click="goToLogin"
            class="inline-flex items-center justify-center px-6 py-3 border-2 border-indigo-600 text-base font-medium rounded-xl text-indigo-600 bg-white hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500 transition-all duration-200 shadow-md hover:shadow-lg transform hover:-translate-y-0.5"
          >
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path>
            </svg>
            Login
          </button>
        </div>

        <!-- Home Link -->
        <div class="pt-4">
          <button
            @click="goToHome"
            class="text-indigo-600 hover:text-indigo-800 font-medium transition-colors duration-200 underline decoration-2 underline-offset-4 hover:decoration-indigo-800"
          >
            Return to Homepage
          </button>
        </div>
      </div>

      <!-- Additional Help -->
      <div class="pt-8 border-t border-gray-200">
        <div class="bg-indigo-50 rounded-xl p-6 shadow-sm">
          <h3 class="text-lg font-semibold text-indigo-900 mb-2 font-roboto">
            Need Help?
          </h3>
          <p class="text-indigo-700 text-sm mb-4 font-sans">
            If you believe this is an error, please contact support with the following details:
          </p>
          <div class="bg-white rounded-lg p-3 text-left">
            <div class="text-xs text-gray-600 space-y-1 font-mono">
              <div><span class="font-semibold">Error Code:</span> 401</div>
              <div><span class="font-semibold">Timestamp:</span> {{ currentTime }}</div>
              <div><span class="font-semibold">Page:</span> {{ currentPath }}</div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'

const router = useRouter()
const route = useRoute()

const currentTime = ref('')
const currentPath = ref('')

onMounted(() => {
  currentTime.value = new Date().toLocaleString()
  currentPath.value = route.fullPath
})

const goBack = () => {
  if (window.history.length > 1) {
    router.go(-1)
  } else {
    router.push('/')
  }
}

const goToLogin = () => {
  router.push('/login')
}

const goToHome = () => {
  router.push('/')
}
</script>

<style scoped>
/* Custom fonts - Add to your main CSS or import */
.font-roboto {
  font-family: 'Roboto', sans-serif;
}

.font-sans {
  font-family: system-ui, -apple-system, sans-serif;
}

/* Custom animations */
@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.space-y-8 > * {
  animation: fadeInUp 0.6s ease-out forwards;
}

.space-y-8 > *:nth-child(2) {
  animation-delay: 0.1s;
}

.space-y-8 > *:nth-child(3) {
  animation-delay: 0.2s;
}

.space-y-8 > *:nth-child(4) {
  animation-delay: 0.3s;
}

.space-y-8 > *:nth-child(5) {
  animation-delay: 0.4s;
}

.space-y-8 > *:nth-child(6) {
  animation-delay: 0.5s;
}

/* Responsive adjustments */
@media (max-width: 640px) {
  .text-8xl {
    font-size: 4rem;
  }
  
  .text-3xl {
    font-size: 1.875rem;
  }
  
  .text-lg {
    font-size: 1rem;
  }
}
</style>