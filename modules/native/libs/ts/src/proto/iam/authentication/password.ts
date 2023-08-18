/* eslint-disable */
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

export interface CreateOrUpdateResponse {
  /** True if password was created. False if it was updated and existed before this operation. */
  created: boolean;
}

export interface DeleteRequest {
  namespace: string;
  identity: string;
}

export interface DeleteResponse {
  /** Indicates if password existed before thi request or not. */
  existed: boolean;
}

export interface ExistsRequest {
  namespace: string;
  identity: string;
}

export interface ExistsResponse {
  exists: boolean;
}

function createBaseAuthenticateRequest(): AuthenticateRequest {
  return { namespace: "", identity: "", password: "" };
}

export const AuthenticateRequest = {
  encode(message: AuthenticateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticateRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.password !== "") {
      obj.password = message.password;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthenticateRequest>, I>>(base?: I): AuthenticateRequest {
    return AuthenticateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AuthenticateRequest>, I>>(object: I): AuthenticateRequest {
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
  encode(message: AuthenticateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.authenticated === true) {
      writer.uint32(8).bool(message.authenticated);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): AuthenticateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAuthenticateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.authenticated = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): AuthenticateResponse {
    return { authenticated: isSet(object.authenticated) ? Boolean(object.authenticated) : false };
  },

  toJSON(message: AuthenticateResponse): unknown {
    const obj: any = {};
    if (message.authenticated === true) {
      obj.authenticated = message.authenticated;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<AuthenticateResponse>, I>>(base?: I): AuthenticateResponse {
    return AuthenticateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<AuthenticateResponse>, I>>(object: I): AuthenticateResponse {
    const message = createBaseAuthenticateResponse();
    message.authenticated = object.authenticated ?? false;
    return message;
  },
};

function createBaseCreateOrUpdateRequest(): CreateOrUpdateRequest {
  return { namespace: "", identity: "", password: "" };
}

export const CreateOrUpdateRequest = {
  encode(message: CreateOrUpdateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateOrUpdateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateOrUpdateRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
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
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.password !== "") {
      obj.password = message.password;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateOrUpdateRequest>, I>>(base?: I): CreateOrUpdateRequest {
    return CreateOrUpdateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateOrUpdateRequest>, I>>(object: I): CreateOrUpdateRequest {
    const message = createBaseCreateOrUpdateRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.password = object.password ?? "";
    return message;
  },
};

function createBaseCreateOrUpdateResponse(): CreateOrUpdateResponse {
  return { created: false };
}

export const CreateOrUpdateResponse = {
  encode(message: CreateOrUpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.created === true) {
      writer.uint32(8).bool(message.created);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateOrUpdateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateOrUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.created = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateOrUpdateResponse {
    return { created: isSet(object.created) ? Boolean(object.created) : false };
  },

  toJSON(message: CreateOrUpdateResponse): unknown {
    const obj: any = {};
    if (message.created === true) {
      obj.created = message.created;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateOrUpdateResponse>, I>>(base?: I): CreateOrUpdateResponse {
    return CreateOrUpdateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateOrUpdateResponse>, I>>(object: I): CreateOrUpdateResponse {
    const message = createBaseCreateOrUpdateResponse();
    message.created = object.created ?? false;
    return message;
  },
};

function createBaseDeleteRequest(): DeleteRequest {
  return { namespace: "", identity: "" };
}

export const DeleteRequest = {
  encode(message: DeleteRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
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

  fromJSON(object: any): DeleteRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: DeleteRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteRequest>, I>>(base?: I): DeleteRequest {
    return DeleteRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteRequest>, I>>(object: I): DeleteRequest {
    const message = createBaseDeleteRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
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

function createBaseExistsRequest(): ExistsRequest {
  return { namespace: "", identity: "" };
}

export const ExistsRequest = {
  encode(message: ExistsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsRequest();
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
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExistsRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: ExistsRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistsRequest>, I>>(base?: I): ExistsRequest {
    return ExistsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistsRequest>, I>>(object: I): ExistsRequest {
    const message = createBaseExistsRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseExistsResponse(): ExistsResponse {
  return { exists: false };
}

export const ExistsResponse = {
  encode(message: ExistsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.exists === true) {
      writer.uint32(8).bool(message.exists);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.exists = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ExistsResponse {
    return { exists: isSet(object.exists) ? Boolean(object.exists) : false };
  },

  toJSON(message: ExistsResponse): unknown {
    const obj: any = {};
    if (message.exists === true) {
      obj.exists = message.exists;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ExistsResponse>, I>>(base?: I): ExistsResponse {
    return ExistsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ExistsResponse>, I>>(object: I): ExistsResponse {
    const message = createBaseExistsResponse();
    message.exists = object.exists ?? false;
    return message;
  },
};

/** Provides API to authentificate identities using password */
export interface IAMAuthenticationPasswordService {
  /** Tries to find identity and compare its password. */
  Authenticate(request: AuthenticateRequest): Promise<AuthenticateResponse>;
  /** Creates or updates identity password for authentification. */
  CreateOrUpdate(request: CreateOrUpdateRequest): Promise<CreateOrUpdateResponse>;
  /** Deletes idenity password. After this action, action can not be authentificated using password. */
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  /** Checks if password authentification method is defined for specified entity */
  Exists(request: ExistsRequest): Promise<ExistsResponse>;
}

export const IAMAuthenticationPasswordServiceServiceName =
  "native_iam_authentication_password.IAMAuthenticationPasswordService";
export class IAMAuthenticationPasswordServiceClientImpl implements IAMAuthenticationPasswordService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMAuthenticationPasswordServiceServiceName;
    this.rpc = rpc;
    this.Authenticate = this.Authenticate.bind(this);
    this.CreateOrUpdate = this.CreateOrUpdate.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Exists = this.Exists.bind(this);
  }
  Authenticate(request: AuthenticateRequest): Promise<AuthenticateResponse> {
    const data = AuthenticateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Authenticate", data);
    return promise.then((data) => AuthenticateResponse.decode(_m0.Reader.create(data)));
  }

  CreateOrUpdate(request: CreateOrUpdateRequest): Promise<CreateOrUpdateResponse> {
    const data = CreateOrUpdateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CreateOrUpdate", data);
    return promise.then((data) => CreateOrUpdateResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteRequest): Promise<DeleteResponse> {
    const data = DeleteRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteResponse.decode(_m0.Reader.create(data)));
  }

  Exists(request: ExistsRequest): Promise<ExistsResponse> {
    const data = ExistsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Exists", data);
    return promise.then((data) => ExistsResponse.decode(_m0.Reader.create(data)));
  }
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
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
