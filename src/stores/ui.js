// src/stores/ui.js
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useUIStore = defineStore('ui', () => {
  // ==================== SIDEBAR STATE ====================
  const sidebarOpen = ref(false)
  const sidebarCollapsed = ref(false)
  const sidebarPinned = ref(false)

  // ==================== MODAL STATE ====================
  const modals = ref(new Map())
  const modalZIndex = ref(1000)

  // ==================== LOADING STATE ====================
  const globalLoading = ref(false)
  const loadingTasks = ref(new Set())
  const pageLoading = ref(false)

  // ==================== NOTIFICATION STATE ====================
  const notifications = ref([])
  const notificationId = ref(0)

  // ==================== THEME STATE ====================
  const darkMode = ref(false)
  const colorScheme = ref('indigo') // Primary color scheme
  const fontSize = ref('base') // text-sm, text-base, text-lg
  const compactMode = ref(false)

  // ==================== BREADCRUMB STATE ====================
  const breadcrumbs = ref([])

  // ==================== MOBILE STATE ====================
  const isMobile = ref(window.innerWidth < 768)
  const isTablet = ref(window.innerWidth >= 768 && window.innerWidth < 1024)
  const isDesktop = ref(window.innerWidth >= 1024)

  // ==================== SEARCH STATE ====================
  const globalSearchOpen = ref(false)
  const globalSearchQuery = ref('')

  // ==================== FORM STATE ====================
  const unsavedChanges = ref(false)
  const formErrors = ref({})

  // ==================== PAGE STATE ====================
  const pageTitle = ref('')
  const pageDescription = ref('')
  const showBackButton = ref(false)

  // ==================== COMPUTED ====================
  const sidebarClasses = computed(() => {
    return [
      'fixed inset-y-0 left-0 z-50 w-64 bg-white shadow-lg border-r border-gray-200 transform transition-transform duration-300 ease-in-out',
      {
        'translate-x-0': sidebarOpen.value || isDesktop.value,
        '-translate-x-full': !sidebarOpen.value && !isDesktop.value,
        'w-16': sidebarCollapsed.value,
        'w-64': !sidebarCollapsed.value
      }
    ]
  })

  const overlayClasses = computed(() => {
    return [
      'fixed inset-0 bg-gray-900 bg-opacity-50 transition-opacity duration-300 z-40',
      {
        'opacity-100 pointer-events-auto': sidebarOpen.value && !isDesktop.value,
        'opacity-0 pointer-events-none': !sidebarOpen.value || isDesktop.value
      }
    ]
  })

  const mainContentClasses = computed(() => {
    return [
      'transition-all duration-300 ease-in-out',
      {
        'ml-64': isDesktop.value && !sidebarCollapsed.value,
        'ml-16': isDesktop.value && sidebarCollapsed.value,
        'ml-0': !isDesktop.value
      }
    ]
  })

  const currentBreakpoint = computed(() => {
    if (isMobile.value) return 'mobile'
    if (isTablet.value) return 'tablet'
    return 'desktop'
  })

  const themeClasses = computed(() => {
    const baseClasses = 'transition-colors duration-200'
    if (darkMode.value) {
      return `${baseClasses} dark bg-gray-900 text-white`
    }
    return `${baseClasses} bg-gray-50 text-gray-900`
  })

  // ==================== SIDEBAR ACTIONS ====================
  const toggleSidebar = () => {
    sidebarOpen.value = !sidebarOpen.value
  }

  const closeSidebar = () => {
    sidebarOpen.value = false
  }

  const openSidebar = () => {
    sidebarOpen.value = true
  }

  const toggleSidebarCollapse = () => {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }

  const setSidebarCollapsed = (collapsed) => {
    sidebarCollapsed.value = collapsed
  }

  const toggleSidebarPin = () => {
    sidebarPinned.value = !sidebarPinned.value
  }

  // ==================== MODAL ACTIONS ====================
  const openModal = (id, props = {}) => {
    const zIndex = modalZIndex.value++
    modals.value.set(id, {
      id,
      isOpen: true,
      zIndex,
      props,
      timestamp: Date.now()
    })
  }

  const closeModal = (id) => {
    if (modals.value.has(id)) {
      const modal = modals.value.get(id)
      modal.isOpen = false
      // Remove after animation completes
      setTimeout(() => {
        modals.value.delete(id)
      }, 300)
    }
  }

  const closeAllModals = () => {
    modals.value.forEach((modal) => {
      modal.isOpen = false
    })
    setTimeout(() => {
      modals.value.clear()
    }, 300)
  }

  const isModalOpen = (id) => {
    const modal = modals.value.get(id)
    return modal ? modal.isOpen : false
  }

  const getModalProps = (id) => {
    const modal = modals.value.get(id)
    return modal ? modal.props : {}
  }

  // ==================== LOADING ACTIONS ====================
  const setGlobalLoading = (loading) => {
    globalLoading.value = loading
  }

  const addLoadingTask = (taskId) => {
    loadingTasks.value.add(taskId)
    globalLoading.value = loadingTasks.value.size > 0
  }

  const removeLoadingTask = (taskId) => {
    loadingTasks.value.delete(taskId)
    globalLoading.value = loadingTasks.value.size > 0
  }

  const setPageLoading = (loading) => {
    pageLoading.value = loading
  }

  // ==================== NOTIFICATION ACTIONS ====================
  const addNotification = ({ type = 'info', title, message, duration = 5000, persistent = false, actions = [] }) => {
    const id = ++notificationId.value
    const notification = {
      id,
      type, // success, error, warning, info
      title,
      message,
      duration,
      persistent,
      actions,
      timestamp: Date.now(),
      dismissed: false
    }
    
    notifications.value.unshift(notification)
    
    // Auto dismiss non-persistent notifications
    if (!persistent && duration > 0) {
      setTimeout(() => {
        dismissNotification(id)
      }, duration)
    }
    
    return id
  }

  const dismissNotification = (id) => {
    const index = notifications.value.findIndex(n => n.id === id)
    if (index > -1) {
      notifications.value[index].dismissed = true
      // Remove after animation
      setTimeout(() => {
        notifications.value.splice(index, 1)
      }, 300)
    }
  }

  const clearAllNotifications = () => {
    notifications.value.forEach(n => n.dismissed = true)
    setTimeout(() => {
      notifications.value = []
    }, 300)
  }

  // Helper methods for common notification types
  const showSuccess = (title, message, options = {}) => {
    return addNotification({ type: 'success', title, message, ...options })
  }

  const showError = (title, message, options = {}) => {
    return addNotification({ type: 'error', title, message, persistent: true, ...options })
  }

  const showWarning = (title, message, options = {}) => {
    return addNotification({ type: 'warning', title, message, ...options })
  }

  const showInfo = (title, message, options = {}) => {
    return addNotification({ type: 'info', title, message, ...options })
  }

  // ==================== THEME ACTIONS ====================
  const toggleDarkMode = () => {
    darkMode.value = !darkMode.value
    localStorage.setItem('darkMode', darkMode.value.toString())
  }

  const setColorScheme = (scheme) => {
    colorScheme.value = scheme
    localStorage.setItem('colorScheme', scheme)
  }

  const setFontSize = (size) => {
    fontSize.value = size
    localStorage.setItem('fontSize', size)
  }

  const toggleCompactMode = () => {
    compactMode.value = !compactMode.value
    localStorage.setItem('compactMode', compactMode.value.toString())
  }

  // ==================== BREADCRUMB ACTIONS ====================
  const setBreadcrumbs = (crumbs) => {
    breadcrumbs.value = crumbs
  }

  const addBreadcrumb = (crumb) => {
    breadcrumbs.value.push(crumb)
  }

  const clearBreadcrumbs = () => {
    breadcrumbs.value = []
  }

  // ==================== RESPONSIVE ACTIONS ====================
  const updateBreakpoints = () => {
    const width = window.innerWidth
    isMobile.value = width < 768
    isTablet.value = width >= 768 && width < 1024
    isDesktop.value = width >= 1024
  }

  // ==================== SEARCH ACTIONS ====================
  const toggleGlobalSearch = () => {
    globalSearchOpen.value = !globalSearchOpen.value
    if (globalSearchOpen.value) {
      // Focus search input after DOM update
      setTimeout(() => {
        const searchInput = document.querySelector('#global-search-input')
        if (searchInput) searchInput.focus()
      }, 100)
    }
  }

  const setGlobalSearchQuery = (query) => {
    globalSearchQuery.value = query
  }

  const clearGlobalSearch = () => {
    globalSearchQuery.value = ''
    globalSearchOpen.value = false
  }

  // ==================== FORM ACTIONS ====================
  const setUnsavedChanges = (hasChanges) => {
    unsavedChanges.value = hasChanges
  }

  const setFormErrors = (errors) => {
    formErrors.value = errors
  }

  const clearFormErrors = () => {
    formErrors.value = {}
  }

  const addFormError = (field, error) => {
    formErrors.value[field] = error
  }

  const removeFormError = (field) => {
    delete formErrors.value[field]
  }

  // ==================== PAGE ACTIONS ====================
  const setPageTitle = (title) => {
    pageTitle.value = title
    document.title = title ? `${title} - Temple Management` : 'Temple Management'
  }

  const setPageDescription = (description) => {
    pageDescription.value = description
  }

  const setShowBackButton = (show) => {
    showBackButton.value = show
  }

  const setPageMeta = ({ title, description, showBack = false }) => {
    if (title) setPageTitle(title)
    if (description) setPageDescription(description)
    setShowBackButton(showBack)
  }

  // ==================== UTILITY ACTIONS ====================
  const resetUIState = () => {
    // Reset all non-persistent UI state
    closeSidebar()
    closeAllModals()
    setGlobalLoading(false)
    loadingTasks.value.clear()
    setPageLoading(false)
    clearAllNotifications()
    clearGlobalSearch()
    setUnsavedChanges(false)
    clearFormErrors()
    clearBreadcrumbs()
  }

  // ==================== INITIALIZATION ====================
  const initializeUI = () => {
    // Load theme preferences from localStorage
    const savedDarkMode = localStorage.getItem('darkMode')
    if (savedDarkMode !== null) {
      darkMode.value = savedDarkMode === 'true'
    }

    const savedColorScheme = localStorage.getItem('colorScheme')
    if (savedColorScheme) {
      colorScheme.value = savedColorScheme
    }

    const savedFontSize = localStorage.getItem('fontSize')
    if (savedFontSize) {
      fontSize.value = savedFontSize
    }

    const savedCompactMode = localStorage.getItem('compactMode')
    if (savedCompactMode !== null) {
      compactMode.value = savedCompactMode === 'true'
    }

    // Initialize responsive state
    updateBreakpoints()

    // Listen for window resize
    window.addEventListener('resize', updateBreakpoints)

    // Listen for beforeunload if there are unsaved changes
    window.addEventListener('beforeunload', (e) => {
      if (unsavedChanges.value) {
        e.preventDefault()
        e.returnValue = 'You have unsaved changes. Are you sure you want to leave?'
        return e.returnValue
      }
    })
  }

  return {
    // State
    sidebarOpen,
    sidebarCollapsed,
    sidebarPinned,
    modals,
    globalLoading,
    loadingTasks,
    pageLoading,
    notifications,
    darkMode,
    colorScheme,
    fontSize,
    compactMode,
    breadcrumbs,
    isMobile,
    isTablet,
    isDesktop,
    globalSearchOpen,
    globalSearchQuery,
    unsavedChanges,
    formErrors,
    pageTitle,
    pageDescription,
    showBackButton,

    // Computed
    sidebarClasses,
    overlayClasses,
    mainContentClasses,
    currentBreakpoint,
    themeClasses,

    // Actions
    toggleSidebar,
    closeSidebar,
    openSidebar,
    toggleSidebarCollapse,
    setSidebarCollapsed,
    toggleSidebarPin,
    openModal,
    closeModal,
    closeAllModals,
    isModalOpen,
    getModalProps,
    setGlobalLoading,
    addLoadingTask,
    removeLoadingTask,
    setPageLoading,
    addNotification,
    dismissNotification,
    clearAllNotifications,
    showSuccess,
    showError,
    showWarning,
    showInfo,
    toggleDarkMode,
    setColorScheme,
    setFontSize,
    toggleCompactMode,
    setBreadcrumbs,
    addBreadcrumb,
    clearBreadcrumbs,
    updateBreakpoints,
    toggleGlobalSearch,
    setGlobalSearchQuery,
    clearGlobalSearch,
    setUnsavedChanges,
    setFormErrors,
    clearFormErrors,
    addFormError,
    removeFormError,
    setPageTitle,
    setPageDescription,
    setShowBackButton,
    setPageMeta,
    resetUIState,
    initializeUI
  }
})