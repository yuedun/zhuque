//建立websocket连接
socket = io.connect('/',{
    // transports: ['polling']
});

socket.emit('msg', JSON.stringify({content:">>>>>>>>>>>"}), function(data){
    // 发送消息后直接返回的消息内容
    console.log(">>>>>",data);
})
// 接收消息
socket.on('msg', function (json) {
    console.log(">>>>>>>>>>接收", json);
    document.querySelector("#recive").textContent = json;
});