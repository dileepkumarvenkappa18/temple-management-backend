<template>
  <div class="grid grid-cols-3 gap-4 p-4 bg-white rounded shadow">
    <div class="stat-item">
      <h3 class="text-sm font-semibold text-gray-500">Total Devotees</h3>
      <p class="text-2xl font-bold text-gray-900">{{ stats.total_devotees || 0 }}</p>
    </div>
    <div class="stat-item">
      <h3 class="text-sm font-semibold text-gray-500">Active Devotees</h3>
      <p class="text-2xl font-bold text-green-600">{{ stats.active_devotees || 0 }}</p>
    </div>
    <div class="stat-item">
      <h3 class="text-sm font-semibold text-gray-500">New This Month</h3>
      <p class="text-2xl font-bold text-indigo-600">{{ stats.new_this_month || 0 }}</p>
    </div>
  </div>
</template>

<script>
export default {
  name: 'DevoteeStats',
  props: {
    stats: {
      type: Object,
      default: () => ({})
    },
    loading: {
      type: Boolean,
      default: false
    }
  }
}
</script>

<style scoped>
.stat-item {
  text-align: center;
}
</style>