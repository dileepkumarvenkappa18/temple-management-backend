import { ref, computed } from 'vue'
import { apiClient } from '@/plugins/axios'
import { useToast } from './useToast'

export function useApi() {
  const { showToast } = useToast()
  
  // Global loading state
  const loading = ref(false)
  const error = ref(null)
  
  // Create a reactive API call wrapper
  const createApiCall = (apiFunction, options = {}) => {
    const data = ref(null)
    const isLoading = ref(false)
    const error = ref(null)
    const isSuccess = ref(false)
    
    const {
      onSuccess = null,
      onError = null,
      showSuccessToast = false,
      showErrorToast = true,
      successMessage = 'Operation completed successfully',
      transform = null // Function to transform response data
    } = options
    
    const execute = async (...args) => {
      try {
        isLoading.value = true
        error.value = null
        isSuccess.value = false
        
        const response = await apiFunction(...args)
        
        // Transform data if transformer provided
        const responseData = transform ? transform(response.data) : response.data
        data.value = responseData
        isSuccess.value = true
        
        // Show success toast if enabled
        if (showSuccessToast) {
          const message = response.data?.message || successMessage
          showToast(message, 'success')
        }
        
        // Call success callback
        if (onSuccess) {
          onSuccess(responseData, response)
        }
        
        return { success: true, data: responseData }
      } catch (err) {
        error.value = err.response?.data?.message || err.message || 'An error occurred'
        isSuccess.value = false
        
        // Show error toast if enabled
        if (showErrorToast) {
          showToast(error.value, 'error')
        }
        
        // Call error callback
        if (onError) {
          onError(err)
        }
        
        return { success: false, error: error.value }
      } finally {
        isLoading.value = false
      }
    }
    
    // Reset function
    const reset = () => {
      data.value = null
      error.value = null
      isLoading.value = false
      isSuccess.value = false
    }
    
    return {
      data: computed(() => data.value),
      isLoading: computed(() => isLoading.value),
      error: computed(() => error.value),
      isSuccess: computed(() => isSuccess.value),
      execute,
      reset
    }
  }
  
  // Pre-configured API calls for common operations
  
  // Authentication APIs
  const useLogin = (options = {}) => {
    return createApiCall(apiClient.auth.login, {
      showSuccessToast: false, // Handled in useAuth
      ...options
    })
  }
  
  const useRegister = (options = {}) => {
    return createApiCall(apiClient.auth.register, {
      showSuccessToast: false, // Handled in useAuth
      ...options
    })
  }
  
  const useForgotPassword = (options = {}) => {
    return createApiCall(apiClient.auth.forgotPassword, {
      showSuccessToast: true,
      successMessage: 'Password reset link sent to your email',
      ...options
    })
  }
  
  // Tenant APIs
  const useCreateTemple = (options = {}) => {
    return createApiCall(apiClient.tenant.createTemple, {
      showSuccessToast: true,
      successMessage: 'Temple created successfully! It will be reviewed for approval.',
      ...options
    })
  }
  
  const useGetTemples = (options = {}) => {
    return createApiCall(apiClient.tenant.getTemples, {
      showErrorToast: true,
      ...options
    })
  }
  
  const useUploadDocuments = (options = {}) => {
    return createApiCall(apiClient.tenant.uploadDocuments, {
      showSuccessToast: true,
      successMessage: 'Documents uploaded successfully',
      ...options
    })
  }
  
  // Entity (Temple Admin) APIs
  const useGetDashboardStats = (entityId, options = {}) => {
    return createApiCall(() => apiClient.entity.getDashboard(entityId), {
      showErrorToast: false, // Dashboard should handle errors gracefully
      ...options
    })
  }
  
  const useGetDevotees = (entityId, options = {}) => {
    return createApiCall((params) => apiClient.entity.getDevotees(entityId, params), {
      showErrorToast: false,
      ...options
    })
  }
  
  const useCreateSeva = (entityId, options = {}) => {
    return createApiCall((sevaData) => apiClient.entity.createSeva(entityId, sevaData), {
      showSuccessToast: true,
      successMessage: 'Seva created successfully',
      ...options
    })
  }
  
  const useApproveSeva = (entityId, options = {}) => {
    return createApiCall((sevaId) => apiClient.entity.approveSeva(entityId, sevaId), {
      showSuccessToast: true,
      successMessage: 'Seva approved successfully',
      ...options
    })
  }
  
  const useCreateEvent = (entityId, options = {}) => {
    return createApiCall((eventData) => apiClient.entity.createEvent(entityId, eventData), {
      showSuccessToast: true,
      successMessage: 'Event created successfully',
      ...options
    })
  }
  
  const useSendBulkMessage = (entityId, options = {}) => {
    return createApiCall((messageData) => apiClient.entity.sendBulkMessage(entityId, messageData), {
      showSuccessToast: true,
      successMessage: 'Message sent successfully',
      ...options
    })
  }
  
  // Devotee APIs
  const useSearchTemples = (options = {}) => {
    return createApiCall((searchParams) => apiClient.devotee.searchTemples(searchParams), {
      showErrorToast: false,
      ...options
    })
  }
  
  const useJoinTemple = (options = {}) => {
    return createApiCall(apiClient.devotee.joinTemple, {
      showSuccessToast: false, // Handled in useAuth
      ...options
    })
  }
  
  const useGetProfileCompletion = (options = {}) => {
    return createApiCall(apiClient.devotee.getProfileCompletion, {
      showErrorToast: false,
      ...options
    })
  }
  
  const useBookSeva = (entityId, options = {}) => {
    return createApiCall((sevaData) => apiClient.devotee.bookSeva(entityId, sevaData), {
      showSuccessToast: true,
      successMessage: 'Seva booked successfully',
      ...options
    })
  }
  
  const useMakeDonation = (entityId, options = {}) => {
    return createApiCall((donationData) => apiClient.devotee.makeDonation(entityId, donationData), {
      showSuccessToast: true,
      successMessage: 'Donation made successfully',
      ...options
    })
  }
  
  const useGetDonationHistory = (entityId, options = {}) => {
    return createApiCall(() => apiClient.devotee.getDonationHistory(entityId), {
      showErrorToast: false,
      ...options
    })
  }
  
  // Volunteer APIs
  const useGetVolunteerAssignments = (entityId, options = {}) => {
    return createApiCall(() => apiClient.volunteer.getAssignments(entityId), {
      showErrorToast: false,
      ...options
    })
  }
  
  const useUpdateAssignmentStatus = (entityId, options = {}) => {
    return createApiCall((assignmentId, status) => apiClient.volunteer.updateAssignmentStatus(entityId, assignmentId, status), {
      showSuccessToast: true,
      successMessage: 'Assignment status updated',
      ...options
    })
  }
  
  // Generic file upload with progress
  const useFileUpload = (uploadFunction, options = {}) => {
    const progress = ref(0)
    const isUploading = ref(false)
    
    const upload = async (files, ...args) => {
      try {
        isUploading.value = true
        progress.value = 0
        
        const formData = new FormData()
        if (Array.isArray(files)) {
          files.forEach(file => formData.append('files', file))
        } else {
          formData.append('file', files)
        }
        
        const response = await uploadFunction(formData, ...args, {
          onUploadProgress: (progressEvent) => {
            progress.value = Math.round((progressEvent.loaded * 100) / progressEvent.total)
          }
        })
        
        if (options.showSuccessToast) {
          showToast(options.successMessage || 'File uploaded successfully', 'success')
        }
        
        return { success: true, data: response.data }
      } catch (error) {
        if (options.showErrorToast !== false) {
          showToast(error.response?.data?.message || 'File upload failed', 'error')
        }
        return { success: false, error: error.response?.data?.message || 'Upload failed' }
      } finally {
        isUploading.value = false
        progress.value = 0
      }
    }
    
    return {
      upload,
      progress: computed(() => progress.value),
      isUploading: computed(() => isUploading.value)
    }
  }
  
  // Pagination helper
  const usePagination = (fetchFunction, options = {}) => {
    const items = ref([])
    const currentPage = ref(1)
    const totalPages = ref(1)
    const totalItems = ref(0)
    const perPage = ref(options.perPage || 10)
    const isLoading = ref(false)
    const error = ref(null)
    
    const fetchPage = async (page = 1) => {
      try {
        isLoading.value = true
        error.value = null
        
        const params = {
          page,
          per_page: perPage.value,
          ...options.defaultParams
        }
        
        const response = await fetchFunction(params)
        
        items.value = response.data.data || response.data.items || []
        currentPage.value = response.data.current_page || page
        totalPages.value = response.data.last_page || response.data.total_pages || 1
        totalItems.value = response.data.total || 0
        
        return { success: true, data: response.data }
      } catch (err) {
        error.value = err.response?.data?.message || err.message
        if (options.showErrorToast !== false) {
          showToast(error.value, 'error')
        }
        return { success: false, error: error.value }
      } finally {
        isLoading.value = false
      }
    }
    
    const nextPage = () => {
      if (currentPage.value < totalPages.value) {
        fetchPage(currentPage.value + 1)
      }
    }
    
    const prevPage = () => {
      if (currentPage.value > 1) {
        fetchPage(currentPage.value - 1)
      }
    }
    
    const goToPage = (page) => {
      if (page >= 1 && page <= totalPages.value) {
        fetchPage(page)
      }
    }
    
    return {
      items: computed(() => items.value),
      currentPage: computed(() => currentPage.value),
      totalPages: computed(() => totalPages.value),
      totalItems: computed(() => totalItems.value),
      perPage: computed(() => perPage.value),
      isLoading: computed(() => isLoading.value),
      error: computed(() => error.value),
      fetchPage,
      nextPage,
      prevPage,
      goToPage
    }
  }
  
  return {
    // State
    loading: computed(() => loading.value),
    error: computed(() => error.value),
    
    // Generic API wrapper
    createApiCall,
    
    // Authentication
    useLogin,
    useRegister,
    useForgotPassword,
    
    // Tenant
    useCreateTemple,
    useGetTemples,
    useUploadDocuments,
    
    // Entity (Temple Admin)
    useGetDashboardStats,
    useGetDevotees,
    useCreateSeva,
    useApproveSeva,
    useCreateEvent,
    useSendBulkMessage,
    
    // Devotee
    useSearchTemples,
    useJoinTemple,
    useGetProfileCompletion,
    useBookSeva,
    useMakeDonation,
    useGetDonationHistory,
    
    // Volunteer
    useGetVolunteerAssignments,
    useUpdateAssignmentStatus,
    
    // Utilities
    useFileUpload,
    usePagination
  }
}