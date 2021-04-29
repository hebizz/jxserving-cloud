import React from 'react';
import styles from './index.module.scss'
import { Button } from 'antd';
import {NavLink} from 'react-router-dom'
export const AppContext = React.createContext({})
export default function () {
	return (
		<div className={styles.wrapper}>
			<Button>
				<NavLink to='/analyst'>
					Analyst 分析器
				</NavLink>
			</Button>
			<Button>
				<NavLink to='/model'>Model 模型管理</NavLink>
			</Button>
			{/* <Button></Button>
			<Button></Button> */}
		</div>
	);
}
