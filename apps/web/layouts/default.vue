<template>
  <div class="min-h-screen bg-background">
    <aside
      class="hidden md:flex fixed left-0 top-0 bottom-0 w-64 flex-col z-30 overflow-hidden"
      style="background: hsl(185, 40%, 93%);"
    >
      <div
        class="absolute left-0 top-0 bottom-0 w-1 rounded-full"
        style="background: linear-gradient(to bottom, #10b981, #2dd4bf);"
      />
      <div class="p-6 pl-8">
        <h1 class="text-xl font-bold tracking-tight" style="color: hsl(200, 30%, 12%);">
          Compra<span style="color: #0d9488;">Certa</span>
        </h1>
        <p class="text-xs mt-1" style="color: hsl(200, 15%, 50%);">Controle de compras</p>
      </div>
      <nav class="flex-1 px-4 space-y-1 relative z-10">
        <NuxtLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="flex items-center gap-3 px-4 py-3 rounded-2xl text-sm font-medium transition-colors duration-200"
          :class="isActive(item.path) ? 'text-foreground' : 'text-muted-foreground hover:text-foreground hover:bg-emerald-50/80'"
          :style="isActive(item.path) ? { background: 'rgba(16, 185, 129, 0.18)' } : undefined"
        >
          <component :is="item.icon" class="h-4 w-4" />
          <span>{{ item.label }}</span>
        </NuxtLink>
      </nav>
    </aside>

    <nav
      class="md:hidden fixed bottom-0 left-0 right-0 border-t z-30 px-2 pb-safe"
      style="background: hsl(185, 40%, 93%); border-color: hsl(168, 20%, 86%);"
    >
      <div class="flex justify-around py-2">
        <NuxtLink
          v-for="item in navItems"
          :key="item.path"
          :to="item.path"
          class="flex flex-col items-center gap-1 px-4 py-2 rounded-xl text-xs font-medium transition-all"
          :class="isActive(item.path) ? 'text-primary' : 'text-muted-foreground'"
        >
          <component
            :is="item.icon"
            class="h-5 w-5"
            :class="isActive(item.path) ? 'stroke-[2.5]' : ''"
          />
          {{ item.label }}
        </NuxtLink>
      </div>
    </nav>

    <main class="md:ml-64 min-h-screen pb-20 md:pb-0">
      <div class="max-w-6xl mx-auto p-4 md:p-8">
        <slot />
      </div>
    </main>
  </div>
</template>

<script setup lang="ts">
import { LayoutDashboard, Package, ShoppingCart } from 'lucide-vue-next'

const route = useRoute()

const navItems = [
  { path: '/', label: 'Painel', icon: LayoutDashboard },
  { path: '/products', label: 'Produtos', icon: Package },
  { path: '/shopping', label: 'Compras', icon: ShoppingCart }
]

const isActive = (path: string) => route.path === path
</script>
