<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header Section -->
    <div class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-6">
        <div class="flex items-center justify-between">
          <div>
            <nav class="flex items-center space-x-2 text-sm text-gray-500 mb-2">
              <router-link to="/tenant/dashboard" class="hover:text-indigo-600">Dashboard</router-link>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
              <router-link to="/tenant/entities" class="hover:text-indigo-600">Manage Temples</router-link>
              <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
              </svg>
              <span class="text-gray-900">Edit Temple</span>
            </nav>
            <h1 class="text-2xl font-bold text-gray-900">Edit Temple</h1>
            <p class="text-gray-600 mt-1">Update temple information</p>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content -->
    <div class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
      <!-- Progress Indicator -->
      <div class="mb-8">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-indigo-600">Step {{ currentStep }} of 3</span>
          <span class="text-sm text-gray-500">{{ Math.round((currentStep / 3) * 100) }}% Complete</span>
        </div>
        <div class="w-full bg-gray-200 rounded-full h-2">
          <div 
            class="bg-indigo-600 h-2 rounded-full transition-all duration-300"
            :style="{ width: `${(currentStep / 3) * 100}%` }"
          ></div>
        </div>
      </div>

      <!-- Loading State -->
      <div v-if="isLoading" class="flex justify-center items-center py-20">
        <div class="animate-spin rounded-full h-12 w-12 border-t-2 border-b-2 border-indigo-600"></div>
        <span class="ml-3 text-gray-600">Loading temple data...</span>
      </div>

      <!-- Error State -->
      <div v-else-if="loadError" class="bg-red-50 border border-red-200 rounded-lg p-6 text-center">
        <svg class="w-12 h-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z"/>
        </svg>
        <h3 class="text-lg font-medium text-red-800 mb-2">Failed to load temple data</h3>
        <p class="text-red-700 mb-4">{{ loadError }}</p>
        <button 
          @click="fetchTempleData" 
          class="inline-flex items-center px-4 py-2 bg-red-600 text-white rounded-lg hover:bg-red-700"
        >
          Try Again
        </button>
      </div>

      <!-- Form Container -->
      <div v-else class="bg-white rounded-2xl shadow-sm border border-gray-200 overflow-hidden">
        <form @submit.prevent="handleSubmit">
          <!-- Step 1: Basic Information -->
          <div v-show="currentStep === 1" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Temple Basic Information</h2>
              <p class="text-gray-600">Update the essential details about your temple</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Temple Name -->
              <div class="md:col-span-2">
                <label for="templeName" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Name *
                </label>
                <input
                  id="templeName"
                  v-model="form.name"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter temple name"
                />
              </div>

              <!-- Deity -->
              <div>
                <label for="deity" class="block text-sm font-medium text-gray-700 mb-2">
                  Main Deity *
                </label>
                <input
                  id="deity"
                  v-model="form.mainDeity"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="e.g., Lord Vishnu, Goddess Lakshmi"
                />
              </div>

              <!-- Temple Type -->
              <div>
                <label for="templeType" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Type *
                </label>
                <select
                  id="templeType"
                  v-model="form.templeType"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                >
                  <option value="">Select temple type</option>
                  <option value="traditional">Traditional</option>
                  <option value="modern">Modern</option>
                  <option value="heritage">Heritage</option>
                  <option value="community">Community</option>
                </select>
              </div>

              <!-- Established Year -->
              <div>
                <label for="establishedYear" class="block text-sm font-medium text-gray-700 mb-2">
                  Established Year
                </label>
                <input
                  id="establishedYear"
                  v-model="form.establishedYear"
                  type="number"
                  min="1800"
                  :max="new Date().getFullYear()"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="YYYY"
                />
              </div>

              <!-- Phone -->
              <div>
                <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
                  Contact Phone *
                </label>
                <input
                  id="phone"
                  v-model="form.phone"
                  type="tel"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="+91 XXXXX XXXXX"
                />
              </div>

              <!-- Email -->
              <div class="md:col-span-2">
                <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Email
                </label>
                <input
                  id="email"
                  v-model="form.email"
                  type="email"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="temple@example.com"
                />
              </div>

              <!-- Description -->
              <div class="md:col-span-2">
                <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
                  Temple Description
                </label>
                <textarea
                  id="description"
                  v-model="form.description"
                  rows="4"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Brief description of the temple, its history, and significance..."
                ></textarea>
              </div>
            </div>
          </div>

          <!-- Step 2: Address Information -->
          <div v-show="currentStep === 2" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Temple Address</h2>
              <p class="text-gray-600">Update the complete address details</p>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Street Address -->
              <div class="md:col-span-2">
                <label for="address" class="block text-sm font-medium text-gray-700 mb-2">
                  Street Address *
                </label>
                <input
                  id="address"
                  v-model="form.streetAddress"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter complete street address"
                />
              </div>

              <!-- City -->
              <div>
                <label for="city" class="block text-sm font-medium text-gray-700 mb-2">
                  City *
                </label>
                <input
                  id="city"
                  v-model="form.city"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter city name"
                />
              </div>

              <!-- State -->
              <div>
                <label for="state" class="block text-sm font-medium text-gray-700 mb-2">
                  State *
                </label>
                <select
                  id="state"
                  v-model="form.state"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                >
                  <option value="">Select state</option>
                  <option value="Karnataka">Karnataka</option>
                  <option value="Tamil Nadu">Tamil Nadu</option>
                  <option value="Andhra Pradesh">Andhra Pradesh</option>
                  <option value="Kerala">Kerala</option>
                  <option value="Maharashtra">Maharashtra</option>
                  <option value="Gujarat">Gujarat</option>
                  <option value="Rajasthan">Rajasthan</option>
                  <option value="Uttar Pradesh">Uttar Pradesh</option>
                  <option value="West Bengal">West Bengal</option>
                  <option value="Odisha">Odisha</option>
                </select>
              </div>

              <!-- District -->
              <div>
                <label for="district" class="block text-sm font-medium text-gray-700 mb-2">
                  District *
                </label>
                <input
                  id="district"
                  v-model="form.district"
                  type="text"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter district name"
                />
              </div>

              <!-- PIN Code -->
              <div>
                <label for="pinCode" class="block text-sm font-medium text-gray-700 mb-2">
                  PIN Code *
                </label>
                <input
                  id="pinCode"
                  v-model="form.pincode"
                  type="text"
                  pattern="[0-9]{6}"
                  required
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Enter 6-digit PIN code"
                />
              </div>

              <!-- Landmark -->
              <div class="md:col-span-2">
                <label for="landmark" class="block text-sm font-medium text-gray-700 mb-2">
                  Landmark (Optional)
                </label>
                <input
                  id="landmark"
                  v-model="form.landmark"
                  type="text"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="Nearby landmark for easy identification"
                />
              </div>

              <!-- Google Maps Link -->
              <div class="md:col-span-2">
                <label for="googleMapsLink" class="block text-sm font-medium text-gray-700 mb-2">
                  Google Maps Link (Optional)
                </label>
                <input
                  id="googleMapsLink"
                  v-model="form.mapLink"
                  type="url"
                  class="w-full px-4 py-3 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500 transition-colors duration-200"
                  placeholder="https://maps.google.com/..."
                />
              </div>
            </div>
          </div>

          <!-- Step 3: Documents Upload -->
          <div v-show="currentStep === 3" class="p-8">
            <div class="mb-6">
              <h2 class="text-xl font-bold text-gray-900 mb-2">Update Documents</h2>
              <p class="text-gray-600">Upload or update temple verification documents</p>
            </div>

            <div class="space-y-6">
              <!-- Registration Certificate -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Registration Certificate *</h3>
                  <p class="text-gray-600 mb-4">Upload temple registration certificate</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200">
                      <span>{{ existingFiles.registration ? 'Replace File' : 'Choose File' }}</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'registration')" />
                    </label>
                  </div>
                  <p v-if="form.documents.registration" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.registration.name }}
                  </p>
                  <p v-else-if="existingFiles.registration" class="text-sm text-gray-600 mt-2">
                    Current: {{ existingFiles.registration }}
                  </p>
                </div>
              </div>

              <!-- Trust Deed -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Trust Deed *</h3>
                  <p class="text-gray-600 mb-4">Upload temple trust deed document</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-indigo-600 text-white px-6 py-2 rounded-lg hover:bg-indigo-700 transition-colors duration-200">
                      <span>{{ existingFiles.trustDeed ? 'Replace File' : 'Choose File' }}</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'trustDeed')" />
                    </label>
                  </div>
                  <p v-if="form.documents.trustDeed" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.trustDeed.name }}
                  </p>
                  <p v-else-if="existingFiles.trustDeed" class="text-sm text-gray-600 mt-2">
                    Current: {{ existingFiles.trustDeed }}
                  </p>
                </div>
              </div>

              <!-- Property Documents -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2H5a2 2 0 00-2-2z"/>
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 5a2 2 0 012-2h4a2 2 0 012 2v2H8V5z"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Property Documents</h3>
                  <p class="text-gray-600 mb-4">Upload property ownership documents (Optional)</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-gray-600 text-white px-6 py-2 rounded-lg hover:bg-gray-700 transition-colors duration-200">
                      <span>{{ existingFiles.property ? 'Replace File' : 'Choose File' }}</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" @change="handleFileUpload($event, 'property')" />
                    </label>
                  </div>
                  <p v-if="form.documents.property" class="text-sm text-green-600 mt-2">
                    ✓ {{ form.documents.property.name }}
                  </p>
                  <p v-else-if="existingFiles.property" class="text-sm text-gray-600 mt-2">
                    Current: {{ existingFiles.property }}
                  </p>
                </div>
              </div>

              <!-- Additional Documents -->
              <div class="border-2 border-dashed border-gray-300 rounded-lg p-6 hover:border-indigo-400 transition-colors duration-200">
                <div class="text-center">
                  <svg class="w-12 h-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15.172 7l-6.586 6.586a2 2 0 102.828 2.828l6.414-6.586a4 4 0 00-5.656-5.656l-6.415 6.585a6 6 0 108.486 8.486L20.5 13"/>
                  </svg>
                  <h3 class="text-lg font-medium text-gray-900 mb-2">Additional Documents</h3>
                  <p class="text-gray-600 mb-4">Upload any additional supporting documents (Optional)</p>
                  <div class="flex justify-center">
                    <label class="cursor-pointer bg-gray-600 text-white px-6 py-2 rounded-lg hover:bg-gray-700 transition-colors duration-200">
                      <span>{{ existingFiles.additional && existingFiles.additional.length > 0 ? 'Add More Files' : 'Choose Files' }}</span>
                      <input type="file" class="hidden" accept=".pdf,.jpg,.jpeg,.png" multiple @change="handleFileUpload($event, 'additional')" />
                    </label>
                  </div>
                  <div v-if="form.documents.additional && form.documents.additional.length > 0" class="mt-2">
                    <p v-for="(file, index) in form.documents.additional" :key="'new-'+index" class="text-sm text-green-600">
                      ✓ {{ file.name }}
                    </p>
                  </div>
                  <div v-if="existingFiles.additional && existingFiles.additional.length > 0" class="mt-2">
                    <p v-for="(file, index) in existingFiles.additional" :key="'existing-'+index" class="text-sm text-gray-600">
                      Current: {{ file }}
                    </p>
                  </div>
                </div>
              </div>
            </div>

            <!-- Status Information -->
            <div class="mt-8 p-4 bg-gray-50 rounded-lg">
              <div class="flex items-start">
                <svg class="w-5 h-5 text-gray-600 mt-0.5 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                </svg>
                <div>
                  <h3 class="text-sm font-medium text-gray-900 mb-1">Current Status: <span :class="statusClass">{{ form.status }}</span></h3>
                  <p v-if="form.rejectionReason" class="text-sm text-red-600 mt-1">
                    <strong>Rejection Reason:</strong> {{ form.rejectionReason }}
                  </p>
                  <p v-if="form.adminNotes" class="text-sm text-gray-700 mt-1">
                    <strong>Admin Notes:</strong> {{ form.adminNotes }}
                  </p>
                </div>
              </div>
            </div>
          </div>

          <!-- Navigation Buttons -->
          <div class="px-8 py-6 bg-gray-50 border-t border-gray-200 flex justify-between">
            <button
              v-if="currentStep > 1"
              type="button"
              @click="previousStep"
              class="inline-flex items-center px-6 py-3 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors duration-200"
            >
              <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7"/>
              </svg>
              Previous
            </button>
            <div v-else></div>

            <div class="flex space-x-3">
              <router-link
                to="/tenant/entities"
                class="inline-flex items-center px-6 py-3 border border-gray-300 rounded-lg text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 transition-colors duration-200"
              >
                Cancel
              </router-link>

              <button
                v-if="currentStep < 3"
                type="button"
                @click="nextStep"
                class="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 transition-colors duration-200"
              >
                Next
                <svg class="w-4 h-4 ml-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7"/>
                </svg>
              </button>

              <button
                v-else
                type="submit"
                :disabled="isSubmitting"
                class="inline-flex items-center px-6 py-3 bg-indigo-600 text-white rounded-lg text-sm font-medium hover:bg-indigo-700 disabled:opacity-50 disabled:cursor-not-allowed transition-colors duration-200"
              >
                <svg v-if="isSubmitting" class="animate-spin -ml-1 mr-2 h-4 w-4 text-white" fill="none" viewBox="0 0 24 24">
                  <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                  <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                </svg>
                {{ isSubmitting ? 'Saving Changes...' : 'Update Temple' }}
              </button>
            </div>
          </div>
        </form>
      </div>

      <!-- Help Section -->
      <div class="mt-8 bg-blue-50 border border-blue-200 rounded-lg p-6">
        <div class="flex items-start">
          <svg class="w-6 h-6 text-blue-600 mt-0.5 mr-3 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
          </svg>
          <div>
            <h3 class="text-sm font-medium text-blue-800 mb-1">Need Help?</h3>
            <p class="text-sm text-blue-700">
              If you need assistance with updating temple information, please contact our support team at 
              <a href="mailto:support@templesaas.com" class="underline">support@templesaas.com</a> 
              or call us at <a href="tel:+91-800-123-4567" class="underline">+91-800-123-4567</a>
            </p>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, reactive, computed, onMounted, defineComponent } from 'vue';
