<!-- src/layouts/PublicLayout.vue -->
<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100">
    

    <!-- Main Content -->
    <main class="flex-1">
      <router-view />
    </main>

    
  </div>
</template>

<script setup>
import { ref } from 'vue'

const mobileMenuOpen = ref(false)
</script>