<template>
  <div v-if="show" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50" @click.self="closeModal">
    <div class="bg-white rounded-lg shadow-xl w-80 p-5 mx-4">
      <div class="flex justify-between items-center mb-4">
        <h3 class="text-lg font-medium text-gray-900">Download {{ title }}</h3>
        <button @click="closeModal" class="text-gray-400 hover:text-gray-500">
          <svg class="h-5 w-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
      </div>
      
      <p class="text-sm text-gray-600 mb-4">Choose your preferred file format:</p>
      
      <div class="space-y-2">
        <label class="flex items-center p-3 border rounded-md cursor-pointer hover:bg-gray-50" :class="modelValue === 'pdf' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200'">
          <input type="radio" name="format" value="pdf" v-model="localValue" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300">
          <span class="ml-3 flex items-center">
            <svg class="h-5 w-5 mr-2 text-red-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M4 4a2 2 0 012-2h8a2 2 0 012 2v12a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 0v12h8V4H6z" clip-rule="evenodd" />
            </svg>
            <span class="text-sm font-medium text-gray-700">PDF Format</span>
          </span>
        </label>
        
        <label class="flex items-center p-3 border rounded-md cursor-pointer hover:bg-gray-50" :class="modelValue === 'csv' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200'">
          <input type="radio" name="format" value="csv" v-model="localValue" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300">
          <span class="ml-3 flex items-center">
            <svg class="h-5 w-5 mr-2 text-green-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M3 5a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 3a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 3a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1zm0 3a1 1 0 011-1h12a1 1 0 110 2H4a1 1 0 01-1-1z" clip-rule="evenodd" />
            </svg>
            <span class="text-sm font-medium text-gray-700">CSV Format</span>
          </span>
        </label>
        
        <label class="flex items-center p-3 border rounded-md cursor-pointer hover:bg-gray-50" :class="modelValue === 'excel' ? 'border-indigo-500 bg-indigo-50' : 'border-gray-200'">
          <input type="radio" name="format" value="excel" v-model="localValue" class="h-4 w-4 text-indigo-600 focus:ring-indigo-500 border-gray-300">
          <span class="ml-3 flex items-center">
            <svg class="h-5 w-5 mr-2 text-blue-500" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M3 5a2 2 0 012-2h10a2 2 0 012 2v10a2 2 0 01-2 2H5a2 2 0 01-2-2V5zm11 1H6v8h8V6z" clip-rule="evenodd" />
            </svg>
            <span class="text-sm font-medium text-gray-700">Excel Format</span>
          </span>
        </label>
      </div>
      
      <div class="mt-5">
        <button @click="downloadFile" class="w-full flex justify-center items-center px-4 py-2 bg-indigo-600 text-white rounded-md hover:bg-indigo-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-indigo-500">
          <svg class="mr-2 h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-4l-4 4m0 0l-4-4m4 4V4" />
          </svg>
          Download Now
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed, ref, watch } from 'vue';

const props = defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Report'
  },
  modelValue: {
    type: String,
    default: 'pdf'
  },
  filter: {
    type: String,
    default: 'monthly'
  }
});

const emit = defineEmits(['update:modelValue', 'close', 'download']);

const localValue = computed({
  get: () => props.modelValue,
  set: (value) => emit('update:modelValue', value)
});

const closeModal = () => {
  emit('close');
};

const downloadFile = () => {
  emit('download', {
    format: localValue.value,
    title: props.title,
    filter: props.filter
  });
  closeModal();
};
</script>