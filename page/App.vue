<template>
  <main class="container">
    <h1>
      TiDB SQL Parser
    </h1>
    <b-card>
      <div>
        <b-form-group
            id="input-group-ddl"
            label="DDL"
            label-for="ddl"
            description="Write some TiDB DDL SQL here"
        >
          <b-form-textarea
              class="textarea"
              id="ddl"
              v-model="ddl"
              placeholder="Enter DDL"
              required
          />
        </b-form-group>
        <b-form-group
            id="input-group-query"
            label="Query"
            label-for="query"
            description="Write some TiDB query SQL here"
        >
          <b-form-textarea
              class="textarea"
              id="query"
              v-model="query"
              placeholder="Enter Query"
              required
          />
        </b-form-group>
        <b-button variant="primary" @click="parse">
          Parse!
        </b-button>
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
      </div>
    </b-card>
  </main>
</template>
<script lang="ts" setup>
import {ref, shallowRef} from "vue";
import {Column, EvalTypes, EvalTypeNames, init} from "./index";

const ddl = ref('')
const query = ref('')
const columns = shallowRef<Column[]>()

Promise.all([
  fetch('/schema.sql').then(res => res.text()).then(res => ddl.value = res),
  fetch('/query.sql').then(res => res.text()).then(res => query.value = res),
])
    .then(parse)

async function parse() {
  const program = await init()
  const p = program.newParser();

  p.DefineTransparentFunc("IFNULL")
  p.DefineTransparentFunc("ABS")
  p.DefineTransparentFunc("SUM")
  p.DefineTransparentFunc("GREATEST")
  p.DefineTransparentFunc("LEAST")
  p.DefineFunc("DATE_SUB", {Type: EvalTypes.ETDatetime, Nullable: false})
  p.DefineFunc("COUNT", {Type: EvalTypes.ETInt, Nullable: false})
  p.DefineFunc("TIMESTAMPDIFF", {Type: EvalTypes.ETReal, Nullable: false})
  p.AddDdl(ddl.value)

  // go value would be unavailable after program stopped.
  columns.value = JSON.parse(JSON.stringify(p.Parse(query.value)))
  program.stop()
}

</script>
<style scoped>
.textarea {
  min-height: 250px;
}
</style>