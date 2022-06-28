import { Client } from '../../../../system/libs/ts/grpc'

import { FileServiceClientImpl } from './proto/file'

const grpcClient = new Client("native_file:80")
export const client = new FileServiceClientImpl(grpcClient)

export async function connect() {
    await grpcClient.connect()
}

export async function close() {
    grpcClient.close()
}