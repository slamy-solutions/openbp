/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_configuration";

/** Configuration of the IAM */
export interface Configuration {
  /** Time to live of access token in milliseconds */
  accessTokenTTL: number;
  /** Time to live ot refresh token in milliseconds */
  refreshTokenTTL: number;
  /** Password authentication configuration */
  passwordAuth: Configuration_PasswordAuth | undefined;
  /** Google oauth2 configuration */
  googleOAuth2: Configuration_OAuth2 | undefined;
  /** Facebook oauth2 configuration */
  facebookOAuth2: Configuration_OAuth2 | undefined;
  /** Github oauth2 configuration */
  githubOAuth2: Configuration_OAuth2 | undefined;
  /** Github oauth2 configuration */
  gitlabOAuth2: Configuration_OAuth2 | undefined;
}

/** Configuration of specific OAuth2 provider */
export interface Configuration_OAuth2 {
  /** Enable or disable this provider of OAuth2 */
  enabled: boolean;
  /** OAuth2 client ID */
  clientId: string;
  /** OAuth2 client secret */
  clientSecret: string;
  /** Allow registration using this OAuth2 provider */
  allowRegistration: boolean;
}

export interface Configuration_PasswordAuth {
  /** Allow password authorization or not */
  enabled: boolean;
  /** Allow registration using password method */
  allowRegistration: boolean;
}

export interface GetConfigRequest {
  /** Use cache or not. Cache have a very low chance to be invalid. Cache invalidates after short period of thime (60 seconds). Cache can only be invalid on multiple simultanious read and writes. Its safe to use cache in most of the cases. */
  useCache: boolean;
}

export interface GetConfigresponse {
  /** Current configuration */
  configuration: Configuration | undefined;
}

export interface SetConfigRequest {
  /** Configuration to set */
  configuration: Configuration | undefined;
}

export interface SetConfigResponse {}

function createBaseConfiguration(): Configuration {
  return {
    accessTokenTTL: 0,
    refreshTokenTTL: 0,
    passwordAuth: undefined,
    googleOAuth2: undefined,
    facebookOAuth2: undefined,
    githubOAuth2: undefined,
    gitlabOAuth2: undefined,
  };
}

