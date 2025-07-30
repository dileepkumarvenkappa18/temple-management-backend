<!-- src/components/common/BaseInput.vue -->
<template>
  <div class="space-y-1">
    <!-- Label -->
    <label
      v-if="label"
      :for="inputId"
      class="block text-sm font-medium text-gray-700"
      :class="{ 'text-red-700': hasError }"
    >
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>

    <!-- Input container -->
    <div class="relative">
      <!-- Left icon -->
      <div
        v-if="iconLeft"
        class="absolute inset-y-0 left-0 pl-3 flex items-center pointer-events-none"
      >
        <component
          :is="iconLeft"
          class="h-5 w-5 text-gray-400"
          :class="{ 'text-red-400': hasError }"
        />
      </div>

      <!-- Input field -->
      <input
        :id="inputId"
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :readonly="readonly"
        :required="required"
        :autocomplete="autocomplete"
        :class="inputClasses"
        @input="handleInput"
        @blur="handleBlur"
        @focus="handleFocus"
        @keydown="handleKeydown"
      />

      <!-- Right icon / Toggle password visibility -->
      <div
        v-if="iconRight || isPassword"
        class="absolute inset-y-0 right-0 pr-3 flex items-center"
      >
        <!-- Password toggle -->
        <button
          v-if="isPassword"
          type="button"
          class="text-gray-400 hover:text-gray-600 focus:outline-none"
          @click="togglePasswordVisibility"
        >
          <EyeIcon v-if="showPassword" class="h-5 w-5" />
          <EyeSlashIcon v-else class="h-5 w-5" />
        </button>
        
        <!-- Regular right icon -->
        <component
          v-else-if="iconRight"
          :is="iconRight"
          class="h-5 w-5 text-gray-400"
          :class="{ 'text-red-400': hasError }"
        />
      </div>
    </div>

    <!-- Help text -->
    <p
      v-if="helpText && !hasError"
      class="text-sm text-gray-500"
    >
      {{ helpText }}
    </p>

    <!-- Error message -->
    <p
      v-if="hasError"
      class="text-sm text-red-600"
    >
      {{ error }}
    </p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { EyeIcon, EyeSlashIcon } from '@heroicons/vue/24/outline'

const props = defineProps({
  // v-model
  modelValue: {
    type: [String, Number],
    default: ''
  },
  
  // Input attributes
  type: {
    type: String,
    default: 'text',
    validator: (value) => [
      'text', 'email', 'password', 'number', 'tel', 'url', 'search'
    ].includes(value)
  },
  
  label: {
    type: String,
    default: ''
  },
  
  placeholder: {
    type: String,
    default: ''
  },
  
  helpText: {
    type: String,
    default: ''
  },
  
  // Validation
  error: {
    type: String,
    default: ''
  },
  
  required: {
    type: Boolean,
    default: false
  },
  
  // States
  disabled: {
    type: Boolean,
    default: false
  },
  
  readonly: {
    type: Boolean,
    default: false
  },
  
  // Icons
  iconLeft: {
    type: [String, Object],
    default: null
  },
  
  iconRight: {
    type: [String, Object],
    default: null
  },
  
  // Autocomplete
  autocomplete: {
    type: String,
    default: 'off'
  },
  
  // Size
  size: {
    type: String,
    default: 'md',
    validator: (value) => ['sm', 'md', 'lg'].includes(value)
  }
})

const emit = defineEmits(['update:modelValue', 'blur', 'focus', 'keydown'])

// Reactive state
const inputId = ref(`input-${Math.random().toString(36).substr(2, 9)}`)
const showPassword = ref(false)

// Computed properties
const hasError = computed(() => !!props.error)

const isPassword = computed(() => props.type === 'password')

const inputType = computed(() => {
  if (isPassword.value) {
    return showPassword.value ? 'text' : 'password'
  }
  return props.type
})

const inputClasses = computed(() => {
  const baseClasses = [
    'block w-full border rounded-md shadow-sm',
    'placeholder-gray-400',
    'focus:outline-none focus:ring-2 focus:ring-offset-2',
    'transition-colors duration-200'
  ]

  // Size classes
  const sizeClasses = {
    sm: 'px-3 py-2 text-sm',
    md: 'px-4 py-2.5 text-sm',
    lg: 'px-4 py-3 text-base'
  }
  baseClasses.push(sizeClasses[props.size] || sizeClasses.md)

  // Icon padding
  if (props.iconLeft) {
    baseClasses.push('pl-10')
  }
  if (props.iconRight || isPassword.value) {
    baseClasses.push('pr-10')
  }

  // State classes
  if (hasError.value) {
    baseClasses.push(
      'border-red-300 text-red-900',
      'focus:ring-red-500 focus:border-red-500'
    )
  } else {
    baseClasses.push(
      'border-gray-300 text-gray-900',
      'focus:ring-indigo-500 focus:border-indigo-500'
    )
  }

  if (props.disabled) {
    baseClasses.push('bg-gray-50 text-gray-500 cursor-not-allowed')
  } else if (props.readonly) {
    baseClasses.push('bg-gray-50')
  } else {
    baseClasses.push('bg-white')
  }

  return baseClasses.join(' ')
})

// Methods
const handleInput = (event) => {
  emit('update:modelValue', event.target.value)
}

const handleBlur = (event) => {
  emit('blur', event)
}

const handleFocus = (event) => {
  emit('focus', event)
}

const handleKeydown = (event) => {
  emit('keydown', event)
}

const togglePasswordVisibility = () => {
  showPassword.value = !showPassword.value
}

// Set unique ID on mount
onMounted(() => {
  if (!inputId.value.startsWith('input-')) {
    inputId.value = `input-${Math.random().toString(36).substr(2, 9)}`
  }
})
</script>