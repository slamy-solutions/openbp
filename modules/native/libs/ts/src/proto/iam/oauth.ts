/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_oauth";

/** Scope of the requested access. Check native_iam_policy for more information. */
export interface Scope {
  /** Namespace where this scope applies */
  namespace: string;
  /** Resources that can be accessed using token */
  resources: string[];
  /** Actions that can be done on the resources */
  actions: string[];
}

export interface CreateTokenWithPasswordRequest {
  /** Namespace where identity located. May be empty for global identity */
  namespace: string;
  /** Identity UUID */
  identity: string;
  /** Identity password */
  password: string;
  /** Arbitrary metadata. For example MAC/IP/information of the actor/application/browser/machine that created this token. The exact format of metadata is not defined, but JSON is suggested. */
  metadata: string;
  /** Scopes of the created token. Empty for creating token with all possible scopes for identity. */
  scopes: Scope[];
}

export interface CreateTokenWithPasswordResponse {
  /** Status of the token creation */
  status: CreateTokenWithPasswordResponse_Status;
  /** Token used for authentication and authorization. If status is not OK - empty string */
  accessToken: string;
  /** Token used for refreshing accessToken. If status is not OK - empty string */
  refreshToken: string;
}

export enum CreateTokenWithPasswordResponse_Status {
  /** OK - Everything is ok. Access and refresh tokens were successfully created */
  OK = 0,
  /** CREDENTIALS_INVALID - Login or password is not valid */
  CREDENTIALS_INVALID = 1,
  /** IDENTITY_NOT_ACTIVE - Identity was manually disabled */
  IDENTITY_NOT_ACTIVE = 2,
  /** UNAUTHORIZED - Not enough privileges to create token with specified scopes */
  UNAUTHORIZED = 3,
  UNRECOGNIZED = -1,
}

export function createTokenWithPasswordResponse_StatusFromJSON(
  object: any
): CreateTokenWithPasswordResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CreateTokenWithPasswordResponse_Status.OK;
    case 1:
    case "CREDENTIALS_INVALID":
      return CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID;
    case 2:
    case "IDENTITY_NOT_ACTIVE":
      return CreateTokenWithPasswordResponse_Status.IDENTITY_NOT_ACTIVE;
    case 3:
    case "UNAUTHORIZED":
      return CreateTokenWithPasswordResponse_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CreateTokenWithPasswordResponse_Status.UNRECOGNIZED;
  }
}

