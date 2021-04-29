import React, { useEffect, useState, useRef, useReducer } from 'react';
import styles from './index.module.scss'
import { Table, Button, Upload, Modal, Form, Input, message, Spin, Select, DatePicker } from 'antd';
import ButtonGroup from 'antd/lib/button/button-group';
import api from 'utils/api';
import urls from 'asserts/urls';
import { SpiningForm } from 'components/TheForm';
import showInfoFormModal from 'utils/showInfoModal';
import moment from 'moment';
import { withKeyColumn } from 'sections';
import {withRouter} from 'react-router-dom'
import qs from 'querystring'
import { usingSpecialPath } from 'utils/whenTestEnv';
const columns = [
	{ dataIndex:'_id', title:'#', render:(_, _2, index) => index+1},
	{ dataIndex:'name', title:'名称', ...withKeyColumn('name')},
	{ dataIndex:'marker/total', title:'已标注数量/总数量', render:(_, row) => (row.marker===undefined || row.total===undefined) ? '' : (row.marker+' / '+row.total)},
	{ dataIndex:'operation', title:'操作'}
]
const mockData = [
	{ name:'sgcc', framework:'tensorflow-serving', version:'1.0.2', state:'craft', title:'a1b2c3', size:'2.5MB' }
]

function promising(cb){
	return new Promise((r, j) => {
		cb(r, j)
	})
}

export default withRouter(function ({history}) {
	const [showModal, setShowModal] = useState(false)
	const [uploading, setUploading] = useState(false)
	const [updateValue, update] = useReducer(value=>!value, false)
	const [file, setFile] = useState(null)
	const action = (file) => new Promise(()=>{
		setFile(file)
		setShowModal(true)
	})

	const download = (row, times) => {
		const data = {
			name:row.name,
			type:1,
			downType:1
		}
		if(times){
			data.startTime = times[0]
			data.endTime = times[1]
			data.downType = 2
		}
		setUploading(true)
		api({
			method:'post',
			url:urls.downloadDataSet,
			data:data,
		}).then(res=>{
			if(!res.data || !res.data.path)message.error(res.msg)
			const path = res.data.path
			const link = document.createElement('a')
			link.download = 'file.zip'
			link.href = usingSpecialPath('http://srv.cloud.jiangxingai.com:8088') + path
			link.click()
		}).finally(() =>setUploading(false))
	}

	const selectTimeForDownload = async row => {
		const formData = [
			{
				label:'选择时间',
				key:'times',
				options:{
					rules:[{required:true, message:'请输入datasettype'}],
					initialValue:[moment().subtract(1, 'months'), moment()]
				},
				component:<DatePicker.RangePicker />
			}
		]
		showInfoFormModal.render(<SpiningForm
			formData = {formData}
			onSubmit = {(values)=>{
				showInfoFormModal.destroy()
				download(row, values.times.map(time=>time.format('YYYY-MM-DD')))
			}}
		/>)
	}

	const deleteModel = async row => {
		if(! await new Promise(r=>Modal.confirm({
			content:'确定移除'+row.name+'数据集嘛？',
			onCancel:()=>r(0), 
			onOk:()=>r(1)
		})))return
		api({
			method:'post',
			url:urls.deleteDataSet,
			data:{ name:row.name }
		}).then(()=>{
			update()
			message.success('移除成功')
		}).catch(()=>message.error('移除失败'))
	}

	const operationColumn = columns.find(item=>item.dataIndex === 'operation')
	operationColumn.render = (_, row) => {
		const query = qs.stringify(row)
		return <span>
			<ButtonGroup>
				<Button onClick = {()=>deleteModel(row)} disabled = {row.name === 'default'?true:false}>移除</Button>
				<Button onClick = {()=>history.push('/datasetmanage?'+query)}>
					管理
				</Button>
				<Button onClick = {()=>download(row)}>下载</Button>
				<Button onClick = {()=>selectTimeForDownload(row)}>部分下载</Button>
			</ButtonGroup>
		</span>
	}

	const formRef = useRef(null)
	const [data, setData] = useState([])
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getDataSet
		}).then(res=>{
			setData(res.data)
		})
	}, [updateValue])
	return (
		<Spin spinning={uploading}>
			<div className={styles.wrapper}>
				<Table 
					columns = {columns}
					pagination={false}
					dataSource = {data}
					rowKey = {record=>record.name+record.marker}
				/>
				<div className = {styles.tool}>
					<Upload
						showUploadList = {false}
						name = 'file'
						action = {action}
					>
						<Button type="primary">
							新增数据集
						</Button>
					</Upload>
				</div>
				<Modal
					title = '请选择类型'
					visible = {showModal}
					onCancel = {async () => { 
						const confirmCancel = await promising((resolve)=>Modal.confirm({
							content:'确定取消上传吗？',
							onOk:()=>resolve(1),
							onCancel:()=>resolve(0)
						}))
						if(confirmCancel){
							setShowModal(false)
						}
					}}
					footer = {null}
				>
					<InfoForm file={file} onSubmit = {()=>{
						setShowModal(false)
						update()
					}}/>
				</Modal>
			</div>
		</Spin>
	);
})
const InfoForm = ({file, onSubmit}) => {
	const formData = [
		{
			label:'类型',
			key:'type',
			options:{
				rules: [{ required: true, message: '请输入类型' }]
			},
			component:<Select>
				<Select.Option value={1}>voc</Select.Option>
				<Select.Option value={2}>coco</Select.Option>
			</Select>
		},
		{
			label:'名字',
			key:'name',
			options:{
				rules: [{ required: true, message: '请输入类型' }]
			},
			component:<Input />
		},
	]
	const handleSubmit = values => {
		const formData = new FormData()
		formData.append('file', file)
		formData.append('type', values.type)
		formData.append('name', values.name)
		return api({
			method:'post',
			url:urls.uploadDataSet,
			headers:{'Content-Type':'multipart/form-data'},
			data:formData,
			needErrMsg:false
		}).then(res=>{
			message.success('上传信息成功')
			onSubmit && onSubmit()
		}).catch(err=>{
			if(err.status===500)message.error('上传文件格式不对')
			else message.error('网络错误')
		})
	}
	return <SpiningForm onSubmit = {handleSubmit} formData = {formData}/>
}