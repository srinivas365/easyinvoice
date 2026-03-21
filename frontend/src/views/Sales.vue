<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Sales (POS)</h1>
    </div>

    <div class="grid grid-cols-1 lg:grid-cols-3 gap-6">
      <!-- Medicine Selection -->
      <div class="lg:col-span-2">
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-xl font-semibold mb-4">Select Medicine</h2>
          <input
            v-model="medicineSearch"
            type="text"
            placeholder="Search medicine..."
            class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 mb-4"
            @input="searchMedicines"
          />

          <div v-if="searchLoading" class="text-center py-4">
            <div class="inline-block animate-spin rounded-full h-6 w-6 border-b-2 border-blue-600"></div>
          </div>

          <div v-else-if="availableMedicines.length === 0" class="text-center py-4 text-gray-500">
            No medicines found
          </div>

          <div v-else class="space-y-2 max-h-96 overflow-y-auto">
            <div
              v-for="medicine in availableMedicines"
              :key="medicine.id"
              class="border border-gray-200 rounded-md p-3 hover:bg-gray-50 cursor-pointer"
              @click="addToCart(medicine)"
            >
              <div class="flex justify-between items-center">
                <div>
                  <p class="font-medium">{{ medicine.name }}</p>
                  <p class="text-sm text-gray-500">Stock: {{ medicine.quantity }} | Price: ₹{{ medicine.selling_price }}</p>
                </div>
                <button
                  class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm"
                  @click.stop="addToCart(medicine)"
                >
                  Add
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>

      <!-- Cart -->
      <div class="lg:col-span-1">
        <div class="bg-white shadow rounded-lg p-6">
          <h2 class="text-xl font-semibold mb-4">Cart</h2>

          <div v-if="cart.length === 0" class="text-center py-8 text-gray-500">
            Cart is empty
          </div>

          <div v-else>
            <div class="space-y-3 mb-4 max-h-64 overflow-y-auto">
              <div
                v-for="(item, index) in cart"
                :key="index"
                class="border border-gray-200 rounded-md p-3"
              >
                <div class="flex justify-between items-start mb-2">
                  <div class="flex-1">
                    <p class="font-medium">{{ item.name }}</p>
                    <p class="text-sm text-gray-500">₹{{ item.price }} each</p>
                  </div>
                  <button
                    @click="removeFromCart(index)"
                    class="text-red-600 hover:text-red-800"
                  >
                    ×
                  </button>
                </div>
                <div class="flex items-center space-x-2">
                  <button
                    @click="updateQuantity(index, -1)"
                    class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                    :disabled="item.quantity <= 1"
                  >
                    -
                  </button>
                  <input
                    v-model.number="item.quantity"
                    type="number"
                    min="1"
                    :max="item.maxQuantity"
                    class="w-16 text-center border border-gray-300 rounded py-1"
                    @change="validateQuantity(index)"
                  />
                  <button
                    @click="updateQuantity(index, 1)"
                    class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                    :disabled="item.quantity >= item.maxQuantity"
                  >
                    +
                  </button>
                  <span class="ml-auto font-semibold">₹{{ (item.price * item.quantity).toFixed(2) }}</span>
                </div>
              </div>
            </div>

            <div class="border-t pt-4">
              <div class="flex justify-between text-lg font-semibold mb-4">
                <span>Total:</span>
                <span>₹{{ totalAmount.toFixed(2) }}</span>
              </div>
              <button
                @click="processSale"
                :disabled="processing"
                class="w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md font-medium disabled:opacity-50"
              >
                {{ processing ? 'Processing...' : 'Complete Sale' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Confirm Dialog -->
  <ConfirmDialog
    v-model:show="showConfirm"
    :title="confirmTitle"
    :message="confirmMessage"
    @confirm="handleConfirm"
    @cancel="handleCancel"
  />
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import { useConfirm } from '../composables/useConfirm'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { medicineService } from '../services/medicineService'
import { saleService } from '../services/saleService'

const toast = useToast()
const { showConfirm, confirmMessage, confirmTitle, confirm, handleConfirm, handleCancel } = useConfirm()

const medicineSearch = ref('')
const availableMedicines = ref([])
const searchLoading = ref(false)
const cart = ref([])
const processing = ref(false)

const totalAmount = computed(() => {
  return cart.value.reduce((sum, item) => sum + (item.price * item.quantity), 0)
})

const searchMedicines = async () => {
  if (medicineSearch.value.length < 2) {
    availableMedicines.value = []
    return
  }
  try {
    searchLoading.value = true
    const medicines = await medicineService.getAll(medicineSearch.value)
    availableMedicines.value = medicines.filter(m => m.quantity > 0)
  } catch (error) {
    console.error('Failed to search medicines:', error)
  } finally {
    searchLoading.value = false
  }
}

const addToCart = (medicine) => {
  const existingIndex = cart.value.findIndex(item => item.id === medicine.id)
  if (existingIndex >= 0) {
    if (cart.value[existingIndex].quantity < medicine.quantity) {
      cart.value[existingIndex].quantity++
    }
  } else {
    cart.value.push({
      id: medicine.id,
      name: medicine.name,
      price: medicine.selling_price,
      quantity: 1,
      maxQuantity: medicine.quantity
    })
  }
}

const removeFromCart = (index) => {
  cart.value.splice(index, 1)
}

const updateQuantity = (index, change) => {
  const item = cart.value[index]
  const newQuantity = item.quantity + change
  if (newQuantity >= 1 && newQuantity <= item.maxQuantity) {
    item.quantity = newQuantity
  }
}

const validateQuantity = (index) => {
  const item = cart.value[index]
  if (item.quantity < 1) item.quantity = 1
  if (item.quantity > item.maxQuantity) item.quantity = item.maxQuantity
}

const processSale = async () => {
  if (cart.value.length === 0) return

  try {
    processing.value = true
    const saleData = {
      items: cart.value.map(item => ({
        medicine_id: item.id,
        quantity: item.quantity
      }))
    }

    const sale = await saleService.create(saleData)
    
    // Print bill (bonus feature)
    const shouldPrint = await confirm('Sale completed! Would you like to print the bill?', 'Print Bill')
    if (shouldPrint) {
      printBill(sale)
    }

    // Clear cart
    cart.value = []
    medicineSearch.value = ''
    availableMedicines.value = []
    toast.success('Sale completed successfully!')
  } catch (error) {
    toast.error('Failed to process sale: ' + (error.response?.data?.error || error.message))
  } finally {
    processing.value = false
  }
}

const printBill = (sale) => {
  const printWindow = window.open('', '_blank')
  const billContent = `
    <html>
      <head>
        <title>Bill #${sale.id}</title>
        <style>
          body { font-family: Arial, sans-serif; padding: 20px; }
          h1 { text-align: center; }
          table { width: 100%; border-collapse: collapse; margin: 20px 0; }
          th, td { border: 1px solid #ddd; padding: 8px; text-align: left; }
          th { background-color: #f2f2f2; }
          .total { font-size: 18px; font-weight: bold; text-align: right; }
        </style>
      </head>
      <body>
        <h1>Pharmacy ERP - Bill</h1>
        <p><strong>Bill #:</strong> ${sale.id}</p>
        <p><strong>Date:</strong> ${new Date(sale.created_at).toLocaleString()}</p>
        <table>
          <thead>
            <tr>
              <th>Medicine</th>
              <th>Quantity</th>
              <th>Price</th>
              <th>Total</th>
            </tr>
          </thead>
          <tbody>
            ${sale.items.map(item => `
              <tr>
                <td>${item.medicine.name}</td>
                <td>${item.quantity}</td>
                <td>₹${item.price.toFixed(2)}</td>
                <td>₹${item.total.toFixed(2)}</td>
              </tr>
            `).join('')}
          </tbody>
        </table>
        <div class="total">Total Amount: ₹${sale.total_amount.toFixed(2)}</div>
      </body>
    </html>
  `
  printWindow.document.write(billContent)
  printWindow.document.close()
  printWindow.print()
}

onMounted(() => {
  // Load all medicines initially
  medicineService.getAll().then(medicines => {
    availableMedicines.value = medicines.filter(m => m.quantity > 0)
  })
})
</script>
