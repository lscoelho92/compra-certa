<template>
  <div v-if="chartData.labels.length === 0" class="flex items-center justify-center h-64 text-muted-foreground text-sm">
    Nenhuma compra registrada neste mês
  </div>
  <ClientOnly v-else>
    <div class="h-[300px]">
      <Bar :data="chartData" :options="chartOptions" />
    </div>
  </ClientOnly>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { Bar } from 'vue-chartjs'
import {
  Chart as ChartJS,
  BarElement,
  CategoryScale,
  LinearScale,
  Tooltip,
  Legend
} from 'chart.js'
import type { ChartOptions, TooltipItem } from 'chart.js'
import type { Purchase } from '~/composables/usePurchases'

ChartJS.register(BarElement, CategoryScale, LinearScale, Tooltip, Legend)

const props = defineProps<{
  purchases: Purchase[]
  currentMonth: string
}>()

const chartData = computed(() => {
  const monthPurchases = props.purchases.filter((p) => p.month === props.currentMonth)

  const quantities: Record<string, number> = {}
  monthPurchases.forEach((purchase) => {
    ;(purchase.items ?? []).forEach((item) => {
      quantities[item.product_name] = (quantities[item.product_name] ?? 0) + item.quantity
    })
  })

  const entries = Object.entries(quantities).sort(([, a], [, b]) => b - a).slice(0, 10)

  const colors = [
    '#10b981', '#3b82f6', '#f59e0b', '#ef4444', '#8b5cf6',
    '#ec4899', '#14b8a6', '#f97316', '#6366f1', '#84cc16'
  ]

  return {
    labels: entries.map(([name]) => name),
    datasets: [
      {
        label: 'Quantidade',
        data: entries.map(([, qty]) => qty),
        backgroundColor: entries.map((_, i) => colors[i % colors.length]),
        borderRadius: 6
      }
    ]
  }
})

const chartOptions = computed<ChartOptions<'bar'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      ticks: {
        color: 'hsl(200, 15%, 45%)',
        stepSize: 1
      },
      grid: {
        color: 'hsl(168, 20%, 86%)'
      }
    },
    x: {
      ticks: {
        color: 'hsl(200, 15%, 45%)'
      },
      grid: {
        color: 'hsl(168, 20%, 86%)'
      }
    }
  },
  plugins: {
    legend: {
      display: false
    },
    tooltip: {
      callbacks: {
        label(context: TooltipItem<'bar'>) {
          return `Quantidade: ${context.parsed.y}`
        }
      }
    }
  }
}))
</script>
