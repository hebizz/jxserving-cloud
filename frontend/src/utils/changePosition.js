export default function({
    changeTo,
    x1,x2,y1,y2,
    left_x,left_y,right_x,right_y,
    showWidth, showHeight, naturalWidth, naturalHeight
}){
    if(changeTo==='natural'){//需要参数x1,x2,y1,y2做替换
        const widthRatio = naturalWidth / showWidth
        const heightRatio = naturalHeight / showHeight
        const left_x = Math.min(x1, x2) * widthRatio
        const left_y = Math.min(y1, y2) * heightRatio
        const right_x = Math.max(x1, x2) * widthRatio
        const right_y = Math.max(y1, y2) * heightRatio
        const width = Math.abs(left_x - right_x)
        const height = Math.abs(left_y - right_y)
        return {
            left_x, left_y, right_x, right_y,
            x1,y1,x2,y2,
            x:left_x, y:left_y, width, height
        }
    }else if(changeTo === 'show'){//需要参数left_x,left_y,right_x,right_y
        const widthRatio = showWidth / naturalWidth
        const heightRatio = showHeight / naturalHeight
        const x1 = left_x * widthRatio
        const y1 = left_y * heightRatio
        const x2 = right_x * widthRatio
        const y2 = right_y * heightRatio
        const width = Math.abs(x1 - x2)
        const height = Math.abs(y1 - y2)
        return {
            left_x,left_y,right_x,right_y,
            x1, y1, x2 ,y2,
            x:x1, y:y1,width, height
        }
    }
}