import api from './api'

export const distributorService = {
  async getAll() {
    const response = await api.get('/distributors')
    return response.data
  },
  
  async getById(id) {
    const response = await api.get(`/distributors/${id}`)
    return response.data
  },
  
  async create(distributor) {
    const response = await api.post('/distributors', distributor)
    return response.data
  },
  
  async update(id, distributor) {
    const response = await api.put(`/distributors/${id}`, distributor)
    return response.data
  },
  
  async delete(id) {
    const response = await api.delete(`/distributors/${id}`)
    return response.data
  }
}
