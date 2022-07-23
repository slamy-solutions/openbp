/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_auth";

export interface AuthenticationResponse {
  /** 2FA information if it ise required */
  twoFactorAuth: AuthenticationResponse_TwoFactorAuthorization | undefined;
  /** Authorization data if 2FA is not required */
  authData: AuthenticationResponse_AuthData | undefined;
  /** Identity UUID */
  identity: string;
}

export interface AuthenticationResponse_TwoFactorAuthorization {
  /** 2FA method */
  method: string;
  /** Token used for 2FA */
  token: string;
}

export interface AuthenticationResponse_AuthData {
  accessToken: string;
  refreshToken: string;
}

export interface LoginWithPasswordRequest {
  /** Identity email to autheticate */
  email: string;
  /** Identity password */
  password: string;
}

export interface LoginWithPasswordResponse {
  authData: AuthenticationResponse | undefined;
}

export interface LoginWithOAuth2Request {
  /** OAuth2 provider */
  provider: string;
  /** Token issued by OAuth2 provider`` */
  token: string;
}

export interface LoginWithOAuth2Response {
  authData: AuthenticationResponse | undefined;
}

export interface CompleteTwoFactorTOTPRequest {
  token: string;
  totpKey: string;
}

export interface CompleteTwoFactorTOTPResponse {
  accessToken: string;
  refreshToken: string;
}

export interface RefreshTokenRequest {
  refreshToken: string;
}

export interface RefreshTokenResponse {
  /** New access token */
  accessToken: string;
}

export interface InvalidateTokenRequest {
  /** Refresh or access token to invalidate. Both token will be invalidated */
  token: string;
}

export interface InvalidateTokenResponse {}

export interface VerifyAccessRequest {
  /** Token to verify */
  accessToken: string;
  /** What to verify */
  policies: string;
}

export interface VerifyAccessRequest_VerifyPolicy {
  /** Namespace where to verify. Namespace can be empty for global policy. */
  namespace: string;
  /** List of privileges to verify */
  privileges: string[];
}

export interface VerifyAccessResponse {
  hasAccess: boolean;
}

function createBaseAuthenticationResponse(): AuthenticationResponse {
  return { twoFactorAuth: undefined, authData: undefined, identity: "" };
}

