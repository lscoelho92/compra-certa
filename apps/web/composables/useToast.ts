export type ToastVariant = 'default' | 'destructive'

export type ToastMessage = {
  id: string
  title: string
  description?: string
  variant?: ToastVariant
}

const TOAST_DURATION_MS = 3200

export const useToast = () => {
  const toasts = useState<ToastMessage[]>('toasts', () => [])

  const remove = (id: string) => {
    toasts.value = toasts.value.filter((toast) => toast.id !== id)
  }

  const push = (message: Omit<ToastMessage, 'id'>) => {
    if (!process.client) return
    const id = crypto.randomUUID()
    const entry: ToastMessage = { id, ...message }
    toasts.value = [...toasts.value, entry]

    setTimeout(() => {
      remove(id)
    }, TOAST_DURATION_MS)
  }

  const success = (title: string, description?: string) => {
    push({ title, description })
  }

  const error = (title: string, description?: string) => {
    push({ title, description, variant: 'destructive' })
  }

  return { toasts, push, success, error, remove }
}
