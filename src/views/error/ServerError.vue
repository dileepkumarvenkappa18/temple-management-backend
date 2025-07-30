<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100 flex items-center justify-center px-4 sm:px-6 lg:px-8">
    <div class="max-w-lg w-full text-center">
      <!-- Error Icon -->
      <div class="mb-8">
        <div class="mx-auto flex items-center justify-center h-24 w-24 rounded-full bg-red-100 shadow-lg">
          <svg 
            class="h-12 w-12 text-red-600" 
            fill="none" 
            stroke="currentColor" 
            viewBox="0 0 24 24" 
            xmlns="http://www.w3.org/2000/svg"
          >
            <path 
              stroke-linecap="round" 
              stroke-linejoin="round" 
              stroke-width="2" 
              d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z"
            />
          </svg>
        </div>
      </div>

      <!-- Error Content -->
      <div class="bg-white rounded-2xl shadow-lg p-8 sm:p-10 border border-gray-100">
        <!-- Error Code -->
        <div class="mb-6">
          <h1 class="text-6xl sm:text-7xl font-bold text-indigo-600 mb-2">
            500
          </h1>
          <div class="h-1 w-20 bg-gradient-to-r from-indigo-500 to-indigo-700 mx-auto rounded-full"></div>
        </div>

        <!-- Main Message -->
        <h2 class="text-2xl sm:text-3xl font-bold text-gray-900 mb-4" style="font-family: 'Roboto', sans-serif;">
          Server Error
        </h2>
        
        <p class="text-gray-600 text-lg mb-6 leading-relaxed" style="font-family: sans-serif;">
          Oops! Something went wrong on our end. Our technical team has been notified and is working to resolve this issue.
        </p>

        <!-- Additional Info -->
        <div class="bg-indigo-50 rounded-xl p-4 mb-8 border border-indigo-100">
          <p class="text-indigo-800 text-sm font-medium mb-2">
            What happened?
          </p>
          <p class="text-indigo-700 text-sm leading-relaxed">
            The server encountered an unexpected condition that prevented it from fulfilling your request. This is typically a temporary issue.
          </p>
        </div>

        <!-- Action Buttons -->
        <div class="space-y-4 sm:space-y-0 sm:space-x-4 sm:flex sm:justify-center">
          <!-- Retry Button -->
          <button
            @click="retryRequest"
            :disabled="isRetrying"
            class="w-full sm:w-auto inline-flex items-center justify-center px-6 py-3 bg-indigo-600 text-white font-semibold rounded-xl shadow-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200 transform hover:scale-105"
          >
            <svg 
              v-if="isRetrying" 
              class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" 
              xmlns="http://www.w3.org/2000/svg" 
              fill="none" 
              viewBox="0 0 24 24"
            >
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <svg 
              v-else
              class="mr-2 h-5 w-5" 
              fill="none" 
              stroke="currentColor" 
              viewBox="0 0 24 24"
            >
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {{ isRetrying ? 'Retrying...' : 'Try Again' }}
          </button>

          <!-- Go Home Button -->
          <button
            @click="goHome"
            class="w-full sm:w-auto inline-flex items-center justify-center px-6 py-3 bg-white text-indigo-600 font-semibold rounded-xl border-2 border-indigo-600 hover:bg-indigo-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 transition-all duration-200 transform hover:scale-105"
          >
            <svg class="mr-2 h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Go Home
          </button>
        </div>

        <!-- Support Contact -->
        <div class="mt-8 pt-6 border-t border-gray-200">
          <p class="text-gray-500 text-sm mb-3">
            Need immediate assistance?
          </p>
          <div class="flex flex-col sm:flex-row items-center justify-center space-y-2 sm:space-y-0 sm:space-x-6">
            <a 
              href="mailto:support@templesaas.com" 
              class="inline-flex items-center text-indigo-600 hover:text-indigo-800 text-sm font-medium transition-colors duration-200"
            >
              <svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
              </svg>
              Email Support
            </a>
            <a 
              href="tel:+1234567890" 
              class="inline-flex items-center text-indigo-600 hover:text-indigo-800 text-sm font-medium transition-colors duration-200"
            >
              <svg class="mr-2 h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
              </svg>
              Call Support
            </a>
          </div>
        </div>
      </div>

      <!-- Footer Message -->
      <p class="mt-8 text-gray-400 text-sm">
        Error ID: {{ errorId }} • {{ new Date().toLocaleString() }}
      </p>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const isRetrying = ref(false)
const errorId = ref('')

// Generate random error ID for tracking
onMounted(() => {
  errorId.value = Math.random().toString(36).substr(2, 9).toUpperCase()
})

const retryRequest = async () => {
  isRetrying.value = true
  
  try {
    // Simulate retry delay
    await new Promise(resolve => setTimeout(resolve, 2000))
    
    // Reload the current page or redirect to last known route
    window.location.reload()
  } catch (error) {
    console.error('Retry failed:', error)
  } finally {
    isRetrying.value = false
  }
}

const goHome = () => {
  router.push('/')
}
</script>

<style scoped>
/* Custom animations for enhanced UX */
.animate-bounce-slow {
  animation: bounce 2s infinite;
}

@keyframes bounce {
  0%, 100% {
    transform: translateY(-5%);
    animation-timing-function: cubic-bezier(0.8, 0, 1, 1);
  }
  50% {
    transform: translateY(0);
    animation-timing-function: cubic-bezier(0, 0, 0.2, 1);
  }
}

/* Gradient text effect */
.gradient-text {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
</style>