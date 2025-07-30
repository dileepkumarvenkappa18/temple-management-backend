<template>
  <div class="min-h-screen bg-gradient-to-br from-indigo-50 via-white to-indigo-100">
    <!-- Header -->
    <header class="bg-white shadow-md border-b border-indigo-100">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-16">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2L2 7v10c0 5.55 3.84 9.739 9 11 5.16-1.261 9-5.45 9-11V7l-10-5z"/>
              </svg>
            </div>
            <span class="text-xl font-bold text-gray-900">TempleSaaS</span>
          </div>
          
          <router-link 
            to="/" 
            class="inline-flex items-center px-4 py-2 border border-indigo-200 rounded-xl text-sm font-medium text-indigo-700 bg-indigo-50 hover:bg-indigo-100 hover:border-indigo-300 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-indigo-500 focus:ring-offset-2"
          >
            <svg class="w-4 h-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"/>
            </svg>
            Back to Home
          </router-link>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12">
      <div class="bg-white rounded-2xl shadow-lg border border-indigo-100 overflow-hidden">
        <!-- Hero Section -->
        <div class="bg-gradient-to-r from-indigo-600 to-indigo-800 px-8 py-12 text-center">
          <div class="w-16 h-16 bg-white bg-opacity-20 rounded-2xl flex items-center justify-center mx-auto mb-4">
            <svg class="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m5.618-4.016A11.955 11.955 0 0112 2.944a11.955 11.955 0 01-8.618 3.04A12.02 12.02 0 003 9c0 5.591 3.824 10.29 9 11.622 5.176-1.332 9-6.031 9-11.622 0-1.042-.13-2.052-.382-3.016z"/>
            </svg>
          </div>
          <h1 class="text-3xl sm:text-4xl font-bold text-white mb-3">Privacy Policy</h1>
          <p class="text-indigo-100 text-lg max-w-2xl mx-auto">
            Your privacy and spiritual data security are sacred to us
          </p>
          <div class="mt-6 text-indigo-200 text-sm">
            Last updated: {{ lastUpdated }}
          </div>
        </div>

        <!-- Content Sections -->
        <div class="px-8 py-12">
          <!-- Navigation -->
          <nav class="mb-12">
            <div class="flex flex-wrap gap-2">
              <button
                v-for="section in sections"
                :key="section.id"
                @click="scrollToSection(section.id)"
                class="px-4 py-2 text-sm font-medium rounded-xl border transition-all duration-200"
                :class="activeSection === section.id 
                  ? 'bg-indigo-600 text-white border-indigo-600' 
                  : 'text-indigo-600 border-indigo-200 hover:bg-indigo-50'"
              >
                {{ section.title }}
              </button>
            </div>
          </nav>

          <!-- Content -->
          <div class="space-y-12">
            <!-- Information We Collect -->
            <section :id="sections[0].id" class="scroll-mt-8">
              <div class="flex items-start space-x-4 mb-6">
                <div class="w-12 h-12 bg-indigo-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <svg class="w-6 h-6 text-indigo-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z"/>
                  </svg>
                </div>
                <div class="flex-1">
                  <h2 class="text-2xl font-bold text-gray-900 mb-4">Information We Collect</h2>
                  <div class="prose prose-indigo max-w-none">
                    <p class="text-gray-700 mb-4">We collect information to provide better spiritual services and temple management:</p>
                    
                    <div class="grid md:grid-cols-2 gap-6 mb-6">
                      <div class="bg-indigo-50 rounded-xl p-6 border borwhom:der-indigo-100">
                        <h3 class="font-semibold text-indigo-900 mb-3">Personal Information</h3>
                        <ul class="text-sm text-indigo-800 space-y-2">
                          <li>• Full name and contact details</li>
                          <li>• Spiritual information (Gotra, Nakshatra, Rashi)</li>
                          <li>• Family lineage and relationships</li>
                          <li>• Seva preferences and donations</li>
                        </ul>
                      </div>
                      
                      <div class="bg-green-50 rounded-xl p-6 border border-green-100">
                        <h3 class="font-semibold text-green-900 mb-3">Usage Information</h3>
                        <ul class="text-sm text-green-800 space-y-2">
                          <li>• Temple visits and seva bookings</li>
                          <li>• Event participation and RSVPs</li>
                          <li>• Communication preferences</li>
                          <li>• Device and browser information</li>
                        </ul>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- How We Use Information -->
            <section :id="sections[1].id" class="scroll-mt-8">
              <div class="flex items-start space-x-4 mb-6">
                <div class="w-12 h-12 bg-green-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <svg class="w-6 h-6 text-green-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9.663 17h4.673M12 3v1m6.364 1.636l-.707.707M21 12h-1M4 12H3m3.343-5.657l-.707-.707m2.828 9.9a5 5 0 117.072 0l-.548.547A3.374 3.374 0 0014 18.469V19a2 2 0 11-4 0v-.531c0-.895-.356-1.754-.988-2.386l-.548-.547z"/>
                  </svg>
                </div>
                <div class="flex-1">
                  <h2 class="text-2xl font-bold text-gray-900 mb-4">How We Use Your Information</h2>
                  <div class="grid gap-4">
                    <div class="bg-gradient-to-r from-indigo-50 to-purple-50 rounded-xl p-6 border border-indigo-100">
                      <h3 class="font-semibold text-indigo-900 mb-2">🛐 Spiritual Services</h3>
                      <p class="text-indigo-800 text-sm">Facilitate seva bookings, manage donations, and coordinate temple events according to your spiritual preferences.</p>
                    </div>
                    
                    <div class="bg-gradient-to-r from-blue-50 to-indigo-50 rounded-xl p-6 border border-blue-100">
                      <h3 class="font-semibold text-blue-900 mb-2">📧 Communications</h3>
                      <p class="text-blue-800 text-sm">Send you important temple announcements, event notifications, and seva reminders via your preferred communication channels.</p>
                    </div>
                    
                    <div class="bg-gradient-to-r from-green-50 to-teal-50 rounded-xl p-6 border border-green-100">
                      <h3 class="font-semibold text-green-900 mb-2">🔒 Account Security</h3>
                      <p class="text-green-800 text-sm">Protect your account, prevent fraud, and ensure the security of your spiritual and financial data.</p>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- Data Security -->
            <section :id="sections[2].id" class="scroll-mt-8">
              <div class="flex items-start space-x-4 mb-6">
                <div class="w-12 h-12 bg-red-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <svg class="w-6 h-6 text-red-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 15v2m-6 4h12a2 2 0 002-2v-6a2 2 0 00-2-2H6a2 2 0 00-2 2v6a2 2 0 002 2zm10-10V7a4 4 0 00-8 0v4h8z"/>
                  </svg>
                </div>
                <div class="flex-1">
                  <h2 class="text-2xl font-bold text-gray-900 mb-4">Data Security & Storage</h2>
                  <div class="bg-red-50 rounded-xl p-6 border border-red-100 mb-6">
                    <h3 class="font-semibold text-red-900 mb-3">🔐 Security Measures</h3>
                    <div class="grid md:grid-cols-2 gap-4 text-sm text-red-800">
                      <div>• End-to-end encryption for all data</div>
                      <div>• Secure cloud storage with backups</div>
                      <div>• Multi-factor authentication</div>
                      <div>• Regular security audits</div>
                      <div>• Role-based access controls</div>
                      <div>• GDPR compliance measures</div>
                    </div>
                  </div>
                  
                  <div class="bg-yellow-50 rounded-xl p-6 border border-yellow-100">
                    <h3 class="font-semibold text-yellow-900 mb-2">⚠️ Important Note</h3>
                    <p class="text-yellow-800 text-sm">Your spiritual data is treated with the utmost respect and confidentiality. We never share personal spiritual information with third parties without explicit consent.</p>
                  </div>
                </div>
              </div>
            </section>

            <!-- Your Rights -->
            <section :id="sections[3].id" class="scroll-mt-8">
              <div class="flex items-start space-x-4 mb-6">
                <div class="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                  <svg class="w-6 h-6 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 20h5v-2a3 3 0 00-5.356-1.857M17 20H7m10 0v-2c0-.656-.126-1.283-.356-1.857M7 20H2v-2a3 3 0 015.356-1.857M7 20v-2c0-.656.126-1.283.356-1.857m0 0a5.002 5.002 0 019.288 0M15 7a3 3 0 11-6 0 3 3 0 016 0zm6 3a2 2 0 11-4 0 2 2 0 014 0zM7 10a2 2 0 11-4 0 2 2 0 014 0z"/>
                  </svg>
                </div>
                <div class="flex-1">
                  <h2 class="text-2xl font-bold text-gray-900 mb-4">Your Privacy Rights</h2>
                  <div class="grid gap-4">
                    <div class="flex items-start space-x-4 p-4 bg-gray-50 rounded-xl border border-gray-200">
                      <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center flex-shrink-0">
                        <span class="text-indigo-600 font-bold">1</span>
                      </div>
                      <div>
                        <h3 class="font-semibold text-gray-900 mb-1">Access & Download</h3>
                        <p class="text-gray-700 text-sm">Request a copy of all your personal and spiritual data stored in our system.</p>
                      </div>
                    </div>
                    
                    <div class="flex items-start space-x-4 p-4 bg-gray-50 rounded-xl border border-gray-200">
                      <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center flex-shrink-0">
                        <span class="text-indigo-600 font-bold">2</span>
                      </div>
                      <div>
                        <h3 class="font-semibold text-gray-900 mb-1">Correction & Updates</h3>
                        <p class="text-gray-700 text-sm">Update or correct any inaccurate information in your profile at any time.</p>
                      </div>
                    </div>
                    
                    <div class="flex items-start space-x-4 p-4 bg-gray-50 rounded-xl border border-gray-200">
                      <div class="w-8 h-8 bg-indigo-100 rounded-lg flex items-center justify-center flex-shrink-0">
                        <span class="text-indigo-600 font-bold">3</span>
                      </div>
                      <div>
                        <h3 class="font-semibold text-gray-900 mb-1">Account Deletion</h3>
                        <p class="text-gray-700 text-sm">Request complete deletion of your account and associated data from our platform.</p>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </section>

            <!-- Contact Information -->
            <section :id="sections[4].id" class="scroll-mt-8">
              <div class="bg-gradient-to-r from-indigo-600 to-purple-600 rounded-2xl p-8 text-white">
                <div class="text-center mb-8">
                  <div class="w-16 h-16 bg-opacity-20 rounded-2xl flex items-center justify-center mx-auto mb-4">
                    <svg class="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                    </svg>
                  </div>
                  <h2 class="text-2xl font-bold mb-2">Questions About Privacy?</h2>
                  <p class="text-indigo-100">We're here to help with any privacy concerns</p>
                </div>
                
                <div class="grid md:grid-cols-2 gap-6">
                  <div class="bg-opacity-10 rounded-xl p-6 backdrop-blur-sm">
                    <h3 class="font-semibold mb-3 flex items-center">
                      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"/>
                      </svg>
                      Email Support
                    </h3>
                    <p class="text-indigo-100 text-sm">privacy@TempleSaaS.com</p>
                  </div>
                  
                  <div class="bg-opacity-10 rounded-xl p-6 backdrop-blur-sm">
                    <h3 class="font-semibold mb-3 flex items-center">
                      <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z"/>
                      </svg>
                      Response Time
                    </h3>
                    <p class="text-indigo-100 text-sm">Within 48 hours</p>
                  </div>
                </div>
              </div>
            </section>
          </div>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-indigo-100 mt-12">
      <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div class="text-center text-gray-600 text-sm">
          <p>&copy; 2024 TempleSaaS. All rights reserved. | Built with devotion for spiritual communities.</p>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted } from 'vue'

// Data
const lastUpdated = ref('December 15, 2024')
const activeSection = ref('information-collection')

const sections = [
  { id: 'information-collection', title: 'Information We Collect' },
  { id: 'how-we-use', title: 'How We Use Data' },
  { id: 'data-security', title: 'Data Security' },
  { id: 'your-rights', title: 'Your Rights' },
  { id: 'contact-info', title: 'Contact Us' }
]

// Methods
const scrollToSection = (sectionId) => {
  const element = document.getElementById(sectionId)
  if (element) {
    element.scrollIntoView({ 
      behavior: 'smooth',
      block: 'start'
    })
  }
}

const handleScroll = () => {
  const scrollPosition = window.scrollY + 200
  
  for (const section of sections) {
    const element = document.getElementById(section.id)
    if (element) {
      const { offsetTop, offsetHeight } = element
      if (scrollPosition >= offsetTop && scrollPosition < offsetTop + offsetHeight) {
        activeSection.value = section.id
        break
      }
    }
  }
}

// Lifecycle
onMounted(() => {
  window.addEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  window.removeEventListener('scroll', handleScroll)
})
</script>