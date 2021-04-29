export default class DrawRect {
    constructor(canvas) {
        this.canvas = canvas
        this.cxt = canvas.getContext('2d')
    }
    clearRects = () => {
        this.cxt.clearRect(0, 0, this.canvas.width, this.canvas.height)
    }
    drawRects = (datas) => {
        this.cxt.strokeStyle = 'red'
        datas.forEach(item => this.cxt.strokeRect(
            Math.min(item.x1, item.x2) * this.canvas.width,
            Math.min(item.y1, item.y2) * this.canvas.height,
            Math.abs(item.x1 - item.x2) * this.canvas.width,
            Math.abs(item.y1 - item.y2) * this.canvas.height
        ))
    }
    fillRect = (item) => {
        this.cxt.fillStyle = 'rgba(255,255,255,0.5)'
        item && this.cxt.fillRect(
            Math.min(item.x1, item.x2) * this.canvas.width,
            Math.min(item.y1, item.y2) * this.canvas.height,
            Math.abs(item.x1 - item.x2) * this.canvas.width,
            Math.abs(item.y1 - item.y2) * this.canvas.height
        )
    }
    drawDots = item => {
        this.cxt.fillStyle = 'red'
        if(!item)return
        const coordinates = [
            [item.x1, item.y1],
            [(item.x1 + item.x2) / 2, item.y1],
            [item.x2, item.y1],
            [item.x1, (item.y1 + item.y2) / 2],
            [item.x2, (item.y1 + item.y2) / 2],
            [item.x1, item.y2],
            [(item.x1 + item.x2) / 2, item.y2],
            [item.x2, item.y2]
        ]
        coordinates.map(coordinate => {
            this.cxt.fillRect(coordinate[0]*this.canvas.width -3, coordinate[1]*this.canvas.height - 3, 6, 6)
        })
    }
    drawWords = (datas) => {
        this.cxt.fillStyle = 'red'
        this.cxt.font = 'bold 16px Arial'
        datas.forEach(item => {
            item.name !== undefined && this.cxt.fillText(item.name, Math.min(item.x1, item.x2) * this.canvas.width,  Math.min(item.y1, item.y2) * this.canvas.height - 5)
        })
    }
    drawRectsAndWords = (datas) => {
        this.clearRects()
        this.drawWords(datas)
        this.drawRects(datas)
    }
    drawRectsAndWordsAndFillSelectedItem = (datas, item) => {
        this.drawRectsAndWords(datas)
        this.fillRect(item)
        this.drawDots(item)
    }
}