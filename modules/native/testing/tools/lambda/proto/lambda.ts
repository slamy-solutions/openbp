/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";

export const protobufPackage = "native_lambda";

export interface Lambda {
  /** Namespace where lambda was created */
  namespace: string;
  /** Unique identifier inside UUID. */
  uuid: string;
  /** Type of the lambda runtime, that will run lambda */
  runtime: string;
  /** Unique identifier of the bundle. It is used for better lambda provisioning. If two lambdas have same bundle they will be executed using same execution pod */
  bundle: Buffer;
  /** Make sure, that when this lambda will be executed, it will be executed only ones. If set to True, it will be very slow, because it has to use distributed lock to prevent multiple executions. Multiple executions are VERY rare, only at 0.1% of the times or less. Better implement you functions in such way, that several executions of the function with same data is ok. */
  ensureExactlyOneDelivery: boolean;
}

export interface CodeBundle {
  /** Unique indentifier of the code bundle. Must be binded to the code. Better to use hash of the code or binary file. */
  uuid: Buffer;
  /** Number of refenrences pointing to this code bundle. When number of references will be zero, bundle will be deleted. */
  references: number;
}

/** Structure in wich tasks are getting submited to the AMQP and received by runtime */
export interface AMQPLambdaTaskRequest {
  lambda: Lambda | undefined;
  data: Buffer;
}

/** Response from the lambda runtime in the amqp */
export interface AMQPLambdaTaskResponse {
  statusCode: number;
  message: string;
  data: Buffer;
}

export interface CreateLambdaRequest {
  /** Namespace where to create lambda */
  namespace: string;
  /** Unique identifier. */
  uuid: string;
  /** Type of the lambda runtime, that will run lambda */
  runtime: string;
  /** Unique identifier of the bundle. Must be strictly binded to the code. SHA256 of the code/binary will be used if empty */
  bundle: Buffer;
  /** Bundle binary data */
  data: Buffer;
  /** Make sure, that when this lambda will be executed, it will be executed only ones. If set to True, it will be very slow, because it has to use distributed lock to prevent multiple executions. Multiple executions are VERY rare, only at 0.1% of the times or less. Better implement you functions in such way, that several executions of the function with same data is ok. */
  ensureExactlyOneDelivery: boolean;
}

export interface CreateLambdaResponse {
  /** Created lambda */
  lambda: Lambda | undefined;
}

export interface DeleteLambdaRequest {
  /** Namespace of the lambda */
  namespace: string;
  /** Unique identifier of the lambda to delete. */
  uuid: string;
}

export interface DeleteLambdaResponse {}

export interface ExistsLambdaRequest {
  /** Namespace of the lambda */
  namespace: string;
  /** Unique identifier of the lambda inside namespace. */
  uuid: string;
}

export interface ExistsLambdaResponse {
  /** Lambda exists or not */
  exists: boolean;
}

export interface GetLambdaRequest {
  /** Namespace of the lambda */
  namespace: string;
  /** Unique identifier of the lambda to get. */
  uuid: string;
}

export interface GetLambdaResponse {
  Lambda: Lambda | undefined;
}

export interface GetBundleRequest {
  /** Unique identifier of the bundle to get. */
  bundle: Buffer;
}

export interface GetBundleResponse {
  /** Bundle data */
  data: Buffer;
}

export interface ExecuteLambdaRequest {
  /** Namespace of the lambda */
  namespace: string;
  /** UUID of the lambda */
  lambda: string;
  /** Data that will be passed to the function */
  data: Buffer;
  /** Execution timeout in miliseconds */
  timeout: number;
}

export interface ExecuteLambdaResponse {
  /** Execution result */
  result: Buffer;
}

export interface CallLambdaRequest {
  /** Namespace of the lambda */
  namespace: string;
  /** UUID of the lambda */
  lambda: string;
  /** Data that will be passed to the function */
  data: Buffer;
}

export interface CallLambdaResponse {}

function createBaseLambda(): Lambda {
  return {
    namespace: "",
    uuid: "",
    runtime: "",
    bundle: Buffer.alloc(0),
    ensureExactlyOneDelivery: false,
  };
}

