import api from './api'

export const purchaseInvoiceService = {
  async getAll() {
    const response = await api.get('/purchase-invoices')
    return response.data
  },

  async getById(id) {
    const response = await api.get(`/purchase-invoices/${id}`)
    return response.data
  },

  async create(invoice) {
    const response = await api.post('/purchase-invoices', invoice)
    return response.data
  },

  async update(id, invoice) {
    const response = await api.put(`/purchase-invoices/${id}`, invoice)
    return response.data
  },

  async delete(id) {
    const response = await api.delete(`/purchase-invoices/${id}`)
    return response.data
  }
}
