<template>
  <div class="px-0 py-4 sm:px-0 sm:py-6">
    <div class="mb-4 sm:mb-6">
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900">Medicines</h1>
      <p class="mt-1 text-sm text-gray-500">Products from your catalog. Add products when adding purchase invoice items.</p>
    </div>

    <div class="mb-4">
      <input
        v-model="searchQuery"
        type="text"
        placeholder="Search by name, manufacturer or HSN..."
        class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        @input="loadMedicines"
      />
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else-if="!(medicines && medicines.length)" class="text-center py-12 text-gray-500">
      No products yet. Add products when adding items in a purchase invoice.
    </div>

    <div v-else class="bg-white shadow overflow-hidden sm:rounded-md">
      <ul class="divide-y divide-gray-200">
        <li v-for="medicine in (medicines || [])" :key="medicine.id">
          <div class="px-4 py-4 sm:px-6">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center">
                  <p class="text-lg font-medium text-gray-900">{{ medicine.name }}</p>
                  <span
                    :class="[
                      'ml-2 px-2 py-1 text-xs font-semibold rounded-full',
                      (medicine.quantity || 0) <= 10 ? 'bg-red-100 text-red-800' : 'bg-green-100 text-green-800'
                    ]"
                    title="Total stock (paid qty + free qty)"
                  >
                    {{ (medicine.quantity || 0) }} units
                  </span>
                </div>
                <div class="mt-2 sm:flex sm:justify-between">
                  <div class="sm:flex flex-wrap gap-x-6 gap-y-1">
                    <p class="text-sm text-gray-500">MFG: {{ medicine.mfg || '—' }}</p>
                    <p class="text-sm text-gray-500">HSN: {{ medicine.hsn_code || '—' }}</p>
                    <p class="text-sm text-gray-500">Pack: {{ medicine.pack || '—' }}</p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { productService } from '../services/productService'

const medicines = ref([])
const searchQuery = ref('')
const loading = ref(true)

const loadMedicines = async () => {
  try {
    loading.value = true
    const data = await productService.getAllWithQuantity(searchQuery.value)
    medicines.value = Array.isArray(data) ? data : []
  } catch (error) {
    console.error('Failed to load medicines:', error)
    medicines.value = []
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadMedicines()
})
</script>