import { useRouter, useRoute } from 'vue-router';
import { useTempleStore } from '@/stores/temple';
import { useToast } from '@/composables/useToast';
import templeService from '@/services/temple.service.js';

export default defineComponent({
 name: 'EditTemple',
 setup() {
   const router = useRouter();
   const route = useRoute();
   const templeStore = useTempleStore();
   const toast = useToast();

   // Temple ID from route params
   const templeId = computed(() => route.params.id);

   // UI state
   const currentStep = ref(1);
   const isSubmitting = ref(false);
   const isLoading = ref(true);
   const loadError = ref(null);

   // Form data - FIXED FIELD NAMES to match backend structure
   const form = reactive({
     name: '',
     mainDeity: '', // Changed from deity to mainDeity
     templeType: '', // Changed from type to templeType
     establishedYear: '',
     phone: '',
     email: '',
     description: '',
     streetAddress: '', // Changed from address to streetAddress
     city: '',
     state: '',
     district: '',
     pincode: '', // Changed from pinCode to pincode
     landmark: '',
     mapLink: '', // Changed from googleMapsLink to mapLink
     documents: {
       registration: null,
       trustDeed: null,
       property: null,
       additional: []
     },
     status: '',
     rejectionReason: '',
     adminNotes: ''
   });

   // Track existing files
   const existingFiles = reactive({
     registration: '',
     trustDeed: '',
     property: '',
     additional: []
   });

   // Helper function to safely show toast messages
   const showToast = (message, type = 'success') => {
     try {
       if (type === 'success' && toast.success) {
         toast.success(message);
       } else if (type === 'error' && toast.error) {
         toast.error(message);
       } else if (toast.addToast) {
         // Alternative API using addToast
         toast.addToast({
           content: message,
           type: type
         });
       } else {
         // Fallback to console and alert for errors
         console.log(`${type.toUpperCase()}: ${message}`);
         if (type === 'error') {
           alert(message);
         }
       }
     } catch (e) {
       console.error('Toast error:', e);
       if (type === 'error') {
         alert(message);
       } else {
         console.log(`${type.toUpperCase()}: ${message}`);
       }
     }
   };

   // Computed properties
   const statusClass = computed(() => {
     const status = form.status?.toUpperCase() || '';
     switch(status) {
       case 'APPROVED':
         return 'text-green-600 font-medium';
       case 'PENDING':
         return 'text-yellow-600 font-medium';
       case 'REJECTED':
         return 'text-red-600 font-medium';
       default:
         return 'text-gray-600';
     }
   });

   // Fetch temple data
   const fetchTempleData = async () => {
     isLoading.value = true;
     loadError.value = null;
     
     try {
       // Get the temple data from API
       const response = await templeService.getTempleById(templeId.value);
       const templeData = response.data || response;
       
       // Populate form with temple data - USING CORRECT FIELD NAMES
       form.name = templeData.name || '';
       form.mainDeity = templeData.main_deity || '';
       form.templeType = templeData.temple_type || '';
       form.establishedYear = templeData.established_year || '';
       form.phone = templeData.phone || '';
       form.email = templeData.email || '';
       form.description = templeData.description || '';
       
       // Address details
       form.streetAddress = templeData.street_address || '';
       form.city = templeData.city || '';
       form.state = templeData.state || '';
       form.district = templeData.district || '';
       form.pincode = templeData.pincode || '';
       form.landmark = templeData.landmark || '';
       form.mapLink = templeData.map_link || '';
       
       // Status information
       form.status = templeData.status || 'PENDING';
       form.rejectionReason = templeData.rejection_reason || '';
       form.adminNotes = templeData.admin_notes || '';
       
       // Document information
       if (templeData.registration_cert_url) {
         existingFiles.registration = templeData.registration_cert_url.split('/').pop() || 'registration.pdf';
       }
       
       if (templeData.trust_deed_url) {
         existingFiles.trustDeed = templeData.trust_deed_url.split('/').pop() || 'trust_deed.pdf';
       }
       
       if (templeData.property_docs_url) {
         existingFiles.property = templeData.property_docs_url.split('/').pop() || 'property.pdf';
       }
       
       if (templeData.additional_docs_urls) {
         try {
           // Try to parse if it's a JSON string
           const additionalDocs = JSON.parse(templeData.additional_docs_urls);
           if (Array.isArray(additionalDocs)) {
             existingFiles.additional = additionalDocs.map(url => url.split('/').pop());
           } else {
             existingFiles.additional = Object.values(additionalDocs).map(url => url.split('/').pop());
           }
         } catch (e) {
           // If not JSON, handle as a comma-separated string
           existingFiles.additional = templeData.additional_docs_urls.split(',').map(url => url.trim().split('/').pop());
         }
       }
       
     } catch (error) {
       console.error('Failed to fetch temple data:', error);
       loadError.value = typeof error === 'object' ? (error.message || 'Failed to load temple data') : 'Failed to load temple data';
     } finally {
       isLoading.value = false;
     }
   };

   // Navigation methods
   const nextStep = () => {
     if (validateCurrentStep()) {
       currentStep.value++;
     }
   };

   const previousStep = () => {
     currentStep.value--;
   };

   // Validation
   const validateCurrentStep = () => {
     if (currentStep.value === 1) {
       const requiredFields = [
         { field: 'name', label: 'Temple Name' },
         { field: 'mainDeity', label: 'Main Deity' },
         { field: 'templeType', label: 'Temple Type' },
         { field: 'phone', label: 'Phone Number' }
       ];
       
       const missingFields = requiredFields.filter(({ field }) => !form[field]?.trim());
       
       if (missingFields.length > 0) {
         showToast(`Please fill in all required fields: ${missingFields.map(f => f.label).join(', ')}`, 'error');
         return false;
       }
     } else if (currentStep.value === 2) {
       const requiredFields = [
         { field: 'streetAddress', label: 'Street Address' },
         { field: 'city', label: 'City' },
         { field: 'state', label: 'State' },
         { field: 'district', label: 'District' },
         { field: 'pincode', label: 'PIN Code' }
       ];
       
       const missingFields = requiredFields.filter(({ field }) => !form[field]?.trim());
       
       if (missingFields.length > 0) {
         showToast(`Please fill in all required address fields: ${missingFields.map(f => f.label).join(', ')}`, 'error');
         return false;
       }
       
       // Validate PIN code format
       if (form.pincode && !/^\d{6}$/.test(form.pincode)) {
         showToast('PIN Code must be a 6-digit number', 'error');
         return false;
       }
     }
     
     return true;
   };

   // Handle file uploads
   const handleFileUpload = (event, type) => {
     const files = Array.from(event.target.files);
     
     if (type === 'additional') {
       form.documents.additional = files;
     } else {
       form.documents[type] = files[0];
     }
   };

   // Form submission - COMPLETELY REWRITTEN to properly format data
  const handleSubmit = async () => {
  if (!validateCurrentStep()) return;

  try {
    isSubmitting.value = true;

    // Ensure the established year is a valid number
    let establishedYear = null;
    if (form.establishedYear && !isNaN(parseInt(form.establishedYear))) {
      establishedYear = parseInt(form.establishedYear);
    }

    // Prepare data in the format expected by the backend
    const payload = {
      id: parseInt(templeId.value),
      name: form.name.trim(),
      main_deity: form.mainDeity.trim(),
      temple_type: form.templeType,
      established_year: establishedYear,
      phone: form.phone.trim(),
      email: form.email.trim(),
      description: form.description.trim(),
      street_address: form.streetAddress.trim(),
      city: form.city.trim(),
      state: form.state.trim(),
      district: form.district.trim(),
      pincode: form.pincode.trim(),
      landmark: form.landmark.trim(),
      map_link: form.mapLink.trim()
    };

    console.log('Submitting update with data:', payload);
    
    await templeService.updateTemple(templeId.value, payload);

    showToast('Temple information updated successfully!');
    
    // Use the same routing as create temple - redirect to tenant dashboard
    router.push('/tenant/dashboard');
  } catch (error) {
    console.error('Failed to update temple:', error);
    const errorMessage = typeof error === 'object' ? (error.message || 'Failed to update temple information. Please try again.') : 'Failed to update temple information. Please try again.';
    showToast(errorMessage, 'error');
  } finally {
    isSubmitting.value = false;
  }
};

   // Initialize component
   onMounted(() => {
     fetchTempleData();
   });

   return {
     // UI state
     currentStep,
     isSubmitting,
     isLoading,
     loadError,
     
     // Form data
     form,
     existingFiles,
     
     // Computed
     statusClass,
     
     // Methods
     fetchTempleData,
     nextStep,
     previousStep,
     handleFileUpload,
     handleSubmit
   };
 }
});
</script>