<template>
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    <div
      v-for="stat in stats"
      :key="stat.id"
      class="bg-white rounded-xl shadow-md p-6 border border-gray-100 hover:shadow-lg transition-all duration-200"
    >
      <div class="flex items-center justify-between">
        <div>
          <p class="text-sm font-medium text-gray-600 mb-1">{{ stat.label }}</p>
          <p class="text-3xl font-bold text-gray-900 mb-2">{{ formatNumber(stat.value) }}</p>
          <div class="flex items-center text-sm">
            <span
              :class="[
                'flex items-center font-medium',
                stat.trend > 0 ? 'text-green-600' : stat.trend < 0 ? 'text-red-600' : 'text-gray-500'
              ]"
            >
              <svg
                v-if="stat.trend > 0"
                class="w-4 h-4 mr-1"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 10l7-7m0 0l7 7m-7-7v18" />
              </svg>
              <svg
                v-else-if="stat.trend < 0"
                class="w-4 h-4 mr-1"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24"
              >
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 14l-7 7m0 0l-7-7m7 7V3" />
              </svg>
              {{ Math.abs(stat.trend) }}%
            </span>
            <span class="text-gray-500 ml-1">{{ stat.period }}</span>
          </div>
        </div>
        <div 
          :class="[
            'flex items-center justify-center w-12 h-12 rounded-lg',
            stat.iconBg
          ]"
        >
          <component :is="stat.icon" :class="['w-6 h-6', stat.iconColor]" />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  stats: {
    type: Array,
    required: true,
    validator: (stats) => {
      return stats.every(stat => 
        stat.hasOwnProperty('id') &&
        stat.hasOwnProperty('label') &&
        stat.hasOwnProperty('value') &&
        stat.hasOwnProperty('trend') &&
        stat.hasOwnProperty('period') &&
        stat.hasOwnProperty('icon') &&
        stat.hasOwnProperty('iconBg') &&
        stat.hasOwnProperty('iconColor')
      )
    }
  }
})

const formatNumber = (value) => {
  if (value >= 1000000) {
    return (value / 1000000).toFixed(1) + 'M'
  }
  if (value >= 1000) {
    return (value / 1000).toFixed(1) + 'K'
  }
  return value.toLocaleString()
}
</script>