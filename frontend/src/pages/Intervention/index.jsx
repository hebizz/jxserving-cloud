import React, { useEffect, useState, useReducer, useCallback } from 'react';
import styles from './index.module.scss'
import LeftSide from 'pages/DatasetMange/LeftSide';
import Main from 'pages/DatasetMange/Main';
import { Icon, Select, Spin, message } from 'antd';
import { useStateStream } from 'utils/useStateSream';
import showInfoFormModal from 'utils/showInfoModal';
import TheForm from 'components/TheForm';
import api from 'utils/api';
import urls from 'asserts/urls';

import Img1 from 'asserts/imgs/example1.jpeg'
import Img2 from 'asserts/imgs/example2.jpg'
import { Observable, timer, of, forkJoin } from 'rxjs';
import { expand, withLatestFrom, filter, map, delay, tap, concatMap, catchError, concat, retry } from 'rxjs/operators';
import rxapi from 'utils/rxapi';

const alarmTypes = {
    "construction_machinery": "附近存在施工机械",
    "crane": "附近存在吊车",
    "tower_crane": "附近存在塔吊",
    "foreign_conductor": "导线异物提醒",
    "fire": "烟雾山火提醒",
}

const unicodeToStr = unicode => {
	try{
		return unicode.split('\\u').slice(1).map(v=>(JSON.parse('{"str": "'+ '\\u'+v +'"}')).str).join('')
	}catch(err){
		return unicode
	}
}
const strToUnicode = str => {
	return str.split('').map(v=>"\\u"+v.charCodeAt().toString(16)).join('')
}
const getUrlBase64 = (url, ext='jpeg') => {
	return new Promise((resolve, reject)=>{
		var canvas = document.createElement('canvas');
		var ctx = canvas.getContext('2d')
		var img = new Image
		img.crossOrigin = 'Anonymous'
		img.src = url
		img.onload = () => {
			const a = {aa:img}
			console.log(a)
			canvas.width = img.width
			canvas.height = img.height
			ctx.drawImage(img, 0, 0, canvas.width, canvas.height)
			var dataURL = canvas.toDataURL('image/'+ext);
			resolve(dataURL)
			canvas = null
		}
	})
}

// const mockStream = () => {//人工干预接口返回的格式
// 	return Observable.create(observer => {
// 		observer.next({
// 			data: [
// 				{
// 					"cluster_id": "sgcc",
// 					"title": "\\u9644\\u8fd1\\u5b58\\u5728\\u540a\\u8f66",
// 					"alert_type": "crane",
// 					"description": "reliability: 0.29938915371894836",
// 					"alert_position": [
// 						{
// 							"left_x": 100,
// 							"left_y": 100,
// 							"right_x": 200,
// 							"right_y": 200
// 						}
// 					],
// 					"event_id": "5fc8752c-c599-11e9-94fe-0242ac110005"+Math.random(),//此字段可作为唯一id
// 					"device_id": "5d4a9c1d97e91b6667bcec93",
// 					"image": Math.random()>0.5?Img1:Img2,
// 					"image_path":Math.random()>0.5?Img1:Img2,
// 					// "image": huolizi,
// 				}
// 			]
// 		})
// 		observer.complete()
// 	})
// }

// const imgList = [//组件需要的格式
// 	{
// 		path:'/sdf/sdfd.jpg',
// 		label:[
// 			{
// 				name:'xxx',
// 				type:'sdsa',
// 				position:['x1','y1','x2','y2']
// 			}
// 		]
// 	}
// ]
const changeImgFormatFromOrigin = originImgItem => {
	return ({
		event_id:originImgItem.event_id,
		path:originImgItem.image_path,
		label:originImgItem.alert_position.map(item=>({
			name:originImgItem.alert_type,
			type:originImgItem.alert_type,
			position:[
				item.left_x,
				item.left_y,
				item.right_x,
				item.right_y,
			]
		}))
	})
}
const changeImgListFormatFromOrigin = (imgList) => {
	return imgList.map(item=>changeImgFormatFromOrigin(item))
}

