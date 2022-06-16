import { Client } from '../../system/grpc'

import { NamespaceServiceClientImpl } from './proto/namespace'

const grpcClient = new Client("")
export const client = new NamespaceServiceClientImpl(grpcClient)

export async function connect() {
    await grpcClient.connect()
}

export async function close() {
    grpcClient.close()
}