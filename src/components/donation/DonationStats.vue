<template>
  <div class="grid grid-cols-2 md:grid-cols-3 gap-4">
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">Total Amount</h3>
      <p class="text-2xl font-bold text-green-600">₹ {{ dashboardData.totalAmount }}</p>
    </div>
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">Average Amount</h3>
      <p class="text-2xl font-bold text-blue-600">₹ {{ dashboardData.averageAmount }}</p>
    </div>
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">This Month</h3>
      <p class="text-2xl font-bold text-purple-600">₹ {{ dashboardData.thisMonth }}</p>
    </div>
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">Total Donors</h3>
      <p class="text-2xl font-bold text-orange-600">{{ dashboardData.totalDonors }}</p>
    </div>
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">Completed</h3>
      <p class="text-2xl font-bold text-teal-600">{{ dashboardData.completed }}</p>
    </div>
    <div class="bg-white p-4 shadow rounded-xl text-center">
      <h3 class="text-xl font-semibold text-gray-700">Pending</h3>
      <p class="text-2xl font-bold text-red-600">{{ dashboardData.pending }}</p>
    </div>
  </div>
</template>

<script setup>
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { useDonationStore } from '@/stores/donation'

const route = useRoute()
const donationStore = useDonationStore()

// Dashboard data reference
const dashboardData = ref({
  totalAmount: 0,
  averageAmount: 0,
  thisMonth: 0,
  totalDonors: 0,
  completed: 0,
  pending: 0
})

onMounted(async () => {
  // Get entity ID from route and set to localStorage
  const entityId = route.params.entity_id || route.params.id
  if (entityId) {
    localStorage.setItem('current_entity_id', entityId)
  }

  try {
    await donationStore.fetchDashboard()
    dashboardData.value = donationStore.dashboardData
  } catch (err) {
    console.error('Error fetching dashboard:', err)
  }
})
</script>
