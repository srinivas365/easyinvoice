<template>
  <div class="px-0 py-4 sm:py-6">
    <div class="mb-4 sm:mb-6 flex flex-col sm:flex-row justify-between items-stretch sm:items-center gap-3 sm:gap-4">
      <h1 class="text-2xl sm:text-3xl font-bold text-gray-900 shrink-0">Invoices</h1>
      <div class="flex flex-wrap items-stretch sm:items-center gap-2 w-full sm:w-auto min-w-0">
        <input
          v-model="searchQuery"
          type="text"
          placeholder="Search by Bill no., Customer Name, or Phone..."
          class="min-w-0 w-full sm:w-auto sm:min-w-[12rem] flex-1 px-3 sm:px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          @input="filterInvoices"
        />
        <input
          v-model="dateFilter"
          type="date"
          class="px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 text-sm"
          @change="filterInvoices"
        />
        <button
          v-if="dateFilter"
          @click="dateFilter = ''; filterInvoices()"
          class="px-3 py-2 bg-gray-200 hover:bg-gray-300 text-gray-700 rounded-md text-sm font-medium transition-colors"
          title="Clear date filter"
        >
          <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
        </button>
        <button
          @click="showArchive = !showArchive"
          :class="[
            'px-4 py-2 rounded-md text-sm font-medium transition-colors',
            showArchive 
              ? 'bg-gray-600 hover:bg-gray-700 text-white' 
              : 'bg-gray-200 hover:bg-gray-300 text-gray-700'
          ]"
        >
          {{ showArchive ? 'Active Invoices' : 'Archive' }}
        </button>
        <button
          @click="openAddModal"
          class="bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-md text-sm font-medium shadow-md hover:shadow-lg transition-shadow"
        >
          <span class="flex items-center">
            <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Create Invoice
          </span>
        </button>
      </div>
    </div>

    <div v-if="loading" class="text-center py-12">
      <div class="inline-block animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
    </div>

    <div v-else-if="groupedInvoices.length === 0" class="text-center py-12 text-gray-500">
      No invoices found
    </div>

    <div v-else class="space-y-4 sm:space-y-6 max-h-[calc(100dvh-12rem)] sm:max-h-[calc(100vh-250px)] overflow-y-auto -mx-1 px-1 sm:mx-0 sm:px-0">
      <!-- Date-wise grouped invoices -->
      <div v-for="group in paginatedGroups" :key="group.date" class="bg-white shadow rounded-lg overflow-hidden">
        <!-- Date Header -->
        <div class="bg-gray-50 px-4 py-3 border-b border-gray-200 sticky top-0 z-10">
          <h3 class="text-lg font-semibold text-gray-900">
            {{ formatDateHeader(group.date) }}
            <span class="ml-2 text-sm font-normal text-gray-500">({{ group.invoices.length }} invoice{{ group.invoices.length !== 1 ? 's' : '' }})</span>
          </h3>
        </div>
        
        <!-- Invoices for this date -->
        <ul class="divide-y divide-gray-200">
        <li v-for="invoice in group.invoices" :key="invoice.id" :class="showArchive ? 'bg-gray-50' : ''">
          <div class="px-4 py-4 sm:px-6">
            <div class="flex flex-col lg:flex-row lg:items-center lg:justify-between gap-4">
              <!-- Left Section: Invoice Info -->
              <div class="flex-1 min-w-0">
                <div class="flex items-center flex-wrap gap-2 mb-2">
                  <p class="text-lg font-semibold text-gray-900">Bill no. {{ invoice.id }}</p>
                  <span :class="[
                    'px-2 py-1 text-xs font-semibold rounded-full whitespace-nowrap',
                    showArchive ? 'bg-yellow-100 text-yellow-800' : 'bg-blue-100 text-blue-800'
                  ]">
                    {{ invoice.items?.length || 0 }} items
                  </span>
                  <span v-if="showArchive" class="px-2 py-1 text-xs font-semibold rounded-full bg-red-100 text-red-800 whitespace-nowrap">
                    Archived
                  </span>
                </div>
                <div class="flex flex-col sm:flex-row sm:items-center gap-2 sm:gap-4">
                  <p class="text-sm text-gray-600">
                    <span class="font-medium">Date:</span> {{ formatDate(invoice.created_at) }}
                  </p>
                  <p v-if="showArchive && invoice.deleted_at" class="text-sm text-red-600">
                    <span class="font-medium">Deleted:</span> {{ formatDate(invoice.deleted_at) }}
                  </p>
                </div>
              </div>
              
              <!-- Right Section: Total and Actions -->
              <div class="flex items-center justify-between lg:justify-end gap-4 lg:gap-6">
                <!-- Total Amount -->
                <div class="flex-shrink-0">
                  <p class="text-lg font-bold text-gray-900 whitespace-nowrap">Total: ₹{{ invoice.total_amount?.toFixed(2) || '0.00' }}</p>
                </div>
                
                <!-- Action Buttons -->
                <div class="flex-shrink-0 flex flex-wrap items-center gap-2">
                <button
                  @click="viewInvoice(invoice)"
                  class="px-3 py-1.5 bg-blue-50 hover:bg-blue-100 text-blue-700 rounded-md text-sm font-medium transition-colors shadow-sm hover:shadow"
                  title="View Invoice"
                >
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                    View
                  </span>
                </button>
                <button
                  v-if="!showArchive"
                  @click="editInvoice(invoice)"
                  class="px-3 py-1.5 bg-green-50 hover:bg-green-100 text-green-700 rounded-md text-sm font-medium transition-colors shadow-sm hover:shadow"
                  title="Edit Invoice"
                >
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                    </svg>
                    Edit
                  </span>
                </button>
                <button
                  v-if="!showArchive"
                  @click="deleteInvoice(invoice.id)"
                  class="px-3 py-1.5 bg-red-50 hover:bg-red-100 text-red-700 rounded-md text-sm font-medium transition-colors shadow-sm hover:shadow"
                  title="Delete Invoice"
                >
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                    </svg>
                    Delete
                  </span>
                </button>
                <button
                  @click="printInvoice(invoice)"
                  class="px-3 py-1.5 bg-purple-50 hover:bg-purple-100 text-purple-700 rounded-md text-sm font-medium transition-colors shadow-sm hover:shadow"
                  title="Print Invoice"
                >
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17 17h2a2 2 0 002-2v-4a2 2 0 00-2-2H5a2 2 0 00-2 2v4a2 2 0 002 2h2m2 4h6a2 2 0 002-2v-4a2 2 0 00-2-2H9a2 2 0 00-2 2v4a2 2 0 002 2zm8-12V5a2 2 0 00-2-2H9a2 2 0 00-2 2v4h10z" />
                    </svg>
                    Print
                  </span>
                </button>
                <button
                  v-if="showArchive"
                  @click="restoreInvoice(invoice.id)"
                  class="px-3 py-1.5 bg-green-50 hover:bg-green-100 text-green-700 rounded-md text-sm font-medium transition-colors shadow-sm hover:shadow"
                  title="Restore Invoice"
                >
                  <span class="flex items-center">
                    <svg class="w-4 h-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                    </svg>
                    Restore
                  </span>
                </button>
                </div>
              </div>
            </div>
          </div>
        </li>
        </ul>
      </div>
      
      <!-- Pagination -->
      <div v-if="totalPages > 1" class="flex items-center justify-between bg-white px-4 py-3 border border-gray-200 rounded-lg sticky bottom-0">
        <div class="flex-1 flex justify-between sm:hidden">
          <button
            @click="currentPage = Math.max(1, currentPage - 1)"
            :disabled="currentPage === 1"
            class="relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Previous
          </button>
          <button
            @click="currentPage = Math.min(totalPages, currentPage + 1)"
            :disabled="currentPage === totalPages"
            class="ml-3 relative inline-flex items-center px-4 py-2 border border-gray-300 text-sm font-medium rounded-md text-gray-700 bg-white hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
          >
            Next
          </button>
        </div>
        <div class="hidden sm:flex-1 sm:flex sm:items-center sm:justify-between">
          <div>
            <p class="text-sm text-gray-700">
              Showing <span class="font-medium">{{ startIndex + 1 }}</span> to <span class="font-medium">{{ endIndex }}</span> of <span class="font-medium">{{ groupedInvoices.length }}</span> date groups
            </p>
          </div>
          <div>
            <nav class="relative z-0 inline-flex rounded-md shadow-sm -space-x-px" aria-label="Pagination">
              <button
                @click="currentPage = Math.max(1, currentPage - 1)"
                :disabled="currentPage === 1"
                class="relative inline-flex items-center px-2 py-2 rounded-l-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span class="sr-only">Previous</span>
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M12.707 5.293a1 1 0 010 1.414L9.414 10l3.293 3.293a1 1 0 01-1.414 1.414l-4-4a1 1 0 010-1.414l4-4a1 1 0 011.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
              <button
                v-for="page in visiblePages"
                :key="page"
                @click="goToPage(page)"
                :disabled="page === '...'"
                :class="[
                  'relative inline-flex items-center px-4 py-2 border text-sm font-medium',
                  page === currentPage
                    ? 'z-10 bg-blue-50 border-blue-500 text-blue-600'
                    : page === '...'
                    ? 'bg-white border-gray-300 text-gray-500 cursor-default'
                    : 'bg-white border-gray-300 text-gray-500 hover:bg-gray-50'
                ]"
              >
                {{ page }}
              </button>
              <button
                @click="currentPage = Math.min(totalPages, currentPage + 1)"
                :disabled="currentPage === totalPages"
                class="relative inline-flex items-center px-2 py-2 rounded-r-md border border-gray-300 bg-white text-sm font-medium text-gray-500 hover:bg-gray-50 disabled:opacity-50 disabled:cursor-not-allowed"
              >
                <span class="sr-only">Next</span>
                <svg class="h-5 w-5" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M7.293 14.707a1 1 0 010-1.414L10.586 10 7.293 6.707a1 1 0 011.414-1.414l4 4a1 1 0 010 1.414l-4 4a1 1 0 01-1.414 0z" clip-rule="evenodd" />
                </svg>
              </button>
            </nav>
          </div>
        </div>
      </div>
    </div>

    <!-- View Invoice Modal -->
    <div
      v-if="showViewModal && selectedInvoice"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click.self="showViewModal = false"
    >
      <div class="relative top-10 mx-auto p-5 border w-[95%] max-w-6xl shadow-lg rounded-md bg-white max-h-[90vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-bold">Bill no. {{ selectedInvoice.id }}</h3>
          <button @click="showViewModal = false" class="text-gray-500 hover:text-gray-700">×</button>
        </div>
        <div class="mb-4 space-y-1">
          <p class="text-sm text-gray-600">Date: {{ formatDate(selectedInvoice.created_at) }}</p>
          <p class="text-sm text-gray-600" v-if="selectedInvoice.customer_name"><strong>Customer Name:</strong> {{ selectedInvoice.customer_name }}</p>
          <p class="text-sm text-gray-600" v-if="selectedInvoice.customer_number"><strong>Customer Phone:</strong> {{ selectedInvoice.customer_number }}</p>
          <p class="text-sm text-gray-600" v-if="selectedInvoice.prescribed_by">Prescribed by: Dr. {{ selectedInvoice.prescribed_by }}</p>
        </div>
        <table class="min-w-full divide-y divide-gray-200 mb-4">
          <thead class="bg-gray-50">
            <tr>
              <th class="px-2 py-3 text-left text-xs font-medium text-gray-500 uppercase">S.No.</th>
              <th class="px-2 py-3 text-left text-xs font-medium text-gray-500 uppercase">Particulars</th>
              <th class="px-2 py-3 text-left text-xs font-medium text-gray-500 uppercase">Pack</th>
              <th class="px-2 py-3 text-right text-xs font-medium text-gray-500 uppercase">Quantity</th>
              <th class="px-2 py-3 text-left text-xs font-medium text-gray-500 uppercase">Batch</th>
              <th class="px-2 py-3 text-left text-xs font-medium text-gray-500 uppercase">Expiry</th>
              <th class="px-2 py-3 text-right text-xs font-medium text-gray-500 uppercase">MRP</th>
              <th class="px-2 py-3 text-right text-xs font-medium text-gray-500 uppercase">Amount</th>
            </tr>
          </thead>
          <tbody class="bg-white divide-y divide-gray-200">
            <tr v-for="(item, index) in selectedInvoice.items" :key="item.id">
              <td class="px-2 py-3 text-sm text-gray-500">{{ item.s_no || index + 1 }}</td>
              <td class="px-2 py-3 text-sm text-gray-900">{{ item.medicine?.name || item.medicine_name || 'N/A' }}</td>
              <td class="px-2 py-3 text-sm text-gray-500">{{ item.pack || '—' }}</td>
              <td class="px-2 py-3 text-sm text-gray-500 text-right">{{ item.quantity }}</td>
              <td class="px-2 py-3 text-sm text-gray-500">{{ item.batch || '—' }}</td>
              <td class="px-2 py-3 text-sm text-gray-500">{{ item.expiry ? formatDate(item.expiry) : '—' }}</td>
              <td class="px-2 py-3 text-sm text-gray-500 text-right">₹{{ (item.mrp ?? item.price ?? 0).toFixed(2) }}</td>
              <td class="px-2 py-3 text-sm text-gray-900 text-right">₹{{ item.total?.toFixed(2) }}</td>
            </tr>
          </tbody>
        </table>
        <div class="flex justify-end">
          <p class="text-lg font-semibold">Total Amount: ₹{{ selectedInvoice.total_amount?.toFixed(2) }}</p>
        </div>
      </div>
    </div>

    <!-- Edit Invoice Modal -->
    <div
      v-if="showEditModal && selectedInvoice"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click.self="showEditModal = false"
    >
      <div class="relative top-6 mx-auto p-6 border w-[96%] max-w-7xl shadow-lg rounded-md bg-white max-h-[94vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-lg font-bold">Edit Bill no. {{ selectedInvoice.id }}</h3>
          <button @click="showEditModal = false" class="text-gray-500 hover:text-gray-700">×</button>
        </div>

        <!-- Customer Information -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Customer Name</label>
            <input
              v-model="customerName"
              type="text"
              placeholder="Enter customer name"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Customer Phone Number</label>
            <input
              v-model="customerNumber"
              type="text"
              placeholder="Enter customer phone number"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Prescribed by Dr.</label>
            <input
              v-model="prescribedBy"
              type="text"
              placeholder="Enter doctor name"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Medicine Selection -->
          <div>
            <div class="mb-4">
              <button
                @click="showManualEntry = !showManualEntry"
                class="mb-2 text-sm text-blue-600 hover:text-blue-800"
              >
                {{ showManualEntry ? '← Select from Catalog' : '+ Add Manual Entry' }}
              </button>
            </div>

            <div v-if="!showManualEntry">
              <h4 class="font-semibold mb-2">Add Medicine from Catalog</h4>
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

            <div v-else class="space-y-2 max-h-64 overflow-y-auto">
              <div
                v-for="medicine in availableMedicines"
                :key="medicine.id"
                class="border border-gray-200 rounded-md p-3 hover:bg-gray-50 cursor-pointer"
                @click="addToInvoice(medicine)"
              >
                <div class="flex justify-between items-center">
                  <div>
                    <p class="font-medium">{{ medicine.name }}</p>
                    <p class="text-sm text-gray-500">Stock: {{ medicine.quantity }} | Price: ₹{{ medicine.selling_price }}</p>
                  </div>
                  <button
                    class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm"
                    @click.stop="addToInvoice(medicine)"
                  >
                    Add
                  </button>
                </div>
              </div>
            </div>
            </div>

            <!-- Manual Entry Form -->
            <div v-else class="border border-gray-300 rounded-md p-4">
              <h4 class="font-semibold mb-3">Manual Entry</h4>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div class="sm:col-span-2 lg:col-span-4">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Particulars *</label>
                  <input
                    v-model="manualEntry.particulars"
                    type="text"
                    placeholder="Item name / description"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Pack</label>
                  <input
                    v-model="manualEntry.pack"
                    type="text"
                    placeholder="Pack"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Quantity *</label>
                  <input
                    v-model.number="manualEntry.quantity"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Batch</label>
                  <input
                    v-model="manualEntry.batch"
                    type="text"
                    placeholder="Batch"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">MRP *</label>
                  <input
                    v-model.number="manualEntry.mrp"
                    type="number"
                    step="0.01"
                    min="0"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div class="sm:col-span-2 lg:col-span-1">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Expiry</label>
                  <input
                    v-model="manualEntry.expiry"
                    type="date"
                    class="w-full max-w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div class="sm:col-span-2 lg:col-span-4 text-sm text-gray-600">
                  Amount = Quantity × MRP = ₹{{ ((manualEntry.quantity || 0) * (manualEntry.mrp || 0)).toFixed(2) }}
                </div>
              </div>
              <button
                @click="addManualEntry"
                class="mt-4 w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Add to Invoice
              </button>
            </div>
          </div>

          <!-- Invoice Items -->
          <div>
            <h4 class="font-semibold mb-2">Invoice Items</h4>
            <div v-if="invoiceItems.length === 0" class="text-center py-8 text-gray-500">
              No items in invoice
            </div>
            <div v-else class="space-y-3 max-h-64 overflow-y-auto mb-4">
              <div
                v-for="(item, index) in invoiceItems"
                :key="index"
                class="border border-gray-200 rounded-md p-3"
              >
                <template v-if="editingItemIndex === index">
                  <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-sm mb-3">
                    <div class="sm:col-span-2">
                      <label class="block text-xs text-gray-600">Particulars</label>
                      <input v-model="item.name" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Pack</label>
                      <input v-model="item.pack" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Quantity</label>
                      <input v-model.number="item.quantity" type="number" min="1" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Batch</label>
                      <input v-model="item.batch" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Expiry</label>
                      <input v-model="item.expiry" type="date" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">MRP</label>
                      <input v-model.number="item.mrp" type="number" step="0.01" min="0" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div class="sm:col-span-2 flex items-end gap-2">
                      <button type="button" @click="saveEditItem(index)" class="px-3 py-1.5 bg-green-600 text-white rounded text-sm">Save</button>
                      <button type="button" @click="cancelEditItem(index)" class="px-3 py-1.5 bg-gray-300 rounded text-sm">Cancel</button>
                    </div>
                  </div>
                </template>
                <template v-else>
                <div class="flex justify-between items-start mb-2">
                  <div class="flex-1">
                    <p class="font-medium">{{ item.name }}</p>
                    <p class="text-sm text-gray-500" v-if="item.maxQuantity == null">
                      {{ item.pack || '—' }} · Batch {{ item.batch || '—' }} · Expiry: {{ formatExpiryDisplay(item.expiry) }} · ₹{{ (item.amount != null ? item.amount : (item.quantity * ((item.mrp ?? item.price)))).toFixed(2) }}
                    </p>
                    <p class="text-sm text-gray-500" v-else>₹{{ item.price }} each <span v-if="item.expiry">· Expiry: {{ formatExpiryDisplay(item.expiry) }}</span></p>
                  </div>
                  <div class="flex items-center gap-1">
                    <button type="button" @click="startEditItem(index)" class="text-blue-600 hover:text-blue-800 text-sm">Edit</button>
                    <button type="button" @click="removeFromInvoice(index)" class="text-red-600 hover:text-red-800">×</button>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <button
                      v-if="item.maxQuantity != null"
                      @click="updateItemQuantity(index, -1)"
                      class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                      :disabled="item.quantity <= 1"
                    >
                      -
                    </button>
                    <input
                      v-if="item.maxQuantity != null"
                      v-model.number="item.quantity"
                      type="number"
                      min="1"
                      :max="item.maxQuantity || 999999"
                      class="w-16 text-center border border-gray-300 rounded py-1"
                      @change="validateItemQuantity(index)"
                    />
                    <span v-else class="text-sm">Qty: {{ item.quantity }}</span>
                    <button
                      v-if="item.maxQuantity != null"
                      @click="updateItemQuantity(index, 1)"
                      class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                      :disabled="item.maxQuantity && item.quantity >= item.maxQuantity"
                    >
                      +
                    </button>
                    <div class="ml-auto flex flex-col items-end">
                      <span class="text-xs text-gray-500" v-if="item.discount > 0 && item.maxQuantity != null">
                        Discount: {{ item.discountType === 'percentage' ? item.discount + '%' : '₹' + (item.discount || 0).toFixed(2) }}
                      </span>
                      <span class="font-semibold">₹{{ (item.amount != null ? item.amount : calculateItemTotal(item)).toFixed(2) }}</span>
                    </div>
                  </div>
                  <div v-if="item.maxQuantity != null" class="flex items-center space-x-2 text-xs flex-wrap">
                    <label class="text-gray-600">Price:</label>
                    <input
                      v-model.number="item.price"
                      type="number"
                      step="0.01"
                      min="0"
                      class="w-20 px-2 py-1 border border-gray-300 rounded"
                      @change="updateItemTotal(index)"
                    />
                    <label class="text-gray-600 ml-2">Discount Type:</label>
                    <select
                      v-model="item.discountType"
                      class="w-24 px-2 py-1 border border-gray-300 rounded text-xs"
                      @change="updateItemTotal(index)"
                    >
                      <option value="amount">Amount</option>
                      <option value="percentage">%</option>
                    </select>
                    <label class="text-gray-600 ml-2">Discount:</label>
                    <input
                      v-model.number="item.discount"
                      type="number"
                      step="0.01"
                      min="0"
                      :max="item.discountType === 'percentage' ? 100 : undefined"
                      class="w-20 px-2 py-1 border border-gray-300 rounded"
                      @change="updateItemTotal(index)"
                    />
                    <span class="text-gray-500">{{ item.discountType === 'percentage' ? '%' : '₹' }}</span>
                  </div>
                </div>
                </template>
              </div>
            </div>
            <div class="border-t pt-4">
              <div class="flex justify-between text-lg font-semibold mb-4">
                <span>Total:</span>
                <span>₹{{ invoiceTotal.toFixed(2) }}</span>
              </div>
              <button
                @click="saveInvoice"
                :disabled="saving || invoiceItems.length === 0"
                class="w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md font-medium disabled:opacity-50"
              >
                {{ saving ? 'Saving...' : 'Update Invoice' }}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Create Invoice Modal -->
    <div
      v-if="showAddModal"
      class="fixed inset-0 bg-gray-600 bg-opacity-50 overflow-y-auto h-full w-full z-50"
      @click.self="showAddModal = false"
    >
      <div class="relative top-6 mx-auto p-6 border w-[96%] max-w-7xl shadow-lg rounded-md bg-white max-h-[94vh] overflow-y-auto">
        <div class="flex justify-between items-center mb-4">
          <h3 class="text-2xl font-bold">Create New Invoice</h3>
          <button @click="showAddModal = false" class="text-gray-500 hover:text-gray-700 text-2xl">×</button>
        </div>

        <!-- Customer Information -->
        <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Customer Name</label>
            <input
              v-model="customerName"
              type="text"
              placeholder="Enter customer name"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Customer Phone Number</label>
            <input
              v-model="customerNumber"
              type="text"
              placeholder="Enter customer phone number"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">Prescribed by Dr.</label>
            <input
              v-model="prescribedBy"
              type="text"
              placeholder="Enter doctor name"
              class="w-full px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>

        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
          <!-- Medicine Selection -->
          <div>
            <div class="mb-4">
              <button
                @click="showManualEntry = !showManualEntry"
                class="mb-2 text-sm text-blue-600 hover:text-blue-800"
              >
                {{ showManualEntry ? '← Select from Catalog' : '+ Add Manual Entry' }}
              </button>
            </div>

            <div v-if="!showManualEntry">
              <h4 class="font-semibold mb-2">Add Medicine from Catalog</h4>
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

            <div v-else class="space-y-2 max-h-64 overflow-y-auto">
              <div
                v-for="medicine in availableMedicines"
                :key="medicine.id"
                class="border border-gray-200 rounded-md p-3 hover:bg-gray-50 cursor-pointer"
                @click="addToInvoice(medicine)"
              >
                <div class="flex justify-between items-center">
                  <div>
                    <p class="font-medium">{{ medicine.name }}</p>
                    <p class="text-sm text-gray-500">Stock: {{ medicine.quantity }} | Price: ₹{{ medicine.selling_price }}</p>
                  </div>
                  <button
                    class="bg-blue-600 hover:bg-blue-700 text-white px-3 py-1 rounded text-sm"
                    @click.stop="addToInvoice(medicine)"
                  >
                    Add
                  </button>
                </div>
              </div>
            </div>
            </div>

            <!-- Manual Entry Form -->
            <div v-else class="border border-gray-300 rounded-md p-4">
              <h4 class="font-semibold mb-3">Manual Entry</h4>
              <div class="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-4">
                <div class="sm:col-span-2 lg:col-span-4">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Particulars *</label>
                  <input
                    v-model="manualEntry.particulars"
                    type="text"
                    placeholder="Item name / description"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Pack</label>
                  <input
                    v-model="manualEntry.pack"
                    type="text"
                    placeholder="Pack"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Quantity *</label>
                  <input
                    v-model.number="manualEntry.quantity"
                    type="number"
                    min="1"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">Batch</label>
                  <input
                    v-model="manualEntry.batch"
                    type="text"
                    placeholder="Batch"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div>
                  <label class="block text-sm font-medium text-gray-700 mb-1">MRP *</label>
                  <input
                    v-model.number="manualEntry.mrp"
                    type="number"
                    step="0.01"
                    min="0"
                    class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div class="sm:col-span-2 lg:col-span-1">
                  <label class="block text-sm font-medium text-gray-700 mb-1">Expiry</label>
                  <input
                    v-model="manualEntry.expiry"
                    type="date"
                    class="w-full max-w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  />
                </div>
                <div class="sm:col-span-2 lg:col-span-4 text-sm text-gray-600">
                  Amount = Quantity × MRP = ₹{{ ((manualEntry.quantity || 0) * (manualEntry.mrp || 0)).toFixed(2) }}
                </div>
              </div>
              <button
                @click="addManualEntry"
                class="mt-4 w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md text-sm font-medium"
              >
                Add to Invoice
              </button>
            </div>
          </div>

          <!-- Invoice Items -->
          <div>
            <h4 class="font-semibold mb-2">Invoice Items</h4>
            <div v-if="invoiceItems.length === 0" class="text-center py-8 text-gray-500">
              No items in invoice
            </div>
            <div v-else class="space-y-3 max-h-64 overflow-y-auto mb-4">
              <div
                v-for="(item, index) in invoiceItems"
                :key="index"
                class="border border-gray-200 rounded-md p-3"
              >
                <template v-if="editingItemIndex === index">
                  <div class="grid grid-cols-2 sm:grid-cols-4 gap-2 text-sm mb-3">
                    <div class="sm:col-span-2">
                      <label class="block text-xs text-gray-600">Particulars</label>
                      <input v-model="item.name" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Pack</label>
                      <input v-model="item.pack" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Quantity</label>
                      <input v-model.number="item.quantity" type="number" min="1" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Batch</label>
                      <input v-model="item.batch" type="text" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">Expiry</label>
                      <input v-model="item.expiry" type="date" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div>
                      <label class="block text-xs text-gray-600">MRP</label>
                      <input v-model.number="item.mrp" type="number" step="0.01" min="0" class="w-full px-2 py-1.5 border rounded" />
                    </div>
                    <div class="sm:col-span-2 flex items-end gap-2">
                      <button type="button" @click="saveEditItem(index)" class="px-3 py-1.5 bg-green-600 text-white rounded text-sm">Save</button>
                      <button type="button" @click="cancelEditItem(index)" class="px-3 py-1.5 bg-gray-300 rounded text-sm">Cancel</button>
                    </div>
                  </div>
                </template>
                <template v-else>
                <div class="flex justify-between items-start mb-2">
                  <div class="flex-1">
                    <p class="font-medium">{{ item.name }}</p>
                    <p class="text-sm text-gray-500" v-if="item.maxQuantity == null">
                      {{ item.pack || '—' }} · Batch {{ item.batch || '—' }} · Expiry: {{ formatExpiryDisplay(item.expiry) }} · ₹{{ (item.amount != null ? item.amount : (item.quantity * ((item.mrp ?? item.price)))).toFixed(2) }}</p>
                    <p class="text-sm text-gray-500" v-else>₹{{ item.price }} each <span v-if="item.expiry">· Expiry: {{ formatExpiryDisplay(item.expiry) }}</span></p>
                  </div>
                  <div class="flex items-center gap-1">
                    <button type="button" @click="startEditItem(index)" class="text-blue-600 hover:text-blue-800 text-sm">Edit</button>
                    <button type="button" @click="removeFromInvoice(index)" class="text-red-600 hover:text-red-800">×</button>
                  </div>
                </div>
                <div class="space-y-2">
                  <div class="flex items-center space-x-2">
                    <button
                      v-if="item.maxQuantity != null"
                      @click="updateItemQuantity(index, -1)"
                      class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                      :disabled="item.quantity <= 1"
                    >
                      -
                    </button>
                    <input
                      v-if="item.maxQuantity != null"
                      v-model.number="item.quantity"
                      type="number"
                      min="1"
                      :max="item.maxQuantity || 999999"
                      class="w-16 text-center border border-gray-300 rounded py-1"
                      @change="validateItemQuantity(index)"
                    />
                    <span v-else class="text-sm">Qty: {{ item.quantity }}</span>
                    <button
                      v-if="item.maxQuantity != null"
                      @click="updateItemQuantity(index, 1)"
                      class="bg-gray-200 hover:bg-gray-300 px-2 py-1 rounded"
                      :disabled="item.maxQuantity && item.quantity >= item.maxQuantity"
                    >
                      +
                    </button>
                    <div class="ml-auto flex flex-col items-end">
                      <span class="text-xs text-gray-500" v-if="item.discount > 0 && item.maxQuantity != null">
                        Discount: {{ item.discountType === 'percentage' ? item.discount + '%' : '₹' + (item.discount || 0).toFixed(2) }}
                      </span>
                      <span class="font-semibold">₹{{ (item.amount != null ? item.amount : calculateItemTotal(item)).toFixed(2) }}</span>
                    </div>
                  </div>
                  <div v-if="item.maxQuantity != null" class="flex items-center space-x-2 text-xs flex-wrap">
                    <label class="text-gray-600">Price:</label>
                    <input
                      v-model.number="item.price"
                      type="number"
                      step="0.01"
                      min="0"
                      class="w-20 px-2 py-1 border border-gray-300 rounded"
                      @change="updateItemTotal(index)"
                    />
                    <label class="text-gray-600 ml-2">Discount Type:</label>
                    <select
                      v-model="item.discountType"
                      class="w-24 px-2 py-1 border border-gray-300 rounded text-xs"
                      @change="updateItemTotal(index)"
                    >
                      <option value="amount">Amount</option>
                      <option value="percentage">%</option>
                    </select>
                    <label class="text-gray-600 ml-2">Discount:</label>
                    <input
                      v-model.number="item.discount"
                      type="number"
                      step="0.01"
                      min="0"
                      :max="item.discountType === 'percentage' ? 100 : undefined"
                      class="w-20 px-2 py-1 border border-gray-300 rounded"
                      @change="updateItemTotal(index)"
                    />
                    <span class="text-gray-500">{{ item.discountType === 'percentage' ? '%' : '₹' }}</span>
                  </div>
                </div>
                </template>
              </div>
            </div>
            <div class="border-t pt-4">
              <div class="flex justify-between text-lg font-semibold mb-4">
                <span>Total:</span>
                <span>₹{{ invoiceTotal.toFixed(2) }}</span>
              </div>
              <button
                @click="createInvoice"
                :disabled="saving || invoiceItems.length === 0"
                class="w-full bg-green-600 hover:bg-green-700 text-white px-4 py-2 rounded-md font-medium disabled:opacity-50"
              >
                {{ saving ? 'Creating...' : 'Create Invoice' }}
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
import { ref, computed, onMounted, watch } from 'vue'
import { useToast } from 'vue-toastification'
import { useConfirm } from '../composables/useConfirm'
import ConfirmDialog from '../components/ConfirmDialog.vue'
import { saleService } from '../services/saleService'
import { medicineService } from '../services/medicineService'

