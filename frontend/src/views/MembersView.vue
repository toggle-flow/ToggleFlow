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
        <div class="flex items-center justify-between">
          <div class="space-y-0.5">
            <h1 class="text-base font-semibold">{{ $t('members.title') }}</h1>
            <p class="text-sm text-muted-foreground">{{ $t('members.subtitle') }}</p>
          </div>
          <Button v-if="authStore.isOwner" size="sm" @click="addDialogOpen = true">
            <UserPlus class="size-3.5" />
            {{ $t('members.add') }}
          </Button>
        </div>

        <Separator />
      </div>

      <!-- Scrollable body -->
      <div class="min-h-0 flex-1 overflow-y-auto px-6 pb-6">
        <div v-if="loading" class="flex h-full items-center justify-center">
          <Loader2 class="size-6 animate-spin text-muted-foreground/40" />
        </div>

        <div
          v-else-if="members.length === 0"
          class="flex h-full flex-col items-center justify-center text-center"
        >
          <Users class="size-8 text-muted-foreground/30 mb-3" />
          <p class="text-sm font-medium">{{ $t('members.emptyTitle') }}</p>
          <p class="mt-1 text-xs text-muted-foreground max-w-xs">
            {{ $t('members.emptyDescription') }}
          </p>
        </div>

        <div v-else class="space-y-2">
          <div
            v-for="member in members"
            :key="member.user_id"
            class="rounded-lg border bg-card p-4"
          >
            <div class="flex items-center justify-between gap-4">
              <div class="flex items-center gap-3 min-w-0">
                <div
                  class="flex size-8 shrink-0 items-center justify-center rounded-full bg-muted text-xs font-semibold uppercase"
                >
                  {{ member.name[0] }}
                </div>
                <div class="min-w-0">
                  <div class="flex flex-wrap items-center gap-2">
                    <p class="text-sm font-medium leading-none">{{ member.name }}</p>
                    <span
                      v-if="member.user_id === authStore.user?.id"
                      class="text-[10px] text-muted-foreground"
                      >({{ $t('users.you') }})</span
                    >
                    <span
                      :class="roleBadgeClass(member.role)"
                      class="inline-flex items-center rounded px-1.5 py-0.5 text-[10px] font-medium"
                    >
                      {{ $t(`roles.${member.role}`) }}
                    </span>
                  </div>
                  <p class="mt-0.5 text-xs text-muted-foreground">{{ member.email }}</p>
                </div>
              </div>

              <div class="flex items-center gap-2 shrink-0">
                <span class="text-xs text-muted-foreground">
                  <CalendarDays class="inline size-3 mr-0.5" />{{ timeAgo(member.joined_at) }}
                </span>
                <Tooltip v-if="canRemove(member)" :text="$t('members.remove')">
                  <button
                    class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
                    @click="remove(member)"
                  >
                    <UserMinus class="size-3.5" />
                  </button>
                </Tooltip>
              </div>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>

  <AddMemberDialog
    v-if="projectStore.current"
    v-model:open="addDialogOpen"
    :project-id="projectStore.current.id"
    :existing-member-ids="members.map((m) => m.user_id)"
    @added="onAdded"
  />
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { FolderOpen, Loader2, UserPlus, UserMinus, Users, CalendarDays } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { Tooltip } from '@/components/ui/tooltip'
import { useProjectStore } from '@/stores/project'
import { useAuthStore } from '@/stores/auth'
import { membersApi, type Member } from '@/api/members'
import { timeAgo } from '@/lib/utils'
import AddMemberDialog from '@/components/AddMemberDialog.vue'

const projectStore = useProjectStore()
const authStore = useAuthStore()
const members = ref<Member[]>([])
const loading = ref(false)
const addDialogOpen = ref(false)

function roleBadgeClass(role: string) {
  switch (role) {
    case 'superuser':
      return 'bg-destructive/10 text-destructive'
    case 'admin':
      return 'bg-orange-500/10 text-orange-600 dark:text-orange-400'
    case 'owner':
      return 'bg-blue-500/10 text-blue-600 dark:text-blue-400'
    case 'editor':
      return 'bg-green-500/10 text-green-600 dark:text-green-400'
    default:
      return 'bg-muted text-muted-foreground'
  }
}

function canRemove(member: Member) {
  if (!authStore.isOwner) return false
  if (member.user_id === authStore.user?.id) return false
  return true
}

async function load() {
  if (!projectStore.current) return
  loading.value = true
  try {
    members.value = await membersApi.list(projectStore.current.id)
  } finally {
    loading.value = false
  }
}

watch(() => projectStore.current, load, { immediate: true })

function onAdded(member: Member) {
  members.value.push(member)
}

async function remove(member: Member) {
  if (!projectStore.current) return
  try {
    await membersApi.remove(projectStore.current.id, member.user_id)
    members.value = members.value.filter((m) => m.user_id !== member.user_id)
  } catch {
    // leave the list unchanged on error
  }
}
</script>
