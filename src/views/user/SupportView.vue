<template>
  <div class="support-container">
    <BaseCard class="mb-8">
      <template #header>
        <h1 class="text-2xl font-bold text-indigo-800">Support Center</h1>
      </template>
      <div class="p-4">
        <p class="mb-4 text-gray-700">Welcome to the Temple SaaS Support Center. We're here to help you with any questions or issues you may have.</p>
        
        <div class="bg-indigo-50 p-4 rounded-lg mb-6">
          <h2 class="text-lg font-semibold text-indigo-700 mb-2">Contact Support</h2>
          <p class="mb-4">Our support team is available Monday-Friday, 9am-5pm IST</p>
          <div class="flex items-center mb-2">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-600 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
            </svg>
            <span>support@templesaas.com</span>
          </div>
          <div class="flex items-center">
            <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 text-indigo-600 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
            </svg>
            <span>+91 8000 123456</span>
          </div>
        </div>
        
        <div class="mb-8">
          <h2 class="text-lg font-semibold text-indigo-700 mb-4">Submit a Support Ticket</h2>
          <form @submit.prevent="submitSupportTicket">
            <div class="mb-4">
              <label for="subject" class="block text-sm font-medium text-gray-700 mb-1">Subject</label>
              <BaseInput
                id="subject"
                v-model="ticketForm.subject"
                placeholder="Brief description of your issue"
                :error="ticketErrors.subject"
                required
              />
            </div>
            
            <div class="mb-4">
              <label for="category" class="block text-sm font-medium text-gray-700 mb-1">Category</label>
              <BaseSelect
                id="category"
                v-model="ticketForm.category"
                :options="supportCategories"
                :error="ticketErrors.category"
                required
              />
            </div>
            
            <div class="mb-4">
              <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
              <textarea
                id="description"
                v-model="ticketForm.description"
                rows="4"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-indigo-500"
                placeholder="Please provide details about your issue"
                :class="{'border-red-500': ticketErrors.description}"
                required
              ></textarea>
              <p v-if="ticketErrors.description" class="mt-1 text-sm text-red-600">{{ ticketErrors.description }}</p>
            </div>
            
            <div class="mb-4">
              <label class="block text-sm font-medium text-gray-700 mb-1">Attachments (Optional)</label>
              <BaseFileUpload
                v-model="ticketForm.attachments"
                accept="image/*,.pdf,.doc,.docx"
                multiple
              />
            </div>
            
            <div class="flex justify-end">
              <BaseButton type="submit" :loading="loading" variant="primary">
                Submit Ticket
              </BaseButton>
            </div>
          </form>
        </div>
        
        <div>
          <h2 class="text-lg font-semibold text-indigo-700 mb-4">Frequently Asked Questions</h2>
          <div v-for="(faq, index) in faqs" :key="index" class="mb-4 border-b border-gray-200 pb-4">
            <h3 @click="toggleFaq(index)" class="flex justify-between items-center text-md font-medium text-gray-800 cursor-pointer hover:text-indigo-600">
              {{ faq.question }}
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 transform transition-transform" :class="faq.open ? 'rotate-180' : ''" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </h3>
            <div v-if="faq.open" class="mt-2 text-gray-600 transition-all">
              {{ faq.answer }}
            </div>
          </div>
        </div>
      </div>
    </BaseCard>
  </div>
</template>

<script>
import { ref, reactive } from 'vue';
import { useToast } from '@/composables/useToast';

// Import base components
import BaseCard from '@/components/common/BaseCard.vue';
import BaseInput from '@/components/common/BaseInput.vue';
import BaseSelect from '@/components/common/BaseSelect.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseFileUpload from '@/components/common/BaseFileUpload.vue';

export default {
  name: 'SupportView',
  
  components: {
    BaseCard,
    BaseInput,
    BaseSelect,
    BaseButton,
    BaseFileUpload
  },
  
  setup() {
    const { showToast } = useToast();
    const loading = ref(false);
    
    const ticketForm = reactive({
      subject: '',
      category: '',
      description: '',
      attachments: []
    });
    
    const ticketErrors = reactive({
      subject: '',
      category: '',
      description: ''
    });
    
    const supportCategories = [
      { value: 'account', label: 'Account Issues' },
      { value: 'billing', label: 'Billing & Payments' },
      { value: 'temple', label: 'Temple Management' },
      { value: 'devotee', label: 'Devotee Registration' },
      { value: 'seva', label: 'Seva Booking' },
      { value: 'donation', label: 'Donations' },
      { value: 'technical', label: 'Technical Support' },
      { value: 'other', label: 'Other' }
    ];
    
    const faqs = reactive([
      {
        question: 'How do I reset my password?',
        answer: 'You can reset your password by clicking on the "Forgot Password" link on the login page. You will receive an email with instructions to reset your password.',
        open: false
      },
      {
        question: 'How can I add a new temple to my account?',
        answer: 'If you are a Temple Administrator, you can add a new temple by navigating to the "Manage Temples" section in your dashboard and clicking on the "Add Temple" button.',
        open: false
      },
      {
        question: 'How do I book a seva?',
        answer: 'As a devotee, you can book a seva by selecting your temple, navigating to the "Sevas" section, choosing the seva you want to book, and following the booking process.',
        open: false
      },
      {
        question: 'How do I make a donation?',
        answer: 'You can make a donation by selecting your temple, navigating to the "Donations" section, choosing the donation type, and following the payment process.',
        open: false
      },
      {
        question: 'How can I volunteer for temple activities?',
        answer: 'You can volunteer by selecting your temple, navigating to the "Volunteer" section, and registering as a volunteer. Temple administrators will review your application.',
        open: false
      }
    ]);
    
    const toggleFaq = (index) => {
      faqs[index].open = !faqs[index].open;
    };
    
    const validateTicketForm = () => {
      let isValid = true;
      
      // Reset errors
      Object.keys(ticketErrors).forEach(key => {
        ticketErrors[key] = '';
      });
      
      if (!ticketForm.subject.trim()) {
        ticketErrors.subject = 'Subject is required';
        isValid = false;
      }
      
      if (!ticketForm.category) {
        ticketErrors.category = 'Please select a category';
        isValid = false;
      }
      
      if (!ticketForm.description.trim()) {
        ticketErrors.description = 'Description is required';
        isValid = false;
      } else if (ticketForm.description.trim().length < 20) {
        ticketErrors.description = 'Description must be at least 20 characters';
        isValid = false;
      }
      
      return isValid;
    };
    
    const submitSupportTicket = async () => {
      if (!validateTicketForm()) return;
      
      loading.value = true;
      
      try {
        // Simulate API call
        await new Promise(resolve => setTimeout(resolve, 1500));
        
        showToast({
          type: 'success',
          message: 'Support ticket submitted successfully! Our team will contact you soon.',
          duration: 5000
        });
        
        // Reset form
        ticketForm.subject = '';
        ticketForm.category = '';
        ticketForm.description = '';
        ticketForm.attachments = [];
        
      } catch (error) {
        showToast({
          type: 'error',
          message: 'Failed to submit support ticket. Please try again later.',
          duration: 5000
        });
        console.error('Error submitting support ticket:', error);
      } finally {
        loading.value = false;
      }
    };
    
    return {
      loading,
      ticketForm,
      ticketErrors,
      supportCategories,
      faqs,
      toggleFaq,
      submitSupportTicket
    };
  }
};
</script>

<style scoped>
.support-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 1.5rem;
}

@media (max-width: 768px) {
  .support-container {
    padding: 1rem;
  }
}
</style>