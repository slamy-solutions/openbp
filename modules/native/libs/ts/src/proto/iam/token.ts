/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../google/protobuf/timestamp";

export const protobufPackage = "native_iam_token";

/** Scope defines what can be accessed by token and what actions can be performed on accessible resources. Scope is bounded the namespace where resources are located. */
export interface Scope {
  /** Namespace to which scope is bounded. Empty string if scope is not bounded to any namespace (is global) */
  namespace: string;
  /** Resources that can be accessed using this token */
  resources: string[];
  /** Actions that can be performed on accessible resources */
  actions: string[];
  /** Should scope work in all namespaces */
  namespaceIndependent: boolean;
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
  expiresAt:
    | Date
    | undefined;
  /** List of token scopes. Describes what actions can token perform on what resources */
  scopes: Scope[];
  /** Datetime when token was created */
  createdAt:
    | Date
    | undefined;
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

export interface GetRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
  /** Use cache for faster authorization. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) */
  useCache: boolean;
}

export interface GetResponse {
  /** Actual token data */
  tokenData: TokenData | undefined;
}

export interface RawGetRequest {
  /** Refresh or access token */
  token: string;
  /** Use cache for faster authorization. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) */
  useCache: boolean;
}

export interface RawGetResponse {
  /** Actual token data */
  tokenData: TokenData | undefined;
}

export interface DeleteRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
}

export interface DeleteResponse {
  /** Indicates if token existed before request or it was already deleted. */
  existed: boolean;
}

export interface DisableRequest {
  /** Namespace of the token. Empty for global token. */
  namespace: string;
  /** Unique identifier of the token inside namespace */
  uuid: string;
}

export interface DisableResponse {
}

export interface ValidateRequest {
  /** Token to validate */
  token: string;
  /** Use cache for faster validation. Cache has a very low chance to not be valid. If cache is not valid it will be deleted after short period of time (30 seconds by default) */
  useCache: boolean;
}

export interface ValidateResponse {
  status: ValidateResponse_Status;
  /** Token data. Null if status is not OK */
  tokenData: TokenData | undefined;
}

export enum ValidateResponse_Status {
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

export function validateResponse_StatusFromJSON(object: any): ValidateResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return ValidateResponse_Status.OK;
    case 1:
    case "INVALID":
      return ValidateResponse_Status.INVALID;
    case 2:
    case "NOT_FOUND":
      return ValidateResponse_Status.NOT_FOUND;
    case 3:
    case "DISABLED":
      return ValidateResponse_Status.DISABLED;
    case 4:
    case "EXPIRED":
      return ValidateResponse_Status.EXPIRED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ValidateResponse_Status.UNRECOGNIZED;
  }
}

