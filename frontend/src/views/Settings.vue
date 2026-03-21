<template>
  <div class="px-4 py-6 sm:px-0">
    <div class="mb-6">
      <h1 class="text-3xl font-bold text-gray-900">Settings</h1>
      <p class="mt-2 text-sm text-gray-600">Configure company details for invoices</p>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      <!-- Settings Section -->
      <div class="bg-white shadow rounded-lg p-6">
        <div class="mb-6">
          <h2 class="text-2xl font-bold text-gray-900">Company Settings</h2>
          <p class="mt-1 text-sm text-gray-600">Configure company details for invoices</p>
        </div>
        <form @submit.prevent="handleSubmit">
          <div class="space-y-6">
            <div>
              <label for="company_name" class="block text-sm font-medium text-gray-700 mb-1">
                Company Name *
              </label>
              <input
                id="company_name"
                v-model="form.company_name"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter company name"
              />
            </div>

            <div>
              <label for="address" class="block text-sm font-medium text-gray-700 mb-1">
                Address
              </label>
              <textarea
                id="address"
                v-model="form.address"
                rows="3"
                class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter shop address"
              ></textarea>
            </div>

            <div>
              <label for="phone_number" class="block text-sm font-medium text-gray-700 mb-1">
                Phone Number
              </label>
              <input
                id="phone_number"
                v-model="form.phone_number"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter shop phone number"
              />
            </div>

            <div>
              <label for="license_id" class="block text-sm font-medium text-gray-700 mb-1">
                D.L.No *
              </label>
              <input
                id="license_id"
                v-model="form.license_id"
                type="text"
                required
                class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter D.L.No"
              />
            </div>

            <div>
              <label for="gstin_number" class="block text-sm font-medium text-gray-700 mb-1">
                GSTIN Number
              </label>
              <input
                id="gstin_number"
                v-model="form.gstin_number"
                type="text"
                class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                placeholder="Enter GSTIN number (optional)"
              />
              <p class="mt-1 text-sm text-gray-500">
                Leave empty if not applicable. GSTIN will only appear on invoices if provided.
              </p>
            </div>
          </div>

          <div class="mt-6 flex justify-end space-x-3">
            <button
              type="button"
              @click="loadSettings"
              class="bg-gray-300 hover:bg-gray-400 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
            >
              Reset
            </button>
            <button
              type="submit"
              :disabled="saving"
              class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
            >
              {{ saving ? 'Saving...' : 'Save Settings' }}
            </button>
          </div>
        </form>
      </div>

      <!-- Distributors Section -->
      <div class="bg-white shadow rounded-lg p-6">
      <div class="mb-6 flex justify-between items-center">
        <div>
          <h2 class="text-2xl font-bold text-gray-900">Distributors</h2>
          <p class="mt-1 text-sm text-gray-600">Manage your distributor contacts</p>
        </div>
        <button
          @click="openAddModal"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium flex items-center gap-2"
        >
          <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          Add Distributor
        </button>
      </div>

      <div v-if="distributorsLoading" class="text-center py-12">
        <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
      </div>

      <div v-else-if="distributors.length === 0" class="text-center py-12 text-gray-500">
        No distributors found. Click "Add Distributor" to get started.
      </div>

      <div v-else class="overflow-hidden sm:rounded-md">
        <ul class="divide-y divide-gray-200">
          <li v-for="distributor in distributors" :key="distributor.id">
            <div class="px-4 py-4 sm:px-6">
              <div class="flex items-center justify-between">
                <div class="flex-1">
                  <div class="flex items-center">
                    <p class="text-lg font-medium text-gray-900">{{ distributor.name }}</p>
                  </div>
                  <div class="mt-2 sm:flex sm:justify-between">
                    <div class="sm:flex flex-wrap gap-4">
                      <p v-if="distributor.contact_person" class="flex items-center text-sm text-gray-500">
                        Contact: {{ distributor.contact_person }}
                      </p>
                      <p v-if="distributor.phone" class="flex items-center text-sm text-gray-500">
                        Phone: {{ distributor.phone }}
                      </p>
                      <p v-if="distributor.email" class="flex items-center text-sm text-gray-500">
                        Email: {{ distributor.email }}
                      </p>
                    </div>
                  </div>
                  <div v-if="distributor.address" class="mt-2">
                    <p class="text-sm text-gray-500">Address: {{ distributor.address }}</p>
                  </div>
                </div>
                <div class="ml-4 flex-shrink-0 flex gap-2">
                  <button
                    @click="openEditModal(distributor)"
                    class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    Edit
                  </button>
                  <button
                    @click="deleteDistributor(distributor.id)"
                    class="bg-red-600 hover:bg-red-700 text-white px-3 py-1.5 rounded-md text-sm font-medium flex items-center gap-1"
                  >
                    <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </button>
                </div>
              </div>
            </div>
          </li>
        </ul>
      </div>
      </div>
    </div>

    <!-- Add/Edit Distributor Modal -->
    <div v-if="showDistributorModal" class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50" @click.self="closeDistributorModal">
      <div class="relative top-20 mx-auto p-5 border w-96 shadow-lg rounded-md bg-white">
        <div class="mt-3">
          <div class="flex justify-between items-center mb-4">
            <h3 class="text-lg font-medium text-gray-900">
              {{ editingDistributor ? 'Edit Distributor' : 'Add Distributor' }}
            </h3>
            <button @click="closeDistributorModal" class="text-gray-400 hover:text-gray-500">
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          </div>

          <form @submit.prevent="handleDistributorSubmit">
            <div class="space-y-4">
              <div>
                <label for="dist_name" class="block text-sm font-medium text-gray-700 mb-1">
                  Name *
                </label>
                <input
                  id="dist_name"
                  v-model="distributorForm.name"
                  type="text"
                  required
                  class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter distributor name"
                />
              </div>

              <div>
                <label for="dist_contact" class="block text-sm font-medium text-gray-700 mb-1">
                  Contact Person
                </label>
                <input
                  id="dist_contact"
                  v-model="distributorForm.contact_person"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter contact person name"
                />
              </div>

              <div>
                <label for="dist_phone" class="block text-sm font-medium text-gray-700 mb-1">
                  Phone
                </label>
                <input
                  id="dist_phone"
                  v-model="distributorForm.phone"
                  type="text"
                  class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter phone number"
                />
              </div>

              <div>
                <label for="dist_email" class="block text-sm font-medium text-gray-700 mb-1">
                  Email
                </label>
                <input
                  id="dist_email"
                  v-model="distributorForm.email"
                  type="email"
                  class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter email address"
                />
              </div>

              <div>
                <label for="dist_address" class="block text-sm font-medium text-gray-700 mb-1">
                  Address
                </label>
                <textarea
                  id="dist_address"
                  v-model="distributorForm.address"
                  rows="3"
                  class="w-full px-4 py-2 border border-gray-300 rounded-md shadow-sm focus:outline-none focus:ring-blue-500 focus:border-blue-500"
                  placeholder="Enter address"
                ></textarea>
              </div>
            </div>

            <div class="mt-6 flex justify-end space-x-3">
              <button
                type="button"
                @click="closeDistributorModal"
                class="bg-gray-300 hover:bg-gray-400 text-gray-800 px-4 py-2 rounded-md text-sm font-medium"
              >
                Cancel
              </button>
              <button
                type="submit"
                :disabled="distributorSaving"
                class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium disabled:opacity-50"
              >
                {{ distributorSaving ? 'Saving...' : (editingDistributor ? 'Update' : 'Add') }}
              </button>
            </div>
          </form>
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
import { ref, onMounted } from 'vue'
import { useToast } from 'vue-toastification'
import { useConfirm } from '../composables/useConfirm'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { settingsService } from '../services/settingsService'
import { distributorService } from '../services/distributorService'

