<template>
  <main class="container">
    <h1 class="display-1">
      TiDB SQL Parser
    </h1>

    <br/>

    <b-tabs>
      <b-tab v-for="s in sections" :key="s.key" :title="s.title" :active="section === s.key"
             @click="() => section = s.key">
        <keep-alive>
          <Codes
              v-if="section === s.key"
              :language="s.language"
              :model-value="s.value.value"
              @update:model-value="val => s.value.value = val"
          />
        </keep-alive>
      </b-tab>
    </b-tabs>

    <section>
      <b-button variant="primary" @click="parse">
        Parse!
      </b-button>
    </section>

    <section v-if="warns.length">
      <b-alert show variant="warning">
        <h4 class="alert-heading">Warns</h4>
        <ul>
          <li v-for="warn in warns">{{ warn }}</li>
        </ul>
      </b-alert>
    </section>

    <section v-if="columns">
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

    <section v-if="error">
      <b-alert show variant="danger">
        <h4 class="alert-heading">Somethings wrong</h4>
        <pre>{{ error?.message }}</pre>
      </b-alert>
    </section>
  </main>
</template>
<script lang="ts" setup>
import {defineAsyncComponent, Ref, ref, shallowRef} from "vue";
import {Column, EvalTypeNames, init} from "./index";
import DDL_SQL from './schema.sql?raw'
import DEFINES_JS from './defines.js?raw'
import isMobile from "ismobilejs";

const Codes = defineAsyncComponent(() => {
  if (isMobile(navigator.userAgent).any) {
    return import("./components/highlight")
  } else {
    return import("./components/monaco-editor")
  }
});

const usp = new URLSearchParams(window.location.search)

const query = usp.get('query') || 'trending-repos'

const QUERY_URL = `https://raw.githubusercontent.com/pingcap/ossinsight/main/api/queries/${query}/template.sql`


const section = ref<'ddl' | 'defines' | 'query'>('query')
const sections: {
  key: string
  title: string
  language: string
  value: Ref<string>
}[] = [
  {
    key: 'query',
    title: 'Query',
    language: 'mysql',
    value: ref('-- Loading...'),
  },
  {
    key: 'ddl',
    title: 'DDL',
    language: 'mysql',
    value: ref(DDL_SQL),
  },
  {
    key: 'defines',
    title: 'Defines',
    language: 'javascript',
    value: ref(DEFINES_JS),
  },
]

fetch(QUERY_URL).then(res => res.text()).then(text => {
  sections[0].value.value = `-- ${QUERY_URL}\n\n${text}`
})

const columns = shallowRef<Column[]>()
const error = ref<unknown>()
const warns = shallowRef<string[]>([])

async function parse() {
  try {
    columns.value = undefined
    error.value = undefined
    warns.value = []
    const program = await init()
    const p = program.newParser();

    p.AddDdl(sections.find(s => s.key === 'ddl')!.value.value)
    eval(sections.find(s => s.key === 'defines')!.value.value)

    const rawColumns = p.Parse(sections.find(s => s.key === 'query')!.value.value)
    // go value would be unavailable after program stopped.
    columns.value = JSON.parse(JSON.stringify(rawColumns))
    const w = warns.value = p.Warns()

    console.warn(w)

    program.stop()
  } catch (e) {
    error.value = e
  }
}

</script>
<style scoped>

::v-deep(.tab-content) {
  border-width: 1px;
  border-style: solid;
  border-color: #dee2e6;
  border-top: none;
  padding-right: 1px;
}

@media (prefers-color-scheme: dark) {
  ::v-deep(.tab-content) {
    border-color: #515151;
  }
}

section {
  margin-top: 16px;
}
</style>