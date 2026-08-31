<template>
  <Teleport to="body">
    <div v-if="visible" class="app-dialog-overlay" @click.self="emit('close')">
      <section
        ref="dialogRef"
        class="app-dialog"
        :class="`size-${size}`"
        role="dialog"
        aria-modal="true"
        :aria-labelledby="titleId"
        tabindex="-1"
        @keydown="onKeydown"
      >
        <slot />
      </section>
    </div>
  </Teleport>
</template>

<script setup lang="ts">
import { nextTick, onBeforeUnmount, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    visible: boolean
    titleId: string
    size?: 'sm' | 'md' | 'lg'
  }>(),
  { size: 'md' },
)

const emit = defineEmits<{ close: [] }>()
const dialogRef = ref<HTMLElement | null>(null)
let restoreElement: HTMLElement | null = null
let previousOverflow = ''

const focusableSelector = [
  'a[href]',
  'button:not([disabled])',
  'input:not([disabled])',
  'select:not([disabled])',
  'textarea:not([disabled])',
  '[tabindex]:not([tabindex="-1"])',
].join(', ')

watch(
  () => props.visible,
  (visible) => {
    if (visible) void focusDialog()
    else restoreFocus()
  },
  { immediate: true },
)

function restoreFocus() {
  document.body.style.overflow = previousOverflow
  restoreElement?.focus()
  restoreElement = null
}

async function focusDialog() {
  restoreElement = document.activeElement instanceof HTMLElement ? document.activeElement : null
  previousOverflow = document.body.style.overflow
  document.body.style.overflow = 'hidden'
  await nextTick()
  const dialog = dialogRef.value
  const firstFocusable = dialog?.querySelector<HTMLElement>(focusableSelector)
  ;(firstFocusable ?? dialog)?.focus()
}

function onKeydown(event: KeyboardEvent) {
  if (event.key === 'Escape') {
    event.preventDefault()
    emit('close')
    return
  }

  if (event.key !== 'Tab') return

  const dialog = dialogRef.value
  if (!dialog) return
  const focusable = [...dialog.querySelectorAll<HTMLElement>(focusableSelector)]
  if (!focusable.length) {
    event.preventDefault()
    dialog.focus()
    return
  }

  const first = focusable[0]
  const last = focusable[focusable.length - 1]
  if (event.shiftKey && document.activeElement === first) {
    event.preventDefault()
    last.focus()
  } else if (!event.shiftKey && document.activeElement === last) {
    event.preventDefault()
    first.focus()
  }
}

onBeforeUnmount(() => {
  restoreFocus()
})
</script>

<style scoped>
.app-dialog-overlay {
  align-items: center;
  background: var(--color-overlay);
  display: flex;
  inset: 0;
  justify-content: center;
  padding: 20px;
  position: fixed;
  z-index: 1000;
}

.app-dialog {
  animation: dialog-in 0.15s ease-out;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-card);
  max-height: min(720px, calc(100vh - 40px));
  max-width: 520px;
  outline: none;
  overflow-y: auto;
  padding: 24px;
  width: 100%;
}

@keyframes dialog-in {
  from {
    opacity: 0;
    transform: scale(0.98) translateY(4px);
  }

  to {
    opacity: 1;
    transform: scale(1) translateY(0);
  }
}

.app-dialog.size-sm {
  max-width: 420px;
}

.app-dialog.size-lg {
  max-width: 640px;
}

@media (max-width: 640px) {
  .app-dialog-overlay {
    align-items: flex-end;
    padding: 12px;
  }

  .app-dialog {
    max-height: calc(100vh - 24px);
    padding: 20px;
  }
}

@media (prefers-reduced-motion: reduce) {
  .app-dialog {
    animation: none;
  }
}
</style>
