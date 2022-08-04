/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { Timestamp } from "./google/protobuf/timestamp";
import { map } from "rxjs/operators";

export const protobufPackage = "native_iam_token";

/** Scope defines what can be accessed by token and what actions can be performed on accessible resources. Scope is bounded the namespace where resources are located. */
export interface Scope {
  /** Namespace to which scope is bounded. Empty string if scope is not bounded to any namespace (is global) */
  namespace: string;
  /** Resources that can be accessed using this token */
  resources: string[];
  /** Actions that can be performed on accessible resources */
  actions: string[];
}

export interface TokenData {
  /** Namespace where token and identity are located. Epmty for global token (without namespace) */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
  /** Token identity unique identifier inside namespace */
  identity: string;
  /** Identifies if token was manually disabled. Disabled token always fails on authorization and can not be reenabled */
  disabled: boolean;
  /** Datetime after with token will not be valid and will fail on Refresh and Authorize attempts */
  expiresAt: Date | undefined;
  /** List of token scopes. Describes what actions can token perform on what resources */
  scopes: Scope[];
  /** Datetime when token was created */
  createdAt: Date | undefined;
  /** Arbitrary metadata added on token creation. For example MAC/IP/information of the actor/application/browser/machine that created this token. The exact format of metadata is not defined, but JSON is suggested. */
  creationMetadata: string;
}

export interface CreateRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Identity UUID of the token */
  identity: string;
  /** Scopes that will be applied to the token */
  scopes: Scope[];
  /** Arbitrary metadata. For example MAC/IP/information of the actor/application/browser/machine that created this token. The exact format of metadata is not defined, but JSON is suggested. */
  metadata: string;
}

export interface CreateResponse {
  /** Actual token formated to the string. */
  token: string;
  /** Refreshtoken is used to update token */
  refreshToken: string;
  /** Token data */
  tokenData: TokenData | undefined;
}

export interface GetByUUIDRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
  /** Use cache for faster authorization. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) */
  useCache: boolean;
}

export interface GetByUUIDResponse {
  /** Actual token data */
  tokenData: TokenData | undefined;
}

export interface DisableByUUIDRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
}

export interface DisableByUUIDResponse {}

export interface AuthorizeRequest {
  /** Token to authorize */
  token: string;
  /** Use cache for faster authorization. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) */
  useCache: boolean;
}

export interface AuthorizeResponse {
  status: AuthorizeResponse_Status;
  /** Token data. Null if status is not OK */
  tokenData: TokenData | undefined;
}

export enum AuthorizeResponse_Status {
  /** OK - Token is valid */
  OK = 0,
  /** INVALID - Token has bad format or invalid signature */
  INVALID = 1,
  /** NOT_FOUND - Token not found */
  NOT_FOUND = 2,
  /** DISABLED - Token was manually disabled */
  DISABLED = 3,
  /** EXPIRED - Token expired and is not valid */
  EXPIRED = 4,
  UNRECOGNIZED = -1,
}

export function authorizeResponse_StatusFromJSON(
  object: any
): AuthorizeResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return AuthorizeResponse_Status.OK;
    case 1:
    case "INVALID":
      return AuthorizeResponse_Status.INVALID;
    case 2:
    case "NOT_FOUND":
      return AuthorizeResponse_Status.NOT_FOUND;
    case 3:
    case "DISABLED":
      return AuthorizeResponse_Status.DISABLED;
    case 4:
    case "EXPIRED":
      return AuthorizeResponse_Status.EXPIRED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return AuthorizeResponse_Status.UNRECOGNIZED;
  }
}

