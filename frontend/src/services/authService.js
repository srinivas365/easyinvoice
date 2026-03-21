import api from './api'

export const authService = {
  async login(username, password) {
    const response = await api.post('/login', { username, password })
    return response.data
  },
  
  async getProfile() {
    const response = await api.get('/profile')
    return response.data
  }
}
