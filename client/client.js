var PROTO_PATH = __dirname + '/../api/v1/capturedimage.proto';
var async = require('async');
var grpc = require('@grpc/grpc-js');
var protoLoader = require('@grpc/proto-loader');
var packageDefinition = protoLoader.loadSync(
    PROTO_PATH,
    {keepCase: true,
     longs: String,
     enums: String,
     defaults: true,
     oneofs: true
    });
var messages = require('./capturedimage_pb');    
var img_api = grpc.loadPackageDefinition(packageDefinition).capturedimage_api;
var client = new img_api.MeteoBaltica('127.0.0.1:3144',grpc.credentials.createInsecure());      

function runRoutePath(callback) {
    var reqPath = {
        timeBegin : 0,
        timeEnd : Math.floor(Date.now() / 1000)
    };
    var call = client.getRoutePath(reqPath);
    call.on('data', function(error, route) {
        console.log('Receive route ' +
        route.timeBegin + ' ' + route.timeEnd);
    });
    call.on('end', callback);
    console.log('runRoutePath');
}

function main() {
    async.series([   
        runRoutePath,
    ]);
  }
  
  if (require.main === module) {
    main();
  }
module.exports = runRoutePath;
//call.on('end', callback);