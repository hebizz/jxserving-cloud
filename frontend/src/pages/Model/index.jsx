import React, { useEffect, useState, useRef, useReducer } from 'react';
import styles from './index.module.scss'
import { Table, Button, Upload, Modal, Form, Input, message, Spin, Select } from 'antd';
import ButtonGroup from 'antd/lib/button/button-group';
import api from 'utils/api';
import urls from 'asserts/urls';
import moment from 'moment';
import showInfoFormModal from 'utils/showInfoModal';
import { usingSpecialPath } from 'utils/whenTestEnv';
import { SpiningForm } from 'components/TheForm';
import strategies from 'utils/strategies';
import { _idColumn } from 'sections';
import { withKeyColumn, getFiltersWhenDataing } from 'sections';
import { mergeToOriginObj } from 'sections';
import calSize from 'utils/calSize'
const columns = [
	{ dataIndex:'_id', title:'#', render:text=><span className = {'short1'}>{text}</span>,
	},
	{ dataIndex:'name', title:'名称',
		...withKeyColumn('name')
	},
	{ dataIndex:'framework', title:'框架',
		...withKeyColumn('framework')
	},
	{ dataIndex:'version', title:'版本',
		...withKeyColumn('framework')
	},
	{ dataIndex:'ispublished', title:'发布状态', render:(text) => text?'已发布':'草稿中',
		filters:[
			{ text:'已发布', value:true },
			{ text:'草稿中', value:false },
		],
		onFilter: (value, record) => record.ispublished == value
	},
	{ dataIndex:'relatedlabel', title:'相关标签'},
	{ dataIndex:'createdtimestamp', title:'创建时间',
		render:text => text?moment(text*1000).format('YYYY/MM/DD HH:MM:SS'):text,
		...withKeyColumn('createdtimestamp')
	},
	{ dataIndex:'size', title:'大小',
		render:text=>{
			const size = calSize(Number(text))
			return size
		}
	},
	{ dataIndex:'operation', title:'操作' }
]
const adjustColumns = [
	{ dataIndex:'key', title:'key'},
	{ dataIndex:'value', title:'value'}
]
const historyColumns = [
	{ dataIndex:'_id', title:'#', render:text=><span className = {'short1'}>{text}</span>, width:150 },
	{ dataIndex:'evaldataset', title:'数据集', width:150 },
	{ dataIndex:'datesettype', title:'类型', render:text=>text==1?'voc':'coco', width:150 },
	{ dataIndex:'errorrate', title:'ErrorRate', width:150, render:text => Number(text).toFixed(4) },
	{ dataIndex:'leakrate', title:'LeakRate', width:150, render:text => Number(text).toFixed(4) },
	{ dataIndex:'meanap', title:'MeaNap', width:150, render:text => Number(text).toFixed(4)},
	{ dataIndex:'createdtimestamp', title:'测试时间', render:text => text?moment(text*1000).format('YYYY/MM/DD HH:MM:SS'):text}
]
// const mockData = [
// 	{
// 		name:'sgcc',
// 		framework:'tensorflow-serving',
// 		version:'1.0.2',
// 		state:'craft', title:'a1b2c3',
// 		size:'2.5MB',
// 	}
// ]
function promising(cb){
	return new Promise((r, j) => {
		cb(r, j)
	})
}

