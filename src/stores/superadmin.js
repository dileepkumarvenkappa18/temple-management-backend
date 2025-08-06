import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import superAdminService from '@/services/superadmin.service'

export const useSuperAdminStore = defineStore('superadmin', () => {
  // State
  const tenants = ref([])
  const pendingEntities = ref([])
  
  // NEW: Separate count states for better tracking
  const tenantCounts = ref({
    pending: 0,
    approved: 0,
    rejected: 0
  })
  
  const templeCounts = ref({
    pendingApprovals: 0,
    activeTemples: 0,
    totalUsers: 0,
    rejectedTemples: 0,
    newThisWeek: 0
  })
  
  // Legacy stats for backward compatibility
  const stats = ref({
    pendingApprovals: 0,
    activeTenants: 0,
    totalUsers: 0,
    rejectedTenants: 0
  })
  
  // Loading states
  const loadingTenants = ref(false)
  const loadingEntities = ref(false)
  const loadingStats = ref(false)
  const loadingTenantCounts = ref(false) // NEW
  const loadingTempleCounts = ref(false) // NEW
  
  // Error states
  const tenantError = ref(null)
  const entityError = ref(null)
  const statsError = ref(null)
  const tenantCountsError = ref(null) // NEW
  const templeCountsError = ref(null) // NEW
  
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

  // Dashboard stats getters - UPDATED to use new count endpoints
  const pendingCount = computed(() => {
    // Prioritize API counts over calculated counts
    return tenantCounts.value.pending || pendingTenants.value.length
  })
  
  const activeCount = computed(() => {
    return tenantCounts.value.approved || approvedTenants.value.length
  })
  
  const rejectedCount = computed(() => {
    return tenantCounts.value.rejected || rejectedTenants.value.length
  })

  // Temple stats getters
  const templePendingCount = computed(() => templeCounts.value.pendingApprovals)
  const templeActiveCount = computed(() => templeCounts.value.activeTemples)
  const templeRejectedCount = computed(() => templeCounts.value.rejectedTemples)
  const totalUsersCount = computed(() => templeCounts.value.totalUsers)
  
  // Actions

  // NEW: Fetch tenant counts from API
  async function fetchTenantCounts() {
    loadingTenantCounts.value = true
    tenantCountsError.value = null
    
    try {
      console.log('Store: Fetching tenant approval counts...')
      const response = await superAdminService.getTenantApprovalCounts()
      
      if (response.success && response.data) {
        tenantCounts.value = {
          pending: response.data.pending || 0,
          approved: response.data.approved || 0,
          rejected: response.data.rejected || 0
        }
        
        console.log('Store: Updated tenant counts:', tenantCounts.value)
        
        // Update legacy stats for backward compatibility
        stats.value.activeTenants = tenantCounts.value.approved
        stats.value.rejectedTenants = tenantCounts.value.rejected
      } else {
        console.warn('Store: Failed to fetch tenant counts:', response.message)
        tenantCountsError.value = response.message || 'Failed to fetch tenant counts'
      }
    } catch (error) {
      console.error('Store: Error fetching tenant counts:', error)
      tenantCountsError.value = error.message
    } finally {
      loadingTenantCounts.value = false
    }
  }

  // NEW: Fetch temple counts from API
  async function fetchTempleCounts() {
    loadingTempleCounts.value = true
    templeCountsError.value = null
    
    try {
      console.log('Store: Fetching temple approval counts...')
      const response = await superAdminService.getTempleApprovalCounts()
      
      if (response.success && response.data) {
        templeCounts.value = {
          pendingApprovals: response.data.pendingApprovals || 0,
          activeTemples: response.data.activeTemples || 0,
          totalUsers: response.data.totalUsers || 0,
          rejectedTemples: response.data.rejectedTemples || 0,
          newThisWeek: response.data.newThisWeek || 0
        }
        
        console.log('Store: Updated temple counts:', templeCounts.value)
        
        // Update legacy stats for backward compatibility
        stats.value.pendingApprovals = templeCounts.value.pendingApprovals
        stats.value.totalUsers = templeCounts.value.totalUsers
      } else {
        console.warn('Store: Failed to fetch temple counts:', response.message)
        templeCountsError.value = response.message || 'Failed to fetch temple counts'
      }
    } catch (error) {
      console.error('Store: Error fetching temple counts:', error)
      templeCountsError.value = error.message
    } finally {
      loadingTempleCounts.value = false
    }
  }

  // UPDATED: Legacy stats method now uses the new count endpoints
  async function fetchStats() {
    loadingStats.value = true
    statsError.value = null
    
    try {
      // Fetch both tenant and temple counts
      await Promise.all([
        fetchTenantCounts(),
        fetchTempleCounts()
      ])
      
      // Legacy stats are automatically updated by the individual fetch functions
      console.log('Store: All stats updated')
    } catch (error) {
      console.error('Store: Error fetching combined stats:', error)
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
        
        // Refresh counts after approval
        await fetchTenantCounts()
        
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
        
        // Refresh counts after rejection
        await fetchTenantCounts()
        
        return { success: true }
      }
      return { success: false, error: response.message }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  // NEW: Approve entity (temple) with count refresh
  async function approveEntity(id, data = {}) {
    try {
      const response = await superAdminService.approveEntity(id, data)
      if (response.success) {
        // Update local state
        const index = pendingEntities.value.findIndex(e => e.id === id)
        if (index !== -1) {
          pendingEntities.value[index].status = 'APPROVED'
          pendingEntities.value[index].updatedAt = new Date().toISOString()
        }
        
        // Refresh temple counts after approval
        await fetchTempleCounts()
        
        return { success: true }
      }
      return { success: false, error: response.message }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  // NEW: Reject entity (temple) with count refresh
  async function rejectEntity(id, data = {}) {
    try {
      const response = await superAdminService.rejectEntity(id, data)
      if (response.success) {
        // Update local state
        const index = pendingEntities.value.findIndex(e => e.id === id)
        if (index !== -1) {
          pendingEntities.value[index].status = 'REJECTED'
          pendingEntities.value[index].rejectionNotes = data.notes
          pendingEntities.value[index].updatedAt = new Date().toISOString()
        }
        
        // Refresh temple counts after rejection
        await fetchTempleCounts()
        
        return { success: true }
      }
      return { success: false, error: response.message }
    } catch (error) {
      return { success: false, error: error.message }
    }
  }

  // UPDATED: Initialize with new count fetching
  function initialize() {
    // Fetch counts first (most important for dashboard)
    fetchStats()
    
    // Then fetch detailed data
    fetchTenants()
    fetchPendingEntities()
  }

  // NEW: Refresh all counts
  async function refreshCounts() {
    await Promise.all([
      fetchTenantCounts(),
      fetchTempleCounts()
    ])
  }

  return {
    // State
    tenants,
    pendingEntities,
    stats, // Legacy compatibility
    tenantCounts, // NEW
    templeCounts, // NEW
    
    // Loading states
    loadingTenants,
    loadingEntities,
    loadingStats,
    loadingTenantCounts, // NEW
    loadingTempleCounts, // NEW
    
    // Error states
    tenantError,
    entityError,
    statsError,
    tenantCountsError, // NEW
    templeCountsError, // NEW
    
    // Getters
    pendingTenants,
    approvedTenants,
    rejectedTenants,
    pendingCount,
    activeCount,
    rejectedCount,
    templePendingCount, // NEW
    templeActiveCount, // NEW
    templeRejectedCount, // NEW
    totalUsersCount, // NEW
    
    // Actions
    fetchTenants,
    fetchPendingEntities,
    fetchStats, // UPDATED to use new endpoints
    fetchTenantCounts, // NEW
    fetchTempleCounts, // NEW
    refreshCounts, // NEW
    approveTenant, // UPDATED to refresh counts
    rejectTenant, // UPDATED to refresh counts
    approveEntity, // NEW
    rejectEntity, // NEW
    initialize // UPDATED
  }
})