const toast = useToast()
const { showConfirm, confirmMessage, confirmTitle, confirm, handleConfirm, handleCancel } = useConfirm()

const invoices = ref([])
const filteredInvoices = ref([])
const archivedInvoices = ref([])
const groupedInvoices = ref([])
const loading = ref(true)
const searchQuery = ref('')
const dateFilter = ref('')
const showViewModal = ref(false)
const showEditModal = ref(false)
const showAddModal = ref(false)
const showArchive = ref(false)
const selectedInvoice = ref(null)
const saving = ref(false)

// Pagination
const currentPage = ref(1)
const itemsPerPage = ref(5) // Number of date groups per page

// Invoice editing
const invoiceItems = ref([])
const medicineSearch = ref('')
const availableMedicines = ref([])
const searchLoading = ref(false)
const showManualEntry = ref(false)
const customerName = ref('')
const customerNumber = ref('')
const prescribedBy = ref('')

function defaultExpiryFirstOfMonth() {
  const d = new Date()
  return d.getFullYear() + '-' + String(d.getMonth() + 1).padStart(2, '0') + '-01'
}

const manualEntry = ref({
  particulars: '',
  pack: '',
  quantity: 1,
  batch: '',
  expiry: defaultExpiryFirstOfMonth(),
  mrp: 0
})

