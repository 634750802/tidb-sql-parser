<template>
  <main class="container">
    <h1>
      TiDB SQL Parser
    </h1>
    <section>
      <h2>DDL SQL</h2>
      <div
          class="editor"
          ref="ddlEl"
          id="ddl"
      />
    </section>

    <section>
      <h2>Function defines</h2>
      <div
          class="editor"
          ref="definesEl"
          id="defines"
      />
    </section>

    <section>
      <h2>Query SQL</h2>
      <div
          class="editor"
          ref="queryEl"
          id="query"
      />
    </section>

    <section>
      <b-button variant="primary" @click="parse">
        Parse!
      </b-button>
    </section>

    <section>
      <b-table-simple hover small caption-top responsive>
        <b-thead head-variant="dark">
          <b-tr>
            <b-th>Name</b-th>
            <b-th>Type</b-th>
            <b-th>Nullable (WIP)</b-th>
          </b-tr>
        </b-thead>
        <b-tbody>
          <b-tr v-for="column in columns" :key="column.Name">
            <b-th>{{ column.As || column.Name }}</b-th>
            <b-td>{{ EvalTypeNames[column.Type] }}</b-td>
            <b-td>{{ column.Nullable }}</b-td>
          </b-tr>
        </b-tbody>
      </b-table-simple>
    </section>
  </main>
</template>
<script lang="ts" setup>
import {ref, shallowRef, watch} from "vue";
import {Column, EvalTypeNames, init} from "./index";
import type {editor} from "monaco-editor";
import * as monaco from 'monaco-editor';
import TsWorker from 'monaco-editor/esm/vs/language/typescript/ts.worker.js?worker'
import EditorWorker from 'monaco-editor/esm/vs/editor/editor.worker.js?worker'
import DDL_SQL from './schema.sql?raw'
import QUERY_SQL from './query.sql?raw'
import DEFINES_JS from './defines.js?raw'

const columns = shallowRef<Column[]>()

const ddlEl = ref<HTMLTextAreaElement | null>(null)
const queryEl = ref<HTMLTextAreaElement | null>(null)
const definesEl = ref<HTMLTextAreaElement | null>(null)

const ddlEditor = shallowRef<editor.IStandaloneCodeEditor>()
const queryEditor = shallowRef<editor.IStandaloneCodeEditor>()
const definesEditor = shallowRef<editor.IStandaloneCodeEditor>()

window.MonacoEnvironment = {
  getWorker: function (workerId, label) {
    switch (label) {
      case 'typescript':
      case 'javascript':
        return new TsWorker();
      default:
        return new EditorWorker();
    }
  }
};

async function parse() {
  const program = await init()
  const p = program.newParser();


  p.AddDdl(ddlEditor.value!.getValue())
  eval(definesEditor.value!.getValue())

  // go value would be unavailable after program stopped.
  columns.value = JSON.parse(JSON.stringify(p.Parse(queryEditor.value!.getValue())))
  program.stop()
}

watch(ddlEl, el => {
  if (el) {
    ddlEditor.value = monaco.editor.create(el, {
      value: DDL_SQL,
      language: 'mysql'
    });
  }
})

watch(definesEl, el => {
  if (el) {
    definesEditor.value = monaco.editor.create(el, {
      value: DEFINES_JS,
      language: 'javascript',
    });
  }
})

watch(queryEl, el => {
  if (el) {
    queryEditor.value = monaco.editor.create(el, {
      value: QUERY_SQL,
      language: 'mysql',
    });
  }
})

</script>
<style scoped>
.editor {
  min-height: 400px;
  border: 1px solid lightgray;
}

section {
  margin-top: 16px;
}
</style>