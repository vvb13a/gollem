<script setup lang="ts">
import Button from 'primevue/button'
import Toast from 'primevue/toast'
import { RouterView } from 'vue-router'
import { setThemeMode, theme } from './lib/theme'

function toggleTheme(): void {
  setThemeMode(theme.resolved === 'dark' ? 'light' : 'dark')
}
</script>

<template>
  <div class="min-h-screen bg-slate-50 text-slate-900 dark:bg-slate-950 dark:text-slate-100">
    <header class="border-b border-slate-200 dark:border-slate-800">
      <div class="mx-auto flex max-w-3xl items-center justify-between px-6 py-4">
        <div>
          <h1 class="text-lg font-semibold">gollem</h1>
          <p class="text-sm text-slate-500 dark:text-slate-400">
            Go + SQLite + Vue starter
          </p>
        </div>
        <Button
          :label="theme.resolved === 'dark' ? 'Light mode' : 'Dark mode'"
          severity="secondary"
          text
          @click="toggleTheme"
        />
      </div>
    </header>

    <main class="mx-auto max-w-3xl px-6 py-8">
      <RouterView />
    </main>

    <Toast />
  </div>
</template>
