<template>
  <button
    class="group inline-flex items-center gap-1 rounded bg-muted px-1.5 py-0.5 font-mono text-[11px] text-muted-foreground transition-colors hover:bg-muted/70 hover:text-foreground"
    @click.stop="copy"
  >
    {{ value }}
    <Check v-if="copied" class="size-2.5 text-green-500 shrink-0" />
    <Copy v-else class="size-2.5 shrink-0 opacity-40" />
  </button>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { Copy, Check } from '@lucide/vue'

const props = defineProps<{ value: string }>()

const copied = ref(false)

function copy() {
  navigator.clipboard.writeText(props.value)
  copied.value = true
  setTimeout(() => {
    copied.value = false
  }, 2000)
}
</script>
