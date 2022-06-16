import { setup as setupNamespace, teardown as teardownNamespace } from './namespace'

export async function setup() {
    await setupNamespace()
}

export async function teardown() {
    await teardownNamespace()
}