// src/router/routes/devotee.js
import TempleSelection from '@/views/devotee/TempleSelection.vue' 
import ProfileCreation from '@/views/devotee/ProfileCreation.vue' 
import DevoteeDashboard from '@/views/devotee/DevoteeDashboard.vue' 
import SevaBooking from '@/views/devotee/SevaBooking.vue' 
import DonationHistory from '@/views/devotee/DonationHistory.vue' 
import MyEvents from '@/views/devotee/MyEvents.vue' 
import ProfileEdit from '@/views/devotee/ProfileEdit.vue' 
import { checkProfileCompleted } from '../guards'  

export default [   
  {     
    path: 'temple-selection',     
    name: 'DevoteeTempleSelection',     
    component: TempleSelection,     
    meta: {       
      title: 'Select Temple',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },   
  // Path with entity ID (for after temple selection)
  {     
    path: '/entity/:id/devotee/profile/create',     
    name: 'DevoteeProfileCreation',     
    component: ProfileCreation,     
    meta: {       
      title: 'Complete Your Profile',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  // Path without entity ID (fallback)
  {     
    path: 'profile/create',     
    name: 'DevoteeProfileCreationNoEntity',     
    component: ProfileCreation,     
    meta: {       
      title: 'Complete Your Profile',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },   
  // Path with entity ID
  {     
    path: '/entity/:id/devotee/profile/edit',     
    name: 'ProfileEdit',     
    component: ProfileEdit,     
    meta: {       
      title: 'Edit Profile',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  // Path without entity ID (fallback)   
  {     
    path: 'profile/edit',     
    name: 'ProfileEditNoEntity',     
    component: ProfileEdit,     
    meta: {       
      title: 'Edit Profile',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  // Path with entity ID   
  {     
    path: '/entity/:id/devotee/dashboard',     
    name: 'DevoteeDashboard',     
    component: DevoteeDashboard,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'Devotee Dashboard',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  // Path without entity ID (fallback)   
  {     
    path: 'dashboard',     
    name: 'DevoteeDashboardNoEntity',     
    component: DevoteeDashboard,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'Devotee Dashboard',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  // Same pattern for other routes   
  {     
    path: '/entity/:id/devotee/seva-booking',     
    name: 'SevaBooking',     
    component: SevaBooking,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'Book Seva',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  {     
    path: 'seva-booking',     
    name: 'SevaBookingNoEntity',     
    component: SevaBooking,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'Book Seva',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  {     
    path: '/entity/:id/devotee/donations',     
    name: 'DonationHistory',     
    component: DonationHistory,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'My Donations',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  {     
    path: 'donations',     
    name: 'DonationHistoryNoEntity',     
    component: DonationHistory,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'My Donations',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  {     
    path: '/entity/:id/devotee/events',     
    name: 'MyEvents',     
    component: MyEvents,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'My Events',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  },
  {
    path: '/entity/:id/devotee/my-seva-bookings',
    name: 'DevoteeMySevaBookings',
    component: () => import('@/views/devotee/MySevaBookings.vue'),
    beforeEnter: [checkProfileCompleted],
    meta: {
      title: 'My Seva Bookings',
      requiresAuth: true,
      role: 'devotee'
    }
  },
  {     
    path: 'events',     
    name: 'MyEventsNoEntity',     
    component: MyEvents,     
    beforeEnter: [checkProfileCompleted],     
    meta: {       
      title: 'My Events',       
      requiresAuth: true,       
      role: 'devotee'     
    }   
  }
]