export type Purchase = {
  id: string
  purchase_date: string
  month: string
  total_price: number
  items: PurchaseItem[]
}

export type PurchaseItem = {
  product_id: string
  product_name: string
  quantity: number
  unit_price: number
  category?: string
}

export type PurchasePayload = {
  purchase_date: string
  items: Array<{
    product_id: string
    quantity: number
    unit_price: number
  }>
}

export const usePurchases = () => {
  const purchases = useState<Purchase[]>('purchases', () => [])
  const config = useRuntimeConfig()

  const fetchPurchases = async () => {
    const data = await $fetch<Purchase[]>(`${config.public.apiBase}/purchases`)
    purchases.value = data
    return data
  }

  const createPurchase = async (data: PurchasePayload) => {
    await $fetch(`${config.public.apiBase}/purchases`, {
      method: 'POST',
      body: {
        purchaseDate: data.purchase_date,
        items: data.items.map((item) => ({
          productId: item.product_id,
          quantity: item.quantity,
          unitPrice: item.unit_price
        }))
      }
    })
    return fetchPurchases()
  }

  const updatePurchase = async (id: string, data: PurchasePayload) => {
    await $fetch(`${config.public.apiBase}/purchases/${id}`, {
      method: 'PATCH',
      body: {
        purchaseDate: data.purchase_date,
        items: data.items.map((item) => ({
          productId: item.product_id,
          quantity: item.quantity,
          unitPrice: item.unit_price
        }))
      }
    })
    return fetchPurchases()
  }

  const deletePurchase = async (id: string) => {
    await $fetch(`${config.public.apiBase}/purchases/${id}`, {
      method: 'DELETE'
    })
    purchases.value = purchases.value.filter((item) => item.id !== id)
  }

  return {
    purchases,
    fetchPurchases,
    createPurchase,
    updatePurchase,
    deletePurchase
  }
}
