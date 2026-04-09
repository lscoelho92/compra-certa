export type Product = {
  id: string
  name: string
  description?: string
  categoryId?: string | null
  defaultPrice?: number | null
  createdAt?: string
  updatedAt?: string
}

export type CreateProductPayload = {
  name: string
  description?: string
  categoryId?: string | null
  defaultPrice?: number | null
}

export type UpdateProductPayload = Partial<CreateProductPayload>

export const useProducts = () => {
  const products = useState<Product[]>('products', () => [])
  const config = useRuntimeConfig()

  const fetchProducts = async () => {
    const data = await $fetch<Product[]>(`${config.public.apiBase}/products`)
    products.value = data
    return data
  }

  const createProduct = async (data: CreateProductPayload) => {
    const product = await $fetch<Product>(`${config.public.apiBase}/products`, {
      method: 'POST',
      body: data
    })
    products.value = [product, ...products.value]
    return product
  }

  const updateProduct = async (id: string, data: UpdateProductPayload) => {
    const product = await $fetch<Product>(`${config.public.apiBase}/products/${id}`, {
      method: 'PATCH',
      body: data
    })
    products.value = products.value.map((item) =>
      item.id === id ? { ...item, ...product } : item
    )
    return product
  }

  const deleteProduct = async (id: string) => {
    await $fetch(`${config.public.apiBase}/products/${id}`, {
      method: 'DELETE'
    })
    products.value = products.value.filter((item) => item.id !== id)
  }

  return {
    products,
    fetchProducts,
    createProduct,
    updateProduct,
    deleteProduct
  }
}
