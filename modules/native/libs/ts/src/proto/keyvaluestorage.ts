/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_keyvalueprstorage";

export interface SetRequest {
  /** Namespace where to set key. Use empty for global key */
  namespace: string;
  /** Unique key that will be associated with value */
  key: string;
  /** Value to store. Dont use big value (maximum size is 15 mb). Values wich are more than 1mb in size will not be cached. */
  value: Uint8Array;
}

export interface SetResponse {
}

export interface SetIfNotExistRequest {
  /** Namespace where to set key. Use empty for global key */
  namespace: string;
  /** Unique key that will be associated with value */
  key: string;
  /** Value to store. Dont use big value (maximum size is 15 mb). Values wich are more than 1mb in size will not be cached. */
  value: Uint8Array;
}

export interface SetIfNotExistResponse {
  /** Indicates if value was seted or not */
  seted: boolean;
}

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
  value: Uint8Array;
}

export interface RemoveRequest {
  /** Namespace where to remove key */
  namespace: string;
  /** Key to remove */
  key: string;
}

export interface RemoveResponse {
  /** Indicates if actually removed key (or it wasnt exist in before the request) */
  removed: boolean;
}

export interface ExistRequest {
  /** Namespace where to search for key */
  namespace: string;
  /** Key to search */
  key: string;
  /** Use cache or not. Cache may not be valid at very rare conditions (multiple write and read at same time), but will be invalidated after small period of time. */
  useCache: boolean;
}

export interface ExistResponse {
  /** Indicates if namespace and key exists or not. */
  exist: boolean;
}

function createBaseSetRequest(): SetRequest {
  return { namespace: "", key: "", value: new Uint8Array(0) };
}