export default function ({}) {
	const [showModal, setShowModal] = useState(false)
	const [uploading, setUploading] = useState(false)
	const [updateValue, update] = useReducer(value=>!value, false)
	const formRef = useRef(null)
	const [data, setData] = useState([])
	const [labelData, setLabelData] = useState([])

	const adjustForm = [
		{
			label:'选择测试的数据集',
			key:'name',
			options:{
				rules:[{required:true, message:'请输入datasetname'}]
			},
			component:<Select>
				{labelData.map((item, i)=><Select.Option key = {i} value={item.name}>{item.name}</Select.Option>)}
			</Select>
		},
		{
			label:'选择数据集类型',
			key:'type',
			options:{
				rules:[{required:true, message:'请输入datasettype'}]
			},
			component:<Select>
				<Select.Option value = {1}>voc</Select.Option>
				<Select.Option value = {2}>coco</Select.Option>
			</Select>
		}
	]
	
	const action = (file) => new Promise(()=>{
		setUploading(true)
		const formData = new FormData()
		formData.append('file', file)
		api({
			method:'post',
			url:urls.uploadModel,
			// url:urls.getModel,
			data:formData,
			headers:{'Content-Type':'multipart/form-data'}
		}).then(res=>{
			setUploading(false)
			setShowModal(true)
		})
	})

	const deleteModel = async row => {
		if(! await new Promise(r=>Modal.confirm({
			content:'确定移除'+row.name+'模型嘛？',
			onCancel:()=>r(0), 
			onOk:()=>r(1)
		})))return
		api({
			method:'post',
			url:urls.deleteModel,
			data:{
				modelmd5:row.modelmd5
			}
		}).then(()=>{
			update()
			message.success('删除成功')
		}).catch(()=>message.error('删除失败'))
	}

	const publishModel = async row=> {
		if(! await new Promise(r=>Modal.confirm({
			content:'确定发布'+row.name+'模型吗?',
			onCancel:()=>r(0), 
			onOk:()=>r(1)
		})))return
		api({
			method:'post',
			url:urls.publishModel,
			data:{modelmd5:row.modelmd5},
			needErrMsg:false
		}).then(()=>{
			update()
			message.success('发布成功')
		})
		.catch(()=>message.error('发布失败'))
	}

	const unPublishModel = async row => {
		if(! await new Promise(r => Modal.confirm({
			content:'确定取消发布'+row.name + '模型吗?',
			onCancel:()=>r(0), 
			onOk:()=>r(1)
		})))return
		api({
			method:'post',
			url:urls.unpublishModal,
			data: { modelmd5:row.modelmd5 },
			needErrMsg:false
		}).then(() => {
			update()
			message.success('取消发布成功')
		})
		.catch(()=>message.error('取消发布失败'))
	}
	const adjustModel = async row => {
		showInfoFormModal.render(<SpiningForm
			formData = {adjustForm}
			onSubmit = {(values)=>{
				return api({
					method:'post',
					url:urls.adjust,
					data:{
						modelmd5:row.modelmd5,
						...values
					},
					needErrMsg:false
				}).then(res=>{
					// message.success('模型评价成功')
					const result = JSON.parse(res.data.result)
					const results = Object.keys(result).map( key => ({
						key,
						value:Number(result[key]).toFixed(5)
					}))
					Modal.info({
						title:'评价结果',
						content:<Table pagination={false} dataSource = {results} columns = {adjustColumns}/>
					})
					showInfoFormModal.destroy()
				}).catch((err)=>{
					message.error('模型评价失败')
				})
			}}
		/>)
	}
	const showHistory = async row => {
		const data = (await api({
			method:'post',
			url:urls.modelHistory,
			data:{
				modelmd5:row.modelmd5
			}
		})).data || {}
		Modal.info({
			content:<Table rowKey = {record => record._id} scroll={data.result && data.result.length>=5?{y:250}:{}} pagination={false} columns = {historyColumns} dataSource = {data.result}/>,
			className:styles.historyModal
		})
	}
	const operationColumn = columns.find(item=>item.dataIndex === 'operation')
	operationColumn.render = (_, row) => {
		return <span>
			<ButtonGroup>
				<Button >
					<a href={usingSpecialPath('http://srv.cloud.jiangxingai.com:8088')+row.modelpath} download="file">下载</a>
				</Button>
				<Button disabled = {row.ispublished ? true : false} onClick = {()=>deleteModel(row)}>移除</Button>
				{!row.ispublished ? <Button onClick = {()=>publishModel(row)}>发布</Button>:<Button onClick = {() => unPublishModel(row)}>取消发布</Button>}
				<Button onClick = {()=>adjustModel(row)}>模型评价</Button>
				<Button onClick = {()=>showHistory(row)}>历史评价</Button>
			</ButtonGroup>
		</span>
	}
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getModel
		}).then(res=>{
			setData(res.data.result)
		})
	}, [updateValue])
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getDataSet
		}).then(res=>{
			setLabelData(res.data||[])
		})
	},[])
	

	const frameworkColumn = columns.find(item=>item.dataIndex==='framework')
	if(frameworkColumn)
		mergeToOriginObj(frameworkColumn, getFiltersWhenDataing( data, 'framework')) 
	const versionColumn = columns.find(item=>item.dataIndex === 'version')
	if(versionColumn)mergeToOriginObj(versionColumn, getFiltersWhenDataing(data, 'version'))
	return (
		<Spin spinning={uploading}>
			<div className={styles.wrapper}>
				<Table 
					columns = {columns}
					pagination={false}
					dataSource = {data}
					rowKey = {'_id'}
				/>
				<div className = {styles.tool}>
					<Upload
						showUploadList = {false}
						name = 'file'
						action = {action}
					>
						<Button type="primary">
							上传模型
						</Button>
					</Upload>
				</div>
				<Modal
					title = '上传成功，请输入详细信息'
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
					<InfoForm ref = {formRef} onSubmit = {()=>{
						setShowModal(false)
						update()
					}} closeModal = {() =>setShowModal(false)}/>
				</Modal>
			</div>
		</Spin>
	);
}