//添加模型的组件
const AddIcon = ({whenSelect}) => {
	const [models, setModels] = useState([])
	useEffect(()=>{
		api({
			method:'get',
			url:urls.getModel
		}).then(res=>{
			setModels(res.data.result || [])
		})
	},[])
	const formData = [
		{
			label: 'model',
			key: 'model',
			options: {
				rules: [{ required: false, message: '请选择模型' }]
			},
			component: <Select>
				{models.map( (item, i) => <Select.Option key = {i} value = {JSON.stringify(item)}>{item.name}</Select.Option>)}
			</Select>
		},
	]
	const handleSubmit = values => {
		whenSelect(values.model)
		showInfoFormModal.destroy()
	}
	const handleClick = () => {
		showInfoFormModal.render(<TheForm formData = {formData} onSubmit = {handleSubmit}/>)
	}
	return <div className = {styles.addIcon} onClick = {handleClick}>
		<Icon type="plus-circle" style = {{fontSize:24}} />
	</div>
}
//每个模型列表的组件
const ReferenceModel = ({
	deleteModel, 
	modelData={},//该项的model数据
	detectResult = [],//监测到的结果
	imgSelectedItem = {},//选择到的图片的项目
	serialNum,
	modelSwitched={}, //判断切换的模型是否为此model
	onSwitch, 
	onDetect
}) => {
	// console.log(detectResult, modelData, imgSelectedItem)
	const newItem = {
		path:imgSelectedItem.path,
		label:detectResult
	}
	return <div className = {styles.referenceModel}>
		<span>{modelData.name||''}</span>
		<Icon type="minus-circle" style = {{fontSize:20}} title = '删除' onClick = {() => deleteModel(serialNum)}/>
		{modelData.name === modelSwitched.name ?  
			<Icon type="smile" style = {{fontSize:20}} title = '图片识别' onClick = {()=>onDetect && onDetect()}/>: 
			<Icon type="sync" style = {{fontSize:20}} title = '切换模型' onClick = {()=>onSwitch && onSwitch(modelData)}/>}
		{newItem ? <Main needTool = {false} imgSelectedItem = {newItem}/> : ''}
	</div>
}

