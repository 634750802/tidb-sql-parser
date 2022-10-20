# TiDB SQL Parser

***WIP*** This project provide a wasm library for parsing TiDB SQL.
The wasm module will create a global function `__tidbSqlParse(sql: string): Ast[]`.

## Prepare

- go 1.18
- node

It's recommended to set GO path in this directory.(e.g. `.go`)

```shell
# Copy go wasm runtime script to static dir
npm run prepare:go-wasm-runtime

# Apply some patches to go mods to support wasm
npm run prebuild
```

## Build

```shell
npm run build
```

## Preview

```shell
# Install a simple web server `serve`
npm i

# Normally, the http server would listen 3000 if available
npm start
```

Visit `http://127.0.0.1:3000/` and have fun!

Open browser console and type `__tidbSqlParse('SELECT * FROM TiDB')`
