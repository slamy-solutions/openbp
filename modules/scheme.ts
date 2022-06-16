export interface Module {
    // Public module name
    get name(): string
    // Unique module identifier
    get uuid(): string
}

export interface Testing {
    // This function will be executed before all tests
    setup?: () => Promise<void>
    // This function will be executed after all tests
    teardown?: () => Promise<void>
}