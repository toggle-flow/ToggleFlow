<template>
  <Dialog :open="open" @update:open="onOpenChange">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ $t('members.addTitle') }}</DialogTitle>
        <DialogDescription>{{ $t('members.addDescription') }}</DialogDescription>
      </DialogHeader>

      <div class="space-y-3">
        <div class="space-y-2">
          <Label>{{ $t('members.user') }}</Label>
          <select
            v-model="selectedUserId"
            class="mt-2 w-full rounded-md border border-input bg-background px-3 py-2 text-sm focus:outline-none focus:ring-2 focus:ring-ring"
          >
            <option value="" disabled>{{ $t('members.selectUser') }}</option>
            <option v-for="u in availableUsers" :key="u.id" :value="u.id">
              {{ u.name }} — {{ $t(`roles.${u.role}`) }}
            </option>
          </select>
        </div>

        <Alert v-if="error" variant="destructive">
          <AlertCircle class="size-4" />
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>
      </div>

      <DialogFooter>
        <Button type="button" variant="outline" @click="$emit('update:open', false)">
          {{ $t('common.cancel') }}
        </Button>
        <Button :disabled="!selectedUserId || loading" @click="submit">
          <Loader2 v-if="loading" class="size-4 animate-spin" />
          {{ loading ? $t('members.adding') : $t('members.add') }}
        </Button>
      </DialogFooter>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { AlertCircle, Loader2 } from '@lucide/vue'
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from '@/components/ui/dialog'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { membersApi, type Member } from '@/api/members'
import { usersApi, type User } from '@/api/users'

const props = defineProps<{
  open: boolean
  projectId: number
  existingMemberIds: number[]
}>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  added: [member: Member]
}>()

const allUsers = ref<User[]>([])
const selectedUserId = ref<number | ''>('')
const loading = ref(false)
const error = ref('')

const availableUsers = computed(() =>
  allUsers.value.filter((u) => !props.existingMemberIds.includes(u.id))
)

watch(
  () => props.open,
  async (v) => {
    if (v) {
      selectedUserId.value = ''
      error.value = ''
      try {
        allUsers.value = await usersApi.list()
      } catch {
        allUsers.value = []
      }
    }
  }
)

function onOpenChange(v: boolean) {
  emit('update:open', v)
}

async function submit() {
  if (!selectedUserId.value) return
  error.value = ''
  loading.value = true
  try {
    const member = await membersApi.add(props.projectId, selectedUserId.value as number)
    emit('added', member)
    emit('update:open', false)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Something went wrong'
  } finally {
    loading.value = false
  }
}
</script>
