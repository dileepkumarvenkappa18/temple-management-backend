import TempleSelection from '@/views/volunteer/TempleSelection.vue'
import VolunteerDashboard from '@/views/volunteer/VolunteerDashboard.vue'
import MyAssignments from '@/views/volunteer/MyAssignments.vue'
import MySchedule from '@/views/volunteer/MySchedule.vue'

export default [
  {
    path: 'temple-selection',
    name: 'VolunteerTempleSelection',
    component: TempleSelection,
    meta: {
      title: 'Select Temple',
      requiresAuth: true,
      role: 'volunteer'
    }
  },
  {
    path: 'dashboard',
    name: 'VolunteerDashboard',
    component: VolunteerDashboard,
    meta: {
      title: 'Volunteer Dashboard',
      requiresAuth: true,
      role: 'volunteer'
    }
  },
  {
    path: 'assignments',
    name: 'MyAssignments',
    component: MyAssignments,
    meta: {
      title: 'My Assignments',
      requiresAuth: true,
      role: 'volunteer'
    }
  },
  {
    path: 'schedule',
    name: 'MySchedule',
    component: MySchedule,
    meta: {
      title: 'My Schedule',
      requiresAuth: true,
      role: 'volunteer'
    }
  }
]