<template>
  <div
    class="bg-white rounded-2xl p-5 transition-all hover:shadow-lg"
    :class="className"
    style="box-shadow: 0 4px 20px rgba(16,185,129,0.1); border: 1px solid rgba(16,185,129,0.15);"
  >
    <div class="flex items-center gap-3">
      <div
        class="h-10 w-10 rounded-xl flex items-center justify-center"
        style="background: linear-gradient(135deg, #10b981, #2dd4bf);"
      >
        <component :is="icon" class="h-5 w-5 text-white" />
      </div>
      <div>
        <p class="text-xs font-medium text-muted-foreground uppercase tracking-wider">{{ label }}</p>
        <p class="text-2xl font-bold mt-0.5" style="color: hsl(200, 30%, 12%);">{{ value }}</p>
        <p v-if="subtext" class="text-xs text-muted-foreground mt-0.5">{{ subtext }}</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  icon: Component
  label: string
  value: string | number
  subtext?: string
  className?: string
}>()
</script>
