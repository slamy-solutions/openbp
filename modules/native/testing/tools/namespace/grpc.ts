import { Client } from '../../../../system/libs/ts/grpc'

import { NamespaceServiceClientImpl } from './proto/namespace'

const grpcClient = new Client("native_namespace:80")
export const client = new NamespaceServiceClientImpl(grpcClient)

export async function connect() {
    await grpcClient.connect()
}

export async function close() {
    grpcClient.close()
}   