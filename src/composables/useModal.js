import { ref, reactive, nextTick } from 'vue'

/**
 * Modal Management Composable
 * Provides reactive modal state management with support for multiple modals,
 * confirmation dialogs, and form modals with validation states
 */

// Global modal registry to track all active modals
const modalRegistry = reactive(new Map())
const activeModalStack = ref([])

export function useModal(modalId = null) {
  // Generate unique modal ID if not provided
  const id = modalId || `modal_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`
  
  // Initialize modal state if not exists
  if (!modalRegistry.has(id)) {
    modalRegistry.set(id, {
      isOpen: false,
      isLoading: false,
      data: null,
      config: {
        title: '',
        size: 'md', // sm, md, lg, xl, full
        closable: true,
        persistent: false, // Don't close on backdrop click
        showHeader: true,
        showFooter: true,
        customClass: '',
        zIndex: 50,
        backdrop: 'blur', // blur, dark, transparent
        animation: 'fade', // fade, slide, scale
        position: 'center', // center, top, bottom
        maxHeight: 'auto',
        scrollable: true
      },
      callbacks: {
        onOpen: null,
        onClose: null,
        onConfirm: null,
        onCancel: null,
        beforeClose: null
      }
    })
  }

  const modalState = modalRegistry.get(id)

  /**
   * Open modal with optional configuration and data
   * @param {Object} options - Modal configuration options
   * @param {any} data - Data to pass to modal content
   */
  const openModal = async (options = {}, data = null) => {
    try {
      // Update modal configuration
      Object.assign(modalState.config, options)
      modalState.data = data
      
      // Add to active stack if not already present
      if (!activeModalStack.value.includes(id)) {
        activeModalStack.value.push(id)
      }
      
      // Update z-index based on stack position
      const stackIndex = activeModalStack.value.indexOf(id)
      modalState.config.zIndex = 50 + stackIndex * 10
      
      modalState.isOpen = true
      
      // Execute onOpen callback
      if (modalState.callbacks.onOpen) {
        await modalState.callbacks.onOpen(data)
      }
      
      // Prevent body scroll when modal is open
      await nextTick()
      document.body.classList.add('modal-open')
      
      return true
    } catch (error) {
      console.error('Error opening modal:', error)
      return false
    }
  }

  /**
   * Close modal with optional result data
   * @param {any} result - Result data to return
   * @param {boolean} force - Force close without beforeClose validation
   */
  const closeModal = async (result = null, force = false) => {
    try {
      // Execute beforeClose validation if not forced
      if (!force && modalState.callbacks.beforeClose) {
        const canClose = await modalState.callbacks.beforeClose(result)
        if (!canClose) return false
      }
      
      modalState.isOpen = false
      modalState.isLoading = false
      modalState.data = null
      
      // Remove from active stack
      const stackIndex = activeModalStack.value.indexOf(id)
      if (stackIndex > -1) {
        activeModalStack.value.splice(stackIndex, 1)
      }
      
      // Remove body scroll prevention if no more modals
      if (activeModalStack.value.length === 0) {
        document.body.classList.remove('modal-open')
      }
      
      // Execute onClose callback
      if (modalState.callbacks.onClose) {
        await modalState.callbacks.onClose(result)
      }
      
      return true
    } catch (error) {
      console.error('Error closing modal:', error)
      return false
    }
  }

  /**
   * Toggle modal open/close state
   * @param {Object} options - Modal options when opening
   * @param {any} data - Data when opening
   */
  const toggleModal = async (options = {}, data = null) => {
    if (modalState.isOpen) {
      return await closeModal()
    } else {
      return await openModal(options, data)
    }
  }

  /**
   * Set modal loading state
   * @param {boolean} loading - Loading state
   */
  const setLoading = (loading = true) => {
    modalState.isLoading = loading
  }

  /**
   * Update modal data
   * @param {any} newData - New data to set
   */
  const updateData = (newData) => {
    modalState.data = newData
  }

  /**
   * Update modal configuration
   * @param {Object} newConfig - New configuration options
   */
  const updateConfig = (newConfig) => {
    Object.assign(modalState.config, newConfig)
  }

  /**
   * Set callback functions
   * @param {Object} callbacks - Callback functions
   */
  const setCallbacks = (callbacks) => {
    Object.assign(modalState.callbacks, callbacks)
  }

  /**
   * Confirm modal - specialized for confirmation dialogs
   * @param {Object} config - Confirmation dialog config
   */
  const confirmModal = (config = {}) => {
    return new Promise((resolve) => {
      const confirmConfig = {
        title: config.title || 'Confirm Action',
        size: config.size || 'sm',
        persistent: true,
        showHeader: true,
        showFooter: true,
        customClass: 'confirm-modal',
        ...config
      }
      
      setCallbacks({
        onConfirm: () => {
          closeModal(true)
          resolve(true)
        },
        onCancel: () => {
          closeModal(false)
          resolve(false)
        },
        onClose: (result) => {
          resolve(result || false)
        }
      })
      
      openModal(confirmConfig, {
        type: 'confirm',
        message: config.message || 'Are you sure you want to proceed?',
        confirmText: config.confirmText || 'Confirm',
        cancelText: config.cancelText || 'Cancel',
        variant: config.variant || 'primary', // primary, danger, warning
        icon: config.icon || null
      })
    })
  }

  /**
   * Alert modal - specialized for alert/info dialogs
   * @param {Object} config - Alert dialog config
   */
  const alertModal = (config = {}) => {
    return new Promise((resolve) => {
      const alertConfig = {
        title: config.title || 'Information',
        size: config.size || 'sm',
        closable: true,
        showHeader: true,
        showFooter: true,
        customClass: 'alert-modal',
        ...config
      }
      
      setCallbacks({
        onClose: () => resolve(true)
      })
      
      openModal(alertConfig, {
        type: 'alert',
        message: config.message || 'This is an alert message.',
        buttonText: config.buttonText || 'OK',
        variant: config.variant || 'info', // info, success, warning, error
        icon: config.icon || null
      })
    })
  }

  /**
   * Form modal - specialized for form dialogs
   * @param {Object} config - Form dialog config
   */
  const formModal = (config = {}) => {
    return new Promise((resolve) => {
      const formConfig = {
        title: config.title || 'Form',
        size: config.size || 'md',
        persistent: true,
        showHeader: true,
        showFooter: true,
        customClass: 'form-modal',
        scrollable: true,
        ...config
      }
      
      setCallbacks({
        onConfirm: (formData) => {
          resolve({ success: true, data: formData })
        },
        onCancel: () => {
          resolve({ success: false, data: null })
        },
        beforeClose: config.beforeClose || null
      })
      
      openModal(formConfig, {
        type: 'form',
        formData: config.formData || {},
        fields: config.fields || [],
        submitText: config.submitText || 'Submit',
        cancelText: config.cancelText || 'Cancel'
      })
    })
  }

  /**
   * Get size classes for modal
   * @param {string} size - Size key
   */
  const getSizeClasses = (size) => {
    const sizeMap = {
      sm: 'max-w-md',
      md: 'max-w-lg',
      lg: 'max-w-2xl',
      xl: 'max-w-4xl',
      full: 'max-w-full mx-4 h-full'
    }
    return sizeMap[size] || sizeMap.md
  }

  /**
   * Get animation classes
   * @param {string} animation - Animation type
   */
  const getAnimationClasses = (animation) => {
    const animationMap = {
      fade: {
        enter: 'transition-opacity duration-300 ease-out',
        enterFrom: 'opacity-0',
        enterTo: 'opacity-100',
        leave: 'transition-opacity duration-200 ease-in',
        leaveFrom: 'opacity-100',
        leaveTo: 'opacity-0'
      },
      slide: {
        enter: 'transition-all duration-300 ease-out',
        enterFrom: 'opacity-0 translate-y-4',
        enterTo: 'opacity-100 translate-y-0',
        leave: 'transition-all duration-200 ease-in',
        leaveFrom: 'opacity-100 translate-y-0',
        leaveTo: 'opacity-0 translate-y-4'
      },
      scale: {
        enter: 'transition-all duration-300 ease-out',
        enterFrom: 'opacity-0 scale-95',
        enterTo: 'opacity-100 scale-100',
        leave: 'transition-all duration-200 ease-in',
        leaveFrom: 'opacity-100 scale-100',
        leaveTo: 'opacity-0 scale-95'
      }
    }
    return animationMap[animation] || animationMap.fade
  }

  /**
   * Handle backdrop click
   */
  const handleBackdropClick = () => {
    if (!modalState.config.persistent && modalState.config.closable) {
      closeModal()
    }
  }

  /**
   * Handle escape key
   */
  const handleEscapeKey = (event) => {
    if (event.key === 'Escape' && modalState.config.closable && modalState.isOpen) {
      closeModal()
    }
  }

  // Cleanup function to remove modal from registry
  const cleanup = () => {
    if (modalRegistry.has(id)) {
      closeModal(null, true)
      modalRegistry.delete(id)
    }
  }

  return {
    // State
    id,
    isOpen: ref(() => modalState.isOpen),
    isLoading: ref(() => modalState.isLoading),
    data: ref(() => modalState.data),
    config: ref(() => modalState.config),
    
    // Actions
    openModal,
    closeModal,
    toggleModal,
    setLoading,
    updateData,
    updateConfig,
    setCallbacks,
    
    // Specialized modals
    confirmModal,
    alertModal,
    formModal,
    
    // Utilities
    getSizeClasses,
    getAnimationClasses,
    handleBackdropClick,
    handleEscapeKey,
    cleanup,
    
    // Global state
    activeModals: activeModalStack,
    modalCount: ref(() => activeModalStack.value.length)
  }
}

/**
 * Global modal utilities
 */
export const useGlobalModal = () => {
  const closeAllModals = async () => {
    const promises = Array.from(modalRegistry.keys()).map(id => {
      const modal = modalRegistry.get(id)
      if (modal.isOpen) {
        modal.isOpen = false
        modal.isLoading = false
        modal.data = null
      }
    })
    
    activeModalStack.value = []
    document.body.classList.remove('modal-open')
    
    return Promise.all(promises)
  }

  const getActiveModals = () => {
    return activeModalStack.value.map(id => ({
      id,
      ...modalRegistry.get(id)
    }))
  }

  const hasActiveModals = ref(() => activeModalStack.value.length > 0)

  return {
    closeAllModals,
    getActiveModals,
    hasActiveModals,
    activeCount: ref(() => activeModalStack.value.length)
  }
}