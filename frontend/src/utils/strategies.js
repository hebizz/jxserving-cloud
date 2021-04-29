const strategies = {
    isVersion: function (value, errorMsg = '请输入正确的版本号!') {
        if (!/^\d+\.\d+\.\d+$/.test(value)) {
            return errorMsg
        }
    },
}

export default strategies