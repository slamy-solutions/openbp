/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_iam_identity";

/** Empty information to indicate that management for this identity is not defined. */
export interface NotManagedData {
}

/** Information about other identity that manages this identity */
export interface IdentityManagedData {
  /** Namespace where identity is located */
  identityNamespace: string;
  /** Identity UUID inside this namespace */
  identityUUID: string;
}

/** Handles information about the service that manages this identity */
export interface ServiceManagedData {
  /** Name of the service */
  service: string;
  /** Reason why this service created this identity */
  reason: string;
  /** This is an ID that can be defined by managed service to find this identity. Set to empty string if you dont need this. If this is not empty Service+ID combination is unique. */
  managementId: string;
}

export interface Identity {
  /** Namespaces of the identity. Can be empty for global identities. */
  namespace: string;
  /** Unique identity identifier */
  uuid: string;
  /** Public identity name */
  name: string;
  /** If identity is not active, it will not be able to login and perform any actions. */
  active: boolean;
  /** Identity is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Identity is managed by other identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Identity is managed by service */
  service?:
    | ServiceManagedData
    | undefined;
  /** Security policies assigned to the identity */
  policies: Identity_PolicyReference[];
  /** Security roles assigned to the identity */
  roles: Identity_RoleReference[];
  /** When the policy was created */
  created:
    | Date
    | undefined;
  /** Last time when the policy information was updated. */
  updated:
    | Date
    | undefined;
  /** Counter that increases after every update of the policy */
  version: number;
}

/** Holds information about specific policy assigned to the identity */
export interface Identity_PolicyReference {
  /** Policy namespace. Empty for global policy */
  namespace: string;
  /** Policy uuid (unique identifier) inside namespace */
  uuid: string;
}

/** Holds information about specific role assigned to the indentity */
export interface Identity_RoleReference {
  /** Role namespace. Empty for global role */
  namespace: string;
  /** Role uuid */
  uuid: string;
}

