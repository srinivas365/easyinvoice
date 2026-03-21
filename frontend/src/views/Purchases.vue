<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="mb-6 flex justify-between items-center">
      <h1 class="text-3xl font-bold text-gray-900">Purchase Invoices</h1>
      <button
        @click="openAddModal()"
        class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium"
      >
        Add Purchase Invoice
      </button>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else-if="invoices.length === 0" class="text-center py-12 text-gray-500">
      No purchase invoices found. Create one by selecting a distributor and adding items.
    </div>

    <div v-else class="bg-white shadow overflow-hidden sm:rounded-md">
      <table class="min-w-full divide-y divide-gray-200">
        <thead class="bg-gray-50">
          <tr>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Distributor</th>
            <th class="px-4 py-3 text-left text-xs font-medium text-gray-500 uppercase">Date</th>
            <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Total</th>
            <th class="px-4 py-3 text-right text-xs font-medium text-gray-500 uppercase">Actions</th>
          </tr>
        </thead>
        <tbody class="bg-white divide-y divide-gray-200">
          <tr v-for="inv in invoices" :key="inv.id">
            <td class="px-4 py-3 text-sm text-gray-900">{{ inv.distributor?.name || '—' }}</td>
            <td class="px-4 py-3 text-sm text-gray-500">{{ formatDate(inv.purchase_date) }}</td>
            <td class="px-4 py-3 text-sm text-gray-900 text-right font-medium">₹{{ (inv.total_amount || 0).toFixed(2) }}</td>
            <td class="px-4 py-3 text-right space-x-2">
              <button
                @click="openEditModal(inv.id)"
                class="text-blue-600 hover:text-blue-800 text-sm font-medium"
              >
                Edit
              </button>
              <button
                @click="confirmDelete(inv)"
                class="text-red-600 hover:text-red-800 text-sm font-medium"
              >
                Delete
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Add / Edit Invoice Modal -->
    <div
      v-if="showFormModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click.self="showFormModal = false"
    >
      <div class="relative top-4 mx-auto p-8 border w-full max-w-7xl shadow-lg rounded-md bg-white mb-8">
        <h3 class="text-lg font-bold mb-4">{{ editingId ? 'Edit' : 'Add' }} Purchase Invoice</h3>

        <div class="grid grid-cols-1 md:grid-cols-2 gap-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700">Distributor *</label>
            <select
              v-model="form.distributor_id"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            >
              <option value="">Select distributor</option>
              <option v-for="d in distributors" :key="d.id" :value="d.id">{{ d.name }}</option>
            </select>
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700">Purchase Date *</label>
            <input
              v-model="form.purchase_date"
              type="date"
              required
              class="mt-1 block w-full border border-gray-300 rounded-md shadow-sm py-2 px-3 focus:outline-none focus:ring-blue-500 focus:border-blue-500"
            />
          </div>
        </div>

        <!-- Add item form -->
        <div class="border rounded-xl p-6 mb-6 bg-gray-50 min-h-[280px]">
          <h4 class="text-base font-semibold text-gray-800 mb-4">{{ editingItemIndex >= 0 ? 'Edit' : 'Add' }} Item</h4>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-6 gap-x-4 gap-y-4">
            <div class="col-span-2">
              <label class="block text-sm font-medium text-gray-700 mb-1">Product *</label>
              <div class="relative" ref="productDropdownRef">
                <input
                  v-model="productSearchQuery"
                  type="text"
                  autocomplete="off"
                  placeholder="Search by name, MFG, HSN..."
                  class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]"
                  @focus="productDropdownOpen = true"
                  @input="productDropdownOpen = true"
                />
                <div
                  v-show="productDropdownOpen"
                  class="absolute z-10 mt-1 w-full bg-white border border-gray-300 rounded-md shadow-lg max-h-56 overflow-auto"
                >
                  <button
                    v-for="p in filteredProducts"
                    :key="p.id"
                    type="button"
                    class="w-full text-left px-3 py-2 text-sm hover:bg-blue-50 focus:bg-blue-50 focus:outline-none"
                    @click="selectProduct(p)"
                  >
                    {{ p.name }}{{ p.mfg ? ' (' + p.mfg + ')' : '' }}
                    <span v-if="p.hsn_code || p.pack" class="text-gray-500 text-xs ml-1">— {{ [p.hsn_code, p.pack].filter(Boolean).join(' · ') }}</span>
                  </button>
                  <p v-if="filteredProducts.length === 0" class="px-3 py-2 text-sm text-gray-500">No products match.</p>
                </div>
              </div>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">MFG</label>
              <input :value="selectedProductForForm?.mfg ?? ''" type="text" readonly class="block w-full border border-gray-200 rounded-md py-2 px-3 text-base bg-gray-50 text-gray-700 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Pack</label>
              <input :value="selectedProductForForm?.pack ?? ''" type="text" readonly class="block w-full border border-gray-200 rounded-md py-2 px-3 text-base bg-gray-50 text-gray-700 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">HSN Code</label>
              <input :value="selectedProductForForm?.hsn_code ?? ''" type="text" readonly class="block w-full border border-gray-200 rounded-md py-2 px-3 text-base bg-gray-50 text-gray-700 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Batch No</label>
              <input v-model="itemForm.batch_no" type="text" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Expiry</label>
              <input v-model="itemForm.expiry" type="date" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Quantity *</label>
              <input v-model.number="itemForm.quantity" type="number" step="any" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Free (qty)</label>
              <input v-model.number="itemForm.free" type="number" step="any" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" placeholder="Extra units at no cost" />
              <p class="mt-0.5 text-xs text-gray-500">Additional quantity received free of cost (not charged)</p>
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">MRP</label>
              <input v-model.number="itemForm.mrp" type="number" step="0.01" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Rate *</label>
              <input v-model.number="itemForm.rate" type="number" step="0.01" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Discount %</label>
              <input v-model.number="itemForm.discount_percent" type="number" step="0.01" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">GST %</label>
              <input v-model.number="itemForm.gst" type="number" step="0.01" min="0" class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base focus:ring-blue-500 focus:border-blue-500 h-[42px]" />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-1">Amount</label>
              <input :value="computedItemAmount.toFixed(2)" type="text" readonly class="block w-full border border-gray-300 rounded-md py-2 px-3 text-base bg-gray-100 text-gray-700 h-[42px]" />
            </div>
          </div>
          <div class="mt-5 flex gap-3">
            <button
              v-if="editingItemIndex >= 0"
              type="button"
              @click="updateItemInList()"
              class="bg-amber-600 hover:bg-amber-700 text-white px-4 py-2 rounded-md text-sm font-medium"
            >
              Update Item
            </button>
            <button
              v-else
              type="button"
              @click="addItemToList()"
              class="bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium"
            >
              Add Item
            </button>
            <button
              v-if="editingItemIndex >= 0"
              type="button"
              @click="cancelEditItem()"
              class="bg-gray-400 hover:bg-gray-500 text-white px-4 py-2 rounded-md text-sm font-medium"
            >
              Cancel
            </button>
          </div>
        </div>

        <!-- Items list -->
        <div class="mb-4">
          <h4 class="font-medium text-gray-800 mb-2">Items ({{ form.items.length }})</h4>
          <div v-if="form.items.length === 0" class="text-sm text-gray-500 py-2">No items added yet.</div>
          <div v-else class="overflow-x-auto">
            <table class="min-w-full divide-y divide-gray-200 text-sm">
              <thead class="bg-gray-100">
                <tr>
                  <th class="px-2 py-2 text-left">Product</th>
                  <th class="px-2 py-2 text-left">Batch</th>
                  <th class="px-2 py-2 text-right">Qty</th>
                  <th class="px-2 py-2 text-right">Rate</th>
                  <th class="px-2 py-2 text-right">Amount</th>
                  <th class="px-2 py-2 text-right w-24">Actions</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                <tr v-for="(it, idx) in form.items" :key="idx">
                  <td class="px-2 py-1">{{ productNameForItem(it) }}</td>
                  <td class="px-2 py-1">{{ it.batch_no }}</td>
                  <td class="px-2 py-1 text-right">{{ it.quantity }}</td>
                  <td class="px-2 py-1 text-right">₹{{ (it.rate || 0).toFixed(2) }}</td>
                  <td class="px-2 py-1 text-right">₹{{ (it.amount || 0).toFixed(2) }}</td>
                  <td class="px-2 py-1 text-right">
                    <button type="button" @click="editItemAtIndex(idx)" class="text-blue-600 hover:underline mr-2">Edit</button>
                    <button type="button" @click="removeItemAtIndex(idx)" class="text-red-600 hover:underline">Delete</button>
                  </td>
                </tr>
              </tbody>
            </table>
          </div>
        </div>

        <!-- Invoice-level discount and total -->
        <div class="flex flex-wrap items-center gap-6 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700">Additional discount (invoice level)</label>
            <input
              v-model.number="form.additional_discount"
              type="number"
              step="0.01"
              min="0"
              class="mt-1 w-40 border border-gray-300 rounded-md shadow-sm py-2 px-3"
            />
          </div>
          <div class="text-lg font-semibold text-gray-900">
            Total: ₹{{ invoiceTotal.toFixed(2) }}
          </div>
        </div>

        <div class="flex justify-end gap-3">
          <button
            type="button"
            @click="showFormModal = false"
            class="bg-gray-300 hover:bg-gray-400 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
          >
            Cancel
          </button>
          <button
            type="button"
            @click="submitForm()"
            :disabled="saving || form.items.length === 0"
            class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
          >
            {{ saving ? 'Saving...' : (editingId ? 'Update' : 'Save') }}
          </button>
        </div>
      </div>
    </div>

    <!-- Delete confirm -->
    <ConfirmDialog
      v-model:show="showConfirm"
      :title="confirmTitle"
      :message="confirmMessage"
      @confirm="handleConfirmDelete"
      @cancel="showConfirm = false"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useToast } from 'vue-toastification'
