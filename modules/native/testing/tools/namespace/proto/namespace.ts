/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

export const protobufPackage = "native_namespace";

export interface Namespace {
  /** Unique name of the namespace */
  name: string;
}

export interface EnsureNamespaceRequest {
  /** Unique name of the namespace */
  name: string;
}

export interface EnsureNamespaceResponse {
  /** Created namespace */
  namespace: Namespace | undefined;
}

export interface DeleteNamespaceRequest {
  /** Name of the namespace to delete */
  name: string;
}

export interface DeleteNamespaceResponse {}

export interface GetNamespaceRequest {
  /** Name of the namespace to get */
  name: string;
  /** Use cache or not. Cache have very small chance to be inconsisten on frequent read/writes operations to same namespace. Concurrent reads are safe. Inconsistent cache will be deleted after some period of time. */
  useCache: boolean;
}

export interface GetNamespaceResponse {
  namespace: Namespace | undefined;
}

export interface GetAllNamespacesRequest {
  /** Use cache or not. Cache have very small chance to be inconsisten on frequent read/writes operations to any namespace. Concurrent reads are safe. Inconsistent cache will be deleted after some period of time. */
  useCache: boolean;
}

export interface GetAllNamespacesResponse {
  namespace: Namespace | undefined;
}

export interface IsNamespaceExistRequest {
  /** Name of the namespace to get */
  name: string;
  /** Use cache or not. Cache have very small chance to be inconsisten on frequent read/writes operations to same namespace. Inconsistent cache will be deleted after some period of time. */
  useCache: boolean;
}

export interface IsNamespaceExistResponse {
  /** True if namespace exist, else - False */
  exist: boolean;
}

function createBaseNamespace(): Namespace {
  return { name: "" };
}

