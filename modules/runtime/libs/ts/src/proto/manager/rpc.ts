/* eslint-disable */
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "runtime_manager_runtime";

/** Message that is passedthrought NATS during the invocation of the RPC */
export interface RPCRequestMesasge {
  /** JSON foramted payload that will be passd to the invocated runtime method */
  data: string;
}

/** Response that is sended throught the NATS after the invocation of the RPC */
export interface RPCResponseMessage {
  /** JSON formated error that accured during the invocation of the RPC method. Empty if no error */
  error: string;
  /** Short message that describes the error. Empty if no error */
  errorMessage: string;
  /** JSON formated response from the RPC method */
  response: string;
}

export interface CallRequest {
  /** Namespace where runtime is located */
  namespace: string;
  /** Name of runtime */
  runtimeName: string;
  /** Name of method to call (without runtime name) */
  methodName: string;
  /** JSON payload */
  payload: string;
}

export interface CallResponse {
  /** JSON formated error accured during call. Empty if no error */
  error: string;
  /** Short message that describes the error. Empty if no error */
  errorMessage: string;
  /** JSON formated response from runtime */
  response: string;
}

function createBaseRPCRequestMesasge(): RPCRequestMesasge {
  return { data: "" };
}

export const RPCRequestMesasge = {
  encode(message: RPCRequestMesasge, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data !== "") {
      writer.uint32(10).string(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RPCRequestMesasge {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRPCRequestMesasge();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RPCRequestMesasge {
    return { data: isSet(object.data) ? String(object.data) : "" };
  },

  toJSON(message: RPCRequestMesasge): unknown {
    const obj: any = {};
    if (message.data !== "") {
      obj.data = message.data;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RPCRequestMesasge>, I>>(base?: I): RPCRequestMesasge {
    return RPCRequestMesasge.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RPCRequestMesasge>, I>>(object: I): RPCRequestMesasge {
    const message = createBaseRPCRequestMesasge();
    message.data = object.data ?? "";
    return message;
  },
};

function createBaseRPCResponseMessage(): RPCResponseMessage {
  return { error: "", errorMessage: "", response: "" };
}

export const RPCResponseMessage = {
  encode(message: RPCResponseMessage, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.error !== "") {
      writer.uint32(10).string(message.error);
    }
    if (message.errorMessage !== "") {
      writer.uint32(18).string(message.errorMessage);
    }
    if (message.response !== "") {
      writer.uint32(26).string(message.response);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RPCResponseMessage {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRPCResponseMessage();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.error = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.errorMessage = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.response = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RPCResponseMessage {
    return {
      error: isSet(object.error) ? String(object.error) : "",
      errorMessage: isSet(object.errorMessage) ? String(object.errorMessage) : "",
      response: isSet(object.response) ? String(object.response) : "",
    };
  },

  toJSON(message: RPCResponseMessage): unknown {
    const obj: any = {};
    if (message.error !== "") {
      obj.error = message.error;
    }
    if (message.errorMessage !== "") {
      obj.errorMessage = message.errorMessage;
    }
    if (message.response !== "") {
      obj.response = message.response;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RPCResponseMessage>, I>>(base?: I): RPCResponseMessage {
    return RPCResponseMessage.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RPCResponseMessage>, I>>(object: I): RPCResponseMessage {
    const message = createBaseRPCResponseMessage();
    message.error = object.error ?? "";
    message.errorMessage = object.errorMessage ?? "";
    message.response = object.response ?? "";
    return message;
  },
};

function createBaseCallRequest(): CallRequest {
  return { namespace: "", runtimeName: "", methodName: "", payload: "" };
}

export const CallRequest = {
  encode(message: CallRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.runtimeName !== "") {
      writer.uint32(18).string(message.runtimeName);
    }
    if (message.methodName !== "") {
      writer.uint32(26).string(message.methodName);
    }
    if (message.payload !== "") {
      writer.uint32(34).string(message.payload);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CallRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCallRequest();
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

          message.runtimeName = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.methodName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.payload = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CallRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      runtimeName: isSet(object.runtimeName) ? String(object.runtimeName) : "",
      methodName: isSet(object.methodName) ? String(object.methodName) : "",
      payload: isSet(object.payload) ? String(object.payload) : "",
    };
  },

  toJSON(message: CallRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.runtimeName !== "") {
      obj.runtimeName = message.runtimeName;
    }
    if (message.methodName !== "") {
      obj.methodName = message.methodName;
    }
    if (message.payload !== "") {
      obj.payload = message.payload;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CallRequest>, I>>(base?: I): CallRequest {
    return CallRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CallRequest>, I>>(object: I): CallRequest {
    const message = createBaseCallRequest();
    message.namespace = object.namespace ?? "";
    message.runtimeName = object.runtimeName ?? "";
    message.methodName = object.methodName ?? "";
    message.payload = object.payload ?? "";
    return message;
  },
};

function createBaseCallResponse(): CallResponse {
  return { error: "", errorMessage: "", response: "" };
}

export const CallResponse = {
  encode(message: CallResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.error !== "") {
      writer.uint32(10).string(message.error);
    }
    if (message.errorMessage !== "") {
      writer.uint32(18).string(message.errorMessage);
    }
    if (message.response !== "") {
      writer.uint32(26).string(message.response);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CallResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCallResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.error = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.errorMessage = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.response = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CallResponse {
    return {
      error: isSet(object.error) ? String(object.error) : "",
      errorMessage: isSet(object.errorMessage) ? String(object.errorMessage) : "",
      response: isSet(object.response) ? String(object.response) : "",
    };
  },

  toJSON(message: CallResponse): unknown {
    const obj: any = {};
    if (message.error !== "") {
      obj.error = message.error;
    }
    if (message.errorMessage !== "") {
      obj.errorMessage = message.errorMessage;
    }
    if (message.response !== "") {
      obj.response = message.response;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CallResponse>, I>>(base?: I): CallResponse {
    return CallResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CallResponse>, I>>(object: I): CallResponse {
    const message = createBaseCallResponse();
    message.error = object.error ?? "";
    message.errorMessage = object.errorMessage ?? "";
    message.response = object.response ?? "";
    return message;
  },
};

export interface RPCService {
  Call(request: CallRequest): Promise<CallResponse>;
}

export const RPCServiceServiceName = "runtime_manager_runtime.RPCService";
export class RPCServiceClientImpl implements RPCService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || RPCServiceServiceName;
    this.rpc = rpc;
    this.Call = this.Call.bind(this);
  }
  Call(request: CallRequest): Promise<CallResponse> {
    const data = CallRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Call", data);
    return promise.then((data) => CallResponse.decode(_m0.Reader.create(data)));
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
