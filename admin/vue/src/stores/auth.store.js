import { defineStore } from 'pinia'

import axios from 'axios'
import router from '@/router'

// const baseUrl = `${import.meta.env.VITE_API_URL}/users`;
const axiosInstance = axios.create({
  withCredentials: true
})

const useAuthStore = defineStore({
  id: 'auth',
  state: () => ({
    // initialize state from local storage to enable user to stay logged in
    isAuthenticated: localStorage.getItem('isAuthenticated') === 'true',
    returnUrl: null
  }),
  actions: {
    setAuth() {
      this.isAuthenticated = true
      localStorage.setItem('isAuthenticated', 'true')
      return router.push(this.returnUrl || '/')
    },
    // TODO
    async login(username, password) {
      axiosInstance
        .post('/admin/v1/api/auth/login', {
          login: username,
          password: password
        })
        .then((response) => {
          if (!response) {
            return
          }
          this.setAuth()
        })
        .catch(function (error) {
          return error
        })
    },

    logout() {
      this.isAuthenticated = false
      localStorage.removeItem('isAuthenticated')
      router.push('/login')
    }
  }
})

export default useAuthStore
