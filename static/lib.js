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

const SYMBOL_PARSER_ID = Symbol('Parser#id')

class Parser {
    [SYMBOL_PARSER_ID] = -1

    constructor() {
        this.open()
    }

    #cmd(...args) {
        const {Code, Data} = __tidbSqlParserExecuteCmd.apply(this[SYMBOL_PARSER_ID], args)
        if (Code !== 0) {
            throw new Error(Data)
        }
        return Data
    }

    get isOpen() {
        return this.#cmd(0)
    }

    open() {
        console.debug('[tidb-sql-parser] open parser')
        const res = this.#cmd(1)
        this[SYMBOL_PARSER_ID] = res
        return res
    }

    addDdl(sql) {
        console.debug('[tidb-sql-parser] add ddl:', sql)
        return this.#cmd(3, sql)
    }

    defineFunc(func, type, nullable) {
        console.debug('[tidb-sql-parser] define function:', func, type, nullable)
        return this.#cmd(4, func, evalTypes[type], nullable)
    }

    parse(sql) {
        console.debug('[tidb-sql-parser] parse:', sql)
        return this.#cmd(5, sql)
    }

    getTable(name) {
        console.debug('[tidb-sql-parser] get table:', name)
        return this.#cmd(6, name)
    }

    defineTransparentFunc(name) {
        console.debug('[tidb-sql-parser] define transparent func:', name)
        return this.#cmd(7, name)
    }

    close() {
        console.debug('[tidb-sql-parser] close')
        return this.#cmd(2)
    }

    static normalizeDigest (sql) {
        return __tidbSqlParserExecuteCmd(8, sql)
    }
}
