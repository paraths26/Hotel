$(document).ready(function(){
	
    console.log(location.host)

    var orderSocket = new WebSocket("ws://"+location.host+"/order/alexa")
    orderSocket.onopen = function(event){
        orderSocket.send("I got connected")
    }
    orderSocket.onclose = function(event){
        console.log("closing connection")
    }
    orderSocket.onerror = function(event){
        console.log("error in connection")
    }
    
    orderSocket.onmessage = function(event){
        console.log("message received ", event.data)
        orderData = JSON.parse(event.data)
      
        showOrder(orderData)
    }
    
    function showOrder(data){
        // document.getElementById('')
        console.log("Rceived oder data: ",data)

        var rows = "<tr><td>" + data.room + "</td><td>" + data.command + "</td><td>" + data.data + "</td></tr>";
       
        $( rows ).appendTo( "#orderTable tbody" );
    }

   

});