export function authorizeResponse_StatusToJSON(
  object: AuthorizeResponse_Status
): string {
  switch (object) {
    case AuthorizeResponse_Status.OK:
      return "OK";
    case AuthorizeResponse_Status.INVALID:
      return "INVALID";
    case AuthorizeResponse_Status.NOT_FOUND:
      return "NOT_FOUND";
    case AuthorizeResponse_Status.DISABLED:
      return "DISABLED";
    case AuthorizeResponse_Status.EXPIRED:
      return "EXPIRED";
    case AuthorizeResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface RefreshRequest {
  /** Refresh token, based on which, new token will be returned */
  refreshToken: string;
}

export interface RefreshResponse {
  status: RefreshResponse_Status;
  /** New token if status is OK. Null otherwise */
  token: string;
  /** New token data if status is OK. Null otherwise */
  tokenData: TokenData | undefined;
}

export enum RefreshResponse_Status {
  /** OK - Token is valid */
  OK = 0,
  /** INVALID - Token has bad format or invalid signature */
  INVALID = 1,
  /** NOT_FOUND - Token not found */
  NOT_FOUND = 2,
  /** DISABLED - Token was manually disabled */
  DISABLED = 3,
  /** EXPIRED - Token expired and is not valid */
  EXPIRED = 4,
  /** NOT_REFRESH_TOKEN - This token is valid but this is not refresh token */
  NOT_REFRESH_TOKEN = 5,
  UNRECOGNIZED = -1,
}

export function refreshResponse_StatusFromJSON(
  object: any
): RefreshResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return RefreshResponse_Status.OK;
    case 1:
    case "INVALID":
      return RefreshResponse_Status.INVALID;
    case 2:
    case "NOT_FOUND":
      return RefreshResponse_Status.NOT_FOUND;
    case 3:
    case "DISABLED":
      return RefreshResponse_Status.DISABLED;
    case 4:
    case "EXPIRED":
      return RefreshResponse_Status.EXPIRED;
    case 5:
    case "NOT_REFRESH_TOKEN":
      return RefreshResponse_Status.NOT_REFRESH_TOKEN;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RefreshResponse_Status.UNRECOGNIZED;
  }
}