export default function () {
	const [modelList, setModelList] = useStateStream([])//下面对比的模型列表
	const [imgList, setImgList, imgList$] = useStateStream([])
	const [isAdding, setIsAdding, isAdding$] = useStateStream(false)
	const [loading, setLoading, isLoading$] = useStateStream(false)
	const [imgSelectedItem, setImgSelectedItem] = useState('')
	const [_, update] = useReducer(value=>!value, false)
	const [modelSelected, setModelSelected, modelSelected$] = useStateStream(null)
	const [modelSwitched, setModelSwitched, modelSwitched$] = useStateStream({})
	const [detectResult, setDetectResult] = useStateStream({})
	useEffect(()=>{
		const getImgList = rxapi({url:urls.verifyInfo}).pipe(
		// const getImgList = mockStream().pipe(
			expand(()=>timer(3000).pipe(concatMap(()=>rxapi({url:urls.verifyInfo})))),
			// expand(()=>timer(3000).pipe(concatMap(()=>mockStream()))),
			withLatestFrom(isAdding$),
			filter(([_, isAdding]) => !isAdding),
			map(([res,_])=> res.data),
			filter(data=>data!==''),
			withLatestFrom(imgList$),
			map(([data,imgList]) => imgList.concat(data)),
			tap(data=>setImgList(data)),
			map(data=>data.filter(item=>item.time*1000>Date.now()))
		).subscribe(data=>{
			// console.log(data)
		})

		//切换模型
		const switchModel = modelSelected$.pipe(
			filter(v=>v),
			tap(()=>setLoading(true)),
			concatMap( modelSelected=>rxapi({
				method:'post',
				url:urls.modelSwitch,
				data:{
					id:modelSelected.name,
					mode:'frozen'
				}
			}).pipe(map(()=>modelSelected))),
			catchError(err=>of(true)),
			tap(modelSelected=>setModelSwitched(modelSelected))
		).subscribe(()=>setLoading(false))
		return () => {
			getImgList.unsubscribe()
			switchModel.unsubscribe()
		}
	},[])

	//图片列表的操作
	const deleteImgSelected = (selectedItem) => {
		jumpToNextItem(selectedItem)
		const index = imgList.findIndex((item)=>item.event_id===selectedItem.event_id)
		imgList.splice(index,1)
		update()
	}
	const jumpToNextItem = selectedItem => {
		const index = imgList.findIndex((item)=>item.event_id===selectedItem.event_id)
		if(!imgList.length)return setImgSelectedItem('')
		if(index===-1){
			return setImgSelectedItem(changeImgFormatFromOrigin(imgList[0]))
		}
		const imgs = Array.from(document.getElementsByClassName('sideImg'))
		if(index===imgList.length-1){
			setImgSelectedItem(changeImgFormatFromOrigin(imgList[0])||'')
			imgs[0].focus()
		} else {
			const item = imgList[index+1]
			setImgSelectedItem(changeImgFormatFromOrigin(item))
			imgs[index+1].focus()
		}
	}
	const onImgSelect = item => {
		setImgSelectedItem(item === imgSelectedItem ? '' : item)
		setDetectResult({})
	}

	//当画完框后弹出的表单
	const whenDrawFinishAndSelectType = (data) => {
		return Observable.create(observer => {
			const onConfirm = (values) => {
				observer.next({ ...data, ...values, type:values.name })
				observer.complete()
				showInfoFormModal.destroy()
			}
			const onCancel = () => {
				observer.next('cancel')
				observer.complete()
			}
			const formItems = [
				{
					label: '告警类型',
					key: 'name',
					options: {
						rules: [{ required: true, message: '请选择类型' }]
					},
					component: <Select autoFocus={true}>
						{Object.keys(alarmTypes).map(key=><Select.Option value = {key}>{key}</Select.Option>)}
					</Select>
				}
			]
			showInfoFormModal.render(
				<TheForm onSubmit={onConfirm} formData={formItems} />,
				{ onCancel }
			)
		})
	}
	//提交干预的结果
	const onSubmit = (imgSelectedItem, alreadyEdit) => {
		const submit = {
			event_id:imgSelectedItem.event_id,
			ignore:false
		}
		if(alreadyEdit){
			const resultsSplitByName = imgSelectedItem.label.reduce((pre, cur)=>{
				if(!pre.hasOwnProperty(cur.name))pre[cur.name] = []
				pre[cur.name].push(cur)
                return pre
			},{})
			submit.override = Object.keys(resultsSplitByName).map(name=>({
				alert_type:name,
				title:alarmTypes[name],
				alert_position:resultsSplitByName[name].map(item=>({
					left_x: item.position[0],
					left_y: item.position[1],
					right_x: item.position[2],
					right_y: item.position[3],
				}))
			}))
		}
		api({
			method:'post',
			url:urls.update,
			data:submit
		}).then(()=>deleteImgSelected(imgSelectedItem))
	}
	//无需告警的回调
	const onCancel = () => {
        api({
            method:'post',
            url:urls.update,
            data:{
                event_id:imgSelectedItem.event_id,
                ignore:true
            }
        }).then(()=>deleteImgSelected(imgSelectedItem))
	}


	//添加模型
	const addModel = (item) => {
		// console.log(item)
		const itemParse = JSON.parse(item)
		setModelList(modelList.concat([itemParse||{}]))
		setModelSelected(itemParse||{})
	}
	const deleteModel = (theIndex) => {
		setModelList(modelList.slice(0, theIndex).concat(modelList.slice(theIndex+1)))
	}
	const checkIsSameLogic = (originItem, newItem) => {
		return originItem.event_id === newItem.event_id ? true : false
	}

	const onDetect = async id => {
		if(!imgSelectedItem)return message.error('请先选择图片')
		const imgBase64 = await getUrlBase64(imgSelectedItem.path)
		setLoading(true)
		api({
			method:'post',
			url:urls.modelDetect,
			data:{
				image:imgBase64
			}
		})
		.then(res=>{
			const resultList = res.data.result.map(item=>({
				name:item[0]+' '+Number(item[1]).toFixed(2),
				// type:item[0],
				position:[
					item[2],
					item[3],
					item[4],
					item[5],
				]
			}))
			setDetectResult({...detectResult, [modelSelected.name]:resultList})
		})
		.finally(res=>{
			setLoading(false)
		})
	}
	const onSwitch = (item) => {
		setModelSelected(item)
	}
	return (
		<Spin spinning={loading} wrapperClassName = {styles.spinWrapper}>
			<div className={styles.wrapper}>
				<LeftSide 
					onImgSelect = {onImgSelect} 
					imgList = {changeImgListFormatFromOrigin(imgList)} 
					imgSelectedItem = {imgSelectedItem}
					checkIsSameLogic = {checkIsSameLogic}
					whenItemChange = {setImgSelectedItem}
				/>
				<Main 
					isAdding = {isAdding} 
					setIsAdding = {setIsAdding} 
					imgSelectedItem = {imgSelectedItem} 
					onSubmit = {onSubmit} 
					onCancel = {onCancel}
					whenDrawFinishAndSelectType = {whenDrawFinishAndSelectType}
				/>
				<div className = {styles.reference}>
					{modelList.map((item, i)=> <ReferenceModel 
						key={i} 
						modelSwitched = {modelSwitched} 
						modelData={item} 
						detectResult = {detectResult[item.name]}
						deleteModel = {deleteModel} 
						serialNum = {i}
						onDetect = {onDetect}
						onSwitch = {onSwitch}
						imgSelectedItem = {imgSelectedItem}
					/>)}
					<AddIcon whenSelect = {addModel}/>
				</div>
			</div>
		</Spin>
	);
}
