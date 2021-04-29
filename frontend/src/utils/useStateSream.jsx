import {useState, useEffect, useMemo} from 'react'
import { BehaviorSubject } from 'rxjs';
export function useStateStream(initValue){
    const [theValue, setTheValue]= useState(initValue)
    const stream$ = useMemo(()=>new BehaviorSubject(theValue), [])
    useEffect(()=>{
        stream$.next(theValue)
    }, [theValue])
    return [theValue, setTheValue, stream$]
}