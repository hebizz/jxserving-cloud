export function hasRoutePermission(item, role={}) {
    // if (item.path.startsWith('/app') && role.app === 'no') return null
    // if (item.path.startsWith('/alarm') && role.alert === 'no') return null
    // if (item.path.startsWith('/device') && role.device === 'no') return null
    // if (item.path.startsWith('/developer') && role.developer === 'no') return null
    // if (item.path.startsWith('/user') && !role.isAdmin && !role.isSubAdmin) return null
    // if (item.path.startsWith('/account') && !role.isAdmin) return null
    return true
}
export function hasAsideItemPermission(item, role) {
    if (item.name === '应用管理' && role.app === 'no') return null
    if (item.name === '设备管理' && role.device === 'no') return null
    if (item.name === '告警管理' && role.alert === 'no') return null
    if (item.name === '开发管理' && role.developer === 'no') return null
    if (item.name === '用户管理' && !role.isAdmin && !role.isSubAdmin) return null
    if (item.name === '子账户管理' && !role.isAdmin) return null
    return true
}

const roleStr2Obj = s => s.map(role => role.slice(5).split('_')).map(x => ({ [x[0]]: x[1] || 'admin' })).reduce((a, b) => ({ ...a, ...b }), {})

export function Role(roleStrArr) {
    this.hasInited = roleStrArr instanceof Array ? true : false
    this.roleObj = roleStrArr ? roleStr2Obj(roleStrArr) : {}
    this.isAdmin = this.roleObj.hasOwnProperty('super')
    this.isSubAdmin = this.roleObj.hasOwnProperty('sub')
    this.app = this.isAdmin || this.isSubAdmin ? 'admin' : this.roleObj.app || 'no'
    this.device = this.isAdmin || this.isSubAdmin ? 'admin' : this.roleObj.device || 'no'
    this.alert = this.isAdmin || this.isSubAdmin ? 'admin' : this.roleObj.alert || 'no'
    this.developer = this.isAdmin || this.isSubAdmin ? 'admin' : this.roleObj.hasOwnProperty('developer') ? 'admin' : 'no'
}