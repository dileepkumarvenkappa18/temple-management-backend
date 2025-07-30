// src/composables/useForm.js
import { ref, reactive, computed, watch } from 'vue'
import { useToast } from './useToast'
import { validators } from '../utils/validators'

export const useForm = (initialData = {}, validationRules = {}) => {
  const { error: showError, success: showSuccess } = useToast()
  
  // Form state
  const formData = reactive({ ...initialData })
  const errors = reactive({})
  const touched = reactive({})
  const isSubmitting = ref(false)
  const isValid = ref(false)

  // Computed properties
  const hasErrors = computed(() => Object.keys(errors).length > 0)
  const touchedFields = computed(() => Object.keys(touched))
  const isFormTouched = computed(() => touchedFields.value.length > 0)

  // Validate single field
  const validateField = (fieldName, value = formData[fieldName]) => {
    const rules = validationRules[fieldName]
    if (!rules) return true

    // Clear previous error
    delete errors[fieldName]

    // Run validation rules
    for (const rule of rules) {
      const result = rule(value)
      if (result !== true) {
        errors[fieldName] = result
        return false
      }
    }
    return true
  }

  // Validate all fields
  const validateForm = () => {
    let isFormValid = true
    
    // Clear all errors first
    Object.keys(errors).forEach(key => delete errors[key])
    
    // Validate each field with rules
    Object.keys(validationRules).forEach(fieldName => {
      const fieldValid = validateField(fieldName)
      if (!fieldValid) {
        isFormValid = false
      }
    })
    
    isValid.value = isFormValid
    return isFormValid
  }

  // Set field value and validate
  const setFieldValue = (fieldName, value) => {
    formData[fieldName] = value
    touched[fieldName] = true
    validateField(fieldName, value)
  }

  // Handle field blur
  const handleFieldBlur = (fieldName) => {
    touched[fieldName] = true
    validateField(fieldName)
  }

  // Handle field focus
  const handleFieldFocus = (fieldName) => {
    // Clear error on focus for better UX
    if (errors[fieldName]) {
      delete errors[fieldName]
    }
  }

  // Reset form
  const resetForm = (newData = initialData) => {
    // Reset form data
    Object.keys(formData).forEach(key => delete formData[key])
    Object.assign(formData, { ...newData })
    
    // Clear errors and touched state
    Object.keys(errors).forEach(key => delete errors[key])
    Object.keys(touched).forEach(key => delete touched[key])
    
    isSubmitting.value = false
    isValid.value = false
  }

  // Handle form submission
  const handleSubmit = async (submitFn, options = {}) => {
    const {
      showSuccessToast = true,
      successMessage = 'Form submitted successfully',
      validateBeforeSubmit = true
    } = options

    try {
      isSubmitting.value = true
      
      // Validate form before submission
      if (validateBeforeSubmit && !validateForm()) {
        showError('Please fix the errors before submitting')
        return false
      }

      // Execute submit function
      const result = await submitFn(formData)
      
      if (showSuccessToast) {
        showSuccess(successMessage)
      }
      
      return result
    } catch (error) {
      showError(error.message || 'Form submission failed')
      throw error
    } finally {
      isSubmitting.value = false
    }
  }

  // Get field error message
  const getFieldError = (fieldName) => {
    return errors[fieldName] || null
  }

  // Check if field has error
  const hasFieldError = (fieldName) => {
    return Boolean(errors[fieldName])
  }

  // Check if field is touched
  const isFieldTouched = (fieldName) => {
    return Boolean(touched[fieldName])
  }

  // Get field props for input components
  const getFieldProps = (fieldName) => {
    return {
      modelValue: formData[fieldName],
      error: getFieldError(fieldName),
      touched: isFieldTouched(fieldName),
      'onUpdate:modelValue': (value) => setFieldValue(fieldName, value),
      onBlur: () => handleFieldBlur(fieldName),
      onFocus: () => handleFieldFocus(fieldName)
    }
  }

  // Common validation rules builder
  const buildRules = (fieldConfig) => {
    const rules = []
    
    if (fieldConfig.required) {
      rules.push(validators.required)
    }
    
    if (fieldConfig.email) {
      rules.push(validators.email)
    }
    
    if (fieldConfig.minLength) {
      rules.push(validators.minLength(fieldConfig.minLength))
    }
    
    if (fieldConfig.maxLength) {
      rules.push(validators.maxLength(fieldConfig.maxLength))
    }
    
    if (fieldConfig.pattern) {
      rules.push(validators.pattern(fieldConfig.pattern, fieldConfig.patternMessage))
    }
    
    if (fieldConfig.custom) {
      rules.push(...fieldConfig.custom)
    }
    
    return rules
  }

  // Watch for form data changes to update validity
  watch(
    () => ({ ...formData }),
    () => {
      if (isFormTouched.value) {
        validateForm()
      }
    },
    { deep: true }
  )

  return {
    // Form state
    formData,
    errors,
    touched,
    isSubmitting,
    isValid,
    
    // Computed
    hasErrors,
    isFormTouched,
    touchedFields,
    
    // Methods
    validateField,
    validateForm,
    setFieldValue,
    handleFieldBlur,
    handleFieldFocus,
    handleSubmit,
    resetForm,
    
    // Field helpers
    getFieldError,
    hasFieldError,
    isFieldTouched,
    getFieldProps,
    
    // Utilities
    buildRules
  }
}