export function validateResponse_StatusToJSON(object: ValidateResponse_Status): string {
  switch (object) {
    case ValidateResponse_Status.OK:
      return "OK";
    case ValidateResponse_Status.INVALID:
      return "INVALID";
    case ValidateResponse_Status.NOT_FOUND:
      return "NOT_FOUND";
    case ValidateResponse_Status.DISABLED:
      return "DISABLED";
    case ValidateResponse_Status.EXPIRED:
      return "EXPIRED";
    case ValidateResponse_Status.UNRECOGNIZED:
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

export function refreshResponse_StatusFromJSON(object: any): RefreshResponse_Status {
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

export function refreshResponse_StatusToJSON(object: RefreshResponse_Status): string {
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

export interface GetTokensForIdentityRequest {
  /** Namespace where token is located. Empty for global token */
  namespace: string;
  /** Identity unique identifier inside namespace */
  identity: string;
  /** Perform results filtering on "active" property of the token */
  activeFilter: GetTokensForIdentityRequest_ActiveFilter;
  /** Skip number of results before returning actual tokens. Set to 0 in order not to skip */
  skip: number;
  /** Limit number of returned results. Set to 0 in order to remove limit and return all possible results up to the end. */
  limit: number;
}

export enum GetTokensForIdentityRequest_ActiveFilter {
  /** ALL - Get all token */
  ALL = 0,
  /** ONLY_ACTIVE - Only get tokens that wasnt disabled and not expired */
  ONLY_ACTIVE = 1,
  /** ONLY_NOT_ACTIVE - Only get tokens that are disabled or expired */
  ONLY_NOT_ACTIVE = 2,
  UNRECOGNIZED = -1,
}

export function getTokensForIdentityRequest_ActiveFilterFromJSON(
  object: any,
): GetTokensForIdentityRequest_ActiveFilter {
  switch (object) {
    case 0:
    case "ALL":
      return GetTokensForIdentityRequest_ActiveFilter.ALL;
    case 1:
    case "ONLY_ACTIVE":
      return GetTokensForIdentityRequest_ActiveFilter.ONLY_ACTIVE;
    case 2:
    case "ONLY_NOT_ACTIVE":
      return GetTokensForIdentityRequest_ActiveFilter.ONLY_NOT_ACTIVE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return GetTokensForIdentityRequest_ActiveFilter.UNRECOGNIZED;
  }
}

export function getTokensForIdentityRequest_ActiveFilterToJSON(
  object: GetTokensForIdentityRequest_ActiveFilter,
): string {
  switch (object) {
    case GetTokensForIdentityRequest_ActiveFilter.ALL:
      return "ALL";
    case GetTokensForIdentityRequest_ActiveFilter.ONLY_ACTIVE:
      return "ONLY_ACTIVE";
    case GetTokensForIdentityRequest_ActiveFilter.ONLY_NOT_ACTIVE:
      return "ONLY_NOT_ACTIVE";
    case GetTokensForIdentityRequest_ActiveFilter.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface GetTokensForIdentityResponse {
  /** Actual token data */
  tokenData: TokenData | undefined;
}

function createBaseScope(): Scope {
  return { namespace: "", resources: [], actions: [], namespaceIndependent: false };
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
    if (message.namespaceIndependent === true) {
      writer.uint32(32).bool(message.namespaceIndependent);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Scope {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseScope();
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

          message.resources.push(reader.string());
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.actions.push(reader.string());
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.namespaceIndependent = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Scope {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      resources: Array.isArray(object?.resources) ? object.resources.map((e: any) => String(e)) : [],
      actions: Array.isArray(object?.actions) ? object.actions.map((e: any) => String(e)) : [],
      namespaceIndependent: isSet(object.namespaceIndependent) ? Boolean(object.namespaceIndependent) : false,
    };
  },

  toJSON(message: Scope): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.resources?.length) {
      obj.resources = message.resources;
    }
    if (message.actions?.length) {
      obj.actions = message.actions;
    }
    if (message.namespaceIndependent === true) {
      obj.namespaceIndependent = message.namespaceIndependent;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Scope>, I>>(base?: I): Scope {
    return Scope.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Scope>, I>>(object: I): Scope {
    const message = createBaseScope();
    message.namespace = object.namespace ?? "";
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    message.namespaceIndependent = object.namespaceIndependent ?? false;
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
  encode(message: TokenData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
      Timestamp.encode(toTimestamp(message.expiresAt), writer.uint32(42).fork()).ldelim();
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(50).fork()).ldelim();
    }
    if (message.createdAt !== undefined) {
      Timestamp.encode(toTimestamp(message.createdAt), writer.uint32(58).fork()).ldelim();
    }
    if (message.creationMetadata !== "") {
      writer.uint32(66).string(message.creationMetadata);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): TokenData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseTokenData();
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

          message.identity = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.disabled = reader.bool();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.expiresAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
        case 7:
          if (tag !== 58) {
            break;
          }

          message.createdAt = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.creationMetadata = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): TokenData {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      disabled: isSet(object.disabled) ? Boolean(object.disabled) : false,
      expiresAt: isSet(object.expiresAt) ? fromJsonTimestamp(object.expiresAt) : undefined,
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
      createdAt: isSet(object.createdAt) ? fromJsonTimestamp(object.createdAt) : undefined,
      creationMetadata: isSet(object.creationMetadata) ? String(object.creationMetadata) : "",
    };
  },

  toJSON(message: TokenData): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.disabled === true) {
      obj.disabled = message.disabled;
    }
    if (message.expiresAt !== undefined) {
      obj.expiresAt = message.expiresAt.toISOString();
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    if (message.createdAt !== undefined) {
      obj.createdAt = message.createdAt.toISOString();
    }
    if (message.creationMetadata !== "") {
      obj.creationMetadata = message.creationMetadata;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<TokenData>, I>>(base?: I): TokenData {
    return TokenData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<TokenData>, I>>(object: I): TokenData {
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
  encode(message: CreateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

          message.identity = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
        case 8:
          if (tag !== 66) {
            break;
          }

          message.metadata = reader.string();
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
      identity: isSet(object.identity) ? String(object.identity) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
      metadata: isSet(object.metadata) ? String(object.metadata) : "",
    };
  },

  toJSON(message: CreateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    if (message.metadata !== "") {
      obj.metadata = message.metadata;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateRequest>, I>>(base?: I): CreateRequest {
    return CreateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateRequest>, I>>(object: I): CreateRequest {
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
  encode(message: CreateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

          message.token = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.refreshToken = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.tokenData = TokenData.decode(reader, reader.uint32());
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
    return {
      token: isSet(object.token) ? String(object.token) : "",
      refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "",
      tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined,
    };
  },

  toJSON(message: CreateResponse): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateResponse>, I>>(base?: I): CreateResponse {
    return CreateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateResponse>, I>>(object: I): CreateResponse {
    const message = createBaseCreateResponse();
    message.token = object.token ?? "";
    message.refreshToken = object.refreshToken ?? "";
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
      ? TokenData.fromPartial(object.tokenData)
      : undefined;
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
  return { tokenData: undefined };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(10).fork()).ldelim();
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

          message.tokenData = TokenData.decode(reader, reader.uint32());
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
    return { tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
      ? TokenData.fromPartial(object.tokenData)
      : undefined;
    return message;
  },
};

function createBaseRawGetRequest(): RawGetRequest {
  return { token: "", useCache: false };
}

export const RawGetRequest = {
  encode(message: RawGetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RawGetRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRawGetRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.token = reader.string();
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

  fromJSON(object: any): RawGetRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: RawGetRequest): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RawGetRequest>, I>>(base?: I): RawGetRequest {
    return RawGetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RawGetRequest>, I>>(object: I): RawGetRequest {
    const message = createBaseRawGetRequest();
    message.token = object.token ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseRawGetResponse(): RawGetResponse {
  return { tokenData: undefined };
}

export const RawGetResponse = {
  encode(message: RawGetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RawGetResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRawGetResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.tokenData = TokenData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RawGetResponse {
    return { tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined };
  },

  toJSON(message: RawGetResponse): unknown {
    const obj: any = {};
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RawGetResponse>, I>>(base?: I): RawGetResponse {
    return RawGetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RawGetResponse>, I>>(object: I): RawGetResponse {
    const message = createBaseRawGetResponse();
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
      ? TokenData.fromPartial(object.tokenData)
      : undefined;
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

function createBaseDisableRequest(): DisableRequest {
  return { namespace: "", uuid: "" };
}

export const DisableRequest = {
  encode(message: DisableRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DisableRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisableRequest();
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

  fromJSON(object: any): DisableRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DisableRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DisableRequest>, I>>(base?: I): DisableRequest {
    return DisableRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DisableRequest>, I>>(object: I): DisableRequest {
    const message = createBaseDisableRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDisableResponse(): DisableResponse {
  return {};
}

export const DisableResponse = {
  encode(_: DisableResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DisableResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisableResponse();
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

  fromJSON(_: any): DisableResponse {
    return {};
  },

  toJSON(_: DisableResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<DisableResponse>, I>>(base?: I): DisableResponse {
    return DisableResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DisableResponse>, I>>(_: I): DisableResponse {
    const message = createBaseDisableResponse();
    return message;
  },
};

function createBaseValidateRequest(): ValidateRequest {
  return { token: "", useCache: false };
}

export const ValidateRequest = {
  encode(message: ValidateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.useCache === true) {
      writer.uint32(16).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.token = reader.string();
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

  fromJSON(object: any): ValidateRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ValidateRequest): unknown {
    const obj: any = {};
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ValidateRequest>, I>>(base?: I): ValidateRequest {
    return ValidateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ValidateRequest>, I>>(object: I): ValidateRequest {
    const message = createBaseValidateRequest();
    message.token = object.token ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseValidateResponse(): ValidateResponse {
  return { status: 0, tokenData: undefined };
}

export const ValidateResponse = {
  encode(message: ValidateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.status = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.tokenData = TokenData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ValidateResponse {
    return {
      status: isSet(object.status) ? validateResponse_StatusFromJSON(object.status) : 0,
      tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined,
    };
  },

  toJSON(message: ValidateResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = validateResponse_StatusToJSON(message.status);
    }
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ValidateResponse>, I>>(base?: I): ValidateResponse {
    return ValidateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ValidateResponse>, I>>(object: I): ValidateResponse {
    const message = createBaseValidateResponse();
    message.status = object.status ?? 0;
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
      ? TokenData.fromPartial(object.tokenData)
      : undefined;
    return message;
  },
};

function createBaseRefreshRequest(): RefreshRequest {
  return { refreshToken: "" };
}

export const RefreshRequest = {
  encode(message: RefreshRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.refreshToken !== "") {
      writer.uint32(10).string(message.refreshToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.refreshToken = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RefreshRequest {
    return { refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "" };
  },

  toJSON(message: RefreshRequest): unknown {
    const obj: any = {};
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshRequest>, I>>(base?: I): RefreshRequest {
    return RefreshRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshRequest>, I>>(object: I): RefreshRequest {
    const message = createBaseRefreshRequest();
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseRefreshResponse(): RefreshResponse {
  return { status: 0, token: "", tokenData: undefined };
}

export const RefreshResponse = {
  encode(message: RefreshResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.status = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.token = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.tokenData = TokenData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RefreshResponse {
    return {
      status: isSet(object.status) ? refreshResponse_StatusFromJSON(object.status) : 0,
      token: isSet(object.token) ? String(object.token) : "",
      tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined,
    };
  },

  toJSON(message: RefreshResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = refreshResponse_StatusToJSON(message.status);
    }
    if (message.token !== "") {
      obj.token = message.token;
    }
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshResponse>, I>>(base?: I): RefreshResponse {
    return RefreshResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshResponse>, I>>(object: I): RefreshResponse {
    const message = createBaseRefreshResponse();
    message.status = object.status ?? 0;
    message.token = object.token ?? "";
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
      ? TokenData.fromPartial(object.tokenData)
      : undefined;
    return message;
  },
};

function createBaseGetTokensForIdentityRequest(): GetTokensForIdentityRequest {
  return { namespace: "", identity: "", activeFilter: 0, skip: 0, limit: 0 };
}

export const GetTokensForIdentityRequest = {
  encode(message: GetTokensForIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTokensForIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTokensForIdentityRequest();
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

          message.activeFilter = reader.int32() as any;
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.skip = reader.uint32();
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.limit = reader.uint32();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetTokensForIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      activeFilter: isSet(object.activeFilter)
        ? getTokensForIdentityRequest_ActiveFilterFromJSON(object.activeFilter)
        : 0,
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: GetTokensForIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.activeFilter !== 0) {
      obj.activeFilter = getTokensForIdentityRequest_ActiveFilterToJSON(message.activeFilter);
    }
    if (message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetTokensForIdentityRequest>, I>>(base?: I): GetTokensForIdentityRequest {
    return GetTokensForIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetTokensForIdentityRequest>, I>>(object: I): GetTokensForIdentityRequest {
    const message = createBaseGetTokensForIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.activeFilter = object.activeFilter ?? 0;
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseGetTokensForIdentityResponse(): GetTokensForIdentityResponse {
  return { tokenData: undefined };
}

export const GetTokensForIdentityResponse = {
  encode(message: GetTokensForIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.tokenData !== undefined) {
      TokenData.encode(message.tokenData, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetTokensForIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetTokensForIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.tokenData = TokenData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetTokensForIdentityResponse {
    return { tokenData: isSet(object.tokenData) ? TokenData.fromJSON(object.tokenData) : undefined };
  },

  toJSON(message: GetTokensForIdentityResponse): unknown {
    const obj: any = {};
    if (message.tokenData !== undefined) {
      obj.tokenData = TokenData.toJSON(message.tokenData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetTokensForIdentityResponse>, I>>(base?: I): GetTokensForIdentityResponse {
    return GetTokensForIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetTokensForIdentityResponse>, I>>(object: I): GetTokensForIdentityResponse {
    const message = createBaseGetTokensForIdentityResponse();
    message.tokenData = (object.tokenData !== undefined && object.tokenData !== null)
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
  Get(request: GetRequest): Promise<GetResponse>;
  /** Get token data using raw access/refresh token. Validates if token still exists in the system. */
  RawGet(request: RawGetRequest): Promise<RawGetResponse>;
  /** Delete token using token UUID (unique identifier) */
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  /** Disable token using its unique identifier */
  Disable(request: DisableRequest): Promise<DisableResponse>;
  /** Validates token and gets its data */
  Validate(request: ValidateRequest): Promise<ValidateResponse>;
  /** Validates refresh token and create new token based on it. New token will have same scopes */
  Refresh(request: RefreshRequest): Promise<RefreshResponse>;
  /** Returns list of tokens for specified identity */
  GetTokensForIdentity(request: GetTokensForIdentityRequest): Observable<GetTokensForIdentityResponse>;
}

export const IAMTokenServiceServiceName = "native_iam_token.IAMTokenService";
export class IAMTokenServiceClientImpl implements IAMTokenService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMTokenServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Get = this.Get.bind(this);
    this.RawGet = this.RawGet.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Disable = this.Disable.bind(this);
    this.Validate = this.Validate.bind(this);
    this.Refresh = this.Refresh.bind(this);
    this.GetTokensForIdentity = this.GetTokensForIdentity.bind(this);
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

  RawGet(request: RawGetRequest): Promise<RawGetResponse> {
    const data = RawGetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RawGet", data);
    return promise.then((data) => RawGetResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteResponse.decode(_m0.Reader.create(data)));
  }

  Disable(request: DisableRequest): Promise<DisableResponse> {
    const data = DisableRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Disable", data);
    return promise.then((data) => DisableResponse.decode(_m0.Reader.create(data)));
  }

  Validate(request: ValidateRequest): Promise<ValidateResponse> {
    const data = ValidateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Validate", data);
    return promise.then((data) => ValidateResponse.decode(_m0.Reader.create(data)));
  }

  Refresh(request: RefreshRequest): Promise<RefreshResponse> {
    const data = RefreshRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Refresh", data);
    return promise.then((data) => RefreshResponse.decode(_m0.Reader.create(data)));
  }

  GetTokensForIdentity(request: GetTokensForIdentityRequest): Observable<GetTokensForIdentityResponse> {
    const data = GetTokensForIdentityRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "GetTokensForIdentity", data);
    return result.pipe(map((data) => GetTokensForIdentityResponse.decode(_m0.Reader.create(data))));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
  clientStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Promise<Uint8Array>;
  serverStreamingRequest(service: string, method: string, data: Uint8Array): Observable<Uint8Array>;
  bidirectionalStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Observable<Uint8Array>;
}

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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
