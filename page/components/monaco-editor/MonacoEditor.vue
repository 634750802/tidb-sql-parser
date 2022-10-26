<template>
  <div class="editor" ref="el">
  </div>
</template>
<script lang="ts" setup>
import * as monaco from 'monaco-editor';
import {onMounted, onUnmounted, ref, shallowRef, watch} from "vue";
import {preferDarkScheme} from "../../hooks/media-query.js";

const props = defineProps<{
  language: string,
  modelValue: string,
}>()

const emit = defineEmits<{
  (event: 'update:model-value', value: string): void
}>()

const el = ref<HTMLDivElement | null>(null)
const editor = shallowRef<monaco.editor.IStandaloneCodeEditor>()
const dark = preferDarkScheme()

onMounted(() => {
  const instance = editor.value = monaco.editor.create(el.value!, {
    language: props.language,
    value: props.modelValue,
    theme: dark.value ? 'vs-dark' : 'vs',
    scrollbar: {
      horizontalScrollbarSize: 8,
      verticalScrollbarSize: 8,
      useShadows: true,
    },
    smoothScrolling: true,
  })

  onUnmounted(() => {
    instance.dispose();
  })
})

watch(editor, (editor, _, onCleanup) => {
  if (editor) {
    const handleUpdate = () => emit('update:model-value', editor.getValue());
    const disposables: monaco.IDisposable[] = []
    disposables.push(
        editor.onDidBlurEditorText(handleUpdate),
    )
    onCleanup(() => disposables.forEach(disposable => disposable.dispose()))
  }
})

watch(() => props.modelValue, (value) => {
  console.log(value)
  editor.value?.setValue(value)
})

watch(dark, (dark) => {
  editor.value?.updateOptions({
    theme: dark ? 'vs-dark' : 'vs'
  })
}, {immediate: true})

</script>
<style scoped>
.editor {
  min-height: 400px;
}
</style>