export function createTokenWithPasswordResponse_StatusToJSON(
  object: CreateTokenWithPasswordResponse_Status
): string {
  switch (object) {
    case CreateTokenWithPasswordResponse_Status.OK:
      return "OK";
    case CreateTokenWithPasswordResponse_Status.CREDENTIALS_INVALID:
      return "CREDENTIALS_INVALID";
    case CreateTokenWithPasswordResponse_Status.IDENTITY_NOT_ACTIVE:
      return "IDENTITY_NOT_ACTIVE";
    case CreateTokenWithPasswordResponse_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CreateTokenWithPasswordResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface CreateTokenWithOAuth2Request {
  /** OAuth2 provider */
  provider: string;
  /** Token issued by OAuth2 provider`` */
  token: string;
}

export interface CreateTokenWithOAuth2Response {
  /** Token used for authentication and authorization */
  accessToken: string;
  /** Token used for refreshing accessToken */
  refreshToken: string;
  /** Identity UUID */
  identity: string;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface RefreshTokenResponse {
  /** Status of the refresh */
  status: RefreshTokenResponse_Status;
  /** New access token */
  accessToken: string;
}

export enum RefreshTokenResponse_Status {
  /** OK - Everything is ok. New access token was successfully created */
  OK = 0,
  /** TOKEN_INVALID - Received token has bad format or its signature doesnt match */
  TOKEN_INVALID = 1,
  /** TOKEN_NOT_FOUND - Most probably token was deleted after its creation */
  TOKEN_NOT_FOUND = 2,
  /** TOKEN_DISABLED - Token was manually disabled */
  TOKEN_DISABLED = 3,
  /** TOKEN_EXPIRED - Token expired */
  TOKEN_EXPIRED = 4,
  /** TOKEN_IS_NOT_REFRESH_TOKEN - Provided token was recognized but most probably it is normal access token (not refresh one) */
  TOKEN_IS_NOT_REFRESH_TOKEN = 5,
  /** IDENTITY_NOT_FOUND - Identity wasnt founded. Most probably it was deleted after token creation */
  IDENTITY_NOT_FOUND = 6,
  /** IDENTITY_NOT_ACTIVE - Identity was manually disabled. */
  IDENTITY_NOT_ACTIVE = 7,
  /** IDENTITY_UNAUTHENTICATED - Most probably indentity policies changed and now its not possible to create token with same scopes */
  IDENTITY_UNAUTHENTICATED = 8,
  UNRECOGNIZED = -1,
}

export function refreshTokenResponse_StatusFromJSON(
  object: any
): RefreshTokenResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return RefreshTokenResponse_Status.OK;
    case 1:
    case "TOKEN_INVALID":
      return RefreshTokenResponse_Status.TOKEN_INVALID;
    case 2:
    case "TOKEN_NOT_FOUND":
      return RefreshTokenResponse_Status.TOKEN_NOT_FOUND;
    case 3:
    case "TOKEN_DISABLED":
      return RefreshTokenResponse_Status.TOKEN_DISABLED;
    case 4:
    case "TOKEN_EXPIRED":
      return RefreshTokenResponse_Status.TOKEN_EXPIRED;
    case 5:
    case "TOKEN_IS_NOT_REFRESH_TOKEN":
      return RefreshTokenResponse_Status.TOKEN_IS_NOT_REFRESH_TOKEN;
    case 6:
    case "IDENTITY_NOT_FOUND":
      return RefreshTokenResponse_Status.IDENTITY_NOT_FOUND;
    case 7:
    case "IDENTITY_NOT_ACTIVE":
      return RefreshTokenResponse_Status.IDENTITY_NOT_ACTIVE;
    case 8:
    case "IDENTITY_UNAUTHENTICATED":
      return RefreshTokenResponse_Status.IDENTITY_UNAUTHENTICATED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RefreshTokenResponse_Status.UNRECOGNIZED;
  }
}

export function refreshTokenResponse_StatusToJSON(
  object: RefreshTokenResponse_Status
): string {
  switch (object) {
    case RefreshTokenResponse_Status.OK:
      return "OK";
    case RefreshTokenResponse_Status.TOKEN_INVALID:
      return "TOKEN_INVALID";
    case RefreshTokenResponse_Status.TOKEN_NOT_FOUND:
      return "TOKEN_NOT_FOUND";
    case RefreshTokenResponse_Status.TOKEN_DISABLED:
      return "TOKEN_DISABLED";
    case RefreshTokenResponse_Status.TOKEN_EXPIRED:
      return "TOKEN_EXPIRED";
    case RefreshTokenResponse_Status.TOKEN_IS_NOT_REFRESH_TOKEN:
      return "TOKEN_IS_NOT_REFRESH_TOKEN";
    case RefreshTokenResponse_Status.IDENTITY_NOT_FOUND:
      return "IDENTITY_NOT_FOUND";
    case RefreshTokenResponse_Status.IDENTITY_NOT_ACTIVE:
      return "IDENTITY_NOT_ACTIVE";
    case RefreshTokenResponse_Status.IDENTITY_UNAUTHENTICATED:
      return "IDENTITY_UNAUTHENTICATED";
    case RefreshTokenResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface CheckAccessRequest {
  /** Token to verify */
  accessToken: string;
  /** Scopes for with to validate access */
  scopes: Scope[];
}

export interface CheckAccessResponse {
  /** Status of the verification */
  status: CheckAccessResponse_Status;
  /** Details of the status, that can be safelly returned and displayed to the user */
  message: string;
}

export enum CheckAccessResponse_Status {
  /** OK - Provided token allows to access scopes */
  OK = 0,
  /** TOKEN_INVALID - Received token has bad format or its signature doesnt match */
  TOKEN_INVALID = 1,
  /** TOKEN_NOT_FOUND - Most probably token was deleted after its creation */
  TOKEN_NOT_FOUND = 2,
  /** TOKEN_DISABLED - Token was manually disabled */
  TOKEN_DISABLED = 3,
  /** TOKEN_EXPIRED - Token expired */
  TOKEN_EXPIRED = 4,
  /** UNAUTHORIZED - Token has not enought privileges to access specified scopes */
  UNAUTHORIZED = 5,
  UNRECOGNIZED = -1,
}

export function checkAccessResponse_StatusFromJSON(
  object: any
): CheckAccessResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CheckAccessResponse_Status.OK;
    case 1:
    case "TOKEN_INVALID":
      return CheckAccessResponse_Status.TOKEN_INVALID;
    case 2:
    case "TOKEN_NOT_FOUND":
      return CheckAccessResponse_Status.TOKEN_NOT_FOUND;
    case 3:
    case "TOKEN_DISABLED":
      return CheckAccessResponse_Status.TOKEN_DISABLED;
    case 4:
    case "TOKEN_EXPIRED":
      return CheckAccessResponse_Status.TOKEN_EXPIRED;
    case 5:
    case "UNAUTHORIZED":
      return CheckAccessResponse_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CheckAccessResponse_Status.UNRECOGNIZED;
  }
}

export function checkAccessResponse_StatusToJSON(
  object: CheckAccessResponse_Status
): string {
  switch (object) {
    case CheckAccessResponse_Status.OK:
      return "OK";
    case CheckAccessResponse_Status.TOKEN_INVALID:
      return "TOKEN_INVALID";
    case CheckAccessResponse_Status.TOKEN_NOT_FOUND:
      return "TOKEN_NOT_FOUND";
    case CheckAccessResponse_Status.TOKEN_DISABLED:
      return "TOKEN_DISABLED";
    case CheckAccessResponse_Status.TOKEN_EXPIRED:
      return "TOKEN_EXPIRED";
    case CheckAccessResponse_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CheckAccessResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
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

function createBaseCreateTokenWithPasswordRequest(): CreateTokenWithPasswordRequest {
  return {
    namespace: "",
    identity: "",
    password: "",
    metadata: "",
    scopes: [],
  };
}

export const CreateTokenWithPasswordRequest = {
  encode(
    message: CreateTokenWithPasswordRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    if (message.password !== "") {
      writer.uint32(26).string(message.password);
    }
    if (message.metadata !== "") {
      writer.uint32(34).string(message.metadata);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(42).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTokenWithPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithPasswordRequest();
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
          message.password = reader.string();
          break;
        case 4:
          message.metadata = reader.string();
          break;
        case 5:
          message.scopes.push(Scope.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithPasswordRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      password: isSet(object.password) ? String(object.password) : "",
      metadata: isSet(object.metadata) ? String(object.metadata) : "",
      scopes: Array.isArray(object?.scopes)
        ? object.scopes.map((e: any) => Scope.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CreateTokenWithPasswordRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    message.password !== undefined && (obj.password = message.password);
    message.metadata !== undefined && (obj.metadata = message.metadata);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => (e ? Scope.toJSON(e) : undefined));
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTokenWithPasswordRequest>, I>>(
    object: I
  ): CreateTokenWithPasswordRequest {
    const message = createBaseCreateTokenWithPasswordRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.password = object.password ?? "";
    message.metadata = object.metadata ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCreateTokenWithPasswordResponse(): CreateTokenWithPasswordResponse {
  return { status: 0, accessToken: "", refreshToken: "" };
}

export const CreateTokenWithPasswordResponse = {
  encode(
    message: CreateTokenWithPasswordResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.accessToken !== "") {
      writer.uint32(18).string(message.accessToken);
    }
    if (message.refreshToken !== "") {
      writer.uint32(26).string(message.refreshToken);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTokenWithPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithPasswordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.accessToken = reader.string();
          break;
        case 3:
          message.refreshToken = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithPasswordResponse {
    return {
      status: isSet(object.status)
        ? createTokenWithPasswordResponse_StatusFromJSON(object.status)
        : 0,
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
    };
  },

  toJSON(message: CreateTokenWithPasswordResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = createTokenWithPasswordResponse_StatusToJSON(
        message.status
      ));
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTokenWithPasswordResponse>, I>>(
    object: I
  ): CreateTokenWithPasswordResponse {
    const message = createBaseCreateTokenWithPasswordResponse();
    message.status = object.status ?? 0;
    message.accessToken = object.accessToken ?? "";
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseCreateTokenWithOAuth2Request(): CreateTokenWithOAuth2Request {
  return { provider: "", token: "" };
}

export const CreateTokenWithOAuth2Request = {
  encode(
    message: CreateTokenWithOAuth2Request,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.provider !== "") {
      writer.uint32(10).string(message.provider);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTokenWithOAuth2Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithOAuth2Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.provider = reader.string();
          break;
        case 2:
          message.token = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithOAuth2Request {
    return {
      provider: isSet(object.provider) ? String(object.provider) : "",
      token: isSet(object.token) ? String(object.token) : "",
    };
  },

  toJSON(message: CreateTokenWithOAuth2Request): unknown {
    const obj: any = {};
    message.provider !== undefined && (obj.provider = message.provider);
    message.token !== undefined && (obj.token = message.token);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTokenWithOAuth2Request>, I>>(
    object: I
  ): CreateTokenWithOAuth2Request {
    const message = createBaseCreateTokenWithOAuth2Request();
    message.provider = object.provider ?? "";
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseCreateTokenWithOAuth2Response(): CreateTokenWithOAuth2Response {
  return { accessToken: "", refreshToken: "", identity: "" };
}

export const CreateTokenWithOAuth2Response = {
  encode(
    message: CreateTokenWithOAuth2Response,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.refreshToken !== "") {
      writer.uint32(18).string(message.refreshToken);
    }
    if (message.identity !== "") {
      writer.uint32(26).string(message.identity);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateTokenWithOAuth2Response {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithOAuth2Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.refreshToken = reader.string();
          break;
        case 3:
          message.identity = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithOAuth2Response {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: CreateTokenWithOAuth2Response): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    message.identity !== undefined && (obj.identity = message.identity);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateTokenWithOAuth2Response>, I>>(
    object: I
  ): CreateTokenWithOAuth2Response {
    const message = createBaseCreateTokenWithOAuth2Response();
    message.accessToken = object.accessToken ?? "";
    message.refreshToken = object.refreshToken ?? "";
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseRefreshTokenRequest(): RefreshTokenRequest {
  return { refreshToken: "" };
}

export const RefreshTokenRequest = {
  encode(
    message: RefreshTokenRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.refreshToken !== "") {
      writer.uint32(10).string(message.refreshToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshTokenRequest();
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

  fromJSON(object: any): RefreshTokenRequest {
    return {
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
    };
  },

  toJSON(message: RefreshTokenRequest): unknown {
    const obj: any = {};
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefreshTokenRequest>, I>>(
    object: I
  ): RefreshTokenRequest {
    const message = createBaseRefreshTokenRequest();
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseRefreshTokenResponse(): RefreshTokenResponse {
  return { status: 0, accessToken: "" };
}

export const RefreshTokenResponse = {
  encode(
    message: RefreshTokenResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.accessToken !== "") {
      writer.uint32(18).string(message.accessToken);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): RefreshTokenResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshTokenResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.accessToken = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): RefreshTokenResponse {
    return {
      status: isSet(object.status)
        ? refreshTokenResponse_StatusFromJSON(object.status)
        : 0,
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
    };
  },

  toJSON(message: RefreshTokenResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = refreshTokenResponse_StatusToJSON(message.status));
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefreshTokenResponse>, I>>(
    object: I
  ): RefreshTokenResponse {
    const message = createBaseRefreshTokenResponse();
    message.status = object.status ?? 0;
    message.accessToken = object.accessToken ?? "";
    return message;
  },
};

function createBaseCheckAccessRequest(): CheckAccessRequest {
  return { accessToken: "", scopes: [] };
}

export const CheckAccessRequest = {
  encode(
    message: CheckAccessRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.scopes.push(Scope.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CheckAccessRequest {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      scopes: Array.isArray(object?.scopes)
        ? object.scopes.map((e: any) => Scope.fromJSON(e))
        : [],
    };
  },

  toJSON(message: CheckAccessRequest): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    if (message.scopes) {
      obj.scopes = message.scopes.map((e) => (e ? Scope.toJSON(e) : undefined));
    } else {
      obj.scopes = [];
    }
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CheckAccessRequest>, I>>(
    object: I
  ): CheckAccessRequest {
    const message = createBaseCheckAccessRequest();
    message.accessToken = object.accessToken ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCheckAccessResponse(): CheckAccessResponse {
  return { status: 0, message: "" };
}

export const CheckAccessResponse = {
  encode(
    message: CheckAccessResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.status = reader.int32() as any;
          break;
        case 2:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CheckAccessResponse {
    return {
      status: isSet(object.status)
        ? checkAccessResponse_StatusFromJSON(object.status)
        : 0,
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: CheckAccessResponse): unknown {
    const obj: any = {};
    message.status !== undefined &&
      (obj.status = checkAccessResponse_StatusToJSON(message.status));
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CheckAccessResponse>, I>>(
    object: I
  ): CheckAccessResponse {
    const message = createBaseCheckAccessResponse();
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    return message;
  },
};

/** Provides API for OAuth ("Open Authorization") style access control */
export interface IAMOAuthService {
  /** Create access token and refresh token using password */
  CreateTokenWithPassword(
    request: CreateTokenWithPasswordRequest
  ): Promise<CreateTokenWithPasswordResponse>;
  /** Creates new access token using refresh tokenna */
  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse>;
  /**
   * rpc VerifyResoureAccess(VerifyResourceAccessRequest) returns (VerifyResourceAccessResponse);
   * Checks if token is allowed to perform actions from the specified scopes
   */
  CheckAccess(request: CheckAccessRequest): Promise<CheckAccessResponse>;
}

export class IAMOAuthServiceClientImpl implements IAMOAuthService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateTokenWithPassword = this.CreateTokenWithPassword.bind(this);
    this.RefreshToken = this.RefreshToken.bind(this);
    this.CheckAccess = this.CheckAccess.bind(this);
  }
  CreateTokenWithPassword(
    request: CreateTokenWithPasswordRequest
  ): Promise<CreateTokenWithPasswordResponse> {
    const data = CreateTokenWithPasswordRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_oauth.IAMOAuthService",
      "CreateTokenWithPassword",
      data
    );
    return promise.then((data) =>
      CreateTokenWithPasswordResponse.decode(new _m0.Reader(data))
    );
  }

  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    const data = RefreshTokenRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_oauth.IAMOAuthService",
      "RefreshToken",
      data
    );
    return promise.then((data) =>
      RefreshTokenResponse.decode(new _m0.Reader(data))
    );
  }

  CheckAccess(request: CheckAccessRequest): Promise<CheckAccessResponse> {
    const data = CheckAccessRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_oauth.IAMOAuthService",
      "CheckAccess",
      data
    );
    return promise.then((data) =>
      CheckAccessResponse.decode(new _m0.Reader(data))
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