export const Lambda = {
  encode(
    message: Lambda,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.runtime !== "") {
      writer.uint32(26).string(message.runtime);
    }
    if (message.bundle.length !== 0) {
      writer.uint32(34).bytes(message.bundle);
    }
    if (message.ensureExactlyOneDelivery === true) {
      writer.uint32(40).bool(message.ensureExactlyOneDelivery);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Lambda {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseLambda();
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
          message.runtime = reader.string();
          break;
        case 4:
          message.bundle = reader.bytes() as Buffer;
          break;
        case 5:
          message.ensureExactlyOneDelivery = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): Lambda {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      runtime: isSet(object.runtime) ? String(object.runtime) : "",
      bundle: isSet(object.bundle)
        ? Buffer.from(bytesFromBase64(object.bundle))
        : Buffer.alloc(0),
      ensureExactlyOneDelivery: isSet(object.ensureExactlyOneDelivery)
        ? Boolean(object.ensureExactlyOneDelivery)
        : false,
    };
  },

  toJSON(message: Lambda): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.runtime !== undefined && (obj.runtime = message.runtime);
    message.bundle !== undefined &&
      (obj.bundle = base64FromBytes(
        message.bundle !== undefined ? message.bundle : Buffer.alloc(0)
      ));
    message.ensureExactlyOneDelivery !== undefined &&
      (obj.ensureExactlyOneDelivery = message.ensureExactlyOneDelivery);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<Lambda>, I>>(object: I): Lambda {
    const message = createBaseLambda();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.runtime = object.runtime ?? "";
    message.bundle = object.bundle ?? Buffer.alloc(0);
    message.ensureExactlyOneDelivery = object.ensureExactlyOneDelivery ?? false;
    return message;
  },
};

function createBaseCodeBundle(): CodeBundle {
  return { uuid: Buffer.alloc(0), references: 0 };
}

