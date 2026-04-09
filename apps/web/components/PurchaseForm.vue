<template>
  <Teleport to="body">
    <div v-if="open" class="fixed inset-0 z-50 flex items-center justify-center">
      <div class="absolute inset-0 bg-slate-900/40" @click="close" />
      <div class="relative w-full max-w-2xl max-h-[90vh] rounded-2xl bg-white p-6 shadow-2xl flex flex-col">
        <div class="flex items-center justify-between">
          <h2 class="text-lg font-semibold">{{ isEditing ? 'Editar Compra' : 'Registrar Compra' }}</h2>
          <button class="text-sm text-muted-foreground hover:text-foreground" @click="close">
            Fechar
          </button>
        </div>
        <form class="mt-4 flex min-h-0 flex-1 flex-col space-y-4" @submit.prevent="handleSubmit">
          <div class="rounded-2xl border border-border/70 p-4">
            <div class="grid grid-cols-1 gap-4 sm:grid-cols-[2fr_1fr_1fr]">
              <div class="space-y-2">
                <label class="text-sm font-medium">Produto *</label>
                <ProductSearchSelect
                  :model-value="draft.product_id"
                  :products="products"
                  placeholder="Selecione um produto"
                  @update:modelValue="handleDraftProductChange"
                />
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium" for="draft-qty">Quantidade *</label>
                <input
                  id="draft-qty"
                  v-model.number="draft.quantity"
                  type="number"
                  min="1"
                  class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
              </div>
              <div class="space-y-2">
                <label class="text-sm font-medium" for="draft-up">Preço Unit. (€)</label>
                <input
                  id="draft-up"
                  :value="draft.unit_price.toFixed(2)"
                  type="number"
                  step="0.01"
                  min="0"
                  readonly
                  class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
              </div>
            </div>
            <button
              type="button"
              class="mt-4 w-full rounded-xl border border-dashed border-border px-4 py-2 text-sm text-muted-foreground hover:border-primary hover:text-primary disabled:cursor-not-allowed disabled:opacity-60"
              @click="addItem"
              :disabled="!canAddItem"
            >
              Adicionar item
            </button>
          </div>

          <div class="space-y-2 overflow-auto pr-1" style="max-height: 40vh;">
            <div
              v-if="form.items.length === 0"
              class="rounded-2xl border border-dashed border-border px-4 py-3 text-sm text-muted-foreground"
            >
              Nenhum item adicionado ainda.
            </div>
            <div
              v-for="(item, index) in form.items"
              :key="item.key"
              class="flex flex-wrap items-center justify-between gap-3 rounded-2xl border border-border/70 px-4 py-3"
            >
              <div class="min-w-0 flex-1">
                <p class="text-sm font-medium">{{ item.product_name }}</p>
                <p class="text-xs text-muted-foreground">€ {{ item.unit_price.toFixed(2) }}</p>
              </div>
              <div class="flex items-center gap-3">
                <input
                  :id="`qty-${index}`"
                  v-model.number="item.quantity"
                  type="number"
                  min="1"
                  class="w-24 rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
                />
                <span class="text-sm font-semibold">€ {{ (item.unit_price * item.quantity).toFixed(2) }}</span>
                <button
                  type="button"
                  class="text-xs text-muted-foreground hover:text-destructive"
                  @click="removeItem(index)"
                >
                  Remover
                </button>
              </div>
            </div>
          </div>
          <div class="space-y-2">
            <label class="text-sm font-medium" for="date">Data da Compra *</label>
            <input
              id="date"
              v-model="form.purchase_date"
              type="date"
              class="w-full rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
              required
            />
          </div>
          <div class="flex gap-3 pt-2">
            <button type="button" class="flex-1 rounded-xl border border-border px-4 py-2 text-sm" @click="close">
              Cancelar
            </button>
            <button
              type="submit"
              class="flex-1 rounded-xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground"
              :disabled="!canSubmit"
            >
              {{ isEditing ? 'Salvar' : 'Registrar' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { computed, reactive, watch } from 'vue'
import { format } from 'date-fns'
import ProductSearchSelect from '~/components/ProductSearchSelect.vue'
import type { Product } from '~/composables/useProducts'
import type { Purchase, PurchasePayload } from '~/composables/usePurchases'

const props = defineProps<{
  open: boolean
  products: Product[]
  purchase?: Purchase | null
}>()

const emit = defineEmits<{
  (event: 'update:open', value: boolean): void
  (event: 'submit', payload: PurchasePayload & { id?: string }): void
}>()

type FormItem = {
  key: string
  product_id: string
  product_name: string
  quantity: number
  unit_price: number
}

const createItem = (): FormItem => ({
  key: crypto.randomUUID(),
  product_id: '',
  product_name: '',
  quantity: 1,
  unit_price: 0
})

const form = reactive({
  purchase_date: format(new Date(), 'yyyy-MM-dd'),
  items: [createItem()]
})

const draft = reactive({
  product_id: '',
  product_name: '',
  quantity: 1,
  unit_price: 0
})

const isEditing = computed(() => Boolean(props.purchase))

const canSubmit = computed(() => {
  if (!form.purchase_date) return false
  if (form.items.length === 0) return false
  return form.items.every((item) => item.product_id && item.quantity > 0 && item.unit_price >= 0)
})

const canAddItem = computed(() => {
  return Boolean(draft.product_id) && draft.quantity > 0
})

const resetForm = () => {
  form.purchase_date = format(new Date(), 'yyyy-MM-dd')
  form.items = []
  draft.product_id = ''
  draft.product_name = ''
  draft.quantity = 1
  draft.unit_price = 0
}

const fillFromPurchase = (purchase: Purchase) => {
  form.purchase_date = purchase.purchase_date
  form.items = purchase.items.map((item) => ({
    key: crypto.randomUUID(),
    product_id: item.product_id,
    product_name: item.product_name,
    quantity: item.quantity,
    unit_price: item.unit_price
  }))
  draft.product_id = ''
  draft.product_name = ''
  draft.quantity = 1
  draft.unit_price = 0
}

watch(
  () => props.open,
  (value) => {
    if (!value) return
    if (props.purchase) {
      fillFromPurchase(props.purchase)
    } else {
      resetForm()
    }
  }
)

watch(
  () => props.purchase,
  (value) => {
    if (!props.open) return
    if (value) {
      fillFromPurchase(value)
    } else {
      resetForm()
    }
  }
)

const close = () => emit('update:open', false)

const handleDraftProductChange = (value: string) => {
  draft.product_id = value
  const selected = props.products.find((product) => product.id === value)
  draft.product_name = selected?.name ?? ''
  draft.unit_price = selected?.defaultPrice ?? 0
}

const addItem = () => {
  if (!canAddItem.value) return
  const existing = form.items.find((item) => item.product_id === draft.product_id)
  if (existing) {
    existing.quantity += draft.quantity
  } else {
    form.items.push({
      key: crypto.randomUUID(),
      product_id: draft.product_id,
      product_name: draft.product_name,
      quantity: draft.quantity,
      unit_price: draft.unit_price
    })
  }
  draft.product_id = ''
  draft.product_name = ''
  draft.quantity = 1
  draft.unit_price = 0
}

const removeItem = (index: number) => {
  form.items.splice(index, 1)
}

const handleSubmit = () => {
  const payload: PurchasePayload & { id?: string } = {
    id: props.purchase?.id,
    purchase_date: form.purchase_date,
    items: form.items.map((item) => ({
      product_id: item.product_id,
      quantity: Number(item.quantity),
      unit_price: Number(item.unit_price)
    }))
  }
  emit('submit', payload)
}
</script>
