/* eslint-disable */
import Long from "long";
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { Timestamp } from "./google/protobuf/timestamp";
import { map } from "rxjs/operators";

export const protobufPackage = "native_files";

/** Represents file information (metadata) that is directly binded to the binary data */
export interface File {
  /** Namespace where file is located */
  namespace: string;
  /** Unique file id */
  uuid: string;
  /** Indicates if file can be modified after creation. Making file readonly allows caching. In most cases its better to reacreate file on changes, than make it writeable. */
  readonly: boolean;
  /** Mime type of the file */
  mimeType: string;
  /** Size of the file in bytes */
  size: number;
  /** Indicates if SHA512 hash for this file was calculated or not. When file is readoly, hash calculates on creation. For non readonly files, calculation of the hash has to be manually trigerred. */
  SHA512HashCalculated: boolean;
  /** Unique SHA512 hash of the data in the file */
  SHA512Hash: Buffer;
  /** Disables caching for this file. Use it when you think that file will not be used often, so caching of this file will not pollute cache servers. If file is not readonly, does nothing. */
  disableCache: boolean;
  /** Forcelly apply cache. By default, service will not cache files with size more, than 50 megabytes. If you think, that this file will be used very often, you can force service to always cache it. Remember, that in this case it will use space on cache service, to other important files may not be cached. This parameter doesnt have effect when file is writeable or caching is disabled. */
  forceCaching: boolean;
  /** When file was creted */
  Created: Date | undefined;
  /** When file was updated last time */
  Updated: Date | undefined;
  /** Version of the file. Automatically increases on every update to the file binary data or file information */
  Version: number;
}

export interface FileCreateRequest {
  info: FileCreateRequest_FileInfo | undefined;
  chunk: FileCreateRequest_FileChunk | undefined;
}

export interface FileCreateRequest_FileInfo {
  /** Namespace where file is located */
  namespace: string;
  /** Indicates if file can be modified after creation. Making file readonly allows caching. In most cases its better to reacreate file on changes, than make it writeable. */
  readonly: boolean;
  /** Mime type of the file */
  mimeType: string;
  /** Disables caching for this file. Use it when you think that file will not be used often, so caching of this file will not pollute cache servers. If file is not readonly, does nothing. */
  disableCache: boolean;
  /** Forcelly apply cache. By default, service will not cache files with size more, than 50 megabytes. If you think, that this file will be used very often, you can force service to always cache it. Remember, that in this case it will use space on cache service, to other important files may not be cached. This parameter doesnt have effect when file is writeable or caching is disabled. */
  forceCaching: boolean;
}

export interface FileCreateRequest_FileChunk {
  data: Buffer;
}

export interface FileCreateResponse {
  /** Created file */
  file: File | undefined;
}

export interface FileExistRequest {
  /** File namespace */
  namespace: string;
  /** File identifier in the namespace */
  uuid: string;
  /** Use cache or not. This cache is not the same as file data cache. This cache can be invalid under very rare circumstances (Race condition can occure on reading and updating at same time). Cache automatically invalidates after 60 seconds */
  useCache: boolean;
}

export interface FileExistResponse {
  /** True if file exists, else false */
  exist: boolean;
}

export interface StatFileRequest {
  /** File namespace */
  namespace: string;
  /** File identifier in the namespace */
  uuid: string;
  /** Use cache or not. This cache is not the same as file data cache. This cache can be invalid under very rare circumstances (Race condition can occure on reading and updating at same time). Cache automatically invalidates after 60 seconds */
  useCache: boolean;
}

export interface StatFileResponse {
  file: File | undefined;
}

export interface ReadFileRequest {
  /** File namespace */
  namespace: string;
  /** File identifier in the namespace */
  uuid: string;
  /** Location in the file (in bytes) where to start read data. 0 to start from the begining of the file */
  start: number;
  /** Number of bytes to read. 0 to read up to the end */
  toRead: number;
}

export interface ReadFileResponse {
  /** Total ammount of data (in bytes) to be transfered. May differ from file size if transfer started from the middle of the file. */
  totalSize: number;
  /** Total ammount of data already transfered including current chunk */
  transfered: number;
  /** Starting index of the chunk in the original file. */
  chunkStart: number;
  /** Chunk of data */
  chunk: Buffer;
}

