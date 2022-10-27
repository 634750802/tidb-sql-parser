<template>
  <pre><code class="highlight" :class="`language-${language}`" v-html="html"/></pre>
</template>
<script lang="ts" setup>
import Highlight from 'highlight.js'
import {ref, watch} from "vue";

const props = defineProps<{
  modelValue: string
  language: string
  ignoreIllegals?: boolean
}>()

const html = ref('')

watch(() => props.modelValue, value => {
  html.value = Highlight.highlight(value, {
    language: props.language,
    ignoreIllegals: props.ignoreIllegals ?? true,
  }).value
}, { immediate: true })
</script>
<style scoped>
pre {
  padding: 8px;
  max-height: 400px;
  margin: 0;
  scroll-behavior: smooth;
}
</style>

