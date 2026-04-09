<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-bold tracking-tight">Produtos</h1>
        <p class="text-muted-foreground text-sm mt-1">{{ products.length }} produtos cadastrados</p>
      </div>
      <div class="flex flex-wrap items-center gap-2">
        <button
          class="inline-flex items-center gap-2 rounded-2xl border border-border bg-white/80 px-4 py-2 text-sm font-semibold text-slate-700 hover:bg-white"
          @click="categoriesOpen = true"
        >
          Categorias
        </button>
        <button
          class="inline-flex items-center gap-2 rounded-2xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground"
          @click="openCreate"
        >
          <Plus class="h-4 w-4" /> Novo Produto
        </button>
      </div>
    </div>

    <div class="relative max-w-sm">
      <Search class="absolute right-3 top-1/2 -translate-y-1/2 h-4 w-4" style="color: #0d9488;" />
      <input
        v-model.trim="search"
        placeholder="Buscar produtos..."
        class="w-full pr-10 border-0 border-b-2 rounded-none bg-transparent focus:outline-none focus:ring-0"
        style="border-bottom-color: hsl(168, 76%, 42%);"
      />
    </div>

    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="w-8 h-8 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>

    <template v-else>
      <div
        v-if="filtered.length === 0"
        class="text-center py-16 rounded-3xl relative"
        style="box-shadow: 0 0 60px 10px rgba(16,185,129,0.18), 0 0 0 2px rgba(45,212,191,0.2); background: #fff;"
      >
        <p class="text-muted-foreground">Nenhum produto encontrado</p>
        <button
          class="mt-4 rounded-2xl border-0 px-4 py-2 text-sm"
          style="background: rgba(16,185,129,0.15); color: #0d9488;"
          @click="openCreate"
        >
          Criar primeiro produto
        </button>
      </div>
      <div v-else class="grid sm:grid-cols-2 lg:grid-cols-3 gap-4">
        <ProductCard
          v-for="product in filtered"
          :key="product.id"
          :product="product"
          :category-label="product.categoryId ? categoryById.get(product.categoryId) : undefined"
          @edit="openEdit"
          @delete="confirmDelete"
        />
      </div>
    </template>

    <ProductForm
      v-model:open="formOpen"
      :product="editing"
      :categories="categories"
      @submit="handleSubmit"
    />

    <Teleport to="body">
      <div v-if="categoriesOpen" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-slate-900/40" @click="closeCategories" />
        <div class="relative w-full max-w-lg rounded-2xl bg-white p-6 shadow-2xl">
          <div class="flex items-center justify-between">
            <h2 class="text-lg font-semibold">Categorias</h2>
            <button class="text-sm text-muted-foreground hover:text-foreground" @click="closeCategories">
              Fechar
            </button>
          </div>
          <form class="mt-4 flex gap-2" @submit.prevent="handleCreateCategory">
            <input
              v-model.trim="newCategoryName"
              class="flex-1 rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
              placeholder="Nova categoria"
              required
            />
            <button
              type="submit"
              class="rounded-xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground"
              :disabled="creatingCategory"
            >
              Adicionar
            </button>
          </form>

          <div class="mt-4 space-y-2">
            <div
              v-for="category in categories"
              :key="category.id"
              class="flex items-center justify-between rounded-xl border border-border/70 px-3 py-2 text-sm"
            >
              <span>{{ category.name }}</span>
              <button
                class="text-xs text-muted-foreground hover:text-destructive"
                :disabled="deletingCategoryId === category.id"
                @click="handleDeleteCategory(category.id)"
              >
                Remover
              </button>
            </div>
            <p v-if="categories.length === 0" class="text-sm text-muted-foreground">
              Nenhuma categoria cadastrada.
            </p>
          </div>
        </div>
      </div>
      <div v-if="deleting" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-slate-900/40" @click="closeDelete" />
        <div class="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
          <h2 class="text-lg font-semibold">Excluir produto?</h2>
          <p class="text-sm text-muted-foreground mt-2">
            Tem certeza que deseja excluir "{{ deleting.name }}"? Esta ação não pode ser desfeita.
          </p>
          <div class="flex gap-3 pt-6">
            <button class="flex-1 rounded-xl border border-border px-4 py-2 text-sm" @click="closeDelete">
              Cancelar
            </button>
            <button
              class="flex-1 rounded-xl bg-destructive px-4 py-2 text-sm font-semibold text-destructive-foreground"
              @click="handleDelete"
            >
              Excluir
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Plus, Search } from 'lucide-vue-next'
import ProductCard from '~/components/ProductCard.vue'
import ProductForm from '~/components/ProductForm.vue'
import { useCategories } from '~/composables/useCategories'
import { useProducts, type Product } from '~/composables/useProducts'

const { products, fetchProducts, createProduct, updateProduct, deleteProduct } = useProducts()
const { categories, fetchCategories, createCategory, deleteCategory } = useCategories()
const { success, error } = useToast()

const { pending: productsPending } = await useAsyncData('products', () => fetchProducts())
const { pending: categoriesPending } = await useAsyncData('categories', () => fetchCategories())

const search = ref('')
const formOpen = ref(false)
const categoriesOpen = ref(false)
const newCategoryName = ref('')
const creatingCategory = ref(false)
const deletingCategoryId = ref<string | null>(null)
const editing = ref<Product | null>(null)
const deleting = ref<Product | null>(null)

const loading = computed(() => productsPending.value || categoriesPending.value)

const categoryById = computed(() =>
  new Map(categories.value.map((category) => [category.id, category.name]))
)

const filtered = computed(() =>
  products.value.filter((product) => {
    const categoryLabel = product.categoryId
      ? categoryById.value.get(product.categoryId) || ''
      : ''
    const query = search.value.toLowerCase()
    return (
      product.name.toLowerCase().includes(query) ||
      categoryLabel.toLowerCase().includes(query)
    )
  })
)

const openCreate = () => {
  editing.value = null
  formOpen.value = true
}

const closeCategories = () => {
  categoriesOpen.value = false
  newCategoryName.value = ''
}

const openEdit = (product: Product) => {
  editing.value = product
  formOpen.value = true
}

const handleSubmit = async (data: Omit<Product, 'id' | 'createdAt' | 'updatedAt'>) => {
  try {
    if (editing.value) {
      await updateProduct(editing.value.id, data)
      success('Produto atualizado!')
    } else {
      await createProduct(data)
      success('Produto criado!')
    }
    formOpen.value = false
    editing.value = null
  } catch (error) {
    console.error(error)
  }
}

const handleCreateCategory = async () => {
  if (!newCategoryName.value.trim()) return
  try {
    creatingCategory.value = true
    await createCategory({ name: newCategoryName.value.trim() })
    newCategoryName.value = ''
    success('Categoria criada!')
  } catch (err) {
    console.error(err)
    error('Nao foi possivel criar a categoria')
  } finally {
    creatingCategory.value = false
  }
}

const handleDeleteCategory = async (id: string) => {
  try {
    deletingCategoryId.value = id
    await deleteCategory(id)
    success('Categoria removida')
  } catch (err) {
    console.error(err)
    error('Nao foi possivel remover a categoria')
  } finally {
    deletingCategoryId.value = null
  }
}

const confirmDelete = (product: Product) => {
  deleting.value = product
}

const closeDelete = () => {
  deleting.value = null
}

const handleDelete = async () => {
  if (!deleting.value) return
  try {
    await deleteProduct(deleting.value.id)
    success('Produto removido')
    deleting.value = null
  } catch (err) {
    console.error(err)
    error('Nao foi possivel remover o produto')
  }
}
</script>
