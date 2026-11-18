import { reactive } from 'vue'

export type ResolvedTheme = 'light' | 'dark'

const THEME_STORAGE_KEY = 'app:theme'

function loadTheme(): ResolvedTheme {
  try {
    const raw = localStorage.getItem(THEME_STORAGE_KEY)
    if (raw === 'light' || raw === 'dark') return raw
  } catch {
    // Storage unavailable: fall back to the system preference.
  }
  return window.matchMedia('(prefers-color-scheme: dark)').matches
    ? 'dark'
    : 'light'
}

export const theme = reactive({
  resolved: loadTheme() as ResolvedTheme,
})

function apply(): void {
  const root = document.documentElement
  root.classList.toggle('dark', theme.resolved === 'dark')
  root.classList.toggle('p-dark', theme.resolved === 'dark')
  root.style.colorScheme = theme.resolved
}

export function setThemeMode(mode: ResolvedTheme): void {
  theme.resolved = mode
  try {
    localStorage.setItem(THEME_STORAGE_KEY, mode)
  } catch {
    // Storage unavailable: the choice simply won't persist.
  }
  apply()
}

export function initTheme(): void {
  apply()
}
