<template>
  <div
    class="group bg-white rounded-2xl p-4 hover:shadow-xl transition-all duration-300 hover:-translate-y-0.5"
    style="box-shadow: 0 4px 20px rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.15);"
  >
    <div class="flex items-start justify-between">
      <div class="flex items-center gap-3 min-w-0">
        <div
          class="h-10 w-10 rounded-xl flex items-center justify-center flex-shrink-0"
          style="background: linear-gradient(135deg, #10b981, #2dd4bf);"
        >
          <Package class="h-5 w-5 text-white" />
        </div>
        <div class="min-w-0">
          <h3 class="font-semibold text-sm truncate">{{ product.name }}</h3>
          <p v-if="product.description" class="text-xs text-muted-foreground truncate mt-0.5">
            {{ product.description }}
          </p>
        </div>
      </div>
      <div class="flex gap-1 opacity-0 group-hover:opacity-100 transition-opacity">
        <button
          class="h-8 w-8 inline-flex items-center justify-center rounded-lg text-muted-foreground hover:text-foreground"
          @click="emit('edit', product)"
        >
          <Pencil class="h-3.5 w-3.5" />
        </button>
        <button
          class="h-8 w-8 inline-flex items-center justify-center rounded-lg text-muted-foreground hover:text-destructive"
          @click="emit('delete', product)"
        >
          <Trash2 class="h-3.5 w-3.5" />
        </button>
      </div>
    </div>
    <div class="flex items-center justify-between mt-3 pt-3" style="border-top: 1px solid rgba(16,185,129,0.12);">
      <span
        class="text-xs font-medium px-2.5 py-1 rounded-full"
        :class="categoryClass"
      >
        {{ categoryLabel || 'Sem categoria' }}
      </span>
      <span v-if="product.defaultPrice != null" class="text-sm font-bold" style="color: #0d9488;">
        € {{ Number(product.defaultPrice).toFixed(2) }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Package, Pencil, Trash2 } from 'lucide-vue-next'
import type { Product } from '~/composables/useProducts'

const props = defineProps<{
  product: Product
  categoryLabel?: string
}>()

const emit = defineEmits<{
  (event: 'edit', product: Product): void
  (event: 'delete', product: Product): void
}>()

const categoryPalette = [
  'bg-sky-100 text-sky-700',
  'bg-indigo-100 text-indigo-700',
  'bg-purple-100 text-purple-700',
  'bg-rose-100 text-rose-700',
  'bg-pink-100 text-pink-700',
  'bg-emerald-100 text-emerald-700',
  'bg-teal-100 text-teal-700',
  'bg-amber-100 text-amber-700',
  'bg-blue-100 text-blue-700',
  'bg-lime-100 text-lime-700',
  'bg-orange-100 text-orange-700',
]

const hashLabel = (value: string) => {
  let hash = 0
  for (let i = 0; i < value.length; i += 1) {
    hash = (hash * 31 + value.charCodeAt(i)) % 2147483647
  }
  return Math.abs(hash)
}

const categoryClass = computed(() => {
  const label = (props.categoryLabel || 'Outros').trim()
  if (!label) return 'bg-slate-100 text-slate-700'
  const index = hashLabel(label) % categoryPalette.length
  return categoryPalette[index]
})
</script>
