function socketInit($) {
    //建立websocket连接
    socket = io.connect('/', {
        // transports: ['polling']
    });

    socket.emit('msg', JSON.stringify({ content: ">>>>>>>>>>>" }), function (data) {
        // 发送消息后直接返回的消息内容
        console.log(">>>>>", data);
    })
    // 接收消息
    socket.on('msg', function (data) {
        console.log(">>>>>>>>>>接收", JSON.stringify(data));
        var htmldata = `
        <li class="layui-timeline-item">
            <i class="layui-icon layui-timeline-axis">&#xe63f;</i>
            <div class="layui-timeline-content layui-text">
                <h3 class="layui-timeline-title">${data.Time}</h3>
                <p>
                    <i class="layui-icon"></i>${data.Msg}
                </p>
            </div>
        </li>
        `
        $("#waterfall-step ul").prepend(htmldata);
    });

}