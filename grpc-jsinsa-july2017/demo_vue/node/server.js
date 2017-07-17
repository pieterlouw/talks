
var events = require('events');
var bookStream = new events.EventEmitter();

var grpc = require('grpc');
var service = grpc.load('../proto/carservice.proto');
var server = new grpc.Server();

server.addService(service.proto.CarServiceDepartment.service, {
    makeBooking: function(call, callback) {
        var booking = call.request;
        console.log("New Booking received");
        bookStream.emit('new_booking', booking);
        callback(null, {});
    }, 
    watch: function(stream) {
        bookStream.on('new_booking', function(booking){
            stream.write(booking);
        });
    }
});

console.log("gRPC server started on 0.0.0.0:9090")
server.bind('0.0.0.0:9090', grpc.ServerCredentials.createInsecure());
server.start();