const editingItemIndex = ref(-1)
const itemEditCopy = ref(null)

function formatExpiryDisplay(expiry) {
  if (!expiry) return '—'
  const s = typeof expiry === 'string' ? expiry.slice(0, 10) : String(expiry).slice(0, 10)
  if (s.length >= 7) return s.slice(5, 7) + '-' + s.slice(0, 4)
  return s || '—'
}

function startEditItem(index) {
  const item = invoiceItems.value[index]
  itemEditCopy.value = {
    name: item.name,
    pack: item.pack ?? '',
    quantity: item.quantity,
    batch: item.batch ?? '',
    expiry: item.expiry ?? '',
    mrp: item.mrp ?? item.price ?? 0,
    amount: item.amount,
    price: item.price,
    discount: item.discount,
    discountType: item.discountType,
    maxQuantity: item.maxQuantity,
    id: item.id
  }
  editingItemIndex.value = index
}

function saveEditItem(index) {
  const item = invoiceItems.value[index]
  if (item.maxQuantity == null) {
    item.amount = (Number(item.quantity) || 0) * ((Number(item.mrp) ?? Number(item.price)) || 0)
    item.price = item.mrp
  } else {
    if (item.mrp != null) item.price = item.mrp
  }
  editingItemIndex.value = -1
  itemEditCopy.value = null
}

