import fsp from 'fs/promises';
import crypto from 'crypto';

globalThis.crypto = crypto.webcrypto

await import('./static/wasm_exec.js');

const go = new Go();

const wasmBuffer = await fsp.readFile('static/main.wasm');
WebAssembly.instantiate(wasmBuffer, go.importObject).then(result => {
    go.run(result.instance);
    console.log(globalThis.__tidbSqlParse); // Outputs: 11

});
