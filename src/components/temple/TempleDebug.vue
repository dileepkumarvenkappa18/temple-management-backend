<template>
  <div class="fixed top-0 left-0 w-full bg-black bg-opacity-90 text-white p-4 overflow-auto" style="z-index: 9999; max-height: 50vh;">
    <h2 class="text-xl font-bold mb-2">Temple Debug Panel</h2>
    <button @click="$emit('close')" class="absolute top-2 right-2 text-white bg-red-600 px-2 py-1 rounded">Close</button>
    
    <div class="mb-4">
      <h3 class="font-bold">Temple Data</h3>
      <pre class="bg-gray-800 p-2 rounded text-xs overflow-auto max-h-40">{{ JSON.stringify(temples, null, 2) }}</pre>
    </div>
    
    <div class="mb-4">
      <h3 class="font-bold">Temple Statuses</h3>
      <div v-for="(temple, index) in temples" :key="index" class="p-2 border-b border-gray-700">
        <div><b>Name:</b> {{ temple.name }}</div>
        <div><b>ID:</b> {{ temple.id }}</div>
        <div><b>Status:</b> "{{ temple.status }}" ({{ typeof temple.status }})</div>
        <div class="mt-1">
          <button 
            @click="forceApproved(temple)" 
            class="px-2 py-1 bg-green-600 text-xs rounded mr-2"
          >
            Force Approved
          </button>
          <button 
            @click="manageDirect(temple)" 
            class="px-2 py-1 bg-blue-600 text-xs rounded"
          >
            Direct Manage
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  temples: {
    type: Array,
    default: () => []
  }
})

const emit = defineEmits(['close', 'update'])
const router = useRouter()

const forceApproved = (temple) => {
  if (temple) {
    temple.status = 'APPROVED'
    console.log(`Forcing temple ${temple.id} status to APPROVED`)
    emit('update')
  }
}

const manageDirect = (temple) => {
  if (temple && temple.id) {
    console.log(`Direct management of temple ${temple.id}`)
    router.push(`/entity/${temple.id}/dashboard`)
  }
}
</script>