<template>
  <Dialog :open="open" @update:open="$emit('update:open', $event)">
    <DialogContent class="sm:max-w-sm">
      <DialogHeader>
        <DialogTitle>{{ $t('flags.createTitle') }}</DialogTitle>
        <DialogDescription>{{ $t('flags.createDescription') }}</DialogDescription>
      </DialogHeader>

      <form class="space-y-3" @submit.prevent="submit">
        <div class="space-y-1.5">
          <Label for="flag-name">{{ $t('flags.name') }}</Label>
          <Input
            id="flag-name"
            v-model="name"
            :placeholder="$t('flags.namePlaceholder')"
            required
            autofocus
            @input="syncKey"
          />
        </div>

        <div class="space-y-1.5">
          <Label for="flag-key">{{ $t('flags.key') }}</Label>
          <Input
            id="flag-key"
            v-model="key"
            :placeholder="$t('flags.keyPlaceholder')"
            required
            class="font-mono"
          />
        </div>

        <div class="space-y-1.5">
          <Label for="flag-description">{{ $t('flags.description') }}</Label>
          <Input
            id="flag-description"
            v-model="description"
            :placeholder="$t('flags.descriptionPlaceholder')"
          />
        </div>

        <Alert v-if="error" variant="destructive">
          <AlertCircle class="size-4" />
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>

        <DialogFooter>
          <Button type="button" variant="outline" @click="$emit('update:open', false)">
            {{ $t('common.cancel') }}
          </Button>
          <Button type="submit" :disabled="loading">
            <Loader2 v-if="loading" class="size-4 animate-spin" />
            {{ loading ? $t('flags.creating') : $t('flags.create') }}
          </Button>
        </DialogFooter>
      </form>
    </DialogContent>
  </Dialog>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { AlertCircle, Loader2 } from '@lucide/vue'
import { Dialog, DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle } from '@/components/ui/dialog'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { flagsApi, type Flag } from '@/api/flags'

const props = defineProps<{ open: boolean; projectId: number }>()
const emit = defineEmits<{
  'update:open': [value: boolean]
  created: [flag: Flag]
}>()

const name = ref('')
const key = ref('')
const description = ref('')
const loading = ref(false)
const error = ref('')
let keyTouched = false

watch(() => props.open, (v) => {
  if (v) { name.value = ''; key.value = ''; description.value = ''; error.value = ''; keyTouched = false }
})

function slugify(s: string) {
  return s.toLowerCase().replace(/[^a-z0-9]+/g, '-').replace(/^-|-$/g, '')
}

function syncKey() {
  if (!keyTouched) key.value = slugify(name.value)
}

async function submit() {
  error.value = ''
  loading.value = true
  try {
    const flag = await flagsApi.create(props.projectId, {
      name: name.value,
      key: key.value,
      description: description.value || undefined,
    })
    emit('created', flag)
    emit('update:open', false)
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : 'Something went wrong'
  } finally {
    loading.value = false
  }
}
</script>
