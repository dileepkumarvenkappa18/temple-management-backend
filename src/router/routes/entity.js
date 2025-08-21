import EntityDashboard from '@/views/entity/EntityDashboard.vue'
import DevoteeManagement from '@/views/entity/DevoteeManagement.vue'
import SevaManagement from '@/views/entity/SevaManagement.vue'
import DonationManagement from '@/views/entity/DonationManagement.vue'
import EventManagement from '@/views/entity/EventManagement.vue'
import CommunicationCenter from '@/views/entity/CommunicationCenter.vue'
import VolunteerManagement from '@/views/entity/VolunteerManagement.vue'
import EntityDetail from '@/views/entity/EntityDetail.vue' // Add this import
import { checkRole, checkTempleApproved } from '../guards'

export default [
  {
    path: '/entity/:id/dashboard',
    name: 'EntityDashboard',
    component: () => import('@/views/entity/EntityDashboard.vue'),
    meta: { 
      requiresAuth: true,
      allowedRoles: ['superadmin', 'standarduser', 'monitoringuser']
    }
  },
  
  {
    path: 'dashboard',
    name: 'EntityDashboard',
    component: EntityDashboard,
    beforeEnter: [(to, from, next) => checkRole(to, from, next, 'tenant')],
    meta: {
      title: 'Temple Dashboard',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  // Add this new route to handle the entities/:id pattern
  {
    path: 'entities/:id',
    name: 'EntityDetail',
    component: EntityDetail,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Entity Details',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'devotees',
    name: 'DevoteeManagement',
    component: DevoteeManagement,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Devotee Management',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'sevas',
    name: 'SevaManagement',
    component: SevaManagement,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Seva Management',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'donations',
    name: 'DonationManagement',
    component: DonationManagement,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Donation Management',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'events',
    name: 'EventManagement',
    component: EventManagement,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Event Management',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'communication',
    name: 'CommunicationCenter',
    component: CommunicationCenter,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Communication Center',
      requiresAuth: true,
      role: 'tenant'
    }
  },
  {
    path: 'volunteers',
    name: 'VolunteerManagement',
    component: VolunteerManagement,
    beforeEnter: [
      (to, from, next) => checkRole(to, from, next, 'tenant'),
      checkTempleApproved
    ],
    meta: {
      title: 'Volunteer Management',
      requiresAuth: true,
      role: 'tenant'
    }
  }
]