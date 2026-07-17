const API_BASE = import.meta.env.VITE_API_URL || ''

export const apiUrl = (path: string) => {
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${cleanPath}`
}

export const handleUnauthorized = () => {
  localStorage.removeItem('token')
  window.location.href = '/sign-in'
}