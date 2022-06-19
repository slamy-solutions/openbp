import { Client as GRPCClient, requestCallback, ServiceError } from '@grpc/grpc-js'
import { Observable } from 'rxjs'

export interface Rpc {
    request(
        service: string,
        method: string,
        data: Uint8Array
    ): Promise<Uint8Array>;
    clientStreamingRequest(
        service: string,
        method: string,
        data: Observable<Uint8Array>
    ): Promise<Uint8Array>;
    serverStreamingRequest(
        service: string,
        method: string,
        data: Uint8Array
    ): Observable<Uint8Array>;
    bidirectionalStreamingRequest(
        service: string,
        method: string,
        data: Observable<Uint8Array>
    ): Observable<Uint8Array>;
}

export type RequestError = ServiceError

export class Client implements Rpc {
    private client: GRPCClient

    constructor(client: GRPCClient) {
        this.client = client
    }

    request(service: string, method: string, data: Uint8Array): Promise<Uint8Array> {
        const path = `${service}/${method}`
        return new Promise((resolve, reject) => {
            const resultCallback: requestCallback<any> = (err, res) => {
              if (err) {
                return reject(err)
              }
              resolve(res)
            };
        
            function passThrough(argument: any) {
              return argument
            }
            this.client.makeUnaryRequest(path, passThrough, passThrough, data, resultCallback)
        });
    }
    clientStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Promise<Uint8Array> {
        const path = `${service}/${method}`
        
        return new Promise(async (resolve, reject) => {
            const resultCallback: requestCallback<any> = (err, res) => {
                if (err) {
                  return reject(err)
                }
                resolve(res)
              };
            
            function passThrough(argument: any) {
                return argument
            }
            const stream = this.client.makeClientStreamRequest(path, passThrough, passThrough, resultCallback)
            data.subscribe({
                next: stream.write,
                error: stream.destroy,
                complete: stream.end
            })
        })
    }
    serverStreamingRequest(service: string, method: string, data: Uint8Array): Observable<Uint8Array> {
        const path = `${service}/${method}`

        function passThrough(argument: any) {
            return argument
        }
        return new Observable((obs) => {
            const stream = this.client.makeServerStreamRequest(path, passThrough, passThrough, data)
            stream.on('data', (data) => obs.next(data))
            stream.on('error', (err) => obs.error(err))
            stream.on('end', () => obs.complete())
        })
    }
    bidirectionalStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Observable<Uint8Array> {
        const path = `${service}/${method}`

        function passThrough(argument: any) {
            return argument
        }

        return new Observable((obs) => {
            const stream = this.client.makeBidiStreamRequest(path, passThrough, passThrough)
            stream.on('data', obs.next)
            stream.on('error', obs.error)
            stream.on('end', obs.complete)
            data.subscribe({
                next: stream.write,
                error: stream.destroy,
                complete: stream.end
            })
        })
    }

}