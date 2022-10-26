declare global {
    class Go {
        readonly importObject: WebAssembly.Imports
        readonly argv: string[]
        readonly env: Record<string, string>
        readonly exited: boolean
        readonly mem: DataView

        async run(instance: WebAssembly.Instance): Promise<void>

        exit (code: number)
    }
}

export {}
