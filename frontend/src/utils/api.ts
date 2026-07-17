const API_BASE = import.meta.env.VITE_API_URL || ''

export const apiUrl = (path: string) => {
  const cleanPath = path.startsWith('/') ? path : `/${path}`
  return `${API_BASE}${cleanPath}`
}

let isHandling401 = false

const handleUnauthorized = () => {
  if (isHandling401) return
  isHandling401 = true

  localStorage.clear()
  sessionStorage.clear()
  
  const cookies = document.cookie.split(';')
  for (const cookie of cookies) {
    const [name] = cookie.trim().split('=')
    if (name) {
      document.cookie = `${name}=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;`
    }
  }

  const clerk = (window as any).__clerk__
  if (clerk?.signOut) {
    clerk.signOut().then(() => {
      window.location.href = '/sign-in'
    })
  } else {
    window.location.href = '/sign-in'
  }
}

export const apiFetch = async (url: string, options: RequestInit = {}): Promise<Response> => {
  if (!options.headers) {
    options.headers = {}
  }
  
  const res = await fetch(url, options)

  if (res.status === 401) {
    handleUnauthorized()
  }

  return res
}