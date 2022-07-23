/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_identity";

export interface Identity {
  /** Unique identity identifier */
  uuid: string;
  /** Public identity name */
  name: string;
  /** If identity is not active, it will not be able to login and perform any actions. */
  active: boolean;
  /** Security policies assigned to the identity */
  policies: string[];
}

export interface CreateIdentityRequest {
  /** Public name for newly created identity. It may not be unique - this is just human-readable name. */
  name: string;
  /** Should the identity be active on the start or not */
  initiallyActive: boolean;
}

export interface CreateIdentityResponse {
  /** Created identity */
  identity: Identity | undefined;
}

export interface GetIdentityRequest {
  /** Identity unique identifier */
  uuid: string;
}

export interface GetIdentityResponse {
  /** Identity information */
  identity: Identity | undefined;
}

export interface AddPolicyRequest {
  /** Identity UUID to add policy */
  identity: string;
  /** Policy UUID to add */
  policy: string;
}

export interface AddPolicyResponse {}

export interface RemovePolicyRequest {
  /** Identity UUID to remove policy */
  identity: string;
  /** Policy UUID to remove */
  policy: string;
}

export interface RemovePolicyResponse {}

export interface SetIdentityActiveRequest {
  /** Identity UUID */
  identity: string;
  /** Set active or not */
  active: boolean;
}

export interface SetIdentityActiveResponse {}

function createBaseIdentity(): Identity {
  return { uuid: "", name: "", active: false, policies: [] };
}

