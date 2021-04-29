const urls = {
    //人工干预
    'verifyInfo':'/api/v1/aiManual/verifyInfo',
    'update':'/api/v1/aiManual/update',
    'create':'/api/v1/aiManual/create',

    //阈值分析
    'analysisQuery':'/api/v1/analyst/query',

    //模型
    'uploadModel':'/api/v1/model/upload',
    'uploadModelMsg':'/api/v1/model/getmsg',
    'deleteModel':'/api/v1/model/delete',
    'getModel':'/api/v1/model/getmodel',
    'publishModel':'/api/v1/model/publish',
    'adjust':'/api/v1/model/evaluate',
    'modelHistory':'/api/v1/model/historyevaluate',
    'getFrameWork':'/api/v1/model/getframework',
    'unpublishModal':'/api/v1/model/cancelpublish',

    //标签
    'uploadLabel':'/api/v1/labelSet/upload',
    'deleteLabel':'/api/v1/labelSet/delete',
    'updateLabel':'/api/v1/labelSet/update',
    'getLabel':'/api/v1/labelSet/query',

    //数据集
    'downloadDataSet':'/api/v1/dataSet/download',
    'uploadDataSet':'/api/v1/dataSet/upload',
    'getDataSet':'/api/v1/dataSet/query',
    'deleteDataSet':'api/v1/dataSet/delete',
    'datasetManageQuery':'/api/v1/dataSet/manage/query',
    'datasetManageUpdate':'/api/v1/dataSet/manage/update',

    //人工干预
    'modelSwitch':'/api/v1/aiManual/jxserving/switch',
    'modelDetect':'/api/v1/aiManual/jxserving/detect',

    //on-premises
    'getOnPremises':'/api/v1/on-premises/getjxsboard'
}
export default urls