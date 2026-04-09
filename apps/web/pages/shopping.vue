<template>
  <div class="space-y-6">
    <div class="flex flex-col sm:flex-row sm:items-center justify-between gap-4">
      <div>
        <h1 class="text-2xl md:text-3xl font-bold tracking-tight">Compras</h1>
        <p class="text-muted-foreground text-sm mt-1">{{ purchases.length }} compras registradas</p>
      </div>
      <button
        class="inline-flex items-center gap-2 rounded-2xl bg-primary px-4 py-2 text-sm font-semibold text-primary-foreground"
        :disabled="products.length === 0"
        @click="openCreate"
      >
        <Plus class="h-4 w-4" /> Registrar Compra
      </button>
    </div>

    <div v-if="products.length === 0" class="bg-amber-50 border border-amber-200 text-amber-800 rounded-lg p-4 text-sm">
      Cadastre produtos primeiro antes de registrar compras.
    </div>

    <div v-if="months.length > 0" class="flex items-center gap-3">
      <Calendar class="h-4 w-4 text-muted-foreground" />
      <select
        v-model="monthFilter"
        class="rounded-xl border border-input bg-transparent px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-primary/30"
      >
        <option value="all">Todos os meses</option>
        <option v-for="month in months" :key="month" :value="month">
          {{ formatMonthLabel(month) }}
        </option>
      </select>
    </div>

    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="w-8 h-8 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>

    <template v-else>
      <div v-if="filteredPurchases.length === 0" class="text-center py-16">
        <p class="text-muted-foreground">Nenhuma compra registrada</p>
      </div>
      <div v-else class="bg-white rounded-2xl overflow-hidden" style="box-shadow: 0 4px 20px rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.15);">
        <div class="overflow-x-auto">
          <table class="w-full">
            <thead>
              <tr class="border-b border-border bg-muted/50">
                <th class="text-left text-xs font-medium text-muted-foreground uppercase tracking-wider px-5 py-3">Data</th>
                <th class="text-center text-xs font-medium text-muted-foreground uppercase tracking-wider px-5 py-3">Qtd Produtos</th>
                <th class="text-center text-xs font-medium text-muted-foreground uppercase tracking-wider px-5 py-3">Valor Total</th>
                <th class="px-5 py-3 w-10"></th>
              </tr>
            </thead>
            <tbody>
              <tr
                v-for="purchase in filteredPurchases"
                :key="purchase.id"
                class="border-b border-border/50 hover:bg-muted/30 transition-colors"
              >
                <td class="px-5 py-3.5 text-sm text-muted-foreground">
                  {{ formatDate(purchase.purchase_date) }}
                </td>
                <td class="px-5 py-3.5 text-sm text-center">
                  {{ purchase.items.length }}
                </td>
                <td class="px-5 py-3.5 text-sm font-semibold text-center">
                  € {{ purchase.total_price.toFixed(2) }}
                </td>
                <td class="px-5 py-3.5">
                  <button
                    class="h-8 w-8 inline-flex items-center justify-center rounded-lg text-muted-foreground hover:text-primary"
                    @click="startEdit(purchase)"
                  >
                    <Pencil class="h-3.5 w-3.5" />
                  </button>
                  <button
                    class="h-8 w-8 inline-flex items-center justify-center rounded-lg text-muted-foreground hover:text-destructive"
                    @click="confirmDelete(purchase)"
                  >
                    <Trash2 class="h-3.5 w-3.5" />
                  </button>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </template>

    <PurchaseForm
      v-model:open="formOpen"
      :products="products"
      :purchase="editing"
      @submit="handleSubmit"
    />

    <Teleport to="body">
      <div v-if="deleting" class="fixed inset-0 z-50 flex items-center justify-center">
        <div class="absolute inset-0 bg-slate-900/40" @click="closeDelete" />
        <div class="relative w-full max-w-md rounded-2xl bg-white p-6 shadow-2xl">
          <h2 class="text-lg font-semibold">Excluir compra?</h2>
          <p class="text-sm text-muted-foreground mt-2">
            Tem certeza que deseja excluir esta compra com {{ deleting.items.length }} itens?
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
import { format } from 'date-fns'
import { Plus, Trash2, Calendar, Pencil } from 'lucide-vue-next'
import PurchaseForm from '~/components/PurchaseForm.vue'
import { useProducts } from '~/composables/useProducts'
import { usePurchases, type Purchase, type PurchasePayload } from '~/composables/usePurchases'

const { products, fetchProducts } = useProducts()
const { purchases, fetchPurchases, createPurchase, updatePurchase, deletePurchase } = usePurchases()
const { success, error } = useToast()

const { pending: productsPending } = await useAsyncData('products', () => fetchProducts())
const { pending: purchasesPending } = await useAsyncData(
  'purchases',
  () => fetchPurchases(),
  { server: false }
)

const formOpen = ref(false)
const deleting = ref<Purchase | null>(null)
const editing = ref<Purchase | null>(null)
const monthFilter = ref('all')

const loading = computed(() => productsPending.value || purchasesPending.value)

const months = computed(() => {
  const set = new Set(purchases.value.map((purchase) => purchase.month))
  return [...set].sort().reverse()
})

const filteredPurchases = computed(() => {
  if (monthFilter.value === 'all') return purchases.value
  return purchases.value.filter((purchase) => purchase.month === monthFilter.value)
})


const formatMonthLabel = (value: string) => {
  const [year, month] = value.split('-')
  const monthNames: Record<string, string> = {
    '01': 'Janeiro', '02': 'Fevereiro', '03': 'Março', '04': 'Abril',
    '05': 'Maio', '06': 'Junho', '07': 'Julho', '08': 'Agosto',
    '09': 'Setembro', '10': 'Outubro', '11': 'Novembro', '12': 'Dezembro'
  }
  return `${monthNames[month]} ${year}`
}

const formatDate = (value: string) => {
  if (!value) return '—'
  return format(new Date(`${value}T12:00:00`), 'dd/MM/yyyy')
}

const openCreate = () => {
  editing.value = null
  formOpen.value = true
}

const startEdit = (purchase: Purchase) => {
  editing.value = purchase
  formOpen.value = true
}

const handleSubmit = async (payload: PurchasePayload & { id?: string }) => {
  try {
    if (payload.id) {
      await updatePurchase(payload.id, payload)
    } else {
      await createPurchase(payload)
    }
    editing.value = null
    formOpen.value = false
    success(payload.id ? 'Compra atualizada!' : 'Compra registrada!')
  } catch (err) {
    error('Nao foi possivel salvar a compra')
  }
}

const confirmDelete = (purchase: Purchase) => {
  deleting.value = purchase
}

const closeDelete = () => {
  deleting.value = null
}

const handleDelete = async () => {
  if (!deleting.value) return
  try {
    await deletePurchase(deleting.value.id)
    deleting.value = null
    success('Compra removida')
  } catch (err) {
    error('Nao foi possivel remover a compra')
  }
}
</script>
