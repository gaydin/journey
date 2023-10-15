import axios from 'axios'
import useAuthStore from '@/stores/auth.store'

const axiosInstance = axios.create({
  withCredentials: true
})
axiosInstance.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    if (error.response.status === 401) {
      const authStore = useAuthStore()
      return authStore.logout()
    }
    return error
  }
)

export default axiosInstance
