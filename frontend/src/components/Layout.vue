<template>
  <div class="min-h-screen min-h-dvh bg-gray-100 flex flex-col">
    <nav class="bg-white shadow-lg shrink-0 safe-top">
      <div class="max-w-7xl mx-auto px-3 sm:px-6 lg:px-8">
        <div class="flex justify-between items-center h-14 sm:h-16 gap-2">
          <div class="flex items-center min-w-0 gap-2">
            <button
              type="button"
              class="sm:hidden shrink-0 inline-flex items-center justify-center rounded-md p-2.5 text-gray-600 hover:bg-gray-100 hover:text-gray-900 focus:outline-none focus:ring-2 focus:ring-blue-500"
              aria-label="Open menu"
              @click="mobileNavOpen = true"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 12h16M4 18h16" />
              </svg>
            </button>
            <h1 class="text-lg sm:text-xl font-bold text-blue-600 truncate">Pharmacy ERP</h1>
          </div>
          <div class="hidden sm:flex sm:items-center sm:space-x-6 lg:space-x-8">
            <router-link
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="border-transparent text-gray-500 hover:border-gray-300 hover:text-gray-700 inline-flex items-center px-1 pt-1 pb-0.5 border-b-2 text-sm font-medium whitespace-nowrap"
              active-class="border-blue-500 text-gray-900"
            >
              {{ link.label }}
            </router-link>
          </div>
          <div class="flex items-center gap-1.5 sm:gap-3 shrink-0">
            <span class="hidden sm:inline text-gray-700 text-xs sm:text-sm truncate max-w-[10rem] md:max-w-none">
              {{ authStore.user?.username }}
              <span class="text-gray-500">({{ authStore.user?.role }})</span>
            </span>
            <button
              type="button"
              @click="handleLogout"
              class="bg-red-500 hover:bg-red-600 text-white px-3 py-2 sm:px-4 rounded-md text-xs sm:text-sm font-medium touch-manipulation"
            >
              Logout
            </button>
          </div>
        </div>
      </div>
    </nav>

    <!-- Mobile slide-over -->
    <Teleport to="body">
      <div
        v-show="mobileNavOpen"
        class="fixed inset-0 z-[100] sm:hidden"
        aria-modal="true"
        role="dialog"
      >
        <div
          class="absolute inset-0 bg-black/40 backdrop-blur-[1px]"
          aria-hidden="true"
          @click="mobileNavOpen = false"
        />
        <div
          class="absolute top-0 right-0 h-full w-[min(18rem,100vw)] bg-white shadow-xl flex flex-col pt-[env(safe-area-inset-top,0px)] pb-[env(safe-area-inset-bottom,0px)]"
        >
          <div class="flex items-center justify-between px-4 py-3 border-b border-gray-200">
            <span class="font-semibold text-gray-900">Menu</span>
            <button
              type="button"
              class="p-2 rounded-md text-gray-500 hover:bg-gray-100"
              aria-label="Close menu"
              @click="mobileNavOpen = false"
            >
              <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>
          <nav class="flex-1 overflow-y-auto px-2 py-3 space-y-1">
            <router-link
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="block rounded-lg px-4 py-3 text-base font-medium text-gray-700 hover:bg-blue-50 hover:text-blue-700 active:bg-blue-100 touch-manipulation"
              active-class="bg-blue-50 text-blue-700"
              @click="mobileNavOpen = false"
            >
              {{ link.label }}
            </router-link>
          </nav>
          <div class="border-t border-gray-200 p-4 text-sm text-gray-600">
            <p class="truncate font-medium text-gray-900">{{ authStore.user?.username }}</p>
            <p class="text-gray-500 capitalize">{{ authStore.user?.role }}</p>
          </div>
        </div>
      </div>
    </Teleport>

    <main class="max-w-7xl mx-auto w-full flex-1 py-4 sm:py-6 px-3 sm:px-6 lg:px-8 pb-[max(1rem,env(safe-area-inset-bottom))]">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useAuthStore } from '../store/auth'

const authStore = useAuthStore()
const router = useRouter()
const route = useRoute()
const mobileNavOpen = ref(false)

const navLinks = [
  { to: '/', label: 'Dashboard' },
  { to: '/medicines', label: 'Medicines' },
  { to: '/invoices', label: 'Invoices' },
  { to: '/purchases', label: 'Purchases' },
  { to: '/settings', label: 'Settings' }
]

watch(
  () => route.path,
  () => {
    mobileNavOpen.value = false
  }
)

const handleLogout = () => {
  authStore.logout()
  router.push('/login')
}
</script>
