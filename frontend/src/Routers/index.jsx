import React, { useContext, useEffect } from 'react';
import { HashRouter, Switch, Route, Redirect } from 'react-router-dom';
import routerConfig from './routerConfig';
import { hasRoutePermission } from './permissionCheck';
import { AppContext } from 'App'
import BasisLayout from 'Layouts/BasisLayout'

// const Login = loadable(()=>import('pages/Login'))
const useRenderNormalRoute = (item, index) => {
    const { role } = useContext(AppContext)
    return hasRoutePermission(item, role) ? <Route
        key={index}
        path={item.path}
        component={() => <item.component />}
        exact={item.exact ? true : false}
    /> : null;
}
export default function () {
    const hash = window.location.hash.replace('#/','').replace(/\/.*/,'').replace(/\?.*/,'').toLowerCase()
    useEffect(()=>{
        if(!hash)window.location.hash= '/glance'
    })
    return <HashRouter>
        <Route path='/' component={() => (
            <BasisLayout>
                <Switch>
                    {routerConfig.map(useRenderNormalRoute)}
                </Switch>
            </BasisLayout>
        )} />
    </HashRouter>
}