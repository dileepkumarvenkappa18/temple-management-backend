<template>
  <div class="w-full">
    <!-- Progress Bar -->
    <div class="relative mb-8">
      <div class="absolute top-1/2 left-0 w-full h-0.5 bg-gray-200 -translate-y-1/2"></div>
      <div 
        class="absolute top-1/2 left-0 h-0.5 bg-indigo-600 -translate-y-1/2 transition-all duration-300"
        :style="{ width: `${progressPercentage}%` }"
      ></div>
      
      <!-- Steps -->
      <div class="relative flex justify-between">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="flex flex-col items-center group"
        >
          <!-- Step Circle -->
          <div
            class="w-10 h-10 rounded-full border-2 flex items-center justify-center bg-white transition-all duration-300 relative z-10"
            :class="getStepCircleClass(index)"
          >
            <!-- Completed Icon -->
            <svg
              v-if="index < currentStep"
              class="w-5 h-5 text-white"
              fill="currentColor"
              viewBox="0 0 20 20"
            >
              <path
                fill-rule="evenodd"
                d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                clip-rule="evenodd"
              />
            </svg>
            <!-- Current/Future Step Number -->
            <span
              v-else
              class="text-sm font-medium transition-colors duration-300"
              :class="index === currentStep ? 'text-white' : 'text-gray-400'"
            >
              {{ index + 1 }}
            </span>
          </div>
          
          <!-- Step Label -->
          <div class="mt-3 text-center">
            <p
              class="text-sm font-medium transition-colors duration-300"
              :class="getStepLabelClass(index)"
            >
              {{ step.title }}
            </p>
            <p
              v-if="step.description"
              class="text-xs text-gray-500 mt-1 max-w-24 leading-tight"
            >
              {{ step.description }}
            </p>
          </div>
        </div>
      </div>
    </div>

    <!-- Progress Percentage Display -->
    <div class="flex justify-between items-center mb-6">
      <div class="text-sm text-gray-600">
        Step {{ currentStep + 1 }} of {{ steps.length }}
      </div>
      <div class="flex items-center space-x-2">
        <div class="text-sm text-gray-600">Progress:</div>
        <div class="text-sm font-semibold text-indigo-600">
          {{ Math.round(progressPercentage) }}%
        </div>
      </div>
    </div>

    <!-- Current Step Content -->
    <div class="bg-white rounded-xl shadow-sm border border-gray-100 p-6">
      <div class="mb-4">
        <h3 class="text-lg font-semibold text-gray-900">
          {{ steps[currentStep]?.title }}
        </h3>
        <p v-if="steps[currentStep]?.description" class="text-sm text-gray-600 mt-1">
          {{ steps[currentStep]?.description }}
        </p>
      </div>
      
      <slot :current-step="currentStep" :step="steps[currentStep]"></slot>
    </div>

    <!-- Navigation Buttons -->
    <div class="flex justify-between mt-6">
      <button
        @click="goToPrevious"
        :disabled="currentStep === 0"
        class="px-4 py-2 text-sm font-medium text-gray-700 bg-white border border-gray-300 rounded-lg hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
      >
        Previous
      </button>
      
      <div class="flex space-x-3">
        <button
          v-if="allowSkip && currentStep < steps.length - 1"
          @click="skipStep"
          class="px-4 py-2 text-sm font-medium text-gray-600 bg-gray-100 rounded-lg hover:bg-gray-200 focus:outline-none focus:ring-2 focus:ring-gray-500 focus:ring-offset-2 transition-colors duration-200"
        >
          Skip
        </button>
        
        <button
          @click="goToNext"
          :disabled="!canProceed"
          class="px-6 py-2 text-sm font-medium text-white bg-indigo-600 rounded-lg hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
        >
          {{ currentStep === steps.length - 1 ? 'Complete' : 'Next' }}
        </button>
      </div>
    </div>

    <!-- Step Indicators (Mobile) -->
    <div class="md:hidden mt-6">
      <div class="flex justify-center space-x-2">
        <div
          v-for="(step, index) in steps"
          :key="index"
          class="w-2 h-2 rounded-full transition-colors duration-300"
          :class="index <= currentStep ? 'bg-indigo-600' : 'bg-gray-300'"
        ></div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  steps: {
    type: Array,
    required: true,
    validator: (steps) => steps.every(step => step.title)
  },
  currentStep: {
    type: Number,
    default: 0
  },
  allowSkip: {
    type: Boolean,
    default: true
  },
  canProceed: {
    type: Boolean,
    default: true
  }
})

const emit = defineEmits(['update:currentStep', 'next', 'previous', 'skip', 'complete'])

const progressPercentage = computed(() => {
  return ((props.currentStep + 1) / props.steps.length) * 100
})

const getStepCircleClass = (index) => {
  if (index < props.currentStep) {
    return 'border-indigo-600 bg-indigo-600'
  } else if (index === props.currentStep) {
    return 'border-indigo-600 bg-indigo-600'
  } else {
    return 'border-gray-300 bg-white'
  }
}

const getStepLabelClass = (index) => {
  if (index < props.currentStep) {
    return 'text-indigo-600'
  } else if (index === props.currentStep) {
    return 'text-indigo-600'
  } else {
    return 'text-gray-500'
  }
}

const goToNext = () => {
  if (props.currentStep < props.steps.length - 1) {
    emit('update:currentStep', props.currentStep + 1)
    emit('next', props.currentStep + 1)
  } else {
    emit('complete')
  }
}

const goToPrevious = () => {
  if (props.currentStep > 0) {
    emit('update:currentStep', props.currentStep - 1)
    emit('previous', props.currentStep - 1)
  }
}

const skipStep = () => {
  emit('skip', props.currentStep)
  goToNext()
}
</script>