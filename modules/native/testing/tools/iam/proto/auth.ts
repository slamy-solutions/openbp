/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_auth";

export interface CreateTokenWithPasswordRequest {
  /** Namespace where identity located. May be empty for global identity */
  namespace: string;
  /** Identity UUID */
  identity: string;
  /** Identity password */
  password: string;
  /** Arbitrary metadata. For example MAC/IP/information of the actor/application/browser/machine that created this token. The exact format of metadata is not defined, but JSON is suggested. */
  metadata: string;
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
  OK = 0,
  UNAUTHENTICATED = 401,
  UNAUTHORIZED = 403,
  UNRECOGNIZED = -1,
}

export function createTokenWithPasswordResponse_StatusFromJSON(
  object: any
): CreateTokenWithPasswordResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CreateTokenWithPasswordResponse_Status.OK;
    case 401:
    case "UNAUTHENTICATED":
      return CreateTokenWithPasswordResponse_Status.UNAUTHENTICATED;
    case 403:
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
    case CreateTokenWithPasswordResponse_Status.UNAUTHENTICATED:
      return "UNAUTHENTICATED";
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
  /** New access token */
  accessToken: string;
}

export interface InvalidateTokenRequest {
  /** Refresh or access token to invalidate. Both tokens will be invalidated */
  token: string;
}

export interface InvalidateTokenResponse {}

export interface VerifyTokenAccessRequest {
  /** Token to verify */
  accessToken: string;
  /** Namespace where to access resources */
  namespace: string;
  /** What resources to theck */
  resources: string[];
  /** What actions token must be able to perform for resources */
  actions: string[];
}

export interface VerifyTokenAccessResponse {
  hasAccess: boolean;
}

function createBaseCreateTokenWithPasswordRequest(): CreateTokenWithPasswordRequest {
  return { namespace: "", identity: "", password: "", metadata: "" };
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
    };
  },

  toJSON(message: CreateTokenWithPasswordRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    message.password !== undefined && (obj.password = message.password);
    message.metadata !== undefined && (obj.metadata = message.metadata);
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
  return { accessToken: "" };
}

export const RefreshTokenResponse = {
  encode(
    message: RefreshTokenResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
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
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
    };
  },

  toJSON(message: RefreshTokenResponse): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<RefreshTokenResponse>, I>>(
    object: I
  ): RefreshTokenResponse {
    const message = createBaseRefreshTokenResponse();
    message.accessToken = object.accessToken ?? "";
    return message;
  },
};

function createBaseInvalidateTokenRequest(): InvalidateTokenRequest {
  return { token: "" };
}

