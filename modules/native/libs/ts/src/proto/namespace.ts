/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "./google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_namespace";

export interface Namespace {
  /** Unique name of the namespace */
  name: string;
  /** User-friendly public name of the namespace. */
  fullName: string;
  /** Arbitrary description */
  description: string;
  /** When the namespace was created */
  created:
    | Date
    | undefined;
  /** Last time when the namespace information was updated. */
  updated:
    | Date
    | undefined;
  /** Counter that increases after every update of the namespace */
  version: number;
}

export interface EnsureNamespaceRequest {
  /** Unique name of the namespace. Name must match the regex "[A-Za-z0-9]+$". Max length is 32 symbols. */
  name: string;
  /** User-friendly public name of the namespace. Can be changed in the future. Max length is 128 symbols. */
  fullName: string;
  /** Arbitrary description. Can be changed in the future. Max length is 512 symbols. */
  description: string;
}

export interface EnsureNamespaceResponse {
  /** Created namespace */
  namespace:
    | Namespace
    | undefined;
  /** If true - new namespace was created */
  created: boolean;
}

export interface CreateNamespaceRequest {
  /** Unique name of the namespace. Name must match the regex "[A-Za-z0-9]+$". Max length is 32 symbols. */
  name: string;
  /** User-friendly public name of the namespace. Can be changed in the future. Max length is 128 symbols. */
  fullName: string;
  /** Arbitrary description. Can be changed in the future. Max length is 512 symbols. */
  description: string;
}

export interface CreateNamespaceResponse {
  /** Created namespace */
  namespace: Namespace | undefined;
}

export interface UpdateNamespaceRequest {
  /** Unique name of the namespace. */
  name: string;
  /** User-friendly public name of the namespace. May be changed in the future. */
  fullName: string;
  /** Arbitrary description. May be changed in the future. */
  description: string;
}

export interface UpdateNamespaceResponse {
  /** Created namespace */
  namespace: Namespace | undefined;
}

export interface DeleteNamespaceRequest {
  /** Name of the namespace to delete */
  name: string;
}

export interface DeleteNamespaceResponse {
  /** Indicates if namespace existed before this operation */
  existed: boolean;
}

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

export interface GetNamespaceStatisticsRequest {
  /** Name of the namespace to get statistics. Required. */
  name: string;
  /** Cached data will be returned faster but may not be realtime. */
  useCache: boolean;
}

export interface GetNamespaceStatisticsResponse {
  /** Statistics from the database related to the namespace */
  db: GetNamespaceStatisticsResponse_Db | undefined;
}

export interface GetNamespaceStatisticsResponse_Db {
  /** Number ob objects stored in the database */
  objects: number;
  /** Total size of the raw data stored (without pre-alocated allocated space and indexes) */
  dataSize: number;
  /** Total memory usage of the namespace */
  totalSize: number;
}

function createBaseNamespace(): Namespace {
  return { name: "", fullName: "", description: "", created: undefined, updated: undefined, version: 0 };
}

