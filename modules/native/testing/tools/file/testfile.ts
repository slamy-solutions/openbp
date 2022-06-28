import { randomBytes } from 'crypto'
import { Observable } from 'rxjs'

import { FileCreateRequest } from './proto/file'
import { client } from './grpc'
import { ObjectId } from 'mongodb'

const SEND_CHUNK_SIZE = 32000

export class TestFile {
    private uuid: string
    private buf: Buffer
    private namespace: string
    private disableCache: boolean
    private forceCaching: boolean
    private mimeType: string
    private readonly: boolean

    constructor(size: number, namespace: string, disableCache: boolean, forceCaching: boolean, mimeType: string, readonly: boolean) {
        this.uuid = ""
        this.buf = randomBytes(size)
        this.namespace = namespace
        this.disableCache = disableCache
        this.forceCaching = forceCaching
        this.mimeType = mimeType
        this.readonly = readonly
    }

    get data() {
        return this.buf
    }

    get UUID() {
        return this.uuid
    }

    get mongoId() {
        return ObjectId.createFromHexString(this.UUID)
    }

    get obs(): Observable<FileCreateRequest> {
        return new Observable<FileCreateRequest>((subscriber) => {
            subscriber.next({info: {
                namespace: this.namespace,
                disableCache: this.disableCache,
                forceCaching: this.forceCaching,
                mimeType: this.mimeType,
                readonly: this.readonly
            },
            chunk: undefined
            })

            for (let index = 0; index < this.buf.length; index += SEND_CHUNK_SIZE) {
                subscriber.next({
                    chunk: {
                        data: this.buf.slice(index, Math.min(index + SEND_CHUNK_SIZE, this.buf.length))
                    },
                    info: undefined
                })
            }
        })
    }

    async with<T>(call: (file: TestFile) => Promise<T>) {
        const createResponse = await this.create()
        this.uuid = createResponse.file?.uuid as string
        try {
            return await call(this)
        } finally {
            await client.Delete({ namespace: this.namespace, uuid: this.uuid })
        }
    }

    async create() {
        const inputStream = new Observable<FileCreateRequest>((subscriber) => {
            subscriber.next({info: {
                namespace: this.namespace,
                disableCache: this.disableCache,
                forceCaching: this.forceCaching,
                mimeType: this.mimeType,
                readonly: this.readonly
            },
            chunk: undefined
            })

            for (let i = 0; i < this.buf.length; i+= 32000) {
                subscriber.next({
                    chunk: {
                        data: this.buf.slice(i, Math.min(i + 32000, this.buf.length))
                    },
                    info: undefined
                })
            }
        })
        return client.Create(inputStream)
    }
}