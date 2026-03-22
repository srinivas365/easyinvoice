<template>
  <div class="px-0 py-4 sm:py-6">
    <div class="mb-4 sm:mb-6">
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Edit Medicine</h1>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else class="bg-white shadow rounded-lg p-4 sm:p-6 w-full max-w-2xl">
      <form @submit.prevent="handleSubmit">
        <div class="grid grid-cols-1 gap-6 sm:grid-cols-2">
          <div>
            <label for="name" class="block text-sm font-medium text-gray-700">Name *</label>
            <input
              id="name"
              v-model="form.name"
              type="text"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="manufacturer" class="block text-sm font-medium text-gray-700">Manufacturer</label>
            <input
              id="manufacturer"
              v-model="form.manufacturer"
              type="text"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="batch_number" class="block text-sm font-medium text-gray-700">Batch Number</label>
            <input
              id="batch_number"
              v-model="form.batch_number"
              type="text"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="expiry_date" class="block text-sm font-medium text-gray-700">Expiry Date</label>
            <input
              id="expiry_date"
              v-model="form.expiry_date"
              type="date"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="purchase_price" class="block text-sm font-medium text-gray-700">Purchase Price *</label>
            <input
              id="purchase_price"
              v-model.number="form.purchase_price"
              type="number"
              step="0.01"
              required
              min="0"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="selling_price" class="block text-sm font-medium text-gray-700">Selling Price *</label>
            <input
              id="selling_price"
              v-model.number="form.selling_price"
              type="number"
              step="0.01"
              required
              min="0"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>

          <div>
            <label for="quantity" class="block text-sm font-medium text-gray-700">Quantity *</label>
            <input
              id="quantity"
              v-model.number="form.quantity"
              type="number"
              required
              min="0"
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <div class="mt-6 flex justify-end space-x-3">
          <router-link
            to="/medicines"
            class="bg-gray-300 hover:bg-gray-400 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
          >
            Cancel
          </router-link>
          <button
            type="submit"
            :disabled="saving"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
          >
            {{ saving ? 'Saving...' : 'Update' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useToast } from 'vue-toastification'
import { medicineService } from '../services/medicineService'

const toast = useToast()

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const saving = ref(false)

const form = ref({
  name: '',
  manufacturer: '',
  batch_number: '',
  expiry_date: '',
  purchase_price: 0,
  selling_price: 0,
  quantity: 0
})

const loadMedicine = async () => {
  try {
    loading.value = true
    const medicine = await medicineService.getById(route.params.id)
    form.value = {
      name: medicine.name,
      manufacturer: medicine.manufacturer || '',
      batch_number: medicine.batch_number || '',
      expiry_date: medicine.expiry_date ? new Date(medicine.expiry_date).toISOString().split('T')[0] : '',
      purchase_price: medicine.purchase_price,
      selling_price: medicine.selling_price,
      quantity: medicine.quantity
    }
  } catch (error) {
    toast.error('Failed to load medicine: ' + (error.response?.data?.error || error.message))
    router.push('/medicines')
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
    const medicineData = {
      ...form.value,
      expiry_date: form.value.expiry_date || null
    }
    await medicineService.update(route.params.id, medicineData)
    toast.success('Medicine updated successfully!')
    router.push('/medicines')
  } catch (error) {
    toast.error('Failed to update medicine: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

onMounted(() => {
  loadMedicine()
})
</script>