import { useConfirm } from '../composables/useConfirm'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { purchaseInvoiceService } from '../services/purchaseInvoiceService'
import { distributorService } from '../services/distributorService'
import { productService } from '../services/productService'

const toast = useToast()
const { showConfirm, confirmMessage, confirmTitle, confirm } = useConfirm()

const invoices = ref([])
const distributors = ref([])
const products = ref([])
const loading = ref(true)
const showFormModal = ref(false)
const saving = ref(false)
const editingId = ref(null)
const editingItemIndex = ref(-1)
const deleteTarget = ref(null)
const productSearchQuery = ref('')
const productDropdownOpen = ref(false)
const productDropdownRef = ref(null)

const form = ref({
  distributor_id: '',
  purchase_date: '',
  additional_discount: 0,
  items: []
})

const emptyItemForm = () => ({
  product_id: '',
  batch_no: '',
  expiry: '',
  quantity: 0,
  free: 0,
  mrp: 0,
  rate: 0,
  amount: 0,
  gst: 0,
  discount_percent: 0
})

const itemForm = ref(emptyItemForm())

const invoiceTotal = computed(() => {
  const sum = form.value.items.reduce((s, it) => s + (Number(it.amount) || 0), 0)
  const disc = Number(form.value.additional_discount) || 0
  return Math.max(0, sum - disc)
})

