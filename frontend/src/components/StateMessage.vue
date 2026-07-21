<script setup>
defineProps({ loading:Boolean, error:String })
defineEmits(['retry'])
</script>
<template><div v-if="loading||error" class="rounded-3xl border border-ink/10 bg-white/60 p-8 text-center"><div v-if="loading" class="mx-auto mb-3 h-7 w-7 animate-spin rounded-full border-2 border-ink/20 border-t-coral"/><p>{{ error||'Загружаем…' }}</p><button v-if="error" class="button button-light mt-5" @click="$emit('retry')">Повторить</button></div></template>
