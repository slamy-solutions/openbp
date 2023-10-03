/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

export const protobufPackage = "runtime_manager_runtime";

export interface Runtime {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime. Unique within namespace */
  name: string;
  /** Should runtime be running */
  run: boolean;
}

export interface GetRuntimesForNamespaceReqeust {
  /** Namespace where runtimes are located */
  namespace: string;
}

export interface GetRuntimesForNamespaceResponse {
  /** Runtimes in namespace */
  runtimes: Runtime[];
}

export interface GetRuntimeRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime */
  name: string;
}

export interface GetRuntimeResponse {
  /** Runtime */
  runtime: Runtime | undefined;
}

export interface CreateRuntimeRequest {
  /** Runtime to create */
  runtime: Runtime | undefined;
}

export interface CreateRuntimeResponse {
}

export interface UpdateRuntimeRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime */
  name: string;
  newRun: boolean;
}

export interface UpdateRuntimeResponse {
  /** Updated runtime */
  runtime: Runtime | undefined;
}

export interface DeleteRuntimeRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime */
  name: string;
}

export interface DeleteRuntimeResponse {
}

export interface UploadRuntimeBinaryRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime */
  name: string;
  /** Chunk of binary data */
  binary: Uint8Array;
}

export interface UploadRuntimeBinaryResponse {
}

export interface DownloadRuntimeBinaryRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of the runtime */
  name: string;
}

export interface DownloadRuntimeBinaryResponse {
  /** Chunk of binary data */
  binary: Uint8Array;
}

function createBaseRuntime(): Runtime {
  return { namespace: "", name: "", run: false };
}

