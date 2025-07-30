<template>
  <div class="bg-white rounded-xl shadow-lg overflow-hidden">
    <!-- Profile Header -->
    <div class="relative bg-gradient-to-r from-indigo-500 to-indigo-600 px-6 py-8 text-white">
      <div class="flex items-start space-x-6">
        <!-- Profile Photo -->
        <div class="relative">
          <div class="w-24 h-24 bg-white/20 rounded-full border-4 border-white/30 flex items-center justify-center overflow-hidden">
            <img 
              v-if="profileData.profilePhoto" 
              :src="profileData.profilePhoto" 
              :alt="profileData.name"
              class="w-full h-full object-cover"
            >
            <svg v-else class="w-12 h-12 text-white/70" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z"></path>
            </svg>
          </div>
          
          <!-- Upload Photo Button (if own profile) -->
          <button 
            v-if="isOwnProfile" 
            @click="showImageUpload = true"
            class="absolute -bottom-1 -right-1 bg-indigo-700 hover:bg-indigo-800 text-white p-2 rounded-full shadow-lg transition-colors duration-200"
          >
            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 9a2 2 0 012-2h.93a2 2 0 001.664-.89l.812-1.22A2 2 0 0110.07 4h3.86a2 2 0 011.664.89l.812 1.22A2 2 0 0018.07 7H19a2 2 0 012 2v9a2 2 0 01-2 2H5a2 2 0 01-2-2V9z"></path>
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 13a3 3 0 11-6 0 3 3 0 016 0z"></path>
            </svg>
          </button>
        </div>

        <!-- Basic Info -->
        <div class="flex-1">
          <div class="flex items-center justify-between">
            <h1 class="text-2xl font-bold">{{ profileData.name }}</h1>
            <button 
              v-if="isOwnProfile"
              @click="toggleEdit"
              class="bg-white/20 hover:bg-white/30 text-white px-4 py-2 rounded-lg transition-colors duration-200 flex items-center space-x-2"
            >
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z"></path>
              </svg>
              <span>Edit Profile</span>
            </button>
          </div>
          
          <div class="mt-2 space-y-1 text-indigo-100">
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
              </svg>
              <span>{{ profileData.email }}</span>
            </div>
            <div class="flex items-center space-x-2">
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z"></path>
              </svg>
              <span>{{ profileData.phone }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Profile Completion Progress (if own profile) -->
      <div v-if="isOwnProfile" class="mt-6">
        <div class="flex items-center justify-between text-sm text-indigo-100 mb-2">
          <span>Profile Completion</span>
          <span>{{ completionPercentage }}%</span>
        </div>
        <div class="w-full bg-white/20 rounded-full h-2">
          <div 
            class="bg-white h-2 rounded-full transition-all duration-300"
            :style="{ width: completionPercentage + '%' }"
          ></div>
        </div>
      </div>
    </div>

    <!-- Profile Details Tabs -->
    <div class="border-b border-gray-200">
      <nav class="flex space-x-8 px-6" aria-label="Tabs">
        <button
          v-for="tab in tabs"
          :key="tab.id"
          @click="activeTab = tab.id"
          :class="[
            'py-4 px-1 border-b-2 font-medium text-sm whitespace-nowrap transition-colors duration-200',
            activeTab === tab.id
              ? 'border-indigo-500 text-indigo-600'
              : 'border-transparent text-gray-500 hover:text-gray-700 hover:border-gray-300'
          ]"
        >
          {{ tab.name }}
        </button>
      </nav>
    </div>

    <!-- Tab Content -->
    <div class="p-6">
      <!-- Personal Details Tab -->
      <div v-if="activeTab === 'personal'" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Date of Birth</label>
            <p class="text-gray-900">{{ formatDate(profileData.dateOfBirth) }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Gender</label>
            <p class="text-gray-900">{{ profileData.gender }}</p>
          </div>
        </div>
        
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Address</label>
          <div class="text-gray-900">
            <p>{{ profileData.address.street }}</p>
            <p>{{ profileData.address.city }}, {{ profileData.address.state }} - {{ profileData.address.pincode }}</p>
          </div>
        </div>
      </div>

      <!-- Spiritual Information Tab -->
      <div v-if="activeTab === 'spiritual'" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Gotra</label>
            <p class="text-gray-900">{{ profileData.gotra || 'Not specified' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Nakshatra</label>
            <p class="text-gray-900">{{ profileData.nakshatra || 'Not specified' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Rashi</label>
            <p class="text-gray-900">{{ profileData.rashi || 'Not specified' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Lagna</label>
            <p class="text-gray-900">{{ profileData.lagna || 'Not specified' }}</p>
          </div>
          <div class="md:col-span-2">
            <label class="block text-sm font-medium text-gray-700 mb-1">Veda Shaka</label>
            <p class="text-gray-900">{{ profileData.vedaShaka || 'Not specified' }}</p>
          </div>
        </div>
      </div>

      <!-- Family Tab -->
      <div v-if="activeTab === 'family'" class="space-y-6">
        <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Father</label>
            <p class="text-gray-900">{{ profileData.father || 'Not specified' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Mother</label>
            <p class="text-gray-900">{{ profileData.mother || 'Not specified' }}</p>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Spouse</label>
            <p class="text-gray-900">{{ profileData.spouse || 'Not specified' }}</p>
          </div>
        </div>
        
        <div v-if="profileData.children && profileData.children.length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-2">Children</label>
          <div class="space-y-2">
            <span 
              v-for="child in profileData.children" 
              :key="child"
              class="inline-block bg-indigo-100 text-indigo-800 text-sm px-3 py-1 rounded-full mr-2"
            >
              {{ child }}
            </span>
          </div>
        </div>
      </div>

      <!-- Seva Preferences Tab -->
      <div v-if="activeTab === 'seva'" class="space-y-6">
        <div v-if="profileData.favoriteSevas && profileData.favoriteSevas.length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-2">Favorite Sevas</label>
          <div class="flex flex-wrap gap-2">
            <span 
              v-for="seva in profileData.favoriteSevas" 
              :key="seva"
              class="inline-block bg-indigo-100 text-indigo-800 text-sm px-3 py-2 rounded-lg"
            >
              {{ seva }}
            </span>
          </div>
        </div>
        
        <div v-if="profileData.donationInterests && profileData.donationInterests.length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-2">Donation Interests</label>
          <div class="flex flex-wrap gap-2">
            <span 
              v-for="interest in profileData.donationInterests" 
              :key="interest"
              class="inline-block bg-green-100 text-green-800 text-sm px-3 py-2 rounded-lg"
            >
              {{ interest }}
            </span>
          </div>
        </div>
      </div>

      <!-- Notes Tab -->
      <div v-if="activeTab === 'notes'" class="space-y-6">
        <div v-if="profileData.notes">
          <label class="block text-sm font-medium text-gray-700 mb-1">Additional Notes</label>
          <p class="text-gray-900 whitespace-pre-line">{{ profileData.notes }}</p>
        </div>
        
        <div v-if="profileData.attachments && profileData.attachments.length > 0">
          <label class="block text-sm font-medium text-gray-700 mb-2">Attachments</label>
          <div class="space-y-2">
            <div 
              v-for="attachment in profileData.attachments" 
              :key="attachment.name"
              class="flex items-center space-x-3 p-3 bg-gray-50 rounded-lg"
            >
              <svg class="w-5 h-5 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"></path>
              </svg>
              <span class="text-sm text-gray-900">{{ attachment.name }}</span>
              <button class="text-indigo-600 hover:text-indigo-500 text-sm">Download</button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Image Upload Modal -->
  <div v-if="showImageUpload" class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50 p-4">
    <div class="bg-white rounded-xl max-w-md w-full p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-semibold text-gray-900">Upload Profile Photo</h3>
        <button @click="showImageUpload = false" class="text-gray-400 hover:text-gray-600">
          <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12"></path>
          </svg>
        </button>
      </div>
      
      <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 text-center">
        <svg class="mx-auto h-12 w-12 text-gray-400" stroke="currentColor" fill="none" viewBox="0 0 48 48">
          <path d="M28 8H12a4 4 0 00-4 4v20m32-12v8m0 0v8a4 4 0 01-4 4H12a4 4 0 01-4-4v-4m32-4l-3.172-3.172a4 4 0 00-5.656 0L28 28M8 32l9.172-9.172a4 4 0 015.656 0L28 28m0 0l4 4m4-24h8m-4-4v8m-12 4h.02" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"></path>
        </svg>
        <div class="mt-4">
          <label class="cursor-pointer">
            <span class="mt-2 block text-sm font-medium text-gray-900">
              Click to upload or drag and drop
            </span>
            <input type="file" class="sr-only" accept="image/*" @change="handleImageUpload">
          </label>
          <p class="mt-1 text-xs text-gray-500">PNG, JPG, GIF up to 10MB</p>
        </div>
      </div>
      
      <div class="mt-6 flex justify-end space-x-3">
        <button 
          @click="showImageUpload = false"
          class="px-4 py-2 text-sm font-medium text-gray-700 bg-gray-100 hover:bg-gray-200 rounded-lg transition-colors duration-200"
        >
          Cancel
        </button>
        <button class="px-4 py-2 text-sm font-medium text-white bg-indigo-600 hover:bg-indigo-700 rounded-lg transition-colors duration-200">
          Upload
        </button>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, computed, onMounted } from 'vue'

export default {
  name: 'DevoteeProfile',
  props: {
    devotee: {
      type: Object,
      default: () => ({})
    },
    isOwnProfile: {
      type: Boolean,
      default: false
    }
  },
  setup(props, { emit }) {
    const isEditing = ref(false)
    const showImageUpload = ref(false)
    const activeTab = ref('personal')
    
    // Sample profile data - in real app, this comes from props or store
    const profileData = ref({
      id: props.devotee.id || '1',
      profilePhoto: props.devotee.profilePhoto || '',
      name: props.devotee.name || 'Rajesh Kumar',
      email: props.devotee.email || 'rajesh@example.com',
      phone: props.devotee.phone || '+91 9876543210',
      dateOfBirth: props.devotee.dateOfBirth || '1985-06-15',
      gender: props.devotee.gender || 'Male',
      address: props.devotee.address || {
        street: '123 Temple Street',
        city: 'Bangalore',
        state: 'Karnataka',
        pincode: '560001'
      },
      // Spiritual Information
      gotra: props.devotee.gotra || 'Bharadwaja',
      nakshatra: props.devotee.nakshatra || 'Rohini',
      rashi: props.devotee.rashi || 'Vrishabha',
      lagna: props.devotee.lagna || 'Simha',
      vedaShaka: props.devotee.vedaShaka || 'Rig Veda',
      // Family Lineage
      father: props.devotee.father || 'Ramesh Kumar',
      mother: props.devotee.mother || 'Lakshmi Devi',
      spouse: props.devotee.spouse || 'Priya Kumar',
      children: props.devotee.children || ['Arjun Kumar', 'Anita Kumar'],
      // Seva Preferences
      favoriteSevas: props.devotee.favoriteSevas || ['Abhisheka', 'Archana', 'Annadanam'],
      donationInterests: props.devotee.donationInterests || ['Temple Maintenance', 'Anna Danam', 'Education Fund'],
      // Notes and attachments
      notes: props.devotee.notes || 'Regular devotee, prefers morning sevas.',
      attachments: props.devotee.attachments || [
        { name: 'Family_Photo.jpg', type: 'image' },
        { name: 'Health_Certificate.pdf', type: 'document' }
      ]
    })

    const tabs = ref([
      { id: 'personal', name: 'Personal Details' },
      { id: 'spiritual', name: 'Spiritual Info' },
      { id: 'family', name: 'Family' },
      { id: 'seva', name: 'Seva Preferences' },
      { id: 'notes', name: 'Notes & Files' }
    ])

    // Calculate profile completion percentage
    const completionPercentage = computed(() => {
      const fields = [
        profileData.value.name,
        profileData.value.email,
        profileData.value.phone,
        profileData.value.dateOfBirth,
        profileData.value.address?.street,
        profileData.value.gotra,
        profileData.value.nakshatra,
        profileData.value.father,
        profileData.value.mother,
        profileData.value.favoriteSevas?.length > 0
      ]
      
      const completedFields = fields.filter(field => field).length
      return Math.round((completedFields / fields.length) * 100)
    })

    const formatDate = (dateString) => {
      if (!dateString) return 'Not specified'
      const date = new Date(dateString)
      return date.toLocaleDateString('en-IN', {
        year: 'numeric',
        month: 'long',
        day: 'numeric'
      })
    }

    const toggleEdit = () => {
      isEditing.value = !isEditing.value
      emit('edit-profile', profileData.value)
    }

    const handleImageUpload = (event) => {
      const file = event.target.files[0]
      if (file) {
        // In real app, upload to server and get URL
        const reader = new FileReader()
        reader.onload = (e) => {
          profileData.value.profilePhoto = e.target.result
          showImageUpload.value = false
          emit('photo-updated', e.target.result)
        }
        reader.readAsDataURL(file)
      }
    }

    return {
      isEditing,
      showImageUpload,
      activeTab,
      tabs,
      profileData,
      completionPercentage,
      formatDate,
      toggleEdit,
      handleImageUpload
    }
  }
}
</script>