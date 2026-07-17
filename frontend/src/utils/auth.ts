import { useAuth } from '@clerk/vue'

export const useApiHeaders = () => {
  const { getToken } = useAuth()

  const getHeaders = async () => {
    const tokenFn = typeof getToken.value === 'function' ? getToken.value : getToken
    const token = (await (tokenFn as any)({ template: 'arishem-backend' })) || localStorage.getItem('token') || ''
    return { 
      'Authorization': `Bearer ${token}`,
      'Content-Type': 'application/json'
    }
  }

  return { getHeaders, getToken }
}