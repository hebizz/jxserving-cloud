import React, { useEffect, useState, useRef, useReducer } from 'react';
import ReactDOM from 'react-dom'
import styles from './index.module.scss'
import { Table, Button, Upload, Modal, Form, Input, message, Spin, Tag } from 'antd';
import ButtonGroup from 'antd/lib/button/button-group';
import api from 'utils/api';
import urls from 'asserts/urls';
import moment from 'moment';
import showInfoFormModal from 'utils/showInfoModal';
import { withKeyColumn } from 'sections';
const columns = [
	{
		dataIndex:'Id',
		title:'#',
		render:text=><span className = {'short1'}>{text}</span>
	},
	{
		dataIndex:'name',
		title:'名称',
		...withKeyColumn('name')
	},
	{
		dataIndex:'label',
		title:'标签',
		render:text=><span>{text.map((item, i)=><Tag key = {i}>{item}</Tag>)}</span>
	},
	{
		dataIndex:'operation',
		title:'编辑'
	}
]
function promising(cb){
	return new Promise((r, j) => {
		cb(r, j)
	})
}

export default function ({}) {
	const [showModal, setShowModal] = useState(false)
	const [uploading, setUploading] = useState(false)
	const [updateValue, update] = useReducer(value=>!value, false)

	const deleteModel = async row => {
		if(! await new Promise(r=>Modal.confirm({
			content:'确定删除此'+row.name+'标签集嘛？',
			onCancel:()=>r(0), 
			onOk:()=>r(1)
		})))return
		api({
			method:'post',
			url:urls.deleteLabel,
			data:{
				id:row.Id
			}
		}).then(()=>{
			update()
			message.success('删除成功')
		}).catch(()=>message.error('删除失败'))
	}

	const updateModel = async row=> {
		const onSubmit = (datas)=>{
			api({
				method:'post',
				url:urls.updateLabel,
				data:{
					id:row.Id,
					label:datas.label.split(','),
					name:datas.name
				}
			}).then(()=>{
				update()
				message.success('修改成功')
				showInfoFormModal.destroy()
			}).catch(err=>message.error('修改失败'))
		}
		showInfoFormModal.render(<InfoForm 
			onSubmit = {onSubmit} 
			defaultValues = {{label:row.label.join(','),name:row.name}}
		/>)
	}
	const addLabels = async () => {
		const onSubmit = (datas)=>{
			api({
				method:'post',
				url:urls.uploadLabel,
				data:{
					label:datas.label.split(','),
					name:datas.name
				}
			}).then(()=>{
				update()
				message.success('新建成功')
				showInfoFormModal.destroy()
			}).catch(err=>message.error('新建失败'))
		}
		showInfoFormModal.render(<InfoForm 
			onSubmit = {onSubmit} 
			defaultValues = {{}}
		/>)
	}
	const operationColumn = columns.find(item=>item.dataIndex === 'operation')
	operationColumn.render = (_, row) => {
		return <span>
			<ButtonGroup>
				<Button onClick = {()=>updateModel(row)}>编辑</Button>
				<Button onClick = {()=>deleteModel(row)}>删除</Button>
			</ButtonGroup>
		</span>
	}
	const [data, setData] = useState([])
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getLabel
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
					rowKey = {'Id'}
				/>
				<div className = {styles.tool}>
					<Button onClick = {addLabels}>
						新建标签集
					</Button>
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
					<InfoForm  onSubmit = {()=>{
						setShowModal(false)
						update()
					}}/>
				</Modal>
			</div>
		</Spin>
	);
}

const InfoForm = Form.create()(function({form, onSubmit, defaultValues}){
	const {getFieldDecorator} = form
	const handleSubmit = e => {
		e.preventDefault()
		form.validateFields((err, values)=>{
			console.log('err', err)
			if(!err){
				onSubmit && onSubmit(values)
			}
		})
	}
	return <Form labelCol={{ span: 5 }} wrapperCol={{ span: 12 }} onSubmit = {handleSubmit}  >
		<Form.Item label = '名字'>
			{
				getFieldDecorator('name', {
					rules:[{required:true, message:'请输入名字'}],
					initialValue:defaultValues.name
				})(<Input/>)
			}
		</Form.Item>
		<Form.Item label = '标签'>
			{
				getFieldDecorator('label', {
					rules:[{required:true, message:'多个标签用英文逗号,分隔'}],
					initialValue:defaultValues.label
				})(<Input/>)
			}
		</Form.Item>
		<Form.Item wrapperCol={{ span: 12, offset: 5 }}>
          <Button type="primary" htmlType="submit">
            提交
          </Button>
        </Form.Item>
	</Form>
})