import React from 'react'
import Loadable from 'react-loadable';
export default function loadable(loader, loadingComponent=()=><div></div>){
    return Loadable({
        loader,
        loading:loadingComponent
    });
}