export const AuthenticationResponse = {
  encode(
    message: AuthenticationResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.twoFactorAuth !== undefined) {
      AuthenticationResponse_TwoFactorAuthorization.encode(
        message.twoFactorAuth,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.authData !== undefined) {
      AuthenticationResponse_AuthData.encode(
        message.authData,
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.identity !== "") {
      writer.uint32(26).string(message.identity);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AuthenticationResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticationResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.twoFactorAuth =
            AuthenticationResponse_TwoFactorAuthorization.decode(
              reader,
              reader.uint32()
            );
          break;
        case 2:
          message.authData = AuthenticationResponse_AuthData.decode(
            reader,
            reader.uint32()
          );
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

  fromJSON(object: any): AuthenticationResponse {
    return {
      twoFactorAuth: isSet(object.twoFactorAuth)
        ? AuthenticationResponse_TwoFactorAuthorization.fromJSON(
            object.twoFactorAuth
          )
        : undefined,
      authData: isSet(object.authData)
        ? AuthenticationResponse_AuthData.fromJSON(object.authData)
        : undefined,
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: AuthenticationResponse): unknown {
    const obj: any = {};
    message.twoFactorAuth !== undefined &&
      (obj.twoFactorAuth = message.twoFactorAuth
        ? AuthenticationResponse_TwoFactorAuthorization.toJSON(
            message.twoFactorAuth
          )
        : undefined);
    message.authData !== undefined &&
      (obj.authData = message.authData
        ? AuthenticationResponse_AuthData.toJSON(message.authData)
        : undefined);
    message.identity !== undefined && (obj.identity = message.identity);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthenticationResponse>, I>>(
    object: I
  ): AuthenticationResponse {
    const message = createBaseAuthenticationResponse();
    message.twoFactorAuth =
      object.twoFactorAuth !== undefined && object.twoFactorAuth !== null
        ? AuthenticationResponse_TwoFactorAuthorization.fromPartial(
            object.twoFactorAuth
          )
        : undefined;
    message.authData =
      object.authData !== undefined && object.authData !== null
        ? AuthenticationResponse_AuthData.fromPartial(object.authData)
        : undefined;
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseAuthenticationResponse_TwoFactorAuthorization(): AuthenticationResponse_TwoFactorAuthorization {
  return { method: "", token: "" };
}

export const AuthenticationResponse_TwoFactorAuthorization = {
  encode(
    message: AuthenticationResponse_TwoFactorAuthorization,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.method !== "") {
      writer.uint32(10).string(message.method);
    }
    if (message.token !== "") {
      writer.uint32(18).string(message.token);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AuthenticationResponse_TwoFactorAuthorization {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticationResponse_TwoFactorAuthorization();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.method = reader.string();
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

  fromJSON(object: any): AuthenticationResponse_TwoFactorAuthorization {
    return {
      method: isSet(object.method) ? String(object.method) : "",
      token: isSet(object.token) ? String(object.token) : "",
    };
  },

  toJSON(message: AuthenticationResponse_TwoFactorAuthorization): unknown {
    const obj: any = {};
    message.method !== undefined && (obj.method = message.method);
    message.token !== undefined && (obj.token = message.token);
    return obj;
  },

  fromPartial<
    I extends Exact<
      DeepPartial<AuthenticationResponse_TwoFactorAuthorization>,
      I
    >
  >(object: I): AuthenticationResponse_TwoFactorAuthorization {
    const message = createBaseAuthenticationResponse_TwoFactorAuthorization();
    message.method = object.method ?? "";
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseAuthenticationResponse_AuthData(): AuthenticationResponse_AuthData {
  return { accessToken: "", refreshToken: "" };
}

export const AuthenticationResponse_AuthData = {
  encode(
    message: AuthenticationResponse_AuthData,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.refreshToken !== "") {
      writer.uint32(18).string(message.refreshToken);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AuthenticationResponse_AuthData {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticationResponse_AuthData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.refreshToken = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AuthenticationResponse_AuthData {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
    };
  },

  toJSON(message: AuthenticationResponse_AuthData): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthenticationResponse_AuthData>, I>>(
    object: I
  ): AuthenticationResponse_AuthData {
    const message = createBaseAuthenticationResponse_AuthData();
    message.accessToken = object.accessToken ?? "";
    message.refreshToken = object.refreshToken ?? "";
    return message;
  },
};

function createBaseLoginWithPasswordRequest(): LoginWithPasswordRequest {
  return { email: "", password: "" };
}

export const LoginWithPasswordRequest = {
  encode(
    message: LoginWithPasswordRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.email !== "") {
      writer.uint32(10).string(message.email);
    }
    if (message.password !== "") {
      writer.uint32(18).string(message.password);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LoginWithPasswordRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginWithPasswordRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.email = reader.string();
          break;
        case 2:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginWithPasswordRequest {
    return {
      email: isSet(object.email) ? String(object.email) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: LoginWithPasswordRequest): unknown {
    const obj: any = {};
    message.email !== undefined && (obj.email = message.email);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<LoginWithPasswordRequest>, I>>(
    object: I
  ): LoginWithPasswordRequest {
    const message = createBaseLoginWithPasswordRequest();
    message.email = object.email ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseLoginWithPasswordResponse(): LoginWithPasswordResponse {
  return { authData: undefined };
}

export const LoginWithPasswordResponse = {
  encode(
    message: LoginWithPasswordResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.authData !== undefined) {
      AuthenticationResponse.encode(
        message.authData,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LoginWithPasswordResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginWithPasswordResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authData = AuthenticationResponse.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginWithPasswordResponse {
    return {
      authData: isSet(object.authData)
        ? AuthenticationResponse.fromJSON(object.authData)
        : undefined,
    };
  },

  toJSON(message: LoginWithPasswordResponse): unknown {
    const obj: any = {};
    message.authData !== undefined &&
      (obj.authData = message.authData
        ? AuthenticationResponse.toJSON(message.authData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<LoginWithPasswordResponse>, I>>(
    object: I
  ): LoginWithPasswordResponse {
    const message = createBaseLoginWithPasswordResponse();
    message.authData =
      object.authData !== undefined && object.authData !== null
        ? AuthenticationResponse.fromPartial(object.authData)
        : undefined;
    return message;
  },
};

function createBaseLoginWithOAuth2Request(): LoginWithOAuth2Request {
  return { provider: "", token: "" };
}

export const LoginWithOAuth2Request = {
  encode(
    message: LoginWithOAuth2Request,
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
  ): LoginWithOAuth2Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginWithOAuth2Request();
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

  fromJSON(object: any): LoginWithOAuth2Request {
    return {
      provider: isSet(object.provider) ? String(object.provider) : "",
      token: isSet(object.token) ? String(object.token) : "",
    };
  },

  toJSON(message: LoginWithOAuth2Request): unknown {
    const obj: any = {};
    message.provider !== undefined && (obj.provider = message.provider);
    message.token !== undefined && (obj.token = message.token);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<LoginWithOAuth2Request>, I>>(
    object: I
  ): LoginWithOAuth2Request {
    const message = createBaseLoginWithOAuth2Request();
    message.provider = object.provider ?? "";
    message.token = object.token ?? "";
    return message;
  },
};

function createBaseLoginWithOAuth2Response(): LoginWithOAuth2Response {
  return { authData: undefined };
}

export const LoginWithOAuth2Response = {
  encode(
    message: LoginWithOAuth2Response,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.authData !== undefined) {
      AuthenticationResponse.encode(
        message.authData,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): LoginWithOAuth2Response {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLoginWithOAuth2Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authData = AuthenticationResponse.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginWithOAuth2Response {
    return {
      authData: isSet(object.authData)
        ? AuthenticationResponse.fromJSON(object.authData)
        : undefined,
    };
  },

  toJSON(message: LoginWithOAuth2Response): unknown {
    const obj: any = {};
    message.authData !== undefined &&
      (obj.authData = message.authData
        ? AuthenticationResponse.toJSON(message.authData)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<LoginWithOAuth2Response>, I>>(
    object: I
  ): LoginWithOAuth2Response {
    const message = createBaseLoginWithOAuth2Response();
    message.authData =
      object.authData !== undefined && object.authData !== null
        ? AuthenticationResponse.fromPartial(object.authData)
        : undefined;
    return message;
  },
};

function createBaseCompleteTwoFactorTOTPRequest(): CompleteTwoFactorTOTPRequest {
  return { token: "", totpKey: "" };
}

export const CompleteTwoFactorTOTPRequest = {
  encode(
    message: CompleteTwoFactorTOTPRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.token !== "") {
      writer.uint32(10).string(message.token);
    }
    if (message.totpKey !== "") {
      writer.uint32(18).string(message.totpKey);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CompleteTwoFactorTOTPRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCompleteTwoFactorTOTPRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.token = reader.string();
          break;
        case 2:
          message.totpKey = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CompleteTwoFactorTOTPRequest {
    return {
      token: isSet(object.token) ? String(object.token) : "",
      totpKey: isSet(object.totpKey) ? String(object.totpKey) : "",
    };
  },

  toJSON(message: CompleteTwoFactorTOTPRequest): unknown {
    const obj: any = {};
    message.token !== undefined && (obj.token = message.token);
    message.totpKey !== undefined && (obj.totpKey = message.totpKey);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CompleteTwoFactorTOTPRequest>, I>>(
    object: I
  ): CompleteTwoFactorTOTPRequest {
    const message = createBaseCompleteTwoFactorTOTPRequest();
    message.token = object.token ?? "";
    message.totpKey = object.totpKey ?? "";
    return message;
  },
};

function createBaseCompleteTwoFactorTOTPResponse(): CompleteTwoFactorTOTPResponse {
  return { accessToken: "", refreshToken: "" };
}

export const CompleteTwoFactorTOTPResponse = {
  encode(
    message: CompleteTwoFactorTOTPResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.refreshToken !== "") {
      writer.uint32(18).string(message.refreshToken);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CompleteTwoFactorTOTPResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCompleteTwoFactorTOTPResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.refreshToken = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CompleteTwoFactorTOTPResponse {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      refreshToken: isSet(object.refreshToken)
        ? String(object.refreshToken)
        : "",
    };
  },

  toJSON(message: CompleteTwoFactorTOTPResponse): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    message.refreshToken !== undefined &&
      (obj.refreshToken = message.refreshToken);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CompleteTwoFactorTOTPResponse>, I>>(
    object: I
  ): CompleteTwoFactorTOTPResponse {
    const message = createBaseCompleteTwoFactorTOTPResponse();
    message.accessToken = object.accessToken ?? "";
    message.refreshToken = object.refreshToken ?? "";
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

function createBaseVerifyAccessRequest(): VerifyAccessRequest {
  return { accessToken: "", policies: "" };
}

export const VerifyAccessRequest = {
  encode(
    message: VerifyAccessRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessToken !== "") {
      writer.uint32(10).string(message.accessToken);
    }
    if (message.policies !== "") {
      writer.uint32(18).string(message.policies);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): VerifyAccessRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyAccessRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessToken = reader.string();
          break;
        case 2:
          message.policies = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyAccessRequest {
    return {
      accessToken: isSet(object.accessToken) ? String(object.accessToken) : "",
      policies: isSet(object.policies) ? String(object.policies) : "",
    };
  },

  toJSON(message: VerifyAccessRequest): unknown {
    const obj: any = {};
    message.accessToken !== undefined &&
      (obj.accessToken = message.accessToken);
    message.policies !== undefined && (obj.policies = message.policies);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VerifyAccessRequest>, I>>(
    object: I
  ): VerifyAccessRequest {
    const message = createBaseVerifyAccessRequest();
    message.accessToken = object.accessToken ?? "";
    message.policies = object.policies ?? "";
    return message;
  },
};

function createBaseVerifyAccessRequest_VerifyPolicy(): VerifyAccessRequest_VerifyPolicy {
  return { namespace: "", privileges: [] };
}

export const VerifyAccessRequest_VerifyPolicy = {
  encode(
    message: VerifyAccessRequest_VerifyPolicy,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    for (const v of message.privileges) {
      writer.uint32(34).string(v!);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): VerifyAccessRequest_VerifyPolicy {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyAccessRequest_VerifyPolicy();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 4:
          message.privileges.push(reader.string());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): VerifyAccessRequest_VerifyPolicy {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      privileges: Array.isArray(object?.privileges)
        ? object.privileges.map((e: any) => String(e))
        : [],
    };
  },

  toJSON(message: VerifyAccessRequest_VerifyPolicy): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    if (message.privileges) {
      obj.privileges = message.privileges.map((e) => e);
    } else {
      obj.privileges = [];
    }
    return obj;
  },

  fromPartial<
    I extends Exact<DeepPartial<VerifyAccessRequest_VerifyPolicy>, I>
  >(object: I): VerifyAccessRequest_VerifyPolicy {
    const message = createBaseVerifyAccessRequest_VerifyPolicy();
    message.namespace = object.namespace ?? "";
    message.privileges = object.privileges?.map((e) => e) || [];
    return message;
  },
};

function createBaseVerifyAccessResponse(): VerifyAccessResponse {
  return { hasAccess: false };
}

export const VerifyAccessResponse = {
  encode(
    message: VerifyAccessResponse,
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
  ): VerifyAccessResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseVerifyAccessResponse();
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

  fromJSON(object: any): VerifyAccessResponse {
    return {
      hasAccess: isSet(object.hasAccess) ? Boolean(object.hasAccess) : false,
    };
  },

  toJSON(message: VerifyAccessResponse): unknown {
    const obj: any = {};
    message.hasAccess !== undefined && (obj.hasAccess = message.hasAccess);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<VerifyAccessResponse>, I>>(
    object: I
  ): VerifyAccessResponse {
    const message = createBaseVerifyAccessResponse();
    message.hasAccess = object.hasAccess ?? false;
    return message;
  },
};

/** Provides API to verify identity and determine access rights of the identity */
export interface IAMAuthService {
  /** Create access token and refresh token using password. Creates identity if not exist */
  LoginWithPassword(
    request: LoginWithPasswordRequest
  ): Promise<LoginWithPasswordResponse>;
  /** Create access token and refresh token using thrid party OAuth2 provider. Creates identity if not exist */
  LoginWithOAuth2(
    request: LoginWithOAuth2Request
  ): Promise<LoginWithOAuth2Response>;
  /**
   * Create access token and refresh token using SSO (Single Sign On)
   * rpc LoginWithSSO() returns ();
   * Completes started two factor TOTP (Time-based one-time password) authetication and returns actual access asn refresh tokens
   */
  CompleteTwoFactorTOTP(
    request: CompleteTwoFactorTOTPRequest
  ): Promise<CompleteTwoFactorTOTPResponse>;
  /** Creates new access token using refresh token */
  RefreshToken(request: RefreshTokenRequest): Promise<RefreshTokenResponse>;
  /** Invalidates pare of access token and refresh tokens */
  InvalidateToken(
    request: InvalidateTokenRequest
  ): Promise<InvalidateTokenResponse>;
  /** Verifies if token has access to provided resources */
  VerifyAccess(request: VerifyAccessRequest): Promise<VerifyAccessResponse>;
}

export class IAMAuthServiceClientImpl implements IAMAuthService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.LoginWithPassword = this.LoginWithPassword.bind(this);
    this.LoginWithOAuth2 = this.LoginWithOAuth2.bind(this);
    this.CompleteTwoFactorTOTP = this.CompleteTwoFactorTOTP.bind(this);
    this.RefreshToken = this.RefreshToken.bind(this);
    this.InvalidateToken = this.InvalidateToken.bind(this);
    this.VerifyAccess = this.VerifyAccess.bind(this);
  }
  LoginWithPassword(
    request: LoginWithPasswordRequest
  ): Promise<LoginWithPasswordResponse> {
    const data = LoginWithPasswordRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "LoginWithPassword",
      data
    );
    return promise.then((data) =>
      LoginWithPasswordResponse.decode(new _m0.Reader(data))
    );
  }

  LoginWithOAuth2(
    request: LoginWithOAuth2Request
  ): Promise<LoginWithOAuth2Response> {
    const data = LoginWithOAuth2Request.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "LoginWithOAuth2",
      data
    );
    return promise.then((data) =>
      LoginWithOAuth2Response.decode(new _m0.Reader(data))
    );
  }

  CompleteTwoFactorTOTP(
    request: CompleteTwoFactorTOTPRequest
  ): Promise<CompleteTwoFactorTOTPResponse> {
    const data = CompleteTwoFactorTOTPRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "CompleteTwoFactorTOTP",
      data
    );
    return promise.then((data) =>
      CompleteTwoFactorTOTPResponse.decode(new _m0.Reader(data))
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

  VerifyAccess(request: VerifyAccessRequest): Promise<VerifyAccessResponse> {
    const data = VerifyAccessRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_auth.IAMAuthService",
      "VerifyAccess",
      data
    );
    return promise.then((data) =>
      VerifyAccessResponse.decode(new _m0.Reader(data))
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
