import { InternalOpenBPStub } from '@openbp/sdk'

export class Executor {
    private api: InternalOpenBPStub
    private _running = false;
    get running() {
        return this._running;
    }

    readonly runtimeNamespace: string
    readonly runtimeName: string;
    
    constructor(api: InternalOpenBPStub, runtimeNamespace: string, runtimeName: string) {
        this.api = api;
        this.runtimeNamespace = runtimeNamespace;
        this.runtimeName = runtimeName;
    }

    async load() {

    }

    async start() {

    }

    async stop() {

    }

    async dispose() {

    }
}