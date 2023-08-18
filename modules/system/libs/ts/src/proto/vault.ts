/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";

export const protobufPackage = "system_vault";

export enum RSASignMechanism {
  /** DEFAULT - System will use best algorithm for HSM in use */
  DEFAULT = 0,
  SHA256_RSA = 1,
  SHA512_RSA = 2,
  RSA_PKCS = 3,
  UNRECOGNIZED = -1,
}

export function rSASignMechanismFromJSON(object: any): RSASignMechanism {
  switch (object) {
    case 0:
    case "DEFAULT":
      return RSASignMechanism.DEFAULT;
    case 1:
    case "SHA256_RSA":
      return RSASignMechanism.SHA256_RSA;
    case 2:
    case "SHA512_RSA":
      return RSASignMechanism.SHA512_RSA;
    case 3:
    case "RSA_PKCS":
      return RSASignMechanism.RSA_PKCS;
    case -1:
    case "UNRECOGNIZED":
    default:
      return RSASignMechanism.UNRECOGNIZED;
  }
}

export function rSASignMechanismToJSON(object: RSASignMechanism): string {
  switch (object) {
    case RSASignMechanism.DEFAULT:
      return "DEFAULT";
    case RSASignMechanism.SHA256_RSA:
      return "SHA256_RSA";
    case RSASignMechanism.SHA512_RSA:
      return "SHA512_RSA";
    case RSASignMechanism.RSA_PKCS:
      return "RSA_PKCS";
    case RSASignMechanism.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface SealRequest {
}

export interface SealResponse {
}

export interface UnsealRequest {
  /** Secret used for decrypting the vault */
  secret: string;
}

export interface UnsealResponse {
  /** Indicates if vault was unsealed or not */
  success: boolean;
}

export interface UpdateSealSecretRequest {
  /** Current administrator password. Administrator acess needed to modify secret. Administrator password will not be saved in the system. */
  currentAdminPassword: string;
  /** New administrator password. Administrator password will not be saved in the system. */
  newAdminPassword: string;
  /** New secret. */
  newSecret: string;
}

export interface UpdateSealSecretResponse {
}

export interface GetStatusRequest {
}

export interface GetStatusResponse {
  /** Is vault sealed or not */
  sealed: boolean;
}

export interface EnsureRSAKeyPairRequest {
  /** Unique name of the key pair */
  keyName: string;
}

export interface EnsureRSAKeyPairResponse {
}

export interface GetRSAPublicKeyRequest {
  /** Name of the key pair for which to get key */
  keyName: string;
}

export interface GetRSAPublicKeyResponse {
  /** Public key of the RSA key pair in PKCS #1, ASN.1 DER format */
  publicKey: Uint8Array;
}

export interface RSASignStreamRequest {
  /** Name of the key pair to use */
  keyName: string;
  /** Data chunk to sign */
  data: Uint8Array;
  /** Mechanism to use. If not sure - leave default (SHA512_RSA) */
  mechanism: RSASignMechanism;
}

export interface RSASignStreamResponse {
  /** Signature of the provided data */
  signature: Uint8Array;
}

export interface RSAVerifyStreamRequest {
  /** Name of the key pair to use for validation */
  keyName: string;
  /** Data chunk to validate */
  data: Uint8Array;
  /** Signature to validate */
  signature: Uint8Array;
  /** Mechanism to use. If not sure - leave default (SHA512_RSA) */
  mechanism: RSASignMechanism;
}

export interface RSAVerifyStreamResponse {
  /** Returns True if and only if provided data and its signature matches provided key-pair */
  valid: boolean;
}

export interface RSASignRequest {
  /** Name of the key pair to use */
  keyName: string;
  /** Data to sign */
  data: Uint8Array;
  /** Mechanism to use. If not sure - leave default (SHA512_RSA) */
  mechanism: RSASignMechanism;
}

export interface RSASignResponse {
  /** Signature of the provided data */
  signature: Uint8Array;
}

export interface RSAVerifyRequest {
  /** Name of the key pair to use for validation */
  keyName: string;
  /** Data to validate */
  data: Uint8Array;
  /** Signature to validate */
  signature: Uint8Array;
  /** Mechanism to use. If not sure - leave default (SHA512_RSA) */
  mechanism: RSASignMechanism;
}

export interface RSAVerifyResponse {
  /** Returns True if and only if provided data and its signature matches provided key-pair */
  valid: boolean;
}

export interface HMACSignStreamRequest {
  /** Block of data to sign */
  data: Uint8Array;
}

export interface HMACSignStreamResponse {
  /** Signature of the provided data */
  signature: Uint8Array;
}

export interface HMACVerifyStreamRequest {
  /** Block of data to sign */
  data: Uint8Array;
  /** Signature to validate */
  signature: Uint8Array;
}

export interface HMACVerifyStreamResponse {
  /** Returns True if and only if provided data and its signature is valid */
  valid: boolean;
}

export interface HMACSignRequest {
  /** Data to sign */
  data: Uint8Array;
}

export interface HMACSignResponse {
  /** Signature of the provided data */
  signature: Uint8Array;
}

export interface HMACVerifyRequest {
  /** Data to sign */
  data: Uint8Array;
  /** Signature to validate */
  signature: Uint8Array;
}

export interface HMACVerifyResponse {
  /** Returns True if and only if provided data and its signature is valid */
  valid: boolean;
}

export interface EncryptStreamRequest {
  /** Block of data to encrypt */
  plainData: Uint8Array;
}

export interface EncryptStreamResponse {
  /** Encrypted block of the data */
  encryptedData: Uint8Array;
}

export interface DecryptStreamRequest {
  /** Encrypted block of the data */
  encryptedData: Uint8Array;
}

export interface DecryptStreamResponse {
  /** Decrypted block of the data */
  plainData: Uint8Array;
}

export interface EncryptRequest {
  /** Data to encrypt */
  plainData: Uint8Array;
}

export interface EncryptResponse {
  /** Encrypted data */
  encryptedData: Uint8Array;
}

export interface DecryptRequest {
  /** Encrypted data */
  encryptedData: Uint8Array;
}

export interface DecryptResponse {
  /** Decrypted data */
  plainData: Uint8Array;
}

function createBaseSealRequest(): SealRequest {
  return {};
}

export const SealRequest = {
  encode(_: SealRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SealRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSealRequest();
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

  fromJSON(_: any): SealRequest {
    return {};
  },

  toJSON(_: SealRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<SealRequest>, I>>(base?: I): SealRequest {
    return SealRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SealRequest>, I>>(_: I): SealRequest {
    const message = createBaseSealRequest();
    return message;
  },
};

function createBaseSealResponse(): SealResponse {
  return {};
}

export const SealResponse = {
  encode(_: SealResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): SealResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseSealResponse();
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

  fromJSON(_: any): SealResponse {
    return {};
  },

  toJSON(_: SealResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<SealResponse>, I>>(base?: I): SealResponse {
    return SealResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<SealResponse>, I>>(_: I): SealResponse {
    const message = createBaseSealResponse();
    return message;
  },
};

function createBaseUnsealRequest(): UnsealRequest {
  return { secret: "" };
}

export const UnsealRequest = {
  encode(message: UnsealRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.secret !== "") {
      writer.uint32(10).string(message.secret);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnsealRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnsealRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.secret = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UnsealRequest {
    return { secret: isSet(object.secret) ? String(object.secret) : "" };
  },

  toJSON(message: UnsealRequest): unknown {
    const obj: any = {};
    if (message.secret !== "") {
      obj.secret = message.secret;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UnsealRequest>, I>>(base?: I): UnsealRequest {
    return UnsealRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UnsealRequest>, I>>(object: I): UnsealRequest {
    const message = createBaseUnsealRequest();
    message.secret = object.secret ?? "";
    return message;
  },
};

function createBaseUnsealResponse(): UnsealResponse {
  return { success: false };
}

export const UnsealResponse = {
  encode(message: UnsealResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.success === true) {
      writer.uint32(8).bool(message.success);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UnsealResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUnsealResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.success = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UnsealResponse {
    return { success: isSet(object.success) ? Boolean(object.success) : false };
  },

  toJSON(message: UnsealResponse): unknown {
    const obj: any = {};
    if (message.success === true) {
      obj.success = message.success;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UnsealResponse>, I>>(base?: I): UnsealResponse {
    return UnsealResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UnsealResponse>, I>>(object: I): UnsealResponse {
    const message = createBaseUnsealResponse();
    message.success = object.success ?? false;
    return message;
  },
};

function createBaseUpdateSealSecretRequest(): UpdateSealSecretRequest {
  return { currentAdminPassword: "", newAdminPassword: "", newSecret: "" };
}

export const UpdateSealSecretRequest = {
  encode(message: UpdateSealSecretRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.currentAdminPassword !== "") {
      writer.uint32(10).string(message.currentAdminPassword);
    }
    if (message.newAdminPassword !== "") {
      writer.uint32(18).string(message.newAdminPassword);
    }
    if (message.newSecret !== "") {
      writer.uint32(26).string(message.newSecret);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSealSecretRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSealSecretRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.currentAdminPassword = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.newAdminPassword = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.newSecret = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateSealSecretRequest {
    return {
      currentAdminPassword: isSet(object.currentAdminPassword) ? String(object.currentAdminPassword) : "",
      newAdminPassword: isSet(object.newAdminPassword) ? String(object.newAdminPassword) : "",
      newSecret: isSet(object.newSecret) ? String(object.newSecret) : "",
    };
  },

  toJSON(message: UpdateSealSecretRequest): unknown {
    const obj: any = {};
    if (message.currentAdminPassword !== "") {
      obj.currentAdminPassword = message.currentAdminPassword;
    }
    if (message.newAdminPassword !== "") {
      obj.newAdminPassword = message.newAdminPassword;
    }
    if (message.newSecret !== "") {
      obj.newSecret = message.newSecret;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateSealSecretRequest>, I>>(base?: I): UpdateSealSecretRequest {
    return UpdateSealSecretRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateSealSecretRequest>, I>>(object: I): UpdateSealSecretRequest {
    const message = createBaseUpdateSealSecretRequest();
    message.currentAdminPassword = object.currentAdminPassword ?? "";
    message.newAdminPassword = object.newAdminPassword ?? "";
    message.newSecret = object.newSecret ?? "";
    return message;
  },
};

function createBaseUpdateSealSecretResponse(): UpdateSealSecretResponse {
  return {};
}

export const UpdateSealSecretResponse = {
  encode(_: UpdateSealSecretResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateSealSecretResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateSealSecretResponse();
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

  fromJSON(_: any): UpdateSealSecretResponse {
    return {};
  },

  toJSON(_: UpdateSealSecretResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateSealSecretResponse>, I>>(base?: I): UpdateSealSecretResponse {
    return UpdateSealSecretResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateSealSecretResponse>, I>>(_: I): UpdateSealSecretResponse {
    const message = createBaseUpdateSealSecretResponse();
    return message;
  },
};

function createBaseGetStatusRequest(): GetStatusRequest {
  return {};
}

export const GetStatusRequest = {
  encode(_: GetStatusRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetStatusRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetStatusRequest();
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

  fromJSON(_: any): GetStatusRequest {
    return {};
  },

  toJSON(_: GetStatusRequest): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<GetStatusRequest>, I>>(base?: I): GetStatusRequest {
    return GetStatusRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetStatusRequest>, I>>(_: I): GetStatusRequest {
    const message = createBaseGetStatusRequest();
    return message;
  },
};

function createBaseGetStatusResponse(): GetStatusResponse {
  return { sealed: false };
}

export const GetStatusResponse = {
  encode(message: GetStatusResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.sealed === true) {
      writer.uint32(8).bool(message.sealed);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetStatusResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetStatusResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.sealed = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetStatusResponse {
    return { sealed: isSet(object.sealed) ? Boolean(object.sealed) : false };
  },

  toJSON(message: GetStatusResponse): unknown {
    const obj: any = {};
    if (message.sealed === true) {
      obj.sealed = message.sealed;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetStatusResponse>, I>>(base?: I): GetStatusResponse {
    return GetStatusResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetStatusResponse>, I>>(object: I): GetStatusResponse {
    const message = createBaseGetStatusResponse();
    message.sealed = object.sealed ?? false;
    return message;
  },
};

function createBaseEnsureRSAKeyPairRequest(): EnsureRSAKeyPairRequest {
  return { keyName: "" };
}

export const EnsureRSAKeyPairRequest = {
  encode(message: EnsureRSAKeyPairRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureRSAKeyPairRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureRSAKeyPairRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EnsureRSAKeyPairRequest {
    return { keyName: isSet(object.keyName) ? String(object.keyName) : "" };
  },

  toJSON(message: EnsureRSAKeyPairRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureRSAKeyPairRequest>, I>>(base?: I): EnsureRSAKeyPairRequest {
    return EnsureRSAKeyPairRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureRSAKeyPairRequest>, I>>(object: I): EnsureRSAKeyPairRequest {
    const message = createBaseEnsureRSAKeyPairRequest();
    message.keyName = object.keyName ?? "";
    return message;
  },
};

function createBaseEnsureRSAKeyPairResponse(): EnsureRSAKeyPairResponse {
  return {};
}

export const EnsureRSAKeyPairResponse = {
  encode(_: EnsureRSAKeyPairResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureRSAKeyPairResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureRSAKeyPairResponse();
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

  fromJSON(_: any): EnsureRSAKeyPairResponse {
    return {};
  },

  toJSON(_: EnsureRSAKeyPairResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureRSAKeyPairResponse>, I>>(base?: I): EnsureRSAKeyPairResponse {
    return EnsureRSAKeyPairResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureRSAKeyPairResponse>, I>>(_: I): EnsureRSAKeyPairResponse {
    const message = createBaseEnsureRSAKeyPairResponse();
    return message;
  },
};

function createBaseGetRSAPublicKeyRequest(): GetRSAPublicKeyRequest {
  return { keyName: "" };
}

export const GetRSAPublicKeyRequest = {
  encode(message: GetRSAPublicKeyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRSAPublicKeyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRSAPublicKeyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRSAPublicKeyRequest {
    return { keyName: isSet(object.keyName) ? String(object.keyName) : "" };
  },

  toJSON(message: GetRSAPublicKeyRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRSAPublicKeyRequest>, I>>(base?: I): GetRSAPublicKeyRequest {
    return GetRSAPublicKeyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRSAPublicKeyRequest>, I>>(object: I): GetRSAPublicKeyRequest {
    const message = createBaseGetRSAPublicKeyRequest();
    message.keyName = object.keyName ?? "";
    return message;
  },
};

function createBaseGetRSAPublicKeyResponse(): GetRSAPublicKeyResponse {
  return { publicKey: new Uint8Array(0) };
}

export const GetRSAPublicKeyResponse = {
  encode(message: GetRSAPublicKeyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.publicKey.length !== 0) {
      writer.uint32(10).bytes(message.publicKey);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetRSAPublicKeyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetRSAPublicKeyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.publicKey = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetRSAPublicKeyResponse {
    return { publicKey: isSet(object.publicKey) ? bytesFromBase64(object.publicKey) : new Uint8Array(0) };
  },

  toJSON(message: GetRSAPublicKeyResponse): unknown {
    const obj: any = {};
    if (message.publicKey.length !== 0) {
      obj.publicKey = base64FromBytes(message.publicKey);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetRSAPublicKeyResponse>, I>>(base?: I): GetRSAPublicKeyResponse {
    return GetRSAPublicKeyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetRSAPublicKeyResponse>, I>>(object: I): GetRSAPublicKeyResponse {
    const message = createBaseGetRSAPublicKeyResponse();
    message.publicKey = object.publicKey ?? new Uint8Array(0);
    return message;
  },
};

function createBaseRSASignStreamRequest(): RSASignStreamRequest {
  return { keyName: "", data: new Uint8Array(0), mechanism: 0 };
}

export const RSASignStreamRequest = {
  encode(message: RSASignStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    if (message.mechanism !== 0) {
      writer.uint32(24).int32(message.mechanism);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSASignStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSASignStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.mechanism = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSASignStreamRequest {
    return {
      keyName: isSet(object.keyName) ? String(object.keyName) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      mechanism: isSet(object.mechanism) ? rSASignMechanismFromJSON(object.mechanism) : 0,
    };
  },

  toJSON(message: RSASignStreamRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.mechanism !== 0) {
      obj.mechanism = rSASignMechanismToJSON(message.mechanism);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSASignStreamRequest>, I>>(base?: I): RSASignStreamRequest {
    return RSASignStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSASignStreamRequest>, I>>(object: I): RSASignStreamRequest {
    const message = createBaseRSASignStreamRequest();
    message.keyName = object.keyName ?? "";
    message.data = object.data ?? new Uint8Array(0);
    message.mechanism = object.mechanism ?? 0;
    return message;
  },
};

function createBaseRSASignStreamResponse(): RSASignStreamResponse {
  return { signature: new Uint8Array(0) };
}

export const RSASignStreamResponse = {
  encode(message: RSASignStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.signature.length !== 0) {
      writer.uint32(10).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSASignStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSASignStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSASignStreamResponse {
    return { signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0) };
  },

  toJSON(message: RSASignStreamResponse): unknown {
    const obj: any = {};
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSASignStreamResponse>, I>>(base?: I): RSASignStreamResponse {
    return RSASignStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSASignStreamResponse>, I>>(object: I): RSASignStreamResponse {
    const message = createBaseRSASignStreamResponse();
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseRSAVerifyStreamRequest(): RSAVerifyStreamRequest {
  return { keyName: "", data: new Uint8Array(0), signature: new Uint8Array(0), mechanism: 0 };
}

export const RSAVerifyStreamRequest = {
  encode(message: RSAVerifyStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    if (message.signature.length !== 0) {
      writer.uint32(26).bytes(message.signature);
    }
    if (message.mechanism !== 0) {
      writer.uint32(32).int32(message.mechanism);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSAVerifyStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSAVerifyStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.signature = reader.bytes();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.mechanism = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSAVerifyStreamRequest {
    return {
      keyName: isSet(object.keyName) ? String(object.keyName) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0),
      mechanism: isSet(object.mechanism) ? rSASignMechanismFromJSON(object.mechanism) : 0,
    };
  },

  toJSON(message: RSAVerifyStreamRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    if (message.mechanism !== 0) {
      obj.mechanism = rSASignMechanismToJSON(message.mechanism);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSAVerifyStreamRequest>, I>>(base?: I): RSAVerifyStreamRequest {
    return RSAVerifyStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSAVerifyStreamRequest>, I>>(object: I): RSAVerifyStreamRequest {
    const message = createBaseRSAVerifyStreamRequest();
    message.keyName = object.keyName ?? "";
    message.data = object.data ?? new Uint8Array(0);
    message.signature = object.signature ?? new Uint8Array(0);
    message.mechanism = object.mechanism ?? 0;
    return message;
  },
};

function createBaseRSAVerifyStreamResponse(): RSAVerifyStreamResponse {
  return { valid: false };
}

export const RSAVerifyStreamResponse = {
  encode(message: RSAVerifyStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valid === true) {
      writer.uint32(8).bool(message.valid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSAVerifyStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSAVerifyStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.valid = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSAVerifyStreamResponse {
    return { valid: isSet(object.valid) ? Boolean(object.valid) : false };
  },

  toJSON(message: RSAVerifyStreamResponse): unknown {
    const obj: any = {};
    if (message.valid === true) {
      obj.valid = message.valid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSAVerifyStreamResponse>, I>>(base?: I): RSAVerifyStreamResponse {
    return RSAVerifyStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSAVerifyStreamResponse>, I>>(object: I): RSAVerifyStreamResponse {
    const message = createBaseRSAVerifyStreamResponse();
    message.valid = object.valid ?? false;
    return message;
  },
};

function createBaseRSASignRequest(): RSASignRequest {
  return { keyName: "", data: new Uint8Array(0), mechanism: 0 };
}

export const RSASignRequest = {
  encode(message: RSASignRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    if (message.mechanism !== 0) {
      writer.uint32(24).int32(message.mechanism);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSASignRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSASignRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.mechanism = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSASignRequest {
    return {
      keyName: isSet(object.keyName) ? String(object.keyName) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      mechanism: isSet(object.mechanism) ? rSASignMechanismFromJSON(object.mechanism) : 0,
    };
  },

  toJSON(message: RSASignRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.mechanism !== 0) {
      obj.mechanism = rSASignMechanismToJSON(message.mechanism);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSASignRequest>, I>>(base?: I): RSASignRequest {
    return RSASignRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSASignRequest>, I>>(object: I): RSASignRequest {
    const message = createBaseRSASignRequest();
    message.keyName = object.keyName ?? "";
    message.data = object.data ?? new Uint8Array(0);
    message.mechanism = object.mechanism ?? 0;
    return message;
  },
};

function createBaseRSASignResponse(): RSASignResponse {
  return { signature: new Uint8Array(0) };
}

export const RSASignResponse = {
  encode(message: RSASignResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.signature.length !== 0) {
      writer.uint32(10).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSASignResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSASignResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSASignResponse {
    return { signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0) };
  },

  toJSON(message: RSASignResponse): unknown {
    const obj: any = {};
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSASignResponse>, I>>(base?: I): RSASignResponse {
    return RSASignResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSASignResponse>, I>>(object: I): RSASignResponse {
    const message = createBaseRSASignResponse();
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseRSAVerifyRequest(): RSAVerifyRequest {
  return { keyName: "", data: new Uint8Array(0), signature: new Uint8Array(0), mechanism: 0 };
}

export const RSAVerifyRequest = {
  encode(message: RSAVerifyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.keyName !== "") {
      writer.uint32(10).string(message.keyName);
    }
    if (message.data.length !== 0) {
      writer.uint32(18).bytes(message.data);
    }
    if (message.signature.length !== 0) {
      writer.uint32(26).bytes(message.signature);
    }
    if (message.mechanism !== 0) {
      writer.uint32(32).int32(message.mechanism);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSAVerifyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSAVerifyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.keyName = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.signature = reader.bytes();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.mechanism = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSAVerifyRequest {
    return {
      keyName: isSet(object.keyName) ? String(object.keyName) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0),
      mechanism: isSet(object.mechanism) ? rSASignMechanismFromJSON(object.mechanism) : 0,
    };
  },

  toJSON(message: RSAVerifyRequest): unknown {
    const obj: any = {};
    if (message.keyName !== "") {
      obj.keyName = message.keyName;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    if (message.mechanism !== 0) {
      obj.mechanism = rSASignMechanismToJSON(message.mechanism);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSAVerifyRequest>, I>>(base?: I): RSAVerifyRequest {
    return RSAVerifyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSAVerifyRequest>, I>>(object: I): RSAVerifyRequest {
    const message = createBaseRSAVerifyRequest();
    message.keyName = object.keyName ?? "";
    message.data = object.data ?? new Uint8Array(0);
    message.signature = object.signature ?? new Uint8Array(0);
    message.mechanism = object.mechanism ?? 0;
    return message;
  },
};

function createBaseRSAVerifyResponse(): RSAVerifyResponse {
  return { valid: false };
}

export const RSAVerifyResponse = {
  encode(message: RSAVerifyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valid === true) {
      writer.uint32(8).bool(message.valid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RSAVerifyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRSAVerifyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.valid = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RSAVerifyResponse {
    return { valid: isSet(object.valid) ? Boolean(object.valid) : false };
  },

  toJSON(message: RSAVerifyResponse): unknown {
    const obj: any = {};
    if (message.valid === true) {
      obj.valid = message.valid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RSAVerifyResponse>, I>>(base?: I): RSAVerifyResponse {
    return RSAVerifyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RSAVerifyResponse>, I>>(object: I): RSAVerifyResponse {
    const message = createBaseRSAVerifyResponse();
    message.valid = object.valid ?? false;
    return message;
  },
};

function createBaseHMACSignStreamRequest(): HMACSignStreamRequest {
  return { data: new Uint8Array(0) };
}

export const HMACSignStreamRequest = {
  encode(message: HMACSignStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACSignStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACSignStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACSignStreamRequest {
    return { data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0) };
  },

  toJSON(message: HMACSignStreamRequest): unknown {
    const obj: any = {};
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACSignStreamRequest>, I>>(base?: I): HMACSignStreamRequest {
    return HMACSignStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACSignStreamRequest>, I>>(object: I): HMACSignStreamRequest {
    const message = createBaseHMACSignStreamRequest();
    message.data = object.data ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACSignStreamResponse(): HMACSignStreamResponse {
  return { signature: new Uint8Array(0) };
}

export const HMACSignStreamResponse = {
  encode(message: HMACSignStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.signature.length !== 0) {
      writer.uint32(10).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACSignStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACSignStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACSignStreamResponse {
    return { signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0) };
  },

  toJSON(message: HMACSignStreamResponse): unknown {
    const obj: any = {};
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACSignStreamResponse>, I>>(base?: I): HMACSignStreamResponse {
    return HMACSignStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACSignStreamResponse>, I>>(object: I): HMACSignStreamResponse {
    const message = createBaseHMACSignStreamResponse();
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACVerifyStreamRequest(): HMACVerifyStreamRequest {
  return { data: new Uint8Array(0), signature: new Uint8Array(0) };
}

export const HMACVerifyStreamRequest = {
  encode(message: HMACVerifyStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    if (message.signature.length !== 0) {
      writer.uint32(18).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACVerifyStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACVerifyStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACVerifyStreamRequest {
    return {
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0),
    };
  },

  toJSON(message: HMACVerifyStreamRequest): unknown {
    const obj: any = {};
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACVerifyStreamRequest>, I>>(base?: I): HMACVerifyStreamRequest {
    return HMACVerifyStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACVerifyStreamRequest>, I>>(object: I): HMACVerifyStreamRequest {
    const message = createBaseHMACVerifyStreamRequest();
    message.data = object.data ?? new Uint8Array(0);
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACVerifyStreamResponse(): HMACVerifyStreamResponse {
  return { valid: false };
}

export const HMACVerifyStreamResponse = {
  encode(message: HMACVerifyStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valid === true) {
      writer.uint32(8).bool(message.valid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACVerifyStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACVerifyStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.valid = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACVerifyStreamResponse {
    return { valid: isSet(object.valid) ? Boolean(object.valid) : false };
  },

  toJSON(message: HMACVerifyStreamResponse): unknown {
    const obj: any = {};
    if (message.valid === true) {
      obj.valid = message.valid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACVerifyStreamResponse>, I>>(base?: I): HMACVerifyStreamResponse {
    return HMACVerifyStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACVerifyStreamResponse>, I>>(object: I): HMACVerifyStreamResponse {
    const message = createBaseHMACVerifyStreamResponse();
    message.valid = object.valid ?? false;
    return message;
  },
};

function createBaseHMACSignRequest(): HMACSignRequest {
  return { data: new Uint8Array(0) };
}

export const HMACSignRequest = {
  encode(message: HMACSignRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACSignRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACSignRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACSignRequest {
    return { data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0) };
  },

  toJSON(message: HMACSignRequest): unknown {
    const obj: any = {};
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACSignRequest>, I>>(base?: I): HMACSignRequest {
    return HMACSignRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACSignRequest>, I>>(object: I): HMACSignRequest {
    const message = createBaseHMACSignRequest();
    message.data = object.data ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACSignResponse(): HMACSignResponse {
  return { signature: new Uint8Array(0) };
}

export const HMACSignResponse = {
  encode(message: HMACSignResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.signature.length !== 0) {
      writer.uint32(10).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACSignResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACSignResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACSignResponse {
    return { signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0) };
  },

  toJSON(message: HMACSignResponse): unknown {
    const obj: any = {};
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACSignResponse>, I>>(base?: I): HMACSignResponse {
    return HMACSignResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACSignResponse>, I>>(object: I): HMACSignResponse {
    const message = createBaseHMACSignResponse();
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACVerifyRequest(): HMACVerifyRequest {
  return { data: new Uint8Array(0), signature: new Uint8Array(0) };
}

export const HMACVerifyRequest = {
  encode(message: HMACVerifyRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(10).bytes(message.data);
    }
    if (message.signature.length !== 0) {
      writer.uint32(18).bytes(message.signature);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACVerifyRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACVerifyRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.signature = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACVerifyRequest {
    return {
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      signature: isSet(object.signature) ? bytesFromBase64(object.signature) : new Uint8Array(0),
    };
  },

  toJSON(message: HMACVerifyRequest): unknown {
    const obj: any = {};
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    if (message.signature.length !== 0) {
      obj.signature = base64FromBytes(message.signature);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACVerifyRequest>, I>>(base?: I): HMACVerifyRequest {
    return HMACVerifyRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACVerifyRequest>, I>>(object: I): HMACVerifyRequest {
    const message = createBaseHMACVerifyRequest();
    message.data = object.data ?? new Uint8Array(0);
    message.signature = object.signature ?? new Uint8Array(0);
    return message;
  },
};

function createBaseHMACVerifyResponse(): HMACVerifyResponse {
  return { valid: false };
}

export const HMACVerifyResponse = {
  encode(message: HMACVerifyResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.valid === true) {
      writer.uint32(8).bool(message.valid);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): HMACVerifyResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseHMACVerifyResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.valid = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): HMACVerifyResponse {
    return { valid: isSet(object.valid) ? Boolean(object.valid) : false };
  },

  toJSON(message: HMACVerifyResponse): unknown {
    const obj: any = {};
    if (message.valid === true) {
      obj.valid = message.valid;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<HMACVerifyResponse>, I>>(base?: I): HMACVerifyResponse {
    return HMACVerifyResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<HMACVerifyResponse>, I>>(object: I): HMACVerifyResponse {
    const message = createBaseHMACVerifyResponse();
    message.valid = object.valid ?? false;
    return message;
  },
};

function createBaseEncryptStreamRequest(): EncryptStreamRequest {
  return { plainData: new Uint8Array(0) };
}

export const EncryptStreamRequest = {
  encode(message: EncryptStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.plainData.length !== 0) {
      writer.uint32(10).bytes(message.plainData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EncryptStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEncryptStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.plainData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EncryptStreamRequest {
    return { plainData: isSet(object.plainData) ? bytesFromBase64(object.plainData) : new Uint8Array(0) };
  },

  toJSON(message: EncryptStreamRequest): unknown {
    const obj: any = {};
    if (message.plainData.length !== 0) {
      obj.plainData = base64FromBytes(message.plainData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EncryptStreamRequest>, I>>(base?: I): EncryptStreamRequest {
    return EncryptStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EncryptStreamRequest>, I>>(object: I): EncryptStreamRequest {
    const message = createBaseEncryptStreamRequest();
    message.plainData = object.plainData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseEncryptStreamResponse(): EncryptStreamResponse {
  return { encryptedData: new Uint8Array(0) };
}

export const EncryptStreamResponse = {
  encode(message: EncryptStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.encryptedData.length !== 0) {
      writer.uint32(10).bytes(message.encryptedData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EncryptStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEncryptStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.encryptedData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EncryptStreamResponse {
    return { encryptedData: isSet(object.encryptedData) ? bytesFromBase64(object.encryptedData) : new Uint8Array(0) };
  },

  toJSON(message: EncryptStreamResponse): unknown {
    const obj: any = {};
    if (message.encryptedData.length !== 0) {
      obj.encryptedData = base64FromBytes(message.encryptedData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EncryptStreamResponse>, I>>(base?: I): EncryptStreamResponse {
    return EncryptStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EncryptStreamResponse>, I>>(object: I): EncryptStreamResponse {
    const message = createBaseEncryptStreamResponse();
    message.encryptedData = object.encryptedData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseDecryptStreamRequest(): DecryptStreamRequest {
  return { encryptedData: new Uint8Array(0) };
}

export const DecryptStreamRequest = {
  encode(message: DecryptStreamRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.encryptedData.length !== 0) {
      writer.uint32(10).bytes(message.encryptedData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DecryptStreamRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDecryptStreamRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.encryptedData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DecryptStreamRequest {
    return { encryptedData: isSet(object.encryptedData) ? bytesFromBase64(object.encryptedData) : new Uint8Array(0) };
  },

  toJSON(message: DecryptStreamRequest): unknown {
    const obj: any = {};
    if (message.encryptedData.length !== 0) {
      obj.encryptedData = base64FromBytes(message.encryptedData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DecryptStreamRequest>, I>>(base?: I): DecryptStreamRequest {
    return DecryptStreamRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DecryptStreamRequest>, I>>(object: I): DecryptStreamRequest {
    const message = createBaseDecryptStreamRequest();
    message.encryptedData = object.encryptedData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseDecryptStreamResponse(): DecryptStreamResponse {
  return { plainData: new Uint8Array(0) };
}

export const DecryptStreamResponse = {
  encode(message: DecryptStreamResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.plainData.length !== 0) {
      writer.uint32(10).bytes(message.plainData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DecryptStreamResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDecryptStreamResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.plainData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DecryptStreamResponse {
    return { plainData: isSet(object.plainData) ? bytesFromBase64(object.plainData) : new Uint8Array(0) };
  },

  toJSON(message: DecryptStreamResponse): unknown {
    const obj: any = {};
    if (message.plainData.length !== 0) {
      obj.plainData = base64FromBytes(message.plainData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DecryptStreamResponse>, I>>(base?: I): DecryptStreamResponse {
    return DecryptStreamResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DecryptStreamResponse>, I>>(object: I): DecryptStreamResponse {
    const message = createBaseDecryptStreamResponse();
    message.plainData = object.plainData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseEncryptRequest(): EncryptRequest {
  return { plainData: new Uint8Array(0) };
}

export const EncryptRequest = {
  encode(message: EncryptRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.plainData.length !== 0) {
      writer.uint32(10).bytes(message.plainData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EncryptRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEncryptRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.plainData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EncryptRequest {
    return { plainData: isSet(object.plainData) ? bytesFromBase64(object.plainData) : new Uint8Array(0) };
  },

  toJSON(message: EncryptRequest): unknown {
    const obj: any = {};
    if (message.plainData.length !== 0) {
      obj.plainData = base64FromBytes(message.plainData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EncryptRequest>, I>>(base?: I): EncryptRequest {
    return EncryptRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EncryptRequest>, I>>(object: I): EncryptRequest {
    const message = createBaseEncryptRequest();
    message.plainData = object.plainData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseEncryptResponse(): EncryptResponse {
  return { encryptedData: new Uint8Array(0) };
}

export const EncryptResponse = {
  encode(message: EncryptResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.encryptedData.length !== 0) {
      writer.uint32(10).bytes(message.encryptedData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EncryptResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEncryptResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.encryptedData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EncryptResponse {
    return { encryptedData: isSet(object.encryptedData) ? bytesFromBase64(object.encryptedData) : new Uint8Array(0) };
  },

  toJSON(message: EncryptResponse): unknown {
    const obj: any = {};
    if (message.encryptedData.length !== 0) {
      obj.encryptedData = base64FromBytes(message.encryptedData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EncryptResponse>, I>>(base?: I): EncryptResponse {
    return EncryptResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EncryptResponse>, I>>(object: I): EncryptResponse {
    const message = createBaseEncryptResponse();
    message.encryptedData = object.encryptedData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseDecryptRequest(): DecryptRequest {
  return { encryptedData: new Uint8Array(0) };
}

export const DecryptRequest = {
  encode(message: DecryptRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.encryptedData.length !== 0) {
      writer.uint32(10).bytes(message.encryptedData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DecryptRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDecryptRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.encryptedData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DecryptRequest {
    return { encryptedData: isSet(object.encryptedData) ? bytesFromBase64(object.encryptedData) : new Uint8Array(0) };
  },

  toJSON(message: DecryptRequest): unknown {
    const obj: any = {};
    if (message.encryptedData.length !== 0) {
      obj.encryptedData = base64FromBytes(message.encryptedData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DecryptRequest>, I>>(base?: I): DecryptRequest {
    return DecryptRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DecryptRequest>, I>>(object: I): DecryptRequest {
    const message = createBaseDecryptRequest();
    message.encryptedData = object.encryptedData ?? new Uint8Array(0);
    return message;
  },
};

function createBaseDecryptResponse(): DecryptResponse {
  return { plainData: new Uint8Array(0) };
}

export const DecryptResponse = {
  encode(message: DecryptResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.plainData.length !== 0) {
      writer.uint32(10).bytes(message.plainData);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DecryptResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDecryptResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.plainData = reader.bytes();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DecryptResponse {
    return { plainData: isSet(object.plainData) ? bytesFromBase64(object.plainData) : new Uint8Array(0) };
  },

  toJSON(message: DecryptResponse): unknown {
    const obj: any = {};
    if (message.plainData.length !== 0) {
      obj.plainData = base64FromBytes(message.plainData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DecryptResponse>, I>>(base?: I): DecryptResponse {
    return DecryptResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DecryptResponse>, I>>(object: I): DecryptResponse {
    const message = createBaseDecryptResponse();
    message.plainData = object.plainData ?? new Uint8Array(0);
    return message;
  },
};

export interface VaultService {
  /** Close and encrypt vault. After sealing, most of the operations will not be accessible. */
  Seal(request: SealRequest): Promise<SealResponse>;
  /** Decrypt and open vault. Must be done before most of the operations with vault secrets. */
  Unseal(request: UnsealRequest): Promise<UnsealResponse>;
  /**
   * Set up new seal secret and reincrypt vault. The vault must be unsealed before this operation. You don't need to unseal vault after this operation.
   * This operation requires you to have administrator access to the HSM. Check PKCS11 spec. If you are using emulated HSM (by default) this will be the same as the seal/unseal secret by default ("12345678"). Change it.
   */
  UpdateSealSecret(request: UpdateSealSecretRequest): Promise<UpdateSealSecretResponse>;
  /** Get current status of the vault. */
  GetStatus(request: GetStatusRequest): Promise<GetStatusResponse>;
  /** Creates RSA key pair if it doesnt exist. Private key never leaves the HSM (hardware security module). */
  EnsureRSAKeyPair(request: EnsureRSAKeyPairRequest): Promise<EnsureRSAKeyPairResponse>;
  /** Get public key of the RSA keypair. */
  GetRSAPublicKey(request: GetRSAPublicKeyRequest): Promise<GetRSAPublicKeyResponse>;
  /** Sign message stream with RSA. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message. */
  RSASignStream(request: Observable<RSASignStreamRequest>): Promise<RSASignStreamResponse>;
  /** Validate signature of the message stream using RSA key-pairs public key. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message. */
  RSAVerifyStream(request: Observable<RSAVerifyStreamRequest>): Promise<RSAVerifyStreamResponse>;
  /** Sign message with RSA. The data must be short (max several kilobytes). If i is longer - use `RSASignStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to sign the message. */
  RSASign(request: RSASignRequest): Promise<RSASignResponse>;
  /** Validate signature of the message using RSA key-pairs public key. The data must be short (max several kilobytes). If i is longer - use `RSAVerifyStream` instead. It will use SHA512_RSA_PKCS (RS512) algorithm to verify the message. */
  RSAVerify(request: RSAVerifyRequest): Promise<RSAVerifyResponse>;
  /** Calculate HMAC signature for input data stream. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM. */
  HMACSignStream(request: Observable<HMACSignStreamRequest>): Promise<HMACSignStreamResponse>;
  /** Verify HMAC signature of the data stream. */
  HMACVerifyStream(request: Observable<HMACVerifyStreamRequest>): Promise<HMACVerifyStreamResponse>;
  /** Calculate HMAC signature for input data. The data must be short (max several kilobytes). If it is longer - use `HMACSignStream` instead. HMAC secret never leaves the HSM (hardware security module). It automatically uses the best available HMAC algorithm for currently used HSM. */
  HMACSign(request: HMACSignRequest): Promise<HMACSignResponse>;
  /** Verify HMAC signature of the data. The data must be short (max several kilobytes). If it is longer - use `HMACVerifyStream` instead. */
  HMACVerify(request: HMACVerifyRequest): Promise<HMACVerifyResponse>;
  /** Encrypt data stream. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption). */
  EncryptStream(request: Observable<EncryptStreamRequest>): Observable<EncryptStreamResponse>;
  /** Decrypt data stream. If you want to decrypt part of the information from the middle of the whole data - ensure, that the first chunk of the data, that you are sending is padded by 2048 bit. */
  DecryptStream(request: Observable<DecryptStreamRequest>): Observable<DecryptStreamResponse>;
  /** Encrypt data. The data must be short (max several kilobytes). If it is longer - use `EncryptStream` instead. Encryption secret never leaves the HSM (hardware security module). It automatically uses the best available encryption algorithm for currently used HSM. It will only select the algorithm that is capable of decrypting from the middle of the whole data (partial decryption). */
  Encrypt(request: EncryptRequest): Promise<EncryptResponse>;
  /** Decrypt data. The data must be short (max several kilobytes). If it is longer - use `DecryptStream` instead. */
  Decrypt(request: DecryptRequest): Promise<DecryptResponse>;
}

export const VaultServiceServiceName = "system_vault.VaultService";
export class VaultServiceClientImpl implements VaultService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || VaultServiceServiceName;
    this.rpc = rpc;
    this.Seal = this.Seal.bind(this);
    this.Unseal = this.Unseal.bind(this);
    this.UpdateSealSecret = this.UpdateSealSecret.bind(this);
    this.GetStatus = this.GetStatus.bind(this);
    this.EnsureRSAKeyPair = this.EnsureRSAKeyPair.bind(this);
    this.GetRSAPublicKey = this.GetRSAPublicKey.bind(this);
    this.RSASignStream = this.RSASignStream.bind(this);
    this.RSAVerifyStream = this.RSAVerifyStream.bind(this);
    this.RSASign = this.RSASign.bind(this);
    this.RSAVerify = this.RSAVerify.bind(this);
    this.HMACSignStream = this.HMACSignStream.bind(this);
    this.HMACVerifyStream = this.HMACVerifyStream.bind(this);
    this.HMACSign = this.HMACSign.bind(this);
    this.HMACVerify = this.HMACVerify.bind(this);
    this.EncryptStream = this.EncryptStream.bind(this);
    this.DecryptStream = this.DecryptStream.bind(this);
    this.Encrypt = this.Encrypt.bind(this);
    this.Decrypt = this.Decrypt.bind(this);
  }
  Seal(request: SealRequest): Promise<SealResponse> {
    const data = SealRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Seal", data);
    return promise.then((data) => SealResponse.decode(_m0.Reader.create(data)));
  }

  Unseal(request: UnsealRequest): Promise<UnsealResponse> {
    const data = UnsealRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Unseal", data);
    return promise.then((data) => UnsealResponse.decode(_m0.Reader.create(data)));
  }

  UpdateSealSecret(request: UpdateSealSecretRequest): Promise<UpdateSealSecretResponse> {
    const data = UpdateSealSecretRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "UpdateSealSecret", data);
    return promise.then((data) => UpdateSealSecretResponse.decode(_m0.Reader.create(data)));
  }

  GetStatus(request: GetStatusRequest): Promise<GetStatusResponse> {
    const data = GetStatusRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetStatus", data);
    return promise.then((data) => GetStatusResponse.decode(_m0.Reader.create(data)));
  }

  EnsureRSAKeyPair(request: EnsureRSAKeyPairRequest): Promise<EnsureRSAKeyPairResponse> {
    const data = EnsureRSAKeyPairRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "EnsureRSAKeyPair", data);
    return promise.then((data) => EnsureRSAKeyPairResponse.decode(_m0.Reader.create(data)));
  }

  GetRSAPublicKey(request: GetRSAPublicKeyRequest): Promise<GetRSAPublicKeyResponse> {
    const data = GetRSAPublicKeyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetRSAPublicKey", data);
    return promise.then((data) => GetRSAPublicKeyResponse.decode(_m0.Reader.create(data)));
  }

  RSASignStream(request: Observable<RSASignStreamRequest>): Promise<RSASignStreamResponse> {
    const data = request.pipe(map((request) => RSASignStreamRequest.encode(request).finish()));
    const promise = this.rpc.clientStreamingRequest(this.service, "RSASignStream", data);
    return promise.then((data) => RSASignStreamResponse.decode(_m0.Reader.create(data)));
  }

  RSAVerifyStream(request: Observable<RSAVerifyStreamRequest>): Promise<RSAVerifyStreamResponse> {
    const data = request.pipe(map((request) => RSAVerifyStreamRequest.encode(request).finish()));
    const promise = this.rpc.clientStreamingRequest(this.service, "RSAVerifyStream", data);
    return promise.then((data) => RSAVerifyStreamResponse.decode(_m0.Reader.create(data)));
  }

  RSASign(request: RSASignRequest): Promise<RSASignResponse> {
    const data = RSASignRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RSASign", data);
    return promise.then((data) => RSASignResponse.decode(_m0.Reader.create(data)));
  }

  RSAVerify(request: RSAVerifyRequest): Promise<RSAVerifyResponse> {
    const data = RSAVerifyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RSAVerify", data);
    return promise.then((data) => RSAVerifyResponse.decode(_m0.Reader.create(data)));
  }

  HMACSignStream(request: Observable<HMACSignStreamRequest>): Promise<HMACSignStreamResponse> {
    const data = request.pipe(map((request) => HMACSignStreamRequest.encode(request).finish()));
    const promise = this.rpc.clientStreamingRequest(this.service, "HMACSignStream", data);
    return promise.then((data) => HMACSignStreamResponse.decode(_m0.Reader.create(data)));
  }

  HMACVerifyStream(request: Observable<HMACVerifyStreamRequest>): Promise<HMACVerifyStreamResponse> {
    const data = request.pipe(map((request) => HMACVerifyStreamRequest.encode(request).finish()));
    const promise = this.rpc.clientStreamingRequest(this.service, "HMACVerifyStream", data);
    return promise.then((data) => HMACVerifyStreamResponse.decode(_m0.Reader.create(data)));
  }

  HMACSign(request: HMACSignRequest): Promise<HMACSignResponse> {
    const data = HMACSignRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "HMACSign", data);
    return promise.then((data) => HMACSignResponse.decode(_m0.Reader.create(data)));
  }

  HMACVerify(request: HMACVerifyRequest): Promise<HMACVerifyResponse> {
    const data = HMACVerifyRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "HMACVerify", data);
    return promise.then((data) => HMACVerifyResponse.decode(_m0.Reader.create(data)));
  }

  EncryptStream(request: Observable<EncryptStreamRequest>): Observable<EncryptStreamResponse> {
    const data = request.pipe(map((request) => EncryptStreamRequest.encode(request).finish()));
    const result = this.rpc.bidirectionalStreamingRequest(this.service, "EncryptStream", data);
    return result.pipe(map((data) => EncryptStreamResponse.decode(_m0.Reader.create(data))));
  }

  DecryptStream(request: Observable<DecryptStreamRequest>): Observable<DecryptStreamResponse> {
    const data = request.pipe(map((request) => DecryptStreamRequest.encode(request).finish()));
    const result = this.rpc.bidirectionalStreamingRequest(this.service, "DecryptStream", data);
    return result.pipe(map((data) => DecryptStreamResponse.decode(_m0.Reader.create(data))));
  }

  Encrypt(request: EncryptRequest): Promise<EncryptResponse> {
    const data = EncryptRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Encrypt", data);
    return promise.then((data) => EncryptResponse.decode(_m0.Reader.create(data)));
  }

  Decrypt(request: DecryptRequest): Promise<DecryptResponse> {
    const data = DecryptRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Decrypt", data);
    return promise.then((data) => DecryptResponse.decode(_m0.Reader.create(data)));
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

function isSet(value: any): boolean {
  return value !== null && value !== undefined;
}
