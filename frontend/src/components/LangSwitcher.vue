<template>
  <DropdownMenu>
    <DropdownMenuTrigger as-child>
      <Button variant="ghost" size="sm" class="gap-1.5 text-xs text-muted-foreground">
        <Globe class="size-3.5" />
        {{ currentLabel }}
        <ChevronDown class="size-3" />
      </Button>
    </DropdownMenuTrigger>
    <DropdownMenuContent align="center">
      <DropdownMenuItem
        v-for="loc in SUPPORTED_LOCALES"
        :key="loc"
        class="gap-2"
        @click="setLocale(loc)"
      >
        <Check v-if="locale === loc" class="size-3.5" />
        <span v-else class="size-3.5" />
        {{ $t(`lang.${loc}`) }}
      </DropdownMenuItem>
    </DropdownMenuContent>
  </DropdownMenu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useI18n } from 'vue-i18n'
import { Globe, ChevronDown, Check } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { setLocale, SUPPORTED_LOCALES } from '@/plugins/i18n'

const { t, locale } = useI18n()

const currentLabel = computed(() => t(`lang.${locale.value}`))
</script>
