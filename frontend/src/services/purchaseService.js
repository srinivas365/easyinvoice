import api from './api'

export const purchaseService = {
  async getAll() {
    const response = await api.get('/purchases')
    return response.data
  },
  
  async getById(id) {
    const response = await api.get(`/purchases/${id}`)
    return response.data
  },
  
  async create(purchase) {
    const response = await api.post('/purchases', purchase)
    return response.data
  }
}
