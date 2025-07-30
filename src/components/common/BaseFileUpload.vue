<template>
  <div class="w-full">
    <!-- Label -->
    <label v-if="label" class="block text-sm font-medium text-gray-700 mb-2">
      {{ label }}
      <span v-if="required" class="text-red-500 ml-1">*</span>
    </label>

    <!-- Upload Area -->
    <div
      @drop="handleDrop"
      @dragover="handleDragOver"
      @dragenter="handleDragEnter"
      @dragleave="handleDragLeave"
      @click="triggerFileInput"
      class="relative border-2 border-dashed rounded-xl p-6 transition-all duration-300 cursor-pointer"
      :class="[
        isDragging ? 'border-indigo-400 bg-indigo-50' : 'border-gray-300 bg-gray-50',
        disabled ? 'opacity-50 cursor-not-allowed' : 'hover:border-indigo-400 hover:bg-indigo-50'
      ]"
    >
      <!-- Upload Icon and Text -->
      <div v-if="!files.length" class="text-center">
        <svg
          class="mx-auto h-12 w-12 text-gray-400 mb-4"
          stroke="currentColor"
          fill="none"
          viewBox="0 0 48 48"
        >
          <path
            d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        
        <div class="text-sm text-gray-600">
          <span class="font-medium text-indigo-600">Click to upload</span>
          or drag and drop
        </div>
        
        <p class="text-xs text-gray-500 mt-1">
          {{ acceptText }} (max {{ formatFileSize(maxSize) }})
        </p>
      </div>

      <!-- File Preview -->
      <div v-else class="space-y-3">
        <div
          v-for="(file, index) in files"
          :key="index"
          class="flex items-center justify-between p-3 bg-white rounded-lg border border-gray-200"
        >
          <div class="flex items-center space-x-3">
            <!-- File Icon -->
            <div class="flex-shrink-0">
              <div
                v-if="isImage(file)"
                class="w-10 h-10 rounded-lg overflow-hidden bg-gray-100"
              >
                <img
                  :src="getImagePreview(file)"
                  :alt="file.name"
                  class="w-full h-full object-cover"
                />
              </div>
              <div
                v-else
                class="w-10 h-10 rounded-lg bg-indigo-100 flex items-center justify-center"
              >
                <svg class="w-5 h-5 text-indigo-600" fill="currentColor" viewBox="0 0 20 20">
                  <path
                    fill-rule="evenodd"
                    d="M4 4a2 2 0 012-2h4.586A2 2 0 0112 2.586L15.414 6A2 2 0 0116 7.414V16a2 2 0 01-2 2H6a2 2 0 01-2-2V4zm2 6a1 1 0 011-1h6a1 1 0 110 2H7a1 1 0 01-1-1zm1 3a1 1 0 100 2h6a1 1 0 100-2H7z"
                    clip-rule="evenodd"
                  />
                </svg>
              </div>
            </div>

            <!-- File Info -->
            <div class="flex-1 min-w-0">
              <p class="text-sm font-medium text-gray-900 truncate">
                {{ file.name }}
              </p>
              <p class="text-sm text-gray-500">
                {{ formatFileSize(file.size) }}
                <span v-if="file.progress !== undefined && file.progress < 100">
                  • {{ file.progress }}% uploaded
                </span>
              </p>
            </div>
          </div>

          <!-- Actions -->
          <div class="flex items-center space-x-2">
            <!-- Upload Progress -->
            <div
              v-if="file.progress !== undefined && file.progress < 100"
              class="w-16 bg-gray-200 rounded-full h-2"
            >
              <div
                class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
                :style="{ width: `${file.progress}%` }"
              ></div>
            </div>

            <!-- Success Icon -->
            <div
              v-else-if="file.progress === 100"
              class="text-green-500"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M16.707 5.293a1 1 0 010 1.414l-8 8a1 1 0 01-1.414 0l-4-4a1 1 0 011.414-1.414L8 12.586l7.293-7.293a1 1 0 011.414 0z"
                  clip-rule="evenodd"
                />
              </svg>
            </div>

            <!-- Remove Button -->
            <button
              @click.stop="removeFile(index)"
              :disabled="disabled"
              class="text-gray-400 hover:text-red-500 transition-colors duration-200 disabled:opacity-50"
            >
              <svg class="w-5 h-5" fill="currentColor" viewBox="0 0 20 20">
                <path
                  fill-rule="evenodd"
                  d="M4.293 4.293a1 1 0 011.414 0L10 8.586l4.293-4.293a1 1 0 111.414 1.414L11.414 10l4.293 4.293a1 1 0 01-1.414 1.414L10 11.414l-4.293 4.293a1 1 0 01-1.414-1.414L8.586 10 4.293 5.707a1 1 0 010-1.414z"
                  clip-rule="evenodd"
                />
              </svg>
            </button>
          </div>
        </div>

        <!-- Add More Files -->
        <button
          v-if="multiple && files.length < maxFiles"
          @click="triggerFileInput"
          class="w-full py-2 px-4 border border-dashed border-gray-300 rounded-lg text-sm text-gray-600 hover:border-indigo-400 hover:text-indigo-600 transition-colors duration-200 disabled:opacity-50"
          :disabled="disabled"
        >
          + Add More Files
        </button>
      </div>

      <!-- Loading Overlay -->
      <div
        v-if="uploading"
        class="absolute inset-0 bg-white bg-opacity-90 rounded-xl flex items-center justify-center"
      >
        <div class="flex items-center space-x-2 text-indigo-600">
          <svg class="animate-spin w-5 h-5" fill="none" viewBox="0 0 24 24">
            <circle
              class="opacity-25"
              cx="12"
              cy="12"
              r="10"
              stroke="currentColor"
              stroke-width="4"
            ></circle>
            <path
              class="opacity-75"
              fill="currentColor"
              d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
            ></path>
          </svg>
          <span class="text-sm font-medium">Uploading...</span>
        </div>
      </div>
    </div>

    <!-- Hidden File Input -->
    <input
      ref="fileInput"
      type="file"
      :accept="accept"
      :multiple="multiple"
      @change="handleFileSelect"
      class="hidden"
      :disabled="disabled"
    />

    <!-- Error Message -->
    <div v-if="error" class="mt-2 text-sm text-red-600">
      {{ error }}
    </div>

    <!-- Help Text -->
    <div v-if="helpText" class="mt-2 text-sm text-gray-500">
      {{ helpText }}
    </div>
  </div>
