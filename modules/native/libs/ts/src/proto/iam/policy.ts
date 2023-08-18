/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_iam_policy";

export enum BuiltInPolicyType {
  /** GLOBAL_ROOT - Policy with full access to everything */
  GLOBAL_ROOT = 0,
  /** NAMESPACE_ROOT - Policy with full access to the namespace */
  NAMESPACE_ROOT = 1,
  /** EMPTY - Empty policy that gives nothing. Use it if you want to correlate something with namespace but dont want to give any permisions. */
  EMPTY = 2,
  UNRECOGNIZED = -1,
}

export function builtInPolicyTypeFromJSON(object: any): BuiltInPolicyType {
  switch (object) {
    case 0:
    case "GLOBAL_ROOT":
      return BuiltInPolicyType.GLOBAL_ROOT;
    case 1:
    case "NAMESPACE_ROOT":
      return BuiltInPolicyType.NAMESPACE_ROOT;
    case 2:
    case "EMPTY":
      return BuiltInPolicyType.EMPTY;
    case -1:
    case "UNRECOGNIZED":
    default:
      return BuiltInPolicyType.UNRECOGNIZED;
  }
}

export function builtInPolicyTypeToJSON(object: BuiltInPolicyType): string {
  switch (object) {
    case BuiltInPolicyType.GLOBAL_ROOT:
      return "GLOBAL_ROOT";
    case BuiltInPolicyType.NAMESPACE_ROOT:
      return "NAMESPACE_ROOT";
    case BuiltInPolicyType.EMPTY:
      return "EMPTY";
    case BuiltInPolicyType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** Empty information to indicate that management for this role is not defined. */
export interface NotManagedData {
}

/** Information about built in policy */
export interface BuiltInManagedData {
  /** Type of the builtin policy */
  type: BuiltInPolicyType;
}

/** Information about identity that manages this policy */
export interface IdentityManagedData {
  /** Namespace where identity is located */
  identityNamespace: string;
  /** Identity UUID inside this namespace */
  identityUUID: string;
}

/** Handles information about the service that manages this policy */
export interface ServiceManagedData {
  /** Name of the service */
  service: string;
  /** Reason why this service created this policy */
  reason: string;
  /** This is an ID that can be defined by managed service to find this policy. Set to empty string if you dont need this. If this is not empty Service+ID combination is unique. */
  managementId: string;
}

export interface Policy {
  /** Namespace where policy was created. Namespace can be empty for global policy. */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Public name */
  name: string;
  /** Arbitrary description */
  description: string;
  /** Policy is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Policy is builtIn and predifined */
  builtIn?:
    | BuiltInManagedData
    | undefined;
  /** Policy is managed by identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Policy is managed by service */
  service?:
    | ServiceManagedData
    | undefined;
  /** Indicates if this policy works in all namespaces or only in the namespace where it is defined */
  namespaceIndependent: boolean;
  /** List of resource for wich actions will be performed */
  resources: string[];
  /** List of actions that can be performed */
  actions: string[];
  /** List of tags associated with this policy */
  tags: string[];
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

export interface CreatePolicyRequest {
  /** Namespace where policy will be created. Namespace can be empty for global policy. */
  namespace: string;
  /** Public name. May not be unique. */
  name: string;
  /** Arbitrary description */
  description: string;
  /** Policy is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Policy is managed by identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Policy is managed by service */
  service?:
    | ServiceManagedData
    | undefined;
  /** Indicates if this policy works in all namespaces or only in the namespace where it is defined */
  namespaceIndependent: boolean;
  /** List of resource for wich actions will be performed */
  resources: string[];
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

export interface GetMultiplePoliciesRequest {
  /** List of policies to get */
  policies: GetMultiplePoliciesRequest_RequestedPolicy[];
}

/** Hold information on where to find the policy */
export interface GetMultiplePoliciesRequest_RequestedPolicy {
  /** Namespace where to search for policy. Leave empty for global policy. */
  namespace: string;
  /** Unique identifier of the policy inside searched namespace */
  uuid: string;
}

export interface GetMultiplePoliciesResponse {
  /** Founded policy. The ordering is random. */
  policy: Policy | undefined;
}

export interface ExistPolicyRequest {
  /** Namespace of the policy */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Use cache or not. Cache may be invalid under very rare conditions (simultanious read and writes to the policy while it is not in cache). Cache automatically deletes after short period of time (30 seconds by default). */
  useCache: boolean;
}

export interface ExistPolicyResponse {
  /** True if policy exists, false if not */
  exist: boolean;
}

export interface UpdatePolicyRequest {
  /** Namespace of the policy */
  namespace: string;
  /** Unique identifier of the policy in the namespace */
  uuid: string;
  /** Public name */
  name: string;
  /** Arbitrary description */
  description: string;
  /** Indicates if this policy works in all namespaces or only in the namespace where it is defined */
  namespaceIndependent: boolean;
  /** List of resource for wich actions will be performed */
  resources: string[];
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

export interface DeletePolicyResponse {
  /** Indicates if policy existed before this request */
  existed: boolean;
}

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

export interface CountPoliciesRequest {
  /** Namespace where to count policies */
  namespace: string;
  /** Use cache or not. Cached policy data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface CountPoliciesResponse {
  /** Count of policies in specified namespace */
  count: number;
}

export interface GetServiceManagedPolicyRequest {
  /** Namespace where to search for policy */
  namespace: string;
  /** Service which manages this policy */
  service: string;
  /** Special ID for this policy defined by this service */
  managedId: string;
  /** Use cache or not. Cached policy data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface GetServiceManagedPolicyResponse {
  /** Founded policy */
  policy: Policy | undefined;
}

export interface GetBuiltInPolicyRequest {
  /** Namespace where to get builtin policy. */
  namespace: string;
  /** Type of the policy to search */
  type: BuiltInPolicyType;
}

export interface GetBuiltInPolicyResponse {
  policy: Policy | undefined;
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

function createBaseBuiltInManagedData(): BuiltInManagedData {
  return { type: 0 };
}

export const BuiltInManagedData = {
  encode(message: BuiltInManagedData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.type !== 0) {
      writer.uint32(8).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): BuiltInManagedData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseBuiltInManagedData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): BuiltInManagedData {
    return { type: isSet(object.type) ? builtInPolicyTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: BuiltInManagedData): unknown {
    const obj: any = {};
    if (message.type !== 0) {
      obj.type = builtInPolicyTypeToJSON(message.type);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<BuiltInManagedData>, I>>(base?: I): BuiltInManagedData {
    return BuiltInManagedData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<BuiltInManagedData>, I>>(object: I): BuiltInManagedData {
    const message = createBaseBuiltInManagedData();
    message.type = object.type ?? 0;
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

function createBasePolicy(): Policy {
  return {
    namespace: "",
    uuid: "",
    name: "",
    description: "",
    no: undefined,
    builtIn: undefined,
    identity: undefined,
    service: undefined,
    namespaceIndependent: false,
    resources: [],
    actions: [],
    tags: [],
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const Policy = {
  encode(message: Policy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    if (message.no !== undefined) {
      NotManagedData.encode(message.no, writer.uint32(162).fork()).ldelim();
    }
    if (message.builtIn !== undefined) {
      BuiltInManagedData.encode(message.builtIn, writer.uint32(170).fork()).ldelim();
    }
    if (message.identity !== undefined) {
      IdentityManagedData.encode(message.identity, writer.uint32(178).fork()).ldelim();
    }
    if (message.service !== undefined) {
      ServiceManagedData.encode(message.service, writer.uint32(186).fork()).ldelim();
    }
    if (message.namespaceIndependent === true) {
      writer.uint32(40).bool(message.namespaceIndependent);
    }
    for (const v of message.resources) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.actions) {
      writer.uint32(58).string(v!);
    }
    for (const v of message.tags) {
      writer.uint32(66).string(v!);
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

  decode(input: _m0.Reader | Uint8Array, length?: number): Policy {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBasePolicy();
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
          if (tag !== 34) {
            break;
          }

          message.description = reader.string();
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

          message.builtIn = BuiltInManagedData.decode(reader, reader.uint32());
          continue;
        case 22:
          if (tag !== 178) {
            break;
          }

          message.identity = IdentityManagedData.decode(reader, reader.uint32());
          continue;
        case 23:
          if (tag !== 186) {
            break;
          }

          message.service = ServiceManagedData.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.namespaceIndependent = reader.bool();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.resources.push(reader.string());
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.actions.push(reader.string());
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.tags.push(reader.string());
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

  fromJSON(object: any): Policy {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      builtIn: isSet(object.builtIn) ? BuiltInManagedData.fromJSON(object.builtIn) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
      namespaceIndependent: isSet(object.namespaceIndependent) ? Boolean(object.namespaceIndependent) : false,
      resources: Array.isArray(object?.resources) ? object.resources.map((e: any) => String(e)) : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => String(e)) : [],
      tags: Array.isArray(object?.tags) ? object.tags.map((e: any) => String(e)) : [],
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Policy): unknown {
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
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.no !== undefined) {
      obj.no = NotManagedData.toJSON(message.no);
    }
    if (message.builtIn !== undefined) {
      obj.builtIn = BuiltInManagedData.toJSON(message.builtIn);
    }
    if (message.identity !== undefined) {
      obj.identity = IdentityManagedData.toJSON(message.identity);
    }
    if (message.service !== undefined) {
      obj.service = ServiceManagedData.toJSON(message.service);
    }
    if (message.namespaceIndependent === true) {
      obj.namespaceIndependent = message.namespaceIndependent;
    }
    if (message.resources?.length) {
      obj.resources = message.resources;
    }
    if (message.actions?.length) {
      obj.actions = message.actions;
    }
    if (message.tags?.length) {
      obj.tags = message.tags;
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

  create<I extends Exact<DeepPartial<Policy>, I>>(base?: I): Policy {
    return Policy.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Policy>, I>>(object: I): Policy {
    const message = createBasePolicy();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.no = (object.no !== undefined && object.no !== null) ? NotManagedData.fromPartial(object.no) : undefined;
    message.builtIn = (object.builtIn !== undefined && object.builtIn !== null)
      ? BuiltInManagedData.fromPartial(object.builtIn)
      : undefined;
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? IdentityManagedData.fromPartial(object.identity)
      : undefined;
    message.service = (object.service !== undefined && object.service !== null)
      ? ServiceManagedData.fromPartial(object.service)
      : undefined;
    message.namespaceIndependent = object.namespaceIndependent ?? false;
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    message.tags = object.tags?.map((e) => e) || [];
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseCreatePolicyRequest(): CreatePolicyRequest {
  return {
    namespace: "",
    name: "",
    description: "",
    no: undefined,
    identity: undefined,
    service: undefined,
    namespaceIndependent: false,
    resources: [],
    actions: [],
  };
}

export const CreatePolicyRequest = {
  encode(message: CreatePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(26).string(message.description);
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
    if (message.namespaceIndependent === true) {
      writer.uint32(32).bool(message.namespaceIndependent);
    }
    for (const v of message.resources) {
      writer.uint32(42).string(v!);
    }
    for (const v of message.actions) {
      writer.uint32(50).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreatePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreatePolicyRequest();
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

          message.description = reader.string();
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
        case 4:
          if (tag !== 32) {
            break;
          }

          message.namespaceIndependent = reader.bool();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.resources.push(reader.string());
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.actions.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreatePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
      namespaceIndependent: isSet(object.namespaceIndependent) ? Boolean(object.namespaceIndependent) : false,
      resources: Array.isArray(object?.resources) ? object.resources.map((e: any) => String(e)) : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: CreatePolicyRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.description !== "") {
      obj.description = message.description;
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
    if (message.namespaceIndependent === true) {
      obj.namespaceIndependent = message.namespaceIndependent;
    }
    if (message.resources?.length) {
      obj.resources = message.resources;
    }
    if (message.actions?.length) {
      obj.actions = message.actions;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreatePolicyRequest>, I>>(base?: I): CreatePolicyRequest {
    return CreatePolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreatePolicyRequest>, I>>(object: I): CreatePolicyRequest {
    const message = createBaseCreatePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.no = (object.no !== undefined && object.no !== null) ? NotManagedData.fromPartial(object.no) : undefined;
    message.identity = (object.identity !== undefined && object.identity !== null)
      ? IdentityManagedData.fromPartial(object.identity)
      : undefined;
    message.service = (object.service !== undefined && object.service !== null)
      ? ServiceManagedData.fromPartial(object.service)
      : undefined;
    message.namespaceIndependent = object.namespaceIndependent ?? false;
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseCreatePolicyResponse(): CreatePolicyResponse {
  return { policy: undefined };
}

export const CreatePolicyResponse = {
  encode(message: CreatePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreatePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreatePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreatePolicyResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: CreatePolicyResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreatePolicyResponse>, I>>(base?: I): CreatePolicyResponse {
    return CreatePolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreatePolicyResponse>, I>>(object: I): CreatePolicyResponse {
    const message = createBaseCreatePolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetPolicyRequest(): GetPolicyRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetPolicyRequest = {
  encode(message: GetPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPolicyRequest();
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

  fromJSON(object: any): GetPolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetPolicyRequest): unknown {
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

  create<I extends Exact<DeepPartial<GetPolicyRequest>, I>>(base?: I): GetPolicyRequest {
    return GetPolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetPolicyRequest>, I>>(object: I): GetPolicyRequest {
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
  encode(message: GetPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetPolicyResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetPolicyResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetPolicyResponse>, I>>(base?: I): GetPolicyResponse {
    return GetPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetPolicyResponse>, I>>(object: I): GetPolicyResponse {
    const message = createBaseGetPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetMultiplePoliciesRequest(): GetMultiplePoliciesRequest {
  return { policies: [] };
}

export const GetMultiplePoliciesRequest = {
  encode(message: GetMultiplePoliciesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.policies) {
      GetMultiplePoliciesRequest_RequestedPolicy.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultiplePoliciesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultiplePoliciesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policies.push(GetMultiplePoliciesRequest_RequestedPolicy.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetMultiplePoliciesRequest {
    return {
      policies: Array.isArray(object?.policies)
        ? object.policies.map((e: any) => GetMultiplePoliciesRequest_RequestedPolicy.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetMultiplePoliciesRequest): unknown {
    const obj: any = {};
    if (message.policies?.length) {
      obj.policies = message.policies.map((e) => GetMultiplePoliciesRequest_RequestedPolicy.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultiplePoliciesRequest>, I>>(base?: I): GetMultiplePoliciesRequest {
    return GetMultiplePoliciesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultiplePoliciesRequest>, I>>(object: I): GetMultiplePoliciesRequest {
    const message = createBaseGetMultiplePoliciesRequest();
    message.policies = object.policies?.map((e) => GetMultiplePoliciesRequest_RequestedPolicy.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetMultiplePoliciesRequest_RequestedPolicy(): GetMultiplePoliciesRequest_RequestedPolicy {
  return { namespace: "", uuid: "" };
}

export const GetMultiplePoliciesRequest_RequestedPolicy = {
  encode(message: GetMultiplePoliciesRequest_RequestedPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultiplePoliciesRequest_RequestedPolicy {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultiplePoliciesRequest_RequestedPolicy();
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

  fromJSON(object: any): GetMultiplePoliciesRequest_RequestedPolicy {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: GetMultiplePoliciesRequest_RequestedPolicy): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultiplePoliciesRequest_RequestedPolicy>, I>>(
    base?: I,
  ): GetMultiplePoliciesRequest_RequestedPolicy {
    return GetMultiplePoliciesRequest_RequestedPolicy.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultiplePoliciesRequest_RequestedPolicy>, I>>(
    object: I,
  ): GetMultiplePoliciesRequest_RequestedPolicy {
    const message = createBaseGetMultiplePoliciesRequest_RequestedPolicy();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseGetMultiplePoliciesResponse(): GetMultiplePoliciesResponse {
  return { policy: undefined };
}

export const GetMultiplePoliciesResponse = {
  encode(message: GetMultiplePoliciesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultiplePoliciesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultiplePoliciesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetMultiplePoliciesResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetMultiplePoliciesResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultiplePoliciesResponse>, I>>(base?: I): GetMultiplePoliciesResponse {
    return GetMultiplePoliciesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultiplePoliciesResponse>, I>>(object: I): GetMultiplePoliciesResponse {
    const message = createBaseGetMultiplePoliciesResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseExistPolicyRequest(): ExistPolicyRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const ExistPolicyRequest = {
  encode(message: ExistPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistPolicyRequest();
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

  fromJSON(object: any): ExistPolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ExistPolicyRequest): unknown {
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

  create<I extends Exact<DeepPartial<ExistPolicyRequest>, I>>(base?: I): ExistPolicyRequest {
    return ExistPolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistPolicyRequest>, I>>(object: I): ExistPolicyRequest {
    const message = createBaseExistPolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseExistPolicyResponse(): ExistPolicyResponse {
  return { exist: false };
}

export const ExistPolicyResponse = {
  encode(message: ExistPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistPolicyResponse();
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

  fromJSON(object: any): ExistPolicyResponse {
    return { exist: isSet(object.exist) ? Boolean(object.exist) : false };
  },

  toJSON(message: ExistPolicyResponse): unknown {
    const obj: any = {};
    if (message.exist === true) {
      obj.exist = message.exist;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistPolicyResponse>, I>>(base?: I): ExistPolicyResponse {
    return ExistPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistPolicyResponse>, I>>(object: I): ExistPolicyResponse {
    const message = createBaseExistPolicyResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

function createBaseUpdatePolicyRequest(): UpdatePolicyRequest {
  return {
    namespace: "",
    uuid: "",
    name: "",
    description: "",
    namespaceIndependent: false,
    resources: [],
    actions: [],
  };
}

export const UpdatePolicyRequest = {
  encode(message: UpdatePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    if (message.namespaceIndependent === true) {
      writer.uint32(40).bool(message.namespaceIndependent);
    }
    for (const v of message.resources) {
      writer.uint32(50).string(v!);
    }
    for (const v of message.actions) {
      writer.uint32(58).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePolicyRequest();
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
          if (tag !== 34) {
            break;
          }

          message.description = reader.string();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.namespaceIndependent = reader.bool();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.resources.push(reader.string());
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.actions.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdatePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      namespaceIndependent: isSet(object.namespaceIndependent) ? Boolean(object.namespaceIndependent) : false,
      resources: Array.isArray(object?.resources) ? object.resources.map((e: any) => String(e)) : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => String(e)) : [],
    };
  },

  toJSON(message: UpdatePolicyRequest): unknown {
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
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.namespaceIndependent === true) {
      obj.namespaceIndependent = message.namespaceIndependent;
    }
    if (message.resources?.length) {
      obj.resources = message.resources;
    }
    if (message.actions?.length) {
      obj.actions = message.actions;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdatePolicyRequest>, I>>(base?: I): UpdatePolicyRequest {
    return UpdatePolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdatePolicyRequest>, I>>(object: I): UpdatePolicyRequest {
    const message = createBaseUpdatePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.name = object.name ?? "";
    message.description = object.description ?? "";
    message.namespaceIndependent = object.namespaceIndependent ?? false;
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseUpdatePolicyResponse(): UpdatePolicyResponse {
  return { policy: undefined };
}

export const UpdatePolicyResponse = {
  encode(message: UpdatePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdatePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdatePolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdatePolicyResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: UpdatePolicyResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdatePolicyResponse>, I>>(base?: I): UpdatePolicyResponse {
    return UpdatePolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdatePolicyResponse>, I>>(object: I): UpdatePolicyResponse {
    const message = createBaseUpdatePolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseDeletePolicyRequest(): DeletePolicyRequest {
  return { namespace: "", uuid: "" };
}

export const DeletePolicyRequest = {
  encode(message: DeletePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeletePolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePolicyRequest();
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

  fromJSON(object: any): DeletePolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeletePolicyRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeletePolicyRequest>, I>>(base?: I): DeletePolicyRequest {
    return DeletePolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeletePolicyRequest>, I>>(object: I): DeletePolicyRequest {
    const message = createBaseDeletePolicyRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeletePolicyResponse(): DeletePolicyResponse {
  return { existed: false };
}

export const DeletePolicyResponse = {
  encode(message: DeletePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.existed === true) {
      writer.uint32(8).bool(message.existed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeletePolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeletePolicyResponse();
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

  fromJSON(object: any): DeletePolicyResponse {
    return { existed: isSet(object.existed) ? Boolean(object.existed) : false };
  },

  toJSON(message: DeletePolicyResponse): unknown {
    const obj: any = {};
    if (message.existed === true) {
      obj.existed = message.existed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeletePolicyResponse>, I>>(base?: I): DeletePolicyResponse {
    return DeletePolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeletePolicyResponse>, I>>(object: I): DeletePolicyResponse {
    const message = createBaseDeletePolicyResponse();
    message.existed = object.existed ?? false;
    return message;
  },
};

function createBaseListPoliciesRequest(): ListPoliciesRequest {
  return { namespace: "", skip: 0, limit: 0 };
}

export const ListPoliciesRequest = {
  encode(message: ListPoliciesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ListPoliciesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListPoliciesRequest();
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

  fromJSON(object: any): ListPoliciesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListPoliciesRequest): unknown {
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

  create<I extends Exact<DeepPartial<ListPoliciesRequest>, I>>(base?: I): ListPoliciesRequest {
    return ListPoliciesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListPoliciesRequest>, I>>(object: I): ListPoliciesRequest {
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
  encode(message: ListPoliciesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListPoliciesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListPoliciesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListPoliciesResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: ListPoliciesResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListPoliciesResponse>, I>>(base?: I): ListPoliciesResponse {
    return ListPoliciesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListPoliciesResponse>, I>>(object: I): ListPoliciesResponse {
    const message = createBaseListPoliciesResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseCountPoliciesRequest(): CountPoliciesRequest {
  return { namespace: "", useCache: false };
}

export const CountPoliciesRequest = {
  encode(message: CountPoliciesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountPoliciesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountPoliciesRequest();
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

  fromJSON(object: any): CountPoliciesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: CountPoliciesRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountPoliciesRequest>, I>>(base?: I): CountPoliciesRequest {
    return CountPoliciesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountPoliciesRequest>, I>>(object: I): CountPoliciesRequest {
    const message = createBaseCountPoliciesRequest();
    message.namespace = object.namespace ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseCountPoliciesResponse(): CountPoliciesResponse {
  return { count: 0 };
}

export const CountPoliciesResponse = {
  encode(message: CountPoliciesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountPoliciesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountPoliciesResponse();
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

  fromJSON(object: any): CountPoliciesResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountPoliciesResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountPoliciesResponse>, I>>(base?: I): CountPoliciesResponse {
    return CountPoliciesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountPoliciesResponse>, I>>(object: I): CountPoliciesResponse {
    const message = createBaseCountPoliciesResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseGetServiceManagedPolicyRequest(): GetServiceManagedPolicyRequest {
  return { namespace: "", service: "", managedId: "", useCache: false };
}

export const GetServiceManagedPolicyRequest = {
  encode(message: GetServiceManagedPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedPolicyRequest();
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

  fromJSON(object: any): GetServiceManagedPolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      service: isSet(object.service) ? String(object.service) : "",
      managedId: isSet(object.managedId) ? String(object.managedId) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetServiceManagedPolicyRequest): unknown {
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

  create<I extends Exact<DeepPartial<GetServiceManagedPolicyRequest>, I>>(base?: I): GetServiceManagedPolicyRequest {
    return GetServiceManagedPolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedPolicyRequest>, I>>(
    object: I,
  ): GetServiceManagedPolicyRequest {
    const message = createBaseGetServiceManagedPolicyRequest();
    message.namespace = object.namespace ?? "";
    message.service = object.service ?? "";
    message.managedId = object.managedId ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetServiceManagedPolicyResponse(): GetServiceManagedPolicyResponse {
  return { policy: undefined };
}

export const GetServiceManagedPolicyResponse = {
  encode(message: GetServiceManagedPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetServiceManagedPolicyResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetServiceManagedPolicyResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetServiceManagedPolicyResponse>, I>>(base?: I): GetServiceManagedPolicyResponse {
    return GetServiceManagedPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedPolicyResponse>, I>>(
    object: I,
  ): GetServiceManagedPolicyResponse {
    const message = createBaseGetServiceManagedPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
      ? Policy.fromPartial(object.policy)
      : undefined;
    return message;
  },
};

function createBaseGetBuiltInPolicyRequest(): GetBuiltInPolicyRequest {
  return { namespace: "", type: 0 };
}

export const GetBuiltInPolicyRequest = {
  encode(message: GetBuiltInPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBuiltInPolicyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBuiltInPolicyRequest();
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

          message.type = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBuiltInPolicyRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      type: isSet(object.type) ? builtInPolicyTypeFromJSON(object.type) : 0,
    };
  },

  toJSON(message: GetBuiltInPolicyRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.type !== 0) {
      obj.type = builtInPolicyTypeToJSON(message.type);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBuiltInPolicyRequest>, I>>(base?: I): GetBuiltInPolicyRequest {
    return GetBuiltInPolicyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBuiltInPolicyRequest>, I>>(object: I): GetBuiltInPolicyRequest {
    const message = createBaseGetBuiltInPolicyRequest();
    message.namespace = object.namespace ?? "";
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseGetBuiltInPolicyResponse(): GetBuiltInPolicyResponse {
  return { policy: undefined };
}

export const GetBuiltInPolicyResponse = {
  encode(message: GetBuiltInPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.policy !== undefined) {
      Policy.encode(message.policy, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBuiltInPolicyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBuiltInPolicyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.policy = Policy.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBuiltInPolicyResponse {
    return { policy: isSet(object.policy) ? Policy.fromJSON(object.policy) : undefined };
  },

  toJSON(message: GetBuiltInPolicyResponse): unknown {
    const obj: any = {};
    if (message.policy !== undefined) {
      obj.policy = Policy.toJSON(message.policy);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBuiltInPolicyResponse>, I>>(base?: I): GetBuiltInPolicyResponse {
    return GetBuiltInPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBuiltInPolicyResponse>, I>>(object: I): GetBuiltInPolicyResponse {
    const message = createBaseGetBuiltInPolicyResponse();
    message.policy = (object.policy !== undefined && object.policy !== null)
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
  /** Get multiple policies. */
  GetMultiple(request: GetMultiplePoliciesRequest): Observable<GetMultiplePoliciesResponse>;
  /** Check if policy exist or not */
  Exist(request: ExistPolicyRequest): Promise<ExistPolicyResponse>;
  /** Update policy */
  Update(request: UpdatePolicyRequest): Promise<UpdatePolicyResponse>;
  /** Delete policy */
  Delete(request: DeletePolicyRequest): Promise<DeletePolicyResponse>;
  /** List policies in namespace */
  List(request: ListPoliciesRequest): Observable<ListPoliciesResponse>;
  /** Count policies in namespace */
  Count(request: CountPoliciesRequest): Promise<CountPoliciesResponse>;
  /** Get policy that is managed by service */
  GetServiceManagedPolicy(request: GetServiceManagedPolicyRequest): Promise<GetServiceManagedPolicyResponse>;
  /** Get one of the builtin policies */
  GetBuiltInPolicy(request: GetBuiltInPolicyRequest): Promise<GetBuiltInPolicyResponse>;
}

export const IAMPolicyServiceServiceName = "native_iam_policy.IAMPolicyService";
export class IAMPolicyServiceClientImpl implements IAMPolicyService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMPolicyServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.GetMultiple = this.GetMultiple.bind(this);
    this.Exist = this.Exist.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.List = this.List.bind(this);
    this.Count = this.Count.bind(this);
    this.GetServiceManagedPolicy = this.GetServiceManagedPolicy.bind(this);
    this.GetBuiltInPolicy = this.GetBuiltInPolicy.bind(this);
  }
  Create(request: CreatePolicyRequest): Promise<CreatePolicyResponse> {
    const data = CreatePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreatePolicyResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetPolicyRequest): Promise<GetPolicyResponse> {
    const data = GetPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetPolicyResponse.decode(_m0.Reader.create(data)));
  }

  GetMultiple(request: GetMultiplePoliciesRequest): Observable<GetMultiplePoliciesResponse> {
    const data = GetMultiplePoliciesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "GetMultiple", data);
    return result.pipe(map((data) => GetMultiplePoliciesResponse.decode(_m0.Reader.create(data))));
  }

  Exist(request: ExistPolicyRequest): Promise<ExistPolicyResponse> {
    const data = ExistPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exist", data);
    return promise.then((data) => ExistPolicyResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdatePolicyRequest): Promise<UpdatePolicyResponse> {
    const data = UpdatePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdatePolicyResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeletePolicyRequest): Promise<DeletePolicyResponse> {
    const data = DeletePolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeletePolicyResponse.decode(_m0.Reader.create(data)));
  }

  List(request: ListPoliciesRequest): Observable<ListPoliciesResponse> {
    const data = ListPoliciesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListPoliciesResponse.decode(_m0.Reader.create(data))));
  }

  Count(request: CountPoliciesRequest): Promise<CountPoliciesResponse> {
    const data = CountPoliciesRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountPoliciesResponse.decode(_m0.Reader.create(data)));
  }

  GetServiceManagedPolicy(request: GetServiceManagedPolicyRequest): Promise<GetServiceManagedPolicyResponse> {
    const data = GetServiceManagedPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetServiceManagedPolicy", data);
    return promise.then((data) => GetServiceManagedPolicyResponse.decode(_m0.Reader.create(data)));
  }

  GetBuiltInPolicy(request: GetBuiltInPolicyRequest): Promise<GetBuiltInPolicyResponse> {
    const data = GetBuiltInPolicyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetBuiltInPolicy", data);
    return promise.then((data) => GetBuiltInPolicyResponse.decode(_m0.Reader.create(data)));
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