export const Identity = {
  encode(
    message: Identity,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.active === true) {
      writer.uint32(24).bool(message.active);
    }
    for (const v of message.policies) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Identity {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentity();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.active = reader.bool();
          break;
        case 4:
          message.policies.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Identity {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
      policies: Array.isArray(object?.policies)
        ? object.policies.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: Identity): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.name !== undefined && (obj.name = message.name);
    message.active !== undefined && (obj.active = message.active);
    if (message.policies) {
      obj.policies = message.policies.map((e) => e);
    } else {
      obj.policies = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Identity>, I>>(object: I): Identity {
    const message = createBaseIdentity();
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.active = object.active ?? false;
    message.policies = object.policies?.map((e) => e) || [];
    return message;
  },
};

function createBaseCreateIdentityRequest(): CreateIdentityRequest {
  return { name: "", initiallyActive: false };
}

export const CreateIdentityRequest = {
  encode(
    message: CreateIdentityRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.initiallyActive === true) {
      writer.uint32(16).bool(message.initiallyActive);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateIdentityRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.name = reader.string();
          break;
        case 2:
          message.initiallyActive = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateIdentityRequest {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      initiallyActive: isSet(object.initiallyActive)
        ? Boolean(object.initiallyActive)
        : false,
    };
  },

  toJSON(message: CreateIdentityRequest): unknown {
    const obj: any = {};
    message.name !== undefined && (obj.name = message.name);
    message.initiallyActive !== undefined &&
      (obj.initiallyActive = message.initiallyActive);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateIdentityRequest>, I>>(
    object: I
  ): CreateIdentityRequest {
    const message = createBaseCreateIdentityRequest();
    message.name = object.name ?? "";
    message.initiallyActive = object.initiallyActive ?? false;
    return message;
  },
};

function createBaseCreateIdentityResponse(): CreateIdentityResponse {
  return { identity: undefined };
}

export const CreateIdentityResponse = {
  encode(
    message: CreateIdentityResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = Identity.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateIdentityResponse {
    return {
      identity: isSet(object.identity)
        ? Identity.fromJSON(object.identity)
        : undefined,
    };
  },

  toJSON(message: CreateIdentityResponse): unknown {
    const obj: any = {};
    message.identity !== undefined &&
      (obj.identity = message.identity
        ? Identity.toJSON(message.identity)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateIdentityResponse>, I>>(
    object: I
  ): CreateIdentityResponse {
    const message = createBaseCreateIdentityResponse();
    message.identity =
      object.identity !== undefined && object.identity !== null
        ? Identity.fromPartial(object.identity)
        : undefined;
    return message;
  },
};

function createBaseGetIdentityRequest(): GetIdentityRequest {
  return { uuid: "" };
}

export const GetIdentityRequest = {
  encode(
    message: GetIdentityRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIdentityRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetIdentityRequest {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: GetIdentityRequest): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetIdentityRequest>, I>>(
    object: I
  ): GetIdentityRequest {
    const message = createBaseGetIdentityRequest();
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseGetIdentityResponse(): GetIdentityResponse {
  return { identity: undefined };
}

export const GetIdentityResponse = {
  encode(
    message: GetIdentityResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = Identity.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetIdentityResponse {
    return {
      identity: isSet(object.identity)
        ? Identity.fromJSON(object.identity)
        : undefined,
    };
  },

  toJSON(message: GetIdentityResponse): unknown {
    const obj: any = {};
    message.identity !== undefined &&
      (obj.identity = message.identity
        ? Identity.toJSON(message.identity)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetIdentityResponse>, I>>(
    object: I
  ): GetIdentityResponse {
    const message = createBaseGetIdentityResponse();
    message.identity =
      object.identity !== undefined && object.identity !== null
        ? Identity.fromPartial(object.identity)
        : undefined;
    return message;
  },
};

function createBaseAddPolicyRequest(): AddPolicyRequest {
  return { identity: "", policy: "" };
}

export const AddPolicyRequest = {
  encode(
    message: AddPolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== "") {
      writer.uint32(10).string(message.identity);
    }
    if (message.policy !== "") {
      writer.uint32(18).string(message.policy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = reader.string();
          break;
        case 2:
          message.policy = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddPolicyRequest {
    return {
      identity: isSet(object.identity) ? String(object.identity) : "",
      policy: isSet(object.policy) ? String(object.policy) : "",
    };
  },

  toJSON(message: AddPolicyRequest): unknown {
    const obj: any = {};
    message.identity !== undefined && (obj.identity = message.identity);
    message.policy !== undefined && (obj.policy = message.policy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AddPolicyRequest>, I>>(
    object: I
  ): AddPolicyRequest {
    const message = createBaseAddPolicyRequest();
    message.identity = object.identity ?? "";
    message.policy = object.policy ?? "";
    return message;
  },
};

function createBaseAddPolicyResponse(): AddPolicyResponse {
  return {};
}

export const AddPolicyResponse = {
  encode(
    _: AddPolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddPolicyResponse();
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

  fromJSON(_: any): AddPolicyResponse {
    return {};
  },

  toJSON(_: AddPolicyResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(
    _: I
  ): AddPolicyResponse {
    const message = createBaseAddPolicyResponse();
    return message;
  },
};

function createBaseRemovePolicyRequest(): RemovePolicyRequest {
  return { identity: "", policy: "" };
}

export const RemovePolicyRequest = {
  encode(
    message: RemovePolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== "") {
      writer.uint32(10).string(message.identity);
    }
    if (message.policy !== "") {
      writer.uint32(18).string(message.policy);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemovePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemovePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = reader.string();
          break;
        case 2:
          message.policy = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RemovePolicyRequest {
    return {
      identity: isSet(object.identity) ? String(object.identity) : "",
      policy: isSet(object.policy) ? String(object.policy) : "",
    };
  },

  toJSON(message: RemovePolicyRequest): unknown {
    const obj: any = {};
    message.identity !== undefined && (obj.identity = message.identity);
    message.policy !== undefined && (obj.policy = message.policy);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RemovePolicyRequest>, I>>(
    object: I
  ): RemovePolicyRequest {
    const message = createBaseRemovePolicyRequest();
    message.identity = object.identity ?? "";
    message.policy = object.policy ?? "";
    return message;
  },
};

function createBaseRemovePolicyResponse(): RemovePolicyResponse {
  return {};
}

export const RemovePolicyResponse = {
  encode(
    _: RemovePolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): RemovePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemovePolicyResponse();
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

  fromJSON(_: any): RemovePolicyResponse {
    return {};
  },

  toJSON(_: RemovePolicyResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(
    _: I
  ): RemovePolicyResponse {
    const message = createBaseRemovePolicyResponse();
    return message;
  },
};

function createBaseSetIdentityActiveRequest(): SetIdentityActiveRequest {
  return { identity: "", active: false };
}

export const SetIdentityActiveRequest = {
  encode(
    message: SetIdentityActiveRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== "") {
      writer.uint32(10).string(message.identity);
    }
    if (message.active === true) {
      writer.uint32(16).bool(message.active);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SetIdentityActiveRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIdentityActiveRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = reader.string();
          break;
        case 2:
          message.active = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetIdentityActiveRequest {
    return {
      identity: isSet(object.identity) ? String(object.identity) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
    };
  },

  toJSON(message: SetIdentityActiveRequest): unknown {
    const obj: any = {};
    message.identity !== undefined && (obj.identity = message.identity);
    message.active !== undefined && (obj.active = message.active);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveRequest>, I>>(
    object: I
  ): SetIdentityActiveRequest {
    const message = createBaseSetIdentityActiveRequest();
    message.identity = object.identity ?? "";
    message.active = object.active ?? false;
    return message;
  },
};

function createBaseSetIdentityActiveResponse(): SetIdentityActiveResponse {
  return {};
}

export const SetIdentityActiveResponse = {
  encode(
    _: SetIdentityActiveResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): SetIdentityActiveResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIdentityActiveResponse();
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

  fromJSON(_: any): SetIdentityActiveResponse {
    return {};
  },

  toJSON(_: SetIdentityActiveResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveResponse>, I>>(
    _: I
  ): SetIdentityActiveResponse {
    const message = createBaseSetIdentityActiveResponse();
    return message;
  },
};

/** Provides API to manage IAM identities */
export interface IAMIdentityService {
  /** Create new identity */
  Create(request: CreateIdentityRequest): Promise<CreateIdentityResponse>;
  /** Get identity */
  Get(request: GetIdentityRequest): Promise<GetIdentityResponse>;
  /** Add policy to the identity. If policy was already added - does nothing. */
  AddPolicy(request: AddPolicyRequest): Promise<AddPolicyResponse>;
  /** Removes policy from the identity. If policy was already removed - does nothing. */
  RemovePolicy(request: RemovePolicyRequest): Promise<RemovePolicyResponse>;
  /** Set identity active or not. */
  SetActive(
    request: SetIdentityActiveRequest
  ): Promise<SetIdentityActiveResponse>;
}

export class IAMIdentityServiceClientImpl implements IAMIdentityService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.AddPolicy = this.AddPolicy.bind(this);
    this.RemovePolicy = this.RemovePolicy.bind(this);
    this.SetActive = this.SetActive.bind(this);
  }
  Create(request: CreateIdentityRequest): Promise<CreateIdentityResponse> {
    const data = CreateIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "Create",
      data
    );
    return promise.then((data) =>
      CreateIdentityResponse.decode(new _m0.Reader(data))
    );
  }

  Get(request: GetIdentityRequest): Promise<GetIdentityResponse> {
    const data = GetIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "Get",
      data
    );
    return promise.then((data) =>
      GetIdentityResponse.decode(new _m0.Reader(data))
    );
  }

  AddPolicy(request: AddPolicyRequest): Promise<AddPolicyResponse> {
    const data = AddPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "AddPolicy",
      data
    );
    return promise.then((data) =>
      AddPolicyResponse.decode(new _m0.Reader(data))
    );
  }

  RemovePolicy(request: RemovePolicyRequest): Promise<RemovePolicyResponse> {
    const data = RemovePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "RemovePolicy",
      data
    );
    return promise.then((data) =>
      RemovePolicyResponse.decode(new _m0.Reader(data))
    );
  }

  SetActive(
    request: SetIdentityActiveRequest
  ): Promise<SetIdentityActiveResponse> {
    const data = SetIdentityActiveRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "SetActive",
      data
    );
    return promise.then((data) =>
      SetIdentityActiveResponse.decode(new _m0.Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
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
