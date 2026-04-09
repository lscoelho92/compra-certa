<template>
  <div v-if="chartData.labels.length === 0" class="flex items-center justify-center h-64 text-muted-foreground text-sm">
    Nenhuma compra registrada ainda
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
}>()

const monthNames: Record<string, string> = {
  '01': 'Jan', '02': 'Fev', '03': 'Mar', '04': 'Abr',
  '05': 'Mai', '06': 'Jun', '07': 'Jul', '08': 'Ago',
  '09': 'Set', '10': 'Out', '11': 'Nov', '12': 'Dez'
}

const chartData = computed(() => {
  const grouped: Record<string, { quantidade: number; total: number }> = {}
  props.purchases.forEach((purchase) => {
    if (!purchase.month) return
    if (!grouped[purchase.month]) grouped[purchase.month] = { quantidade: 0, total: 0 }
    grouped[purchase.month].quantidade += (purchase.items ?? []).reduce((sum, item) => sum + item.quantity, 0)
    grouped[purchase.month].total += purchase.total_price
  })

  const items = Object.entries(grouped)
    .sort(([a], [b]) => a.localeCompare(b))
    .slice(-12)

  const labels = items.map(([month]) => {
    const [year, monthValue] = month.split('-')
    return `${monthNames[monthValue]}/${year.slice(2)}`
  })

  return {
    labels,
    datasets: [
      {
        label: 'Quantidade',
        data: items.map(([, value]) => value.quantidade),
        backgroundColor: '#ef4444',
        borderRadius: 6,
        yAxisID: 'y'
      },
      {
        label: 'Total (€)',
        data: items.map(([, value]) => Math.round(value.total * 100) / 100),
        backgroundColor: '#3b82f6',
        borderRadius: 6,
        yAxisID: 'y1'
      }
    ]
  }
})

const chartOptions = computed<ChartOptions<'bar'>>(() => ({
  responsive: true,
  maintainAspectRatio: false,
  scales: {
    y: {
      position: 'left',
      ticks: {
        color: 'hsl(200, 15%, 45%)'
      },
      grid: {
        color: 'hsl(168, 20%, 86%)'
      }
    },
    y1: {
      position: 'right',
      ticks: {
        color: 'hsl(200, 15%, 45%)'
      },
      grid: {
        drawOnChartArea: false
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
      labels: {
        color: 'hsl(200, 30%, 12%)'
      }
    },
    tooltip: {
      callbacks: {
        label(context: TooltipItem<'bar'>) {
          if (context.dataset.label === 'Total (€)') {
            return `Total (€): € ${Number(context.parsed.y).toFixed(2)}`
          }
          return `Quantidade: ${context.parsed.y}`
        }
      }
    }
  }
}))
</script>