function cancelEditItem(index) {
  if (itemEditCopy.value) {
    const item = invoiceItems.value[index]
    Object.assign(item, itemEditCopy.value)
  }
  editingItemIndex.value = -1
  itemEditCopy.value = null
}

const calculateItemTotal = (item) => {
  if (item.amount != null && item.maxQuantity == null) {
    return item.amount
  }
  const subtotal = item.price * item.quantity
  const discountType = item.discountType || 'amount'
  let discountAmount = 0
  if (discountType === 'percentage') {
    discountAmount = subtotal * ((item.discount || 0) / 100)
  } else {
    discountAmount = item.discount || 0
  }
  return subtotal - discountAmount
}

const invoiceTotal = computed(() => {
  return invoiceItems.value.reduce((sum, item) => {
    return sum + calculateItemTotal(item)
  }, 0)
})

const loadInvoices = async () => {
  try {
    loading.value = true
    if (showArchive.value) {
      archivedInvoices.value = await saleService.getArchived()
    } else {
      invoices.value = await saleService.getAll()
    }
    filterInvoices()
  } catch (error) {
    console.error('Failed to load invoices:', error)
  } finally {
    loading.value = false
  }
}

const filterInvoices = () => {
  let filtered = showArchive.value ? archivedInvoices.value : invoices.value
  
  // Apply search filter
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase().trim()
    filtered = filtered.filter(inv => {
      // Search by invoice ID
      if (inv.id.toString().includes(query)) return true
      
      // Search by customer name
      if (inv.customer_name && inv.customer_name.toLowerCase().includes(query)) return true
      
      // Search by customer phone number
      if (inv.customer_number && inv.customer_number.toLowerCase().includes(query)) return true
      
      return false
    })
  }
  
  // Apply date filter
  if (dateFilter.value) {
    const filterDate = dateFilter.value
    filtered = filtered.filter(inv => {
      const invoiceDate = new Date(inv.created_at).toISOString().split('T')[0]
      return invoiceDate === filterDate
    })
  }
  
  filteredInvoices.value = filtered
  groupInvoicesByDate()
  currentPage.value = 1 // Reset to first page when filtering
}

