export type Category = {
  id: string
  name: string
  createdAt?: string
}

export type CreateCategoryPayload = {
  name: string
}

export const useCategories = () => {
  const categories = useState<Category[]>('categories', () => [])
  const config = useRuntimeConfig()

  const fetchCategories = async () => {
    const data = await $fetch<Category[]>(`${config.public.apiBase}/categories`)
    categories.value = data
    return data
  }

  const createCategory = async (data: CreateCategoryPayload) => {
    const category = await $fetch<Category>(`${config.public.apiBase}/categories`, {
      method: 'POST',
      body: data
    })
    categories.value = [...categories.value, category].sort((a, b) => a.name.localeCompare(b.name))
    return category
  }

  const deleteCategory = async (id: string) => {
    await $fetch(`${config.public.apiBase}/categories/${id}`, { method: 'DELETE' })
    categories.value = categories.value.filter((category) => category.id !== id)
  }

  return {
    categories,
    fetchCategories,
    createCategory,
    deleteCategory
  }
}
