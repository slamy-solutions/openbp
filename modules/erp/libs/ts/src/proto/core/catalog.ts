/* eslint-disable */
import * as _m0 from "protobufjs/minimal";
import { Observable } from "rxjs";
import { map } from "rxjs/operators";
import { NullValue, nullValueFromJSON, nullValueToJSON } from "../google/protobuf/struct";
import { Timestamp } from "../google/protobuf/timestamp";
import Long = require("long");

export const protobufPackage = "erp_catalog";

export interface FieldSchema {
  name: string;
  publicName: string;
  intData?: FieldSchema_IntData | undefined;
  floatData?: FieldSchema_FloatData | undefined;
  stringData?: FieldSchema_StringData | undefined;
  booleanData?: FieldSchema_BooleanData | undefined;
  tableData?: FieldSchema_TableData | undefined;
  objectData?: FieldSchema_ObjectData | undefined;
  catalogLinkData?: FieldSchema_CatalogLinkData | undefined;
}

export interface FieldSchema_IntData {
}

export interface FieldSchema_FloatData {
}

export interface FieldSchema_StringData {
}

export interface FieldSchema_BooleanData {
}

export interface FieldSchema_TableData {
  columns: FieldSchema[];
}

export interface FieldSchema_ObjectData {
  fields: FieldSchema[];
}

export interface FieldSchema_CatalogLinkData {
  catalogName: string;
}

export interface FieldSchema_CompoundCatalogLinkData {
  catalogNames: string[];
}

export interface Catalog {
  /** Namespace where catalog is located */
  namespace: string;
  /** Private short name used only for development. Can not be changed after creation and always stays the same */
  name: string;
  /** Public name. Can be changed at any time and should be used when displaying catalog to the user */
  publicName: string;
  /** Entrty fields schema */
  fields: FieldSchema[];
  /** When catalog was creted */
  created:
    | Date
    | undefined;
  /** When catalog was updated last time */
  updated:
    | Date
    | undefined;
  /** Version of the catalog. After each update, version increases by 1 */
  version: number;
}

export interface CreateCatalogRequest {
  /** Namespace to use */
  namespace: string;
  /** Private short name used only for development. Can not be changed after creation and always stays the same */
  name: string;
  /** Can be changed at any time and should be used when displaying catalog to the user */
  publicName: string;
  /** Entrty fields schema */
  fields: FieldSchema[];
}

export interface CreateCatalogResponse {
  catalog: Catalog | undefined;
}

export interface DeleteCatalogRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog name */
  name: string;
}

export interface DeleteCatalogResponse {
}

export interface UpdateCatalogRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog name */
  name: string;
  /** Can be changed at any time and should be used when displaying catalog to the user */
  publicName: string;
  /** Entrty fields schema */
  fields: FieldSchema[];
}

export interface UpdateCatalogReponse {
  catalog: Catalog | undefined;
}

export interface GetCatalogRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog name */
  name: string;
  /** Use cache or not. Greatly speeds up request if catalog is inside cache. Cache invalidates on updates, but still has very low chance to be inconsistent. Inconsistency occures on concurrent write-reads. Cocurrent reads are 100% safe. Incosistent cache will be removed after some period of time. */
  useCache: boolean;
}

export interface GetCatalogResponse {
  catalog: Catalog | undefined;
}

export interface GetCatalogIfChangedRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog name */
  name: string;
  /** Only return catalog if current version is not same */
  version: number;
  /** Use cache or not. Greatly speeds up request if catalog is inside cache. Cache invalidates on updates, but still has very low chance to be inconsistent. Inconsistency occures on concurrent write-reads. Cocurrent reads are 100% safe. Incosistent cache will be removed after some period of time. */
  useCache: boolean;
}

export interface GetCatalogIfChangedResponse {
  null?: NullValue | undefined;
  catalog?: Catalog | undefined;
}

export interface GetAllCatalogsRequest {
  /** Namespace to use */
  namespace: string;
  /** Use cache or not. Greately increases speed. Cache invalidates on updates in this namespace, but still has very low chance to be inconsistent. Inconsistency occures on concurrent write-reads. Cocurrent reads are 100% safe. Incosistent cache will be removed after some period of time. */
  useCache: boolean;
}

export interface GetAllCatalogsResponse {
  catalogs: Catalog[];
}

export interface CatalogEntry {
  /** Namespace where catalog is located */
  namespace: string;
  /** Name of the entries catalog */
  catalog: string;
  /** Unique identifier of the entry in catalog */
  uuid: string;
  /** JSON encoded entry data */
  data: Uint8Array;
  /** When entry was creted */
  created:
    | Date
    | undefined;
  /** When entry was updated last time */
  updated:
    | Date
    | undefined;
  /** Version of the entry. After each update, version increases by 1 */
  version: number;
}

export interface CreateCatalogEntryRequest {
  /** Namespace where catalog is located */
  namespace: string;
  /** Name of the catalog where create entry */
  catalog: string;
  /** JSON encoded data, that will be added on entry creation */
  data: Uint8Array;
}

export interface CreateCatalogEntryResponse {
  entry: CatalogEntry | undefined;
}

export interface DeleteCatalogEntryRequest {
  /** Namespace where catalog entry is located */
  namespace: string;
  /** Name of the catalog where delete entry */
  catalog: string;
  /** Unique identifier of the entry in this catalog */
  entry: string;
}

export interface DeleteCatalogEntryResponse {
}

export interface UpdateCatalogEntryRequest {
  entry: CatalogEntry | undefined;
}

export interface UpdateCatalogEntryResponse {
  /** Updated entry */
  entry: CatalogEntry | undefined;
}

export interface GetCatalogEntryRequest {
  /** Namespace where catalog entry is located */
  namespace: string;
  /** Name of the catalog where entry is located */
  catalog: string;
  /** Unique identifier of the entry in catalog */
  entry: string;
}

export interface GetCatalogEntryResponse {
  entry: CatalogEntry | undefined;
}

export interface ListCatalogEntriesRequest {
  /** Namespace where catalog entries are located */
  namespace: string;
  /** Name of the catalog where entries are located */
  catalog: string;
  /** How many values to skip before return. If you dont want to skip, use 0. */
  skip: number;
  /** How many values to return. Use 0 to ignore limit. */
  limit: number;
}

export interface ListCatalogEntriesResponse {
  entry: CatalogEntry | undefined;
}

export interface CountCatalogEntriesRequest {
  /** Namespace where catalog entries are located */
  namespace: string;
  /** Name of the catalog where entries are located */
  catalog: string;
}

export interface CountCatalogEntriesResponse {
  /** Number of entries in catalog */
  count: number;
}

export interface CatalogIndex {
  /** Namespace where catalog index is located */
  namespace: string;
  /** Catalog name */
  catalog: string;
  /** Unique name of the index */
  name: string;
  /** Information about field on wich index is applied */
  fields: CatalogIndex_IndexField[];
  /** Is this index unique for this catalog or not. If index unique, there will be no two entries with values that satisfy same position in index. Usefull for making unique values */
  unique: boolean;
}

export interface CatalogIndex_IndexField {
  name: string;
  type: CatalogIndex_IndexField_IndexType;
}

export enum CatalogIndex_IndexField_IndexType {
  HASHED = 0,
  ASCENDING = 1,
  DESCENDING = 2,
  UNRECOGNIZED = -1,
}

export function catalogIndex_IndexField_IndexTypeFromJSON(object: any): CatalogIndex_IndexField_IndexType {
  switch (object) {
    case 0:
    case "HASHED":
      return CatalogIndex_IndexField_IndexType.HASHED;
    case 1:
    case "ASCENDING":
      return CatalogIndex_IndexField_IndexType.ASCENDING;
    case 2:
    case "DESCENDING":
      return CatalogIndex_IndexField_IndexType.DESCENDING;
    case -1:
    case "UNRECOGNIZED":
    default:
      return CatalogIndex_IndexField_IndexType.UNRECOGNIZED;
  }
}

export function catalogIndex_IndexField_IndexTypeToJSON(object: CatalogIndex_IndexField_IndexType): string {
  switch (object) {
    case CatalogIndex_IndexField_IndexType.HASHED:
      return "HASHED";
    case CatalogIndex_IndexField_IndexType.ASCENDING:
      return "ASCENDING";
    case CatalogIndex_IndexField_IndexType.DESCENDING:
      return "DESCENDING";
    case CatalogIndex_IndexField_IndexType.UNRECOGNIZED:
    default:
      return "UNRECOGNIZED";
  }
}

