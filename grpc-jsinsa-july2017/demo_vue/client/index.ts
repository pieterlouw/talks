import {grpc, BrowserHeaders} from "grpc-web-client";
import {CarServiceDepartment} from "./_proto/carservice_pb_service";
import {Booking, Empty} from "./_proto/carservice_pb";

import { Component } from 'vue-typed'
import * as Vue from 'vue'

const template = require('./app.jade')();

@Component({
	template
})
class App extends Vue {
    host: string = 'https://localhost:8443'; // grpc endpoint

    regNo: string =  '';
    odoMeter: number = 0;
    customerName: string = '';

    mounted(){  
        console.log("Vue app started.. Will connect to:" + this.host);


        /*this.regNo = "PYB 268 GP";
        this.odoMeter = 1234;
        this.customerName = "Pieter";
        this.registerBooking();*/
    }

    registerBooking(){
        
        //validation
        if (!this.regNo || !this.customerName || this.odoMeter == 0){
            alert("Please populate all fields!");
            return;
        }

        
        const request = new Booking();
        request.setReg(this.regNo);
        request.setOdo(this.odoMeter);
        request.setName(this.customerName);

        grpc.invoke(CarServiceDepartment.MakeBooking, {
            host: this.host,
            request: request,

            onMessage: (empty: Empty) => {
                this.regNo = ''; 
                this.customerName = ''; 
                this.odoMeter = 0;

                alert("Booking successful");
            },

            onEnd(code: grpc.Code, message: string, trailers: BrowserHeaders){
                if (code != 0) {
                    alert("Booking error:" + code + " " + message);
                }

                console.log(code, message, trailers);
            }
        })

    }
}

new App().$mount('#app');