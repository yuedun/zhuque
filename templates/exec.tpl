<!DOCTYPE html>
<html lang="zh-cn-Hans">

<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <meta http-equiv="X-UA-Compatible" content="ie=edge">
    <title>{{.title}}</title>
    <!-- 生产环境版本，优化了尺寸和速度 -->
    <script src="https://cdn.jsdelivr.net/npm/vue@2.6.11"></script>
</head>

<body>
    <div id="app">
        <input type="text" name="cmd" v-model="cmd">
        <button type="submit" @click="send">提交</button>
    </div>
    <script>
        var app = new Vue({
            el: '#app',
            data: {
                cmd: 'Hello Vue!'
            },
            methods: {
                send() {
                    var formData = new FormData()
                    formData.append("cmd", this.cmd)
                    return fetch("/exec/send", {
                        method: "POST",
                        body: formData
                    })
                }
            }
        })
    </script>
</body>

</html>