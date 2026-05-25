<template>
  <div class="p-6 space-y-4">
    <!-- No project selected -->
    <template v-if="!projectStore.current">
      <div class="flex flex-col items-center justify-center py-24 text-center">
        <FolderOpen class="size-8 text-muted-foreground/30 mb-3" />
        <p class="text-sm font-medium">{{ $t('projects.noProject') }}</p>
        <p class="mt-1 text-xs text-muted-foreground max-w-xs">{{ $t('projects.noProjectDescription') }}</p>
      </div>
    </template>

    <!-- Project selected -->
    <template v-else>
      <div class="flex items-center justify-between">
        <div class="space-y-0.5">
          <h1 class="text-base font-semibold">{{ $t('flags.title') }}</h1>
          <p class="text-sm text-muted-foreground">{{ $t('flags.subtitle') }}</p>
        </div>
        <Button size="sm" disabled>
          <Plus class="size-3.5" />
          {{ $t('flags.create') }}
        </Button>
      </div>

      <Separator />

      <div class="flex flex-col items-center justify-center py-24 text-center">
        <Flag class="size-8 text-muted-foreground/30 mb-3" />
        <p class="text-sm font-medium">{{ $t('flags.emptyTitle') }}</p>
        <p class="mt-1 text-xs text-muted-foreground max-w-xs">{{ $t('flags.emptyDescription') }}</p>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { Flag, Plus, FolderOpen } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { useProjectStore } from '@/stores/project'

const projectStore = useProjectStore()
</script>
