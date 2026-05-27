<template>
  <div class="flex h-full flex-col overflow-hidden">
    <!-- No project selected -->
    <div
      v-if="!projectStore.current"
      class="flex flex-1 flex-col items-center justify-center text-center"
    >
      <FolderOpen class="size-8 text-muted-foreground/30 mb-3" />
      <p class="text-sm font-medium">{{ $t('projects.noProject') }}</p>
      <p class="mt-1 text-xs text-muted-foreground max-w-xs">
        {{ $t('projects.noProjectDescription') }}
      </p>
    </div>

    <template v-else>
      <!-- Fixed header -->
      <div class="shrink-0 space-y-4 px-6 pt-6 pb-4">
        <div class="space-y-0.5">
          <h1 class="text-base font-semibold">{{ $t('audit.title') }}</h1>
          <p class="text-sm text-muted-foreground">{{ $t('audit.subtitle') }}</p>
        </div>
        <Separator />
      </div>

      <!-- Scrollable body + pinned pagination -->
      <div class="flex min-h-0 flex-1 flex-col overflow-hidden px-6 pb-6">
        <div v-if="loading" class="flex flex-1 items-center justify-center">
          <Loader2 class="size-6 animate-spin text-muted-foreground/40" />
        </div>

        <div
          v-else-if="total === 0"
          class="flex flex-1 flex-col items-center justify-center text-center"
        >
          <ClipboardList class="size-8 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-medium">{{ $t('audit.emptyTitle') }}</p>
          <p class="mt-1 text-xs text-muted-foreground max-w-xs">
            {{ $t('audit.emptyDescription') }}
          </p>
        </div>

        <template v-else>
          <div class="section-audit min-h-0 flex-1 overflow-y-auto space-y-2">
            <div
              v-for="entry in entries"
              :key="entry.id"
              class="rounded-lg border bg-card px-4 py-3"
            >
              <div class="flex items-start justify-between gap-4">
                <div class="flex items-start gap-3 min-w-0">
                  <!-- Actor avatar -->
                  <div
                    class="flex size-7 shrink-0 items-center justify-center rounded-full bg-muted text-[11px] font-semibold uppercase mt-0.5"
                  >
                    {{ entry.actor[0] }}
                  </div>

                  <div class="min-w-0">
                    <!-- Actor + action badge + resource -->
                    <div class="flex flex-wrap items-center gap-2">
                      <span class="text-sm font-medium">{{ entry.actor }}</span>
                      <span
                        :class="actionBadgeClass(entry.action)"
                        class="inline-flex items-center rounded px-1.5 py-0.5 text-[10px] font-medium"
                      >
                        {{ actionLabel(entry.action) }}
                      </span>
                      <span class="font-mono text-[11px] text-muted-foreground">{{
                        entry.resource
                      }}</span>
                    </div>

                    <!-- Diff line -->
                    <div v-if="diffLines(entry).length" class="mt-1.5 space-y-0.5">
                      <div
                        v-for="(line, i) in diffLines(entry)"
                        :key="i"
                        class="text-xs text-muted-foreground font-mono"
                      >
                        {{ line }}
                      </div>
                    </div>
                  </div>
                </div>

                <span class="shrink-0 text-xs text-muted-foreground whitespace-nowrap">
                  {{ timeAgo(entry.created_at) }}
                </span>
              </div>
            </div>
          </div>

          <div class="shrink-0 pt-4">
            <Pagination :page="page" :total="total" :limit="limit" @change="goTo" />
          </div>
        </template>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { ClipboardList, FolderOpen, Loader2 } from '@lucide/vue'
import { Separator } from '@/components/ui/separator'
import Pagination from '@/components/ui/pagination/Pagination.vue'
import { useProjectStore } from '@/stores/project'
import { auditApi, type AuditEntry } from '@/api/audit'
import { timeAgo } from '@/lib/utils'

const LIMIT = 30

const projectStore = useProjectStore()
const entries = ref<AuditEntry[]>([])
const total = ref(0)
const page = ref(1)
const limit = ref(LIMIT)
const loading = ref(false)

async function load() {
  if (!projectStore.current) return
  loading.value = true
  try {
    const res = await auditApi.list(projectStore.current.id, {
      limit: limit.value,
      offset: (page.value - 1) * limit.value,
    })
    entries.value = res.data ?? []
    total.value = res.total
  } finally {
    loading.value = false
  }
}

watch(
  () => projectStore.current,
  () => {
    page.value = 1
    load()
  },
  { immediate: true }
)

function goTo(p: number) {
  page.value = p
  load()
}

function actionLabel(action: string): string {
  const labels: Record<string, string> = {
    'flag.created': 'created flag',
    'flag.updated': 'updated flag',
    'flag.toggled': 'toggled flag',
    'flag.deleted': 'deleted flag',
    'env.created': 'created env',
    'env.updated': 'updated env',
    'env.deleted': 'deleted env',
  }
  return labels[action] ?? action
}

function actionBadgeClass(action: string): string {
  if (action.endsWith('.created')) return 'bg-green-500/10 text-green-600 dark:text-green-400'
  if (action.endsWith('.deleted')) return 'bg-destructive/10 text-destructive'
  if (action === 'flag.toggled') return 'bg-purple-500/10 text-purple-600 dark:text-purple-400'
  return 'bg-blue-500/10 text-blue-600 dark:text-blue-400'
}

function diffLines(entry: AuditEntry): string[] {
  const lines: string[] = []
  try {
    const oldV = entry.old_value ? JSON.parse(entry.old_value) : null
    const newV = entry.new_value ? JSON.parse(entry.new_value) : null

    if (entry.action === 'flag.toggled' && newV) {
      const arrow = `${oldV?.enabled ?? false} → ${newV.enabled}`
      lines.push(`${newV.env}: ${arrow}`)
      return lines
    }

    if (!oldV && newV) {
      for (const [k, v] of Object.entries(newV)) {
        if (v !== '' && v !== null) lines.push(`${k}: ${v}`)
      }
      return lines
    }

    if (oldV && !newV) {
      for (const [k, v] of Object.entries(oldV)) {
        if (v !== '' && v !== null) lines.push(`${k}: ${v}`)
      }
      return lines
    }

    if (oldV && newV) {
      for (const k of Object.keys(newV)) {
        if (String(oldV[k]) !== String(newV[k])) {
          lines.push(`${k}: ${oldV[k]} → ${newV[k]}`)
        }
      }
    }
  } catch {
    // malformed JSON — show nothing
  }
  return lines
}
</script>
