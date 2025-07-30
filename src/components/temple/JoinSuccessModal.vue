<template>
  <div
    v-if="show"
    class="fixed inset-0 bg-gray-500 bg-opacity-75 flex items-center justify-center p-4 z-50"
  >
    <div
      class="bg-white rounded-xl shadow-xl max-w-md w-full p-6"
      @click.stop
    >
      <div class="text-center">
        <!-- Success Icon -->
        <div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-green-100 mb-4">
          <svg class="h-10 w-10 text-green-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 13l4 4L19 7"></path>
          </svg>
        </div>

        <!-- Success Message -->
        <h3 class="text-xl font-bold text-indigo-900 mb-2 font-heading">
          Successfully Joined!
        </h3>
        <p class="text-indigo-700 mb-6 font-side">
          You've successfully joined <span class="font-semibold">{{ templeName }}</span>!
          Start building your profile to get the most out of your temple experience.
        </p>

        <!-- Action Buttons -->
        <div class="flex flex-col sm:flex-row gap-3">
          <button
            @click="completeProfile"
            class="flex-1 bg-indigo-600 hover:bg-indigo-700 text-white font-medium py-2 px-4 rounded-lg"
          >
            Complete Profile
          </button>
          <button
            @click="goToDashboard"
            class="flex-1 bg-white hover:bg-gray-50 text-indigo-700 font-medium py-2 px-4 rounded-lg border border-indigo-200"
          >
            Go to Dashboard
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { defineComponent } from 'vue';
import { useRouter } from 'vue-router';

export default defineComponent({
  name: 'JoinSuccessModal',
  props: {
    show: {
      type: Boolean,
      default: false
    },
    templeId: {
      type: [Number, String],
      required: true
    },
    templeName: {
      type: String,
      default: 'Temple'
    }
  },
  emits: ['close'],
  setup(props, { emit }) {
    const router = useRouter();

    const completeProfile = () => {
      // Close the modal
      emit('close');
      
      // Navigate to profile creation page
      router.push('/devotee/profile/create');
    };

    const goToDashboard = () => {
      // Close the modal
      emit('close');
      
      // Navigate to devotee dashboard for the specific temple
      router.push(`/entity/${props.templeId}/devotee/dashboard`);
    };

    return {
      completeProfile,
      goToDashboard
    };
  }
});
</script>