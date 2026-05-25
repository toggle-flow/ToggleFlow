<template>
  <RouterView />
</template>

<script setup lang="ts">
import { useThemeStore } from '@/stores/theme'
useThemeStore() // initializes on app mount — applies dark class before first render
</script>
