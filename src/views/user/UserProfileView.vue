<template>
  <div class="p-6">
    <h1 class="text-2xl font-bold mb-6">My Profile</h1>
    
    <div class="bg-white shadow-md rounded-lg p-6 mb-6">
      <h2 class="text-xl font-semibold mb-4 border-b pb-2">Personal Information</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex items-center">
          <span class="font-medium w-32">Name:</span>
          <span>John Devotee</span>
        </div>
        <div class="flex items-center">
          <span class="font-medium w-32">Email:</span>
          <span>john@example.com</span>
        </div>
        <div class="flex items-center">
          <span class="font-medium w-32">Role:</span>
          <span class="px-2 py-1 bg-indigo-100 text-indigo-800 rounded-full text-sm">
            Devotee
          </span>
        </div>
        <div class="flex items-center">
          <span class="font-medium w-32">Phone:</span>
          <span>+1 (555) 123-4567</span>
        </div>
      </div>
    </div>
    
    <div class="bg-white shadow-md rounded-lg p-6">
      <h2 class="text-xl font-semibold mb-4 border-b pb-2">Account Information</h2>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
        <div class="flex items-center">
          <span class="font-medium w-32">Member Since:</span>
          <span>January 1, 2023</span>
        </div>
        <div class="flex items-center">
          <span class="font-medium w-32">Last Login:</span>
          <span>June 29, 2025</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  name: 'UserProfileView'
}
</script>