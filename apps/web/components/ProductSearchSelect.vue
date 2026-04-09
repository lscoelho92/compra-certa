<template>
  <div ref="root" class="relative">
    <input
      v-model="query"
      type="text"
      class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
      :placeholder="placeholder"
      :disabled="disabled"
      @focus="openDropdown"
      @input="handleInput"
      @keydown.down.prevent="moveHighlight(1)"
      @keydown.up.prevent="moveHighlight(-1)"
      @keydown.enter.prevent="chooseHighlighted"
      @keydown.escape="closeDropdown"
    />

    <div
      v-if="open"
      class="absolute z-20 mt-2 w-full overflow-hidden rounded-xl border border-border bg-white shadow-xl"
    >
      <div class="max-h-56 overflow-auto">
        <button
          v-for="(product, index) in filteredProducts"
          :key="product.id"
          type="button"
          class="flex w-full items-center justify-between px-3 py-2 text-left text-sm transition-colors"
          :class="index === highlightedIndex ? 'bg-muted/70 text-foreground' : 'hover:bg-muted/40'"
          @click="selectProduct(product)"
        >
          <span>{{ product.name }}</span>
        </button>
        <div v-if="filteredProducts.length === 0" class="px-3 py-2 text-sm text-muted-foreground">
          Nenhum produto encontrado
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, nextTick, onBeforeUnmount, onMounted, ref, watch } from 'vue'
import type { Product } from '~/composables/useProducts'

const props = defineProps<{
  modelValue: string
  products: Product[]
  placeholder?: string
  disabled?: boolean
}>()

const emit = defineEmits<{
  (event: 'update:modelValue', value: string): void
}>()

const root = ref<HTMLElement | null>(null)
const open = ref(false)
const query = ref('')
const highlightedIndex = ref(0)
const suppressModelSync = ref(false)

const selectedProduct = computed(() =>
  props.products.find((product) => product.id === props.modelValue) || null
)

const filteredProducts = computed(() => {
  const term = query.value.trim().toLowerCase()
  if (!term) return props.products
  return props.products.filter((product) => product.name.toLowerCase().includes(term))
})

const syncQueryFromModel = () => {
  if (suppressModelSync.value) return
  const selected = selectedProduct.value
  query.value = selected ? selected.name : ''
}

watch(() => props.modelValue, syncQueryFromModel, { immediate: true })

watch(query, () => {
  highlightedIndex.value = 0
})

const openDropdown = () => {
  if (props.disabled) return
  open.value = true
}

const closeDropdown = () => {
  open.value = false
}

const handleInput = () => {
  if (props.disabled) return
  open.value = true
  suppressModelSync.value = true
  if (selectedProduct.value && query.value !== selectedProduct.value.name) {
    emit('update:modelValue', '')
  }
  nextTick(() => {
    suppressModelSync.value = false
  })
}

const selectProduct = (product: Product) => {
  emit('update:modelValue', product.id)
  query.value = product.name
  open.value = false
}

const moveHighlight = (step: number) => {
  if (!open.value) {
    open.value = true
    return
  }
  const total = filteredProducts.value.length
  if (total === 0) return
  highlightedIndex.value = (highlightedIndex.value + step + total) % total
}

const chooseHighlighted = () => {
  if (!open.value) return
  const product = filteredProducts.value[highlightedIndex.value]
  if (product) selectProduct(product)
}

const handleClickOutside = (event: MouseEvent) => {
  if (!root.value) return
  if (!root.value.contains(event.target as Node)) closeDropdown()
}

onMounted(() => {
  document.addEventListener('mousedown', handleClickOutside)
})

onBeforeUnmount(() => {
  document.removeEventListener('mousedown', handleClickOutside)
})
</script>
