<template>
  <div class="min-h-screen flex flex-col items-center justify-center bg-background px-8 py-10">
    <div class="w-full max-w-sm space-y-5">
      <div class="flex items-center gap-2">
        <ToggleRight class="size-5 text-primary" />
        <span class="font-semibold tracking-tight">ToggleFlow</span>
      </div>

      <div class="space-y-1">
        <h2 class="text-xl font-semibold tracking-tight">{{ $t('reset.title') }}</h2>
        <p v-if="inviteName" class="text-sm font-medium">{{ inviteName }}</p>
        <p v-if="inviteEmail" class="text-xs text-muted-foreground">{{ inviteEmail }}</p>
        <p class="text-sm text-muted-foreground" :class="{ 'mt-2': inviteName || inviteEmail }">
          {{ $t('reset.subtitle') }}
        </p>
      </div>

      <form class="space-y-3" @submit.prevent="submit">
        <div class="space-y-2">
          <Label for="secret">{{ $t('activate.secret') }}</Label>
          <Input
            id="secret"
            v-model="token"
            :placeholder="$t('activate.secretPlaceholder')"
            class="mt-2 font-mono tracking-widest"
            required
            autofocus
            autocomplete="off"
            maxlength="24"
          />
        </div>

        <div class="space-y-2">
          <Label for="password">{{ $t('reset.newPassword') }}</Label>
          <Input
            id="password"
            v-model="password"
            type="password"
            :placeholder="$t('activate.passwordPlaceholder')"
            class="mt-2"
            required
            minlength="8"
          />
        </div>

        <div class="space-y-2">
          <Label for="confirm-password">{{ $t('activate.confirmPassword') }}</Label>
          <Input
            id="confirm-password"
            v-model="confirmPassword"
            type="password"
            :placeholder="$t('activate.confirmPasswordPlaceholder')"
            class="mt-2"
            required
          />
        </div>

        <Alert v-if="error" variant="destructive">
          <AlertCircle class="size-4" />
          <AlertDescription>{{ error }}</AlertDescription>
        </Alert>

        <Button type="submit" class="w-full" :disabled="loading">
          <Loader2 v-if="loading" class="size-4 animate-spin" />
          {{ loading ? $t('reset.submitting') : $t('reset.submit') }}
        </Button>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useI18n } from 'vue-i18n'
import { ToggleRight, Loader2, AlertCircle } from '@lucide/vue'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Button } from '@/components/ui/button'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { authApi } from '@/api/auth'
import { useAuthStore } from '@/stores/auth'

const { t } = useI18n()
const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const inviteName = ref('')
const inviteEmail = ref('')
const token = ref('')
const password = ref('')
const confirmPassword = ref('')
const loading = ref(false)
const error = ref('')

onMounted(async () => {
  const id = route.query.id as string
  if (!id) return
  try {
    const info = await authApi.getResetInfo(id)
    inviteName.value = info.name
    inviteEmail.value = info.email
  } catch {
    // invalid or already used — the form will surface the error on submit
  }
})

async function submit() {
  if (password.value !== confirmPassword.value) {
    error.value = t('activate.passwordMismatch')
    return
  }
  error.value = ''
  loading.value = true
  try {
    const { token: jwt, user } = await authApi.resetPassword(token.value, password.value)
    authStore.setAuth(jwt, user)
    router.push('/')
  } catch (e: unknown) {
    error.value = e instanceof Error ? e.message : t('common.error')
  } finally {
    loading.value = false
  }
}
</script>
