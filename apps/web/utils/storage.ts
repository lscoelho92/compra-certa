export function readStorage<T>(key: string, fallback: T): T {
  if (!process.client) return fallback
  try {
    const raw = localStorage.getItem(key)
    if (!raw) return fallback
    return JSON.parse(raw) as T
  } catch {
    return fallback
  }
}

export function writeStorage<T>(key: string, value: T): void {
  if (!process.client) return
  localStorage.setItem(key, JSON.stringify(value))
}