export const Runtime = {
  encode(message: Runtime, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.run === true) {
      writer.uint32(24).bool(message.run);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Runtime {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRuntime();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.run = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Runtime {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      run: isSet(object.run) ? Boolean(object.run) : false,
    };
  },

  toJSON(message: Runtime): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.run === true) {
      obj.run = message.run;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Runtime>, I>>(base?: I): Runtime {
    return Runtime.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Runtime>, I>>(object: I): Runtime {
    const message = createBaseRuntime();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.run = object.run ?? false;
    return message;
  },
};

function createBaseGetRuntimesForNamespaceReqeust(): GetRuntimesForNamespaceReqeust {
  return { namespace: "" };
}

export const GetRuntimesForNamespaceReqeust = {
  encode(message: GetRuntimesForNamespaceReqeust, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRuntimesForNamespaceReqeust {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRuntimesForNamespaceReqeust();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRuntimesForNamespaceReqeust {
    return { namespace: isSet(object.namespace) ? String(object.namespace) : "" };
  },

  toJSON(message: GetRuntimesForNamespaceReqeust): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRuntimesForNamespaceReqeust>, I>>(base?: I): GetRuntimesForNamespaceReqeust {
    return GetRuntimesForNamespaceReqeust.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRuntimesForNamespaceReqeust>, I>>(
    object: I,
  ): GetRuntimesForNamespaceReqeust {
    const message = createBaseGetRuntimesForNamespaceReqeust();
    message.namespace = object.namespace ?? "";
    return message;
  },
};

function createBaseGetRuntimesForNamespaceResponse(): GetRuntimesForNamespaceResponse {
  return { runtimes: [] };
}

export const GetRuntimesForNamespaceResponse = {
  encode(message: GetRuntimesForNamespaceResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.runtimes) {
      Runtime.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRuntimesForNamespaceResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRuntimesForNamespaceResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.runtimes.push(Runtime.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRuntimesForNamespaceResponse {
    return { runtimes: Array.isArray(object?.runtimes) ? object.runtimes.map((e: any) => Runtime.fromJSON(e)) : [] };
  },

  toJSON(message: GetRuntimesForNamespaceResponse): unknown {
    const obj: any = {};
    if (message.runtimes?.length) {
      obj.runtimes = message.runtimes.map((e) => Runtime.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRuntimesForNamespaceResponse>, I>>(base?: I): GetRuntimesForNamespaceResponse {
    return GetRuntimesForNamespaceResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRuntimesForNamespaceResponse>, I>>(
    object: I,
  ): GetRuntimesForNamespaceResponse {
    const message = createBaseGetRuntimesForNamespaceResponse();
    message.runtimes = object.runtimes?.map((e) => Runtime.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetRuntimeRequest(): GetRuntimeRequest {
  return { namespace: "", name: "" };
}

export const GetRuntimeRequest = {
  encode(message: GetRuntimeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRuntimeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRuntimeRequest();
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

  fromJSON(object: any): GetRuntimeRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: GetRuntimeRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRuntimeRequest>, I>>(base?: I): GetRuntimeRequest {
    return GetRuntimeRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRuntimeRequest>, I>>(object: I): GetRuntimeRequest {
    const message = createBaseGetRuntimeRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseGetRuntimeResponse(): GetRuntimeResponse {
  return { runtime: undefined };
}

export const GetRuntimeResponse = {
  encode(message: GetRuntimeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.runtime !== undefined) {
      Runtime.encode(message.runtime, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRuntimeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRuntimeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.runtime = Runtime.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRuntimeResponse {
    return { runtime: isSet(object.runtime) ? Runtime.fromJSON(object.runtime) : undefined };
  },

  toJSON(message: GetRuntimeResponse): unknown {
    const obj: any = {};
    if (message.runtime !== undefined) {
      obj.runtime = Runtime.toJSON(message.runtime);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRuntimeResponse>, I>>(base?: I): GetRuntimeResponse {
    return GetRuntimeResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRuntimeResponse>, I>>(object: I): GetRuntimeResponse {
    const message = createBaseGetRuntimeResponse();
    message.runtime = (object.runtime !== undefined && object.runtime !== null)
      ? Runtime.fromPartial(object.runtime)
      : undefined;
    return message;
  },
};

function createBaseCreateRuntimeRequest(): CreateRuntimeRequest {
  return { runtime: undefined };
}

export const CreateRuntimeRequest = {
  encode(message: CreateRuntimeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.runtime !== undefined) {
      Runtime.encode(message.runtime, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRuntimeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRuntimeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.runtime = Runtime.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateRuntimeRequest {
    return { runtime: isSet(object.runtime) ? Runtime.fromJSON(object.runtime) : undefined };
  },

  toJSON(message: CreateRuntimeRequest): unknown {
    const obj: any = {};
    if (message.runtime !== undefined) {
      obj.runtime = Runtime.toJSON(message.runtime);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRuntimeRequest>, I>>(base?: I): CreateRuntimeRequest {
    return CreateRuntimeRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRuntimeRequest>, I>>(object: I): CreateRuntimeRequest {
    const message = createBaseCreateRuntimeRequest();
    message.runtime = (object.runtime !== undefined && object.runtime !== null)
      ? Runtime.fromPartial(object.runtime)
      : undefined;
    return message;
  },
};

function createBaseCreateRuntimeResponse(): CreateRuntimeResponse {
  return {};
}

export const CreateRuntimeResponse = {
  encode(_: CreateRuntimeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRuntimeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRuntimeResponse();
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

  fromJSON(_: any): CreateRuntimeResponse {
    return {};
  },

  toJSON(_: CreateRuntimeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRuntimeResponse>, I>>(base?: I): CreateRuntimeResponse {
    return CreateRuntimeResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRuntimeResponse>, I>>(_: I): CreateRuntimeResponse {
    const message = createBaseCreateRuntimeResponse();
    return message;
  },
};

function createBaseUpdateRuntimeRequest(): UpdateRuntimeRequest {
  return { namespace: "", name: "", newRun: false };
}

export const UpdateRuntimeRequest = {
  encode(message: UpdateRuntimeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.newRun === true) {
      writer.uint32(24).bool(message.newRun);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRuntimeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRuntimeRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.newRun = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRuntimeRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      newRun: isSet(object.newRun) ? Boolean(object.newRun) : false,
    };
  },

  toJSON(message: UpdateRuntimeRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.newRun === true) {
      obj.newRun = message.newRun;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRuntimeRequest>, I>>(base?: I): UpdateRuntimeRequest {
    return UpdateRuntimeRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRuntimeRequest>, I>>(object: I): UpdateRuntimeRequest {
    const message = createBaseUpdateRuntimeRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.newRun = object.newRun ?? false;
    return message;
  },
};

function createBaseUpdateRuntimeResponse(): UpdateRuntimeResponse {
  return { runtime: undefined };
}

export const UpdateRuntimeResponse = {
  encode(message: UpdateRuntimeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.runtime !== undefined) {
      Runtime.encode(message.runtime, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRuntimeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRuntimeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.runtime = Runtime.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRuntimeResponse {
    return { runtime: isSet(object.runtime) ? Runtime.fromJSON(object.runtime) : undefined };
  },

  toJSON(message: UpdateRuntimeResponse): unknown {
    const obj: any = {};
    if (message.runtime !== undefined) {
      obj.runtime = Runtime.toJSON(message.runtime);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRuntimeResponse>, I>>(base?: I): UpdateRuntimeResponse {
    return UpdateRuntimeResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRuntimeResponse>, I>>(object: I): UpdateRuntimeResponse {
    const message = createBaseUpdateRuntimeResponse();
    message.runtime = (object.runtime !== undefined && object.runtime !== null)
      ? Runtime.fromPartial(object.runtime)
      : undefined;
    return message;
  },
};

function createBaseDeleteRuntimeRequest(): DeleteRuntimeRequest {
  return { namespace: "", name: "" };
}

export const DeleteRuntimeRequest = {
  encode(message: DeleteRuntimeRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRuntimeRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRuntimeRequest();
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

  fromJSON(object: any): DeleteRuntimeRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: DeleteRuntimeRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRuntimeRequest>, I>>(base?: I): DeleteRuntimeRequest {
    return DeleteRuntimeRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRuntimeRequest>, I>>(object: I): DeleteRuntimeRequest {
    const message = createBaseDeleteRuntimeRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseDeleteRuntimeResponse(): DeleteRuntimeResponse {
  return {};
}

export const DeleteRuntimeResponse = {
  encode(_: DeleteRuntimeResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRuntimeResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRuntimeResponse();
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

  fromJSON(_: any): DeleteRuntimeResponse {
    return {};
  },

  toJSON(_: DeleteRuntimeResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRuntimeResponse>, I>>(base?: I): DeleteRuntimeResponse {
    return DeleteRuntimeResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRuntimeResponse>, I>>(_: I): DeleteRuntimeResponse {
    const message = createBaseDeleteRuntimeResponse();
    return message;
  },
};

function createBaseUploadRuntimeBinaryRequest(): UploadRuntimeBinaryRequest {
  return { namespace: "", name: "", binary: new Uint8Array(0) };
}

export const UploadRuntimeBinaryRequest = {
  encode(message: UploadRuntimeBinaryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.binary.length !== 0) {
      writer.uint32(26).bytes(message.binary);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UploadRuntimeBinaryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUploadRuntimeBinaryRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.binary = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UploadRuntimeBinaryRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      binary: isSet(object.binary) ? bytesFromBase64(object.binary) : new Uint8Array(0),
    };
  },

  toJSON(message: UploadRuntimeBinaryRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.binary.length !== 0) {
      obj.binary = base64FromBytes(message.binary);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UploadRuntimeBinaryRequest>, I>>(base?: I): UploadRuntimeBinaryRequest {
    return UploadRuntimeBinaryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UploadRuntimeBinaryRequest>, I>>(object: I): UploadRuntimeBinaryRequest {
    const message = createBaseUploadRuntimeBinaryRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.binary = object.binary ?? new Uint8Array(0);
    return message;
  },
};

function createBaseUploadRuntimeBinaryResponse(): UploadRuntimeBinaryResponse {
  return {};
}

export const UploadRuntimeBinaryResponse = {
  encode(_: UploadRuntimeBinaryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UploadRuntimeBinaryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUploadRuntimeBinaryResponse();
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

  fromJSON(_: any): UploadRuntimeBinaryResponse {
    return {};
  },

  toJSON(_: UploadRuntimeBinaryResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<UploadRuntimeBinaryResponse>, I>>(base?: I): UploadRuntimeBinaryResponse {
    return UploadRuntimeBinaryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UploadRuntimeBinaryResponse>, I>>(_: I): UploadRuntimeBinaryResponse {
    const message = createBaseUploadRuntimeBinaryResponse();
    return message;
  },
};

function createBaseDownloadRuntimeBinaryRequest(): DownloadRuntimeBinaryRequest {
  return { namespace: "", name: "" };
}

export const DownloadRuntimeBinaryRequest = {
  encode(message: DownloadRuntimeBinaryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DownloadRuntimeBinaryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDownloadRuntimeBinaryRequest();
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

  fromJSON(object: any): DownloadRuntimeBinaryRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: DownloadRuntimeBinaryRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DownloadRuntimeBinaryRequest>, I>>(base?: I): DownloadRuntimeBinaryRequest {
    return DownloadRuntimeBinaryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DownloadRuntimeBinaryRequest>, I>>(object: I): DownloadRuntimeBinaryRequest {
    const message = createBaseDownloadRuntimeBinaryRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseDownloadRuntimeBinaryResponse(): DownloadRuntimeBinaryResponse {
  return { binary: new Uint8Array(0) };
}

export const DownloadRuntimeBinaryResponse = {
  encode(message: DownloadRuntimeBinaryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.binary.length !== 0) {
      writer.uint32(10).bytes(message.binary);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DownloadRuntimeBinaryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDownloadRuntimeBinaryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.binary = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DownloadRuntimeBinaryResponse {
    return { binary: isSet(object.binary) ? bytesFromBase64(object.binary) : new Uint8Array(0) };
  },

  toJSON(message: DownloadRuntimeBinaryResponse): unknown {
    const obj: any = {};
    if (message.binary.length !== 0) {
      obj.binary = base64FromBytes(message.binary);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DownloadRuntimeBinaryResponse>, I>>(base?: I): DownloadRuntimeBinaryResponse {
    return DownloadRuntimeBinaryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DownloadRuntimeBinaryResponse>, I>>(
    object: I,
  ): DownloadRuntimeBinaryResponse {
    const message = createBaseDownloadRuntimeBinaryResponse();
    message.binary = object.binary ?? new Uint8Array(0);
    return message;
  },
};

export interface RuntimeService {
  GetRuntimesForNamespace(request: GetRuntimesForNamespaceReqeust): Promise<GetRuntimesForNamespaceResponse>;
  GetRuntime(request: GetRuntimeRequest): Promise<GetRuntimeResponse>;
  CreateRuntime(request: CreateRuntimeRequest): Promise<CreateRuntimeResponse>;
  UpdateRuntime(request: UpdateRuntimeRequest): Promise<UpdateRuntimeResponse>;
  DeleteRuntime(request: DeleteRuntimeRequest): Promise<DeleteRuntimeResponse>;
  UploadRuntimeBinary(request: Observable<UploadRuntimeBinaryRequest>): Promise<UploadRuntimeBinaryResponse>;
  DownloadRuntimeBinary(request: DownloadRuntimeBinaryRequest): Observable<DownloadRuntimeBinaryResponse>;
}

export const RuntimeServiceServiceName = "runtime_manager_runtime.RuntimeService";
export class RuntimeServiceClientImpl implements RuntimeService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || RuntimeServiceServiceName;
    this.rpc = rpc;
    this.GetRuntimesForNamespace = this.GetRuntimesForNamespace.bind(this);
    this.GetRuntime = this.GetRuntime.bind(this);
    this.CreateRuntime = this.CreateRuntime.bind(this);
    this.UpdateRuntime = this.UpdateRuntime.bind(this);
    this.DeleteRuntime = this.DeleteRuntime.bind(this);
    this.UploadRuntimeBinary = this.UploadRuntimeBinary.bind(this);
    this.DownloadRuntimeBinary = this.DownloadRuntimeBinary.bind(this);
  }
  GetRuntimesForNamespace(request: GetRuntimesForNamespaceReqeust): Promise<GetRuntimesForNamespaceResponse> {
    const data = GetRuntimesForNamespaceReqeust.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetRuntimesForNamespace", data);
    return promise.then((data) => GetRuntimesForNamespaceResponse.decode(_m0.Reader.create(data)));
  }

  GetRuntime(request: GetRuntimeRequest): Promise<GetRuntimeResponse> {
    const data = GetRuntimeRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetRuntime", data);
    return promise.then((data) => GetRuntimeResponse.decode(_m0.Reader.create(data)));
  }

  CreateRuntime(request: CreateRuntimeRequest): Promise<CreateRuntimeResponse> {
    const data = CreateRuntimeRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CreateRuntime", data);
    return promise.then((data) => CreateRuntimeResponse.decode(_m0.Reader.create(data)));
  }

  UpdateRuntime(request: UpdateRuntimeRequest): Promise<UpdateRuntimeResponse> {
    const data = UpdateRuntimeRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "UpdateRuntime", data);
    return promise.then((data) => UpdateRuntimeResponse.decode(_m0.Reader.create(data)));
  }

  DeleteRuntime(request: DeleteRuntimeRequest): Promise<DeleteRuntimeResponse> {
    const data = DeleteRuntimeRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "DeleteRuntime", data);
    return promise.then((data) => DeleteRuntimeResponse.decode(_m0.Reader.create(data)));
  }

  UploadRuntimeBinary(request: Observable<UploadRuntimeBinaryRequest>): Promise<UploadRuntimeBinaryResponse> {
    const data = request.pipe(map((request) => UploadRuntimeBinaryRequest.encode(request).finish()));
    const promise = this.rpc.clientStreamingRequest(this.service, "UploadRuntimeBinary", data);
    return promise.then((data) => UploadRuntimeBinaryResponse.decode(_m0.Reader.create(data)));
  }

  DownloadRuntimeBinary(request: DownloadRuntimeBinaryRequest): Observable<DownloadRuntimeBinaryResponse> {
    const data = DownloadRuntimeBinaryRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "DownloadRuntimeBinary", data);
    return result.pipe(map((data) => DownloadRuntimeBinaryResponse.decode(_m0.Reader.create(data))));
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
