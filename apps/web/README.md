# Compra Certa (Nuxt 3)

Frontend migrado de React para Nuxt 3 + TypeScript. Toda a integração Base44 foi removida e substituída por stores locais (localStorage) para manter o fluxo funcionando enquanto o novo backend não estiver pronto.

## Requisitos

- Node.js >= 20.19 ou >= 22.12
- npm

## Setup

```bash
npm install
```

## Desenvolvimento

```bash
npm run dev
```

## Build

```bash
npm run build
npm run preview
```

## Pontos de integração futura

- Produtos: [composables/useProducts.ts](composables/useProducts.ts)
- Compras: [composables/usePurchases.ts](composables/usePurchases.ts)
- Toasts: [composables/useToast.ts](composables/useToast.ts)

Substitua as funções de criação/leitura/remoção desses composables pelo seu cliente de API quando o backend estiver pronto.
