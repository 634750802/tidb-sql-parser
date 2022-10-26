import {
    getType,
    newObject,
    ObjectHandle,
    ObjectType,
    reflectCallMethod,
    reflectDescribeType,
    reflectGetMethod,
    reflectGetProperty,
    reflectSetProperty,
    RuntimeSharedDataContext,
    setType
} from "./cmd.js";
import {registerAutoGc} from "./gc.js";
import * as styles from "./console.styles.js";

export const KEY_INTERNAL_ID = Symbol('JsBridgeType#internalId')
export const KEY_INITIALIZE = Symbol('JsBridgeType#initialize')
export const KEY_TYPE = Symbol('JsBridgeType#type')
const KEY_TP = Symbol('JsBridgeType#tp')
const KEY_CTX = Symbol('JsBridgeType#ctx')

export function native(): never {
    throw new Error('Native implementation not provided')
}

export type TypeInfo = {
    id: string
    otp: number
    size: number
    fields: Record<string, number[]>
    methods: Record<string, number>
}

export class BridgeTypeInfo {

    readonly construct: (ctx: RuntimeSharedDataContext, handle: ObjectHandle) => JsOpaqueBridgeType

    constructor(private tp: TypeInfo) {
        console.debug(`%c[wasm:debug]%c new type registered: %c${tp.id}`, styles.tag, '', styles.bold)
        this.construct = (ctx, handle) => new JsOpaqueBridgeType(ctx, handle)
    }

    containsField(name: string): boolean {
        return name in this.tp.fields
    }

    fieldIndex(name: string): number[] {
        const id = this.tp.fields[name]
        if (id === undefined) {
            throw new Error(`Unknown field '${name}' for type '${this.tp.id}'`)
        }
        return id
    }

    methodIndex(name: string): number {
        const id = this.tp.methods[name]
        if (id === undefined) {
            throw new Error(`Unknown method '${name}' for type '${this.tp.id}'`)
        }
        return id
    }

    fields() {
        return Object.keys(this.tp.fields)
    }

    methods() {
        return Object.keys(this.tp.methods)
    }

    get id(): string {
        return this.tp.id
    }
}

export default abstract class JsBridgeType {
    private [KEY_INTERNAL_ID]!: bigint
    private [KEY_TP]!: BridgeTypeInfo
    private [KEY_CTX]!: RuntimeSharedDataContext

    protected constructor(ctx: RuntimeSharedDataContext, objectHandle?: ObjectHandle) {
        Object.defineProperty(this, KEY_CTX, {
            value: ctx,
            writable: false,
            enumerable: false,
            configurable: false
        })
        if (objectHandle === undefined || (objectHandle & 0x00ffffffn) === 0n) {
            this[KEY_INITIALIZE](ctx, objectHandle)
        } else {
            this[KEY_INTERNAL_ID] = objectHandle
        }
        const id = this[KEY_INTERNAL_ID]
        Object.defineProperty(this, KEY_INTERNAL_ID, {
            value: id,
            writable: false,
            enumerable: false,
            configurable: false,
        })

        let tp = getType(ctx, id)
        if (!tp) {
            tp = new BridgeTypeInfo(reflectDescribeType(ctx, id))
            setType(ctx, id, tp)
        }
        Object.defineProperty(this, KEY_TP, {
            value: tp,
            writable: false,
            enumerable: false,
            configurable: false,
        })

        registerAutoGc(ctx, this)
        for (let f of this[KEY_TP].fields()) {
            field(f)(this, f)
        }
        for (let m of this[KEY_TP].methods()) {
            func(m)(this, m)
        }

    }

    [KEY_INITIALIZE](ctx: RuntimeSharedDataContext, type?: ObjectType) {
        this[KEY_INTERNAL_ID] = newObject(ctx, type ?? this[KEY_TYPE])
    }

    abstract get [KEY_TYPE](): ObjectType

    //
    // [KEY_KEYS](): string[] {
    //     return reflectGetPropertyNames(this[KEY_CTX], this[KEY_INTERNAL_ID])
    // }

    [Symbol.toStringTag]() {
        return `JsBridgeType<${this[KEY_TP].id}>@{${this[KEY_INTERNAL_ID] & 0x00ffffffn}}`
    }
}

class JsOpaqueBridgeType extends JsBridgeType {
    constructor(ctx: RuntimeSharedDataContext, objectHandle: ObjectHandle) {
        super(ctx, objectHandle);
    }

    get [KEY_TYPE](): ObjectType {
        return this[KEY_INTERNAL_ID] >> 24n;
    }

    toString() {
        return `JsBridgeType<${this[KEY_TP].id}>@{${this[KEY_INTERNAL_ID] & 0x00ffffffn}}`
    }
}

export function createJsBridgeType<T>(ctx: RuntimeSharedDataContext, h: ObjectHandle): T
export function createJsBridgeType(ctx: RuntimeSharedDataContext, h: ObjectHandle): any {
    let tp = getType(ctx, h)
    if (!tp) {
        tp = new BridgeTypeInfo(reflectDescribeType(ctx, h))
        setType(ctx, h, tp)
    }
    return tp.construct(ctx, h)
}

export function field(name: string) {
    return function decorate(target: JsBridgeType, key: string): void {
        Object.defineProperty(target, key, {
            configurable: false,
            enumerable: true,
            get(): any {
                if (this[KEY_TP].containsField(name)) {
                    return reflectGetProperty(this[KEY_CTX], this[KEY_INTERNAL_ID], this[KEY_TP].fieldIndex(name))
                } else {
                    return reflectGetMethod(this[KEY_CTX], this[KEY_INTERNAL_ID], this[KEY_TP].methodIndex(name))
                }
            },
            set(value): void {
                reflectSetProperty(this[KEY_CTX], this[KEY_INTERNAL_ID], this[KEY_TP].fieldIndex(name), value)
            }
        })
    }
}

export function func(name: string) {
    return function decorate(target: JsBridgeType, key: string): void {
        Object.defineProperty(target, key, {
            configurable: false,
            enumerable: false,
            value(...args: any[]) {
                return reflectCallMethod(this[KEY_CTX], this[KEY_INTERNAL_ID], this[KEY_TP].methodIndex(name), args)
            },
        })
    }
}

