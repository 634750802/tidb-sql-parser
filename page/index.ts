import {initWasm} from "../js";

type Globals = {
    NewParser(): Parser

    readonly go: Go
}

export const enum EvalTypes {
    ETInt,
    ETReal,
    ETDecimal,
    ETString,
    ETDatetime,
    ETTimestamp,
    ETDuration,
    ETJson,
}

export const EvalTypeNames = [
    'ETInt',
    'ETReal',
    'ETDecimal',
    'ETString',
    'ETDatetime',
    'ETTimestamp',
    'ETDuration',
    'ETJson',
]

export interface Tp {
    Type: EvalTypes,
    Nullable: boolean
}

export interface Column extends Tp {
    Name: string
    As: string
}

export interface TableDefine {
    Name(): string

    Columns(): Column[]
}

export interface Parser {
    DefineFunc(name: string, type: Tp): void

    DefineTransparentFunc(name: string): void

    AddDdl(sql: string): void

    Parse(sql: string): Column[]

    GetTable(name: string): TableDefine

    Warns(): string[]
}

export interface Program {
    newParser(): Parser
    stop (): void
}

declare global {
    interface Window {
        EvalTypes: Record<string, number>
    }
}

window.EvalTypes = {
    ETInt: EvalTypes.ETInt,
    ETReal: EvalTypes.ETReal,
    ETDecimal: EvalTypes.ETDecimal,
    ETString: EvalTypes.ETString,
    ETDatetime: EvalTypes.ETDatetime,
    ETTimestamp: EvalTypes.ETTimestamp,
    ETDuration: EvalTypes.ETDuration,
    ETJson: EvalTypes.ETJson,
}

export async function init (): Promise<Program> {
    const Globals = await initWasm<Globals>()
    return {
        newParser () {
            return Globals.NewParser()
        },
        stop () {
            return Globals.go.exit(0)
        }
    }
}

