<template>
  <div class="fixed top-4 right-4 z-50 space-y-3">
    <div
      v-for="toast in toasts"
      :key="toast.id"
      class="w-72 rounded-2xl border bg-white px-4 py-3 shadow-xl"
      :class="toast.variant === 'destructive' ? 'border-destructive/40' : 'border-border/60'"
    >
      <div class="flex items-start justify-between gap-2">
        <div>
          <p class="text-sm font-semibold text-foreground">
            {{ toast.title }}
          </p>
          <p v-if="toast.description" class="text-xs text-muted-foreground mt-1">
            {{ toast.description }}
          </p>
        </div>
        <button
          class="text-xs text-muted-foreground hover:text-foreground"
          @click="remove(toast.id)"
        >
          Fechar
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
const { toasts, remove } = useToast()
</script>