// Amount: first apply discount % on (rate × quantity), then apply GST % on that discounted amount only
const computedItemAmount = computed(() => {
  const q = Number(itemForm.value.quantity) || 0
  const r = Number(itemForm.value.rate) || 0
  const disc = Number(itemForm.value.discount_percent) || 0
  const gst = Number(itemForm.value.gst) || 0
  const base = q * r
  const amountAfterDiscount = base * (1 - disc / 100)  // amount after removing discount %
  const amountWithGst = amountAfterDiscount * (1 + gst / 100)  // GST % on amount after discount
  return amountWithGst
})

function productDisplayName(p) {
  if (!p) return ''
  return p.name + (p.mfg ? ' (' + p.mfg + ')' : '')
}

function productNameForItem(it) {
  if (it.product?.name) return it.product.name
  if (it.product_name) return it.product_name
  const p = products.value.find(x => x.id === it.product_id)
  return p ? p.name : '—'
}

const selectedProductForForm = computed(() => {
  const id = itemForm.value.product_id
  if (!id) return null
  const p = products.value.find(x => x.id === id)
  if (p) return p
  if (editingItemIndex.value >= 0) {
    const it = form.value.items[editingItemIndex.value]
    return it.product || null
  }
  return null
})

const filteredProducts = computed(() => {
  const q = (productSearchQuery.value || '').trim().toLowerCase()
  const list = products.value || []
  if (!q) return list
  return list.filter(p => {
    const name = (p.name || '').toLowerCase()
    const mfg = (p.mfg || '').toLowerCase()
    const hsn = (p.hsn_code || '').toLowerCase()
    const pack = (p.pack || '').toLowerCase()
    return name.includes(q) || mfg.includes(q) || hsn.includes(q) || pack.includes(q)
  })
})

