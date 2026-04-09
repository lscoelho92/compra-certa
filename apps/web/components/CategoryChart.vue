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
    <div class="space-y-2 mt-2">
      <div v-for="(item, index) in legendItems" :key="item.name" class="flex items-center justify-between text-sm">
        <div class="flex items-center gap-2">
          <div class="h-3 w-3 rounded-full" :style="{ backgroundColor: colors[index % colors.length] }" />
          <span class="text-muted-foreground">{{ item.name }}</span>
        </div>
        <span class="font-medium">{{ item.value }}</span>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Doughnut } from 'vue-chartjs'
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

const legendItems = computed(() => {
  const grouped: Record<string, number> = {}
  props.purchases
    .filter((purchase) => purchase.month === props.currentMonth)
    .forEach((purchase) => {
      ;(purchase.items ?? []).forEach((item) => {
        const category = item.category || 'Outros'
        grouped[category] = (grouped[category] || 0) + item.quantity
      })
    })

  return Object.entries(grouped)
    .map(([name, value]) => ({ name, value }))
    .sort((a, b) => b.value - a.value)
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
