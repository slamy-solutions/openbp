/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { Timestamp } from "../../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "native_iam_authentication_x509";

export interface Certificate {
  /** Namespace where indetity and its certificate are located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
  /** Unique identifier of the identity */
  identity: string;
  /** Indicates if certificate was manually disabled. Disabled certificate connot be used. */
  disabled: boolean;
  /** Arbitrary, human-readable desription of the certificate */
  description: string;
  /** RSA public key in DER format */
  publicKey: Uint8Array;
  /** When the certificate was created */
  created:
    | Date
    | undefined;
  /** Last time when the certificate information was updated. */
  updated:
    | Date
    | undefined;
  /** Counter that increases after every update of the certificate */
  version: number;
}

export interface GetRootCAInfoRequest {
}

export interface GetRootCAInfoResponse {
  /** x509 certificate in the DER format */
  certificate: Uint8Array;
}

export interface RegisterAndGenerateRequest {
  /** Namespace where identity is located and where to generate certificate */
  namespace: string;
  /** Identity unique identifier */
  identity: string;
  /** Public key of the identity in the PEN format. Should be generated externaly in order not to share the public key with service. */
  publicKey: Uint8Array;
  /** Arbitrary, human-readable desription of the certificate */
  description: string;
}

export interface RegisterAndGenerateResponse {
  /** Certificate in DER format signed by CA of the service */
  raw: Uint8Array;
  /** Certificate information */
  info: Certificate | undefined;
}

export interface RegenerateRequest {
  /** Namespace where certificate is located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
}

export interface RegenerateResponse {
  /** New X509 in DER format signed by CA of the service */
  certificate: Uint8Array;
}

export interface ValidateAndGetFromRawX509Request {
  /** Certificate in DER format */
  raw: Uint8Array;
}

export interface ValidateAndGetFromRawX509Response {
  /** Status of the validation and search */
  status: ValidateAndGetFromRawX509Response_Status;
  /** Certificate information if status is OK */
  certificate: Certificate | undefined;
}

export enum ValidateAndGetFromRawX509Response_Status {
  /** OK - Everything is ok */
  OK = 0,
  /** INVALID_FORMAT - Certificate is corrupted or it wasnt supplied in DER format */
  INVALID_FORMAT = 1,
  /** SIGNATURE_INVALID - Certificate has invalid signature. */
  SIGNATURE_INVALID = 2,
  /** NOT_FOUND - Cant find certificate. */
  NOT_FOUND = 3,
  UNRECOGNIZED = -1,
}

export function validateAndGetFromRawX509Response_StatusFromJSON(
  object: any,
): ValidateAndGetFromRawX509Response_Status {
  switch (object) {
    case 0:
    case "OK":
      return ValidateAndGetFromRawX509Response_Status.OK;
    case 1:
    case "INVALID_FORMAT":
      return ValidateAndGetFromRawX509Response_Status.INVALID_FORMAT;
    case 2:
    case "SIGNATURE_INVALID":
      return ValidateAndGetFromRawX509Response_Status.SIGNATURE_INVALID;
    case 3:
    case "NOT_FOUND":
      return ValidateAndGetFromRawX509Response_Status.NOT_FOUND;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ValidateAndGetFromRawX509Response_Status.UNRECOGNIZED;
  }
}

