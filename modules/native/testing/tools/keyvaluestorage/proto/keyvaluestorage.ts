/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_keyvalueprstorage";

export interface SetRequest {
  /** Namespace where to set key. Use empty for global key */
  namespace: string;
  /** Unique key that will be associated with value */
  key: string;
  /** Value to store. Dont use big value (maximum size is 15 mb). Values wich are more than 1mb in size will not be cached. */
  value: Buffer;
}

export interface SetResponse {}

export interface GetRequest {
  /** Namespace where to get value. Use empty for global value. */
  namespace: string;
  /** Key associated with value */
  key: string;
  /** Use cache or not. Cache may not be valid at very rare conditions (multiple write and read at same time), but will be invalidated after small period of time. */
  useCache: boolean;
}

export interface GetResponse {
  /** Value that was stored under specified namespace and key */
  value: Buffer;
}

function createBaseSetRequest(): SetRequest {
  return { namespace: "", key: "", value: Buffer.alloc(0) };
}

export const SetRequest = {
  encode(
    message: SetRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    if (message.value.length !== 0) {
      writer.uint32(26).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.key = reader.string();
          break;
        case 3:
          message.value = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value)
        ? Buffer.from(bytesFromBase64(object.value))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: SetRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.key !== undefined && (obj.key = message.key);
    message.value !== undefined &&
      (obj.value = base64FromBytes(
        message.value !== undefined ? message.value : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetRequest>, I>>(
    object: I
  ): SetRequest {
    const message = createBaseSetRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseSetResponse(): SetResponse {
  return {};
}

export const SetResponse = {
  encode(_: SetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(_: any): SetResponse {
    return {};
  },

  toJSON(_: SetResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetResponse>, I>>(_: I): SetResponse {
    const message = createBaseSetResponse();
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { namespace: "", key: "", useCache: false };
}

export const GetRequest = {
  encode(
    message: GetRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.key = reader.string();
          break;
        case 3:
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.key !== undefined && (obj.key = message.key);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(
    object: I
  ): GetRequest {
    const message = createBaseGetRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { value: Buffer.alloc(0) };
}

export const GetResponse = {
  encode(
    message: GetResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.value.length !== 0) {
      writer.uint32(10).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.value = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetResponse {
    return {
      value: isSet(object.value)
        ? Buffer.from(bytesFromBase64(object.value))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    message.value !== undefined &&
      (obj.value = base64FromBytes(
        message.value !== undefined ? message.value : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(
    object: I
  ): GetResponse {
    const message = createBaseGetResponse();
    message.value = object.value ?? Buffer.alloc(0);
    return message;
  },
};

/** Provides API for persistent key-value storage. Unlike system redis service it guaratees, that data will not be lost. Uses system_db to store data. Value+key size is limited to 15mb. */
export interface KeyValueStorageService {
  /** Sets value under the key in specified namespace. */
  Set(request: SetRequest): Promise<SetResponse>;
  /** Gets value for specified key. */
  Get(request: GetRequest): Promise<GetResponse>;
}

export class KeyValueStorageServiceClientImpl
  implements KeyValueStorageService
{
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Set = this.Set.bind(this);
    this.Get = this.Get.bind(this);
  }
  Set(request: SetRequest): Promise<SetResponse> {
    const data = SetRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_keyvalueprstorage.KeyValueStorageService",
      "Set",
      data
    );
    return promise.then((data) => SetResponse.decode(new _m0.Reader(data)));
  }

  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_keyvalueprstorage.KeyValueStorageService",
      "Get",
      data
    );
    return promise.then((data) => GetResponse.decode(new _m0.Reader(data)));
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  arr.forEach((byte) => {
    bin.push(String.fromCharCode(byte));
  });
  return btoa(bin.join(""));
}

type Builtin =
  | Date
  | Function
  | Uint8Array
  | string
  | number
  | boolean
  | undefined;

export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin
  ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & Record<
        Exclude<keyof I, KeysOfUnion<P>>,
        never
      >;

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