export const InvalidateTokenRequest = {
  encode(
    message: InvalidateTokenRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): InvalidateTokenRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInvalidateTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.token = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): InvalidateTokenRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
    };
  },

  toJSON(message: InvalidateTokenRequest): unknown {
    const obj: any = {};
    message.token !== undefined && (obj.token = message.token);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<InvalidateTokenRequest>, I>>(
    object: I
  ): InvalidateTokenRequest {
    const message = createBaseInvalidateTokenRequest();
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseInvalidateTokenResponse(): InvalidateTokenResponse {
  return {};
}

export const InvalidateTokenResponse = {
  encode(
    _: InvalidateTokenResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): InvalidateTokenResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseInvalidateTokenResponse();
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

  fromJSON(_: any): InvalidateTokenResponse {
    return {};
  },

  toJSON(_: InvalidateTokenResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<InvalidateTokenResponse>, I>>(
    _: I
  ): InvalidateTokenResponse {
    const message = createBaseInvalidateTokenResponse();
    return message;
  },
};

function createBaseVerifyTokenAccessRequest(): VerifyTokenAccessRequest {
  return { accessToken: "", namespace: "", resources: [], actions: [] };
}

export const VerifyTokenAccessRequest = {
  encode(
    message: VerifyTokenAccessRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.namespace !== "") {
      writer.uint32(18).string(message.namespace);
    }
    for (const v of message.resources) {
      writer.uint32(26).string(v!);
    }
    for (const v of message.actions) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): VerifyTokenAccessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyTokenAccessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.namespace = reader.string();
          break;
        case 3:
          message.resources.push(reader.string());
          break;
        case 4:
          message.actions.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyTokenAccessRequest {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      resources: Array.isArray(object?.resources)
        ? object.resources.map((e: any) => String(e))
        : [],
      actions: Array.isArray(object?.actions)
        ? object.actions.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: VerifyTokenAccessRequest): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
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

  fromPartial<I extends Exact<DeepPartial<VerifyTokenAccessRequest>, I>>(
    object: I
  ): VerifyTokenAccessRequest {
    const message = createBaseVerifyTokenAccessRequest();
    message.accessToken = object.accessToken ?? "";
    message.namespace = object.namespace ?? "";
    message.resources = object.resources?.map((e) => e) || [];
    message.actions = object.actions?.map((e) => e) || [];
    return message;
  },
};

function createBaseVerifyTokenAccessResponse(): VerifyTokenAccessResponse {
  return { hasAccess: false };
}

export const VerifyTokenAccessResponse = {
  encode(
    message: VerifyTokenAccessResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.hasAccess === true) {
      writer.uint32(8).bool(message.hasAccess);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): VerifyTokenAccessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyTokenAccessResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.hasAccess = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyTokenAccessResponse {
    return {
      hasAccess: isSet(object.hasAccess) ? Boolean(object.hasAccess) : false,
    };
  },

  toJSON(message: VerifyTokenAccessResponse): unknown {
    const obj: any = {};
    message.hasAccess !== undefined && (obj.hasAccess = message.hasAccess);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VerifyTokenAccessResponse>, I>>(
    object: I
  ): VerifyTokenAccessResponse {
    const message = createBaseVerifyTokenAccessResponse();
    message.hasAccess = object.hasAccess ?? false;
    return message;
  },
};

/** Provides API to verify identity and determine access rights of the identity */
export interface IAMAuthService {
  /** Create access token and refresh token using password. Creates identity if not exist */
  CreateTokenWithPassword(
    request: CreateTokenWithPasswordRequest
  ): Promise<CreateTokenWithPasswordResponse>;
  /** Creates new access token using refresh token */
  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse>;
  /** Invalidates pare of access token and refresh tokens */
  InvalidateToken(
    request: InvalidateTokenRequest
  ): Promise<InvalidateTokenResponse>;
  /** Verifies if token can perform actions on the resources */
  VerifyTokenAccess(
    request: VerifyTokenAccessRequest
  ): Promise<VerifyTokenAccessResponse>;
}

export class IAMAuthServiceClientImpl implements IAMAuthService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.CreateTokenWithPassword = this.CreateTokenWithPassword.bind(this);
    this.RefreshToken = this.RefreshToken.bind(this);
    this.InvalidateToken = this.InvalidateToken.bind(this);
    this.VerifyTokenAccess = this.VerifyTokenAccess.bind(this);
  }
  CreateTokenWithPassword(
    request: CreateTokenWithPasswordRequest
  ): Promise<CreateTokenWithPasswordResponse> {
    const data = CreateTokenWithPasswordRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
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
      "native_iam_auth.IAMAuthService",
      "RefreshToken",
      data
    );
    return promise.then((data) =>
      RefreshTokenResponse.decode(new _m0.Reader(data))
    );
  }

  InvalidateToken(
    request: InvalidateTokenRequest
  ): Promise<InvalidateTokenResponse> {
    const data = InvalidateTokenRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "InvalidateToken",
      data
    );
    return promise.then((data) =>
      InvalidateTokenResponse.decode(new _m0.Reader(data))
    );
  }

  VerifyTokenAccess(
    request: VerifyTokenAccessRequest
  ): Promise<VerifyTokenAccessResponse> {
    const data = VerifyTokenAccessRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "VerifyTokenAccess",
      data
    );
    return promise.then((data) =>
      VerifyTokenAccessResponse.decode(new _m0.Reader(data))
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