export const Namespace = {
  encode(
    message: Namespace,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Namespace {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNamespace();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Namespace {
    return {
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: Namespace): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Namespace>, I>>(
    object: I
  ): Namespace {
    const message = createBaseNamespace();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseEnsureNamespaceRequest(): EnsureNamespaceRequest {
  return { name: "" };
}

export const EnsureNamespaceRequest = {
  encode(
    message: EnsureNamespaceRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): EnsureNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EnsureNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: EnsureNamespaceRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EnsureNamespaceRequest>, I>>(
    object: I
  ): EnsureNamespaceRequest {
    const message = createBaseEnsureNamespaceRequest();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseEnsureNamespaceResponse(): EnsureNamespaceResponse {
  return { namespace: undefined };
}

export const EnsureNamespaceResponse = {
  encode(
    message: EnsureNamespaceResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): EnsureNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = Namespace.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): EnsureNamespaceResponse {
    return {
      namespace: isSet(object.namespace)
        ? Namespace.fromJSON(object.namespace)
        : undefined,
    };
  },

  toJSON(message: EnsureNamespaceResponse): unknown {
    const obj: any = {};
    message.namespace !== undefined &&
      (obj.namespace = message.namespace
        ? Namespace.toJSON(message.namespace)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<EnsureNamespaceResponse>, I>>(
    object: I
  ): EnsureNamespaceResponse {
    const message = createBaseEnsureNamespaceResponse();
    message.namespace =
      object.namespace !== undefined && object.namespace !== null
        ? Namespace.fromPartial(object.namespace)
        : undefined;
    return message;
  },
};

function createBaseDeleteNamespaceRequest(): DeleteNamespaceRequest {
  return { name: "" };
}

export const DeleteNamespaceRequest = {
  encode(
    message: DeleteNamespaceRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeleteNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeleteNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: DeleteNamespaceRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteNamespaceRequest>, I>>(
    object: I
  ): DeleteNamespaceRequest {
    const message = createBaseDeleteNamespaceRequest();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseDeleteNamespaceResponse(): DeleteNamespaceResponse {
  return {};
}

export const DeleteNamespaceResponse = {
  encode(
    _: DeleteNamespaceResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeleteNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteNamespaceResponse();
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

  fromJSON(_: any): DeleteNamespaceResponse {
    return {};
  },

  toJSON(_: DeleteNamespaceResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteNamespaceResponse>, I>>(
    _: I
  ): DeleteNamespaceResponse {
    const message = createBaseDeleteNamespaceResponse();
    return message;
  },
};

function createBaseGetNamespaceRequest(): GetNamespaceRequest {
  return { name: "", useCache: false };
}

export const GetNamespaceRequest = {
  encode(
    message: GetNamespaceRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetNamespaceRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetNamespaceRequest>, I>>(
    object: I
  ): GetNamespaceRequest {
    const message = createBaseGetNamespaceRequest();
    message.name = object.name ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetNamespaceResponse(): GetNamespaceResponse {
  return { namespace: undefined };
}

export const GetNamespaceResponse = {
  encode(
    message: GetNamespaceResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): GetNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = Namespace.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetNamespaceResponse {
    return {
      namespace: isSet(object.namespace)
        ? Namespace.fromJSON(object.namespace)
        : undefined,
    };
  },

  toJSON(message: GetNamespaceResponse): unknown {
    const obj: any = {};
    message.namespace !== undefined &&
      (obj.namespace = message.namespace
        ? Namespace.toJSON(message.namespace)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetNamespaceResponse>, I>>(
    object: I
  ): GetNamespaceResponse {
    const message = createBaseGetNamespaceResponse();
    message.namespace =
      object.namespace !== undefined && object.namespace !== null
        ? Namespace.fromPartial(object.namespace)
        : undefined;
    return message;
  },
};

function createBaseGetAllNamespacesRequest(): GetAllNamespacesRequest {
  return { useCache: false };
}

export const GetAllNamespacesRequest = {
  encode(
    message: GetAllNamespacesRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.useCache === true) {
      writer.uint32(8).bool(message.useCache);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): GetAllNamespacesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllNamespacesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetAllNamespacesRequest {
    return {
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetAllNamespacesRequest): unknown {
    const obj: any = {};
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetAllNamespacesRequest>, I>>(
    object: I
  ): GetAllNamespacesRequest {
    const message = createBaseGetAllNamespacesRequest();
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetAllNamespacesResponse(): GetAllNamespacesResponse {
  return { namespace: undefined };
}

export const GetAllNamespacesResponse = {
  encode(
    message: GetAllNamespacesResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): GetAllNamespacesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllNamespacesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = Namespace.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetAllNamespacesResponse {
    return {
      namespace: isSet(object.namespace)
        ? Namespace.fromJSON(object.namespace)
        : undefined,
    };
  },

  toJSON(message: GetAllNamespacesResponse): unknown {
    const obj: any = {};
    message.namespace !== undefined &&
      (obj.namespace = message.namespace
        ? Namespace.toJSON(message.namespace)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetAllNamespacesResponse>, I>>(
    object: I
  ): GetAllNamespacesResponse {
    const message = createBaseGetAllNamespacesResponse();
    message.namespace =
      object.namespace !== undefined && object.namespace !== null
        ? Namespace.fromPartial(object.namespace)
        : undefined;
    return message;
  },
};

function createBaseIsNamespaceExistRequest(): IsNamespaceExistRequest {
  return { name: "", useCache: false };
}

export const IsNamespaceExistRequest = {
  encode(
    message: IsNamespaceExistRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): IsNamespaceExistRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsNamespaceExistRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsNamespaceExistRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: IsNamespaceExistRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsNamespaceExistRequest>, I>>(
    object: I
  ): IsNamespaceExistRequest {
    const message = createBaseIsNamespaceExistRequest();
    message.name = object.name ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseIsNamespaceExistResponse(): IsNamespaceExistResponse {
  return { exist: false };
}

export const IsNamespaceExistResponse = {
  encode(
    message: IsNamespaceExistResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): IsNamespaceExistResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsNamespaceExistResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.exist = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): IsNamespaceExistResponse {
    return {
      exist: isSet(object.exist) ? Boolean(object.exist) : false,
    };
  },

  toJSON(message: IsNamespaceExistResponse): unknown {
    const obj: any = {};
    message.exist !== undefined && (obj.exist = message.exist);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<IsNamespaceExistResponse>, I>>(
    object: I
  ): IsNamespaceExistResponse {
    const message = createBaseIsNamespaceExistResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

export interface NamespaceService {
  /** Create new namespace if it doesnt exist */
  Ensure(request: EnsureNamespaceRequest): Promise<EnsureNamespaceResponse>;
  /** Deletes namespace and all its data */
  Delete(request: DeleteNamespaceRequest): Promise<DeleteNamespaceResponse>;
  /** Returns namespace information by its name */
  Get(request: GetNamespaceRequest): Promise<GetNamespaceResponse>;
  /** Streams list of all namespaces */
  GetAll(
    request: GetAllNamespacesRequest
  ): Observable<GetAllNamespacesResponse>;
  /** Checks if namespace exists */
  Exists(request: IsNamespaceExistRequest): Promise<IsNamespaceExistResponse>;
}

export class NamespaceServiceClientImpl implements NamespaceService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Ensure = this.Ensure.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Get = this.Get.bind(this);
    this.GetAll = this.GetAll.bind(this);
    this.Exists = this.Exists.bind(this);
  }
  Ensure(request: EnsureNamespaceRequest): Promise<EnsureNamespaceResponse> {
    const data = EnsureNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_namespace.NamespaceService",
      "Ensure",
      data
    );
    return promise.then((data) =>
      EnsureNamespaceResponse.decode(new _m0.Reader(data))
    );
  }

  Delete(request: DeleteNamespaceRequest): Promise<DeleteNamespaceResponse> {
    const data = DeleteNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_namespace.NamespaceService",
      "Delete",
      data
    );
    return promise.then((data) =>
      DeleteNamespaceResponse.decode(new _m0.Reader(data))
    );
  }

  Get(request: GetNamespaceRequest): Promise<GetNamespaceResponse> {
    const data = GetNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_namespace.NamespaceService",
      "Get",
      data
    );
    return promise.then((data) =>
      GetNamespaceResponse.decode(new _m0.Reader(data))
    );
  }

  GetAll(
    request: GetAllNamespacesRequest
  ): Observable<GetAllNamespacesResponse> {
    const data = GetAllNamespacesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(
      "native_namespace.NamespaceService",
      "GetAll",
      data
    );
    return result.pipe(
      map((data) => GetAllNamespacesResponse.decode(new _m0.Reader(data)))
    );
  }

  Exists(request: IsNamespaceExistRequest): Promise<IsNamespaceExistResponse> {
    const data = IsNamespaceExistRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_namespace.NamespaceService",
      "Exists",
      data
    );
    return promise.then((data) =>
      IsNamespaceExistResponse.decode(new _m0.Reader(data))
    );
  }
}

interface Rpc {
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
