<template>
  <div
    v-if="show"
    class="fixed inset-0 z-50 flex items-end justify-center sm:items-start sm:pt-16 sm:pt-24 p-3 sm:p-4 bg-gray-600/50 overflow-y-auto"
    @click.self="handleCancel"
  >
    <div
      class="relative w-full max-w-md rounded-lg border bg-white p-5 shadow-lg mb-[max(0.5rem,env(safe-area-inset-bottom))] sm:mb-0 max-h-[min(90dvh,100%)] overflow-y-auto touch-manipulation"
      @click.stop
    >
      <div class="flex items-center mb-4">
        <div class="mx-auto flex h-12 w-12 items-center justify-center rounded-full bg-yellow-100">
          <svg class="h-6 w-6 text-yellow-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
          </svg>
        </div>
      </div>
      <div class="mb-4 text-center">
        <h3 class="mb-2 text-lg font-medium text-gray-900">{{ title }}</h3>
        <p class="text-sm text-gray-500 break-words">{{ message }}</p>
      </div>
      <div class="flex flex-col-reverse gap-2 sm:flex-row sm:justify-end sm:gap-3">
        <button
          type="button"
          @click="handleCancel"
          class="w-full sm:w-auto rounded-md bg-gray-300 px-4 py-2.5 text-sm font-medium text-gray-800 hover:bg-gray-400 min-h-[44px] sm:min-h-0"
        >
          Cancel
        </button>
        <button
          type="button"
          @click="handleConfirm"
          class="w-full sm:w-auto rounded-md bg-blue-600 px-4 py-2.5 text-sm font-medium text-white hover:bg-blue-700 min-h-[44px] sm:min-h-0"
        >
          Confirm
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps({
  show: {
    type: Boolean,
    default: false
  },
  title: {
    type: String,
    default: 'Confirm Action'
  },
  message: {
    type: String,
    required: true
  }
})

const emit = defineEmits(['confirm', 'cancel', 'update:show'])

const handleConfirm = () => {
  emit('confirm')
  emit('update:show', false)
}

const handleCancel = () => {
  emit('cancel')
  emit('update:show', false)
}
</script>
