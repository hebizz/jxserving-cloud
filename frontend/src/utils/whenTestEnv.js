export const usingSpecialPath = (path) => {
    if(process.env.NODE_ENV === 'production')return ''
    else return path
}