function selectProduct(p) {
  itemForm.value.product_id = p.id
  productSearchQuery.value = productDisplayName(p)
  productDropdownOpen.value = false
}

function addItemToList() {
  if (!itemForm.value.product_id) {
    toast.error('Select a product')
    return
  }
  const name = productNameForItem({ product_id: itemForm.value.product_id })
  form.value.items.push({
    product_id: Number(itemForm.value.product_id),
    product_name: name,
    batch_no: itemForm.value.batch_no,
    expiry: itemForm.value.expiry || null,
    quantity: Number(itemForm.value.quantity) || 0,
    free: Number(itemForm.value.free) || 0,
    mrp: Number(itemForm.value.mrp) || 0,
    rate: Number(itemForm.value.rate) || 0,
    amount: computedItemAmount.value,
    gst: Number(itemForm.value.gst) || 0,
    discount_percent: Number(itemForm.value.discount_percent) || 0
  })
  itemForm.value = emptyItemForm()
  productSearchQuery.value = ''
  productDropdownOpen.value = false
  toast.success('Item added')
}

function editItemAtIndex(idx) {
  const it = form.value.items[idx]
  itemForm.value = {
    product_id: it.product_id || '',
    batch_no: it.batch_no || '',
    expiry: it.expiry ? (typeof it.expiry === 'string' ? it.expiry.slice(0, 10) : it.expiry) : '',
    quantity: it.quantity || 0,
    free: it.free || 0,
    mrp: it.mrp || 0,
    rate: it.rate || 0,
    amount: it.amount || 0,
    gst: it.gst || 0,
    discount_percent: it.discount_percent || 0
  }
  const p = products.value.find(x => x.id === it.product_id) || it.product
  productSearchQuery.value = p ? productDisplayName(p) : ''
  editingItemIndex.value = idx
}

function updateItemInList() {
  if (editingItemIndex.value < 0) return
  if (!itemForm.value.product_id) {
    toast.error('Select a product')
    return
  }
  const it = form.value.items[editingItemIndex.value]
  const name = productNameForItem({ product_id: itemForm.value.product_id })
  it.product_id = Number(itemForm.value.product_id)
  it.product_name = name
  it.batch_no = itemForm.value.batch_no
  it.expiry = itemForm.value.expiry || null
  it.quantity = Number(itemForm.value.quantity) || 0
  it.free = Number(itemForm.value.free) || 0
  it.mrp = Number(itemForm.value.mrp) || 0
  it.rate = Number(itemForm.value.rate) || 0
  it.amount = computedItemAmount.value
  it.gst = Number(itemForm.value.gst) || 0
  it.discount_percent = Number(itemForm.value.discount_percent) || 0
  itemForm.value = emptyItemForm()
  productSearchQuery.value = ''
  productDropdownOpen.value = false
  editingItemIndex.value = -1
  toast.success('Item updated')
}

function cancelEditItem() {
  itemForm.value = emptyItemForm()
  productSearchQuery.value = ''
  productDropdownOpen.value = false
  editingItemIndex.value = -1
}

function removeItemAtIndex(idx) {
  form.value.items.splice(idx, 1)
  if (editingItemIndex.value === idx) {
    editingItemIndex.value = -1
    itemForm.value = emptyItemForm()
  } else if (editingItemIndex.value > idx) {
    editingItemIndex.value--
  }
}

function formatDate(d) {
  if (!d) return '—'
  const x = new Date(d)
  return x.toLocaleDateString()
}

async function openAddModal() {
  editingId.value = null
  form.value = {
    distributor_id: '',
    purchase_date: new Date().toISOString().slice(0, 10),
    additional_discount: 0,
    items: []
  }
  itemForm.value = emptyItemForm()
  productSearchQuery.value = ''
  productDropdownOpen.value = false
  editingItemIndex.value = -1
  await loadProducts()
  showFormModal.value = true
}

