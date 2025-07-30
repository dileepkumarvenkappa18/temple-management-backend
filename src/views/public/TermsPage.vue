<template>
  <div class="min-h-screen bg-gray-50">
    <!-- Header -->
    <header class="bg-white shadow-sm border-b border-gray-200">
      <div class="max-w-4xl mx-auto px-4 py-3">
        <div class="flex items-center justify-between">
          <router-link to="/" class="inline-flex items-center">
          <div class="flex items-center space-x-3">
            <div class="w-10 h-10 bg-gradient-to-br from-indigo-600 to-indigo-800 rounded-xl flex items-center justify-center">
              <svg class="w-6 h-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M12 2L2 7v10c0 5.55 3.84 9.739 9 11 5.16-1.261 9-5.45 9-11V7l-10-5z"/>
              </svg>
            </div>
            <span class="text-xl font-bold text-gray-900">TempleSaaS</span>
          </div>
          </router-link>
          <router-link 
            to="/" 
            class="text-sm font-medium text-indigo-600 hover:text-indigo-500 transition-colors duration-200"
          >
            Back to Home
          </router-link>
        </div>
      </div>
    </header>

    <!-- Main Content -->
    <main class="max-w-4xl mx-auto px-4 py-12">
      <div class="bg-white rounded-2xl shadow-lg border border-gray-200 p-8 lg:p-12">
        <!-- Title -->
        <div class="mb-8">
          <h1 class="text-3xl font-bold text-gray-900 mb-4">Terms of Service</h1>
          <p class="text-gray-600">Last updated: {{ lastUpdated }}</p>
        </div>

        <!-- Content -->
        <div class="prose prose-indigo max-w-none">
          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">1. Acceptance of Terms</h2>
            <p class="text-gray-700 mb-4">
              By accessing and using TempleSaaS ("Service"), you accept and agree to be bound by the terms and provision of this agreement.
            </p>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">2. Description of Service</h2>
            <p class="text-gray-700 mb-4">
              TempleSaaS is a multi-tenant SaaS platform designed to help temples manage their activities including:
            </p>
            <ul class="list-disc pl-6 text-gray-700 mb-4 space-y-2">
              <li>Devotee management and profile creation</li>
              <li>Seva bookings and scheduling</li>
              <li>Donation tracking and receipt generation</li>
              <li>Event management and communication</li>
              <li>Volunteer coordination</li>
            </ul>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">3. User Accounts and Registration</h2>
            <p class="text-gray-700 mb-4">
              To access certain features of the Service, you may be required to register for an account. You agree to:
            </p>
            <ul class="list-disc pl-6 text-gray-700 mb-4 space-y-2">
              <li>Provide accurate, current, and complete information during registration</li>
              <li>Maintain and promptly update your account information</li>
              <li>Maintain the security of your password and account</li>
              <li>Accept responsibility for all activities under your account</li>
            </ul>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">4. User Roles and Responsibilities</h2>
            <div class="space-y-4">
              <div class="bg-indigo-50 rounded-lg p-4">
                <h3 class="font-medium text-indigo-900 mb-2">Temple Administrators (Tenants)</h3>
                <p class="text-indigo-800 text-sm">
                  Responsible for temple registration, approval process, and managing temple operations including devotees, sevas, and events.
                </p>
              </div>
              <div class="bg-green-50 rounded-lg p-4">
                <h3 class="font-medium text-green-900 mb-2">Devotees</h3>
                <p class="text-green-800 text-sm">
                  Can join approved temples, complete profiles, book sevas, make donations, and participate in temple events.
                </p>
              </div>
              <div class="bg-purple-50 rounded-lg p-4">
                <h3 class="font-medium text-purple-900 mb-2">Volunteers</h3>
                <p class="text-purple-800 text-sm">
                  Can join approved temples, receive assignments, and support temple operations under admin guidance.
                </p>
              </div>
            </div>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">5. Temple Approval Process</h2>
            <p class="text-gray-700 mb-4">
              Temple registrations are subject to administrative approval. We reserve the right to:
            </p>
            <ul class="list-disc pl-6 text-gray-700 mb-4 space-y-2">
              <li>Review and verify temple documentation</li>
              <li>Approve or reject temple applications</li>
              <li>Request additional documentation as needed</li>
              <li>Suspend or deactivate temple accounts for policy violations</li>
            </ul>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">6. Privacy and Data Protection</h2>
            <p class="text-gray-700 mb-4">
              Your privacy is important to us. Please review our Privacy Policy to understand how we collect, use, and protect your information.
            </p>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">7. Prohibited Uses</h2>
            <p class="text-gray-700 mb-4">You may not use our Service:</p>
            <ul class="list-disc pl-6 text-gray-700 mb-4 space-y-2">
              <li>For any unlawful purpose or to solicit others to perform unlawful acts</li>
              <li>To violate any international, federal, provincial, or state regulations, rules, laws, or local ordinances</li>
              <li>To infringe upon or violate our intellectual property rights or the intellectual property rights of others</li>
              <li>To harass, abuse, insult, harm, defame, slander, disparage, intimidate, or discriminate</li>
              <li>To submit false or misleading information</li>
            </ul>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">8. Limitation of Liability</h2>
            <p class="text-gray-700 mb-4">
              In no event shall TempleSaaS, nor its directors, employees, partners, agents, suppliers, or affiliates, be liable for any indirect, incidental, special, consequential, or punitive damages, including without limitation, loss of profits, data, use, goodwill, or other intangible losses, resulting from your use of the Service.
            </p>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">9. Termination</h2>
            <p class="text-gray-700 mb-4">
              We may terminate or suspend your account and bar access to the Service immediately, without prior notice or liability, under our sole discretion, for any reason whatsoever and without limitation, including but not limited to a breach of the Terms.
            </p>
          </section>

          <section class="mb-8">
            <h2 class="text-xl font-semibold text-gray-900 mb-4">10. Changes to Terms</h2>
            <p class="text-gray-700 mb-4">
              We reserve the right, at our sole discretion, to modify or replace these Terms at any time. If a revision is material, we will try to provide at least 30 days notice prior to any new terms taking effect.
            </p>
          </section>

          <section>
            <h2 class="text-xl font-semibold text-gray-900 mb-4">11. Contact Information</h2>
            <p class="text-gray-700 mb-4">
              If you have any questions about these Terms of Service, please contact us at:
            </p>
            <div class="bg-gray-50 rounded-lg p-4">
              <p class="text-gray-700">
                <strong>Email:</strong> support@TempleSaaS.com<br>
                <strong>Address:</strong> Temple Management Solutions<br>
                123 Spiritual Avenue, Sacred City, SC 12345
              </p>
            </div>
          </section>
        </div>
      </div>
    </main>

    <!-- Footer -->
    <footer class="bg-white border-t border-gray-200 mt-12">
      <div class="max-w-4xl mx-auto px-4 py-8">
        <div class="flex items-center justify-between">
          <p class="text-sm text-gray-500">
            © 2025 TempleSaaS. All rights reserved.
          </p>
          <div class="flex space-x-6">
            <router-link to="/privacy" class="text-sm text-gray-500 hover:text-indigo-600 transition-colors duration-200">
              Privacy Policy
            </router-link>
            <router-link to="/" class="text-sm text-gray-500 hover:text-indigo-600 transition-colors duration-200">
              Home
            </router-link>
          </div>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { ref } from 'vue'

// Last updated date
const lastUpdated = ref('January 15, 2025')
</script>