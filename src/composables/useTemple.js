// src/composables/useTemple.js
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

export function useTemple() {
  const router = useRouter()
  const selectedEntityId = ref(localStorage.getItem('selectedEntityId') || null)
  const selectedTempleName = ref(localStorage.getItem('selectedTempleName') || null)
  
  const hasSelectedTemple = computed(() => !!selectedEntityId.value)
  
  function selectTemple(temple) {
    if (temple && temple.id) {
      selectedEntityId.value = temple.id.toString()
      selectedTempleName.value = temple.name
      
      localStorage.setItem('selectedEntityId', temple.id.toString())
      localStorage.setItem('selectedTempleName', temple.name)
      
      // Get user role and redirect to appropriate dashboard
      const userRole = localStorage.getItem('userRole')
      if (userRole === 'devotee') {
        router.push(`/entity/${temple.id}/devotee/dashboard`)
      } else if (userRole === 'volunteer') {
        router.push(`/entity/${temple.id}/volunteer/dashboard`)
      }
    }
  }
  
  function clearSelectedTemple() {
    selectedEntityId.value = null
    selectedTempleName.value = null
    localStorage.removeItem('selectedEntityId')
    localStorage.removeItem('selectedTempleName')
  }
  
  return {
    selectedEntityId,
    selectedTempleName,
    hasSelectedTemple,
    selectTemple,
    clearSelectedTemple
  }
}