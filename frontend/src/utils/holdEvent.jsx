import React, {useEffect, useRef} from 'react'
import { filter, map, concatAll, withLatestFrom, takeUntil, tap,} from "rxjs/operators";
import { fromEvent } from "rxjs";

function getPos(node){//获得节点相对网页的x，y坐标值
    const position = {
        left:node.offsetLeft,
        top:node.offsetTop
    }
    var p = node
    while(p.offsetParent){
        p = p.offsetParent
        position.left += p.offsetLeft||0
        position.top += p.offsetTop||0
    }
    return position
}

function originToRatio(e, nodePosition, node){

}

/**
 *node  在此拖拽的节点
 *
 * @class HoldEvents
 */
export default class HoldEvents{
    constructor(node){
        if(!node)throw Error('need correct dom node')
        this.node = node
        this.nodePosition = getPos(node)
        const filterFunc = e => {
            return e.target===node || node.contains(e.target)
        }
        const mouseDown$ = fromEvent(document.body, 'mousedown')
        const mouseMove$ = fromEvent(document.body, 'mousemove')
        const mouseUp$ = fromEvent(document.body, 'mouseup')

        this.filterMouseDown$ = mouseDown$.pipe(filter(filterFunc))
        this.filterMouseUp$ = mouseUp$.pipe(filter(filterFunc))

        this.mouseMovingAndCalRatio$ =  mouseMove$.pipe(
            filter(filterFunc),
            map(e => {
                const nodePosition = getPos(node)
                const x = (e.pageX - nodePosition.left) / node.offsetWidth
                const y = ( e.pageY - nodePosition.top) / node.offsetHeight
                return {x, y}
            })
        )
        this.drawingBefore$ = this.filterMouseDown$.pipe(
            map(e => ({x1:e.pageX, y1:e.pageY}))
        )
        this.drawing$ = this.filterMouseDown$.pipe(//返回的坐标值相对网页定位
            map(() => mouseMove$.pipe(takeUntil(mouseUp$))),
            concatAll(),
            withLatestFrom(this.drawingBefore$),
            map(([e, from]) => {
                const nodePosition = getPos(this.node)
                const boundaryLeft = nodePosition.left
                const boundaryRight = nodePosition.left + node.offsetWidth
                const boundaryTop = nodePosition.top
                const boundaryBottom = nodePosition.top + node.offsetHeight
                const x2 = boundaryLeft<e.pageX ? e.pageX<boundaryRight ? e.pageX : boundaryRight : boundaryLeft
                const y2 = boundaryTop<e.pageY ? e.pageY<boundaryBottom ? e.pageY : boundaryBottom : boundaryTop
                // const minX = Math.min(from.x1, x2)
                // const maxX = Math.max(from.x1, x2)
                // const minY = Math.min(from.y1, y2)
                // const maxY = Math.max(from.y1, y2)
                // return {x1:minX, y1:minY, x2:maxX, y2:maxY}
                return {...from, x2, y2}
            })
        )
        this.drawingRelative$ = this.drawing$.pipe(//返回的坐标值相对该节点定位
            map(position=> {
                const nodePosition = getPos(this.node)
                const x1 = position.x1 - nodePosition.left
                const y1 = position.y1 - nodePosition.top
                const x2 = position.x2 - nodePosition.left
                const y2 = position.y2 - nodePosition.top
                const x = Math.min(x1, x2)
                const y = Math.min(y1, y2)
                const width = Math.abs(x2 - x1)
                const height = Math.abs(y2 - y1)
                return {
                    x1, y1, x2, y2,
                    x, y, width, height
                }
            })
        )
        this.drawingAndCalRatio$ = this.drawingRelative$.pipe(//返回的坐标值为比例值，范围在[0, 1]间
            map(position => ({
                x1: position.x1  / node.offsetWidth,
                y1: position.y1 / node.offsetHeight,
                x2: position.x2  / node.offsetWidth,
                y2: position.y2 / node.offsetHeight,
                x:position.x/node.offsetWidth,
                y:position.y/node.offsetHeight,
                width:position.width/node.offsetWidth,
                height:position.height/node.offsetHeight
            }))
        )
            
        this.drawingAndCalNatural$ = this.drawingAndCalRatio$.pipe(//返回的坐标值，相对于原始大小定位，适用于图片大小上
            map(position => ({
                x1:position.x1*node.naturalWidth,
                y1:position.y1*node.naturalHeight,
                x2:position.x2*node.naturalWidth,
                y2:position.y2*node.naturalHeight,
            }))
        ) 

        //返回坐标值是相对于网页的(x1,x2,y1,y2)，并返回当前容器的宽高(containerWidth, containerHeight)
        this.drawingFinish$ = mouseUp$.pipe(
            withLatestFrom(mouseDown$),
            filter(([up, down])=>filterFunc(down)),
            filter(([{pageX:pageXUp, pageY:pageYUp}, {pageX:pageXDown, pageY:pageYDown}]) => pageXUp !== pageXDown && pageYUp !== pageYDown),//防止鼠标没有触发滑动事件
            withLatestFrom(this.drawing$),
            map(([e, position])=>({
                ...position,
                containerWidth:node.offsetWidth,
                containerHeight:node.offsetHeight,
            }))
        )
        //sdfsdfsdfsdfds
        
        //返回的坐标值相对于容器
        this.drawingFinishAndCalRelative$ = this.drawingFinish$.pipe(
            map(position=> {
                const nodePosition = getPos(this.node)
                const x1 = position.x1 - nodePosition.left
                const y1 = position.y1 - nodePosition.top
                const x2 = position.x2 - nodePosition.left
                const y2 = position.y2 - nodePosition.top
                const x = Math.min(x1, x2)
                const y = Math.min(y1, y2)
                const width = Math.abs(x2 - x1)
                const height = Math.abs(y2 - y1)
                return {
                    ...position,
                    x1, y1, x2, y2,
                    x, y, width, height
                }
            })
        )

        //放回的坐标值是按比例的
        this.drawingFinishAndCalRatio$ = this.drawingFinishAndCalRelative$.pipe(
            map(position => ({
                ...position,
                x1: position.x1 / node.offsetWidth,
                y1: position.y1 / node.offsetHeight,
                x2: position.x2 / node.offsetWidth,
                y2: position.y2 / node.offsetHeight,
                x:position.x /  node.offsetWidth,
                y:position.y / node.offsetHeight,
                width: position.width / node.offsetWidth,
                height: position.height / node.offsetHeight
            }))
        )
        
        //返回的坐标值是相对于node原始大小的，适合图片
        this.drawingFinishAndCalNatural$ = this.drawingFinishAndCalRatio$.pipe(
            map(position => ({
                ...position,
                x1:position.x1*node.naturalWidth,
                y1:position.y1*node.naturalHeight,
                x2:position.x2*node.naturalWidth,
                y2:position.y2*node.naturalHeight,
            }))
        ) 
    }

    refreshNodePosition = () => this.nodePosition = getPos(this.node)

    drawingBeforeSubscribe = (cb) =>this.drawingBefore$.subscribe(cb)

    drawingSubscribe = cb => this.drawing$.subscribe(cb)
    drawingRelativeSubscribe = cb => this.drawingRelative$.subscribe(cb)
    drawingAndCalRatioSubscribe = cb => this.drawingAndCalRatio$.subscribe(cb)
    drawingAndCalNaturalSubscribe = cb => this.drawingAndCalNatural$.subscribe(cb)


    drawingFinishSubscribe = cb => this.drawingFinish$.subscribe(cb)
    drawingFinishRelativeSubscribe = cb => this.drawingFinishAndCalRelative$.subscribe(cb)
    drawingFinishRatioSubscribe = cb => this.drawingFinishAndCalRatio$.subscribe(cb)
    drawingFinishAndCalNaturalSubscribe = cb => this.drawingFinishAndCalNatural$.subscribe(cb)
}