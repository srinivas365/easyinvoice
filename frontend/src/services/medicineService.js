import api from './api'

export const medicineService = {
  async getAll(search = '') {
    const response = await api.get('/medicines', { params: { search } })
    return response.data
  },

  async getById(id) {
    const response = await api.get(`/medicines/${id}`)
    return response.data
  },

  async create(medicine) {
    const response = await api.post('/medicines', medicine)
    return response.data
  },

  async update(id, medicine) {
    const response = await api.put(`/medicines/${id}`, medicine)
    return response.data
  },

  async delete(id) {
    const response = await api.delete(`/medicines/${id}`)
    return response.data
  }
}
