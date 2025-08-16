<template>
  <button 
    v-if="shouldShowAssignButton" 
    @click="navigateToAssignPage"
    class="text-indigo-600 hover:text-indigo-900 ml-2"
  >
    Assign
  </button>
</template>

<script>
export default {
  name: 'UserAssignButton',
  props: {
    user: {
      type: Object,
      required: true
    }
  },
  computed: {
    shouldShowAssignButton() {
      if (!this.user || !this.user.role) return false;
      
      // Get the role name, handling different possible formats
      let roleName = '';
      if (typeof this.user.role === 'string') {
        roleName = this.user.role.toLowerCase();
      } else if (typeof this.user.role === 'object') {
        roleName = (this.user.role.role_name || this.user.role.roleName || '').toLowerCase();
      }
      
      // Check if the role matches any of the allowed roles (with or without spaces)
      return ['standarduser', 'standard user', 'Standard User', 'monitoringuser', 'monitoring user', 'Monitor User'].includes(roleName);
    }
  },
  methods: {
    navigateToAssignPage() {
      this.$router.push({
        name: 'AssignTenants',
        params: { userId: this.user.id || this.user.ID }
      });
    }
  }
}
</script>