async function openEditModal(id) {
  try {
    await loadProducts()
    const inv = await purchaseInvoiceService.getById(id)
    editingId.value = id
    form.value = {
      distributor_id: inv.distributor_id,
      purchase_date: inv.purchase_date ? inv.purchase_date.slice(0, 10) : '',
      additional_discount: inv.additional_discount || 0,
      items: (inv.items || []).map(it => ({
        product_id: it.product_id,
        product_name: it.product?.name || '—',
        product: it.product,
        batch_no: it.batch_no || '',
        expiry: it.expiry ? (typeof it.expiry === 'string' ? it.expiry.slice(0, 10) : it.expiry) : '',
        quantity: it.quantity || 0,
        free: it.free || 0,
        mrp: it.mrp || 0,
        rate: it.rate || 0,
        amount: it.amount || 0,
        gst: it.gst || 0,
        discount_percent: it.discount_percent || 0
      }))
    }
    itemForm.value = emptyItemForm()
    productSearchQuery.value = ''
    productDropdownOpen.value = false
    editingItemIndex.value = -1
    showFormModal.value = true
  } catch (e) {
    toast.error('Failed to load invoice: ' + (e.response?.data?.error || e.message))
  }
}

async function submitForm() {
  if (form.value.items.length === 0) {
    toast.error('Add at least one item')
    return
  }
  if (!form.value.distributor_id || !form.value.purchase_date) {
    toast.error('Select distributor and date')
    return
  }
  const payload = {
    distributor_id: Number(form.value.distributor_id),
    purchase_date: form.value.purchase_date,
    additional_discount: Number(form.value.additional_discount) || 0,
    items: form.value.items.map(it => ({
      product_id: Number(it.product_id),
      batch_no: it.batch_no || '',
      expiry: it.expiry || null,
      quantity: Number(it.quantity) || 0,
      free: Number(it.free) || 0,
      mrp: Number(it.mrp) || 0,
      rate: Number(it.rate) || 0,
      amount: Number(it.amount) || 0,
      gst: Number(it.gst) || 0,
      discount_percent: Number(it.discount_percent) || 0
    }))
  }
  try {
    saving.value = true
    if (editingId.value) {
      await purchaseInvoiceService.update(editingId.value, payload)
      toast.success('Purchase invoice updated')
    } else {
      await purchaseInvoiceService.create(payload)
      toast.success('Purchase invoice created')
    }
    showFormModal.value = false
    loadInvoices()
  } catch (e) {
    toast.error(e.response?.data?.error || e.message || 'Failed to save')
  } finally {
    saving.value = false
  }
}

function confirmDelete(inv) {
  deleteTarget.value = inv
  confirmTitle.value = 'Delete Purchase Invoice'
  confirmMessage.value = `Delete invoice from ${inv.distributor?.name || 'distributor'} dated ${formatDate(inv.purchase_date)}?`
  showConfirm.value = true
}

async function handleConfirmDelete() {
  if (!deleteTarget.value) return
  try {
    await purchaseInvoiceService.delete(deleteTarget.value.id)
    toast.success('Purchase invoice deleted')
    loadInvoices()
  } catch (e) {
    toast.error(e.response?.data?.error || e.message || 'Failed to delete')
  }
  deleteTarget.value = null
  showConfirm.value = false
}

async function loadInvoices() {
  try {
    loading.value = true
    invoices.value = await purchaseInvoiceService.getAll()
  } catch (e) {
    toast.error('Failed to load invoices')
    invoices.value = []
  } finally {
    loading.value = false
  }
}

async function loadDistributors() {
  try {
    distributors.value = await distributorService.getAll()
  } catch (e) {
    distributors.value = []
  }
}

async function loadProducts() {
  try {
    products.value = await productService.getAll()
  } catch (e) {
    products.value = []
  }
}

function handleClickOutside(e) {
  if (productDropdownRef.value && !productDropdownRef.value.contains(e.target)) {
    productDropdownOpen.value = false
  }
}

onMounted(() => {
  loadInvoices()
  loadDistributors()
  document.addEventListener('click', handleClickOutside)
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>
