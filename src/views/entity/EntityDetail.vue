<template>
  <div class="entity-detail">
    <h1 class="text-2xl font-bold mb-4">Entity Details</h1>
    <div class="bg-white shadow rounded p-4">
      <p>Entity ID: {{ $route.params.id }}</p>
      <!-- Add your entity details content here -->
      <p class="text-gray-500 mt-4">Entity details loading...</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()

onMounted(() => {
  console.log('Entity ID:', route.params.id)
  // Add your entity fetch logic here
})
</script>