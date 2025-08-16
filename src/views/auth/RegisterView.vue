<!-- src/views/auth/RegisterView.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-b from-indigo-50 to-white py-12 px-4 sm:px-6 lg:px-8 flex flex-col justify-center">
    <div class="sm:mx-auto sm:w-full sm:max-w-md">
      <!-- Logo -->
      <div class="mx-auto w-20 h-20 bg-indigo-600 rounded-full flex items-center justify-center">
        <svg xmlns="http://www.w3.org/2000/svg" class="h-12 w-12 text-white" viewBox="0 0 20 20" fill="currentColor">
          <path d="M10.707 2.293a1 1 0 00-1.414 0l-7 7a1 1 0 001.414 1.414L4 10.414V17a1 1 0 001 1h2a1 1 0 001-1v-2a1 1 0 011-1h2a1 1 0 011 1v2a1 1 0 001 1h2a1 1 0 001-1v-6.586l.293.293a1 1 0 001.414-1.414l-7-7z" />
        </svg>
      </div>
      
      <h2 class="mt-6 text-center text-3xl font-extrabold text-gray-900">Join Temple SaaS</h2>
      <p class="mt-2 text-center text-sm text-gray-600">
        Create your account to get started
      </p>
    </div>

    <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
      <!-- Register Form Component -->
      <RegisterForm />
      
      <!-- Footer -->
      <div class="mt-6 text-center text-xs text-gray-500">
        © 2025 Temple SaaS. All rights reserved.
      </div>
    </div>
  </div>
</template>

<script setup>
import RegisterForm from '@/components/auth/RegisterForm.vue'
</script>