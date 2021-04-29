import React, { useState, useReducer, useMemo, useCallback, useEffect } from 'react';
import styles from './index.module.scss'
import Select from 'components/Select'
import Input from 'components/Input'
import { Table, DatePicker, Button, message } from 'antd';
import moment from 'moment'
import api from 'utils/api';
import urls from 'asserts/urls';

import Charts from './Charts'

const selectDataSource = [
	{
		value:'123',
		text:'haha'
	},
	{
		value:'1222',
		text:'eee'
	}
]
const columns = [
	{
		title:'逻辑',
		dataIndex:'logic'
	},
	{
		title:'字段',
		dataIndex:'field',
		render:(text) => (typeof text === 'string') ? fieldMap[text] : text
	},
	{
		title:'过滤值',
		dataIndex:'value'
	},
	{
		title:'操作',
		dataIndex:'operation'
	}
]
const logicDataSource = [
	{
		text:'AND',
		value:'AND'
	},
	// {
	// 	text:'OR',
	// 	value:'OR'
	// },
	// {
	// 	text:'NOT',
	// 	value:'NOT'
	// },
]
const fieldDataSource = [
	{
		text:'集群Id',
		value:'clusterid'
	},
	{
		text:'应用Id',
		value:'applicationid'
	},
	{
		text:'设备Id',
		value:'deviceid'
	},
	{
		text:'标签类型',
		value:'labeltype'
	}
]
const fieldMap = {
	timestamp:'时间区间',
	targetvalue:'分析项',
	clusterid:'集群Id',
	applicationid:'应用Id',
	deviceid:'设备Id',
	labeltype:'标签类型'
}

let mockData = []
const originData = [
	{
		logic:'REQ',
		field:'timestamp',
		value:[moment().subtract(1,'months').valueOf(), moment().valueOf()],
	},
	{
		logic:'REQ',
		field:'targetvalue',
		value:'reliability',

	}
]
const fromUser = localStorage.getItem('config')
if(fromUser){
	mockData = JSON.parse(fromUser)
}else mockData = JSON.parse(JSON.stringify(originData))
const defaultValue = {
	logic:'AND',
	field:'clusterid',
	value:''
}
export default function () {
	const [updateValue, update] = useReducer((updateValue)=>!updateValue, false)
	const [state, dispatch] = useReducer((state, action) => {
		if(action.type)return {...state, [action.type]:action.value}
		else return defaultValue
	}, defaultValue)
	const [needAnalysis, setNeedAnalysis] = useState(false)
	const [analysisResult, setAnalysisResult] = useState({data:[], stat:{}})
	const [dataSet, setDataSet] = useState([])
	const [selectedData, setSeletecdData] = useState('')
	const onChange = (row, value) => {
		row.value = value
		setNeedAnalysis(true)
	}
	const mockDataOperation = columns.find(item=>item.dataIndex==='operation')
	mockDataOperation.render = (text, row, index)=>(typeof row.logic==='string') && row.logic !=='REQ'?
		<Button onClick = {()=>{
			mockData.splice(index,1)
			update()
		}}>移除</Button>
		:text||''
	const valueColumn = columns.find(item=>item.dataIndex==='value')
	valueColumn.render = (text, row) => {
		if( (typeof text !== 'string') && !(text instanceof Array))return text
		switch(row.field){
			case 'timestamp':
				return <DatePicker.RangePicker style={{width:'400px'}} showTime={true} defaultValue={text.map(time=>moment(time))} onChange = {dates => onChange(row, dates.map(moment=>moment.valueOf()))}/>
			case 'targetvalue':
				return <Input value  = {text} disabled={true}/>
			case 'clusterid':
			default:
				return <Input defaultValue={text} onChange = {e => onChange(row, e.target.value)}/>
		}
	}


	const onAdd = ()=>{
		for(let key in state){
			if(!state[key]){
				return message.warn('请在'+key+'中输入正确的值')}
		}
		mockData.push(state)
		dispatch({})
	}
	const clear = () => {
		mockData = JSON.parse(JSON.stringify(originData))
		setNeedAnalysis(true)
		update()
	}
	const beginAnalysis = () => {
		//api  request
		const reqData = {
			timestamp:mockData[0].value.map(v=>~~(v/1000)),
			value:'reliability',
			cond:[]
		}
		if(selectedData)reqData.dataSet = selectedData
		reqData.cond = mockData.slice(2).map(item=>{
			delete item.key
			return item
		})
		const oneItemPush = (arr, item, value, type) => {
			arr.push({
				id:item.id,
				value:value,
				type
			})
		}
		api({
			method:'post',
			url:urls.analysisQuery,
			data: reqData
		}).then(res=>{
			const theData = res.data
			const data = Object.keys(theData.data).map(key=>({
				id:key,
				value:theData.data[key],
				type:'data'
			}))
			const ave = [],
				// cov=[],
				max=[],
				min=[],
				p95=[],
				p99=[]
				// std=[]
			data.forEach(item=>{
				oneItemPush(ave, item, theData.stat.ave, 'ave')
				// oneItemPush(cov, item, theData.stat.cov, 'cov')
				oneItemPush(max, item, theData.stat.max, 'max')
				oneItemPush(min, item, theData.stat.min, 'min')
				oneItemPush(p95, item, theData.stat.p95, 'p95')
				oneItemPush(p99, item, theData.stat.p99, 'p99')
				// oneItemPush(std, item, theData.stat.std, 'std')
			})
			const drawData = data.concat(ave,max,min,p95,p99)
			setAnalysisResult({data:drawData, stat:theData.stat})
			setNeedAnalysis(false)
		})
	}
	useEffect(()=>{
		window.onbeforeunload=()=>{
			localStorage.setItem('config', JSON.stringify(mockData))
		}
		api({
			method:'get',
			url:urls.getDataSet
		}).then(res=>{
			if(res && res.data && res.data.length){
				setDataSet(res.data.map(item=>({
					text:item.name,
					value:item.name
				})))
				setSeletecdData(res.data[0].name)
			}
		})
	},[])
	useEffect(()=>{
		setNeedAnalysis(true)
	},[mockData.toString()])

	const showData = mockData.concat({space:true},{
		logic:<Select dataSource={logicDataSource} value={state.logic} onChange={value=>dispatch({type:'logic', value})}/>,
		field:<Select dataSource = {fieldDataSource} value={state.field} onChange = {value=> dispatch({type:'field', value})}/>,
		value:<Input value={state.value} onChange = {e => dispatch({type:'value', value:e.target.value})}/>,
		operation:<Button onClick = {onAdd}>添加</Button>
	})
	showData.forEach((item, index)=>item.key=index)	
	return (
		<div className={styles.wrapper}>
			<div className = {styles.chooseData}>
				<Select desc='请选择要分析的数据源' onChange = {(v)=>setSeletecdData(v)} dataSource={dataSet} className = {styles.chooseDataSelect}/>				
			</div>
			<div>
				<Table columns={columns}
					dataSource={showData}
					pagination={false}
				/>
			</div>
			<div className = {styles.tool}>
				<Button type={needAnalysis?'primary':'default'} onClick = {beginAnalysis}>分析</Button>
				<Button onClick = {clear}>清除</Button>
			</div>
			<div>
				{analysisResult.stat.cov !== undefined ? <span className = {styles.label}>变异系数：{(+analysisResult.stat.cov)}</span>:''}
				{analysisResult.stat.std !== undefined ? <span className = {styles.label}>标准差：{(+analysisResult.stat.std)}</span>:''}
				<Charts data = {analysisResult.data} stat = {analysisResult.stat} />
			</div>
    </div>
	);
}
