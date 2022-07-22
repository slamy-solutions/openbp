import { Client } from '../../../../system/libs/ts/grpc'

import { IAMConfigServiceClientImpl } from './proto/iam'

const grpcConfigClient = new Client("native_iam_config:80")
export const configClient = new IAMConfigServiceClientImpl(grpcConfigClient)

export async function connect() {
    await grpcConfigClient.connect()
}

export async function close() {
    grpcConfigClient.close()
}