const go = new Go();
WebAssembly.instantiateStreaming(fetch("/main.wasm", { headers: { 'accept': 'application/wasm' }}), go.importObject).then((result) => {
    go.run(result.instance);
});
