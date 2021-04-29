import React, {useState, useEffect} from 'react';
import styles from './index.module.scss'
import {NavLink, withRouter} from 'react-router-dom'
import { Menu, Icon, Progress } from 'antd';
import { useStateStream } from 'utils/useStateSream';
import { skip, tap } from 'rxjs/operators';
const titleMap = {
	'/usermanagement':'用户管理',
	'/glance':'概览',
	'/jxservingonpremises':'On-premises',
	'/dataset':'数据集管理',
	'/datasetmanage':'数据集管理',
	'/label':'标签管理',
	'/model':'Model 模型管理',
	'/analyst':'Analyst 分析器',
	'/intervention':'人工干预'
}
const datas = [
	{
		title:'概览',
		icon:'smile',
		path:'/glance'
	},
	{
		title:'On-premises',
		icon:'smile',
		path:'/jxservingonpremises',
		// disabled:true
	},
	{
		title:'数据集管理',
		icon:'smile',
		path:'/dataset'
	},
	{
		title:'标签管理',
		icon:'smile',
		path:'/label'
	},
	{
		title:'模型管理',
		icon:'smile',
		path:'/model'
	},
	{
		title:'阈值分析',
		icon:'smile',
		path:'/analyst'
	},
	{
		title:'人工干预',
		icon:'smile',
		path:'/intervention'
	},
	// {
	// 	title:'数据备份',
	// 	icon:'smile',
	// 	path:'/databackup'
	// },
]
export default withRouter(function ({children, history}) {
	const key = window.location.hash.replace('#/','').replace(/\/.*/,'').replace(/\?.*/,'').toLowerCase()
	// const item = datas.find( item => item.path === ('/'+key))
	const title = titleMap['/'+key]
	const [selectedKeys, setSelectedKeys] = useState(['/'+key])
	const [theTab, setTheTab, theTab$] = useStateStream(0)
	useEffect(()=>{
		const sub = theTab$.pipe(
			skip(1)
		).subscribe(theTab =>{
			if(theTab===1){
				history.push('/usermanagement')
			}else history.goBack()
		})
		return ()=>{
			sub.unsubscribe()
		}
	},[])
	return (
		<div className={styles.wrapper}>
			<header>
				<span className = {styles.title}>
					{title}
				</span>
				<span className = {styles.toolbar}>
					<span onClick = {()=>setTheTab(0)} className = { theTab===0 ? styles.selected : ''}>工作台</span>
					<span onClick = {()=>setTheTab(1)} className = { theTab===1 ? styles.selected : ''} >用户管理</span>
				</span>
			</header>
			<div className = {styles.content}>
				<Menu className = {styles.menu} selectedKeys = {selectedKeys} onSelect = {({key})=>setSelectedKeys([key])}>
					{datas.map(item=><Menu.Item key = {item.path} disabled = {item.disabled?true:false}>
						<NavLink to = {item.path} >
							<Icon type={item.icon}></Icon>
							<span>{item.title}</span>
						</NavLink>
					</Menu.Item>)}
					<Menu.Item className = {styles.dataBackup} disabled = {true}>
						<span>
							<Icon type={'smile'}></Icon>
							<span>数据备份</span>
						</span>
						<Progress percent = {60}/>
					</Menu.Item>
				</Menu>
				{children}
			</div>
    </div>
	);
})