export interface CreateIdentityRequest {
  /** Namespace where to create identity */
  namespace: string;
  /** Public name for newly created identity. It may not be unique - this is just human-readable name. */
  name: string;
  /** Should the identity be active on the start or not */
  initiallyActive: boolean;
  /** Identity is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Identity is managed by other identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Identity is managed by service */
  service?: ServiceManagedData | undefined;
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

export interface DeleteIdentityResponse {
  /** Indicates if identity existed before request or it was already deleted earlier */
  existed: boolean;
}

export interface ExistsIdentityRequest {
  /** Identity namespace */
  namespace: string;
  /** Identity unique identifier inside namespace */
  uuid: string;
  /** Use cache or not. Cache may not be valid under very rare conditions (simultaniour read and writes). Cache automatically clears after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface ExistsIdentityResponse {
  /** True if identity exists, false if not */
  exists: boolean;
}

export interface ListIdentityRequest {
  /** Namespace where to get identitites */
  namespace: string;
  /** The returned identities will always be sorted. So you can load them with pagination. Service will not return first "skip" identities. Use "0" if you want to ignore pagination, */
  skip: number;
  /** Limit the number of returned entries. Use "0" to ignore limit. */
  limit: number;
}

export interface ListIdentityResponse {
  /** One of the founded identities */
  identity: Identity | undefined;
}

export interface CountIdentityRequest {
  /** Namespace where to count identities */
  namespace: string;
  /** Use cache or not. Cache may not be valid under very rare conditions (simultaniour read and writes). Cache automatically clears after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface CountIdentityResponse {
  /** Number of the identities in the provided namespace */
  count: number;
}

export interface GetServiceManagedIdentityRequest {
  /** Namespace where to search for identity */
  namespace: string;
  /** Service which manages this identity */
  service: string;
  /** Special ID for this identity defined by this service */
  managedId: string;
  /** Use cache or not. Cache may not be valid under very rare conditions (simultaniour read and writes). Cache automatically clears after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface GetServiceManagedIdentityResponse {
  /** Founded identity */
  identity: Identity | undefined;
}

export interface UpdateIdentityRequest {
  /** Namespace where identity is located */
  namespace: string;
  /** Unique identifier of the identity */
  uuid: string;
  /** New name */
  newName: string;
}

export interface UpdateIdentityResponse {
  /** Identity information after update */
  identity: Identity | undefined;
}

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

export interface AddRoleRequest {
  /** Identity namespace */
  identityNamespace: string;
  /** Identity identifier inside identity namespace */
  identityUUID: string;
  /** Role namespace */
  roleNamespace: string;
  /** Role UUID inside role namespace */
  roleUUID: string;
}

export interface AddRoleResponse {
  /** Updated identity (after adding role) */
  identity: Identity | undefined;
}

export interface RemoveRoleRequest {
  /** Identity namespace */
  identityNamespace: string;
  /** Identity unique identifier inside identity namespace */
  identityUUID: string;
  /** Role namespace */
  roleNamespace: string;
  /** Role UUID inside role namespace */
  roleUUID: string;
}

export interface RemoveRoleResponse {
  /** Updated identity (after removing role) */
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

function createBaseNotManagedData(): NotManagedData {
  return {};
}

export const NotManagedData = {
  encode(_: NotManagedData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): NotManagedData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseNotManagedData();
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

  fromJSON(_: any): NotManagedData {
    return {};
  },

  toJSON(_: NotManagedData): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<NotManagedData>, I>>(base?: I): NotManagedData {
    return NotManagedData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<NotManagedData>, I>>(_: I): NotManagedData {
    const message = createBaseNotManagedData();
    return message;
  },
};

function createBaseIdentityManagedData(): IdentityManagedData {
  return { identityNamespace: "", identityUUID: "" };
}

export const IdentityManagedData = {
  encode(message: IdentityManagedData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identityNamespace !== "") {
      writer.uint32(10).string(message.identityNamespace);
    }
    if (message.identityUUID !== "") {
      writer.uint32(18).string(message.identityUUID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): IdentityManagedData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentityManagedData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identityNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.identityUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): IdentityManagedData {
    return {
      identityNamespace: isSet(object.identityNamespace) ? String(object.identityNamespace) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
    };
  },

  toJSON(message: IdentityManagedData): unknown {
    const obj: any = {};
    if (message.identityNamespace !== "") {
      obj.identityNamespace = message.identityNamespace;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<IdentityManagedData>, I>>(base?: I): IdentityManagedData {
    return IdentityManagedData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<IdentityManagedData>, I>>(object: I): IdentityManagedData {
    const message = createBaseIdentityManagedData();
    message.identityNamespace = object.identityNamespace ?? "";
    message.identityUUID = object.identityUUID ?? "";
    return message;
  },
};

function createBaseServiceManagedData(): ServiceManagedData {
  return { service: "", reason: "", managementId: "" };
}

export const ServiceManagedData = {
  encode(message: ServiceManagedData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.service !== "") {
      writer.uint32(10).string(message.service);
    }
    if (message.reason !== "") {
      writer.uint32(18).string(message.reason);
    }
    if (message.managementId !== "") {
      writer.uint32(26).string(message.managementId);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ServiceManagedData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseServiceManagedData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.service = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.reason = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.managementId = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ServiceManagedData {
    return {
      service: isSet(object.service) ? String(object.service) : "",
      reason: isSet(object.reason) ? String(object.reason) : "",
      managementId: isSet(object.managementId) ? String(object.managementId) : "",
    };
  },

  toJSON(message: ServiceManagedData): unknown {
    const obj: any = {};
    if (message.service !== "") {
      obj.service = message.service;
    }
    if (message.reason !== "") {
      obj.reason = message.reason;
    }
    if (message.managementId !== "") {
      obj.managementId = message.managementId;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ServiceManagedData>, I>>(base?: I): ServiceManagedData {
    return ServiceManagedData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ServiceManagedData>, I>>(object: I): ServiceManagedData {
    const message = createBaseServiceManagedData();
    message.service = object.service ?? "";
    message.reason = object.reason ?? "";
    message.managementId = object.managementId ?? "";
    return message;
  },
};

function createBaseIdentity(): Identity {
  return {
    namespace: "",
    uuid: "",
    name: "",
    active: false,
    no: undefined,
    identity: undefined,
    service: undefined,
    policies: [],
    roles: [],
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const Identity = {
  encode(message: Identity, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    if (message.no !== undefined) {
      NotManagedData.encode(message.no, writer.uint32(162).fork()).ldelim();
    }
    if (message.identity !== undefined) {
      IdentityManagedData.encode(message.identity, writer.uint32(170).fork()).ldelim();
    }
    if (message.service !== undefined) {
      ServiceManagedData.encode(message.service, writer.uint32(178).fork()).ldelim();
    }
    for (const v of message.policies) {
      Identity_PolicyReference.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.roles) {
      Identity_RoleReference.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(74).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(82).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(88).uint64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Identity {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentity();
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

          message.uuid = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.name = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.active = reader.bool();
          continue;
        case 20:
          if (tag !== 162) {
            break;
          }

          message.no = NotManagedData.decode(reader, reader.uint32());
          continue;
        case 21:
          if (tag !== 170) {
            break;
          }

          message.identity = IdentityManagedData.decode(reader, reader.uint32());
          continue;
        case 22:
          if (tag !== 178) {
            break;
          }

          message.service = ServiceManagedData.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.policies.push(Identity_PolicyReference.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.roles.push(Identity_RoleReference.decode(reader, reader.uint32()));
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 88) {
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

  fromJSON(object: any): Identity {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      active: isSet(object.active) ? Boolean(object.active) : false,
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
      policies: Array.isArray(object?.policies)
        ? object.policies.map((e: any) => Identity_PolicyReference.fromJSON(e))
        : [],
      roles: Array.isArray(object?.roles) ? object.roles.map((e: any) => Identity_RoleReference.fromJSON(e)) : [],
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Identity): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.active === true) {
      obj.active = message.active;
    }
    if (message.no !== undefined) {
      obj.no = NotManagedData.toJSON(message.no);
    }
    if (message.identity !== undefined) {
      obj.identity = IdentityManagedData.toJSON(message.identity);
    }
    if (message.service !== undefined) {
      obj.service = ServiceManagedData.toJSON(message.service);
    }
    if (message.policies?.length) {
      obj.policies = message.policies.map((e) => Identity_PolicyReference.toJSON(e));
    }
    if (message.roles?.length) {
      obj.roles = message.roles.map((e) => Identity_RoleReference.toJSON(e));
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

  create<I extends Exact<DeepPartial<Identity>, I>>(base?: I): Identity {
    return Identity.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Identity>, I>>(object: I): Identity {
    const message = createBaseIdentity();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.active = object.active ?? false;
    message.no = (object.no !== undefined && object.no !== null) ? NotManagedData.fromPartial(object.no) : undefined;
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? IdentityManagedData.fromPartial(object.identity)
      : undefined;
    message.service = (object.service !== undefined && object.service !== null)
      ? ServiceManagedData.fromPartial(object.service)
      : undefined;
    message.policies = object.policies?.map((e) => Identity_PolicyReference.fromPartial(e)) || [];
    message.roles = object.roles?.map((e) => Identity_RoleReference.fromPartial(e)) || [];
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseIdentity_PolicyReference(): Identity_PolicyReference {
  return { namespace: "", uuid: "" };
}

export const Identity_PolicyReference = {
  encode(message: Identity_PolicyReference, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Identity_PolicyReference {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentity_PolicyReference();
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

          message.uuid = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Identity_PolicyReference>, I>>(base?: I): Identity_PolicyReference {
    return Identity_PolicyReference.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Identity_PolicyReference>, I>>(object: I): Identity_PolicyReference {
    const message = createBaseIdentity_PolicyReference();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseIdentity_RoleReference(): Identity_RoleReference {
  return { namespace: "", uuid: "" };
}

export const Identity_RoleReference = {
  encode(message: Identity_RoleReference, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Identity_RoleReference {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseIdentity_RoleReference();
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

          message.uuid = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Identity_RoleReference {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: Identity_RoleReference): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Identity_RoleReference>, I>>(base?: I): Identity_RoleReference {
    return Identity_RoleReference.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Identity_RoleReference>, I>>(object: I): Identity_RoleReference {
    const message = createBaseIdentity_RoleReference();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseCreateIdentityRequest(): CreateIdentityRequest {
  return { namespace: "", name: "", initiallyActive: false, no: undefined, identity: undefined, service: undefined };
}

export const CreateIdentityRequest = {
  encode(message: CreateIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.initiallyActive === true) {
      writer.uint32(24).bool(message.initiallyActive);
    }
    if (message.no !== undefined) {
      NotManagedData.encode(message.no, writer.uint32(162).fork()).ldelim();
    }
    if (message.identity !== undefined) {
      IdentityManagedData.encode(message.identity, writer.uint32(170).fork()).ldelim();
    }
    if (message.service !== undefined) {
      ServiceManagedData.encode(message.service, writer.uint32(178).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateIdentityRequest();
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

          message.initiallyActive = reader.bool();
          continue;
        case 20:
          if (tag !== 162) {
            break;
          }

          message.no = NotManagedData.decode(reader, reader.uint32());
          continue;
        case 21:
          if (tag !== 170) {
            break;
          }

          message.identity = IdentityManagedData.decode(reader, reader.uint32());
          continue;
        case 22:
          if (tag !== 178) {
            break;
          }

          message.service = ServiceManagedData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      initiallyActive: isSet(object.initiallyActive) ? Boolean(object.initiallyActive) : false,
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
    };
  },

  toJSON(message: CreateIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.initiallyActive === true) {
      obj.initiallyActive = message.initiallyActive;
    }
    if (message.no !== undefined) {
      obj.no = NotManagedData.toJSON(message.no);
    }
    if (message.identity !== undefined) {
      obj.identity = IdentityManagedData.toJSON(message.identity);
    }
    if (message.service !== undefined) {
      obj.service = ServiceManagedData.toJSON(message.service);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateIdentityRequest>, I>>(base?: I): CreateIdentityRequest {
    return CreateIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateIdentityRequest>, I>>(object: I): CreateIdentityRequest {
    const message = createBaseCreateIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.initiallyActive = object.initiallyActive ?? false;
    message.no = (object.no !== undefined && object.no !== null) ? NotManagedData.fromPartial(object.no) : undefined;
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? IdentityManagedData.fromPartial(object.identity)
      : undefined;
    message.service = (object.service !== undefined && object.service !== null)
      ? ServiceManagedData.fromPartial(object.service)
      : undefined;
    return message;
  },
};

function createBaseCreateIdentityResponse(): CreateIdentityResponse {
  return { identity: undefined };
}

export const CreateIdentityResponse = {
  encode(message: CreateIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateIdentityResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: CreateIdentityResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateIdentityResponse>, I>>(base?: I): CreateIdentityResponse {
    return CreateIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateIdentityResponse>, I>>(object: I): CreateIdentityResponse {
    const message = createBaseCreateIdentityResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseGetIdentityRequest(): GetIdentityRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetIdentityRequest = {
  encode(message: GetIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIdentityRequest();
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

          message.uuid = reader.string();
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

  fromJSON(object: any): GetIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetIdentityRequest>, I>>(base?: I): GetIdentityRequest {
    return GetIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetIdentityRequest>, I>>(object: I): GetIdentityRequest {
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
  encode(message: GetIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetIdentityResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: GetIdentityResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetIdentityResponse>, I>>(base?: I): GetIdentityResponse {
    return GetIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetIdentityResponse>, I>>(object: I): GetIdentityResponse {
    const message = createBaseGetIdentityResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseDeleteIdentityRequest(): DeleteIdentityRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteIdentityRequest = {
  encode(message: DeleteIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteIdentityRequest();
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

          message.uuid = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteIdentityRequest>, I>>(base?: I): DeleteIdentityRequest {
    return DeleteIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteIdentityRequest>, I>>(object: I): DeleteIdentityRequest {
    const message = createBaseDeleteIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteIdentityResponse(): DeleteIdentityResponse {
  return { existed: false };
}

export const DeleteIdentityResponse = {
  encode(message: DeleteIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.existed === true) {
      writer.uint32(8).bool(message.existed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteIdentityResponse();
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

  fromJSON(object: any): DeleteIdentityResponse {
    return { existed: isSet(object.existed) ? Boolean(object.existed) : false };
  },

  toJSON(message: DeleteIdentityResponse): unknown {
    const obj: any = {};
    if (message.existed === true) {
      obj.existed = message.existed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteIdentityResponse>, I>>(base?: I): DeleteIdentityResponse {
    return DeleteIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteIdentityResponse>, I>>(object: I): DeleteIdentityResponse {
    const message = createBaseDeleteIdentityResponse();
    message.existed = object.existed ?? false;
    return message;
  },
};

function createBaseExistsIdentityRequest(): ExistsIdentityRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const ExistsIdentityRequest = {
  encode(message: ExistsIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistsIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsIdentityRequest();
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

          message.uuid = reader.string();
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

  fromJSON(object: any): ExistsIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ExistsIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistsIdentityRequest>, I>>(base?: I): ExistsIdentityRequest {
    return ExistsIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistsIdentityRequest>, I>>(object: I): ExistsIdentityRequest {
    const message = createBaseExistsIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseExistsIdentityResponse(): ExistsIdentityResponse {
  return { exists: false };
}

export const ExistsIdentityResponse = {
  encode(message: ExistsIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exists === true) {
      writer.uint32(8).bool(message.exists);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistsIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.exists = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExistsIdentityResponse {
    return { exists: isSet(object.exists) ? Boolean(object.exists) : false };
  },

  toJSON(message: ExistsIdentityResponse): unknown {
    const obj: any = {};
    if (message.exists === true) {
      obj.exists = message.exists;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistsIdentityResponse>, I>>(base?: I): ExistsIdentityResponse {
    return ExistsIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistsIdentityResponse>, I>>(object: I): ExistsIdentityResponse {
    const message = createBaseExistsIdentityResponse();
    message.exists = object.exists ?? false;
    return message;
  },
};

function createBaseListIdentityRequest(): ListIdentityRequest {
  return { namespace: "", skip: 0, limit: 0 };
}

export const ListIdentityRequest = {
  encode(message: ListIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.skip !== 0) {
      writer.uint32(16).uint64(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(24).uint64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIdentityRequest();
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
          if (tag !== 16) {
            break;
          }

          message.skip = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.limit = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListIdentityRequest>, I>>(base?: I): ListIdentityRequest {
    return ListIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListIdentityRequest>, I>>(object: I): ListIdentityRequest {
    const message = createBaseListIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListIdentityResponse(): ListIdentityResponse {
  return { identity: undefined };
}

export const ListIdentityResponse = {
  encode(message: ListIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListIdentityResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: ListIdentityResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListIdentityResponse>, I>>(base?: I): ListIdentityResponse {
    return ListIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListIdentityResponse>, I>>(object: I): ListIdentityResponse {
    const message = createBaseListIdentityResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseCountIdentityRequest(): CountIdentityRequest {
  return { namespace: "", useCache: false };
}

export const CountIdentityRequest = {
  encode(message: CountIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountIdentityRequest();
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

  fromJSON(object: any): CountIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: CountIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountIdentityRequest>, I>>(base?: I): CountIdentityRequest {
    return CountIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountIdentityRequest>, I>>(object: I): CountIdentityRequest {
    const message = createBaseCountIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseCountIdentityResponse(): CountIdentityResponse {
  return { count: 0 };
}

export const CountIdentityResponse = {
  encode(message: CountIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.count = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountIdentityResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountIdentityResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountIdentityResponse>, I>>(base?: I): CountIdentityResponse {
    return CountIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountIdentityResponse>, I>>(object: I): CountIdentityResponse {
    const message = createBaseCountIdentityResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseGetServiceManagedIdentityRequest(): GetServiceManagedIdentityRequest {
  return { namespace: "", service: "", managedId: "", useCache: false };
}

export const GetServiceManagedIdentityRequest = {
  encode(message: GetServiceManagedIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.service !== "") {
      writer.uint32(18).string(message.service);
    }
    if (message.managedId !== "") {
      writer.uint32(26).string(message.managedId);
    }
    if (message.useCache === true) {
      writer.uint32(32).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedIdentityRequest();
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

          message.service = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.managedId = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
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

  fromJSON(object: any): GetServiceManagedIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      service: isSet(object.service) ? String(object.service) : "",
      managedId: isSet(object.managedId) ? String(object.managedId) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetServiceManagedIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.service !== "") {
      obj.service = message.service;
    }
    if (message.managedId !== "") {
      obj.managedId = message.managedId;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetServiceManagedIdentityRequest>, I>>(
    base?: I,
  ): GetServiceManagedIdentityRequest {
    return GetServiceManagedIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedIdentityRequest>, I>>(
    object: I,
  ): GetServiceManagedIdentityRequest {
    const message = createBaseGetServiceManagedIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.service = object.service ?? "";
    message.managedId = object.managedId ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetServiceManagedIdentityResponse(): GetServiceManagedIdentityResponse {
  return { identity: undefined };
}

export const GetServiceManagedIdentityResponse = {
  encode(message: GetServiceManagedIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetServiceManagedIdentityResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: GetServiceManagedIdentityResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetServiceManagedIdentityResponse>, I>>(
    base?: I,
  ): GetServiceManagedIdentityResponse {
    return GetServiceManagedIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedIdentityResponse>, I>>(
    object: I,
  ): GetServiceManagedIdentityResponse {
    const message = createBaseGetServiceManagedIdentityResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseUpdateIdentityRequest(): UpdateIdentityRequest {
  return { namespace: "", uuid: "", newName: "" };
}

export const UpdateIdentityRequest = {
  encode(message: UpdateIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.newName !== "") {
      writer.uint32(26).string(message.newName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIdentityRequest();
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

          message.uuid = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.newName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      newName: isSet(object.newName) ? String(object.newName) : "",
    };
  },

  toJSON(message: UpdateIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.newName !== "") {
      obj.newName = message.newName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateIdentityRequest>, I>>(base?: I): UpdateIdentityRequest {
    return UpdateIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateIdentityRequest>, I>>(object: I): UpdateIdentityRequest {
    const message = createBaseUpdateIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.newName = object.newName ?? "";
    return message;
  },
};

function createBaseUpdateIdentityResponse(): UpdateIdentityResponse {
  return { identity: undefined };
}

export const UpdateIdentityResponse = {
  encode(message: UpdateIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateIdentityResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: UpdateIdentityResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateIdentityResponse>, I>>(base?: I): UpdateIdentityResponse {
    return UpdateIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateIdentityResponse>, I>>(object: I): UpdateIdentityResponse {
    const message = createBaseUpdateIdentityResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseAddPolicyRequest(): AddPolicyRequest {
  return { identityNamespace: "", identityUUID: "", policyNamespace: "", policyUUID: "" };
}

export const AddPolicyRequest = {
  encode(message: AddPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddPolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identityNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.identityUUID = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.policyNamespace = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.policyUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddPolicyRequest {
    return {
      identityNamespace: isSet(object.identityNamespace) ? String(object.identityNamespace) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
      policyNamespace: isSet(object.policyNamespace) ? String(object.policyNamespace) : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: AddPolicyRequest): unknown {
    const obj: any = {};
    if (message.identityNamespace !== "") {
      obj.identityNamespace = message.identityNamespace;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    if (message.policyNamespace !== "") {
      obj.policyNamespace = message.policyNamespace;
    }
    if (message.policyUUID !== "") {
      obj.policyUUID = message.policyUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddPolicyRequest>, I>>(base?: I): AddPolicyRequest {
    return AddPolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddPolicyRequest>, I>>(object: I): AddPolicyRequest {
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
  encode(message: AddPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddPolicyResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: AddPolicyResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(base?: I): AddPolicyResponse {
    return AddPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(object: I): AddPolicyResponse {
    const message = createBaseAddPolicyResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseRemovePolicyRequest(): RemovePolicyRequest {
  return { identityNamespace: "", identityUUID: "", policyNamespace: "", policyUUID: "" };
}

export const RemovePolicyRequest = {
  encode(message: RemovePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemovePolicyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identityNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.identityUUID = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.policyNamespace = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.policyUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemovePolicyRequest {
    return {
      identityNamespace: isSet(object.identityNamespace) ? String(object.identityNamespace) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
      policyNamespace: isSet(object.policyNamespace) ? String(object.policyNamespace) : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: RemovePolicyRequest): unknown {
    const obj: any = {};
    if (message.identityNamespace !== "") {
      obj.identityNamespace = message.identityNamespace;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    if (message.policyNamespace !== "") {
      obj.policyNamespace = message.policyNamespace;
    }
    if (message.policyUUID !== "") {
      obj.policyUUID = message.policyUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemovePolicyRequest>, I>>(base?: I): RemovePolicyRequest {
    return RemovePolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemovePolicyRequest>, I>>(object: I): RemovePolicyRequest {
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
  encode(message: RemovePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemovePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemovePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemovePolicyResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: RemovePolicyResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(base?: I): RemovePolicyResponse {
    return RemovePolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(object: I): RemovePolicyResponse {
    const message = createBaseRemovePolicyResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseAddRoleRequest(): AddRoleRequest {
  return { identityNamespace: "", identityUUID: "", roleNamespace: "", roleUUID: "" };
}

export const AddRoleRequest = {
  encode(message: AddRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identityNamespace !== "") {
      writer.uint32(10).string(message.identityNamespace);
    }
    if (message.identityUUID !== "") {
      writer.uint32(18).string(message.identityUUID);
    }
    if (message.roleNamespace !== "") {
      writer.uint32(26).string(message.roleNamespace);
    }
    if (message.roleUUID !== "") {
      writer.uint32(34).string(message.roleUUID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddRoleRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identityNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.identityUUID = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.roleNamespace = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.roleUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddRoleRequest {
    return {
      identityNamespace: isSet(object.identityNamespace) ? String(object.identityNamespace) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
      roleNamespace: isSet(object.roleNamespace) ? String(object.roleNamespace) : "",
      roleUUID: isSet(object.roleUUID) ? String(object.roleUUID) : "",
    };
  },

  toJSON(message: AddRoleRequest): unknown {
    const obj: any = {};
    if (message.identityNamespace !== "") {
      obj.identityNamespace = message.identityNamespace;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    if (message.roleNamespace !== "") {
      obj.roleNamespace = message.roleNamespace;
    }
    if (message.roleUUID !== "") {
      obj.roleUUID = message.roleUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddRoleRequest>, I>>(base?: I): AddRoleRequest {
    return AddRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddRoleRequest>, I>>(object: I): AddRoleRequest {
    const message = createBaseAddRoleRequest();
    message.identityNamespace = object.identityNamespace ?? "";
    message.identityUUID = object.identityUUID ?? "";
    message.roleNamespace = object.roleNamespace ?? "";
    message.roleUUID = object.roleUUID ?? "";
    return message;
  },
};

function createBaseAddRoleResponse(): AddRoleResponse {
  return { identity: undefined };
}

export const AddRoleResponse = {
  encode(message: AddRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AddRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAddRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AddRoleResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: AddRoleResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddRoleResponse>, I>>(base?: I): AddRoleResponse {
    return AddRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddRoleResponse>, I>>(object: I): AddRoleResponse {
    const message = createBaseAddRoleResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseRemoveRoleRequest(): RemoveRoleRequest {
  return { identityNamespace: "", identityUUID: "", roleNamespace: "", roleUUID: "" };
}

export const RemoveRoleRequest = {
  encode(message: RemoveRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identityNamespace !== "") {
      writer.uint32(10).string(message.identityNamespace);
    }
    if (message.identityUUID !== "") {
      writer.uint32(18).string(message.identityUUID);
    }
    if (message.roleNamespace !== "") {
      writer.uint32(26).string(message.roleNamespace);
    }
    if (message.roleUUID !== "") {
      writer.uint32(34).string(message.roleUUID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveRoleRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identityNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.identityUUID = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.roleNamespace = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.roleUUID = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveRoleRequest {
    return {
      identityNamespace: isSet(object.identityNamespace) ? String(object.identityNamespace) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
      roleNamespace: isSet(object.roleNamespace) ? String(object.roleNamespace) : "",
      roleUUID: isSet(object.roleUUID) ? String(object.roleUUID) : "",
    };
  },

  toJSON(message: RemoveRoleRequest): unknown {
    const obj: any = {};
    if (message.identityNamespace !== "") {
      obj.identityNamespace = message.identityNamespace;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    if (message.roleNamespace !== "") {
      obj.roleNamespace = message.roleNamespace;
    }
    if (message.roleUUID !== "") {
      obj.roleUUID = message.roleUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveRoleRequest>, I>>(base?: I): RemoveRoleRequest {
    return RemoveRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveRoleRequest>, I>>(object: I): RemoveRoleRequest {
    const message = createBaseRemoveRoleRequest();
    message.identityNamespace = object.identityNamespace ?? "";
    message.identityUUID = object.identityUUID ?? "";
    message.roleNamespace = object.roleNamespace ?? "";
    message.roleUUID = object.roleUUID ?? "";
    return message;
  },
};

function createBaseRemoveRoleResponse(): RemoveRoleResponse {
  return { identity: undefined };
}

export const RemoveRoleResponse = {
  encode(message: RemoveRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveRoleResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: RemoveRoleResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveRoleResponse>, I>>(base?: I): RemoveRoleResponse {
    return RemoveRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveRoleResponse>, I>>(object: I): RemoveRoleResponse {
    const message = createBaseRemoveRoleResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? Identity.fromPartial(object.identity)
      : undefined;
    return message;
  },
};

function createBaseSetIdentityActiveRequest(): SetIdentityActiveRequest {
  return { namespace: "", uuid: "", active: false };
}

export const SetIdentityActiveRequest = {
  encode(message: SetIdentityActiveRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): SetIdentityActiveRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIdentityActiveRequest();
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

          message.uuid = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.active = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.active === true) {
      obj.active = message.active;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetIdentityActiveRequest>, I>>(base?: I): SetIdentityActiveRequest {
    return SetIdentityActiveRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveRequest>, I>>(object: I): SetIdentityActiveRequest {
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
  encode(message: SetIdentityActiveResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.identity !== undefined) {
      Identity.encode(message.identity, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetIdentityActiveResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetIdentityActiveResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.identity = Identity.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SetIdentityActiveResponse {
    return { identity: isSet(object.identity) ? Identity.fromJSON(object.identity) : undefined };
  },

  toJSON(message: SetIdentityActiveResponse): unknown {
    const obj: any = {};
    if (message.identity !== undefined) {
      obj.identity = Identity.toJSON(message.identity);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SetIdentityActiveResponse>, I>>(base?: I): SetIdentityActiveResponse {
    return SetIdentityActiveResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SetIdentityActiveResponse>, I>>(object: I): SetIdentityActiveResponse {
    const message = createBaseSetIdentityActiveResponse();
    message.identity = (object.identity !== undefined && object.identity !== null)
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
  /** Check if identity exists */
  Exists(request: ExistsIdentityRequest): Promise<ExistsIdentityResponse>;
  /** Get list of the identities */
  List(request: ListIdentityRequest): Observable<ListIdentityResponse>;
  /** Get number of the identities in the namespace */
  Count(request: CountIdentityRequest): Promise<CountIdentityResponse>;
  /** Get policy that is managed by service */
  GetServiceManagedIdentity(request: GetServiceManagedIdentityRequest): Promise<GetServiceManagedIdentityResponse>;
  /** Update identity information */
  Update(request: UpdateIdentityRequest): Promise<UpdateIdentityResponse>;
  /** Add policy to the identity. If policy was already added - does nothing. */
  AddPolicy(request: AddPolicyRequest): Promise<AddPolicyResponse>;
  /** Remove policy from the identity. If policy was already removed - does nothing. */
  RemovePolicy(request: RemovePolicyRequest): Promise<RemovePolicyResponse>;
  /** Add role to the identity. If role was already added = does nothing */
  AddRole(request: AddRoleRequest): Promise<AddRoleResponse>;
  /** Remove role from the identity. If role was already removed - does nothing. */
  RemoveRole(request: RemoveRoleRequest): Promise<RemoveRoleResponse>;
  /** Set identity active or not. */
  SetActive(request: SetIdentityActiveRequest): Promise<SetIdentityActiveResponse>;
}

export const IAMIdentityServiceServiceName = "native_iam_identity.IAMIdentityService";
export class IAMIdentityServiceClientImpl implements IAMIdentityService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMIdentityServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Exists = this.Exists.bind(this);
    this.List = this.List.bind(this);
    this.Count = this.Count.bind(this);
    this.GetServiceManagedIdentity = this.GetServiceManagedIdentity.bind(this);
    this.Update = this.Update.bind(this);
    this.AddPolicy = this.AddPolicy.bind(this);
    this.RemovePolicy = this.RemovePolicy.bind(this);
    this.AddRole = this.AddRole.bind(this);
    this.RemoveRole = this.RemoveRole.bind(this);
    this.SetActive = this.SetActive.bind(this);
  }
  Create(request: CreateIdentityRequest): Promise<CreateIdentityResponse> {
    const data = CreateIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateIdentityResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetIdentityRequest): Promise<GetIdentityResponse> {
    const data = GetIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetIdentityResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteIdentityRequest): Promise<DeleteIdentityResponse> {
    const data = DeleteIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteIdentityResponse.decode(_m0.Reader.create(data)));
  }

  Exists(request: ExistsIdentityRequest): Promise<ExistsIdentityResponse> {
    const data = ExistsIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exists", data);
    return promise.then((data) => ExistsIdentityResponse.decode(_m0.Reader.create(data)));
  }

  List(request: ListIdentityRequest): Observable<ListIdentityResponse> {
    const data = ListIdentityRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListIdentityResponse.decode(_m0.Reader.create(data))));
  }

  Count(request: CountIdentityRequest): Promise<CountIdentityResponse> {
    const data = CountIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountIdentityResponse.decode(_m0.Reader.create(data)));
  }

  GetServiceManagedIdentity(request: GetServiceManagedIdentityRequest): Promise<GetServiceManagedIdentityResponse> {
    const data = GetServiceManagedIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetServiceManagedIdentity", data);
    return promise.then((data) => GetServiceManagedIdentityResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateIdentityRequest): Promise<UpdateIdentityResponse> {
    const data = UpdateIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateIdentityResponse.decode(_m0.Reader.create(data)));
  }

  AddPolicy(request: AddPolicyRequest): Promise<AddPolicyResponse> {
    const data = AddPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "AddPolicy", data);
    return promise.then((data) => AddPolicyResponse.decode(_m0.Reader.create(data)));
  }

  RemovePolicy(request: RemovePolicyRequest): Promise<RemovePolicyResponse> {
    const data = RemovePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RemovePolicy", data);
    return promise.then((data) => RemovePolicyResponse.decode(_m0.Reader.create(data)));
  }

  AddRole(request: AddRoleRequest): Promise<AddRoleResponse> {
    const data = AddRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "AddRole", data);
    return promise.then((data) => AddRoleResponse.decode(_m0.Reader.create(data)));
  }

  RemoveRole(request: RemoveRoleRequest): Promise<RemoveRoleResponse> {
    const data = RemoveRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RemoveRole", data);
    return promise.then((data) => RemoveRoleResponse.decode(_m0.Reader.create(data)));
  }

  SetActive(request: SetIdentityActiveRequest): Promise<SetIdentityActiveResponse> {
    const data = SetIdentityActiveRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "SetActive", data);
    return promise.then((data) => SetIdentityActiveResponse.decode(_m0.Reader.create(data)));
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
