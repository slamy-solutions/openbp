/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_auth";

/** Scope of the requested access. Check native_iam_policy for more information. */
export interface Scope {
  /** Namespace where this scope applies */
  namespace: string;
  /** Resources that can be accessed using token */
  resources: string[];
  /** Actions that can be done on resources */
  actions: string[];
  /** If this scope applies to all namespaces */
  namespaceIndependent: boolean;
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

export function createTokenWithPasswordResponse_StatusFromJSON(object: any): CreateTokenWithPasswordResponse_Status {
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

export function createTokenWithPasswordResponse_StatusToJSON(object: CreateTokenWithPasswordResponse_Status): string {
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

export function refreshTokenResponse_StatusFromJSON(object: any): RefreshTokenResponse_Status {
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

export function refreshTokenResponse_StatusToJSON(object: RefreshTokenResponse_Status): string {
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

export interface CheckAccessWithTokenRequest {
  /** Token to verify */
  accessToken: string;
  /** Scopes for with to validate access */
  scopes: Scope[];
}

export interface CheckAccessWithTokenResponse {
  /** Status of the verification */
  status: CheckAccessWithTokenResponse_Status;
  /** Details of the status, that can be safelly returned and displayed to the requester */
  message: string;
  /** Namespace where token and identity are located */
  namespace: string;
  /** Unique token id */
  tokenUUID: string;
  /** Unique identity id */
  identityUUID: string;
}

export enum CheckAccessWithTokenResponse_Status {
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

export function checkAccessWithTokenResponse_StatusFromJSON(object: any): CheckAccessWithTokenResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CheckAccessWithTokenResponse_Status.OK;
    case 1:
    case "TOKEN_INVALID":
      return CheckAccessWithTokenResponse_Status.TOKEN_INVALID;
    case 2:
    case "TOKEN_NOT_FOUND":
      return CheckAccessWithTokenResponse_Status.TOKEN_NOT_FOUND;
    case 3:
    case "TOKEN_DISABLED":
      return CheckAccessWithTokenResponse_Status.TOKEN_DISABLED;
    case 4:
    case "TOKEN_EXPIRED":
      return CheckAccessWithTokenResponse_Status.TOKEN_EXPIRED;
    case 5:
    case "UNAUTHORIZED":
      return CheckAccessWithTokenResponse_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CheckAccessWithTokenResponse_Status.UNRECOGNIZED;
  }
}

export function checkAccessWithTokenResponse_StatusToJSON(object: CheckAccessWithTokenResponse_Status): string {
  switch (object) {
    case CheckAccessWithTokenResponse_Status.OK:
      return "OK";
    case CheckAccessWithTokenResponse_Status.TOKEN_INVALID:
      return "TOKEN_INVALID";
    case CheckAccessWithTokenResponse_Status.TOKEN_NOT_FOUND:
      return "TOKEN_NOT_FOUND";
    case CheckAccessWithTokenResponse_Status.TOKEN_DISABLED:
      return "TOKEN_DISABLED";
    case CheckAccessWithTokenResponse_Status.TOKEN_EXPIRED:
      return "TOKEN_EXPIRED";
    case CheckAccessWithTokenResponse_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CheckAccessWithTokenResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface CheckAccessWithPasswordRequest {
  /** Namespace where identity is located */
  namespace: string;
  /** Identity UUID inside namespace */
  identity: string;
  /** Identity secret key */
  password: string;
  /** Arbitrary metadata. For example MAC/IP/information of the actor/application/browser/machine that provided this indentity and password. The exact format of metadata is not defined, but JSON is suggested. */
  metadata: string;
  /** Scopes to check */
  scopes: Scope[];
}

export interface CheckAccessWithPasswordResponse {
  /** Status of the check */
  status: CheckAccessWithPasswordResponse_Status;
  /** Details of the status, that can be safelly returned and displayed to the requester */
  message: string;
}

export enum CheckAccessWithPasswordResponse_Status {
  /** OK - Provided identity with provided password is allows to access scopes */
  OK = 0,
  /** UNAUTHENTICATED - Identity or password doesnt match */
  UNAUTHENTICATED = 1,
  /** UNAUTHORIZED - Identity dont have enought priviliges to perform actions from provided scopes */
  UNAUTHORIZED = 5,
  UNRECOGNIZED = -1,
}

export function checkAccessWithPasswordResponse_StatusFromJSON(object: any): CheckAccessWithPasswordResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CheckAccessWithPasswordResponse_Status.OK;
    case 1:
    case "UNAUTHENTICATED":
      return CheckAccessWithPasswordResponse_Status.UNAUTHENTICATED;
    case 5:
    case "UNAUTHORIZED":
      return CheckAccessWithPasswordResponse_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CheckAccessWithPasswordResponse_Status.UNRECOGNIZED;
  }
}

export function checkAccessWithPasswordResponse_StatusToJSON(object: CheckAccessWithPasswordResponse_Status): string {
  switch (object) {
    case CheckAccessWithPasswordResponse_Status.OK:
      return "OK";
    case CheckAccessWithPasswordResponse_Status.UNAUTHENTICATED:
      return "UNAUTHENTICATED";
    case CheckAccessWithPasswordResponse_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CheckAccessWithPasswordResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface CheckAccessWithX509Request {
  /** X509 certificate in DER format */
  certificate: Uint8Array;
  /** Scopes to check */
  scopes: Scope[];
}

export interface CheckAccessWithX509Response {
  /** Status of the check */
  status: CheckAccessWithX509Response_Status;
  /** Details of the status, that can be safelly returned and displayed to the requester */
  message: string;
  /** Certificate information. Only available if status is one of the {OK; CERTIFICATE_DISABLED; IDENTITY_NOT_FOUND; IDENTITY_NOT_ACTIVE; UNAUTHORIZED} */
  certificateInfo: CheckAccessWithX509Response_CertificateInfo | undefined;
}

export enum CheckAccessWithX509Response_Status {
  /** OK - Provided identity with provided certificate is allows to access scopes */
  OK = 0,
  /** CERTIFICATE_INVALID_FORMAT - Certificate corrupted or was supplied not in the DER format */
  CERTIFICATE_INVALID_FORMAT = 1,
  /** CERTIFICATE_INVALID - Signature or other aspects of the certificate are invalid */
  CERTIFICATE_INVALID = 2,
  /** CERTIFICATE_NOT_FOUND - Certificate wasnt founded. Most probably certificate or entire namespace was deleted */
  CERTIFICATE_NOT_FOUND = 3,
  /** CERTIFICATE_DISABLED - Certificate was manually disable and cont be used in auth mechanisms */
  CERTIFICATE_DISABLED = 4,
  /** IDENTITY_NOT_FOUND - Identity wasnt founded. Most probably it was deleted and certificate will be deleted soon */
  IDENTITY_NOT_FOUND = 5,
  /** IDENTITY_NOT_ACTIVE - Identity was manually disabled. */
  IDENTITY_NOT_ACTIVE = 6,
  /** UNAUTHORIZED - Certificate is valid, but identity dont have enought priviliges to perform actions from provided scopes */
  UNAUTHORIZED = 7,
  UNRECOGNIZED = -1,
}

export function checkAccessWithX509Response_StatusFromJSON(object: any): CheckAccessWithX509Response_Status {
  switch (object) {
    case 0:
    case "OK":
      return CheckAccessWithX509Response_Status.OK;
    case 1:
    case "CERTIFICATE_INVALID_FORMAT":
      return CheckAccessWithX509Response_Status.CERTIFICATE_INVALID_FORMAT;
    case 2:
    case "CERTIFICATE_INVALID":
      return CheckAccessWithX509Response_Status.CERTIFICATE_INVALID;
    case 3:
    case "CERTIFICATE_NOT_FOUND":
      return CheckAccessWithX509Response_Status.CERTIFICATE_NOT_FOUND;
    case 4:
    case "CERTIFICATE_DISABLED":
      return CheckAccessWithX509Response_Status.CERTIFICATE_DISABLED;
    case 5:
    case "IDENTITY_NOT_FOUND":
      return CheckAccessWithX509Response_Status.IDENTITY_NOT_FOUND;
    case 6:
    case "IDENTITY_NOT_ACTIVE":
      return CheckAccessWithX509Response_Status.IDENTITY_NOT_ACTIVE;
    case 7:
    case "UNAUTHORIZED":
      return CheckAccessWithX509Response_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CheckAccessWithX509Response_Status.UNRECOGNIZED;
  }
}

export function checkAccessWithX509Response_StatusToJSON(object: CheckAccessWithX509Response_Status): string {
  switch (object) {
    case CheckAccessWithX509Response_Status.OK:
      return "OK";
    case CheckAccessWithX509Response_Status.CERTIFICATE_INVALID_FORMAT:
      return "CERTIFICATE_INVALID_FORMAT";
    case CheckAccessWithX509Response_Status.CERTIFICATE_INVALID:
      return "CERTIFICATE_INVALID";
    case CheckAccessWithX509Response_Status.CERTIFICATE_NOT_FOUND:
      return "CERTIFICATE_NOT_FOUND";
    case CheckAccessWithX509Response_Status.CERTIFICATE_DISABLED:
      return "CERTIFICATE_DISABLED";
    case CheckAccessWithX509Response_Status.IDENTITY_NOT_FOUND:
      return "IDENTITY_NOT_FOUND";
    case CheckAccessWithX509Response_Status.IDENTITY_NOT_ACTIVE:
      return "IDENTITY_NOT_ACTIVE";
    case CheckAccessWithX509Response_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CheckAccessWithX509Response_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

/** Detailed information about certificate */
export interface CheckAccessWithX509Response_CertificateInfo {
  /** Namespace where certificate and identity are located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
  /** Unique identifier of the identity */
  identity: string;
}

export interface CheckAccessRequest {
  /** Namespace where identity is located */
  namespace: string;
  /** Unique identifier of the identity */
  identity: string;
  /** Scopes to check */
  scopes: Scope[];
}

export interface CheckAccessResponse {
  /** Status of the check */
  status: CheckAccessResponse_Status;
  /** Details of the status, that can be safelly returned and displayed to the requester */
  message: string;
}

export enum CheckAccessResponse_Status {
  /** OK - Provided identity with provided certificate is allows to access scopes */
  OK = 0,
  /** IDENTITY_NOT_FOUND - Identity wasnt founded. */
  IDENTITY_NOT_FOUND = 1,
  /** IDENTITY_NOT_ACTIVE - Identity was manually disabled. */
  IDENTITY_NOT_ACTIVE = 2,
  /** UNAUTHORIZED - Identity dont have enought priviliges to perform actions from provided scopes */
  UNAUTHORIZED = 3,
  UNRECOGNIZED = -1,
}

export function checkAccessResponse_StatusFromJSON(object: any): CheckAccessResponse_Status {
  switch (object) {
    case 0:
    case "OK":
      return CheckAccessResponse_Status.OK;
    case 1:
    case "IDENTITY_NOT_FOUND":
      return CheckAccessResponse_Status.IDENTITY_NOT_FOUND;
    case 2:
    case "IDENTITY_NOT_ACTIVE":
      return CheckAccessResponse_Status.IDENTITY_NOT_ACTIVE;
    case 3:
    case "UNAUTHORIZED":
      return CheckAccessResponse_Status.UNAUTHORIZED;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CheckAccessResponse_Status.UNRECOGNIZED;
  }
}

export function checkAccessResponse_StatusToJSON(object: CheckAccessResponse_Status): string {
  switch (object) {
    case CheckAccessResponse_Status.OK:
      return "OK";
    case CheckAccessResponse_Status.IDENTITY_NOT_FOUND:
      return "IDENTITY_NOT_FOUND";
    case CheckAccessResponse_Status.IDENTITY_NOT_ACTIVE:
      return "IDENTITY_NOT_ACTIVE";
    case CheckAccessResponse_Status.UNAUTHORIZED:
      return "UNAUTHORIZED";
    case CheckAccessResponse_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
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

function createBaseCreateTokenWithPasswordRequest(): CreateTokenWithPasswordRequest {
  return { namespace: "", identity: "", password: "", metadata: "", scopes: [] };
}

export const CreateTokenWithPasswordRequest = {
  encode(message: CreateTokenWithPasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTokenWithPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithPasswordRequest();
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

          message.password = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.metadata = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithPasswordRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      password: isSet(object.password) ? String(object.password) : "",
      metadata: isSet(object.metadata) ? String(object.metadata) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
    };
  },

  toJSON(message: CreateTokenWithPasswordRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.password !== "") {
      obj.password = message.password;
    }
    if (message.metadata !== "") {
      obj.metadata = message.metadata;
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateTokenWithPasswordRequest>, I>>(base?: I): CreateTokenWithPasswordRequest {
    return CreateTokenWithPasswordRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateTokenWithPasswordRequest>, I>>(
    object: I,
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
  encode(message: CreateTokenWithPasswordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTokenWithPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithPasswordResponse();
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

          message.accessToken = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
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

  fromJSON(object: any): CreateTokenWithPasswordResponse {
    return {
      status: isSet(object.status) ? createTokenWithPasswordResponse_StatusFromJSON(object.status) : 0,
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "",
    };
  },

  toJSON(message: CreateTokenWithPasswordResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = createTokenWithPasswordResponse_StatusToJSON(message.status);
    }
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateTokenWithPasswordResponse>, I>>(base?: I): CreateTokenWithPasswordResponse {
    return CreateTokenWithPasswordResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateTokenWithPasswordResponse>, I>>(
    object: I,
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
  encode(message: CreateTokenWithOAuth2Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.provider !== "") {
      writer.uint32(10).string(message.provider);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTokenWithOAuth2Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithOAuth2Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.provider = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.token = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.provider !== "") {
      obj.provider = message.provider;
    }
    if (message.token !== "") {
      obj.token = message.token;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateTokenWithOAuth2Request>, I>>(base?: I): CreateTokenWithOAuth2Request {
    return CreateTokenWithOAuth2Request.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateTokenWithOAuth2Request>, I>>(object: I): CreateTokenWithOAuth2Request {
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
  encode(message: CreateTokenWithOAuth2Response, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateTokenWithOAuth2Response {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateTokenWithOAuth2Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.accessToken = reader.string();
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

          message.identity = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateTokenWithOAuth2Response {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: CreateTokenWithOAuth2Response): unknown {
    const obj: any = {};
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateTokenWithOAuth2Response>, I>>(base?: I): CreateTokenWithOAuth2Response {
    return CreateTokenWithOAuth2Response.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateTokenWithOAuth2Response>, I>>(
    object: I,
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
  encode(message: RefreshTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.refreshToken !== "") {
      writer.uint32(10).string(message.refreshToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshTokenRequest();
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

  fromJSON(object: any): RefreshTokenRequest {
    return { refreshToken: isSet(object.refreshToken) ? String(object.refreshToken) : "" };
  },

  toJSON(message: RefreshTokenRequest): unknown {
    const obj: any = {};
    if (message.refreshToken !== "") {
      obj.refreshToken = message.refreshToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshTokenRequest>, I>>(base?: I): RefreshTokenRequest {
    return RefreshTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshTokenRequest>, I>>(object: I): RefreshTokenRequest {
    const message = createBaseRefreshTokenRequest();
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseRefreshTokenResponse(): RefreshTokenResponse {
  return { status: 0, accessToken: "" };
}

export const RefreshTokenResponse = {
  encode(message: RefreshTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.accessToken !== "") {
      writer.uint32(18).string(message.accessToken);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RefreshTokenResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRefreshTokenResponse();
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

          message.accessToken = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RefreshTokenResponse {
    return {
      status: isSet(object.status) ? refreshTokenResponse_StatusFromJSON(object.status) : 0,
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
    };
  },

  toJSON(message: RefreshTokenResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = refreshTokenResponse_StatusToJSON(message.status);
    }
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RefreshTokenResponse>, I>>(base?: I): RefreshTokenResponse {
    return RefreshTokenResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RefreshTokenResponse>, I>>(object: I): RefreshTokenResponse {
    const message = createBaseRefreshTokenResponse();
    message.status = object.status ?? 0;
    message.accessToken = object.accessToken ?? "";
    return message;
  },
};

function createBaseCheckAccessWithTokenRequest(): CheckAccessWithTokenRequest {
  return { accessToken: "", scopes: [] };
}

export const CheckAccessWithTokenRequest = {
  encode(message: CheckAccessWithTokenRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithTokenRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithTokenRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.accessToken = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithTokenRequest {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
    };
  },

  toJSON(message: CheckAccessWithTokenRequest): unknown {
    const obj: any = {};
    if (message.accessToken !== "") {
      obj.accessToken = message.accessToken;
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithTokenRequest>, I>>(base?: I): CheckAccessWithTokenRequest {
    return CheckAccessWithTokenRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithTokenRequest>, I>>(object: I): CheckAccessWithTokenRequest {
    const message = createBaseCheckAccessWithTokenRequest();
    message.accessToken = object.accessToken ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCheckAccessWithTokenResponse(): CheckAccessWithTokenResponse {
  return { status: 0, message: "", namespace: "", tokenUUID: "", identityUUID: "" };
}

export const CheckAccessWithTokenResponse = {
  encode(message: CheckAccessWithTokenResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    if (message.namespace !== "") {
      writer.uint32(26).string(message.namespace);
    }
    if (message.tokenUUID !== "") {
      writer.uint32(34).string(message.tokenUUID);
    }
    if (message.identityUUID !== "") {
      writer.uint32(42).string(message.identityUUID);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithTokenResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithTokenResponse();
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

          message.message = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.tokenUUID = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
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

  fromJSON(object: any): CheckAccessWithTokenResponse {
    return {
      status: isSet(object.status) ? checkAccessWithTokenResponse_StatusFromJSON(object.status) : 0,
      message: isSet(object.message) ? String(object.message) : "",
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      tokenUUID: isSet(object.tokenUUID) ? String(object.tokenUUID) : "",
      identityUUID: isSet(object.identityUUID) ? String(object.identityUUID) : "",
    };
  },

  toJSON(message: CheckAccessWithTokenResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = checkAccessWithTokenResponse_StatusToJSON(message.status);
    }
    if (message.message !== "") {
      obj.message = message.message;
    }
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.tokenUUID !== "") {
      obj.tokenUUID = message.tokenUUID;
    }
    if (message.identityUUID !== "") {
      obj.identityUUID = message.identityUUID;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithTokenResponse>, I>>(base?: I): CheckAccessWithTokenResponse {
    return CheckAccessWithTokenResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithTokenResponse>, I>>(object: I): CheckAccessWithTokenResponse {
    const message = createBaseCheckAccessWithTokenResponse();
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    message.namespace = object.namespace ?? "";
    message.tokenUUID = object.tokenUUID ?? "";
    message.identityUUID = object.identityUUID ?? "";
    return message;
  },
};

function createBaseCheckAccessWithPasswordRequest(): CheckAccessWithPasswordRequest {
  return { namespace: "", identity: "", password: "", metadata: "", scopes: [] };
}

export const CheckAccessWithPasswordRequest = {
  encode(message: CheckAccessWithPasswordRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithPasswordRequest();
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

          message.password = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.metadata = reader.string();
          continue;
        case 5:
          if (tag !== 42) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithPasswordRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      password: isSet(object.password) ? String(object.password) : "",
      metadata: isSet(object.metadata) ? String(object.metadata) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
    };
  },

  toJSON(message: CheckAccessWithPasswordRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.password !== "") {
      obj.password = message.password;
    }
    if (message.metadata !== "") {
      obj.metadata = message.metadata;
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithPasswordRequest>, I>>(base?: I): CheckAccessWithPasswordRequest {
    return CheckAccessWithPasswordRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithPasswordRequest>, I>>(
    object: I,
  ): CheckAccessWithPasswordRequest {
    const message = createBaseCheckAccessWithPasswordRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.password = object.password ?? "";
    message.metadata = object.metadata ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCheckAccessWithPasswordResponse(): CheckAccessWithPasswordResponse {
  return { status: 0, message: "" };
}

export const CheckAccessWithPasswordResponse = {
  encode(message: CheckAccessWithPasswordResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithPasswordResponse();
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

          message.message = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithPasswordResponse {
    return {
      status: isSet(object.status) ? checkAccessWithPasswordResponse_StatusFromJSON(object.status) : 0,
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: CheckAccessWithPasswordResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = checkAccessWithPasswordResponse_StatusToJSON(message.status);
    }
    if (message.message !== "") {
      obj.message = message.message;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithPasswordResponse>, I>>(base?: I): CheckAccessWithPasswordResponse {
    return CheckAccessWithPasswordResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithPasswordResponse>, I>>(
    object: I,
  ): CheckAccessWithPasswordResponse {
    const message = createBaseCheckAccessWithPasswordResponse();
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    return message;
  },
};

function createBaseCheckAccessWithX509Request(): CheckAccessWithX509Request {
  return { certificate: new Uint8Array(0), scopes: [] };
}

export const CheckAccessWithX509Request = {
  encode(message: CheckAccessWithX509Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate.length !== 0) {
      writer.uint32(10).bytes(message.certificate);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithX509Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithX509Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.scopes.push(Scope.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithX509Request {
    return {
      certificate: isSet(object.certificate) ? bytesFromBase64(object.certificate) : new Uint8Array(0),
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
    };
  },

  toJSON(message: CheckAccessWithX509Request): unknown {
    const obj: any = {};
    if (message.certificate.length !== 0) {
      obj.certificate = base64FromBytes(message.certificate);
    }
    if (message.scopes?.length) {
      obj.scopes = message.scopes.map((e) => Scope.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithX509Request>, I>>(base?: I): CheckAccessWithX509Request {
    return CheckAccessWithX509Request.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithX509Request>, I>>(object: I): CheckAccessWithX509Request {
    const message = createBaseCheckAccessWithX509Request();
    message.certificate = object.certificate ?? new Uint8Array(0);
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCheckAccessWithX509Response(): CheckAccessWithX509Response {
  return { status: 0, message: "", certificateInfo: undefined };
}

export const CheckAccessWithX509Response = {
  encode(message: CheckAccessWithX509Response, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    if (message.certificateInfo !== undefined) {
      CheckAccessWithX509Response_CertificateInfo.encode(message.certificateInfo, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithX509Response {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithX509Response();
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

          message.message = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.certificateInfo = CheckAccessWithX509Response_CertificateInfo.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithX509Response {
    return {
      status: isSet(object.status) ? checkAccessWithX509Response_StatusFromJSON(object.status) : 0,
      message: isSet(object.message) ? String(object.message) : "",
      certificateInfo: isSet(object.certificateInfo)
        ? CheckAccessWithX509Response_CertificateInfo.fromJSON(object.certificateInfo)
        : undefined,
    };
  },

  toJSON(message: CheckAccessWithX509Response): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = checkAccessWithX509Response_StatusToJSON(message.status);
    }
    if (message.message !== "") {
      obj.message = message.message;
    }
    if (message.certificateInfo !== undefined) {
      obj.certificateInfo = CheckAccessWithX509Response_CertificateInfo.toJSON(message.certificateInfo);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithX509Response>, I>>(base?: I): CheckAccessWithX509Response {
    return CheckAccessWithX509Response.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithX509Response>, I>>(object: I): CheckAccessWithX509Response {
    const message = createBaseCheckAccessWithX509Response();
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    message.certificateInfo = (object.certificateInfo !== undefined && object.certificateInfo !== null)
      ? CheckAccessWithX509Response_CertificateInfo.fromPartial(object.certificateInfo)
      : undefined;
    return message;
  },
};

function createBaseCheckAccessWithX509Response_CertificateInfo(): CheckAccessWithX509Response_CertificateInfo {
  return { namespace: "", uuid: "", identity: "" };
}

export const CheckAccessWithX509Response_CertificateInfo = {
  encode(message: CheckAccessWithX509Response_CertificateInfo, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.identity !== "") {
      writer.uint32(26).string(message.identity);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessWithX509Response_CertificateInfo {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessWithX509Response_CertificateInfo();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessWithX509Response_CertificateInfo {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: CheckAccessWithX509Response_CertificateInfo): unknown {
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
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessWithX509Response_CertificateInfo>, I>>(
    base?: I,
  ): CheckAccessWithX509Response_CertificateInfo {
    return CheckAccessWithX509Response_CertificateInfo.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessWithX509Response_CertificateInfo>, I>>(
    object: I,
  ): CheckAccessWithX509Response_CertificateInfo {
    const message = createBaseCheckAccessWithX509Response_CertificateInfo();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseCheckAccessRequest(): CheckAccessRequest {
  return { namespace: "", identity: "", scopes: [] };
}

export const CheckAccessRequest = {
  encode(message: CheckAccessRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    for (const v of message.scopes) {
      Scope.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      scopes: Array.isArray(object?.scopes) ? object.scopes.map((e: any) => Scope.fromJSON(e)) : [],
    };
  },

  toJSON(message: CheckAccessRequest): unknown {
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
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessRequest>, I>>(base?: I): CheckAccessRequest {
    return CheckAccessRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessRequest>, I>>(object: I): CheckAccessRequest {
    const message = createBaseCheckAccessRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.scopes = object.scopes?.map((e) => Scope.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCheckAccessResponse(): CheckAccessResponse {
  return { status: 0, message: "" };
}

export const CheckAccessResponse = {
  encode(message: CheckAccessResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CheckAccessResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCheckAccessResponse();
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

          message.message = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CheckAccessResponse {
    return {
      status: isSet(object.status) ? checkAccessResponse_StatusFromJSON(object.status) : 0,
      message: isSet(object.message) ? String(object.message) : "",
    };
  },

  toJSON(message: CheckAccessResponse): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = checkAccessResponse_StatusToJSON(message.status);
    }
    if (message.message !== "") {
      obj.message = message.message;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CheckAccessResponse>, I>>(base?: I): CheckAccessResponse {
    return CheckAccessResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CheckAccessResponse>, I>>(object: I): CheckAccessResponse {
    const message = createBaseCheckAccessResponse();
    message.status = object.status ?? 0;
    message.message = object.message ?? "";
    return message;
  },
};

/** Provides API for Basic, X509 and OAuth ("Open Authorization") style access control */
export interface IAMAuthService {
  /** OAuth. Create access token and refresh token using password */
  CreateTokenWithPassword(request: CreateTokenWithPasswordRequest): Promise<CreateTokenWithPasswordResponse>;
  /** OAuth. Creates new access token using refresh tokenna */
  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse>;
  /**
   * rpc VerifyResoureAccess(VerifyResourceAccessRequest) returns (VerifyResourceAccessResponse);
   * OAuth. Check if token is allowed to perform actions from the specified scopes
   */
  CheckAccessWithToken(request: CheckAccessWithTokenRequest): Promise<CheckAccessWithTokenResponse>;
  /** Basic Auth. Check if provided identity with proposed password is allowed to perform actions from the provided scopes */
  CheckAccessWithPassword(request: CheckAccessWithPasswordRequest): Promise<CheckAccessWithPasswordResponse>;
  /** Authorization with X509 certificates. Check if provided identity identified by proposed certificate is allowed to perform actions from the provided scopes */
  CheckAccessWithX509(request: CheckAccessWithX509Request): Promise<CheckAccessWithX509Response>;
  /** Check if provided identity is allowed to perform actions from the provided scopes */
  CheckAccess(request: CheckAccessRequest): Promise<CheckAccessResponse>;
}

export const IAMAuthServiceServiceName = "native_iam_auth.IAMAuthService";
export class IAMAuthServiceClientImpl implements IAMAuthService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMAuthServiceServiceName;
    this.rpc = rpc;
    this.CreateTokenWithPassword = this.CreateTokenWithPassword.bind(this);
    this.RefreshToken = this.RefreshToken.bind(this);
    this.CheckAccessWithToken = this.CheckAccessWithToken.bind(this);
    this.CheckAccessWithPassword = this.CheckAccessWithPassword.bind(this);
    this.CheckAccessWithX509 = this.CheckAccessWithX509.bind(this);
    this.CheckAccess = this.CheckAccess.bind(this);
  }
  CreateTokenWithPassword(request: CreateTokenWithPasswordRequest): Promise<CreateTokenWithPasswordResponse> {
    const data = CreateTokenWithPasswordRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CreateTokenWithPassword", data);
    return promise.then((data) => CreateTokenWithPasswordResponse.decode(_m0.Reader.create(data)));
  }

  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse> {
    const data = RefreshTokenRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RefreshToken", data);
    return promise.then((data) => RefreshTokenResponse.decode(_m0.Reader.create(data)));
  }

  CheckAccessWithToken(request: CheckAccessWithTokenRequest): Promise<CheckAccessWithTokenResponse> {
    const data = CheckAccessWithTokenRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CheckAccessWithToken", data);
    return promise.then((data) => CheckAccessWithTokenResponse.decode(_m0.Reader.create(data)));
  }

  CheckAccessWithPassword(request: CheckAccessWithPasswordRequest): Promise<CheckAccessWithPasswordResponse> {
    const data = CheckAccessWithPasswordRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CheckAccessWithPassword", data);
    return promise.then((data) => CheckAccessWithPasswordResponse.decode(_m0.Reader.create(data)));
  }

  CheckAccessWithX509(request: CheckAccessWithX509Request): Promise<CheckAccessWithX509Response> {
    const data = CheckAccessWithX509Request.encode(request).finish();
    const promise = this.rpc.request(this.service, "CheckAccessWithX509", data);
    return promise.then((data) => CheckAccessWithX509Response.decode(_m0.Reader.create(data)));
  }

  CheckAccess(request: CheckAccessRequest): Promise<CheckAccessResponse> {
    const data = CheckAccessRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CheckAccess", data);
    return promise.then((data) => CheckAccessResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
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

function bytesFromBase64(b64: string): Uint8Array {
  if (tsProtoGlobalThis.Buffer) {
    return Uint8Array.from(tsProtoGlobalThis.Buffer.from(b64, "base64"));
  } else {
    const bin = tsProtoGlobalThis.atob(b64);
    const arr = new Uint8Array(bin.length);
    for (let i = 0; i < bin.length; ++i) {
      arr[i] = bin.charCodeAt(i);
    }
    return arr;
  }
}

function base64FromBytes(arr: Uint8Array): string {
  if (tsProtoGlobalThis.Buffer) {
    return tsProtoGlobalThis.Buffer.from(arr).toString("base64");
  } else {
    const bin: string[] = [];
    arr.forEach((byte) => {
      bin.push(String.fromCharCode(byte));
    });
    return tsProtoGlobalThis.btoa(bin.join(""));
  }
}

type Builtin = Date | Function | Uint8Array | string | number | boolean | undefined;

export type DeepPartial<T> = T extends Builtin ? T
  : T extends Array<infer U> ? Array<DeepPartial<U>> : T extends ReadonlyArray<infer U> ? ReadonlyArray<DeepPartial<U>>
  : T extends {} ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

type KeysOfUnion<T> = T extends T ? keyof T : never;
export type Exact<P, I extends P> = P extends Builtin ? P
  : P & { [K in keyof P]: Exact<P[K], I[K]> } & { [K in Exclude<keyof I, KeysOfUnion<P>>]: never };

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