export const SetRequest = {
  encode(message: SetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.value = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? bytesFromBase64(object.value) : new Uint8Array(0),
    };
  },

  toJSON(message: SetRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value.length !== 0) {
      obj.value = base64FromBytes(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetRequest>, I>>(base?: I): SetRequest {
    return SetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetRequest>, I>>(object: I): SetRequest {
    const message = createBaseSetRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.value = object.value ?? new Uint8Array(0);
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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

  create<I extends Exact<DeepPartial<SetResponse>, I>>(base?: I): SetResponse {
    return SetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetResponse>, I>>(_: I): SetResponse {
    const message = createBaseSetResponse();
    return message;
  },
};

function createBaseSetIfNotExistRequest(): SetIfNotExistRequest {
  return { namespace: "", key: "", value: new Uint8Array(0) };
}

export const SetIfNotExistRequest = {
  encode(message: SetIfNotExistRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): SetIfNotExistRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIfNotExistRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.value = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetIfNotExistRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
      value: isSet(object.value) ? bytesFromBase64(object.value) : new Uint8Array(0),
    };
  },

  toJSON(message: SetIfNotExistRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.value.length !== 0) {
      obj.value = base64FromBytes(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetIfNotExistRequest>, I>>(base?: I): SetIfNotExistRequest {
    return SetIfNotExistRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetIfNotExistRequest>, I>>(object: I): SetIfNotExistRequest {
    const message = createBaseSetIfNotExistRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.value = object.value ?? new Uint8Array(0);
    return message;
  },
};

function createBaseSetIfNotExistResponse(): SetIfNotExistResponse {
  return { seted: false };
}

export const SetIfNotExistResponse = {
  encode(message: SetIfNotExistResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.seted === true) {
      writer.uint32(8).bool(message.seted);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetIfNotExistResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIfNotExistResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.seted = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetIfNotExistResponse {
    return { seted: isSet(object.seted) ? Boolean(object.seted) : false };
  },

  toJSON(message: SetIfNotExistResponse): unknown {
    const obj: any = {};
    if (message.seted === true) {
      obj.seted = message.seted;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetIfNotExistResponse>, I>>(base?: I): SetIfNotExistResponse {
    return SetIfNotExistResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetIfNotExistResponse>, I>>(object: I): SetIfNotExistResponse {
    const message = createBaseSetIfNotExistResponse();
    message.seted = object.seted ?? false;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { namespace: "", key: "", useCache: false };
}

export const GetRequest = {
  encode(message: GetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRequest>, I>>(base?: I): GetRequest {
    return GetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(object: I): GetRequest {
    const message = createBaseGetRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { value: new Uint8Array(0) };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.value.length !== 0) {
      writer.uint32(10).bytes(message.value);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.value = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetResponse {
    return { value: isSet(object.value) ? bytesFromBase64(object.value) : new Uint8Array(0) };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.value.length !== 0) {
      obj.value = base64FromBytes(message.value);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.value = object.value ?? new Uint8Array(0);
    return message;
  },
};

function createBaseRemoveRequest(): RemoveRequest {
  return { namespace: "", key: "" };
}

export const RemoveRequest = {
  encode(message: RemoveRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.key !== "") {
      writer.uint32(18).string(message.key);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
    };
  },

  toJSON(message: RemoveRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.key !== "") {
      obj.key = message.key;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveRequest>, I>>(base?: I): RemoveRequest {
    return RemoveRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveRequest>, I>>(object: I): RemoveRequest {
    const message = createBaseRemoveRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    return message;
  },
};

function createBaseRemoveResponse(): RemoveResponse {
  return { removed: false };
}

export const RemoveResponse = {
  encode(message: RemoveResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.removed === true) {
      writer.uint32(8).bool(message.removed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.removed = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveResponse {
    return { removed: isSet(object.removed) ? Boolean(object.removed) : false };
  },

  toJSON(message: RemoveResponse): unknown {
    const obj: any = {};
    if (message.removed === true) {
      obj.removed = message.removed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveResponse>, I>>(base?: I): RemoveResponse {
    return RemoveResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveResponse>, I>>(object: I): RemoveResponse {
    const message = createBaseRemoveResponse();
    message.removed = object.removed ?? false;
    return message;
  },
};

function createBaseExistRequest(): ExistRequest {
  return { namespace: "", key: "", useCache: false };
}

export const ExistRequest = {
  encode(message: ExistRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.key = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExistRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      key: isSet(object.key) ? String(object.key) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ExistRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.key !== "") {
      obj.key = message.key;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistRequest>, I>>(base?: I): ExistRequest {
    return ExistRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistRequest>, I>>(object: I): ExistRequest {
    const message = createBaseExistRequest();
    message.namespace = object.namespace ?? "";
    message.key = object.key ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseExistResponse(): ExistResponse {
  return { exist: false };
}

export const ExistResponse = {
  encode(message: ExistResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.exist = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExistResponse {
    return { exist: isSet(object.exist) ? Boolean(object.exist) : false };
  },

  toJSON(message: ExistResponse): unknown {
    const obj: any = {};
    if (message.exist === true) {
      obj.exist = message.exist;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistResponse>, I>>(base?: I): ExistResponse {
    return ExistResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistResponse>, I>>(object: I): ExistResponse {
    const message = createBaseExistResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

/** Provides API for persistent key-value storage. Unlike system redis service it guaratees, that data will not be lost. Uses system_db to store data. Value+key size is limited to 15mb. */
export interface KeyValueStorageService {
  /** Sets value under the key in specified namespace. */
  Set(request: SetRequest): Promise<SetResponse>;
  /** Sets value under the key in specified namespace only if it is not set. */
  SetIfNotExist(request: SetIfNotExistRequest): Promise<SetIfNotExistResponse>;
  /** Gets value for specified key. */
  Get(request: GetRequest): Promise<GetResponse>;
  /** Remove key in specified namespace */
  Remove(request: RemoveRequest): Promise<RemoveResponse>;
  /** Checks if key exists in specified namespace */
  Exist(request: ExistRequest): Promise<ExistResponse>;
}

export const KeyValueStorageServiceServiceName = "native_keyvalueprstorage.KeyValueStorageService";
export class KeyValueStorageServiceClientImpl implements KeyValueStorageService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || KeyValueStorageServiceServiceName;
    this.rpc = rpc;
    this.Set = this.Set.bind(this);
    this.SetIfNotExist = this.SetIfNotExist.bind(this);
    this.Get = this.Get.bind(this);
    this.Remove = this.Remove.bind(this);
    this.Exist = this.Exist.bind(this);
  }
  Set(request: SetRequest): Promise<SetResponse> {
    const data = SetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Set", data);
    return promise.then((data) => SetResponse.decode(_m0.Reader.create(data)));
  }

  SetIfNotExist(request: SetIfNotExistRequest): Promise<SetIfNotExistResponse> {
    const data = SetIfNotExistRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "SetIfNotExist", data);
    return promise.then((data) => SetIfNotExistResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetResponse.decode(_m0.Reader.create(data)));
  }

  Remove(request: RemoveRequest): Promise<RemoveResponse> {
    const data = RemoveRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Remove", data);
    return promise.then((data) => RemoveResponse.decode(_m0.Reader.create(data)));
  }

  Exist(request: ExistRequest): Promise<ExistResponse> {
    const data = ExistRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exist", data);
    return promise.then((data) => ExistResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
}

declare const self: any | undefined;
declare const window: any | undefined;
declare const global: any | undefined;
const tsProtoGlobalThis: any = (() => {
  if (typeof globalThis !== "undefined") {
    return globalThis;
  }
  if (typeof self !== "undefined") {
    return self;
  }
  if (typeof window !== "undefined") {
    return window;
  }
  if (typeof global !== "undefined") {
    return global;
  }
  throw "Unable to locate global object";
})();

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
