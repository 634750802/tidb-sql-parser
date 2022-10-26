import 'reflect-metadata';
import {RuntimeSharedDataContext} from "./wasm/cmd.js";
import {createJsBridgeType} from "./wasm/JsBridgeType.js";

export async function initWasm<T>() {
    const WASM_PREFIX = `__wasm_go_v1_${Date.now()}_`

    const go = new Go()
    go.env.WASM_PREFIX = WASM_PREFIX

    const result = await WebAssembly.instantiateStreaming(fetch('/tidb-sql-parser.wasm'), go.importObject)

    go.run(result.instance).catch((err) => {
        console.error('wasm panic:', err)
        console.error('current ctx:', ctx)
    })
    const ctx: RuntimeSharedDataContext = [0, 0n, [], undefined, undefined, []] as any;
    (globalThis as any)[`${WASM_PREFIX}init`](
        ctx,
        (h: number) => createJsBridgeType(ctx, BigInt(h))
    );
    ctx.cmd = (globalThis as any)[`${WASM_PREFIX}cmd`]

    const global = createJsBridgeType<T>(ctx, 0x00000000n)
    Object.defineProperty(global, 'go', {
        value: go,
        configurable: false,
        writable: false,
    })
    return global
}
