<template>
  <div v-if="chartData.labels.length === 0" class="flex items-center justify-center h-48 text-muted-foreground text-sm">
    Sem dados este mês
  </div>
  <div v-else>
    <ClientOnly>
      <div class="h-[200px]">
        <Doughnut :data="chartData" :options="chartOptions" />
      </div>
    </ClientOnly>
    <div class="space-y-1 mt-2 max-h-64 overflow-y-auto pr-2">
      <div v-for="(item, index) in legendItems" :key="item.name">
        <button
          class="w-full flex items-center justify-between text-sm py-1.5 px-2 rounded-lg hover:bg-muted/50 transition-colors cursor-pointer"
          @click="toggleCategory(item.name)"
        >
          <div class="flex items-center gap-2">
            <div class="h-3 w-3 shrink-0 rounded-full" :style="{ backgroundColor: colors[index % colors.length] }" />
            <span class="text-muted-foreground">{{ item.name }}</span>
          </div>
          <div class="flex items-center gap-2">
            <span class="font-medium">{{ item.value }}</span>
            <ChevronDown
              class="h-3.5 w-3.5 text-muted-foreground transition-transform duration-200"
              :class="{ 'rotate-180': selectedCategory === item.name }"
            />
          </div>
        </button>
        <div
          v-if="selectedCategory === item.name"
          class="mt-1 mb-1 ml-3.5 space-y-1 border-l-2 pl-3"
          :style="{ borderColor: colors[index % colors.length] }"
        >
          <div
            v-for="product in itemsByCategory[item.name]"
            :key="product.name"
            class="flex items-center justify-between text-xs text-muted-foreground py-0.5"
          >
            <span>{{ product.name }}</span>
            <span class="font-medium tabular-nums">{{ product.quantity }}x</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import { Doughnut } from 'vue-chartjs'
import { ChevronDown } from 'lucide-vue-next'
import {
  Chart as ChartJS,
  ArcElement,
  Tooltip
} from 'chart.js'
import type { Purchase } from '~/composables/usePurchases'

ChartJS.register(ArcElement, Tooltip)

const props = defineProps<{
  purchases: Purchase[]
  currentMonth: string
}>()

const colors = [
  '#16a34a',
  '#38bdf8',
  '#facc15',
  '#f59e0b',
  '#a855f7',
  '#ec4899',
  '#14b8a6',
  '#22c55e',
  '#f472b6',
  '#7c3aed'
]

const selectedCategory = ref<string | null>(null)

const toggleCategory = (name: string) => {
  selectedCategory.value = selectedCategory.value === name ? null : name
}

const monthPurchases = computed(() =>
  props.purchases.filter((purchase) => purchase.month === props.currentMonth)
)

const legendItems = computed(() => {
  const grouped: Record<string, number> = {}
  monthPurchases.value.forEach((purchase) => {
    ;(purchase.items ?? []).forEach((item) => {
      const category = item.category || 'Outros'
      grouped[category] = (grouped[category] || 0) + item.quantity
    })
  })

  return Object.entries(grouped)
    .map(([name, value]) => ({ name, value }))
    .sort((a, b) => b.value - a.value)
})

const itemsByCategory = computed(() => {
  const result: Record<string, { name: string; quantity: number }[]> = {}
  monthPurchases.value.forEach((purchase) => {
    ;(purchase.items ?? []).forEach((item) => {
      const category = item.category || 'Outros'
      if (!result[category]) result[category] = []
      const existing = result[category].find((p) => p.name === item.product_name)
      if (existing) {
        existing.quantity += item.quantity
      } else {
        result[category].push({ name: item.product_name, quantity: item.quantity })
      }
    })
  })
  Object.values(result).forEach((products) => products.sort((a, b) => b.quantity - a.quantity))
  return result
})

const chartData = computed(() => ({
  labels: legendItems.value.map((item) => item.name),
  datasets: [
    {
      data: legendItems.value.map((item) => item.value),
      backgroundColor: legendItems.value.map((_, index) => colors[index % colors.length]),
      borderWidth: 0
    }
  ]
}))

const chartOptions = computed(() => ({
  responsive: true,
  maintainAspectRatio: false,
  plugins: {
    tooltip: {
      callbacks: {
        label(context: { label?: string; parsed: number }) {
          return `${context.label}: ${context.parsed} itens`
        }
      }
    }
  }
}))
</script>
