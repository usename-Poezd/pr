<script setup>
import { ref,onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { api } from '../api'
import StateMessage from '../components/StateMessage.vue'
import { getCreatedPolls } from '../utils/storage'
const route=useRoute(),results=ref(null),loading=ref(true),error=ref(''),copied=ref(false)
const token=getCreatedPolls().find(p=>p.id===route.params.id)?.adminToken
const createdPoll=getCreatedPolls().find(p=>p.id===route.params.id)
const publicLink=createdPoll?.publicLink||`/polls/${route.params.id}`
const publicUrl=window.location.origin+publicLink
function total(question){return question.options.reduce((n,o)=>n+(o.votes||0),0)}
async function load(){loading.value=true;error.value='';try{results.value=await api.getResults(route.params.id,token)}catch(e){error.value=e.message}finally{loading.value=false}}
async function share(){if(navigator.share){await navigator.share({title:results.value.title,url:publicUrl});return}await navigator.clipboard.writeText(publicUrl);copied.value=true;setTimeout(()=>copied.value=false,1800)}
onMounted(load)
</script>
<template>
  <section class="mx-auto max-w-2xl px-5 pb-16 pt-10 sm:px-8 sm:pt-20">
    <RouterLink to="/" class="text-sm font-bold text-ink/50 hover:text-coral">← На главную</RouterLink>
    <StateMessage :loading="loading" :error="error" @retry="load"/>
    <div v-if="route.query.created" class="mt-10 rounded-3xl bg-mint p-5 sm:p-6" role="status">
        <p class="text-lg font-extrabold">Ваш опрос успешно создан</p>
        <p class="mt-1 text-sm text-ink/70">Поделитесь ссылкой — результаты обновляются здесь.</p>
        <div class="mt-4 flex flex-col gap-3 sm:flex-row sm:items-center">
          <a :href="publicUrl" class="min-w-0 flex-1 truncate rounded-xl bg-white px-3 py-2 text-sm font-semibold text-ink underline decoration-ink/20 underline-offset-4">{{publicUrl}}</a>
          <button type="button" class="button button-dark shrink-0" @click="share">{{copied?'Скопировано':'Поделиться / копировать'}}</button>
        </div>
    </div>
    <div v-if="results&&!loading" class="mt-10">
      <p class="text-sm font-bold uppercase tracking-[.2em] text-coral">Результаты</p>
      <h1 class="mt-5 font-display text-5xl leading-tight sm:text-6xl">{{results.title}}</h1>
      <div v-for="(question,i) in results.questions" :key="question.id" class="mt-10">
        <h2 class="font-display text-3xl">{{i+1}}. {{question.text}}</h2>
        <p class="mt-2 text-ink/55">{{total(question)}} {{total(question)===1?'голос':'голосов'}}</p>
        <div class="mt-6 space-y-6">
          <div v-for="option in question.options" :key="option.id">
            <div class="mb-2 flex justify-between gap-4 font-semibold"><span>{{option.text}}</span><span>{{Math.round((option.votes/Math.max(total(question),1))*100)}}%</span></div>
            <div class="h-3 overflow-hidden rounded-full bg-ink/10"><div class="h-full rounded-full bg-coral transition-all duration-700" :style="{width:`${Math.round((option.votes/Math.max(total(question),1))*100)}%`}"/></div>
          </div>
        </div>
      </div>
      <RouterLink :to="`/poll/${route.params.id}`" class="button button-light mt-10">Проголосовать ещё раз</RouterLink>
    </div>
  </section>
</template>
