const API_BASE = import.meta.env.VITE_API_URL || ''

export const apiUrl = (path: string) => {
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${cleanPath}`
}

let isRedirecting = false

export const handleUnauthorized = () => {
  if (isRedirecting) return
  if (window.location.pathname === '/sign-in' || window.location.pathname === '/sign-up') return
  isRedirecting = true
  localStorage.removeItem('token')
  window.location.replace('/sign-in')
}