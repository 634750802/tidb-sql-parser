// H: 2B <type>; L: 6B <id>

import {BridgeTypeInfo, TypeInfo} from "./JsBridgeType.js";

export type ObjectHandle = bigint // int64
export type ObjectType = bigint
export type ObjectId = bigint

export const TYPE_MASK = 0xff000000n;
export const ID_MASK = 0x00ffffffn;

const CODE_MASK = TYPE_MASK;
const DATA_MASK = ID_MASK;

export const enum GeneralCmd {
    NEW,
    DELETE,
    ReflectGetProperty,
    ReflectGetMethod,
    ReflectSetProperty,
    ReflectCallMethod,
    ReflectGetPropertyNames,
    ReflectDescribeType,
}

const GENERAL_CMD_CODE_OK = 0x00000000n
const GENERAL_CMD_CODE_BAD_CMD = 0xff000000n

class GeneralCmdError extends Error {
    constructor(public code: bigint, message: string) {
        super(message);
    }
}

const enum RuntimeSharedDataContextIndex {
    cmd,
    handle,
    args,
    returnedValue,
    errorMessage,
    types,
}

export type RuntimeSharedDataContext = [
    cmd: GeneralCmd,
    handle: number,
    args: any[],
    /**
     * If general cmd returns some data, then it would be stored in ret.
     * You should read this value immediately after `sendGeneralCmd`
     * */
    returnedValue: unknown,

    /**
     * If general cmd failed, then it would be stored in ret.
     * You should read this value immediately after `sendGeneralCmd`
     */
    errorMessage: string | undefined,
    types: BridgeTypeInfo[],
] & {
    cmd(cmd: number): bigint
    debugString(): string
}

type CmdResult<T> = {
    id: bigint
    data: T
}

export function getType(ctx: RuntimeSharedDataContext, tp: ObjectType): BridgeTypeInfo | undefined {
    return ctx[RuntimeSharedDataContextIndex.types][Number(tp >> 24n)]
}

export function setType(ctx: RuntimeSharedDataContext, tp: ObjectType, typeInfo: BridgeTypeInfo) {
    ctx[RuntimeSharedDataContextIndex.types][Number(tp >> 24n)] = typeInfo
}


function sendGeneralCmd<T = void>(ctx: RuntimeSharedDataContext, handle: ObjectHandle, cmd: GeneralCmd, ...args: any[]): CmdResult<T> {
    ctx[RuntimeSharedDataContextIndex.cmd] = cmd
    ctx[RuntimeSharedDataContextIndex.handle] = Number(handle)
    ctx[RuntimeSharedDataContextIndex.args] = args

    const r = ctx.cmd(Number(cmd))
    if (!isFinite(r as any)) {
        console.error('ctx', ctx)
        throw new Error('Bad cmd return value')
    }
    const res = BigInt(r)

    const code = res & CODE_MASK
    if (code === GENERAL_CMD_CODE_OK) {
        return {
            id: res & ID_MASK,
            data: ctx[RuntimeSharedDataContextIndex.returnedValue] as T
        }
    }

    throw new GeneralCmdError(code >> 24n, ctx[RuntimeSharedDataContextIndex.errorMessage] ?? 'unknown message')
}

export function newObject(ctx: RuntimeSharedDataContext, type: ObjectType): ObjectHandle {
    type &= TYPE_MASK
    const {id} = sendGeneralCmd(ctx, type, GeneralCmd.NEW)
    return type | id
}

export function deleteObject(ctx: RuntimeSharedDataContext, handle: ObjectHandle): boolean {
    const id = handle & ID_MASK
    const {id: deletedId} = sendGeneralCmd(ctx, handle, GeneralCmd.DELETE)
    return deletedId === id && id !== 0n
}

export function reflectGetProperty<T = unknown>(ctx: RuntimeSharedDataContext, handle: ObjectHandle, index: number[]): T {
    const {data} = sendGeneralCmd<T>(ctx, handle, GeneralCmd.ReflectGetProperty, index)
    return data
}

export function reflectGetMethod<T = unknown>(ctx: RuntimeSharedDataContext, handle: ObjectHandle, index: number): T {
    const {data} = sendGeneralCmd<T>(ctx, handle, GeneralCmd.ReflectGetMethod, index)
    return data
}

export function reflectSetProperty(ctx: RuntimeSharedDataContext, handle: ObjectHandle, index: number[], value: any): void {
    const {data} = sendGeneralCmd(ctx, handle, GeneralCmd.ReflectSetProperty, index, value)
    return data
}

export function reflectCallMethod<T>(ctx: RuntimeSharedDataContext, handle: ObjectHandle, index: number, args: any[]): T {
    const {data} = sendGeneralCmd(ctx, handle, GeneralCmd.ReflectCallMethod, index, args)
    return data as T
}

export function reflectGetPropertyNames(ctx: RuntimeSharedDataContext, handle: ObjectHandle): string[] {
    const {data} = sendGeneralCmd<string[]>(ctx, handle, GeneralCmd.ReflectGetPropertyNames)
    return data
}

export function reflectDescribeType(ctx: RuntimeSharedDataContext, tp: ObjectType) {
    const {data} = sendGeneralCmd<TypeInfo>(ctx, tp, GeneralCmd.ReflectDescribeType)
    return data
}
