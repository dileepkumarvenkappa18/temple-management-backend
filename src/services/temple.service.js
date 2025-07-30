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
      
      // SIMPLE LOGIC: If on tenant or admin path, use entities endpoint
      if (currentPath.includes('/tenant/') || 
          currentPath.includes('/entity/') || 
          currentPath.includes('/admin/') || 
          currentPath.includes('/superadmin/')) {
        console.log('🔒 Using admin endpoint: /v1/entities')
        response = await api.get('/v1/entities')
      } else {
        // Otherwise use devotee endpoint with search params
        console.log('🔍 Using devotee endpoint: /v1/temples/search')
        const queryString = new URLSearchParams()
        if (searchParams.query) queryString.append('query', searchParams.query)
        if (searchParams.state) queryString.append('state', searchParams.state)
        if (searchParams.type) queryString.append('type', searchParams.type)
        
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

  async createTemple(templeData) {
    try {
      console.log('📡 Creating new temple entity')
      console.log('📋 Raw form data received:', templeData)

      const mainDeity = templeData.MainDeity || templeData.mainDeity || templeData['main_deity'] || ''
      const establishedYear = parseInt(templeData.EstablishedYear || templeData.establishedYear || templeData['established_year'] || '')

      const payload = {
        name: templeData.Name || templeData.name || "",
        main_deity: mainDeity,
        temple_type: templeData.TempleType || templeData.templeType || templeData['temple_type'] || "",
        established_year: establishedYear,
        phone: templeData.Phone || templeData.phone || "",
        email: templeData.Email || templeData.email || "",
        description: templeData.Description || templeData.description || '',
        street_address: templeData.StreetAddress || templeData.addressLine1 || '',
        city: templeData.City || templeData.city || "",
        district: templeData.District || templeData.district || "",
        state: templeData.State || templeData.state || "",
        pincode: templeData.Pincode || templeData.pincode || "",
        landmark: templeData.Landmark || templeData.landmark || '',
        map_link: templeData.MapLink || templeData.googleMapsLink || '',
        accepted_terms: true,
        status: 'pending'
      }

      console.log('🔍 FIELD VALIDATION CHECK:')
      console.log(`name: "${payload.name}" - ${payload.name ? '✅' : '❌'}`)
      console.log(`main_deity: "${payload.main_deity}" - ${payload.main_deity ? '✅' : '❌'}`)
      console.log(`temple_type: "${payload.temple_type}" - ${payload.temple_type ? '✅' : '❌'}`)
      console.log(`established_year: ${payload.established_year} - ${payload.established_year ? '✅' : '❌'}`)
      console.log(`phone: "${payload.phone}" - ${payload.phone ? '✅' : '❌'}`)
      console.log(`email: "${payload.email}" - ${payload.email ? '✅' : '❌'}`)
      console.log(`city: "${payload.city}" - ${payload.city ? '✅' : '❌'}`)
      console.log(`district: "${payload.district}" - ${payload.district ? '✅' : '❌'}`)
      console.log(`state: "${payload.state}" - ${payload.state ? '✅' : '❌'}`)
      console.log(`pincode: "${payload.pincode}" - ${payload.pincode ? '✅' : '❌'}`)

      const missingFields = []
      if (!payload.name) missingFields.push('Temple Name')
      if (!payload.main_deity) missingFields.push('Main Deity')
      if (!payload.temple_type) missingFields.push('Temple Type')
      if (!payload.established_year) missingFields.push('Established Year')
      if (!payload.phone) missingFields.push('Phone')
      if (!payload.email) missingFields.push('Email')
      if (!payload.city) missingFields.push('City')
      if (!payload.district) missingFields.push('District')
      if (!payload.state) missingFields.push('State')
      if (!payload.pincode) missingFields.push('Pincode')

      if (missingFields.length > 0) {
        const errorMessage = `Missing required fields: ${missingFields.join(', ')}. Please fill in all required fields.`
        console.error('❌ VALIDATION ERROR:', errorMessage)
        throw new Error(errorMessage)
      }

      const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
      if (!emailRegex.test(payload.email)) {
        throw new Error('Please enter a valid email address.')
      }

      const phoneRegex = /^\d{10}$/
      if (!phoneRegex.test(payload.phone.replace(/[^0-9]/g, ''))) {
        throw new Error('Please enter a valid 10-digit phone number.')
      }

      console.log('✅ All required fields present!')
      console.log('📦 Payload being sent to API:', payload)

      const response = await api.post('/v1/entities', payload)
      console.log('📥 Create temple response:', response)

      return response.data || response
    } catch (error) {
      console.error('❌ Error creating temple:', error)
      console.error('Error details:', error.message || error.response?.data)
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
      console.log(`📡 Updating temple with ID: ${id}`)

      const mainDeity = updates.MainDeity || updates.mainDeity || updates['main_deity'] || ''
      const establishedYear = parseInt(updates.EstablishedYear || updates.establishedYear || updates['established_year'] || '')

      const payload = {
        name: updates.Name || updates.name || "",
        main_deity: mainDeity,
        temple_type: updates.TempleType || updates.templeType || updates['temple_type'] || "",
        established_year: establishedYear,
        phone: updates.Phone || updates.phone || "",
        email: updates.Email || updates.email || "",
        description: updates.Description || updates.description || '',
        street_address: updates.StreetAddress || updates.addressLine1 || '',
        city: updates.City || updates.city || "",
        district: updates.District || updates.district || "",
        state: updates.State || updates.state || "",
        pincode: updates.Pincode || updates.pincode || "",
        landmark: updates.Landmark || updates.landmark || '',
        map_link: updates.MapLink || updates.googleMapsLink || ''
      }

      const missingFields = []
      if (!payload.name) missingFields.push('Temple Name')
      if (!payload.main_deity) missingFields.push('Main Deity')
      if (!payload.temple_type) missingFields.push('Temple Type')
      if (!payload.established_year) missingFields.push('Established Year')
      if (!payload.phone) missingFields.push('Phone')
      if (!payload.email) missingFields.push('Email')
      if (!payload.city) missingFields.push('City')
      if (!payload.district) missingFields.push('District')
      if (!payload.state) missingFields.push('State')
      if (!payload.pincode) missingFields.push('Pincode')

      if (missingFields.length > 0) {
        const errorMessage = `Missing required fields: ${missingFields.join(', ')}. Please fill in all required fields.`
        console.error('❌ VALIDATION ERROR:', errorMessage)
        throw new Error(errorMessage)
      }

      const response = await api.put(`/v1/entities/${id}`, payload)
      console.log('📥 Update temple response:', response)

      return response.data || response
    } catch (error) {
      console.error(`❌ Error updating temple ID ${id}:`, error)
      console.error('Error details:', error.message || error.response?.data)
      throw error
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