const proxy = require('http-proxy-middleware');

const onPremisses = {
    target:'http://10.55.2.80:8000',
    apis:[
        '/api/v1/on-premises'
    ]
}

const dataSet = {//数据集
    target:'http://srv.cloud.jiangxingai.com:9002',
    apis:[
        '/api/v1/dataSet',
        '/data'
    ]
}

const label = {//标签
    target: 'http://srv.cloud.jiangxingai.com:9002',
    apis: [
        '/api/v1/labelSet'
    ]
}

const model = {//模型
    target: 'http://srv.cloud.jiangxingai.com:9004',
    apis: [
        '/api/v1/model'
    ]
}


const usingModel = {//模型
    target:'http://srv.cloud.jiangxingai.com:9004',
    apis:[
        '/api/v1/jxserving',
    ]
}

const analysis = {//阈值分析
    target:'http://srv.cloud.jiangxingai.com:9006',
    apis:[
        '/api/v1/analyst/query'
    ]
}

const interference = {//人工干预
    target:'http://srv.cloud.jiangxingai.com:10000',
    apis:[
        '/api/v1/aiManual'
    ]
}


const urlss = [
    dataSet,
    label,
    model,
    usingModel,
    analysis,
    interference,
    onPremisses
]
module.exports = function(app) {
    urlss.forEach(urls=>{
        urls.apis.forEach(path=>{
            app.use(path, proxy({target:urls.target}))
        })
    })
};