import { Client } from '../../../../system/libs/ts/grpc'

import { ActorUserServiceClientImpl } from './proto/user'

const grpcClient = new Client("native_actor_user:80")
export const userClient = new ActorUserServiceClientImpl(grpcClient)

export async function connect() {
    await grpcClient.connect()
}

export async function close() {
    grpcClient.close()
}