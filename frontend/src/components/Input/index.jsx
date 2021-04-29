import React from 'react'
import {Select, Input} from 'antd'
import styles from './index.module.scss'

export default function JXSelect({wrapperStyle, dataSource, desc, ...rest}){
    return(
        <span className={styles.wrapper} style={wrapperStyle}>
             {desc&&<span className={styles.desc}>
                 {desc}
             </span>}
             <Input {...rest}/>
        </span>
    )
}