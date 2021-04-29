import React from 'react'
import ReactDOM from 'react-dom'
import { Modal } from 'antd'

const showInfoFormModal = {
	mask:null,
	render:(component, {onCancel, ...rest}={}) => {
		showInfoFormModal.mask = document.createElement('div')
		document.body.appendChild(showInfoFormModal.mask)
		const modal = <Modal
			visible = {true}
			onCancel = {()=>{
				onCancel && onCancel()
				showInfoFormModal.destroy()
			}}
			footer={null}
			maskClosable = {false}
			{...rest}
		>
			{component}
		</Modal>
		ReactDOM.render(modal, showInfoFormModal.mask)
	},
	destroy: ()=>ReactDOM.unmountComponentAtNode(showInfoFormModal.mask)
}

export default showInfoFormModal