export const CodeBundle = {
  encode(
    message: CodeBundle,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.uuid.length !== 0) {
      writer.uint32(10).bytes(message.uuid);
    }
    if (message.references !== 0) {
      writer.uint32(16).int32(message.references);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CodeBundle {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCodeBundle();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.uuid = reader.bytes() as Buffer;
          break;
        case 2:
          message.references = reader.int32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CodeBundle {
    return {
      uuid: isSet(object.uuid)
        ? Buffer.from(bytesFromBase64(object.uuid))
        : Buffer.alloc(0),
      references: isSet(object.references) ? Number(object.references) : 0,
    };
  },

  toJSON(message: CodeBundle): unknown {
    const obj: any = {};
    message.uuid !== undefined &&
      (obj.uuid = base64FromBytes(
        message.uuid !== undefined ? message.uuid : Buffer.alloc(0)
      ));
    message.references !== undefined &&
      (obj.references = Math.round(message.references));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CodeBundle>, I>>(
    object: I
  ): CodeBundle {
    const message = createBaseCodeBundle();
    message.uuid = object.uuid ?? Buffer.alloc(0);
    message.references = object.references ?? 0;
    return message;
  },
};

function createBaseAMQPLambdaTaskRequest(): AMQPLambdaTaskRequest {
  return { lambda: undefined, data: Buffer.alloc(0) };
}

export const AMQPLambdaTaskRequest = {
  encode(
    message: AMQPLambdaTaskRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.lambda !== undefined) {
      Lambda.encode(message.lambda, writer.uint32(10).fork()).ldelim();
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AMQPLambdaTaskRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAMQPLambdaTaskRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.lambda = Lambda.decode(reader, reader.uint32());
          break;
        case 2:
          message.data = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AMQPLambdaTaskRequest {
    return {
      lambda: isSet(object.lambda) ? Lambda.fromJSON(object.lambda) : undefined,
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: AMQPLambdaTaskRequest): unknown {
    const obj: any = {};
    message.lambda !== undefined &&
      (obj.lambda = message.lambda ? Lambda.toJSON(message.lambda) : undefined);
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AMQPLambdaTaskRequest>, I>>(
    object: I
  ): AMQPLambdaTaskRequest {
    const message = createBaseAMQPLambdaTaskRequest();
    message.lambda =
      object.lambda !== undefined && object.lambda !== null
        ? Lambda.fromPartial(object.lambda)
        : undefined;
    message.data = object.data ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseAMQPLambdaTaskResponse(): AMQPLambdaTaskResponse {
  return { statusCode: 0, message: "", data: Buffer.alloc(0) };
}

export const AMQPLambdaTaskResponse = {
  encode(
    message: AMQPLambdaTaskResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.statusCode !== 0) {
      writer.uint32(8).uint32(message.statusCode);
    }
    if (message.message !== "") {
      writer.uint32(18).string(message.message);
    }
    if (message.data.length !== 0) {
      writer.uint32(26).bytes(message.data);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): AMQPLambdaTaskResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseAMQPLambdaTaskResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.statusCode = reader.uint32();
          break;
        case 2:
          message.message = reader.string();
          break;
        case 3:
          message.data = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AMQPLambdaTaskResponse {
    return {
      statusCode: isSet(object.statusCode) ? Number(object.statusCode) : 0,
      message: isSet(object.message) ? String(object.message) : "",
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: AMQPLambdaTaskResponse): unknown {
    const obj: any = {};
    message.statusCode !== undefined &&
      (obj.statusCode = Math.round(message.statusCode));
    message.message !== undefined && (obj.message = message.message);
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<AMQPLambdaTaskResponse>, I>>(
    object: I
  ): AMQPLambdaTaskResponse {
    const message = createBaseAMQPLambdaTaskResponse();
    message.statusCode = object.statusCode ?? 0;
    message.message = object.message ?? "";
    message.data = object.data ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseCreateLambdaRequest(): CreateLambdaRequest {
  return {
    namespace: "",
    uuid: "",
    runtime: "",
    bundle: Buffer.alloc(0),
    data: Buffer.alloc(0),
    ensureExactlyOneDelivery: false,
  };
}

export const CreateLambdaRequest = {
  encode(
    message: CreateLambdaRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.runtime !== "") {
      writer.uint32(26).string(message.runtime);
    }
    if (message.bundle.length !== 0) {
      writer.uint32(34).bytes(message.bundle);
    }
    if (message.data.length !== 0) {
      writer.uint32(42).bytes(message.data);
    }
    if (message.ensureExactlyOneDelivery === true) {
      writer.uint32(48).bool(message.ensureExactlyOneDelivery);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateLambdaRequest();
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
          message.runtime = reader.string();
          break;
        case 4:
          message.bundle = reader.bytes() as Buffer;
          break;
        case 5:
          message.data = reader.bytes() as Buffer;
          break;
        case 6:
          message.ensureExactlyOneDelivery = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      runtime: isSet(object.runtime) ? String(object.runtime) : "",
      bundle: isSet(object.bundle)
        ? Buffer.from(bytesFromBase64(object.bundle))
        : Buffer.alloc(0),
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
      ensureExactlyOneDelivery: isSet(object.ensureExactlyOneDelivery)
        ? Boolean(object.ensureExactlyOneDelivery)
        : false,
    };
  },

  toJSON(message: CreateLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.runtime !== undefined && (obj.runtime = message.runtime);
    message.bundle !== undefined &&
      (obj.bundle = base64FromBytes(
        message.bundle !== undefined ? message.bundle : Buffer.alloc(0)
      ));
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    message.ensureExactlyOneDelivery !== undefined &&
      (obj.ensureExactlyOneDelivery = message.ensureExactlyOneDelivery);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateLambdaRequest>, I>>(
    object: I
  ): CreateLambdaRequest {
    const message = createBaseCreateLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.runtime = object.runtime ?? "";
    message.bundle = object.bundle ?? Buffer.alloc(0);
    message.data = object.data ?? Buffer.alloc(0);
    message.ensureExactlyOneDelivery = object.ensureExactlyOneDelivery ?? false;
    return message;
  },
};

function createBaseCreateLambdaResponse(): CreateLambdaResponse {
  return { lambda: undefined };
}

export const CreateLambdaResponse = {
  encode(
    message: CreateLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.lambda !== undefined) {
      Lambda.encode(message.lambda, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CreateLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateLambdaResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.lambda = Lambda.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CreateLambdaResponse {
    return {
      lambda: isSet(object.lambda) ? Lambda.fromJSON(object.lambda) : undefined,
    };
  },

  toJSON(message: CreateLambdaResponse): unknown {
    const obj: any = {};
    message.lambda !== undefined &&
      (obj.lambda = message.lambda ? Lambda.toJSON(message.lambda) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CreateLambdaResponse>, I>>(
    object: I
  ): CreateLambdaResponse {
    const message = createBaseCreateLambdaResponse();
    message.lambda =
      object.lambda !== undefined && object.lambda !== null
        ? Lambda.fromPartial(object.lambda)
        : undefined;
    return message;
  },
};

function createBaseDeleteLambdaRequest(): DeleteLambdaRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteLambdaRequest = {
  encode(
    message: DeleteLambdaRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteLambdaRequest();
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

  fromJSON(object: any): DeleteLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteLambdaRequest>, I>>(
    object: I
  ): DeleteLambdaRequest {
    const message = createBaseDeleteLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteLambdaResponse(): DeleteLambdaResponse {
  return {};
}

export const DeleteLambdaResponse = {
  encode(
    _: DeleteLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): DeleteLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteLambdaResponse();
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

  fromJSON(_: any): DeleteLambdaResponse {
    return {};
  },

  toJSON(_: DeleteLambdaResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteLambdaResponse>, I>>(
    _: I
  ): DeleteLambdaResponse {
    const message = createBaseDeleteLambdaResponse();
    return message;
  },
};

function createBaseExistsLambdaRequest(): ExistsLambdaRequest {
  return { namespace: "", uuid: "" };
}

export const ExistsLambdaRequest = {
  encode(
    message: ExistsLambdaRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): ExistsLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsLambdaRequest();
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

  fromJSON(object: any): ExistsLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: ExistsLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExistsLambdaRequest>, I>>(
    object: I
  ): ExistsLambdaRequest {
    const message = createBaseExistsLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseExistsLambdaResponse(): ExistsLambdaResponse {
  return { exists: false };
}

export const ExistsLambdaResponse = {
  encode(
    message: ExistsLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.exists === true) {
      writer.uint32(8).bool(message.exists);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ExistsLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExistsLambdaResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.exists = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExistsLambdaResponse {
    return {
      exists: isSet(object.exists) ? Boolean(object.exists) : false,
    };
  },

  toJSON(message: ExistsLambdaResponse): unknown {
    const obj: any = {};
    message.exists !== undefined && (obj.exists = message.exists);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExistsLambdaResponse>, I>>(
    object: I
  ): ExistsLambdaResponse {
    const message = createBaseExistsLambdaResponse();
    message.exists = object.exists ?? false;
    return message;
  },
};

function createBaseGetLambdaRequest(): GetLambdaRequest {
  return { namespace: "", uuid: "" };
}

export const GetLambdaRequest = {
  encode(
    message: GetLambdaRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLambdaRequest();
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

  fromJSON(object: any): GetLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: GetLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetLambdaRequest>, I>>(
    object: I
  ): GetLambdaRequest {
    const message = createBaseGetLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseGetLambdaResponse(): GetLambdaResponse {
  return { Lambda: undefined };
}

export const GetLambdaResponse = {
  encode(
    message: GetLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.Lambda !== undefined) {
      Lambda.encode(message.Lambda, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetLambdaResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.Lambda = Lambda.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetLambdaResponse {
    return {
      Lambda: isSet(object.Lambda) ? Lambda.fromJSON(object.Lambda) : undefined,
    };
  },

  toJSON(message: GetLambdaResponse): unknown {
    const obj: any = {};
    message.Lambda !== undefined &&
      (obj.Lambda = message.Lambda ? Lambda.toJSON(message.Lambda) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetLambdaResponse>, I>>(
    object: I
  ): GetLambdaResponse {
    const message = createBaseGetLambdaResponse();
    message.Lambda =
      object.Lambda !== undefined && object.Lambda !== null
        ? Lambda.fromPartial(object.Lambda)
        : undefined;
    return message;
  },
};

function createBaseGetBundleRequest(): GetBundleRequest {
  return { bundle: Buffer.alloc(0) };
}

export const GetBundleRequest = {
  encode(
    message: GetBundleRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.bundle.length !== 0) {
      writer.uint32(10).bytes(message.bundle);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBundleRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBundleRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.bundle = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetBundleRequest {
    return {
      bundle: isSet(object.bundle)
        ? Buffer.from(bytesFromBase64(object.bundle))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: GetBundleRequest): unknown {
    const obj: any = {};
    message.bundle !== undefined &&
      (obj.bundle = base64FromBytes(
        message.bundle !== undefined ? message.bundle : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetBundleRequest>, I>>(
    object: I
  ): GetBundleRequest {
    const message = createBaseGetBundleRequest();
    message.bundle = object.bundle ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseGetBundleResponse(): GetBundleResponse {
  return { data: Buffer.alloc(0) };
}

export const GetBundleResponse = {
  encode(
    message: GetBundleResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetBundleResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetBundleResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.data = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): GetBundleResponse {
    return {
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: GetBundleResponse): unknown {
    const obj: any = {};
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<GetBundleResponse>, I>>(
    object: I
  ): GetBundleResponse {
    const message = createBaseGetBundleResponse();
    message.data = object.data ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseExecuteLambdaRequest(): ExecuteLambdaRequest {
  return { namespace: "", lambda: "", data: Buffer.alloc(0), timeout: 0 };
}

export const ExecuteLambdaRequest = {
  encode(
    message: ExecuteLambdaRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.lambda !== "") {
      writer.uint32(18).string(message.lambda);
    }
    if (message.data.length !== 0) {
      writer.uint32(26).bytes(message.data);
    }
    if (message.timeout !== 0) {
      writer.uint32(32).uint64(message.timeout);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ExecuteLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExecuteLambdaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.lambda = reader.string();
          break;
        case 3:
          message.data = reader.bytes() as Buffer;
          break;
        case 4:
          message.timeout = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExecuteLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      lambda: isSet(object.lambda) ? String(object.lambda) : "",
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
      timeout: isSet(object.timeout) ? Number(object.timeout) : 0,
    };
  },

  toJSON(message: ExecuteLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.lambda !== undefined && (obj.lambda = message.lambda);
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    message.timeout !== undefined &&
      (obj.timeout = Math.round(message.timeout));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExecuteLambdaRequest>, I>>(
    object: I
  ): ExecuteLambdaRequest {
    const message = createBaseExecuteLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.lambda = object.lambda ?? "";
    message.data = object.data ?? Buffer.alloc(0);
    message.timeout = object.timeout ?? 0;
    return message;
  },
};

function createBaseExecuteLambdaResponse(): ExecuteLambdaResponse {
  return { result: Buffer.alloc(0) };
}

export const ExecuteLambdaResponse = {
  encode(
    message: ExecuteLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.result.length !== 0) {
      writer.uint32(10).bytes(message.result);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): ExecuteLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseExecuteLambdaResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.result = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ExecuteLambdaResponse {
    return {
      result: isSet(object.result)
        ? Buffer.from(bytesFromBase64(object.result))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: ExecuteLambdaResponse): unknown {
    const obj: any = {};
    message.result !== undefined &&
      (obj.result = base64FromBytes(
        message.result !== undefined ? message.result : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ExecuteLambdaResponse>, I>>(
    object: I
  ): ExecuteLambdaResponse {
    const message = createBaseExecuteLambdaResponse();
    message.result = object.result ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseCallLambdaRequest(): CallLambdaRequest {
  return { namespace: "", lambda: "", data: Buffer.alloc(0) };
}

export const CallLambdaRequest = {
  encode(
    message: CallLambdaRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.lambda !== "") {
      writer.uint32(18).string(message.lambda);
    }
    if (message.data.length !== 0) {
      writer.uint32(26).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CallLambdaRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCallLambdaRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.namespace = reader.string();
          break;
        case 2:
          message.lambda = reader.string();
          break;
        case 3:
          message.data = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CallLambdaRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      lambda: isSet(object.lambda) ? String(object.lambda) : "",
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: CallLambdaRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.lambda !== undefined && (obj.lambda = message.lambda);
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CallLambdaRequest>, I>>(
    object: I
  ): CallLambdaRequest {
    const message = createBaseCallLambdaRequest();
    message.namespace = object.namespace ?? "";
    message.lambda = object.lambda ?? "";
    message.data = object.data ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseCallLambdaResponse(): CallLambdaResponse {
  return {};
}

export const CallLambdaResponse = {
  encode(
    _: CallLambdaResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CallLambdaResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCallLambdaResponse();
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

  fromJSON(_: any): CallLambdaResponse {
    return {};
  },

  toJSON(_: CallLambdaResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CallLambdaResponse>, I>>(
    _: I
  ): CallLambdaResponse {
    const message = createBaseCallLambdaResponse();
    return message;
  },
};

/** Creates/Deletes/Updates/Agregates information about lambdas */
export interface LambdaManagerService {
  /** Create new lambda */
  Create(request: CreateLambdaRequest): Promise<CreateLambdaResponse>;
  /** Deletes lambda */
  Delete(request: DeleteLambdaRequest): Promise<DeleteLambdaResponse>;
  /** Checks if lambda exists or not */
  Exists(request: ExistsLambdaRequest): Promise<ExistsLambdaResponse>;
  /** Get lambda information */
  Get(request: GetLambdaRequest): Promise<GetLambdaResponse>;
  /** Gets lambda bundle */
  GetBundle(request: GetBundleRequest): Promise<GetBundleResponse>;
}

export class LambdaManagerServiceClientImpl implements LambdaManagerService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Exists = this.Exists.bind(this);
    this.Get = this.Get.bind(this);
    this.GetBundle = this.GetBundle.bind(this);
  }
  Create(request: CreateLambdaRequest): Promise<CreateLambdaResponse> {
    const data = CreateLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaManagerService",
      "Create",
      data
    );
    return promise.then((data) =>
      CreateLambdaResponse.decode(new _m0.Reader(data))
    );
  }

  Delete(request: DeleteLambdaRequest): Promise<DeleteLambdaResponse> {
    const data = DeleteLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaManagerService",
      "Delete",
      data
    );
    return promise.then((data) =>
      DeleteLambdaResponse.decode(new _m0.Reader(data))
    );
  }

  Exists(request: ExistsLambdaRequest): Promise<ExistsLambdaResponse> {
    const data = ExistsLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaManagerService",
      "Exists",
      data
    );
    return promise.then((data) =>
      ExistsLambdaResponse.decode(new _m0.Reader(data))
    );
  }

  Get(request: GetLambdaRequest): Promise<GetLambdaResponse> {
    const data = GetLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaManagerService",
      "Get",
      data
    );
    return promise.then((data) =>
      GetLambdaResponse.decode(new _m0.Reader(data))
    );
  }

  GetBundle(request: GetBundleRequest): Promise<GetBundleResponse> {
    const data = GetBundleRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaManagerService",
      "GetBundle",
      data
    );
    return promise.then((data) =>
      GetBundleResponse.decode(new _m0.Reader(data))
    );
  }
}

/** Provides API to execute lambda functions */
export interface LambdaEntrypointService {
  /** Runs lambda and doesnt wait for response. Call will be resheduled on internal error, also it can be resheduled by lambda code. */
  Call(request: CallLambdaRequest): Promise<CallLambdaResponse>;
  /** Runs function and returns its response. Returns error if something went wrong during the execution. */
  Execute(request: ExecuteLambdaRequest): Promise<ExecuteLambdaResponse>;
}

export class LambdaEntrypointServiceClientImpl
  implements LambdaEntrypointService
{
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Call = this.Call.bind(this);
    this.Execute = this.Execute.bind(this);
  }
  Call(request: CallLambdaRequest): Promise<CallLambdaResponse> {
    const data = CallLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaEntrypointService",
      "Call",
      data
    );
    return promise.then((data) =>
      CallLambdaResponse.decode(new _m0.Reader(data))
    );
  }

  Execute(request: ExecuteLambdaRequest): Promise<ExecuteLambdaResponse> {
    const data = ExecuteLambdaRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_lambda.LambdaEntrypointService",
      "Execute",
      data
    );
    return promise.then((data) =>
      ExecuteLambdaResponse.decode(new _m0.Reader(data))
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

declare var self: any | undefined;
declare var window: any | undefined;
declare var global: any | undefined;
var globalThis: any = (() => {
  if (typeof globalThis !== "undefined") return globalThis;
  if (typeof self !== "undefined") return self;
  if (typeof window !== "undefined") return window;
  if (typeof global !== "undefined") return global;
  throw "Unable to locate global object";
})();

const atob: (b64: string) => string =
  globalThis.atob ||
  ((b64) => globalThis.Buffer.from(b64, "base64").toString("binary"));
function bytesFromBase64(b64: string): Uint8Array {
  const bin = atob(b64);
  const arr = new Uint8Array(bin.length);
  for (let i = 0; i < bin.length; ++i) {
    arr[i] = bin.charCodeAt(i);
  }
  return arr;
}

const btoa: (bin: string) => string =
  globalThis.btoa ||
  ((bin) => globalThis.Buffer.from(bin, "binary").toString("base64"));
function base64FromBytes(arr: Uint8Array): string {
  const bin: string[] = [];
  arr.forEach((byte) => {
    bin.push(String.fromCharCode(byte));
  });
  return btoa(bin.join(""));
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

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new globalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
  }
  return long.toNumber();
}

if (_m0.util.Long !== Long) {
  _m0.util.Long = Long as any;
  _m0.configure();
}

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}