// src/stores/dashboard.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import dashboardService from '@/services/dashboard.service'
import { useToast } from '@/composables/useToast'

export const useDashboardStore = defineStore('dashboard', () => {
  // State
  const dashboardData = ref({
    devotees: { total: 0, newThisMonth: 0 },
    sevas: { today: 0, thisMonth: 0 },
    donations: { thisMonth: 0, growth: 0 },
    events: { upcoming: 0, thisWeek: 0 }
  })
  const loading = ref(false)
  const error = ref(null)
  const toast = useToast()

  // Actions
  const fetchDashboardData = async (entityId) => {
    try {
      loading.value = true
      error.value = null
      
      console.log(`📊 Fetching dashboard data for entity ID: ${entityId}`)
      
      // Add timestamp for cache busting
      const timestamp = Date.now()
      const data = await dashboardService.getDashboardSummary(entityId, timestamp)
      
      dashboardData.value = data
      console.log('📊 Dashboard data set in store:', dashboardData.value)
      
      return data
    } catch (err) {
      const errorMessage = err.message || 'Failed to fetch dashboard data'
      error.value = errorMessage
      toast.error('Failed to load dashboard statistics')
      console.error('Error fetching dashboard data:', err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Reset dashboard data
  const resetDashboardData = () => {
    dashboardData.value = {
      devotees: { total: 0, newThisMonth: 0 },
      sevas: { today: 0, thisMonth: 0 },
      donations: { thisMonth: 0, growth: 0 },
      events: { upcoming: 0, thisWeek: 0 }
    }
  }

  return {
    // State
    dashboardData,
    loading,
    error,
    
    // Actions
    fetchDashboardData,
    resetDashboardData
  }
})