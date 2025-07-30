import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import superAdminService from '@/services/superadmin.service'

export const useSuperAdminStore = defineStore('superadmin', () => {
  // State
  const tenants = ref([])
  const pendingEntities = ref([])
  const stats = ref({
    pendingApprovals: 0,
    activeTenants: 0,
    totalUsers: 0,
    rejectedTenants: 0
  })
  const loadingTenants = ref(false)
  const loadingEntities = ref(false)
  const loadingStats = ref(false)
  const tenantError = ref(null)
  const entityError = ref(null)
  const statsError = ref(null)
  
  // Getters
  const pendingTenants = computed(() => 
    tenants.value.filter(t => t.status === 'PENDING' || t.status === 'pending')
  )
  
  const approvedTenants = computed(() => 
    tenants.value.filter(t => t.status === 'APPROVED' || t.status === 'approved')
  )
  
  const rejectedTenants = computed(() => 
    tenants.value.filter(t => t.status === 'REJECTED' || t.status === 'rejected')
  )

  // Dashboard stats getters
  const pendingCount = computed(() => stats.value.pendingApprovals || pendingTenants.value.length)
  const activeCount = computed(() => stats.value.activeTenants || approvedTenants.value.length)
  const rejectedCount = computed(() => stats.value.rejectedTenants || rejectedTenants.value.length)
  
  // Actions
  async function fetchStats() {
    loadingStats.value = true
    statsError.value = null
    
    try {
      const response = await superAdminService.getDashboardStats()
      if (response.success && response.data) {
        stats.value = {
          pendingApprovals: response.data.pendingApprovals || 0,
          activeTenants: response.data.activeTenants || 0,
          totalUsers: response.data.totalUsers || 0,
          rejectedTenants: response.data.rejectedTenants || 0
        }
      }
    } catch (error) {
      console.error('Error fetching dashboard stats:', error)
      statsError.value = error.message
    } finally {
      loadingStats.value = false
    }
  }

  async function fetchTenants() {
    loadingTenants.value = true
    tenantError.value = null
    
    try {
      console.log('Store: Fetching tenants...')
      const response = await superAdminService.getPendingTenants()
      console.log('Store: Got response:', response)
      
      if (response.success && Array.isArray(response.data)) {
        console.log('Store: Setting', response.data.length, 'tenants')
        tenants.value = response.data.map(tenant => ({
          id: tenant.id,
          fullName: tenant.full_name || tenant.fullName || tenant.name || '',
          name: tenant.name || tenant.fullName || tenant.full_name || '',
          email: tenant.email || '',
          phone: tenant.phone || '',
          status: (tenant.status || 'PENDING').toUpperCase(),
          createdAt: tenant.created_at || tenant.createdAt || new Date().toISOString(),
          updatedAt: tenant.updated_at || tenant.updatedAt,
          rejectionNotes: tenant.rejection_notes || tenant.rejectionNotes,
          temple: tenant.temple ? {
            name: tenant.temple.name,
            type: tenant.temple.type || 'Hindu Temple',
            address: tenant.temple.address,
            city: tenant.temple.city,
            state: tenant.temple.state
          } : {
            name: tenant.name || tenant.fullName || tenant.full_name || 'Unknown Temple',
            type: 'Hindu Temple',
            address: tenant.address || 'No address provided',
            city: tenant.city || 'Unknown',
            state: tenant.state || 'Unknown'
          },
          documents: tenant.documents || []
        }))
      } else {
        console.log('Store: Empty or invalid response')
        tenants.value = []
        tenantError.value = 'No tenant data available'
      }
    } catch (error) {
      console.error('Store: Error fetching tenants:', error)
      tenantError.value = error.message
    } finally {
      loadingTenants.value = false
    }
  }

  async function fetchPendingEntities() {
    loadingEntities.value = true
    entityError.value = null
    
    try {
      const response = await superAdminService.getPendingEntities()
      
      if (response.success && Array.isArray(response.data)) {
        pendingEntities.value = response.data
      } else {
        pendingEntities.value = []
        entityError.value = 'No entity data available'
      }
    } catch (error) {
      console.error('Error fetching pending entities:', error)
      entityError.value = error.message
    } finally {
      loadingEntities.value = false
    }
  }
  
  async function approveTenant(id, data = {}) {
    try {
      const response = await superAdminService.approveTenant(id, data)
      if (response.success) {
        // Update local state
        const index = tenants.value.findIndex(t => t.id === id)
        if (index !== -1) {
          tenants.value[index].status = 'APPROVED'
          tenants.value[index].updatedAt = new Date().toISOString()
        }
        return { success: true }
      }
      return { success: false, error: response.message }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }
  
  async function rejectTenant(id, data = {}) {
    try {
      const response = await superAdminService.rejectTenant(id, data)
      if (response.success) {
        // Update local state
        const index = tenants.value.findIndex(t => t.id === id)
        if (index !== -1) {
          tenants.value[index].status = 'REJECTED'
          tenants.value[index].rejectionNotes = data.notes
          tenants.value[index].updatedAt = new Date().toISOString()
        }
        return { success: true }
      }
      return { success: false, error: response.message }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  // Initialize data on store creation
  function initialize() {
    fetchStats()
    fetchTenants()
    fetchPendingEntities()
  }

  return {
    // State
    tenants,
    pendingEntities,
    stats,
    loadingTenants,
    loadingEntities,
    loadingStats,
    tenantError,
    entityError,
    statsError,
    
    // Getters
    pendingTenants,
    approvedTenants,
    rejectedTenants,
    pendingCount,
    activeCount,
    rejectedCount,
    
    // Actions
    fetchTenants,
    fetchPendingEntities,
    fetchStats,
    approveTenant,
    rejectTenant,
    initialize
  }
})