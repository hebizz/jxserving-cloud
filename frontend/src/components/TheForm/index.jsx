import React from 'react'
import { Form, Button, Spin} from 'antd'
import { useStateStream } from 'utils/useStateSream'

// const adjustForm = [
// 	{
// 		label:'datasetname',
// 		key:'name',
// 		options:{
// 			rules:[{required:true, message:'请输入datasetname'}]
// 		},
// 		component:<Input/>
// 	},
// 	{
// 		label:'datasettype',
// 		key:'type',
// 		options:{
// 			rules:[{required:true, message:'请输入datasettype'}]
// 		},
// 		component:<Select>
// 			<Select.Option value = {1}>voc</Select.Option>
// 			<Select.Option value = {2}>coco</Select.Option>
// 		</Select>
// 	}
// ]
const TheForm = Form.create()(function({formData, form, onSubmit}){
	const {getFieldDecorator, setFieldsValue, validateFields} = form
	const handleSubmit = e => {
		e.preventDefault()
		validateFields((err, values)=>{
			if(!err){
				onSubmit && onSubmit(values)
			}
		})
	}
	return <Form labelCol = {{span:8}} wrapperCol = {{span:12}} onSubmit = {handleSubmit}>
		{formData.map((item, i)=><Form.Item key={i} label = {item.label}>
			{
				getFieldDecorator(item.key, item.options)(item.component)
			}
		</Form.Item>)}
		<Form.Item wrapperCol={{ span: 12, offset: 8}}>
          <Button type="primary" htmlType="submit">
            提交
          </Button>
        </Form.Item>
	</Form>
})
export const SpiningForm = ({onSubmit, ...rest}) => {
	const [spining, setSpining] = useStateStream(false)
	const submit = async (values) => {
		setSpining(true)
		if(onSubmit)await onSubmit(values)
		setSpining(false)
	}
	return <Spin spinning = {spining}>
		<TheForm onSubmit = {submit} {...rest} />
	</Spin>
}
export default TheForm