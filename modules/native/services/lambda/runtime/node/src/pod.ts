export interface ExecutionParameters {
    lambda: any
    data: Buffer
    setupData: unknown
}

export interface Module {
    setUp(): Promise<unknown>;
    tearDown(): Promise<void>;
    execute(params: ExecutionParameters): Promise<Buffer>;
}

export class Pod {
    protected scriptModule: Module
    protected setUpData: unknown = undefined


    constructor(code: string) {
        this.scriptModule = Function(code)() as Module
        if (!(this.scriptModule && this.scriptModule.execute === undefined)) {
            throw new Error("Bad script formatting. Execute function for found.")
        }
    }

    async setUp() {
        if (this.scriptModule.setUp) {
            this.setUpData = await this.scriptModule.setUp()
        }
    }

    async tearDown() {
        if (this.scriptModule.tearDown) {
            await this.scriptModule.tearDown()
        }
    }

    async execute(params: ExecutionParameters) {
        await this.scriptModule.execute(params)
    }
}