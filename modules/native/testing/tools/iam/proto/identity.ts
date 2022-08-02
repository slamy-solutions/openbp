/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_identity";

export interface Identity {
  /** Namespaces of the identity. Can be empty for global identities. */
  namespace: string;
  /** Unique identity identifier */
  uuid: string;
  /** Public identity name */
  name: string;
  /** If identity is not active, it will not be able to login and perform any actions. */
  active: boolean;
  /** Security policies assigned to the identity */
  policies: Identity_PolicyReference[];
}

export interface Identity_PolicyReference {
  /** Policy namespace. Empty for global policy */
  namespace: string;
  /** Policy uuid (unique identifier) inside namespace */
  uuid: string;
}

export interface CreateIdentityRequest {
  /** Namespace where to create identity */
  namespace: string;
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
  /** Identity namespace */
  namespace: string;
  /** Identity unique identifier inside namespace */
  uuid: string;
  /** Use cache or not. Cache may not be valid under very rare conditions (simultaniour read and writes). Cache automatically clears after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface GetIdentityResponse {
  /** Identity information */
  identity: Identity | undefined;
}

export interface DeleteIdentityRequest {
  /** Identity namespace */
  namespace: string;
  /** Identity unique identifier inside namespace */
  uuid: string;
}

export interface DeleteIdentityResponse {}

export interface AddPolicyRequest {
  /** Identity namespace */
  identityNamespace: string;
  /** Identity identifier inside identity namespace */
  identityUUID: string;
  /** Policy namespace */
  policyNamespace: string;
  /** Policy UUID inside policy namespace */
  policyUUID: string;
}

export interface AddPolicyResponse {
  /** Updated identity (after adding policy) */
  identity: Identity | undefined;
}

export interface RemovePolicyRequest {
  /** Identity namespace */
  identityNamespace: string;
  /** Identity unique identifier inside identity namespace */
  identityUUID: string;
  /** Policy namespace */
  policyNamespace: string;
  /** Policy UUID inside policy namespace */
  policyUUID: string;
}

export interface RemovePolicyResponse {
  /** Updated identity (after removing policy) */
  identity: Identity | undefined;
}

export interface SetIdentityActiveRequest {
  /** Namespace of the identity */
  namespace: string;
  /** Identity unique identifier inside namespace */
  uuid: string;
  /** Set active or not */
  active: boolean;
}

export interface SetIdentityActiveResponse {
  /** Identity after update */
  identity: Identity | undefined;
}

function createBaseIdentity(): Identity {
  return { namespace: "", uuid: "", name: "", active: false, policies: [] };
}

