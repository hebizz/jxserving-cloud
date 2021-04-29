import axios from 'axios'
import { message } from 'antd'

export default function({method='get', url, params, data, successMsg, needErrMsg = true, ...rest}){
    return new Promise((resolve, reject) => {
        // console.log(method, url)
        axios({
            method,
            url,
            params,
            data,
            ...rest
        }).then(data=>{
            if(!data.data)throw Error('not response data')
            successMsg && message.success(successMsg)
            resolve(data.data)
            return data
        }).catch(err => {
            let response = err.response
            if (response && response.status === 499) {
                window.location.href = window.location.href.replace(/#\/.*/,'#/login')
            } else {
                needErrMsg && message.error('网络错误')
                reject(response)
            }
           
        })
    })
}