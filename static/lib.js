const go = new Go();
let stream = fetch(document.baseURI + "main.wasm", {headers: {'accept': 'application/wasm'}})
    .then(res => {
        window.dispatchEvent(new CustomEvent('parserloaded'))
        return res
    })

let evalTypes = {}

WebAssembly.instantiateStreaming(stream, go.importObject).then((result) => {
    go.run(result.instance);
    for (const i in __tidbSqlParserEvalTypes) {
        evalTypes[__tidbSqlParserEvalTypes[i]] = parseInt(i)
    }
    window.__onInit?.()
});

class Parser {
    constructor() {
        this.open()
    }

    get isOpen() {
        return __tidbSqlParserExecuteCmd(0).Data
    }

    open() {
        console.debug('[tidb-sql-parser] open parser')
        return __tidbSqlParserExecuteCmd(1)
    }

    addDdl(sql) {
        console.debug('[tidb-sql-parser] add ddl:', sql)
        return __tidbSqlParserExecuteCmd(3, sql)
    }

    defineFunc(func, type, nullable) {
        console.debug('[tidb-sql-parser] define function:', func, type, nullable)
        return __tidbSqlParserExecuteCmd(4, func, evalTypes[type], nullable)
    }

    parse(sql) {
        console.debug('[tidb-sql-parser] parse:', sql)
        return __tidbSqlParserExecuteCmd(5, sql)
    }

    getTable (name) {
        console.debug('[tidb-sql-parser] get table:', name)
        return __tidbSqlParserExecuteCmd(6, name)
    }

    defineTransparentFunc (name) {
        console.debug('[tidb-sql-parser] define transparent func:', name)
        return __tidbSqlParserExecuteCmd(7, name)
    }


    close() {
        console.debug('[tidb-sql-parser] open close')
        return __tidbSqlParserExecuteCmd(2)
    }
}
