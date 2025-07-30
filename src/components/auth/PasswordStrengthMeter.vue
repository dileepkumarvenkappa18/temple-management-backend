<template>
  <div v-if="password" class="mt-2 space-y-2">
    <!-- Strength Bar -->
    <div class="flex space-x-1">
      <div
        v-for="i in 4"
        :key="i"
        :class="[
          'h-2 flex-1 rounded-full transition-all duration-300',
          getBarColor(i)
        ]"
      ></div>
    </div>

    <!-- Strength Text -->
    <div class="flex items-center justify-between">
      <span :class="['text-sm font-medium', getTextColor()]">
        {{ strengthText }}
      </span>
      <span class="text-xs text-gray-500">
        {{ password.length }}/8+ characters
      </span>
    </div>

    <!-- Requirements Checklist -->
    <div class="space-y-1">
      <div
        v-for="requirement in requirements"
        :key="requirement.key"
        :class="[
          'flex items-center text-xs transition-colors duration-200',
          requirement.met ? 'text-green-600' : 'text-gray-500'
        ]"
      >
        <svg
          :class="[
            'w-3 h-3 mr-2 transition-colors duration-200',
            requirement.met ? 'text-green-500' : 'text-gray-400'
          ]"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            v-if="requirement.met"
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M5 13l4 4L19 7"
          />
          <circle
            v-else
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            stroke-width="2"
          />
        </svg>
        {{ requirement.text }}
      </div>
    </div>

    <!-- Additional Security Tips -->
    <div v-if="strength < 3" class="bg-yellow-50 border border-yellow-200 rounded-lg p-3 mt-3">
      <div class="flex items-start">
        <svg class="w-4 h-4 text-yellow-600 mt-0.5 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L5.082 16.5c-.77.833.192 2.5 1.732 2.5z" />
        </svg>
        <div>
          <p class="text-sm text-yellow-800 font-medium">Security Tip</p>
          <p class="text-xs text-yellow-700 mt-1">
            Use a mix of uppercase, lowercase, numbers, and symbols for better security
          </p>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  password: {
    type: String,
    default: ''
  }
})

// Calculate password strength
const strength = computed(() => {
  if (!props.password) return 0
  
  let score = 0
  
  // Length check
  if (props.password.length >= 8) score += 1
  if (props.password.length >= 12) score += 1
  
  // Character variety checks
  if (/[a-z]/.test(props.password)) score += 1
  if (/[A-Z]/.test(props.password)) score += 1
  if (/\d/.test(props.password)) score += 1
  if (/[^a-zA-Z\d]/.test(props.password)) score += 1
  
  // Bonus for very long passwords
  if (props.password.length >= 16) score += 1
  
  return Math.min(score, 4)
})

// Password requirements
const requirements = computed(() => [
  {
    key: 'length',
    text: 'At least 8 characters',
    met: props.password.length >= 8
  },
  {
    key: 'lowercase',
    text: 'One lowercase letter',
    met: /[a-z]/.test(props.password)
  },
  {
    key: 'uppercase',
    text: 'One uppercase letter',
    met: /[A-Z]/.test(props.password)
  },
  {
    key: 'number',
    text: 'One number',
    met: /\d/.test(props.password)
  },
  {
    key: 'special',
    text: 'One special character',
    met: /[^a-zA-Z\d]/.test(props.password)
  }
])

// Strength text labels
const strengthText = computed(() => {
  const labels = ['Very Weak', 'Weak', 'Fair', 'Good', 'Strong']
  return labels[strength.value] || 'Very Weak'
})

// Bar colors based on strength
const getBarColor = (index) => {
  if (index <= strength.value) {
    switch (strength.value) {
      case 1:
        return 'bg-red-500'
      case 2:
        return 'bg-orange-500'
      case 3:
        return 'bg-yellow-500'
      case 4:
        return 'bg-green-500'
      default:
        return 'bg-gray-200'
    }
  }
  return 'bg-gray-200'
}

// Text colors based on strength
const getTextColor = () => {
  switch (strength.value) {
    case 1:
      return 'text-red-600'
    case 2:
      return 'text-orange-600'
    case 3:
      return 'text-yellow-600'
    case 4:
      return 'text-green-600'
    default:
      return 'text-gray-500'
  }
}

// Expose strength for parent component validation
defineExpose({
  strength: strength,
  isValid: computed(() => strength.value >= 3)
})
</script>