const toast = useToast()
const { showConfirm, confirmMessage, confirmTitle, confirm, handleConfirm, handleCancel } = useConfirm()

const loading = ref(true)
const saving = ref(false)

const form = ref({
  company_name: '',
  address: '',
  phone_number: '',
  license_id: '',
  gstin_number: ''
})

// Distributors
const distributors = ref([])
const distributorsLoading = ref(false)
const showDistributorModal = ref(false)
const editingDistributor = ref(null)
const distributorSaving = ref(false)

const distributorForm = ref({
  name: '',
  contact_person: '',
  phone: '',
  email: '',
  address: ''
})

const loadSettings = async () => {
  try {
    loading.value = true
    const settings = await settingsService.get()
    form.value = {
      company_name: settings.company_name || '',
      address: settings.address || '',
      phone_number: settings.phone_number || '',
      license_id: settings.license_id || '',
      gstin_number: settings.gstin_number || ''
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
    toast.error('Failed to load settings: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const handleSubmit = async () => {
  try {
    saving.value = true
    await settingsService.update(form.value)
    toast.success('Settings saved successfully!')
  } catch (error) {
    toast.error('Failed to save settings: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

const loadDistributors = async () => {
  try {
    distributorsLoading.value = true
    distributors.value = await distributorService.getAll()
  } catch (error) {
    console.error('Failed to load distributors:', error)
    toast.error('Failed to load distributors: ' + (error.response?.data?.error || error.message))
  } finally {
    distributorsLoading.value = false
  }
}

const openAddModal = () => {
  editingDistributor.value = null
  distributorForm.value = {
    name: '',
    contact_person: '',
    phone: '',
    email: '',
    address: ''
  }
  showDistributorModal.value = true
}

const openEditModal = (distributor) => {
  editingDistributor.value = distributor
  distributorForm.value = {
    name: distributor.name || '',
    contact_person: distributor.contact_person || '',
    phone: distributor.phone || '',
    email: distributor.email || '',
    address: distributor.address || ''
  }
  showDistributorModal.value = true
}

const closeDistributorModal = () => {
  showDistributorModal.value = false
  editingDistributor.value = null
  distributorForm.value = {
    name: '',
    contact_person: '',
    phone: '',
    email: '',
    address: ''
  }
}

const handleDistributorSubmit = async () => {
  try {
    distributorSaving.value = true
    if (editingDistributor.value) {
      await distributorService.update(editingDistributor.value.id, distributorForm.value)
      toast.success('Distributor updated successfully!')
    } else {
      await distributorService.create(distributorForm.value)
      toast.success('Distributor added successfully!')
    }
    closeDistributorModal()
    loadDistributors()
  } catch (error) {
    toast.error('Failed to save distributor: ' + (error.response?.data?.error || error.message))
  } finally {
    distributorSaving.value = false
  }
}

const deleteDistributor = async (id) => {
  const confirmed = await confirm('Are you sure you want to delete this distributor?', 'Delete Distributor')
  if (!confirmed) {
    return
  }

  try {
    await distributorService.delete(id)
    toast.success('Distributor deleted successfully!')
    loadDistributors()
  } catch (error) {
    toast.error('Failed to delete distributor: ' + (error.response?.data?.error || error.message))
  }
}

onMounted(() => {
  loadSettings()
  loadDistributors()
})
</script>
