import React, { useContext, useState, useEffect, useRef, useReducer } from 'react'
import styles from './index.module.scss'
import { useStateStream } from 'utils/useStateSream'
import { filter, tap, concatMap, map, debounceTime, withLatestFrom, catchError, scan, pairwise, take } from 'rxjs/operators'
import rxapi from 'utils/rxapi'
import { timer, fromEvent, of, merge, concat, Observable } from 'rxjs'
import { Spin, Select, Input } from 'antd'
import urls from 'asserts/urls'
import api from 'utils/api'
import WithTitle from 'components/WithTitle'
import Main from './Main'
import { useEventCallback } from 'rxjs-hooks'
import {withRouter} from 'react-router-dom'
import qs from 'querystring'
import exampleImg1 from 'asserts/imgs/example1.jpeg'
import exampleImg2 from 'asserts/imgs/example2.jpg'
import { Button } from 'antd/lib/radio'
import LeftSide from './LeftSide'
import showInfoFormModal from 'utils/showInfoModal'
import TheForm from 'components/TheForm'
var mockData = [
	{
		id: "xxxx1", name: 'exampleImg1', path: exampleImg1,
		timestamp: 2333333, label: [{ "name":"yanhuo","type": 1, "position": ['10', '10', '50','50'] }]
	},
	{
		id: "xxxx1", name: 'exampleImg1', path: exampleImg2,
		timestamp: 2333333, label: [{ "name":"yanhuo","type": 1, "position": ['10', '10', '50','50'] }]
	}
] 

export default withRouter(function ({row={}, history, location}) {
	const queryObj = qs.parse(location.search.slice(1))
	const [streamLoading, setStreamLoading, streamLoading$] = useStateStream(false)
	const [streamEnd, setStreamEnd, streamEnd$] = useStateStream(false)
	const [imgSelectedItem, setImgSelectedItem] = useState(mockData[0])
	const [imgList, setImgList, imgList$] = useStateStream(mockData)
	const listRef = useRef(null)
	const [isAdding, setIsAdding] = useState(false)
	const [_, update] = useReducer(value=>!value, false)
	useEffect(() => {
		api({
            method: 'post',
            url: urls.datasetManageQuery,
            data: {
                name: queryObj.name,
                startIndex: 1,
                offset: 10
            }
		})
		.then(res=>res.data || [])
		.then(list => {
			setImgList(list)
			setImgSelectedItem(list[0]||'')
		})

		const scrollListener = concat(fromEvent(listRef.current, 'scroll'), imgList$).pipe(
			filter(e => e.type?((e.target.scrollTop + e.target.clientHeight) >= (e.target.scrollHeight)):(e.length<10)),
			withLatestFrom(streamLoading$),
			withLatestFrom(streamEnd$),
			filter(([[e, streamLoading], streamEnd]) => !streamLoading && !streamEnd),
			scan((acc) => acc + 1, 0),
			filter(times => {
				if((10*times)>=Number(queryObj.total)){
					setStreamEnd(true)
					return false
				}else return true
			}),
			tap(() => setStreamLoading(true)),
			concatMap((times) => {
				return rxapi({
					method: 'post',
					url: urls.datasetManageQuery,
					data: {
						name: queryObj.name,
						startIndex: 1 + 10*times,
						offset: 10
					}
				})
			}),
			catchError(()=>{
				setStreamLoading(false)
				return of([])
			}),
			tap(() =>setStreamLoading(false)),
			tap(({data})=>!data || !data.length && setStreamEnd(true)),
			withLatestFrom(imgList$),
			map(([{data}, imgList]) => imgList.concat(data||[])),
		).subscribe(newList=>setImgList(newList))
		return () => {
			scrollListener.unsubscribe()
		}
	}, [])

	const onSubmit = data => api({
		method: 'post',
		url: urls.datasetManageUpdate,
		data:{
			labelName:queryObj.name,
			...data
		},
		successMsg:'保存成功'
	}).then(() => {
		const index = imgList.findIndex(item => item === imgSelectedItem) 
		api({
			method:'post',
			url: urls.datasetManageQuery,
            data: {
                name: queryObj.name,
                startIndex: index+1,
                offset: 1
            }
		}).then(({data})=> {
			// console.log(123,index, data[0], imgList.slice(0, index).concat(data, imgList.slice(index+1))) &&
			 setImgList(imgList.slice(0, index).concat(data, imgList.slice(index+1)))
			setImgSelectedItem(data[0])
		})
	})

	//当画完框后弹出的表单
	const whenDrawFinishAndSelectType = (data) => {
		return Observable.create(observer => {
			const onConfirm = (values) => {
				observer.next({ ...data, ...values })
				observer.complete()
				showInfoFormModal.destroy()
			}
			const onCancel = () => {
				observer.next('cancel')
				observer.complete()
			}
			const formItems = [
				{
					label: '框名字',
					key: 'name',
					options: { rules: [{ required: true, message: '请输入名字' }] },
					component: <Input autoFocus={true} />
				}
			]
			showInfoFormModal.render(
				<TheForm onSubmit={onConfirm} formData={formItems} />,
				{ onCancel }
			)
		})
	}
	const onImgSelect = item => setImgSelectedItem(item === imgSelectedItem ? '' : item)
	return <div className = {styles.wrapper}>
		<Button onClick = {()=>history.goBack()} className ={styles.goBack} >返回</Button>
		<LeftSide whenItemChange = {setImgSelectedItem} streamLoading = {streamLoading} onImgSelect = {onImgSelect} imgList = {imgList} listRef = {listRef} imgSelectedItem = {imgSelectedItem}/>
		<Main onCancel={null} isAdding = {isAdding} setIsAdding={setIsAdding} onSubmit = {onSubmit} imgSelectedItem = {imgSelectedItem} whenDrawFinishAndSelectType = {whenDrawFinishAndSelectType}/>
	</div>
})