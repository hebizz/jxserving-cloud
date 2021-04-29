const units = ['B', 'KB', 'MB', 'GB','TB','PB', 'EB', 'ZB', 'YB', 'BB']
const calSize = (size, index=0) => {
	return size<1024 ? Number(size).toFixed(2)+units[index] : calSize(size/1024, index+1)
}
export default calSize