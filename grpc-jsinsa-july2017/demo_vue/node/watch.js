var grpc = require('grpc');
var fs = require('fs');
var path = require('path');

var service = grpc.load('../proto/carservice.proto');

//var keyPath = path.join('..','..', 'misc', 'localhost.key');
//var certPath = path.join('..','..', 'misc', 'localhost.crt');
//var ssl_creds = grpc.credentials.createSsl(fs.readFileSync(keyPath),
 //                                           fs.readFileSync(certPath));

var client = new service.proto.CarServiceDepartment('0.0.0.0:9090', 
                    grpc.credentials.createInsecure());
                    //ssl_creds);

console.log("Waiting for new bookings...");
var call = client.watch({});
    call.on('data', function(booking) {
        console.log('Car Service booking made! ');
        // do lookup to see who the sales person was
        // send notification
    });
/*
// In-memory array of history of sales objects
var salesHistory = [{
    car: {
        regNo : "ABC123GP",
        make : "SmallCar",
        model : "i99",
        year : 2014,

    },
    salesPerson : "Mrs. Jones"
},
{
    car: {
        regNo : "ZZZ777GP",
        make : "BigCar",
        model : "320i",
        year : 2017,

    },
    salesPerson : "Mr. Smith"
}];*/