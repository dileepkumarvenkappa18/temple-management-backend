<template>
  <div class="w-full max-w-full">
    <!-- Header -->
    <div class="mb-6">
      <AppBreadcrumb :items="breadcrumbItems" />
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Audit Logs</h1>
          <p class="mt-1 text-sm text-gray-600">
            Track and monitor all system activities across your platform
          </p>
        </div>
        <div class="mt-4 sm:mt-0">
          <BaseButton
            variant="primary"
            @click="refreshLogs"
            class="w-full sm:w-auto"
          >
            <template #icon>
              <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
              </svg>
            </template>
            Refresh
          </BaseButton>
        </div>
      </div>
    </div>
    
    <!-- Main Content -->
    <AuditLogsList @select-log="openLogDetail" />
    
    <!-- Detail Drawer -->
    <AuditLogDetail :log-id="selectedLogId" @close="closeLogDetail" />
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue';
import { useAuditLogStore } from '@/stores/auditlog';
import AppBreadcrumb from '@/components/layout/AppBreadcrumb.vue';
import BaseButton from '@/components/common/BaseButton.vue';
import BaseCard from '@/components/common/BaseCard.vue';
import AuditLogsList from '@/components/superadmin/AuditLogsList.vue';
import AuditLogDetail from '@/components/superadmin/AuditLogDetail.vue';

// Store
const auditLogStore = useAuditLogStore();

// Local state
const selectedLogId = ref(null);
const stats = ref({
  totalLogs: 0,
  successLogs: 0,
  failureLogs: 0,
  todayLogs: 0
});

// Computed
const breadcrumbItems = computed(() => [
  { name: 'Dashboard', path: '/superadmin' },
  { name: 'Audit Logs', path: '/superadmin/audit-logs' }
]);

// Methods
function openLogDetail(id) {
  selectedLogId.value = id;
}

function closeLogDetail() {
  selectedLogId.value = null;
}

function refreshLogs() {
  auditLogStore.fetchAuditLogs();
  calculateStats();
}

function calculateStats() {
  // This would ideally come from the API, but we'll calculate from current data for now
  const logs = auditLogStore.logs;
  
  // Count total logs
  stats.value.totalLogs = auditLogStore.total || logs.length;
  
  // Count success and failure logs
  stats.value.successLogs = logs.filter(log => log.status === 'success').length;
  stats.value.failureLogs = logs.filter(log => log.status === 'failure').length;
  
  // Count today's logs
  const today = new Date();
  today.setHours(0, 0, 0, 0);
  
  stats.value.todayLogs = logs.filter(log => {
    const logDate = new Date(log.createdAt);
    return logDate >= today;
  }).length;
}

// Lifecycle hooks
onMounted(() => {
  refreshLogs();
});
</script>

<style scoped>
/* Add hover effects to cards */
.hover\:shadow-md {
  transition: box-shadow 0.3s ease-in-out;
}
</style>