const groupInvoicesByDate = () => {
  const groups = {}
  
  filteredInvoices.value.forEach(invoice => {
    const date = new Date(invoice.created_at).toISOString().split('T')[0]
    if (!groups[date]) {
      groups[date] = []
    }
    groups[date].push(invoice)
  })
  
  // Convert to array and sort by date (newest first)
  groupedInvoices.value = Object.keys(groups)
    .sort((a, b) => new Date(b) - new Date(a))
    .map(date => ({
      date,
      invoices: groups[date].sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    }))
}

const formatDateHeader = (dateString) => {
  const date = new Date(dateString)
  const today = new Date()
  const yesterday = new Date(today)
  yesterday.setDate(yesterday.getDate() - 1)
  
  const dateStr = date.toISOString().split('T')[0]
  const todayStr = today.toISOString().split('T')[0]
  const yesterdayStr = yesterday.toISOString().split('T')[0]
  
  if (dateStr === todayStr) {
    return 'Today'
  } else if (dateStr === yesterdayStr) {
    return 'Yesterday'
  } else {
    return date.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    })
  }
}

// Pagination computed properties
const totalPages = computed(() => {
  return Math.ceil(groupedInvoices.value.length / itemsPerPage.value)
})

const startIndex = computed(() => {
  return (currentPage.value - 1) * itemsPerPage.value
})