export function refreshResponse_StatusToJSON(
  object: RefreshResponse_Status
): string {
  switch (object) {
    case RefreshResponse_Status.OK:
      return "OK";
    case RefreshResponse_Status.INVALID:
      return "INVALID";
    case RefreshResponse_Status.NOT_FOUND:
      return "NOT_FOUND";
    case RefreshResponse_Status.DISABLED:
      return "DISABLED";
    case RefreshResponse_Status.EXPIRED:
      return "EXPIRED";
    case RefreshResponse_Status.NOT_REFRESH_TOKEN:
      return "NOT_REFRESH_TOKEN";
    case RefreshResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface TokensForIdentityRequest {
  /** Namespace where token is located. Empty for global token */
  namespace: string;
  /** Identity unique identifier inside namespace */
  identity: string;
  /** Perform results filtering on "active" property of the token */
  activeFilter: TokensForIdentityRequest_ActiveFilter;
  /** Skip number of results before returning actual tokens. Set to 0 in order not to skip */
  skip: number;
  /** Limit number of returned results. Set to 0 in order to remove limit and return all possible results up to the end. */
  limit: number;
}

export enum TokensForIdentityRequest_ActiveFilter {
  /** ALL - Get all token */
  ALL = 0,
  /** ONLY_ACTIVE - Only get tokens that wasnt disabled and not expired */
  ONLY_ACTIVE = 1,
  /** ONLY_NOT_ACTIVE - Only get tokens that are disabled or expired */
  ONLY_NOT_ACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function tokensForIdentityRequest_ActiveFilterFromJSON(
  object: any
): TokensForIdentityRequest_ActiveFilter {
  switch (object) {
    case 0:
    case "ALL":
      return TokensForIdentityRequest_ActiveFilter.ALL;
    case 1:
    case "ONLY_ACTIVE":
      return TokensForIdentityRequest_ActiveFilter.ONLY_ACTIVE;
    case 2:
    case "ONLY_NOT_ACTIVE":
      return TokensForIdentityRequest_ActiveFilter.ONLY_NOT_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return TokensForIdentityRequest_ActiveFilter.UNRECOGNIZED;
  }
}

export function tokensForIdentityRequest_ActiveFilterToJSON(
  object: TokensForIdentityRequest_ActiveFilter
): string {
  switch (object) {
    case TokensForIdentityRequest_ActiveFilter.ALL:
      return "ALL";
    case TokensForIdentityRequest_ActiveFilter.ONLY_ACTIVE:
      return "ONLY_ACTIVE";
    case TokensForIdentityRequest_ActiveFilter.ONLY_NOT_ACTIVE:
      return "ONLY_NOT_ACTIVE";
    case TokensForIdentityRequest_ActiveFilter.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface TokensForIdentityResponse {
  /** Actual token data */
  tokenData: TokenData | undefined;
}

function createBaseScope(): Scope {
  return { namespace: "", resources: [], actions: [] };
}

export const Scope = {
  encode(message: Scope, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    for (const v of message.resources) {
      writer.uint32(18).string(v!);
    }
    for (const v of message.actions) {
      writer.uint32(26).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Scope {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseScope();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.resources.push(reader.string());
          break;
        case 3:
          message.actions.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Scope {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      resources: Array.isArray(object?.resources)
        ? object.resources.map((e: any) => String(e))
        : [],
      actions: Array.isArray(object?.actions)
        ? object.actions.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: Scope): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    if (message.resources) {
      obj.resources = message.resources.map((e) => e);
    } else {
      obj.resources = [];
    }
    if (message.actions) {
      obj.actions = message.actions.map((e) => e);
    } else {
      obj.actions = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Scope>, I>>(object: I): Scope {
    const message = createBaseScope();
    message.namespace = object.namespace ?? "";
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseTokenData(): TokenData {
  return {
    namespace: "",
    uuid: "",
    identity: "",
    disabled: false,
    expiresAt: undefined,
    scopes: [],
    createdAt: undefined,
    creationMetadata: "",
  };
}

export const TokenData = {
  encode(
    message: TokenData,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.identity !== "") {
      writer.uint32(26).string(message.identity);
    }
    if (message.disabled === true) {
      writer.uint32(32).bool(message.disabled);
    }
    if (message.expiresAt !== undefined) {
      Timestamp.encode(
        toTimestamp(message.expiresAt),
        writer.uint32(42).fork()
      ).ldelim();
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.createdAt !== undefined) {
      Timestamp.encode(
        toTimestamp(message.createdAt),
        writer.uint32(58).fork()
      ).ldelim();
    }
    if (message.creationMetadata !== "") {
      writer.uint32(66).string(message.creationMetadata);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TokenData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTokenData();
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
          message.identity = reader.string();
          break;
        case 4:
          message.disabled = reader.bool();
          break;
        case 5:
          message.expiresAt = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 6:
          message.scopes.push(Scope.decode(reader, reader.uint32()));
          break;
        case 7:
          message.createdAt = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 8:
          message.creationMetadata = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TokenData {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      disabled: isSet(object.disabled) ? Boolean(object.disabled) : false,
      expiresAt: isSet(object.expiresAt)
        ? fromJsonTimestamp(object.expiresAt)
        : undefined,
      scopes: Array.isArray(object?.scopes)
        ? object.scopes.map((e: any) => Scope.fromJSON(e))
        : [],
      createdAt: isSet(object.createdAt)
        ? fromJsonTimestamp(object.createdAt)
        : undefined,
      creationMetadata: isSet(object.creationMetadata)
        ? String(object.creationMetadata)
        : "",
    };
  },

  toJSON(message: TokenData): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.identity !== undefined && (obj.identity = message.identity);
    message.disabled !== undefined && (obj.disabled = message.disabled);
    message.expiresAt !== undefined &&
      (obj.expiresAt = message.expiresAt.toISOString());
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => (e ? Scope.toJSON(e) : undefined));
    } else {
      obj.scopes = [];
    }
    message.createdAt !== undefined &&
      (obj.createdAt = message.createdAt.toISOString());
    message.creationMetadata !== undefined &&
      (obj.creationMetadata = message.creationMetadata);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TokenData>, I>>(
    object: I
  ): TokenData {
    const message = createBaseTokenData();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.identity = object.identity ?? "";
    message.disabled = object.disabled ?? false;
    message.expiresAt = object.expiresAt ?? undefined;
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    message.createdAt = object.createdAt ?? undefined;
    message.creationMetadata = object.creationMetadata ?? "";
    return message;
  },
};

function createBaseCreateRequest(): CreateRequest {
  return { namespace: "", identity: "", scopes: [], metadata: "" };
}

export const CreateRequest = {
  encode(
    message: CreateRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    if (message.metadata !== "") {
      writer.uint32(66).string(message.metadata);
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
          message.namespace = reader.string();
          break;
        case 2:
          message.identity = reader.string();
          break;
        case 3:
          message.scopes.push(Scope.decode(reader, reader.uint32()));
          break;
        case 8:
          message.metadata = reader.string();
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
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      scopes: Array.isArray(object?.scopes)
        ? object.scopes.map((e: any) => Scope.fromJSON(e))
        : [],
      metadata: isSet(object.metadata) ? String(object.metadata) : "",
    };
  },

  toJSON(message: CreateRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => (e ? Scope.toJSON(e) : undefined));
    } else {
      obj.scopes = [];
    }
    message.metadata !== undefined && (obj.metadata = message.metadata);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateRequest>, I>>(
    object: I
  ): CreateRequest {
    const message = createBaseCreateRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    message.metadata = object.metadata ?? "";
    return message;
  },
};

function createBaseCreateResponse(): CreateResponse {
  return { token: "", refreshToken: "", tokenData: undefined };
}

export const CreateResponse = {
  encode(
    message: CreateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.refreshToken !== "") {
      writer.uint32(18).string(message.refreshToken);
    }
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(26).fork()).ldelim();
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
          message.token = reader.string();
          break;
        case 2:
          message.refreshToken = reader.string();
          break;
        case 3:
          message.tokenData = TokenData.decode(reader, reader.uint32());
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
      token: isSet(object.token) ? String(object.token) : "",
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
      tokenData: isSet(object.tokenData)
        ? TokenData.fromJSON(object.tokenData)
        : undefined,
    };
  },

  toJSON(message: CreateResponse): unknown {
    const obj: any = {};
    message.token !== undefined && (obj.token = message.token);
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    message.tokenData !== undefined &&
      (obj.tokenData = message.tokenData
        ? TokenData.toJSON(message.tokenData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateResponse>, I>>(
    object: I
  ): CreateResponse {
    const message = createBaseCreateResponse();
    message.token = object.token ?? "";
    message.refreshToken = object.refreshToken ?? "";
    message.tokenData =
      object.tokenData !== undefined && object.tokenData !== null
        ? TokenData.fromPartial(object.tokenData)
        : undefined;
    return message;
  },
};

function createBaseGetByUUIDRequest(): GetByUUIDRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const GetByUUIDRequest = {
  encode(
    message: GetByUUIDRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByUUIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByUUIDRequest();
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

  fromJSON(object: any): GetByUUIDRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetByUUIDRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByUUIDRequest>, I>>(
    object: I
  ): GetByUUIDRequest {
    const message = createBaseGetByUUIDRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetByUUIDResponse(): GetByUUIDResponse {
  return { tokenData: undefined };
}

export const GetByUUIDResponse = {
  encode(
    message: GetByUUIDResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetByUUIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetByUUIDResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenData = TokenData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetByUUIDResponse {
    return {
      tokenData: isSet(object.tokenData)
        ? TokenData.fromJSON(object.tokenData)
        : undefined,
    };
  },

  toJSON(message: GetByUUIDResponse): unknown {
    const obj: any = {};
    message.tokenData !== undefined &&
      (obj.tokenData = message.tokenData
        ? TokenData.toJSON(message.tokenData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetByUUIDResponse>, I>>(
    object: I
  ): GetByUUIDResponse {
    const message = createBaseGetByUUIDResponse();
    message.tokenData =
      object.tokenData !== undefined && object.tokenData !== null
        ? TokenData.fromPartial(object.tokenData)
        : undefined;
    return message;
  },
};

function createBaseDisableByUUIDRequest(): DisableByUUIDRequest {
  return { namespace: "", uuid: "" };
}

export const DisableByUUIDRequest = {
  encode(
    message: DisableByUUIDRequest,
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
  ): DisableByUUIDRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisableByUUIDRequest();
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

  fromJSON(object: any): DisableByUUIDRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DisableByUUIDRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DisableByUUIDRequest>, I>>(
    object: I
  ): DisableByUUIDRequest {
    const message = createBaseDisableByUUIDRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDisableByUUIDResponse(): DisableByUUIDResponse {
  return {};
}

export const DisableByUUIDResponse = {
  encode(
    _: DisableByUUIDResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DisableByUUIDResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisableByUUIDResponse();
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

  fromJSON(_: any): DisableByUUIDResponse {
    return {};
  },

  toJSON(_: DisableByUUIDResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DisableByUUIDResponse>, I>>(
    _: I
  ): DisableByUUIDResponse {
    const message = createBaseDisableByUUIDResponse();
    return message;
  },
};

function createBaseAuthorizeRequest(): AuthorizeRequest {
  return { token: "", useCache: false };
}

export const AuthorizeRequest = {
  encode(
    message: AuthorizeRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthorizeRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthorizeRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.token = reader.string();
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

  fromJSON(object: any): AuthorizeRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: AuthorizeRequest): unknown {
    const obj: any = {};
    message.token !== undefined && (obj.token = message.token);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthorizeRequest>, I>>(
    object: I
  ): AuthorizeRequest {
    const message = createBaseAuthorizeRequest();
    message.token = object.token ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseAuthorizeResponse(): AuthorizeResponse {
  return { status: 0, tokenData: undefined };
}

export const AuthorizeResponse = {
  encode(
    message: AuthorizeResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthorizeResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthorizeResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.tokenData = TokenData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AuthorizeResponse {
    return {
      status: isSet(object.status)
        ? authorizeResponse_StatusFromJSON(object.status)
        : 0,
      tokenData: isSet(object.tokenData)
        ? TokenData.fromJSON(object.tokenData)
        : undefined,
    };
  },

  toJSON(message: AuthorizeResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = authorizeResponse_StatusToJSON(message.status));
    message.tokenData !== undefined &&
      (obj.tokenData = message.tokenData
        ? TokenData.toJSON(message.tokenData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthorizeResponse>, I>>(
    object: I
  ): AuthorizeResponse {
    const message = createBaseAuthorizeResponse();
    message.status = object.status ?? 0;
    message.tokenData =
      object.tokenData !== undefined && object.tokenData !== null
        ? TokenData.fromPartial(object.tokenData)
        : undefined;
    return message;
  },
};

function createBaseRefreshRequest(): RefreshRequest {
  return { refreshToken: "" };
}

export const RefreshRequest = {
  encode(
    message: RefreshRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.refreshToken !== "") {
      writer.uint32(10).string(message.refreshToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.refreshToken = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RefreshRequest {
    return {
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
    };
  },

  toJSON(message: RefreshRequest): unknown {
    const obj: any = {};
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefreshRequest>, I>>(
    object: I
  ): RefreshRequest {
    const message = createBaseRefreshRequest();
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseRefreshResponse(): RefreshResponse {
  return { status: 0, token: "", tokenData: undefined };
}

export const RefreshResponse = {
  encode(
    message: RefreshResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.token = reader.string();
          break;
        case 3:
          message.tokenData = TokenData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RefreshResponse {
    return {
      status: isSet(object.status)
        ? refreshResponse_StatusFromJSON(object.status)
        : 0,
      token: isSet(object.token) ? String(object.token) : "",
      tokenData: isSet(object.tokenData)
        ? TokenData.fromJSON(object.tokenData)
        : undefined,
    };
  },

  toJSON(message: RefreshResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = refreshResponse_StatusToJSON(message.status));
    message.token !== undefined && (obj.token = message.token);
    message.tokenData !== undefined &&
      (obj.tokenData = message.tokenData
        ? TokenData.toJSON(message.tokenData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefreshResponse>, I>>(
    object: I
  ): RefreshResponse {
    const message = createBaseRefreshResponse();
    message.status = object.status ?? 0;
    message.token = object.token ?? "";
    message.tokenData =
      object.tokenData !== undefined && object.tokenData !== null
        ? TokenData.fromPartial(object.tokenData)
        : undefined;
    return message;
  },
};

function createBaseTokensForIdentityRequest(): TokensForIdentityRequest {
  return { namespace: "", identity: "", activeFilter: 0, skip: 0, limit: 0 };
}

export const TokensForIdentityRequest = {
  encode(
    message: TokensForIdentityRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    if (message.activeFilter !== 0) {
      writer.uint32(24).int32(message.activeFilter);
    }
    if (message.skip !== 0) {
      writer.uint32(32).uint32(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(40).uint32(message.limit);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): TokensForIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTokensForIdentityRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.identity = reader.string();
          break;
        case 3:
          message.activeFilter = reader.int32() as any;
          break;
        case 4:
          message.skip = reader.uint32();
          break;
        case 5:
          message.limit = reader.uint32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TokensForIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      activeFilter: isSet(object.activeFilter)
        ? tokensForIdentityRequest_ActiveFilterFromJSON(object.activeFilter)
        : 0,
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: TokensForIdentityRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    message.activeFilter !== undefined &&
      (obj.activeFilter = tokensForIdentityRequest_ActiveFilterToJSON(
        message.activeFilter
      ));
    message.skip !== undefined && (obj.skip = Math.round(message.skip));
    message.limit !== undefined && (obj.limit = Math.round(message.limit));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TokensForIdentityRequest>, I>>(
    object: I
  ): TokensForIdentityRequest {
    const message = createBaseTokensForIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.activeFilter = object.activeFilter ?? 0;
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseTokensForIdentityResponse(): TokensForIdentityResponse {
  return { tokenData: undefined };
}

export const TokensForIdentityResponse = {
  encode(
    message: TokensForIdentityResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): TokensForIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTokensForIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.tokenData = TokenData.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): TokensForIdentityResponse {
    return {
      tokenData: isSet(object.tokenData)
        ? TokenData.fromJSON(object.tokenData)
        : undefined,
    };
  },

  toJSON(message: TokensForIdentityResponse): unknown {
    const obj: any = {};
    message.tokenData !== undefined &&
      (obj.tokenData = message.tokenData
        ? TokenData.toJSON(message.tokenData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<TokensForIdentityResponse>, I>>(
    object: I
  ): TokensForIdentityResponse {
    const message = createBaseTokensForIdentityResponse();
    message.tokenData =
      object.tokenData !== undefined && object.tokenData !== null
        ? TokenData.fromPartial(object.tokenData)
        : undefined;
    return message;
  },
};

/** Provides API to manage auth tokens */
export interface IAMTokenService {
  /** Create new token */
  Create(request: CreateRequest): Promise<CreateResponse>;
  /** Get token data using token UUID (unique identifier) */
  GetByUUID(request: GetByUUIDRequest): Promise<GetByUUIDResponse>;
  /** Disable token using its unique identifier */
  DisableByUUID(request: DisableByUUIDRequest): Promise<DisableByUUIDResponse>;
  /** Validates token and gets its scopes */
  Authorize(request: AuthorizeRequest): Promise<AuthorizeResponse>;
  /** Validates refresh token and create new token based on it. New token will have same scopes */
  Refresh(request: RefreshRequest): Promise<RefreshResponse>;
  /** Returns list of tokens for specified identity */
  TokensForIdentity(
    request: TokensForIdentityRequest
  ): Observable<TokensForIdentityResponse>;
}

export class IAMTokenServiceClientImpl implements IAMTokenService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.GetByUUID = this.GetByUUID.bind(this);
    this.DisableByUUID = this.DisableByUUID.bind(this);
    this.Authorize = this.Authorize.bind(this);
    this.Refresh = this.Refresh.bind(this);
    this.TokensForIdentity = this.TokensForIdentity.bind(this);
  }
  Create(request: CreateRequest): Promise<CreateResponse> {
    const data = CreateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_token.IAMTokenService",
      "Create",
      data
    );
    return promise.then((data) => CreateResponse.decode(new _m0.Reader(data)));
  }

  GetByUUID(request: GetByUUIDRequest): Promise<GetByUUIDResponse> {
    const data = GetByUUIDRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_token.IAMTokenService",
      "GetByUUID",
      data
    );
    return promise.then((data) =>
      GetByUUIDResponse.decode(new _m0.Reader(data))
    );
  }

  DisableByUUID(request: DisableByUUIDRequest): Promise<DisableByUUIDResponse> {
    const data = DisableByUUIDRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_token.IAMTokenService",
      "DisableByUUID",
      data
    );
    return promise.then((data) =>
      DisableByUUIDResponse.decode(new _m0.Reader(data))
    );
  }

  Authorize(request: AuthorizeRequest): Promise<AuthorizeResponse> {
    const data = AuthorizeRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_token.IAMTokenService",
      "Authorize",
      data
    );
    return promise.then((data) =>
      AuthorizeResponse.decode(new _m0.Reader(data))
    );
  }

  Refresh(request: RefreshRequest): Promise<RefreshResponse> {
    const data = RefreshRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_token.IAMTokenService",
      "Refresh",
      data
    );
    return promise.then((data) => RefreshResponse.decode(new _m0.Reader(data)));
  }

  TokensForIdentity(
    request: TokensForIdentityRequest
  ): Observable<TokensForIdentityResponse> {
    const data = TokensForIdentityRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(
      "native_iam_token.IAMTokenService",
      "TokensForIdentity",
      data
    );
    return result.pipe(
      map((data) => TokensForIdentityResponse.decode(new _m0.Reader(data)))
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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
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

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