export const Namespace = {
  encode(message: Namespace, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.fullName !== "") {
      writer.uint32(18).string(message.fullName);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(34).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(42).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(48).uint64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Namespace {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNamespace();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.description = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 48) {
            break;
          }

          message.version = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Namespace {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      description: isSet(object.description) ? String(object.description) : "",
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Namespace): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.created !== undefined) {
      obj.created = message.created.toISOString();
    }
    if (message.updated !== undefined) {
      obj.updated = message.updated.toISOString();
    }
    if (message.version !== 0) {
      obj.version = Math.round(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Namespace>, I>>(base?: I): Namespace {
    return Namespace.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Namespace>, I>>(object: I): Namespace {
    const message = createBaseNamespace();
    message.name = object.name ?? "";
    message.fullName = object.fullName ?? "";
    message.description = object.description ?? "";
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseEnsureNamespaceRequest(): EnsureNamespaceRequest {
  return { name: "", fullName: "", description: "" };
}

export const EnsureNamespaceRequest = {
  encode(message: EnsureNamespaceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.fullName !== "") {
      writer.uint32(18).string(message.fullName);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EnsureNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: EnsureNamespaceRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureNamespaceRequest>, I>>(base?: I): EnsureNamespaceRequest {
    return EnsureNamespaceRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureNamespaceRequest>, I>>(object: I): EnsureNamespaceRequest {
    const message = createBaseEnsureNamespaceRequest();
    message.name = object.name ?? "";
    message.fullName = object.fullName ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseEnsureNamespaceResponse(): EnsureNamespaceResponse {
  return { namespace: undefined, created: false };
}

export const EnsureNamespaceResponse = {
  encode(message: EnsureNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    if (message.created === true) {
      writer.uint32(16).bool(message.created);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = Namespace.decode(reader, reader.uint32());
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.created = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EnsureNamespaceResponse {
    return {
      namespace: isSet(object.namespace) ? Namespace.fromJSON(object.namespace) : undefined,
      created: isSet(object.created) ? Boolean(object.created) : false,
    };
  },

  toJSON(message: EnsureNamespaceResponse): unknown {
    const obj: any = {};
    if (message.namespace !== undefined) {
      obj.namespace = Namespace.toJSON(message.namespace);
    }
    if (message.created === true) {
      obj.created = message.created;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureNamespaceResponse>, I>>(base?: I): EnsureNamespaceResponse {
    return EnsureNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureNamespaceResponse>, I>>(object: I): EnsureNamespaceResponse {
    const message = createBaseEnsureNamespaceResponse();
    message.namespace = (object.namespace !== undefined && object.namespace !== null)
      ? Namespace.fromPartial(object.namespace)
      : undefined;
    message.created = object.created ?? false;
    return message;
  },
};

function createBaseCreateNamespaceRequest(): CreateNamespaceRequest {
  return { name: "", fullName: "", description: "" };
}

export const CreateNamespaceRequest = {
  encode(message: CreateNamespaceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.fullName !== "") {
      writer.uint32(18).string(message.fullName);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: CreateNamespaceRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateNamespaceRequest>, I>>(base?: I): CreateNamespaceRequest {
    return CreateNamespaceRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateNamespaceRequest>, I>>(object: I): CreateNamespaceRequest {
    const message = createBaseCreateNamespaceRequest();
    message.name = object.name ?? "";
    message.fullName = object.fullName ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseCreateNamespaceResponse(): CreateNamespaceResponse {
  return { namespace: undefined };
}

export const CreateNamespaceResponse = {
  encode(message: CreateNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = Namespace.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateNamespaceResponse {
    return { namespace: isSet(object.namespace) ? Namespace.fromJSON(object.namespace) : undefined };
  },

  toJSON(message: CreateNamespaceResponse): unknown {
    const obj: any = {};
    if (message.namespace !== undefined) {
      obj.namespace = Namespace.toJSON(message.namespace);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateNamespaceResponse>, I>>(base?: I): CreateNamespaceResponse {
    return CreateNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateNamespaceResponse>, I>>(object: I): CreateNamespaceResponse {
    const message = createBaseCreateNamespaceResponse();
    message.namespace = (object.namespace !== undefined && object.namespace !== null)
      ? Namespace.fromPartial(object.namespace)
      : undefined;
    return message;
  },
};

function createBaseUpdateNamespaceRequest(): UpdateNamespaceRequest {
  return { name: "", fullName: "", description: "" };
}

export const UpdateNamespaceRequest = {
  encode(message: UpdateNamespaceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.fullName !== "") {
      writer.uint32(18).string(message.fullName);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: UpdateNamespaceRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateNamespaceRequest>, I>>(base?: I): UpdateNamespaceRequest {
    return UpdateNamespaceRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateNamespaceRequest>, I>>(object: I): UpdateNamespaceRequest {
    const message = createBaseUpdateNamespaceRequest();
    message.name = object.name ?? "";
    message.fullName = object.fullName ?? "";
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseUpdateNamespaceResponse(): UpdateNamespaceResponse {
  return { namespace: undefined };
}

export const UpdateNamespaceResponse = {
  encode(message: UpdateNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = Namespace.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateNamespaceResponse {
    return { namespace: isSet(object.namespace) ? Namespace.fromJSON(object.namespace) : undefined };
  },

  toJSON(message: UpdateNamespaceResponse): unknown {
    const obj: any = {};
    if (message.namespace !== undefined) {
      obj.namespace = Namespace.toJSON(message.namespace);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateNamespaceResponse>, I>>(base?: I): UpdateNamespaceResponse {
    return UpdateNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateNamespaceResponse>, I>>(object: I): UpdateNamespaceResponse {
    const message = createBaseUpdateNamespaceResponse();
    message.namespace = (object.namespace !== undefined && object.namespace !== null)
      ? Namespace.fromPartial(object.namespace)
      : undefined;
    return message;
  },
};

function createBaseDeleteNamespaceRequest(): DeleteNamespaceRequest {
  return { name: "" };
}

export const DeleteNamespaceRequest = {
  encode(message: DeleteNamespaceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteNamespaceRequest {
    return { name: isSet(object.name) ? String(object.name) : "" };
  },

  toJSON(message: DeleteNamespaceRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteNamespaceRequest>, I>>(base?: I): DeleteNamespaceRequest {
    return DeleteNamespaceRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteNamespaceRequest>, I>>(object: I): DeleteNamespaceRequest {
    const message = createBaseDeleteNamespaceRequest();
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseDeleteNamespaceResponse(): DeleteNamespaceResponse {
  return { existed: false };
}

export const DeleteNamespaceResponse = {
  encode(message: DeleteNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.existed === true) {
      writer.uint32(8).bool(message.existed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.existed = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteNamespaceResponse {
    return { existed: isSet(object.existed) ? Boolean(object.existed) : false };
  },

  toJSON(message: DeleteNamespaceResponse): unknown {
    const obj: any = {};
    if (message.existed === true) {
      obj.existed = message.existed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteNamespaceResponse>, I>>(base?: I): DeleteNamespaceResponse {
    return DeleteNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteNamespaceResponse>, I>>(object: I): DeleteNamespaceResponse {
    const message = createBaseDeleteNamespaceResponse();
    message.existed = object.existed ?? false;
    return message;
  },
};

function createBaseGetNamespaceRequest(): GetNamespaceRequest {
  return { name: "", useCache: false };
}

export const GetNamespaceRequest = {
  encode(message: GetNamespaceRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
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

  fromJSON(object: any): GetNamespaceRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetNamespaceRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetNamespaceRequest>, I>>(base?: I): GetNamespaceRequest {
    return GetNamespaceRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetNamespaceRequest>, I>>(object: I): GetNamespaceRequest {
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
  encode(message: GetNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = Namespace.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetNamespaceResponse {
    return { namespace: isSet(object.namespace) ? Namespace.fromJSON(object.namespace) : undefined };
  },

  toJSON(message: GetNamespaceResponse): unknown {
    const obj: any = {};
    if (message.namespace !== undefined) {
      obj.namespace = Namespace.toJSON(message.namespace);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetNamespaceResponse>, I>>(base?: I): GetNamespaceResponse {
    return GetNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetNamespaceResponse>, I>>(object: I): GetNamespaceResponse {
    const message = createBaseGetNamespaceResponse();
    message.namespace = (object.namespace !== undefined && object.namespace !== null)
      ? Namespace.fromPartial(object.namespace)
      : undefined;
    return message;
  },
};

function createBaseGetAllNamespacesRequest(): GetAllNamespacesRequest {
  return { useCache: false };
}

export const GetAllNamespacesRequest = {
  encode(message: GetAllNamespacesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.useCache === true) {
      writer.uint32(8).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllNamespacesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllNamespacesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
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

  fromJSON(object: any): GetAllNamespacesRequest {
    return { useCache: isSet(object.useCache) ? Boolean(object.useCache) : false };
  },

  toJSON(message: GetAllNamespacesRequest): unknown {
    const obj: any = {};
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllNamespacesRequest>, I>>(base?: I): GetAllNamespacesRequest {
    return GetAllNamespacesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllNamespacesRequest>, I>>(object: I): GetAllNamespacesRequest {
    const message = createBaseGetAllNamespacesRequest();
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetAllNamespacesResponse(): GetAllNamespacesResponse {
  return { namespace: undefined };
}

export const GetAllNamespacesResponse = {
  encode(message: GetAllNamespacesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== undefined) {
      Namespace.encode(message.namespace, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllNamespacesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllNamespacesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = Namespace.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAllNamespacesResponse {
    return { namespace: isSet(object.namespace) ? Namespace.fromJSON(object.namespace) : undefined };
  },

  toJSON(message: GetAllNamespacesResponse): unknown {
    const obj: any = {};
    if (message.namespace !== undefined) {
      obj.namespace = Namespace.toJSON(message.namespace);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllNamespacesResponse>, I>>(base?: I): GetAllNamespacesResponse {
    return GetAllNamespacesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllNamespacesResponse>, I>>(object: I): GetAllNamespacesResponse {
    const message = createBaseGetAllNamespacesResponse();
    message.namespace = (object.namespace !== undefined && object.namespace !== null)
      ? Namespace.fromPartial(object.namespace)
      : undefined;
    return message;
  },
};

function createBaseIsNamespaceExistRequest(): IsNamespaceExistRequest {
  return { name: "", useCache: false };
}

export const IsNamespaceExistRequest = {
  encode(message: IsNamespaceExistRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsNamespaceExistRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsNamespaceExistRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
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

  fromJSON(object: any): IsNamespaceExistRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: IsNamespaceExistRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IsNamespaceExistRequest>, I>>(base?: I): IsNamespaceExistRequest {
    return IsNamespaceExistRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IsNamespaceExistRequest>, I>>(object: I): IsNamespaceExistRequest {
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
  encode(message: IsNamespaceExistResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IsNamespaceExistResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIsNamespaceExistResponse();
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

  fromJSON(object: any): IsNamespaceExistResponse {
    return { exist: isSet(object.exist) ? Boolean(object.exist) : false };
  },

  toJSON(message: IsNamespaceExistResponse): unknown {
    const obj: any = {};
    if (message.exist === true) {
      obj.exist = message.exist;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IsNamespaceExistResponse>, I>>(base?: I): IsNamespaceExistResponse {
    return IsNamespaceExistResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IsNamespaceExistResponse>, I>>(object: I): IsNamespaceExistResponse {
    const message = createBaseIsNamespaceExistResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

function createBaseGetNamespaceStatisticsRequest(): GetNamespaceStatisticsRequest {
  return { name: "", useCache: false };
}

export const GetNamespaceStatisticsRequest = {
  encode(message: GetNamespaceStatisticsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceStatisticsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceStatisticsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
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

  fromJSON(object: any): GetNamespaceStatisticsRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetNamespaceStatisticsRequest): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetNamespaceStatisticsRequest>, I>>(base?: I): GetNamespaceStatisticsRequest {
    return GetNamespaceStatisticsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetNamespaceStatisticsRequest>, I>>(
    object: I,
  ): GetNamespaceStatisticsRequest {
    const message = createBaseGetNamespaceStatisticsRequest();
    message.name = object.name ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetNamespaceStatisticsResponse(): GetNamespaceStatisticsResponse {
  return { db: undefined };
}

export const GetNamespaceStatisticsResponse = {
  encode(message: GetNamespaceStatisticsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.db !== undefined) {
      GetNamespaceStatisticsResponse_Db.encode(message.db, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceStatisticsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceStatisticsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.db = GetNamespaceStatisticsResponse_Db.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetNamespaceStatisticsResponse {
    return { db: isSet(object.db) ? GetNamespaceStatisticsResponse_Db.fromJSON(object.db) : undefined };
  },

  toJSON(message: GetNamespaceStatisticsResponse): unknown {
    const obj: any = {};
    if (message.db !== undefined) {
      obj.db = GetNamespaceStatisticsResponse_Db.toJSON(message.db);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetNamespaceStatisticsResponse>, I>>(base?: I): GetNamespaceStatisticsResponse {
    return GetNamespaceStatisticsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetNamespaceStatisticsResponse>, I>>(
    object: I,
  ): GetNamespaceStatisticsResponse {
    const message = createBaseGetNamespaceStatisticsResponse();
    message.db = (object.db !== undefined && object.db !== null)
      ? GetNamespaceStatisticsResponse_Db.fromPartial(object.db)
      : undefined;
    return message;
  },
};

function createBaseGetNamespaceStatisticsResponse_Db(): GetNamespaceStatisticsResponse_Db {
  return { objects: 0, dataSize: 0, totalSize: 0 };
}

export const GetNamespaceStatisticsResponse_Db = {
  encode(message: GetNamespaceStatisticsResponse_Db, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.objects !== 0) {
      writer.uint32(8).uint64(message.objects);
    }
    if (message.dataSize !== 0) {
      writer.uint32(16).uint64(message.dataSize);
    }
    if (message.totalSize !== 0) {
      writer.uint32(24).uint64(message.totalSize);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetNamespaceStatisticsResponse_Db {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetNamespaceStatisticsResponse_Db();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.objects = longToNumber(reader.uint64() as Long);
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.dataSize = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.totalSize = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetNamespaceStatisticsResponse_Db {
    return {
      objects: isSet(object.objects) ? Number(object.objects) : 0,
      dataSize: isSet(object.dataSize) ? Number(object.dataSize) : 0,
      totalSize: isSet(object.totalSize) ? Number(object.totalSize) : 0,
    };
  },

  toJSON(message: GetNamespaceStatisticsResponse_Db): unknown {
    const obj: any = {};
    if (message.objects !== 0) {
      obj.objects = Math.round(message.objects);
    }
    if (message.dataSize !== 0) {
      obj.dataSize = Math.round(message.dataSize);
    }
    if (message.totalSize !== 0) {
      obj.totalSize = Math.round(message.totalSize);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetNamespaceStatisticsResponse_Db>, I>>(
    base?: I,
  ): GetNamespaceStatisticsResponse_Db {
    return GetNamespaceStatisticsResponse_Db.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetNamespaceStatisticsResponse_Db>, I>>(
    object: I,
  ): GetNamespaceStatisticsResponse_Db {
    const message = createBaseGetNamespaceStatisticsResponse_Db();
    message.objects = object.objects ?? 0;
    message.dataSize = object.dataSize ?? 0;
    message.totalSize = object.totalSize ?? 0;
    return message;
  },
};

export interface NamespaceService {
  /** Create new namespace if it doesnt exist. If namespace exist, its data will not be updated. */
  Ensure(request: EnsureNamespaceRequest): Promise<EnsureNamespaceResponse>;
  /** Creates new namespace. If namespace already exist, will return error. */
  Create(request: CreateNamespaceRequest): Promise<CreateNamespaceResponse>;
  /** Updates namespace information. If namespace doesnt exist, will return error. */
  Update(request: UpdateNamespaceRequest): Promise<UpdateNamespaceResponse>;
  /** Deletes namespace and all its data */
  Delete(request: DeleteNamespaceRequest): Promise<DeleteNamespaceResponse>;
  /** Returns namespace information by its name */
  Get(request: GetNamespaceRequest): Promise<GetNamespaceResponse>;
  /** Streams list of all namespaces */
  GetAll(request: GetAllNamespacesRequest): Observable<GetAllNamespacesResponse>;
  /** Checks if namespace exists */
  Exists(request: IsNamespaceExistRequest): Promise<IsNamespaceExistResponse>;
  /** Gets namespace statistics. If namespace doesnt exist, will return error. */
  Stat(request: GetNamespaceStatisticsRequest): Promise<GetNamespaceStatisticsResponse>;
}

export const NamespaceServiceServiceName = "native_namespace.NamespaceService";
export class NamespaceServiceClientImpl implements NamespaceService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || NamespaceServiceServiceName;
    this.rpc = rpc;
    this.Ensure = this.Ensure.bind(this);
    this.Create = this.Create.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Get = this.Get.bind(this);
    this.GetAll = this.GetAll.bind(this);
    this.Exists = this.Exists.bind(this);
    this.Stat = this.Stat.bind(this);
  }
  Ensure(request: EnsureNamespaceRequest): Promise<EnsureNamespaceResponse> {
    const data = EnsureNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Ensure", data);
    return promise.then((data) => EnsureNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  Create(request: CreateNamespaceRequest): Promise<CreateNamespaceResponse> {
    const data = CreateNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateNamespaceRequest): Promise<UpdateNamespaceResponse> {
    const data = UpdateNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteNamespaceRequest): Promise<DeleteNamespaceResponse> {
    const data = DeleteNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetNamespaceRequest): Promise<GetNamespaceResponse> {
    const data = GetNamespaceRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  GetAll(request: GetAllNamespacesRequest): Observable<GetAllNamespacesResponse> {
    const data = GetAllNamespacesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "GetAll", data);
    return result.pipe(map((data) => GetAllNamespacesResponse.decode(_m0.Reader.create(data))));
  }

  Exists(request: IsNamespaceExistRequest): Promise<IsNamespaceExistResponse> {
    const data = IsNamespaceExistRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exists", data);
    return promise.then((data) => IsNamespaceExistResponse.decode(_m0.Reader.create(data)));
  }

  Stat(request: GetNamespaceStatisticsRequest): Promise<GetNamespaceStatisticsResponse> {
    const data = GetNamespaceStatisticsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Stat", data);
    return promise.then((data) => GetNamespaceStatisticsResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
  clientStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Promise<Uint8Array>;
  serverStreamingRequest(service: string, method: string, data: Uint8Array): Observable<Uint8Array>;
  bidirectionalStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Observable<Uint8Array>;
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

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = (t.seconds || 0) * 1_000;
  millis += (t.nanos || 0) / 1_000_000;
  return new Date(millis);
}

function fromJsonTimestamp(o: any): Date {
  if (o instanceof Date) {
    return o;
  } else if (typeof o === "string") {
    return new Date(o);
  } else {
    return fromTimestamp(Timestamp.fromJSON(o));
  }
}

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