export const Identity = {
  encode(
    message: Identity,
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
    if (message.active === true) {
      writer.uint32(32).bool(message.active);
    }
    for (const v of message.policies) {
      Identity_PolicyReference.encode(v!, writer.uint32(42).fork()).ldelim();
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
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
          break;
        case 3:
          message.name = reader.string();
          break;
        case 4:
          message.active = reader.bool();
          break;
        case 5:
          message.policies.push(
            Identity_PolicyReference.decode(reader, reader.uint32())
          );
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
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
      policies: Array.isArray(object?.policies)
        ? object.policies.map((e: any) => Identity_PolicyReference.fromJSON(e))
        : [],
    };
  },

  toJSON(message: Identity): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.name !== undefined && (obj.name = message.name);
    message.active !== undefined && (obj.active = message.active);
    if (message.policies) {
      obj.policies = message.policies.map((e) =>
        e ? Identity_PolicyReference.toJSON(e) : undefined
      );
    } else {
      obj.policies = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Identity>, I>>(object: I): Identity {
    const message = createBaseIdentity();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.active = object.active ?? false;
    message.policies =
      object.policies?.map((e) => Identity_PolicyReference.fromPartial(e)) ||
      [];
    return message;
  },
};

function createBaseIdentity_PolicyReference(): Identity_PolicyReference {
  return { namespace: "", uuid: "" };
}

export const Identity_PolicyReference = {
  encode(
    message: Identity_PolicyReference,
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

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): Identity_PolicyReference {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentity_PolicyReference();
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

  fromJSON(object: any): Identity_PolicyReference {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: Identity_PolicyReference): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Identity_PolicyReference>, I>>(
    object: I
  ): Identity_PolicyReference {
    const message = createBaseIdentity_PolicyReference();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseCreateIdentityRequest(): CreateIdentityRequest {
  return { namespace: "", name: "", initiallyActive: false };
}

export const CreateIdentityRequest = {
  encode(
    message: CreateIdentityRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.initiallyActive === true) {
      writer.uint32(24).bool(message.initiallyActive);
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
          message.namespace = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
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
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      initiallyActive: isSet(object.initiallyActive)
        ? Boolean(object.initiallyActive)
        : false,
    };
  },

  toJSON(message: CreateIdentityRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.name !== undefined && (obj.name = message.name);
    message.initiallyActive !== undefined &&
      (obj.initiallyActive = message.initiallyActive);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateIdentityRequest>, I>>(
    object: I
  ): CreateIdentityRequest {
    const message = createBaseCreateIdentityRequest();
    message.namespace = object.namespace ?? "";
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
  return { namespace: "", uuid: "", useCache: false };
}

export const GetIdentityRequest = {
  encode(
    message: GetIdentityRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIdentityRequest();
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

  fromJSON(object: any): GetIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetIdentityRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetIdentityRequest>, I>>(
    object: I
  ): GetIdentityRequest {
    const message = createBaseGetIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
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

function createBaseDeleteIdentityRequest(): DeleteIdentityRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteIdentityRequest = {
  encode(
    message: DeleteIdentityRequest,
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

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeleteIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteIdentityRequest();
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

  fromJSON(object: any): DeleteIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteIdentityRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteIdentityRequest>, I>>(
    object: I
  ): DeleteIdentityRequest {
    const message = createBaseDeleteIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteIdentityResponse(): DeleteIdentityResponse {
  return {};
}

export const DeleteIdentityResponse = {
  encode(
    _: DeleteIdentityResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeleteIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteIdentityResponse();
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

  fromJSON(_: any): DeleteIdentityResponse {
    return {};
  },

  toJSON(_: DeleteIdentityResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteIdentityResponse>, I>>(
    _: I
  ): DeleteIdentityResponse {
    const message = createBaseDeleteIdentityResponse();
    return message;
  },
};

function createBaseAddPolicyRequest(): AddPolicyRequest {
  return {
    identityNamespace: "",
    identityUUID: "",
    policyNamespace: "",
    policyUUID: "",
  };
}

export const AddPolicyRequest = {
  encode(
    message: AddPolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identityNamespace !== "") {
      writer.uint32(10).string(message.identityNamespace);
    }
    if (message.identityUUID !== "") {
      writer.uint32(18).string(message.identityUUID);
    }
    if (message.policyNamespace !== "") {
      writer.uint32(26).string(message.policyNamespace);
    }
    if (message.policyUUID !== "") {
      writer.uint32(34).string(message.policyUUID);
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
          message.identityNamespace = reader.string();
          break;
        case 2:
          message.identityUUID = reader.string();
          break;
        case 3:
          message.policyNamespace = reader.string();
          break;
        case 4:
          message.policyUUID = reader.string();
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
      identityNamespace: isSet(object.identityNamespace)
        ? String(object.identityNamespace)
        : "",
      identityUUID: isSet(object.identityUUID)
        ? String(object.identityUUID)
        : "",
      policyNamespace: isSet(object.policyNamespace)
        ? String(object.policyNamespace)
        : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: AddPolicyRequest): unknown {
    const obj: any = {};
    message.identityNamespace !== undefined &&
      (obj.identityNamespace = message.identityNamespace);
    message.identityUUID !== undefined &&
      (obj.identityUUID = message.identityUUID);
    message.policyNamespace !== undefined &&
      (obj.policyNamespace = message.policyNamespace);
    message.policyUUID !== undefined && (obj.policyUUID = message.policyUUID);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AddPolicyRequest>, I>>(
    object: I
  ): AddPolicyRequest {
    const message = createBaseAddPolicyRequest();
    message.identityNamespace = object.identityNamespace ?? "";
    message.identityUUID = object.identityUUID ?? "";
    message.policyNamespace = object.policyNamespace ?? "";
    message.policyUUID = object.policyUUID ?? "";
    return message;
  },
};

function createBaseAddPolicyResponse(): AddPolicyResponse {
  return { identity: undefined };
}

export const AddPolicyResponse = {
  encode(
    message: AddPolicyResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddPolicyResponse();
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

  fromJSON(object: any): AddPolicyResponse {
    return {
      identity: isSet(object.identity)
        ? Identity.fromJSON(object.identity)
        : undefined,
    };
  },

  toJSON(message: AddPolicyResponse): unknown {
    const obj: any = {};
    message.identity !== undefined &&
      (obj.identity = message.identity
        ? Identity.toJSON(message.identity)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(
    object: I
  ): AddPolicyResponse {
    const message = createBaseAddPolicyResponse();
    message.identity =
      object.identity !== undefined && object.identity !== null
        ? Identity.fromPartial(object.identity)
        : undefined;
    return message;
  },
};

function createBaseRemovePolicyRequest(): RemovePolicyRequest {
  return {
    identityNamespace: "",
    identityUUID: "",
    policyNamespace: "",
    policyUUID: "",
  };
}

export const RemovePolicyRequest = {
  encode(
    message: RemovePolicyRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identityNamespace !== "") {
      writer.uint32(10).string(message.identityNamespace);
    }
    if (message.identityUUID !== "") {
      writer.uint32(18).string(message.identityUUID);
    }
    if (message.policyNamespace !== "") {
      writer.uint32(26).string(message.policyNamespace);
    }
    if (message.policyUUID !== "") {
      writer.uint32(34).string(message.policyUUID);
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
          message.identityNamespace = reader.string();
          break;
        case 2:
          message.identityUUID = reader.string();
          break;
        case 3:
          message.policyNamespace = reader.string();
          break;
        case 4:
          message.policyUUID = reader.string();
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
      identityNamespace: isSet(object.identityNamespace)
        ? String(object.identityNamespace)
        : "",
      identityUUID: isSet(object.identityUUID)
        ? String(object.identityUUID)
        : "",
      policyNamespace: isSet(object.policyNamespace)
        ? String(object.policyNamespace)
        : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: RemovePolicyRequest): unknown {
    const obj: any = {};
    message.identityNamespace !== undefined &&
      (obj.identityNamespace = message.identityNamespace);
    message.identityUUID !== undefined &&
      (obj.identityUUID = message.identityUUID);
    message.policyNamespace !== undefined &&
      (obj.policyNamespace = message.policyNamespace);
    message.policyUUID !== undefined && (obj.policyUUID = message.policyUUID);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RemovePolicyRequest>, I>>(
    object: I
  ): RemovePolicyRequest {
    const message = createBaseRemovePolicyRequest();
    message.identityNamespace = object.identityNamespace ?? "";
    message.identityUUID = object.identityUUID ?? "";
    message.policyNamespace = object.policyNamespace ?? "";
    message.policyUUID = object.policyUUID ?? "";
    return message;
  },
};

function createBaseRemovePolicyResponse(): RemovePolicyResponse {
  return { identity: undefined };
}

export const RemovePolicyResponse = {
  encode(
    message: RemovePolicyResponse,
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
  ): RemovePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemovePolicyResponse();
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

  fromJSON(object: any): RemovePolicyResponse {
    return {
      identity: isSet(object.identity)
        ? Identity.fromJSON(object.identity)
        : undefined,
    };
  },

  toJSON(message: RemovePolicyResponse): unknown {
    const obj: any = {};
    message.identity !== undefined &&
      (obj.identity = message.identity
        ? Identity.toJSON(message.identity)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(
    object: I
  ): RemovePolicyResponse {
    const message = createBaseRemovePolicyResponse();
    message.identity =
      object.identity !== undefined && object.identity !== null
        ? Identity.fromPartial(object.identity)
        : undefined;
    return message;
  },
};

function createBaseSetIdentityActiveRequest(): SetIdentityActiveRequest {
  return { namespace: "", uuid: "", active: false };
}

export const SetIdentityActiveRequest = {
  encode(
    message: SetIdentityActiveRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.active === true) {
      writer.uint32(24).bool(message.active);
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
          message.namespace = reader.string();
          break;
        case 2:
          message.uuid = reader.string();
          break;
        case 3:
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
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
    };
  },

  toJSON(message: SetIdentityActiveRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.active !== undefined && (obj.active = message.active);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveRequest>, I>>(
    object: I
  ): SetIdentityActiveRequest {
    const message = createBaseSetIdentityActiveRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.active = object.active ?? false;
    return message;
  },
};

function createBaseSetIdentityActiveResponse(): SetIdentityActiveResponse {
  return { identity: undefined };
}

export const SetIdentityActiveResponse = {
  encode(
    message: SetIdentityActiveResponse,
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
  ): SetIdentityActiveResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIdentityActiveResponse();
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

  fromJSON(object: any): SetIdentityActiveResponse {
    return {
      identity: isSet(object.identity)
        ? Identity.fromJSON(object.identity)
        : undefined,
    };
  },

  toJSON(message: SetIdentityActiveResponse): unknown {
    const obj: any = {};
    message.identity !== undefined &&
      (obj.identity = message.identity
        ? Identity.toJSON(message.identity)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveResponse>, I>>(
    object: I
  ): SetIdentityActiveResponse {
    const message = createBaseSetIdentityActiveResponse();
    message.identity =
      object.identity !== undefined && object.identity !== null
        ? Identity.fromPartial(object.identity)
        : undefined;
    return message;
  },
};

/** Provides API to manage IAM identities */
export interface IAMIdentityService {
  /** Create new identity */
  Create(request: CreateIdentityRequest): Promise<CreateIdentityResponse>;
  /** Get identity */
  Get(request: GetIdentityRequest): Promise<GetIdentityResponse>;
  /** Delete identity */
  Delete(request: DeleteIdentityRequest): Promise<DeleteIdentityResponse>;
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
    this.Delete = this.Delete.bind(this);
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

  Delete(request: DeleteIdentityRequest): Promise<DeleteIdentityResponse> {
    const data = DeleteIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_identity.IAMIdentityService",
      "Delete",
      data
    );
    return promise.then((data) =>
      DeleteIdentityResponse.decode(new _m0.Reader(data))
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
