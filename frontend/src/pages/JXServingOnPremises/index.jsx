import React, { useEffect, useState, useRef, useReducer } from 'react';
import styles from './index.module.scss'
import { Descriptions } from 'antd';
import api from 'utils/api';
import urls from 'asserts/urls';
import { from, timer } from 'rxjs';
import { expand, concatMap, map, tap } from 'rxjs/operators';
const {Item} = Descriptions
const initData = {
	backend:[],
	status:'',
	monitorstatus:{cpu:'', memory:''},
	ip:''
}
export default function () {
	const [ data, setData] = useState(initData)
	useEffect(()=>{
		const s = from(api({ method:'get', url:urls.getOnPremises})).pipe(
			expand(() => timer(3000).pipe(
				concatMap( () => from(api({method:'get', url:urls.getOnPremises})))
			)),
			map(res => res.data && res.data.result && res.data.result.length ? res.data.result[0] : initData),
			tap( data => setData(data))
		).subscribe()
		return () => {
			s.unsubscribe()
		}
	},[])
	return (
		<div className={styles.wrapper}>
			<Descriptions column={1} title='状态' bordered>
				<Item label='ip'>{data.ip}</Item>
				<Item label='运行状态'>{data.status}</Item>
				<Item label='后端数量'>{`${data.backend.length}(${data.backend.join(', ')})`}</Item>
				<Item label='监控状态'>{Object.keys(data.monitorstatus).map(key => <span key = {key}>
					<span>{key}: </span>
					<span>{data.monitorstatus[key]}</span>
					&nbsp;&nbsp;&nbsp;
				</span>)}</Item>
			</Descriptions>
		</div>
	)
}