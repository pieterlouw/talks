// package: proto
// file: carservice.proto

import * as carservice_pb from "./carservice_pb";
export class CarServiceDepartment {
  static serviceName = "proto.CarServiceDepartment";
}
export namespace CarServiceDepartment {
  export class MakeBooking {
    static methodName = "MakeBooking";
    static service = CarServiceDepartment;
    static requestStream = false;
    static responseStream = false;
    static requestType = carservice_pb.Booking;
    static responseType = carservice_pb.Empty;
  }
  export class Watch {
    static methodName = "Watch";
    static service = CarServiceDepartment;
    static requestStream = false;
    static responseStream = true;
    static requestType = carservice_pb.Empty;
    static responseType = carservice_pb.Booking;
  }
}
