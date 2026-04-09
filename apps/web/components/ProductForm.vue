<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-slate-900/40" @click="close" />
      <div class="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">
            {{ product ? 'Editar Produto' : 'Novo Produto' }}
          </h2>
          <button class="text-sm text-muted-foreground hover:text-foreground" @click="close">
            Fechar
          </button>
        </div>
        <form class="space-y-4 mt-4" @submit.prevent="handleSubmit">
          <div class="space-y-2">
            <label class="text-sm font-medium" for="name">Nome *</label>
            <input
              id="name"
              v-model.trim="form.name"
              class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
              placeholder="Nome do produto"
              required
            />
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium" for="description">Descrição</label>
            <textarea
              id="description"
              v-model.trim="form.description"
              class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
              placeholder="Descrição opcional"
              rows="2"
            />
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label class="text-sm font-medium">Categoria</label>
              <select
                v-model="form.categoryId"
                class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
              >
                <option value="">Sem categoria</option>
                <option v-for="category in categories" :key="category.id" :value="category.id">
                  {{ category.name }}
                </option>
              </select>
            </div>
            <div class="space-y-2">
              <label class="text-sm font-medium" for="price">Preço (€)</label>
              <input
                id="price"
                v-model="form.defaultPrice"
                type="number"
                step="0.01"
                min="0"
                class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                placeholder="0,00"
              />
            </div>
          </div>
          <div class="flex gap-3 pt-2">
            <button type="button" class="flex-1 rounded-xl border border-border px-4 py-2 text-sm" @click="close">
              Cancelar
            </button>
            <button type="submit" class="flex-1 rounded-xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground">
              {{ product ? 'Salvar' : 'Criar' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import type { Category } from '~/composables/useCategories'
import type { Product } from '~/composables/useProducts'

const props = defineProps<{
  open: boolean
  product?: Product | null
  categories: Category[]
}>()

const emit = defineEmits<{
  (event: 'update:open', value: boolean): void
  (event: 'submit', payload: Omit<Product, 'id' | 'createdAt' | 'updatedAt'>): void
}>()

const form = reactive({
  name: '',
  description: '',
  categoryId: '',
  defaultPrice: '' as string | number
})

const resetForm = () => {
  form.name = ''
  form.description = ''
  form.categoryId = ''
  form.defaultPrice = ''
}

watch(
  () => props.product,
  (value) => {
    if (value) {
      form.name = value.name || ''
      form.description = value.description || ''
      form.categoryId = value.categoryId || ''
      form.defaultPrice = value.defaultPrice ?? ''
      return
    }
    resetForm()
  },
  { immediate: true }
)

watch(
  () => props.open,
  (value) => {
    if (!value && !props.product) resetForm()
  }
)

const close = () => emit('update:open', false)

const handleSubmit = () => {
  const priceValue = form.defaultPrice === '' ? undefined : Number(form.defaultPrice)
  emit('submit', {
    name: form.name,
    description: form.description || undefined,
    categoryId: form.categoryId || undefined,
    defaultPrice: priceValue
  })
}
</script>
