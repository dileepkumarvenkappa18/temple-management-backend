<template>
  <div class="relative w-full">
    <!-- Label -->
    <label 
      v-if="label" 
      :for="inputId"
      class="block text-sm font-medium text-gray-900 mb-2"
    >
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>

    <!-- Select Container -->
    <div class="relative">
      <!-- Select Button -->
      <button
        :id="inputId"
        type="button"
        class="relative w-full bg-white border rounded-xl px-4 py-3 text-left cursor-pointer focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-all duration-200"
        :class="[
          error 
            ? 'border-red-300 focus:ring-red-500 focus:border-red-500' 
            : 'border-gray-300 hover:border-gray-400',
          disabled 
            ? 'bg-gray-100 cursor-not-allowed opacity-50' 
            : 'hover:bg-gray-50'
        ]"
        :disabled="disabled"
        @click="toggleDropdown"
        @blur="handleBlur"
      >
        <!-- Selected Value Display -->
        <span 
          class="block truncate text-sm"
          :class="selectedOption ? 'text-gray-900' : 'text-gray-500'"
        >
          {{ selectedOption ? selectedOption.label : placeholder }}
        </span>

        <!-- Chevron Icon -->
        <span class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
          <svg 
            class="h-5 w-5 text-gray-400 transition-transform duration-200"
            :class="{ 'rotate-180': isOpen }"
            xmlns="http://www.w3.org/2000/svg" 
            viewBox="0 0 20 20" 
            fill="currentColor"
          >
            <path fill-rule="evenodd" d="M5.293 7.293a1 1 0 011.414 0L10 10.586l3.293-3.293a1 1 0 111.414 1.414l-4 4a1 1 0 01-1.414 0l-4-4a1 1 0 010-1.414z" clip-rule="evenodd" />
          </svg>
        </span>
      </button>

      <!-- Dropdown Menu -->
      <Transition
        enter-active-class="transition ease-out duration-200"
        enter-from-class="opacity-0 translate-y-1"
        enter-to-class="opacity-100 translate-y-0"
        leave-active-class="transition ease-in duration-150"
        leave-from-class="opacity-100 translate-y-0"
        leave-to-class="opacity-0 translate-y-1"
      >
        <div
          v-if="isOpen"
          class="absolute z-50 mt-2 w-full bg-white shadow-lg rounded-xl border border-gray-200 py-1 max-h-60 overflow-auto"
        >
          <!-- Search Input (if searchable) -->
          <div v-if="searchable" class="p-2 border-b border-gray-100">
            <input
              v-model="searchQuery"
              type="text"
              class="w-full px-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
              :placeholder="`Search ${label?.toLowerCase() || 'options'}...`"
              @click.stop
            />
          </div>

          <!-- Options List -->
          <div class="max-h-48 overflow-y-auto">
            <!-- No Results -->
            <div 
              v-if="filteredOptions.length === 0" 
              class="px-4 py-3 text-sm text-gray-500 text-center"
            >
              {{ searchQuery ? 'No matches found' : 'No options available' }}
            </div>

            <!-- Option Items -->
            <button
              v-for="option in filteredOptions"
              :key="option.value"
              type="button"
              class="relative w-full text-left px-4 py-3 text-sm hover:bg-indigo-50 focus:outline-none focus:bg-indigo-50 transition-colors duration-150"
              :class="{
                'bg-indigo-100 text-indigo-900': selectedOption?.value === option.value,
                'text-gray-900': selectedOption?.value !== option.value
              }"
              @click="selectOption(option)"
            >
              <span class="block truncate">{{ option.label }}</span>
              
              <!-- Selected Checkmark -->
              <span 
                v-if="selectedOption?.value === option.value"
                class="absolute inset-y-0 right-0 flex items-center pr-4"
              >
                <svg class="h-5 w-5 text-indigo-600" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20" fill="currentColor">
                  <path fill-rule="evenodd" d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </span>
            </button>
          </div>
        </div>
      </Transition>
    </div>

    <!-- Error Message -->
    <p v-if="error" class="mt-2 text-sm text-red-600">
      {{ error }}
    </p>

    <!-- Help Text -->
    <p v-if="helpText && !error" class="mt-2 text-sm text-gray-500">
      {{ helpText }}
    </p>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'

const props = defineProps({
  modelValue: {
    type: [String, Number, Object],
    default: null
  },
  options: {
    type: Array,
    default: () => []
  },
  label: {
    type: String,
    default: ''
  },
  placeholder: {
    type: String,
    default: 'Select an option'
  },
  error: {
    type: String,
    default: ''
  },
  helpText: {
    type: String,
    default: ''
  },
  disabled: {
    type: Boolean,
    default: false
  },
  required: {
    type: Boolean,
    default: false
  },
  searchable: {
    type: Boolean,
    default: false
  },
  clearable: {
    type: Boolean,
    default: false
  }
})

const emit = defineEmits(['update:modelValue', 'change'])

// Component state
const isOpen = ref(false)
const searchQuery = ref('')
const inputId = ref(`select-${Math.random().toString(36).substr(2, 9)}`)

// Computed properties
const selectedOption = computed(() => {
  if (!props.modelValue) return null
  return props.options.find(option => option.value === props.modelValue)
})

const filteredOptions = computed(() => {
  if (!props.searchable || !searchQuery.value) {
    return props.options
  }
  
  return props.options.filter(option =>
    option.label.toLowerCase().includes(searchQuery.value.toLowerCase())
  )
})

// Methods
const toggleDropdown = () => {
  if (props.disabled) return
  isOpen.value = !isOpen.value
}

const selectOption = (option) => {
  emit('update:modelValue', option.value)
  emit('change', option)
  isOpen.value = false
  searchQuery.value = ''
}

const handleBlur = (event) => {
  // Close dropdown when clicking outside
  setTimeout(() => {
    if (!event.relatedTarget?.closest('.relative')) {
      isOpen.value = false
    }
  }, 100)
}

const handleClickOutside = (event) => {
  if (!event.target.closest(`#${inputId.value}`)?.parentElement) {
    isOpen.value = false
  }
}

// Lifecycle
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>