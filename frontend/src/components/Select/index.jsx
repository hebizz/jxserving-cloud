import React from 'react'
import {Select} from 'antd'
import styles from './index.module.scss'

export default function JXSelect({wrapperStyle, dataSource, desc, ...rest}){
    // console.log(dataSource)
    return(
        <span className={styles.container} style={wrapperStyle}>
             {desc&&<span className={styles.desc}>
                 {desc}
             </span>}
             {dataSource && dataSource.length&&<Select defaultValue={dataSource[0].value} dropdownStyle={{fontSize:'12px'}}  {...rest}>
                 {dataSource.map(item=><Select.Option key={item.value} value={item.value}>{item.text}</Select.Option>)}
             </Select>}
        </span>
    )
}