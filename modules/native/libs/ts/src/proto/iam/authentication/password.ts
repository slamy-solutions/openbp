/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_iam_authentication_password";

export interface AuthenticateRequest {
  namespace: string;
  identity: string;
  password: string;
}

export interface AuthenticateResponse {
  authenticated: boolean;
}

export interface CreateOrUpdateRequest {
  namespace: string;
  identity: string;
  password: string;
}

export interface CreateOrUpdateResponse {}

export interface DeleteRequest {
  namespace: string;
  identity: string;
}

export interface DeleteResponse {}

export interface ExistRequest {
  namespace: string;
  identity: string;
}

export interface ExistResponse {
  exist: boolean;
}

function createBaseAuthenticateRequest(): AuthenticateRequest {
  return { namespace: "", identity: "", password: "" };
}

export const AuthenticateRequest = {
  encode(
    message: AuthenticateRequest,
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
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticateRequest();
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AuthenticateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: AuthenticateRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthenticateRequest>, I>>(
    object: I
  ): AuthenticateRequest {
    const message = createBaseAuthenticateRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseAuthenticateResponse(): AuthenticateResponse {
  return { authenticated: false };
}

export const AuthenticateResponse = {
  encode(
    message: AuthenticateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.authenticated === true) {
      writer.uint32(8).bool(message.authenticated);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AuthenticateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authenticated = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AuthenticateResponse {
    return {
      authenticated: isSet(object.authenticated)
        ? Boolean(object.authenticated)
        : false,
    };
  },

  toJSON(message: AuthenticateResponse): unknown {
    const obj: any = {};
    message.authenticated !== undefined &&
      (obj.authenticated = message.authenticated);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AuthenticateResponse>, I>>(
    object: I
  ): AuthenticateResponse {
    const message = createBaseAuthenticateResponse();
    message.authenticated = object.authenticated ?? false;
    return message;
  },
};

function createBaseCreateOrUpdateRequest(): CreateOrUpdateRequest {
  return { namespace: "", identity: "", password: "" };
}

export const CreateOrUpdateRequest = {
  encode(
    message: CreateOrUpdateRequest,
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
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateOrUpdateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateOrUpdateRequest();
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
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateOrUpdateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      password: isSet(object.password) ? String(object.password) : "",
    };
  },

  toJSON(message: CreateOrUpdateRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateOrUpdateRequest>, I>>(
    object: I
  ): CreateOrUpdateRequest {
    const message = createBaseCreateOrUpdateRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseCreateOrUpdateResponse(): CreateOrUpdateResponse {
  return {};
}

export const CreateOrUpdateResponse = {
  encode(
    _: CreateOrUpdateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateOrUpdateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateOrUpdateResponse();
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

  fromJSON(_: any): CreateOrUpdateResponse {
    return {};
  },

  toJSON(_: CreateOrUpdateResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateOrUpdateResponse>, I>>(
    _: I
  ): CreateOrUpdateResponse {
    const message = createBaseCreateOrUpdateResponse();
    return message;
  },
};

function createBaseDeleteRequest(): DeleteRequest {
  return { namespace: "", identity: "" };
}

export const DeleteRequest = {
  encode(
    message: DeleteRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
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
          message.namespace = reader.string();
          break;
        case 2:
          message.identity = reader.string();
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
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: DeleteRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteRequest>, I>>(
    object: I
  ): DeleteRequest {
    const message = createBaseDeleteRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
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

function createBaseExistRequest(): ExistRequest {
  return { namespace: "", identity: "" };
}

export const ExistRequest = {
  encode(
    message: ExistRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.identity = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExistRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: ExistRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.identity !== undefined && (obj.identity = message.identity);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExistRequest>, I>>(
    object: I
  ): ExistRequest {
    const message = createBaseExistRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseExistResponse(): ExistResponse {
  return { exist: false };
}

export const ExistResponse = {
  encode(
    message: ExistResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.exist = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExistResponse {
    return {
      exist: isSet(object.exist) ? Boolean(object.exist) : false,
    };
  },

  toJSON(message: ExistResponse): unknown {
    const obj: any = {};
    message.exist !== undefined && (obj.exist = message.exist);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExistResponse>, I>>(
    object: I
  ): ExistResponse {
    const message = createBaseExistResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

/** Provides API to authentificate identities using password */
export interface IAMAuthenticationPasswordService {
  /** Tries to find identity and compare its password. */
  Authenticate(request: AuthenticateRequest): Promise<AuthenticateResponse>;
  /** Creates or updates identity password for authentification. */
  CreateOrUpdate(
    request: CreateOrUpdateRequest
  ): Promise<CreateOrUpdateResponse>;
  /** Deletes idenity password. After this action, action can not be authentificated using password. */
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  /** Checks if password authentification method is defined for specified entity */
  Exist(request: ExistRequest): Promise<ExistResponse>;
}

export class IAMAuthenticationPasswordServiceClientImpl
  implements IAMAuthenticationPasswordService
{
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Authenticate = this.Authenticate.bind(this);
    this.CreateOrUpdate = this.CreateOrUpdate.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Exist = this.Exist.bind(this);
  }
  Authenticate(request: AuthenticateRequest): Promise<AuthenticateResponse> {
    const data = AuthenticateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_authentication_password.IAMAuthenticationPasswordService",
      "Authenticate",
      data
    );
    return promise.then((data) =>
      AuthenticateResponse.decode(new _m0.Reader(data))
    );
  }

  CreateOrUpdate(
    request: CreateOrUpdateRequest
  ): Promise<CreateOrUpdateResponse> {
    const data = CreateOrUpdateRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_authentication_password.IAMAuthenticationPasswordService",
      "CreateOrUpdate",
      data
    );
    return promise.then((data) =>
      CreateOrUpdateResponse.decode(new _m0.Reader(data))
    );
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_authentication_password.IAMAuthenticationPasswordService",
      "Delete",
      data
    );
    return promise.then((data) => DeleteResponse.decode(new _m0.Reader(data)));
  }

  Exist(request: ExistRequest): Promise<ExistResponse> {
    const data = ExistRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_iam_authentication_password.IAMAuthenticationPasswordService",
      "Exist",
      data
    );
    return promise.then((data) => ExistResponse.decode(new _m0.Reader(data)));
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
