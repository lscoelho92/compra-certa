<template>
  <div class="space-y-8">
    <div>
      <h1 class="text-2xl md:text-3xl font-bold tracking-tight">Painel</h1>
      <p class="text-muted-foreground text-sm mt-1">Resumo de {{ monthLabel }}</p>
    </div>

    <div v-if="loading" class="flex items-center justify-center h-64">
      <div class="w-8 h-8 border-4 border-primary/20 border-t-primary rounded-full animate-spin" />
    </div>

    <template v-else>
      <div class="grid grid-cols-2 lg:grid-cols-3 gap-5">
        <StatCard
          :icon="Package"
          label="Produtos Cadastrados"
          :value="products.length"
        />
        <StatCard
          :icon="ShoppingCart"
          label="Compras do Mês"
          :value="stats.thisMonthPurchases"
          :subtext="stats.countDiff !== null ? `${stats.countDiff > 0 ? '+' : ''}${stats.countDiff}% vs mês anterior` : ''"
        />
        <StatCard
          :icon="DollarSign"
          label="Gasto do Mês"
          :value="`€ ${stats.totalSpent.toFixed(2)}`"
          :subtext="stats.spentDiff !== null ? `${stats.spentDiff > 0 ? '+' : ''}${stats.spentDiff}% vs mês anterior` : ''"
        />
      </div>

      <div class="grid lg:grid-cols-3 gap-6">
        <div
          class="lg:col-span-2 bg-white rounded-2xl p-5"
          style="box-shadow: 0 4px 20px rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.15);"
        >
          <h2 class="text-base font-semibold mb-4">Compras por Mês</h2>
          <MonthlyChart :purchases="purchases" />
        </div>
        <div
          class="bg-white rounded-2xl p-5"
          style="box-shadow: 0 4px 20px rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.15);"
        >
          <h2 class="text-base font-semibold mb-4">Por Categoria</h2>
          <CategoryChart :purchases="purchases" :current-month="currentMonth" />
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { format } from 'date-fns'
import { Package, ShoppingCart, DollarSign, TrendingUp } from 'lucide-vue-next'
import { useProducts } from '~/composables/useProducts'
import { usePurchases } from '~/composables/usePurchases'

const { products, fetchProducts } = useProducts()
const { purchases, fetchPurchases } = usePurchases()

const { pending: productsPending } = await useAsyncData('products', () => fetchProducts())
const { pending: purchasesPending } = await useAsyncData(
  'purchases',
  () => fetchPurchases(),
  { server: false }
)

const loading = computed(() => productsPending.value || purchasesPending.value)
const currentMonth = format(new Date(), 'yyyy-MM')

const stats = computed(() => {
  const thisMonth = purchases.value.filter((purchase) => purchase.month === currentMonth)
  const totalQty = thisMonth.reduce(
    (sum, purchase) => sum + (purchase.items ?? []).reduce((itemSum, item) => itemSum + item.quantity, 0),
    0
  )
  const totalSpent = thisMonth.reduce((sum, purchase) => sum + purchase.total_price, 0)

  const date = new Date()
  date.setMonth(date.getMonth() - 1)
  const prevMonth = format(date, 'yyyy-MM')
  const prevData = purchases.value.filter((purchase) => purchase.month === prevMonth)
  const prevQty = prevData.reduce(
    (sum, purchase) => sum + (purchase.items ?? []).reduce((itemSum, item) => itemSum + item.quantity, 0),
    0
  )
  const prevSpent = prevData.reduce((sum, purchase) => sum + purchase.total_price, 0)
  const prevCount = prevData.length

  const qtyDiff = prevQty > 0 ? Math.round(((totalQty - prevQty) / prevQty) * 100) : null
  const spentDiff = prevSpent > 0 ? Math.round(((totalSpent - prevSpent) / prevSpent) * 100) : null
  const countDiff = prevCount > 0
    ? Math.round(((thisMonth.length - prevCount) / prevCount) * 100)
    : null

  return {
    totalQty,
    totalSpent,
    qtyDiff,
    spentDiff,
    countDiff,
    thisMonthPurchases: thisMonth.length
  }
})

const monthLabel = new Date().toLocaleDateString('pt-BR', { month: 'long', year: 'numeric' })
</script>
