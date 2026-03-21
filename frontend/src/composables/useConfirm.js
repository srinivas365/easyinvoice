import { ref } from 'vue'

export function useConfirm() {
  const showConfirm = ref(false)
  const confirmMessage = ref('')
  const confirmTitle = ref('Confirm Action')
  const confirmResolve = ref(null)

  const confirm = (message, title = 'Confirm Action') => {
    return new Promise((resolve) => {
      confirmMessage.value = message
      confirmTitle.value = title
      showConfirm.value = true
      confirmResolve.value = resolve
    })
  }

  const handleConfirm = () => {
    showConfirm.value = false
    if (confirmResolve.value) {
      confirmResolve.value(true)
      confirmResolve.value = null
    }
  }

  const handleCancel = () => {
    showConfirm.value = false
    if (confirmResolve.value) {
      confirmResolve.value(false)
      confirmResolve.value = null
    }
  }

  return {
    showConfirm,
    confirmMessage,
    confirmTitle,
    confirm,
    handleConfirm,
    handleCancel
  }
}
