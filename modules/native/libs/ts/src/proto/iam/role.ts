/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_iam_role";

export enum BuiltInRoleType {
  /** GLOBAL_ROOT - Role with full access to everything */
  GLOBAL_ROOT = 0,
  /** NAMESPACE_ROOT - Role with full access to namespace */
  NAMESPACE_ROOT = 1,
  /** EMPTY - Empty role that gives nothing. Use it if you want to correlate something with namespace but dont want to give any permisions. */
  EMPTY = 2,
  UNRECOGNIZED = -1,
}

export function builtInRoleTypeFromJSON(object: any): BuiltInRoleType {
  switch (object) {
    case 0:
    case "GLOBAL_ROOT":
      return BuiltInRoleType.GLOBAL_ROOT;
    case 1:
    case "NAMESPACE_ROOT":
      return BuiltInRoleType.NAMESPACE_ROOT;
    case 2:
    case "EMPTY":
      return BuiltInRoleType.EMPTY;
    case -1:
    case "UNRECOGNIZED":
    default:
      return BuiltInRoleType.UNRECOGNIZED;
  }
}

export function builtInRoleTypeToJSON(object: BuiltInRoleType): string {
  switch (object) {
    case BuiltInRoleType.GLOBAL_ROOT:
      return "GLOBAL_ROOT";
    case BuiltInRoleType.NAMESPACE_ROOT:
      return "NAMESPACE_ROOT";
    case BuiltInRoleType.EMPTY:
      return "EMPTY";
    case BuiltInRoleType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** Respresents pointer to the assigned policy */
export interface AssignedPolicy {
  /** Namespace of the assigned policy */
  namespace: string;
  /** Unique identifier of the assigned policy inside namespace */
  uuid: string;
}

/** Empty information to indicate that management for this role is not defined. */
export interface NotManagedData {
}

/** Information about built in role */
export interface BuiltInManagedData {
  /** Type of the builtin role */
  type: BuiltInRoleType;
}

/** Information about identity that manages this role */
export interface IdentityManagedData {
  /** Namespace where identity is located */
  identityNamespace: string;
  /** Identity UUID inside this namespace */
  identityUUID: string;
}

/** Handles information about the service that manages this role */
export interface ServiceManagedData {
  /** Name of the service */
  service: string;
  /** Reason why this service created this role */
  reason: string;
  /** This is an ID that can be defined by managed service to find this role. Set to empty string if you dont need this. If this is not empty Service+ID combination is unique. */
  managementId: string;
}

/** Role is a group of scopes */
export interface Role {
  /** Namespace of the role. Empty for global roles */
  namespace: string;
  /** Unique identifier of the role in the namespace */
  uuid: string;
  /** Short, human-readable name */
  name: string;
  /** Arbitrary description */
  description: string;
  /** Role is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Role is builtIn and predifined */
  builtIn?:
    | BuiltInManagedData
    | undefined;
  /** Role is managed by identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Role is managed by service */
  service?:
    | ServiceManagedData
    | undefined;
  /** List of the policies assigned to the role */
  policies: AssignedPolicy[];
  /** List of tags associated with this role */
  tags: string[];
  /** When the role was created */
  created:
    | Date
    | undefined;
  /** Last time when the role information was updated. */
  updated:
    | Date
    | undefined;
  /** Counter that increases after every update of the role */
  version: number;
}

export interface CreateRoleRequest {
  /** Namespace where role will be located. Leave empty for global role. */
  namespace: string;
  /** Short, human-readable name */
  name: string;
  /** Arbitrary description */
  description: string;
  /** Role is not managed */
  no?:
    | NotManagedData
    | undefined;
  /** Role is managed by identity */
  identity?:
    | IdentityManagedData
    | undefined;
  /** Role is managed by service */
  service?: ServiceManagedData | undefined;
}

export interface CreateRoleResponse {
  /** Created role */
  role: Role | undefined;
}

export interface GetRoleRequest {
  /** Namespace where to search for role. Leave empty for global role. */
  namespace: string;
  /** Unique identifier of the role inside searched namespace */
  uuid: string;
  /** Use cache or not. Cached role data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface GetRoleResponse {
  /** Founded role */
  role: Role | undefined;
}

export interface GetMultipleRolesRequest {
  /** List of roles to get */
  roles: GetMultipleRolesRequest_RequestedRole[];
}

/** Hold information on where to find the role */
export interface GetMultipleRolesRequest_RequestedRole {
  /** Namespace where to search for role. Leave empty for global role. */
  namespace: string;
  /** Unique identifier of the role inside searched namespace */
  uuid: string;
}

export interface GetMultipleRolesResponse {
  /** Founded role. The ordering is random. */
  role: Role | undefined;
}

export interface ListRolesRequest {
  /** Namesapce where to search for roles */
  namespace: string;
  /** Skip specified number of roles before return response. All the roleas are returned sorted, so this parameter can be used for list pagination. */
  skip: number;
  /** Limits the count of returned responses to specified value */
  limit: number;
}

export interface ListRolesResponse {
  role: Role | undefined;
}

export interface CountRolesRequest {
  /** Namespace where to count roles */
  namespace: string;
  /** Use cache or not. Cached role data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface CountRolesResponse {
  /** Total count of roles in the namespace */
  count: number;
}

export interface UpdateRoleRequest {
  /** Namespace where the role is located */
  namespace: string;
  /** Unique identifier of the role to update */
  uuid: string;
  newName: string;
  newDescription: string;
}

export interface UpdateRoleResponse {
  /** Updated role */
  role: Role | undefined;
}

export interface DeleteRoleRequest {
  /** Namespace where to search for role. Leave empty for global role. */
  namespace: string;
  /** Unique identifier of the role inside searched namespace */
  uuid: string;
}

export interface DeleteRoleResponse {
  /** Indicates if role existed before this request. */
  existed: boolean;
}

export interface GetServiceManagedRoleRequest {
  /** Namespace where to search for role */
  namespace: string;
  /** Service which manages this role */
  service: string;
  /** Special ID for this role defined by this service */
  managedId: string;
  /** Use cache or not. Cached role data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface GetServiceManagedRoleResponse {
  /** Founded role */
  role: Role | undefined;
}

export interface GetBuiltInRoleRequest {
  /** Namespace where to get builtin role. */
  namespace: string;
  /** Type of the role to search */
  type: BuiltInRoleType;
}

export interface GetBuiltInRoleResponse {
  role: Role | undefined;
}

export interface AddPolicyRequest {
  /** Namespace where role is located */
  roleNamespace: string;
  /** Unique identifier of the role */
  roleUUID: string;
  /** Namespace where policy is located */
  policyNamespace: string;
  /** Unique identifier of the policy */
  policyUUID: string;
}

export interface AddPolicyResponse {
  /** Role after this operation */
  role: Role | undefined;
}

export interface RemovePolicyRequest {
  /** Namespace where role is located */
  roleNamespace: string;
  /** Unique identifier of the role */
  roleUUID: string;
  /** Namespace where policy is located */
  policyNamespace: string;
  /** Unique identifier of the policy */
  policyUUID: string;
}

export interface RemovePolicyResponse {
  /** Role after this operation */
  role: Role | undefined;
}

export interface ExistRoleRequest {
  /** Namespace where to search for role. Leave empty for global role. */
  namespace: string;
  /** Unique identifier of the role inside searched namespace */
  uuid: string;
  /** Use cache or not. Cached role data may not be actual under very rare conditions. Invalid cache data is automatically clear after short period of time. */
  useCache: boolean;
}

export interface ExistRoleResponse {
  /** Indicates if role exists or not */
  exist: boolean;
}

function createBaseAssignedPolicy(): AssignedPolicy {
  return { namespace: "", uuid: "" };
}

export const AssignedPolicy = {
  encode(message: AssignedPolicy, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AssignedPolicy {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAssignedPolicy();
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

  fromJSON(object: any): AssignedPolicy {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: AssignedPolicy): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AssignedPolicy>, I>>(base?: I): AssignedPolicy {
    return AssignedPolicy.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AssignedPolicy>, I>>(object: I): AssignedPolicy {
    const message = createBaseAssignedPolicy();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

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
    return { type: isSet(object.type) ? builtInRoleTypeFromJSON(object.type) : 0 };
  },

  toJSON(message: BuiltInManagedData): unknown {
    const obj: any = {};
    if (message.type !== 0) {
      obj.type = builtInRoleTypeToJSON(message.type);
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

function createBaseRole(): Role {
  return {
    namespace: "",
    uuid: "",
    name: "",
    description: "",
    no: undefined,
    builtIn: undefined,
    identity: undefined,
    service: undefined,
    policies: [],
    tags: [],
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const Role = {
  encode(message: Role, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
      NotManagedData.encode(message.no, writer.uint32(82).fork()).ldelim();
    }
    if (message.builtIn !== undefined) {
      BuiltInManagedData.encode(message.builtIn, writer.uint32(90).fork()).ldelim();
    }
    if (message.identity !== undefined) {
      IdentityManagedData.encode(message.identity, writer.uint32(98).fork()).ldelim();
    }
    if (message.service !== undefined) {
      ServiceManagedData.encode(message.service, writer.uint32(106).fork()).ldelim();
    }
    for (const v of message.policies) {
      AssignedPolicy.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.tags) {
      writer.uint32(50).string(v!);
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(58).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(66).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(72).uint64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Role {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRole();
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
        case 10:
          if (tag !== 82) {
            break;
          }

          message.no = NotManagedData.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.builtIn = BuiltInManagedData.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.identity = IdentityManagedData.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.service = ServiceManagedData.decode(reader, reader.uint32());
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.policies.push(AssignedPolicy.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.tags.push(reader.string());
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 9:
          if (tag !== 72) {
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

  fromJSON(object: any): Role {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      builtIn: isSet(object.builtIn) ? BuiltInManagedData.fromJSON(object.builtIn) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
      policies: Array.isArray(object?.policies) ? object.policies.map((e: any) => AssignedPolicy.fromJSON(e)) : [],
      tags: Array.isArray(object?.tags) ? object.tags.map((e: any) => String(e)) : [],
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Role): unknown {
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
    if (message.policies?.length) {
      obj.policies = message.policies.map((e) => AssignedPolicy.toJSON(e));
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

  create<I extends Exact<DeepPartial<Role>, I>>(base?: I): Role {
    return Role.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Role>, I>>(object: I): Role {
    const message = createBaseRole();
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
    message.policies = object.policies?.map((e) => AssignedPolicy.fromPartial(e)) || [];
    message.tags = object.tags?.map((e) => e) || [];
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseCreateRoleRequest(): CreateRoleRequest {
  return { namespace: "", name: "", description: "", no: undefined, identity: undefined, service: undefined };
}

export const CreateRoleRequest = {
  encode(message: CreateRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    if (message.no !== undefined) {
      NotManagedData.encode(message.no, writer.uint32(82).fork()).ldelim();
    }
    if (message.identity !== undefined) {
      IdentityManagedData.encode(message.identity, writer.uint32(90).fork()).ldelim();
    }
    if (message.service !== undefined) {
      ServiceManagedData.encode(message.service, writer.uint32(98).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRoleRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
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
        case 10:
          if (tag !== 82) {
            break;
          }

          message.no = NotManagedData.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.identity = IdentityManagedData.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 98) {
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

  fromJSON(object: any): CreateRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      description: isSet(object.description) ? String(object.description) : "",
      no: isSet(object.no) ? NotManagedData.fromJSON(object.no) : undefined,
      identity: isSet(object.identity) ? IdentityManagedData.fromJSON(object.identity) : undefined,
      service: isSet(object.service) ? ServiceManagedData.fromJSON(object.service) : undefined,
    };
  },

  toJSON(message: CreateRoleRequest): unknown {
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
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRoleRequest>, I>>(base?: I): CreateRoleRequest {
    return CreateRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRoleRequest>, I>>(object: I): CreateRoleRequest {
    const message = createBaseCreateRoleRequest();
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
    return message;
  },
};

function createBaseCreateRoleResponse(): CreateRoleResponse {
  return { role: undefined };
}

export const CreateRoleResponse = {
  encode(message: CreateRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateRoleResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: CreateRoleResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRoleResponse>, I>>(base?: I): CreateRoleResponse {
    return CreateRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRoleResponse>, I>>(object: I): CreateRoleResponse {
    const message = createBaseCreateRoleResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseGetRoleRequest(): GetRoleRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetRoleRequest = {
  encode(message: GetRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRoleRequest();
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

  fromJSON(object: any): GetRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetRoleRequest): unknown {
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

  create<I extends Exact<DeepPartial<GetRoleRequest>, I>>(base?: I): GetRoleRequest {
    return GetRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRoleRequest>, I>>(object: I): GetRoleRequest {
    const message = createBaseGetRoleRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetRoleResponse(): GetRoleResponse {
  return { role: undefined };
}

export const GetRoleResponse = {
  encode(message: GetRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRoleResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: GetRoleResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRoleResponse>, I>>(base?: I): GetRoleResponse {
    return GetRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRoleResponse>, I>>(object: I): GetRoleResponse {
    const message = createBaseGetRoleResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseGetMultipleRolesRequest(): GetMultipleRolesRequest {
  return { roles: [] };
}

export const GetMultipleRolesRequest = {
  encode(message: GetMultipleRolesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.roles) {
      GetMultipleRolesRequest_RequestedRole.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultipleRolesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultipleRolesRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.roles.push(GetMultipleRolesRequest_RequestedRole.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetMultipleRolesRequest {
    return {
      roles: Array.isArray(object?.roles)
        ? object.roles.map((e: any) => GetMultipleRolesRequest_RequestedRole.fromJSON(e))
        : [],
    };
  },

  toJSON(message: GetMultipleRolesRequest): unknown {
    const obj: any = {};
    if (message.roles?.length) {
      obj.roles = message.roles.map((e) => GetMultipleRolesRequest_RequestedRole.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultipleRolesRequest>, I>>(base?: I): GetMultipleRolesRequest {
    return GetMultipleRolesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultipleRolesRequest>, I>>(object: I): GetMultipleRolesRequest {
    const message = createBaseGetMultipleRolesRequest();
    message.roles = object.roles?.map((e) => GetMultipleRolesRequest_RequestedRole.fromPartial(e)) || [];
    return message;
  },
};

function createBaseGetMultipleRolesRequest_RequestedRole(): GetMultipleRolesRequest_RequestedRole {
  return { namespace: "", uuid: "" };
}

export const GetMultipleRolesRequest_RequestedRole = {
  encode(message: GetMultipleRolesRequest_RequestedRole, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultipleRolesRequest_RequestedRole {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultipleRolesRequest_RequestedRole();
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

  fromJSON(object: any): GetMultipleRolesRequest_RequestedRole {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: GetMultipleRolesRequest_RequestedRole): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultipleRolesRequest_RequestedRole>, I>>(
    base?: I,
  ): GetMultipleRolesRequest_RequestedRole {
    return GetMultipleRolesRequest_RequestedRole.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultipleRolesRequest_RequestedRole>, I>>(
    object: I,
  ): GetMultipleRolesRequest_RequestedRole {
    const message = createBaseGetMultipleRolesRequest_RequestedRole();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseGetMultipleRolesResponse(): GetMultipleRolesResponse {
  return { role: undefined };
}

export const GetMultipleRolesResponse = {
  encode(message: GetMultipleRolesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetMultipleRolesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetMultipleRolesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetMultipleRolesResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: GetMultipleRolesResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetMultipleRolesResponse>, I>>(base?: I): GetMultipleRolesResponse {
    return GetMultipleRolesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetMultipleRolesResponse>, I>>(object: I): GetMultipleRolesResponse {
    const message = createBaseGetMultipleRolesResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseListRolesRequest(): ListRolesRequest {
  return { namespace: "", skip: 0, limit: 0 };
}

export const ListRolesRequest = {
  encode(message: ListRolesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRolesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRolesRequest();
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

  fromJSON(object: any): ListRolesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListRolesRequest): unknown {
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

  create<I extends Exact<DeepPartial<ListRolesRequest>, I>>(base?: I): ListRolesRequest {
    return ListRolesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRolesRequest>, I>>(object: I): ListRolesRequest {
    const message = createBaseListRolesRequest();
    message.namespace = object.namespace ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListRolesResponse(): ListRolesResponse {
  return { role: undefined };
}

export const ListRolesResponse = {
  encode(message: ListRolesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRolesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRolesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListRolesResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: ListRolesResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRolesResponse>, I>>(base?: I): ListRolesResponse {
    return ListRolesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRolesResponse>, I>>(object: I): ListRolesResponse {
    const message = createBaseListRolesResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseCountRolesRequest(): CountRolesRequest {
  return { namespace: "", useCache: false };
}

export const CountRolesRequest = {
  encode(message: CountRolesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountRolesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountRolesRequest();
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

  fromJSON(object: any): CountRolesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: CountRolesRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountRolesRequest>, I>>(base?: I): CountRolesRequest {
    return CountRolesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountRolesRequest>, I>>(object: I): CountRolesRequest {
    const message = createBaseCountRolesRequest();
    message.namespace = object.namespace ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseCountRolesResponse(): CountRolesResponse {
  return { count: 0 };
}

export const CountRolesResponse = {
  encode(message: CountRolesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountRolesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountRolesResponse();
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

  fromJSON(object: any): CountRolesResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountRolesResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountRolesResponse>, I>>(base?: I): CountRolesResponse {
    return CountRolesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountRolesResponse>, I>>(object: I): CountRolesResponse {
    const message = createBaseCountRolesResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseUpdateRoleRequest(): UpdateRoleRequest {
  return { namespace: "", uuid: "", newName: "", newDescription: "" };
}

export const UpdateRoleRequest = {
  encode(message: UpdateRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.newName !== "") {
      writer.uint32(26).string(message.newName);
    }
    if (message.newDescription !== "") {
      writer.uint32(34).string(message.newDescription);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRoleRequest();
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
        case 4:
          if (tag !== 34) {
            break;
          }

          message.newDescription = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      newName: isSet(object.newName) ? String(object.newName) : "",
      newDescription: isSet(object.newDescription) ? String(object.newDescription) : "",
    };
  },

  toJSON(message: UpdateRoleRequest): unknown {
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
    if (message.newDescription !== "") {
      obj.newDescription = message.newDescription;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRoleRequest>, I>>(base?: I): UpdateRoleRequest {
    return UpdateRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRoleRequest>, I>>(object: I): UpdateRoleRequest {
    const message = createBaseUpdateRoleRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.newName = object.newName ?? "";
    message.newDescription = object.newDescription ?? "";
    return message;
  },
};

function createBaseUpdateRoleResponse(): UpdateRoleResponse {
  return { role: undefined };
}

export const UpdateRoleResponse = {
  encode(message: UpdateRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRoleResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: UpdateRoleResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRoleResponse>, I>>(base?: I): UpdateRoleResponse {
    return UpdateRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRoleResponse>, I>>(object: I): UpdateRoleResponse {
    const message = createBaseUpdateRoleResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseDeleteRoleRequest(): DeleteRoleRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteRoleRequest = {
  encode(message: DeleteRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRoleRequest();
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

  fromJSON(object: any): DeleteRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteRoleRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRoleRequest>, I>>(base?: I): DeleteRoleRequest {
    return DeleteRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRoleRequest>, I>>(object: I): DeleteRoleRequest {
    const message = createBaseDeleteRoleRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteRoleResponse(): DeleteRoleResponse {
  return { existed: false };
}

export const DeleteRoleResponse = {
  encode(message: DeleteRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.existed === true) {
      writer.uint32(8).bool(message.existed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRoleResponse();
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

  fromJSON(object: any): DeleteRoleResponse {
    return { existed: isSet(object.existed) ? Boolean(object.existed) : false };
  },

  toJSON(message: DeleteRoleResponse): unknown {
    const obj: any = {};
    if (message.existed === true) {
      obj.existed = message.existed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRoleResponse>, I>>(base?: I): DeleteRoleResponse {
    return DeleteRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRoleResponse>, I>>(object: I): DeleteRoleResponse {
    const message = createBaseDeleteRoleResponse();
    message.existed = object.existed ?? false;
    return message;
  },
};

function createBaseGetServiceManagedRoleRequest(): GetServiceManagedRoleRequest {
  return { namespace: "", service: "", managedId: "", useCache: false };
}

export const GetServiceManagedRoleRequest = {
  encode(message: GetServiceManagedRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedRoleRequest();
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

  fromJSON(object: any): GetServiceManagedRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      service: isSet(object.service) ? String(object.service) : "",
      managedId: isSet(object.managedId) ? String(object.managedId) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetServiceManagedRoleRequest): unknown {
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

  create<I extends Exact<DeepPartial<GetServiceManagedRoleRequest>, I>>(base?: I): GetServiceManagedRoleRequest {
    return GetServiceManagedRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedRoleRequest>, I>>(object: I): GetServiceManagedRoleRequest {
    const message = createBaseGetServiceManagedRoleRequest();
    message.namespace = object.namespace ?? "";
    message.service = object.service ?? "";
    message.managedId = object.managedId ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetServiceManagedRoleResponse(): GetServiceManagedRoleResponse {
  return { role: undefined };
}

export const GetServiceManagedRoleResponse = {
  encode(message: GetServiceManagedRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetServiceManagedRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetServiceManagedRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetServiceManagedRoleResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: GetServiceManagedRoleResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetServiceManagedRoleResponse>, I>>(base?: I): GetServiceManagedRoleResponse {
    return GetServiceManagedRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetServiceManagedRoleResponse>, I>>(
    object: I,
  ): GetServiceManagedRoleResponse {
    const message = createBaseGetServiceManagedRoleResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseGetBuiltInRoleRequest(): GetBuiltInRoleRequest {
  return { namespace: "", type: 0 };
}

export const GetBuiltInRoleRequest = {
  encode(message: GetBuiltInRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBuiltInRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBuiltInRoleRequest();
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

  fromJSON(object: any): GetBuiltInRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      type: isSet(object.type) ? builtInRoleTypeFromJSON(object.type) : 0,
    };
  },

  toJSON(message: GetBuiltInRoleRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.type !== 0) {
      obj.type = builtInRoleTypeToJSON(message.type);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBuiltInRoleRequest>, I>>(base?: I): GetBuiltInRoleRequest {
    return GetBuiltInRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBuiltInRoleRequest>, I>>(object: I): GetBuiltInRoleRequest {
    const message = createBaseGetBuiltInRoleRequest();
    message.namespace = object.namespace ?? "";
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseGetBuiltInRoleResponse(): GetBuiltInRoleResponse {
  return { role: undefined };
}

export const GetBuiltInRoleResponse = {
  encode(message: GetBuiltInRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBuiltInRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBuiltInRoleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.role = Role.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetBuiltInRoleResponse {
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: GetBuiltInRoleResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetBuiltInRoleResponse>, I>>(base?: I): GetBuiltInRoleResponse {
    return GetBuiltInRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetBuiltInRoleResponse>, I>>(object: I): GetBuiltInRoleResponse {
    const message = createBaseGetBuiltInRoleResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseAddPolicyRequest(): AddPolicyRequest {
  return { roleNamespace: "", roleUUID: "", policyNamespace: "", policyUUID: "" };
}

export const AddPolicyRequest = {
  encode(message: AddPolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.roleNamespace !== "") {
      writer.uint32(10).string(message.roleNamespace);
    }
    if (message.roleUUID !== "") {
      writer.uint32(18).string(message.roleUUID);
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

          message.roleNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.roleUUID = reader.string();
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
      roleNamespace: isSet(object.roleNamespace) ? String(object.roleNamespace) : "",
      roleUUID: isSet(object.roleUUID) ? String(object.roleUUID) : "",
      policyNamespace: isSet(object.policyNamespace) ? String(object.policyNamespace) : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: AddPolicyRequest): unknown {
    const obj: any = {};
    if (message.roleNamespace !== "") {
      obj.roleNamespace = message.roleNamespace;
    }
    if (message.roleUUID !== "") {
      obj.roleUUID = message.roleUUID;
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
    message.roleNamespace = object.roleNamespace ?? "";
    message.roleUUID = object.roleUUID ?? "";
    message.policyNamespace = object.policyNamespace ?? "";
    message.policyUUID = object.policyUUID ?? "";
    return message;
  },
};

function createBaseAddPolicyResponse(): AddPolicyResponse {
  return { role: undefined };
}

export const AddPolicyResponse = {
  encode(message: AddPolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
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

          message.role = Role.decode(reader, reader.uint32());
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
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: AddPolicyResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(base?: I): AddPolicyResponse {
    return AddPolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AddPolicyResponse>, I>>(object: I): AddPolicyResponse {
    const message = createBaseAddPolicyResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseRemovePolicyRequest(): RemovePolicyRequest {
  return { roleNamespace: "", roleUUID: "", policyNamespace: "", policyUUID: "" };
}

export const RemovePolicyRequest = {
  encode(message: RemovePolicyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.roleNamespace !== "") {
      writer.uint32(10).string(message.roleNamespace);
    }
    if (message.roleUUID !== "") {
      writer.uint32(18).string(message.roleUUID);
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

          message.roleNamespace = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.roleUUID = reader.string();
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
      roleNamespace: isSet(object.roleNamespace) ? String(object.roleNamespace) : "",
      roleUUID: isSet(object.roleUUID) ? String(object.roleUUID) : "",
      policyNamespace: isSet(object.policyNamespace) ? String(object.policyNamespace) : "",
      policyUUID: isSet(object.policyUUID) ? String(object.policyUUID) : "",
    };
  },

  toJSON(message: RemovePolicyRequest): unknown {
    const obj: any = {};
    if (message.roleNamespace !== "") {
      obj.roleNamespace = message.roleNamespace;
    }
    if (message.roleUUID !== "") {
      obj.roleUUID = message.roleUUID;
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
    message.roleNamespace = object.roleNamespace ?? "";
    message.roleUUID = object.roleUUID ?? "";
    message.policyNamespace = object.policyNamespace ?? "";
    message.policyUUID = object.policyUUID ?? "";
    return message;
  },
};

function createBaseRemovePolicyResponse(): RemovePolicyResponse {
  return { role: undefined };
}

export const RemovePolicyResponse = {
  encode(message: RemovePolicyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.role !== undefined) {
      Role.encode(message.role, writer.uint32(10).fork()).ldelim();
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

          message.role = Role.decode(reader, reader.uint32());
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
    return { role: isSet(object.role) ? Role.fromJSON(object.role) : undefined };
  },

  toJSON(message: RemovePolicyResponse): unknown {
    const obj: any = {};
    if (message.role !== undefined) {
      obj.role = Role.toJSON(message.role);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(base?: I): RemovePolicyResponse {
    return RemovePolicyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemovePolicyResponse>, I>>(object: I): RemovePolicyResponse {
    const message = createBaseRemovePolicyResponse();
    message.role = (object.role !== undefined && object.role !== null) ? Role.fromPartial(object.role) : undefined;
    return message;
  },
};

function createBaseExistRoleRequest(): ExistRoleRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const ExistRoleRequest = {
  encode(message: ExistRoleRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistRoleRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistRoleRequest();
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

  fromJSON(object: any): ExistRoleRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ExistRoleRequest): unknown {
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

  create<I extends Exact<DeepPartial<ExistRoleRequest>, I>>(base?: I): ExistRoleRequest {
    return ExistRoleRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistRoleRequest>, I>>(object: I): ExistRoleRequest {
    const message = createBaseExistRoleRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseExistRoleResponse(): ExistRoleResponse {
  return { exist: false };
}

export const ExistRoleResponse = {
  encode(message: ExistRoleResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistRoleResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistRoleResponse();
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

  fromJSON(object: any): ExistRoleResponse {
    return { exist: isSet(object.exist) ? Boolean(object.exist) : false };
  },

  toJSON(message: ExistRoleResponse): unknown {
    const obj: any = {};
    if (message.exist === true) {
      obj.exist = message.exist;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistRoleResponse>, I>>(base?: I): ExistRoleResponse {
    return ExistRoleResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistRoleResponse>, I>>(object: I): ExistRoleResponse {
    const message = createBaseExistRoleResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

/** Provides API to manage IAM roles */
export interface IAMRoleService {
  /** Create new role */
  Create(request: CreateRoleRequest): Promise<CreateRoleResponse>;
  /** Get role */
  Get(request: GetRoleRequest): Promise<GetRoleResponse>;
  /** Get multiple roles in one request. */
  GetMultiple(request: GetMultipleRolesRequest): Observable<GetMultipleRolesResponse>;
  /** Get list of roles in the namespace */
  List(request: ListRolesRequest): Observable<ListRolesResponse>;
  /** Count roles in the namespace */
  Count(request: CountRolesRequest): Promise<CountRolesResponse>;
  /** Update role information */
  Update(request: UpdateRoleRequest): Promise<UpdateRoleResponse>;
  /** Delete role */
  Delete(request: DeleteRoleRequest): Promise<DeleteRoleResponse>;
  /** Get role that is managed by service */
  GetServiceManagedRole(request: GetServiceManagedRoleRequest): Promise<GetServiceManagedRoleResponse>;
  /** Get one of the builtin roles */
  GetBuiltInRole(request: GetBuiltInRoleRequest): Promise<GetBuiltInRoleResponse>;
  /** Add policy to the role. If policy was already added - does nothing. */
  AddPolicy(request: AddPolicyRequest): Promise<AddPolicyResponse>;
  /** Removes policy from the role. If policy was already removed - does nothing. */
  RemovePolicy(request: RemovePolicyRequest): Promise<RemovePolicyResponse>;
  Exist(request: ExistRoleRequest): Promise<ExistRoleResponse>;
}

export const IAMRoleServiceServiceName = "native_iam_role.IAMRoleService";
export class IAMRoleServiceClientImpl implements IAMRoleService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMRoleServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.GetMultiple = this.GetMultiple.bind(this);
    this.List = this.List.bind(this);
    this.Count = this.Count.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.GetServiceManagedRole = this.GetServiceManagedRole.bind(this);
    this.GetBuiltInRole = this.GetBuiltInRole.bind(this);
    this.AddPolicy = this.AddPolicy.bind(this);
    this.RemovePolicy = this.RemovePolicy.bind(this);
    this.Exist = this.Exist.bind(this);
  }
  Create(request: CreateRoleRequest): Promise<CreateRoleResponse> {
    const data = CreateRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateRoleResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetRoleRequest): Promise<GetRoleResponse> {
    const data = GetRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetRoleResponse.decode(_m0.Reader.create(data)));
  }

  GetMultiple(request: GetMultipleRolesRequest): Observable<GetMultipleRolesResponse> {
    const data = GetMultipleRolesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "GetMultiple", data);
    return result.pipe(map((data) => GetMultipleRolesResponse.decode(_m0.Reader.create(data))));
  }

  List(request: ListRolesRequest): Observable<ListRolesResponse> {
    const data = ListRolesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListRolesResponse.decode(_m0.Reader.create(data))));
  }

  Count(request: CountRolesRequest): Promise<CountRolesResponse> {
    const data = CountRolesRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountRolesResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateRoleRequest): Promise<UpdateRoleResponse> {
    const data = UpdateRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateRoleResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteRoleRequest): Promise<DeleteRoleResponse> {
    const data = DeleteRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteRoleResponse.decode(_m0.Reader.create(data)));
  }

  GetServiceManagedRole(request: GetServiceManagedRoleRequest): Promise<GetServiceManagedRoleResponse> {
    const data = GetServiceManagedRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetServiceManagedRole", data);
    return promise.then((data) => GetServiceManagedRoleResponse.decode(_m0.Reader.create(data)));
  }

  GetBuiltInRole(request: GetBuiltInRoleRequest): Promise<GetBuiltInRoleResponse> {
    const data = GetBuiltInRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetBuiltInRole", data);
    return promise.then((data) => GetBuiltInRoleResponse.decode(_m0.Reader.create(data)));
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

  Exist(request: ExistRoleRequest): Promise<ExistRoleResponse> {
    const data = ExistRoleRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exist", data);
    return promise.then((data) => ExistRoleResponse.decode(_m0.Reader.create(data)));
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
