import { Client } from '../../../../system/libs/ts/grpc'

import { LambdaManagerServiceClientImpl, LambdaEntrypointServiceClientImpl } from './proto/lambda'

const grpcManagerClient = new Client("native_lambda_manager:80")
export const managerClient = new LambdaManagerServiceClientImpl(grpcManagerClient)
const grpcEntrypointClient = new Client("native_lambda_entrypoint:80")
export const entrypointClient = new LambdaEntrypointServiceClientImpl(grpcEntrypointClient)

export async function connect() {
    await grpcManagerClient.connect()
    await grpcEntrypointClient.connect()
}

export async function close() {
    grpcManagerClient.close()
    grpcEntrypointClient.close()
}