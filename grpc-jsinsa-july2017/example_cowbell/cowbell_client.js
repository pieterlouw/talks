
var PROTO_PATH = __dirname + '/../service/service.proto';

var grpc = require('grpc');
var cowbell = grpc.load(PROTO_PATH).cowbell;

function main() {
  var client = new cowbell.CowbellService('localhost:9090',
                                       grpc.credentials.createInsecure());

  client.moreCowbell({qty: 1}, function(err, response) {
    if (err)
        console.log('Error: ', err);
    else
        console.log('Total Cowbells:', response.total);
  });
}

main();