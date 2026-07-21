<script setup>
import { ref,onMounted } from 'vue'
import { useRoute,useRouter } from 'vue-router'
import { api } from '../api'
import StateMessage from '../components/StateMessage.vue'
const route=useRoute(),router=useRouter(),poll=ref(null),selected=ref({}),loading=ref(true),busy=ref(false),error=ref('')
async function load(){loading.value=true;error.value='';try{poll.value=await api.getPoll(route.params.id)}catch(e){error.value=e.message}finally{loading.value=false}}
async function vote(){const answers=poll.value.questions.map(q=>({question_id:q.id,option_id:selected.value[q.id]}));if(answers.some(a=>!a.option_id))return;busy.value=true;try{await api.vote(route.params.id,answers);router.push(`/poll/${route.params.id}/results`)}catch(e){error.value=e.message}finally{busy.value=false}}
onMounted(load)
</script>
<template><section class="mx-auto max-w-2xl px-5 pb-16 pt-10 sm:px-8 sm:pt-20"><RouterLink to="/" class="text-sm font-bold text-ink/50 hover:text-coral">← На главную</RouterLink><StateMessage :loading="loading" :error="error" @retry="load"/><form v-if="poll&&!loading" class="mt-10" @submit.prevent="vote"><p class="text-sm font-bold uppercase tracking-[.2em] text-coral">Ваш голос</p><h1 class="mt-5 font-display text-5xl leading-tight sm:text-6xl">{{poll.title}}</h1><div v-for="(question,i) in poll.questions" :key="question.id" class="mt-10"><h2 class="font-display text-3xl">{{i+1}}. {{question.text}}</h2><fieldset class="mt-4 space-y-3"><legend class="sr-only">Выберите вариант ответа</legend><label v-for="option in question.options" :key="option.id" class="flex cursor-pointer items-center gap-4 rounded-2xl border border-ink/10 bg-white/70 p-4 transition hover:border-coral has-[:checked]:border-coral has-[:checked]:bg-mint/40"><input v-model="selected[question.id]" class="h-5 w-5 accent-coral" type="radio" :name="`answer-${question.id}`" :value="option.id"><span class="font-semibold">{{option.text}}</span></label></fieldset></div><button class="button button-dark mt-8 w-full" :disabled="busy||poll.questions.some(q=>!selected[q.id])">{{busy?'Отправляем…':'Проголосовать ↗'}}</button></form></section></template>
