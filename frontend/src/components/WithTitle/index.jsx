import React, {useEffect, useState, useReducer} from 'react'
import styles from './index.module.scss'

export default function({title='', children}){
    return <div className = {styles.wrapper}>
        <span className = {styles.title}>{title}</span>
        {children}
    </div>
}