export interface ListCatalogIndexesRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog name */
  catalog: string;
  /** Use cache or not. Greately increases speed. Cache invalidates on updates in this catalog, but still has very low chance to be inconsistent. Inconsistency occures on concurrent write-reads. Cocurrent reads are 100% safe. Incosistent cache will be removed after some period of time. */
  useCache: boolean;
}

export interface ListCatalogIndexesResponse {
  indexes: CatalogIndex[];
}

export interface EnsureCatalogIndexRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog on wich to ensure index */
  catalog: string;
  /** Index information */
  index: CatalogIndex | undefined;
}

export interface EnsureCatalogIndexResponse {
}

export interface RemoveCatalogIndexRequest {
  /** Namespace to use */
  namespace: string;
  /** Catalog where to remove index */
  catalog: string;
  /** Index name to remove */
  index: string;
}

export interface RemoveCatalogIndexResponse {
}

function createBaseFieldSchema(): FieldSchema {
  return {
    name: "",
    publicName: "",
    intData: undefined,
    floatData: undefined,
    stringData: undefined,
    booleanData: undefined,
    tableData: undefined,
    objectData: undefined,
    catalogLinkData: undefined,
  };
}

export const FieldSchema = {
  encode(message: FieldSchema, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.publicName !== "") {
      writer.uint32(18).string(message.publicName);
    }
    if (message.intData !== undefined) {
      FieldSchema_IntData.encode(message.intData, writer.uint32(82).fork()).ldelim();
    }
    if (message.floatData !== undefined) {
      FieldSchema_FloatData.encode(message.floatData, writer.uint32(90).fork()).ldelim();
    }
    if (message.stringData !== undefined) {
      FieldSchema_StringData.encode(message.stringData, writer.uint32(98).fork()).ldelim();
    }
    if (message.booleanData !== undefined) {
      FieldSchema_BooleanData.encode(message.booleanData, writer.uint32(106).fork()).ldelim();
    }
    if (message.tableData !== undefined) {
      FieldSchema_TableData.encode(message.tableData, writer.uint32(114).fork()).ldelim();
    }
    if (message.objectData !== undefined) {
      FieldSchema_ObjectData.encode(message.objectData, writer.uint32(122).fork()).ldelim();
    }
    if (message.catalogLinkData !== undefined) {
      FieldSchema_CatalogLinkData.encode(message.catalogLinkData, writer.uint32(130).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.publicName = reader.string();
          continue;
        case 10:
          if (tag !== 82) {
            break;
          }

          message.intData = FieldSchema_IntData.decode(reader, reader.uint32());
          continue;
        case 11:
          if (tag !== 90) {
            break;
          }

          message.floatData = FieldSchema_FloatData.decode(reader, reader.uint32());
          continue;
        case 12:
          if (tag !== 98) {
            break;
          }

          message.stringData = FieldSchema_StringData.decode(reader, reader.uint32());
          continue;
        case 13:
          if (tag !== 106) {
            break;
          }

          message.booleanData = FieldSchema_BooleanData.decode(reader, reader.uint32());
          continue;
        case 14:
          if (tag !== 114) {
            break;
          }

          message.tableData = FieldSchema_TableData.decode(reader, reader.uint32());
          continue;
        case 15:
          if (tag !== 122) {
            break;
          }

          message.objectData = FieldSchema_ObjectData.decode(reader, reader.uint32());
          continue;
        case 16:
          if (tag !== 130) {
            break;
          }

          message.catalogLinkData = FieldSchema_CatalogLinkData.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FieldSchema {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      publicName: isSet(object.publicName) ? String(object.publicName) : "",
      intData: isSet(object.intData) ? FieldSchema_IntData.fromJSON(object.intData) : undefined,
      floatData: isSet(object.floatData) ? FieldSchema_FloatData.fromJSON(object.floatData) : undefined,
      stringData: isSet(object.stringData) ? FieldSchema_StringData.fromJSON(object.stringData) : undefined,
      booleanData: isSet(object.booleanData) ? FieldSchema_BooleanData.fromJSON(object.booleanData) : undefined,
      tableData: isSet(object.tableData) ? FieldSchema_TableData.fromJSON(object.tableData) : undefined,
      objectData: isSet(object.objectData) ? FieldSchema_ObjectData.fromJSON(object.objectData) : undefined,
      catalogLinkData: isSet(object.catalogLinkData)
        ? FieldSchema_CatalogLinkData.fromJSON(object.catalogLinkData)
        : undefined,
    };
  },

  toJSON(message: FieldSchema): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.publicName !== "") {
      obj.publicName = message.publicName;
    }
    if (message.intData !== undefined) {
      obj.intData = FieldSchema_IntData.toJSON(message.intData);
    }
    if (message.floatData !== undefined) {
      obj.floatData = FieldSchema_FloatData.toJSON(message.floatData);
    }
    if (message.stringData !== undefined) {
      obj.stringData = FieldSchema_StringData.toJSON(message.stringData);
    }
    if (message.booleanData !== undefined) {
      obj.booleanData = FieldSchema_BooleanData.toJSON(message.booleanData);
    }
    if (message.tableData !== undefined) {
      obj.tableData = FieldSchema_TableData.toJSON(message.tableData);
    }
    if (message.objectData !== undefined) {
      obj.objectData = FieldSchema_ObjectData.toJSON(message.objectData);
    }
    if (message.catalogLinkData !== undefined) {
      obj.catalogLinkData = FieldSchema_CatalogLinkData.toJSON(message.catalogLinkData);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema>, I>>(base?: I): FieldSchema {
    return FieldSchema.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema>, I>>(object: I): FieldSchema {
    const message = createBaseFieldSchema();
    message.name = object.name ?? "";
    message.publicName = object.publicName ?? "";
    message.intData = (object.intData !== undefined && object.intData !== null)
      ? FieldSchema_IntData.fromPartial(object.intData)
      : undefined;
    message.floatData = (object.floatData !== undefined && object.floatData !== null)
      ? FieldSchema_FloatData.fromPartial(object.floatData)
      : undefined;
    message.stringData = (object.stringData !== undefined && object.stringData !== null)
      ? FieldSchema_StringData.fromPartial(object.stringData)
      : undefined;
    message.booleanData = (object.booleanData !== undefined && object.booleanData !== null)
      ? FieldSchema_BooleanData.fromPartial(object.booleanData)
      : undefined;
    message.tableData = (object.tableData !== undefined && object.tableData !== null)
      ? FieldSchema_TableData.fromPartial(object.tableData)
      : undefined;
    message.objectData = (object.objectData !== undefined && object.objectData !== null)
      ? FieldSchema_ObjectData.fromPartial(object.objectData)
      : undefined;
    message.catalogLinkData = (object.catalogLinkData !== undefined && object.catalogLinkData !== null)
      ? FieldSchema_CatalogLinkData.fromPartial(object.catalogLinkData)
      : undefined;
    return message;
  },
};

function createBaseFieldSchema_IntData(): FieldSchema_IntData {
  return {};
}

export const FieldSchema_IntData = {
  encode(_: FieldSchema_IntData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_IntData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_IntData();
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

  fromJSON(_: any): FieldSchema_IntData {
    return {};
  },

  toJSON(_: FieldSchema_IntData): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_IntData>, I>>(base?: I): FieldSchema_IntData {
    return FieldSchema_IntData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_IntData>, I>>(_: I): FieldSchema_IntData {
    const message = createBaseFieldSchema_IntData();
    return message;
  },
};

function createBaseFieldSchema_FloatData(): FieldSchema_FloatData {
  return {};
}

export const FieldSchema_FloatData = {
  encode(_: FieldSchema_FloatData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_FloatData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_FloatData();
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

  fromJSON(_: any): FieldSchema_FloatData {
    return {};
  },

  toJSON(_: FieldSchema_FloatData): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_FloatData>, I>>(base?: I): FieldSchema_FloatData {
    return FieldSchema_FloatData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_FloatData>, I>>(_: I): FieldSchema_FloatData {
    const message = createBaseFieldSchema_FloatData();
    return message;
  },
};

function createBaseFieldSchema_StringData(): FieldSchema_StringData {
  return {};
}

export const FieldSchema_StringData = {
  encode(_: FieldSchema_StringData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_StringData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_StringData();
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

  fromJSON(_: any): FieldSchema_StringData {
    return {};
  },

  toJSON(_: FieldSchema_StringData): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_StringData>, I>>(base?: I): FieldSchema_StringData {
    return FieldSchema_StringData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_StringData>, I>>(_: I): FieldSchema_StringData {
    const message = createBaseFieldSchema_StringData();
    return message;
  },
};

function createBaseFieldSchema_BooleanData(): FieldSchema_BooleanData {
  return {};
}

export const FieldSchema_BooleanData = {
  encode(_: FieldSchema_BooleanData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_BooleanData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_BooleanData();
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

  fromJSON(_: any): FieldSchema_BooleanData {
    return {};
  },

  toJSON(_: FieldSchema_BooleanData): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_BooleanData>, I>>(base?: I): FieldSchema_BooleanData {
    return FieldSchema_BooleanData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_BooleanData>, I>>(_: I): FieldSchema_BooleanData {
    const message = createBaseFieldSchema_BooleanData();
    return message;
  },
};

function createBaseFieldSchema_TableData(): FieldSchema_TableData {
  return { columns: [] };
}

export const FieldSchema_TableData = {
  encode(message: FieldSchema_TableData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.columns) {
      FieldSchema.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_TableData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_TableData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.columns.push(FieldSchema.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FieldSchema_TableData {
    return { columns: Array.isArray(object?.columns) ? object.columns.map((e: any) => FieldSchema.fromJSON(e)) : [] };
  },

  toJSON(message: FieldSchema_TableData): unknown {
    const obj: any = {};
    if (message.columns?.length) {
      obj.columns = message.columns.map((e) => FieldSchema.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_TableData>, I>>(base?: I): FieldSchema_TableData {
    return FieldSchema_TableData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_TableData>, I>>(object: I): FieldSchema_TableData {
    const message = createBaseFieldSchema_TableData();
    message.columns = object.columns?.map((e) => FieldSchema.fromPartial(e)) || [];
    return message;
  },
};

function createBaseFieldSchema_ObjectData(): FieldSchema_ObjectData {
  return { fields: [] };
}

export const FieldSchema_ObjectData = {
  encode(message: FieldSchema_ObjectData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.fields) {
      FieldSchema.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_ObjectData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_ObjectData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.fields.push(FieldSchema.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FieldSchema_ObjectData {
    return { fields: Array.isArray(object?.fields) ? object.fields.map((e: any) => FieldSchema.fromJSON(e)) : [] };
  },

  toJSON(message: FieldSchema_ObjectData): unknown {
    const obj: any = {};
    if (message.fields?.length) {
      obj.fields = message.fields.map((e) => FieldSchema.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_ObjectData>, I>>(base?: I): FieldSchema_ObjectData {
    return FieldSchema_ObjectData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_ObjectData>, I>>(object: I): FieldSchema_ObjectData {
    const message = createBaseFieldSchema_ObjectData();
    message.fields = object.fields?.map((e) => FieldSchema.fromPartial(e)) || [];
    return message;
  },
};

function createBaseFieldSchema_CatalogLinkData(): FieldSchema_CatalogLinkData {
  return { catalogName: "" };
}

export const FieldSchema_CatalogLinkData = {
  encode(message: FieldSchema_CatalogLinkData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.catalogName !== "") {
      writer.uint32(10).string(message.catalogName);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_CatalogLinkData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_CatalogLinkData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalogName = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FieldSchema_CatalogLinkData {
    return { catalogName: isSet(object.catalogName) ? String(object.catalogName) : "" };
  },

  toJSON(message: FieldSchema_CatalogLinkData): unknown {
    const obj: any = {};
    if (message.catalogName !== "") {
      obj.catalogName = message.catalogName;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_CatalogLinkData>, I>>(base?: I): FieldSchema_CatalogLinkData {
    return FieldSchema_CatalogLinkData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_CatalogLinkData>, I>>(object: I): FieldSchema_CatalogLinkData {
    const message = createBaseFieldSchema_CatalogLinkData();
    message.catalogName = object.catalogName ?? "";
    return message;
  },
};

function createBaseFieldSchema_CompoundCatalogLinkData(): FieldSchema_CompoundCatalogLinkData {
  return { catalogNames: [] };
}

export const FieldSchema_CompoundCatalogLinkData = {
  encode(message: FieldSchema_CompoundCatalogLinkData, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.catalogNames) {
      writer.uint32(10).string(v!);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): FieldSchema_CompoundCatalogLinkData {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseFieldSchema_CompoundCatalogLinkData();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalogNames.push(reader.string());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): FieldSchema_CompoundCatalogLinkData {
    return { catalogNames: Array.isArray(object?.catalogNames) ? object.catalogNames.map((e: any) => String(e)) : [] };
  },

  toJSON(message: FieldSchema_CompoundCatalogLinkData): unknown {
    const obj: any = {};
    if (message.catalogNames?.length) {
      obj.catalogNames = message.catalogNames;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<FieldSchema_CompoundCatalogLinkData>, I>>(
    base?: I,
  ): FieldSchema_CompoundCatalogLinkData {
    return FieldSchema_CompoundCatalogLinkData.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<FieldSchema_CompoundCatalogLinkData>, I>>(
    object: I,
  ): FieldSchema_CompoundCatalogLinkData {
    const message = createBaseFieldSchema_CompoundCatalogLinkData();
    message.catalogNames = object.catalogNames?.map((e) => e) || [];
    return message;
  },
};

function createBaseCatalog(): Catalog {
  return { namespace: "", name: "", publicName: "", fields: [], created: undefined, updated: undefined, version: 0 };
}

export const Catalog = {
  encode(message: Catalog, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.publicName !== "") {
      writer.uint32(26).string(message.publicName);
    }
    for (const v of message.fields) {
      FieldSchema.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(802).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(810).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(816).int64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): Catalog {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCatalog();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.publicName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.fields.push(FieldSchema.decode(reader, reader.uint32()));
          continue;
        case 100:
          if (tag !== 802) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 101:
          if (tag !== 810) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 102:
          if (tag !== 816) {
            break;
          }

          message.version = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): Catalog {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      publicName: isSet(object.publicName) ? String(object.publicName) : "",
      fields: Array.isArray(object?.fields) ? object.fields.map((e: any) => FieldSchema.fromJSON(e)) : [],
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: Catalog): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.publicName !== "") {
      obj.publicName = message.publicName;
    }
    if (message.fields?.length) {
      obj.fields = message.fields.map((e) => FieldSchema.toJSON(e));
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

  create<I extends Exact<DeepPartial<Catalog>, I>>(base?: I): Catalog {
    return Catalog.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<Catalog>, I>>(object: I): Catalog {
    const message = createBaseCatalog();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.publicName = object.publicName ?? "";
    message.fields = object.fields?.map((e) => FieldSchema.fromPartial(e)) || [];
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseCreateCatalogRequest(): CreateCatalogRequest {
  return { namespace: "", name: "", publicName: "", fields: [] };
}

export const CreateCatalogRequest = {
  encode(message: CreateCatalogRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.publicName !== "") {
      writer.uint32(26).string(message.publicName);
    }
    for (const v of message.fields) {
      FieldSchema.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCatalogRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCatalogRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.publicName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.fields.push(FieldSchema.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateCatalogRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      publicName: isSet(object.publicName) ? String(object.publicName) : "",
      fields: Array.isArray(object?.fields) ? object.fields.map((e: any) => FieldSchema.fromJSON(e)) : [],
    };
  },

  toJSON(message: CreateCatalogRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.publicName !== "") {
      obj.publicName = message.publicName;
    }
    if (message.fields?.length) {
      obj.fields = message.fields.map((e) => FieldSchema.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateCatalogRequest>, I>>(base?: I): CreateCatalogRequest {
    return CreateCatalogRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateCatalogRequest>, I>>(object: I): CreateCatalogRequest {
    const message = createBaseCreateCatalogRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.publicName = object.publicName ?? "";
    message.fields = object.fields?.map((e) => FieldSchema.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCreateCatalogResponse(): CreateCatalogResponse {
  return { catalog: undefined };
}

export const CreateCatalogResponse = {
  encode(message: CreateCatalogResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.catalog !== undefined) {
      Catalog.encode(message.catalog, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCatalogResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCatalogResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalog = Catalog.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateCatalogResponse {
    return { catalog: isSet(object.catalog) ? Catalog.fromJSON(object.catalog) : undefined };
  },

  toJSON(message: CreateCatalogResponse): unknown {
    const obj: any = {};
    if (message.catalog !== undefined) {
      obj.catalog = Catalog.toJSON(message.catalog);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateCatalogResponse>, I>>(base?: I): CreateCatalogResponse {
    return CreateCatalogResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateCatalogResponse>, I>>(object: I): CreateCatalogResponse {
    const message = createBaseCreateCatalogResponse();
    message.catalog = (object.catalog !== undefined && object.catalog !== null)
      ? Catalog.fromPartial(object.catalog)
      : undefined;
    return message;
  },
};

function createBaseDeleteCatalogRequest(): DeleteCatalogRequest {
  return { namespace: "", name: "" };
}

export const DeleteCatalogRequest = {
  encode(message: DeleteCatalogRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteCatalogRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteCatalogRequest();
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

          message.name = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteCatalogRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
    };
  },

  toJSON(message: DeleteCatalogRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteCatalogRequest>, I>>(base?: I): DeleteCatalogRequest {
    return DeleteCatalogRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteCatalogRequest>, I>>(object: I): DeleteCatalogRequest {
    const message = createBaseDeleteCatalogRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    return message;
  },
};

function createBaseDeleteCatalogResponse(): DeleteCatalogResponse {
  return {};
}

export const DeleteCatalogResponse = {
  encode(_: DeleteCatalogResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteCatalogResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteCatalogResponse();
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

  fromJSON(_: any): DeleteCatalogResponse {
    return {};
  },

  toJSON(_: DeleteCatalogResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteCatalogResponse>, I>>(base?: I): DeleteCatalogResponse {
    return DeleteCatalogResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteCatalogResponse>, I>>(_: I): DeleteCatalogResponse {
    const message = createBaseDeleteCatalogResponse();
    return message;
  },
};

function createBaseUpdateCatalogRequest(): UpdateCatalogRequest {
  return { namespace: "", name: "", publicName: "", fields: [] };
}

export const UpdateCatalogRequest = {
  encode(message: UpdateCatalogRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.publicName !== "") {
      writer.uint32(26).string(message.publicName);
    }
    for (const v of message.fields) {
      FieldSchema.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCatalogRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCatalogRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.publicName = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.fields.push(FieldSchema.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateCatalogRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      publicName: isSet(object.publicName) ? String(object.publicName) : "",
      fields: Array.isArray(object?.fields) ? object.fields.map((e: any) => FieldSchema.fromJSON(e)) : [],
    };
  },

  toJSON(message: UpdateCatalogRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.publicName !== "") {
      obj.publicName = message.publicName;
    }
    if (message.fields?.length) {
      obj.fields = message.fields.map((e) => FieldSchema.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateCatalogRequest>, I>>(base?: I): UpdateCatalogRequest {
    return UpdateCatalogRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateCatalogRequest>, I>>(object: I): UpdateCatalogRequest {
    const message = createBaseUpdateCatalogRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.publicName = object.publicName ?? "";
    message.fields = object.fields?.map((e) => FieldSchema.fromPartial(e)) || [];
    return message;
  },
};

function createBaseUpdateCatalogReponse(): UpdateCatalogReponse {
  return { catalog: undefined };
}

export const UpdateCatalogReponse = {
  encode(message: UpdateCatalogReponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.catalog !== undefined) {
      Catalog.encode(message.catalog, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCatalogReponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCatalogReponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalog = Catalog.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateCatalogReponse {
    return { catalog: isSet(object.catalog) ? Catalog.fromJSON(object.catalog) : undefined };
  },

  toJSON(message: UpdateCatalogReponse): unknown {
    const obj: any = {};
    if (message.catalog !== undefined) {
      obj.catalog = Catalog.toJSON(message.catalog);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateCatalogReponse>, I>>(base?: I): UpdateCatalogReponse {
    return UpdateCatalogReponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateCatalogReponse>, I>>(object: I): UpdateCatalogReponse {
    const message = createBaseUpdateCatalogReponse();
    message.catalog = (object.catalog !== undefined && object.catalog !== null)
      ? Catalog.fromPartial(object.catalog)
      : undefined;
    return message;
  },
};

function createBaseGetCatalogRequest(): GetCatalogRequest {
  return { namespace: "", name: "", useCache: false };
}

export const GetCatalogRequest = {
  encode(message: GetCatalogRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.useCache === true) {
      writer.uint32(24).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetCatalogRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogRequest>, I>>(base?: I): GetCatalogRequest {
    return GetCatalogRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogRequest>, I>>(object: I): GetCatalogRequest {
    const message = createBaseGetCatalogRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetCatalogResponse(): GetCatalogResponse {
  return { catalog: undefined };
}

export const GetCatalogResponse = {
  encode(message: GetCatalogResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.catalog !== undefined) {
      Catalog.encode(message.catalog, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalog = Catalog.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogResponse {
    return { catalog: isSet(object.catalog) ? Catalog.fromJSON(object.catalog) : undefined };
  },

  toJSON(message: GetCatalogResponse): unknown {
    const obj: any = {};
    if (message.catalog !== undefined) {
      obj.catalog = Catalog.toJSON(message.catalog);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogResponse>, I>>(base?: I): GetCatalogResponse {
    return GetCatalogResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogResponse>, I>>(object: I): GetCatalogResponse {
    const message = createBaseGetCatalogResponse();
    message.catalog = (object.catalog !== undefined && object.catalog !== null)
      ? Catalog.fromPartial(object.catalog)
      : undefined;
    return message;
  },
};

function createBaseGetCatalogIfChangedRequest(): GetCatalogIfChangedRequest {
  return { namespace: "", name: "", version: 0, useCache: false };
}

export const GetCatalogIfChangedRequest = {
  encode(message: GetCatalogIfChangedRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.version !== 0) {
      writer.uint32(24).int64(message.version);
    }
    if (message.useCache === true) {
      writer.uint32(32).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogIfChangedRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogIfChangedRequest();
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

          message.name = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.version = longToNumber(reader.int64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogIfChangedRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      name: isSet(object.name) ? String(object.name) : "",
      version: isSet(object.version) ? Number(object.version) : 0,
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetCatalogIfChangedRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.version !== 0) {
      obj.version = Math.round(message.version);
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogIfChangedRequest>, I>>(base?: I): GetCatalogIfChangedRequest {
    return GetCatalogIfChangedRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogIfChangedRequest>, I>>(object: I): GetCatalogIfChangedRequest {
    const message = createBaseGetCatalogIfChangedRequest();
    message.namespace = object.namespace ?? "";
    message.name = object.name ?? "";
    message.version = object.version ?? 0;
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetCatalogIfChangedResponse(): GetCatalogIfChangedResponse {
  return { null: undefined, catalog: undefined };
}

export const GetCatalogIfChangedResponse = {
  encode(message: GetCatalogIfChangedResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.null !== undefined) {
      writer.uint32(8).int32(message.null);
    }
    if (message.catalog !== undefined) {
      Catalog.encode(message.catalog, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogIfChangedResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogIfChangedResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.null = reader.int32() as any;
          continue;
        case 2:
          if (tag !== 18) {
            break;
          }

          message.catalog = Catalog.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogIfChangedResponse {
    return {
      null: isSet(object.null) ? nullValueFromJSON(object.null) : undefined,
      catalog: isSet(object.catalog) ? Catalog.fromJSON(object.catalog) : undefined,
    };
  },

  toJSON(message: GetCatalogIfChangedResponse): unknown {
    const obj: any = {};
    if (message.null !== undefined) {
      obj.null = nullValueToJSON(message.null);
    }
    if (message.catalog !== undefined) {
      obj.catalog = Catalog.toJSON(message.catalog);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogIfChangedResponse>, I>>(base?: I): GetCatalogIfChangedResponse {
    return GetCatalogIfChangedResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogIfChangedResponse>, I>>(object: I): GetCatalogIfChangedResponse {
    const message = createBaseGetCatalogIfChangedResponse();
    message.null = object.null ?? undefined;
    message.catalog = (object.catalog !== undefined && object.catalog !== null)
      ? Catalog.fromPartial(object.catalog)
      : undefined;
    return message;
  },
};

function createBaseGetAllCatalogsRequest(): GetAllCatalogsRequest {
  return { namespace: "", useCache: false };
}

export const GetAllCatalogsRequest = {
  encode(message: GetAllCatalogsRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.useCache === true) {
      writer.uint32(32).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllCatalogsRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllCatalogsRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.namespace = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAllCatalogsRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: GetAllCatalogsRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllCatalogsRequest>, I>>(base?: I): GetAllCatalogsRequest {
    return GetAllCatalogsRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllCatalogsRequest>, I>>(object: I): GetAllCatalogsRequest {
    const message = createBaseGetAllCatalogsRequest();
    message.namespace = object.namespace ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseGetAllCatalogsResponse(): GetAllCatalogsResponse {
  return { catalogs: [] };
}

export const GetAllCatalogsResponse = {
  encode(message: GetAllCatalogsResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.catalogs) {
      Catalog.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetAllCatalogsResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetAllCatalogsResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.catalogs.push(Catalog.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetAllCatalogsResponse {
    return { catalogs: Array.isArray(object?.catalogs) ? object.catalogs.map((e: any) => Catalog.fromJSON(e)) : [] };
  },

  toJSON(message: GetAllCatalogsResponse): unknown {
    const obj: any = {};
    if (message.catalogs?.length) {
      obj.catalogs = message.catalogs.map((e) => Catalog.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetAllCatalogsResponse>, I>>(base?: I): GetAllCatalogsResponse {
    return GetAllCatalogsResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetAllCatalogsResponse>, I>>(object: I): GetAllCatalogsResponse {
    const message = createBaseGetAllCatalogsResponse();
    message.catalogs = object.catalogs?.map((e) => Catalog.fromPartial(e)) || [];
    return message;
  },
};

function createBaseCatalogEntry(): CatalogEntry {
  return {
    namespace: "",
    catalog: "",
    uuid: "",
    data: new Uint8Array(0),
    created: undefined,
    updated: undefined,
    version: 0,
  };
}

export const CatalogEntry = {
  encode(message: CatalogEntry, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.uuid !== "") {
      writer.uint32(26).string(message.uuid);
    }
    if (message.data.length !== 0) {
      writer.uint32(34).bytes(message.data);
    }
    if (message.created !== undefined) {
      Timestamp.encode(toTimestamp(message.created), writer.uint32(802).fork()).ldelim();
    }
    if (message.updated !== undefined) {
      Timestamp.encode(toTimestamp(message.updated), writer.uint32(810).fork()).ldelim();
    }
    if (message.version !== 0) {
      writer.uint32(816).int64(message.version);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CatalogEntry {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCatalogEntry();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.uuid = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.data = reader.bytes();
          continue;
        case 100:
          if (tag !== 802) {
            break;
          }

          message.created = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 101:
          if (tag !== 810) {
            break;
          }

          message.updated = fromTimestamp(Timestamp.decode(reader, reader.uint32()));
          continue;
        case 102:
          if (tag !== 816) {
            break;
          }

          message.version = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CatalogEntry {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      uuid: isSet(object.uuid) ? String(object.uuid) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
      created: isSet(object.created) ? fromJsonTimestamp(object.created) : undefined,
      updated: isSet(object.updated) ? fromJsonTimestamp(object.updated) : undefined,
      version: isSet(object.version) ? Number(object.version) : 0,
    };
  },

  toJSON(message: CatalogEntry): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.uuid !== "") {
      obj.uuid = message.uuid;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
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

  create<I extends Exact<DeepPartial<CatalogEntry>, I>>(base?: I): CatalogEntry {
    return CatalogEntry.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CatalogEntry>, I>>(object: I): CatalogEntry {
    const message = createBaseCatalogEntry();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.uuid = object.uuid ?? "";
    message.data = object.data ?? new Uint8Array(0);
    message.created = object.created ?? undefined;
    message.updated = object.updated ?? undefined;
    message.version = object.version ?? 0;
    return message;
  },
};

function createBaseCreateCatalogEntryRequest(): CreateCatalogEntryRequest {
  return { namespace: "", catalog: "", data: new Uint8Array(0) };
}

export const CreateCatalogEntryRequest = {
  encode(message: CreateCatalogEntryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.data.length !== 0) {
      writer.uint32(26).bytes(message.data);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCatalogEntryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCatalogEntryRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
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

  fromJSON(object: any): CreateCatalogEntryRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      data: isSet(object.data) ? bytesFromBase64(object.data) : new Uint8Array(0),
    };
  },

  toJSON(message: CreateCatalogEntryRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.data.length !== 0) {
      obj.data = base64FromBytes(message.data);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateCatalogEntryRequest>, I>>(base?: I): CreateCatalogEntryRequest {
    return CreateCatalogEntryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateCatalogEntryRequest>, I>>(object: I): CreateCatalogEntryRequest {
    const message = createBaseCreateCatalogEntryRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.data = object.data ?? new Uint8Array(0);
    return message;
  },
};

function createBaseCreateCatalogEntryResponse(): CreateCatalogEntryResponse {
  return { entry: undefined };
}

export const CreateCatalogEntryResponse = {
  encode(message: CreateCatalogEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      CatalogEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CreateCatalogEntryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCreateCatalogEntryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entry = CatalogEntry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CreateCatalogEntryResponse {
    return { entry: isSet(object.entry) ? CatalogEntry.fromJSON(object.entry) : undefined };
  },

  toJSON(message: CreateCatalogEntryResponse): unknown {
    const obj: any = {};
    if (message.entry !== undefined) {
      obj.entry = CatalogEntry.toJSON(message.entry);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CreateCatalogEntryResponse>, I>>(base?: I): CreateCatalogEntryResponse {
    return CreateCatalogEntryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CreateCatalogEntryResponse>, I>>(object: I): CreateCatalogEntryResponse {
    const message = createBaseCreateCatalogEntryResponse();
    message.entry = (object.entry !== undefined && object.entry !== null)
      ? CatalogEntry.fromPartial(object.entry)
      : undefined;
    return message;
  },
};

function createBaseDeleteCatalogEntryRequest(): DeleteCatalogEntryRequest {
  return { namespace: "", catalog: "", entry: "" };
}

export const DeleteCatalogEntryRequest = {
  encode(message: DeleteCatalogEntryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.entry !== "") {
      writer.uint32(26).string(message.entry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteCatalogEntryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteCatalogEntryRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.entry = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): DeleteCatalogEntryRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      entry: isSet(object.entry) ? String(object.entry) : "",
    };
  },

  toJSON(message: DeleteCatalogEntryRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.entry !== "") {
      obj.entry = message.entry;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteCatalogEntryRequest>, I>>(base?: I): DeleteCatalogEntryRequest {
    return DeleteCatalogEntryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteCatalogEntryRequest>, I>>(object: I): DeleteCatalogEntryRequest {
    const message = createBaseDeleteCatalogEntryRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.entry = object.entry ?? "";
    return message;
  },
};

function createBaseDeleteCatalogEntryResponse(): DeleteCatalogEntryResponse {
  return {};
}

export const DeleteCatalogEntryResponse = {
  encode(_: DeleteCatalogEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): DeleteCatalogEntryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseDeleteCatalogEntryResponse();
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

  fromJSON(_: any): DeleteCatalogEntryResponse {
    return {};
  },

  toJSON(_: DeleteCatalogEntryResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<DeleteCatalogEntryResponse>, I>>(base?: I): DeleteCatalogEntryResponse {
    return DeleteCatalogEntryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<DeleteCatalogEntryResponse>, I>>(_: I): DeleteCatalogEntryResponse {
    const message = createBaseDeleteCatalogEntryResponse();
    return message;
  },
};

function createBaseUpdateCatalogEntryRequest(): UpdateCatalogEntryRequest {
  return { entry: undefined };
}

export const UpdateCatalogEntryRequest = {
  encode(message: UpdateCatalogEntryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      CatalogEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCatalogEntryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCatalogEntryRequest();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entry = CatalogEntry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateCatalogEntryRequest {
    return { entry: isSet(object.entry) ? CatalogEntry.fromJSON(object.entry) : undefined };
  },

  toJSON(message: UpdateCatalogEntryRequest): unknown {
    const obj: any = {};
    if (message.entry !== undefined) {
      obj.entry = CatalogEntry.toJSON(message.entry);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateCatalogEntryRequest>, I>>(base?: I): UpdateCatalogEntryRequest {
    return UpdateCatalogEntryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateCatalogEntryRequest>, I>>(object: I): UpdateCatalogEntryRequest {
    const message = createBaseUpdateCatalogEntryRequest();
    message.entry = (object.entry !== undefined && object.entry !== null)
      ? CatalogEntry.fromPartial(object.entry)
      : undefined;
    return message;
  },
};

function createBaseUpdateCatalogEntryResponse(): UpdateCatalogEntryResponse {
  return { entry: undefined };
}

export const UpdateCatalogEntryResponse = {
  encode(message: UpdateCatalogEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      CatalogEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): UpdateCatalogEntryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseUpdateCatalogEntryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entry = CatalogEntry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): UpdateCatalogEntryResponse {
    return { entry: isSet(object.entry) ? CatalogEntry.fromJSON(object.entry) : undefined };
  },

  toJSON(message: UpdateCatalogEntryResponse): unknown {
    const obj: any = {};
    if (message.entry !== undefined) {
      obj.entry = CatalogEntry.toJSON(message.entry);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<UpdateCatalogEntryResponse>, I>>(base?: I): UpdateCatalogEntryResponse {
    return UpdateCatalogEntryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<UpdateCatalogEntryResponse>, I>>(object: I): UpdateCatalogEntryResponse {
    const message = createBaseUpdateCatalogEntryResponse();
    message.entry = (object.entry !== undefined && object.entry !== null)
      ? CatalogEntry.fromPartial(object.entry)
      : undefined;
    return message;
  },
};

function createBaseGetCatalogEntryRequest(): GetCatalogEntryRequest {
  return { namespace: "", catalog: "", entry: "" };
}

export const GetCatalogEntryRequest = {
  encode(message: GetCatalogEntryRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.entry !== "") {
      writer.uint32(26).string(message.entry);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogEntryRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogEntryRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.entry = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogEntryRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      entry: isSet(object.entry) ? String(object.entry) : "",
    };
  },

  toJSON(message: GetCatalogEntryRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.entry !== "") {
      obj.entry = message.entry;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogEntryRequest>, I>>(base?: I): GetCatalogEntryRequest {
    return GetCatalogEntryRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogEntryRequest>, I>>(object: I): GetCatalogEntryRequest {
    const message = createBaseGetCatalogEntryRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.entry = object.entry ?? "";
    return message;
  },
};

function createBaseGetCatalogEntryResponse(): GetCatalogEntryResponse {
  return { entry: undefined };
}

export const GetCatalogEntryResponse = {
  encode(message: GetCatalogEntryResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      CatalogEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): GetCatalogEntryResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseGetCatalogEntryResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entry = CatalogEntry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): GetCatalogEntryResponse {
    return { entry: isSet(object.entry) ? CatalogEntry.fromJSON(object.entry) : undefined };
  },

  toJSON(message: GetCatalogEntryResponse): unknown {
    const obj: any = {};
    if (message.entry !== undefined) {
      obj.entry = CatalogEntry.toJSON(message.entry);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<GetCatalogEntryResponse>, I>>(base?: I): GetCatalogEntryResponse {
    return GetCatalogEntryResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<GetCatalogEntryResponse>, I>>(object: I): GetCatalogEntryResponse {
    const message = createBaseGetCatalogEntryResponse();
    message.entry = (object.entry !== undefined && object.entry !== null)
      ? CatalogEntry.fromPartial(object.entry)
      : undefined;
    return message;
  },
};

function createBaseListCatalogEntriesRequest(): ListCatalogEntriesRequest {
  return { namespace: "", catalog: "", skip: 0, limit: 0 };
}

export const ListCatalogEntriesRequest = {
  encode(message: ListCatalogEntriesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.skip !== 0) {
      writer.uint32(24).int64(message.skip);
    }
    if (message.limit !== 0) {
      writer.uint32(32).int64(message.limit);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListCatalogEntriesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListCatalogEntriesRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 24) {
            break;
          }

          message.skip = longToNumber(reader.int64() as Long);
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.limit = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListCatalogEntriesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      skip: isSet(object.skip) ? Number(object.skip) : 0,
      limit: isSet(object.limit) ? Number(object.limit) : 0,
    };
  },

  toJSON(message: ListCatalogEntriesRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.skip !== 0) {
      obj.skip = Math.round(message.skip);
    }
    if (message.limit !== 0) {
      obj.limit = Math.round(message.limit);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListCatalogEntriesRequest>, I>>(base?: I): ListCatalogEntriesRequest {
    return ListCatalogEntriesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListCatalogEntriesRequest>, I>>(object: I): ListCatalogEntriesRequest {
    const message = createBaseListCatalogEntriesRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.skip = object.skip ?? 0;
    message.limit = object.limit ?? 0;
    return message;
  },
};

function createBaseListCatalogEntriesResponse(): ListCatalogEntriesResponse {
  return { entry: undefined };
}

export const ListCatalogEntriesResponse = {
  encode(message: ListCatalogEntriesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.entry !== undefined) {
      CatalogEntry.encode(message.entry, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListCatalogEntriesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListCatalogEntriesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.entry = CatalogEntry.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListCatalogEntriesResponse {
    return { entry: isSet(object.entry) ? CatalogEntry.fromJSON(object.entry) : undefined };
  },

  toJSON(message: ListCatalogEntriesResponse): unknown {
    const obj: any = {};
    if (message.entry !== undefined) {
      obj.entry = CatalogEntry.toJSON(message.entry);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListCatalogEntriesResponse>, I>>(base?: I): ListCatalogEntriesResponse {
    return ListCatalogEntriesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListCatalogEntriesResponse>, I>>(object: I): ListCatalogEntriesResponse {
    const message = createBaseListCatalogEntriesResponse();
    message.entry = (object.entry !== undefined && object.entry !== null)
      ? CatalogEntry.fromPartial(object.entry)
      : undefined;
    return message;
  },
};

function createBaseCountCatalogEntriesRequest(): CountCatalogEntriesRequest {
  return { namespace: "", catalog: "" };
}

export const CountCatalogEntriesRequest = {
  encode(message: CountCatalogEntriesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountCatalogEntriesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountCatalogEntriesRequest();
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

          message.catalog = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountCatalogEntriesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
    };
  },

  toJSON(message: CountCatalogEntriesRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountCatalogEntriesRequest>, I>>(base?: I): CountCatalogEntriesRequest {
    return CountCatalogEntriesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountCatalogEntriesRequest>, I>>(object: I): CountCatalogEntriesRequest {
    const message = createBaseCountCatalogEntriesRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    return message;
  },
};

function createBaseCountCatalogEntriesResponse(): CountCatalogEntriesResponse {
  return { count: 0 };
}

export const CountCatalogEntriesResponse = {
  encode(message: CountCatalogEntriesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.count !== 0) {
      writer.uint32(8).int64(message.count);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CountCatalogEntriesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCountCatalogEntriesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 8) {
            break;
          }

          message.count = longToNumber(reader.int64() as Long);
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CountCatalogEntriesResponse {
    return { count: isSet(object.count) ? Number(object.count) : 0 };
  },

  toJSON(message: CountCatalogEntriesResponse): unknown {
    const obj: any = {};
    if (message.count !== 0) {
      obj.count = Math.round(message.count);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CountCatalogEntriesResponse>, I>>(base?: I): CountCatalogEntriesResponse {
    return CountCatalogEntriesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CountCatalogEntriesResponse>, I>>(object: I): CountCatalogEntriesResponse {
    const message = createBaseCountCatalogEntriesResponse();
    message.count = object.count ?? 0;
    return message;
  },
};

function createBaseCatalogIndex(): CatalogIndex {
  return { namespace: "", catalog: "", name: "", fields: [], unique: false };
}

export const CatalogIndex = {
  encode(message: CatalogIndex, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.name !== "") {
      writer.uint32(26).string(message.name);
    }
    for (const v of message.fields) {
      CatalogIndex_IndexField.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    if (message.unique === true) {
      writer.uint32(40).bool(message.unique);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CatalogIndex {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCatalogIndex();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.name = reader.string();
          continue;
        case 4:
          if (tag !== 34) {
            break;
          }

          message.fields.push(CatalogIndex_IndexField.decode(reader, reader.uint32()));
          continue;
        case 5:
          if (tag !== 40) {
            break;
          }

          message.unique = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CatalogIndex {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      name: isSet(object.name) ? String(object.name) : "",
      fields: Array.isArray(object?.fields) ? object.fields.map((e: any) => CatalogIndex_IndexField.fromJSON(e)) : [],
      unique: isSet(object.unique) ? Boolean(object.unique) : false,
    };
  },

  toJSON(message: CatalogIndex): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.fields?.length) {
      obj.fields = message.fields.map((e) => CatalogIndex_IndexField.toJSON(e));
    }
    if (message.unique === true) {
      obj.unique = message.unique;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CatalogIndex>, I>>(base?: I): CatalogIndex {
    return CatalogIndex.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CatalogIndex>, I>>(object: I): CatalogIndex {
    const message = createBaseCatalogIndex();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.name = object.name ?? "";
    message.fields = object.fields?.map((e) => CatalogIndex_IndexField.fromPartial(e)) || [];
    message.unique = object.unique ?? false;
    return message;
  },
};

function createBaseCatalogIndex_IndexField(): CatalogIndex_IndexField {
  return { name: "", type: 0 };
}

export const CatalogIndex_IndexField = {
  encode(message: CatalogIndex_IndexField, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.name !== "") {
      writer.uint32(10).string(message.name);
    }
    if (message.type !== 0) {
      writer.uint32(16).int32(message.type);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): CatalogIndex_IndexField {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseCatalogIndex_IndexField();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.name = reader.string();
          continue;
        case 2:
          if (tag !== 16) {
            break;
          }

          message.type = reader.int32() as any;
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): CatalogIndex_IndexField {
    return {
      name: isSet(object.name) ? String(object.name) : "",
      type: isSet(object.type) ? catalogIndex_IndexField_IndexTypeFromJSON(object.type) : 0,
    };
  },

  toJSON(message: CatalogIndex_IndexField): unknown {
    const obj: any = {};
    if (message.name !== "") {
      obj.name = message.name;
    }
    if (message.type !== 0) {
      obj.type = catalogIndex_IndexField_IndexTypeToJSON(message.type);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<CatalogIndex_IndexField>, I>>(base?: I): CatalogIndex_IndexField {
    return CatalogIndex_IndexField.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<CatalogIndex_IndexField>, I>>(object: I): CatalogIndex_IndexField {
    const message = createBaseCatalogIndex_IndexField();
    message.name = object.name ?? "";
    message.type = object.type ?? 0;
    return message;
  },
};

function createBaseListCatalogIndexesRequest(): ListCatalogIndexesRequest {
  return { namespace: "", catalog: "", useCache: false };
}

export const ListCatalogIndexesRequest = {
  encode(message: ListCatalogIndexesRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.useCache === true) {
      writer.uint32(32).bool(message.useCache);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListCatalogIndexesRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListCatalogIndexesRequest();
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

          message.catalog = reader.string();
          continue;
        case 4:
          if (tag !== 32) {
            break;
          }

          message.useCache = reader.bool();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListCatalogIndexesRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      useCache: isSet(object.useCache) ? Boolean(object.useCache) : false,
    };
  },

  toJSON(message: ListCatalogIndexesRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.useCache === true) {
      obj.useCache = message.useCache;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListCatalogIndexesRequest>, I>>(base?: I): ListCatalogIndexesRequest {
    return ListCatalogIndexesRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListCatalogIndexesRequest>, I>>(object: I): ListCatalogIndexesRequest {
    const message = createBaseListCatalogIndexesRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.useCache = object.useCache ?? false;
    return message;
  },
};

function createBaseListCatalogIndexesResponse(): ListCatalogIndexesResponse {
  return { indexes: [] };
}

export const ListCatalogIndexesResponse = {
  encode(message: ListCatalogIndexesResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    for (const v of message.indexes) {
      CatalogIndex.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): ListCatalogIndexesResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseListCatalogIndexesResponse();
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          if (tag !== 10) {
            break;
          }

          message.indexes.push(CatalogIndex.decode(reader, reader.uint32()));
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): ListCatalogIndexesResponse {
    return { indexes: Array.isArray(object?.indexes) ? object.indexes.map((e: any) => CatalogIndex.fromJSON(e)) : [] };
  },

  toJSON(message: ListCatalogIndexesResponse): unknown {
    const obj: any = {};
    if (message.indexes?.length) {
      obj.indexes = message.indexes.map((e) => CatalogIndex.toJSON(e));
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<ListCatalogIndexesResponse>, I>>(base?: I): ListCatalogIndexesResponse {
    return ListCatalogIndexesResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<ListCatalogIndexesResponse>, I>>(object: I): ListCatalogIndexesResponse {
    const message = createBaseListCatalogIndexesResponse();
    message.indexes = object.indexes?.map((e) => CatalogIndex.fromPartial(e)) || [];
    return message;
  },
};

function createBaseEnsureCatalogIndexRequest(): EnsureCatalogIndexRequest {
  return { namespace: "", catalog: "", index: undefined };
}

export const EnsureCatalogIndexRequest = {
  encode(message: EnsureCatalogIndexRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.index !== undefined) {
      CatalogIndex.encode(message.index, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureCatalogIndexRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureCatalogIndexRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.index = CatalogIndex.decode(reader, reader.uint32());
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): EnsureCatalogIndexRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      index: isSet(object.index) ? CatalogIndex.fromJSON(object.index) : undefined,
    };
  },

  toJSON(message: EnsureCatalogIndexRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.index !== undefined) {
      obj.index = CatalogIndex.toJSON(message.index);
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureCatalogIndexRequest>, I>>(base?: I): EnsureCatalogIndexRequest {
    return EnsureCatalogIndexRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureCatalogIndexRequest>, I>>(object: I): EnsureCatalogIndexRequest {
    const message = createBaseEnsureCatalogIndexRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.index = (object.index !== undefined && object.index !== null)
      ? CatalogIndex.fromPartial(object.index)
      : undefined;
    return message;
  },
};

function createBaseEnsureCatalogIndexResponse(): EnsureCatalogIndexResponse {
  return {};
}

export const EnsureCatalogIndexResponse = {
  encode(_: EnsureCatalogIndexResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): EnsureCatalogIndexResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseEnsureCatalogIndexResponse();
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

  fromJSON(_: any): EnsureCatalogIndexResponse {
    return {};
  },

  toJSON(_: EnsureCatalogIndexResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<EnsureCatalogIndexResponse>, I>>(base?: I): EnsureCatalogIndexResponse {
    return EnsureCatalogIndexResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<EnsureCatalogIndexResponse>, I>>(_: I): EnsureCatalogIndexResponse {
    const message = createBaseEnsureCatalogIndexResponse();
    return message;
  },
};

function createBaseRemoveCatalogIndexRequest(): RemoveCatalogIndexRequest {
  return { namespace: "", catalog: "", index: "" };
}

export const RemoveCatalogIndexRequest = {
  encode(message: RemoveCatalogIndexRequest, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    if (message.namespace !== "") {
      writer.uint32(10).string(message.namespace);
    }
    if (message.catalog !== "") {
      writer.uint32(18).string(message.catalog);
    }
    if (message.index !== "") {
      writer.uint32(26).string(message.index);
    }
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveCatalogIndexRequest {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveCatalogIndexRequest();
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

          message.catalog = reader.string();
          continue;
        case 3:
          if (tag !== 26) {
            break;
          }

          message.index = reader.string();
          continue;
      }
      if ((tag & 7) === 4 || tag === 0) {
        break;
      }
      reader.skipType(tag & 7);
    }
    return message;
  },

  fromJSON(object: any): RemoveCatalogIndexRequest {
    return {
      namespace: isSet(object.namespace) ? String(object.namespace) : "",
      catalog: isSet(object.catalog) ? String(object.catalog) : "",
      index: isSet(object.index) ? String(object.index) : "",
    };
  },

  toJSON(message: RemoveCatalogIndexRequest): unknown {
    const obj: any = {};
    if (message.namespace !== "") {
      obj.namespace = message.namespace;
    }
    if (message.catalog !== "") {
      obj.catalog = message.catalog;
    }
    if (message.index !== "") {
      obj.index = message.index;
    }
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveCatalogIndexRequest>, I>>(base?: I): RemoveCatalogIndexRequest {
    return RemoveCatalogIndexRequest.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveCatalogIndexRequest>, I>>(object: I): RemoveCatalogIndexRequest {
    const message = createBaseRemoveCatalogIndexRequest();
    message.namespace = object.namespace ?? "";
    message.catalog = object.catalog ?? "";
    message.index = object.index ?? "";
    return message;
  },
};

function createBaseRemoveCatalogIndexResponse(): RemoveCatalogIndexResponse {
  return {};
}

export const RemoveCatalogIndexResponse = {
  encode(_: RemoveCatalogIndexResponse, writer: _m0.Writer = _m0.Writer.create()): _m0.Writer {
    return writer;
  },

  decode(input: _m0.Reader | Uint8Array, length?: number): RemoveCatalogIndexResponse {
    const reader = input instanceof _m0.Reader ? input : _m0.Reader.create(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = createBaseRemoveCatalogIndexResponse();
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

  fromJSON(_: any): RemoveCatalogIndexResponse {
    return {};
  },

  toJSON(_: RemoveCatalogIndexResponse): unknown {
    const obj: any = {};
    return obj;
  },

  create<I extends Exact<DeepPartial<RemoveCatalogIndexResponse>, I>>(base?: I): RemoveCatalogIndexResponse {
    return RemoveCatalogIndexResponse.fromPartial(base ?? ({} as any));
  },
  fromPartial<I extends Exact<DeepPartial<RemoveCatalogIndexResponse>, I>>(_: I): RemoveCatalogIndexResponse {
    const message = createBaseRemoveCatalogIndexResponse();
    return message;
  },
};

export interface CatalogService {
  /** Create new catalog */
  Create(request: CreateCatalogRequest): Promise<CreateCatalogResponse>;
  /** Deletes catalog and all its entries */
  Delete(request: DeleteCatalogRequest): Promise<DeleteCatalogResponse>;
  /** Updates catalog with new information */
  Update(request: UpdateCatalogRequest): Promise<UpdateCatalogReponse>;
  /** Returns catalog by its name */
  Get(request: GetCatalogRequest): Promise<GetCatalogResponse>;
  /** Returns catalog by its name only if provided version differs from the actual. In other case returns NULL. More optimized version, than Get */
  GetIfChanged(request: GetCatalogIfChangedRequest): Promise<GetCatalogIfChangedResponse>;
  /** Streams list of all catalogs */
  GetAll(request: GetAllCatalogsRequest): Promise<GetAllCatalogsResponse>;
}

export const CatalogServiceServiceName = "erp_catalog.CatalogService";
export class CatalogServiceClientImpl implements CatalogService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || CatalogServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Update = this.Update.bind(this);
    this.Get = this.Get.bind(this);
    this.GetIfChanged = this.GetIfChanged.bind(this);
    this.GetAll = this.GetAll.bind(this);
  }
  Create(request: CreateCatalogRequest): Promise<CreateCatalogResponse> {
    const data = CreateCatalogRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateCatalogResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteCatalogRequest): Promise<DeleteCatalogResponse> {
    const data = DeleteCatalogRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteCatalogResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateCatalogRequest): Promise<UpdateCatalogReponse> {
    const data = UpdateCatalogRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateCatalogReponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetCatalogRequest): Promise<GetCatalogResponse> {
    const data = GetCatalogRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetCatalogResponse.decode(_m0.Reader.create(data)));
  }

  GetIfChanged(request: GetCatalogIfChangedRequest): Promise<GetCatalogIfChangedResponse> {
    const data = GetCatalogIfChangedRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetIfChanged", data);
    return promise.then((data) => GetCatalogIfChangedResponse.decode(_m0.Reader.create(data)));
  }

  GetAll(request: GetAllCatalogsRequest): Promise<GetAllCatalogsResponse> {
    const data = GetAllCatalogsRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "GetAll", data);
    return promise.then((data) => GetAllCatalogsResponse.decode(_m0.Reader.create(data)));
  }
}

export interface CatalogEntryService {
  /** Creates new entry in the specified catalog. Entry will receive uuid after successfull creation. Returns newly created entry */
  Create(request: CreateCatalogEntryRequest): Promise<CreateCatalogEntryResponse>;
  /** Deletes catalog entry */
  Delete(request: DeleteCatalogEntryRequest): Promise<DeleteCatalogEntryResponse>;
  /** Updates catalog entry with new data */
  Update(request: UpdateCatalogEntryRequest): Promise<UpdateCatalogEntryResponse>;
  /** Get catalog entry. Uses cache and works much faster than Query operation */
  Get(request: GetCatalogEntryRequest): Promise<GetCatalogEntryResponse>;
  /** List catalog entries. */
  List(request: ListCatalogEntriesRequest): Observable<ListCatalogEntriesResponse>;
  /** Count catalog entries. */
  Count(request: CountCatalogEntriesRequest): Promise<CountCatalogEntriesResponse>;
}

export const CatalogEntryServiceServiceName = "erp_catalog.CatalogEntryService";
export class CatalogEntryServiceClientImpl implements CatalogEntryService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || CatalogEntryServiceServiceName;
    this.rpc = rpc;
    this.Create = this.Create.bind(this);
    this.Delete = this.Delete.bind(this);
    this.Update = this.Update.bind(this);
    this.Get = this.Get.bind(this);
    this.List = this.List.bind(this);
    this.Count = this.Count.bind(this);
  }
  Create(request: CreateCatalogEntryRequest): Promise<CreateCatalogEntryResponse> {
    const data = CreateCatalogEntryRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Create", data);
    return promise.then((data) => CreateCatalogEntryResponse.decode(_m0.Reader.create(data)));
  }

  Delete(request: DeleteCatalogEntryRequest): Promise<DeleteCatalogEntryResponse> {
    const data = DeleteCatalogEntryRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Delete", data);
    return promise.then((data) => DeleteCatalogEntryResponse.decode(_m0.Reader.create(data)));
  }

  Update(request: UpdateCatalogEntryRequest): Promise<UpdateCatalogEntryResponse> {
    const data = UpdateCatalogEntryRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Update", data);
    return promise.then((data) => UpdateCatalogEntryResponse.decode(_m0.Reader.create(data)));
  }

  Get(request: GetCatalogEntryRequest): Promise<GetCatalogEntryResponse> {
    const data = GetCatalogEntryRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Get", data);
    return promise.then((data) => GetCatalogEntryResponse.decode(_m0.Reader.create(data)));
  }

  List(request: ListCatalogEntriesRequest): Observable<ListCatalogEntriesResponse> {
    const data = ListCatalogEntriesRequest.encode(request).finish();
    const result = this.rpc.serverStreamingRequest(this.service, "List", data);
    return result.pipe(map((data) => ListCatalogEntriesResponse.decode(_m0.Reader.create(data))));
  }

  Count(request: CountCatalogEntriesRequest): Promise<CountCatalogEntriesResponse> {
    const data = CountCatalogEntriesRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "Count", data);
    return promise.then((data) => CountCatalogEntriesResponse.decode(_m0.Reader.create(data)));
  }
}

export interface CatalogIndexService {
  /** Lists all indexes in the catalog */
  ListIndexes(request: ListCatalogIndexesRequest): Promise<ListCatalogIndexesResponse>;
  /** Creates or updates index in the catalog */
  EnsureIndex(request: EnsureCatalogIndexRequest): Promise<EnsureCatalogIndexResponse>;
  /** Removes index from the catalog */
  RemoveIndex(request: RemoveCatalogIndexRequest): Promise<RemoveCatalogIndexResponse>;
}

export const CatalogIndexServiceServiceName = "erp_catalog.CatalogIndexService";
export class CatalogIndexServiceClientImpl implements CatalogIndexService {
  private readonly rpc: Rpc;
  private readonly service: string;
  constructor(rpc: Rpc, opts?: { service?: string }) {
    this.service = opts?.service || CatalogIndexServiceServiceName;
    this.rpc = rpc;
    this.ListIndexes = this.ListIndexes.bind(this);
    this.EnsureIndex = this.EnsureIndex.bind(this);
    this.RemoveIndex = this.RemoveIndex.bind(this);
  }
  ListIndexes(request: ListCatalogIndexesRequest): Promise<ListCatalogIndexesResponse> {
    const data = ListCatalogIndexesRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "ListIndexes", data);
    return promise.then((data) => ListCatalogIndexesResponse.decode(_m0.Reader.create(data)));
  }

  EnsureIndex(request: EnsureCatalogIndexRequest): Promise<EnsureCatalogIndexResponse> {
    const data = EnsureCatalogIndexRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "EnsureIndex", data);
    return promise.then((data) => EnsureCatalogIndexResponse.decode(_m0.Reader.create(data)));
  }

  RemoveIndex(request: RemoveCatalogIndexRequest): Promise<RemoveCatalogIndexResponse> {
    const data = RemoveCatalogIndexRequest.encode(request).finish();
    const promise = this.rpc.request(this.service, "RemoveIndex", data);
    return promise.then((data) => RemoveCatalogIndexResponse.decode(_m0.Reader.create(data)));
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