export function validateAndGetFromRawX509Response_StatusToJSON(
  object: ValidateAndGetFromRawX509Response_Status,
): string {
  switch (object) {
    case ValidateAndGetFromRawX509Response_Status.OK:
      return "OK";
    case ValidateAndGetFromRawX509Response_Status.INVALID_FORMAT:
      return "INVALID_FORMAT";
    case ValidateAndGetFromRawX509Response_Status.SIGNATURE_INVALID:
      return "SIGNATURE_INVALID";
    case ValidateAndGetFromRawX509Response_Status.NOT_FOUND:
      return "NOT_FOUND";
    case ValidateAndGetFromRawX509Response_Status.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface GetRequest {
  /** Namespace where certificate is located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
}

export interface GetResponse {
  certificate: Certificate | undefined;
}

export interface CountRequest {
  /** Namespace where to count certificates */
  namespace: string;
}

export interface CountResponse {
  count: number;
}

export interface ListRequest {
  /** Namespace where to list certificates */
  namespace: string;
  /** How much entries to skip before returning actual entries */
  skip: number;
  /** Limit response to specified count of entries. Use 0 to ignore this and return all the possible entries. */
  limit: number;
}

export interface ListResponse {
  /** One of the certificates */
  certificate: Certificate | undefined;
}

export interface CountForIdentityRequest {
  /** Namespace where to count certificates */
  namespace: string;
  /** Identity unique identifier for each to count certificates */
  identity: string;
}

export interface CountForIdentityResponse {
  count: number;
}

export interface ListForIdentityRequest {
  /** Namespace where to list certificates */
  namespace: string;
  /** Identity unique identifier for each to list certificates */
  identity: string;
  /** How much entries to skip before returning actual entries */
  skip: number;
  /** Limit response to specified count of entries. Use 0 to ignore this and return all the possible entries. */
  limit: number;
}

export interface ListForIdentityResponse {
  /** One of the certificates */
  certificate: Certificate | undefined;
}

export interface UpdateRequest {
  /** Namespace where certificate is located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
  /** New description */
  newDescription: string;
}

export interface UpdateResponse {
  /** Updated certificate information */
  certificate: Certificate | undefined;
}

export interface DeleteRequest {
  /** Namespace where certificate is located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
}

export interface DeleteResponse {
  /** indicates if certificate existed before this operation */
  existed: boolean;
}

export interface DisableRequest {
  /** Namespace where certificate is located */
  namespace: string;
  /** Unique identifier of the certificate */
  uuid: string;
}

export interface DisableResponse {
  /** indicates if certificate was active before this operation */
  wasActive: boolean;
}

function createBaseCertificate(): Certificate {
  return {
    namespace: "",
    uuid: "",
    identity: "",
    disabled: false,
    description: "",
    publicKey: new Uint8Array(0),
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const Certificate = {
  encode(message: Certificate, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
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
    if (message.description !== "") {
      writer.uint32(42).string(message.description);
    }
    if (message.publicKey.length !== 0) {
      writer.uint32(50).bytes(message.publicKey);
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(74).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(82).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(88).uint64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Certificate {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCertificate();
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

          message.description = reader.string();
          continue;
        case 6:
          if (tag !== 50) {
            break;
          }

          message.publicKey = reader.bytes();
          continue;
        case 9:
          if (tag !== 74) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 11:
          if (tag !== 88) {
            break;
          }

          message.version = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Certificate {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      disabled: isSet(object.disabled) ? Boolean(object.disabled) : false,
      description: isSet(object.description) ? String(object.description) : "",
      publicKey: isSet(object.publicKey) ? bytesFromBase64(object.publicKey) : new Uint8Array(0),
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Certificate): unknown {
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
    if (message.description !== "") {
      obj.description = message.description;
    }
    if (message.publicKey.length !== 0) {
      obj.publicKey = base64FromBytes(message.publicKey);
    }
    if (message.created !== undefined) {
      obj.created = message.created.toISOString();
    }
    if (message.updated !== undefined) {
      obj.updated = message.updated.toISOString();
    }
    if (message.version !== 0) {
      obj.version = Math.round(message.version);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<Certificate>, I>>(base?: I): Certificate {
    return Certificate.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Certificate>, I>>(object: I): Certificate {
    const message = createBaseCertificate();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.identity = object.identity ?? "";
    message.disabled = object.disabled ?? false;
    message.description = object.description ?? "";
    message.publicKey = object.publicKey ?? new Uint8Array(0);
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseGetRootCAInfoRequest(): GetRootCAInfoRequest {
  return {};
}

export const GetRootCAInfoRequest = {
  encode(_: GetRootCAInfoRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRootCAInfoRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRootCAInfoRequest();
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

  fromJSON(_: any): GetRootCAInfoRequest {
    return {};
  },

  toJSON(_: GetRootCAInfoRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRootCAInfoRequest>, I>>(base?: I): GetRootCAInfoRequest {
    return GetRootCAInfoRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRootCAInfoRequest>, I>>(_: I): GetRootCAInfoRequest {
    const message = createBaseGetRootCAInfoRequest();
    return message;
  },
};

function createBaseGetRootCAInfoResponse(): GetRootCAInfoResponse {
  return { certificate: new Uint8Array(0) };
}

export const GetRootCAInfoResponse = {
  encode(message: GetRootCAInfoResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate.length !== 0) {
      writer.uint32(10).bytes(message.certificate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRootCAInfoResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRootCAInfoResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRootCAInfoResponse {
    return { certificate: isSet(object.certificate) ? bytesFromBase64(object.certificate) : new Uint8Array(0) };
  },

  toJSON(message: GetRootCAInfoResponse): unknown {
    const obj: any = {};
    if (message.certificate.length !== 0) {
      obj.certificate = base64FromBytes(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRootCAInfoResponse>, I>>(base?: I): GetRootCAInfoResponse {
    return GetRootCAInfoResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRootCAInfoResponse>, I>>(object: I): GetRootCAInfoResponse {
    const message = createBaseGetRootCAInfoResponse();
    message.certificate = object.certificate ?? new Uint8Array(0);
    return message;
  },
};

function createBaseRegisterAndGenerateRequest(): RegisterAndGenerateRequest {
  return { namespace: "", identity: "", publicKey: new Uint8Array(0), description: "" };
}

export const RegisterAndGenerateRequest = {
  encode(message: RegisterAndGenerateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    if (message.publicKey.length !== 0) {
      writer.uint32(26).bytes(message.publicKey);
    }
    if (message.description !== "") {
      writer.uint32(34).string(message.description);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterAndGenerateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterAndGenerateRequest();
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

          message.publicKey = reader.bytes();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.description = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RegisterAndGenerateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      publicKey: isSet(object.publicKey) ? bytesFromBase64(object.publicKey) : new Uint8Array(0),
      description: isSet(object.description) ? String(object.description) : "",
    };
  },

  toJSON(message: RegisterAndGenerateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.publicKey.length !== 0) {
      obj.publicKey = base64FromBytes(message.publicKey);
    }
    if (message.description !== "") {
      obj.description = message.description;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RegisterAndGenerateRequest>, I>>(base?: I): RegisterAndGenerateRequest {
    return RegisterAndGenerateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RegisterAndGenerateRequest>, I>>(object: I): RegisterAndGenerateRequest {
    const message = createBaseRegisterAndGenerateRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.publicKey = object.publicKey ?? new Uint8Array(0);
    message.description = object.description ?? "";
    return message;
  },
};

function createBaseRegisterAndGenerateResponse(): RegisterAndGenerateResponse {
  return { raw: new Uint8Array(0), info: undefined };
}

export const RegisterAndGenerateResponse = {
  encode(message: RegisterAndGenerateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.raw.length !== 0) {
      writer.uint32(10).bytes(message.raw);
    }
    if (message.info !== undefined) {
      Certificate.encode(message.info, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegisterAndGenerateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegisterAndGenerateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.raw = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.info = Certificate.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RegisterAndGenerateResponse {
    return {
      raw: isSet(object.raw) ? bytesFromBase64(object.raw) : new Uint8Array(0),
      info: isSet(object.info) ? Certificate.fromJSON(object.info) : undefined,
    };
  },

  toJSON(message: RegisterAndGenerateResponse): unknown {
    const obj: any = {};
    if (message.raw.length !== 0) {
      obj.raw = base64FromBytes(message.raw);
    }
    if (message.info !== undefined) {
      obj.info = Certificate.toJSON(message.info);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RegisterAndGenerateResponse>, I>>(base?: I): RegisterAndGenerateResponse {
    return RegisterAndGenerateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RegisterAndGenerateResponse>, I>>(object: I): RegisterAndGenerateResponse {
    const message = createBaseRegisterAndGenerateResponse();
    message.raw = object.raw ?? new Uint8Array(0);
    message.info = (object.info !== undefined && object.info !== null)
      ? Certificate.fromPartial(object.info)
      : undefined;
    return message;
  },
};

function createBaseRegenerateRequest(): RegenerateRequest {
  return { namespace: "", uuid: "" };
}

export const RegenerateRequest = {
  encode(message: RegenerateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegenerateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegenerateRequest();
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

  fromJSON(object: any): RegenerateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: RegenerateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RegenerateRequest>, I>>(base?: I): RegenerateRequest {
    return RegenerateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RegenerateRequest>, I>>(object: I): RegenerateRequest {
    const message = createBaseRegenerateRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseRegenerateResponse(): RegenerateResponse {
  return { certificate: new Uint8Array(0) };
}

export const RegenerateResponse = {
  encode(message: RegenerateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate.length !== 0) {
      writer.uint32(10).bytes(message.certificate);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RegenerateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRegenerateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RegenerateResponse {
    return { certificate: isSet(object.certificate) ? bytesFromBase64(object.certificate) : new Uint8Array(0) };
  },

  toJSON(message: RegenerateResponse): unknown {
    const obj: any = {};
    if (message.certificate.length !== 0) {
      obj.certificate = base64FromBytes(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RegenerateResponse>, I>>(base?: I): RegenerateResponse {
    return RegenerateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RegenerateResponse>, I>>(object: I): RegenerateResponse {
    const message = createBaseRegenerateResponse();
    message.certificate = object.certificate ?? new Uint8Array(0);
    return message;
  },
};

function createBaseValidateAndGetFromRawX509Request(): ValidateAndGetFromRawX509Request {
  return { raw: new Uint8Array(0) };
}

export const ValidateAndGetFromRawX509Request = {
  encode(message: ValidateAndGetFromRawX509Request, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.raw.length !== 0) {
      writer.uint32(10).bytes(message.raw);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidateAndGetFromRawX509Request {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidateAndGetFromRawX509Request();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.raw = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ValidateAndGetFromRawX509Request {
    return { raw: isSet(object.raw) ? bytesFromBase64(object.raw) : new Uint8Array(0) };
  },

  toJSON(message: ValidateAndGetFromRawX509Request): unknown {
    const obj: any = {};
    if (message.raw.length !== 0) {
      obj.raw = base64FromBytes(message.raw);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ValidateAndGetFromRawX509Request>, I>>(
    base?: I,
  ): ValidateAndGetFromRawX509Request {
    return ValidateAndGetFromRawX509Request.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ValidateAndGetFromRawX509Request>, I>>(
    object: I,
  ): ValidateAndGetFromRawX509Request {
    const message = createBaseValidateAndGetFromRawX509Request();
    message.raw = object.raw ?? new Uint8Array(0);
    return message;
  },
};

function createBaseValidateAndGetFromRawX509Response(): ValidateAndGetFromRawX509Response {
  return { status: 0, certificate: undefined };
}

export const ValidateAndGetFromRawX509Response = {
  encode(message: ValidateAndGetFromRawX509Response, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.status !== 0) {
      writer.uint32(8).int32(message.status);
    }
    if (message.certificate !== undefined) {
      Certificate.encode(message.certificate, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ValidateAndGetFromRawX509Response {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseValidateAndGetFromRawX509Response();
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

          message.certificate = Certificate.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ValidateAndGetFromRawX509Response {
    return {
      status: isSet(object.status) ? validateAndGetFromRawX509Response_StatusFromJSON(object.status) : 0,
      certificate: isSet(object.certificate) ? Certificate.fromJSON(object.certificate) : undefined,
    };
  },

  toJSON(message: ValidateAndGetFromRawX509Response): unknown {
    const obj: any = {};
    if (message.status !== 0) {
      obj.status = validateAndGetFromRawX509Response_StatusToJSON(message.status);
    }
    if (message.certificate !== undefined) {
      obj.certificate = Certificate.toJSON(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ValidateAndGetFromRawX509Response>, I>>(
    base?: I,
  ): ValidateAndGetFromRawX509Response {
    return ValidateAndGetFromRawX509Response.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ValidateAndGetFromRawX509Response>, I>>(
    object: I,
  ): ValidateAndGetFromRawX509Response {
    const message = createBaseValidateAndGetFromRawX509Response();
    message.status = object.status ?? 0;
    message.certificate = (object.certificate !== undefined && object.certificate !== null)
      ? Certificate.fromPartial(object.certificate)
      : undefined;
    return message;
  },
};

function createBaseGetRequest(): GetRequest {
  return { namespace: "", uuid: "" };
}

export const GetRequest = {
  encode(message: GetRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
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
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRequest>, I>>(base?: I): GetRequest {
    return GetRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRequest>, I>>(object: I): GetRequest {
    const message = createBaseGetRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseGetResponse(): GetResponse {
  return { certificate: undefined };
}

export const GetResponse = {
  encode(message: GetResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate !== undefined) {
      Certificate.encode(message.certificate, writer.uint32(10).fork()).ldelim();
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

          message.certificate = Certificate.decode(reader, reader.uint32());
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
    return { certificate: isSet(object.certificate) ? Certificate.fromJSON(object.certificate) : undefined };
  },

  toJSON(message: GetResponse): unknown {
    const obj: any = {};
    if (message.certificate !== undefined) {
      obj.certificate = Certificate.toJSON(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetResponse>, I>>(base?: I): GetResponse {
    return GetResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetResponse>, I>>(object: I): GetResponse {
    const message = createBaseGetResponse();
    message.certificate = (object.certificate !== undefined && object.certificate !== null)
      ? Certificate.fromPartial(object.certificate)
      : undefined;
    return message;
  },
};

function createBaseCountRequest(): CountRequest {
  return { namespace: "" };
}

export const CountRequest = {
  encode(message: CountRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountRequest {
    return { namespace: isSet(object.namespace) ? String(object.namespace) : "" };
  },

  toJSON(message: CountRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountRequest>, I>>(base?: I): CountRequest {
    return CountRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountRequest>, I>>(object: I): CountRequest {
    const message = createBaseCountRequest();
    message.namespace = object.namespace ?? "";
    return message;
  },
};

function createBaseCountResponse(): CountResponse {
  return { count: 0 };
}

export const CountResponse = {
  encode(message: CountResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.count = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountResponse>, I>>(base?: I): CountResponse {
    return CountResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountResponse>, I>>(object: I): CountResponse {
    const message = createBaseCountResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseListRequest(): ListRequest {
  return { namespace: "", skip: 0, limit: 0 };
}

export const ListRequest = {
  encode(message: ListRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.skip !== 0) {
      writer.uint32(16).uint64(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(24).uint64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListRequest();
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
          if (tag !== 16) {
            break;
          }

          message.skip = longToNumber(reader.uint64() as Long);
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.limit = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListRequest>, I>>(base?: I): ListRequest {
    return ListRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListRequest>, I>>(object: I): ListRequest {
    const message = createBaseListRequest();
    message.namespace = object.namespace ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListResponse(): ListResponse {
  return { certificate: undefined };
}

export const ListResponse = {
  encode(message: ListResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate !== undefined) {
      Certificate.encode(message.certificate, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = Certificate.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListResponse {
    return { certificate: isSet(object.certificate) ? Certificate.fromJSON(object.certificate) : undefined };
  },

  toJSON(message: ListResponse): unknown {
    const obj: any = {};
    if (message.certificate !== undefined) {
      obj.certificate = Certificate.toJSON(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListResponse>, I>>(base?: I): ListResponse {
    return ListResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListResponse>, I>>(object: I): ListResponse {
    const message = createBaseListResponse();
    message.certificate = (object.certificate !== undefined && object.certificate !== null)
      ? Certificate.fromPartial(object.certificate)
      : undefined;
    return message;
  },
};

function createBaseCountForIdentityRequest(): CountForIdentityRequest {
  return { namespace: "", identity: "" };
}

export const CountForIdentityRequest = {
  encode(message: CountForIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountForIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountForIdentityRequest();
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

  fromJSON(object: any): CountForIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
    };
  },

  toJSON(message: CountForIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountForIdentityRequest>, I>>(base?: I): CountForIdentityRequest {
    return CountForIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountForIdentityRequest>, I>>(object: I): CountForIdentityRequest {
    const message = createBaseCountForIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    return message;
  },
};

function createBaseCountForIdentityResponse(): CountForIdentityResponse {
  return { count: 0 };
}

export const CountForIdentityResponse = {
  encode(message: CountForIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).uint64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountForIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountForIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.count = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountForIdentityResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountForIdentityResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountForIdentityResponse>, I>>(base?: I): CountForIdentityResponse {
    return CountForIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountForIdentityResponse>, I>>(object: I): CountForIdentityResponse {
    const message = createBaseCountForIdentityResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseListForIdentityRequest(): ListForIdentityRequest {
  return { namespace: "", identity: "", skip: 0, limit: 0 };
}

export const ListForIdentityRequest = {
  encode(message: ListForIdentityRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.identity !== "") {
      writer.uint32(18).string(message.identity);
    }
    if (message.skip !== 0) {
      writer.uint32(24).uint64(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(32).uint64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListForIdentityRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListForIdentityRequest();
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

          message.skip = longToNumber(reader.uint64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.limit = longToNumber(reader.uint64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListForIdentityRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      identity: isSet(object.identity) ? String(object.identity) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListForIdentityRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.identity !== "") {
      obj.identity = message.identity;
    }
    if (message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListForIdentityRequest>, I>>(base?: I): ListForIdentityRequest {
    return ListForIdentityRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListForIdentityRequest>, I>>(object: I): ListForIdentityRequest {
    const message = createBaseListForIdentityRequest();
    message.namespace = object.namespace ?? "";
    message.identity = object.identity ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListForIdentityResponse(): ListForIdentityResponse {
  return { certificate: undefined };
}

export const ListForIdentityResponse = {
  encode(message: ListForIdentityResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate !== undefined) {
      Certificate.encode(message.certificate, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListForIdentityResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListForIdentityResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = Certificate.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListForIdentityResponse {
    return { certificate: isSet(object.certificate) ? Certificate.fromJSON(object.certificate) : undefined };
  },

  toJSON(message: ListForIdentityResponse): unknown {
    const obj: any = {};
    if (message.certificate !== undefined) {
      obj.certificate = Certificate.toJSON(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListForIdentityResponse>, I>>(base?: I): ListForIdentityResponse {
    return ListForIdentityResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListForIdentityResponse>, I>>(object: I): ListForIdentityResponse {
    const message = createBaseListForIdentityResponse();
    message.certificate = (object.certificate !== undefined && object.certificate !== null)
      ? Certificate.fromPartial(object.certificate)
      : undefined;
    return message;
  },
};

function createBaseUpdateRequest(): UpdateRequest {
  return { namespace: "", uuid: "", newDescription: "" };
}

export const UpdateRequest = {
  encode(message: UpdateRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.newDescription !== "") {
      writer.uint32(42).string(message.newDescription);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateRequest();
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
        case 5:
          if (tag !== 42) {
            break;
          }

          message.newDescription = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      newDescription: isSet(object.newDescription) ? String(object.newDescription) : "",
    };
  },

  toJSON(message: UpdateRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.newDescription !== "") {
      obj.newDescription = message.newDescription;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateRequest>, I>>(base?: I): UpdateRequest {
    return UpdateRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateRequest>, I>>(object: I): UpdateRequest {
    const message = createBaseUpdateRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.newDescription = object.newDescription ?? "";
    return message;
  },
};

function createBaseUpdateResponse(): UpdateResponse {
  return { certificate: undefined };
}

export const UpdateResponse = {
  encode(message: UpdateResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.certificate !== undefined) {
      Certificate.encode(message.certificate, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.certificate = Certificate.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateResponse {
    return { certificate: isSet(object.certificate) ? Certificate.fromJSON(object.certificate) : undefined };
  },

  toJSON(message: UpdateResponse): unknown {
    const obj: any = {};
    if (message.certificate !== undefined) {
      obj.certificate = Certificate.toJSON(message.certificate);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateResponse>, I>>(base?: I): UpdateResponse {
    return UpdateResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateResponse>, I>>(object: I): UpdateResponse {
    const message = createBaseUpdateResponse();
    message.certificate = (object.certificate !== undefined && object.certificate !== null)
      ? Certificate.fromPartial(object.certificate)
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
  return { wasActive: false };
}

export const DisableResponse = {
  encode(message: DisableResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.wasActive === true) {
      writer.uint32(8).bool(message.wasActive);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DisableResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDisableResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.wasActive = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DisableResponse {
    return { wasActive: isSet(object.wasActive) ? Boolean(object.wasActive) : false };
  },

  toJSON(message: DisableResponse): unknown {
    const obj: any = {};
    if (message.wasActive === true) {
      obj.wasActive = message.wasActive;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DisableResponse>, I>>(base?: I): DisableResponse {
    return DisableResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DisableResponse>, I>>(object: I): DisableResponse {
    const message = createBaseDisableResponse();
    message.wasActive = object.wasActive ?? false;
    return message;
  },
};

/** Provides API to authentificate identities using x509 certificates */
export interface IAMAuthenticationX509Service {
  /** Get root CA certificate and public key in DER format. It can be used to validate all the certificates created by this service */
  GetRootCAInfo(request: GetRootCAInfoRequest): Promise<GetRootCAInfoResponse>;
  /** Register public key for identity and generate x509 certificate for it. Sign certificate using internal CA. */
  RegisterAndGenerate(request: RegisterAndGenerateRequest): Promise<RegisterAndGenerateResponse>;
  /** Regenerate x509 certificate. Return new x509 certificate signed with CA. */
  Regenerate(request: RegenerateRequest): Promise<RegenerateResponse>;
  /** Get certificate information from RAW X509 certificate. */
  ValidateAndGetFromRawX509(request: ValidateAndGetFromRawX509Request): Promise<ValidateAndGetFromRawX509Response>;
  /** Get certificate information using its unique identifier */
  Get(request: GetRequest): Promise<GetResponse>;
  /** Count all the registered certificates in the namespace */
  Count(request: CountRequest): Promise<CountResponse>;
  /** List all the registered certificates in the namespace */
  List(request: ListRequest): Observable<ListResponse>;
  /** List all the registered certificates for specified identity */
  CountForIdentity(request: CountForIdentityRequest): Promise<CountForIdentityResponse>;
  /** List all the registered certificates for specified identity */
  ListForIdentity(request: ListForIdentityRequest): Observable<ListForIdentityResponse>;
  /** Update certificate information */
  Update(request: UpdateRequest): Promise<UpdateResponse>;
  /** Delete certificate. Note, that previously generated X509 certificate is still valid. Thats why you have to check if they still exists and wasnt disabled. */
  Delete(request: DeleteRequest): Promise<DeleteResponse>;
  /** Mark certificate as manually disabled. Disabled certificated cant be used. */
  Disable(request: DisableRequest): Promise<DisableResponse>;
}

export const IAMAuthenticationX509ServiceServiceName = "native_iam_authentication_x509.IAMAuthenticationX509Service";
export class IAMAuthenticationX509ServiceClientImpl implements IAMAuthenticationX509Service {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || IAMAuthenticationX509ServiceServiceName;
    this.rpc = rpc;
    this.GetRootCAInfo = this.GetRootCAInfo.bind(this);
    this.RegisterAndGenerate = this.RegisterAndGenerate.bind(this);
    this.Regenerate = this.Regenerate.bind(this);
    this.ValidateAndGetFromRawX509 = this.ValidateAndGetFromRawX509.bind(this);
    this.Get = this.Get.bind(this);
    this.Count = this.Count.bind(this);
    this.List = this.List.bind(this);
    this.CountForIdentity = this.CountForIdentity.bind(this);
    this.ListForIdentity = this.ListForIdentity.bind(this);
    this.Update = this.Update.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Disable = this.Disable.bind(this);
  }
  GetRootCAInfo(request: GetRootCAInfoRequest): Promise<GetRootCAInfoResponse> {
    const data = GetRootCAInfoRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetRootCAInfo", data);
    return promise.then((data) => GetRootCAInfoResponse.decode(_m0.Reader.create(data)));
  }

  RegisterAndGenerate(request: RegisterAndGenerateRequest): Promise<RegisterAndGenerateResponse> {
    const data = RegisterAndGenerateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RegisterAndGenerate", data);
    return promise.then((data) => RegisterAndGenerateResponse.decode(_m0.Reader.create(data)));
  }

  Regenerate(request: RegenerateRequest): Promise<RegenerateResponse> {
    const data = RegenerateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Regenerate", data);
    return promise.then((data) => RegenerateResponse.decode(_m0.Reader.create(data)));
  }

  ValidateAndGetFromRawX509(request: ValidateAndGetFromRawX509Request): Promise<ValidateAndGetFromRawX509Response> {
    const data = ValidateAndGetFromRawX509Request.encode(request).finish();
    const promise = this.rpc.request(this.service, "ValidateAndGetFromRawX509", data);
    return promise.then((data) => ValidateAndGetFromRawX509Response.decode(_m0.Reader.create(data)));
  }

  Get(request: GetRequest): Promise<GetResponse> {
    const data = GetRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetResponse.decode(_m0.Reader.create(data)));
  }

  Count(request: CountRequest): Promise<CountResponse> {
    const data = CountRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountResponse.decode(_m0.Reader.create(data)));
  }

  List(request: ListRequest): Observable<ListResponse> {
    const data = ListRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListResponse.decode(_m0.Reader.create(data))));
  }

  CountForIdentity(request: CountForIdentityRequest): Promise<CountForIdentityResponse> {
    const data = CountForIdentityRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "CountForIdentity", data);
    return promise.then((data) => CountForIdentityResponse.decode(_m0.Reader.create(data)));
  }

  ListForIdentity(request: ListForIdentityRequest): Observable<ListForIdentityResponse> {
    const data = ListForIdentityRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "ListForIdentity", data);
    return result.pipe(map((data) => ListForIdentityResponse.decode(_m0.Reader.create(data))));
  }

  Update(request: UpdateRequest): Promise<UpdateResponse> {
    const data = UpdateRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateResponse.decode(_m0.Reader.create(data)));
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
}

interface Rpc {
  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;
  clientStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Promise<Uint8Array>;
  serverStreamingRequest(service: string, method: string, data: Uint8Array): Observable<Uint8Array>;
  bidirectionalStreamingRequest(service: string, method: string, data: Observable<Uint8Array>): Observable<Uint8Array>;
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

function longToNumber(long: Long): number {
  if (long.gt(Number.MAX_SAFE_INTEGER)) {
    throw new tsProtoGlobalThis.Error("Value is larger than Number.MAX_SAFE_INTEGER");
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
