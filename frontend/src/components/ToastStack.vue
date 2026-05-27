<template>
  <div class="fixed bottom-4 right-4 z-50 pointer-events-none">
    <div v-if="toastStore.toasts.length > 0" class="relative w-72">
      <!-- Ghost card 2 — deepest, inset 8px each side, peeks 16px below -->
      <Transition
        enter-active-class="transition-all duration-200"
        enter-from-class="opacity-0"
        leave-active-class="transition-all duration-150"
        leave-to-class="opacity-0"
      >
        <div
          v-if="depth >= 3"
          class="absolute left-3 right-3 top-0 rounded-lg border bg-card"
          style="height: 100%; transform: translateY(14px); z-index: 1; opacity: 0.5"
        />
      </Transition>

      <!-- Ghost card 1 — middle, inset 4px each side, peeks 8px below -->
      <Transition
        enter-active-class="transition-all duration-200"
        enter-from-class="opacity-0"
        leave-active-class="transition-all duration-150"
        leave-to-class="opacity-0"
      >
        <div
          v-if="depth >= 2"
          class="absolute left-1.5 right-1.5 top-0 rounded-lg border bg-card shadow-sm"
          style="height: 100%; transform: translateY(8px); z-index: 2; opacity: 0.7"
        />
      </Transition>

      <!-- Front card — newest toast, in normal flow so it sizes the container -->
      <Transition
        mode="out-in"
        enter-active-class="transition-all duration-150 ease-out"
        enter-from-class="opacity-0 scale-95"
        leave-active-class="transition-all duration-100 ease-in"
        leave-to-class="opacity-0 scale-95"
      >
        <div
          :key="newest!.id"
          class="pointer-events-auto relative z-[3] flex items-center gap-2.5 rounded-lg border bg-card px-3.5 py-3 shadow-lg text-sm cursor-pointer select-none"
          @click="toastStore.dismiss(newest!.id)"
        >
          <CheckCircle2 v-if="newest!.type === 'success'" class="size-4 shrink-0 text-green-500" />
          <XCircle v-else-if="newest!.type === 'error'" class="size-4 shrink-0 text-destructive" />
          <Info v-else class="size-4 shrink-0 text-blue-500" />
          <span class="flex-1 text-foreground leading-snug">{{ newest!.message }}</span>
          <X class="size-3 shrink-0 opacity-30 hover:opacity-60" />
        </div>
      </Transition>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { CheckCircle2, XCircle, Info, X } from '@lucide/vue'
import { useToastStore } from '@/stores/toast'

const toastStore = useToastStore()

const newest = computed(() => toastStore.toasts.at(-1))
const depth = computed(() => Math.min(toastStore.toasts.length, 3))
</script>