const endIndex = computed(() => {
  return Math.min(startIndex.value + itemsPerPage.value, groupedInvoices.value.length)
})

const paginatedGroups = computed(() => {
  return groupedInvoices.value.slice(startIndex.value, endIndex.value)
})

const visiblePages = computed(() => {
  const pages = []
  const total = totalPages.value
  const current = currentPage.value
  
  if (total <= 7) {
    // Show all pages if 7 or fewer
    for (let i = 1; i <= total; i++) {
      pages.push(i)
    }
  } else {
    // Show first page, last page, current page, and pages around current
    if (current <= 3) {
      for (let i = 1; i <= 5; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    } else if (current >= total - 2) {
      pages.push(1)
      pages.push('...')
      for (let i = total - 4; i <= total; i++) pages.push(i)
    } else {
      pages.push(1)
      pages.push('...')
      for (let i = current - 1; i <= current + 1; i++) pages.push(i)
      pages.push('...')
      pages.push(total)
    }
  }
  
  return pages.filter((page, index, arr) => {
    if (page === '...') return true
    return arr.indexOf(page) === index
  })
})

const goToPage = (page) => {
  if (page !== '...' && typeof page === 'number') {
    currentPage.value = page
  }
}

const viewInvoice = async (invoice) => {
  // Load full invoice details to ensure prescribed_by is included
  try {
    const fullInvoice = await saleService.getById(invoice.id)
    selectedInvoice.value = fullInvoice
    showViewModal.value = true
  } catch (error) {
    selectedInvoice.value = invoice
    showViewModal.value = true
  }
}

const editInvoice = async (invoice) => {
  try {
    const fullInvoice = await saleService.getById(invoice.id)
    selectedInvoice.value = fullInvoice
    customerName.value = fullInvoice.customer_name || ''
    customerNumber.value = fullInvoice.customer_number || ''
    prescribedBy.value = fullInvoice.prescribed_by || ''
    // When editing, maxQuantity should be current stock + quantity already in invoice
    invoiceItems.value = fullInvoice.items.map(item => ({
      id: item.medicine_id,
      name: item.medicine?.name || item.medicine_name || 'N/A',
      price: item.price,
      quantity: item.quantity,
      discount: item.discount || 0,
      discountType: item.discount_type || 'amount',
      maxQuantity: item.medicine_id ? ((item.medicine?.quantity || 0) + item.quantity) : null,
      s_no: item.s_no,
      pack: item.pack ?? '',
      batch: item.batch ?? '',
      expiry: item.expiry ? (typeof item.expiry === 'string' ? item.expiry.slice(0, 10) : item.expiry) : '',
      mrp: item.mrp ?? 0,
      amount: item.medicine_id ? null : (item.total ?? 0)
    }))
    editingItemIndex.value = -1
    showEditModal.value = true
    loadAllMedicines()
  } catch (error) {
    toast.error('Failed to load invoice: ' + (error.response?.data?.error || error.message))
  }
}

const deleteInvoice = async (id) => {
  const confirmed = await confirm('Are you sure you want to delete this invoice? It will be moved to archive and stock will be restored.', 'Delete Invoice')
  if (!confirmed) {
    return
  }
  try {
    await saleService.delete(id)
    loadInvoices()
    toast.success('Invoice moved to archive successfully!')
  } catch (error) {
    toast.error('Failed to delete invoice: ' + (error.response?.data?.error || error.message))
  }
}

const restoreInvoice = async (id) => {
  const confirmed = await confirm('Are you sure you want to restore this invoice? Stock will be reduced again.', 'Restore Invoice')
  if (!confirmed) {
    return
  }
  try {
    await saleService.restore(id)
    loadInvoices()
    toast.success('Invoice restored successfully!')
  } catch (error) {
    toast.error('Failed to restore invoice: ' + (error.response?.data?.error || error.message))
  }
}

const printInvoice = async (invoice) => {
  // Load full invoice details to ensure prescribed_by is included
  let fullInvoice = invoice
  try {
    fullInvoice = await saleService.getById(invoice.id)
  } catch (error) {
    console.error('Failed to load full invoice:', error)
  }
  
  // Load settings for company details
  let settings = { company_name: '', address: '', phone_number: '', license_id: '', gstin_number: '' }
  try {
    const { settingsService } = await import('../services/settingsService')
    settings = await settingsService.get()
  } catch (error) {
    console.error('Failed to load settings:', error)
  }

  const printWindow = window.open('', '_blank')
  const billContent = `
    <html>
      <head>
        <title>Bill no. ${fullInvoice.id}</title>
        <style>
          body { font-family: Arial, sans-serif; padding: 20px; }
          .header { margin-bottom: 20px; }
          .header-top { display: flex; justify-content: space-between; align-items: flex-start; margin-bottom: 10px; }
          .header-left { flex: 1; }
          .header-right { text-align: right; }
          h1 { margin: 0 0 10px 0; font-size: 24px; }
          .company-details { font-size: 14px; line-height: 1.6; }
          .company-details p { margin: 4px 0; }
          .invoice-info { margin: 15px 0; display: flex; justify-content: space-between; }
          .invoice-info-left { flex: 1; }
          .invoice-info-right { text-align: right; }
          table { width: 100%; border-collapse: collapse; margin: 20px 0; font-size: 11px; }
          th, td { border: 1px solid #ddd; padding: 5px 6px; text-align: left; }
          th { background-color: #f2f2f2; }
          td.expiry-cell { min-width: 52px; white-space: nowrap; }
          .total { font-size: 16px; font-weight: bold; text-align: right; margin-top: 20px; }
          .signature-block { margin-top: 40px; text-align: right; }
          .signature-block img { display: block; margin-left: auto; margin-bottom: 6px; max-height: 48px; max-width: 160px; object-fit: contain; }
          .signature-label { font-size: 12px; font-weight: 600; }
        </style>
      </head>
      <body>
        <div class="header">
          <div class="header-top">
            <div class="header-left">
              <h1>${settings.company_name || 'Pharmacy ERP'}</h1>
              <p><strong>D.L.No:</strong> ${settings.license_id || 'N/A'}</p>
              <div class="company-details">
                ${settings.address ? `<p>${settings.address}</p>` : ''}
                ${settings.phone_number ? `<p>Phone: ${settings.phone_number}</p>` : ''}
              </div>
            </div>
            <div class="header-right">
              ${settings.gstin_number ? `<p><strong>GSTIN:</strong> ${settings.gstin_number}</p>` : ''}
            </div>
          </div>
        </div>
        <div class="invoice-info">
          <div class="invoice-info-left">
            <p><strong>Bill No:</strong> ${fullInvoice.id}</p>
            <p><strong>Date:</strong> ${formatDate(fullInvoice.created_at)}</p>
            ${fullInvoice.customer_name ? `<p><strong>Customer Name:</strong> ${fullInvoice.customer_name}</p>` : ''}
            ${fullInvoice.customer_number ? `<p><strong>Customer Phone:</strong> ${fullInvoice.customer_number}</p>` : ''}
            ${fullInvoice.prescribed_by ? `<p><strong>Prescribed by:</strong> Dr. ${fullInvoice.prescribed_by}</p>` : ''}
          </div>
        </div>
        <table>
          <thead>
            <tr>
              <th>S.No.</th>
              <th>Particulars</th>
              <th>Pack</th>
              <th>Qty</th>
              <th>Batch</th>
              <th>Expiry</th>
              <th>MRP</th>
              <th>Amount</th>
            </tr>
          </thead>
          <tbody>
            ${fullInvoice.items?.map((item, index) => {
              let expiryStr = '—';
              if (item.expiry) {
                const s = typeof item.expiry === 'string' ? item.expiry.slice(0, 10) : String(item.expiry).slice(0, 10);
                if (s.length >= 7) expiryStr = s.slice(5, 7) + '-' + s.slice(0, 4);
              }
              return `
              <tr>
                <td>${item.s_no || index + 1}</td>
                <td>${item.medicine?.name || item.medicine_name || 'N/A'}</td>
                <td>${item.pack || '—'}</td>
                <td>${item.quantity}</td>
                <td>${item.batch || '—'}</td>
                <td class="expiry-cell">${expiryStr}</td>
                <td>₹${(item.mrp ?? item.price ?? 0).toFixed(2)}</td>
                <td>₹${item.total?.toFixed(2)}</td>
              </tr>
            `}).join('') || ''}
          </tbody>
        </table>
        <div class="total">Total Amount: ₹${fullInvoice.total_amount?.toFixed(2)}</div>
        <div class="signature-block">
          <img src="/signature.png" alt="Signature" />
          <div class="signature-label">Swasth Pharmacy</div>
        </div>
      </body>
    </html>
  `
  printWindow.document.write(billContent)
  printWindow.document.close()
  printWindow.print()
}

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

const loadAllMedicines = async () => {
  try {
    const medicines = await medicineService.getAll()
    availableMedicines.value = medicines.filter(m => m.quantity > 0)
  } catch (error) {
    console.error('Failed to load medicines:', error)
  }
}

const addToInvoice = (medicine) => {
  const existingIndex = invoiceItems.value.findIndex(item => item.id === medicine.id)
  if (existingIndex >= 0) {
    if (!invoiceItems.value[existingIndex].maxQuantity || invoiceItems.value[existingIndex].quantity < invoiceItems.value[existingIndex].maxQuantity) {
      invoiceItems.value[existingIndex].quantity++
    }
  } else {
    invoiceItems.value.push({
      id: medicine.id,
      name: medicine.name,
      price: medicine.selling_price,
      quantity: 1,
      discount: 0,
      discountType: 'amount',
      maxQuantity: medicine.quantity
    })
  }
}

const openAddModal = () => {
  showAddModal.value = true
  invoiceItems.value = []
  showManualEntry.value = false
  editingItemIndex.value = -1
  manualEntry.value = { particulars: '', pack: '', quantity: 1, batch: '', expiry: defaultExpiryFirstOfMonth(), mrp: 0 }
}

const addManualEntry = () => {
  const e = manualEntry.value
  if (!e.particulars?.trim() || !e.quantity || (e.mrp ?? 0) < 0) {
    toast.warning('Please fill in Particulars, Quantity and MRP')
    return
  }
  const qty = Number(e.quantity) || 1
  const mrp = Number(e.mrp) || 0
  const amount = qty * mrp
  invoiceItems.value.push({
    id: null,
    name: e.particulars.trim(),
    pack: e.pack || '',
    quantity: qty,
    batch: e.batch || '',
    expiry: e.expiry || '',
    mrp,
    amount,
    price: mrp,
    discount: 0,
    discountType: 'amount',
    maxQuantity: null
  })
  manualEntry.value = { particulars: '', pack: '', quantity: 1, batch: '', expiry: defaultExpiryFirstOfMonth(), mrp: 0 }
  showManualEntry.value = false
}

const updateItemTotal = (index) => {
  // Recalculate is handled by computed property
}

const removeFromInvoice = (index) => {
  invoiceItems.value.splice(index, 1)
}

const updateItemQuantity = (index, change) => {
  const item = invoiceItems.value[index]
  const newQuantity = item.quantity + change
  if (newQuantity >= 1 && newQuantity <= item.maxQuantity) {
    item.quantity = newQuantity
  }
}

const validateItemQuantity = (index) => {
  const item = invoiceItems.value[index]
  if (item.quantity < 1) item.quantity = 1
  if (item.quantity > item.maxQuantity) item.quantity = item.maxQuantity
}

const createInvoice = async () => {
  if (invoiceItems.value.length === 0) return

  try {
    saving.value = true
    const saleData = {
      customer_name: customerName.value,
      customer_number: customerNumber.value,
      prescribed_by: prescribedBy.value,
      items: invoiceItems.value.map((item, index) => {
        const payload = {
          medicine_id: item.id || null,
          medicine_name: item.name,
          s_no: item.s_no != null ? item.s_no : index + 1,
          pack: item.pack ?? '',
          quantity: item.quantity,
          batch: item.batch ?? '',
          expiry: item.expiry || '',
          mrp: item.mrp ?? 0,
          price: item.price,
          discount: item.discount || 0,
          discount_type: item.discountType || 'amount'
        }
        if (item.amount != null && item.maxQuantity == null) {
          payload.total = item.amount
        }
        return payload
      })
    }

    await saleService.create(saleData)
    showAddModal.value = false
    invoiceItems.value = []
    medicineSearch.value = ''
    availableMedicines.value = []
    customerName.value = ''
    customerNumber.value = ''
    prescribedBy.value = ''
    showManualEntry.value = false
    loadInvoices()
    toast.success('Invoice created successfully!')
  } catch (error) {
    toast.error('Failed to create invoice: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

const saveInvoice = async () => {
  if (invoiceItems.value.length === 0 || !selectedInvoice.value) return

  try {
    saving.value = true
    const saleData = {
      customer_name: customerName.value,
      customer_number: customerNumber.value,
      prescribed_by: prescribedBy.value,
      items: invoiceItems.value.map((item, index) => {
        const payload = {
          medicine_id: item.id || null,
          medicine_name: item.name,
          s_no: item.s_no != null ? item.s_no : index + 1,
          pack: item.pack ?? '',
          quantity: item.quantity,
          batch: item.batch ?? '',
          expiry: item.expiry || '',
          mrp: item.mrp ?? 0,
          price: item.price,
          discount: item.discount || 0,
          discount_type: item.discountType || 'amount'
        }
        if (item.amount != null && item.maxQuantity == null) {
          payload.total = item.amount
        }
        return payload
      })
    }

    await saleService.update(selectedInvoice.value.id, saleData)
    showEditModal.value = false
    invoiceItems.value = []
    editingItemIndex.value = -1
    medicineSearch.value = ''
    availableMedicines.value = []
    customerName.value = ''
    customerNumber.value = ''
    prescribedBy.value = ''
    selectedInvoice.value = null
    showManualEntry.value = false
    loadInvoices()
    toast.success('Invoice updated successfully!')
  } catch (error) {
    toast.error('Failed to update invoice: ' + (error.response?.data?.error || error.message))
  } finally {
    saving.value = false
  }
}

const formatDate = (dateString) => {
  if (!dateString) return 'N/A'
  const date = new Date(dateString)
  return date.toLocaleString()
}

// Watch for archive toggle
watch(showArchive, () => {
  loadInvoices()
})

onMounted(() => {
  loadInvoices()
})
</script>
