import React, { useEffect, useMemo } from 'react'
import styles from './index.module.scss'
import {Spin} from 'antd'
import { usingSpecialPath } from 'utils/whenTestEnv'
import { fromEvent, of, BehaviorSubject } from 'rxjs'
import { filter, concatMap, withLatestFrom, tap } from 'rxjs/operators'

export default function LeftSide({
	streamLoading=false, 
	onImgSelect, 
	imgList=[], 
	listRef, 
	imgSelectedItem,
	whenItemChange = (item) => {console.log(item)},
	checkIsSameLogic = (originItem, newItem) => originItem===newItem
}){
	const imgList$ = useMemo(()=>new BehaviorSubject(imgList),[])
	useEffect(()=>{
		imgList$.next(imgList)
	},[imgList])
	useEffect(()=>{
		const s = fromEvent(document, 'keydown').pipe(
			filter(e=>e.code==='ArrowUp'||e.code==='ArrowDown'),
			concatMap(e=>{
				if(e.code==='ArrowUp')return jumpToPrevItem(imgSelectedItem, whenItemChange)
				else return jumpToNextItem(imgSelectedItem, whenItemChange)
				return of(1)
			})
		).subscribe()
		return () => s.unsubscribe()
	},[imgSelectedItem])

	const jumpToNextItem = (selectedItem, cb) => {
		return of(selectedItem).pipe(
			withLatestFrom(imgList$),
			tap(([selectedItem, imgList]) => {
				const index = imgList.findIndex(item=>item===selectedItem)
				if(!imgList.length)return cb('')
				if(index===-1){
					return cb(imgList[0])
				}
				const imgs = Array.from(document.getElementsByClassName('sideImg'))
				if(index===imgList.length-1){
					cb(imgList[0]||'')
					imgs[0].focus()
				} else {
					const item = imgList[index+1]
					cb(item)
					imgs[index+1].focus()
				}
			})
		)
	}

	const jumpToPrevItem = (selectedItem, cb) => {
		return of(selectedItem).pipe(
			withLatestFrom(imgList$),
			tap(([selectedItem, imgList]) => {
				const index = imgList.findIndex(item=>item===selectedItem)
				if(!imgList.length)return cb('')
				if(index===-1)return cb(imgList[0])
				const imgs = Array.from(document.getElementsByClassName('sideImg'))
				if(index===0){
					cb(imgList[imgList.length-1]||'')
					imgs[imgList.length-1].focus()
				}
				else {
					const item = imgList[index-1]
					cb(item)
					imgs[index-1].focus()
				}
			})
		)
	}

	return <div className={styles.leftSide}>
		<div className={styles.tool}>
		</div>
		<div className={styles.title}>
		</div>
		<div className = {styles.content}>
			<Spin spinning={streamLoading} >
				<div className={styles.list} ref={listRef}>
					{imgList.map((item, index) => <img
						tabIndex={-1}
						key={index}
						onClick={() => onImgSelect(item)}
						className={`${styles.item} ${checkIsSameLogic(imgSelectedItem, item) ? styles.selected : ''} sideImg`}
						src={item.path.includes('http')?item.path:usingSpecialPath('http://srv.cloud.jiangxingai.com:8088')+item.path}
						// src={item.path}
					/>)}
				</div>
			</Spin>
		</div>
	</div>
}