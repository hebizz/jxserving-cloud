import React from 'react';
import styles from './index.module.scss'
import { Card, Descriptions } from 'antd';
const {Item} = Descriptions
export default function () {
	return (
		<div className={styles.wrapper}>
			<Descriptions column = {1} title = '状态' bordered>
				<Item label='On-premises'>已启动</Item>
				<Item label='人工干预'>已启动</Item>
				<Item label='接入应用'>sgcc</Item>
				<Item label='系统负载'>正常</Item>
			</Descriptions>
		</div>
	);
}
