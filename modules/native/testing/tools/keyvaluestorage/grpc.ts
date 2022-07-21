import { Client } from '../../../../system/libs/ts/grpc'

import { KeyValueStorageServiceClientImpl } from './proto/keyvaluestorage'

const grpcClient = new Client("native_keyvaluestorage:80")
export const client = new KeyValueStorageServiceClientImpl(grpcClient)

export async function connect() {
    await grpcClient.connect()
}

export async function close() {
    grpcClient.close()
}