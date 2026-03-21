<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Add Medicine</h1>
    </div>

    <div class="bg-white shadow rounded-lg p-6 max-w-2xl">
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
            :disabled="loading"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
          >
            {{ loading ? 'Saving...' : 'Save' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useToast } from 'vue-toastification'
import { medicineService } from '../services/medicineService'

const toast = useToast()

const router = useRouter()
const loading = ref(false)

const form = ref({
  name: '',
  manufacturer: '',
  batch_number: '',
  expiry_date: '',
  purchase_price: 0,
  selling_price: 0,
  quantity: 0
})

const handleSubmit = async () => {
  try {
    loading.value = true
    const medicineData = {
      ...form.value,
      expiry_date: form.value.expiry_date || null
    }
    await medicineService.create(medicineData)
    toast.success('Medicine created successfully!')
    router.push('/medicines')
  } catch (error) {
    toast.error('Failed to create medicine: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}
</script>