</template>

<script setup>
import { ref, computed, watch } from 'vue'

const props = defineProps({
  modelValue: {
    type: Array,
    default: () => []
  },
  label: {
    type: String,
    default: ''
  },
  accept: {
    type: String,
    default: '*/*'
  },
  multiple: {
    type: Boolean,
    default: false
  },
  maxSize: {
    type: Number,
    default: 10 * 1024 * 1024 // 10MB
  },
  maxFiles: {
    type: Number,
    default: 5
  },
  required: {
    type: Boolean,
    default: false
  },
  disabled: {
    type: Boolean,
    default: false
  },
  uploading: {
    type: Boolean,
    default: false
  },
  helpText: {
    type: String,
    default: ''
  }
})

const emit = defineEmits(['update:modelValue', 'upload', 'remove', 'error'])

const fileInput = ref(null)
const files = ref([...props.modelValue])
const isDragging = ref(false)
const error = ref('')

const acceptText = computed(() => {
  if (props.accept === '*/*') return 'Any file type'
  if (props.accept.includes('image/*')) return 'Images'
  if (props.accept.includes('application/pdf')) return 'PDF files'
  return props.accept
})

watch(() => props.modelValue, (newFiles) => {
  files.value = [...newFiles]
}, { deep: true })

watch(files, (newFiles) => {
  emit('update:modelValue', newFiles)
}, { deep: true })

const triggerFileInput = () => {
  if (!props.disabled) {
    fileInput.value?.click()
  }
}

const handleFileSelect = (event) => {
  const selectedFiles = Array.from(event.target.files)
  processFiles(selectedFiles)
  // Reset input
  event.target.value = ''
}

const handleDrop = (event) => {
  event.preventDefault()
  isDragging.value = false
  
  if (props.disabled) return
  
  const droppedFiles = Array.from(event.dataTransfer.files)
  processFiles(droppedFiles)
}

const handleDragOver = (event) => {
  event.preventDefault()
}

const handleDragEnter = (event) => {
  event.preventDefault()
  if (!props.disabled) {
    isDragging.value = true
  }
}

const handleDragLeave = (event) => {
  event.preventDefault()
  // Only set isDragging to false if we're leaving the drop zone entirely
  if (!event.currentTarget.contains(event.relatedTarget)) {
    isDragging.value = false
  }
}

const processFiles = (newFiles) => {
  error.value = ''
  
  // Check file count limit
  const totalFiles = files.value.length + newFiles.length
  if (totalFiles > props.maxFiles) {
    error.value = `Maximum ${props.maxFiles} files allowed`
    return
  }

  const validFiles = []
  
  for (const file of newFiles) {
    // Check file size
    if (file.size > props.maxSize) {
      error.value = `File "${file.name}" is too large. Maximum size is ${formatFileSize(props.maxSize)}`
      continue
    }

    // Check file type
    if (!isFileTypeAllowed(file)) {
      error.value = `File type not allowed: ${file.name}`
      continue
    }

    validFiles.push({
      ...file,
      name: file.name,
      size: file.size,
      type: file.type,
      progress: undefined // For upload progress tracking
    })
  }
  
  if (validFiles.length) {
    if (props.multiple) {
      files.value = [...files.value, ...validFiles]
    } else {
      files.value = [validFiles[0]]
    }
    
    emit('upload', validFiles)
  }
}

const removeFile = (index) => {
  const removedFile = files.value[index]
  files.value.splice(index, 1)
  emit('remove', removedFile)
}

const isFileTypeAllowed = (file) => {
  if (props.accept === '*/*') return true
  
  const acceptTypes = props.accept.split(',').map(type => type.trim())
  
  return acceptTypes.some(acceptType => {
    if (acceptType.endsWith('/*')) {
      const baseType = acceptType.slice(0, -2)
      return file.type.startsWith(baseType)
    }
    return file.type === acceptType
  })
}

const isImage = (file) => {
  return file.type?.startsWith('image/')
}

const getImagePreview = (file) => {
  if (file instanceof File) {
    return URL.createObjectURL(file)
  }
  return file.url || file.preview || ''
}

const formatFileSize = (bytes) => {
  if (bytes === 0) return '0 Bytes'
  
  const k = 1024
  const sizes = ['Bytes', 'KB', 'MB', 'GB']
  const i = Math.floor(Math.log(bytes) / Math.log(k))
  
  return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i]
}

// Method to update file progress (for external use)
const updateFileProgress = (fileName, progress) => {
  const fileIndex = files.value.findIndex(f => f.name === fileName)
  if (fileIndex !== -1) {
    files.value[fileIndex].progress = progress
  }
}

defineExpose({
  updateFileProgress
})
</script>