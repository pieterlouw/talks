// package: proto
// file: carservice.proto

import * as jspb from "google-protobuf";

export class Booking extends jspb.Message {
  getReg(): string;
  setReg(value: string): void;

  getOdo(): number;
  setOdo(value: number): void;

  getName(): string;
  setName(value: string): void;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Booking.AsObject;
  static toObject(includeInstance: boolean, msg: Booking): Booking.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Booking, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Booking;
  static deserializeBinaryFromReader(message: Booking, reader: jspb.BinaryReader): Booking;
}

export namespace Booking {
  export type AsObject = {
    reg: string,
    odo: number,
    name: string,
  }
}

export class Empty extends jspb.Message {
  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): Empty.AsObject;
  static toObject(includeInstance: boolean, msg: Empty): Empty.AsObject;
  static extensions: {[key: number]: jspb.ExtensionFieldInfo<jspb.Message>};
  static extensionsBinary: {[key: number]: jspb.ExtensionFieldBinaryInfo<jspb.Message>};
  static serializeBinaryToWriter(message: Empty, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): Empty;
  static deserializeBinaryFromReader(message: Empty, reader: jspb.BinaryReader): Empty;
}

export namespace Empty {
  export type AsObject = {
  }
}

