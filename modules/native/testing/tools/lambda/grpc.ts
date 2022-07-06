import { Client } from '../../../../system/libs/ts/grpc'

import { LambdaManagerServiceClientImpl } from './proto/lambda'

const grpcManagerClient = new Client("native_lambda_manager:80")
export const managerClient = new LambdaManagerServiceClientImpl(grpcManagerClient)

export async function connect() {
    await grpcManagerClient.connect()
}

export async function close() {
    grpcManagerClient.close()
}