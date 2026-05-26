<template>
  <div class="p-6 space-y-4">
    <div class="flex items-center justify-between">
      <div class="space-y-0.5">
        <h1 class="text-base font-semibold">{{ $t('nav.projects') }}</h1>
        <p class="text-sm text-muted-foreground">{{ $t('projects.subtitle') }}</p>
      </div>
      <Button size="sm" @click="createDialogOpen = true">
        <Plus class="size-3.5" />
        {{ $t('projects.new') }}
      </Button>
    </div>

    <Separator />

    <div v-if="loading" class="flex flex-col items-center justify-center py-24">
      <Loader2 class="size-6 animate-spin text-muted-foreground/40" />
    </div>

    <div v-else-if="projects.length === 0" class="flex flex-col items-center justify-center py-24 text-center">
      <FolderOpen class="size-8 text-muted-foreground/30 mb-3" />
      <p class="text-sm font-medium">{{ $t('projects.emptyTitle') }}</p>
      <p class="mt-1 text-xs text-muted-foreground max-w-xs">{{ $t('projects.emptyDescription') }}</p>
    </div>

    <div v-else class="space-y-2">
      <div
        v-for="project in projects"
        :key="project.id"
        class="rounded-lg border bg-card p-4"
      >
        <div class="flex items-center justify-between gap-4">
          <div class="min-w-0 flex-1">
            <p class="text-sm font-medium">{{ project.name }}</p>
            <p class="text-xs font-mono text-muted-foreground mt-0.5">{{ project.slug }}</p>
          </div>
          <div class="flex items-center gap-1 shrink-0">
            <button
              class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-accent hover:text-foreground"
              @click="openEdit(project)"
            >
              <Pencil class="size-3.5" />
            </button>
            <button
              class="rounded-md p-1.5 text-muted-foreground transition-colors hover:bg-destructive/10 hover:text-destructive"
              @click="openDelete(project)"
            >
              <Trash2 class="size-3.5" />
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>

  <CreateProjectDialog v-model:open="createDialogOpen" @created="onCreated" />
  <EditProjectDialog v-model:open="editDialogOpen" :project="editTarget" @updated="onUpdated" />
  <DeleteProjectDialog v-model:open="deleteDialogOpen" :project="deleteTarget" @deleted="onDeleted" />
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Plus, FolderOpen, Loader2, Pencil, Trash2 } from '@lucide/vue'
import { Button } from '@/components/ui/button'
import { Separator } from '@/components/ui/separator'
import { projectsApi, type Project } from '@/api/projects'
import { useProjectStore } from '@/stores/project'
import CreateProjectDialog from '@/components/CreateProjectDialog.vue'
import EditProjectDialog from '@/components/EditProjectDialog.vue'
import DeleteProjectDialog from '@/components/DeleteProjectDialog.vue'

const projectStore = useProjectStore()
const projects = ref<Project[]>([])
const loading = ref(false)
const createDialogOpen = ref(false)
const editDialogOpen = ref(false)
const editTarget = ref<Project | null>(null)
const deleteDialogOpen = ref(false)
const deleteTarget = ref<Project | null>(null)

onMounted(async () => {
  loading.value = true
  try {
    projects.value = await projectsApi.list() ?? []
  } finally {
    loading.value = false
  }
})

function openEdit(project: Project) {
  editTarget.value = project
  editDialogOpen.value = true
}

function openDelete(project: Project) {
  deleteTarget.value = project
  deleteDialogOpen.value = true
}

function onCreated(project: Project) {
  projects.value.push(project)
  projectStore.projects.push(project)
  if (!projectStore.current) projectStore.setCurrent(project)
}

function onUpdated(updated: Project) {
  const i = projects.value.findIndex(p => p.id === updated.id)
  if (i !== -1) projects.value[i] = updated

  // Keep the sidebar store in sync
  const si = projectStore.projects.findIndex(p => p.id === updated.id)
  if (si !== -1) projectStore.projects[si] = updated
  if (projectStore.current?.id === updated.id) projectStore.setCurrent(updated)
}

function onDeleted(project: Project) {
  const remaining = projects.value.filter(p => p.id !== project.id)
  projects.value = remaining
  projectStore.setProjects(remaining)
}
</script>
