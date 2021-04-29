export const withKeyColumn = key => {
    return {
        sorter:(a, b) => (a[key] > b[key]) ? 1 : -1
    }
}

export const getFiltersWhenDataing = (dataSource, key) => {
    if(!(dataSource instanceof Array))return {}
    const filtersValue = dataSource.map(item => item[key]) // [value1, value2, value3, ...]
    const uniqFilters = Array.from(new Set(filtersValue))
    const filters = uniqFilters.map(value=>({
        text:value, 
        value
    }))
    const onFilter = (value, record) => record[key].toString().indexOf(value) === 0
    return {
        filters,
        onFilter
    }
}
export const mergeToOriginObj = ((originObj, newObj) => {
    Object.keys(newObj).forEach(key=>{
        originObj[key] = newObj[key]
    })
}) 