const InfoForm = Form.create()(function({form, onSubmit, closeModal}){
	const {getFieldDecorator, setFieldsValue} = form
	const [labels, setLabels] = useState([])
	const [frameWorks, setFrameWorks] = useState([])
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getLabel
		}).then(res=>{
			setLabels(res.data)
			setFieldsValue({
				relatedlabel:res.data[0].name
			})
		})
		api({
			method:'get',
			url:urls.getFrameWork,
		}).then(res => {
			if(res.data && res.data.result){
				setFrameWorks( res.data.result)
				setFieldsValue({
					framework:res.data.result[0]
				})
			}
		})
	},[])
	const handleSubmit = e => {
		e.preventDefault()
		form.validateFields((err, values)=>{
			if(!err){
				api({
					method:'post',
					url:urls.uploadModelMsg,
					data:values,
					needErrMsg:false
				}).then(res=>{
					onSubmit && onSubmit()
					message.success('上传信息成功')
				}).catch(err => {
					if(err.status === 500){
						message.error(err.data.msg)
						closeModal && closeModal()
					}
				})
			}
		})
	}
	return <Form labelCol={{ span: 5 }} wrapperCol={{ span: 12 }} onSubmit = {handleSubmit}  >
		<Form.Item label = '名称'>
			{
				getFieldDecorator('name', {
					rules:[{required:true, message:'请输入名称'}]
				})(<Input/>)
			}
		</Form.Item>
		<Form.Item label = '框架'>
			{
				getFieldDecorator('framework', {
					rules:[{required:true, message:'请输入框架'}]
				})(<Select>
					{frameWorks.map((item, i)=><Select.Option value = {item} key = {i}>
						{item}
					</Select.Option>)}
				</Select>)
			}
		</Form.Item>
		<Form.Item label = '版本'>
			{
				getFieldDecorator('version', {
					rules:[
						{required:true, message:'请输入版本'},
						{
							validator: (rule, value, callback) => {
								callback(strategies.isVersion(value))
							}
						}
					],
				})(<Input placeholder='格式为X.X.X, X为数字'/>)
			}
		</Form.Item>
		<Form.Item label = '标签集'>
			{
				getFieldDecorator('relatedlabel', {
					rules:[{required:true, message:'请输入标签集'}]
				})(<Select>
					{labels.map((item, i)=><Select.Option value = {item.name} key = {i}>
						{item.name}
					</Select.Option>)}
				</Select>)
			}
		</Form.Item>
		<Form.Item wrapperCol={{ span: 12, offset: 5 }}>
          <Button type="primary" htmlType="submit">
            提交
          </Button>
        </Form.Item>
	</Form>
})
