// Define function types
/*
interface Tp {
    Type: EvalTypes,
    Nullable: boolean
}

interface Parser {
    DefineFunc(name: string, type: Tp): void

    DefineTransparentFunc(name: string): void
}

const enum EvalTypes {
    ETInt,
    ETReal,
    ETDecimal,
    ETString,
    ETDatetime,
    ETTimestamp,
    ETDuration,
    ETJson,
}
*/
p.DefineTransparentFunc("IFNULL")
p.DefineTransparentFunc("ABS")
p.DefineTransparentFunc("SUM")
p.DefineTransparentFunc("GREATEST")
p.DefineTransparentFunc("LEAST")
p.DefineFunc("DATE_SUB", {Type: EvalTypes.ETDatetime, Nullable: false})
p.DefineFunc("COUNT", {Type: EvalTypes.ETInt, Nullable: false})
p.DefineFunc("TIMESTAMPDIFF", {Type: EvalTypes.ETReal, Nullable: false})
