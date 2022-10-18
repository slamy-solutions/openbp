/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_actor_user";

export interface User {
  /** Unique identifier */
  uuid: string;
  /** Login unique identifies user (like UUID) but user defined and can be changed. */
  login: string;
  /** Nnative_iam_identity UUID */
  identity: string;
  /** User-defined name that will be displayed instead of login */
  fullName: string;
  /** Link to the user avatar image */
  avatar: string;
  /** Email address */
  email: string;
}

export interface CreateRequest {
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
  /** User unique identifier to get */
  uuid: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface GetResponse {
  /** User with specified UUID */
  user: User | undefined;
}

export interface GetByLoginRequest {
  /** Search for user that has this login */
  login: string;
  /** Use cache for this request or not. Cache has a very small chance to be invalid. Invalid cache deletes after small period of time (60 seconds by default) */
  useCache: boolean;
}

export interface GetByLoginResponse {
  /** User with specified login */
  user: User | undefined;
}

export interface GetByIdentityRequest {
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
  /** Unique identifier of user to delete */
  uuid: string;
}

export interface DeleteResponse {}

function createBaseUser(): User {
  return {
    uuid: "",
    login: "",
    identity: "",
    fullName: "",
    avatar: "",
    email: "",
  };
}

export const User = {
  encode(message: User, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
    }
    if (message.login !== "") {
      writer.uint32(18).string(message.login);
    }
    if (message.identity !== "") {
      writer.uint32(26).string(message.identity);
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

  decode(input: _m0.Reader | Uint8Array, length?: number): User {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUser();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.string();
          break;
        case 2:
          message.login = reader.string();
          break;
        case 3:
          message.identity = reader.string();
          break;
        case 4:
          message.fullName = reader.string();
          break;
        case 5:
          message.avatar = reader.string();
          break;
        case 6:
          message.email = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): User {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      login: isSet(object.login) ? String(object.login) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: User): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.login !== undefined && (obj.login = message.login);
    message.identity !== undefined && (obj.identity = message.identity);
    message.fullName !== undefined && (obj.fullName = message.fullName);
    message.avatar !== undefined && (obj.avatar = message.avatar);
    message.email !== undefined && (obj.email = message.email);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<User>, I>>(object: I): User {
    const message = createBaseUser();
    message.uuid = object.uuid ?? "";
    message.login = object.login ?? "";
    message.identity = object.identity ?? "";
    message.fullName = object.fullName ?? "";
    message.avatar = object.avatar ?? "";
    message.email = object.email ?? "";
    return message;
  },
};

function createBaseCreateRequest(): CreateRequest {
  return { login: "", fullName: "", avatar: "", email: "" };
}

export const CreateRequest = {
  encode(
    message: CreateRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.login !== "") {
      writer.uint32(10).string(message.login);
    }
    if (message.fullName !== "") {
      writer.uint32(18).string(message.fullName);
    }
    if (message.avatar !== "") {
      writer.uint32(26).string(message.avatar);
    }
    if (message.email !== "") {
      writer.uint32(34).string(message.email);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.login = reader.string();
          break;
        case 2:
          message.fullName = reader.string();
          break;
        case 3:
          message.avatar = reader.string();
          break;
        case 4:
          message.email = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateRequest {
    return {
      login: isSet(object.login) ? String(object.login) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: CreateRequest): unknown {
    const obj: any = {};
    message.login !== undefined && (obj.login = message.login);
    message.fullName !== undefined && (obj.fullName = message.fullName);
    message.avatar !== undefined && (obj.avatar = message.avatar);
    message.email !== undefined && (obj.email = message.email);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateRequest>, I>>(
    object: I
  ): CreateRequest {
    const message = createBaseCreateRequest();
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
  encode(
    message: CreateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = User.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: CreateResponse): unknown {
    const obj: any = {};
    message.user !== undefined &&
      (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateResponse>, I>>(
    object: I
  ): CreateResponse {
    const message = createBaseCreateResponse();
    message.user =
      object.user !== undefined && object.user !== null
        ? User.fromPartial(object.user)
        : undefined;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { uuid: "", useCache: false };
}

export const GetRequest = {
  encode(
    message: GetRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.string();
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

  fromJSON(object: any): GetRequest {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetRequest): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(
    object: I
  ): GetRequest {
    const message = createBaseGetRequest();
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { user: undefined };
}

export const GetResponse = {
  encode(
    message: GetResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 2:
          message.user = User.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    message.user !== undefined &&
      (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(
    object: I
  ): GetResponse {
    const message = createBaseGetResponse();
    message.user =
      object.user !== undefined && object.user !== null
        ? User.fromPartial(object.user)
        : undefined;
    return message;
  },
};

function createBaseGetByLoginRequest(): GetByLoginRequest {
  return { login: "", useCache: false };
}

export const GetByLoginRequest = {
  encode(
    message: GetByLoginRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.login !== "") {
      writer.uint32(10).string(message.login);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByLoginRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByLoginRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.login = reader.string();
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

  fromJSON(object: any): GetByLoginRequest {
    return {
      login: isSet(object.login) ? String(object.login) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetByLoginRequest): unknown {
    const obj: any = {};
    message.login !== undefined && (obj.login = message.login);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByLoginRequest>, I>>(
    object: I
  ): GetByLoginRequest {
    const message = createBaseGetByLoginRequest();
    message.login = object.login ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetByLoginResponse(): GetByLoginResponse {
  return { user: undefined };
}

export const GetByLoginResponse = {
  encode(
    message: GetByLoginResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByLoginResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByLoginResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = User.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetByLoginResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: GetByLoginResponse): unknown {
    const obj: any = {};
    message.user !== undefined &&
      (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByLoginResponse>, I>>(
    object: I
  ): GetByLoginResponse {
    const message = createBaseGetByLoginResponse();
    message.user =
      object.user !== undefined && object.user !== null
        ? User.fromPartial(object.user)
        : undefined;
    return message;
  },
};

function createBaseGetByIdentityRequest(): GetByIdentityRequest {
  return { identity: "", useCache: false };
}

export const GetByIdentityRequest = {
  encode(
    message: GetByIdentityRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.identity !== "") {
      writer.uint32(10).string(message.identity);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): GetByIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByIdentityRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.identity = reader.string();
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

  fromJSON(object: any): GetByIdentityRequest {
    return {
      identity: isSet(object.identity) ? String(object.identity) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetByIdentityRequest): unknown {
    const obj: any = {};
    message.identity !== undefined && (obj.identity = message.identity);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByIdentityRequest>, I>>(
    object: I
  ): GetByIdentityRequest {
    const message = createBaseGetByIdentityRequest();
    message.identity = object.identity ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetByIdentityResponse(): GetByIdentityResponse {
  return { user: undefined };
}

export const GetByIdentityResponse = {
  encode(
    message: GetByIdentityResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): GetByIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = User.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetByIdentityResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: GetByIdentityResponse): unknown {
    const obj: any = {};
    message.user !== undefined &&
      (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByIdentityResponse>, I>>(
    object: I
  ): GetByIdentityResponse {
    const message = createBaseGetByIdentityResponse();
    message.user =
      object.user !== undefined && object.user !== null
        ? User.fromPartial(object.user)
        : undefined;
    return message;
  },
};

function createBaseUpdateRequest(): UpdateRequest {
  return { uuid: "", login: "", fullName: "", avatar: "", email: "" };
}

export const UpdateRequest = {
  encode(
    message: UpdateRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
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

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.string();
          break;
        case 2:
          message.login = reader.string();
          break;
        case 3:
          message.fullName = reader.string();
          break;
        case 4:
          message.avatar = reader.string();
          break;
        case 5:
          message.email = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateRequest {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      login: isSet(object.login) ? String(object.login) : "",
      fullName: isSet(object.fullName) ? String(object.fullName) : "",
      avatar: isSet(object.avatar) ? String(object.avatar) : "",
      email: isSet(object.email) ? String(object.email) : "",
    };
  },

  toJSON(message: UpdateRequest): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.login !== undefined && (obj.login = message.login);
    message.fullName !== undefined && (obj.fullName = message.fullName);
    message.avatar !== undefined && (obj.avatar = message.avatar);
    message.email !== undefined && (obj.email = message.email);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdateRequest>, I>>(
    object: I
  ): UpdateRequest {
    const message = createBaseUpdateRequest();
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
  encode(
    message: UpdateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.user !== undefined) {
      User.encode(message.user, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.user = User.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UpdateResponse {
    return {
      user: isSet(object.user) ? User.fromJSON(object.user) : undefined,
    };
  },

  toJSON(message: UpdateResponse): unknown {
    const obj: any = {};
    message.user !== undefined &&
      (obj.user = message.user ? User.toJSON(message.user) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<UpdateResponse>, I>>(
    object: I
  ): UpdateResponse {
    const message = createBaseUpdateResponse();
    message.user =
      object.user !== undefined && object.user !== null
        ? User.fromPartial(object.user)
        : undefined;
    return message;
  },
};

function createBaseDeleteRequest(): DeleteRequest {
  return { uuid: "" };
}

export const DeleteRequest = {
  encode(
    message: DeleteRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid !== "") {
      writer.uint32(10).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteRequest();
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

  fromJSON(object: any): DeleteRequest {
    return {
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteRequest): unknown {
    const obj: any = {};
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteRequest>, I>>(
    object: I
  ): DeleteRequest {
    const message = createBaseDeleteRequest();
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteResponse(): DeleteResponse {
  return {};
}

export const DeleteResponse = {
  encode(
    _: DeleteResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteResponse();
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

  fromJSON(_: any): DeleteResponse {
    return {};
  },

  toJSON(_: DeleteResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteResponse>, I>>(
    _: I
  ): DeleteResponse {
    const message = createBaseDeleteResponse();
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
}

export class ActorUserServiceClientImpl implements ActorUserService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.GetByLogin = this.GetByLogin.bind(this);
    this.GetByIdentity = this.GetByIdentity.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
  }
  Create(request: CreateRequest): Promise<CreateResponse> {
    const data = CreateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "Create",
      data
    );
    return promise.then((data) => CreateResponse.decode(new _m0.Reader(data)));
  }

  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "Get",
      data
    );
    return promise.then((data) => GetResponse.decode(new _m0.Reader(data)));
  }

  GetByLogin(request: GetByLoginRequest): Promise<GetByLoginResponse> {
    const data = GetByLoginRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "GetByLogin",
      data
    );
    return promise.then((data) =>
      GetByLoginResponse.decode(new _m0.Reader(data))
    );
  }

  GetByIdentity(request: GetByIdentityRequest): Promise<GetByIdentityResponse> {
    const data = GetByIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "GetByIdentity",
      data
    );
    return promise.then((data) =>
      GetByIdentityResponse.decode(new _m0.Reader(data))
    );
  }

  Update(request: UpdateRequest): Promise<UpdateResponse> {
    const data = UpdateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "Update",
      data
    );
    return promise.then((data) => UpdateResponse.decode(new _m0.Reader(data)));
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_actor_user.ActorUserService",
      "Delete",
      data
    );
    return promise.then((data) => DeleteResponse.decode(new _m0.Reader(data)));
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
