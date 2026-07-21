<script setup>
import { RouterLink } from 'vue-router'
</script>
<template><div class="min-h-screen overflow-hidden bg-cream text-ink"><header class="mx-auto flex max-w-6xl items-center justify-between px-5 py-6 sm:px-8"><RouterLink to="/" class="font-display text-3xl tracking-tight">пульс<span class="text-coral">.</span></RouterLink><RouterLink to="/create" class="button button-dark hidden sm:inline-flex">Создать опрос <span>↗</span></RouterLink></header><main><RouterView/></main><footer class="mx-auto max-w-6xl px-5 py-12 text-sm text-ink/50 sm:px-8">Простые вопросы. Живые ответы.</footer></div></template>
