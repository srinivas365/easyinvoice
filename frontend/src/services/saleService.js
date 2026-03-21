import api from './api'

export const saleService = {
  async getAll() {
    const response = await api.get('/sales')
    return response.data
  },
  
  async getById(id) {
    const response = await api.get(`/sales/${id}`)
    return response.data
  },
  
  async create(sale) {
    const response = await api.post('/sales', sale)
    return response.data
  },
  
  async update(id, sale) {
    const response = await api.put(`/sales/${id}`, sale)
    return response.data
  },
  
  async delete(id) {
    const response = await api.delete(`/sales/${id}`)
    return response.data
  },
  
  async getArchived() {
    const response = await api.get('/sales/archived')
    return response.data
  },
  
  async restore(id) {
    const response = await api.post(`/sales/${id}/restore`)
    return response.data
  }
}
