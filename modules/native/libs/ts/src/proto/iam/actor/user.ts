/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_iam_actor_user";

export interface User {
  /** Namespace of the user. If namespace is empty, user is global. */
  namespace: string;
  /** Unique identifier */
  uuid: string;
  /** Login unique identifies user (like UUID) but user defined and can be changed. */
  login: string;
  /** Native_iam_identity UUID */
  identity: string;
  /** User-defined name that will be displayed instead of login */
  fullName: string;
  /** Link to the user avatar image */
  avatar: string;
  /** Email address */
  email: string;
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

export interface CreateRequest {
  /** Namespace where to create user. Empty for global user */
  namespace: string;
  /** User-defined unique identifier */
  login: string;
  /** User-defined name that will be displayed instead of login */
  fullName: string;
  /** Link to the user avatar image */
  avatar: string;
  /** Email address */
  email: string;
}

export interface CreateResponse {
  /** Created user */
  user: User | undefined;
}

export interface GetRequest {
  /** Namespace where to search for user. Empty for global user. */
  namespace: string;
  /** User unique identifier inside searched namespace */
  uuid: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface GetResponse {
  /** Founded user */
  user: User | undefined;
}

export interface GetByLoginRequest {
  /** Namespace where to search for user. Empty for global user. */
  namespace: string;
  /** Login of the user to search */
  login: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface GetByLoginResponse {
  /** User with specified login */
  user: User | undefined;
}

export interface GetByIdentityRequest {
  /** Namespace where to search for user. Empty for global user. */
  namespace: string;
  /** Search for user which has this identity uuid assigned to it */
  identity: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface GetByIdentityResponse {
  /** User wich has specified identity */
  user: User | undefined;
}

export interface UpdateRequest {
  /** Namespace of the user. Empty for global user. */
  namespace: string;
  /** Unique identifier of user, that will be updated */
  uuid: string;
  /** User-defined unique identifier */
  login: string;
  /** User-defined name that will be displayed instead of login */
  fullName: string;
  /** Link to the user avatar image */
  avatar: string;
  /** Email address */
  email: string;
}

export interface UpdateResponse {
  /** User after update */
  user: User | undefined;
}

export interface DeleteRequest {
  /** Namespace of the user. Empty for global user. */
  namespace: string;
  /** Unique identifier of user inside namespace */
  uuid: string;
}

export interface DeleteResponse {
  /** Indicates if user existed before this request or it was already deleted earlier. */
  existed: boolean;
}

export interface ListRequest {
  /** Namespace where to search for users */
  namespace: string;
}

export interface ListResponse {
  /** Founded user */
  user: User | undefined;
}

export interface CountRequest {
  /** Namespace where to count users */
  namespace: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface CountResponse {
  /** Total number of users */
  count: number;
}

export interface SearchRequest {
  /** Namespace where to search for users. Empty for global user. */
  namespace: string;
  /** String to search in the user information */
  match: string;
  /** How much values to return */
  limit: number;
}

export interface SearchResponse {
  /** User data */
  user: User | undefined;
}

function createBaseUser(): User {
  return {
    namespace: "",
    uuid: "",
    login: "",
    identity: "",
    fullName: "",
    avatar: "",
    email: "",
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const User = {
  encode(message: User, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.login !== "") {
      writer.uint32(26).string(message.login);
    }
    if (message.identity !== "") {
      writer.uint32(34).string(message.identity);
    }
    if (message.fullName !== "") {
      writer.uint32(42).string(message.fullName);
    }
    if (message.avatar !== "") {
      writer.uint32(50).string(message.avatar);
    }
    if (message.email !== "") {
      writer.uint32(58).string(message.email);
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

  decode(input: _m0.Reader | Uint8Array, length?: number): User {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUser();
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

          message.login = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.identity = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.avatar = reader.string();
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.email = reader.string();
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

  fromJSON(object: any): User {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      login: isSet(object.login) ? String(object.login) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.login !== "") {
      obj.login = message.login;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.avatar !== "") {
      obj.avatar = message.avatar;
    }
    if (message.email !== "") {
      obj.email = message.email;
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

  create<I extends Exact<DeepPartial<User>, I>>(base?: I): User {
    return User.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<User>, I>>(object: I): User {
    const message = createBaseUser();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.login = object.login ?? "";
    message.identity = object.identity ?? "";
    message.fullName = object.fullName ?? "";
    message.avatar = object.avatar ?? "";
    message.email = object.email ?? "";
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseCreateRequest(): CreateRequest {
  return { namespace: "", login: "", fullName: "", avatar: "", email: "" };
}

export const CreateRequest = {
  encode(message: CreateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.login !== "") {
      writer.uint32(18).string(message.login);
    }
    if (message.fullName !== "") {
      writer.uint32(26).string(message.fullName);
    }
    if (message.avatar !== "") {
      writer.uint32(34).string(message.avatar);
    }
    if (message.email !== "") {
      writer.uint32(42).string(message.email);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRequest();
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

          message.login = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.avatar = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.email = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      login: isSet(object.login) ? String(object.login) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: CreateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.login !== "") {
      obj.login = message.login;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.avatar !== "") {
      obj.avatar = message.avatar;
    }
    if (message.email !== "") {
      obj.email = message.email;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRequest>, I>>(base?: I): CreateRequest {
    return CreateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRequest>, I>>(object: I): CreateRequest {
    const message = createBaseCreateRequest();
    message.namespace = object.namespace ?? "";
    message.login = object.login ?? "";
    message.fullName = object.fullName ?? "";
    message.avatar = object.avatar ?? "";
    message.email = object.email ?? "";
    return message;
  },
};

function createBaseCreateResponse(): CreateResponse {
  return { user: undefined };
}

export const CreateResponse = {
  encode(message: CreateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: CreateResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateResponse>, I>>(base?: I): CreateResponse {
    return CreateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateResponse>, I>>(object: I): CreateResponse {
    const message = createBaseCreateResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetRequest = {
  encode(message: GetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  fromJSON(object: any): GetRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetRequest): unknown {
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

  create<I extends Exact<DeepPartial<GetRequest>, I>>(base?: I): GetRequest {
    return GetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(object: I): GetRequest {
    const message = createBaseGetRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { user: undefined };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
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

          message.user = User.decode(reader, reader.uint32());
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
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseGetByLoginRequest(): GetByLoginRequest {
  return { namespace: "", login: "", useCache: false };
}

export const GetByLoginRequest = {
  encode(message: GetByLoginRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.login !== "") {
      writer.uint32(18).string(message.login);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByLoginRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByLoginRequest();
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

          message.login = reader.string();
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

  fromJSON(object: any): GetByLoginRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      login: isSet(object.login) ? String(object.login) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetByLoginRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.login !== "") {
      obj.login = message.login;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetByLoginRequest>, I>>(base?: I): GetByLoginRequest {
    return GetByLoginRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetByLoginRequest>, I>>(object: I): GetByLoginRequest {
    const message = createBaseGetByLoginRequest();
    message.namespace = object.namespace ?? "";
    message.login = object.login ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetByLoginResponse(): GetByLoginResponse {
  return { user: undefined };
}

export const GetByLoginResponse = {
  encode(message: GetByLoginResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByLoginResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByLoginResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetByLoginResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: GetByLoginResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetByLoginResponse>, I>>(base?: I): GetByLoginResponse {
    return GetByLoginResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetByLoginResponse>, I>>(object: I): GetByLoginResponse {
    const message = createBaseGetByLoginResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseGetByIdentityRequest(): GetByIdentityRequest {
  return { namespace: "", identity: "", useCache: false };
}

export const GetByIdentityRequest = {
  encode(message: GetByIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByIdentityRequest();
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

          message.identity = reader.string();
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

  fromJSON(object: any): GetByIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetByIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetByIdentityRequest>, I>>(base?: I): GetByIdentityRequest {
    return GetByIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetByIdentityRequest>, I>>(object: I): GetByIdentityRequest {
    const message = createBaseGetByIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetByIdentityResponse(): GetByIdentityResponse {
  return { user: undefined };
}

export const GetByIdentityResponse = {
  encode(message: GetByIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetByIdentityResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: GetByIdentityResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetByIdentityResponse>, I>>(base?: I): GetByIdentityResponse {
    return GetByIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetByIdentityResponse>, I>>(object: I): GetByIdentityResponse {
    const message = createBaseGetByIdentityResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseUpdateRequest(): UpdateRequest {
  return { namespace: "", uuid: "", login: "", fullName: "", avatar: "", email: "" };
}

export const UpdateRequest = {
  encode(message: UpdateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.login !== "") {
      writer.uint32(26).string(message.login);
    }
    if (message.fullName !== "") {
      writer.uint32(34).string(message.fullName);
    }
    if (message.avatar !== "") {
      writer.uint32(42).string(message.avatar);
    }
    if (message.email !== "") {
      writer.uint32(50).string(message.email);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRequest();
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

          message.login = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.fullName = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.avatar = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.email = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      login: isSet(object.login) ? String(object.login) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: UpdateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.login !== "") {
      obj.login = message.login;
    }
    if (message.fullName !== "") {
      obj.fullName = message.fullName;
    }
    if (message.avatar !== "") {
      obj.avatar = message.avatar;
    }
    if (message.email !== "") {
      obj.email = message.email;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRequest>, I>>(base?: I): UpdateRequest {
    return UpdateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRequest>, I>>(object: I): UpdateRequest {
    const message = createBaseUpdateRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.login = object.login ?? "";
    message.fullName = object.fullName ?? "";
    message.avatar = object.avatar ?? "";
    message.email = object.email ?? "";
    return message;
  },
};

function createBaseUpdateResponse(): UpdateResponse {
  return { user: undefined };
}

export const UpdateResponse = {
  encode(message: UpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: UpdateResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateResponse>, I>>(base?: I): UpdateResponse {
    return UpdateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateResponse>, I>>(object: I): UpdateResponse {
    const message = createBaseUpdateResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseDeleteRequest(): DeleteRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteRequest = {
  encode(message: DeleteRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRequest();
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

  fromJSON(object: any): DeleteRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRequest>, I>>(base?: I): DeleteRequest {
    return DeleteRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRequest>, I>>(object: I): DeleteRequest {
    const message = createBaseDeleteRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteResponse(): DeleteResponse {
  return { existed: false };
}

export const DeleteResponse = {
  encode(message: DeleteResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.existed === true) {
      writer.uint32(8).bool(message.existed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteResponse();
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

  fromJSON(object: any): DeleteResponse {
    return { existed: isSet(object.existed) ? Boolean(object.existed) : false };
  },

  toJSON(message: DeleteResponse): unknown {
    const obj: any = {};
    if (message.existed === true) {
      obj.existed = message.existed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteResponse>, I>>(base?: I): DeleteResponse {
    return DeleteResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteResponse>, I>>(object: I): DeleteResponse {
    const message = createBaseDeleteResponse();
    message.existed = object.existed ?? false;
    return message;
  },
};

function createBaseListRequest(): ListRequest {
  return { namespace: "" };
}

export const ListRequest = {
  encode(message: ListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRequest();
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

  fromJSON(object: any): ListRequest {
    return { namespace: isSet(object.namespace) ? String(object.namespace) : "" };
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRequest>, I>>(base?: I): ListRequest {
    return ListRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRequest>, I>>(object: I): ListRequest {
    const message = createBaseListRequest();
    message.namespace = object.namespace ?? "";
    return message;
  },
};

function createBaseListResponse(): ListResponse {
  return { user: undefined };
}

export const ListResponse = {
  encode(message: ListResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListResponse>, I>>(base?: I): ListResponse {
    return ListResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListResponse>, I>>(object: I): ListResponse {
    const message = createBaseListResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

function createBaseCountRequest(): CountRequest {
  return { namespace: "", useCache: false };
}

export const CountRequest = {
  encode(message: CountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountRequest();
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

  fromJSON(object: any): CountRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: CountRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountRequest>, I>>(base?: I): CountRequest {
    return CountRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountRequest>, I>>(object: I): CountRequest {
    const message = createBaseCountRequest();
    message.namespace = object.namespace ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseCountResponse(): CountResponse {
  return { count: 0 };
}

export const CountResponse = {
  encode(message: CountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountResponse();
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

  fromJSON(object: any): CountResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountResponse>, I>>(base?: I): CountResponse {
    return CountResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountResponse>, I>>(object: I): CountResponse {
    const message = createBaseCountResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseSearchRequest(): SearchRequest {
  return { namespace: "", match: "", limit: 0 };
}

export const SearchRequest = {
  encode(message: SearchRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.match !== "") {
      writer.uint32(18).string(message.match);
    }
    if (message.limit !== 0) {
      writer.uint32(24).uint64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SearchRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSearchRequest();
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

          message.match = reader.string();
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

  fromJSON(object: any): SearchRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      match: isSet(object.match) ? String(object.match) : "",
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: SearchRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.match !== "") {
      obj.match = message.match;
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SearchRequest>, I>>(base?: I): SearchRequest {
    return SearchRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SearchRequest>, I>>(object: I): SearchRequest {
    const message = createBaseSearchRequest();
    message.namespace = object.namespace ?? "";
    message.match = object.match ?? "";
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseSearchResponse(): SearchResponse {
  return { user: undefined };
}

export const SearchResponse = {
  encode(message: SearchResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SearchResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSearchResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.user = User.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): SearchResponse {
    return { user: isSet(object.user) ? User.fromJSON(object.user) : undefined };
  },

  toJSON(message: SearchResponse): unknown {
    const obj: any = {};
    if (message.user !== undefined) {
      obj.user = User.toJSON(message.user);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<SearchResponse>, I>>(base?: I): SearchResponse {
    return SearchResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SearchResponse>, I>>(object: I): SearchResponse {
    const message = createBaseSearchResponse();
    message.user = (object.user !== undefined && object.user !== null) ? User.fromPartial(object.user) : undefined;
    return message;
  },
};

export interface ActorUserService {
  /** Create new user and assign identity to it. */
  Create(request: CreateRequest): Promise<CreateResponse>;
  /** Get user by its unique identifier */
  Get(request: GetRequest): Promise<GetResponse>;
  /** Get user by its login */
  GetByLogin(request: GetByLoginRequest): Promise<GetByLoginResponse>;
  /** Get user by the identity uuid that was assigned to it */
  GetByIdentity(request: GetByIdentityRequest): Promise<GetByIdentityResponse>;
  /** Update user information */
  Update(request: UpdateRequest): Promise<UpdateResponse>;
  /** Delete user */
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  /** Get all users in the namespace. */
  List(request: ListRequest): Observable<ListResponse>;
  /** Get total number of users inside namespace */
  Count(request: CountRequest): Promise<CountResponse>;
  /**
   * Searches for user using some "matching" string. Much faster than find operation. Searches for matches in login/fullName/email.
   * Matches may be not ideal and its not possible to predict how much users matched provided string.
   */
  Search(request: SearchRequest): Observable<SearchResponse>;
}

export const ActorUserServiceServiceName = "native_iam_actor_user.ActorUserService";
export class ActorUserServiceClientImpl implements ActorUserService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || ActorUserServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.GetByLogin = this.GetByLogin.bind(this);
    this.GetByIdentity = this.GetByIdentity.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.List = this.List.bind(this);
    this.Count = this.Count.bind(this);
    this.Search = this.Search.bind(this);
  }
  Create(request: CreateRequest): Promise<CreateResponse> {
    const data = CreateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetResponse.decode(_m0.Reader.create(data)));
  }

  GetByLogin(request: GetByLoginRequest): Promise<GetByLoginResponse> {
    const data = GetByLoginRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetByLogin", data);
    return promise.then((data) => GetByLoginResponse.decode(_m0.Reader.create(data)));
  }

  GetByIdentity(request: GetByIdentityRequest): Promise<GetByIdentityResponse> {
    const data = GetByIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetByIdentity", data);
    return promise.then((data) => GetByIdentityResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateRequest): Promise<UpdateResponse> {
    const data = UpdateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteResponse.decode(_m0.Reader.create(data)));
  }

  List(request: ListRequest): Observable<ListResponse> {
    const data = ListRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListResponse.decode(_m0.Reader.create(data))));
  }

  Count(request: CountRequest): Promise<CountResponse> {
    const data = CountRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountResponse.decode(_m0.Reader.create(data)));
  }

  Search(request: SearchRequest): Observable<SearchResponse> {
    const data = SearchRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "Search", data);
    return result.pipe(map((data) => SearchResponse.decode(_m0.Reader.create(data))));
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
