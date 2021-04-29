import React, { useContext, useRef, useEffect, useState, useCallback, useReducer, useMemo } from 'react'
import styles from './index.module.scss'
import HoldEvent from 'utils/holdEvent';
import { useEventCallback } from 'rxjs-hooks';
import { filter, withLatestFrom, map, tap, concatMap, delay, takeUntil } from 'rxjs/operators';
import { fromEvent, of, BehaviorSubject, combineLatest } from 'rxjs';
import { Button, Input, Switch, Icon } from 'antd';
import WithTitle from 'components/WithTitle';
import DrawRect from 'utils/DrawRect';
import { useStateStream } from 'utils/useStateSream';
import { usingSpecialPath } from 'utils/whenTestEnv';

export default function (props) {
    const {
        isAdding = false,
        setIsAdding = () => { },
        imgSelectedItem = '', 
        onSubmit = (data) => {},
        onCancel = (imgSelectedItem) => {},
        whenDrawFinishAndSelectType = (data)=>of(data),
        needTool = true
    } = props
    const isAdding$ = useMemo(() => new BehaviorSubject(isAdding),[])
    const imgRef = useRef(null)
    const canvasRef = useRef(null)
    const [{ drawRect, holdEvent }, setInitData] = useState({ drawRect: null })
    const [results, setResults, results$] = useStateStream([])
    const [resizing, resizingUpdate] = useReducer(v => !v, false)
    const [redoList, redoDispatch] = useReducer((v, action)=>{
        switch(action.type){
            case 'push':
                v.push(action.value)
                return v
            case 'reset':
                return [action.value]
        }
    },[])
    const [theResultIndex, setTheResultIndex, theResultIndex$] = useStateStream(null)//选中的某个指定的框框

    useEffect(()=>{
        isAdding$.next(isAdding)
    },[isAdding])

    //当选择的图片改变时，读取数据:1.放到results列表中  2.重置撤销栈  3.清除选中的框框
    const [recalSelectedImg] = useEventCallback((e$, input$) => e$.pipe(
        delay(100),
        withLatestFrom(input$),
        map(([e, [imgSelectedItem, img]]) => {
            if (imgSelectedItem === '') return []
            else {
                const naturalWidth = img.naturalWidth
                const naturalHeight = img.naturalHeight
                return imgSelectedItem.label ? imgSelectedItem.label.map(item => {
                    item.position.forEach(v => Number(v))
                    const position = {
                        x1: item.position[0] / naturalWidth,
                        y1: item.position[1] / naturalHeight,
                        x2: item.position[2] / naturalWidth,
                        y2: item.position[3] / naturalHeight,
                        x: Math.min(item.position[0], item.position[2]) / naturalWidth,
                        y: Math.min(item.position[1], item.position[3]) / naturalHeight,
                        width: Math.abs(item.position[0] - item.position[2]) / naturalWidth,
                        height: Math.abs(item.position[1] - item.position[3]) / naturalHeight,
                    }
                    return { name: item.name, type: item.type, ...position }
                }): []
            }
        }),
        tap(data => setResults(data)),
        tap(data=>redoDispatch({type:'reset',value:data})),//若改变了选择项，则重置redolist
        tap(() => setTheResultIndex(-1))
    ), '', [imgSelectedItem, imgRef.current])

    //监听selecteditem是否改变
    useEffect(() => {
        recalSelectedImg()
    }, [imgSelectedItem])
    //当绘图结果发生变化以及窗口大小发生变化时，canvas内容重绘
    useEffect(() => {
        if (!drawRect || !holdEvent) return
        if (!results && !results.length) drawRect.clearRects()
        else drawRect && drawRect.drawRectsAndWordsAndFillSelectedItem(results, results[theResultIndex])
    }, [results, resizing, theResultIndex, drawRect, holdEvent])

    useEffect(() => {
        const holdEvent = new HoldEvent(canvasRef.current)
        const drawRect = new DrawRect(canvasRef.current)


        const isInBorderRange = (x, y, x1, x2, y1, y2) => {
            const theX = 3 / drawRect.canvas.width
            const theY = 3 / drawRect.canvas.height 
            const minX = Math.min(x1, x2) - theX
            const maxX = Math.max(x1, x2) + theX
            const minY = Math.min(y1, y2) - theY
            const maxY = Math.max(y1, y2) + theY
            return isInRange(x, y, minX, maxX, minY, maxY)
        }
        const isInRange = (x, y, x1, x2, y1, y2) => {
            const minX = Math.min(x1, x2)
            const maxX = Math.max(x1, x2)
            const minY = Math.min(y1, y2)
            const maxY = Math.max(y1, y2)
            if((minX <= x && x <= maxX) && (minY <= y && y <= maxY)){
                return true
            }else return false
        }

        //过滤正在选中的框框
        const filterTheResult = (data) => {
            var isIn = false
            const sub = combineLatest(results$, theResultIndex$).pipe(
                map(([results, index]) => index !==-1 && results[index] && isInBorderRange(data.x1, data.y1, results[index].x1, results[index].x2, results[index].y1, results[index].y2))
            ).subscribe(data => isIn = !data)
            sub.unsubscribe()
            return isIn
        }

        //过滤不在编辑的情况下
        const filterIsAdding = () => {
            let filterV = false
            const sub = combineLatest(isAdding$).pipe(
                filter(([isAdding]) => isAdding),
            ).subscribe(([data])=>filterV = data)
            sub.unsubscribe()
            return filterV
        }

        //监听画图时的事件
        const drawing = holdEvent.drawingAndCalRatio$.pipe(
            withLatestFrom(isAdding$),
            filter(([_, isAdding]) => isAdding),
            map(([data]) =>data),
            filter(data => filterTheResult(data)),
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            tap(([[data, results], theResultIndex]) => drawRect.drawRectsAndWordsAndFillSelectedItem(results.concat([data]), results[theResultIndex])),//实时绘制当前鼠标移动的框
            map(([[data]]) => data),
        ).subscribe()

        //监听画框完成,抬起鼠标的事件
        const drawingFinish = holdEvent.drawingFinishAndCalRatio$.pipe(
            withLatestFrom(isAdding$),
            filter(([_, isAdding]) => isAdding),
            map(([data, _])=>data),
            filter(data => filterTheResult(data)),
            concatMap(data => whenDrawFinishAndSelectType ? whenDrawFinishAndSelectType(data):of(data)),
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([[data, results], index]) => {
                if (data === 'cancel')//取消此次的绘图结果，并绘制上次的结果
                    drawRect.drawRectsAndWordsAndFillSelectedItem(results, results[index])
                return data !== 'cancel'
            }),
            map(([[data, results]]) => results.concat([data])),
            tap((newResults) => setResults(newResults)),//将画的框放到结果中
            tap(newResults => redoDispatch({type:'push', value:newResults}))//将执行的动作放到撤销栈中
        ).subscribe()

        //当鼠标抬起时监听框框是否可选中
        const seletectingItem = holdEvent.filterMouseUp$.pipe(
            withLatestFrom(isAdding$),
            filter(([_, isAdding]) => isAdding),
            withLatestFrom(holdEvent.filterMouseDown$),
            filter(([[upEvent], downEvent]) => upEvent.pageX === downEvent.pageX && upEvent.pageY === downEvent.pageY),
            map(([[upEvent]]) => upEvent),
            withLatestFrom(results$),//获取最新的results
            map(([e, results]) => [e, results]),
            withLatestFrom(theResultIndex$),//获取当前选中的框框index
            tap(([[e, results], index]) => {
                const x = Number(e.offsetX / e.target.width)
                const y = Number(e.offsetY / e.target.height)
                if(index !== -1 && isInBorderRange(x, y, results[index].x1, results[index].x2, results[index].y1, results[index].y2))return
                for(let i = 0; i < results.length; i++){
                    if(isInBorderRange(x, y, results[i].x1, results[i].x2, results[i].y1, results[i].y2))
                        setTheResultIndex(i)
                }
            })
        ).subscribe()
    
        const getCoordinatesAndCursor = (theResult) => {
            if(!theResult)return []
            const minX = Math.min(theResult.x1, theResult.x2)
            const minY = Math.min(theResult.y1, theResult.y2)
            const maxX = Math.max(theResult.x1, theResult.x2)
            const maxY = Math.max(theResult.y1, theResult.y2)
            const cursorType = ['nw-resize', 'n-resize', 'ne-resize', 'w-resize', 'e-resize', 'sw-resize', 's-resize', 'se-resize']
            const centerCoordinates = [
                [minX, minY],
                [(minX + maxX) / 2, minY],
                [maxX, minY],
                [minX, (minY + maxY) / 2],
                [maxX, (minY + maxY) / 2],
                [minX, maxY],
                [(minX + maxX) / 2, maxY],
                [maxX, maxY]
            ]
            const dX = 3 / drawRect.canvas.width
            const dY = 3 / drawRect.canvas.height
            return centerCoordinates.map((item, i) => [
                cursorType[i], item[0] - dX, item[1] - dY, item[0] + dX, item[1] + dY
            ])
        }

        const isInEdgeArea = (x, y, result) => {
            if(!result)return false
            const coordinates = getCoordinatesAndCursor(result)
            if(coordinates.find(coordinate => isInRange(x, y, coordinate[1], coordinate[3], coordinate[2], coordinate[4]))){
                return true
            }
            else return false
        }


        //监听鼠标移动时cursor的变化
        const isMouseInTheResult = holdEvent.mouseMovingAndCalRatio$.pipe(
            filter(() => {
                if(filterIsAdding()){
                    drawRect.canvas.style.cursor = 'crosshair'
                    return true
                } else {
                    drawRect.canvas.style.cursor = 'unset'
                    return false
                }
            }),
            tap(({x, y}) => {
                if(!filterTheResult({x1:x, y1:y})){//在选中框的范围内
                    drawRect.canvas.style.cursor = 'move'
                }
            }),
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([_, index]) => index!==-1),
            tap(([[moveEvent, results], index]) => {
                const coordinates = getCoordinatesAndCursor(results[index])
                coordinates.map(coordinate => {
                    if(isInRange(moveEvent.x, moveEvent.y, coordinate[1], coordinate[3], coordinate[2], coordinate[4])){
                        drawRect.canvas.style.cursor = coordinate[0]
                    }
                })
            })
        ).subscribe()

        //正在移动框框时
        const movingTheResult = holdEvent.drawingAndCalRatio$.pipe(
            filter(data => !filterTheResult(data)),
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([[data, results], index])=>{
                return !isInEdgeArea(data.x1, data.y1, results[index])
            }),
            tap(([[data, results], index]) => {
                const dX = data.x2 - data.x1
                const dY = data.y2 - data.y1
                const theResult = results[index]
                const theNewResult = {
                    ...theResult,
                    x1:theResult.x1 + dX,
                    y1:theResult.y1 + dY,
                    x2:theResult.x2 + dX,
                    y2:theResult.y2 + dY,
                    x:theResult.x + dX,
                    y:theResult.y + dY
                }
                drawRect.drawRectsAndWordsAndFillSelectedItem(results.slice(0,index).concat([theNewResult], results.slice(index+1)), theNewResult)
            }),
        ).subscribe()

        //移动指定框框结束后的操作
        const movingFinishTheResult = holdEvent.drawingFinishAndCalRatio$.pipe(
            filter(data => !filterTheResult(data)),
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([[data, results], index])=>{
                return !isInEdgeArea(data.x1, data.y1, results[index])
            }),
            tap(([[data, results], index]) => {
                const dX = data.x2 - data.x1
                const dY = data.y2 - data.y1
                const theResult = results[index]
                const theNewResult = {
                    ...theResult,
                    x1:theResult.x1 + dX,
                    y1:theResult.y1 + dY,
                    x2:theResult.x2 + dX,
                    y2:theResult.y2 + dY,
                    x:theResult.x + dX,
                    y:theResult.y + dY
                }
                const newResults = results.slice(0,index).concat([theNewResult], results.slice(index+1))
                redoDispatch({type:'push', value:newResults})
                setResults(newResults)
            })
        ).subscribe()

        const dealSizeChanging =  (direction, data, results, index) => {
            const theResult = results[index]
            const minX = Math.min(theResult.x1, theResult.x2)
            const minY = Math.min(theResult.y1, theResult.y2)
            const maxX = Math.max(theResult.x1, theResult.x2)
            const maxY = Math.max(theResult.y1, theResult.y2)
            let newTheResult = results[index]
            const getResults = () => ({
                newTheResult,
                newResults : [].concat(results.slice(0, index), [newTheResult], results.slice(index+1))
            })
            switch(direction){
                case 'nw-resize':
                    newTheResult = {x1:data.x2, y1:data.y2, x2:maxX, y2:maxY}
                   return getResults()
                case 'n-resize':
                    newTheResult = {...theResult, y2:maxY, y1:data.y2}
                    return getResults()
                case 'ne-resize':
                    newTheResult = { x1:minX,  y1:data.y2, x2:data.x2, y2:maxY}
                    return getResults()
                case 'w-resize':
                    newTheResult = {...theResult, x1:data.x2, x2:maxX}
                    return getResults()
                case 'e-resize':
                    newTheResult = { ...theResult, x1:minX, x2:data.x2}
                    return getResults()
                case 'sw-resize':
                    newTheResult = {x1:data.x2, y1:minY, x2:maxX, y2:data.y2}
                    return getResults()
                case 's-resize':
                    newTheResult = { ...theResult, y1:minY, y2:data.y2}
                    return getResults()
                case 'se-resize':
                    newTheResult = {x1:minX, y1:minY, x2:data.x2,y2:data.y2}
                    return getResults()
            }
        }

        //改变框框的大小
        const changeSize = holdEvent.drawingAndCalRatio$.pipe(
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([_, index]) => index !== -1),
            filter(([[data, results], index]) => {
                return isInEdgeArea(data.x1, data.y1, results[index])
            }),
            tap(([[data, results], index]) => {
                const coordinates = getCoordinatesAndCursor(results[index])
                const theCoordinate = coordinates.find(coordinate => isInRange(data.x1, data.y1, coordinate[1], coordinate[3], coordinate[2], coordinate[4]))
                const {newTheResult} = dealSizeChanging(theCoordinate[0], data, results, index)
                drawRect.drawRectsAndWordsAndFillSelectedItem([].concat(results.slice(0, index), [newTheResult], results.slice(index + 1)), newTheResult)
            })
        ).subscribe()

        const changeSizeFinish = holdEvent.drawingFinishAndCalRatio$.pipe(
            withLatestFrom(results$),
            withLatestFrom(theResultIndex$),
            filter(([_, index]) => index!==-1),
            filter(([[data, results], index]) => {
                return isInEdgeArea(data.x1, data.y1, results[index])
            }),
            tap(([[data, results], index]) => {
                const coordinates = getCoordinatesAndCursor(results[index])
                const theCoordinate = coordinates.find(coordinate => isInRange(data.x1, data.y1, coordinate[1], coordinate[3], coordinate[2], coordinate[4]))
                const {newResults} =  dealSizeChanging(theCoordinate[0], data, results, index)
                setResults(newResults)
                redoDispatch({type:'push', value:newResults})
            })  
        ).subscribe()

        setInitData({ holdEvent, drawRect })
        imgRef.current.addEventListener('mousedown', e => e.preventDefault())
        imgRef.current.onload = (e) => holdEvent.refreshNodePosition()
        
        //监听屏幕大小的变化
        const listenResize = fromEvent(window, 'resize').pipe(
            tap(() => holdEvent.refreshNodePosition()),
            tap(() => resizingUpdate())
        ).subscribe()

        return () => {
            listenResize.unsubscribe()
            drawing.unsubscribe()
            drawingFinish.unsubscribe()
            seletectingItem.unsubscribe()
            isMouseInTheResult.unsubscribe()
            movingFinishTheResult.unsubscribe()
            movingTheResult.unsubscribe()
        }
    }, [])
    const clearResult = () => {
        const emptyArr = []
        setResults(emptyArr)
        redoList.push(emptyArr)
    }
    //确定按钮逻辑，将框转成发送给服务端的数据
    const submit = () => { 
        const data = JSON.parse(JSON.stringify(imgSelectedItem))
        const alreadyEdit = redoList.length >=2 ? true : false
        if (alreadyEdit) {
            const naturalWidth = imgRef.current.naturalWidth
            const naturalHeight = imgRef.current.naturalHeight
            data.label = results.map(item => ({
                name: item.name,
                type: 1,
                position: [
                    (item.x1 * naturalWidth) + '',
                    (item.y1 * naturalHeight) + '',
                    (item.x2 * naturalWidth) + '',
                    (item.y2 * naturalHeight) + ''
                ]
            }))
        }
        onSubmit && onSubmit(data, alreadyEdit)
    }

    const onNameChange = (e, i) => {
        const newResult = { ...results[i], name:e.target.value}
        const newResults = [].concat(results.slice(0, i), [newResult], results.slice(i+1))
        setResults(newResults)
    }
    const onResultDelete = (e, i) => {
        e.preventDefault()
        const newResults = results.slice(0, i).concat(results.slice(i+1))
        if(theResultIndex === i)setTheResultIndex(-1)
        redoDispatch({type:'push', value:newResults})
        setResults(newResults)
    }
    //撤销操作
    const redo = () => {
        if(!redoList.length)return
        redoList.splice(-1)
        setResults(redoList[redoList.length-1])
    }
    const isDisable = imgSelectedItem === '' || !imgRef.current 
    const imgWidth = imgRef.current ? imgRef.current.offsetWidth : 0
    const imgHeight = imgRef.current ? imgRef.current.offsetHeight : 0
    return <div className={styles.wrapper}>
        <div className={styles.content}>
            <canvas width={imgWidth} height={imgHeight} ref={canvasRef} className = {`${isAdding ? styles.isAdding : ''}`} />
            <img
                className={`${styles.img} ${imgSelectedItem === '' ? styles.invisibleImg : ''}`}//TODO: 这里变成模拟的了
                src={( (imgSelectedItem && imgSelectedItem.path ) ? (imgSelectedItem.path.includes('http')?imgSelectedItem.path:usingSpecialPath('http://srv.cloud.jiangxingai.com:8088')+imgSelectedItem.path) : '')}
                // src={( (imgSelectedItem && imgSelectedItem.path ) ? (imgSelectedItem.path) : '')}
                ref={imgRef}
            />
        </div>
        <div className = {styles.dataShow}>
            {results.map((item, i) =><p className = {`${styles.dataItem} ${theResultIndex === i ? styles.dataShowActive:''}`} key = {i} onClick = {() => setTheResultIndex(theResultIndex === i ? -1 : i)}>
                <Input 
                    addonBefore = 'name: ' 
                    value = {item.name} 
                    onChange = {e => onNameChange(e, i)} 
                    addonAfter = {<Icon type="delete" onClick = {(e)=>onResultDelete(e, i)}/>}
                />
                {/* <span>type: {item.type}</span> */}
            </p>)}
        </div>
        {needTool ? <div className={styles.buttonGroup}>
            <div className = {styles.addingTool}>
            </div>
            <WithTitle title='增加结果'><Switch disabled={isDisable} checked={isAdding} onClick={value => setIsAdding(value)} /></WithTitle>
            <Button disabled = {redoList.length <= 1} onClick = {redo}>撤销</Button>
            <Button disabled={isDisable} onClick={clearResult}>清空结果</Button>
            {onCancel && <Button disabled={isDisable} onClick = {()=>onCancel(imgSelectedItem)}>无需告警</Button>}
            <Button disabled={isDisable} onClick={submit}>确定</Button>
        </div> :''}
    </div>
}