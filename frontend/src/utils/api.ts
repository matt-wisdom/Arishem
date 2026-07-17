const API_BASE = import.meta.env.VITE_API_URL || ''

export const apiUrl = (path: string) => {
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${cleanPath}`
}