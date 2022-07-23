/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

export const protobufPackage = "native_iam_policy";

export interface Policy {
  /** Namespace where policy was created. Namespace can be empty for global policy. */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Public name */
  name: string;
  /** List of actions that can be performed */
  actions: string[];
}

export interface CreatePolicyRequest {
  /** Namespace where policy will be created. Namespace can be empty for global policy. */
  namespace: string;
  /** Public name. May not be unique. */
  name: string;
  /** List of actions that can be performed with this policy */
  actions: string[];
}

export interface CreatePolicyResponse {
  policy: Policy | undefined;
}

export interface GetPolicyRequest {
  /** Namespace of the policy */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Use cache or not. Cache may be invalid under very rare conditions (simultanious read and writes to the policy while it is not in cache). Cache automatically deletes after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface GetPolicyResponse {
  policy: Policy | undefined;
}

export interface UpdatePolicyRequest {
  /** Namespace of the policy */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Public name */
  name: string;
  /** List of actions that can be performed */
  actions: string[];
}

export interface UpdatePolicyResponse {
  /** Updated policy */
  policy: Policy | undefined;
}

export interface DeletePolicyRequest {
  /** Namespace of the policy */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
}

export interface DeletePolicyResponse {}

export interface ListPoliciesRequest {
  /** Namespace from where to list policies */
  namespace: string;
  /** How many values to skip before returning result */
  skip: number;
  /** Maximum number of values to return. Use 0 to return all up to the end. */
  limit: number;
}

export interface ListPoliciesResponse {
  policy: Policy | undefined;
}

function createBasePolicy(): Policy {
  return { namespace: "", uuid: "", name: "", actions: [] };
}

export const Policy = {
  encode(
    message: Policy,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    for (const v of message.actions) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Policy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
          break;
        case 3:
          message.name = reader.string();
          break;
        case 4:
          message.actions.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Policy {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      actions: Array.isArray(object?.actions)
        ? object.actions.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: Policy): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.name !== undefined && (obj.name = message.name);
    if (message.actions) {
      obj.actions = message.actions.map((e) => e);
    } else {
      obj.actions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Policy>, I>>(object: I): Policy {
    const message = createBasePolicy();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseCreatePolicyRequest(): CreatePolicyRequest {
  return { namespace: "", name: "", actions: [] };
}

export const CreatePolicyRequest = {
  encode(
    message: CreatePolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    for (const v of message.actions) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreatePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreatePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 3:
          message.name = reader.string();
          break;
        case 4:
          message.actions.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreatePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      actions: Array.isArray(object?.actions)
        ? object.actions.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: CreatePolicyRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.name !== undefined && (obj.name = message.name);
    if (message.actions) {
      obj.actions = message.actions.map((e) => e);
    } else {
      obj.actions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreatePolicyRequest>, I>>(
    object: I
  ): CreatePolicyRequest {
    const message = createBaseCreatePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseCreatePolicyResponse(): CreatePolicyResponse {
  return { policy: undefined };
}

export const CreatePolicyResponse = {
  encode(
    message: CreatePolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreatePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreatePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = Policy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreatePolicyResponse {
    return {
      policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined,
    };
  },

  toJSON(message: CreatePolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? Policy.toJSON(message.policy) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreatePolicyResponse>, I>>(
    object: I
  ): CreatePolicyResponse {
    const message = createBaseCreatePolicyResponse();
    message.policy =
      object.policy !== undefined && object.policy !== null
        ? Policy.fromPartial(object.policy)
        : undefined;
    return message;
  },
};

function createBaseGetPolicyRequest(): GetPolicyRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetPolicyRequest = {
  encode(
    message: GetPolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
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

  fromJSON(object: any): GetPolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetPolicyRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetPolicyRequest>, I>>(
    object: I
  ): GetPolicyRequest {
    const message = createBaseGetPolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetPolicyResponse(): GetPolicyResponse {
  return { policy: undefined };
}

export const GetPolicyResponse = {
  encode(
    message: GetPolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = Policy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetPolicyResponse {
    return {
      policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined,
    };
  },

  toJSON(message: GetPolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? Policy.toJSON(message.policy) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetPolicyResponse>, I>>(
    object: I
  ): GetPolicyResponse {
    const message = createBaseGetPolicyResponse();
    message.policy =
      object.policy !== undefined && object.policy !== null
        ? Policy.fromPartial(object.policy)
        : undefined;
    return message;
  },
};

function createBaseUpdatePolicyRequest(): UpdatePolicyRequest {
  return { namespace: "", uuid: "", name: "", actions: [] };
}

export const UpdatePolicyRequest = {
  encode(
    message: UpdatePolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    for (const v of message.actions) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
          break;
        case 3:
          message.name = reader.string();
          break;
        case 4:
          message.actions.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      actions: Array.isArray(object?.actions)
        ? object.actions.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: UpdatePolicyRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.name !== undefined && (obj.name = message.name);
    if (message.actions) {
      obj.actions = message.actions.map((e) => e);
    } else {
      obj.actions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdatePolicyRequest>, I>>(
    object: I
  ): UpdatePolicyRequest {
    const message = createBaseUpdatePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseUpdatePolicyResponse(): UpdatePolicyResponse {
  return { policy: undefined };
}

export const UpdatePolicyResponse = {
  encode(
    message: UpdatePolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): UpdatePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = Policy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdatePolicyResponse {
    return {
      policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined,
    };
  },

  toJSON(message: UpdatePolicyResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? Policy.toJSON(message.policy) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdatePolicyResponse>, I>>(
    object: I
  ): UpdatePolicyResponse {
    const message = createBaseUpdatePolicyResponse();
    message.policy =
      object.policy !== undefined && object.policy !== null
        ? Policy.fromPartial(object.policy)
        : undefined;
    return message;
  },
};

function createBaseDeletePolicyRequest(): DeletePolicyRequest {
  return { namespace: "", uuid: "" };
}

export const DeletePolicyRequest = {
  encode(
    message: DeletePolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeletePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): DeletePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeletePolicyRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeletePolicyRequest>, I>>(
    object: I
  ): DeletePolicyRequest {
    const message = createBaseDeletePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeletePolicyResponse(): DeletePolicyResponse {
  return {};
}

export const DeletePolicyResponse = {
  encode(
    _: DeletePolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeletePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePolicyResponse();
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

  fromJSON(_: any): DeletePolicyResponse {
    return {};
  },

  toJSON(_: DeletePolicyResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeletePolicyResponse>, I>>(
    _: I
  ): DeletePolicyResponse {
    const message = createBaseDeletePolicyResponse();
    return message;
  },
};

function createBaseListPoliciesRequest(): ListPoliciesRequest {
  return { namespace: "", skip: 0, limit: 0 };
}

export const ListPoliciesRequest = {
  encode(
    message: ListPoliciesRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.skip !== 0) {
      writer.uint32(16).uint32(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(24).uint32(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListPoliciesRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListPoliciesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.skip = reader.uint32();
          break;
        case 3:
          message.limit = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListPoliciesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListPoliciesRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.skip !== undefined && (obj.skip = Math.round(message.skip));
    message.limit !== undefined && (obj.limit = Math.round(message.limit));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ListPoliciesRequest>, I>>(
    object: I
  ): ListPoliciesRequest {
    const message = createBaseListPoliciesRequest();
    message.namespace = object.namespace ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListPoliciesResponse(): ListPoliciesResponse {
  return { policy: undefined };
}

export const ListPoliciesResponse = {
  encode(
    message: ListPoliciesResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ListPoliciesResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListPoliciesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.policy = Policy.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ListPoliciesResponse {
    return {
      policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined,
    };
  },

  toJSON(message: ListPoliciesResponse): unknown {
    const obj: any = {};
    message.policy !== undefined &&
      (obj.policy = message.policy ? Policy.toJSON(message.policy) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ListPoliciesResponse>, I>>(
    object: I
  ): ListPoliciesResponse {
    const message = createBaseListPoliciesResponse();
    message.policy =
      object.policy !== undefined && object.policy !== null
        ? Policy.fromPartial(object.policy)
        : undefined;
    return message;
  },
};

/** Provides API to manage policies */
export interface IAMPolicyService {
  /** Create new policy */
  Create(request: CreatePolicyRequest): Promise<CreatePolicyResponse>;
  /** Get existing policy by uuid */
  Get(request: GetPolicyRequest): Promise<GetPolicyResponse>;
  /** Update policy */
  Update(request: UpdatePolicyRequest): Promise<UpdatePolicyResponse>;
  /** Delete policy */
  Delete(request: DeletePolicyRequest): Promise<DeletePolicyResponse>;
  /** List policies in namespace */
  List(request: ListPoliciesRequest): Observable<ListPoliciesResponse>;
}

export class IAMPolicyServiceClientImpl implements IAMPolicyService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.List = this.List.bind(this);
  }
  Create(request: CreatePolicyRequest): Promise<CreatePolicyResponse> {
    const data = CreatePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_policy.IAMPolicyService",
      "Create",
      data
    );
    return promise.then((data) =>
      CreatePolicyResponse.decode(new _m0.Reader(data))
    );
  }

  Get(request: GetPolicyRequest): Promise<GetPolicyResponse> {
    const data = GetPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_policy.IAMPolicyService",
      "Get",
      data
    );
    return promise.then((data) =>
      GetPolicyResponse.decode(new _m0.Reader(data))
    );
  }

  Update(request: UpdatePolicyRequest): Promise<UpdatePolicyResponse> {
    const data = UpdatePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_policy.IAMPolicyService",
      "Update",
      data
    );
    return promise.then((data) =>
      UpdatePolicyResponse.decode(new _m0.Reader(data))
    );
  }

  Delete(request: DeletePolicyRequest): Promise<DeletePolicyResponse> {
    const data = DeletePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_policy.IAMPolicyService",
      "Delete",
      data
    );
    return promise.then((data) =>
      DeletePolicyResponse.decode(new _m0.Reader(data))
    );
  }

  List(request: ListPoliciesRequest): Observable<ListPoliciesResponse> {
    const data = ListPoliciesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(
      "native_iam_policy.IAMPolicyService",
      "List",
      data
    );
    return result.pipe(
      map((data) => ListPoliciesResponse.decode(new _m0.Reader(data)))
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
