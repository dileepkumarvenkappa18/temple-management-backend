import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import superAdminService from '@/services/superadmin.service'

export const useSuperAdminStore = defineStore('superadmin', () => {
  // State
  const tenants = ref([])
  const pendingEntities = ref([])
  
  // NEW: User Management State
  const users = ref([])
  const userRoles = ref([])
  
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
  const loadingTenantCounts = ref(false)
  const loadingTempleCounts = ref(false)
  
  // NEW: User Management Loading States
  const loadingUsers = ref(false)
  const loadingUserRoles = ref(false)
  const loadingUserAction = ref(false) // For create/update/delete operations
  
  // Error states
  const tenantError = ref(null)
  const entityError = ref(null)
  const statsError = ref(null)
  const tenantCountsError = ref(null)
  const templeCountsError = ref(null)
  
  // NEW: User Management Error States
  const userError = ref(null)
  const userRolesError = ref(null)
  const userActionError = ref(null)
  
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
  
  // NEW: User Management Getters
  const activeUsers = computed(() => 
    users.value.filter(u => u.status === 'active')
  )
  
  const inactiveUsers = computed(() => 
    users.value.filter(u => u.status === 'inactive')
  )
  
  const usersByRole = computed(() => {
    const roleGroups = {}
    users.value.forEach(user => {
      const roleName = user.role?.role_name || 'unknown'
      if (!roleGroups[roleName]) {
        roleGroups[roleName] = []
      }
      roleGroups[roleName].push(user)
    })
    return roleGroups
  })
  
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

  // ==========================================
  // NEW: USER MANAGEMENT ACTIONS
  // ==========================================

  // Fetch user roles
  async function fetchUserRoles() {
    loadingUserRoles.value = true
    userRolesError.value = null
    
    try {
      console.log('Store: Fetching user roles...')
      const response = await superAdminService.getUserRoles()
      
      if (response.success && Array.isArray(response.data)) {
        userRoles.value = response.data
        console.log('Store: User roles loaded:', userRoles.value)
      } else {
        userRoles.value = []
        userRolesError.value = response.message || 'No user roles available'
      }
    } catch (error) {
      console.error('Store: Error fetching user roles:', error)
      userRolesError.value = error.message
      userRoles.value = []
    } finally {
      loadingUserRoles.value = false
    }
  }

  // Fetch users with optional filters
  async function fetchUsers(filters = {}) {
    loadingUsers.value = true
    userError.value = null
    
    try {
      console.log('Store: Fetching users with filters:', filters)
      const response = await superAdminService.getUsers(filters)
      
      if (response.success && Array.isArray(response.data)) {
        users.value = response.data
        console.log('Store: Users loaded:', users.value.length, 'users')
      } else {
        users.value = []
        userError.value = response.message || 'No users available'
      }
    } catch (error) {
      console.error('Store: Error fetching users:', error)
      userError.value = error.message
      users.value = []
    } finally {
      loadingUsers.value = false
    }
  }

  // Get user by ID
  async function fetchUserById(userId) {
    loadingUserAction.value = true
    userActionError.value = null
    
    try {
      console.log('Store: Fetching user by ID:', userId)
      const response = await superAdminService.getUserById(userId)
      
      if (response.success && response.data) {
        return {
          success: true,
          data: response.data
        }
      } else {
        return {
          success: false,
          error: response.message || 'User not found'
        }
      }
    } catch (error) {
      console.error('Store: Error fetching user by ID:', error)
      userActionError.value = error.message
      return {
        success: false,
        error: error.message
      }
    } finally {
      loadingUserAction.value = false
    }
  }

  // Create new user
  async function createUser(userData) {
    loadingUserAction.value = true
    userActionError.value = null
    
    try {
      console.log('Store: Creating user:', userData)
      const response = await superAdminService.createUser(userData)
      
      if (response.success) {
        // Refresh user list after successful creation
        await fetchUsers()
        
        return {
          success: true,
          message: 'User created successfully'
        }
      } else {
        userActionError.value = response.message
        return {
          success: false,
          error: response.message
        }
      }
    } catch (error) {
      console.error('Store: Error creating user:', error)
      userActionError.value = error.message
      return {
        success: false,
        error: error.message
      }
    } finally {
      loadingUserAction.value = false
    }
  }

  // Update user
  async function updateUser(userId, userData) {
    loadingUserAction.value = true
    userActionError.value = null
    
    try {
      console.log('Store: Updating user:', userId, userData)
      const response = await superAdminService.updateUser(userId, userData)
      
      if (response.success) {
        // Update local state
        const userIndex = users.value.findIndex(u => u.id === userId)
        if (userIndex !== -1) {
          // Refetch user data to ensure we have latest data
          const userResponse = await superAdminService.getUserById(userId)
          if (userResponse.success) {
            users.value[userIndex] = userResponse.data
          }
        }
        
        return {
          success: true,
          message: 'User updated successfully'
        }
      } else {
        userActionError.value = response.message
        return {
          success: false,
          error: response.message
        }
      }
    } catch (error) {
      console.error('Store: Error updating user:', error)
      userActionError.value = error.message
      return {
        success: false,
        error: error.message
      }
    } finally {
      loadingUserAction.value = false
    }
  }

  // Update user status
  async function updateUserStatus(userId, status) {
    loadingUserAction.value = true
    userActionError.value = null
    
    try {
      console.log('Store: Updating user status:', userId, status)
      const response = await superAdminService.updateUserStatus(userId, status)
      
      if (response.success) {
        // Update local state
        const userIndex = users.value.findIndex(u => u.id === userId)
        if (userIndex !== -1) {
          users.value[userIndex].status = status
        }
        
        return {
          success: true,
          message: `User status updated to ${status}`
        }
      } else {
        userActionError.value = response.message
        return {
          success: false,
          error: response.message
        }
      }
    } catch (error) {
      console.error('Store: Error updating user status:', error)
      userActionError.value = error.message
      return {
        success: false,
        error: error.message
      }
    } finally {
      loadingUserAction.value = false
    }
  }

  // UPDATED: Initialize with new count fetching
  function initialize() {
    // Fetch counts first (most important for dashboard)
    fetchStats()
    
    // Then fetch detailed data
    fetchTenants()
    fetchPendingEntities()
    
    // NEW: Initialize user management data
    fetchUserRoles()
    fetchUsers()
  }

  // NEW: Refresh all counts
  async function refreshCounts() {
    await Promise.all([
      fetchTenantCounts(),
      fetchTempleCounts()
    ])
  }

  // NEW: Refresh all user data
  async function refreshUserData() {
    await Promise.all([
      fetchUserRoles(),
      fetchUsers()
    ])
  }

  return {
    // State
    tenants,
    pendingEntities,
    stats, // Legacy compatibility
    tenantCounts,
    templeCounts,
    
    // NEW: User Management State
    users,
    userRoles,
    
    // Loading states
    loadingTenants,
    loadingEntities,
    loadingStats,
    loadingTenantCounts,
    loadingTempleCounts,
    
    // NEW: User Management Loading States
    loadingUsers,
    loadingUserRoles,
    loadingUserAction,
    
    // Error states
    tenantError,
    entityError,
    statsError,
    tenantCountsError,
    templeCountsError,
    
    // NEW: User Management Error States
    userError,
    userRolesError,
    userActionError,
    
    // Getters
    pendingTenants,
    approvedTenants,
    rejectedTenants,
    pendingCount,
    activeCount,
    rejectedCount,
    templePendingCount,
    templeActiveCount,
    templeRejectedCount,
    totalUsersCount,
    
    // NEW: User Management Getters
    activeUsers,
    inactiveUsers,
    usersByRole,
    
    // Actions
    fetchTenants,
    fetchPendingEntities,
    fetchStats,
    fetchTenantCounts,
    fetchTempleCounts,
    refreshCounts,
    approveTenant,
    rejectTenant,
    approveEntity,
    rejectEntity,
    
    // NEW: User Management Actions
    fetchUserRoles,
    fetchUsers,
    fetchUserById,
    createUser,
    updateUser,
    updateUserStatus,
    refreshUserData,
    
    initialize // UPDATED to include user management
  }
})