export interface CalculateFileSHA512Request {
  /** File namespace */
  namespace: string;
  /** File identifier in the namespace */
  uuid: string;
}

export interface CalculateFileSHA512Response {
  /** Calculated SHA512 */
  SHA512: Buffer;
}

export interface DeleteFileRequest {
  /** File namespace */
  namespace: string;
  /** File identifier in the namespace */
  uuid: string;
}

export interface DeleteFileResponse {}

function createBaseFile(): File {
  return {
    namespace: "",
    uuid: "",
    readonly: false,
    mimeType: "",
    size: 0,
    SHA512HashCalculated: false,
    SHA512Hash: Buffer.alloc(0),
    disableCache: false,
    forceCaching: false,
    Created: undefined,
    Updated: undefined,
    Version: 0,
  };
}

export const File = {
  encode(message: File, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.readonly === true) {
      writer.uint32(24).bool(message.readonly);
    }
    if (message.mimeType !== "") {
      writer.uint32(34).string(message.mimeType);
    }
    if (message.size !== 0) {
      writer.uint32(40).uint64(message.size);
    }
    if (message.SHA512HashCalculated === true) {
      writer.uint32(48).bool(message.SHA512HashCalculated);
    }
    if (message.SHA512Hash.length !== 0) {
      writer.uint32(58).bytes(message.SHA512Hash);
    }
    if (message.disableCache === true) {
      writer.uint32(64).bool(message.disableCache);
    }
    if (message.forceCaching === true) {
      writer.uint32(72).bool(message.forceCaching);
    }
    if (message.Created !== undefined) {
      Timestamp.encode(
        toTimestamp(message.Created),
        writer.uint32(802).fork()
      ).ldelim();
    }
    if (message.Updated !== undefined) {
      Timestamp.encode(
        toTimestamp(message.Updated),
        writer.uint32(810).fork()
      ).ldelim();
    }
    if (message.Version !== 0) {
      writer.uint32(816).int64(message.Version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): File {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFile();
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
          message.readonly = reader.bool();
          break;
        case 4:
          message.mimeType = reader.string();
          break;
        case 5:
          message.size = longToNumber(reader.uint64() as Long);
          break;
        case 6:
          message.SHA512HashCalculated = reader.bool();
          break;
        case 7:
          message.SHA512Hash = reader.bytes() as Buffer;
          break;
        case 8:
          message.disableCache = reader.bool();
          break;
        case 9:
          message.forceCaching = reader.bool();
          break;
        case 100:
          message.Created = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 101:
          message.Updated = fromTimestamp(
            Timestamp.decode(reader, reader.uint32())
          );
          break;
        case 102:
          message.Version = longToNumber(reader.int64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): File {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      readonly: isSet(object.readonly) ? Boolean(object.readonly) : false,
      mimeType: isSet(object.mimeType) ? String(object.mimeType) : "",
      size: isSet(object.size) ? Number(object.size) : 0,
      SHA512HashCalculated: isSet(object.SHA512HashCalculated)
        ? Boolean(object.SHA512HashCalculated)
        : false,
      SHA512Hash: isSet(object.SHA512Hash)
        ? Buffer.from(bytesFromBase64(object.SHA512Hash))
        : Buffer.alloc(0),
      disableCache: isSet(object.disableCache)
        ? Boolean(object.disableCache)
        : false,
      forceCaching: isSet(object.forceCaching)
        ? Boolean(object.forceCaching)
        : false,
      Created: isSet(object.Created)
        ? fromJsonTimestamp(object.Created)
        : undefined,
      Updated: isSet(object.Updated)
        ? fromJsonTimestamp(object.Updated)
        : undefined,
      Version: isSet(object.Version) ? Number(object.Version) : 0,
    };
  },

  toJSON(message: File): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.readonly !== undefined && (obj.readonly = message.readonly);
    message.mimeType !== undefined && (obj.mimeType = message.mimeType);
    message.size !== undefined && (obj.size = Math.round(message.size));
    message.SHA512HashCalculated !== undefined &&
      (obj.SHA512HashCalculated = message.SHA512HashCalculated);
    message.SHA512Hash !== undefined &&
      (obj.SHA512Hash = base64FromBytes(
        message.SHA512Hash !== undefined ? message.SHA512Hash : Buffer.alloc(0)
      ));
    message.disableCache !== undefined &&
      (obj.disableCache = message.disableCache);
    message.forceCaching !== undefined &&
      (obj.forceCaching = message.forceCaching);
    message.Created !== undefined &&
      (obj.Created = message.Created.toISOString());
    message.Updated !== undefined &&
      (obj.Updated = message.Updated.toISOString());
    message.Version !== undefined &&
      (obj.Version = Math.round(message.Version));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<File>, I>>(object: I): File {
    const message = createBaseFile();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.readonly = object.readonly ?? false;
    message.mimeType = object.mimeType ?? "";
    message.size = object.size ?? 0;
    message.SHA512HashCalculated = object.SHA512HashCalculated ?? false;
    message.SHA512Hash = object.SHA512Hash ?? Buffer.alloc(0);
    message.disableCache = object.disableCache ?? false;
    message.forceCaching = object.forceCaching ?? false;
    message.Created = object.Created ?? undefined;
    message.Updated = object.Updated ?? undefined;
    message.Version = object.Version ?? 0;
    return message;
  },
};

function createBaseFileCreateRequest(): FileCreateRequest {
  return { info: undefined, chunk: undefined };
}

export const FileCreateRequest = {
  encode(
    message: FileCreateRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.info !== undefined) {
      FileCreateRequest_FileInfo.encode(
        message.info,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.chunk !== undefined) {
      FileCreateRequest_FileChunk.encode(
        message.chunk,
        writer.uint32(18).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileCreateRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileCreateRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.info = FileCreateRequest_FileInfo.decode(
            reader,
            reader.uint32()
          );
          break;
        case 2:
          message.chunk = FileCreateRequest_FileChunk.decode(
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

  fromJSON(object: any): FileCreateRequest {
    return {
      info: isSet(object.info)
        ? FileCreateRequest_FileInfo.fromJSON(object.info)
        : undefined,
      chunk: isSet(object.chunk)
        ? FileCreateRequest_FileChunk.fromJSON(object.chunk)
        : undefined,
    };
  },

  toJSON(message: FileCreateRequest): unknown {
    const obj: any = {};
    message.info !== undefined &&
      (obj.info = message.info
        ? FileCreateRequest_FileInfo.toJSON(message.info)
        : undefined);
    message.chunk !== undefined &&
      (obj.chunk = message.chunk
        ? FileCreateRequest_FileChunk.toJSON(message.chunk)
        : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileCreateRequest>, I>>(
    object: I
  ): FileCreateRequest {
    const message = createBaseFileCreateRequest();
    message.info =
      object.info !== undefined && object.info !== null
        ? FileCreateRequest_FileInfo.fromPartial(object.info)
        : undefined;
    message.chunk =
      object.chunk !== undefined && object.chunk !== null
        ? FileCreateRequest_FileChunk.fromPartial(object.chunk)
        : undefined;
    return message;
  },
};

function createBaseFileCreateRequest_FileInfo(): FileCreateRequest_FileInfo {
  return {
    namespace: "",
    readonly: false,
    mimeType: "",
    disableCache: false,
    forceCaching: false,
  };
}

export const FileCreateRequest_FileInfo = {
  encode(
    message: FileCreateRequest_FileInfo,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(90).string(message.namespace);
    }
    if (message.readonly === true) {
      writer.uint32(96).bool(message.readonly);
    }
    if (message.mimeType !== "") {
      writer.uint32(106).string(message.mimeType);
    }
    if (message.disableCache === true) {
      writer.uint32(112).bool(message.disableCache);
    }
    if (message.forceCaching === true) {
      writer.uint32(120).bool(message.forceCaching);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): FileCreateRequest_FileInfo {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileCreateRequest_FileInfo();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 11:
          message.namespace = reader.string();
          break;
        case 12:
          message.readonly = reader.bool();
          break;
        case 13:
          message.mimeType = reader.string();
          break;
        case 14:
          message.disableCache = reader.bool();
          break;
        case 15:
          message.forceCaching = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileCreateRequest_FileInfo {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      readonly: isSet(object.readonly) ? Boolean(object.readonly) : false,
      mimeType: isSet(object.mimeType) ? String(object.mimeType) : "",
      disableCache: isSet(object.disableCache)
        ? Boolean(object.disableCache)
        : false,
      forceCaching: isSet(object.forceCaching)
        ? Boolean(object.forceCaching)
        : false,
    };
  },

  toJSON(message: FileCreateRequest_FileInfo): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.readonly !== undefined && (obj.readonly = message.readonly);
    message.mimeType !== undefined && (obj.mimeType = message.mimeType);
    message.disableCache !== undefined &&
      (obj.disableCache = message.disableCache);
    message.forceCaching !== undefined &&
      (obj.forceCaching = message.forceCaching);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileCreateRequest_FileInfo>, I>>(
    object: I
  ): FileCreateRequest_FileInfo {
    const message = createBaseFileCreateRequest_FileInfo();
    message.namespace = object.namespace ?? "";
    message.readonly = object.readonly ?? false;
    message.mimeType = object.mimeType ?? "";
    message.disableCache = object.disableCache ?? false;
    message.forceCaching = object.forceCaching ?? false;
    return message;
  },
};

function createBaseFileCreateRequest_FileChunk(): FileCreateRequest_FileChunk {
  return { data: Buffer.alloc(0) };
}

export const FileCreateRequest_FileChunk = {
  encode(
    message: FileCreateRequest_FileChunk,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.data.length !== 0) {
      writer.uint32(170).bytes(message.data);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): FileCreateRequest_FileChunk {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileCreateRequest_FileChunk();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 21:
          message.data = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileCreateRequest_FileChunk {
    return {
      data: isSet(object.data)
        ? Buffer.from(bytesFromBase64(object.data))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: FileCreateRequest_FileChunk): unknown {
    const obj: any = {};
    message.data !== undefined &&
      (obj.data = base64FromBytes(
        message.data !== undefined ? message.data : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileCreateRequest_FileChunk>, I>>(
    object: I
  ): FileCreateRequest_FileChunk {
    const message = createBaseFileCreateRequest_FileChunk();
    message.data = object.data ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseFileCreateResponse(): FileCreateResponse {
  return { file: undefined };
}

export const FileCreateResponse = {
  encode(
    message: FileCreateResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.file !== undefined) {
      File.encode(message.file, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileCreateResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileCreateResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.file = File.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileCreateResponse {
    return {
      file: isSet(object.file) ? File.fromJSON(object.file) : undefined,
    };
  },

  toJSON(message: FileCreateResponse): unknown {
    const obj: any = {};
    message.file !== undefined &&
      (obj.file = message.file ? File.toJSON(message.file) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileCreateResponse>, I>>(
    object: I
  ): FileCreateResponse {
    const message = createBaseFileCreateResponse();
    message.file =
      object.file !== undefined && object.file !== null
        ? File.fromPartial(object.file)
        : undefined;
    return message;
  },
};

function createBaseFileExistRequest(): FileExistRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const FileExistRequest = {
  encode(
    message: FileExistRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileExistRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileExistRequest();
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
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FileExistRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: FileExistRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileExistRequest>, I>>(
    object: I
  ): FileExistRequest {
    const message = createBaseFileExistRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseFileExistResponse(): FileExistResponse {
  return { exist: false };
}

export const FileExistResponse = {
  encode(
    message: FileExistResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.exist === true) {
      writer.uint32(8).bool(message.exist);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FileExistResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFileExistResponse();
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

  fromJSON(object: any): FileExistResponse {
    return {
      exist: isSet(object.exist) ? Boolean(object.exist) : false,
    };
  },

  toJSON(message: FileExistResponse): unknown {
    const obj: any = {};
    message.exist !== undefined && (obj.exist = message.exist);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<FileExistResponse>, I>>(
    object: I
  ): FileExistResponse {
    const message = createBaseFileExistResponse();
    message.exist = object.exist ?? false;
    return message;
  },
};

function createBaseStatFileRequest(): StatFileRequest {
  return { namespace: "", uuid: "", useCache: false };
}

export const StatFileRequest = {
  encode(
    message: StatFileRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatFileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatFileRequest();
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
          message.useCache = reader.bool();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StatFileRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: StatFileRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.useCache !== undefined && (obj.useCache = message.useCache);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StatFileRequest>, I>>(
    object: I
  ): StatFileRequest {
    const message = createBaseStatFileRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseStatFileResponse(): StatFileResponse {
  return { file: undefined };
}

export const StatFileResponse = {
  encode(
    message: StatFileResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.file !== undefined) {
      File.encode(message.file, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): StatFileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseStatFileResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.file = File.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): StatFileResponse {
    return {
      file: isSet(object.file) ? File.fromJSON(object.file) : undefined,
    };
  },

  toJSON(message: StatFileResponse): unknown {
    const obj: any = {};
    message.file !== undefined &&
      (obj.file = message.file ? File.toJSON(message.file) : undefined);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<StatFileResponse>, I>>(
    object: I
  ): StatFileResponse {
    const message = createBaseStatFileResponse();
    message.file =
      object.file !== undefined && object.file !== null
        ? File.fromPartial(object.file)
        : undefined;
    return message;
  },
};

function createBaseReadFileRequest(): ReadFileRequest {
  return { namespace: "", uuid: "", start: 0, toRead: 0 };
}

export const ReadFileRequest = {
  encode(
    message: ReadFileRequest,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.uuid !== "") {
      writer.uint32(18).string(message.uuid);
    }
    if (message.start !== 0) {
      writer.uint32(24).uint64(message.start);
    }
    if (message.toRead !== 0) {
      writer.uint32(32).uint64(message.toRead);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReadFileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReadFileRequest();
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
          message.start = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.toRead = longToNumber(reader.uint64() as Long);
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ReadFileRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      start: isSet(object.start) ? Number(object.start) : 0,
      toRead: isSet(object.toRead) ? Number(object.toRead) : 0,
    };
  },

  toJSON(message: ReadFileRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    message.start !== undefined && (obj.start = Math.round(message.start));
    message.toRead !== undefined && (obj.toRead = Math.round(message.toRead));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ReadFileRequest>, I>>(
    object: I
  ): ReadFileRequest {
    const message = createBaseReadFileRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    message.start = object.start ?? 0;
    message.toRead = object.toRead ?? 0;
    return message;
  },
};

function createBaseReadFileResponse(): ReadFileResponse {
  return { totalSize: 0, transfered: 0, chunkStart: 0, chunk: Buffer.alloc(0) };
}

export const ReadFileResponse = {
  encode(
    message: ReadFileResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.totalSize !== 0) {
      writer.uint32(8).uint64(message.totalSize);
    }
    if (message.transfered !== 0) {
      writer.uint32(16).uint64(message.transfered);
    }
    if (message.chunkStart !== 0) {
      writer.uint32(24).uint64(message.chunkStart);
    }
    if (message.chunk.length !== 0) {
      writer.uint32(34).bytes(message.chunk);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ReadFileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseReadFileResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.totalSize = longToNumber(reader.uint64() as Long);
          break;
        case 2:
          message.transfered = longToNumber(reader.uint64() as Long);
          break;
        case 3:
          message.chunkStart = longToNumber(reader.uint64() as Long);
          break;
        case 4:
          message.chunk = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ReadFileResponse {
    return {
      totalSize: isSet(object.totalSize) ? Number(object.totalSize) : 0,
      transfered: isSet(object.transfered) ? Number(object.transfered) : 0,
      chunkStart: isSet(object.chunkStart) ? Number(object.chunkStart) : 0,
      chunk: isSet(object.chunk)
        ? Buffer.from(bytesFromBase64(object.chunk))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: ReadFileResponse): unknown {
    const obj: any = {};
    message.totalSize !== undefined &&
      (obj.totalSize = Math.round(message.totalSize));
    message.transfered !== undefined &&
      (obj.transfered = Math.round(message.transfered));
    message.chunkStart !== undefined &&
      (obj.chunkStart = Math.round(message.chunkStart));
    message.chunk !== undefined &&
      (obj.chunk = base64FromBytes(
        message.chunk !== undefined ? message.chunk : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<ReadFileResponse>, I>>(
    object: I
  ): ReadFileResponse {
    const message = createBaseReadFileResponse();
    message.totalSize = object.totalSize ?? 0;
    message.transfered = object.transfered ?? 0;
    message.chunkStart = object.chunkStart ?? 0;
    message.chunk = object.chunk ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseCalculateFileSHA512Request(): CalculateFileSHA512Request {
  return { namespace: "", uuid: "" };
}

export const CalculateFileSHA512Request = {
  encode(
    message: CalculateFileSHA512Request,
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

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CalculateFileSHA512Request {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCalculateFileSHA512Request();
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

  fromJSON(object: any): CalculateFileSHA512Request {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: CalculateFileSHA512Request): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CalculateFileSHA512Request>, I>>(
    object: I
  ): CalculateFileSHA512Request {
    const message = createBaseCalculateFileSHA512Request();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseCalculateFileSHA512Response(): CalculateFileSHA512Response {
  return { SHA512: Buffer.alloc(0) };
}

export const CalculateFileSHA512Response = {
  encode(
    message: CalculateFileSHA512Response,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    if (message.SHA512.length !== 0) {
      writer.uint32(10).bytes(message.SHA512);
    }
    return writer;
  },

  decode(
    input: _m0.Reader | Uint8Array,
    length?: number
  ): CalculateFileSHA512Response {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCalculateFileSHA512Response();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.SHA512 = reader.bytes() as Buffer;
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CalculateFileSHA512Response {
    return {
      SHA512: isSet(object.SHA512)
        ? Buffer.from(bytesFromBase64(object.SHA512))
        : Buffer.alloc(0),
    };
  },

  toJSON(message: CalculateFileSHA512Response): unknown {
    const obj: any = {};
    message.SHA512 !== undefined &&
      (obj.SHA512 = base64FromBytes(
        message.SHA512 !== undefined ? message.SHA512 : Buffer.alloc(0)
      ));
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<CalculateFileSHA512Response>, I>>(
    object: I
  ): CalculateFileSHA512Response {
    const message = createBaseCalculateFileSHA512Response();
    message.SHA512 = object.SHA512 ?? Buffer.alloc(0);
    return message;
  },
};

function createBaseDeleteFileRequest(): DeleteFileRequest {
  return { namespace: "", uuid: "" };
}

export const DeleteFileRequest = {
  encode(
    message: DeleteFileRequest,
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

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteFileRequest {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteFileRequest();
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

  fromJSON(object: any): DeleteFileRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
    };
  },

  toJSON(message: DeleteFileRequest): unknown {
    const obj: any = {};
    message.namespace !== undefined && (obj.namespace = message.namespace);
    message.uuid !== undefined && (obj.uuid = message.uuid);
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteFileRequest>, I>>(
    object: I
  ): DeleteFileRequest {
    const message = createBaseDeleteFileRequest();
    message.namespace = object.namespace ?? "";
    message.uuid = object.uuid ?? "";
    return message;
  },
};

function createBaseDeleteFileResponse(): DeleteFileResponse {
  return {};
}

export const DeleteFileResponse = {
  encode(
    _: DeleteFileResponse,
    writer: _m0.Writer = _m0.Writer.create()
  ): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteFileResponse {
    const reader = input instanceof _m0.Reader ? input : new _m0.Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteFileResponse();
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

  fromJSON(_: any): DeleteFileResponse {
    return {};
  },

  toJSON(_: DeleteFileResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial<I extends Exact<DeepPartial<DeleteFileResponse>, I>>(
    _: I
  ): DeleteFileResponse {
    const message = createBaseDeleteFileResponse();
    return message;
  },
};

export interface FileService {
  /** Creates new file. First message in the client stream must have file information. All next messages must have bytes chunk data. Information about created file file will be returned only on EOF of client stream. */
  Create(request: Observable<FileCreateRequest>): Promise<FileCreateResponse>;
  /** Check if file at specified location exists */
  Exists(request: FileExistRequest): Promise<FileExistResponse>;
  /** Get file information. Has same meaning as POSIX stat function */
  Stat(request: StatFileRequest): Promise<StatFileResponse>;
  /** Read data from the file */
  Read(request: ReadFileRequest): Observable<ReadFileResponse>;
  /**
   * rpc Append();
   * rpc Write();
   */
  Delete(request: DeleteFileRequest): Promise<DeleteFileResponse>;
  /** Calculates SHA512 for the file data and adds it to the file information. If SHA512 was already calculated, returns stored hash value. This method is not concurrency safe. Writing to file while calculating hash will result in wrong hash stored in file information. */
  CalculateSHA512(
    request: CalculateFileSHA512Request
  ): Promise<CalculateFileSHA512Response>;
}

export class FileServiceClientImpl implements FileService {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Exists = this.Exists.bind(this);
    this.Stat = this.Stat.bind(this);
    this.Read = this.Read.bind(this);
    this.Delete = this.Delete.bind(this);
    this.CalculateSHA512 = this.CalculateSHA512.bind(this);
  }
  Create(request: Observable<FileCreateRequest>): Promise<FileCreateResponse> {
    const data = request.pipe(
      map((request) => FileCreateRequest.encode(request).finish())
    );
    const promise = this.rpc.clientStreamingRequest(
      "native_files.FileService",
      "Create",
      data
    );
    return promise.then((data) =>
      FileCreateResponse.decode(new _m0.Reader(data))
    );
  }

  Exists(request: FileExistRequest): Promise<FileExistResponse> {
    const data = FileExistRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_files.FileService",
      "Exists",
      data
    );
    return promise.then((data) =>
      FileExistResponse.decode(new _m0.Reader(data))
    );
  }

  Stat(request: StatFileRequest): Promise<StatFileResponse> {
    const data = StatFileRequest.encode(request).finish();
    const promise = this.rpc.request("native_files.FileService", "Stat", data);
    return promise.then((data) =>
      StatFileResponse.decode(new _m0.Reader(data))
    );
  }

  Read(request: ReadFileRequest): Observable<ReadFileResponse> {
    const data = ReadFileRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(
      "native_files.FileService",
      "Read",
      data
    );
    return result.pipe(
      map((data) => ReadFileResponse.decode(new _m0.Reader(data)))
    );
  }

  Delete(request: DeleteFileRequest): Promise<DeleteFileResponse> {
    const data = DeleteFileRequest.encode(request).finish();
    const promise = this.rpc.request(
      "native_files.FileService",
      "Delete",
      data
    );
    return promise.then((data) =>
      DeleteFileResponse.decode(new _m0.Reader(data))
    );
  }

  CalculateSHA512(
    request: CalculateFileSHA512Request
  ): Promise<CalculateFileSHA512Response> {
    const data = CalculateFileSHA512Request.encode(request).finish();
    const promise = this.rpc.request(
      "native_files.FileService",
      "CalculateSHA512",
      data
    );
    return promise.then((data) =>
      CalculateFileSHA512Response.decode(new _m0.Reader(data))
    );
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
  clientStreamingRequest(
    service: string,
    method: string,
    data: Observable<Uint8Array>
  ): Promise<Uint8Array>;
  serverStreamingRequest(
    service: string,
    method: string,
    data: Uint8Array
  ): Observable<Uint8Array>;
  bidirectionalStreamingRequest(
    service: string,
    method: string,
    data: Observable<Uint8Array>
  ): Observable<Uint8Array>;
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

function toTimestamp(date: Date): Timestamp {
  const seconds = date.getTime() / 1_000;
  const nanos = (date.getTime() % 1_000) * 1_000_000;
  return { seconds, nanos };
}

function fromTimestamp(t: Timestamp): Date {
  let millis = t.seconds * 1_000;
  millis += t.nanos / 1_000_000;
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