export const Configuration = {
  encode(
    message: Configuration,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.accessTokenTTL !== 0) {
      writer.uint32(8).uint32(message.accessTokenTTL);
    }
    if (message.refreshTokenTTL !== 0) {
      writer.uint32(16).uint32(message.refreshTokenTTL);
    }
    if (message.passwordAuth !== undefined) {
      Configuration_PasswordAuth.encode(
        message.passwordAuth,
        writer.uint32(82).fork()
      ).ldelim();
    }
    if (message.googleOAuth2 !== undefined) {
      Configuration_OAuth2.encode(
        message.googleOAuth2,
        writer.uint32(90).fork()
      ).ldelim();
    }
    if (message.facebookOAuth2 !== undefined) {
      Configuration_OAuth2.encode(
        message.facebookOAuth2,
        writer.uint32(98).fork()
      ).ldelim();
    }
    if (message.githubOAuth2 !== undefined) {
      Configuration_OAuth2.encode(
        message.githubOAuth2,
        writer.uint32(106).fork()
      ).ldelim();
    }
    if (message.gitlabOAuth2 !== undefined) {
      Configuration_OAuth2.encode(
        message.gitlabOAuth2,
        writer.uint32(114).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Configuration {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfiguration();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.accessTokenTTL = reader.uint32();
          break;
        case 2:
          message.refreshTokenTTL = reader.uint32();
          break;
        case 10:
          message.passwordAuth = Configuration_PasswordAuth.decode(
            reader,
            reader.uint32()
          );
          break;
        case 11:
          message.googleOAuth2 = Configuration_OAuth2.decode(
            reader,
            reader.uint32()
          );
          break;
        case 12:
          message.facebookOAuth2 = Configuration_OAuth2.decode(
            reader,
            reader.uint32()
          );
          break;
        case 13:
          message.githubOAuth2 = Configuration_OAuth2.decode(
            reader,
            reader.uint32()
          );
          break;
        case 14:
          message.gitlabOAuth2 = Configuration_OAuth2.decode(
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

  fromJSON(object: any): Configuration {
    return {
      accessTokenTTL: isSet(object.accessTokenTTL)
        ? Number(object.accessTokenTTL)
        : 0,
      refreshTokenTTL: isSet(object.refreshTokenTTL)
        ? Number(object.refreshTokenTTL)
        : 0,
      passwordAuth: isSet(object.passwordAuth)
        ? Configuration_PasswordAuth.fromJSON(object.passwordAuth)
        : undefined,
      googleOAuth2: isSet(object.googleOAuth2)
        ? Configuration_OAuth2.fromJSON(object.googleOAuth2)
        : undefined,
      facebookOAuth2: isSet(object.facebookOAuth2)
        ? Configuration_OAuth2.fromJSON(object.facebookOAuth2)
        : undefined,
      githubOAuth2: isSet(object.githubOAuth2)
        ? Configuration_OAuth2.fromJSON(object.githubOAuth2)
        : undefined,
      gitlabOAuth2: isSet(object.gitlabOAuth2)
        ? Configuration_OAuth2.fromJSON(object.gitlabOAuth2)
        : undefined,
    };
  },

  toJSON(message: Configuration): unknown {
    const obj: any = {};
    message.accessTokenTTL !== undefined &&
      (obj.accessTokenTTL = Math.round(message.accessTokenTTL));
    message.refreshTokenTTL !== undefined &&
      (obj.refreshTokenTTL = Math.round(message.refreshTokenTTL));
    message.passwordAuth !== undefined &&
      (obj.passwordAuth = message.passwordAuth
        ? Configuration_PasswordAuth.toJSON(message.passwordAuth)
        : undefined);
    message.googleOAuth2 !== undefined &&
      (obj.googleOAuth2 = message.googleOAuth2
        ? Configuration_OAuth2.toJSON(message.googleOAuth2)
        : undefined);
    message.facebookOAuth2 !== undefined &&
      (obj.facebookOAuth2 = message.facebookOAuth2
        ? Configuration_OAuth2.toJSON(message.facebookOAuth2)
        : undefined);
    message.githubOAuth2 !== undefined &&
      (obj.githubOAuth2 = message.githubOAuth2
        ? Configuration_OAuth2.toJSON(message.githubOAuth2)
        : undefined);
    message.gitlabOAuth2 !== undefined &&
      (obj.gitlabOAuth2 = message.gitlabOAuth2
        ? Configuration_OAuth2.toJSON(message.gitlabOAuth2)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Configuration>, I>>(
    object: I
  ): Configuration {
    const message = createBaseConfiguration();
    message.accessTokenTTL = object.accessTokenTTL ?? 0;
    message.refreshTokenTTL = object.refreshTokenTTL ?? 0;
    message.passwordAuth =
      object.passwordAuth !== undefined && object.passwordAuth !== null
        ? Configuration_PasswordAuth.fromPartial(object.passwordAuth)
        : undefined;
    message.googleOAuth2 =
      object.googleOAuth2 !== undefined && object.googleOAuth2 !== null
        ? Configuration_OAuth2.fromPartial(object.googleOAuth2)
        : undefined;
    message.facebookOAuth2 =
      object.facebookOAuth2 !== undefined && object.facebookOAuth2 !== null
        ? Configuration_OAuth2.fromPartial(object.facebookOAuth2)
        : undefined;
    message.githubOAuth2 =
      object.githubOAuth2 !== undefined && object.githubOAuth2 !== null
        ? Configuration_OAuth2.fromPartial(object.githubOAuth2)
        : undefined;
    message.gitlabOAuth2 =
      object.gitlabOAuth2 !== undefined && object.gitlabOAuth2 !== null
        ? Configuration_OAuth2.fromPartial(object.gitlabOAuth2)
        : undefined;
    return message;
  },
};

function createBaseConfiguration_OAuth2(): Configuration_OAuth2 {
  return {
    enabled: false,
    clientId: "",
    clientSecret: "",
    allowRegistration: false,
  };
}

export const Configuration_OAuth2 = {
  encode(
    message: Configuration_OAuth2,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }
    if (message.clientId !== "") {
      writer.uint32(18).string(message.clientId);
    }
    if (message.clientSecret !== "") {
      writer.uint32(26).string(message.clientSecret);
    }
    if (message.allowRegistration === true) {
      writer.uint32(32).bool(message.allowRegistration);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): Configuration_OAuth2 {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfiguration_OAuth2();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.enabled = reader.bool();
          break;
        case 2:
          message.clientId = reader.string();
          break;
        case 3:
          message.clientSecret = reader.string();
          break;
        case 4:
          message.allowRegistration = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Configuration_OAuth2 {
    return {
      enabled: isSet(object.enabled) ? Boolean(object.enabled) : false,
      clientId: isSet(object.clientId) ? String(object.clientId) : "",
      clientSecret: isSet(object.clientSecret)
        ? String(object.clientSecret)
        : "",
      allowRegistration: isSet(object.allowRegistration)
        ? Boolean(object.allowRegistration)
        : false,
    };
  },

  toJSON(message: Configuration_OAuth2): unknown {
    const obj: any = {};
    message.enabled !== undefined && (obj.enabled = message.enabled);
    message.clientId !== undefined && (obj.clientId = message.clientId);
    message.clientSecret !== undefined &&
      (obj.clientSecret = message.clientSecret);
    message.allowRegistration !== undefined &&
      (obj.allowRegistration = message.allowRegistration);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Configuration_OAuth2>, I>>(
    object: I
  ): Configuration_OAuth2 {
    const message = createBaseConfiguration_OAuth2();
    message.enabled = object.enabled ?? false;
    message.clientId = object.clientId ?? "";
    message.clientSecret = object.clientSecret ?? "";
    message.allowRegistration = object.allowRegistration ?? false;
    return message;
  },
};

function createBaseConfiguration_PasswordAuth(): Configuration_PasswordAuth {
  return { enabled: false, allowRegistration: false };
}

export const Configuration_PasswordAuth = {
  encode(
    message: Configuration_PasswordAuth,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.enabled === true) {
      writer.uint32(8).bool(message.enabled);
    }
    if (message.allowRegistration === true) {
      writer.uint32(16).bool(message.allowRegistration);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): Configuration_PasswordAuth {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseConfiguration_PasswordAuth();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.enabled = reader.bool();
          break;
        case 2:
          message.allowRegistration = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Configuration_PasswordAuth {
    return {
      enabled: isSet(object.enabled) ? Boolean(object.enabled) : false,
      allowRegistration: isSet(object.allowRegistration)
        ? Boolean(object.allowRegistration)
        : false,
    };
  },

  toJSON(message: Configuration_PasswordAuth): unknown {
    const obj: any = {};
    message.enabled !== undefined && (obj.enabled = message.enabled);
    message.allowRegistration !== undefined &&
      (obj.allowRegistration = message.allowRegistration);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Configuration_PasswordAuth>, I>>(
    object: I
  ): Configuration_PasswordAuth {
    const message = createBaseConfiguration_PasswordAuth();
    message.enabled = object.enabled ?? false;
    message.allowRegistration = object.allowRegistration ?? false;
    return message;
  },
};

function createBaseGetConfigRequest(): GetConfigRequest {
  return { useCache: false };
}

export const GetConfigRequest = {
  encode(
    message: GetConfigRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.useCache === true) {
      writer.uint32(8).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetConfigRequest {
    return {
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetConfigRequest): unknown {
    const obj: any = {};
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetConfigRequest>, I>>(
    object: I
  ): GetConfigRequest {
    const message = createBaseGetConfigRequest();
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetConfigresponse(): GetConfigresponse {
  return { configuration: undefined };
}

export const GetConfigresponse = {
  encode(
    message: GetConfigresponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.configuration !== undefined) {
      Configuration.encode(
        message.configuration,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetConfigresponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetConfigresponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.configuration = Configuration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetConfigresponse {
    return {
      configuration: isSet(object.configuration)
        ? Configuration.fromJSON(object.configuration)
        : undefined,
    };
  },

  toJSON(message: GetConfigresponse): unknown {
    const obj: any = {};
    message.configuration !== undefined &&
      (obj.configuration = message.configuration
        ? Configuration.toJSON(message.configuration)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetConfigresponse>, I>>(
    object: I
  ): GetConfigresponse {
    const message = createBaseGetConfigresponse();
    message.configuration =
      object.configuration !== undefined && object.configuration !== null
        ? Configuration.fromPartial(object.configuration)
        : undefined;
    return message;
  },
};

function createBaseSetConfigRequest(): SetConfigRequest {
  return { configuration: undefined };
}

export const SetConfigRequest = {
  encode(
    message: SetConfigRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.configuration !== undefined) {
      Configuration.encode(
        message.configuration,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetConfigRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetConfigRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.configuration = Configuration.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): SetConfigRequest {
    return {
      configuration: isSet(object.configuration)
        ? Configuration.fromJSON(object.configuration)
        : undefined,
    };
  },

  toJSON(message: SetConfigRequest): unknown {
    const obj: any = {};
    message.configuration !== undefined &&
      (obj.configuration = message.configuration
        ? Configuration.toJSON(message.configuration)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetConfigRequest>, I>>(
    object: I
  ): SetConfigRequest {
    const message = createBaseSetConfigRequest();
    message.configuration =
      object.configuration !== undefined && object.configuration !== null
        ? Configuration.fromPartial(object.configuration)
        : undefined;
    return message;
  },
};

function createBaseSetConfigResponse(): SetConfigResponse {
  return {};
}

export const SetConfigResponse = {
  encode(
    _: SetConfigResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SetConfigResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSetConfigResponse();
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

  fromJSON(_: any): SetConfigResponse {
    return {};
  },

  toJSON(_: SetConfigResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<SetConfigResponse>, I>>(
    _: I
  ): SetConfigResponse {
    const message = createBaseSetConfigResponse();
    return message;
  },
};

/** Provides general configuration API for IAM */
export interface IAMConfigService {
  Get(request: GetConfigRequest): Promise<GetConfigresponse>;
  Set(request: SetConfigRequest): Promise<SetConfigResponse>;
}

export class IAMConfigServiceClientImpl implements IAMConfigService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Get = this.Get.bind(this);
    this.Set = this.Set.bind(this);
  }
  Get(request: GetConfigRequest): Promise<GetConfigresponse> {
    const data = GetConfigRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_configuration.IAMConfigService",
      "Get",
      data
    );
    return promise.then((data) =>
      GetConfigresponse.decode(new _m0.Reader(data))
    );
  }

  Set(request: SetConfigRequest): Promise<SetConfigResponse> {
    const data = SetConfigRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_configuration.IAMConfigService",
      "Set",
      data
    );
    return promise.then((data) =>
      SetConfigResponse.decode(new _m0.Reader(data))
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
