// src/services/temple.service.js

import api from '@/plugins/axios'

const templeService = {
  async getTemples(searchParams = {}) {
    try {
      console.log('📡 Making API call to fetch available temples')
      console.log('🔍 Search params:', searchParams)

      // DIRECT FIX: Check the current URL path directly
      const currentPath = window.location.pathname
      console.log('📍 Current path:', currentPath)
      
      let response
      
      // Handle SuperAdmin temple fetching with fallback strategy
      if (searchParams.superAdmin) {
        console.log(`🔍 Using SuperAdmin endpoint for tenant ${searchParams.tenantId}`)
        try {
          // Add cache busting timestamp to ensure fresh data
          const timestamp = Date.now()
          
          // First attempt with specific tenant parameter and cache busting
          response = await api.get(`/v1/entities?tenant_id=${searchParams.tenantId}&_=${timestamp}`)
          console.log('✅ First attempt successful with /v1/entities endpoint')
        } catch (err) {
          console.log('⚠️ First endpoint failed, trying fallback...', err.message)
          try {
            const timestamp = Date.now()
            // Fallback to superadmin specific endpoint with cache busting
            response = await api.get(`/v1/superadmin/entities?tenant_id=${searchParams.tenantId}&_=${timestamp}`)
            console.log('✅ Second attempt successful with /v1/superadmin/entities endpoint')
          } catch (err2) {
            console.log('⚠️ Second endpoint failed, trying third endpoint...', err2.message)
            try {
              const timestamp = Date.now()
              // Try a different superadmin endpoint format with cache busting
              response = await api.get(`/v1/superadmin/tenants/${searchParams.tenantId}/entities?_=${timestamp}`)
              console.log('✅ Third attempt successful with /v1/superadmin/tenants/[id]/entities endpoint')
            } catch (err3) {
              console.log('⚠️ All SuperAdmin tenant-specific endpoints failed, using generic endpoint', err3.message)
              // Last resort, general entities endpoint
              const timestamp = Date.now()
              response = await api.get(`/v1/entities?_=${timestamp}`)
            }
          }
        }
      }
      // NEW LOGIC: Special case for temple admin dashboard
      else if (currentPath.includes('/tenant/dashboard')) {
        // Use the special endpoint for temple admins to see their created temples
        const timestamp = Date.now()
        console.log('🔑 Using temple admin special endpoint: /v1/entities/by-creator')
        response = await api.get(`/v1/entities/by-creator?_=${timestamp}`)
      }
      // For other admin paths
      else if (currentPath.includes('/tenant/') || 
          currentPath.includes('/entity/') || 
          currentPath.includes('/admin/') || 
          currentPath.includes('/superadmin/')) {
        const timestamp = Date.now()
        console.log('🔒 Using admin endpoint: /v1/entities')
        response = await api.get(`/v1/entities?_=${timestamp}`)
      } else {
        // Otherwise use devotee endpoint with search params
        console.log('🔍 Using devotee endpoint: /v1/temples/search')
        const queryString = new URLSearchParams()
        if (searchParams.query) queryString.append('query', searchParams.query)
        if (searchParams.state) queryString.append('state', searchParams.state)
        if (searchParams.type) queryString.append('type', searchParams.type)
        // Add cache busting
        queryString.append('_', Date.now())
        
        response = await api.get(`/v1/temples/search${queryString.toString() ? '?' + queryString.toString() : ''}`)
      }
      
      console.log('📥 Temple API response received:', response)

      // Extract data from response
      let templeData = response.data || response
      if (!Array.isArray(templeData)) {
        if (templeData.data && Array.isArray(templeData.data)) {
          templeData = templeData.data
        } else if (templeData.temples && Array.isArray(templeData.temples)) {
          templeData = templeData.temples
        } else if (templeData.entities && Array.isArray(templeData.entities)) {
          templeData = templeData.entities
        } else if (templeData.items && Array.isArray(templeData.items)) {
          templeData = templeData.items
        }
      }

      if (!Array.isArray(templeData)) {
        console.error('🚨 Could not extract array from response:', response)
        return [] // Return empty array instead of throwing error for better UI handling
      }

      // IMPORTANT FIX: Filter temples by tenant ID for SuperAdmin flow
      if (searchParams.superAdmin && searchParams.tenantId) {
        console.log(`🔍 Filtering temples by tenant ID ${searchParams.tenantId}`)
        templeData = templeData.filter(temple => 
          // Check different possible property names for created_by or tenant_id
          (temple.created_by && temple.created_by.toString() === searchParams.tenantId.toString()) ||
          (temple.tenant_id && temple.tenant_id.toString() === searchParams.tenantId.toString()) ||
          (temple.creator_id && temple.creator_id.toString() === searchParams.tenantId.toString())
        )
        console.log(`✅ After filtering: ${templeData.length} temples match tenant ID ${searchParams.tenantId}`)
      }

      const normalizedTemples = templeData.map(temple => this.normalizeTempleData(temple))
      console.log('✅ Normalized temples:', normalizedTemples)

      return normalizedTemples
    } catch (error) {
      console.error('❌ Error fetching temples:', error)
      console.error('Error response:', error.response?.data)
      // Return empty array instead of throwing error for better UI handling
      return []
    }
  },

  /**
   * Get temples for a specific tenant when logged in as SuperAdmin
   * This is a convenience method that calls getTemples with superAdmin flag
   * @param {string|number} tenantId - The ID of the tenant to fetch temples for
   * @returns {Array} - Normalized list of temples
   */
  async getSuperAdminTemples(tenantId) {
    console.log(`📡 Making API call to fetch temples for SuperAdmin for tenant ${tenantId}`)
    return this.getTemples({
      tenantId: tenantId,
      superAdmin: true
    })
  },

  async createTemple(templeData) {
    try {
      console.log('📡 Creating new temple entity')
      console.log('📋 Raw form data received:', templeData)

      // CRITICAL: Pass the data exactly as it is
      // Don't transform or create a new object
      console.log('🚨 DIRECT API CALL WITH ORIGINAL DATA')
      console.log('Street address check:', templeData.street_address)
      
      // Make sure the field exists
      if (!templeData.street_address) {
        console.warn('⚠️ street_address missing in form data!')
      }

      // Just pass the data directly to the API
      const response = await api.post('/v1/entities', templeData)
      console.log('✅ Direct API response:', response)
      
      return response.data || response
    } catch (error) {
      console.error('❌ Error creating temple:', error)
      throw error
    }
  },

  async getTempleById(id) {
    try {
      console.log(`📡 Fetching temple with ID: ${id}`)

      const response = await api.get(`/v1/entities/${id}`)
      console.log('📥 Temple by ID response:', response)

      return this.normalizeTempleData(response.data || response)
    } catch (error) {
      console.error(`❌ Error fetching temple ID ${id}:`, error)
      console.error('Error response:', error.response?.data)
      throw error
    }
  },

  async updateTemple(id, updates) {
    try {
      console.log(`📡 Updating temple with ID: ${id}`);
      console.log('📦 Update data:', updates);

      // For updates, we also use the direct field names
      const payload = {
        id: parseInt(id),
        name: updates.name || '',
        main_deity: updates.main_deity || '',
        temple_type: updates.temple_type || '',
        established_year: parseInt(updates.established_year || 0),
        phone: updates.phone || '',
        email: updates.email || '',
        description: updates.description || '',
        street_address: updates.street_address || '',
        city: updates.city || '',
        district: updates.district || '',
        state: updates.state || '',
        pincode: updates.pincode || '',
        landmark: updates.landmark || '',
        map_link: updates.map_link || ''
      };
      
      console.log('📦 Final update payload:', payload);
      
      const response = await api.put(`/v1/entities/${id}`, payload);
      console.log('📥 Update temple response:', response);
      return response.data || response;
    } catch (error) {
      console.error(`❌ Error updating temple ID ${id}:`, error);
      console.error('Error details:', error.message || error.response?.data);
      throw error;
    }
  },

  async deleteTemple(id) {
    try {
      console.log(`📡 Deleting temple with ID: ${id}`)

      const response = await api.delete(`/v1/entities/${id}`)
      console.log('✅ Temple deleted successfully')
      return response.data || response
    } catch (error) {
      console.error(`❌ Error deleting temple ID ${id}:`, error)
      console.error('Error response:', error.response?.data)
      throw error
    }
  },

  async joinTemple(templeId) {
    try {
      console.log(`📡 Joining temple ID ${templeId}`)

      // Use the correct endpoint for joining temples
      const response = await api.post('/v1/memberships', {
        entity_id: templeId
      })

      console.log('📥 Join temple response:', response)
      return response.data || response
    } catch (error) {
      console.error(`❌ Error joining temple ID ${templeId}:`, error)
      console.error('Error response:', error.response?.data)
      throw error
    }
  },

  async getRecentTemples() {
    try {
      console.log('📡 Fetching recent temples')
      
      // Use the correct endpoint for recent temples
      const response = await api.get('/v1/temples/recent')
      console.log('📥 Recent temples response:', response)
      
      // Extract data from response
      let templeData = response.data || response;
      if (!Array.isArray(templeData)) {
        if (templeData.data && Array.isArray(templeData.data)) {
          templeData = templeData.data
        } else if (templeData.temples && Array.isArray(templeData.temples)) {
          templeData = templeData.temples
        } else if (templeData.entities && Array.isArray(templeData.entities)) {
          templeData = templeData.entities
        } else if (templeData.items && Array.isArray(templeData.items)) {
          templeData = templeData.items
        }
      }
      
      if (!Array.isArray(templeData)) {
        console.error('🚨 Could not extract array from response:', response)
        return []
      }
      
      return templeData.map(temple => this.normalizeTempleData(temple))
    } catch (error) {
      console.error('❌ Error fetching recent temples:', error)
      console.error('Error response:', error.response?.data)
      return []
    }
  },

  normalizeTempleData(temple) {
    if (!temple) return null

    return {
      id: temple.id || temple.ID || 0,
      name: temple.name || temple.Name || 'Unknown Temple',
      description: temple.description || temple.Description || '',
      templeType: temple.temple_type || temple.TempleType || '',
      category: temple.temple_type || temple.TempleType || '',

      addressLine1: temple.street_address || temple.StreetAddress || '',
      city: temple.city || temple.City || '',
      state: temple.state || temple.State || '',
      district: temple.district || temple.District || '',
      pincode: temple.pincode || temple.Pincode || '',
      country: 'India',

      phone: temple.phone || temple.Phone || '',
      email: temple.email || temple.Email || '',

      status: temple.status || temple.Status || 'pending',
      devoteeCount: temple.devotee_count || temple.DevoteeCount || 0,
      volunteersCount: temple.volunteers_count || temple.VolunteersCount || 0,

      image: temple.image_url || temple.ImageUrl || null,

      mainDeity: temple.main_deity || temple.MainDeity || '',
      establishedYear: temple.established_year || temple.EstablishedYear || null,

      createdAt: temple.created_at || temple.CreatedAt || null,
      updatedAt: temple.updated_at || temple.UpdatedAt || null,
      
      // Add tenant/creator information for filtering
      createdBy: temple.created_by || temple.CreatedBy || null,
      tenantId: temple.tenant_id || temple.TenantId || temple.created_by || temple.CreatedBy || null,

      address: {
        street: temple.street_address || temple.StreetAddress || '',
        city: temple.city || temple.City || '',
        state: temple.state || temple.State || '',
        district: temple.district || temple.District || '',
        pincode: temple.pincode || temple.Pincode || '',
        country: 'India'
      },

      contact: {
        phone: temple.phone || temple.Phone || '',
        email: temple.email || temple.Email || '',
        website: temple.website || temple.Website || ''
      }
    }
  },

  async searchTemples(query) {
    try {
      console.log(`📡 Searching temples with query: ${query}`)

      // Check the current URL path directly
      const currentPath = window.location.pathname
      
      let response
      
      // If on tenant or admin path, use entities endpoint
      if (currentPath.includes('/tenant/') || 
          currentPath.includes('/entity/') || 
          currentPath.includes('/admin/') || 
          currentPath.includes('/superadmin/')) {
        console.log('🔒 Using admin search endpoint: /v1/entities')
        response = await api.get(`/v1/entities?search=${encodeURIComponent(query)}`)
      } else {
        // Otherwise use devotee endpoint
        console.log('🔍 Using devotee search endpoint: /v1/temples/search')
        response = await api.get(`/v1/temples/search?query=${encodeURIComponent(query)}`)
      }

      // Extract data from response
      let templeData = response.data || response;
      if (!Array.isArray(templeData)) {
        if (templeData.data && Array.isArray(templeData.data)) {
          templeData = templeData.data
        } else if (templeData.temples && Array.isArray(templeData.temples)) {
          templeData = templeData.temples
        } else if (templeData.entities && Array.isArray(templeData.entities)) {
          templeData = templeData.entities
        } else if (templeData.items && Array.isArray(templeData.items)) {
          templeData = templeData.items
        }
      }

      if (!Array.isArray(templeData)) {
        console.error('🚨 Could not extract array from response:', response)
        return []
      }

      return templeData.map(temple => this.normalizeTempleData(temple))
    } catch (error) {
      console.error('❌ Error searching temples:', error)
      console.error('Error response:', error.response?.data)
      return []
    }
  }
}

export default templeService